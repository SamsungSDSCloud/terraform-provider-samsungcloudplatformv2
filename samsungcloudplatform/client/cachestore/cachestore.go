package cachestore

import (
	"context"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	cachestore "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/cachestore/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/database"
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

// engine version
func (client *Client) GetEngineVersionList(ctx context.Context) (*cachestore.EngineListResponse, error) {
	req := client.sdkClient.CachestoreV1CacheStoreMasterDataApiAPI.CachestoreListEngineVersions(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

// create (ctx, clusterResource) - (asyncResponse)
func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*cachestore.AsyncResponse, error) {
	req := client.sdkClient.CachestoreV1CacheStoreClustersApiAPI.CachestoreCreateCluster(ctx)

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
	for _, igElem := range request.InstanceGroups.Elements() {
		igObj := igElem.(types.Object)
		ig := database.InstanceGroup{
			Id:                 igObj.Attributes()["id"].(types.String),
			RoleType:           igObj.Attributes()["role_type"].(types.String),
			ServerTypeName:     igObj.Attributes()["server_type_name"].(types.String),
			BlockStorageGroups: igObj.Attributes()["block_storage_groups"].(types.List),
			Instances:          igObj.Attributes()["instances"].(types.List),
		}

		var convertedBlockStorage []cachestore.RedisBlockStorageGroupRequest
		for _, bsElem := range ig.BlockStorageGroups.Elements() {
			bsObj := bsElem.(types.Object)
			convertedBlockStorage = append(convertedBlockStorage, cachestore.RedisBlockStorageGroupRequest{
				RoleType:   cachestore.BlockStorageGroupRoleType(bsObj.Attributes()["role_type"].(types.String).ValueString()),
				SizeGb:     bsObj.Attributes()["size_gb"].(types.Int32).ValueInt32(),
				VolumeType: cachestore.VolumeType(bsObj.Attributes()["volume_type"].(types.String).ValueString()).Ptr(),
			})
		}

		var convertedInstance []cachestore.InstanceRequest
		for _, instElem := range ig.Instances.Elements() {
			instObj := instElem.(types.Object)
			convertedInstance = append(convertedInstance, cachestore.InstanceRequest{
				RoleType:         cachestore.InstanceRoleType(instObj.Attributes()["role_type"].(types.String).ValueString()),
				ServiceIpAddress: *cachestore.NewNullableString(instObj.Attributes()["service_ip_address"].(types.String).ValueStringPointer()),
				PublicIpId:       *cachestore.NewNullableString(instObj.Attributes()["public_ip_id"].(types.String).ValueStringPointer()),
			})
		}

		convertedInstanceGroups = append(convertedInstanceGroups, cachestore.RedisInstanceGroupRequest{
			BlockStorageGroups: convertedBlockStorage,
			Instances:          convertedInstance,
			RoleType:           cachestore.InstanceGroupRoleType(ig.RoleType.ValueString()),
			ServerTypeName:     ig.ServerTypeName.ValueString(),
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

func MapInstanceGroupResponses(sdkResp []cachestore.RedisInstanceGroupResponse) []database.InstanceGroupResponse {
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
				RoleType:         string(it.RoleType),
				ServiceIpAddress: serviceIP,
				PublicIpId:       pubIP,
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
