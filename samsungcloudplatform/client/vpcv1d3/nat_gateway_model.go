package vpcv1d3

import (
	"time"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NatGatewayResource struct {
	Id               types.String `tfsdk:"id"`
	SubnetId         types.String `tfsdk:"subnet_id"`
	PublicipIds      types.List   `tfsdk:"publicip_ids"`
	MultiZoneEnabled types.Bool   `tfsdk:"multi_zone_enabled"`
	Description      types.String `tfsdk:"description"`
	Tags             types.Map    `tfsdk:"tags"`
	NatGateway       types.Object `tfsdk:"nat_gateway"`
}

type NatGatewayValue struct {
	Id               types.String        `tfsdk:"id"`
	Name             types.String        `tfsdk:"name"`
	NatGatewayIps    []NatGatewayIpValue `tfsdk:"nat_gateway_ips"`
	VpcId            types.String        `tfsdk:"vpc_id"`
	VpcName          types.String        `tfsdk:"vpc_name"`
	SubnetId         types.String        `tfsdk:"subnet_id"`
	SubnetName       types.String        `tfsdk:"subnet_name"`
	SubnetCidr       types.String        `tfsdk:"subnet_cidr"`
	AccountId        types.String        `tfsdk:"account_id"`
	State            types.String        `tfsdk:"state"`
	MultiZoneEnabled types.Bool          `tfsdk:"multi_zone_enabled"`
	Description      types.String        `tfsdk:"description"`
	CreatedAt        types.String        `tfsdk:"created_at"`
	CreatedBy        types.String        `tfsdk:"created_by"`
	ModifiedAt       types.String        `tfsdk:"modified_at"`
	ModifiedBy       types.String        `tfsdk:"modified_by"`
}

func (v NatGatewayValue) AttributeTypes() map[string]attr.Type {
	ipAttrTypes := map[string]attr.Type{
		"ip_address":  types.StringType,
		"publicip_id": types.StringType,
	}

	return map[string]attr.Type{
		"id":                 types.StringType,
		"name":               types.StringType,
		"nat_gateway_ips":    types.ListType{ElemType: types.ObjectType{AttrTypes: ipAttrTypes}},
		"vpc_id":             types.StringType,
		"vpc_name":           types.StringType,
		"subnet_id":          types.StringType,
		"subnet_name":        types.StringType,
		"subnet_cidr":        types.StringType,
		"account_id":         types.StringType,
		"state":              types.StringType,
		"multi_zone_enabled": types.BoolType,
		"description":        types.StringType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
	}
}

func ResponseToNatGatewayValue(res vpc.NatGatewayV1Dot3) NatGatewayValue {
	natGatewayTf := NatGatewayValue{
		Id:               types.StringValue(res.Id),
		Name:             types.StringValue(res.Name),
		NatGatewayIps:    mapNatGatewayIps(res.NatGatewayIps),
		VpcId:            types.StringValue(res.VpcId),
		VpcName:          types.StringValue(res.VpcName),
		SubnetId:         types.StringValue(res.SubnetId),
		SubnetName:       types.StringValue(res.SubnetName),
		SubnetCidr:       types.StringValue(res.SubnetCidr),
		AccountId:        types.StringValue(res.AccountId),
		State:            types.StringValue(string(res.State)),
		MultiZoneEnabled: types.BoolValue(res.MultiZoneEnabled),
		Description:      types.StringPointerValue(res.Description.Get()),
		CreatedAt:        types.StringValue(res.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(res.CreatedBy),
		ModifiedAt:       types.StringValue(res.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(res.ModifiedBy),
	}

	return natGatewayTf
}

func mapNatGatewayIps(sdkIps []vpc.NatGatewayIp) []NatGatewayIpValue {
	ips := make([]NatGatewayIpValue, len(sdkIps))

	for pos, sdkIp := range sdkIps {
		ips[pos] = NatGatewayIpValue{
			IpAddress:  types.StringValue(sdkIp.IpAddress),
			PublicipId: types.StringValue(sdkIp.PublicipId),
		}
	}

	return ips
}

// ------------ NatGatewayIpValue ------------

type NatGatewayIpValue struct {
	IpAddress  types.String `tfsdk:"ip_address"`
	PublicipId types.String `tfsdk:"publicip_id"`
}

// ------------ NatGatewayDataSource (Data Source model) ------------

type NatGatewayDataSource struct {
	// Input
	Size                types.Int32  `tfsdk:"size"`
	Page                types.Int32  `tfsdk:"page"`
	Sort                types.String `tfsdk:"sort"`
	Name                types.String `tfsdk:"name"`
	NatGatewayIpAddress types.String `tfsdk:"nat_gateway_ip_address"`
	VpcId               types.String `tfsdk:"vpc_id"`
	VpcName             types.String `tfsdk:"vpc_name"`
	SubnetId            types.String `tfsdk:"subnet_id"`
	SubnetName          types.String `tfsdk:"subnet_name"`
	State               types.String `tfsdk:"state"`
	MultiZoneEnabled    types.Bool   `tfsdk:"multi_zone_enabled"`

	// Output
	NatGateways []NatGatewayDSValue `tfsdk:"nat_gateways"`
	TotalCount  types.Int32         `tfsdk:"total_count"`
	SortFinal   []types.String      `tfsdk:"sort_final"`
}

type NatGatewayDSValue struct {
	Id               types.String        `tfsdk:"id"`
	Name             types.String        `tfsdk:"name"`
	NatGatewayIps    []NatGatewayIpValue `tfsdk:"nat_gateway_ips"`
	VpcId            types.String        `tfsdk:"vpc_id"`
	VpcName          types.String        `tfsdk:"vpc_name"`
	SubnetId         types.String        `tfsdk:"subnet_id"`
	SubnetName       types.String        `tfsdk:"subnet_name"`
	SubnetCidr       types.String        `tfsdk:"subnet_cidr"`
	AccountId        types.String        `tfsdk:"account_id"`
	State            types.String        `tfsdk:"state"`
	MultiZoneEnabled types.Bool          `tfsdk:"multi_zone_enabled"`
	Description      types.String        `tfsdk:"description"`
	CreatedAt        types.String        `tfsdk:"created_at"`
	CreatedBy        types.String        `tfsdk:"created_by"`
	ModifiedAt       types.String        `tfsdk:"modified_at"`
	ModifiedBy       types.String        `tfsdk:"modified_by"`
}

func ResponseToNatGatewayDSValue(res vpc.NatGatewayV1Dot3) NatGatewayDSValue {
	return NatGatewayDSValue{
		Id:               types.StringValue(res.Id),
		Name:             types.StringValue(res.Name),
		NatGatewayIps:    mapNatGatewayIps(res.NatGatewayIps),
		VpcId:            types.StringValue(res.VpcId),
		VpcName:          types.StringValue(res.VpcName),
		SubnetId:         types.StringValue(res.SubnetId),
		SubnetName:       types.StringValue(res.SubnetName),
		SubnetCidr:       types.StringValue(res.SubnetCidr),
		AccountId:        types.StringValue(res.AccountId),
		State:            types.StringValue(string(res.State)),
		MultiZoneEnabled: types.BoolValue(res.MultiZoneEnabled),
		Description:      types.StringPointerValue(res.Description.Get()),
		CreatedAt:        types.StringValue(res.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(res.CreatedBy),
		ModifiedAt:       types.StringValue(res.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(res.ModifiedBy),
	}
}
