package firewall

import (
	scpfirewall "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/firewall/1.0"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-firewall"

//------------------- Firewall -------------------//

type FirewallDataSource struct {
	Id       types.String `tfsdk:"id"`
	Firewall types.Object `tfsdk:"firewall"`
}

type FirewallDataSourceIds struct {
	Page        types.Int32    `tfsdk:"page"`
	Size        types.Int32    `tfsdk:"size"`
	Sort        types.String   `tfsdk:"sort"`
	Name        types.String   `tfsdk:"name"`
	VpcName     types.String   `tfsdk:"vpc_name"`
	ProductType types.List     `tfsdk:"product_type"`
	State       types.List     `tfsdk:"state"`
	Ids         []types.String `tfsdk:"ids"`
}

type Firewall struct {
	Id              types.String `tfsdk:"id"`
	AccountId       types.String `tfsdk:"account_id"`
	Name            types.String `tfsdk:"name"`
	VpcId           types.String `tfsdk:"vpc_id"`
	VpcName         types.String `tfsdk:"vpc_name"`
	Loggable        types.Bool   `tfsdk:"loggable"`
	FwResourceId    types.String `tfsdk:"fw_resource_id"`
	PreProductId    types.String `tfsdk:"pre_product_id"`
	ProductType     types.String `tfsdk:"product_type"`
	State           types.String `tfsdk:"state"`
	Status          types.String `tfsdk:"status"`
	TotalRuleCount  types.Int32  `tfsdk:"total_rule_count"`
	FlavorName      types.String `tfsdk:"flavor_name"`
	FlavorRuleQuota types.Int32  `tfsdk:"flavor_rule_quota"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

func (m Firewall) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                types.StringType,
		"account_id":        types.StringType,
		"name":              types.StringType,
		"vpc_id":            types.StringType,
		"vpc_name":          types.StringType,
		"loggable":          types.BoolType,
		"fw_resource_id":    types.StringType,
		"pre_product_id":    types.StringType,
		"product_type":      types.StringType,
		"state":             types.StringType,
		"status":            types.StringType,
		"total_rule_count":  types.Int32Type,
		"flavor_name":       types.StringType,
		"flavor_rule_quota": types.Int32Type,
		"created_at":        types.StringType,
		"created_by":        types.StringType,
		"modified_at":       types.StringType,
		"modified_by":       types.StringType,
	}
}

//------------------- Firewall Rule -------------------//

type FirewallRuleDataSource struct {
	Id           types.String `tfsdk:"id"`
	FirewallId   types.String `tfsdk:"firewall_id"`
	FirewallRule types.Object `tfsdk:"firewall_rule"`
}

type FirewallRuleDataSourceIds struct {
	Page        types.Int32    `tfsdk:"page"`
	Size        types.Int32    `tfsdk:"size"`
	Sort        types.String   `tfsdk:"sort"`
	FirewallId  types.String   `tfsdk:"firewall_id"`
	SrcIp       types.String   `tfsdk:"src_ip"`
	DstIp       types.String   `tfsdk:"dst_ip"`
	Description types.String   `tfsdk:"description"`
	State       types.List     `tfsdk:"state"`
	Status      types.String   `tfsdk:"status"`
	FetchAll    types.Bool     `tfsdk:"fetch_all"`
	Ids         []types.String `tfsdk:"ids"`
}

type FirewallRuleResource struct {
	Id                 types.String       `tfsdk:"id"`
	FirewallRule       types.Object       `tfsdk:"firewall_rule"`
	FirewallId         types.String       `tfsdk:"firewall_id"`
	FirewallRuleCreate FirewallRuleCreate `tfsdk:"firewall_rule_create"`
}

type FirewallRuleCreate struct {
	SourceAddress      []string       `tfsdk:"source_address"`
	DestinationAddress []string       `tfsdk:"destination_address"`
	Service            []FirewallPort `tfsdk:"service"`
	Action             types.String   `tfsdk:"action"`
	Direction          types.String   `tfsdk:"direction"`
	Description        types.String   `tfsdk:"description"`
	Status             types.String   `tfsdk:"status"`
	OrderRuleId        types.String   `tfsdk:"order_rule_id"`
	OrderDirection     types.String   `tfsdk:"order_direction"`
}

type FirewallRule struct {
	Id                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	FirewallId           types.String   `tfsdk:"firewall_id"`
	Sequence             types.Int32    `tfsdk:"sequence"`
	SourceInterface      types.String   `tfsdk:"source_interface"`
	SourceAddress        []string       `tfsdk:"source_address"`
	DestinationInterface types.String   `tfsdk:"destination_interface"`
	DestinationAddress   []string       `tfsdk:"destination_address"`
	Service              []FirewallPort `tfsdk:"service"`
	Action               types.String   `tfsdk:"action"`
	Direction            types.String   `tfsdk:"direction"`
	VendorRuleId         types.String   `tfsdk:"vendor_rule_id"`
	Description          types.String   `tfsdk:"description"`
	State                types.String   `tfsdk:"state"`
	Status               types.String   `tfsdk:"status"`
	CreatedAt            types.String   `tfsdk:"created_at"`
	CreatedBy            types.String   `tfsdk:"created_by"`
	ModifiedAt           types.String   `tfsdk:"modified_at"`
	ModifiedBy           types.String   `tfsdk:"modified_by"`
}

type FirewallPort struct {
	ServiceType  types.String `tfsdk:"service_type"`
	ServiceValue types.String `tfsdk:"service_value"`
}

func convertFirewallPorts(ports []FirewallPort) []scpfirewall.FirewallPort {
	result := make([]scpfirewall.FirewallPort, len(ports))
	for i, port := range ports {
		sType := scpfirewall.FirewallServiceType(port.ServiceType.ValueString())
		sValue := port.ServiceValue.ValueString()
		result[i] = scpfirewall.FirewallPort{
			ServiceType:  sType,
			ServiceValue: &sValue,
		}
	}
	return result
}

func convertOrderDirection(val *string) *scpfirewall.FirewallRuleOrderDirection {
	if val == nil {
		return nil
	}
	od := scpfirewall.FirewallRuleOrderDirection(*val)
	nullableOd := scpfirewall.NewNullableFirewallRuleOrderDirection(&od)
	if nullableOd == nil {
		return nil
	}
	return nullableOd.Get()
}

func (m FirewallRule) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    types.StringType,
		"name":                  types.StringType,
		"firewall_id":           types.StringType,
		"sequence":              types.Int32Type,
		"source_interface":      types.StringType,
		"source_address":        types.ListType{ElemType: types.StringType},
		"destination_interface": types.StringType,
		"destination_address":   types.ListType{ElemType: types.StringType},
		"service":               types.ListType{ElemType: types.ObjectType{AttrTypes: FirewallPort{}.AttributeTypes()}},
		"action":                types.StringType,
		"direction":             types.StringType,
		"vendor_rule_id":        types.StringType,
		"description":           types.StringType,
		"state":                 types.StringType,
		"status":                types.StringType,
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"modified_at":           types.StringType,
		"modified_by":           types.StringType,
	}
}
func (fp FirewallPort) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"service_type":  types.StringType,
		"service_value": types.StringType,
	}
}
