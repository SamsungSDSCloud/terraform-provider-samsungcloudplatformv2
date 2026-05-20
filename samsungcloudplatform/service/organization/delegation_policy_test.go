package organization_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDelegationPolicyResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDelegationPolicyCreate(),
			},
			{
				Config: testAccDelegationPolicyUpdate(),
			},
		},
	})
}

func testAccDelegationPolicyCreate() string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_organization_delegation_policy" "delegation_policy" {
			document = {
				statement = [
					{
						action = ["organization:CreateServiceControlPolicy",
                    			  "organization:ListAccounts"]
						effect = "Allow"
						resource = ["*"]
						sid = "test-acc-sid-1"
						principal = {
							scp = ["srn:dev2::f045159f40c64125a1fe61bd71d1c14c:::iam:user/e3a6e3f99c1040639e9a3c8f8b7427de"]
						}
					}
				]
				version = "2024-07-01"
			}
			organization_id = "o-2b63982e88b74dbcb71ee972b13e2ce1"
		}`)
}

func testAccDelegationPolicyUpdate() string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_organization_delegation_policy" "delegation_policy" {
			document = {
				statement = [
					{
						action = ["organization:CreateServiceControlPolicy"]
						effect = "Allow"
						resource = ["*"]
						sid = "test-acc-sid-1"
						principal = {
							scp = ["srn:dev2::b219cfc010b04804a6e69a6931b09cc1:::iam:user/dece6618dde444eeb7a4ff4bee84361a"]
						}
					}
				]
				version = "2024-07-01"
			}
			organization_id = "o-2b63982e88b74dbcb71ee972b13e2ce1"
		}`)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_organization_delegation_policy", &resource.Sweeper{
		Name: "samsungcloudplatformv2_organization_delegation_policy",
		F:    sweepDelegationPolicy,
	})
}

func sweepDelegationPolicy(region string) error {
	scpClient, err := SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	orgList, err := scpClient.Client.Organization.GetOrganizationList(nil, organization.OrganizationDataSource{})
	if err != nil {
		return nil
	}

	for _, org := range orgList.GetOrganizations() {
		orgId := org.Id
		policy, err := scpClient.Client.Organization.GetDelegationPolicy(nil, orgId)
		if err != nil {
			continue
		}

		doc := policy.Policy.Document
		for _, stmt := range doc.Statement {
			if strings.HasPrefix(stmt.GetSid(), "test-acc") {
				_, delErr := scpClient.Client.Organization.DeleteDelegationPolicy(nil, orgId)
				if delErr != nil {
					return fmt.Errorf("error deleting delegation policy for org %s: %v", orgId, delErr)
				}
				break
			}
		}
	}

	return nil
}
