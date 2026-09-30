package smp_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSmpWorkspaceConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSmpWorkspaceConfigurationCreate(30),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("samsungcloudplatformv2_smp_workspace_configuration.wc", "retention_period", "30"),
				),
			},
			{
				Config: testAccSmpWorkspaceConfigurationUpdate(60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("samsungcloudplatformv2_smp_workspace_configuration.wc", "retention_period", "60"),
				),
			},
		},
	})
}

func testAccSmpWorkspaceConfigurationCreate(retention int) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_smp_workspace" "ws" {
			name = "test-acc-ws-for-wc"
		}

		resource "samsungcloudplatformv2_smp_workspace_configuration" "wc" {
			workspace_id     = samsungcloudplatformv2_smp_workspace.ws.workspace_id
			retention_period = %d
		}
	`, retention)
}

func testAccSmpWorkspaceConfigurationUpdate(retention int) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_smp_workspace" "ws" {
			name = "test-acc-ws-for-wc"
		}

		resource "samsungcloudplatformv2_smp_workspace_configuration" "wc" {
			workspace_id     = samsungcloudplatformv2_smp_workspace.ws.workspace_id
			retention_period = %d
		}
	`, retention)
}
