package postgresql

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/postgresql/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *postgresql.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: postgresql.NewAPIClient(config),
	}
}

// clusterlist (ctx, ClusterDataSource) - (RdbClusterPageResponse)
func (client *Client) GetClusterList(ctx context.Context, request ClusterDataSource) (*postgresql.RdbClusterPageResponse, error) {
	req := client.sdkClient.PostgresqlV1PostgresqlClustersApiAPI.PostgresqlListClusters(ctx)
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

// create (ctx, clusterResource) - (asyncResponse)
func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*postgresql.AsyncResponse, error) {
	req := client.sdkClient.PostgresqlV1PostgresqlClustersApiAPI.PostgresqlCreateCluster(ctx)

	// AllowableIpAddresses
	var allowableIpAddresses []string

	if request.AllowableIpAddresses.IsNull() || request.AllowableIpAddresses.IsUnknown(){
		allowableIpAddresses = []string{}
	} else {
		for _, elem := range request.AllowableIpAddresses.Elements() {
			strVal := elem.(types.String)
			allowableIpAddresses = append(allowableIpAddresses, strVal.ValueString())
		}
	}

	// InitConfigOption
	var initConfigOption = request.InitConfigOption
	var convertedBackupOption = &postgresql.PostgresqlBackupOption{}

	if client.CheckBackupConfig(initConfigOption) {
		convertedBackupOption = nil
	} else {
		convertedBackupOption = &postgresql.PostgresqlBackupOption{
			ArchiveFrequencyMinute: *postgresql.NewNullableString(initConfigOption.BackupOption.ArchiveFrequencyMinute.ValueStringPointer()),
			RetentionPeriodDay:     *postgresql.NewNullableString(initConfigOption.BackupOption.RetentionPeriodDay.ValueStringPointer()),
			StartingTimeHour:       *postgresql.NewNullableString(initConfigOption.BackupOption.StartingTimeHour.ValueStringPointer()),
		}
	}

	var convertedInitConfigOption = postgresql.PostgresqlInitConfigOptionRequest{
		AuditEnabled:         initConfigOption.AuditEnabled.ValueBoolPointer(),
		BackupOption:         *postgresql.NewNullablePostgresqlBackupOption(convertedBackupOption),
		DatabaseEncoding:     *postgresql.NewNullableString(initConfigOption.DatabaseEncoding.ValueStringPointer()),
		DatabaseLocale:       *postgresql.NewNullableString(initConfigOption.DatabaseLocale.ValueStringPointer()),
		DatabaseName:         initConfigOption.DatabaseName.ValueString(),
		DatabasePort:         *postgresql.NewNullableInt32(initConfigOption.DatabasePort.ValueInt32Pointer()),
		DatabaseUserName:     initConfigOption.DatabaseUserName.ValueString(),
		DatabaseUserPassword: initConfigOption.DatabaseUserPassword.ValueString(),
	}

	// InstanceGroups
	var convertedInstanceGroups []postgresql.InstanceGroupRequest
	for _, instanceGroup := range request.InstanceGroups {
		var convertedBlockStorage []postgresql.BlockStorageGroupRequest
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			convertedBlockStorage = append(convertedBlockStorage, postgresql.BlockStorageGroupRequest{
				RoleType:   postgresql.BlockStorageGroupRoleType(blockStorage.RoleType.ValueString()),
				SizeGb:     blockStorage.SizeGb.ValueInt32(),
				VolumeType: postgresql.VolumeType(blockStorage.VolumeType.ValueString()).Ptr(),
			})
		}

		var convertedInstance []postgresql.InstanceRequest
		for _, instance := range instanceGroup.Instances {
			convertedInstance = append(convertedInstance, postgresql.InstanceRequest{
				RoleType:         postgresql.InstanceRoleType(instance.RoleType.ValueString()),
				ServiceIpAddress: *postgresql.NewNullableString(instance.ServiceIpAddress.ValueStringPointer()),
				PublicIpId:       *postgresql.NewNullableString(instance.PublicIpId.ValueStringPointer()),
			})
		}

		convertedInstanceGroups = append(convertedInstanceGroups, postgresql.InstanceGroupRequest{
			BlockStorageGroups: convertedBlockStorage,
			Instances:          convertedInstance,
			RoleType:           postgresql.InstanceGroupRoleType(instanceGroup.RoleType.ValueString()),
			ServerTypeName:     instanceGroup.ServerTypeName.ValueString(),
		})
	}

	// MaintenanceOption
	var maintenanceOption = request.MaintenanceOption
	var convertedMaintenanceOption = &postgresql.MaintenanceOption{}

	if client.CheckMaintenanceOption(maintenanceOption) {
		convertedMaintenanceOption = nil
	} else {
		var startingDayOfWeek, _ = postgresql.NewDayOfWeekFromValue(maintenanceOption.StartingDayOfWeek.ValueString())
		convertedMaintenanceOption = &postgresql.MaintenanceOption{
			PeriodHour:        maintenanceOption.PeriodHour.ValueStringPointer(),
			StartingTime:      maintenanceOption.StartingTime.ValueStringPointer(),
			StartingDayOfWeek: startingDayOfWeek,
		}
	}

	//Tags
	var TagsObject []postgresql.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := postgresql.Tag{
			Key:   &k,
			Value: *postgresql.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.PostgresqlClusterCreateRequest(postgresql.PostgresqlClusterCreateRequest{
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
		MaintenanceOption:    *postgresql.NewNullableMaintenanceOption(convertedMaintenanceOption),
		Tags:                 TagsObject,
		VipPublicIpId:        *postgresql.NewNullableString(request.VipPublicIpId.ValueStringPointer()),
		VirtualIpAddress:     *postgresql.NewNullableString(request.VirtualIpAddress.ValueStringPointer()),
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

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*postgresql.PostgresqlClusterDetailResponse, error) {
	req := client.sdkClient.PostgresqlV1PostgresqlClustersApiAPI.PostgresqlShowCluster(ctx, clusterId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.PostgresqlV1PostgresqlClustersApiAPI.PostgresqlDeleteCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StopCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.PostgresqlV1PostgresqlClustersApiAPI.PostgresqlStopCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StartCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.PostgresqlV1PostgresqlClustersApiAPI.PostgresqlStartCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBackup(ctx context.Context, clusterId string, archiveFrequencyMinute string, startingTimeHour string, retentionPeriodDay string) error {
	req := client.sdkClient.PostgresqlV1PostgresqlBackupApiAPI.PostgresqlSetBackup(ctx, clusterId)

	req = req.BackupSettingRequest(postgresql.BackupSettingRequest{
		ArchiveFrequencyMinute: archiveFrequencyMinute,
		RetentionPeriodDay:     retentionPeriodDay,
		StartingTimeHour:       startingTimeHour,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) UnSetBackup(ctx context.Context, clusterId string) error {
	req := client.sdkClient.PostgresqlV1PostgresqlBackupApiAPI.PostgresqlUnsetBackup(ctx, clusterId)

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetSecurityGroupRules(ctx context.Context, clusterId string, addedIPs []string, removedIps []string) error {
	req := client.sdkClient.PostgresqlV1PostgresqlInstancesApiAPI.PostgresqlSetSecurityGroupRules(ctx, clusterId)

	req = req.UpdateSecurityGroupRulesRequest(postgresql.UpdateSecurityGroupRulesRequest{
		AddIpAddresses: addedIPs,
		DelIpAddresses: removedIps,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetServerType(ctx context.Context, instanceGroupId string, serverTypeName string) error {
	req := client.sdkClient.PostgresqlV1PostgresqlInstancesApiAPI.PostgresqlSetServerType(ctx, instanceGroupId)
	reqState := &postgresql.InstanceGroupResizeRequest{ServerTypeName: serverTypeName}
	req = req.InstanceGroupResizeRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBlockStorageSize(ctx context.Context, blockStorageGroupId string, sizeGb int32) error {
	req := client.sdkClient.PostgresqlV1PostgresqlInstancesApiAPI.PostgresqlSetBlockStorageSize(ctx, blockStorageGroupId)
	reqState := &postgresql.ResizeBlockStorageGroupRequest{SizeGb: sizeGb}
	req = req.ResizeBlockStorageGroupRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) AddBlockStorages(ctx context.Context, instanceGroupId string, roleType string, sizeGb int32, volumeType string) error {
	req := client.sdkClient.PostgresqlV1PostgresqlInstancesApiAPI.PostgresqlAddBlockStorages(ctx, instanceGroupId)
	reqState := &postgresql.AddBlockStoragesRequest{
		RoleType:   postgresql.BlockStorageGroupRoleType(roleType),
		SizeGb:     sizeGb,
		VolumeType: postgresql.VolumeType(volumeType).Ptr(),
	}
	req = req.AddBlockStoragesRequest(*reqState)
	_, _, err := req.Execute()
	return err
}
