package vpcv1d2

import (
	"context"

	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.2"
)

func (client *Client) ListNatGateways(ctx context.Context, request NatGatewayDataSource) (*scpvpc.NatGatewayListResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.ListNatGateways(ctx)

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
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.SubnetId.IsNull() {
		req = req.SubnetId(request.SubnetId.ValueString())
	}
	if !request.SubnetName.IsNull() {
		req = req.SubnetName(request.SubnetName.ValueString())
	}
	if !request.NatGatewayIpAddress.IsNull() {
		req = req.NatGatewayIpAddress(request.NatGatewayIpAddress.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(scpvpc.NatGatewayState(request.State.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}