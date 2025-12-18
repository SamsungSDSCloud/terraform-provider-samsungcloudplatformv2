package searchengine

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/searchengine/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *searchengine.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: searchengine.NewAPIClient(config),
	}
}

// clusterlist (ctx, ClusterDataSource) - (RdbClusterPageResponse)
func (client *Client) GetClusterList(ctx context.Context, request ClusterDataSource) (*searchengine.ClusterPageResponse, error) {
	req := client.sdkClient.SearchengineV1SearchEngineClustersApiAPI.SearchengineListClusters(ctx)
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

// engine version
func (client *Client) GetEngineVersionList(ctx context.Context) (*searchengine.EngineListResponse, error) {
	req := client.sdkClient.SearchengineV1SearchEngineMasterDataApiAPI.SearchengineListEngineVersions(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

// create (ctx, clusterResource) - (asyncResponse)
func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*searchengine.AsyncResponse, error) {
	req := client.sdkClient.SearchengineV1SearchEngineClustersApiAPI.SearchengineCreateCluster(ctx)

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

	//for _, allowableIpAddress := range request.AllowableIpAddresses {
	//	allowableIpAddresses = append(allowableIpAddresses, allowableIpAddress.ValueString())
	//}

	// InitConfigOption
	var initConfigOption = request.InitConfigOption
	var convertedBackupOption = &searchengine.BackupSettingExcludingArchiveRequest{}

	if client.CheckBackupConfig(initConfigOption) {
		convertedBackupOption = nil
	} else {
		convertedBackupOption = &searchengine.BackupSettingExcludingArchiveRequest{
			RetentionPeriodDay: initConfigOption.BackupOption.RetentionPeriodDay.ValueString(),
			StartingTimeHour:   initConfigOption.BackupOption.StartingTimeHour.ValueString(),
		}
	}

	var convertedInitConfigOption = searchengine.SearchEngineInitConfigOptionRequest{}

	if initConfigOption.DatabasePort.IsUnknown() {
		convertedInitConfigOption = searchengine.SearchEngineInitConfigOptionRequest{
			BackupOption:         *searchengine.NewNullableBackupSettingExcludingArchiveRequest(convertedBackupOption),
			DatabaseUserName:     initConfigOption.DatabaseUserName.ValueString(),
			DatabaseUserPassword: initConfigOption.DatabaseUserPassword.ValueString(),
		}
	} else {
		convertedInitConfigOption = searchengine.SearchEngineInitConfigOptionRequest{
			BackupOption:         *searchengine.NewNullableBackupSettingExcludingArchiveRequest(convertedBackupOption),
			DatabasePort:         *searchengine.NewNullableInt32(initConfigOption.DatabasePort.ValueInt32Pointer()),
			DatabaseUserName:     initConfigOption.DatabaseUserName.ValueString(),
			DatabaseUserPassword: initConfigOption.DatabaseUserPassword.ValueString(),
		}
	}

	// InstanceGroups
	var convertedInstanceGroups []searchengine.InstanceGroupRequest
	for _, instanceGroup := range request.InstanceGroups {
		var convertedBlockStorage []searchengine.BlockStorageGroupRequest
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			convertedBlockStorage = append(convertedBlockStorage, searchengine.BlockStorageGroupRequest{
				RoleType:   searchengine.BlockStorageGroupRoleType(blockStorage.RoleType.ValueString()),
				SizeGb:     blockStorage.SizeGb.ValueInt32(),
				VolumeType: searchengine.VolumeType(blockStorage.VolumeType.ValueString()).Ptr(),
			})
		}

		var convertedInstance []searchengine.InstanceRequest
		for _, instance := range instanceGroup.Instances {
			convertedInstance = append(convertedInstance, searchengine.InstanceRequest{
				RoleType:         searchengine.InstanceRoleType(instance.RoleType.ValueString()),
				ServiceIpAddress: *searchengine.NewNullableString(instance.ServiceIpAddress.ValueStringPointer()),
				PublicIpId:       *searchengine.NewNullableString(instance.PublicIpId.ValueStringPointer()),
			})
		}

		convertedInstanceGroups = append(convertedInstanceGroups, searchengine.InstanceGroupRequest{
			BlockStorageGroups: convertedBlockStorage,
			Instances:          convertedInstance,
			RoleType:           searchengine.InstanceGroupRoleType(instanceGroup.RoleType.ValueString()),
			ServerTypeName:     instanceGroup.ServerTypeName.ValueString(),
		})
	}

	// MaintenanceOption
	var maintenanceOption = request.MaintenanceOption
	var convertedMaintenanceOption = &searchengine.MaintenanceOption{}

	if client.CheckMaintenanceOption(maintenanceOption) {
		convertedMaintenanceOption = nil
	} else {
		var startingDayOfWeek, _ = searchengine.NewDayOfWeekFromValue(maintenanceOption.StartingDayOfWeek.ValueString())
		convertedMaintenanceOption = &searchengine.MaintenanceOption{
			PeriodHour:        maintenanceOption.PeriodHour.ValueStringPointer(),
			StartingTime:      maintenanceOption.StartingTime.ValueStringPointer(),
			StartingDayOfWeek: startingDayOfWeek,
		}
	}

	//Tags
	var TagsObject []searchengine.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := searchengine.Tag{
			Key:   &k,
			Value: *searchengine.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.SearchEngineClusterCreateRequest(searchengine.SearchEngineClusterCreateRequest{
		AllowableIpAddresses: allowableIpAddresses,
		DbaasEngineVersionId: request.DbaasEngineVersionId.ValueString(),
		InitConfigOption:     convertedInitConfigOption,
		InstanceGroups:       convertedInstanceGroups,
		InstanceNamePrefix:   request.InstanceNamePrefix.ValueString(),
		IsCombined:           request.IsCombined.ValueBoolPointer(),
		NatEnabled:           request.NatEnabled.ValueBoolPointer(),
		Name:                 request.Name.ValueString(),
		SubnetId:             request.SubnetId.ValueString(),
		Timezone:             request.Timezone.ValueString(),
		MaintenanceOption:    *searchengine.NewNullableMaintenanceOption(convertedMaintenanceOption),
		Tags:                 TagsObject,
		License:              *searchengine.NewNullableString(request.License.ValueStringPointer()),
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

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*searchengine.SearchEngineClusterDetailResponse, error) {
	req := client.sdkClient.SearchengineV1SearchEngineClustersApiAPI.SearchengineShowCluster(ctx, clusterId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.SearchengineV1SearchEngineClustersApiAPI.SearchengineDeleteCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StopCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.SearchengineV1SearchEngineClustersApiAPI.SearchengineStopCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StartCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.SearchengineV1SearchEngineClustersApiAPI.SearchengineStartCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBackup(ctx context.Context, clusterId string, startingTimeHour string, retentionPeriodDay string) error {
	req := client.sdkClient.SearchengineV1SearchEngineBackupApiAPI.SearchengineSetBackup(ctx, clusterId)

	req = req.BackupSettingExcludingArchiveRequest(searchengine.BackupSettingExcludingArchiveRequest{
		RetentionPeriodDay: retentionPeriodDay,
		StartingTimeHour:   startingTimeHour,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) UnSetBackup(ctx context.Context, clusterId string) error {
	req := client.sdkClient.SearchengineV1SearchEngineBackupApiAPI.SearchengineUnsetBackup(ctx, clusterId)

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetSecurityGroupRules(ctx context.Context, clusterId string, addedIPs []string, removedIps []string) error {
	req := client.sdkClient.SearchengineV1SearchEngineInstancesApiAPI.SearchengineSetSecurityGroupRules(ctx, clusterId)

	req = req.UpdateSecurityGroupRulesRequest(searchengine.UpdateSecurityGroupRulesRequest{
		AddIpAddresses: addedIPs,
		DelIpAddresses: removedIps,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetServerType(ctx context.Context, instanceGroupId string, serverTypeName string) error {
	req := client.sdkClient.SearchengineV1SearchEngineInstancesApiAPI.SearchengineSetServerType(ctx, instanceGroupId)
	reqState := &searchengine.InstanceGroupResizeRequest{ServerTypeName: serverTypeName}
	req = req.InstanceGroupResizeRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBlockStorageSize(ctx context.Context, blockStorageGroupId string, sizeGb int32) error {
	req := client.sdkClient.SearchengineV1SearchEngineInstancesApiAPI.SearchengineSetBlockStorageSize(ctx, blockStorageGroupId)
	reqState := &searchengine.ResizeBlockStorageGroupRequest{SizeGb: sizeGb}
	req = req.ResizeBlockStorageGroupRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) AddBlockStorages(ctx context.Context, instanceGroupId string, roleType string, sizeGb int32, volumeType string) error {
	req := client.sdkClient.SearchengineV1SearchEngineInstancesApiAPI.SearchengineAddBlockStorages(ctx, instanceGroupId)
	reqState := &searchengine.AddBlockStoragesRequest{
		RoleType:   searchengine.BlockStorageGroupRoleType(roleType),
		SizeGb:     sizeGb,
		VolumeType: searchengine.VolumeType(volumeType).Ptr(),
	}
	req = req.AddBlockStoragesRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) AddInstances(ctx context.Context, clusterId string, instanceCount int32, serviceIPAddresses []string) error {
	req := client.sdkClient.SearchengineV1SearchEngineClustersApiAPI.SearchengineAddInstances(ctx, clusterId)
	reqState := &searchengine.SearchEngineClusterAddInstancesRequest{
		InstanceCount:      instanceCount,
		ServiceIpAddresses: serviceIPAddresses,
	}
	req = req.SearchEngineClusterAddInstancesRequest(*reqState)
	_, _, err := req.Execute()
	return err
}
