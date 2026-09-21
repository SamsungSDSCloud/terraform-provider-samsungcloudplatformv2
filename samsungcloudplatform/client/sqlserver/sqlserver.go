package sqlserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/database"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/sqlserver/1.2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *sqlserver.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: sqlserver.NewAPIClient(config),
	}
}

// clusterlist (ctx, ClusterDataSource) - (RdbClusterPageResponse)
func (client *Client) GetClusterList(ctx context.Context, request ClusterDataSource) (*sqlserver.RdbClusterPageResponse, error) {
	req := client.sdkClient.SqlserverV1SqlserverClustersApiAPI.SqlserverListClusters(ctx)
	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.ServiceState.IsNull() {
		req = req.ServiceState(request.ServiceState.ValueString())
	}
	if !request.DatabaseName.IsNull() {
		req = req.DatabaseName(request.DatabaseName.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

// engine version
// 비어있지 않은 필터만 쿼리 파라미터로 전달한다.
// eosIncluded가 nil이면 파라미터를 생략하여 서버 기본값을 따른다.
func (client *Client) GetEngineVersionList(ctx context.Context, id string, productImageType string, eosIncluded *bool) (*sqlserver.EngineListResponse, error) {
	req := client.sdkClient.SqlserverV1SqlserverMasterDataApiAPI.SqlserverListEngineVersions(ctx)
	if id != "" {
		req = req.Id(id)
	}
	if productImageType != "" {
		req = req.ProductImageType(sqlserver.ProductImageType(productImageType))
	}
	if eosIncluded != nil {
		req = req.EosIncluded(*eosIncluded)
	}

	resp, _, err := req.Execute()
	return resp, err
}

// create (ctx, clusterResource) - (asyncResponse)
func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*sqlserver.AsyncResponse, error) {
	req := client.sdkClient.SqlserverV1SqlserverClustersApiAPI.SqlserverCreateCluster(ctx)

	// AllowableIpAddresses
	var allowableIpAddresses []string
	if request.AllowableIpAddresses.IsNull() || request.AllowableIpAddresses.IsUnknown() {
		allowableIpAddresses = []string{}
	} else {
		for _, elem := range request.AllowableIpAddresses.Elements() {
			strVal := elem.(types.String)
			allowableIpAddresses = append(allowableIpAddresses, strVal.ValueString())
		}
	}

	// InitConfigOption
	var initConfigOption = request.InitConfigOption
	var convertedBackupOption = &sqlserver.SqlserverBackupOption{}

	if client.CheckBackupConfig(initConfigOption) {
		convertedBackupOption = nil
	} else {
		// FullBackupDayOfWeek enum
		dayOfWeek, err := sqlserver.NewDayOfWeekFromValue(initConfigOption.BackupOption.FullBackupDayOfWeek.ValueString())
		if err != nil {
			return nil, err
		}

		convertedBackupOption = &sqlserver.SqlserverBackupOption{
			ArchiveFrequencyMinute: initConfigOption.BackupOption.ArchiveFrequencyMinute.ValueString(),
			FullBackupDayOfWeek:    dayOfWeek,
			RetentionPeriodDay:     *sqlserver.NewNullableString(initConfigOption.BackupOption.RetentionPeriodDay.ValueStringPointer()),
			StartingTimeHour:       *sqlserver.NewNullableString(initConfigOption.BackupOption.StartingTimeHour.ValueStringPointer()),
		}
	}

	// DatabaseCollation enum
	databaseCollation, err := sqlserver.NewDbCollationFromValue(initConfigOption.DatabaseCollation.ValueString())
	if err != nil {
		return nil, err
	}

	//
	var convertedDatabases []sqlserver.SqlserverDatabaseOption
	for _, database := range initConfigOption.Databases {
		convertedDatabases = append(convertedDatabases, sqlserver.SqlserverDatabaseOption{
			DatabaseName: database.DatabaseName.ValueString(),
			DriveLetter:  database.DriveLetter.ValueStringPointer(),
		})
	}

	var convertedInitConfigOption = sqlserver.SqlserverInitConfigOptionRequest{
		AdEnabled:            initConfigOption.AdEnabled.ValueBoolPointer(),
		AuditEnabled:         initConfigOption.AuditEnabled.ValueBoolPointer(),
		BackupOption:         *sqlserver.NewNullableSqlserverBackupOption(convertedBackupOption),
		DatabaseCollation:    databaseCollation,
		DatabasePort:         initConfigOption.DatabasePort.ValueInt32Pointer(),
		DatabaseServiceName:  initConfigOption.DatabaseServiceName.ValueString(),
		DatabaseUserName:     initConfigOption.DatabaseUserName.ValueString(),
		DatabaseUserPassword: initConfigOption.DatabaseUserPassword.ValueString(),
		Databases:            convertedDatabases,
		License:              initConfigOption.License.ValueString(),
	}

	// AdConfig는 ad_enabled가 true일 때만 필요하다.
	if initConfigOption.AdEnabled.ValueBool() {
		if initConfigOption.AdConfig == nil {
			return nil, fmt.Errorf("init_config_option.ad_config is required when init_config_option.ad_enabled is true")
		}
		adConfig := sqlserver.SqlserverAdConfigRequest{
			AdDomainName:        initConfigOption.AdConfig.AdDomainName.ValueString(),
			AdNetbiosName:       initConfigOption.AdConfig.AdNetbiosName.ValueString(),
			AdUserId:            initConfigOption.AdConfig.AdUserId.ValueString(),
			AdUserPassword:      initConfigOption.AdConfig.AdUserPassword.ValueString(),
			FailoverClusterName: *sqlserver.NewNullableString(initConfigOption.AdConfig.FailoverClusterName.ValueStringPointer()),
		}
		for _, dns := range initConfigOption.AdConfig.AdDnsServers.Elements() {
			adConfig.AdDnsServers = append(adConfig.AdDnsServers, dns.(types.String).ValueString())
		}
		convertedInitConfigOption.AdConfig = *sqlserver.NewNullableSqlserverAdConfigRequest(&adConfig)
	} else if initConfigOption.AdConfig != nil {
		return nil, fmt.Errorf("init_config_option.ad_config must not be set when init_config_option.ad_enabled is false")
	}

	// InstanceGroups
	var convertedInstanceGroups []sqlserver.SqlserverInstanceGroupRequest
	var igVals []database.InstanceGroup
	request.InstanceGroups.ElementsAs(ctx, &igVals, false)
	for _, instanceGroup := range igVals {
		var convertedBlockStorage []sqlserver.SqlserverBlockStorageGroupRequest
		var bsVals []database.BlockStorageGroup
		instanceGroup.BlockStorageGroups.ElementsAs(ctx, &bsVals, false)
		for _, blockStorage := range bsVals {
			convertedBlockStorage = append(convertedBlockStorage, sqlserver.SqlserverBlockStorageGroupRequest{
				RoleType:   sqlserver.OsDataBlockStorageGroupRoleType(blockStorage.RoleType.ValueString()),
				SizeGb:     blockStorage.SizeGb.ValueInt32(),
				VolumeType: sqlserver.VolumeType(blockStorage.VolumeType.ValueString()).Ptr(),
			})
		}

		var convertedInstance []sqlserver.SqlserverInstanceRequest
		var instVals []database.Instance
		instanceGroup.Instances.ElementsAs(ctx, &instVals, false)
		for _, instance := range instVals {
			convertedInstance = append(convertedInstance, sqlserver.SqlserverInstanceRequest{
				RoleType:         sqlserver.SqlserverInstanceRoleType(instance.RoleType.ValueString()),
				ServiceIpAddress: *sqlserver.NewNullableString(instance.ServiceIpAddress.ValueStringPointer()),
				PublicIpId:       *sqlserver.NewNullableString(instance.PublicIpId.ValueStringPointer()),
			})
		}

		convertedInstanceGroups = append(convertedInstanceGroups, sqlserver.SqlserverInstanceGroupRequest{
			BlockStorageGroups: convertedBlockStorage,
			Instances:          convertedInstance,
			RoleType:           sqlserver.SqlserverInstanceGroupRoleType(instanceGroup.RoleType.ValueString()),
			ServerTypeName:     instanceGroup.ServerTypeName.ValueString(),
		})
	}

	// MaintenanceOption
	var maintenanceOption = request.MaintenanceOption
	var convertedMaintenanceOption = &sqlserver.MaintenanceOption{}

	if client.CheckMaintenanceOption(maintenanceOption) {
		convertedMaintenanceOption = nil
	} else {
		var startingDayOfWeek, _ = sqlserver.NewDayOfWeekFromValue(maintenanceOption.StartingDayOfWeek.ValueString())
		convertedMaintenanceOption = &sqlserver.MaintenanceOption{
			PeriodHour:        maintenanceOption.PeriodHour.ValueStringPointer(),
			StartingTime:      maintenanceOption.StartingTime.ValueStringPointer(),
			StartingDayOfWeek: startingDayOfWeek,
		}
	}

	//Tags
	var TagsObject []sqlserver.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := sqlserver.Tag{
			Key:   &k,
			Value: *sqlserver.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.SqlserverClusterCreateRequestV1Dot1(sqlserver.SqlserverClusterCreateRequestV1Dot1{
		AllowableIpAddresses:      allowableIpAddresses,
		DbaasEngineVersionId:      request.DbaasEngineVersionId.ValueString(),
		NatEnabled:                request.NatEnabled.ValueBoolPointer(),
		HaEnabled:                 request.HaEnabled.ValueBoolPointer(),
		InitConfigOption:          convertedInitConfigOption,
		InstanceGroups:            convertedInstanceGroups,
		InstanceNamePrefix:        request.InstanceNamePrefix.ValueString(),
		Name:                      request.Name.ValueString(),
		SubnetId:                  request.SubnetId.ValueString(),
		Timezone:                  request.Timezone.ValueString(),
		MaintenanceOption:         *sqlserver.NewNullableMaintenanceOption(convertedMaintenanceOption),
		Tags:                      TagsObject,
		VipPublicIpId:             *sqlserver.NewNullableString(request.VipPublicIpId.ValueStringPointer()),
		VirtualIpAddress:          *sqlserver.NewNullableString(request.VirtualIpAddress.ValueStringPointer()),
		ServiceWatchLogCollection: *sqlserver.NewNullableBool(request.ServiceWatchLogCollection.ValueBoolPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CheckBackupConfig(initConfigOption *InitConfigOption) bool {
	return initConfigOption.BackupOption.StartingTimeHour.IsNull() && initConfigOption.BackupOption.RetentionPeriodDay.IsNull() && initConfigOption.BackupOption.ArchiveFrequencyMinute.IsNull() && initConfigOption.BackupOption.FullBackupDayOfWeek.IsNull()
}

func (client *Client) CheckMaintenanceOption(maintenanceOption *MaintenanceOption) bool {
	return !maintenanceOption.UseMaintenanceOption.ValueBool() || (maintenanceOption.StartingDayOfWeek.IsNull() && maintenanceOption.StartingTime.IsNull() && maintenanceOption.PeriodHour.IsNull())
}

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*sqlserver.SqlserverClusterDetailResponseV1Dot1, int, error) {
	req := client.sdkClient.SqlserverV1SqlserverClustersApiAPI.SqlserverShowCluster(ctx, clusterId)
	resp, httpResponse, err := req.Execute()
	if httpResponse == nil {
		return nil, 0, err
	}
	return resp, httpResponse.StatusCode, err
}

func (client *Client) DeleteCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.SqlserverV1SqlserverClustersApiAPI.SqlserverDeleteCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StopCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.SqlserverV1SqlserverClustersApiAPI.SqlserverStopCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StartCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.SqlserverV1SqlserverClustersApiAPI.SqlserverStartCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBackup(ctx context.Context, clusterId string, archiveFrequencyMinute string, startingTimeHour string, retentionPeriodDay string, fullBackupDayOfWeek string) error {
	req := client.sdkClient.SqlserverV1SqlserverBackupApiAPI.SqlserverSetBackup(ctx, clusterId)

	newFullBackupDayOfWeek, _ := sqlserver.NewDayOfWeekFromValue(fullBackupDayOfWeek)

	req = req.SqlserverBackupSettingRequest(sqlserver.SqlserverBackupSettingRequest{
		ArchiveFrequencyMinute: archiveFrequencyMinute,
		FullBackupDayOfWeek:    newFullBackupDayOfWeek,
		RetentionPeriodDay:     retentionPeriodDay,
		StartingTimeHour:       startingTimeHour,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) UnSetBackup(ctx context.Context, clusterId string) error {
	req := client.sdkClient.SqlserverV1SqlserverBackupApiAPI.SqlserverUnsetBackup(ctx, clusterId)

	// OTP 미사용이므로 session_id 는 명시적 null 로 전송한다.
	req = req.OtpSessionIdRequest(sqlserver.OtpSessionIdRequest{
		SessionId: *sqlserver.NewNullableString(nil),
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetSecurityGroupRules(ctx context.Context, clusterId string, addedIPs []string, removedIps []string) error {
	req := client.sdkClient.SqlserverV1SqlserverInstancesApiAPI.SqlserverSetSecurityGroupRules(ctx, clusterId)

	req = req.UpdateSecurityGroupRulesRequest(sqlserver.UpdateSecurityGroupRulesRequest{
		AddIpAddresses: addedIPs,
		DelIpAddresses: removedIps,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetServerType(ctx context.Context, instanceGroupId string, serverTypeName string) error {
	req := client.sdkClient.SqlserverV1SqlserverInstancesApiAPI.SqlserverSetServerType(ctx, instanceGroupId)
	reqState := &sqlserver.InstanceGroupResizeRequest{ServerTypeName: serverTypeName}
	req = req.InstanceGroupResizeRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBlockStorageSize(ctx context.Context, blockStorageGroupId string, sizeGb int32) error {
	req := client.sdkClient.SqlserverV1SqlserverInstancesApiAPI.SqlserverSetBlockStorageSize(ctx, blockStorageGroupId)
	reqState := &sqlserver.ResizeBlockStorageGroupRequest{SizeGb: sizeGb}
	req = req.ResizeBlockStorageGroupRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

// AddBlockStorages는 기존 instance group에 블록 스토리지를 추가한다.
// SDK가 SqlserverExtraBlockStorageGroupRoleType(DATA 단일 값) 타입을 제공하지만,
// roleType은 Terraform 설정에서 온 런타임 문자열이라 형변환만으로는 검증되지 않는다.
// 따라서 변환 전에 허용 값인지 직접 확인한다.
func (client *Client) AddBlockStorages(ctx context.Context, instanceGroupId string, roleType string, sizeGb int32, volumeType string) error {
	if !database.ContainsRoleType(database.BSRoleTypesExtraSqlserver, roleType) {
		return fmt.Errorf("invalid role_type %q for sqlserver block storage addition (allowed: %s)", roleType, strings.Join(database.BSRoleTypesExtraSqlserver, ", "))
	}

	req := client.sdkClient.SqlserverV1SqlserverInstancesApiAPI.SqlserverAddBlockStorages(ctx, instanceGroupId)
	reqState := &sqlserver.SqlserverAddBlockStoragesRequest{
		RoleType:   sqlserver.SqlserverExtraBlockStorageGroupRoleType(roleType),
		SizeGb:     sizeGb,
		VolumeType: sqlserver.VolumeType(volumeType).Ptr(),
	}
	req = req.SqlserverAddBlockStoragesRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func MapInstanceGroupResponses(sdkResp []sqlserver.SqlserverInstanceGroupResponse) []database.InstanceGroupResponse {
	if sdkResp == nil {
		return nil
	}

	result := make([]database.InstanceGroupResponse, len(sdkResp))
	for i, ig := range sdkResp {
		bsGroups := make([]database.BlockStorageGroupResponse, len(ig.BlockStorageGroups))
		for j, bs := range ig.BlockStorageGroups {
			bsGroups[j] = database.BlockStorageGroupResponse{
				Id:         bs.Id,
				Name:       bs.Name,
				RoleType:   string(bs.RoleType),
				SizeGb:     bs.SizeGb,
				VolumeType: string(bs.VolumeType),
			}
		}

		instances := make([]database.InstanceResponse, len(ig.Instances))
		for j, it := range ig.Instances {
			var pubIP, serviceIP string
			if it.ServiceIpAddress.Get() != nil {
				serviceIP = *it.ServiceIpAddress.Get()
			}
			if it.PublicIpId.Get() != nil {
				pubIP = *it.PublicIpId.Get()
			}

			instances[j] = database.InstanceResponse{
				Name:             it.Name,
				PublicIpId:       pubIP,
				RoleType:         string(it.RoleType),
				ServiceIpAddress: serviceIP,
			}
		}

		result[i] = database.InstanceGroupResponse{
			BlockStorageGroups: bsGroups,
			Id:                 ig.Id,
			Instances:          instances,
			RoleType:           string(ig.RoleType),
			ServerTypeName:     ig.ServerTypeName,
		}
	}

	return result
}

// instance
// 클러스터 내 특정 인스턴스의 상세 정보를 조회한다.
func (client *Client) GetInstance(ctx context.Context, clusterId string, instanceName string) (*sqlserver.InstanceDetailResponse, error) {
	req := client.sdkClient.SqlserverV1SqlserverInstancesApiAPI.SqlserverShowInstance(ctx, clusterId, instanceName)
	resp, _, err := req.Execute()
	return resp, err
}
