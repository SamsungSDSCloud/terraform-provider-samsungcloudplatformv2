package vpcv1d3

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ------------ VPC Cidr Resource ------------

type VpcCidrResource struct {
	// Input
	VpcId types.String `tfsdk:"vpc_id"`
	Cidr  types.String `tfsdk:"cidr"`

	// Output
	Id types.String `tfsdk:"id"`
}

// ------------ VPC Cidr Response Value Mapping ------------

type VpcCidrValue struct {
	Cidr types.String `tfsdk:"cidr"`
	Id   types.String `tfsdk:"id"`
}
