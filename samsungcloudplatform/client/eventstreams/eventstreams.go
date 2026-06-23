package eventstreams

import (
	"context"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	eventstreams "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/eventstreams/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/database"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *eventstreams.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: eventstreams.NewAPIClient(config),
	}
}

// clusterlist (ctx, ClusterDataSource) - (RdbClusterPageResponse)
func (client *Client) GetClusterList(ctx context.Context, request ClusterDataSource) (*eventstreams.ClusterPageResponse, error) {
	req := client.sdkClient.EventstreamsV1EventStreamsClustersApiAPI.EventstreamsListClusters(ctx)
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
func (client *Client) GetEngineVersionList(ctx context.Context) (*eventstreams.EngineListResponse, error) {
	req := client.sdkClient.EventstreamsV1EventStreamsMasterDataApiAPI.EventstreamsListEngineVersions(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

// create (ctx, clusterResource) - (asyncResponse)
func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*eventstreams.AsyncResponse, error) {
	req := client.sdkClient.EventstreamsV1EventStreamsClustersApiAPI.EventstreamsCreateCluster(ctx)

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

	var convertedInitConfigOption = eventstreams.EventStreamsInitConfigOptionRequest{
		AkhqId:                *eventstreams.NewNullableString(initConfigOption.AkhqId.ValueStringPointer()),
		AkhqPassword:          *eventstreams.NewNullableString(initConfigOption.AkhqPassword.ValueStringPointer()),
		BrokerPort:            initConfigOption.BrokerPort.ValueInt32Pointer(),
		BrokerSaslId:          initConfigOption.BrokerSaslId.ValueString(),
		BrokerSaslPassword:    initConfigOption.BrokerSaslPassword.ValueString(),
		ZookeeperPort:         initConfigOption.ZookeeperPort.ValueInt32Pointer(),
		ZookeeperSaslId:       initConfigOption.ZookeeperSaslId.ValueString(),
		ZookeeperSaslPassword: initConfigOption.ZookeeperSaslPassword.ValueString(),
	}

	// InstanceGroups
	var convertedInstanceGroups []eventstreams.InstanceGroupRequest
	var igVals []database.InstanceGroup
	request.InstanceGroups.ElementsAs(context.Background(), &igVals, false)
	for _, instanceGroup := range igVals {
		var convertedBlockStorage []eventstreams.BlockStorageGroupRequest
		var bsVals []database.BlockStorageGroup
		instanceGroup.BlockStorageGroups.ElementsAs(context.Background(), &bsVals, false)
		for _, blockStorage := range bsVals {
			convertedBlockStorage = append(convertedBlockStorage, eventstreams.BlockStorageGroupRequest{
				RoleType:   eventstreams.BlockStorageGroupRoleType(blockStorage.RoleType.ValueString()),
				SizeGb:     blockStorage.SizeGb.ValueInt32(),
				VolumeType: eventstreams.VolumeType(blockStorage.VolumeType.ValueString()).Ptr(),
			})
		}

		var convertedInstance []eventstreams.InstanceRequest
		var instVals []database.Instance
		instanceGroup.Instances.ElementsAs(context.Background(), &instVals, false)
		for _, instance := range instVals {
			convertedInstance = append(convertedInstance, eventstreams.InstanceRequest{
				RoleType:         eventstreams.InstanceRoleType(instance.RoleType.ValueString()),
				ServiceIpAddress: *eventstreams.NewNullableString(instance.ServiceIpAddress.ValueStringPointer()),
				PublicIpId:       *eventstreams.NewNullableString(instance.PublicIpId.ValueStringPointer()),
			})
		}

		convertedInstanceGroups = append(convertedInstanceGroups, eventstreams.InstanceGroupRequest{
			BlockStorageGroups: convertedBlockStorage,
			Instances:          convertedInstance,
			RoleType:           eventstreams.InstanceGroupRoleType(instanceGroup.RoleType.ValueString()),
			ServerTypeName:     instanceGroup.ServerTypeName.ValueString(),
		})
	}

	// MaintenanceOption
	var maintenanceOption = request.MaintenanceOption
	var convertedMaintenanceOption = &eventstreams.MaintenanceOption{}

	if client.CheckMaintenanceOption(maintenanceOption) {
		convertedMaintenanceOption = nil
	} else {
		var startingDayOfWeek, _ = eventstreams.NewDayOfWeekFromValue(maintenanceOption.StartingDayOfWeek.ValueString())
		convertedMaintenanceOption = &eventstreams.MaintenanceOption{
			PeriodHour:        maintenanceOption.PeriodHour.ValueStringPointer(),
			StartingTime:      maintenanceOption.StartingTime.ValueStringPointer(),
			StartingDayOfWeek: startingDayOfWeek,
		}
	}

	//Tags
	var TagsObject []eventstreams.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := eventstreams.Tag{
			Key:   &k,
			Value: *eventstreams.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.EventStreamsClusterCreateRequestV1Dot1(eventstreams.EventStreamsClusterCreateRequestV1Dot1{
		AkhqEnabled:               request.AkhqEnabled.ValueBoolPointer(),
		AllowableIpAddresses:      allowableIpAddresses,
		DbaasEngineVersionId:      request.DbaasEngineVersionId.ValueString(),
		InitConfigOption:          convertedInitConfigOption,
		InstanceGroups:            convertedInstanceGroups,
		InstanceNamePrefix:        request.InstanceNamePrefix.ValueString(),
		IsCombined:                request.IsCombined.ValueBoolPointer(),
		NatEnabled:                request.NatEnabled.ValueBoolPointer(),
		Name:                      request.Name.ValueString(),
		SubnetId:                  request.SubnetId.ValueString(),
		Timezone:                  request.Timezone.ValueString(),
		MaintenanceOption:         *eventstreams.NewNullableMaintenanceOption(convertedMaintenanceOption),
		ServiceWatchLogCollection: *eventstreams.NewNullableBool(request.ServiceWatchLogCollection.ValueBoolPointer()),
		Tags:                      TagsObject,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CheckMaintenanceOption(maintenanceOption MaintenanceOption) bool {
	return !maintenanceOption.UseMaintenanceOption.ValueBool() || (maintenanceOption.StartingDayOfWeek.IsNull() && maintenanceOption.StartingTime.IsNull() && maintenanceOption.PeriodHour.IsNull())
}

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*eventstreams.EventStreamsClusterDetailResponseV1Dot1, error) {
	req := client.sdkClient.EventstreamsV1EventStreamsClustersApiAPI.EventstreamsShowCluster(ctx, clusterId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.EventstreamsV1EventStreamsClustersApiAPI.EventstreamsDeleteCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StopCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.EventstreamsV1EventStreamsClustersApiAPI.EventstreamsStopCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) StartCluster(ctx context.Context, clusterId string) error {
	req := client.sdkClient.EventstreamsV1EventStreamsClustersApiAPI.EventstreamsStartCluster(ctx, clusterId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetSecurityGroupRules(ctx context.Context, clusterId string, addedIPs []string, removedIps []string) error {
	req := client.sdkClient.EventstreamsV1EventStreamsInstancesApiAPI.EventstreamsSetSecurityGroupRules(ctx, clusterId)

	req = req.UpdateSecurityGroupRulesRequest(eventstreams.UpdateSecurityGroupRulesRequest{
		AddIpAddresses: addedIPs,
		DelIpAddresses: removedIps,
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) SetServerType(ctx context.Context, instanceGroupId string, serverTypeName string) error {
	req := client.sdkClient.EventstreamsV1EventStreamsInstancesApiAPI.EventstreamsSetServerType(ctx, instanceGroupId)
	reqState := &eventstreams.InstanceGroupResizeRequest{ServerTypeName: serverTypeName}
	req = req.InstanceGroupResizeRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) SetBlockStorageSize(ctx context.Context, blockStorageGroupId string, sizeGb int32) error {
	req := client.sdkClient.EventstreamsV1EventStreamsInstancesApiAPI.EventstreamsSetBlockStorageSize(ctx, blockStorageGroupId)
	reqState := &eventstreams.ResizeBlockStorageGroupRequest{SizeGb: sizeGb}
	req = req.ResizeBlockStorageGroupRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func (client *Client) AddInstances(ctx context.Context, clusterId string, instanceCount int32, serviceIPAddresses []string) error {
	req := client.sdkClient.EventstreamsV1EventStreamsClustersApiAPI.EventstreamsAddInstances(ctx, clusterId)
	reqState := &eventstreams.EventStreamsClusterAddInstancesRequest{
		InstanceCount:      instanceCount,
		ServiceIpAddresses: serviceIPAddresses,
	}
	req = req.EventStreamsClusterAddInstancesRequest(*reqState)
	_, _, err := req.Execute()
	return err
}

func MapInstanceGroupResponses(sdkResp []eventstreams.InstanceGroupResponse) []database.InstanceGroupResponse {
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
			var pubIP, serviceIP, pubIPAddr string
			if it.ServiceIpAddress.Get() != nil {
				serviceIP = *it.ServiceIpAddress.Get()
			}
			if it.PublicIpId.Get() != nil {
				pubIP = *it.PublicIpId.Get()
			}
			if it.PublicIpAddress.Get() != nil {
				pubIPAddr = *it.PublicIpAddress.Get()
			}

			instances[j] = database.InstanceResponse{
				Name:             it.Name,
				PublicIpAddress:  pubIPAddr,
				PublicIpId:       pubIP,
				RoleType:         string(it.RoleType),
				ServiceIpAddress: serviceIP,
				ServiceState:     string(it.ServiceState),
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
