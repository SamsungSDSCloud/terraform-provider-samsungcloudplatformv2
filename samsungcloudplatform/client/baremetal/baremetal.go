package baremetal

import (
	"context"
	"math"
	"net/http"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpbaremetal1d2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/baremetal/1.2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpbaremetal1d2.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpbaremetal1d2.NewAPIClient(config),
	}
}

func (client *Client) GetBaremetalList(ip types.String, serverName types.String, state types.String, vpcId types.String) (*scpbaremetal1d2.BaremetalListResponseV1Dot2, error) {
	ctx := context.Background()

	req := client.sdkClient.BaremetalV1BaremetalsAPIsAPI.ListBaremetals(ctx)
	req = req.Size(math.MaxInt32)

	if !ip.IsNull() {
		req = req.Ip(ip.ValueString())
	}
	if !serverName.IsNull() {
		req = req.ServerName(serverName.ValueString())
	}
	if !state.IsNull() {
		req = req.State(state.ValueString())
	}
	if !vpcId.IsNull() {
		req = req.VpcId(vpcId.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetBaremetal(ctx context.Context, baremetalId string) (*scpbaremetal1d2.BaremetalShowResponseV1Dot2, *http.Response, error) {
	req := client.sdkClient.BaremetalV1BaremetalsAPIsAPI.ShowBaremetal(ctx, baremetalId)

	resp, httpResponse, err := req.Execute()
	return resp, httpResponse, err
}

func (client *Client) CreateBaremetal(ctx context.Context, request BaremetalResource, draft BaremetalResource) (*scpbaremetal1d2.AsyncResponse, error) {
	req := client.sdkClient.BaremetalV1BaremetalsAPIsAPI.CreateBaremetals(ctx)

	tags := make([]scpbaremetal1d2.TagRequest, 0)
	for k, v := range request.Tags.Elements() {
		tag := scpbaremetal1d2.TagRequest{}

		tag.Key = k

		if v != nil {
			value := scpbaremetal1d2.NullableString{}
			value.Set(v.(types.String).ValueStringPointer())
			tag.Value = value
		}
		tags = append(tags, tag)
	}

	initScript := scpbaremetal1d2.NullableString{}
	initScript.Set(request.InitScript.ValueStringPointer())

	placementGroupName := scpbaremetal1d2.NullableString{}
	placementGroupName.Set(request.PlacementGroupName.ValueStringPointer())

	var requestServerDetails []ServerDetails
	serverDetails := make([]scpbaremetal1d2.ServerDetailsRequestV1Dot2, 0)

	request.ServerDetails.ElementsAs(ctx, &requestServerDetails, false)

	for _, requestServerDetail := range requestServerDetails {
		localSubnetId := scpbaremetal1d2.NullableString{}
		localSubnetId.Set(requestServerDetail.BareMetalLocalSubnetId.ValueStringPointer())

		localSubnetIp := scpbaremetal1d2.NullableString{}
		localSubnetIp.Set(requestServerDetail.BareMetalLocalSubnetIpAddress.ValueStringPointer())

		ipAddress := scpbaremetal1d2.NullableString{}
		ipAddress.Set(requestServerDetail.IpAddress.ValueStringPointer())

		publicIpAddressId := scpbaremetal1d2.NullableString{}
		publicIpAddressId.Set(requestServerDetail.PublicIpAddressId.ValueStringPointer())

		serverDetail := scpbaremetal1d2.ServerDetailsRequestV1Dot2{
			BareMetalLocalSubnetId:        localSubnetId,
			BareMetalLocalSubnetIpAddress: localSubnetIp,
			BareMetalServerName:           requestServerDetail.BareMetalServerName.ValueString(),
			IpAddress:                     ipAddress,
			NatEnabled:                    requestServerDetail.NatEnabled.ValueBool(),
			PublicIpAddressId:             publicIpAddressId,
			ServerTypeId:                  requestServerDetail.ServerTypeId.ValueString(),
			Zone:                          requestServerDetail.Zone.ValueString(),
		}

		serverDetails = append(serverDetails, serverDetail)
	}

	req = req.BaremetalCreateRequestV1Dot2(scpbaremetal1d2.BaremetalCreateRequestV1Dot2{
		ImageId:            request.ImageId.ValueString(),
		InitScript:         initScript,
		LockEnabled:        request.LockEnabled.ValueBoolPointer(),
		OsUserId:           request.OsUserId.ValueString(),
		OsUserPassword:     draft.OsUserPassword.ValueString(),
		PlacementGroupName: placementGroupName,
		RegionId:           request.RegionId.ValueString(),
		ServerDetails:      serverDetails,
		SubnetId:           request.SubnetId.ValueString(),
		Tags:               tags,
		UsePlacementGroup:  request.UsePlacementGroup.ValueBoolPointer(),
		VpcId:              request.VpcId.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteBaremetal(ctx context.Context, baremetalIds []string) error {
	req := client.sdkClient.BaremetalV1BaremetalsAPIsAPI.DeleteBaremetals(ctx)

	baremetalIdList := make([]interface{}, 0)
	for _, baremetalId := range baremetalIds {
		baremetalIdList = append(baremetalIdList, baremetalId)
	}
	req = req.BaremetalTerminateRequest(scpbaremetal1d2.BaremetalTerminateRequest{BaremetalIdList: baremetalIdList})

	_, _, err := req.Execute()
	return err
}

func (client *Client) GetImageList(ctx context.Context) (*scpbaremetal1d2.ImageListResponse, error) {
	req := client.sdkClient.BaremetalV1BMImageAPIsAPI.ListBaremetalImages(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) StopBaremetals(ctx context.Context, baremetalIds []string) error {
	req := client.sdkClient.BaremetalV1BaremetalsAPIsAPI.StopBaremetals(ctx)

	baremetalIdList := make([]interface{}, 0)
	for _, baremetalId := range baremetalIds {
		baremetalIdList = append(baremetalIdList, baremetalId)
	}

	req = req.BaremetalsOperationRequest(scpbaremetal1d2.BaremetalsOperationRequest{BareMetalServerIds: baremetalIdList})

	_, _, err := req.Execute()
	return err
}

func (client *Client) StartBaremetals(ctx context.Context, baremetalIds []string) error {
	req := client.sdkClient.BaremetalV1BaremetalsAPIsAPI.StartBaremetals(ctx)

	baremetalIdList := make([]interface{}, 0)
	for _, baremetalId := range baremetalIds {
		baremetalIdList = append(baremetalIdList, baremetalId)
	}

	req = req.BaremetalsOperationRequest(scpbaremetal1d2.BaremetalsOperationRequest{BareMetalServerIds: baremetalIdList})

	_, _, err := req.Execute()
	return err
}
