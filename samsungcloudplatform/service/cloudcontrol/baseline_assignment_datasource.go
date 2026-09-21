package cloudcontrol

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cloudcontrol"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdkcloudcontrol "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/cloudcontrol/1.2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &baselineAssignmentDataSource{}
	_ datasource.DataSourceWithConfigure = &baselineAssignmentDataSource{}
)

func NewBaselineAssignmentDataSource() datasource.DataSource {
	return &baselineAssignmentDataSource{}
}

type baselineAssignmentDataSource struct {
	config  *scpsdk.Configuration
	client  *cloudcontrol.Client
	clients *client.SCPClient
}

func (d *baselineAssignmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudcontrol_baseline_assignments"
}

func (d *baselineAssignmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = inst.Client.CloudControl
	d.clients = inst.Client
}

func (d *baselineAssignmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List CloudControl Baseline Assignments with optional filters.",
		Attributes: map[string]schema.Attribute{
			"landing_zone_id": schema.StringAttribute{
				Optional: true,
				Description: "Landing Zone ID to filter by. \n" +
					"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
			},
			"resource_type": schema.StringAttribute{
				Optional: true,
				Description: "Type of resource to filter by. \n" +
					"  - example : 'OU' \n" +
					"  - allowed_values : ['OU', 'ACCOUNT'] \n",
			},
			"assignment_id": schema.StringAttribute{
				Optional: true,
				Description: "Root/Organization Unit/Account ID to filter by. \n" +
					"  - example : 'ou-b30e9fcc39f84a20bf9e7458e5ec3801' \n",
			},
			"status": schema.StringAttribute{
				Optional: true,
				Description: "Baseline assignment status to filter by. \n" +
					"  - example : 'REGISTERED' \n" +
					"  - allowed_values : ['REGISTERED', 'REGISTRATION_FAILED', 'REGISTERING', 'UNREGISTERING'] \n",
			},
			"baseline_assignments": schema.ListNestedAttribute{
				Description: "List of baseline assignments matching the filter. \n" +
					"  - example : '[{id: ou-b30e9fcc39f84a20bf9e7458e5ec3801, type: OU, status: REGISTERED, ...}]' \n",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the baseline assignment. \n" +
								"  - example : 'ou-b30e9fcc39f84a20bf9e7458e5ec3801' \n",
						},
						"type": schema.StringAttribute{
							Computed: true,
							Description: "Baseline Assignment Type. \n" +
								"  - example : 'OU' \n" +
								"  - allowed_values : ['OU', 'ACCOUNT'] \n",
						},
						"account_count": schema.Int64Attribute{
							Computed: true,
							Description: "Total number of accounts in the baseline. \n" +
								"  - example : '4' \n",
						},
						"account_assigned_count": schema.Int64Attribute{
							Computed: true,
							Description: "Number of accounts assigned to the baseline. \n" +
								"  - example : '3' \n",
						},
						"ou_count": schema.Int64Attribute{
							Computed: true,
							Description: "Total number of organization units in the baseline. \n" +
								"  - example : '10' \n",
						},
						"ou_assigned_count": schema.Int64Attribute{
							Computed: true,
							Description: "Number of organization units assigned to the baseline. \n" +
								"  - example : '5' \n",
						},
						"preventive_guardrail_activated_count": schema.Int64Attribute{
							Computed: true,
							Description: "Number of preventive guardrails currently active. \n" +
								"  - example : '8' \n",
						},
						"detective_guardrail_status": schema.StringAttribute{
							Computed: true,
							Description: "Status of detective guardrails. \n" +
								"  - example : 'ENABLED' \n" +
								"  - allowed_values : ['ENABLED', 'DISABLED', 'ENABLING', 'DISABLING'] \n",
						},
						"detective_guardrail_type": schema.StringAttribute{
							Computed: true,
							Description: "Type of detective guardrails. \n" +
								"  - example : 'SSI' \n",
						},
						"sso_user_name": schema.StringAttribute{
							Computed: true,
							Description: "Username of the Identity Center user. \n" +
								"  - example : 'testuser' \n",
						},
						"status": schema.StringAttribute{
							Computed: true,
							Description: "Baseline Assignment Status. \n" +
								"  - example : 'REGISTERED' \n" +
								"  - allowed_values : ['REGISTERED', 'REGISTRATION_FAILED', 'REGISTERING', 'UNREGISTERING'] \n",
						},
					},
				},
			},
		},
	}
}

func (d *baselineAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudcontrol.BaselineAssignmentDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	landingZoneId := state.LandingZoneId.ValueString()
	resourceType := state.ResourceType.ValueString()
	assignmentId := state.AssignmentId.ValueString()
	status := state.Status.ValueString()

	data, err := d.client.ListBaselineAssignments(ctx, landingZoneId, resourceType, assignmentId, status)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read CloudControl Baseline Assignments",
			"Could not read CloudControl Baseline Assignments, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	results := make([]cloudcontrol.BaselineAssignmentListItem, 0)
	for _, ba := range data.BaselineAssignments {
		results = append(results, d.buildBaselineAssignmentItem(&ba))
	}
	state.BaselineAssignments = results

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *baselineAssignmentDataSource) buildBaselineAssignmentItem(ba *sdkcloudcontrol.BaselineAssignmentV1Dot1) cloudcontrol.BaselineAssignmentListItem {
	item := cloudcontrol.BaselineAssignmentListItem{
		Id:   types.StringValue(ba.Id),
		Type: types.StringValue(ba.Type),
	}

	if ba.AccountAssignedCount.IsSet() && ba.AccountAssignedCount.Get() != nil {
		item.AccountAssignedCount = types.Int64Value(int64(*ba.AccountAssignedCount.Get()))
	} else {
		item.AccountAssignedCount = types.Int64Null()
	}
	if ba.AccountCount.IsSet() && ba.AccountCount.Get() != nil {
		item.AccountCount = types.Int64Value(int64(*ba.AccountCount.Get()))
	} else {
		item.AccountCount = types.Int64Null()
	}
	if ba.OuAssignedCount.IsSet() && ba.OuAssignedCount.Get() != nil {
		item.OuAssignedCount = types.Int64Value(int64(*ba.OuAssignedCount.Get()))
	} else {
		item.OuAssignedCount = types.Int64Null()
	}
	if ba.OuCount.IsSet() && ba.OuCount.Get() != nil {
		item.OuCount = types.Int64Value(int64(*ba.OuCount.Get()))
	} else {
		item.OuCount = types.Int64Null()
	}
	if ba.PreventiveGuardrailActivatedCount.IsSet() && ba.PreventiveGuardrailActivatedCount.Get() != nil {
		item.PreventiveGuardrailActivatedCount = types.Int64Value(int64(*ba.PreventiveGuardrailActivatedCount.Get()))
	} else {
		item.PreventiveGuardrailActivatedCount = types.Int64Null()
	}
	if ba.DetectiveGuardrailStatus.IsSet() && ba.DetectiveGuardrailStatus.Get() != nil {
		item.DetectiveGuardrailStatus = types.StringValue(*ba.DetectiveGuardrailStatus.Get())
	} else {
		item.DetectiveGuardrailStatus = types.StringNull()
	}
	if ba.DetectiveGuardrailType.IsSet() && ba.DetectiveGuardrailType.Get() != nil {
		item.DetectiveGuardrailType = types.StringValue(*ba.DetectiveGuardrailType.Get())
	} else {
		item.DetectiveGuardrailType = types.StringNull()
	}
	if ba.SsoUserName.IsSet() && ba.SsoUserName.Get() != nil {
		item.SsoUserName = types.StringValue(*ba.SsoUserName.Get())
	} else {
		item.SsoUserName = types.StringNull()
	}
	if ba.Status.IsSet() && ba.Status.Get() != nil {
		item.Status = types.StringValue(*ba.Status.Get())
	} else {
		item.Status = types.StringNull()
	}

	return item
}
