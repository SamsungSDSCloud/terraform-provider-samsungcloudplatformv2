package smp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSmpWorkspaceResource(t *testing.T) {
	workspaceName := fmt.Sprintf("test-acc-workspace-%s", time.Now().Format("20060102_150405"))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSmpWorkspaceCreate(workspaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_smp_workspace.workspace", "workspace_id"),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_smp_workspace.workspace", "name", workspaceName),
				),
			},
			{
				Config: testAccSmpWorkspaceUpdate(workspaceName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("samsungcloudplatformv2_smp_workspace.workspace", "name", workspaceName+"-updated"),
				),
			},
		},
	})
}

func testAccSmpWorkspaceCreate(name string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_smp_workspace" "workspace" {
			name = "%s"
		}
	`, name)
}

func testAccSmpWorkspaceUpdate(name string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_smp_workspace" "workspace" {
			name = "%s"
		}
	`, name)
}
