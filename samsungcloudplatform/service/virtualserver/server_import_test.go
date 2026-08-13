package virtualserver

// samsungcloudplatformv2_virtualserver_server 리소스의 import 경로 전체 테스트.
//
// import 시 terraform 은 ImportState(=ImportStatePassthroughID) 로 id 만 채운 state 를
// 만든 뒤 곧바로 Read() 를 호출한다. 즉 Read() 가 보는 prior state 는
// "id 를 제외한 모든 속성이 null" 이다. 이 테스트는 그 state 형태를 그대로 재현해서
// MapGetResponseToState() (networks / security_groups / boot_volume / extra_volumes 매핑)
// 가 panic 없이 동작하는지 검증한다.
//
// 실 API 대신 httptest 로 SCP virtualserver API 의 필요한 엔드포인트만 흉내낸다.
// (cfg.Endpoint 를 로컬 서버로 지정하면 SDK 가 그쪽으로 호출한다.)

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"
	"time"

	vsclient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/client"
	scpvirtualserver "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/library/virtualserver/1.4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	testServerId    = "SERVER-aaaaaaaaaaaaaaaaaaaaaa"
	testServerName  = "imported-server"
	testAccountId   = "ACCOUNT-0000"
	testZone        = "kr-west1-a"
	testStateActive = "ACTIVE"

	testBootVolumeId = "VOLUME-boot0000000000000000000"
	testExtraVol1Id  = "VOLUME-extra111111111111111111"
	testExtraVol2Id  = "VOLUME-extra222222222222222222"

	testSubnet1  = "SUBNET-1"
	testSubnet2  = "SUBNET-2"
	testPort1    = "PORT-1"
	testPort2    = "PORT-2"
	testFixedIp1 = "192.168.0.11"
	testFixedIp2 = "192.168.1.11"

	testSg1Id   = "SG-1"
	testSg1Name = "sg-1"

	volumeTypeSSD            = "SSD"
	volumeTypeSSDProvisioned = "SSD_Provisioned"

	deviceOsDisk    = "/dev/vda"
	deviceDataDisk1 = "/dev/vdb"
	deviceDataDisk2 = "/dev/vdc"

	msgImportFailed = "import 매핑 실패: %v"
)

// ---------------------------------------------------------------------------
// mock API
// ---------------------------------------------------------------------------

// mockAPI 는 import(=Read) 경로가 호출하는 엔드포인트들의 응답을 담는다.
type mockAPI struct {
	interfaces     []scpvirtualserver.InterfaceResponseV1Dot2
	securityGroups []scpvirtualserver.SecurityGroupResponse
	serverVolumes  []scpvirtualserver.ServersVolumeResponse
	volumes        []scpvirtualserver.VolumeShowResponseV1Dot4
	defaultVolType string
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

func (m mockAPI) handler(t *testing.T) http.Handler {
	mux := http.NewServeMux()

	// SDK 는 매 호출마다 basePath 로 버전 확인(GET /) 을 한다.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, scpsdk.VersionsResponse{
			Versions: []scpsdk.Version{{ID: "v1.4", Status: "CURRENT"}},
		})
	})

	mux.HandleFunc("GET /v1/servers/{server_id}/interfaces", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, scpvirtualserver.ServerInterfaceListResponseV1Dot2{Interfaces: m.interfaces})
	})

	mux.HandleFunc("GET /v1/servers/{server_id}/security-groups", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, scpvirtualserver.ServerSecurityGroupListResponse{SecurityGroups: m.securityGroups})
	})

	mux.HandleFunc("GET /v1/servers/{server_id}/volumes", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, scpvirtualserver.ServerVolumesResponse{Volumes: m.serverVolumes})
	})

	mux.HandleFunc("GET /v1/volumes", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, scpvirtualserver.VolumeListResponseV1Dot4{Volumes: m.volumes})
	})

	mux.HandleFunc("GET /v1/volume-types/default", func(w http.ResponseWriter, r *http.Request) {
		name := m.defaultVolType
		writeJSON(w, scpvirtualserver.VolumeTypeDetailResponse{
			Id:   "VOLUMETYPE-0000",
			Name: *scpvirtualserver.NewNullableString(&name),
		})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("mock API: 예상하지 못한 호출 %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected path", http.StatusNotFound)
	})

	return mux
}

// newMockedResource: mock API 를 향하는 virtualServerServerResource 를 만든다.
func newMockedResource(t *testing.T, m mockAPI) *virtualServerServerResource {
	t.Helper()

	srv := httptest.NewServer(m.handler(t))
	t.Cleanup(srv.Close)

	cfg := &scpsdk.Configuration{
		Endpoint:        srv.URL,
		AllowSDKVersion: []string{"v1.4"},
		AccountId:       testAccountId,
		Region:          "kr-west1",
		DefaultRegion:   "kr-west1",
		ServiceType:     "virtualserver",
		Credentials:     &scpsdk.Credentials{AccessKey: "ak", SecretKey: "sk"},
		DefaultHeader:   map[string]string{},
		HTTPClient:      srv.Client(),
	}

	return &virtualServerServerResource{
		config: cfg,
		client: vsclient.NewClient(cfg),
	}
}

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

func str(s string) *string { return &s }

// importState: ImportStatePassthroughID 직후의 state. id 만 known, 나머지는 전부 null.
func importState() vsclient.ServerResource {
	return vsclient.ServerResource{
		Id:                    types.StringValue(testServerId),
		AccountId:             types.StringNull(),
		Networks:              types.MapNull(networkObjectType()),
		AutoScalingGroupId:    types.StringNull(),
		CreatedAt:             types.StringNull(),
		CreatedBy:             types.StringNull(),
		DiskConfig:            types.StringNull(),
		ImageId:               types.StringNull(),
		KeypairName:           types.StringNull(),
		LaunchConfigurationId: types.StringNull(),
		Lock:                  types.BoolNull(),
		Metadata:              types.MapNull(types.StringType),
		ModifiedAt:            types.StringNull(),
		Name:                  types.StringNull(),
		PlannedComputeOsType:  types.StringNull(),
		ProductCategory:       types.StringNull(),
		ProductOffering:       types.StringNull(),
		SecurityGroups:        types.ListNull(types.StringType),
		UserData:              types.StringNull(),
		ServerGroupId:         types.StringNull(),
		ServerTypeId:          types.StringNull(),
		State:                 types.StringNull(),
		BootVolume:            types.ObjectNull(volumeObjectType().AttrTypes),
		ExtraVolumes:          types.MapNull(volumeObjectType()),
		VpcId:                 types.StringNull(),
		PartitionNumber:       types.Int32Null(),
		Tags:                  types.MapNull(types.StringType),
		Zone:                  types.StringNull(),
	}
}

func networkObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"subnet_id":     types.StringType,
		"port_id":       types.StringType,
		"fixed_ip":      types.StringType,
		"public_ip_id":  types.StringType,
		"static_nat_id": types.StringType,
		"is_default":    types.BoolType,
	}}
}

func volumeObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":                    types.StringType,
		"delete_on_termination": types.BoolType,
		"size":                  types.Int32Type,
		"type":                  types.StringType,
		"max_iops":              types.Int32Type,
		"max_throughput":        types.Int32Type,
	}}
}

// showResponse: GET /v1/servers/{id} 응답 (Read 가 이미 받아온 상태로 가정).
func showResponse() *scpvirtualserver.ServerShowResponseV1Dot4 {
	now := time.Date(2026, 8, 11, 3, 4, 5, 0, time.UTC)
	osType := scpvirtualserver.PLANNEDCOMPUTEOSTYPE_OPEN_SOURCE
	category := scpvirtualserver.SERVERPRODUCTCATEGORY_COMPUTE
	offering := scpvirtualserver.SERVERPRODUCTOFFERING_VIRTUAL_SERVER

	return &scpvirtualserver.ServerShowResponseV1Dot4{
		Id:                   testServerId,
		AccountId:            testAccountId,
		CreatedAt:            now,
		CreatedBy:            "tester",
		DiskConfig:           "MANUAL",
		ImageId:              *scpvirtualserver.NewNullableString(str("IMAGE-0000")),
		KeypairName:          *scpvirtualserver.NewNullableString(str("test-keypair")),
		Locked:               false,
		Metadata:             map[string]interface{}{},
		ModifiedAt:           now,
		Name:                 testServerName,
		PlannedComputeOsType: *scpvirtualserver.NewNullablePlannedComputeOsType(&osType),
		ProductCategory:      *scpvirtualserver.NewNullableServerProductCategory(&category),
		ProductOffering:      *scpvirtualserver.NewNullableServerProductOffering(&offering),
		ServerGroupId:        *scpvirtualserver.NewNullableString(str("SERVERGROUP-0000")),
		ServerType: scpvirtualserver.ServerShowResponseServerType{
			Id:   *scpvirtualserver.NewNullableString(str("s1v1m2")),
			Name: "s1v1m2",
		},
		State:   testStateActive,
		VpcId:   *scpvirtualserver.NewNullableString(str("VPC-0000")),
		Zone:    testZone,
		Volumes: []scpvirtualserver.ServerShowResponseVolume{},
	}
}

func iface(portId, subnetId, ip string, isDefault bool, withNat bool) scpvirtualserver.InterfaceResponseV1Dot2 {
	itf := scpvirtualserver.InterfaceResponseV1Dot2{
		FixedIps:  []scpvirtualserver.InterfaceAttachmentFixedIp{{IpAddress: ip}},
		IsDefault: isDefault,
		MacAddr:   "fa:16:3e:00:00:01",
		PortId:    portId,
		PortState: testStateActive,
		SubnetId:  subnetId,
	}
	if withNat {
		itf.StaticNat = &scpvirtualserver.PublicStaticNat{
			ExternalIpAddress: "1.2.3.4",
			Id:                "STATICNAT-0000",
			PublicipId:        *scpvirtualserver.NewNullableString(str("PUBLICIP-0000")),
			State:             testStateActive,
		}
	}
	return itf
}

func attached(volumeId, device string, deleteOnTermination bool) scpvirtualserver.ServersVolumeResponse {
	return scpvirtualserver.ServersVolumeResponse{
		DeleteOnTermination: deleteOnTermination,
		Device:              device,
		Id:                  volumeId,
		ServerId:            testServerId,
		VolumeId:            volumeId,
	}
}

func volume(id string, size int32, volType string, iops, throughput *int32) scpvirtualserver.VolumeShowResponseV1Dot4 {
	return scpvirtualserver.VolumeShowResponseV1Dot4{
		Id:            id,
		Size:          size,
		VolumeType:    volType,
		VolumeTypeId:  "VOLUMETYPE-0000",
		State:         "in-use",
		Zone:          testZone,
		MaxIops:       *scpvirtualserver.NewNullableInt32(iops),
		MaxThroughput: *scpvirtualserver.NewNullableInt32(throughput),
		Servers:       []scpvirtualserver.VolumeServer{},
	}
}

// ---------------------------------------------------------------------------
// tests
// ---------------------------------------------------------------------------

// mapImport: import shape state 로 MapGetResponseToState 를 호출한다.
// panic 이 나면 테스트를 실패시키되 panic 내용을 그대로 보고한다.
func mapImport(t *testing.T, m mockAPI) (state vsclient.ServerResource, err error) {
	t.Helper()
	r := newMockedResource(t, m)

	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("PANIC: %v\n%s", rec, debug.Stack())
		}
	}()

	return r.MapGetResponseToState(context.Background(), showResponse(), importState(), types.MapNull(types.StringType))
}

// 케이스 1: NIC 1개 / SG 1개 / boot volume + extra volume 1개.
// 가장 기본적인 VM import 형태.
func TestImportServer_Basic(t *testing.T) {
	m := mockAPI{
		interfaces:     []scpvirtualserver.InterfaceResponseV1Dot2{iface(testPort1, testSubnet1, testFixedIp1, true, false)},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{{Id: testSg1Id, Name: testSg1Name, AccountId: testAccountId}},
		serverVolumes: []scpvirtualserver.ServersVolumeResponse{
			attached(testBootVolumeId, deviceOsDisk, true),
			attached(testExtraVol1Id, deviceDataDisk1, false),
		},
		volumes: []scpvirtualserver.VolumeShowResponseV1Dot4{
			volume(testBootVolumeId, 104, volumeTypeSSD, nil, nil),
			volume(testExtraVol1Id, 8, volumeTypeSSD, nil, nil),
		},
		defaultVolType: volumeTypeSSD,
	}

	state, err := mapImport(t, m)
	if err != nil {
		t.Fatalf(msgImportFailed, err)
	}
	assertImportedState(t, state, 1, 1, 1)
}

// 케이스 2: security_groups 가 비어 있는 서버.
func TestImportServer_NoSecurityGroups(t *testing.T) {
	m := mockAPI{
		interfaces:     []scpvirtualserver.InterfaceResponseV1Dot2{iface(testPort1, testSubnet1, testFixedIp1, true, false)},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{},
		serverVolumes:  []scpvirtualserver.ServersVolumeResponse{attached(testBootVolumeId, deviceOsDisk, true)},
		volumes:        []scpvirtualserver.VolumeShowResponseV1Dot4{volume(testBootVolumeId, 104, volumeTypeSSD, nil, nil)},
		defaultVolType: volumeTypeSSD,
	}

	state, err := mapImport(t, m)
	if err != nil {
		t.Fatalf(msgImportFailed, err)
	}
	assertImportedState(t, state, 1, 0, 0)
}

// 케이스 2-1: SG 가 없어서 panic 은 피하지만, networks / extra_volumes 가
// import state 에 남는지 확인한다. (panic 뒤에 숨어 있는 2차 문제 노출용)
func TestImportServer_NoSecurityGroupsWithNetworksAndExtraVolumes(t *testing.T) {
	m := mockAPI{
		interfaces: []scpvirtualserver.InterfaceResponseV1Dot2{
			iface(testPort1, testSubnet1, testFixedIp1, true, false),
			iface(testPort2, testSubnet2, testFixedIp2, false, false),
		},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{},
		serverVolumes: []scpvirtualserver.ServersVolumeResponse{
			attached(testBootVolumeId, deviceOsDisk, true),
			attached(testExtraVol1Id, deviceDataDisk1, false),
			attached(testExtraVol2Id, deviceDataDisk2, false),
		},
		volumes: []scpvirtualserver.VolumeShowResponseV1Dot4{
			volume(testBootVolumeId, 104, volumeTypeSSD, nil, nil),
			volume(testExtraVol1Id, 8, volumeTypeSSD, nil, nil),
			volume(testExtraVol2Id, 16, volumeTypeSSD, nil, nil),
		},
		defaultVolType: volumeTypeSSD,
	}

	state, err := mapImport(t, m)
	if err != nil {
		t.Fatalf(msgImportFailed, err)
	}
	assertImportedState(t, state, 2, 0, 2)
}

// 케이스 3: SG 여러 개 (박광우 프로 리포트의 실제 형태에 가까움).
func TestImportServer_MultipleSecurityGroups(t *testing.T) {
	m := mockAPI{
		interfaces: []scpvirtualserver.InterfaceResponseV1Dot2{iface(testPort1, testSubnet1, testFixedIp1, true, false)},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{
			{Id: testSg1Id, Name: testSg1Name, AccountId: testAccountId},
			{Id: "SG-2", Name: "sg-2", AccountId: testAccountId},
			{Id: "SG-3", Name: "sg-3", AccountId: testAccountId},
		},
		serverVolumes:  []scpvirtualserver.ServersVolumeResponse{attached(testBootVolumeId, deviceOsDisk, true)},
		volumes:        []scpvirtualserver.VolumeShowResponseV1Dot4{volume(testBootVolumeId, 104, volumeTypeSSD, nil, nil)},
		defaultVolType: volumeTypeSSD,
	}

	state, err := mapImport(t, m)
	if err != nil {
		t.Fatalf(msgImportFailed, err)
	}
	assertImportedState(t, state, 1, 3, 0)
}

// 케이스 4: NIC 여러 개 + public ip(static nat) 포함.
func TestImportServer_MultipleNetworks(t *testing.T) {
	m := mockAPI{
		interfaces: []scpvirtualserver.InterfaceResponseV1Dot2{
			iface(testPort1, testSubnet1, testFixedIp1, true, true),
			iface(testPort2, testSubnet2, testFixedIp2, false, false),
		},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{{Id: testSg1Id, Name: testSg1Name, AccountId: testAccountId}},
		serverVolumes:  []scpvirtualserver.ServersVolumeResponse{attached(testBootVolumeId, deviceOsDisk, true)},
		volumes:        []scpvirtualserver.VolumeShowResponseV1Dot4{volume(testBootVolumeId, 104, volumeTypeSSD, nil, nil)},
		defaultVolType: volumeTypeSSD,
	}

	state, err := mapImport(t, m)
	if err != nil {
		t.Fatalf(msgImportFailed, err)
	}
	assertImportedState(t, state, 2, 1, 0)
}

// 케이스 5: extra volume 여러 개 + QoS(SSD_Provisioned) 볼륨.
func TestImportServer_ExtraVolumesWithQos(t *testing.T) {
	iops := int32(5000)
	throughput := int32(250)

	m := mockAPI{
		interfaces:     []scpvirtualserver.InterfaceResponseV1Dot2{iface(testPort1, testSubnet1, testFixedIp1, true, false)},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{{Id: testSg1Id, Name: testSg1Name, AccountId: testAccountId}},
		serverVolumes: []scpvirtualserver.ServersVolumeResponse{
			attached(testBootVolumeId, deviceOsDisk, true),
			attached(testExtraVol1Id, deviceDataDisk1, false),
			attached(testExtraVol2Id, deviceDataDisk2, true),
		},
		volumes: []scpvirtualserver.VolumeShowResponseV1Dot4{
			volume(testBootVolumeId, 104, volumeTypeSSDProvisioned, &iops, &throughput),
			volume(testExtraVol1Id, 8, volumeTypeSSD, nil, nil),
			volume(testExtraVol2Id, 16, volumeTypeSSDProvisioned, &iops, &throughput),
		},
		defaultVolType: volumeTypeSSD,
	}

	state, err := mapImport(t, m)
	if err != nil {
		t.Fatalf(msgImportFailed, err)
	}
	assertImportedState(t, state, 1, 1, 2)
}

// 케이스 6: extra volume 이 없는 서버 (boot volume 만).
func TestImportServer_BootVolumeOnly(t *testing.T) {
	m := mockAPI{
		interfaces:     []scpvirtualserver.InterfaceResponseV1Dot2{iface(testPort1, testSubnet1, testFixedIp1, true, false)},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{{Id: testSg1Id, Name: testSg1Name, AccountId: testAccountId}},
		serverVolumes:  []scpvirtualserver.ServersVolumeResponse{attached(testBootVolumeId, "/dev/sda", true)},
		volumes:        []scpvirtualserver.VolumeShowResponseV1Dot4{volume(testBootVolumeId, 104, volumeTypeSSD, nil, nil)},
		defaultVolType: volumeTypeSSD,
	}

	state, err := mapImport(t, m)
	if err != nil {
		t.Fatalf(msgImportFailed, err)
	}
	assertImportedState(t, state, 1, 1, 0)
}

// 케이스 8: 실제 환경(testjihof/terraform.tfstate 의 tf-bastion-lab)에 만들어져 있던
// 서버의 실제 형태를 그대로 재현한 케이스.
//   security_groups: 1개, networks: NIC 1개 + static nat(public ip),
//   boot_volume: 16GB SSD, extra_volumes: 없음, metadata: {"HA_Enabled":"True"}
// 즉 "특별할 것 없는 평범한 VM" 을 import 하면 어떻게 되는지 보여주는 케이스.
func TestImportServer_RealWorldBastionShape(t *testing.T) {
	m := mockAPI{
		interfaces: []scpvirtualserver.InterfaceResponseV1Dot2{
			iface("1440f1c7-f21e-4e95-a4f8-bf9e38d743d8", "183c6c30df0c4fcbb16422aa2d64aa21", "192.168.0.184", true, true),
		},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{
			{Id: "6125c29c-6a0b-4dc8-861c-0bc7430ca2b4", Name: "bastion-sg", AccountId: testAccountId},
		},
		serverVolumes: []scpvirtualserver.ServersVolumeResponse{
			attached("d296a27e-ab67-4ac3-9ffc-4a91226f50ae", deviceOsDisk, false),
		},
		volumes: []scpvirtualserver.VolumeShowResponseV1Dot4{
			volume("d296a27e-ab67-4ac3-9ffc-4a91226f50ae", 16, volumeTypeSSD, nil, nil),
		},
		defaultVolType: volumeTypeSSD,
	}

	r := newMockedResource(t, m)
	resp := showResponse()
	resp.Metadata = map[string]interface{}{"HA_Enabled": "True"}

	var state vsclient.ServerResource
	var err error
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				err = fmt.Errorf("PANIC: %v", rec)
			}
		}()
		state, err = r.MapGetResponseToState(context.Background(), resp, importState(), types.MapNull(types.StringType))
	}()

	if err != nil {
		t.Fatalf("평범한 VM(SG 1개 + NIC 1개) import 가 깨진다: %v", err)
	}
	assertImportedState(t, state, 1, 1, 0)
}

// 케이스 7: API metadata 에 string 이 아닌 값이 섞여 있는 경우.
// processMetadata 가 v.(string) 로 무조건 단정하므로 방어 로직이 없으면 panic 한다.
// (SG 를 비워 1차 panic 을 피한 뒤 metadata 단독으로 확인)
func TestImportServer_NonStringMetadata(t *testing.T) {
	m := mockAPI{
		interfaces:     []scpvirtualserver.InterfaceResponseV1Dot2{iface(testPort1, testSubnet1, testFixedIp1, true, false)},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{},
		serverVolumes:  []scpvirtualserver.ServersVolumeResponse{attached(testBootVolumeId, deviceOsDisk, true)},
		volumes:        []scpvirtualserver.VolumeShowResponseV1Dot4{volume(testBootVolumeId, 104, volumeTypeSSD, nil, nil)},
		defaultVolType: volumeTypeSSD,
	}

	r := newMockedResource(t, m)
	resp := showResponse()
	resp.Metadata = map[string]interface{}{"os_type": "linux", "attempt": float64(1)}

	var err error
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				err = fmt.Errorf("PANIC: %v", rec)
			}
		}()
		_, err = r.MapGetResponseToState(context.Background(), resp, importState(), types.MapNull(types.StringType))
	}()

	if err != nil {
		t.Fatalf("metadata 가 string 이 아닌 값을 포함하면 import 가 깨진다: %v", err)
	}
}

// 케이스 9: import 직후 다시 Read 해도 키/값이 그대로여야 한다.
//
// import 시 생성한 키(interface_1, volume_1 ...)가 다음 refresh 에서 다른 키로 바뀌면
// 매 plan 마다 "원소 삭제 + 추가" diff 가 생긴다. 두 번 매핑해서 동일한지 본다.
func TestImportServer_SecondReadIsStable(t *testing.T) {
	iops := int32(5000)
	throughput := int32(250)

	m := mockAPI{
		interfaces: []scpvirtualserver.InterfaceResponseV1Dot2{
			iface(testPort1, testSubnet1, testFixedIp1, true, true),
			iface(testPort2, testSubnet2, testFixedIp2, false, false),
		},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{
			{Id: testSg1Id, Name: testSg1Name, AccountId: testAccountId},
			{Id: "SG-2", Name: "sg-2", AccountId: testAccountId},
		},
		serverVolumes: []scpvirtualserver.ServersVolumeResponse{
			attached(testBootVolumeId, deviceOsDisk, true),
			attached(testExtraVol1Id, deviceDataDisk1, false),
			attached(testExtraVol2Id, deviceDataDisk2, true),
		},
		volumes: []scpvirtualserver.VolumeShowResponseV1Dot4{
			volume(testBootVolumeId, 104, volumeTypeSSD, nil, nil),
			volume(testExtraVol1Id, 8, volumeTypeSSD, nil, nil),
			volume(testExtraVol2Id, 16, volumeTypeSSDProvisioned, &iops, &throughput),
		},
		defaultVolType: volumeTypeSSD,
	}

	r := newMockedResource(t, m)

	// 1회차: import (prior state 전부 null)
	first, err := r.MapGetResponseToState(context.Background(), showResponse(), importState(), types.MapNull(types.StringType))
	if err != nil {
		t.Fatalf("1회차 import 실패: %v", err)
	}
	assertImportedState(t, first, 2, 2, 2)

	// 2회차: 1회차 결과를 prior state 로 다시 Read
	second, err := r.MapGetResponseToState(context.Background(), showResponse(), first, types.MapNull(types.StringType))
	if err != nil {
		t.Fatalf("2회차 refresh 실패: %v", err)
	}

	if !first.Networks.Equal(second.Networks) {
		t.Errorf("networks 가 재조회에서 바뀐다 (매 plan 마다 diff 발생)\n1회차: %v\n2회차: %v",
			first.Networks, second.Networks)
	}
	if !first.ExtraVolumes.Equal(second.ExtraVolumes) {
		t.Errorf("extra_volumes 가 재조회에서 바뀐다\n1회차: %v\n2회차: %v",
			first.ExtraVolumes, second.ExtraVolumes)
	}
	if !first.SecurityGroups.Equal(second.SecurityGroups) {
		t.Errorf("security_groups 가 재조회에서 바뀐다\n1회차: %v\n2회차: %v",
			first.SecurityGroups, second.SecurityGroups)
	}
	if !first.BootVolume.Equal(second.BootVolume) {
		t.Errorf("boot_volume 가 재조회에서 바뀐다\n1회차: %v\n2회차: %v",
			first.BootVolume, second.BootVolume)
	}
}

// assertImportedState: import 후 state 에 필수 속성이 실제로 채워졌는지 검증.
func assertImportedState(t *testing.T, state vsclient.ServerResource, wantNetworks, wantSGs, wantExtraVolumes int) {
	t.Helper()

	if state.Id.ValueString() != testServerId {
		t.Errorf("id: got %q, want %q", state.Id.ValueString(), testServerId)
	}
	if state.Name.ValueString() != testServerName {
		t.Errorf("name: got %q, want %q", state.Name.ValueString(), testServerName)
	}
	if state.Zone.ValueString() != testZone {
		t.Errorf("zone: got %q, want %q", state.Zone.ValueString(), testZone)
	}

	// networks
	if state.Networks.IsNull() || state.Networks.IsUnknown() {
		t.Errorf("networks: null/unknown 이면 import 후 plan 에서 diff 가 발생한다")
	} else if got := len(state.Networks.Elements()); got != wantNetworks {
		t.Errorf("networks: got %d개, want %d개", got, wantNetworks)
	}

	// security_groups
	if state.SecurityGroups.IsNull() || state.SecurityGroups.IsUnknown() {
		t.Errorf("security_groups: null/unknown")
	} else if got := len(state.SecurityGroups.Elements()); got != wantSGs {
		t.Errorf("security_groups: got %d개, want %d개", got, wantSGs)
	} else {
		for i, e := range state.SecurityGroups.Elements() {
			if e == nil || e.IsNull() {
				t.Errorf("security_groups[%d]: null 원소 (state 에 SG id 가 안 채워짐)", i)
			}
		}
	}

	// boot_volume
	if state.BootVolume.IsNull() || state.BootVolume.IsUnknown() {
		t.Errorf("boot_volume: null/unknown")
	} else if id, ok := state.BootVolume.Attributes()["id"]; !ok || id.IsNull() {
		t.Errorf("boot_volume.id 가 채워지지 않음: %v", state.BootVolume)
	}

	// extra_volumes
	if wantExtraVolumes > 0 {
		if state.ExtraVolumes.IsNull() || state.ExtraVolumes.IsUnknown() {
			t.Errorf("extra_volumes: null/unknown (실제 %d개 붙어 있음)", wantExtraVolumes)
		} else if got := len(state.ExtraVolumes.Elements()); got != wantExtraVolumes {
			t.Errorf("extra_volumes: got %d개, want %d개", got, wantExtraVolumes)
		}
	}

	// import 로 새로 만들어지는 키는 interface_1.., volume_1.. 이어야 한다.
	assertGeneratedKeys(t, "networks", state.Networks, "interface", wantNetworks)
	assertGeneratedKeys(t, "extra_volumes", state.ExtraVolumes, "volume", wantExtraVolumes)
}

// assertGeneratedKeys: 맵 키가 prefix_1 .. prefix_N 집합과 정확히 일치하는지 본다.
func assertGeneratedKeys(t *testing.T, attrName string, m types.Map, prefix string, want int) {
	t.Helper()
	if want == 0 || m.IsNull() || m.IsUnknown() {
		return
	}

	expected := make(map[string]bool, want)
	for i := 1; i <= want; i++ {
		expected[fmt.Sprintf("%s_%d", prefix, i)] = true
	}

	for key := range m.Elements() {
		if !expected[key] {
			t.Errorf("%s: 예상 못한 키 %q (기대: %s_1..%s_%d)", attrName, key, prefix, prefix, want)
			continue
		}
		delete(expected, key)
	}
	for key := range expected {
		t.Errorf("%s: 키 %q 누락", attrName, key)
	}
}
