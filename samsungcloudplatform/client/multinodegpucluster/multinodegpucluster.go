package multinodegpucluster

import (
	"context"
	"fmt"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	multinodegpuclustersdk1d2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/multinodegpucluster/1.2"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math"
	"net/http"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *multinodegpuclustersdk1d2.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: multinodegpuclustersdk1d2.NewAPIClient(config),
	}
}

func (client *Client) GetGpuNodeList(ctx context.Context, gpuNodeName types.String, state types.String, ip types.String, vpcId types.String, clusterFabricName types.String, clusterFabricId types.String) (*multinodegpuclustersdk1d2.GpuNodeListResponse, error) {

	req := client.sdkClient.MultinodegpuclusterV1GpuNodesAPIsAPI.ListGpuNodes(ctx)
	req = req.Size(math.MaxInt32)

	if !gpuNodeName.IsNull() {
		req = req.GpuNodeName(gpuNodeName.ValueString())
	}
	if !state.IsNull() {
		req = req.State(state.ValueString())
	}
	if !ip.IsNull() {
		req = req.Ip(ip.ValueString())
	}
	if !vpcId.IsNull() {
		req = req.VpcId(vpcId.ValueString())
	}
	if !clusterFabricName.IsNull() {
		req = req.ClusterFabricName(clusterFabricName.ValueString())
	}
	if !clusterFabricId.IsNull() {
		req = req.ClusterFabricId(clusterFabricId.ValueString())
	}

	req = req.Sort("gpu_node_name:asc")

	resp, _, err := req.Execute()

	return resp, err

}

func (client *Client) GetGpuNode(ctx context.Context, gpunodeId string) (*multinodegpuclustersdk1d2.GpuNodeShowResponse, *http.Response, error) {

	req := client.sdkClient.MultinodegpuclusterV1GpuNodesAPIsAPI.ShowGpuNode(ctx, gpunodeId)

	resp, httpResponse, err := req.Execute()

	return resp, httpResponse, err

}

func (client *Client) StopGpuNodes(ctx context.Context, ids []string) error {
	req := client.sdkClient.MultinodegpuclusterV1GpuNodesAPIsAPI.StopGpuNodes(ctx)

	nodesIds := make([]interface{}, 0)
	for _, id := range ids {
		nodesIds = append(nodesIds, id)
	}

	req = req.GpuNodeOperationRequest(multinodegpuclustersdk1d2.GpuNodeOperationRequest{GpuNodeIds: nodesIds})

	_, _, err := req.Execute()
	return err
}

func (client *Client) StartGpuNodes(ctx context.Context, ids []string) error {
	req := client.sdkClient.MultinodegpuclusterV1GpuNodesAPIsAPI.StartGpuNodes(ctx)

	nodesIds := make([]interface{}, 0)
	for _, id := range ids {
		nodesIds = append(nodesIds, id)
	}

	req = req.GpuNodeOperationRequest(multinodegpuclustersdk1d2.GpuNodeOperationRequest{GpuNodeIds: nodesIds})

	_, _, err := req.Execute()
	return err
}

func (client *Client) DeleteGpuNodes(ctx context.Context, ids []string) error {
	req := client.sdkClient.MultinodegpuclusterV1GpuNodesAPIsAPI.DeleteGpuNodes(ctx)

	nodesIds := make([]interface{}, 0)
	for _, id := range ids {
		nodesIds = append(nodesIds, id)
	}

	req = req.GpuNodeTerminateRequest(multinodegpuclustersdk1d2.GpuNodeTerminateRequest{GpuNodeIds: nodesIds})

	_, _, err := req.Execute()
	return err
}

func (client *Client) CreateGpuNode(ctx context.Context, request GpuNodeResource, draft GpuNodeResource) (*multinodegpuclustersdk1d2.AsyncResponse, error) {
	req := client.sdkClient.MultinodegpuclusterV1GpuNodesAPIsAPI.CreateGpuNodes(ctx)

	tags := make([]multinodegpuclustersdk1d2.Tag, 0)
	for k, v := range request.Tags.Elements() {
		tag := multinodegpuclustersdk1d2.Tag{}

		key := multinodegpuclustersdk1d2.NullableString{}
		key.Set(&k)
		tag.Key = key

		if v != nil {
			value := multinodegpuclustersdk1d2.NullableString{}
			value.Set(v.(types.String).ValueStringPointer())
			tag.Value = value
		}
		tags = append(tags, tag)
	}

	initScript := multinodegpuclustersdk1d2.NullableString{}
	initScript.Set(request.InitScript.ValueStringPointer())

	var requestServerDetails []ServerDetailsValue
	serverDetails := make([]multinodegpuclustersdk1d2.GpuNodeDetailsRequest, 0)

	request.ServerDetails.ElementsAs(ctx, &requestServerDetails, false)

	for pos, requestServerDetail := range requestServerDetails {

		ipAddress := multinodegpuclustersdk1d2.NullableString{}
		ipAddress.Set(requestServerDetail.IpAddress.ValueStringPointer())

		serverDetail := multinodegpuclustersdk1d2.GpuNodeDetailsRequest{
			GpuNodeName:       fmt.Sprintf("%s-%03d", draft.GpuNodeNamePrefix.ValueString(), pos+1),
			IpAddress:         ipAddress,
			NatEnabled:        false,
			PublicIpAddressId: *multinodegpuclustersdk1d2.NewNullableString(new(string)),
			ServerTypeId:      draft.ServerTypeId.ValueString(),
		}

		serverDetails = append(serverDetails, serverDetail)
	}

	clusterFabricDetailsReq := multinodegpuclustersdk1d2.ClusterFabricDetailsRequest{
		ClusterFabricId:   *multinodegpuclustersdk1d2.NewNullableString(new(string)),
		ClusterFabricName: request.ClusterFabricDetails.ClusterFabricName.ValueString(),
		NodePoolId:        request.ClusterFabricDetails.NodePoolId.ValueString(),
	}

	clusterFabricDetailsReq.ClusterFabricId.Set(request.ClusterFabricDetails.ClusterFabricId.ValueStringPointer())

	req = req.GpuNodeCreateRequest(multinodegpuclustersdk1d2.GpuNodeCreateRequest{
		ClusterFabricDetails: clusterFabricDetailsReq,
		ImageId:              request.ImageId.ValueString(),
		InitScript:           initScript,
		LockEnabled:          request.LockEnabled.ValueBoolPointer(),
		OsUserId:             request.OsUserId.ValueString(),
		OsUserPassword:       draft.OsUserPassword.ValueString(),
		RegionId:             request.RegionId.ValueString(),
		ServerDetails:        serverDetails,
		SubnetId:             request.SubnetId.ValueString(),
		Tags:                 tags,
		VpcId:                request.VpcId.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetImageList(ctx context.Context, regionId string) (*multinodegpuclustersdk1d2.GpuNodeImageListResponse, error) {
	req := client.sdkClient.MultinodegpuclusterV1GpuNodeImageAPIsAPI.ListGpuNodeImages(ctx)
	req = req.RegionId(regionId)
	resp, _, err := req.Execute()
	return resp, err
}

