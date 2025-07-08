package mysql

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
)

const ServiceType = "scp-mysql"

type ClusterDataSource struct {
	Size         types.Int32  `tfsdk:"size"`
	Page         types.Int32  `tfsdk:"page"`
	Sort         types.String `tfsdk:"sort"`
	Name         types.String `tfsdk:"name"`
	ServiceState types.String `tfsdk:"service_state"`
	DatabaseName types.String `tfsdk:"database_name"`
	Clusters     []Cluster    `tfsdk:"clusters"`
}

type ClusterDataSourceDetail struct {
	Id            types.String `tfsdk:"id"`
	ClusterDetail types.Object `tfsdk:"cluster"`
}

type ClusterResource struct {
	Id                   types.String      `tfsdk:"id"`
	AllowableIpAddresses types.List        `tfsdk:"allowable_ip_addresses"`
	DbaasEngineVersionId types.String      `tfsdk:"dbaas_engine_version_id"`
	NatEnabled           types.Bool        `tfsdk:"nat_enabled"`
	HaEnabled            types.Bool        `tfsdk:"ha_enabled"`
	InitConfigOption     InitConfigOption  `tfsdk:"init_config_option"`
	InstanceGroups       []InstanceGroup   `tfsdk:"instance_groups"`
	InstanceNamePrefix   types.String      `tfsdk:"instance_name_prefix"`
	MaintenanceOption    MaintenanceOption `tfsdk:"maintenance_option"`
	Name                 types.String      `tfsdk:"name"`
	ServiceState         types.String      `tfsdk:"service_state"`
	SubnetId             types.String      `tfsdk:"subnet_id"`
	Tags                 types.Map         `tfsdk:"tags"`
	Timezone             types.String      `tfsdk:"timezone"`
	VipPublicIpId        types.String      `tfsdk:"vip_public_ip_id"`
	//VipPublicIpAddress   types.String      `tfsdk:"vip_public_ip_address"`
	VirtualIpAddress types.String `tfsdk:"virtual_ip_address"`
}

type Cluster struct {
	AccountId     types.String `tfsdk:"account_id"`
	DatabaseName  types.String `tfsdk:"database_name"`
	HaEnabled     types.Bool   `tfsdk:"ha_enabled"`
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	InstanceCount types.Int32  `tfsdk:"instance_count"`
	RoleType      types.String `tfsdk:"role_type"`
	ServiceState  types.String `tfsdk:"service_state"`
	CreatedAt     types.String `tfsdk:"created_at"`
	CreatedBy     types.String `tfsdk:"created_by"`
	ModifiedAt    types.String `tfsdk:"modified_at"`
	ModifiedBy    types.String `tfsdk:"modified_by"`
}

type InitConfigOption struct {
	BackupOption          BackupOption `tfsdk:"backup_option"`
	DatabaseCaseSensitive types.Bool   `tfsdk:"database_case_sensitive"`
	DatabaseCharacterSet  types.String `tfsdk:"database_character_set"`
	DatabaseName          types.String `tfsdk:"database_name"`
	DatabasePort          types.Int32  `tfsdk:"database_port"`
	DatabaseUserName      types.String `tfsdk:"database_user_name"`
	DatabaseUserPassword  types.String `tfsdk:"database_user_password"`
}

type BackupOption struct {
	ArchiveFrequencyMinute types.String `tfsdk:"archive_frequency_minute"`
	RetentionPeriodDay     types.String `tfsdk:"retention_period_day"`
	StartingTimeHour       types.String `tfsdk:"starting_time_hour"`
}

type InstanceGroup struct {
	BlockStorageGroups []BlockStorageGroup `tfsdk:"block_storage_groups"`
	Id                 types.String        `tfsdk:"id"`
	Instances          []Instance          `tfsdk:"instances"`
	RoleType           types.String        `tfsdk:"role_type"`
	ServerTypeName     types.String        `tfsdk:"server_type_name"`
}

type BlockStorageGroup struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	RoleType   types.String `tfsdk:"role_type"`
	SizeGb     types.Int32  `tfsdk:"size_gb"`
	VolumeType types.String `tfsdk:"volume_type"`
}

type Instance struct {
	Name             types.String `tfsdk:"name"`
	RoleType         types.String `tfsdk:"role_type"`
	ServiceIpAddress types.String `tfsdk:"service_ip_address"`
	PublicIpId       types.String `tfsdk:"public_ip_id"`
	//PublicIpAddress  types.String `tfsdk:"public_ip_address"`
}

type MaintenanceOption struct {
	PeriodHour           types.String `tfsdk:"period_hour"`
	StartingDayOfWeek    types.String `tfsdk:"starting_day_of_week"`
	StartingTime         types.String `tfsdk:"starting_time"`
	UseMaintenanceOption types.Bool   `tfsdk:"use_maintenance_option"`
}

type ClusterDetail struct {
	AccountId            types.String      `tfsdk:"account_id"`
	AllowableIpAddresses []types.String    `tfsdk:"allowable_ip_addresses"`
	DbaasEngine          types.String      `tfsdk:"dbaas_engine"`
	NatEnabled           types.Bool        `tfsdk:"nat_enabled"`
	HaEnabled            types.Bool        `tfsdk:"ha_enabled"`
	Id                   types.String      `tfsdk:"id"`
	InitConfigOption     InitConfigOption  `tfsdk:"init_config_option"`
	InstanceCount        types.Int32       `tfsdk:"instance_count"`
	InstanceGroups       []InstanceGroup   `tfsdk:"instance_groups"`
	MaintenanceOption    MaintenanceOption `tfsdk:"maintenance_option"`
	Name                 types.String      `tfsdk:"name"`
	OriginClusterId      types.String      `tfsdk:"origin_cluster_id"`
	ProductType          types.String      `tfsdk:"product_type"`
	Replicas             []types.String    `tfsdk:"replicas"`
	RoleType             types.String      `tfsdk:"role_type"`
	ServiceState         types.String      `tfsdk:"service_state"`
	SoftwareVersion      types.String      `tfsdk:"software_version"`
	SubnetId             types.String      `tfsdk:"subnet_id"`
	Timezone             types.String      `tfsdk:"timezone"`
	VipPublicIpId        types.String      `tfsdk:"vip_public_ip_id"`
	//VipPublicIpAddress   types.String      `tfsdk:"vip_public_ip_address"`
	VirtualIpAddress types.String `tfsdk:"virtual_ip_address"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	ModifiedAt       types.String `tfsdk:"modified_at"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
}

func (m ClusterDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id": types.StringType,
		"allowable_ip_addresses": types.ListType{
			ElemType: types.StringType,
		},
		"dbaas_engine": types.StringType,
		"nat_enabled":  types.BoolType,
		"ha_enabled":   types.BoolType,
		"id":           types.StringType,
		"init_config_option": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"database_character_set":  types.StringType,
				"database_case_sensitive": types.BoolType,
				"database_name":           types.StringType,
				"database_port":           types.Int32Type,
				"database_user_name":      types.StringType,
				"database_user_password":  types.StringType,
				"backup_option": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"archive_frequency_minute": types.StringType,
						"retention_period_day":     types.StringType,
						"starting_time_hour":       types.StringType,
					},
				},
			},
		},
		"instance_count": types.Int32Type,
		"instance_groups": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":               types.StringType,
					"role_type":        types.StringType,
					"server_type_name": types.StringType,
					"block_storage_groups": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"id":          types.StringType,
								"name":        types.StringType,
								"role_type":   types.StringType,
								"size_gb":     types.Int32Type,
								"volume_type": types.StringType,
							},
						},
					},
					"instances": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"name":               types.StringType,
								"role_type":          types.StringType,
								"service_ip_address": types.StringType,
								"public_ip_id":       types.StringType,
								//"public_ip_address":  types.StringType,
							},
						},
					},
				},
			},
		},
		"maintenance_option": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"period_hour":            types.StringType,
				"starting_day_of_week":   types.StringType,
				"starting_time":          types.StringType,
				"use_maintenance_option": types.BoolType,
			},
		},
		"name":              types.StringType,
		"origin_cluster_id": types.StringType,
		"product_type":      types.StringType,
		"replicas": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"replica": types.StringType,
				},
			},
		},
		"role_type":        types.StringType,
		"service_state":    types.StringType,
		"software_version": types.StringType,
		"subnet_id":        types.StringType,
		"timezone":         types.StringType,
		"vip_public_ip_id": types.StringType,
		//"vip_public_ip_address": types.StringType,
		"virtual_ip_address": types.StringType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
	}
}

// -------------------- Handler -------------------- //

type UpdateHandler struct {
	Fields  []string
	Handler func(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error
}
