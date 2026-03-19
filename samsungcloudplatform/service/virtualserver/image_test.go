package virtualserver_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImageResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// Instance를 통한 Image 생성
				Config: testAccImageCreate(),
			},
		},
	})
}

func getServerCandidate() string {
	// state가 ACTIVE & product_offering이 virtual_server인 서버 ID 추출
	return fmt.Sprintf(
		`data "samsungcloudplatformv2_virtualserver_servers" "candidate" {
			filter {
				name = "product_offering"
				values = ["virtual_server"]
				use_regex = false
			}
			filter {
				name = "state"
				values = ["ACTIVE"]
				use_regex = false
			}
		}`)
}

func testAccImageCreate() string {
	return fmt.Sprintf(
		`%s
			resource "samsungcloudplatformv2_virtualserver_image" "image" {
				name = "test_terraform_image"
				instance_id = data.samsungcloudplatformv2_virtualserver_servers.candidate.ids.0
			}`, getServerCandidate())
}
