package directconnect

import (
	"context"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	directconnect "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/direct-connect/1.2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *directconnect.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: directconnect.NewAPIClient(config),
	}
}

// ------------ Direct Connect List ------------

func (client *Client) CreateDirectConnect(ctx context.Context, request DirectConnectResource) (*directconnect.DirectConnectShowResponseV1Dot2, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsApiAPI.CreateDirectConnect(ctx)

	var tags []directconnect.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := directconnect.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}

	uplinkActiveZoneNS := directconnect.NullableString{}
	if !request.UplinkActiveZone.IsNull() && !request.UplinkActiveZone.IsUnknown() {
		v := request.UplinkActiveZone.ValueString()
		if v != "" {
			uplinkActiveZoneNS.Set(&v)
		}
	}
	uplinkStandbyZoneNS := directconnect.NullableString{}
	if !request.UplinkStandbyZone.IsNull() && !request.UplinkStandbyZone.IsUnknown() {
		v := request.UplinkStandbyZone.ValueString()
		if v != "" {
			uplinkStandbyZoneNS.Set(&v)
		}
	}

	req = req.DirectConnectCreateRequestV1Dot2(directconnect.DirectConnectCreateRequestV1Dot2{
		Bandwidth:         request.Bandwidth.ValueInt32(),
		Description:       request.Description.ValueStringPointer(),
		FirewallEnabled:   request.FirewallEnabled.ValueBoolPointer(),
		FirewallLoggable:  request.FirewallLoggable.ValueBoolPointer(),
		Name:              request.Name.ValueString(),
		VpcId:             request.VpcId.ValueString(),
		UplinkActiveZone:  uplinkActiveZoneNS,
		UplinkStandbyZone: uplinkStandbyZoneNS,
		Tags:              tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateDirectConnect(ctx context.Context, directConnectId string, request DirectConnectResource) (*directconnect.DirectConnectSetResponse, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsApiAPI.SetDirectConnect(ctx, directConnectId)

	req = req.DirectConnectSetRequest(directconnect.DirectConnectSetRequest{
		Description: request.Description.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetDirectConnect(ctx context.Context, directConnectId string) (*directconnect.DirectConnectShowResponseV1Dot2, error) {
	req := client.sdkClient.DirectConnectV1DirectConnectsApiAPI.ShowDirectConnect(ctx, directConnectId)

	resp, _, err := req.Execute()
	return resp, err
}
