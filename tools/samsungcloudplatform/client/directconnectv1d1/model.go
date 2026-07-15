package directconnect

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-direct-connect"

// ------------ Direct Connect List ------------

type DirectConnectDataSource struct {
	// Input
	Size    types.Int32  `tfsdk:"size"`     // Size
	Page    types.Int32  `tfsdk:"page"`     // Page
	Sort    types.String `tfsdk:"sort"`     // Sort
	Id      types.String `tfsdk:"id"`       // Direct Connect ID
	Name    types.String `tfsdk:"name"`     // Direct Connect Name
	State   types.String `tfsdk:"state"`    // State
	VpcId   types.String `tfsdk:"vpc_id"`   // VPC Id
	VpcName types.String `tfsdk:"vpc_name"` // VPC Name

	// Output
	TotalCount     types.Int32     `tfsdk:"total_count"`     // Total Count
	SortFinal      []types.String  `tfsdk:"sort_final"`      // Sort output
	DirectConnects []DirectConnect `tfsdk:"direct_connects"` // Direct Connects
}

type DirectConnect struct {
	Id          types.String `tfsdk:"id"`          // Direct Connect Id
	CreatedAt   types.String `tfsdk:"created_at"`  // Created At
	CreatedBy   types.String `tfsdk:"created_by"`  // Created By
	ModifiedAt  types.String `tfsdk:"modified_at"` // Modified At
	ModifiedBy  types.String `tfsdk:"modified_by"` // Modified By
	State       types.String `tfsdk:"state"`       // State
	Name        types.String `tfsdk:"name"`        // Direct Connect Name
	AccountId   types.String `tfsdk:"account_id"`  // Account ID
	Description types.String `tfsdk:"description"` // Description
	VpcId       types.String `tfsdk:"vpc_id"`      // VPC ID
	VpcName     types.String `tfsdk:"vpc_name"`    // VPC Name
	Bandwidth   types.Int32  `tfsdk:"bandwidth"`   // Bandwidth
	FirewallId  types.String `tfsdk:"firewall_id"` // Firewall ID
}

func (m DirectConnect) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"state":       types.StringType,
		"name":        types.StringType,
		"account_id":  types.StringType,
		"description": types.StringType,
		"vpc_id":      types.StringType,
		"vpc_name":    types.StringType,
		"bandwidth":   types.Int32Type,
		"firewall_id": types.StringType,
	}
}

// ------------ Routing Rule List ------------

type RoutingRuleDataSource struct {
	// Input
	Size              types.Int32  `tfsdk:"size"`                // Size
	Page              types.Int32  `tfsdk:"page"`                // Page
	Sort              types.String `tfsdk:"sort"`                // Sort
	DirectConnectId   types.String `tfsdk:"direct_connect_id"`   // Direct Connect ID
	Id                types.String `tfsdk:"id"`                  // Routing Rule ID
	DestinationType   types.String `tfsdk:"destination_type"`    // Destination Type
	DestinationCidr   types.String `tfsdk:"destination_cidr"`    // Destination CIDR
	State             types.String `tfsdk:"state"`               // State

	// Output
	TotalCount   types.Int32       `tfsdk:"total_count"`   // Total Count
	SortFinal    []types.String    `tfsdk:"sort_final"`    // Sort output
	RoutingRules []RoutingRule     `tfsdk:"routing_rules"` // Routing Rules
}

type RoutingRule struct {
	Id                      types.String `tfsdk:"id"`                      // Routing Rule Id
	CreatedAt               types.String `tfsdk:"created_at"`              // Created At
	CreatedBy               types.String `tfsdk:"created_by"`              // Created By
	ModifiedAt              types.String `tfsdk:"modified_at"`             // Modified At
	ModifiedBy              types.String `tfsdk:"modified_by"`             // Modified By
	State                   types.String `tfsdk:"state"`                   // State
	AccountId               types.String `tfsdk:"account_id"`              // Account ID
	OwnerId                 types.String `tfsdk:"owner_id"`                // Owner ID
	OwnerType               types.String `tfsdk:"owner_type"`              // Owner Type
	DestinationType         types.String `tfsdk:"destination_type"`        // Destination Type
	DestinationCidr         types.String `tfsdk:"destination_cidr"`        // Destination CIDR
	DestinationResourceId   types.String `tfsdk:"destination_resource_id"` // Destination Resource ID
	DestinationResourceName types.String `tfsdk:"destination_resource_name"` // Destination Resource Name
	Description             types.String `tfsdk:"description"`             // Description
}

func (m RoutingRule) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                        types.StringType,
		"created_at":                types.StringType,
		"created_by":                types.StringType,
		"modified_at":               types.StringType,
		"modified_by":               types.StringType,
		"state":                     types.StringType,
		"account_id":                types.StringType,
		"owner_id":                  types.StringType,
		"owner_type":                types.StringType,
		"destination_type":          types.StringType,
		"destination_cidr":          types.StringType,
		"destination_resource_id":   types.StringType,
		"destination_resource_name": types.StringType,
		"description":               types.StringType,
	}
}
