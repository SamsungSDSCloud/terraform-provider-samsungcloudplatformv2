package cachestore

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/cachestore/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *cachestore.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: cachestore.NewAPIClient(config),
	}
}

// clusterlist (ctx, ClusterDataSource) - (ClusterPageResponse)
func (client *Client) GetClusterList(ctx context.Context, request ClusterDataSource) (*cachestore.ClusterPageResponse, error) {
	req := client.sdkClient.CachestoreV1CacheStoreClustersApiAPI.CachestoreListClusters(ctx)
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
func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*cachestore.AsyncResponse, error) {
	req := client.sdkClient.CachestoreV1CacheStoreClustersApiAPI.CachestoreCreateCluster(ctx)

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
	var convertedBackupOption = &cachestore.BackupOption{}

	if client.CheckBackupConfig(initConfigOption) {
		convertedBackupOption = nil
	} else {
		convertedBackupOption = &cachestore.BackupOption{
			RetentionPeriodDay: *cachestore.NewNullableString(initConfigOption.BackupOption.RetentionPeriodDay.ValueStringPointer()),
			StartingTimeHour:   *cachestore.NewNullableString(initConfigOption.BackupOption.StartingTimeHour.ValueStringPointer()),
		}
	}

	var convertedInitConfigOption = cachestore.RedisInitConfigOption{
		BackupOption:         *cachestore.NewNullableBackupOption(convertedBackupOption),
		DatabasePort:         *cachestore.NewNullableInt32(initConfigOption.DatabasePort.ValueInt32Pointer()),
		DatabaseUserPassword: *cachestore.NewNullableString(initConfigOption.DatabaseUserPassword.ValueStringPointer()),
	}

	// InstanceGroups
	var convertedInstanceGroups []cachestore.RedisInstanceGroupRequest
	for _, instanceGroup := range request.InstanceGroups {
		var convertedBlockStorage []cachestore.RedisBlockStorageGroupRequest
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			convertedBlockStorage = append(convertedBlockStorage, cachestore.RedisBlockStorageGroupRequest{
				RoleType:   cachestore.BlockStorageGroupRoleType(blockStorage.RoleType.ValueString()),
				SizeGb:     blockStorage.SizeGb.ValueInt32(),
				VolumeType: cachestore.VolumeType(blockStorage.VolumeType.ValueString()).Ptr(),
			})
		}

		var convertedInstance []cachestore.InstanceRequest
		for _, instance := range instanceGroup.Instances {
			convertedInstance = append(convertedInstance, cachestore.InstanceRequest{
				RoleType:         cachestore.InstanceRoleType(instance.RoleType.ValueString()),
				ServiceIpAddress: *cachestore.NewNullableString(instance.ServiceIpAddress.ValueStringPointer()),
				PublicIpId:       *cachestore.NewNullableString(instance.PublicIpId.ValueStringPointer()),
			})
		}

		convertedInstanceGroups = append(convertedInstanceGroups, cachestore.RedisInstanceGroupRequest{
			BlockStorageGroups: convertedBlockStorage,
			Instances:          convertedInstance,
			RoleType:           cachestore.InstanceGroupRoleType(instanceGroup.RoleType.ValueString()),
			ServerTypeName:     instanceGroup.ServerTypeName.ValueString(),
		})
	}

	// MaintenanceOption
	var maintenanceOption = request.MaintenanceOption
	var convertedMaintenanceOption = &cachestore.MaintenanceOption{}

	if client.CheckMaintenanceOption(maintenanceOption) {
		convertedMaintenanceOption = nil
	} else {
		var startingDayOfWeek, _ = cachestore.NewDayOfWeekFromValue(maintenanceOption.StartingDayOfWeek.ValueString())
		convertedMaintenanceOption = &cachestore.MaintenanceOption{
			PeriodHour:        maintenanceOption.PeriodHour.ValueStringPointer(),
			StartingTime:      maintenanceOption.StartingTime.ValueStringPointer(),
			StartingDayOfWeek: startingDayOfWeek,
		}
	}

	//Tags
	var TagsObject []cachestore.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := cachestore.Tag{
			Key:   &k,
			Value: *cachestore.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.RedisClusterCreateRequest(cachestore.RedisClusterCreateRequest{
		AllowableIpAddresses: allowableIpAddresses,
		DbaasEngineVersionId: request.DbaasEngineVersionId.ValueString(),
		HaEnabled:            request.HaEnabled.ValueBoolPointer(),
		InitConfigOption:     convertedInitConfigOption,
		InstanceGroups:       convertedInstanceGroups,
		InstanceNamePrefix:   request.InstanceNamePrefix.ValueString(),
		Name:                 request.Name.ValueString(),
		NatEnabled:           request.NatEnabled.ValueBoolPointer(),
		ReplicaCount:         *cachestore.NewNullableInt32(request.ReplicaCount.ValueInt32Pointer()),
		SubnetId:             request.SubnetId.ValueString(),
		Timezone:             request.Timezone.ValueString(),
		MaintenanceOption:    *cachestore.NewNullableMaintenanceOption(convertedMaintenanceOption),
		Tags:                 TagsObject,
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

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*cachestore.RedisClusterDetailResponse, error) {
	req := client.sdkClient.CachestoreV1CacheStoreClustersApiAPI.CachestoreShowCluster(ctx, clusterId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.CachestoreV1CacheStoreClustersApiAPI.CachestoreDeleteCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StopCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.CachestoreV1CacheStoreClustersApiAPI.CachestoreStopCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StartCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.CachestoreV1CacheStoreClustersApiAPI.CachestoreStartCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBackup(ctx context.Context, clusterId string, startingTimeHour string, retentionPeriodDay string) error {
	req := client.sdkClient.CachestoreV1CacheStoreBackupApiAPI.CachestoreSetBackup(ctx, clusterId)

	req = req.BackupSettingExcludingArchiveRequest(cachestore.BackupSettingExcludingArchiveRequest{
		RetentionPeriodDay: retentionPeriodDay,
		StartingTimeHour:   startingTimeHour,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) UnSetBackup(ctx context.Context, clusterId string) error {
	req := client.sdkClient.CachestoreV1CacheStoreBackupApiAPI.CachestoreUnsetBackup(ctx, clusterId)

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetSecurityGroupRules(ctx context.Context, clusterId string, addedIPs []string, removedIps []string) error {
	req := client.sdkClient.CachestoreV1CacheStoreInstancesApiAPI.CachestoreSetSecurityGroupRules(ctx, clusterId)

	req = req.UpdateSecurityGroupRulesRequest(cachestore.UpdateSecurityGroupRulesRequest{
		AddIpAddresses: addedIPs,
		DelIpAddresses: removedIps,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetServerType(ctx context.Context, instanceGroupId string, serverTypeName string) error {
	req := client.sdkClient.CachestoreV1CacheStoreInstancesApiAPI.CachestoreSetServerType(ctx, instanceGroupId)
	reqState := &cachestore.InstanceGroupResizeRequest{ServerTypeName: serverTypeName}
	req = req.InstanceGroupResizeRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBlockStorageSize(ctx context.Context, blockStorageGroupId string, sizeGb int32) error {
	req := client.sdkClient.CachestoreV1CacheStoreInstancesApiAPI.CachestoreSetBlockStorageSize(ctx, blockStorageGroupId)
	reqState := &cachestore.ResizeBlockStorageGroupRequest{SizeGb: sizeGb}
	req = req.ResizeBlockStorageGroupRequest(*reqState)
	_, _, err := req.Execute()
	return err
}
