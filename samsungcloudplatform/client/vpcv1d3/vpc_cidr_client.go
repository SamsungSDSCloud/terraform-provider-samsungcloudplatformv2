package vpcv1d3

import (
	"context"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
)

//------------ VPC CIDRs -------------------//

func (client *Client) AddVpcCidr(ctx context.Context, request VpcCidrResource) (*vpc.VpcCidrCreateResponse, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.AddVpcCidr(ctx, request.VpcId.ValueString())

	req = req.VpcCidrCreateRequest(vpc.VpcCidrCreateRequest{
		Cidr: request.Cidr.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}
