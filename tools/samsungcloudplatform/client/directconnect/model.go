package directconnect

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-direct-connect"

//------------ Direct Connect -------------------//

type DirectConnectDataSource struct {
	Limit          types.Int32     `tfsdk:"limit"`
	Marker         types.String    `tfsdk:"marker"`
	Sort           types.String    `tfsdk:"sort"`
	Id             types.String    `tfsdk:"id"`
	Name           types.String    `tfsdk:"name"`
	State          types.String    `tfsdk:"state"`
	VpcId          types.String    `tfsdk:"vpc_id"`
	VpcName        types.String    `tfsdk:"vpc_name"`
	DirectConnects []DirectConnect `tfsdk:"direct_connects"`
}

type DirectConnectResource struct {
	Id               types.String `tfsdk:"id"`
	Bandwidth        types.Int32  `tfsdk:"bandwidth"`
	Description      types.String `tfsdk:"description"`
	FirewallEnabled  types.Bool   `tfsdk:"firewall_enabled"`
	FirewallLoggable types.Bool   `tfsdk:"firewall_loggable"`
	Name             types.String `tfsdk:"name"`
	VpcId            types.String `tfsdk:"vpc_id"`
	Tags             types.Map    `tfsdk:"tags"`
	DirectConnect    types.Object `tfsdk:"direct_connect"`
}

type DirectConnect struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	AccountId   types.String `tfsdk:"account_id"`
	Description types.String `tfsdk:"description"`
	VpcId       types.String `tfsdk:"vpc_id"`
	VpcName     types.String `tfsdk:"vpc_name"`
	Bandwidth   types.Int32  `tfsdk:"bandwidth"`
	FirewallId  types.String `tfsdk:"firewall_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	State       types.String `tfsdk:"state"`
}

func (m DirectConnect) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"account_id":  types.StringType,
		"description": types.StringType,
		"vpc_id":      types.StringType,
		"vpc_name":    types.StringType,
		"bandwidth":   types.Int32Type,
		"firewall_id": types.StringType,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"state":       types.StringType,
	}
}

// ------------ Routing Rule -------------------//

type RoutingRuleDataSource struct {
	Limit           types.Int32   `tfsdk:"limit"`
	Marker          types.String  `tfsdk:"marker"`
	Sort            types.String  `tfsdk:"sort"`
	DirectConnectId types.String  `tfsdk:"direct_connect_id"`
	Id              types.String  `tfsdk:"id"`
	DestinationType types.String  `tfsdk:"destination_type"`
	DestinationCidr types.String  `tfsdk:"destination_cidr"`
	State           types.String  `tfsdk:"state"`
	RoutingRules    []RoutingRule `tfsdk:"routing_rules"`
}

type RoutingRuleResource struct {
	Id                    types.String `tfsdk:"id"`
	DirectConnectId       types.String `tfsdk:"direct_connect_id"`
	DestinationType       types.String `tfsdk:"destination_type"`
	DestinationCidr       types.String `tfsdk:"destination_cidr"`
	DestinationResourceId types.String `tfsdk:"destination_resource_id"`
	Description           types.String `tfsdk:"description"`
	RoutingRule           types.Object `tfsdk:"routing_rule"`
}

type RoutingRule struct {
	Id                      types.String `tfsdk:"id"`
	AccountId               types.String `tfsdk:"account_id"`
	OwnerId                 types.String `tfsdk:"owner_id"`
	OwnerType               types.String `tfsdk:"owner_type"`
	DestinationType         types.String `tfsdk:"destination_type"`
	DestinationCidr         types.String `tfsdk:"destination_cidr"`
	DestinationResourceId   types.String `tfsdk:"destination_resource_id"`
	DestinationResourceName types.String `tfsdk:"destination_resource_name"`
	Description             types.String `tfsdk:"description"`
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	ModifiedAt              types.String `tfsdk:"modified_at"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	State                   types.String `tfsdk:"state"`
}

func (m RoutingRule) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                        types.StringType,
		"account_id":                types.StringType,
		"owner_id":                  types.StringType,
		"owner_type":                types.StringType,
		"destination_type":          types.StringType,
		"destination_cidr":          types.StringType,
		"destination_resource_id":   types.StringType,
		"destination_resource_name": types.StringType,
		"description":               types.StringType,
		"created_at":                types.StringType,
		"created_by":                types.StringType,
		"modified_at":               types.StringType,
		"modified_by":               types.StringType,
		"state":                     types.StringType,
	}
}
