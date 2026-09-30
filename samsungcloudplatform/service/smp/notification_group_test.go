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

func TestAccSmpNotificationGroupResource(t *testing.T) {
	ngName := fmt.Sprintf("test-acc-ng-%s", time.Now().Format("20060102_150405"))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSmpNotificationGroupCreate(ngName, "test notification group"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_smp_notification_group.ng", "notification_group_id"),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_smp_notification_group.ng", "name", ngName),
				),
			},
			{
				Config: testAccSmpNotificationGroupUpdate(ngName, "updated notification group"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("samsungcloudplatformv2_smp_notification_group.ng", "description", "updated notification group"),
				),
			},
		},
	})
}

func testAccSmpNotificationGroupCreate(name, description string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_smp_notification_group" "ng" {
			name        = "%s"
			description = "%s"
		}
	`, name, description)
}

func testAccSmpNotificationGroupUpdate(name, description string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_smp_notification_group" "ng" {
			name        = "%s"
			description = "%s"
		}
	`, name, description)
}
