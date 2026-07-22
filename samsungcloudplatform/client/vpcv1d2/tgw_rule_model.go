package vpcv1d2

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TransitGatewayRuleDataSource struct {
	TransitGatewayId     types.String  `tfsdk:"transit_gateway_id"`
	Size                 types.Int32   `tfsdk:"size"`
	Page                 types.Int32   `tfsdk:"page"`
	Sort                 types.String  `tfsdk:"sort"`
	Id                   types.String  `tfsdk:"id"`
	TgwConnectionVpcId   types.String  `tfsdk:"tgw_connection_vpc_id"`
	TgwConnectionVpcName types.String  `tfsdk:"tgw_connection_vpc_name"`
	SourceType           types.String  `tfsdk:"source_type"`
	DestinationType      types.String  `tfsdk:"destination_type"`
	DestinationCidr      types.String  `tfsdk:"destination_cidr"`
	State                types.String  `tfsdk:"state"`
	TotalCount           types.Int32   `tfsdk:"total_count"`
	RuleType             types.String  `tfsdk:"rule_type"`
	RoutingRules         []RoutingRule `tfsdk:"routing_rules"`
}

type RoutingRule struct {
	AccountId               types.String `tfsdk:"account_id"`
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	Description             types.String `tfsdk:"description"`
	DestinationCidr         types.String `tfsdk:"destination_cidr"`
	DestinationResourceId   types.String `tfsdk:"destination_resource_id"`
	DestinationResourceName types.String `tfsdk:"destination_resource_name"`
	DestinationType         types.String `tfsdk:"destination_type"`
	Id                      types.String `tfsdk:"id"`
	ModifiedAt              types.String `tfsdk:"modified_at"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	SourceResourceId        types.String `tfsdk:"source_resource_id"`
	SourceResourceName      types.String `tfsdk:"source_resource_name"`
	SourceType              types.String `tfsdk:"source_type"`
	State                   types.String `tfsdk:"state"`
	TgwConnectionVpcId      types.String `tfsdk:"tgw_connection_vpc_id"`
	TgwConnectionVpcName    types.String `tfsdk:"tgw_connection_vpc_name"`
}

type TransitGatewayRuleResource struct {
	// Input
	TransitGatewayId   types.String `tfsdk:"transit_gateway_id"`
	Description        types.String `tfsdk:"description"`
	DestinationCidr    types.String `tfsdk:"destination_cidr"`
	DestinationType    types.String `tfsdk:"destination_type"`
	TgwConnectionVpcId types.String `tfsdk:"tgw_connection_vpc_id"`

	// Output
	Id          types.String `tfsdk:"id"`
	RoutingRule types.Object `tfsdk:"routing_rule"`
}

type CreatedRoutingRule struct {
	AccountId               types.String `tfsdk:"account_id"`
	Description             types.String `tfsdk:"description"`
	DestinationCidr         types.String `tfsdk:"destination_cidr"`
	DestinationResourceId   types.String `tfsdk:"destination_resource_id"`
	DestinationResourceName types.String `tfsdk:"destination_resource_name"`
	DestinationType         types.String `tfsdk:"destination_type"`
	Id                      types.String `tfsdk:"id"`
	SourceResourceId        types.String `tfsdk:"source_resource_id"`
	SourceResourceName      types.String `tfsdk:"source_resource_name"`
	SourceType              types.String `tfsdk:"source_type"`
	State                   types.String `tfsdk:"state"`
	TgwConnectionVpcId      types.String `tfsdk:"tgw_connection_vpc_id"`
	TgwConnectionVpcName    types.String `tfsdk:"tgw_connection_vpc_name"`
}

func (m CreatedRoutingRule) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":                types.StringType,
		"description":               types.StringType,
		"destination_cidr":          types.StringType,
		"destination_resource_id":   types.StringType,
		"destination_resource_name": types.StringType,
		"destination_type":          types.StringType,
		"id":                        types.StringType,
		"source_resource_id":        types.StringType,
		"source_resource_name":      types.StringType,
		"source_type":               types.StringType,
		"state":                     types.StringType,
		"tgw_connection_vpc_id":     types.StringType,
		"tgw_connection_vpc_name":   types.StringType,
	}
}
