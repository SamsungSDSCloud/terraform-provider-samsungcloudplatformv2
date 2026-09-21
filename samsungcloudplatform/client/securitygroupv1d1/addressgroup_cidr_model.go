package securitygroupv1d1

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Address Group CIDR models

type AddressGroupCidrResource struct {
	// Input
	AddressGroupId types.String `tfsdk:"address_group_id"`
	Addresses      types.Set    `tfsdk:"addresses"`
}

type AddressGroupCidrDataSources struct {
	// Input
	AddressGroupId types.String `tfsdk:"address_group_id"`
	Size           types.Int32  `tfsdk:"size"`
	Page           types.Int32  `tfsdk:"page"`
	Sort           types.String `tfsdk:"sort"`
	Address        types.String `tfsdk:"address"`

	// Output
	Addresses  types.Set   `tfsdk:"addresses"`
	TotalCount types.Int32 `tfsdk:"total_count"`
	Links      []Link      `tfsdk:"links"`
}
