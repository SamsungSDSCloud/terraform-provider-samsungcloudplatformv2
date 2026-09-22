package vpcv1d3

import (
	"context"

	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
)

//------------ Public IP Resource Methods ------------

func (client *Client) CreatePublicip(ctx context.Context, request PublicipResource) (*scpvpc.PublicipShowResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.CreatePublicip(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.PublicipCreateRequestV1Dot3(scpvpc.PublicipCreateRequestV1Dot3{
		Type:        scpvpc.PublicipType(request.Type.ValueString()),
		Zone:        request.Zone.ValueString(),
		Description: *scpvpc.NewNullableString(request.Description.ValueStringPointer()),
		Tags:        tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPublicip(ctx context.Context, publicipId string) (*scpvpc.PublicipShowResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.ShowPublicip(ctx, publicipId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePublicip(ctx context.Context, publicipId string, request PublicipResource) (*scpvpc.PublicipSetResponse, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.SetPublicip(ctx, publicipId)

	req = req.PublicipSetRequest(scpvpc.PublicipSetRequest{
		Description: request.Description.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeletePublicip(ctx context.Context, publicipId string) error {
	req := client.sdkClient.VpcV1PublicIpApiAPI.DeletePublicip(ctx, publicipId)

	_, err := req.Execute()
	return err
}

func (client *Client) ListPublicips(ctx context.Context, request PublicipDataSource) (*scpvpc.PublicipListResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.ListPublicip(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.IpAddress.IsNull() {
		req = req.IpAddress(request.IpAddress.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(request.State.ValueString())
	}
	if !request.AttachedResourceType.IsNull() {
		req = req.AttachedResourceType(request.AttachedResourceType.ValueString())
	}
	if !request.AttachedResourceId.IsNull() {
		req = req.AttachedResourceId(request.AttachedResourceId.ValueString())
	}
	if !request.AttachedResourceName.IsNull() {
		req = req.AttachedResourceName(request.AttachedResourceName.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.Type.IsNull() {
		req = req.Type_(scpvpc.PublicipType(request.Type.ValueString()))
	}
	if !request.Zone.IsNull() && !request.Zone.IsUnknown() {
		var zones []string
		diags := request.Zone.ElementsAs(ctx, &zones, false)
		if !diags.HasError() && len(zones) > 0 {
			zone := scpvpc.Zone{ArrayOfString: &zones}
			req = req.Zone(zone)
		}
	}

	resp, _, err := req.Execute()
	return resp, err
}
