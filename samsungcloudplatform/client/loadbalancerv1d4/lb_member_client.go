package loadbalancerv1d4

import (
	"context"

	loadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
)

// GetLbMembersV1d4 retrieves the list of LB members for a given server group using the v1.4 API.
func (client *Client) GetLbMembersV1d4(ctx context.Context, request LbMemberDataSourceV1d4) (*loadbalancer.MemberWithHealthStateListResponseV1Dot4, error) {
	req := client.sdkClient.LoadbalancerV1MemberApiAPI.ListLbServerGroupMembers(ctx, request.LbServerGroupId.ValueString())

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
	if !request.MemberIp.IsNull() {
		req = req.MemberIp(request.MemberIp.ValueString())
	}
	if !request.MemberPort.IsNull() {
		req = req.MemberPort(request.MemberPort.ValueInt32())
	}

	resp, _, err := req.Execute()
	return resp, err
}
