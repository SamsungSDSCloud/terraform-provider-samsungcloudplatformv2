package vpcv1d2

import (
	"context"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpc/1.2"
)

func (client *Client) CreatePort(ctx context.Context, request PortResource) (*vpc.PortShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1PortsApiAPI.CreatePort(ctx)

	sgIDs := make([]string, 0, len(request.SecurityGroups))
	for _, sg := range request.SecurityGroups {
		sgIDs = append(sgIDs, sg.Id.ValueString())
	}

	tags := convertToTags(request.Tags.Elements())

	req = req.PortCreateRequest(vpc.PortCreateRequest{
		Name:           request.Name.ValueString(),
		SubnetId:       request.SubnetId.ValueString(),
		FixedIpAddress: *vpc.NewNullableString(request.FixedIpAddress.ValueStringPointer()),
		SecurityGroups: sgIDs,
		Description:    *vpc.NewNullableString(request.Description.ValueStringPointer()),
		Tags:           tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPort(ctx context.Context, portId string) (*vpc.PortShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1PortsApiAPI.ShowPort(ctx, portId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePort(ctx context.Context, vpcId string, request PortResource) (*vpc.PortShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1PortsApiAPI.SetPort(ctx, vpcId)

	sgIDs := make([]string, 0, len(request.SecurityGroups))
	for _, sg := range request.SecurityGroups {
		sgIDs = append(sgIDs, sg.Id.ValueString())
	}

	req = req.PortSetRequest(vpc.PortSetRequest{
		Description:    *vpc.NewNullableString(request.Description.ValueStringPointer()),
		SecurityGroups: sgIDs,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeletePort(ctx context.Context, portId string) error {
	req := client.sdkClient.VpcV1PortsApiAPI.DeletePort(ctx, portId)

	_, err := req.Execute()
	return err
}
