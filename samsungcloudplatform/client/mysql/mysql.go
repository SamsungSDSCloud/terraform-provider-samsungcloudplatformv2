package mysql

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/mysql/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *mysql.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: mysql.NewAPIClient(config),
	}
}

func (client *Client) GetClusterList(ctx context.Context, request ClusterDataSource) (*mysql.RdbClusterPageResponse, error) {
	req := client.sdkClient.MysqlV1MysqlClustersApiAPI.MysqlListClusters(ctx)
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

func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*mysql.AsyncResponse, error) {
	req := client.sdkClient.MysqlV1MysqlClustersApiAPI.MysqlCreateCluster(ctx)

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
	var convertedBackupOption = &mysql.MysqlBackupOption{}

	if client.CheckBackupConfig(initConfigOption) {
		convertedBackupOption = nil
	} else {
		convertedBackupOption = &mysql.MysqlBackupOption{
			ArchiveFrequencyMinute: *mysql.NewNullableString(initConfigOption.BackupOption.ArchiveFrequencyMinute.ValueStringPointer()),
			RetentionPeriodDay:     *mysql.NewNullableString(initConfigOption.BackupOption.RetentionPeriodDay.ValueStringPointer()),
			StartingTimeHour:       *mysql.NewNullableString(initConfigOption.BackupOption.StartingTimeHour.ValueStringPointer()),
		}
	}

	var convertedInitConfigOption = mysql.MysqlInitConfigOptionRequest{
		DatabaseName:          initConfigOption.DatabaseName.ValueString(),
		DatabaseUserName:      initConfigOption.DatabaseUserName.ValueString(),
		DatabaseUserPassword:  initConfigOption.DatabaseUserPassword.ValueString(),
		DatabasePort:          *mysql.NewNullableInt32(initConfigOption.DatabasePort.ValueInt32Pointer()),
		DatabaseCharacterSet:  *mysql.NewNullableString(initConfigOption.DatabaseCharacterSet.ValueStringPointer()),
		DatabaseCaseSensitive: initConfigOption.DatabaseCaseSensitive.ValueBoolPointer(),
		BackupOption:          *mysql.NewNullableMysqlBackupOption(convertedBackupOption),
	}

	// InstanceGroups
	var convertedInstanceGroups []mysql.InstanceGroupRequest
	for _, instanceGroup := range request.InstanceGroups {
		var convertedBlockStorage []mysql.BlockStorageGroupRequest
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			convertedBlockStorage = append(convertedBlockStorage, mysql.BlockStorageGroupRequest{
				RoleType:   mysql.BlockStorageGroupRoleType(blockStorage.RoleType.ValueString()),
				SizeGb:     blockStorage.SizeGb.ValueInt32(),
				VolumeType: mysql.VolumeType(blockStorage.VolumeType.ValueString()).Ptr(),
			})
		}

		var convertedInstance []mysql.InstanceRequest
		for _, instance := range instanceGroup.Instances {
			convertedInstance = append(convertedInstance, mysql.InstanceRequest{
				RoleType:         mysql.InstanceRoleType(instance.RoleType.ValueString()),
				ServiceIpAddress: *mysql.NewNullableString(instance.ServiceIpAddress.ValueStringPointer()),
				PublicIpId:       *mysql.NewNullableString(instance.PublicIpId.ValueStringPointer()),
			})
		}

		convertedInstanceGroups = append(convertedInstanceGroups, mysql.InstanceGroupRequest{
			BlockStorageGroups: convertedBlockStorage,
			Instances:          convertedInstance,
			RoleType:           mysql.InstanceGroupRoleType(instanceGroup.RoleType.ValueString()),
			ServerTypeName:     instanceGroup.ServerTypeName.ValueString(),
		})
	}

	// MaintenanceOption
	var maintenanceOption = request.MaintenanceOption
	var convertedMaintenanceOption = &mysql.MaintenanceOption{}

	if client.CheckMaintenanceOption(maintenanceOption) {
		convertedMaintenanceOption = nil
	} else {
		var startingDayOfWeek, _ = mysql.NewDayOfWeekFromValue(maintenanceOption.StartingDayOfWeek.ValueString())
		convertedMaintenanceOption = &mysql.MaintenanceOption{
			PeriodHour:        maintenanceOption.PeriodHour.ValueStringPointer(),
			StartingTime:      maintenanceOption.StartingTime.ValueStringPointer(),
			StartingDayOfWeek: startingDayOfWeek,
		}
	}

	//Tags
	var TagsObject []mysql.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := mysql.Tag{
			Key:   &k,
			Value: *mysql.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.MysqlClusterCreateRequest(mysql.MysqlClusterCreateRequest{
		AllowableIpAddresses: allowableIpAddresses,
		DbaasEngineVersionId: request.DbaasEngineVersionId.ValueString(),
		NatEnabled:           request.NatEnabled.ValueBoolPointer(),
		HaEnabled:            request.HaEnabled.ValueBoolPointer(),
		InitConfigOption:     convertedInitConfigOption,
		InstanceGroups:       convertedInstanceGroups,
		InstanceNamePrefix:   request.InstanceNamePrefix.ValueString(),
		Name:                 request.Name.ValueString(),
		SubnetId:             request.SubnetId.ValueString(),
		Timezone:             request.Timezone.ValueString(),
		MaintenanceOption:    *mysql.NewNullableMaintenanceOption(convertedMaintenanceOption),
		Tags:                 TagsObject,
		VipPublicIpId:        *mysql.NewNullableString(request.VipPublicIpId.ValueStringPointer()),
		VirtualIpAddress:     *mysql.NewNullableString(request.VirtualIpAddress.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CheckBackupConfig(initConfigOption InitConfigOption) bool {
	return initConfigOption.BackupOption.StartingTimeHour.IsNull() && initConfigOption.BackupOption.RetentionPeriodDay.IsNull() && initConfigOption.BackupOption.ArchiveFrequencyMinute.IsNull()
}

func (client *Client) CheckMaintenanceOption(maintenanceOption MaintenanceOption) bool {
	return !maintenanceOption.UseMaintenanceOption.ValueBool() || (maintenanceOption.StartingDayOfWeek.IsNull() && maintenanceOption.StartingTime.IsNull() && maintenanceOption.PeriodHour.IsNull())
}

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*mysql.MysqlClusterDetailResponse, error) {
	req := client.sdkClient.MysqlV1MysqlClustersApiAPI.MysqlShowCluster(ctx, clusterId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.MysqlV1MysqlClustersApiAPI.MysqlDeleteCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StopCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.MysqlV1MysqlClustersApiAPI.MysqlStopCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StartCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.MysqlV1MysqlClustersApiAPI.MysqlStartCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetServerType(ctx context.Context, instanceGroupId string, serverTypeName string) error {
	req := client.sdkClient.MysqlV1MysqlInstancesApiAPI.MysqlSetServerType(ctx, instanceGroupId)
	reqState := &mysql.InstanceGroupResizeRequest{ServerTypeName: serverTypeName}
	req = req.InstanceGroupResizeRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBlockStorageSize(ctx context.Context, blockStorageGroupId string, sizeGb int32) error {
	req := client.sdkClient.MysqlV1MysqlInstancesApiAPI.MysqlSetBlockStorageSize(ctx, blockStorageGroupId)
	reqState := &mysql.ResizeBlockStorageGroupRequest{SizeGb: sizeGb}
	req = req.ResizeBlockStorageGroupRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) AddBlockStorages(ctx context.Context, instanceGroupId string, roleType string, sizeGb int32, volumeType string) error {
	req := client.sdkClient.MysqlV1MysqlInstancesApiAPI.MysqlAddBlockStorages(ctx, instanceGroupId)
	reqState := &mysql.AddBlockStoragesRequest{
		RoleType:   mysql.BlockStorageGroupRoleType(roleType),
		SizeGb:     sizeGb,
		VolumeType: mysql.VolumeType(volumeType).Ptr(),
	}
	req = req.AddBlockStoragesRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBackup(ctx context.Context, clusterId string, archiveFrequencyMinute string, startingTimeHour string, retentionPeriodDay string) error {
	req := client.sdkClient.MysqlV1MysqlBackupApiAPI.MysqlSetBackup(ctx, clusterId)

	req = req.BackupSettingRequest(mysql.BackupSettingRequest{
		ArchiveFrequencyMinute: archiveFrequencyMinute,
		RetentionPeriodDay:     retentionPeriodDay,
		StartingTimeHour:       startingTimeHour,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) UnSetBackup(ctx context.Context, clusterId string) error {
	req := client.sdkClient.MysqlV1MysqlBackupApiAPI.MysqlUnsetBackup(ctx, clusterId)

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetSecurityGroupRules(ctx context.Context, clusterId string, addedIPs []string, removedIps []string) error {
	req := client.sdkClient.MysqlV1MysqlInstancesApiAPI.MysqlSetSecurityGroupRules(ctx, clusterId)

	req = req.UpdateSecurityGroupRulesRequest(mysql.UpdateSecurityGroupRulesRequest{
		AddIpAddresses: addedIPs,
		DelIpAddresses: removedIps,
	})

	_, _, err := req.Execute()
	return err
}
