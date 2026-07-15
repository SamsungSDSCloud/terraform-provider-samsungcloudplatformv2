package directconnect

import (
	"context"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	directconnect "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/direct-connect/1.1"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *directconnect.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: directconnect.NewAPIClient(config),
	}
}

// ------------ Direct Connect List ------------

func (client *Client) GetDirectConnectList(ctx context.Context, request DirectConnectDataSource) (*directconnect.DirectConnectListResponseV1Dot1, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsApiAPI.ListDirectConnects(ctx)

	if !request.Size.IsNull() {
		req = req.Size(int32(request.Size.ValueInt32()))
	}
	if !request.Page.IsNull() {
		req = req.Page(int32(request.Page.ValueInt32()))
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
	if !request.State.IsNull() {
		req = req.State(directconnect.DirectConnectState(request.State.ValueString()))
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

// ------------ Routing Rule List ------------

func (client *Client) GetRoutingRuleList(ctx context.Context, request RoutingRuleDataSource) (*directconnect.RoutingRuleListResponseV1Dot1, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsRulesApiAPI.ListRoutingRules(ctx, request.DirectConnectId.ValueString())

	if !request.Size.IsNull() {
		req = req.Size(int32(request.Size.ValueInt32()))
	}
	if !request.Page.IsNull() {
		req = req.Page(int32(request.Page.ValueInt32()))
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}
	if !request.DestinationType.IsNull() {
		req = req.DestinationType(directconnect.DirectConnectRoutingRuleDestinationType(request.DestinationType.ValueString()))
	}
	if !request.DestinationCidr.IsNull() {
		req = req.DestinationCidr(request.DestinationCidr.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(directconnect.RoutingRuleState(request.State.ValueString()))
	}
	resp, _, err := req.Execute()
	return resp, err
}
