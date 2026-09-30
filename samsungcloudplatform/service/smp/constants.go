package smp

import (
	"context"
	"time"

	smp "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/smp/1.0"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ----------------------------
// Time & formatting
// ----------------------------
const (
	TimeFormatDisplay = "2006-01-02 15:04:05"
)

// ----------------------------
// Error & log message templates
// ----------------------------
const (
	ErrUnexpectedConfigure = "Unexpected Data Source Configure Type"

	// ---- Workspace ----
	ErrCreateWorkspace = "Error creating Workspace"
	ErrReadWorkspace   = "Error Reading Workspace"
	ErrUpdateWorkspace = "Error Updating Workspace"
	ErrDeleteWorkspace = "Error Deleting Workspace"

	// ---- RuleNamespace ----
	ErrCreateRuleNamespace   = "Error creating RuleNamespace"
	ErrReadRuleNamespace     = "Error Reading RuleNamespace"
	ErrUpdateRuleNamespace   = "Error Updating RuleNamespace"
	ErrDeleteRuleNamespace   = "Error Deleting RuleNamespace"

	// ---- AlertManager ----
	ErrCreateAlertManager   = "Error creating AlertManager"
	ErrReadAlertManager     = "Error Reading AlertManager"
	ErrUpdateAlertManager   = "Error Updating AlertManager"
	ErrDeleteAlertManager   = "Error Deleting AlertManager"

	// ---- WorkspaceConfiguration ----
	ErrCreateWorkspaceConfiguration   = "Error creating WorkspaceConfiguration"
	ErrReadWorkspaceConfiguration     = "Error Reading WorkspaceConfiguration"
	ErrUpdateWorkspaceConfiguration   = "Error Updating WorkspaceConfiguration"
	ErrDeleteWorkspaceConfiguration   = "Error Deleting WorkspaceConfiguration"

	// ---- NotificationGroup ----
	ErrCreateNotificationGroup   = "Error creating NotificationGroup"
	ErrReadNotificationGroup     = "Error Reading NotificationGroup"
	ErrUpdateNotificationGroup   = "Error Updating NotificationGroup"
	ErrDeleteNotificationGroup   = "Error Deleting NotificationGroup"

	// ---- Template strings ----
	ErrUnexpectedConfigureFmt = "Expected *client.Instance, got: %T. Please report this issue to the provider developers."
	ErrCreateWorkspaceFmt     = "Could not create Workspace, unexpected error: %s\nReason: %s"
	ErrReadWorkspaceFmt       = "Could not read Workspace ID %s: %s\nReason: %s"
	ErrUpdateWorkspaceFmt     = "Could not update Workspace, unexpected error: %s\nReason: %s"
	ErrDeleteWorkspaceFmt     = "Could not delete Workspace, unexpected error: %s\nReason: %s"
	ErrCreateRuleNamespaceFmt = "Could not create RuleNamespace, unexpected error: %s\nReason: %s"
	ErrReadRuleNamespaceFmt   = "Could not read RuleNamespace ID %s: %s\nReason: %s"
	ErrUpdateRuleNamespaceFmt = "Could not update RuleNamespace, unexpected error: %s\nReason: %s"
	ErrDeleteRuleNamespaceFmt = "Could not delete RuleNamespace, unexpected error: %s\nReason: %s"
	ErrCreateAlertManagerFmt  = "Could not create AlertManager, unexpected error: %s\nReason: %s"
	ErrReadAlertManagerFmt    = "Could not read AlertManager ID %s: %s\nReason: %s"
	ErrUpdateAlertManagerFmt  = "Could not update AlertManager, unexpected error: %s\nReason: %s"
	ErrDeleteAlertManagerFmt  = "Could not delete AlertManager, unexpected error: %s\nReason: %s"
	ErrCreateWorkspaceConfigurationFmt  = "Could not create WorkspaceConfiguration, unexpected error: %s\nReason: %s"
	ErrReadWorkspaceConfigurationFmt    = "Could not read WorkspaceConfiguration ID %s: %s\nReason: %s"
	ErrUpdateWorkspaceConfigurationFmt  = "Could not update WorkspaceConfiguration, unexpected error: %s\nReason: %s"
	ErrDeleteWorkspaceConfigurationFmt  = "Could not delete WorkspaceConfiguration, unexpected error: %s\nReason: %s"
	ErrCreateNotificationGroupFmt  = "Could not create NotificationGroup, unexpected error: %s\nReason: %s"
	ErrReadNotificationGroupFmt    = "Could not read NotificationGroup ID %s: %s\nReason: %s"
	ErrUpdateNotificationGroupFmt  = "Could not update NotificationGroup, unexpected error: %s\nReason: %s"
	ErrDeleteNotificationGroupFmt  = "Could not delete NotificationGroup, unexpected error: %s\nReason: %s"
)

// ----------------------------
// Helper functions for nullable types
// ----------------------------

func nullableTimeTypes(value *time.Time, ok bool) types.String {
	if !ok || value == nil {
		return types.StringNull()
	}
	return types.StringValue(value.Format(TimeFormatDisplay))
}

func nullableStringTypes(value *string, ok bool) types.String {
	if !ok || value == nil {
		return types.StringNull()
	}
	return types.StringValue(*value)
}

func nullableInt64Types(value *int64, ok bool) types.Int64 {
	if !ok || value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*value)
}

func nullableBoolTypes(value *bool, ok bool) types.Bool {
	if !ok || value == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*value)
}

func privateAclResourcesToList(ctx context.Context, resources []smp.PrivateAclResource, diags *diag.Diagnostics) types.List {
	if resources == nil || len(resources) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"resource_id":   types.StringType,
				"resource_name": types.StringType,
				"resource_type": types.StringType,
			},
		})
	}

	var objs []types.Object
	for _, r := range resources {
		obj, d := types.ObjectValue(
			map[string]attr.Type{
				"resource_id":   types.StringType,
				"resource_name": types.StringType,
				"resource_type": types.StringType,
			},
			map[string]attr.Value{
				"resource_id":   types.StringValue(r.GetResourceId()),
				"resource_name": nullableStringTypes(r.GetResourceNameOk()),
				"resource_type": types.StringValue(r.GetResourceType()),
			},
		)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"resource_id":   types.StringType,
					"resource_name": types.StringType,
					"resource_type": types.StringType,
				},
			})
		}
		objs = append(objs, obj)
	}

	listVal, d := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"resource_id":   types.StringType,
			"resource_name": types.StringType,
			"resource_type": types.StringType,
		},
	}, objs)
	diags.Append(d...)
	return listVal
}
