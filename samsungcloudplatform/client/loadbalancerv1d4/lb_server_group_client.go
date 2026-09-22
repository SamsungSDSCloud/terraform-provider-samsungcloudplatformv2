package loadbalancerv1d4

import (
	"context"

	loadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
)

// GetLbServerGroupV1d4 retrieves a single LB Server Group by ID using the v1.4 API.
func (client *Client) GetLbServerGroupV1d4(ctx context.Context, lbServerGroupId string) (*loadbalancer.LbServerGroupShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBServerGroupsApiAPI.ShowLbServerGroup(ctx, lbServerGroupId)
	resp, _, err := req.Execute()
	return resp, err
}

// CreateLbServerGroupV1d4 creates a new LB Server Group using the v1.4 API.
func (client *Client) CreateLbServerGroupV1d4(ctx context.Context, request loadbalancer.LbServerGroupCreateRequest) (*loadbalancer.LbServerGroupShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBServerGroupsApiAPI.CreateLbServerGroup(ctx)
	req = req.LbServerGroupCreateRequest(request)
	resp, _, err := req.Execute()
	return resp, err
}
