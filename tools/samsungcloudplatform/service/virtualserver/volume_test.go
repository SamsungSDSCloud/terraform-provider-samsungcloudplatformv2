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

// multiaz: zone(+filter)을 넣어 조회할 때 목록/상세가 정상적으로 나오는지 확인한다.
// (PR 리뷰: zone을 넣으면 volumes 목록이 비어 나오거나, zone+filter 조합에서 volume 상세가 나오지 않던 버그)
// 외부 자원 존재에 의존하지 않도록, 볼륨을 직접 생성한 뒤 그 자원을 zone/filter로 되조회한다.
func TestAccVolumeDataSourceZoneTest(t *testing.T) {
	zone := "kr-west1-y"
	name := "test_terraform_volume_ds_zone"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeDataSourceZone(name, zone),
				Check: resource.ComposeAggregateTestCheckFunc(
					// 목록(plural): zone으로 조회 시 결과가 비어있으면 안 된다.
					resource.TestCheckResourceAttrSet("data.samsungcloudplatformv2_virtualserver_volumes.by_zone", "ids.0"),
					// 상세(single): zone + filter 조합에서도 상세 정보가 나와야 한다.
					resource.TestCheckResourceAttrSet("data.samsungcloudplatformv2_virtualserver_volume.by_zone_filter", "volume.id"),
					resource.TestCheckResourceAttr("data.samsungcloudplatformv2_virtualserver_volume.by_zone_filter", "volume.zone", zone),
				),
			},
		},
	})
}

func testAccVolumeDataSourceZone(name string, zone string) string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_volume" "ds" {
					name = "%s"
					volume_type = "SSD"
					size = 8
					zone = "%s"
				}

				// 목록 조회: zone으로 필터링 (생성한 볼륨을 참조해 조회 순서를 강제한다)
				data "samsungcloudplatformv2_virtualserver_volumes" "by_zone" {
					zone = samsungcloudplatformv2_virtualserver_volume.ds.zone
					name = samsungcloudplatformv2_virtualserver_volume.ds.name
				}

				// 상세 조회: zone + filter 조합
				data "samsungcloudplatformv2_virtualserver_volume" "by_zone_filter" {
					zone = samsungcloudplatformv2_virtualserver_volume.ds.zone
					name = samsungcloudplatformv2_virtualserver_volume.ds.name
					filter {
						name = "volume_type"
						values = ["SSD"]
						use_regex = false
					}
				}`, name, zone)
}

// multiaz: zone은 생성 후 변경할 수 없어야 한다. (volume.go ModifyPlan 에서 immutable 처리)
// 생성 후 zone을 변경하면 apply 이전 plan 단계에서 에러가 발생해야 한다.
func TestAccVolumeZoneImmutableTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// Volume 생성
				Config: testAccVolumeZone("kr-west1-y"),
			},
			{
				// zone 변경 시도 -> immutable 에러
				Config:      testAccVolumeZone("kr-west1-a"),
				ExpectError: regexp.MustCompile(`Immutable fields cannot be modified`),
			},
		},
	})
}

func testAccVolumeZone(zone string) string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_volume" "volume_zone" {
					name = "test_terraform_volume_zone"
					volume_type = "SSD"
					size = 8
					zone = "%s"
				}`, zone)
}

func testAccVolumeCreate() string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_volume" "volume" {
					name = "test_terraform_volume"
					volume_type = "SSD"
					size = 8
                    zone = "kr-west1-y"
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
                    zone = "kr-west1-y"
				  	tags = {
				  	}
				}`)
}

// multiaz + ASG 제약: 볼륨은 볼륨과 동일 zone(kr-west1-y)의 non-ASG 서버에만 attach 가능하다.
// 하지만 servers 데이터소스의 filter는 양성 매칭만 지원하고 ASG 소속 서버(auto_scaling_group_id != null)를
// 제외할 수 없다. 따라서 attach 테스트를 통과시키려면 아래 상수에 "kr-west1-y zone의 attach 가능한
// non-ASG 서버 ID"를 직접 입력해야 한다. (비워두면 데이터소스 ids.0을 쓰며, 그 서버가 ASG VM이면 attach가 400으로 실패한다.)
const attachableServerId = "64fe6123-3906-4d0a-9b4f-635822764d8f" // TODO: kr-west1-y zone의 attach 가능한 non-ASG 서버 ID 입력

func testAccVolumeAttach() string {
	serverIdExpr := "data.samsungcloudplatformv2_virtualserver_servers.candidate.ids.0"
	if attachableServerId != "" {
		serverIdExpr = fmt.Sprintf("%q", attachableServerId)
	}
	return fmt.Sprintf(
		`data "samsungcloudplatformv2_virtualserver_servers" "candidate" {
					// multiaz: 볼륨은 동일 zone의 서버에만 attach 가능하므로 볼륨과 같은 zone의 서버만 후보로 조회한다.
					// (zone이 다르면 attach가 완료(in-use)되지 못하고 실패한다.)
					zone = "kr-west1-y"
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
				}
				resource "samsungcloudplatformv2_virtualserver_volume" "volume" {
					name = "test_terraform_volume_rename"
					volume_type = "SSD"
					size = 16
                    zone = "kr-west1-y"
					servers = [
						{
							id = %s
						}
					]
				  	tags = {
				  	}
				}`, serverIdExpr)
}

func testAccVolumeDetach() string {
	return fmt.Sprintf(
		`resource "samsungcloudplatformv2_virtualserver_volume" "volume" {
					name = "test_terraform_volume_rename"
					volume_type = "SSD"
					size = 16
                    zone = "kr-west1-y"
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
                    zone = "kr-west1-y"
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
                    zone = "kr-west1-y"
					max_iops = 10000
					max_throughput = 500
				}`)
}
