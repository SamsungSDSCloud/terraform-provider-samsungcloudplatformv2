package backup

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-backup"

type BackupDataSourceIds struct {
	Name       types.String    `tfsdk:"name"`
	ServerName types.String    `tfsdk:"server_name"`
	Filter     []filter.Filter `tfsdk:"filter"`
	Ids        []types.String  `tfsdk:"ids"`
	Region     types.String    `tfsdk:"region"`
}

type BackupDataSource struct {
	Id         types.String    `tfsdk:"id"`
	Name       types.String    `tfsdk:"name"`
	ServerName types.String    `tfsdk:"server_name"`
	Filter     []filter.Filter `tfsdk:"filter"`
	Backup     types.Object    `tfsdk:"backup"`
	Region     types.String    `tfsdk:"region"`
}

type BackupResource struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	PolicyCategory  types.String `tfsdk:"policy_category"`
	PolicyType      types.String `tfsdk:"policy_type"`
	ServerUuid      types.String `tfsdk:"server_uuid"`
	ServerCategory  types.String `tfsdk:"server_category"`
	EncryptEnabled  types.Bool   `tfsdk:"encrypt_enabled"`
	RetentionPeriod types.String `tfsdk:"retention_period"`
	Schedules       []Schedule   `tfsdk:"schedules"`
	Tags            types.Map    `tfsdk:"tags"`
	Region          types.String `tfsdk:"region"`
}

type Schedule struct {
	Frequency types.String `tfsdk:"frequency"`
	StartDay  types.String `tfsdk:"start_day"`
	StartTime types.String `tfsdk:"start_time"`
	StartWeek types.String `tfsdk:"start_week"`
	Type      types.String `tfsdk:"type"`
}

type Backup struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	PolicyType      types.String `tfsdk:"policy_type"`
	RetentionPeriod types.String `tfsdk:"retention_period"`
	RoleType        types.String `tfsdk:"role_type"`
	ServerName      types.String `tfsdk:"server_name"`
	State           types.String `tfsdk:"state"`
	EncryptEnabled  types.Bool   `tfsdk:"encrypt_enabled"`
	PolicyCategory  types.String `tfsdk:"policy_category"`
	ServerCategory  types.String `tfsdk:"server_category"`
	ServerUuid      types.String `tfsdk:"server_uuid"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

func (m Backup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"name":             types.StringType,
		"policy_type":      types.StringType,
		"retention_period": types.StringType,
		"role_type":        types.StringType,
		"server_name":      types.StringType,
		"state":            types.StringType,
		"encrypt_enabled":  types.BoolType,
		"policy_category":  types.StringType,
		"server_category":  types.StringType,
		"server_uuid":      types.StringType,
		"created_at":       types.StringType,
		"created_by":       types.StringType,
		"modified_at":      types.StringType,
		"modified_by":      types.StringType,
	}
}

type Tag struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// -------------------- Handler -------------------- //

type UpdateHandler struct {
	Fields  []string
	Handler func(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error
}
