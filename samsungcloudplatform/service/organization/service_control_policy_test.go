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

func TestAccServiceControlPolicyResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServiceControlPolicyCreate(),
			},
			{
				Config: testAccServiceControlPolicyUpdate(),
			},
		},
	})
}

func testAccServiceControlPolicyCreate() string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_organization_service_control_policy" "service_control_policy" {
			name        = "test-acc-service-control-policy"
			description = "test-acc description"
			type        = "USER_DEFINED"
			document = {
				statement = [
					{
						action     = ["organization:CreateServiceControlPolicy", "organization:ListAccounts"]
						effect     = "Allow"
						resource   = ["*"]
						sid        = "test-acc-sid-1"
						not_action = []
						principal  = "*"
						condition  = {}
					}
				]
				version = "2024-07-01"
			}
			organization_id = "o-2b63982e88b74dbcb71ee972b13e2ce1"
		}`)
}

func testAccServiceControlPolicyUpdate() string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_organization_service_control_policy" "service_control_policy" {
			name        = "test-acc-service-control-policy-update"
			description = "test-acc description updated"
			type        = "USER_DEFINED"
			document = {
				statement = [
					{
						action     = ["organization:CreateServiceControlPolicy"]
						effect     = "Allow"
						resource   = ["*"]
						sid        = "test-acc-sid-2"
						not_action = []
						principal  = "*"
						condition  = {}
					}
				]
				version = "2024-07-01"
			}
			organization_id = "o-2b63982e88b74dbcb71ee972b13e2ce1"
		}`)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_organization_service_control_policy", &resource.Sweeper{
		Name: "samsungcloudplatformv2_organization_service_control_policy",
		F:    sweepServiceControlPolicy,
	})
}

func sweepServiceControlPolicy(region string) error {
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

		policies, err := scpClient.Client.Organization.ListServiceControlPolicies(nil, orgId, "", "", "", 0, 100, "", "")
		if err != nil {
			continue
		}

		var deletePolicyIds []string
		for _, policy := range policies.GetPolicies() {
			if strings.HasPrefix(policy.Name, "test-acc") {
				deletePolicyIds = append(deletePolicyIds, policy.Id)
			}
		}

		if len(deletePolicyIds) > 0 {
			_, delErr := scpClient.Client.Organization.DeleteServiceControlPolicies(nil, orgId, deletePolicyIds)
			if delErr != nil {
				return fmt.Errorf("error deleting service control policies for org %s: %v", orgId, delErr)
			}
		}
	}

	return nil
}
