package baremetal

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const ServiceType = "scp-baremetal"

type BaremetalDataSourceIds struct {
	Ids        []types.String  `tfsdk:"ids"`
	Ip         types.String    `tfsdk:"policy_ip"`
	ServerName types.String    `tfsdk:"server_name"`
	State      types.String    `tfsdk:"state"`
	VpcId      types.String    `tfsdk:"vpc_id"`
	Filter     []filter.Filter `tfsdk:"filter"`
}

type BaremetalDataSource struct {
	AccountId          types.String `tfsdk:"account_id"`
	CreatedAt          types.String `tfsdk:"created_at"`
	CreatedBy          types.String `tfsdk:"created_by"`
	HyperThreadingUse  types.Bool   `tfsdk:"hyper_threading_use"`
	Id                 types.String `tfsdk:"id"`
	ImageId            types.String `tfsdk:"image_id"`
	ImageVersion       types.String `tfsdk:"image_version"`
	InitScript         types.String `tfsdk:"init_script"`
	LocalSubnetInfo    types.List   `tfsdk:"local_subnet_info"`
	LockEnabled        types.Bool   `tfsdk:"lock_enabled"`
	ModifiedAt         types.String `tfsdk:"modified_at"`
	ModifiedBy         types.String `tfsdk:"modified_by"`
	NetworkId          types.String `tfsdk:"network_id"`
	OsType             types.String `tfsdk:"os_type"`
	PlacementGroupName types.String `tfsdk:"placement_group_name"`
	PolicyIp           types.String `tfsdk:"policy_ip"`
	PrivateNatInfo     types.Object `tfsdk:"private_nat_info"`
	ProductTypeId      types.String `tfsdk:"product_type_id"`
	PublicNatInfo      types.Object `tfsdk:"public_nat_info"`
	RegionId           types.String `tfsdk:"region_id"`
	RootAccount        types.String `tfsdk:"root_account"`
	ServerName         types.String `tfsdk:"server_name"`
	ServerType         types.String `tfsdk:"server_type"`
	State              types.String `tfsdk:"state"`
	TimeZone           types.String `tfsdk:"time_zone"`
	UseLocalSubnet     types.Bool   `tfsdk:"use_local_subnet"`
	VpcId              types.String `tfsdk:"vpc_id"`
}

type LocalSubnetInfoType struct {
	basetypes.ObjectType
}

type LocalSubnetInfo struct {
	InterfaceName       types.String `tfsdk:"interface_name"`
	LocalSubnetId       types.String `tfsdk:"local_subnet_id"`
	PolicyLocalSubnetIp types.String `tfsdk:"policy_local_subnet_ip"`
	State               types.String `tfsdk:"state"`
	VlanId              types.String `tfsdk:"vlan_id"`
	VniRoleName         types.String `tfsdk:"vni_role_name"`
	state               attr.ValueState
}

func (v LocalSubnetInfo) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface_name":         types.StringType,
		"local_subnet_id":        types.StringType,
		"policy_local_subnet_ip": types.StringType,
		"state":                  types.StringType,
		"vlan_id":                types.StringType,
		"vni_role_name":          types.StringType,
	}
}

type BaremetalResource struct {
	AccountId          types.String   `tfsdk:"account_id"`
	CreatedAt          types.String   `tfsdk:"created_at"`
	CreatedBy          types.String   `tfsdk:"created_by"`
	ImageId            types.String   `tfsdk:"image_id"`
	ImageVersion       types.String   `tfsdk:"image_version"`
	InitScript         types.String   `tfsdk:"init_script"`
	LockEnabled        types.Bool     `tfsdk:"lock_enabled"`
	ModifiedAt         types.String   `tfsdk:"modified_at"`
	ModifiedBy         types.String   `tfsdk:"modified_by"`
	NetworkId          types.String   `tfsdk:"network_id"`
	OsType             types.String   `tfsdk:"os_type"`
	OsUserId           types.String   `tfsdk:"os_user_id"`
	OsUserPassword     types.String   `tfsdk:"os_user_password"`
	PlacementGroupName types.String   `tfsdk:"placement_group_name"`
	RegionId           types.String   `tfsdk:"region_id"`
	RootAccount        types.String   `tfsdk:"root_account"`
	ServerDetails      types.List     `tfsdk:"server_details"`
	SubnetId           types.String   `tfsdk:"subnet_id"`
	TimeZone           types.String   `tfsdk:"time_zone"`
	UsePlacementGroup  types.Bool     `tfsdk:"use_placement_group"`
	VpcId              types.String   `tfsdk:"vpc_id"`
	Tags               types.Map      `tfsdk:"tags"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

type ServerDetailsType struct {
	basetypes.ObjectType
}

type ServerDetails struct {
	BareMetalLocalSubnetId        types.String `tfsdk:"bare_metal_local_subnet_id"`
	BareMetalLocalSubnetIpAddress types.String `tfsdk:"bare_metal_local_subnet_ip_address"`
	BareMetalServerName           types.String `tfsdk:"bare_metal_server_name"`
	IpAddress                     types.String `tfsdk:"ip_address"`
	NatEnabled                    types.Bool   `tfsdk:"nat_enabled"`
	PublicIpAddressId             types.String `tfsdk:"public_ip_address_id"`
	ServerTypeId                  types.String `tfsdk:"server_type_id"`
	UseHyperThreading             types.Bool   `tfsdk:"use_hyper_threading"`
	Id                            types.String `tfsdk:"id"`
	State                         types.String `tfsdk:"state"`
}

func (v ServerDetails) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"bare_metal_local_subnet_id":         types.StringType,
		"bare_metal_local_subnet_ip_address": types.StringType,
		"bare_metal_server_name":             types.StringType,
		"ip_address":                         types.StringType,
		"nat_enabled":                        types.BoolType,
		"public_ip_address_id":               types.StringType,
		"server_type_id":                     types.StringType,
		"use_hyper_threading":                types.BoolType,
		"id":                                 types.StringType,
		"state":                              types.StringType,
	}
}

type PublicNatInfoValue struct {
	NatId       types.String `tfsdk:"nat_id"`
	NatIp       types.String `tfsdk:"nat_ip"`
	NatIpId     types.String `tfsdk:"nat_ip_id"`
	State       types.String `tfsdk:"state"`
	StaticNatId types.String `tfsdk:"static_nat_id"`
}

func (v PublicNatInfoValue) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"nat_id":        types.StringType,
		"nat_ip":        types.StringType,
		"nat_ip_id":     types.StringType,
		"state":         types.StringType,
		"static_nat_id": types.StringType,
	}
}

type PrivateNatInfoValue struct {
	NatId       types.String `tfsdk:"nat_id"`
	NatIp       types.String `tfsdk:"nat_ip"`
	NatIpId     types.String `tfsdk:"nat_ip_id"`
	State       types.String `tfsdk:"state"`
	StaticNatId types.String `tfsdk:"static_nat_id"`
}

func (v PrivateNatInfoValue) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"nat_id":        types.StringType,
		"nat_ip":        types.StringType,
		"nat_ip_id":     types.StringType,
		"state":         types.StringType,
		"static_nat_id": types.StringType,
	}
}
