package vpcv1d2

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VpcPeeringRuleDataSource struct {
	// Input
	VpcPeeringId       types.String `tfsdk:"vpc_peering_id"`
	Size               types.Int32  `tfsdk:"size"`
	Page               types.Int32  `tfsdk:"page"`
	Sort               types.String `tfsdk:"sort"`
	Id                 types.String `tfsdk:"id"`
	SourceVpcId        types.String `tfsdk:"source_vpc_id"`
	SourceVpcType      types.String `tfsdk:"source_vpc_type"`
	DestinationVpcId   types.String `tfsdk:"destination_vpc_id"`
	DestinationVpcType types.String `tfsdk:"destination_vpc_type"`
	DestinationCidr    types.String `tfsdk:"destination_cidr"`
	State              types.String `tfsdk:"state"`

	// Output
	VpcPeeringRules []VpcPeeringRule `tfsdk:"vpc_peering_rules"`
	TotalCount      types.Int32      `tfsdk:"total_count"`
	SortFinal       []types.String   `tfsdk:"sort_final"` // Sort output
}

type VpcPeeringRule struct {
	Id                 types.String `tfsdk:"id"`                   // VPC Peering Rule ID
	CreatedAt          types.String `tfsdk:"created_at"`           // Created At
	CreatedBy          types.String `tfsdk:"created_by"`           // Created By
	ModifiedAt         types.String `tfsdk:"modified_at"`          // Modified At
	ModifiedBy         types.String `tfsdk:"modified_by"`          // Modified By
	VpcPeeringId       types.String `tfsdk:"vpc_peering_id"`       // VPC Peering ID
	SourceVpcId        types.String `tfsdk:"source_vpc_id"`        // Source VPC ID
	SourceVpcName      types.String `tfsdk:"source_vpc_name"`      // Source VPC Name
	SourceVpcType      types.String `tfsdk:"source_vpc_type"`      // Source VPC Type
	DestinationVpcId   types.String `tfsdk:"destination_vpc_id"`   // Destination VPC ID
	DestinationVpcName types.String `tfsdk:"destination_vpc_name"` // Destination VPC Name
	DestinationVpcType types.String `tfsdk:"destination_vpc_type"` // Destination VPC Type
	DestinationCidr    types.String `tfsdk:"destination_cidr"`     // Destination CIDR
	State              types.String `tfsdk:"state"`                // State
}