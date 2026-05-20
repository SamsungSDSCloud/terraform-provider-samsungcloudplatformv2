package vpcv1d2

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-vpc"

// ------------ Subnet VIP List ------------

type SubnetVipDataSources struct {
	// Input
	SubnetId         types.String `tfsdk:"subnet_id"`
	Size             types.Int32  `tfsdk:"size"`
	Page             types.Int32  `tfsdk:"page"`
	Sort             types.String `tfsdk:"sort"`
	VirtualIpAddress types.String `tfsdk:"virtual_ip_address"`
	PublicIpAddress  types.String `tfsdk:"public_ip_address"`

	// Output
	TotalCount types.Int32        `tfsdk:"total_count"`
	SubnetVips []SubnetVipSummary `tfsdk:"subnet_vips"`
}

type SubnetVipSummary struct {
	Id                 types.String      `tfsdk:"id"`                   // Subnet Vip Id
	CreatedAt          types.String      `tfsdk:"created_at"`           // Created At
	CreatedBy          types.String      `tfsdk:"created_by"`           // Created By
	ModifiedAt         types.String      `tfsdk:"modified_at"`          // Modified At
	ModifiedBy         types.String      `tfsdk:"modified_by"`          // Modified By
	State              types.String      `tfsdk:"state"`                // State
	ConnectedPortCount types.Int32       `tfsdk:"connected_port_count"` // Connected Port Count
	StaticNat          *StaticNatSummary `tfsdk:"static_nat"`           // Static NAT Info (Nullable)
	VirtualIpAddress   types.String      `tfsdk:"virtual_ip_address"`   // Virtual IP Address
}

type StaticNatSummary struct {
	ExternalIpAddress types.String `tfsdk:"external_ip_address"` // Static Nat External Ip Address
	Id                types.String `tfsdk:"id"`                  // Static Nat Id
	PublicipId        types.String `tfsdk:"publicip_id"`         // Publicip ID
	State             types.String `tfsdk:"state"`               // Static Nat State
}

// ------------ Subnet VIP Detail ------------

type SubnetVipDataSource struct {
	// Input
	SubnetId types.String `tfsdk:"subnet_id"` // Subnet ID
	VipId    types.String `tfsdk:"vip_id"`    // Subnet Vip Id

	// Output
	SubnetVip *VpcSubnetVipDetail `tfsdk:"subnet_vip"`
}

type VpcSubnetVipDetail struct {
	Id               types.String        `tfsdk:"id"`                 // Subnet Vip Id
	CreatedAt        types.String        `tfsdk:"created_at"`         // Created At
	CreatedBy        types.String        `tfsdk:"created_by"`         // Created By
	ModifiedAt       types.String        `tfsdk:"modified_at"`        // Modified At
	ModifiedBy       types.String        `tfsdk:"modified_by"`        // Modified By
	State            types.String        `tfsdk:"state"`              // State
	SubnetId         types.String        `tfsdk:"subnet_id"`          // Subnet ID
	VipPortId        types.String        `tfsdk:"vip_port_id"`        // Vip Port Id
	VirtualIpAddress types.String        `tfsdk:"virtual_ip_address"` // Virtual IP Address
	Description      types.String        `tfsdk:"description"`        // Description
	ConnectedPorts   []ConnectedPortInfo `tfsdk:"connected_ports"`    // Connected Ports
	StaticNat        *StaticNatSummary   `tfsdk:"static_nat"`         // Static NAT Info
}

type ConnectedPortInfo struct {
	Id                   types.String `tfsdk:"id"`                     // Connected Port Id
	PortId               types.String `tfsdk:"port_id"`                // Port ID
	PortName             types.String `tfsdk:"port_name"`              // Port Name
	PortIpAddress        types.String `tfsdk:"port_ip_address"`        // Port IP Address
	AttachedResourceId   types.String `tfsdk:"attached_resource_id"`   // Connected resource ID
	AttachedResourceType types.String `tfsdk:"attached_resource_type"` // Connected resource Type
}

// ------------ Subnet VIP Create  ------------

type SubnetVipResource struct {
	// Input
	SubnetId         types.String `tfsdk:"subnet_id"`
	VirtualIpAddress types.String `tfsdk:"virtual_ip_address"`
	Description      types.String `tfsdk:"description"`

	// Output
	SubnetVip types.Object `tfsdk:"subnet_vip"`
}

func (m VpcSubnetVipDetail) AttributeTypes() map[string]attr.Type {
	connectedPortAttributeTypes := map[string]attr.Type{
		"id":                     types.StringType,
		"port_id":                types.StringType,
		"port_name":              types.StringType,
		"port_ip_address":        types.StringType,
		"attached_resource_id":   types.StringType,
		"attached_resource_type": types.StringType,
	}

	staticNatAttributeTypes := map[string]attr.Type{
		"external_ip_address": types.StringType,
		"id":                  types.StringType,
		"publicip_id":         types.StringType,
		"state":               types.StringType,
	}

	return map[string]attr.Type{
		"id":                 types.StringType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
		"state":              types.StringType,
		"subnet_id":          types.StringType,
		"vip_port_id":        types.StringType,
		"virtual_ip_address": types.StringType,
		"description":        types.StringType,
		"connected_ports":    types.ListType{ElemType: types.ObjectType{AttrTypes: connectedPortAttributeTypes}},
		"static_nat":         types.ObjectType{AttrTypes: staticNatAttributeTypes},
	}
}

// ------------ Subnet VIP NAT IP Create ------------

type SubnetVipNatIpResource struct {
	// Input
	SubnetId   types.String `tfsdk:"subnet_id"`
	VipId      types.String `tfsdk:"vip_id"`      // Virtual IP Address ID
	PublicipId types.String `tfsdk:"publicip_id"` // Public IP ID
	NatType    types.String `tfsdk:"nat_type"`

	// Output
	Id    types.String `tfsdk:"id"`    // Static Nat Id
	State types.String `tfsdk:"state"` // Static Nat State
}

// ------------ VPC CIDR Create ------------

type VpcCidrResource struct {
	// Input
	VpcId types.String `tfsdk:"vpc_id"`
	Cidr  types.String `tfsdk:"cidr"`

	// Output
	Vpc types.Object `tfsdk:"vpc"`
}

type VpcCidrDetail struct {
	AccountId   types.String  `tfsdk:"account_id"`  // Account ID
	CidrCount   types.Int32   `tfsdk:"cidr_count"`  // CIDR Count
	Cidrs       []VpcCidrInfo `tfsdk:"cidrs"`       // CIDRs
	CreatedAt   types.String  `tfsdk:"created_at"`  // Created At
	CreatedBy   types.String  `tfsdk:"created_by"`  // Created By
	Description types.String  `tfsdk:"description"` // Description
	Id          types.String  `tfsdk:"id"`          // VPC ID
	ModifiedAt  types.String  `tfsdk:"modified_at"` // Modified At
	ModifiedBy  types.String  `tfsdk:"modified_by"` // Modified By
	Name        types.String  `tfsdk:"name"`        // Name
	State       types.String  `tfsdk:"state"`       // State
}

type VpcCidrInfo struct {
	Cidr      types.String `tfsdk:"cidr"`       // CIDR
	CreatedAt types.String `tfsdk:"created_at"` // Created At
	CreatedBy types.String `tfsdk:"created_by"` // Created By
	Id        types.String `tfsdk:"id"`         // CIDR ID
}

func (m VpcCidrDetail) AttributeTypes() map[string]attr.Type {
	cidrInfoAttributeTypes := map[string]attr.Type{
		"cidr":       types.StringType,
		"created_at": types.StringType,
		"created_by": types.StringType,
		"id":         types.StringType,
	}

	return map[string]attr.Type{
		"account_id":  types.StringType,
		"cidr_count":  types.Int32Type,
		"cidrs":       types.ListType{ElemType: types.ObjectType{AttrTypes: cidrInfoAttributeTypes}},
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"description": types.StringType,
		"id":          types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"name":        types.StringType,
		"state":       types.StringType,
	}
}

// ------------ Subnet VIP Port Create ------------

type SubnetVipPortResource struct {
	// Input
	SubnetId types.String `tfsdk:"subnet_id"`
	VipId    types.String `tfsdk:"vip_id"`  // Virtual IP Address ID
	PortId   types.String `tfsdk:"port_id"` // Port ID

	// Output
	Id          types.String `tfsdk:"id"`            // Connected Port Id
	SubnetVipId types.String `tfsdk:"subnet_vip_id"` // Subnet Vip Id
}

// ------------ Subnet VPC Endpoint List ------------

type VpcEndpointDataSource struct {
	// Input
	Size              types.Int32  `tfsdk:"size"`
	Page              types.Int32  `tfsdk:"page"`
	Sort              types.String `tfsdk:"sort"`
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	VpcName           types.String `tfsdk:"vpc_name"`
	VpcId             types.String `tfsdk:"vpc_id"`
	SubnetId          types.String `tfsdk:"subnet_id"`
	ResourceType      types.String `tfsdk:"resource_type"`
	ResourceKey       types.String `tfsdk:"resource_key"`
	EndpointIpAddress types.String `tfsdk:"endpoint_ip_address"`
	State             types.String `tfsdk:"state"`

	// Output
	VpcEndpoints []VpcEndpoint `tfsdk:"vpc_endpoints"`
	TotalCount   types.Int32   `tfsdk:"total_count"`
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

// ------------ Transit gateway firewall ------------

type TransitGatewayFireWallResource struct {
	// Input
	TransitGatewayId types.String `tfsdk:"transit_gateway_id"`
	ProductType      types.String `tfsdk:"product_type"`

	// Output
	TransitGateway types.Object `tfsdk:"transit_gateway"`
}

type TransitGateway struct {
	AccountId               types.String `tfsdk:"account_id"`                // Account ID
	Bandwidth               types.Int32  `tfsdk:"bandwidth"`                 // Transit Gateway Port Bandwidth
	CreatedAt               types.String `tfsdk:"created_at"`                // Created At
	CreatedBy               types.String `tfsdk:"created_by"`                // Created By
	Description             types.String `tfsdk:"description"`               // Transit Gateway Description
	FirewallConnectionState types.String `tfsdk:"firewall_connection_state"` // Firewall Connection State
	FirewallIds             types.String `tfsdk:"firewall_ids"`              // Firewall ID
	FirewallId              types.String `tfsdk:"firewall_id"`               // Firewall ID
	Id                      types.String `tfsdk:"id"`                        // Transit Gateway ID
	ModifiedAt              types.String `tfsdk:"modified_at"`               // Modified At
	ModifiedBy              types.String `tfsdk:"modified_by"`               // Modified By
	Name                    types.String `tfsdk:"name"`                      // Transit Gateway Name
	State                   types.String `tfsdk:"state"`                     // State
	UplinkEnabled           types.Bool   `tfsdk:"uplink_enabled"`            // Uplink Enabled
}

func (m TransitGateway) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":                types.StringType,
		"bandwidth":                 types.Int32Type,
		"created_at":                types.StringType,
		"created_by":                types.StringType,
		"description":               types.StringType,
		"firewall_connection_state": types.StringType,
		"firewall_ids":              types.StringType,
		"firewall_id":               types.StringType,
		"id":                        types.StringType,
		"modified_at":               types.StringType,
		"modified_by":               types.StringType,
		"name":                      types.StringType,
		"state":                     types.StringType,
		"uplink_enabled":            types.BoolType,
	}
}

// ------------ Transit gateway firewall connection ------------

type TransitGatewayFirewallConnectionResource struct {
	// Input
	TransitGatewayId types.String `tfsdk:"transit_gateway_id"`

	// Output
	TransitGateway types.Object `tfsdk:"transit_gateway"`
}
