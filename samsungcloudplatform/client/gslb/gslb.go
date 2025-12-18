package gslb

import (
	"context"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	gslb "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/gslb/1.1" // terraform-sdk-samsungcloudplatformv2 에서 resourcemanager 라이브러리를 import 한다.
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *gslb.APIClient // 서비스의 client 를 구조체에 추가한다.
}

func NewClient(config *scpsdk.Configuration) *Client { // client 생성 함수를 추가한다.
	return &Client{
		Config:    config,
		sdkClient: gslb.NewAPIClient(config),
	}
}

func (client *Client) GetGslbList(ctx context.Context, request GslbDataSource) (*gslb.GslbListResponse, error) {
	req := client.sdkClient.GslbV1GslbsApiAPI.ListGslbs(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(request.State.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetGslb(ctx context.Context, gslbId string) (*gslb.GslbShowResponse, error) {
	req := client.sdkClient.GslbV1GslbsApiAPI.ShowGslb(ctx, gslbId)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) GetGslbResourceList(ctx context.Context, request GslbResourceDataSource) (*gslb.GslbResourceListResponseV1Dot1, error) {
	req := client.sdkClient.GslbV1GslbResourcesApiAPI.ListGslbResources(ctx, request.GslbId.ValueString())

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateGslb(ctx context.Context, request GslbResource) (*gslb.GslbShowResponse, error) {
	req := client.sdkClient.GslbV1GslbsApiAPI.CreateGslb(ctx) // 호출을 위한 구조체를 반환 받는다.

	var GslbTags []gslb.Tag

	for k, v := range request.Tags.Elements() {
		tagObject := gslb.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}

		GslbTags = append(GslbTags, tagObject)
	}

	var healthCheck *gslb.GslbHealthCheck

	gslbCreate := request.GslbCreate
	if gslbCreate.HealthCheck != nil {
		healthCheck = &gslb.GslbHealthCheck{
			HealthCheckInterval:     *gslb.NewNullableInt32(gslbCreate.HealthCheck.HealthCheckInterval.ValueInt32Pointer()),
			HealthCheckProbeTimeout: *gslb.NewNullableInt32(gslbCreate.HealthCheck.HealthCheckProbeTimeout.ValueInt32Pointer()),
			HealthCheckUserId:       *gslb.NewNullableString(gslbCreate.HealthCheck.HealthCheckUserId.ValueStringPointer()),
			HealthCheckUserPassword: *gslb.NewNullableString(gslbCreate.HealthCheck.HealthCheckUserPassword.ValueStringPointer()),
			Protocol:                gslbCreate.HealthCheck.Protocol.ValueString(),
			ReceiveString:           *gslb.NewNullableString(gslbCreate.HealthCheck.ReceiveString.ValueStringPointer()),
			SendString:              *gslb.NewNullableString(gslbCreate.HealthCheck.SendString.ValueStringPointer()),
			ServicePort:             *gslb.NewNullableInt32(gslbCreate.HealthCheck.ServicePort.ValueInt32Pointer()),
			Timeout:                 *gslb.NewNullableInt32(gslbCreate.HealthCheck.Timeout.ValueInt32Pointer()),
		}
	}

	gslbResources := make([]gslb.GslbResourceV1Dot1, len(gslbCreate.Resources))

	for i, gslbResource := range gslbCreate.Resources {
		gslbResources[i] = gslb.GslbResourceV1Dot1{
			Description: *gslb.NewNullableString(gslbResource.Description.ValueStringPointer()),
			Destination: gslbResource.Destination.ValueString(),
			Region:      gslbResource.Region.ValueString(),
			Weight:      *gslb.NewNullableInt32(gslbResource.Weight.ValueInt32Pointer()),
		}
	}

	gslbElement := gslb.GslbCreateRequestV1Dot1{
		Algorithm:   gslbCreate.Algorithm.ValueString(),
		Description: *gslb.NewNullableString(gslbCreate.Description.ValueStringPointer()),
		EnvUsage:    gslbCreate.EnvUsage.ValueString(),
		HealthCheck: *gslb.NewNullableGslbHealthCheck(healthCheck),
		Name:        gslbCreate.Name.ValueString(),
		Resources:   gslbResources,
		Tags:        GslbTags,
	}

	req = req.GslbCreateRequestV1Dot1(gslbElement)

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) UpdateGslb(ctx context.Context, gslbId string, request GslbResource) (*gslb.GslbShowResponse, error) {
	req := client.sdkClient.GslbV1GslbsApiAPI.SetGslb(ctx, gslbId) // 호출을 위한 구조체를 반환 받는다.

	gslbSet := request.GslbCreate

	req = req.GslbSetRequest(gslb.GslbSetRequest{
		Algorithm:   *gslb.NewNullableString(gslbSet.Algorithm.ValueStringPointer()),
		Description: *gslb.NewNullableString(gslbSet.Description.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteGslb(ctx context.Context, gslbId string) (*gslb.GslbShowResponse, error) {
	req := client.sdkClient.GslbV1GslbsApiAPI.DeleteGslb(ctx, gslbId)

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) UpdateGslbResource(ctx context.Context, gslbId string, request GslbResource) (*gslb.GslbShowResponse, error) {
	req := client.sdkClient.GslbV1GslbResourcesApiAPI.SetGslbResources(ctx, gslbId)

	gslbResources := request.GslbCreate.Resources
	var convertedGslbResources []gslb.GslbResourceSetRequestV1Dot1

	for _, gslbResource := range gslbResources {
		convertedGslbResources = append(convertedGslbResources, gslb.GslbResourceSetRequestV1Dot1{
			Description: *gslb.NewNullableString(gslbResource.Description.ValueStringPointer()),
			Destination: gslbResource.Destination.ValueString(),
			Region:      gslbResource.Region.ValueString(),
			Weight:      *gslb.NewNullableInt32(gslbResource.Weight.ValueInt32Pointer()),
		})
	}

	req = req.GslbResourcesSetRequest(gslb.GslbResourcesSetRequest{
		Resources: convertedGslbResources,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateGslbHealthCheck(ctx context.Context, gslbId string, request GslbResource) (*gslb.GslbShowResponse, error) {
	req := client.sdkClient.GslbV1GslbsApiAPI.SetGslbHealthCheck(ctx, gslbId)

	gslbHealthCheck := request.GslbCreate.HealthCheck
	req = req.GslbHealthCheck(gslb.GslbHealthCheck{
		HealthCheckInterval:     *gslb.NewNullableInt32(gslbHealthCheck.HealthCheckInterval.ValueInt32Pointer()),
		HealthCheckProbeTimeout: *gslb.NewNullableInt32(gslbHealthCheck.HealthCheckProbeTimeout.ValueInt32Pointer()),
		HealthCheckUserId:       *gslb.NewNullableString(gslbHealthCheck.HealthCheckUserId.ValueStringPointer()),
		HealthCheckUserPassword: *gslb.NewNullableString(gslbHealthCheck.HealthCheckUserPassword.ValueStringPointer()),
		Protocol:                gslbHealthCheck.Protocol.ValueString(),
		ReceiveString:           *gslb.NewNullableString(gslbHealthCheck.ReceiveString.ValueStringPointer()),
		SendString:              *gslb.NewNullableString(gslbHealthCheck.SendString.ValueStringPointer()),
		ServicePort:             *gslb.NewNullableInt32(gslbHealthCheck.ServicePort.ValueInt32Pointer()),
		Timeout:                 *gslb.NewNullableInt32(gslbHealthCheck.Timeout.ValueInt32Pointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

// #region Regional Routing Control
func (client *Client) GetGslbRegionalRoutingControlList(ctx context.Context, request GslbRegionalRoutingControlListDataSource) (*gslb.GslbRoutingControlListResponse, error) {
	req := client.sdkClient.GslbV1GslbsApiAPI.ListGslbsRegionalRoutingControl(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Region.IsNull() {
		req = req.Region(request.Region.ValueString())
	}
	if !request.Status.IsNull() {
		req = req.Status(request.Status.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateGslbRegionalRoutingControl(ctx context.Context, request GslbRegionalRoutingControlUpdateDataSource) (*gslb.GslbRoutingControlResponse, error) {
	req := client.sdkClient.GslbV1GslbsApiAPI.SetGslbRegionalRoutingControl(ctx, request.GslbId.ValueString())

	req = req.GslbRoutingControlRequest(gslb.GslbRoutingControlRequest{
		Region: request.Region.ValueString(),
		Status: request.Status.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

// #endregion
