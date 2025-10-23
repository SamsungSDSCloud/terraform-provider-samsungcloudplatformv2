package vertica

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/vertica/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *vertica.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: vertica.NewAPIClient(config),
	}
}

// clusterlist (ctx, ClusterDataSource) - (ClusterPageResponse)
func (client *Client) GetClusterList(ctx context.Context, request ClusterDataSource) (*vertica.VerticaClusterPageResponse, error) {
	req := client.sdkClient.VerticaV1VerticaClustersApiAPI.VerticaListClusters(ctx)
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

	resp, _, err := req.Execute()
	return resp, err
}

// create (ctx, clusterResource) - (asyncResponse)
func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*vertica.AsyncResponse, error) {
	req := client.sdkClient.VerticaV1VerticaClustersApiAPI.VerticaCreateCluster(ctx)

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
	var convertedBackupOption = &vertica.BackupSettingExcludingArchiveRequest{}

	if client.CheckBackupConfig(initConfigOption) {
		convertedBackupOption = nil
	} else {
		convertedBackupOption = &vertica.BackupSettingExcludingArchiveRequest{
			RetentionPeriodDay: initConfigOption.BackupOption.RetentionPeriodDay.ValueString(),
			StartingTimeHour:   initConfigOption.BackupOption.StartingTimeHour.ValueString(),
		}
	}

	var convertedInitConfigOption = vertica.VerticaInitConfigOptionRequest{
		BackupOption:         *vertica.NewNullableBackupSettingExcludingArchiveRequest(convertedBackupOption),
		DatabaseLocale:       *vertica.NewNullableString(initConfigOption.DatabaseLocale.ValueStringPointer()),
		DatabaseName:         initConfigOption.DatabaseName.ValueString(),
		DatabaseUserName:     initConfigOption.DatabaseUserName.ValueString(),
		DatabaseUserPassword: initConfigOption.DatabaseUserPassword.ValueString(),
	}

	// InstanceGroups
	var convertedInstanceGroups []vertica.InstanceGroupRequest
	for _, instanceGroup := range request.InstanceGroups {
		var convertedBlockStorage []vertica.BlockStorageGroupRequest
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			convertedBlockStorage = append(convertedBlockStorage, vertica.BlockStorageGroupRequest{
				RoleType:   vertica.BlockStorageGroupRoleType(blockStorage.RoleType.ValueString()),
				SizeGb:     blockStorage.SizeGb.ValueInt32(),
				VolumeType: vertica.VolumeType(blockStorage.VolumeType.ValueString()).Ptr(),
			})
		}

		var convertedInstance []vertica.InstanceRequest
		for _, instance := range instanceGroup.Instances {
			convertedInstance = append(convertedInstance, vertica.InstanceRequest{
				RoleType:         vertica.InstanceRoleType(instance.RoleType.ValueString()),
				ServiceIpAddress: *vertica.NewNullableString(instance.ServiceIpAddress.ValueStringPointer()),
				PublicIpId:       *vertica.NewNullableString(instance.PublicIpId.ValueStringPointer()),
			})
		}

		convertedInstanceGroups = append(convertedInstanceGroups, vertica.InstanceGroupRequest{
			BlockStorageGroups: convertedBlockStorage,
			Instances:          convertedInstance,
			RoleType:           vertica.InstanceGroupRoleType(instanceGroup.RoleType.ValueString()),
			ServerTypeName:     instanceGroup.ServerTypeName.ValueString(),
		})
	}

	// MaintenanceOption
	var maintenanceOption = request.MaintenanceOption
	var convertedMaintenanceOption = &vertica.MaintenanceOption{}

	if client.CheckMaintenanceOption(maintenanceOption) {
		convertedMaintenanceOption = nil
	} else {
		var startingDayOfWeek, _ = vertica.NewDayOfWeekFromValue(maintenanceOption.StartingDayOfWeek.ValueString())
		convertedMaintenanceOption = &vertica.MaintenanceOption{
			PeriodHour:        maintenanceOption.PeriodHour.ValueStringPointer(),
			StartingTime:      maintenanceOption.StartingTime.ValueStringPointer(),
			StartingDayOfWeek: startingDayOfWeek,
		}
	}

	//Tags
	var TagsObject []vertica.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := vertica.Tag{
			Key:   &k,
			Value: *vertica.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.VerticaClusterCreateRequest(vertica.VerticaClusterCreateRequest{
		AllowableIpAddresses: allowableIpAddresses,
		DbaasEngineVersionId: request.DbaasEngineVersionId.ValueString(),
		InitConfigOption:     convertedInitConfigOption,
		InstanceGroups:       convertedInstanceGroups,
		InstanceNamePrefix:   request.InstanceNamePrefix.ValueString(),
		Name:                 request.Name.ValueString(),
		NatEnabled:           *vertica.NewNullableBool(request.NatEnabled.ValueBoolPointer()),
		SubnetId:             request.SubnetId.ValueString(),
		Timezone:             request.Timezone.ValueString(),
		MaintenanceOption:    *vertica.NewNullableMaintenanceOption(convertedMaintenanceOption),
		Tags:                 TagsObject,
		License:              request.License.ValueStringPointer(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CheckBackupConfig(initConfigOption InitConfigOption) bool {
	return initConfigOption.BackupOption.StartingTimeHour.IsNull() && initConfigOption.BackupOption.RetentionPeriodDay.IsNull()

}

func (client *Client) CheckMaintenanceOption(maintenanceOption MaintenanceOption) bool {
	return !maintenanceOption.UseMaintenanceOption.ValueBool() || (maintenanceOption.StartingDayOfWeek.IsNull() && maintenanceOption.StartingTime.IsNull() && maintenanceOption.PeriodHour.IsNull())
}

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*vertica.VerticaClusterDetailResponse, error) {
	req := client.sdkClient.VerticaV1VerticaClustersApiAPI.VerticaShowCluster(ctx, clusterId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.VerticaV1VerticaClustersApiAPI.VerticaDeleteCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StopCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.VerticaV1VerticaClustersApiAPI.VerticaStopCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StartCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.VerticaV1VerticaClustersApiAPI.VerticaStartCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBackup(ctx context.Context, clusterId string, startingTimeHour string, retentionPeriodDay string) error {
	req := client.sdkClient.VerticaV1VerticaBackupApiAPI.VerticaSetBackup(ctx, clusterId)

	req = req.BackupSettingExcludingArchiveRequest(vertica.BackupSettingExcludingArchiveRequest{
		RetentionPeriodDay: retentionPeriodDay,
		StartingTimeHour:   startingTimeHour,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) UnSetBackup(ctx context.Context, clusterId string) error {
	req := client.sdkClient.VerticaV1VerticaBackupApiAPI.VerticaUnsetBackup(ctx, clusterId)

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetSecurityGroupRules(ctx context.Context, clusterId string, addedIPs []string, removedIps []string) error {
	req := client.sdkClient.VerticaV1VerticaInstancesApiAPI.VerticaSetSecurityGroupRules(ctx, clusterId)

	req = req.UpdateSecurityGroupRulesRequest(vertica.UpdateSecurityGroupRulesRequest{
		AddIpAddresses: addedIPs,
		DelIpAddresses: removedIps,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetServerType(ctx context.Context, instanceGroupId string, serverTypeName string) error {
	req := client.sdkClient.VerticaV1VerticaInstancesApiAPI.VerticaSetServerType(ctx, instanceGroupId)
	reqState := &vertica.InstanceGroupResizeRequest{ServerTypeName: serverTypeName}
	req = req.InstanceGroupResizeRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBlockStorageSize(ctx context.Context, blockStorageGroupId string, sizeGb int32) error {
	req := client.sdkClient.VerticaV1VerticaInstancesApiAPI.VerticaSetBlockStorageSize(ctx, blockStorageGroupId)
	reqState := &vertica.ResizeBlockStorageGroupRequest{SizeGb: sizeGb}
	req = req.ResizeBlockStorageGroupRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) AddBlockStorages(ctx context.Context, instanceGroupId string, roleType string, sizeGb int32, volumeType string) error {
	req := client.sdkClient.VerticaV1VerticaInstancesApiAPI.VerticaAddBlockStorages(ctx, instanceGroupId)
	reqState := &vertica.AddBlockStoragesRequest{
		RoleType:   vertica.BlockStorageGroupRoleType(roleType),
		SizeGb:     sizeGb,
		VolumeType: vertica.VolumeType(volumeType).Ptr(),
	}
	req = req.AddBlockStoragesRequest(*reqState)
	_, _, err := req.Execute()
	return err
}
