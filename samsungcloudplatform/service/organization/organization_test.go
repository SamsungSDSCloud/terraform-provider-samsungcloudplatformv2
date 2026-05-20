package organization_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationCreate("test-acc-org", true),
			},
			{
				Config: testAccOrganizationUpdate("test-acc-org-updated", false),
			},
		},
	})
}

func testAccOrganizationCreate(name string, useScpYn bool) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_organization" "org" {
			name                  = "%s"
			use_scp_yn            = %t
		}`, name, useScpYn)
}

func testAccOrganizationUpdate(name string, useScpYn bool) string {
	return fmt.Sprintf(`
resource "samsungcloudplatformv2_organization" "org" {
  name       = "%s"
  use_scp_yn = %t
}`, name, useScpYn)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_organization", &resource.Sweeper{
		Name: "samsungcloudplatformv2_organization",
		F:    sweepOrganization,
	})
}

func sweepOrganization(region string) error {
	ctx := context.Background()
	scpClient, err := SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	orgs, err := scpClient.Client.Organization.GetOrganizationList(ctx, organization.OrganizationDataSource{})
	if err != nil {
		return fmt.Errorf("error getting organizations: %v", err)
	}

	var deleteResourceList []string

	for _, org := range orgs.GetOrganizations() {
		name := org.GetName()
		if strings.HasPrefix(name, "test-acc") {
			deleteResourceList = append(deleteResourceList, org.Id)
		}
	}

	for _, id := range deleteResourceList {
		_, err = scpClient.Client.Organization.DeleteOrganization(ctx, id)
		if err != nil {
			return fmt.Errorf("error deleting organization %s: %v", id, err)
		}
	}

	return nil
}
