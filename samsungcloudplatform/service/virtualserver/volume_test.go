package virtualserver_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVolumeResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// Volume 생성
				Config: testAccVolumeCreate(),
			},
			{
				// Volume 수정 (Name, Size, Tags)
				Config: testAccVolumeUpdate(),
			},
			{
				// Volume Attach
				Config: testAccVolumeAttach(),
			},
			{
				// Volume Detach
				Config: testAccVolumeDetach(),
			},
			{
				// Volume 생성 (SSD_Provisioned)
				Config: testAccVolumeCreateProvisionedSSD(),
			},
			{
				// Volume Qos 수정
				Config: testAccVolumeUpdateQos(),
			},
		},
	})
}

func testAccVolumeCreate() string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_volume" "volume" {
					name = "test_terraform_volume"
					volume_type = "SSD"
					size = 8
				  	tags = {
						"test_terraform_tag_key": "test_terraform_tag_value"
				  	}
				}`)
}

func testAccVolumeUpdate() string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_volume" "volume" {
					name = "test_terraform_volume_rename"
					volume_type = "SSD"
					size = 16
				  	tags = {
				  	}
				}`)
}

func testAccVolumeAttach() string {
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

					// ASG 서버와 같이 Attach 제약이 있는 서버를 현재 필터링 할 수 없으므로 연결 가능한 ID 값 직접 입력하여 테스트 가능
					//filter {
					//	name = "id"
					//	values = ["uuid"]
					//	use_regex = false
					//}
					
				}
				resource "samsungcloudplatformv2_virtualserver_volume" "volume" {
					name = "test_terraform_volume_rename"
					volume_type = "SSD"
					size = 16
					servers = [
						{
							id = data.samsungcloudplatformv2_virtualserver_servers.candidate.ids.0
						}
					]
				  	tags = {
				  	}
				}`)
}

func testAccVolumeDetach() string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_volume" "volume" {
					name = "test_terraform_volume_rename"
					volume_type = "SSD"
					size = 16
					servers = []
				  	tags = {
				  	}
				}`)
}

func testAccVolumeCreateProvisionedSSD() string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_volume" "volume2" {
					name = "test_terraform_volume"
					volume_type = "SSD_Provisioned"
					size = 8
					max_iops = 5000
					max_throughput = 250
				}`)
}

func testAccVolumeUpdateQos() string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_volume" "volume2" {
					name = "test_terraform_volume"
					volume_type = "SSD_Provisioned"
					size = 8
					max_iops = 10000
					max_throughput = 500
				}`)
}
