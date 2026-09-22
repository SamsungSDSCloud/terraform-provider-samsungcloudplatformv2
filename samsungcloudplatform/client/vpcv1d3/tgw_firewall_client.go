package vpcv1d3

import (
	"context"
	"net/http"

	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
)

func (client *Client) CreateTransitGatewayFirewall(ctx context.Context, transitGatewayId string, request TransitGatewayFirewallResource) (*scpvpc.TransitGatewayFirewallCreateResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.CreateTransitGatewayFirewall(ctx, transitGatewayId)

	firewallCreateRequest := scpvpc.TransitGatewayFirewallCreateRequestV1Dot3{
		ProductType:       scpvpc.TransitGatewayFirewallProductType(request.ProductType.ValueString()),
		UplinkActiveZone:  *scpvpc.NewNullableString(request.UplinkActiveZone.ValueStringPointer()),
		UplinkStandbyZone: *scpvpc.NewNullableString(request.UplinkStandbyZone.ValueStringPointer()),
	}

	req = req.TransitGatewayFirewallCreateRequestV1Dot3(firewallCreateRequest)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteTransitGatewayFirewall(ctx context.Context, transitGatewayId, firewallId string) (*http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.DeleteTransitGatewayFirewall(ctx, transitGatewayId, firewallId)

	_, httpResp, err := req.Execute()
	return httpResp, err
}
