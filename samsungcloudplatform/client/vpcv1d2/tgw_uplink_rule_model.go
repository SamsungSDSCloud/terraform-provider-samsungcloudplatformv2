package vpcv1d2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type TransitGatewayUplinkValue struct {
	Description        types.String `tfsdk:"description"`
	DestinationCidr    types.String `tfsdk:"destination_cidr"`
	DestinationType    types.String `tfsdk:"destination_type"`
	TransitGatewayId   types.String `tfsdk:"transit_gateway_id"`
	TransitGatewayRule types.Object `tfsdk:"transit_gateway_rule"`
}

type TransitGatewayRuleValue struct {
	Description     basetypes.StringValue `tfsdk:"description"`
	DestinationCidr basetypes.StringValue `tfsdk:"destination_cidr"`
	DestinationType basetypes.StringValue `tfsdk:"destination_type"`
	Id              basetypes.StringValue `tfsdk:"id"`
	State           basetypes.StringValue `tfsdk:"state"`
}

type TransitGatewayRuleType struct {
	basetypes.ObjectType
}

func (v TransitGatewayRuleValue) Type(ctx context.Context) attr.Type {
	return TransitGatewayRuleType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(),
		},
	}
}

func (v TransitGatewayRuleValue) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"description":      basetypes.StringType{},
		"destination_cidr": basetypes.StringType{},
		"destination_type": basetypes.StringType{},
		"id":               basetypes.StringType{},
		"state":            basetypes.StringType{},
	}
}
