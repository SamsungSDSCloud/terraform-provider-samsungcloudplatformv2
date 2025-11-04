package baremetal

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpbaremetal1d1 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/baremetal/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math"
	"net/http"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpbaremetal1d1.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpbaremetal1d1.NewAPIClient(config),
	}
}

func (client *Client) GetBaremetalList(ip types.String, serverName types.String, state types.String, vpcId types.String) (*scpbaremetal1d1.BaremetalListResponseV1Dot1, error) {
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

func (client *Client) GetBaremetal(ctx context.Context, baremetalId string) (*scpbaremetal1d1.BaremetalShowResponseV1Dot1, *http.Response, error) {
	req := client.sdkClient.BaremetalV1BaremetalsAPIsAPI.ShowBaremetal(ctx, baremetalId)

	resp, httpResponse, err := req.Execute()
	return resp, httpResponse, err
}

func (client *Client) CreateBaremetal(ctx context.Context, request BaremetalResource) (*scpbaremetal1d1.AsyncResponse, error) {
	req := client.sdkClient.BaremetalV1BaremetalsAPIsAPI.CreateBaremetals(ctx)

	tags := make([]scpbaremetal1d1.TagRequest, 0)
	for k, v := range request.Tags.Elements() {
		tag := scpbaremetal1d1.TagRequest{}

		tag.Key = k

		if v != nil {
			value := scpbaremetal1d1.NullableString{}
			value.Set(v.(types.String).ValueStringPointer())
			tag.Value = value
		}
		tags = append(tags, tag)
	}

	initScript := scpbaremetal1d1.NullableString{}
	initScript.Set(request.InitScript.ValueStringPointer())

	placementGroupName := scpbaremetal1d1.NullableString{}
	placementGroupName.Set(request.PlacementGroupName.ValueStringPointer())

	var requestServerDetails []ServerDetails
	serverDetails := make([]scpbaremetal1d1.ServerDetailsRequest, 0)

	request.ServerDetails.ElementsAs(ctx, &requestServerDetails, false)

	for _, requestServerDetail := range requestServerDetails {
		localSubnetId := scpbaremetal1d1.NullableString{}
		localSubnetId.Set(requestServerDetail.BareMetalLocalSubnetId.ValueStringPointer())

		localSubnetIp := scpbaremetal1d1.NullableString{}
		localSubnetIp.Set(requestServerDetail.BareMetalLocalSubnetIpAddress.ValueStringPointer())

		ipAddress := scpbaremetal1d1.NullableString{}
		ipAddress.Set(requestServerDetail.IpAddress.ValueStringPointer())

		publicIpAddressId := scpbaremetal1d1.NullableString{}
		publicIpAddressId.Set(requestServerDetail.PublicIpAddressId.ValueStringPointer())

		serverDetail := scpbaremetal1d1.ServerDetailsRequest{
			BareMetalLocalSubnetId:        localSubnetId,
			BareMetalLocalSubnetIpAddress: localSubnetIp,
			BareMetalServerName:           requestServerDetail.BareMetalServerName.ValueString(),
			IpAddress:                     ipAddress,
			NatEnabled:                    requestServerDetail.NatEnabled.ValueBool(),
			PublicIpAddressId:             publicIpAddressId,
			ServerTypeId:                  requestServerDetail.ServerTypeId.ValueString(),
			UseHyperThreading:             requestServerDetail.UseHyperThreading.ValueBoolPointer(),
		}

		serverDetails = append(serverDetails, serverDetail)
	}

	req = req.BaremetalCreateRequest(scpbaremetal1d1.BaremetalCreateRequest{
		ImageId:            request.ImageId.ValueString(),
		InitScript:         initScript,
		LockEnabled:        request.LockEnabled.ValueBoolPointer(),
		OsUserId:           request.OsUserId.ValueString(),
		OsUserPassword:     request.OsUserPassword.ValueString(),
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
	req = req.BaremetalTerminateRequest(scpbaremetal1d1.BaremetalTerminateRequest{BaremetalIdList: baremetalIdList})

	_, _, err := req.Execute()
	return err
}

func (client *Client) GetImageList(ctx context.Context, regionId string) (*scpbaremetal1d1.ImageListResponse, error) {
	req := client.sdkClient.BaremetalV1BMImageAPIsAPI.ListBaremetalImages(ctx)

	req = req.RegionId(regionId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) StopBaremetals(ctx context.Context, baremetalIds []string) error {
	req := client.sdkClient.BaremetalV1BaremetalsAPIsAPI.StopBaremetals(ctx)

	baremetalIdList := make([]interface{}, 0)
	for _, baremetalId := range baremetalIds {
		baremetalIdList = append(baremetalIdList, baremetalId)
	}

	req = req.BaremetalsOperationRequest(scpbaremetal1d1.BaremetalsOperationRequest{BareMetalServerIds: baremetalIdList})

	_, _, err := req.Execute()
	return err
}

func (client *Client) StartBaremetals(ctx context.Context, baremetalIds []string) error {
	req := client.sdkClient.BaremetalV1BaremetalsAPIsAPI.StartBaremetals(ctx)

	baremetalIdList := make([]interface{}, 0)
	for _, baremetalId := range baremetalIds {
		baremetalIdList = append(baremetalIdList, baremetalId)
	}

	req = req.BaremetalsOperationRequest(scpbaremetal1d1.BaremetalsOperationRequest{BareMetalServerIds: baremetalIdList})

	_, _, err := req.Execute()
	return err
}
