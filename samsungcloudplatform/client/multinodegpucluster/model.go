package multinodegpucluster

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const ServiceType = "scp-multinodegpucluster"

type GpuNodeList struct {
	ClusterFabricId   types.String    `tfsdk:"cluster_fabric_id"`
	ClusterFabricName types.String    `tfsdk:"cluster_fabric_name"`
	Filter            []filter.Filter `tfsdk:"filter"`
	GpuNodeName       types.String    `tfsdk:"gpu_node_name"`
	Ids               []types.String  `tfsdk:"ids"`
	Ip                types.String    `tfsdk:"ip"`
	State             types.String    `tfsdk:"state"`
	VpcId             types.String    `tfsdk:"vpc_id"`
}

type GpuNodeDataSource struct {
	AccountId         types.String `tfsdk:"account_id"`
	ClusterFabricId   types.String `tfsdk:"cluster_fabric_id"`
	ClusterFabricName types.String `tfsdk:"cluster_fabric_name"`
	CreatedAt         types.String `tfsdk:"created_at"`
	CreatedBy         types.String `tfsdk:"created_by"`
	GpuNodeName       types.String `tfsdk:"gpu_node_name"`
	Id                types.String `tfsdk:"id"`
	ImageId           types.String `tfsdk:"image_id"`
	ImageVersion      types.String `tfsdk:"image_version"`
	InitScript        types.String `tfsdk:"init_script"`
	LockEnabled       types.Bool   `tfsdk:"lock_enabled"`
	ModifiedAt        types.String `tfsdk:"modified_at"`
	ModifiedBy        types.String `tfsdk:"modified_by"`
	NetworkId         types.String `tfsdk:"network_id"`
	NodePoolId        types.String `tfsdk:"node_pool_id"`
	OsType            types.String `tfsdk:"os_type"`
	PfsIp             types.List   `tfsdk:"pfs_ip"`
	PolicyIp          types.String `tfsdk:"policy_ip"`
	PolicyNat         types.String `tfsdk:"policy_nat"`
	PolicyUseNat      types.Bool   `tfsdk:"policy_use_nat"`
	ProductTypeId     types.String `tfsdk:"product_type_id"`
	RegionId          types.String `tfsdk:"region_id"`
	RootAccount       types.String `tfsdk:"root_account"`
	ServerType        types.String `tfsdk:"server_type"`
	State             types.String `tfsdk:"state"`
	TimeZone          types.String `tfsdk:"time_zone"`
	VpcId             types.String `tfsdk:"vpc_id"`
}

type GpuNodeResource struct {
	AccountId            types.String              `tfsdk:"account_id"`
	ClusterFabricDetails ClusterFabricDetailsValue `tfsdk:"cluster_fabric_details"`
	CreatedAt            types.String              `tfsdk:"created_at"`
	CreatedBy            types.String              `tfsdk:"created_by"`
	GpuNodeNamePrefix    types.String              `tfsdk:"gpu_node_name_prefix"`
	ImageId              types.String              `tfsdk:"image_id"`
	ImageVersion         types.String              `tfsdk:"image_version"`
	InitScript           types.String              `tfsdk:"init_script"`
	LockEnabled          types.Bool                `tfsdk:"lock_enabled"`
	ModifiedAt           types.String              `tfsdk:"modified_at"`
	ModifiedBy           types.String              `tfsdk:"modified_by"`
	OsType               types.String              `tfsdk:"os_type"`
	OsUserId             types.String              `tfsdk:"os_user_id"`
	OsUserPassword       types.String              `tfsdk:"os_user_password"`
	ProductTypeId        types.String              `tfsdk:"product_type_id"`
	RegionId             types.String              `tfsdk:"region_id"`
	RootAccount          types.String              `tfsdk:"root_account"`
	ServerDetails        types.List                `tfsdk:"server_details"`
	ServerTypeId         types.String              `tfsdk:"server_type_id"`
	SubnetId             types.String              `tfsdk:"subnet_id"`
	Tags                 types.Map                 `tfsdk:"tags"`
	Timeouts             timeouts.Value            `tfsdk:"timeouts"`
	TimeZone             types.String              `tfsdk:"time_zone"`
	VpcId                types.String              `tfsdk:"vpc_id"`
}

type ClusterFabricDetailsValue struct {
	ClusterFabricId   basetypes.StringValue `tfsdk:"cluster_fabric_id"`
	ClusterFabricName basetypes.StringValue `tfsdk:"cluster_fabric_name"`
	NodePoolId        basetypes.StringValue `tfsdk:"node_pool_id"`
}

func (v ClusterFabricDetailsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"cluster_fabric_id":   basetypes.StringType{},
		"cluster_fabric_name": basetypes.StringType{},
		"node_pool_id":        basetypes.StringType{},
	}
}

type ClusterFabricDetailsType struct {
	basetypes.ObjectType
}

type ServerDetailsValue struct {
	GpuNodeName  types.String `tfsdk:"gpu_node_name"`
	Id           types.String `tfsdk:"id"`
	IpAddress    types.String `tfsdk:"ip_address"`
	NatEnabled   types.Bool   `tfsdk:"nat_enabled"`
	PfsIp        types.List   `tfsdk:"pfs_ip"`
	PolicyIp     types.String `tfsdk:"policy_ip"`
	PolicyNat    types.String `tfsdk:"policy_nat"`
	PolicyUseNat types.Bool   `tfsdk:"policy_use_nat"`
	ServerType   types.String `tfsdk:"server_type"`
	State        types.String `tfsdk:"state"`
}

func (v ServerDetailsValue) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"gpu_node_name": types.StringType,
		"id":            types.StringType,
		"ip_address":    types.StringType,
		"nat_enabled":   types.BoolType,
		"pfs_ip": types.ListType{
			ElemType: types.StringType,
		},
		"policy_ip":      types.StringType,
		"policy_nat":     types.StringType,
		"policy_use_nat": types.BoolType,
		"server_type":    types.StringType,
		"state":          types.StringType,
	}
}

type ServerDetailsValueType struct {
	basetypes.ObjectType
}

