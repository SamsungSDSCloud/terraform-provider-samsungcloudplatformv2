package organization_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationUnitResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationUnitCreate("test-acc-org-unit", "test-acc organization unit description", "o-2b63982e88b74dbcb71ee972b13e2ce1", "r-def2433dbbae4f13913d213250ecb533"),
			},
			{
				Config: testAccOrganizationUnitUpdate("test-acc-org-unit", "test-acc organization unit description updated", "o-2b63982e88b74dbcb71ee972b13e2ce1"),
			},
		},
	})
}

func testAccOrganizationUnitCreate(name string, description string, organizationId string, parentUnitId string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_organization_unit" "organization_unit" {
			name = "%s"
			description = "%s"
			organization_id = "%s"
			parent_unit_id = "%s"
		}`, name, description, organizationId, parentUnitId)
}

func testAccOrganizationUnitUpdate(name string, description string, organizationId string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_organization_unit" "organization_unit" {
			name = "%s"
			description = "%s"
			organization_id = "%s"
		}`, name, description, organizationId)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_organization_unit", &resource.Sweeper{
		Name: "samsungcloudplatformv2_organization_unit",
		F:    sweepOrganizationUnit,
	})
}

func sweepOrganizationUnit(region string) error {
	ctx := context.Background()
	scpClient, err := SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	units, err := scpClient.Client.Organization.GetOrganizationUnits(ctx, "", "", "", "")
	if err != nil {
		return fmt.Errorf("error getting organization units: %v", err)
	}

	var deleteResourceList []string

	for _, unit := range units.GetOrganizationUnits() {
		name := unit.GetName()
		if strings.HasPrefix(name, "test-acc") {
			deleteResourceList = append(deleteResourceList, unit.Id)
		}
	}

	for _, id := range deleteResourceList {
		_, err = scpClient.Client.Organization.DeleteOrganizationUnit(ctx, id, "")
		if err != nil {
			return fmt.Errorf("error deleting organization unit %s: %v", id, err)
		}
	}

	return nil
}
