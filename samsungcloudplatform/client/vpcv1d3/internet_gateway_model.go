package vpcv1d3

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	AccountId        types.String `tfsdk:"account_id"`
	Type             types.String `tfsdk:"type"`
	Description      types.String `tfsdk:"description"`
	VpcId            types.String `tfsdk:"vpc_id"`
	VpcName          types.String `tfsdk:"vpc_name"`
	Loggable         types.Bool   `tfsdk:"loggable"`
	FirewallId       types.String `tfsdk:"firewall_id"`
	MultiZoneEnabled types.Bool   `tfsdk:"multi_zone_enabled"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	ModifiedAt       types.String `tfsdk:"modified_at"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
	State            types.String `tfsdk:"state"`
}

func (m InternetGateway) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                 types.StringType,
		"name":               types.StringType,
		"account_id":         types.StringType,
		"type":               types.StringType,
		"description":        types.StringType,
		"vpc_id":             types.StringType,
		"vpc_name":           types.StringType,
		"loggable":           types.BoolType,
		"firewall_id":        types.StringType,
		"multi_zone_enabled": types.BoolType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
		"state":              types.StringType,
	}
}

type InternetGatewayDataSource struct {
	// Input
	Size    types.Int32  `tfsdk:"size"`
	Page    types.Int32  `tfsdk:"page"`
	Sort    types.String `tfsdk:"sort"`
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
	State   types.String `tfsdk:"state"`
	VpcId   types.String `tfsdk:"vpc_id"`
	VpcName types.String `tfsdk:"vpc_name"`

	// Output
	InternetGateways []InternetGateway `tfsdk:"internet_gateways"`
	TotalCount       types.Int32       `tfsdk:"total_count"`
	SortFinal        []types.String    `tfsdk:"sort_final"` // Sort output
}
