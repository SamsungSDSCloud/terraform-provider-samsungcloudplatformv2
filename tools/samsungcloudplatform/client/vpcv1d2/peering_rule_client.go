package vpcv1d2

import (
	"context"

	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpc/1.2"
)

func (client *Client) GetVpcPeeringRuleList(ctx context.Context, request VpcPeeringRuleDataSource) (*scpvpc.VpcPeeringRuleListResponse, error) {
	req := client.sdkClient.VpcV1VpcPeeringRuleApiAPI.ListVpcPeeringRules(ctx, request.VpcPeeringId.ValueString())

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
	if !request.SourceVpcId.IsNull() {
		req = req.SourceVpcId(request.SourceVpcId.ValueString())
	}
	if !request.SourceVpcType.IsNull() {
		req = req.SourceVpcType(scpvpc.VpcPeeringRuleDestinationVpcType(request.SourceVpcType.ValueString()))
	}
	if !request.DestinationVpcId.IsNull() {
		req = req.DestinationVpcId(request.DestinationVpcId.ValueString())
	}
	if !request.DestinationVpcType.IsNull() {
		req = req.DestinationVpcType(scpvpc.VpcPeeringRuleDestinationVpcType(request.DestinationVpcType.ValueString()))
	}
	if !request.DestinationCidr.IsNull() {
		req = req.DestinationCidr(request.DestinationCidr.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(scpvpc.VpcPeeringState(request.State.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}
