package loadbalancerv1d4

import (
	"context"

	loadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
)

// SetLbListenerRule updates the routing rule for a listener.
// PUT /v1/lb-listeners/{listener_id}/rule
func (client *Client) SetLbListenerRule(ctx context.Context, listenerId string,
	body *loadbalancer.LbListenerRuleSetRequest) (*loadbalancer.LbListenerShowResponseV1Dot4, error) {
	req := client.sdkClient.LoadbalancerV1LbListenersApiAPI.SetLbListenerRule(ctx, listenerId)
	req = req.LbListenerRuleSetRequest(*body)
	resp, _, err := req.Execute()
	return resp, err
}
