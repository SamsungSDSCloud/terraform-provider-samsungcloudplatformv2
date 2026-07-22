package organization_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client/organization"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDelegationPolicyResourceTest(t *testing.T) {
	offering, err := getTestOffering()
	if err != nil {
		t.Fatalf("failed to get offering: %v", err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDelegationPolicyCreate(offering),
			},
			{
				Config: testAccDelegationPolicyUpdate(offering),
			},
		},
	})
}

func getTestOffering() (string, error) {
	scpClient, err := SharedClientForRegion("kr-west1")
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}

	u, err := url.Parse(scpClient.Client.Organization.Config.AuthUrl)
	if err != nil {
		return "", fmt.Errorf("failed to parse auth url %q: %w", scpClient.Client.Organization.Config.AuthUrl, err)
	}
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected auth url format: %s", scpClient.Client.Organization.Config.AuthUrl)
	}
	return parts[1], nil
}

func getDelegationPolicyBaseConfig() string {
	return `
		data "samsungcloudplatformv2_organization_organizations" "organizations" {
			size = 1
		}

		resource "samsungcloudplatformv2_iam_user" "delegation_policy_user" {
			account_id         = data.samsungcloudplatformv2_organization_organizations.organizations.organizations[0].master_account_id
			description        = "test-acc delegation policy test user"
			password           = "U2NvcmVTY3AhMjM="
			temporary_password = true
			user_name          = "test-acc-delegation-policy-user"
		}
	`
}

func testAccDelegationPolicyCreate(offering string) string {
	return fmt.Sprintf(`
		%s
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
							scp = ["srn:%s::${data.samsungcloudplatformv2_organization_organizations.organizations.organizations[0].master_account_id}:::iam:user/${samsungcloudplatformv2_iam_user.delegation_policy_user.user_id}"]
						}
					}
				]
				version = "2024-07-01"
			}
			organization_id = data.samsungcloudplatformv2_organization_organizations.organizations.organizations[0].id
		}`, getDelegationPolicyBaseConfig(), offering)
}

func testAccDelegationPolicyUpdate(offering string) string {
	return fmt.Sprintf(`
		%s
		resource "samsungcloudplatformv2_organization_delegation_policy" "delegation_policy" {
			document = {
				statement = [
					{
						action = ["organization:CreateServiceControlPolicy"]
						effect = "Allow"
						resource = ["*"]
						sid = "test-acc-sid-1"
						principal = {
							scp = ["srn:%s::${data.samsungcloudplatformv2_organization_organizations.organizations.organizations[0].master_account_id}:::iam:user/${samsungcloudplatformv2_iam_user.delegation_policy_user.user_id}"]
						}
					}
				]
				version = "2024-07-01"
			}
			organization_id = data.samsungcloudplatformv2_organization_organizations.organizations.organizations[0].id
		}`, getDelegationPolicyBaseConfig(), offering)
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
