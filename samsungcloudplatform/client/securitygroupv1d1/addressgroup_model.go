package securitygroupv1d1

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Address Group models

type AddressGroupResource struct {
	// Input
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Addresses   types.Set    `tfsdk:"addresses"`
	Tags        types.Map    `tfsdk:"tags"`

	//Output
	Id           types.String `tfsdk:"id"`
	AddressGroup types.Object `tfsdk:"address_group"`
}

type AddressGroup struct {
	AccountId    types.String `tfsdk:"account_id"`
	AddressCount types.Int32  `tfsdk:"address_count"`
	AddressLimit types.Int32  `tfsdk:"address_limit"`
	Addresses    types.Set    `tfsdk:"addresses"`
	CreatedAt    types.String `tfsdk:"created_at"`
	CreatedBy    types.String `tfsdk:"created_by"`
	Description  types.String `tfsdk:"description"`
	Id           types.String `tfsdk:"id"`
	ModifiedAt   types.String `tfsdk:"modified_at"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
	Name         types.String `tfsdk:"name"`
	State        types.String `tfsdk:"state"`
}

func (m AddressGroup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":    types.StringType,
		"address_count": types.Int32Type,
		"address_limit": types.Int32Type,
		"addresses":     types.SetType{ElemType: types.StringType},
		"created_at":    types.StringType,
		"created_by":    types.StringType,
		"description":   types.StringType,
		"id":            types.StringType,
		"modified_at":   types.StringType,
		"modified_by":   types.StringType,
		"name":          types.StringType,
		"state":         types.StringType,
	}
}

type AddressGroupDataSource struct {
	Id           types.String `tfsdk:"id"`
	AddressGroup types.Object `tfsdk:"address_group"`
}

type AddressGroupDataSources struct {
	// Input
	Size types.Int32  `tfsdk:"size"`
	Page types.Int32  `tfsdk:"page"`
	Sort types.String `tfsdk:"sort"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	// Output
	TotalCount    types.Int32    `tfsdk:"total_count"`
	Links         []Link         `tfsdk:"links"`
	AddressGroups []AddressGroup `tfsdk:"address_groups"`
}

type Link struct {
	Href types.String `tfsdk:"href"`
	Rel  types.String `tfsdk:"rel"`
}
