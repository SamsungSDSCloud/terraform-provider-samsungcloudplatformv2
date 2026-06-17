package vpcv1d2

import (
	"context"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpc/1.2"
)

func (client *Client) CreatePrivateNat(ctx context.Context, request PrivateNatResource) (*vpc.PrivateNatShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.CreatePrivateNat(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.PrivateNatCreateRequestV1Dot2(vpc.PrivateNatCreateRequestV1Dot2{
		Cidr:              request.Cidr.ValueString(),
		Description:       request.Description.ValueStringPointer(),
		Name:              request.Name.ValueString(),
		ServiceResourceId: request.ServiceResourceId.ValueString(),
		ServiceType:       vpc.PrivateNatServiceType(request.ServiceType.ValueString()),
		Tags:              tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePrivateNat(ctx context.Context, privateNatId string, request PrivateNatResource) (*vpc.PrivateNatShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.SetPrivateNat(ctx, privateNatId)

	req = req.PrivateNatSetRequest(vpc.PrivateNatSetRequest{
		Description: request.Description.ValueStringPointer(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPrivateNat(ctx context.Context, privateNatId string) (*vpc.PrivateNatShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.ShowPrivateNat(ctx, privateNatId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeletePrivateNat(ctx context.Context, privateNatId string) error {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.DeletePrivateNat(ctx, privateNatId)

	_, err := req.Execute()
	return err
}

func (client *Client) ListPrivateNats(ctx context.Context, request PrivateNatDataSources) (*vpc.PrivateNatListResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.ListPrivateNats(ctx)

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
	if !request.Cidr.IsNull() {
		req = req.Cidr(request.Cidr.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.ServiceResourceId.IsNull() {
		req = req.ServiceResourceId(request.ServiceResourceId.ValueString())
	}
	if !request.ServiceType.IsNull() {
		req = req.ServiceType(vpc.PrivateNatServiceType(request.ServiceType.ValueString()))
	}
	if !request.ServiceResourceName.IsNull() {
		req = req.ServiceResourceName(request.ServiceResourceName.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(vpc.PrivateNatState(request.State.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}
