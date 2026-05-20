package vpcv1d2

import (
	"context"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.2"
)

func (client *Client) CreateTgwUplinkRule(ctx context.Context, request TransitGatewayUplinkValue) (*vpc.TransitGatewayUplinkRuleCreateResponse, error) {
	req := client.sdkClient.VpcV1TransitGatewayUplinkRulesApiAPI.CreateTransitGatewayUplinkRule(ctx, request.TransitGatewayId.ValueString())

	req = req.TransitGatewayUplinkRuleCreateRequest(vpc.TransitGatewayUplinkRuleCreateRequest{
		Description:     request.Description.ValueStringPointer(),
		DestinationCidr: request.DestinationCidr.ValueString(),
		DestinationType: vpc.TransitGatewayUplinkRuleDestinationType(request.DestinationType.ValueString()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteTgwUplinkRule(ctx context.Context, transitGatewayId string, uplinkId string) error {
	req := client.sdkClient.VpcV1TransitGatewayUplinkRulesApiAPI.DeleteTransitGatewayUplinkRule(ctx, transitGatewayId, uplinkId)

	_, err := req.Execute()
	return err
}
