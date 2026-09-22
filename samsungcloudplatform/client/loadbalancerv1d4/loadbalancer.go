package loadbalancerv1d4

import (
	"context"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	loadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *loadbalancer.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client { // client 생성 함수를 추가한다.
	return &Client{
		Config:    config,
		sdkClient: loadbalancer.NewAPIClient(config),
	}
}

func convertToTags(elements map[string]attr.Value) []loadbalancer.Tag {
	var tags []loadbalancer.Tag
	for k, v := range elements {
		tagObject := loadbalancer.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}

func listToSlice(elements []attr.Value) []string {
	result := make([]string, 0, len(elements))
	for _, v := range elements {
		result = append(result, v.(types.String).ValueString())
	}
	return result
}

// ------------ Load Balancer -------------------//
func (client *Client) GetLoadbalancerList(ctx context.Context, request LoadbalancerDataSource) (*loadbalancer.LoadbalancerListResponseV1Dot4, error) {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.ListLoadbalancers(ctx)

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
	if !request.ServiceIp.IsNull() {
		req = req.ServiceIp(request.ServiceIp.ValueString())
	}
	if !request.SubnetId.IsNull() {
		req = req.SubnetId(request.SubnetId.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetLoadbalancer(ctx context.Context, loadbalancerId string) (*loadbalancer.LoadbalancerShowResponseV1Dot4, error) {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.ShowLoadbalancer(ctx, loadbalancerId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateLoadbalancer(ctx context.Context, request LoadbalancerResource) (*loadbalancer.LoadbalancerCreateResponse, error) {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.CreateLoadbalancer(ctx)

	loadbalancerCreateRequest := *request.LoadbalancerCreate

	loadbalancerElement := loadbalancer.LoadbalancerCreateRequestV1Dot4{
		Loadbalancer: loadbalancer.LoadbalancerCreateRequestDetailV1Dot4{
			Description:            *loadbalancer.NewNullableString(loadbalancerCreateRequest.Description.ValueStringPointer()),
			FirewallEnabled:        *loadbalancer.NewNullableBool(loadbalancerCreateRequest.FirewallEnabled.ValueBoolPointer()),
			FirewallLoggingEnabled: *loadbalancer.NewNullableBool(loadbalancerCreateRequest.FirewallLoggingEnabled.ValueBoolPointer()),
			HealthCheckIps:         listToSlice(loadbalancerCreateRequest.HealthCheckIps.Elements()),
			Zones:                  listToSlice(loadbalancerCreateRequest.Zones.Elements()),
			LayerType:              loadbalancerCreateRequest.LayerType.ValueString(),
			Name:                   loadbalancerCreateRequest.Name.ValueString(),
			ServiceIp:              *loadbalancer.NewNullableString(loadbalancerCreateRequest.ServiceIp.ValueStringPointer()),
			PublicipId:             *loadbalancer.NewNullableString(loadbalancerCreateRequest.PublicipId.ValueStringPointer()),
			SubnetId:               loadbalancerCreateRequest.SubnetId.ValueString(),
			VpcId:                  loadbalancerCreateRequest.VpcId.ValueString(),
			SourceNatIp:            *loadbalancer.NewNullableString(loadbalancerCreateRequest.SourceNatIp.ValueStringPointer()),
			Tags:                   convertToTags(loadbalancerCreateRequest.Tags.Elements()),
		},
	}

	req = req.LoadbalancerCreateRequestV1Dot4(loadbalancerElement)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateLoadbalancer(ctx context.Context, loadbalancerId string, request LoadbalancerResource) (*loadbalancer.LoadbalancerShowResponse, error) {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.SetLoadbalancer(ctx, loadbalancerId)

	loadbalancerCreateRequest := *request.LoadbalancerCreate

	updateReq := loadbalancer.LoadbalancerUpdateRequest{
		Loadbalancer: loadbalancer.LoadbalancerUpdateRequestDetail{
			Description: loadbalancerCreateRequest.Description.ValueString(),
		},
	}
	req = req.LoadbalancerUpdateRequest(updateReq)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteLoadbalancer(ctx context.Context, loadbalancerId string) error {
	req := client.sdkClient.LoadbalancerV1LoadbalancersApiAPI.DeleteLoadbalancer(ctx, loadbalancerId)

	_, err := req.Execute()
	return err
}
