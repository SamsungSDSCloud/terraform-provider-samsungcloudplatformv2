package vpcv1d2

import (
	"context"
	"fmt"
	"strings"
	"time"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpc/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type VpcDataSource struct {
	Cidr       types.String `tfsdk:"cidr"`
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Page       types.Int32  `tfsdk:"page"`
	Size       types.Int32  `tfsdk:"size"`
	Sort       types.String `tfsdk:"sort"`
	State      types.String `tfsdk:"state"`
	TotalCount types.Int32  `tfsdk:"total_count"`
	Vpcs       []VpcDSValue `tfsdk:"vpcs"`
}

type VpcResource struct {
	Cidr        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Tags        types.Map    `tfsdk:"tags"`
	Vpc         types.Object `tfsdk:"vpc"`
}

func ResponseToVpcDSValue(res vpc.VpcV1Dot2) VpcDSValue {

	vpcTf := VpcDSValue{
		AccountId:   types.StringValue(res.AccountId),
		CidrCount:   types.Int32Value(res.CidrCount),
		Cidrs:       mapCidr(res.Cidrs),
		CreatedAt:   types.StringValue(res.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(res.CreatedBy),
		Description: types.StringPointerValue(res.Description.Get()),
		Id:          types.StringValue(res.Id),
		ModifiedAt:  types.StringValue(res.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(res.ModifiedBy),
		Name:        types.StringValue(res.Name),
		State:       types.StringValue(string(res.State)),
	}

	return vpcTf

}

func mapCidr(sdkCidrs []vpc.VpcCidr) []CidrsValue {

	cidrs := make([]CidrsValue, len(sdkCidrs))

	for pos, sdkCidr := range sdkCidrs {
		cidrs[pos] = CidrsValue{
			Cidr:      types.StringValue(sdkCidr.Cidr),
			CreatedAt: types.StringValue(sdkCidr.CreatedAt.Format(time.RFC3339)),
			CreatedBy: types.StringValue(sdkCidr.CreatedBy),
			Id:        types.StringValue(sdkCidr.Id),
			state:     attr.ValueStateKnown,
		}

	}

	return cidrs

}

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
}

type VpcDSType struct {
	basetypes.ObjectType
}

type CidrsValue struct {
	Cidr      types.String `tfsdk:"cidr"`
	CreatedAt types.String `tfsdk:"created_at"`
	CreatedBy types.String `tfsdk:"created_by"`
	Id        types.String `tfsdk:"id"`
	state     attr.ValueState
}

func (v VpcDSValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"account_id": basetypes.StringType{},
		"cidr_count": basetypes.Int32Type{},
		"cidrs": basetypes.ListType{
			ElemType: CidrsValue{}.Type(ctx),
		},
		"created_at":  basetypes.StringType{},
		"created_by":  basetypes.StringType{},
		"description": basetypes.StringType{},
		"id":          basetypes.StringType{},
		"modified_at": basetypes.StringType{},
		"modified_by": basetypes.StringType{},
		"name":        basetypes.StringType{},
		"state":       basetypes.StringType{},
	}
}

type CidrsType struct {
	types.ObjectType
}

func (t CidrsType) ValueType(ctx context.Context) attr.Value {
	return CidrsValue{}
}

func (v CidrsValue) Type(ctx context.Context) attr.Type {
	return CidrsType{
		types.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CidrsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"cidr":       types.StringType,
		"created_at": types.StringType,
		"created_by": types.StringType,
		"id":         types.StringType,
	}
}

func (t CidrsType) Equal(o attr.Type) bool {
	other, ok := o.(CidrsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CidrsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {

	if in.Type() == nil {
		return NewCidrsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCidrsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCidrsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCidrsValueMust(CidrsValue{}.AttributeTypes(ctx), attributes), nil

}

func (v CidrsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 4)

	var val tftypes.Value
	var err error

	attrTypes["cidr"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["created_at"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["created_by"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["id"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 4)

		val, err = v.Cidr.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["cidr"] = val

		val, err = v.CreatedAt.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["created_at"] = val

		val, err = v.CreatedBy.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["created_by"] = val

		val, err = v.Id.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["id"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CidrsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CidrsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CidrsValue) String() string {
	return "CidrsValue"
}

func (v CidrsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"cidr":       basetypes.StringType{},
		"created_at": basetypes.StringType{},
		"created_by": basetypes.StringType{},
		"id":         basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"cidr":       v.Cidr,
			"created_at": v.CreatedAt,
			"created_by": v.CreatedBy,
			"id":         v.Id,
		})

	return objVal, diags
}

func (v CidrsValue) Equal(o attr.Value) bool {
	other, ok := o.(CidrsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Cidr.Equal(other.Cidr) {
		return false
	}

	if !v.CreatedAt.Equal(other.CreatedAt) {
		return false
	}

	if !v.CreatedBy.Equal(other.CreatedBy) {
		return false
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	return true
}

func NewCidrsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CidrsValue {
	object, diags := NewCidrsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewCidrsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func NewCidrsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CidrsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CidrsValue Attribute Value",
				"While creating a CidrsValue value, an extra attribute value was detected. "+
					"A CidrsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CidrsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCidrsValueUnknown(), diags
	}

	cidrAttribute, ok := attributes["cidr"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`cidr is missing from object`)

		return NewCidrsValueUnknown(), diags
	}

	cidrVal, ok := cidrAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`cidr expected to be basetypes.StringValue, was: %T`, cidrAttribute))
	}

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return NewCidrsValueUnknown(), diags
	}

	createdAtVal, ok := createdAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be basetypes.StringValue, was: %T`, createdAtAttribute))
	}

	createdByAttribute, ok := attributes["created_by"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_by is missing from object`)

		return NewCidrsValueUnknown(), diags
	}

	createdByVal, ok := createdByAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_by expected to be basetypes.StringValue, was: %T`, createdByAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewCidrsValueUnknown(), diags
	}

	idVal, ok := idAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.StringValue, was: %T`, idAttribute))
	}

	if diags.HasError() {
		return NewCidrsValueUnknown(), diags
	}

	return CidrsValue{
		Cidr:      cidrVal,
		CreatedAt: createdAtVal,
		CreatedBy: createdByVal,
		Id:        idVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewCidrsValueUnknown() CidrsValue {
	return CidrsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCidrsValueNull() CidrsValue {
	return CidrsValue{
		state: attr.ValueStateNull,
	}
}
