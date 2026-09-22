package vpcv1d3

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//------------ Public IP Resource -------------------//

type PublicipResource struct {
	Id          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Zone        types.String `tfsdk:"zone"`
	Description types.String `tfsdk:"description"`
	Tags        types.Map    `tfsdk:"tags"`
	Publicip    types.Object `tfsdk:"publicip"`
}

type Publicip struct {
	Id                   types.String `tfsdk:"id"`
	IpAddress            types.String `tfsdk:"ip_address"`
	AccountId            types.String `tfsdk:"account_id"`
	AttachedResourceType types.String `tfsdk:"attached_resource_type"`
	AttachedResourceName types.String `tfsdk:"attached_resource_name"`
	AttachedResourceId   types.String `tfsdk:"attached_resource_id"`
	Type                 types.String `tfsdk:"type"`
	State                types.String `tfsdk:"state"`
	Description          types.String `tfsdk:"description"`
	CreatedAt            types.String `tfsdk:"created_at"`
	CreatedBy            types.String `tfsdk:"created_by"`
	ModifiedAt           types.String `tfsdk:"modified_at"`
	ModifiedBy           types.String `tfsdk:"modified_by"`
	Zone                 types.String `tfsdk:"zone"`
}

func (m Publicip) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                     types.StringType,
		"ip_address":             types.StringType,
		"account_id":             types.StringType,
		"attached_resource_type": types.StringType,
		"attached_resource_name": types.StringType,
		"attached_resource_id":   types.StringType,
		"type":                   types.StringType,
		"state":                  types.StringType,
		"description":            types.StringType,
		"created_at":             types.StringType,
		"created_by":             types.StringType,
		"modified_at":            types.StringType,
		"modified_by":            types.StringType,
		"zone":                   types.StringType,
	}
}

//------------ Public IP Data Source -------------------//

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
	Zone                 types.List   `tfsdk:"zone"`

	// Output
	TotalCount types.Int32  `tfsdk:"total_count"`
	Publicips  []PublicipDS `tfsdk:"publicips"`
}

type PublicipDS struct {
	Id                   types.String `tfsdk:"id"`
	IpAddress            types.String `tfsdk:"ip_address"`
	AccountId            types.String `tfsdk:"account_id"`
	AttachedResourceType types.String `tfsdk:"attached_resource_type"`
	AttachedResourceName types.String `tfsdk:"attached_resource_name"`
	AttachedResourceId   types.String `tfsdk:"attached_resource_id"`
	Type                 types.String `tfsdk:"type"`
	State                types.String `tfsdk:"state"`
	Description          types.String `tfsdk:"description"`
	CreatedAt            types.String `tfsdk:"created_at"`
	CreatedBy            types.String `tfsdk:"created_by"`
	ModifiedAt           types.String `tfsdk:"modified_at"`
	ModifiedBy           types.String `tfsdk:"modified_by"`
	Zone                 types.String `tfsdk:"zone"`
}

func (m PublicipDS) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                     types.StringType,
		"ip_address":             types.StringType,
		"account_id":             types.StringType,
		"attached_resource_type": types.StringType,
		"attached_resource_name": types.StringType,
		"attached_resource_id":   types.StringType,
		"type":                   types.StringType,
		"state":                  types.StringType,
		"description":            types.StringType,
		"created_at":             types.StringType,
		"created_by":             types.StringType,
		"modified_at":            types.StringType,
		"modified_by":            types.StringType,
		"zone":                   types.StringType,
	}
}
