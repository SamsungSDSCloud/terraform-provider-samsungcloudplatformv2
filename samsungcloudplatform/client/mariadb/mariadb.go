package mariadb

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/mariadb/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *mariadb.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: mariadb.NewAPIClient(config),
	}
}

func (client *Client) GetClusterList(ctx context.Context, request ClusterDataSource) (*mariadb.RdbClusterPageResponse, error) {
	req := client.sdkClient.MariadbV1MariadbClustersApiAPI.MariadbListClusters(ctx)
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

func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*mariadb.AsyncResponse, error) {
	req := client.sdkClient.MariadbV1MariadbClustersApiAPI.MariadbCreateCluster(ctx)

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
	var convertedBackupOption = &mariadb.MariadbBackupOption{}

	if client.CheckBackupConfig(initConfigOption) {
		convertedBackupOption = nil
	} else {
		convertedBackupOption = &mariadb.MariadbBackupOption{
			ArchiveFrequencyMinute: *mariadb.NewNullableString(initConfigOption.BackupOption.ArchiveFrequencyMinute.ValueStringPointer()),
			RetentionPeriodDay:     *mariadb.NewNullableString(initConfigOption.BackupOption.RetentionPeriodDay.ValueStringPointer()),
			StartingTimeHour:       *mariadb.NewNullableString(initConfigOption.BackupOption.StartingTimeHour.ValueStringPointer()),
		}
	}

	var convertedInitConfigOption = mariadb.MariadbInitConfigOptionRequest{
		AuditEnabled:         initConfigOption.AuditEnabled.ValueBoolPointer(),
		DatabaseName:         initConfigOption.DatabaseName.ValueString(),
		DatabaseUserName:     initConfigOption.DatabaseUserName.ValueString(),
		DatabaseUserPassword: initConfigOption.DatabaseUserPassword.ValueString(),
		DatabasePort:         *mariadb.NewNullableInt32(initConfigOption.DatabasePort.ValueInt32Pointer()),
		DatabaseCharacterSet: *mariadb.NewNullableString(initConfigOption.DatabaseCharacterSet.ValueStringPointer()),
		BackupOption:         *mariadb.NewNullableMariadbBackupOption(convertedBackupOption),
	}

	// InstanceGroups
	var convertedInstanceGroups []mariadb.InstanceGroupRequest
	for _, instanceGroup := range request.InstanceGroups {
		var convertedBlockStorage []mariadb.BlockStorageGroupRequest
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			convertedBlockStorage = append(convertedBlockStorage, mariadb.BlockStorageGroupRequest{
				RoleType:   mariadb.BlockStorageGroupRoleType(blockStorage.RoleType.ValueString()),
				SizeGb:     blockStorage.SizeGb.ValueInt32(),
				VolumeType: mariadb.VolumeType(blockStorage.VolumeType.ValueString()).Ptr(),
			})
		}

		var convertedInstance []mariadb.InstanceRequest
		for _, instance := range instanceGroup.Instances {
			convertedInstance = append(convertedInstance, mariadb.InstanceRequest{
				RoleType:         mariadb.InstanceRoleType(instance.RoleType.ValueString()),
				ServiceIpAddress: *mariadb.NewNullableString(instance.ServiceIpAddress.ValueStringPointer()),
				PublicIpId:       *mariadb.NewNullableString(instance.PublicIpId.ValueStringPointer()),
			})
		}

		convertedInstanceGroups = append(convertedInstanceGroups, mariadb.InstanceGroupRequest{
			BlockStorageGroups: convertedBlockStorage,
			Instances:          convertedInstance,
			RoleType:           mariadb.InstanceGroupRoleType(instanceGroup.RoleType.ValueString()),
			ServerTypeName:     instanceGroup.ServerTypeName.ValueString(),
		})
	}

	// MaintenanceOption
	var maintenanceOption = request.MaintenanceOption
	var convertedMaintenanceOption = &mariadb.MaintenanceOption{}

	if client.CheckMaintenanceOption(maintenanceOption) {
		convertedMaintenanceOption = nil
	} else {
		var startingDayOfWeek, _ = mariadb.NewDayOfWeekFromValue(maintenanceOption.StartingDayOfWeek.ValueString())
		convertedMaintenanceOption = &mariadb.MaintenanceOption{
			PeriodHour:        maintenanceOption.PeriodHour.ValueStringPointer(),
			StartingTime:      maintenanceOption.StartingTime.ValueStringPointer(),
			StartingDayOfWeek: startingDayOfWeek,
		}
	}

	//Tags
	var TagsObject []mariadb.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := mariadb.Tag{
			Key:   &k,
			Value: *mariadb.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.MariadbClusterCreateRequest(mariadb.MariadbClusterCreateRequest{
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
		MaintenanceOption:    *mariadb.NewNullableMaintenanceOption(convertedMaintenanceOption),
		Tags:                 TagsObject,
		VipPublicIpId:        *mariadb.NewNullableString(request.VipPublicIpId.ValueStringPointer()),
		VirtualIpAddress:     *mariadb.NewNullableString(request.VirtualIpAddress.ValueStringPointer()),
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

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*mariadb.MariadbClusterDetailResponse, error) {
	req := client.sdkClient.MariadbV1MariadbClustersApiAPI.MariadbShowCluster(ctx, clusterId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.MariadbV1MariadbClustersApiAPI.MariadbDeleteCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StopCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.MariadbV1MariadbClustersApiAPI.MariadbStopCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StartCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.MariadbV1MariadbClustersApiAPI.MariadbStartCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetServerType(ctx context.Context, instanceGroupId string, serverTypeName string) error {
	req := client.sdkClient.MariadbV1MariadbInstancesApiAPI.MariadbSetServerType(ctx, instanceGroupId)
	reqState := &mariadb.InstanceGroupResizeRequest{ServerTypeName: serverTypeName}
	req = req.InstanceGroupResizeRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBlockStorageSize(ctx context.Context, blockStorageGroupId string, sizeGb int32) error {
	req := client.sdkClient.MariadbV1MariadbInstancesApiAPI.MariadbSetBlockStorageSize(ctx, blockStorageGroupId)
	reqState := &mariadb.ResizeBlockStorageGroupRequest{SizeGb: sizeGb}
	req = req.ResizeBlockStorageGroupRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) AddBlockStorages(ctx context.Context, instanceGroupId string, roleType string, sizeGb int32, volumeType string) error {
	req := client.sdkClient.MariadbV1MariadbInstancesApiAPI.MariadbAddBlockStorages(ctx, instanceGroupId)
	reqState := &mariadb.AddBlockStoragesRequest{
		RoleType:   mariadb.BlockStorageGroupRoleType(roleType),
		SizeGb:     sizeGb,
		VolumeType: mariadb.VolumeType(volumeType).Ptr(),
	}
	req = req.AddBlockStoragesRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBackup(ctx context.Context, clusterId string, archiveFrequencyMinute string, startingTimeHour string, retentionPeriodDay string) error {
	req := client.sdkClient.MariadbV1MariadbBackupApiAPI.MariadbSetBackup(ctx, clusterId)

	req = req.BackupSettingRequest(mariadb.BackupSettingRequest{
		ArchiveFrequencyMinute: archiveFrequencyMinute,
		RetentionPeriodDay:     retentionPeriodDay,
		StartingTimeHour:       startingTimeHour,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) UnSetBackup(ctx context.Context, clusterId string) error {
	req := client.sdkClient.MariadbV1MariadbBackupApiAPI.MariadbUnsetBackup(ctx, clusterId)

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetSecurityGroupRules(ctx context.Context, clusterId string, addedIPs []string, removedIps []string) error {
	req := client.sdkClient.MariadbV1MariadbInstancesApiAPI.MariadbSetSecurityGroupRules(ctx, clusterId)

	req = req.UpdateSecurityGroupRulesRequest(mariadb.UpdateSecurityGroupRulesRequest{
		AddIpAddresses: addedIPs,
		DelIpAddresses: removedIps,
	})

	_, _, err := req.Execute()
	return err
}
