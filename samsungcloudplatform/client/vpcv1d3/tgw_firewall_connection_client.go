package vpcv1d3

import (
	"context"
	"net/http"

	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
)

func (client *Client) CreateTransitGatewayFirewallConnection(ctx context.Context, transitGatewayId string) (*scpvpc.TransitGatewayFirewallConnectionCreateResponse, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.CreateTransitGatewayFirewallConnection(ctx, transitGatewayId)
	resp, status, err := req.Execute()
	return resp, status, err
}

func (client *Client) DeleteTransitGatewayFirewallConnection(ctx context.Context, transitGatewayId string) (*scpvpc.TransitGatewayFirewallConnectionCreateResponse, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.DeleteTransitGatewayFirewallConnection(ctx, transitGatewayId)
	resp, status, err := req.Execute()
	return resp, status, err
}
