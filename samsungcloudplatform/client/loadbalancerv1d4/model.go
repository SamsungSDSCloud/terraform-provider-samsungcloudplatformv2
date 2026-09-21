package loadbalancerv1d4

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ------------ Load Balancer -------------------//
type LoadbalancerDataSource struct {
	Size          types.Int32    `tfsdk:"size"`
	Page          types.Int32    `tfsdk:"page"`
	Sort          types.String   `tfsdk:"sort"`
	Name          types.String   `tfsdk:"name"`
	ServiceIp     types.String   `tfsdk:"service_ip"`
	SubnetId      types.String   `tfsdk:"subnet_id"`
	VpcId         types.String   `tfsdk:"vpc_id"`
	Loadbalancers []Loadbalancer `tfsdk:"loadbalancers"`
}

// list response
type Loadbalancer struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	ServiceIp        types.String `tfsdk:"service_ip"`
	SourceNatIp      types.String `tfsdk:"source_nat_ip"`
	State            types.String `tfsdk:"state"`
	PublicNatEnabled types.Bool   `tfsdk:"public_nat_enabled"`
	LayerType        types.String `tfsdk:"layer_type"`
	SubnetId         types.String `tfsdk:"subnet_id"`
	VpcId            types.String `tfsdk:"vpc_id"`
	FirewallId       types.String `tfsdk:"firewall_id"`
	ListenerCount    types.Int32  `tfsdk:"listener_count"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	ModifiedAt       types.String `tfsdk:"modified_at"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
	Zones            types.List   `tfsdk:"zones"`
}

type LoadbalancerDataSourceDetail struct {
	Id                 types.String `tfsdk:"id"`
	LoadbalancerDetail types.Object `tfsdk:"loadbalancer"`
}

type LoadbalancerDetail struct {
	AccountId        types.String `tfsdk:"account_id"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	Description      types.String `tfsdk:"description"`
	FirewallId       types.String `tfsdk:"firewall_id"`
	Id               types.String `tfsdk:"id"`
	LayerType        types.String `tfsdk:"layer_type"`
	ModifiedAt       types.String `tfsdk:"modified_at"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
	Name             types.String `tfsdk:"name"`
	PublicNatEnabled types.Bool   `tfsdk:"public_nat_enabled"`
	ServiceIp        types.String `tfsdk:"service_ip"`
	SourceNatIp      types.String `tfsdk:"source_nat_ip"`
	State            types.String `tfsdk:"state"`
	SubnetId         types.String `tfsdk:"subnet_id"`
	VpcId            types.String `tfsdk:"vpc_id"`
	Zones            types.List   `tfsdk:"zones"`
	HealthCheckIps   types.List   `tfsdk:"health_check_ips"`
}

func (m LoadbalancerDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":         types.StringType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"description":        types.StringType,
		"firewall_id":        types.StringType,
		"id":                 types.StringType,
		"layer_type":         types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
		"name":               types.StringType,
		"public_nat_enabled": types.BoolType,
		"service_ip":         types.StringType,
		"source_nat_ip":      types.StringType,
		"state":              types.StringType,
		"subnet_id":          types.StringType,
		"vpc_id":             types.StringType,
		"zones": types.ListType{
			ElemType: types.StringType,
		},
		"health_check_ips": types.ListType{
			ElemType: types.StringType,
		},
	}
}

type LoadbalancerResource struct {
	Id                 types.String        `tfsdk:"id"`
	Loadbalancer       types.Object        `tfsdk:"loadbalancer"`
	LoadbalancerCreate *LoadbalancerCreate `tfsdk:"loadbalancer_create"`
}

type LoadbalancerCreate struct {
	Description            types.String `tfsdk:"description"`
	FirewallEnabled        types.Bool   `tfsdk:"firewall_enabled"`
	FirewallLoggingEnabled types.Bool   `tfsdk:"firewall_logging_enabled"`
	LayerType              types.String `tfsdk:"layer_type"`
	Name                   types.String `tfsdk:"name"`
	ServiceIp              types.String `tfsdk:"service_ip"`
	PublicipId             types.String `tfsdk:"publicip_id"`
	SubnetId               types.String `tfsdk:"subnet_id"`
	VpcId                  types.String `tfsdk:"vpc_id"`
	SourceNatIp            types.String `tfsdk:"source_nat_ip"`
	Tags                   types.Map    `tfsdk:"tags"`
	Zones                  types.List   `tfsdk:"zones"`
	HealthCheckIps         types.List   `tfsdk:"health_check_ips"`
}

type LoadbalancerCreateResponseDetail struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	LayerType        types.String `tfsdk:"layer_type"`
	VpcId            types.String `tfsdk:"vpc_id"`
	SubnetId         types.String `tfsdk:"subnet_id"`
	AccountId        types.String `tfsdk:"account_id"`
	State            types.String `tfsdk:"state"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	ModifiedAt       types.String `tfsdk:"modified_at"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
	FirewallId       types.String `tfsdk:"firewall_id"`
	PublicNatEnabled types.Bool   `tfsdk:"public_nat_enabled"`
	ServiceIp        types.String `tfsdk:"service_ip"`
	SourceNatIp      types.String `tfsdk:"source_nat_ip"`
	Zones            types.List   `tfsdk:"zones"`
	HealthCheckIps   types.List   `tfsdk:"health_check_ips"`
}

func (m LoadbalancerCreateResponseDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                 types.StringType,
		"name":               types.StringType,
		"description":        types.StringType,
		"layer_type":         types.StringType,
		"state":              types.StringType,
		"subnet_id":          types.StringType,
		"vpc_id":             types.StringType,
		"account_id":         types.StringType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
		"firewall_id":        types.StringType,
		"public_nat_enabled": types.BoolType,
		"service_ip":         types.StringType,
		"source_nat_ip":      types.StringType,
		"zones": types.ListType{
			ElemType: types.StringType,
		},
		"health_check_ips": types.ListType{
			ElemType: types.StringType,
		},
	}
}
