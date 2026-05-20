package vpcv1d2

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NatGatewayDataSource struct {
	// Input
	Size                types.Int32  `tfsdk:"size"`
	Page                types.Int32  `tfsdk:"page"`
	Sort                types.String `tfsdk:"sort"`
	Name                types.String `tfsdk:"name"`
	VpcName             types.String `tfsdk:"vpc_name"`
	VpcId               types.String `tfsdk:"vpc_id"`
	SubnetId            types.String `tfsdk:"subnet_id"`
	SubnetName          types.String `tfsdk:"subnet_name"`
	NatGatewayIpAddress types.String `tfsdk:"nat_gateway_ip_address"`
	State               types.String `tfsdk:"state"`

	// Output
	NatGateways []NatGateway   `tfsdk:"nat_gateways"`
	TotalCount  types.Int32    `tfsdk:"total_count"`
	SortFinal   []types.String `tfsdk:"sort_final"` // Sort output
}

type NatGateway struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	NatGatewayIpAddress types.String `tfsdk:"nat_gateway_ip_address"`
	VpcId               types.String `tfsdk:"vpc_id"`
	VpcName             types.String `tfsdk:"vpc_name"`
	SubnetId            types.String `tfsdk:"subnet_id"`
	SubnetName          types.String `tfsdk:"subnet_name"`
	SubnetCidr          types.String `tfsdk:"subnet_cidr"`
	AccountId           types.String `tfsdk:"account_id"`
	State               types.String `tfsdk:"state"`
	Description         types.String `tfsdk:"description"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
	PublicipId          types.String `tfsdk:"publicip_id"`
}
