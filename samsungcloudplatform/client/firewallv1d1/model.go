package firewall

import (
	scpfirewall "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/firewall/1.1"
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
	Page        types.Int32  `tfsdk:"page"`
	Size        types.Int32  `tfsdk:"size"`
	Sort        types.String `tfsdk:"sort"`
	Name        types.String `tfsdk:"name"`
	VpcName     types.String `tfsdk:"vpc_name"`
	ProductType types.List   `tfsdk:"product_type"`
	State       types.List   `tfsdk:"state"`

	//Output
	TotalCount types.Int32 `tfsdk:"total_count"`
	Firewalls  []Firewall  `tfsdk:"firewalls"`
}

type Firewall struct {
	Id              types.String    `tfsdk:"id"`
	AccountId       types.String    `tfsdk:"account_id"`
	Name            types.String    `tfsdk:"name"`
	VpcId           types.String    `tfsdk:"vpc_id"`
	VpcName         types.String    `tfsdk:"vpc_name"`
	Loggable        types.Bool      `tfsdk:"loggable"`
	PreProductId    types.String    `tfsdk:"pre_product_id"`
	ProductType     types.String    `tfsdk:"product_type"`
	State           types.String    `tfsdk:"state"`
	Status          types.String    `tfsdk:"status"`
	TotalRuleCount  types.Int32     `tfsdk:"total_rule_count"`
	FlavorName      types.String    `tfsdk:"flavor_name"`
	FlavorRuleQuota types.Int32     `tfsdk:"flavor_rule_quota"`
	CreatedAt       types.String    `tfsdk:"created_at"`
	CreatedBy       types.String    `tfsdk:"created_by"`
	ModifiedAt      types.String    `tfsdk:"modified_at"`
	ModifiedBy      types.String    `tfsdk:"modified_by"`
	ZoneResources   []ZoneResources `tfsdk:"zone_resources"`
}

func (m Firewall) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                types.StringType,
		"account_id":        types.StringType,
		"name":              types.StringType,
		"vpc_id":            types.StringType,
		"vpc_name":          types.StringType,
		"loggable":          types.BoolType,
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
		"zone_resources": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"allocate_state": types.StringType,
					"fw_resource_id": types.StringType,
					"zone":           types.StringType,
				},
			},
		},
	}
}

type ZoneResources struct {
	AllocateState types.String `tfsdk:"allocate_state"`
	FwResourceId  types.String `tfsdk:"fw_resource_id"`
	Zone          types.String `tfsdk:"zone"`
}

type Resource struct {
	Id         types.String `tfsdk:"id"`
	FlavorName types.String `tfsdk:"flavor_name"`
	Loggable   types.Bool   `tfsdk:"loggable"`
	Firewall   types.Object `tfsdk:"firewall"`
}

//------------------- Firewall Rule -------------------//

type FirewallRuleResource struct {
	Id         types.String `tfsdk:"id"`
	FirewallId types.String `tfsdk:"firewall_id"`

	Action             types.String   `tfsdk:"action"`
	Description        types.String   `tfsdk:"description"`
	DestinationAddress []string       `tfsdk:"destination_address"`
	Direction          types.String   `tfsdk:"direction"`
	OrderDirection     types.String   `tfsdk:"order_direction"`
	OrderRuleId        types.String   `tfsdk:"order_rule_id"`
	Service            []FirewallPort `tfsdk:"service"`
	SourceAddress      []string       `tfsdk:"source_address"`
	Status             types.String   `tfsdk:"status"`

	// Response
	FirewallRule types.Object `tfsdk:"firewall_rule"`
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
