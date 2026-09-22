package vpcv1d3

import (
	"context"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
)

//------------ Subnet -------------------//

func (client *Client) GetSubnetList(ctx context.Context, request SubnetDataSource) (*vpc.SubnetListResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.ListSubnets(ctx)

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
	if len(request.Type) > 0 {
		req = req.Type_(vpc.Type{
			ArrayOfSubnetType: &request.Type,
		})
	}
	if !request.State.IsNull() {
		req = req.State(vpc.SubnetState(request.State.ValueString()))
	}
	if !request.Category.IsNull() {
		req = req.Category(vpc.SubnetCategory(request.Category.ValueString()))
	}
	if !request.Cidr.IsNull() {
		req = req.Cidr(request.Cidr.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}
	if !request.PrimarySubnetId.IsNull() {
		req = req.PrimarySubnetId(request.PrimarySubnetId.ValueString())
	}
	if !request.Zone.IsNull() {
		req = req.Zone(request.Zone.ValueString())
	}

	resp, _, err := req.Execute()

	return resp, err

}

func (client *Client) CreateSubnet(ctx context.Context, request SubnetResource) (*vpc.SubnetShowResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.CreateSubnet(ctx)
	description := request.Description.ValueString()
	descriptionNS := vpc.NullableString{}
	descriptionNS.Set(&description)

	tags := convertToTags(request.Tags.Elements())

	category := vpc.SubnetCategory(request.Category.ValueString())
	primarySubnetId := vpc.NullableString{}
	if category == vpc.SUBNETCATEGORY_SECONDARY {
		pid := request.PrimarySubnetId.ValueString()
		primarySubnetId.Set(&pid)
	}

	req = req.SubnetCreateRequestV1Dot3(vpc.SubnetCreateRequestV1Dot3{
		AllocationPools:  convertAllocationPoolsToInterface(request.AllocationPools),
		Category:         category,
		Cidr:             request.Cidr.ValueString(),
		Description:      descriptionNS,
		DnsNameservers:   convertDnsNameserversToString(request.DnsNameservers),
		GatewayIpAddress: *vpc.NewNullableString(request.GatewayIpAddress.ValueStringPointer()),
		HostRoutes:       convertHostRoutesToInterface(request.HostRoutes),
		Name:             request.Name.ValueString(),
		PrimarySubnetId:  primarySubnetId,
		Tags:             tags,
		Type:             vpc.SubnetType(request.Type.ValueString()),
		VpcId:            request.VpcId.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetSubnet(ctx context.Context, subnetId string) (*vpc.SubnetShowResponseV1Dot3, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.ShowSubnet(ctx, subnetId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateSubnet(ctx context.Context, vpcId string, request SubnetResource, dhcpChanged bool) (*vpc.SubnetSetResponse, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.SetSubnet(ctx, vpcId)

	description := request.Description.ValueString()

	subnetSetReq := vpc.SubnetSetRequestV1Dot2{
		Description: *vpc.NewNullableString(&description),
	}
	if dhcpChanged {
		dhcpIpAddress := request.DhcpIpAddress.ValueString()
		subnetSetReq.DhcpIpAddress = *vpc.NewNullableString(&dhcpIpAddress)
	}

	req = req.SubnetSetRequestV1Dot2(subnetSetReq)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteSubnet(ctx context.Context, subnetId string) error {
	req := client.sdkClient.VpcV1SubnetsApiAPI.DeleteSubnet(ctx, subnetId)

	_, err := req.Execute()
	return err
}
