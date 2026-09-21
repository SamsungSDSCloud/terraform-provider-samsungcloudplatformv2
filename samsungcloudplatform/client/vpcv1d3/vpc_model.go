package vpcv1d3

import (
	"time"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ------------ VPC Resource ------------

type VpcResource struct {
	// Input
	Cidr        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
	Tags        types.Map    `tfsdk:"tags"`

	// Output
	Id  types.String `tfsdk:"id"`
	Vpc types.Object `tfsdk:"vpc"`
}

// ------------ VPC Response Value Mapping ------------

type VpcValue struct {
	AccountId   types.String `tfsdk:"account_id"`
	CidrCount   types.Int32  `tfsdk:"cidr_count"`
	Cidrs       []CidrsValue `tfsdk:"cidrs"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	Name        types.String `tfsdk:"name"`
	State       types.String `tfsdk:"state"`
	ZoneType    types.String `tfsdk:"zone_type"`
	Zones       types.List   `tfsdk:"zones"`
}

func (v VpcValue) AttributeTypes() map[string]attr.Type {
	cidrsAttrTypes := map[string]attr.Type{
		"cidr":       types.StringType,
		"created_at": types.StringType,
		"created_by": types.StringType,
		"id":         types.StringType,
	}

	return map[string]attr.Type{
		"account_id":  types.StringType,
		"cidr_count":  types.Int32Type,
		"cidrs":       types.ListType{ElemType: types.ObjectType{AttrTypes: cidrsAttrTypes}},
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"description": types.StringType,
		"id":          types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"name":        types.StringType,
		"state":       types.StringType,
		"zone_type":   types.StringType,
		"zones":       types.ListType{ElemType: types.StringType},
	}
}

func ResponseToVpcValue(res vpc.VpcV1Dot3) VpcValue {
	vpcTf := VpcValue{
		AccountId:   types.StringValue(res.AccountId),
		CidrCount:   types.Int32Value(res.CidrCount),
		Cidrs:       mapCidrs(res.Cidrs),
		CreatedAt:   types.StringValue(res.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(res.CreatedBy),
		Description: types.StringPointerValue(res.Description.Get()),
		Id:          types.StringValue(res.Id),
		ModifiedAt:  types.StringValue(res.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(res.ModifiedBy),
		Name:        types.StringValue(res.Name),
		State:       types.StringValue(string(res.State)),
		ZoneType:    types.StringPointerValue(res.ZoneType.Get()),
		Zones:       mapZones(res.Zones),
	}

	return vpcTf
}

func mapCidrs(sdkCidrs []vpc.VpcCidr) []CidrsValue {
	cidrs := make([]CidrsValue, len(sdkCidrs))

	for pos, sdkCidr := range sdkCidrs {
		cidrs[pos] = CidrsValue{
			Cidr:      types.StringValue(sdkCidr.Cidr),
			CreatedAt: types.StringValue(sdkCidr.CreatedAt.Format(time.RFC3339)),
			CreatedBy: types.StringValue(sdkCidr.CreatedBy),
			Id:        types.StringValue(sdkCidr.Id),
		}
	}

	return cidrs
}

func mapZones(zones []string) types.List {
	if zones == nil {
		zones = []string{}
	}
	var elements []attr.Value
	for _, z := range zones {
		elements = append(elements, types.StringValue(z))
	}
	return types.ListValueMust(types.StringType, elements)
}

// ------------ CidrsValue ------------

type CidrsValue struct {
	Cidr      types.String `tfsdk:"cidr"`
	CreatedAt types.String `tfsdk:"created_at"`
	CreatedBy types.String `tfsdk:"created_by"`
	Id        types.String `tfsdk:"id"`
}

// ------------ VPC Data Source ------------

type VpcDataSource struct {
	Cidr       types.String `tfsdk:"cidr"`
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Page       types.Int32  `tfsdk:"page"`
	Size       types.Int32  `tfsdk:"size"`
	Sort       types.String `tfsdk:"sort"`
	State      types.String `tfsdk:"state"`
	Zone       types.String `tfsdk:"zone"`
	TotalCount types.Int32  `tfsdk:"total_count"`
	Vpcs       []VpcDSValue `tfsdk:"vpcs"`
}

// ------------ VPC Data Source Response Value ------------

type VpcDSValue struct {
	AccountId   types.String `tfsdk:"account_id"`
	CidrCount   types.Int32  `tfsdk:"cidr_count"`
	Cidrs       []CidrsValue `tfsdk:"cidrs"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	Name        types.String `tfsdk:"name"`
	State       types.String `tfsdk:"state"`
	ZoneType    types.String `tfsdk:"zone_type"`
	Zones       types.List   `tfsdk:"zones"`
}

func (v VpcDSValue) AttributeTypes() map[string]attr.Type {
	cidrsAttrTypes := map[string]attr.Type{
		"cidr":       types.StringType,
		"created_at": types.StringType,
		"created_by": types.StringType,
		"id":         types.StringType,
	}

	return map[string]attr.Type{
		"account_id":  types.StringType,
		"cidr_count":  types.Int32Type,
		"cidrs":       types.ListType{ElemType: types.ObjectType{AttrTypes: cidrsAttrTypes}},
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"description": types.StringType,
		"id":          types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"name":        types.StringType,
		"state":       types.StringType,
		"zone_type":   types.StringType,
		"zones":       types.ListType{ElemType: types.StringType},
	}
}

func ResponseToVpcDSValue(res vpc.VpcV1Dot3) VpcDSValue {
	vpcTf := VpcDSValue{
		AccountId:   types.StringValue(res.AccountId),
		CidrCount:   types.Int32Value(res.CidrCount),
		Cidrs:       mapCidrs(res.Cidrs),
		CreatedAt:   types.StringValue(res.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(res.CreatedBy),
		Description: types.StringPointerValue(res.Description.Get()),
		Id:          types.StringValue(res.Id),
		ModifiedAt:  types.StringValue(res.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(res.ModifiedBy),
		Name:        types.StringValue(res.Name),
		State:       types.StringValue(string(res.State)),
		ZoneType:    types.StringPointerValue(res.ZoneType.Get()),
		Zones:       mapZones(res.Zones),
	}

	return vpcTf
}
