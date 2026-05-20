package vpcv1d2

import (
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
