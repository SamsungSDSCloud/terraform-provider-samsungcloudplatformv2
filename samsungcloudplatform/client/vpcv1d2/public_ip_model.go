package vpcv1d2

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PublicipDataSource struct {
	// Input
	Size                 types.Int32  `tfsdk:"size"`
	Page                 types.Int32  `tfsdk:"page"`
	Sort                 types.String `tfsdk:"sort"`
	IpAddress            types.String `tfsdk:"ip_address"`
	State                types.String `tfsdk:"state"`
	AttachedResourceType types.String `tfsdk:"attached_resource_type"`
	AttachedResourceId   types.String `tfsdk:"attached_resource_id"`
	AttachedResourceName types.String `tfsdk:"attached_resource_name"`
	VpcId                types.String `tfsdk:"vpc_id"`
	Type                 types.String `tfsdk:"type"`

	// Output
	TotalCount types.Int32 `tfsdk:"total_count"`
	Publicips  []PublicIp  `tfsdk:"publicips"`
}

type PublicIp struct {
	Id                   types.String `tfsdk:"id"`                     // PublicIP ID
	CreatedAt            types.String `tfsdk:"created_at"`             // Created At
	CreatedBy            types.String `tfsdk:"created_by"`             // Created By
	ModifiedAt           types.String `tfsdk:"modified_at"`            // Modified At
	ModifiedBy           types.String `tfsdk:"modified_by"`            // Modified By
	IpAddress            types.String `tfsdk:"ip_address"`             // IP Address
	State                types.String `tfsdk:"state"`                  // PublicIP State
	Type                 types.String `tfsdk:"type"`                   // PublicIP Type
	AccountId            types.String `tfsdk:"account_id"`             // Account ID
	Description          types.String `tfsdk:"description"`            // PublicIP Description
	AttachedResourceType types.String `tfsdk:"attached_resource_type"` // PublicIP Attached Resource Type
	AttachedResourceId   types.String `tfsdk:"attached_resource_id"`   // PublicIP Attached Resource ID
	AttachedResourceName types.String `tfsdk:"attached_resource_name"` // PublicIP Attached Resource Name
}