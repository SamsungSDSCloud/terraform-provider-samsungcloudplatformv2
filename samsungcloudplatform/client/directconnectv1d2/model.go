package directconnect

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-direct-connect"

type DirectConnectResource struct {
	Id                types.String `tfsdk:"id"`
	Bandwidth         types.Int32  `tfsdk:"bandwidth"`
	Description       types.String `tfsdk:"description"`
	UplinkActiveZone  types.String `tfsdk:"uplink_active_zone"`
	UplinkStandbyZone types.String `tfsdk:"uplink_standby_zone"`
	FirewallEnabled   types.Bool   `tfsdk:"firewall_enabled"`
	FirewallLoggable  types.Bool   `tfsdk:"firewall_loggable"`
	Name              types.String `tfsdk:"name"`
	VpcId             types.String `tfsdk:"vpc_id"`
	Tags              types.Map    `tfsdk:"tags"`
	DirectConnect     types.Object `tfsdk:"direct_connect"`
}

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
	Id                types.String `tfsdk:"id"`          // Direct Connect Id
	CreatedAt         types.String `tfsdk:"created_at"`  // Created At
	CreatedBy         types.String `tfsdk:"created_by"`  // Created By
	ModifiedAt        types.String `tfsdk:"modified_at"` // Modified At
	ModifiedBy        types.String `tfsdk:"modified_by"` // Modified By
	State             types.String `tfsdk:"state"`       // State
	Name              types.String `tfsdk:"name"`        // Direct Connect Name
	AccountId         types.String `tfsdk:"account_id"`  // Account ID
	Description       types.String `tfsdk:"description"` // Description
	VpcId             types.String `tfsdk:"vpc_id"`      // VPC ID
	VpcName           types.String `tfsdk:"vpc_name"`    // VPC Name
	Bandwidth         types.Int32  `tfsdk:"bandwidth"`   // Bandwidth
	FirewallId        types.String `tfsdk:"firewall_id"` // Firewall ID
	UplinkActiveZone  types.String `tfsdk:"uplink_active_zone"`
	UplinkStandbyZone types.String `tfsdk:"uplink_standby_zone"`
}

func (m DirectConnect) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                  types.StringType,
		"created_at":          types.StringType,
		"created_by":          types.StringType,
		"modified_at":         types.StringType,
		"modified_by":         types.StringType,
		"state":               types.StringType,
		"name":                types.StringType,
		"account_id":          types.StringType,
		"description":         types.StringType,
		"vpc_id":              types.StringType,
		"vpc_name":            types.StringType,
		"bandwidth":           types.Int32Type,
		"firewall_id":         types.StringType,
		"uplink_active_zone":  types.StringType,
		"uplink_standby_zone": types.StringType,
	}
}
