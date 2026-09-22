package vpcv1d3

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/vpc/1.3"
)

type TgwResource struct {
	Id          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
	Tags        types.Map    `tfsdk:"tags"`
	Tgw         types.Object `tfsdk:"tgw"`
}

type TgwDataSource struct {
	FirewallConnectionState types.String `tfsdk:"firewall_connection_state"`
	Size                    types.Int32  `tfsdk:"size"`
	Page                    types.Int32  `tfsdk:"page"`
	Sort                    types.String `tfsdk:"sort"`
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	State                   types.String `tfsdk:"state"`
	UplinkZoneState         types.String `tfsdk:"uplink_zone_state"`
	Tgws                    []Tgw        `tfsdk:"tgws"`
	TotalCount              types.Int32  `tfsdk:"total_count"`
}

type TgwDataSourceDetail struct {
	Id             types.String `tfsdk:"id"`
	TransitGateway types.Object `tfsdk:"transit_gateway"`
}

type Tgw struct {
	Id                      types.String `tfsdk:"id"`
	Description             types.String `tfsdk:"description"`
	Name                    types.String `tfsdk:"name"`
	AccountId               types.String `tfsdk:"account_id"`
	Bandwidth               types.Int32  `tfsdk:"bandwidth"`
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	FirewallConnectionState types.String `tfsdk:"firewall_connection_state"`
	FirewallIds             types.String `tfsdk:"firewall_ids"`
	ModifiedAt              types.String `tfsdk:"modified_at"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	State                   types.String `tfsdk:"state"`
	UplinkEnabled           types.Bool   `tfsdk:"uplink_enabled"`
	UplinkActiveZone        types.String `tfsdk:"uplink_active_zone"`
	UplinkStandbyZone       types.String `tfsdk:"uplink_standby_zone"`
	UplinkZoneState         types.String `tfsdk:"uplink_zone_state"`
}

func MapToTgw(tgwRes vpc.TransitGatewayV1Dot3) Tgw {
	return Tgw{
		Id:                      types.StringValue(tgwRes.Id),
		Description:             types.StringPointerValue(tgwRes.Description.Get()),
		Name:                    types.StringValue(tgwRes.Name),
		AccountId:               types.StringValue(tgwRes.AccountId),
		Bandwidth:               types.Int32PointerValue(tgwRes.Bandwidth.Get()),
		CreatedAt:               types.StringValue(tgwRes.CreatedAt.Format(time.RFC3339)),
		CreatedBy:               types.StringValue(tgwRes.CreatedBy),
		FirewallConnectionState: types.StringValue(string(*(tgwRes.FirewallConnectionState.Get()))),
		FirewallIds:             types.StringPointerValue(tgwRes.FirewallIds.Get()),
		ModifiedAt:              types.StringValue(tgwRes.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:              types.StringValue(tgwRes.ModifiedBy),
		State:                   types.StringValue(string(tgwRes.State)),
		UplinkEnabled:           types.BoolValue(tgwRes.GetUplinkEnabled()),
		UplinkActiveZone:        types.StringPointerValue(tgwRes.UplinkActiveZone.Get()),
		UplinkStandbyZone:       types.StringPointerValue(tgwRes.UplinkStandbyZone.Get()),
		UplinkZoneState:         types.StringValue(string(*(tgwRes.UplinkZoneState.Get()))),
	}
}

func MapToTgwV2(tgwRes vpc.TransitGatewayV1Dot2) Tgw {
	return Tgw{
		Id:                      types.StringValue(tgwRes.Id),
		Description:             types.StringPointerValue(tgwRes.Description.Get()),
		Name:                    types.StringValue(tgwRes.Name),
		AccountId:               types.StringValue(tgwRes.AccountId),
		Bandwidth:               types.Int32PointerValue(tgwRes.Bandwidth.Get()),
		CreatedAt:               types.StringValue(tgwRes.CreatedAt.Format(time.RFC3339)),
		CreatedBy:               types.StringValue(tgwRes.CreatedBy),
		FirewallConnectionState: types.StringValue(string(*(tgwRes.FirewallConnectionState.Get()))),
		FirewallIds:             types.StringPointerValue(tgwRes.FirewallIds.Get()),
		ModifiedAt:              types.StringValue(tgwRes.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:              types.StringValue(tgwRes.ModifiedBy),
		State:                   types.StringValue(string(tgwRes.State)),
		UplinkEnabled:           types.BoolValue(tgwRes.GetUplinkEnabled()),
	}
}

func (m Tgw) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                        types.StringType,
		"description":               types.StringType,
		"name":                      types.StringType,
		"account_id":                types.StringType,
		"bandwidth":                 types.Int32Type,
		"created_at":                types.StringType,
		"created_by":                types.StringType,
		"firewall_ids":              types.StringType,
		"firewall_connection_state": types.StringType,
		"modified_at":               types.StringType,
		"modified_by":               types.StringType,
		"state":                     types.StringType,
		"uplink_enabled":            types.BoolType,
		"uplink_active_zone":        types.StringType,
		"uplink_standby_zone":       types.StringType,
		"uplink_zone_state":         types.StringType,
	}
}
