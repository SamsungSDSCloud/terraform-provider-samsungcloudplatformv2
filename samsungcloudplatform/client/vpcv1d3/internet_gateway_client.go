package vpcv1d3

import (
	"context"

	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
)

func (client *Client) GetInternetGateway(ctx context.Context, internetGatewayId string) (*scpvpc.InternetGatewayShowResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.ShowInternetGateway(ctx, internetGatewayId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateInternetGateway(ctx context.Context, request InternetGatewayResource) (*scpvpc.InternetGatewayShowResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.CreateInternetGateway(ctx)
	description := request.Description.ValueString()
	descriptionNS := scpvpc.NullableString{}
	descriptionNS.Set(&description)

	tags := convertToTags(request.Tags.Elements())

	req = req.InternetGatewayCreateRequest(scpvpc.InternetGatewayCreateRequest{
		Type:             scpvpc.InternetGatewayType(request.Type.ValueString()),
		Description:      descriptionNS,
		FirewallEnabled:  request.FirewallEnabled.ValueBoolPointer(),
		FirewallLoggable: request.FirewallLoggable.ValueBoolPointer(),
		VpcId:            request.VpcId.ValueString(),
		Tags:             tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) ListInternetGateways(ctx context.Context, request InternetGatewayDataSource) (*scpvpc.InternetGatewayListResponseV1Dot3, error) {
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

func (client *Client) UpdateInternetGateway(ctx context.Context, internetGatewayId string, request InternetGatewayResource) (*scpvpc.InternetGatewaySetResponse, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.SetInternetGateway(ctx, internetGatewayId)
	description := request.Description.ValueString()
	descriptionNS := scpvpc.NullableString{}
	descriptionNS.Set(&description)

	req = req.InternetGatewaySetRequest(scpvpc.InternetGatewaySetRequest{
		Description: descriptionNS,
		Loggable:    request.Loggable.ValueBoolPointer(),
	})

	resp, _, err := req.Execute()
	return resp, err
}
