package virtualserver_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform"
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

// import 테스트: 서버를 만든 뒤 terraform import 로 state 를 복원한다.
//
// import 는 ImportStatePassthroughID 로 id 만 채운 state 를 만들고 Read() 를 호출하므로
// Read() 가 보는 prior state 는 "id 외 전부 null" 이다. 이 상태에서
// networks / security_groups / boot_volume / extra_volumes 가 제대로 복원되는지 검증한다.
//
// ImportStateVerify 는 import 된 state 를 생성 시 state 와 전부 비교한다.
// user_data 는 API 가 돌려주지 않아(Read 가 prior state 를 그대로 사용) import 시 null 이 되므로 제외한다.
func TestAccServerImportTest(t *testing.T) {
	const resourceName = "samsungcloudplatformv2_virtualserver_server.server_imp"

	networks := `{
    	interface_1 : {
      		subnet_id : data.samsungcloudplatformv2_vpc_subnets.subnets.subnets.0.id,
    	}
	}`
	bootVolume := `{
		type = "SSD",
		size = 104
	}`
	tags := `{
		"test_terraform_tag_key": "test_terraform_tag_value"
	}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// 1) 기본 형태 (NIC 1개, SG 없음, extra volume 없음)
				Config: testAccServerTemplate("imp", "test_terraform_server_import", "s1v1m2",
					"[]", networks, bootVolume, `{}`, tags, "kr-west1-a"),
			},
			{
				// 2) import (기본 형태)
				Config:                  testAccServerTemplate("imp", "test_terraform_server_import", "s1v1m2", "[]", networks, bootVolume, `{}`, tags, "kr-west1-a"),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_data"},
			},
			{
				// 3) security group + extra volume 추가 (리포트된 실패 형태)
				Config: testAccServerTemplate("imp", "test_terraform_server_import", "s1v1m2",
					"[data.samsungcloudplatformv2_security_group_security_groups.security_groups.ids.0]",
					networks, bootVolume,
					`{
						volume_1 : {
						  type = "SSD"
						  size = 8,
						},
					}`, tags, "kr-west1-a"),
			},
			{
				// 4) import (security group + extra volume 포함)
				Config: testAccServerTemplate("imp", "test_terraform_server_import", "s1v1m2",
					"[data.samsungcloudplatformv2_security_group_security_groups.security_groups.ids.0]",
					networks, bootVolume,
					`{
						volume_1 : {
						  type = "SSD"
						  size = 8,
						},
					}`, tags, "kr-west1-a"),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_data"},
			},
			{
				// 5) import 후 plan 이 비어 있어야 한다 (import 로 복원된 state 에 drift 가 없어야 함)
				Config: testAccServerTemplate("imp", "test_terraform_server_import", "s1v1m2",
					"[data.samsungcloudplatformv2_security_group_security_groups.security_groups.ids.0]",
					networks, bootVolume,
					`{
						volume_1 : {
						  type = "SSD"
						  size = 8,
						},
					}`, tags, "kr-west1-a"),
				PlanOnly: true,
			},
		},
	})
}

// NIC 추가 테스트.
//
// interface 추가/변경 시 IP 자동 할당이 없어졌고 fixed_ip 가 고정 입력이다.
// subnet_id 만 주고 NIC 을 추가하면 API 가 400 을 돌려준다.
// (참고: server.go 의 CreateServerInterface 는 fixed_ip 가 set 일 때만 fixed_ips 를 실어 보낸다)
//
// 사용할 IP 는 환경마다 달라서 하드코딩하면 깨지므로 SCP_FIXED_IP 로 받는다.
// 미지정이면 skip 한다.
func TestAccServerNetworkAddNicTest(t *testing.T) {
	fixedIp := os.Getenv("SCP_FIXED_IP")
	if fixedIp == "" {
		t.Skip("NIC 추가 테스트는 subnet 내 미사용 IP 가 필요하다. SCP_FIXED_IP 를 지정할 것")
	}

	name := "test_terraform_server_nic"
	serverType := "s1v1m2"
	securityGroups := "[]"
	bootVolume := `{
		type = "SSD",
		size = 104
	}`

	oneNic := `{
    	interface_1 : {
      		subnet_id : data.samsungcloudplatformv2_vpc_subnets.subnets.subnets.0.id,
    	}
	}`
	// NIC 추가 시에는 fixed_ip 를 반드시 명시한다.
	twoNics := fmt.Sprintf(`{
    	interface_1 : {
      		subnet_id : data.samsungcloudplatformv2_vpc_subnets.subnets.subnets.0.id,
    	},
    	interface_2 : {
      		subnet_id : data.samsungcloudplatformv2_vpc_subnets.subnets.subnets.0.id,
      		fixed_ip  : "%s",
    	}
	}`, fixedIp)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// NIC 1개로 생성
				Config: testAccServerTemplate("nic", name, serverType, securityGroups, oneNic, bootVolume, `{}`, `{}`, "kr-west1-a"),
			},
			{
				// interface_2 추가 (fixed_ip 명시) — 이것만 바꾼다
				Config: testAccServerTemplate("nic", name, serverType, securityGroups, twoNics, bootVolume, `{}`, `{}`, "kr-west1-a"),
			},
			{
				// NIC 2개 상태로 import 해도 networks 가 그대로 복원되어야 한다
				Config:                  testAccServerTemplate("nic", name, serverType, securityGroups, twoNics, bootVolume, `{}`, `{}`, "kr-west1-a"),
				ResourceName:            "samsungcloudplatformv2_virtualserver_server.server_nic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_data"},
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
