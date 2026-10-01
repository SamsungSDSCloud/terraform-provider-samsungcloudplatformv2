package smp

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-smp"

// ----------------------------
// Workspace
// ----------------------------

type WorkspaceResource struct {
	Name                  types.String `tfsdk:"name"`
	PrivateAclEnabled     types.Bool   `tfsdk:"private_acl_enabled"`
	PrivateAclResources   types.List   `tfsdk:"private_acl_resources"`
	Workspace             types.Object `tfsdk:"workspace"`
	WorkspaceId           types.String `tfsdk:"workspace_id"`
}

type WorkspaceDataSource struct {
	Workspace   types.Object `tfsdk:"workspace"`
	WorkspaceId types.String `tfsdk:"workspace_id"`
}

type WorkspaceDataSources struct {
	Name       types.String `tfsdk:"name"`
	NameLike   types.String `tfsdk:"name_like"`
	Id         types.String `tfsdk:"id"`
	IdLike     types.String `tfsdk:"id_like"`
	Page       types.Int64  `tfsdk:"page"`
	Size       types.Int64  `tfsdk:"size"`
	SortBy     types.String `tfsdk:"sort_by"`
	SortOrder  types.String `tfsdk:"sort_order"`
	TotalCount types.Int64  `tfsdk:"total_count"`
	Sort       types.List   `tfsdk:"sort"`
	Workspaces types.List   `tfsdk:"workspaces"`
}

type PrivateAclResource struct {
	ResourceId   types.String `tfsdk:"resource_id"`
	ResourceName types.String `tfsdk:"resource_name"`
	ResourceType types.String `tfsdk:"resource_type"`
}

func (m PrivateAclResource) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"resource_id":   types.StringType,
		"resource_name": types.StringType,
		"resource_type": types.StringType,
	}
}

type PrivateAclResourcesSetRequest struct {
	PrivateAclEnabled     types.Bool   `tfsdk:"private_acl_enabled"`
	PrivateAclResources   types.List   `tfsdk:"private_acl_resources"`
}

type Workspace struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	State                     types.String `tfsdk:"state"`
	AccountId                 types.String `tfsdk:"account_id"`
	CreatedAt                 types.String `tfsdk:"created_at"`
	CreatedBy                 types.String `tfsdk:"created_by"`
	ModifiedAt                types.String `tfsdk:"modified_at"`
	ModifiedBy                types.String `tfsdk:"modified_by"`
	PrivateAclEnabled         types.Bool   `tfsdk:"private_acl_enabled"`
	PrivateAclResources       types.List   `tfsdk:"private_acl_resources"`
	PrivatePrometheusEndpoint types.String `tfsdk:"private_prometheus_endpoint"`
}

func (m Workspace) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                          types.StringType,
		"name":                        types.StringType,
		"state":                       types.StringType,
		"account_id":                  types.StringType,
		"created_at":                  types.StringType,
		"created_by":                  types.StringType,
		"modified_at":                 types.StringType,
		"modified_by":                 types.StringType,
		"private_acl_enabled":         types.BoolType,
		"private_acl_resources": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: PrivateAclResource{}.AttributeTypes(),
			},
		},
		"private_prometheus_endpoint": types.StringType,
	}
}

// ----------------------------
// RuleNamespace
// ----------------------------

type RuleNamespaceResource struct {
	Name            types.String `tfsdk:"name"`
	ConfigData      types.String `tfsdk:"config_data"`
	RuleNamespace   types.Object `tfsdk:"rule_namespace"`
	RuleNamespaceId types.String `tfsdk:"rule_namespace_id"`
	WorkspaceId     types.String `tfsdk:"workspace_id"`
}

type RuleNamespaceDataSource struct {
	RuleNamespace   types.Object `tfsdk:"rule_namespace"`
	RuleNamespaceId types.String `tfsdk:"rule_namespace_id"`
	WorkspaceId     types.String `tfsdk:"workspace_id"`
}

type RuleNamespaceDataSources struct {
	Name           types.String `tfsdk:"name"`
	NameLike       types.String `tfsdk:"name_like"`
	Id             types.String `tfsdk:"id"`
	IdLike         types.String `tfsdk:"id_like"`
	Page           types.Int64  `tfsdk:"page"`
	Size           types.Int64  `tfsdk:"size"`
	SortBy         types.String `tfsdk:"sort_by"`
	SortOrder      types.String `tfsdk:"sort_order"`
	WorkspaceId    types.String `tfsdk:"workspace_id"`
	TotalCount types.Int64  `tfsdk:"total_count"`
	Sort       types.List   `tfsdk:"sort"`
	RuleNamespaces types.List   `tfsdk:"rule_namespaces"`
}

type RuleNamespace struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	ConfigData  types.String `tfsdk:"config_data"`
	State       types.String `tfsdk:"state"`
	WorkspaceId types.String `tfsdk:"workspace_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
}

func (m RuleNamespace) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           types.StringType,
		"name":         types.StringType,
		"config_data":  types.StringType,
		"state":        types.StringType,
		"workspace_id": types.StringType,
		"created_at":   types.StringType,
		"created_by":   types.StringType,
		"modified_at":  types.StringType,
		"modified_by":  types.StringType,
	}
}

// ----------------------------
// AlertManager
// ----------------------------

type AlertManagerResource struct {
	AlertManager        types.Object `tfsdk:"alert_manager"`
	ConfigData          types.String `tfsdk:"config_data"`
	NotificationGroupId types.String `tfsdk:"notification_group_id"`
	WorkspaceId         types.String `tfsdk:"workspace_id"`
}

type AlertManagerDataSource struct {
	AlertManager types.Object `tfsdk:"alert_manager"`
	WorkspaceId  types.String `tfsdk:"workspace_id"`
}

type AlertManagerDataSources struct {
	WorkspaceId   types.String `tfsdk:"workspace_id"`
	AlertManagers types.List   `tfsdk:"alert_managers"`
}

type AlertManager struct {
	Id                  types.String `tfsdk:"id"`
	ConfigData          types.String `tfsdk:"config_data"`
	State               types.String `tfsdk:"state"`
	WorkspaceId         types.String `tfsdk:"workspace_id"`
	NotificationGroupId types.String `tfsdk:"notification_group_id"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
}

func (m AlertManager) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    types.StringType,
		"config_data":           types.StringType,
		"state":                 types.StringType,
		"workspace_id":          types.StringType,
		"notification_group_id": types.StringType,
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"modified_at":           types.StringType,
		"modified_by":           types.StringType,
	}
}

// ----------------------------
// WorkspaceConfiguration
// ----------------------------

type WorkspaceConfigurationResource struct {
	WorkspaceConfiguration types.Object `tfsdk:"workspace_configuration"`
	WorkspaceId            types.String `tfsdk:"workspace_id"`
	RetentionPeriod        types.Int64  `tfsdk:"retention_period"`
}

type WorkspaceConfigurationDataSource struct {
	WorkspaceConfiguration types.Object `tfsdk:"workspace_configuration"`
	WorkspaceId            types.String `tfsdk:"workspace_id"`
}

type WorkspaceConfigurationDataSources struct {
	WorkspaceId             types.String `tfsdk:"workspace_id"`
	WorkspaceConfigurations types.List   `tfsdk:"workspace_configurations"`
}

type WorkspaceConfiguration struct {
	Id              types.String `tfsdk:"id"`
	State           types.String `tfsdk:"state"`
	WorkspaceId     types.String `tfsdk:"workspace_id"`
	RetentionPeriod types.Int64  `tfsdk:"retention_period"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	ModifiedAt      types.String `tfsdk:"modified_at"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

func (m WorkspaceConfiguration) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"state":            types.StringType,
		"workspace_id":     types.StringType,
		"retention_period": types.Int64Type,
		"created_at":       types.StringType,
		"created_by":       types.StringType,
		"modified_at":      types.StringType,
		"modified_by":      types.StringType,
	}
}

// ----------------------------
// NotificationGroup
// ----------------------------

type NotificationGroupResource struct {
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	NotificationGroup   types.Object `tfsdk:"notification_group"`
	NotificationGroupId types.String `tfsdk:"notification_group_id"`
	RecipientUserIds    types.List   `tfsdk:"recipient_user_ids"`
}

type NotificationGroupDataSource struct {
	NotificationGroup   types.Object `tfsdk:"notification_group"`
	NotificationGroupId types.String `tfsdk:"notification_group_id"`
}

type NotificationGroupDataSources struct {
	Name               types.String `tfsdk:"name"`
	NameLike           types.String `tfsdk:"name_like"`
	Id                 types.String `tfsdk:"id"`
	IdLike             types.String `tfsdk:"id_like"`
	Page               types.Int64  `tfsdk:"page"`
	Size               types.Int64  `tfsdk:"size"`
	SortBy             types.String `tfsdk:"sort_by"`
	SortOrder          types.String `tfsdk:"sort_order"`
	TotalCount types.Int64  `tfsdk:"total_count"`
	Sort       types.List   `tfsdk:"sort"`
	NotificationGroups types.List   `tfsdk:"notification_groups"`
}

type NotificationGroup struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	AccountId   types.String `tfsdk:"account_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	Recipients  types.List   `tfsdk:"recipients"`
}

func (m NotificationGroup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"account_id":  types.StringType,
		"created_at":  types.StringType,
		"created_by":  types.StringType,
		"modified_at": types.StringType,
		"modified_by": types.StringType,
		"recipients": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":                types.StringType,
					"recipient_user_id": types.StringType,
				},
			},
		},
	}
}
