package vpcv1d3

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type TransitGatewayFirewallResource struct {
	// Input
	TransitGatewayId  types.String `tfsdk:"transit_gateway_id"`
	ProductType       types.String `tfsdk:"product_type"`
	UplinkActiveZone  types.String `tfsdk:"uplink_active_zone"`  // ex: kr-west1-a
	UplinkStandbyZone types.String `tfsdk:"uplink_standby_zone"` // ex: kr-west1-b

	// Output
	TransitGatewayFirewall types.Object `tfsdk:"transit_gateway_firewall"`
}

type TransitGatewayFirewall struct {
	State             basetypes.StringValue `tfsdk:"state"`
	UplinkActiveZone  basetypes.StringValue `tfsdk:"uplink_active_zone"`
	UplinkStandbyZone basetypes.StringValue `tfsdk:"uplink_standby_zone"`
	UplinkZoneState   basetypes.StringValue `tfsdk:"uplink_zone_state"`
}

func (v TransitGatewayFirewall) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"state":               basetypes.StringType{},
		"uplink_active_zone":  basetypes.StringType{},
		"uplink_standby_zone": basetypes.StringType{},
		"uplink_zone_state":   basetypes.StringType{},
	}
}
