package multinodegpucluster

import (
	"context"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
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
	Zone              types.String    `tfsdk:"zone"`
}

type GpuNodeImageList struct {
	RegionId types.String   `tfsdk:"region_id"`
	Images   []GpuNodeImage `tfsdk:"images"`
}

type GpuNodeImage struct {
	CreatedAt    types.String `tfsdk:"created_at"`
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	OsDistro     types.String `tfsdk:"os_distro"`
	Priority     types.String `tfsdk:"priority"`
	ScpImageType types.String `tfsdk:"scp_image_type"`
	ScpOsVersion types.String `tfsdk:"scp_os_version"`
}

type GpuNodeProductList struct {
	Type     types.String     `tfsdk:"type"`
	ImageId  types.String     `tfsdk:"image_id"`
	Products []GpuNodeProduct `tfsdk:"products"`
}

type GpuNodeProduct struct {
	CreatedAt    types.String        `tfsdk:"created_at"`
	CreatedBy    types.String        `tfsdk:"created_by"`
	Description  types.String        `tfsdk:"description"`
	Id           types.String        `tfsdk:"id"`
	ModifiedAt   types.String        `tfsdk:"modified_at"`
	ModifiedBy   types.String        `tfsdk:"modified_by"`
	Name         types.String        `tfsdk:"name"`
	ProductAttrs GpuNodeProductAttrs `tfsdk:"product_attrs"`
	State        types.String        `tfsdk:"state"`
	Type         types.String        `tfsdk:"type"`
}

type GpuNodeProductAttrs struct {
	ComputeClassTypeName  types.String `tfsdk:"compute_class_type_name"`
	ComputeClassTypeValue types.String `tfsdk:"compute_class_type_value"`
	CpuValue              types.String `tfsdk:"cpu_value"`
	DiskUnit              types.String `tfsdk:"disk_unit"`
	DiskValue             types.String `tfsdk:"disk_value"`
	GpuModel              types.String `tfsdk:"gpu_model"`
	MemoryValue           types.String `tfsdk:"memory_value"`
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
	Zone              types.String `tfsdk:"zone"`
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

type GpuNodePublicNatIpResource struct {
	GpuNodeId types.String `tfsdk:"gpu_node_id"`
	PublicIpAddressId types.String `tfsdk:"public_ip_address_id"`
	PolicyNat types.String `tfsdk:"policy_nat"`
	PolicyUseNat types.Bool `tfsdk:"policy_use_nat"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
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
	Zone         types.String `tfsdk:"zone"`
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
		"zone":           types.StringType,
	}
}

type ServerDetailsValueType struct {
	basetypes.ObjectType
}

type ClusterFabricList struct {
	ClusterFabricName types.String    `tfsdk:"cluster_fabric_name"`
	State             types.String    `tfsdk:"state"`
	NodePoolId        types.String    `tfsdk:"node_pool_id"`
	Filter            []filter.Filter `tfsdk:"filter"`
	Ids               []types.String  `tfsdk:"ids"`
}

type ClusterFabricDataSource struct {
	AccountId       types.String `tfsdk:"account_id"`
	ClusterName     types.String `tfsdk:"cluster_name"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	Description     types.String `tfsdk:"description"`
	GpuNodeDetails  types.List   `tfsdk:"gpu_node_details"`
	Id              types.String `tfsdk:"id"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	NodePoolId      types.String `tfsdk:"node_pool_id"`
	PirpId          types.String `tfsdk:"pirp_id"`
	ProductId       types.String `tfsdk:"product_id"`
	RegionId        types.String `tfsdk:"region_id"`
	ServerType      types.String `tfsdk:"server_type"`
	State           types.String `tfsdk:"state"`
	UsedServerCount types.Int64  `tfsdk:"used_server_count"`
}

type GpuNodeDetailsValue struct {
	GpuNodeId     types.String `tfsdk:"gpu_node_id"`
	GpuNodeName   types.String `tfsdk:"gpu_node_name"`
	PolicyIp      types.String `tfsdk:"policy_ip"`
	ProductTypeId types.String `tfsdk:"product_type_id"`
	ServerType    types.String `tfsdk:"server_type"`
	State         types.String `tfsdk:"state"`
}

func (v GpuNodeDetailsValue) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"gpu_node_id":     types.StringType,
		"gpu_node_name":   types.StringType,
		"policy_ip":       types.StringType,
		"product_type_id": types.StringType,
		"server_type":     types.StringType,
		"state":           types.StringType,
	}
}

type NodePoolList struct {
	SubnetId        types.String   `tfsdk:"subnet_id"`
	ClusterFabricId types.String   `tfsdk:"cluster_fabric_id"`
	NodePoolId      types.String   `tfsdk:"node_pool_id"`
	Zone            types.String   `tfsdk:"zone"`
	Ids             []types.String `tfsdk:"ids"`
}

type ClusterFabricMember struct {
	BeforeClusterFabricId types.String   `tfsdk:"before_cluster_fabric_id"`
	AfterClusterFabricId  types.String   `tfsdk:"after_cluster_fabric_id"`
	GpuNodeIdList         []types.String `tfsdk:"gpu_node_id_list"`
	Id                    types.String   `tfsdk:"id"`
}
