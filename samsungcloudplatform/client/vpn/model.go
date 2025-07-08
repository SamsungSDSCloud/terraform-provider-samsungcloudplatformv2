package vpn

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-vpn"

type VpnGatewayDataSource struct { // Vpn Gateway List request 모델을 참고하여 구조체를 구성한다.
	Id         types.String `tfsdk:"id"`
	VpnGateway types.Object `tfsdk:"vpn_gateway"`
}

type VpnGatewayDataSourceIds struct { // Vpn Gateway List request 모델을 참고하여 구조체를 구성한다.
	Size      types.Int32    `tfsdk:"size"`
	Page      types.Int32    `tfsdk:"page"`
	Sort      types.String   `tfsdk:"sort"`
	Name      types.String   `tfsdk:"name"`
	IpAddress types.String   `tfsdk:"ip_address"`
	VpcId     types.String   `tfsdk:"vpc_id"`
	VpcName   types.String   `tfsdk:"vpc_name"`
	Ids       []types.String `tfsdk:"ids"`
}

type VpnGatewayResource struct { // Id + VpnGatewayCreateRequest
	Id          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	IpAddress   types.String `tfsdk:"ip_address"`
	IpId        types.String `tfsdk:"ip_id"`
	IpType      types.String `tfsdk:"ip_type"`
	Name        types.String `tfsdk:"name"`
	VpcId       types.String `tfsdk:"vpc_id"`
	Tags        types.Map    `tfsdk:"tags"`
	VpnGateway  types.Object `tfsdk:"vpn_gateway"`
}

type VpnGateway struct { // VpnGateway
	AccountId   types.String `tfsdk:"account_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	IpAddress   types.String `tfsdk:"ip_address"`
	IpId        types.String `tfsdk:"ip_id"`
	IpType      types.String `tfsdk:"ip_type"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	Name        types.String `tfsdk:"name"`
	State       types.String `tfsdk:"state"`
	VpcId       types.String `tfsdk:"vpc_id"`
	VpcName     types.String `tfsdk:"vpc_name"`
}

func (m VpnGateway) AttributeTypes() map[string]attr.Type { // VpnGateway 의 AttributeTypes 메서드를 추가한다.
	return map[string]attr.Type{
		"account_id":  types.StringType,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"description": types.StringType,
		"id":          types.StringType,
		"ip_address":  types.StringType,
		"ip_id":       types.StringType,
		"ip_type":     types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"name":        types.StringType,
		"state":       types.StringType,
		"vpc_id":      types.StringType,
		"vpc_name":    types.StringType,
	}
}

type VpnTunnelDataSource struct { // Vpn Tunnel List request 모델을 참고하여 구조체를 구성한다.
	Id             types.String `tfsdk:"id"`
	VpnTunnel      types.Object `tfsdk:"vpn_tunnel"`
}

type VpnTunnelDataSourceIds struct { // Vpn Tunnel List request 모델을 참고하여 구조체를 구성한다.
	Size           types.Int32    `tfsdk:"size"` // schema mapping
	Page           types.Int32    `tfsdk:"page"`
	Sort           types.String   `tfsdk:"sort"`
	Name           types.String   `tfsdk:"name"`
	VpnGatewayId   types.String   `tfsdk:"vpn_gateway_id"`
	VpnGatewayName types.String   `tfsdk:"vpn_gateway_name"`
	PeerGatewayIp  types.String   `tfsdk:"peer_gateway_ip"`
	Remote_subnet  types.String   `tfsdk:"remote_subnet"`
	Ids            []types.String `tfsdk:"ids"`
}

type VpnTunnelResource struct { // Id + VpnGatewayCreateRequest
	Id           types.String    `tfsdk:"id"`
	Description  types.String    `tfsdk:"description"`
	Name         types.String    `tfsdk:"name"`
	Phase1       VpnPhase1Detail `tfsdk:"phase1"`
	Phase2       VpnPhase2Detail `tfsdk:"phase2"`
	VpnGatewayId types.String    `tfsdk:"vpn_gateway_id"`
	Tags         types.Map       `tfsdk:"tags"`
	VpnTunnel    types.Object    `tfsdk:"vpn_tunnel"`
}

type VpnTunnel struct { // VpnTunnel
	AccountId           types.String    `tfsdk:"account_id"`
	CreatedAt           types.String    `tfsdk:"created_at"`
	CreatedBy           types.String    `tfsdk:"created_by"`
	Description         types.String    `tfsdk:"description"`
	Id                  types.String    `tfsdk:"id"`
	ModifiedAt          types.String    `tfsdk:"modified_at"`
	ModifiedBy          types.String    `tfsdk:"modified_by"`
	Name                types.String    `tfsdk:"name"`
	Phase1              VpnPhase1Detail `tfsdk:"phase1"`
	Phase2              VpnPhase2Detail `tfsdk:"phase2"`
	State               types.String    `tfsdk:"state"`
	VpcId               types.String    `tfsdk:"vpc_id"`
	VpcName             types.String    `tfsdk:"vpc_name"`
	VpnGatewayId        types.String    `tfsdk:"vpn_gateway_id"`
	VpnGatewayIpAddress types.String    `tfsdk:"vpn_gateway_ip_address"`
	VpnGatewayName      types.String    `tfsdk:"vpn_gateway_name"`
}

func (m VpnTunnel) AttributeTypes() map[string]attr.Type { // VpnTunnel 의 AttributeTypes 메서드를 추가한다.
	return map[string]attr.Type{
		"account_id":  types.StringType,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"description": types.StringType,
		"id":          types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"name":        types.StringType,
		"phase1": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"diffie_hellman_groups": types.ListType{
					ElemType: types.Int32Type,
				},
				"encryptions": types.ListType{
					ElemType: types.StringType,
				},
				"dpd_retry_interval": types.Int32Type,
				"ike_version":        types.Int32Type,
				"life_time":          types.Int32Type,
				"peer_gateway_ip":    types.StringType,
				"pre_shared_key":     types.StringType,
			},
		},
		"phase2": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"diffie_hellman_groups": types.ListType{
					ElemType: types.Int32Type,
				},
				"encryptions": types.ListType{
					ElemType: types.StringType,
				},
				"life_time":               types.Int32Type,
				"perfect_forward_secrecy": types.StringType,
				"remote_subnet":           types.StringType,
			},
		},
		"state":                  types.StringType,
		"vpc_id":                 types.StringType,
		"vpc_name":               types.StringType,
		"vpn_gateway_id":         types.StringType,
		"vpn_gateway_ip_address": types.StringType,
		"vpn_gateway_name":       types.StringType,
	}
}

type VpnPhase1Detail struct { // VpnPhase1Detail
	DiffieHellmanGroups []types.Int32  `tfsdk:"diffie_hellman_groups"`
	DpdRetryInterval    types.Int32    `tfsdk:"dpd_retry_interval"`
	Encryptions         []types.String `tfsdk:"encryptions"`
	IkeVersion          types.Int32    `tfsdk:"ike_version"`
	LifeTime            types.Int32    `tfsdk:"life_time"`
	PeerGatewayIp       types.String   `tfsdk:"peer_gateway_ip"`
	PreSharedKey        types.String   `tfsdk:"pre_shared_key"`
}

type VpnPhase2Detail struct { // VpnPhase2Detail
	DiffieHellmanGroups   []types.Int32  `tfsdk:"diffie_hellman_groups"`
	Encryptions           []types.String `tfsdk:"encryptions"`
	LifeTime              types.Int32    `tfsdk:"life_time"`
	PerfectForwardSecrecy types.String   `tfsdk:"perfect_forward_secrecy"`
	RemoteSubnet          types.String   `tfsdk:"remote_subnet"`
}
