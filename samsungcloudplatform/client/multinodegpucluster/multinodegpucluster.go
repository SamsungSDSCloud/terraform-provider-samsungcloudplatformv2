package multinodegpucluster

import (
	"context"
	"fmt"
	"math"
	"net/http"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/client"
	multinodegpuclustersdk1d3 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/library/multinodegpucluster/1.3"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *multinodegpuclustersdk1d3.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: multinodegpuclustersdk1d3.NewAPIClient(config),
	}
}

func (client *Client) GetGpuNodeList(ctx context.Context, state GpuNodeList) (*multinodegpuclustersdk1d3.GpuNodeListResponseV1Dot3, error) {

	req := client.sdkClient.MultinodegpuclusterV1GpuNodesAPIsAPI.ListGpuNodes(ctx)
	req = req.Size(math.MaxInt32)

	if state.GpuNodeName.ValueString() != "" {
		req = req.GpuNodeName(state.GpuNodeName.ValueString())
	}
	if state.State.ValueString() != "" {
		req = req.State(state.State.ValueString())
	}
	if state.Ip.ValueString() != "" {
		req = req.Ip(state.Ip.ValueString())
	}
	if state.VpcId.ValueString() != "" {
		req = req.VpcId(state.VpcId.ValueString())
	}
	if state.ClusterFabricName.ValueString() != "" {
		req = req.ClusterFabricName(state.ClusterFabricName.ValueString())
	}
	if state.ClusterFabricId.ValueString() != "" {
		req = req.ClusterFabricId(state.ClusterFabricId.ValueString())
	}
	if state.Zone.ValueString() != "" {
		req = req.Zone(state.Zone.ValueString())
	}

	req = req.Sort("gpu_node_name:asc")

	resp, _, err := req.Execute()

	return resp, err

}

func (client *Client) GetGpuNode(ctx context.Context, gpunodeId string) (*multinodegpuclustersdk1d3.GpuNodeShowResponseV1Dot3, *http.Response, error) {

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

	req = req.GpuNodeOperationRequest(multinodegpuclustersdk1d3.GpuNodeOperationRequest{GpuNodeIds: nodesIds})

	_, _, err := req.Execute()
	return err
}

func (client *Client) StartGpuNodes(ctx context.Context, ids []string) error {
	req := client.sdkClient.MultinodegpuclusterV1GpuNodesAPIsAPI.StartGpuNodes(ctx)

	nodesIds := make([]interface{}, 0)
	for _, id := range ids {
		nodesIds = append(nodesIds, id)
	}

	req = req.GpuNodeOperationRequest(multinodegpuclustersdk1d3.GpuNodeOperationRequest{GpuNodeIds: nodesIds})

	_, _, err := req.Execute()
	return err
}

func (client *Client) DeleteGpuNodes(ctx context.Context, ids []string) error {
	req := client.sdkClient.MultinodegpuclusterV1GpuNodesAPIsAPI.DeleteGpuNodes(ctx)

	nodesIds := make([]interface{}, 0)
	for _, id := range ids {
		nodesIds = append(nodesIds, id)
	}

	req = req.GpuNodeTerminateRequest(multinodegpuclustersdk1d3.GpuNodeTerminateRequest{GpuNodeIds: nodesIds})

	_, _, err := req.Execute()
	return err
}

func (client *Client) CreateGpuNode(ctx context.Context, request GpuNodeResource, draft GpuNodeResource) (*multinodegpuclustersdk1d3.AsyncResponse, error) {
	req := client.sdkClient.MultinodegpuclusterV1GpuNodesAPIsAPI.CreateGpuNodes(ctx)

	tags := make([]multinodegpuclustersdk1d3.Tag, 0)
	for k, v := range request.Tags.Elements() {
		tag := multinodegpuclustersdk1d3.Tag{}

		key := multinodegpuclustersdk1d3.NullableString{}
		key.Set(&k)
		tag.Key = key

		if v != nil {
			value := multinodegpuclustersdk1d3.NullableString{}
			value.Set(v.(types.String).ValueStringPointer())
			tag.Value = value
		}
		tags = append(tags, tag)
	}

	initScript := multinodegpuclustersdk1d3.NullableString{}
	initScript.Set(request.InitScript.ValueStringPointer())

	var requestServerDetails []ServerDetailsValue
	serverDetails := make([]multinodegpuclustersdk1d3.GpuNodeDetailsRequestV1Dot3, 0)

	request.ServerDetails.ElementsAs(ctx, &requestServerDetails, false)

	for pos, requestServerDetail := range requestServerDetails {

		ipAddress := multinodegpuclustersdk1d3.NullableString{}
		ipAddress.Set(requestServerDetail.IpAddress.ValueStringPointer())

		serverDetail := multinodegpuclustersdk1d3.GpuNodeDetailsRequestV1Dot3{
			GpuNodeName:       fmt.Sprintf("%s-%03d", draft.GpuNodeNamePrefix.ValueString(), pos+1),
			IpAddress:         ipAddress,
			NatEnabled:        false,
			PublicIpAddressId: *multinodegpuclustersdk1d3.NewNullableString(new(string)),
			ServerTypeId:      draft.ServerTypeId.ValueString(),
			Zone:              requestServerDetail.Zone.ValueString(),
		}

		serverDetails = append(serverDetails, serverDetail)
	}

	clusterFabricDetailsReq := multinodegpuclustersdk1d3.ClusterFabricDetailsRequest{
		ClusterFabricId:   *multinodegpuclustersdk1d3.NewNullableString(new(string)),
		ClusterFabricName: request.ClusterFabricDetails.ClusterFabricName.ValueString(),
		NodePoolId:        request.ClusterFabricDetails.NodePoolId.ValueString(),
	}

	clusterFabricDetailsReq.ClusterFabricId.Set(request.ClusterFabricDetails.ClusterFabricId.ValueStringPointer())

	req = req.GpuNodeCreateRequestV1Dot3(multinodegpuclustersdk1d3.GpuNodeCreateRequestV1Dot3{
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

func (client *Client) GetImageList(ctx context.Context, regionId string) (*multinodegpuclustersdk1d3.GpuNodeImageListResponse, error) {
	req := client.sdkClient.MultinodegpuclusterV1GpuNodeImageAPIsAPI.ListGpuNodeImages(ctx)
	req = req.RegionId(regionId)
	resp, _, err := req.Execute()
	return resp, err
}
