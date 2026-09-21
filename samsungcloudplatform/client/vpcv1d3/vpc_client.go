package vpcv1d3

import (
	"context"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func convertToTags(elements map[string]attr.Value) []vpc.Tag {
	var tags []vpc.Tag
	for k, v := range elements {
		tagObject := vpc.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}

func (client *Client) CreateVpc(ctx context.Context, request VpcResource) (*vpc.VpcShowResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.CreateVpc(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.VpcCreateRequest(vpc.VpcCreateRequest{
		Cidr:        request.Cidr.ValueString(),
		Description: *vpc.NewNullableString(request.Description.ValueStringPointer()),
		Name:        request.Name.ValueString(),
		Tags:        tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVpc(ctx context.Context, vpcId string, request VpcResource) (*vpc.VpcSetResponse, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.SetVpc(ctx, vpcId)

	req = req.VpcSetRequest(vpc.VpcSetRequest{
		Description: *vpc.NewNullableString(request.Description.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVpc(ctx context.Context, vpcId string) (*vpc.VpcShowResponseV1Dot3, error) {
	resp, _, err := client.sdkClient.VpcV1VpcsApiAPI.ShowVpc(ctx, vpcId).Execute()
	return resp, err
}

func (client *Client) DeleteVpc(ctx context.Context, vpcId string) error {
	_, err := client.sdkClient.VpcV1VpcsApiAPI.DeleteVpc(ctx, vpcId).Execute()
	return err
}

func (client *Client) ListVpc(ctx context.Context, request VpcDataSource) (*vpc.VpcListResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.ListVpcs(ctx)
	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(vpc.VpcState(request.State.ValueString()))
	}
	if !request.Cidr.IsNull() {
		req = req.Cidr(request.Cidr.ValueString())
	}
	if !request.Zone.IsNull() {
		req = req.Zone(request.Zone.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}
