package vpc_test

import (
	"context"
	"fmt"
	"os"
	sysuser "os/user"
	"path/filepath"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	vpcv1d3 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1d3"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/config"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// createdVpcId records the ID of the VPC created during the test
// so the sweeper can delete it directly by ID.
var createdVpcId string

// -----------------------------------------------------------------------
// VPC test
// -----------------------------------------------------------------------

func TestAccVpcResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC is set")
	}

	const resourceName = "samsungcloudplatformv2_vpc_vpc.vpc"

	name := testVPCResourceName("gdcv-vpc")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{ // step1. create with description = null (verifies null-description fix — API returns default)
				PreConfig: func() { t.Log("[CREATE] creating VPC with null description") },
				Config: testAccVpcConfigNullDesc(vpcParams{
					name: name,
					cidr: "192.168.0.0/16",
					tags: `env = "test"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "vpc.zone_type"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc.zones.#"),
					captureVpcId(resourceName),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{ // step2. update description
				PreConfig: func() { t.Log("[UPDATE] changing description") },
				Config:    testAccVpcConfig(vpcParams{name: name, description: "updated description", cidr: "192.168.0.0/16", tags: `env = "test"`}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{ // step3. import state check
				PreConfig:               func() { t.Log("[IMPORT] importing by ID") },
				Config:                  testAccVpcConfig(vpcParams{name: name, description: "updated description", cidr: "192.168.0.0/16", tags: `env = "test"`}),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cidr", "tags"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return rs.Primary.ID, nil
				},
			},
		},
	})
}

// -----------------------------------------------------------------------
// Config helpers
// -----------------------------------------------------------------------

type vpcParams struct {
	name        string
	description string
	cidr        string
	tags        string // tags block body, e.g. `env = "test"`
}

func testAccVpcConfig(p vpcParams) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_vpc_vpc" "vpc" {
			name        = %q
			description = %q
			cidr        = %q
			tags = {
				%s
			}
		}
	`, p.name, p.description, p.cidr, p.tags)
}

func testAccVpcConfigNullDesc(p vpcParams) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_vpc_vpc" "vpc" {
			name        = %q
			description = null
			cidr        = %q
			tags = {
				%s
			}
		}
	`, p.name, p.cidr, p.tags)
}

// -----------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------

// testVPCResourceName returns a unique name for each test run to avoid conflicts
// with leaked resources from previous runs.
func testVPCResourceName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano()%1000000)
}

// vpcTestClient creates an SCP client for VPC test operations.
func vpcTestClient() (client.Instance, error) {
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

// captureVpcId reads the created resource ID from state and stores it
// for the sweeper.
func captureVpcId(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		createdVpcId = rs.Primary.ID
		return nil
	}
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_vpc_vpc", &resource.Sweeper{
		Name: "samsungcloudplatformv2_vpc_vpc",
		F:    sweepVpc,
	})
}

func sweepVpc(_ string) error {
	if createdVpcId == "" {
		return nil
	}

	inst, err := vpcTestClient()
	if err != nil {
		return fmt.Errorf("error getting client: %v", err)
	}

	if err := inst.Client.VpcV1Dot3.DeleteVpc(context.Background(), createdVpcId); err != nil {
		return fmt.Errorf("error deleting VPC %s: %v", createdVpcId, err)
	}

	return nil
}

// TestVpcAttributeTypes is a unit test that verifies the VpcValue
// model exposes all expected attribute type keys.
func TestVpcAttributeTypes(t *testing.T) {
	t.Log("[UNIT] verifying VpcValue model attribute types")
	v := vpcv1d3.VpcValue{}
	types := v.AttributeTypes()

	expectedKeys := []string{
		"account_id", "cidr_count", "cidrs", "created_at",
		"created_by", "description", "id", "modified_at",
		"modified_by", "name", "state", "zone_type", "zones",
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
