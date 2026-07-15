package virtualserver_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerResourceTest(t *testing.T) {
	name := "test_terraform_server"
	serverType := "s1v1m2"
	securityGroups := "[]"
	networks := `{
    	interface_1 : {
      		subnet_id : data.samsungcloudplatformv2_vpc_subnets.subnets.subnets.0.id,
    	}
	}`
	bootVolume := `{
		type = "SSD",
		size = 104
	}`
	extraVolumes := `{}`
	tags := `{
		"test_terraform_tag_key": "test_terraform_tag_value"
	}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// Server 생성
				Config: testAccServerTemplate("1", name, serverType, securityGroups, networks, bootVolume, extraVolumes, tags, "kr-west1-a"),
			},
			{
				// Server 수정 (name, server type, security group, boot volume size, tags)
				Config: testAccServerTemplate("1", "test_terraform_server_rename", "s1v2m4",
					"[]",
					`{
						interface_1 : {
							subnet_id : data.samsungcloudplatformv2_vpc_subnets.subnets.subnets.0.id,
						},
					}`, `{
							type = "SSD",
							size = 112
						}`,
					`{
						volume_1 : {
						  type = "SSD"
						  size = 8,
						},
					  }`, "{}", "kr-west1-a"),
			},
			{
				// Server 생성 (SSD_Provisioned)
				Config: testAccServerTemplate("2", name+"_2", serverType, securityGroups, networks,
					`{
						type : "SSD_Provisioned",
						size : 104,
						max_iops : 5000,
						max_throughput : 250,
						delete_on_termination: false
					}`,
					extraVolumes, tags, "kr-west1-a"),
			},
			{
				// Server 수정 (QoS)
				Config: testAccServerTemplate("2", name+"_2", serverType, securityGroups, networks,
					`{
						type = "SSD_Provisioned",
						size = 104,
						max_iops = 10000,
						max_throughput = 500,
						delete_on_termination: false
					}`,
					`{
						volume_1 : {
						  size = 8,
						  type = "SSD_Provisioned"
						  max_iops : 5000,
						  max_throughput : 250,
						  delete_on_termination: true
						},
					}`, tags, "kr-west1-a"),
			},
		},
	})
}

// multiaz: zone은 생성 후 변경할 수 없어야 한다. (server.go ModifyPlan 의 immutableFields 에 Zone 포함)
// 생성 후 zone만 변경하면 apply 이전 plan 단계에서 immutable 에러가 발생해야 한다.
func TestAccServerZoneImmutableTest(t *testing.T) {
	name := "test_terraform_server_zone"
	serverType := "s1v1m2"
	securityGroups := "[]"
	networks := `{
    	interface_1 : {
      		subnet_id : data.samsungcloudplatformv2_vpc_subnets.subnets.subnets.0.id,
    	}
	}`
	bootVolume := `{
		type = "SSD",
		size = 104
	}`
	extraVolumes := `{}`
	tags := `{}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// Server 생성 (kr-west1-a)
				Config: testAccServerTemplate("z", name, serverType, securityGroups, networks, bootVolume, extraVolumes, tags, "kr-west1-a"),
			},
			{
				// zone 변경 시도 -> immutable 에러
				Config:      testAccServerTemplate("z", name, serverType, securityGroups, networks, bootVolume, extraVolumes, tags, "kr-west1-b"),
				ExpectError: regexp.MustCompile(`Immutable fields cannot be modified`),
			},
		},
	})
}

func testAccServerTemplate(suffix string, name string, serverType string, securityGroups string, networks string,
	bootVolume string, extraVolumes string, tags string, zone string) string {
	return fmt.Sprintf(
		`
				// 표준 이미지 추출
				data "samsungcloudplatformv2_virtualserver_images" "images" {
					scp_image_type = "standard"
					status = "active"
					name = "Alma 9.6"
				}
				// 키페어 추출
				data "samsungcloudplatformv2_virtualserver_keypairs" "keypairs" {}
				// Security Group 추출
				data "samsungcloudplatformv2_security_group_security_groups" "security_groups" {}
				// Server Group 추출 (partition 제외)
				data "samsungcloudplatformv2_virtualserver_server_groups" "server_groups" {
					filter {
						name = "policy"
						values = ["affinity", "anti-affinity"]
						use_regex = false
					}
				}
				// Subnet 추출 (GENERAL 추출, subnet에 type filter가 적용되어있지 않아서 일단 sort 오름차순으로 확인)
				data "samsungcloudplatformv2_vpc_subnets" "subnets" {
				  state = "ACTIVE"
				  sort = "type:asc"
				}

				resource "samsungcloudplatformv2_virtualserver_server" "server_%s" {
					name            = "%s"
					state           = "ACTIVE"
					image_id        = data.samsungcloudplatformv2_virtualserver_images.images.ids.0
					server_type_id  = "%s"
					keypair_name    = data.samsungcloudplatformv2_virtualserver_keypairs.keypairs.names.0
					security_groups = %s
					server_group_id = data.samsungcloudplatformv2_virtualserver_server_groups.server_groups.ids.0
					networks        = %s
					boot_volume     = %s
					extra_volumes = %s
					tags = %s
                    zone = "%s"
		}`, suffix, name, serverType, securityGroups, networks, bootVolume, extraVolumes, tags, zone)
}
