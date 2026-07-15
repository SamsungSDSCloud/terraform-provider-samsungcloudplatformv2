package vpcv1d2

import (
	"context"
	"net/http"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpc/1.2"
)

func (client *Client) GetTGWRuleList(ctx context.Context, request TransitGatewayRuleDataSource) (*vpc.TransitGatewayRuleListResponseV1Dot2, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayRulesApiAPI.ListTransitGatewayRules(ctx, request.TransitGatewayId.ValueString())

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
	if !request.TgwConnectionVpcId.IsNull() {
		req = req.TgwConnectionVpcId(request.TgwConnectionVpcId.ValueString())
	}
	if !request.TgwConnectionVpcName.IsNull() {
		req = req.TgwConnectionVpcName(request.TgwConnectionVpcName.ValueString())
	}
	if !request.SourceType.IsNull() {
		req = req.SourceType(vpc.TransitGatewayRuleDestinationType(request.SourceType.ValueString()))
	}
	if !request.DestinationType.IsNull() {
		req = req.DestinationType(vpc.TransitGatewayRuleDestinationType(request.DestinationType.ValueString()))
	}
	if !request.DestinationCidr.IsNull() {
		req = req.DestinationCidr(request.DestinationCidr.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(vpc.RoutingRuleState(request.State.ValueString()))
	}
	if !request.RuleType.IsNull() {
		req = req.RuleType(vpc.TransitGatewayRuleType(request.RuleType.ValueString()))
	}

	resp, httpResponse, err := req.Execute()
	return resp, httpResponse, err
}
