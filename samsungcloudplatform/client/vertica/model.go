package vertica

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/common/database"
)

const ServiceType = "scp-vertica"

// List Clusters의 Paramaters
type ClusterDataSource struct {
	Size         types.Int32  `tfsdk:"size"`
	Page         types.Int32  `tfsdk:"page"`
	Sort         types.String `tfsdk:"sort"`
	Name         types.String `tfsdk:"name"`
	ServiceState types.String `tfsdk:"service_state"`
	Clusters     []Cluster    `tfsdk:"clusters"`
}

type ClusterDataSourceDetail struct {
	Id            types.String   `tfsdk:"id"`
	ClusterDetail *ClusterDetail `tfsdk:"cluster"`
}

// Create Cluster의 Request
type ClusterResource struct {
	Id                        types.String       `tfsdk:"id"`
	AllowableIpAddresses      types.Set          `tfsdk:"allowable_ip_addresses"`
	DbaasEngineVersionId      types.String       `tfsdk:"dbaas_engine_version_id"`
	InitConfigOption          *InitConfigOption  `tfsdk:"init_config_option"`
	InstanceGroups            types.List         `tfsdk:"instance_groups"`
	InstanceNamePrefix        types.String       `tfsdk:"instance_name_prefix"`
	License                   types.String       `tfsdk:"license"`
	MaintenanceOption         *MaintenanceOption `tfsdk:"maintenance_option"`
	Name                      types.String       `tfsdk:"name"`
	NatEnabled                types.Bool         `tfsdk:"nat_enabled"`
	ServiceState              types.String       `tfsdk:"service_state"`
	SubnetId                  types.String       `tfsdk:"subnet_id"`
	Tags                      types.Map          `tfsdk:"tags"`
	Timezone                  types.String       `tfsdk:"timezone"`
	ServiceWatchLogCollection types.Bool         `tfsdk:"service_watch_log_collection"`
}

// List Clusters의 Response
type Cluster struct {
	AccountId       types.String `tfsdk:"account_id"`
	ConsoleIncluded types.Bool   `tfsdk:"console_included"`
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	DatabaseName    types.String `tfsdk:"database_name"`
	InstanceCount   types.Int32  `tfsdk:"instance_count"`
	RoleType        types.String `tfsdk:"role_type"`
	ServiceState    types.String `tfsdk:"service_state"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

// InitConfigOptionBase holds the fields the detail API returns. The data source
// embeds it directly; the resource extends it with write-only fields.
type InitConfigOptionBase struct {
	BackupOption     BackupOption `tfsdk:"backup_option"`
	DatabaseLocale   types.String `tfsdk:"database_locale"`
	DatabaseName     types.String `tfsdk:"database_name"`
	DatabasePort     types.Int32  `tfsdk:"database_port"`
	DatabaseUserName types.String `tfsdk:"database_user_name"`
	McPort           types.Int32  `tfsdk:"mc_port"`
}

// InitConfigOption extends the base with the write-only field
// database_user_password that the detail API does not return.
type InitConfigOption struct {
	InitConfigOptionBase
	DatabaseUserPassword types.String `tfsdk:"database_user_password"`
}

type BackupOption struct {
	RetentionPeriodDay types.String `tfsdk:"retention_period_day"`
	StartingTimeHour   types.String `tfsdk:"starting_time_hour"`
}

type MaintenanceOption struct {
	PeriodHour           types.String `tfsdk:"period_hour"`
	StartingDayOfWeek    types.String `tfsdk:"starting_day_of_week"`
	StartingTime         types.String `tfsdk:"starting_time"`
	UseMaintenanceOption types.Bool   `tfsdk:"use_maintenance_option"`
}

type ClusterDetail struct {
	AccountId                 types.String             `tfsdk:"account_id"`
	AllowableIpAddresses      types.Set                `tfsdk:"allowable_ip_addresses"`
	ConsoleIncluded           types.Bool               `tfsdk:"console_included"`
	NatEnabled                types.Bool               `tfsdk:"nat_enabled"`
	DbaasEngine               types.String             `tfsdk:"dbaas_engine"`
	DbaasEngineVersionName    types.String             `tfsdk:"dbaas_engine_version_name"`
	Id                        types.String             `tfsdk:"id"`
	InitConfigOption          *InitConfigOptionBase    `tfsdk:"init_config_option"`
	InstanceCount             types.Int32              `tfsdk:"instance_count"`
	InstanceGroups            []database.InstanceGroup `tfsdk:"instance_groups"`
	MaintenanceOption         *MaintenanceOption       `tfsdk:"maintenance_option"`
	Name                      types.String             `tfsdk:"name"`
	ProductType               types.String             `tfsdk:"product_type"`
	RoleType                  types.String             `tfsdk:"role_type"`
	ServiceState              types.String             `tfsdk:"service_state"`
	SoftwareVersion           types.String             `tfsdk:"software_version"`
	SubnetId                  types.String             `tfsdk:"subnet_id"`
	Timezone                  types.String             `tfsdk:"timezone"`
	CreatedAt                 types.String             `tfsdk:"created_at"`
	CreatedBy                 types.String             `tfsdk:"created_by"`
	ModifiedAt                types.String             `tfsdk:"modified_at"`
	ModifiedBy                types.String             `tfsdk:"modified_by"`
	ServiceWatchLogCollection types.Bool               `tfsdk:"service_watch_log_collection"`
}

// -------------------- Handler -------------------- //

type UpdateHandler struct {
	Fields  []string
	Handler func(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error
}

// --------------- Engine Version ------------ //

type EngineVersionDataSource struct {
	Contents []EngineVersion `tfsdk:"contents"`
}

type EngineVersion struct {
	EndOfService     types.Bool   `tfsdk:"end_of_service"`
	Id               types.String `tfsdk:"id"`
	MajorVersion     types.String `tfsdk:"major_version"`
	Name             types.String `tfsdk:"name"`
	OsType           types.String `tfsdk:"os_type"`
	OsVersion        types.String `tfsdk:"os_version"`
	ProductImageType types.String `tfsdk:"product_image_type"`
	SoftwareVersion  types.String `tfsdk:"software_version"`
}
