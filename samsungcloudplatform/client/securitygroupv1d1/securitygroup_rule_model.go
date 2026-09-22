package securitygroupv1d1

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SecurityGroupRuleDataSource struct {
	Id                types.String `tfsdk:"id"`
	SecurityGroupRule types.Object `tfsdk:"security_group_rule"`
}

type SecurityGroupRuleDataSourceIds struct {
	Size               types.Int32         `tfsdk:"size"`
	Page               types.Int32         `tfsdk:"page"`
	Sort               types.String        `tfsdk:"sort"`
	Id                 types.String        `tfsdk:"id"`
	SecurityGroupId    types.String        `tfsdk:"security_group_id"`
	RemoteIpPrefix     types.String        `tfsdk:"remote_ip_prefix"`
	RemoteGroupId      types.String        `tfsdk:"remote_group_id"`
	Description        types.String        `tfsdk:"description"`
	Direction          types.String        `tfsdk:"direction"`
	Service            types.String        `tfsdk:"service"`
	Ids                []types.String      `tfsdk:"ids"`
	SecurityGroupRules []SecurityGroupRule `tfsdk:"security_group_rules"`
}

type SecurityGroupRule struct {
	Id                     types.String `tfsdk:"id"`
	SecurityGroupId        types.String `tfsdk:"security_group_id"`
	Ethertype              types.String `tfsdk:"ethertype"`
	Protocol               types.String `tfsdk:"protocol"`
	PortRangeMin           types.Int32  `tfsdk:"port_range_min"`
	PortRangeMax           types.Int32  `tfsdk:"port_range_max"`
	RemoteIpPrefix         types.String `tfsdk:"remote_ip_prefix"`
	RemoteGroupId          types.String `tfsdk:"remote_group_id"`
	RemoteGroupName        types.String `tfsdk:"remote_group_name"`
	RemoteAddressGroupId   types.String `tfsdk:"remote_address_group_id"`
	RemoteAddressGroupName types.String `tfsdk:"remote_address_group_name"`
	Description            types.String `tfsdk:"description"`
	Direction              types.String `tfsdk:"direction"`
	CreatedAt              types.String `tfsdk:"created_at"`
	CreatedBy              types.String `tfsdk:"created_by"`
	ModifiedAt             types.String `tfsdk:"modified_at"`
	ModifiedBy             types.String `tfsdk:"modified_by"`
}

type SecurityGroupRuleResource struct {
	Id                   types.String `tfsdk:"id"`
	SecurityGroupId      types.String `tfsdk:"security_group_id"`
	Ethertype            types.String `tfsdk:"ethertype"`
	Protocol             types.String `tfsdk:"protocol"`
	PortRangeMin         types.Int32  `tfsdk:"port_range_min"`
	PortRangeMax         types.Int32  `tfsdk:"port_range_max"`
	RemoteAddressGroupId types.String `tfsdk:"remote_address_group_id"`
	RemoteIpPrefix       types.String `tfsdk:"remote_ip_prefix"`
	RemoteGroupId        types.String `tfsdk:"remote_group_id"`
	Description          types.String `tfsdk:"description"`
	Direction            types.String `tfsdk:"direction"`
	SecurityGroupRule    types.Object `tfsdk:"security_group_rule"`
}

func (m SecurityGroupRule) AttributeTypes() map[string]attr.Type { // SecurityGroupRule 의 AttributeTypes 메서드를 추가한다.
	return map[string]attr.Type{
		"id":                        types.StringType,
		"security_group_id":         types.StringType,
		"ethertype":                 types.StringType,
		"protocol":                  types.StringType,
		"port_range_min":            types.Int32Type,
		"port_range_max":            types.Int32Type,
		"remote_ip_prefix":          types.StringType,
		"remote_group_id":           types.StringType,
		"remote_group_name":         types.StringType,
		"remote_address_group_id":   types.StringType,
		"remote_address_group_name": types.StringType,
		"description":               types.StringType,
		"direction":                 types.StringType,
		"created_at":                types.StringType,
		"created_by":                types.StringType,
		"modified_at":               types.StringType,
		"modified_by":               types.StringType,
	}
}
