package directconnect

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	directconnect "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/direct-connect/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

//------------ Direct Connect -------------------//

func (client *Client) GetDirectConnectList(ctx context.Context, request DirectConnectDataSource) (*directconnect.DirectConnectListResponse, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsApiAPI.ListDirectConnects(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
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

func (client *Client) CreateDirectConnect(ctx context.Context, request DirectConnectResource) (*directconnect.DirectConnectShowResponse, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsApiAPI.CreateDirectConnect(ctx)

	var tags []directconnect.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := directconnect.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}

	req = req.DirectConnectCreateRequest(directconnect.DirectConnectCreateRequest{
		Bandwidth:        request.Bandwidth.ValueInt32(),
		Description:      request.Description.ValueStringPointer(),
		FirewallEnabled:  request.FirewallEnabled.ValueBoolPointer(),
		FirewallLoggable: request.FirewallLoggable.ValueBoolPointer(),
		Name:             request.Name.ValueString(),
		VpcId:            request.VpcId.ValueString(),
		Tags:             tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetDirectConnect(ctx context.Context, directConnectId string) (*directconnect.DirectConnectShowResponse, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsApiAPI.ShowDirectConnect(ctx, directConnectId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateDirectConnect(ctx context.Context, directConnectId string, request DirectConnectResource) (*directconnect.DirectConnectShowResponse, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsApiAPI.SetDirectConnect(ctx, directConnectId)

	req = req.DirectConnectSetRequest(directconnect.DirectConnectSetRequest{
		Description: request.Description.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteDirectConnect(ctx context.Context, directConnectId string) error {
	req := client.sdkClient.DirectConnectV1DirectConnectsApiAPI.DeleteDirectConnect(ctx, directConnectId)

	_, err := req.Execute()
	return err
}

//------------ Routing Rule -------------------//

func (client *Client) GetRoutingRuleList(ctx context.Context, request RoutingRuleDataSource) (*directconnect.RoutingRuleListResponse, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsRulesApiAPI.ListRoutingRules(ctx, request.DirectConnectId.ValueString())
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
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

func (client *Client) CreateRoutingRule(ctx context.Context, request RoutingRuleResource) (*directconnect.RoutingRuleShowResponse, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsRulesApiAPI.CreateRoutingRule(ctx, request.DirectConnectId.ValueString())
	destinationResourceId := request.DestinationResourceId.ValueString()
	destinationResourceIdNS := directconnect.NullableString{}
	if destinationResourceId != "" {
		destinationResourceIdNS.Set(&destinationResourceId)
	}

	req = req.RoutingRuleCreateRequest(directconnect.RoutingRuleCreateRequest{
		Description:           request.Description.ValueStringPointer(),
		DestinationCidr:       request.DestinationCidr.ValueString(),
		DestinationResourceId: destinationResourceIdNS,
		DestinationType:       directconnect.DirectConnectRoutingRuleDestinationType(request.DestinationType.ValueString()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetRoutingRule(ctx context.Context, directConnectId string, routingRuleId string) (*directconnect.RoutingRuleListResponse, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsRulesApiAPI.ListRoutingRules(ctx, directConnectId)

	if routingRuleId != "" {
		req = req.Id(routingRuleId)
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteRoutingRule(ctx context.Context, DirectConnectId string, RoutingRuleId string) error {
	req := client.sdkClient.DirectConnectV1DirectConnectsRulesApiAPI.DeleteRoutingRule(ctx, DirectConnectId, RoutingRuleId)

	_, err := req.Execute()
	return err
}
