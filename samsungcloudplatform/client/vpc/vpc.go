package vpc

import (
	"context"
	"net/http"

	"fmt"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *vpc.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: vpc.NewAPIClient(config),
	}
}

//------------ VPC -------------------//

func (client *Client) GetVpcList(ctx context.Context, request VpcDataSource) (*vpc.VpcListResponse, error) {
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
		req = req.State(vpc.VpcState(request.State.ValueString()))
	}
	if !request.Cidr.IsNull() {
		req = req.Cidr(request.Cidr.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateVpc(ctx context.Context, request VpcResource) (*vpc.VpcShowResponse, error) {
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

func (client *Client) GetVpc(ctx context.Context, vpcId string) (*vpc.VpcShowResponse, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.ShowVpc(ctx, vpcId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVpc(ctx context.Context, vpcId string, request VpcResource) (*vpc.VpcShowResponse, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.SetVpc(ctx, vpcId)

	req = req.VpcSetRequest(vpc.VpcSetRequest{
		Description: *vpc.NewNullableString(request.Description.ValueStringPointer()),
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

func (client *Client) GetSubnetList(ctx context.Context, request SubnetDataSource) (*vpc.SubnetListResponse, error) {
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
		req = req.State(vpc.SubnetState(request.State.ValueString()))
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if len(request.Type) > 0 {
		req = req.Type_(vpc.Type{
			ArrayOfSubnetType: &request.Type,
		})
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateSubnet(ctx context.Context, request SubnetResource) (*vpc.SubnetShowResponse, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.CreateSubnet(ctx)
	description := request.Description.ValueString()
	descriptionNS := vpc.NullableString{}
	descriptionNS.Set(&description)

	tags := convertToTags(request.Tags.Elements())

	req = req.SubnetCreateRequest(vpc.SubnetCreateRequest{
		Name:            request.Name.ValueString(),
		VpcId:           request.VpcId.ValueString(),
		Type:            vpc.SubnetType(request.Type.ValueString()),
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

func (client *Client) GetSubnet(ctx context.Context, subnetId string) (*vpc.SubnetShowResponse, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.ShowSubnet(ctx, subnetId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateSubnet(ctx context.Context, vpcId string, request SubnetResource) (*vpc.SubnetShowResponse, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.SetSubnet(ctx, vpcId)
	description := request.Description.ValueString()

	req = req.SubnetSetRequest(vpc.SubnetSetRequest{
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

func (client *Client) GetPublicipList(ctx context.Context, request PublicipDataSource) (*vpc.PublicipListResponse, error) {
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
		req = req.Type_(vpc.PublicipType(request.Type.ValueString()))
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreatePublicip(ctx context.Context, request PublicipResource) (*vpc.PublicipShowResponse, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.CreatePublicip(ctx)
	description := request.Description.ValueString()
	descriptionNS := vpc.NullableString{}
	descriptionNS.Set(&description)
	tags := convertToTags(request.Tags.Elements())

	req = req.PublicipCreateRequest(vpc.PublicipCreateRequest{
		Type:        vpc.PublicipType(request.Type.ValueString()),
		Description: descriptionNS,
		Tags:        tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPublicip(ctx context.Context, publicipId string) (*vpc.PublicipShowResponse, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.ShowPublicip(ctx, publicipId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeletePublicip(ctx context.Context, publicipId string) error {
	req := client.sdkClient.VpcV1PublicIpApiAPI.DeletePublicip(ctx, publicipId)

	_, err := req.Execute()
	return err
}

func (client *Client) UpdatePublicip(ctx context.Context, publicipId string, request PublicipResource) (*vpc.PublicipShowResponse, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.SetPublicip(ctx, publicipId)

	req = req.PublicipSetRequest(vpc.PublicipSetRequest{
		Description: request.Description.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

//------------ Port -------------------//

func (client *Client) GetPortList(ctx context.Context, request PortDataSource) (*vpc.PortListResponse, error) {
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

func (client *Client) CreatePort(ctx context.Context, request PortResource) (*vpc.PortShowResponse, error) {
	req := client.sdkClient.VpcV1PortsApiAPI.CreatePort(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.PortCreateRequest(vpc.PortCreateRequest{
		Name:           request.Name.ValueString(),
		SubnetId:       request.SubnetId.ValueString(),
		FixedIpAddress: *vpc.NewNullableString(request.FixedIpAddress.ValueStringPointer()),
		SecurityGroups: request.SecurityGroups,
		Description:    *vpc.NewNullableString(request.Description.ValueStringPointer()),
		Tags:           tags,
	})

	resp, _, err := req.Execute()

	return resp, err
}

func (client *Client) GetPort(ctx context.Context, portId string) (*vpc.PortShowResponse, error) {
	req := client.sdkClient.VpcV1PortsApiAPI.ShowPort(ctx, portId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePort(ctx context.Context, vpcId string, request PortResource) (*vpc.PortShowResponse, error) {
	req := client.sdkClient.VpcV1PortsApiAPI.SetPort(ctx, vpcId)

	req = req.PortSetRequest(vpc.PortSetRequest{
		Description:    *vpc.NewNullableString(request.Description.ValueStringPointer()),
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

func (client *Client) GetNatGatewayList(ctx context.Context, request NatGatewayDataSource) (*vpc.NatGatewayListResponse, error) {
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
		req = req.State(vpc.NatGatewayState(request.State.ValueString()))
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateNatGateway(ctx context.Context, request NatGatewayResource) (*vpc.NatGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.CreateNatGateway(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.NatGatewayCreateRequest(vpc.NatGatewayCreateRequest{
		SubnetId:    request.SubnetId.ValueString(),
		PublicipId:  request.PublicipId.ValueString(),
		Description: request.Description.ValueStringPointer(),
		Tags:        tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetNatGateway(ctx context.Context, natGatewayId string) (*vpc.NatGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.ShowNatGateway(ctx, natGatewayId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateNatGateway(ctx context.Context, natGatewayId string, request NatGatewayResource) (*vpc.NatGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1NatGatewayApiAPI.SetNatGateway(ctx, natGatewayId)

	req = req.NatGatewaySetRequest(vpc.NatGatewaySetRequest{
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

func (client *Client) GetInternetGatewayList(ctx context.Context, request InternetGatewayDataSource) (*vpc.InternetGatewayListResponse, error) {
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
		req = req.Type_(vpc.InternetGatewayType(request.Type.ValueString()))
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

func (client *Client) CreateInternetGateway(ctx context.Context, request InternetGatewayResource) (*vpc.InternetGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.CreateInternetGateway(ctx)
	description := request.Description.ValueString()
	descriptionNS := vpc.NullableString{}
	descriptionNS.Set(&description)

	tags := convertToTags(request.Tags.Elements())

	req = req.InternetGatewayCreateRequest(vpc.InternetGatewayCreateRequest{
		Type:             vpc.InternetGatewayType(request.Type.ValueString()),
		Description:      descriptionNS,
		FirewallEnabled:  request.FirewallEnabled.ValueBoolPointer(),
		FirewallLoggable: request.FirewallLoggable.ValueBoolPointer(),
		VpcId:            request.VpcId.ValueString(),
		Tags:             tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetInternetGateway(ctx context.Context, internetGatewayId string) (*vpc.InternetGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.ShowInternetGateway(ctx, internetGatewayId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateInternetGateway(ctx context.Context, internetGatewayId string, request InternetGatewayResource) (*vpc.InternetGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1InternetGatewaysApiAPI.SetInternetGateway(ctx, internetGatewayId)
	description := request.Description.ValueString()
	descriptionNS := vpc.NullableString{}
	descriptionNS.Set(&description)

	req = req.InternetGatewaySetRequest(vpc.InternetGatewaySetRequest{
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

//------------------- VPC Endpoint -------------------//

func (client *Client) GetVpcEndpointList(ctx context.Context, request VpcEndpointDataSource) (*vpc.VpcEndpointListResponse, error) {
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
		req = req.Id(request.Id.ValueString())
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
		req = req.ResourceType(vpc.VpcEndpointResourceType(request.ResourceType.ValueString()))
	}
	if !request.ResourceKey.IsNull() {
		req = req.ResourceKey(request.ResourceKey.ValueString())
	}
	if !request.EndpointIpAddress.IsNull() {
		req = req.EndpointIpAddress(request.EndpointIpAddress.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(vpc.VpcEndpointState(request.State.ValueString()))
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateVpcEndpoint(ctx context.Context, request VpcEndpointResource) (*vpc.VpcEndpointShowResponse, error) {
	req := client.sdkClient.VpcV1VpcEndpointApiAPI.CreateVpcEndpoint(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.VpcEndpointCreateRequest(vpc.VpcEndpointCreateRequest{
		Name:              request.Name.ValueString(),
		VpcId:             request.VpcId.ValueString(),
		SubnetId:          request.SubnetId.ValueString(),
		ResourceType:      vpc.VpcEndpointResourceType(request.ResourceType.ValueString()),
		ResourceKey:       request.ResourceKey.ValueString(),
		ResourceInfo:      request.ResourceInfo.ValueString(),
		EndpointIpAddress: request.EndpointIpAddress.ValueString(),
		Description:       request.Description.ValueStringPointer(),
		Tags:              tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVpcEndpoint(ctx context.Context, vpcEndpointId string) (*vpc.VpcEndpointShowResponse, error) {
	req := client.sdkClient.VpcV1VpcEndpointApiAPI.ShowVpcEndpoint(ctx, vpcEndpointId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVpcEndpoint(ctx context.Context, vpcEndpointId string, request VpcEndpointResource) (*vpc.VpcEndpointShowResponse, error) {
	req := client.sdkClient.VpcV1VpcEndpointApiAPI.SetVpcEndpoint(ctx, vpcEndpointId)

	req = req.VpcEndpointSetRequest(vpc.VpcEndpointSetRequest{
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

func (client *Client) GetPrivateNatList(ctx context.Context, request PrivateNatDataSource) (*vpc.PrivateNatListResponse, error) {
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
		req = req.State(vpc.PrivateNatState(request.State.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreatePrivateNat(ctx context.Context, request PrivateNatResource) (*vpc.PrivateNatShowResponse, error) {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.CreatePrivateNat(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.PrivateNatCreateRequest(vpc.PrivateNatCreateRequest{
		Name:            request.Name.ValueString(),
		DirectConnectId: request.DirectConnectId.ValueString(),
		Cidr:            request.Cidr.ValueString(),
		Description:     request.Description.ValueStringPointer(),
		Tags:            tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPrivateNat(ctx context.Context, privateNatId string) (*vpc.PrivateNatShowResponse, error) {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.ShowPrivateNat(ctx, privateNatId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePrivateNat(ctx context.Context, privateNatId string, request PrivateNatResource) (*vpc.PrivateNatShowResponse, error) {
	req := client.sdkClient.VpcV1PrivateNatApiAPI.SetPrivateNat(ctx, privateNatId)

	req = req.PrivateNatSetRequest(vpc.PrivateNatSetRequest{
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

func (client *Client) GetPrivateNatIpList(ctx context.Context, request PrivateNatIpDataSource) (*vpc.PrivateNatIpListResponse, error) {
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
		req = req.AttachedResourceType(vpc.PrivateNatIpAttachedResourceType(request.AttachedResourceType.ValueString()))
	}
	if !request.AttachedResourceId.IsNull() {
		req = req.AttachedResourceId(request.AttachedResourceId.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(vpc.PrivateNatIpState(request.State.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreatePrivateNatIp(ctx context.Context, request PrivateNatIpResource) (*vpc.PrivateNatIpShowResponse, error) {
	req := client.sdkClient.VpcV1PrivateNatIpApiAPI.CreatePrivateNatIp(ctx, request.PrivateNatId.ValueString())

	req = req.PrivateNatIpCreateRequest(vpc.PrivateNatIpCreateRequest{
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

//---------------------TGW

// Transit Gateway --------------------------->

func (client *Client) GetTransitGatewayInfo(ctx context.Context, transitGatewayId string) (*vpc.TransitGatewayShowResponse, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.ShowTransitGateway(ctx, transitGatewayId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetTransitGatewayList(ctx context.Context, request TgwDataSource) (*vpc.TransitGatewayListResponse, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.ListTransitGateways(ctx)
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
	if !request.State.IsNull() {
		req = req.State(vpc.TransitGatewayState(request.State.ValueString()))
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}

	resp, httpResp, err := req.Execute()
	return resp, httpResp, err
}

func (client *Client) CreateTransitGateway(ctx context.Context, dataReq TgwResource) (*vpc.TransitGatewayShowResponse, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.CreateTransitGateway(ctx)

	req = req.TransitGatewayCreateRequest(vpc.TransitGatewayCreateRequest{
		Description: *vpc.NewNullableString(dataReq.Description.ValueStringPointer()),
		Name:        dataReq.Name.ValueString(),
		Tags:        convertToTags(dataReq.Tags.Elements()),
	})

	resp, httpResp, err := req.Execute()
	return resp, httpResp, err

}

func (client *Client) UpdateTransitGateway(ctx context.Context, dataReq TgwResource) error {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.SetTransitGateway(ctx, dataReq.Id.ValueString())
	if !dataReq.Description.IsNull() {
		req = req.TransitGatewaySetRequest(vpc.TransitGatewaySetRequest{
			Description: dataReq.Description.ValueStringPointer(),
		})
	}

	_, _, err := req.Execute()
	return err
}

func (client *Client) DeleteTransitGateway(ctx context.Context, transitGatewayId string) error {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.DeleteTransitGateway(ctx, transitGatewayId)

	_, err := req.Execute()
	return err
}

// <-----  TGW Rule routing start ------>

func (client *Client) GetRoutingRuleList(ctx context.Context, request RoutingRuleDataSource) (*vpc.TransitGatewayRuleListResponse, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayRulesApiAPI.ListTransitGatewayRules(ctx, request.TransitGatewayId.ValueString())
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
	if !request.TgwConnectionVpcId.IsNull() {
		req = req.TgwConnectionVpcId(request.TgwConnectionVpcId.ValueString())
	}
	if !request.TgwConnectionVpcName.IsNull() {
		req = req.TgwConnectionVpcName(request.TgwConnectionVpcName.ValueString())
	}
	if !request.SourceType.IsNull() {
		req = req.SourceType(vpc.TransitGatewayRuleDestinationType(request.SourceType.ValueString()))
	}
	if !request.DestinationType.IsNull() {
		req = req.DestinationType(vpc.TransitGatewayRuleDestinationType(request.DestinationType.ValueString()))
	}
	if !request.DestinationCidr.IsNull() {
		req = req.DestinationCidr(request.DestinationCidr.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(vpc.RoutingRuleState(request.State.ValueString()))
	}
	resp, httpResp, err := req.Execute()
	return resp, httpResp, err
}

func (client *Client) CreateTgwRule(ctx context.Context, request RoutingRuleResource) (*vpc.TransitGatewayRuleShowResponse, error) {

	req := client.sdkClient.VpcV1TransitGatewayRulesApiAPI.CreateTransitGatewayRule(ctx, request.TransitGatewayId.ValueString())

	req = req.TransitGatewayRuleCreateRequest(vpc.TransitGatewayRuleCreateRequest{
		Description:        request.Description.ValueStringPointer(),
		DestinationCidr:    request.DestinationCidr.ValueString(),
		DestinationType:    vpc.TransitGatewayRuleDestinationType(request.DestinationType.ValueString()),
		TgwConnectionVpcId: request.TgwConnectionVpcId.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetRoutingRule(ctx context.Context, transitGatewayId string, routingRuleId string) (*vpc.TransitGatewayRuleListResponse, error) {
	req := client.sdkClient.VpcV1TransitGatewayRulesApiAPI.ListTransitGatewayRules(ctx, transitGatewayId)

	if routingRuleId != "" {
		req = req.Id(routingRuleId)
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteRoutingRule(ctx context.Context, transitGatewayId string, routingRuleId string) error {
	req := client.sdkClient.VpcV1TransitGatewayRulesApiAPI.DeleteTransitGatewayRule(ctx, transitGatewayId, routingRuleId)
	_, err := req.Execute()
	return err
}

// <-----  TGW Rule routing end ------>

// TGW-VPC  -------------------------->
func (client *Client) GetTgwVpcConnectionList(ctx context.Context, request TgwVpcConnectionDataSource) (*vpc.TransitGatewayVpcConnectionListResponse, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayVpcConnectionApiAPI.ListTransitGatewayVpcConnections(ctx, request.TransitGatewayId.ValueString())
	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(vpc.TransitGatewayVpcConnectionState(request.State.ValueString()))
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}

	resp, httpResp, err := req.Execute()
	return resp, httpResp, err
}

func (client *Client) CreateTransitGatewayVpcConnection(ctx context.Context, request TgwVpcConnectionResource) (*vpc.TransitGatewayVpcConnectionShowResponse, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayVpcConnectionApiAPI.CreateTransitGatewayVpcConnection(ctx, request.TransitGatewayId.ValueString())

	req = req.TransitGatewayVpcConnectionCreateRequest(vpc.TransitGatewayVpcConnectionCreateRequest{
		VpcId: request.VpcId.ValueString(),
	})

	resp, httpResp, err := req.Execute()
	return resp, httpResp, err

}

func (client *Client) GetTransitGatewayVpcConnection(ctx context.Context, transitGatewayId string, vpcConnectionId string) (*vpc.TransitGatewayVpcConnectionListResponse, error) {
	req := client.sdkClient.VpcV1TransitGatewayVpcConnectionApiAPI.ListTransitGatewayVpcConnections(ctx, transitGatewayId)
	if vpcConnectionId != "" {
		req = req.Id(vpcConnectionId)
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteTransitGatewayVpcConnection(ctx context.Context, transitGatewayId string, vpcConnectionId string) error {
	req := client.sdkClient.VpcV1TransitGatewayVpcConnectionApiAPI.DeleteTransitGatewayVpcConnection(ctx, transitGatewayId, vpcConnectionId)
	_, err := req.Execute()
	return err
}

// ------------------- VPC PEERING RULE -------------------//

// Get VPC Peering Rule list
func (client *Client) GetVpcPeeringRuleList(ctx context.Context, request VpcPeeringRuleDataSource) (*vpc.VpcPeeringRuleListResponse, error) {
	req := client.sdkClient.VpcV1VpcPeeringRuleApiAPI.ListVpcPeeringRules(ctx, request.VpcPeeringId.ValueString())
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
	if !request.SourceVpcId.IsNull() {
		req = req.SourceVpcId(request.SourceVpcId.ValueString())
	}
	if !request.SourceVpcType.IsNull() {
		req = req.SourceVpcType(vpc.VpcPeeringRuleDestinationVpcType(request.SourceVpcType.ValueString()))
	}
	if !request.DestinationVpcId.IsNull() {
		req = req.DestinationVpcId(request.DestinationVpcId.ValueString())
	}
	if !request.DestinationVpcType.IsNull() {
		req = req.DestinationVpcType(vpc.VpcPeeringRuleDestinationVpcType(request.DestinationVpcType.ValueString()))
	}
	if !request.DestinationCidr.IsNull() {
		req = req.DestinationCidr(request.DestinationCidr.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(vpc.VpcPeeringState(request.State.ValueString()))
	}
	resp, _, err := req.Execute()
	if err != nil {
		return nil, err
	}
	return resp, err
}

// CreateVpcPeeringRule creates a new VPC peering rule
func (client *Client) CreateVpcPeeringRule(ctx context.Context, request VpcPeeringRuleResource) (*vpc.VpcPeeringRuleShowResponse, error) {
	req := client.sdkClient.VpcV1VpcPeeringRuleApiAPI.CreateVpcPeeringRule(ctx, request.VpcPeeringId.ValueString())

	tags := convertToTags(request.Tags.Elements())

	req = req.VpcPeeringRuleCreateRequest(vpc.VpcPeeringRuleCreateRequest{
		DestinationCidr:    request.DestinationCidr.ValueString(),
		DestinationVpcType: vpc.VpcPeeringRuleDestinationVpcType(request.DestinationVpcType.ValueString()),
		Tags:               tags,
	})

	resp, _, err := req.Execute()
	if err != nil {
		return nil, err
	}
	return resp, err
}

// DeleteVpcPeeringRule deletes a VPC peering rule
func (client *Client) DeleteVpcPeeringRule(ctx context.Context, vpcPeeringId string, routingRuleId string) error {
	_, err := client.sdkClient.VpcV1VpcPeeringRuleApiAPI.DeleteVpcPeeringRule(ctx, vpcPeeringId, routingRuleId).Execute()
	if err != nil {
		return err
	}
	return nil
}

// GetVpcPeeringRule retrieves details of a specific VPC peering rule
func (client *Client) GetVpcPeeringRule(ctx context.Context, vpcPeeringId string, routingRuleId string) (vpc.VpcPeeringRule, int, error) {
	req := client.sdkClient.VpcV1VpcPeeringRuleApiAPI.ListVpcPeeringRules(ctx, vpcPeeringId)
	req = req.Id(routingRuleId)

	resp, status, err := req.Execute()
	if err != nil {
		return vpc.VpcPeeringRule{}, status.StatusCode, err
	}

	// Check if the response contains any items
	if len(resp.VpcPeeringRules) > 0 {
		// Manual check to make sure we get the right record (should be filtered in list api)
		for _, rule := range resp.VpcPeeringRules {
			if rule.Id == routingRuleId {
				return rule, status.StatusCode, nil
			}
		}
	}

	// list has response but not the requested ones => deleted
	return vpc.VpcPeeringRule{}, 404, fmt.Errorf("DELETED")
}
