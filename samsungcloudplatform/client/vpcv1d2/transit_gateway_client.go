package vpcv1d2

import (
	"context"
	"net/http"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.2"
)

func (client *Client) GetTransitGatewayList(ctx context.Context, request TgwDataSource) (*vpc.TransitGatewayListResponseV1Dot2, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.ListTransitGateways(ctx)
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
	if !request.State.IsNull() {
		req = req.State(vpc.TransitGatewayState(request.State.ValueString()))
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}
	if !request.FirewallConnectionState.IsNull() {
		req = req.Id(request.FirewallConnectionState.ValueString())
	}

	resp, httpResp, err := req.Execute()
	return resp, httpResp, err
}

func (client *Client) GetTransitGatewayInfo(ctx context.Context, transitGatewayId string) (*vpc.TransitGatewayShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.ShowTransitGateway(ctx, transitGatewayId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateTransitGateway(ctx context.Context, dataReq TgwResource) (*vpc.TransitGatewayShowResponseV1Dot2, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.CreateTransitGateway(ctx)

	req = req.TransitGatewayCreateRequest(vpc.TransitGatewayCreateRequest{
		Description: *vpc.NewNullableString(dataReq.Description.ValueStringPointer()),
		Name:        dataReq.Name.ValueString(),
		Tags:        convertToTags(dataReq.Tags.Elements()),
	})

	resp, httpResp, err := req.Execute()
	return resp, httpResp, err

}

func (client *Client) UpdateTransitGateway(ctx context.Context, dataReq TgwResource) error {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.SetTransitGateway(ctx, dataReq.Id.ValueString())
	if !dataReq.Description.IsNull() {
		req = req.TransitGatewaySetRequest(vpc.TransitGatewaySetRequest{
			Description: dataReq.Description.ValueStringPointer(),
		})
	}

	_, _, err := req.Execute()
	return err
}

func (client *Client) DeleteTransitGateway(ctx context.Context, transitGatewayId string) error {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.DeleteTransitGateway(ctx, transitGatewayId)

	_, err := req.Execute()
	return err
}
