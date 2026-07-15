package eventstreams

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/database"
)

const ServiceType = "scp-eventstreams"

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
	AkhqEnabled               types.Bool         `tfsdk:"akhq_enabled"`
	AllowableIpAddresses      types.Set          `tfsdk:"allowable_ip_addresses"`
	DbaasEngineVersionId      types.String       `tfsdk:"dbaas_engine_version_id"`
	IsCombined                types.Bool         `tfsdk:"is_combined"`
	InitConfigOption          *InitConfigOption  `tfsdk:"init_config_option"`
	InstanceGroups            types.List         `tfsdk:"instance_groups"`
	InstanceNamePrefix        types.String       `tfsdk:"instance_name_prefix"`
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
	AccountId     types.String `tfsdk:"account_id"`
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

// InitConfigOptionBase holds the API-returned init config fields (used by the data source).
type InitConfigOptionBase struct {
	BrokerPort    types.Int32 `tfsdk:"broker_port"`
	ZookeeperPort types.Int32 `tfsdk:"zookeeper_port"`
}

type InitConfigOption struct {
	InitConfigOptionBase
	AkhqId                types.String `tfsdk:"akhq_id"`
	AkhqPassword          types.String `tfsdk:"akhq_password"`
	BrokerSaslId          types.String `tfsdk:"broker_sasl_id"`
	BrokerSaslPassword    types.String `tfsdk:"broker_sasl_password"`
	ZookeeperSaslId       types.String `tfsdk:"zookeeper_sasl_id"`
	ZookeeperSaslPassword types.String `tfsdk:"zookeeper_sasl_password"`
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
	DbaasEngine               types.String             `tfsdk:"dbaas_engine"`
	IsCombined                types.Bool               `tfsdk:"is_combined"`
	Id                        types.String             `tfsdk:"id"`
	InitConfigOption          *InitConfigOptionBase    `tfsdk:"init_config_option"`
	InstanceCount             types.Int32              `tfsdk:"instance_count"`
	InstanceGroups            []database.InstanceGroup `tfsdk:"instance_groups"`
	MaintenanceOption         *MaintenanceOption       `tfsdk:"maintenance_option"`
	Name                      types.String             `tfsdk:"name"`
	NatEnabled                types.Bool               `tfsdk:"nat_enabled"`
	ProductType               types.String             `tfsdk:"product_type"`
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
