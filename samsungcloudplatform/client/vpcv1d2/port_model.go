package vpcv1d2

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpc/1.2"
)

type SecurityGroup struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type PortResource struct {
	AccountId            types.String    `tfsdk:"account_id"`
	AttachedResourceId   types.String    `tfsdk:"attached_resource_id"`
	AttachedResourceType types.String    `tfsdk:"attached_resource_type"`
	CreatedAt            types.String    `tfsdk:"created_at"`
	Description          types.String    `tfsdk:"description"`
	FixedIpAddress       types.String    `tfsdk:"fixed_ip_address"`
	Id                   types.String    `tfsdk:"id"`
	MacAddress           types.String    `tfsdk:"mac_address"`
	ModifiedAt           types.String    `tfsdk:"modified_at"`
	Name                 types.String    `tfsdk:"name"`
	SecurityGroups       []SecurityGroup `tfsdk:"security_groups"`
	State                types.String    `tfsdk:"state"`
	SubnetId             types.String    `tfsdk:"subnet_id"`
	SubnetName           types.String    `tfsdk:"subnet_name"`
	VirtualIpAddresses   types.List      `tfsdk:"virtual_ip_addresses"`
	Tags                 types.Map       `tfsdk:"tags"`
	VpcId                types.String    `tfsdk:"vpc_id"`
	VpcName              types.String    `tfsdk:"vpc_name"`
}

type Port struct {
	AccountId            types.String `tfsdk:"account_id"`
	AttachedResourceId   types.String `tfsdk:"attached_resource_id"`
	AttachedResourceType types.String `tfsdk:"attached_resource_type"`
	CreatedAt            types.String `tfsdk:"created_at"`
	Description          types.String `tfsdk:"description"`
	FixedIpAddress       types.String `tfsdk:"fixed_ip_address"`
	Id                   types.String `tfsdk:"id"`
	MacAddress           types.String `tfsdk:"mac_address"`
	ModifiedAt           types.String `tfsdk:"modified_at"`
	Name                 types.String `tfsdk:"name"`
	State                types.String `tfsdk:"state"`
	SubnetId             types.String `tfsdk:"subnet_id"`
	SubnetName           types.String `tfsdk:"subnet_name"`
	VpcId                types.String `tfsdk:"vpc_id"`
	VpcName              types.String `tfsdk:"vpc_name"`
}

func MapPort(port *vpc.PortV1Dot2, state *PortResource) {

	state.Id = types.StringValue(port.Id)
	state.AccountId = types.StringValue(port.AccountId)
	state.AttachedResourceId = types.StringValue(port.AttachedResourceId)
	state.AttachedResourceType = types.StringValue(port.AttachedResourceType)
	state.CreatedAt = types.StringValue(port.CreatedAt.Format(time.RFC3339))
	state.Description = types.StringValue(port.Description)
	state.FixedIpAddress = types.StringValue(port.FixedIpAddress)
	state.MacAddress = types.StringValue(port.MacAddress)
	state.ModifiedAt = types.StringValue(port.ModifiedAt.Format(time.RFC3339))
	state.Name = types.StringValue(port.Name)
	state.State = types.StringValue(port.State)
	state.SubnetId = types.StringValue(port.SubnetId)
	state.SubnetName = types.StringValue(port.SubnetName)
	state.VpcId = types.StringValue(port.VpcId)
	state.VpcName = types.StringValue(port.VpcName)

	state.SecurityGroups = make([]SecurityGroup, len(port.SecurityGroups))
	for i, sg := range port.SecurityGroups {
		state.SecurityGroups[i] = SecurityGroup{
			Id:   types.StringValue(*sg.Id.Get()),
			Name: types.StringValue(*sg.Name.Get()),
		}
	}

	var vipElements []attr.Value
	for _, ip := range port.VirtualIpAddresses {
		vipElements = append(vipElements, types.StringValue(ip))
	}
	if vipElements == nil {
		vipElements = []attr.Value{}
	}
	state.VirtualIpAddresses = types.ListValueMust(types.StringType, vipElements)

}
