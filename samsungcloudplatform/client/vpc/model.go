package vpc

import (
	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/vpc/1.0"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-vpc"

//------------ Vpc -------------------//

type VpcDataSource struct {
	Limit  types.Int32  `tfsdk:"limit"`
	Marker types.String `tfsdk:"marker"`
	Sort   types.String `tfsdk:"sort"`
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	State  types.String `tfsdk:"state"`
	Cidr   types.String `tfsdk:"cidr"`
	Vpcs   []Vpc        `tfsdk:"vpcs"`
}

type VpcResource struct {
	Id          types.String `tfsdk:"id"`
	Cidr        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
	Tags        types.Map    `tfsdk:"tags"`
	Vpc         types.Object `tfsdk:"vpc"`
}

type Vpc struct {
	Cidr        types.String `tfsdk:"cidr"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	Name        types.String `tfsdk:"name"`
	AccountId   types.String `tfsdk:"account_id"`
	State       types.String `tfsdk:"state"`
}

func (m Vpc) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cidr":        types.StringType,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"description": types.StringType,
		"id":          types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"name":        types.StringType,
		"account_id":  types.StringType,
		"state":       types.StringType,
	}
}

//------------ Subnet -------------------//

type SubnetDataSource struct {
	Limit   types.Int32      `tfsdk:"limit"`
	Marker  types.String     `tfsdk:"marker"`
	Sort    types.String     `tfsdk:"sort"`
	Id      types.String     `tfsdk:"id"`
	Name    types.String     `tfsdk:"name"`
	VpcName types.String     `tfsdk:"vpc_name"`
	State   types.String     `tfsdk:"state"`
	VpcId   types.String     `tfsdk:"vpc_id"`
	Type    []vpc.SubnetType `tfsdk:"type"`
	Subnets []Subnet         `tfsdk:"subnets"`
}

type SubnetResource struct {
	Id               types.String     `tfsdk:"id"`
	Name             types.String     `tfsdk:"name"`
	AccountId        types.String     `tfsdk:"account_id"`
	VpcId            types.String     `tfsdk:"vpc_id"`
	VpcName          types.String     `tfsdk:"vpc_name"`
	Type             types.String     `tfsdk:"type"`
	Cidr             types.String     `tfsdk:"cidr"`
	GatewayIpAddress types.String     `tfsdk:"gateway_ip_address"`
	AllocationPools  []AllocationPool `tfsdk:"allocation_pools"`
	DnsNameservers   []string         `tfsdk:"dns_nameservers"`
	HostRoutes       []HostRoute      `tfsdk:"host_routes"`
	State            types.String     `tfsdk:"state"`
	Description      types.String     `tfsdk:"description"`
	CreatedAt        types.String     `tfsdk:"created_at"`
	CreatedBy        types.String     `tfsdk:"created_by"`
	ModifiedAt       types.String     `tfsdk:"modified_at"`
	ModifiedBy       types.String     `tfsdk:"modified_by"`
	Tags             types.Map        `tfsdk:"tags"`
}

type Subnet struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	AccountId        types.String `tfsdk:"account_id"`
	VpcId            types.String `tfsdk:"vpc_id"`
	VpcName          types.String `tfsdk:"vpc_name"`
	Type             types.String `tfsdk:"type"`
	Cidr             types.String `tfsdk:"cidr"`
	GatewayIpAddress types.String `tfsdk:"gateway_ip_address"`
	State            types.String `tfsdk:"state"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	ModifiedAt       types.String `tfsdk:"modified_at"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
}

type AllocationPool struct {
	Start types.String `tfsdk:"start"`
	End   types.String `tfsdk:"end"`
}

type HostRoute struct {
	Destination types.String `tfsdk:"destination"`
	Nexthop     types.String `tfsdk:"nexthop"`
}

func convertAllocationPoolsToInterface(pools []AllocationPool) []interface{} {
	result := make([]interface{}, len(pools))
	for i, pool := range pools {
		result[i] = map[string]string{
			"start": pool.Start.ValueString(),
			"end":   pool.End.ValueString(),
		}
	}
	return result
}

func convertHostRoutesToInterface(routes []HostRoute) []interface{} {
	result := make([]interface{}, len(routes))
	for i, route := range routes {
		result[i] = map[string]string{
			"destination": route.Destination.ValueString(),
			"nexthop":     route.Nexthop.ValueString(),
		}
	}
	return result
}

//------------------- Public IP -------------------//

type PublicipResource struct {
	Id          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Tags        types.Map    `tfsdk:"tags"`
	Publicip    types.Object `tfsdk:"publicip"`
}

type PublicipDataSource struct {
	Limit                types.Int32  `tfsdk:"limit"`
	Marker               types.String `tfsdk:"marker"`
	Sort                 types.String `tfsdk:"sort"`
	IpAddress            types.String `tfsdk:"ip_address"`
	State                types.String `tfsdk:"state"`
	AttachedResourceType types.String `tfsdk:"attached_resource_type"`
	AttachedResourceName types.String `tfsdk:"attached_resource_name"`
	AttachedResourceId   types.String `tfsdk:"attached_resource_id"`
	VpcId                types.String `tfsdk:"vpc_id"`
	Type                 types.String `tfsdk:"type"`
	Publicips            []Publicip   `tfsdk:"publicips"`
}

type Publicip struct {
	Id                   types.String `tfsdk:"id"`
	IpAddress            types.String `tfsdk:"ip_address"`
	AccountId            types.String `tfsdk:"account_id"`
	AttachedResourceType types.String `tfsdk:"attached_resource_type"`
	AttachedResourceName types.String `tfsdk:"attached_resource_name"`
	AttachedResourceId   types.String `tfsdk:"attached_resource_id"`
	Type                 types.String `tfsdk:"type"`
	State                types.String `tfsdk:"state"`
	Description          types.String `tfsdk:"description"`
	CreatedAt            types.String `tfsdk:"created_at"`
	CreatedBy            types.String `tfsdk:"created_by"`
	ModifiedAt           types.String `tfsdk:"modified_at"`
	ModifiedBy           types.String `tfsdk:"modified_by"`
}

func (m Publicip) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                     types.StringType,
		"ip_address":             types.StringType,
		"account_id":             types.StringType,
		"attached_resource_type": types.StringType,
		"attached_resource_name": types.StringType,
		"attached_resource_id":   types.StringType,
		"type":                   types.StringType,
		"state":                  types.StringType,
		"description":            types.StringType,
		"created_at":             types.StringType,
		"created_by":             types.StringType,
		"modified_at":            types.StringType,
		"modified_by":            types.StringType,
	}
}

// ------------ Port -------------------//

type PortDataSource struct {
	Limit              types.Int32  `tfsdk:"limit"`
	Marker             types.String `tfsdk:"marker"`
	Sort               types.String `tfsdk:"sort"`
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	SubnetName         types.String `tfsdk:"subnet_name"`
	SubnetId           types.String `tfsdk:"subnet_id"`
	AttachedResourceId types.String `tfsdk:"attached_resource_id"`
	FixedIpAddress     types.String `tfsdk:"fixed_ip_address"`
	MacAddress         types.String `tfsdk:"mac_address"`
	State              types.String `tfsdk:"state"`
	Ports              []Port       `tfsdk:"ports"`
}

type PortResource struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	AccountId            types.String `tfsdk:"account_id"`
	SubnetId             types.String `tfsdk:"subnet_id"`
	SubnetName           types.String `tfsdk:"subnet_name"`
	VpcId                types.String `tfsdk:"vpc_id"`
	VpcName              types.String `tfsdk:"vpc_name"`
	FixedIpAddress       types.String `tfsdk:"fixed_ip_address"`
	MacAddress           types.String `tfsdk:"mac_address"`
	AttachedResourceId   types.String `tfsdk:"attached_resource_id"`
	AttachedResourceType types.String `tfsdk:"attached_resource_type"`
	SecurityGroups       []string     `tfsdk:"security_groups"`
	Description          types.String `tfsdk:"description"`
	State                types.String `tfsdk:"state"`
	CreatedAt            types.String `tfsdk:"created_at"`
	ModifiedAt           types.String `tfsdk:"modified_at"`
	Tags                 types.Map    `tfsdk:"tags"`
}

type Port struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	AccountId            types.String `tfsdk:"account_id"`
	SubnetId             types.String `tfsdk:"subnet_id"`
	SubnetName           types.String `tfsdk:"subnet_name"`
	VpcId                types.String `tfsdk:"vpc_id"`
	VpcName              types.String `tfsdk:"vpc_name"`
	FixedIpAddress       types.String `tfsdk:"fixed_ip_address"`
	MacAddress           types.String `tfsdk:"mac_address"`
	AttachedResourceId   types.String `tfsdk:"attached_resource_id"`
	AttachedResourceType types.String `tfsdk:"attached_resource_type"`
	Description          types.String `tfsdk:"description"`
	State                types.String `tfsdk:"state"`
	CreatedAt            types.String `tfsdk:"created_at"`
	ModifiedAt           types.String `tfsdk:"modified_at"`
}

//------------------- NAT Gateway -------------------//

type NatGatewayDataSource struct {
	Limit               types.Int32  `tfsdk:"limit"`
	Marker              types.String `tfsdk:"marker"`
	Sort                types.String `tfsdk:"sort"`
	Name                types.String `tfsdk:"name"`
	NatGatewayIpAddress types.String `tfsdk:"nat_gateway_ip_address"`
	VpcId               types.String `tfsdk:"vpc_id"`
	VpcName             types.String `tfsdk:"vpc_name"`
	SubnetId            types.String `tfsdk:"subnet_id"`
	SubnetName          types.String `tfsdk:"subnet_name"`
	State               types.String `tfsdk:"state"`
	NatGateways         []NatGateway `tfsdk:"nat_gateways"`
}

type NatGatewayResource struct {
	Id          types.String `tfsdk:"id"`
	SubnetId    types.String `tfsdk:"subnet_id"`
	PublicipId  types.String `tfsdk:"publicip_id"`
	Description types.String `tfsdk:"description"`
	Tags        types.Map    `tfsdk:"tags"`
	NatGateway  types.Object `tfsdk:"nat_gateway"`
}

type NatGateway struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	NatGatewayIpAddress types.String `tfsdk:"nat_gateway_ip_address"`
	VpcId               types.String `tfsdk:"vpc_id"`
	VpcName             types.String `tfsdk:"vpc_name"`
	SubnetId            types.String `tfsdk:"subnet_id"`
	SubnetName          types.String `tfsdk:"subnet_name"`
	SubnetCidr          types.String `tfsdk:"subnet_cidr"`
	AccountId           types.String `tfsdk:"account_id"`
	State               types.String `tfsdk:"state"`
	Description         types.String `tfsdk:"description"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
}

func (m NatGateway) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                     types.StringType,
		"name":                   types.StringType,
		"nat_gateway_ip_address": types.StringType,
		"vpc_id":                 types.StringType,
		"vpc_name":               types.StringType,
		"subnet_id":              types.StringType,
		"subnet_name":            types.StringType,
		"subnet_cidr":            types.StringType,
		"account_id":             types.StringType,
		"state":                  types.StringType,
		"description":            types.StringType,
		"created_at":             types.StringType,
		"created_by":             types.StringType,
		"modified_at":            types.StringType,
		"modified_by":            types.StringType,
	}
}

// ------------ Internet Gateway -------------------//

type InternetGatewayDataSource struct {
	Limit            types.Int32       `tfsdk:"limit"`
	Marker           types.String      `tfsdk:"marker"`
	Sort             types.String      `tfsdk:"sort"`
	Id               types.String      `tfsdk:"id"`
	Name             types.String      `tfsdk:"name"`
	Type             types.String      `tfsdk:"type"`
	State            types.String      `tfsdk:"state"`
	VpcId            types.String      `tfsdk:"vpc_id"`
	VpcName          types.String      `tfsdk:"vpc_name"`
	InternetGateways []InternetGateway `tfsdk:"internet_gateways"`
}

type InternetGatewayResource struct {
	Id               types.String `tfsdk:"id"`
	Type             types.String `tfsdk:"type"`
	Description      types.String `tfsdk:"description"`
	Loggable         types.Bool   `tfsdk:"loggable"`
	FirewallEnabled  types.Bool   `tfsdk:"firewall_enabled"`
	FirewallLoggable types.Bool   `tfsdk:"firewall_loggable"`
	VpcId            types.String `tfsdk:"vpc_id"`
	Tags             types.Map    `tfsdk:"tags"`
	InternetGateway  types.Object `tfsdk:"internet_gateway"`
}

type InternetGateway struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	AccountId   types.String `tfsdk:"account_id"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	VpcId       types.String `tfsdk:"vpc_id"`
	VpcName     types.String `tfsdk:"vpc_name"`
	Loggable    types.Bool   `tfsdk:"loggable"`
	FirewallId  types.String `tfsdk:"firewall_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	State       types.String `tfsdk:"state"`
}

func (m InternetGateway) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"account_id":  types.StringType,
		"type":        types.StringType,
		"description": types.StringType,
		"vpc_id":      types.StringType,
		"vpc_name":    types.StringType,
		"loggable":    types.BoolType,
		"firewall_id": types.StringType,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"state":       types.StringType,
	}
}

//------------------- VPC Endpoint -------------------//

type VpcEndpointDataSource struct {
	Limit             types.Int32   `tfsdk:"limit"`
	Marker            types.String  `tfsdk:"marker"`
	Sort              types.String  `tfsdk:"sort"`
	Id                types.String  `tfsdk:"id"`
	Name              types.String  `tfsdk:"name"`
	VpcName           types.String  `tfsdk:"vpc_name"`
	VpcId             types.String  `tfsdk:"vpc_id"`
	ResourceType      types.String  `tfsdk:"resource_type"`
	ResourceKey       types.String  `tfsdk:"resource_key"`
	EndpointIpAddress types.String  `tfsdk:"endpoint_ip_address"`
	State             types.String  `tfsdk:"state"`
	VpcEndpoints      []VpcEndpoint `tfsdk:"vpc_endpoints"`
}

type VpcEndpointResource struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	VpcId             types.String `tfsdk:"vpc_id"`
	SubnetId          types.String `tfsdk:"subnet_id"`
	ResourceType      types.String `tfsdk:"resource_type"`
	ResourceInfo      types.String `tfsdk:"resource_info"`
	ResourceKey       types.String `tfsdk:"resource_key"`
	EndpointIpAddress types.String `tfsdk:"endpoint_ip_address"`
	VpcEndpoint       types.Object `tfsdk:"vpc_endpoint"`
	Description       types.String `tfsdk:"description"`
	Tags              types.Map    `tfsdk:"tags"`
}

type VpcEndpoint struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	VpcId             types.String `tfsdk:"vpc_id"`
	VpcName           types.String `tfsdk:"vpc_name"`
	SubnetId          types.String `tfsdk:"subnet_id"`
	SubnetName        types.String `tfsdk:"subnet_name"`
	EndpointIpAddress types.String `tfsdk:"endpoint_ip_address"`
	ResourceType      types.String `tfsdk:"resource_type"`
	ResourceKey       types.String `tfsdk:"resource_key"`
	ResourceInfo      types.String `tfsdk:"resource_info"`
	AccountId         types.String `tfsdk:"account_id"`
	State             types.String `tfsdk:"state"`
	Description       types.String `tfsdk:"description"`
	CreatedAt         types.String `tfsdk:"created_at"`
	CreatedBy         types.String `tfsdk:"created_by"`
	ModifiedAt        types.String `tfsdk:"modified_at"`
	ModifiedBy        types.String `tfsdk:"modified_by"`
}

func (m VpcEndpoint) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                  types.StringType,
		"name":                types.StringType,
		"vpc_id":              types.StringType,
		"vpc_name":            types.StringType,
		"subnet_id":           types.StringType,
		"subnet_name":         types.StringType,
		"endpoint_ip_address": types.StringType,
		"resource_type":       types.StringType,
		"resource_key":        types.StringType,
		"resource_info":       types.StringType,
		"account_id":          types.StringType,
		"state":               types.StringType,
		"description":         types.StringType,
		"created_at":          types.StringType,
		"created_by":          types.StringType,
		"modified_at":         types.StringType,
		"modified_by":         types.StringType,
	}
}
