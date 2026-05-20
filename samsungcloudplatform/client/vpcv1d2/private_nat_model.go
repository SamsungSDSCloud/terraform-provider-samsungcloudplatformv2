package vpcv1d2

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateNatDataSource struct {
	PrivateNatId types.String `tfsdk:"private_nat_id"` // Private NAT ID (input)
	PrivateNat   types.Object `tfsdk:"private_nat"`    // Private NAT details (output)
}

type PrivateNatDataSources struct {
	// Input
	Size                types.Int32  `tfsdk:"size"`                  // size
	Page                types.Int32  `tfsdk:"page"`                  // page
	Sort                types.String `tfsdk:"sort"`                  // sort
	Name                types.String `tfsdk:"name"`                  // Private NAT Name
	Cidr                types.String `tfsdk:"cidr"`                  // Private NAT IP range
	VpcId               types.String `tfsdk:"vpc_id"`                // VPC Id
	ServiceResourceId   types.String `tfsdk:"service_resource_id"`   // Private NAT connected Service Resource ID
	ServiceType         types.String `tfsdk:"service_type"`          // Private NAT connected Service Type
	ServiceResourceName types.String `tfsdk:"service_resource_name"` // Private NAT connected Service Resource Name
	State               types.String `tfsdk:"state"`                 // Private NAT State

	//Output
	TotalCount  types.Int32  `tfsdk:"total_count"`
	PrivateNats []PrivateNat `tfsdk:"private_nats"` // List of Private NATs
}

type PrivateNatResource struct {
	Cidr              types.String `tfsdk:"cidr"`                // Private NAT IP range
	Description       types.String `tfsdk:"description"`         // Description
	Name              types.String `tfsdk:"name"`                // Private NAT Name
	ServiceResourceId types.String `tfsdk:"service_resource_id"` // Private NAT connected Service Resource ID
	ServiceType       types.String `tfsdk:"service_type"`        // Private NAT connected Service Type
	Tags              types.Map    `tfsdk:"tags"`

	// Output
	Id         types.String `tfsdk:"id"` // Private NAT ID
	PrivateNat types.Object `tfsdk:"private_nat"`
}

type PrivateNat struct {
	AccountId           types.String `tfsdk:"account_id"`            // Account ID
	Cidr                types.String `tfsdk:"cidr"`                  // Private NAT IP range
	CreatedAt           types.String `tfsdk:"created_at"`            // Created At
	CreatedBy           types.String `tfsdk:"created_by"`            // Created By
	Description         types.String `tfsdk:"description"`           // Description
	Id                  types.String `tfsdk:"id"`                    // Private NAT ID
	ModifiedAt          types.String `tfsdk:"modified_at"`           // Modified At
	ModifiedBy          types.String `tfsdk:"modified_by"`           // Modified By
	Name                types.String `tfsdk:"name"`                  // Private NAT Name
	ServiceResourceId   types.String `tfsdk:"service_resource_id"`   // Private NAT connected Service Resource ID
	ServiceResourceName types.String `tfsdk:"service_resource_name"` // Private NAT connected Service Resource Name
	ServiceType         types.String `tfsdk:"service_type"`          // Private NAT connected Service Type
	State               types.String `tfsdk:"state"`                 // Private NAT State
}

func (m PrivateNat) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":            types.StringType,
		"cidr":                  types.StringType,
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"description":           types.StringType,
		"id":                    types.StringType,
		"modified_at":           types.StringType,
		"modified_by":           types.StringType,
		"name":                  types.StringType,
		"service_resource_id":   types.StringType,
		"service_resource_name": types.StringType,
		"service_type":          types.StringType,
		"state":                 types.StringType,
	}
}
