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

func TestAccSmpRuleNamespaceResource(t *testing.T) {
	rnName := fmt.Sprintf("test-acc-rn-%s", time.Now().Format("20060102_150405"))
	configData := "Z3JvdXBzOgotIG5hbWU6IHRlc3QtcnVsZQogIHJ1bGVzOgogIC0gYWxlcnQ6IFRlc3RBbGVydAogICAgZXhwcjogdXB0aW1lX3NlY29uZHMgPiAzMDAKICAgIGZvcjogNW0K"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSmpRuleNamespaceCreate(rnName, configData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_smp_rule_namespace.rn", "rule_namespace_id"),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_smp_rule_namespace.rn", "name", rnName),
				),
			},
			{
				Config: testAccSmpRuleNamespaceUpdate(rnName + "-updated", configData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("samsungcloudplatformv2_smp_rule_namespace.rn", "name", rnName+"-updated"),
				),
			},
		},
	})
}

func testAccSmpRuleNamespaceCreate(name, configData string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_smp_workspace" "ws" {
			name = "test-acc-ws-for-rn"
		}

		resource "samsungcloudplatformv2_smp_rule_namespace" "rn" {
			name         = "%s"
			config_data  = "%s"
			workspace_id = samsungcloudplatformv2_smp_workspace.ws.workspace_id
		}
	`, name, configData)
}

func testAccSmpRuleNamespaceUpdate(name, configData string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_smp_workspace" "ws" {
			name = "test-acc-ws-for-rn"
		}

		resource "samsungcloudplatformv2_smp_rule_namespace" "rn" {
			name         = "%s"
			config_data  = "%s"
			workspace_id = samsungcloudplatformv2_smp_workspace.ws.workspace_id
		}
	`, name, configData)
}
