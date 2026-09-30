package smp_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSmpAlertManagerResource(t *testing.T) {
	configData := "YWxlcnRtYW5hZ2VyX2NvbmZpZzogfQo="

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSmpAlertManagerCreate(configData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_smp_alert_manager.am", "alert_manager.id"),
				),
			},
			{
				Config: testAccSmpAlertManagerUpdate(configData),
			},
		},
	})
}

func testAccSmpAlertManagerCreate(configData string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_smp_workspace" "ws" {
			name = "test-acc-ws-for-am"
		}

		resource "samsungcloudplatformv2_smp_notification_group" "ng" {
			name        = "test-acc-ng-for-am"
			description = "test"
		}

		resource "samsungcloudplatformv2_smp_alert_manager" "am" {
			config_data           = "%s"
			notification_group_id = samsungcloudplatformv2_smp_notification_group.ng.notification_group_id
			workspace_id          = samsungcloudplatformv2_smp_workspace.ws.workspace_id
		}
	`, configData)
}

func testAccSmpAlertManagerUpdate(configData string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_smp_workspace" "ws" {
			name = "test-acc-ws-for-am"
		}

		resource "samsungcloudplatformv2_smp_notification_group" "ng" {
			name        = "test-acc-ng-for-am"
			description = "test"
		}

		resource "samsungcloudplatformv2_smp_alert_manager" "am" {
			config_data           = "%s"
			notification_group_id = samsungcloudplatformv2_smp_notification_group.ng.notification_group_id
			workspace_id          = samsungcloudplatformv2_smp_workspace.ws.workspace_id
		}
	`, configData)
}
