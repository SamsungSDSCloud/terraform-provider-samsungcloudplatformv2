package securitygroup

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-security-group"

// List Request
type SecurityGroupDataSource struct {
	Id            types.String `tfsdk:"id"`
	SecurityGroup types.Object `tfsdk:"security_group"`
}

type SecurityGroupDataSourceIds struct {
	Size types.Int32    `tfsdk:"size"`
	Page types.Int32    `tfsdk:"page"`
	Sort types.String   `tfsdk:"sort"`
	Id   types.String   `tfsdk:"id"`
	Name types.String   `tfsdk:"name"`
	Ids  []types.String `tfsdk:"ids"`
}

type SecurityGroup struct {
	Id          types.String `tfsdk:"id"`
	AccountId   types.String `tfsdk:"account_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	State       types.String `tfsdk:"state"`
	Loggable    types.Bool   `tfsdk:"loggable"`
	RuleCount   types.Int32  `tfsdk:"rule_count"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
}

type SecurityGroupResource struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Loggable      types.Bool   `tfsdk:"loggable"`
	Tags          types.Map    `tfsdk:"tags"` // tags field 추가
	SecurityGroup types.Object `tfsdk:"security_group"`
}

func (m SecurityGroup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"account_id":  types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"state":       types.StringType,
		"loggable":    types.BoolType,
		"rule_count":  types.Int32Type,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
	}
}

//------------------- Security Group Rule -------------------//

type SecurityGroupRuleDataSource struct {
	Id                types.String `tfsdk:"id"`
	SecurityGroupRule types.Object `tfsdk:"security_group_rule"`
}

type SecurityGroupRuleDataSourceIds struct {
	Size            types.Int32    `tfsdk:"size"`
	Page            types.Int32    `tfsdk:"page"`
	Sort            types.String   `tfsdk:"sort"`
	Id              types.String   `tfsdk:"id"`
	SecurityGroupId types.String   `tfsdk:"security_group_id"`
	RemoteIpPrefix  types.String   `tfsdk:"remote_ip_prefix"`
	RemoteGroupId   types.String   `tfsdk:"remote_group_id"`
	Description     types.String   `tfsdk:"description"`
	Direction       types.String   `tfsdk:"direction"`
	Service         types.String   `tfsdk:"service"`
	Ids             []types.String `tfsdk:"ids"`
}

type SecurityGroupRule struct {
	Id              types.String `tfsdk:"id"`
	SecurityGroupId types.String `tfsdk:"security_group_id"`
	Ethertype       types.String `tfsdk:"ethertype"`
	Protocol        types.String `tfsdk:"protocol"`
	PortRangeMin    types.Int32  `tfsdk:"port_range_min"`
	PortRangeMax    types.Int32  `tfsdk:"port_range_max"`
	RemoteIpPrefix  types.String `tfsdk:"remote_ip_prefix"`
	RemoteGroupId   types.String `tfsdk:"remote_group_id"`
	RemoteGroupName types.String `tfsdk:"remote_group_name"`
	Description     types.String `tfsdk:"description"`
	Direction       types.String `tfsdk:"direction"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

type SecurityGroupRuleResource struct {
	Id                types.String `tfsdk:"id"`
	SecurityGroupId   types.String `tfsdk:"security_group_id"`
	Ethertype         types.String `tfsdk:"ethertype"`
	Protocol          types.String `tfsdk:"protocol"`
	PortRangeMin      types.Int32  `tfsdk:"port_range_min"`
	PortRangeMax      types.Int32  `tfsdk:"port_range_max"`
	RemoteIpPrefix    types.String `tfsdk:"remote_ip_prefix"`
	RemoteGroupId     types.String `tfsdk:"remote_group_id"`
	Description       types.String `tfsdk:"description"`
	Direction         types.String `tfsdk:"direction"`
	SecurityGroupRule types.Object `tfsdk:"security_group_rule"`
}

func (m SecurityGroupRule) AttributeTypes() map[string]attr.Type { // SecurityGroupRule 의 AttributeTypes 메서드를 추가한다.
	return map[string]attr.Type{
		"id":                types.StringType,
		"security_group_id": types.StringType,
		"ethertype":         types.StringType,
		"protocol":          types.StringType,
		"port_range_min":    types.Int32Type,
		"port_range_max":    types.Int32Type,
		"remote_ip_prefix":  types.StringType,
		"remote_group_id":   types.StringType,
		"remote_group_name": types.StringType,
		"description":       types.StringType,
		"direction":         types.StringType,
		"created_at":        types.StringType,
		"created_by":        types.StringType,
		"modified_at":       types.StringType,
		"modified_by":       types.StringType,
	}
}
