// Linux: TF_ACC=1 go test -v -count=1 -timeout=30m ./samsungcloudplatform/service/securitygroup/addressgroup_test.go
// Window CMD: set TF_ACC=1 && go test -v -count=1 -timeout=30m ./samsungcloudplatform/service/securitygroup/addressgroup_test.go
package securitygroup_test

import (
	"context"
	"fmt"
	"os"
	sysuser "os/user"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	securitygroupv1d1 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/securitygroupv1d1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/config"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// createdAddressGroupId records the ID of the address group created during the test
// so the sweeper can delete it directly by ID.
var createdAddressGroupId string

func TestAccAddressGroupResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC is set")
	}

	const resourceName = "samsungcloudplatformv2_security_group_address_group.address_group"

	// Generate a unique name once for all steps.
	name := testResourceName("acctest-addr")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{ // step1. create with description = null (verifies null-description fix — API returns default, no phantom diff)
				PreConfig: func() { t.Log("[CREATE] creating address group with null description") },
				Config: testAccAddressGroupConfigNullDesc(addressGroupParams{
					name:      name,
					addresses: `["10.0.0.1/32"]`,
					tags:      `env = "test"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					captureAddressGroupId(resourceName),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{ // step2. update description
				PreConfig: func() { t.Log("[UPDATE] changing description") },
				Config: testAccAddressGroupConfig(addressGroupParams{
					name:        name,
					description: "updated description",
					addresses:   `["10.0.0.1/32"]`,
					tags:        `env = "test"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{ // step3. remove one address and add two different ones at the same time
				PreConfig: func() { t.Log("[UPDATE] removing 10.0.0.1/32, adding 192.168.1.0/24 and 172.16.0.0/16") },
				Config: testAccAddressGroupConfig(addressGroupParams{
					name:        name,
					description: "updated description",
					addresses:   `["192.168.1.0/24", "172.16.0.0/16"]`,
					tags:        `env = "test"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "addresses.*", "192.168.1.0/24"),
					resource.TestCheckTypeSetElemAttr(resourceName, "addresses.*", "172.16.0.0/16"),
				),
			},
			{ // step4. import state check
				PreConfig: func() { t.Log("[IMPORT] importing by ID") },
				Config: testAccAddressGroupConfig(addressGroupParams{
					name:        name,
					description: "updated description",
					addresses:   `["192.168.1.0/24", "172.16.0.0/16"]`,
					tags:        `env = "test"`,
				}),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return rs.Primary.ID, nil
				},
			},
			// step5. [DELETE] runs implicitly after all steps via Destroy
		},
	})
}

// testResourceName returns a unique name for each test run to avoid conflicts
// with leaked resources from previous runs.
func testResourceName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano()%1000000)
}

type addressGroupParams struct {
	name        string
	description string // "null" emits bare HCL null; otherwise quoted string
	addresses   string // HCL list literal, e.g. `["10.0.0.1/32"]`
	tags        string // tags block body, e.g. `env = "test"`
}

func testAccAddressGroupConfig(p addressGroupParams) string {
	return fmt.Sprintf(`
		provider "samsungcloudplatformv2" {}

		resource "samsungcloudplatformv2_security_group_address_group" "address_group" {
			name        = %q
			description = %q
			addresses   = %s
			tags = {
				%s
			}
		}
	`, p.name, p.description, p.addresses, p.tags)
}

func testAccAddressGroupConfigNullDesc(p addressGroupParams) string {
	return fmt.Sprintf(`
		provider "samsungcloudplatformv2" {}

		resource "samsungcloudplatformv2_security_group_address_group" "address_group" {
			name        = %q
			description = null
			addresses   = %s
			tags = {
				%s
			}
		}
	`, p.name, p.addresses, p.tags)
}

// addressGroupTestClient creates an SCP client for address group test operations.
func addressGroupTestClient() (client.Instance, error) {
	providerConfig := config.ProviderConfig{}
	user, _ := sysuser.Current()

	config.LoadServiceConfig(nil, filepath.Join(user.HomeDir, ".scpconf", config.ServiceConfigFile), &providerConfig)
	config.LoadCredentialsConfig(nil, filepath.Join(user.HomeDir, ".scpconf", config.CredentialConfigFile), &providerConfig)

	scpClient, err := client.NewSCPClient(&providerConfig)
	if err != nil {
		return client.Instance{}, fmt.Errorf("error creating SCP client: %v", err)
	}

	return client.Instance{Client: scpClient}, nil
}

// captureAddressGroupId reads the created resource ID from state and stores it
// for the sweeper.
func captureAddressGroupId(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		createdAddressGroupId = rs.Primary.ID
		return nil
	}
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_security_group_address_group", &resource.Sweeper{
		Name: "samsungcloudplatformv2_security_group_address_group",
		F:    sweepAddressGroup,
	})
}

func sweepAddressGroup(_ string) error {
	if createdAddressGroupId == "" {
		return nil
	}

	inst, err := addressGroupTestClient()
	if err != nil {
		return fmt.Errorf("error getting client: %v", err)
	}

	if err := inst.Client.SecurityGroupV1d1.DeleteAddressGroup(context.Background(), createdAddressGroupId); err != nil {
		return fmt.Errorf("error deleting address group %s: %v", createdAddressGroupId, err)
	}

	return nil
}

// TestAccAddressGroupResourceModifyPlanPreventsChange verifies that ModifyPlan
// blocks in-place changes to immutable fields (name, tags).
func TestAccAddressGroupResourceModifyPlanPreventsChange(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC is set")
	}

	// Generate unique names once for all steps.
	name := testResourceName("acctest-modify")
	renamedName := testResourceName("acctest-modify-renamed")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() { t.Log("[CREATE] creating address group for ModifyPlan test") },
				Config: testAccAddressGroupConfig(addressGroupParams{
					name:        name,
					description: "initial",
					addresses:   `["10.0.0.1/32"]`,
					tags:        `env = "test"`,
				}),
			},
			{
				PreConfig: func() { t.Log("[UPDATE] attempting to change immutable name — should fail") },
				Config: testAccAddressGroupConfig(addressGroupParams{
					name:        renamedName,
					description: "initial",
					addresses:   `["10.0.0.1/32"]`,
					tags:        `env = "test"`,
				}),
				ExpectError: regexp.MustCompile(`Field changes not supported`),
			},
			{
				PreConfig: func() { t.Log("[UPDATE] attempting to change immutable tags — should fail") },
				Config: testAccAddressGroupConfig(addressGroupParams{
					name:        name,
					description: "initial",
					addresses:   `["10.0.0.1/32"]`,
					tags:        `env = "changed"`,
				}),
				ExpectError: regexp.MustCompile(`Field changes not supported`),
			},
		},
	})
}

// TestAddressGroupAttributeTypes is a unit test that verifies the AddressGroup
// model exposes all expected attribute type keys.
func TestAddressGroupAttributeTypes(t *testing.T) {
	t.Log("[UNIT] verifying AddressGroup model attribute types")
	ag := securitygroupv1d1.AddressGroup{}
	types := ag.AttributeTypes()

	expectedKeys := []string{
		"account_id", "address_count", "address_limit", "addresses",
		"created_at", "created_by", "description", "id",
		"modified_at", "modified_by", "name", "state",
	}

	if len(types) != len(expectedKeys) {
		t.Fatalf("expected %d attribute types, got %d", len(expectedKeys), len(types))
	}

	for _, key := range expectedKeys {
		if _, ok := types[key]; !ok {
			t.Errorf("missing attribute type for %q", key)
		}
	}
}
