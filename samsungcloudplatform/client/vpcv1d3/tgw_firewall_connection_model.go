package vpcv1d3

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type TransitGatewayFirewallConnectionResource struct {
	// Input
	TransitGatewayId types.String `tfsdk:"transit_gateway_id"`

	// Output
	TransitGatewayFirewallConnection types.Object `tfsdk:"transit_gateway_firewall_connection"`
}

type TransitGatewayFirewallConnection struct {
	FirewallConnectionState basetypes.StringValue `tfsdk:"firewall_connection_state"`
}

func (v TransitGatewayFirewallConnection) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"firewall_connection_state": basetypes.StringType{},
	}
}
