package vpcv1d3

import (
	"context"
	"fmt"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
)

func (client *Client) GetNatGateway(ctx context.Context, natGatewayId string) (*vpc.NatGatewayShowResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.ShowNatGateway(ctx, natGatewayId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateNatGateway(ctx context.Context, natGatewayId string, request NatGatewayResource) (*vpc.NatGatewaySetResponse, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.SetNatGateway(ctx, natGatewayId)

	req = req.NatGatewaySetRequest(vpc.NatGatewaySetRequest{
		Description: request.Description.ValueStringPointer(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) ListNatGateways(ctx context.Context, request NatGatewayDataSource) (*vpc.NatGatewayListResponseV1Dot3, error) {
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
	if !request.NatGatewayIpAddress.IsNull() {
		req = req.NatGatewayIpAddress(request.NatGatewayIpAddress.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}
	if !request.SubnetId.IsNull() {
		req = req.SubnetId(request.SubnetId.ValueString())
	}
	if !request.SubnetName.IsNull() {
		req = req.SubnetName(request.SubnetName.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(vpc.NatGatewayState(request.State.ValueString()))
	}
	if !request.MultiZoneEnabled.IsNull() {
		req = req.MultiZoneEnabled(request.MultiZoneEnabled.ValueBool())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateNatGateway(ctx context.Context, request NatGatewayResource) (*vpc.NatGatewayShowResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.CreateNatGateway(ctx)

	var publicipIds []string
	diags := request.PublicipIds.ElementsAs(ctx, &publicipIds, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to parse publicip_ids")
	}

	createReq := vpc.NatGatewayCreateRequestV1Dot3{
		SubnetId:    request.SubnetId.ValueString(),
		PublicipIds: publicipIds,
	}

	if !request.Description.IsNull() {
		desc := request.Description.ValueString()
		createReq.Description = &desc
	}
	if !request.MultiZoneEnabled.IsNull() {
		mzEnabled := request.MultiZoneEnabled.ValueBool()
		nullableBool := vpc.NewNullableBool(&mzEnabled)
		createReq.MultiZoneEnabled = *nullableBool
	}

	req = req.NatGatewayCreateRequestV1Dot3(createReq)

	resp, _, err := req.Execute()
	return resp, err
}
