package loadbalancerv1d4

import (
	"context"

	loadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
)

func (client *Client) CreateLbHealthCheck(ctx context.Context, lbHealthCheck *LbHealthCheckCreate) (*loadbalancer.LbHealthCheckShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBHealthCheckApiAPI.CreateLbHealthCheck(ctx)

	body := loadbalancer.LbHealthCheckCreateRequest{
		LbHealthCheck: loadbalancer.LbHealthCheckCreate{
			Name:                lbHealthCheck.Name.ValueString(),
			VpcId:               lbHealthCheck.VpcId.ValueString(),
			SubnetId:            lbHealthCheck.SubnetId.ValueString(),
			Protocol:            loadbalancer.LbMonitorProtocol(lbHealthCheck.Protocol.ValueString()),
			HealthCheckPort:     *loadbalancer.NewNullableInt32(lbHealthCheck.HealthCheckPort.ValueInt32Pointer()),
			HealthCheckInterval: lbHealthCheck.HealthCheckInterval.ValueInt32Pointer(),
			HealthCheckTimeout:  lbHealthCheck.HealthCheckTimeout.ValueInt32Pointer(),
			HealthCheckCount:    lbHealthCheck.HealthCheckCount.ValueInt32Pointer(),
			HttpMethod:          *loadbalancer.NewNullableLbMonitorHttpMethod((*loadbalancer.LbMonitorHttpMethod)(lbHealthCheck.HttpMethod.ValueStringPointer())),
			HealthCheckUrl:      *loadbalancer.NewNullableString(lbHealthCheck.HealthCheckUrl.ValueStringPointer()),
			ResponseCode:        *loadbalancer.NewNullableString(lbHealthCheck.ResponseCode.ValueStringPointer()),
			RequestData:         *loadbalancer.NewNullableString(lbHealthCheck.RequestData.ValueStringPointer()),
			Description:         *loadbalancer.NewNullableString(lbHealthCheck.Description.ValueStringPointer()),
			Tags:                convertToTags(lbHealthCheck.Tags.Elements()),
		},
	}

	req = req.LbHealthCheckCreateRequest(body)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLbHealthCheck(ctx context.Context, lbHealthCheckId string) (*loadbalancer.LbHealthCheckShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBHealthCheckApiAPI.ShowLbHealthCheck(ctx, lbHealthCheckId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateLbHealthCheck(ctx context.Context, lbHealthCheckId string, body *loadbalancer.LbHealthCheckSetRequest) (*loadbalancer.LbHealthCheckShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LBHealthCheckApiAPI.SetLbHealthCheck(ctx, lbHealthCheckId)
	req = req.LbHealthCheckSetRequest(*body)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteLbHealthCheck(ctx context.Context, lbHealthCheckId string) error {
	req := client.sdkClient.LoadbalancerV1LBHealthCheckApiAPI.DeleteLbHealthCheck(ctx, lbHealthCheckId)
	_, err := req.Execute()
	return err
}
