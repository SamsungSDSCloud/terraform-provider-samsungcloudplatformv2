package vpcv1d2

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

type InternetGateway struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	VpcId       types.String `tfsdk:"vpc_id"`
	VpcName     types.String `tfsdk:"vpc_name"`
	AccountId   types.String `tfsdk:"account_id"`
	Type        types.String `tfsdk:"type"`
	State       types.String `tfsdk:"state"`
	Description types.String `tfsdk:"description"`
	Loggable    types.Bool   `tfsdk:"loggable"`
	FirewallId  types.String `tfsdk:"firewall_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
}
