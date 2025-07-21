package vpc

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/vpc/1.0"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpvpc.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpvpc.NewAPIClient(config),
	}
}

//------------ VPC -------------------//

func (client *Client) GetVpcList(ctx context.Context, request VpcDataSource) (*scpvpc.VpcListResponse, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.ListVpcs(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
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
		req = req.State(scpvpc.VpcState(request.State.ValueString()))
	}
	if !request.Cidr.IsNull() {
		req = req.Cidr(request.Cidr.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateVpc(ctx context.Context, request VpcResource) (*scpvpc.VpcShowResponse, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.CreateVpc(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.VpcCreateRequest(scpvpc.VpcCreateRequest{
		Cidr:        request.Cidr.ValueString(),
		Description: *scpvpc.NewNullableString(request.Description.ValueStringPointer()),
		Name:        request.Name.ValueString(),
		Tags:        tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVpc(ctx context.Context, vpcId string) (*scpvpc.VpcShowResponse, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.ShowVpc(ctx, vpcId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVpc(ctx context.Context, vpcId string, request VpcResource) (*scpvpc.VpcShowResponse, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.SetVpc(ctx, vpcId)

	req = req.VpcSetRequest(scpvpc.VpcSetRequest{
		Description: *scpvpc.NewNullableString(request.Description.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteVpc(ctx context.Context, vpcId string) error {
	req := client.sdkClient.VpcV1VpcsApiAPI.DeleteVpc(ctx, vpcId)

	_, err := req.Execute()
	return err
}

//------------ Subnet -------------------//

func (client *Client) GetSubnetList(ctx context.Context, request SubnetDataSource) (*scpvpc.SubnetListResponse, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.ListSubnets(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
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
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(scpvpc.SubnetState(request.State.ValueString()))
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if len(request.Type) > 0 {
		req = req.Type_(scpvpc.Type{
			ArrayOfSubnetType: &request.Type,
		})
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateSubnet(ctx context.Context, request SubnetResource) (*scpvpc.SubnetShowResponse, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.CreateSubnet(ctx)
	description := request.Description.ValueString()
	descriptionNS := scpvpc.NullableString{}
	descriptionNS.Set(&description)

	tags := convertToTags(request.Tags.Elements())

	req = req.SubnetCreateRequest(scpvpc.SubnetCreateRequest{
		Name:            request.Name.ValueString(),
		VpcId:           request.VpcId.ValueString(),
		Type:            scpvpc.SubnetType(request.Type.ValueString()),
		Cidr:            request.Cidr.ValueString(),
		Description:     descriptionNS,
		AllocationPools: convertAllocationPoolsToInterface(request.AllocationPools),
		DnsNameservers:  request.DnsNameservers,
		HostRoutes:      convertHostRoutesToInterface(request.HostRoutes),
		Tags:            tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetSubnet(ctx context.Context, subnetId string) (*scpvpc.SubnetShowResponse, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.ShowSubnet(ctx, subnetId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateSubnet(ctx context.Context, vpcId string, request SubnetResource) (*scpvpc.SubnetShowResponse, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.SetSubnet(ctx, vpcId)
	description := request.Description.ValueString()

	req = req.SubnetSetRequest(scpvpc.SubnetSetRequest{
		Description: description,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteSubnet(ctx context.Context, subnetId string) error {
	req := client.sdkClient.VpcV1SubnetsApiAPI.DeleteSubnet(ctx, subnetId)

	_, err := req.Execute()
	return err
}

//------------------- Public IP -------------------//

func (client *Client) GetPublicipList(ctx context.Context, request PublicipDataSource) (*scpvpc.PublicipListResponse, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.ListPublicip(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
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
	if !request.Type.IsNull() {
		req = req.Type_(scpvpc.PublicipType(request.Type.ValueString()))
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreatePublicip(ctx context.Context, request PublicipResource) (*scpvpc.PublicipShowResponse, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.CreatePublicip(ctx)
	description := request.Description.ValueString()
	descriptionNS := scpvpc.NullableString{}
	descriptionNS.Set(&description)
	tags := convertToTags(request.Tags.Elements())

	req = req.PublicipCreateRequest(scpvpc.PublicipCreateRequest{
		Type:        scpvpc.PublicipType(request.Type.ValueString()),
		Description: descriptionNS,
		Tags:        tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPublicip(ctx context.Context, publicipId string) (*scpvpc.PublicipShowResponse, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.ShowPublicip(ctx, publicipId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeletePublicip(ctx context.Context, publicipId string) error {
	req := client.sdkClient.VpcV1PublicIpApiAPI.DeletePublicip(ctx, publicipId)

	_, err := req.Execute()
	return err
}

func (client *Client) UpdatePublicip(ctx context.Context, publicipId string, request PublicipResource) (*scpvpc.PublicipShowResponse, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.SetPublicip(ctx, publicipId)

	req = req.PublicipSetRequest(scpvpc.PublicipSetRequest{
		Description: request.Description.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

//------------ Port -------------------//

func (client *Client) GetPortList(ctx context.Context, request PortDataSource) (*scpvpc.PortListResponse, error) {
	req := client.sdkClient.VpcV1PortsApiAPI.ListPorts(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
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
	if !request.SubnetName.IsNull() {
		req = req.SubnetName(request.SubnetName.ValueString())
	}
	if !request.SubnetId.IsNull() {
		req = req.SubnetId(request.SubnetId.ValueString())
	}
	if !request.AttachedResourceId.IsNull() {
		req = req.AttachedResourceId(request.AttachedResourceId.ValueString())
	}
	if !request.FixedIpAddress.IsNull() {
		req = req.FixedIpAddress(request.FixedIpAddress.ValueString())
	}
	if !request.MacAddress.IsNull() {
		req = req.MacAddress(request.MacAddress.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(request.State.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreatePort(ctx context.Context, request PortResource) (*scpvpc.PortShowResponse, error) {
	req := client.sdkClient.VpcV1PortsApiAPI.CreatePort(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.PortCreateRequest(scpvpc.PortCreateRequest{
		Name:           request.Name.ValueString(),
		SubnetId:       request.SubnetId.ValueString(),
		FixedIpAddress: *scpvpc.NewNullableString(request.FixedIpAddress.ValueStringPointer()),
		SecurityGroups: request.SecurityGroups,
		Description:    *scpvpc.NewNullableString(request.Description.ValueStringPointer()),
		Tags:           tags,
	})

	resp, _, err := req.Execute()

	return resp, err
}

func (client *Client) GetPort(ctx context.Context, portId string) (*scpvpc.PortShowResponse, error) {
	req := client.sdkClient.VpcV1PortsApiAPI.ShowPort(ctx, portId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePort(ctx context.Context, vpcId string, request PortResource) (*scpvpc.PortShowResponse, error) {
	req := client.sdkClient.VpcV1PortsApiAPI.SetPort(ctx, vpcId)

	req = req.PortSetRequest(scpvpc.PortSetRequest{
		Description:    *scpvpc.NewNullableString(request.Description.ValueStringPointer()),
		SecurityGroups: request.SecurityGroups,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeletePort(ctx context.Context, portId string) error {
	req := client.sdkClient.VpcV1PortsApiAPI.DeletePort(ctx, portId)

	_, err := req.Execute()
	return err
}

//------------------- NAT Gateway -------------------//

func (client *Client) GetNatGatewayList(ctx context.Context, request NatGatewayDataSource) (*scpvpc.NatGatewayListResponse, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.ListNatGateways(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.NatGatewayIpAddress.IsNull() {
		req = req.NatGatewayIpAddress(request.NatGatewayIpAddress.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}
	if !request.SubnetId.IsNull() {
		req = req.SubnetId(request.SubnetId.ValueString())
	}
	if !request.SubnetName.IsNull() {
		req = req.SubnetName(request.SubnetName.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(scpvpc.NatGatewayState(request.State.ValueString()))
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateNatGateway(ctx context.Context, request NatGatewayResource) (*scpvpc.NatGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.CreateNatGateway(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.NatGatewayCreateRequest(scpvpc.NatGatewayCreateRequest{
		SubnetId:    request.SubnetId.ValueString(),
		PublicipId:  request.PublicipId.ValueString(),
		Description: request.Description.ValueStringPointer(),
		Tags:        tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetNatGateway(ctx context.Context, natGatewayId string) (*scpvpc.NatGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.ShowNatGateway(ctx, natGatewayId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateNatGateway(ctx context.Context, natGatewayId string, request NatGatewayResource) (*scpvpc.NatGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.SetNatGateway(ctx, natGatewayId)

	req = req.NatGatewaySetRequest(scpvpc.NatGatewaySetRequest{
		Description: request.Description.ValueStringPointer(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteNatGateway(ctx context.Context, natGatewayId string) error {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.DeleteNatGateway(ctx, natGatewayId)

	_, err := req.Execute()
	return err
}

//------------ Internet Gateway -------------------//

func (client *Client) GetInternetGatewayList(ctx context.Context, request InternetGatewayDataSource) (*scpvpc.InternetGatewayListResponse, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.ListInternetGateways(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
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
	if !request.Type.IsNull() {
		req = req.Type_(scpvpc.InternetGatewayType(request.Type.ValueString()))
	}
	if !request.State.IsNull() {
		req = req.State(request.State.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateInternetGateway(ctx context.Context, request InternetGatewayResource) (*scpvpc.InternetGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.CreateInternetGateway(ctx)
	description := request.Description.ValueString()
	descriptionNS := scpvpc.NullableString{}
	descriptionNS.Set(&description)

	tags := convertToTags(request.Tags.Elements())

	req = req.InternetGatewayCreateRequest(scpvpc.InternetGatewayCreateRequest{
		Type:             scpvpc.InternetGatewayType(request.Type.ValueString()),
		Description:      descriptionNS,
		FirewallEnabled:  request.FirewallEnabled.ValueBoolPointer(),
		FirewallLoggable: request.FirewallLoggable.ValueBoolPointer(),
		VpcId:            request.VpcId.ValueString(),
		Tags:             tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetInternetGateway(ctx context.Context, internetGatewayId string) (*scpvpc.InternetGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.ShowInternetGateway(ctx, internetGatewayId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateInternetGateway(ctx context.Context, internetGatewayId string, request InternetGatewayResource) (*scpvpc.InternetGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.SetInternetGateway(ctx, internetGatewayId)
	description := request.Description.ValueString()
	descriptionNS := scpvpc.NullableString{}
	descriptionNS.Set(&description)

	req = req.InternetGatewaySetRequest(scpvpc.InternetGatewaySetRequest{
		Description: descriptionNS,
		Loggable:    request.Loggable.ValueBoolPointer(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteInternetGateway(ctx context.Context, internetGatewayId string) error {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.DeleteInternetGateway(ctx, internetGatewayId)

	_, err := req.Execute()
	return err
}

func convertToTags(elements map[string]attr.Value) []scpvpc.Tag {
	var tags []scpvpc.Tag
	for k, v := range elements {
		tagObject := scpvpc.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}

//------------------- VPC Endpoint -------------------//

func (client *Client) GetVpcEndpointList(ctx context.Context, request VpcEndpointDataSource) (*scpvpc.VpcEndpointListResponse, error) {
	req := client.sdkClient.VpcV1VpcEndpointApiAPI.ListVpcEndpoints(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.ResourceType.IsNull() {
		req = req.ResourceType(scpvpc.VpcEndpointResourceType(request.ResourceType.ValueString()))
	}
	if !request.ResourceKey.IsNull() {
		req = req.ResourceKey(request.ResourceKey.ValueString())
	}
	if !request.EndpointIpAddress.IsNull() {
		req = req.EndpointIpAddress(request.EndpointIpAddress.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(scpvpc.VpcEndpointState(request.State.ValueString()))
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateVpcEndpoint(ctx context.Context, request VpcEndpointResource) (*scpvpc.VpcEndpointShowResponse, error) {
	req := client.sdkClient.VpcV1VpcEndpointApiAPI.CreateVpcEndpoint(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.VpcEndpointCreateRequest(scpvpc.VpcEndpointCreateRequest{
		Name:              request.Name.ValueString(),
		VpcId:             request.VpcId.ValueString(),
		SubnetId:          request.SubnetId.ValueString(),
		ResourceType:      scpvpc.VpcEndpointResourceType(request.ResourceType.ValueString()),
		ResourceKey:       request.ResourceKey.ValueString(),
		ResourceInfo:      request.ResourceInfo.ValueString(),
		EndpointIpAddress: request.EndpointIpAddress.ValueString(),
		Description:       request.Description.ValueStringPointer(),
		Tags:              tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVpcEndpoint(ctx context.Context, vpcEndpointId string) (*scpvpc.VpcEndpointShowResponse, error) {
	req := client.sdkClient.VpcV1VpcEndpointApiAPI.ShowVpcEndpoint(ctx, vpcEndpointId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVpcEndpoint(ctx context.Context, vpcEndpointId string, request VpcEndpointResource) (*scpvpc.VpcEndpointShowResponse, error) {
	req := client.sdkClient.VpcV1VpcEndpointApiAPI.SetVpcEndpoint(ctx, vpcEndpointId)

	req = req.VpcEndpointSetRequest(scpvpc.VpcEndpointSetRequest{
		Description: request.Description.ValueStringPointer(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteVpcEndpoint(ctx context.Context, vpcEndpointId string) error {
	req := client.sdkClient.VpcV1VpcEndpointApiAPI.DeleteVpcEndpoint(ctx, vpcEndpointId)

	_, err := req.Execute()
	return err
}

//------------------- Private NAT -------------------//

func (client *Client) GetPrivateNatList(ctx context.Context, request PrivateNatDataSource) (*scpvpc.PrivateNatListResponse, error) {
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
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.DirectConnectName.IsNull() {
		req = req.DirectConnectName(request.DirectConnectName.ValueString())
	}
	if !request.DirectConnectId.IsNull() {
		req = req.DirectConnectId(request.DirectConnectId.ValueString())
	}
	if !request.Cidr.IsNull() {
		req = req.Cidr(request.Cidr.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(scpvpc.PrivateNatState(request.State.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreatePrivateNat(ctx context.Context, request PrivateNatResource) (*scpvpc.PrivateNatShowResponse, error) {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.CreatePrivateNat(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.PrivateNatCreateRequest(scpvpc.PrivateNatCreateRequest{
		Name:            request.Name.ValueString(),
		DirectConnectId: request.DirectConnectId.ValueString(),
		Cidr:            request.Cidr.ValueString(),
		Description:     request.Description.ValueStringPointer(),
		Tags:            tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPrivateNat(ctx context.Context, privateNatId string) (*scpvpc.PrivateNatShowResponse, error) {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.ShowPrivateNat(ctx, privateNatId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePrivateNat(ctx context.Context, privateNatId string, request PrivateNatResource) (*scpvpc.PrivateNatShowResponse, error) {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.SetPrivateNat(ctx, privateNatId)

	req = req.PrivateNatSetRequest(scpvpc.PrivateNatSetRequest{
		Description: request.Description.ValueStringPointer(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeletePrivateNat(ctx context.Context, privateNatId string) error {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.DeletePrivateNat(ctx, privateNatId)

	_, err := req.Execute()
	return err
}

//------------------- Private NAT IP -------------------//

func (client *Client) GetPrivateNatIpList(ctx context.Context, request PrivateNatIpDataSource) (*scpvpc.PrivateNatIpListResponse, error) {
	req := client.sdkClient.VpcV1PrivateNatIpApiAPI.ListPrivateNatIps(ctx, request.PrivateNatId.ValueString())
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
	if !request.AttachedResourceName.IsNull() {
		req = req.AttachedResourceName(request.AttachedResourceName.ValueString())
	}
	if !request.AttachedResourceType.IsNull() {
		req = req.AttachedResourceType(scpvpc.PrivateNatIpAttachedResourceType(request.AttachedResourceType.ValueString()))
	}
	if !request.AttachedResourceId.IsNull() {
		req = req.AttachedResourceId(request.AttachedResourceId.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(scpvpc.PrivateNatIpState(request.State.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreatePrivateNatIp(ctx context.Context, request PrivateNatIpResource) (*scpvpc.PrivateNatIpShowResponse, error) {
	req := client.sdkClient.VpcV1PrivateNatIpApiAPI.CreatePrivateNatIp(ctx, request.PrivateNatId.ValueString())

	req = req.PrivateNatIpCreateRequest(scpvpc.PrivateNatIpCreateRequest{
		IpAddress: request.IpAddress.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeletePrivateNatIp(ctx context.Context, privateNatId string, privateNatIpId string) error {
	req := client.sdkClient.VpcV1PrivateNatIpApiAPI.DeletePrivateNatIp(ctx, privateNatId, privateNatIpId)

	_, err := req.Execute()
	return err
}
