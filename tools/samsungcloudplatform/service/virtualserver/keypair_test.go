package virtualserver_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeypairResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// Keypair 생성
				Config: testAccKeypairCreate(),
			},
		},
	})
}

func testAccKeypairCreate() string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_keypair" "keypair" {
			name = "test_terraform_keypair"
			tags = {
				"test_terraform_tag_key": "test_terraform_tag_value"
			}
		}`)
}
