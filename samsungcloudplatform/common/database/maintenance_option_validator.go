package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type maintenanceOptionValidator struct{}

// MaintenanceOptionValidator returns an object validator for the maintenance_option
// attribute that rejects configurations whose schedule the provider would silently drop.
//
// use_maintenance_option is never sent to the API: the create request carries a
// maintenance option only when period_hour, starting_day_of_week and starting_time are
// all present, and the flag merely gates that request. When the flag and the schedule
// disagree the cluster is created without a maintenance option while the configuration
// says otherwise, so state can never match the configuration and every subsequent plan
// shows a diff that apply cannot resolve (maintenance_option is immutable).
func MaintenanceOptionValidator() validator.Object {
	return maintenanceOptionValidator{}
}

func (v maintenanceOptionValidator) Description(_ context.Context) string {
	return "Requires the maintenance schedule to be set when use_maintenance_option is true, and to be absent otherwise."
}

func (v maintenanceOptionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v maintenanceOptionValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attributes := req.ConfigValue.Attributes()

	useMaintenanceOption, ok := attributes[common.ToSnakeCase("UseMaintenanceOption")].(types.Bool)
	if !ok || useMaintenanceOption.IsUnknown() {
		return
	}

	scheduleNames := []string{
		common.ToSnakeCase("PeriodHour"),
		common.ToSnakeCase("StartingDayOfWeek"),
		common.ToSnakeCase("StartingTime"),
	}

	var missing, present []string
	for _, name := range scheduleNames {
		value, ok := attributes[name].(types.String)
		if !ok || value.IsUnknown() {
			// Only known after apply: nothing to validate yet.
			return
		}

		if value.IsNull() {
			missing = append(missing, name)
		} else {
			present = append(present, name)
		}
	}

	if useMaintenanceOption.ValueBool() {
		if len(missing) > 0 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Incomplete Maintenance Option",
				fmt.Sprintf("use_maintenance_option is true, so %s must all be set, but %s %s not set. "+
					"Set them, or set use_maintenance_option to false.",
					strings.Join(scheduleNames, ", "),
					strings.Join(missing, ", "),
					pluralIs(len(missing)),
				),
			)
		}

		return
	}

	if len(present) > 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Unused Maintenance Option",
			fmt.Sprintf("use_maintenance_option is not true, so the maintenance schedule is ignored, but %s %s set. "+
				"Remove them, or set use_maintenance_option to true.",
				strings.Join(present, ", "),
				pluralIs(len(present)),
			),
		)
	}
}

// pluralIs returns the verb matching count so diagnostics read naturally.
func pluralIs(count int) string {
	if count == 1 {
		return "is"
	}

	return "are"
}
