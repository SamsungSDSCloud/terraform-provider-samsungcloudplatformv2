package vpcv1d2

import (
	"context"

	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.2"
)

func (client *Client) ListInternetGateways(ctx context.Context, request InternetGatewayDataSource) (*scpvpc.InternetGatewayListResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.ListInternetGateways(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.Type.IsNull() {
		req = req.Type_(scpvpc.InternetGatewayType(request.Type.ValueString()))
	}
	if !request.State.IsNull() {
		req = req.State(request.State.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}
