package virtualserver_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerGroupResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// Server Group 생성 (affinity)
				Config: testAccServerGroupCreate("1", "affinity"),
			},
			{
				// Server Group 생성 (anti-affinity)
				Config: testAccServerGroupCreate("2", "anti-affinity"),
			},
			{
				// Server Group 생성 (partition)
				Config: testAccServerGroupCreate("3", "partition"),
			},
		},
	})
}

func testAccServerGroupCreate(suffix string, policy string) string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_server_group" "server_group_%s" {
			name = "test_terraform_server_group_%s"
			policy = "%s"
			tags = {
				"test_terraform_tag_key": "test_terraform_tag_value"
			}
		}`, suffix, suffix, policy)
}
