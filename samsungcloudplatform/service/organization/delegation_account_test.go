package organization_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDelegationAccountResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDelegationAccountCreate(),
			},
			{
				Config:                               testAccDelegationAccountCreate(),
				ResourceName:                         "samsungcloudplatformv2_organization_delegation_account.delegation_account",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "account_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["samsungcloudplatformv2_organization_delegation_account.delegation_account"]
					if !ok {
						return "", fmt.Errorf("resource not found: samsungcloudplatformv2_organization_delegation_account.delegation_account")
					}
					accountId := rs.Primary.Attributes["account_id"]
					orgId := rs.Primary.Attributes["organization_id"]
					serviceType := rs.Primary.Attributes["service_type"]
					return fmt.Sprintf("%s:%s:%s", accountId, orgId, serviceType), nil
				},
			},
		},
	})
}

func getDelegationAccountBaseConfig() string {
	return `
		resource "samsungcloudplatformv2_organization" "org" {
			name       = "temp-org"
			use_scp_yn = true
		}
	`
}

func testAccDelegationAccountCreate() string {
	return fmt.Sprintf(`
		%s
		resource "samsungcloudplatformv2_organization_delegation_account" "delegation_account" {
			account_id      = samsungcloudplatformv2_organization.org.master_account_id
			organization_id = samsungcloudplatformv2_organization.org.id
			service_type    = "identity-center"
		}`, getDelegationAccountBaseConfig())
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_organization_delegation_account", &resource.Sweeper{
		Name: "samsungcloudplatformv2_organization_delegation_account",
		F:    sweepDelegationAccount,
	})
}

func sweepDelegationAccount(region string) error {
	ctx := context.Background()
	scpClient, err := SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	orgList, err := scpClient.Client.Organization.GetOrganizationList(ctx, organization.OrganizationDataSource{})
	if err != nil {
		return nil
	}

	for _, org := range orgList.GetOrganizations() {
		if !strings.HasPrefix(org.GetName(), "test-acc") {
			continue
		}

		accounts, err := scpClient.Client.Organization.ListDelegationAccounts(ctx, org.Id, "", 0, 0, "", "")
		if err != nil {
			continue
		}

		for _, account := range accounts.GetDelegationAccounts() {
			_, delErr := scpClient.Client.Organization.DeleteDelegationAccount(ctx, organization.DelegationAccountResource{
				AccountId:      types.StringValue(account.AccountId),
				OrganizationId: types.StringValue(account.OrganizationId),
				ServiceType:    types.StringValue(account.ServiceType),
			})
			if delErr != nil {
				return fmt.Errorf("error deleting delegation account %s for org %s: %v", account.AccountId, org.Id, delErr)
			}
		}
	}

	return nil
}
