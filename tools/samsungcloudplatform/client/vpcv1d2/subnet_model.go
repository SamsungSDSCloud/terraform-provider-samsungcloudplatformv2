package vpcv1d2

import (
	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpc/1.2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SubnetDataSource struct {
	Cidr       types.String     `tfsdk:"cidr"`
	Id         types.String     `tfsdk:"id"`
	Name       types.String     `tfsdk:"name"`
	Page       types.Int32      `tfsdk:"page"`
	Size       types.Int32      `tfsdk:"size"`
	Sort       types.String     `tfsdk:"sort"`
	State      types.String     `tfsdk:"state"`
	Subnets    []Subnet         `tfsdk:"subnets"`
	TotalCount types.Int32      `tfsdk:"total_count"`
	Type       []vpc.SubnetType `tfsdk:"type"`
	VpcId      types.String     `tfsdk:"vpc_id"`
	VpcName    types.String     `tfsdk:"vpc_name"`
}

type Subnet struct {
	AccountId        types.String `tfsdk:"account_id"`
	Cidr             types.String `tfsdk:"cidr"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	GatewayIpAddress types.String `tfsdk:"gateway_ip_address"`
	Id               types.String `tfsdk:"id"`
	ModifiedAt       types.String `tfsdk:"modified_at"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
	Name             types.String `tfsdk:"name"`
	State            types.String `tfsdk:"state"`
	Type             types.String `tfsdk:"type"`
	VpcId            types.String `tfsdk:"vpc_id"`
	VpcName          types.String `tfsdk:"vpc_name"`
}

type SubnetResource struct {
	AccountId        types.String     `tfsdk:"account_id"`
	AllocationPools  []AllocationPool `tfsdk:"allocation_pools"`
	Cidr             types.String     `tfsdk:"cidr"`
	CreatedAt        types.String     `tfsdk:"created_at"`
	CreatedBy        types.String     `tfsdk:"created_by"`
	Description      types.String     `tfsdk:"description"`
	DhcpIpAddress    types.String     `tfsdk:"dhcp_ip_address"`
	DnsNameservers   types.Set        `tfsdk:"dns_nameservers"`
	GatewayIpAddress types.String     `tfsdk:"gateway_ip_address"`
	HostRoutes       []HostRoute      `tfsdk:"host_routes"`
	Id               types.String     `tfsdk:"id"`
	ModifiedAt       types.String     `tfsdk:"modified_at"`
	ModifiedBy       types.String     `tfsdk:"modified_by"`
	Name             types.String     `tfsdk:"name"`
	State            types.String     `tfsdk:"state"`
	Type             types.String     `tfsdk:"type"`
	VpcId            types.String     `tfsdk:"vpc_id"`
	VpcName          types.String     `tfsdk:"vpc_name"`
	Tags             types.Map        `tfsdk:"tags"`
}

type AllocationPool struct {
	Start types.String `tfsdk:"start"`
	End   types.String `tfsdk:"end"`
}

type HostRoute struct {
	Destination types.String `tfsdk:"destination"`
	Nexthop     types.String `tfsdk:"nexthop"`
}

func convertAllocationPoolsToInterface(pools []AllocationPool) []interface{} {
	result := make([]interface{}, len(pools))
	for i, pool := range pools {
		result[i] = map[string]string{
			"start": pool.Start.ValueString(),
			"end":   pool.End.ValueString(),
		}
	}
	return result
}

func convertHostRoutesToInterface(routes []HostRoute) []interface{} {
	result := make([]interface{}, len(routes))
	for i, route := range routes {
		result[i] = map[string]string{
			"destination": route.Destination.ValueString(),
			"nexthop":     route.Nexthop.ValueString(),
		}
	}
	return result
}

func convertDnsNameserversToString(nameservers types.Set) []string {
	if nameservers.IsNull() || nameservers.IsUnknown() {
		return nil
	}
	elements := nameservers.Elements()
	if len(elements) == 0 {
		return nil
	}
	result := make([]string, len(elements))
	for i, elem := range elements {
		result[i] = elem.(types.String).ValueString()
	}
	return result
}
