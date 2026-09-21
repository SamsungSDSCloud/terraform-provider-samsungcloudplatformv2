package cloudcontrol

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cloudcontrol"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &baselineAssignmentResource{}
	_ resource.ResourceWithConfigure   = &baselineAssignmentResource{}
	_ resource.ResourceWithImportState = &baselineAssignmentResource{}
)

func NewBaselineAssignmentResource() resource.Resource {
	return &baselineAssignmentResource{}
}

type baselineAssignmentResource struct {
	config  *scpsdk.Configuration
	client  *cloudcontrol.Client
	clients *client.SCPClient
}

func (r *baselineAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudcontrol_baseline_assignment"
}

func (r *baselineAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Baseline Assignment - assigns baseline configuration to accounts or organization units. \n",
		Attributes: map[string]schema.Attribute{
			"assignment_id": schema.StringAttribute{
				Description: "Root/Organization Unit/Account ID to assign baseline to. \n" +
					"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"landing_zone_id": schema.StringAttribute{
				Description: "Landing Zone ID that contains this baseline assignment. \n" +
					"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
				Required: true,
			},
			"resource_type": schema.StringAttribute{
				Description: "Type of resource to assign baseline to. \n" +
					"  - example : 'OU' \n" +
					"  - allowed_values : ['ACCOUNT', 'OU'] \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"agree_yn": schema.StringAttribute{
				Description: "Agreement to baseline assignment. \n" +
					"  - example : 'Y' \n" +
					"  - allowed_values : ['Y', 'N'] \n",
				Optional: true,
			},
			"reregister_trigger": schema.StringAttribute{
				Description: "Arbitrary value used to force a re-registration (Update) of this baseline assignment. \n" +
					"This value is not sent to the API; changing it simply causes Terraform to detect a diff and call Update, \n" +
					"which re-invokes the baseline re-registration flow. Only meaningful for resource_type = 'OU'. \n" +
					"  - example : '1' \n",
				Optional: true,
			},
			"parent_unit_id": schema.StringAttribute{
				Description: "Parent Organization Unit ID. \n" +
					"  - example : 'ou-fc8c29a138d78e24bf1fa86812fc8b' \n",
				Optional: true,
			},
			"sso_user_name": schema.StringAttribute{
				Description: "SSO User Name for the assignment. \n" +
					"  - example : 'testuser' \n",
				Optional: true,
			},
			"sso_user_real_name": schema.StringAttribute{
				Description: "Real name of the SSO User. \n" +
					"  - example : 'test user' \n",
				Optional: true,
			},
			"job_id": schema.StringAttribute{
				Description: "Job ID for the baseline assignment operation. \n" +
					"  - example : '0a36e0746dbf4908acf0357829701381' \n",
				Computed: true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the baseline assignment. \n" +
					"  - example : 'REGISTERED' \n" +
					"  - allowed_values : ['REGISTERED', 'REGISTRATION_FAILED', 'REGISTERING', 'UNREGISTERING'] \n",
				Computed: true,
			},
			"account_count": schema.Int64Attribute{
				Description: "Number of accounts in this baseline assignment. \n" +
					"  - example : '5' \n",
				Computed: true,
			},
			"account_assigned_count": schema.Int64Attribute{
				Description: "Number of assigned accounts. \n" +
					"  - example : '3' \n",
				Computed: true,
			},
			"ou_count": schema.Int64Attribute{
				Description: "Number of Organization Units in this baseline assignment. \n" +
					"  - example : '2' \n",
				Computed: true,
			},
			"ou_assigned_count": schema.Int64Attribute{
				Description: "Number of assigned Organization Units. \n" +
					"  - example : '1' \n",
				Computed: true,
			},
			"detective_guardrail_status": schema.StringAttribute{
				Description: "Status of detective guardrails. \n" +
					"  - example : 'ENABLED' \n",
				Computed: true,
			},
			"detective_guardrail_type": schema.StringAttribute{
				Description: "Type of detective guardrails. \n" +
					"  - example : 'DETECTIVE' \n",
				Computed: true,
			},
		},
	}
}

func (r *baselineAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.CloudControl
	r.clients = inst.Client
}

func (r *baselineAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cloudcontrol.BaselineAssignmentResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateBaselineAssignment(ctx, plan.AssignmentId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Baseline Assignment",
			"Could not create Baseline Assignment, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.JobId = types.StringValue(data.JobId)

	err = waitForBaselineAssignmentActiveStatus(ctx, r.client, plan.LandingZoneId.ValueString(), plan.AssignmentId.ValueString(), plan.ResourceType.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error waiting for baseline assignment creation",
			"Error waiting for baseline assignment "+plan.AssignmentId.ValueString()+" to become active: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

func waitForBaselineAssignmentActiveStatus(ctx context.Context, ccClient *cloudcontrol.Client, landingZoneId, assignmentId, resourceType string) error {
	return client.WaitForStatus(ctx, nil, []string{"REGISTERING"}, []string{"REGISTERED"}, func() (interface{}, string, error) {
		data, err := ccClient.ListBaselineAssignments(ctx, landingZoneId, resourceType, assignmentId, "")
		if err != nil {
			return nil, "", err
		}
		if len(data.BaselineAssignments) == 0 {
			return nil, "", fmt.Errorf("baseline assignment not found")
		}
		assignment := data.BaselineAssignments[0]
		if assignment.Status.IsSet() && assignment.Status.Get() != nil {
			status := *assignment.Status.Get()
			if status == "REGISTRATION_FAILED" {
				return data, status, fmt.Errorf("baseline assignment registration failed")
			}
			return data, status, nil
		}
		return data, "", nil
	}, -1, -1, -1, -1)
}

func (r *baselineAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cloudcontrol.BaselineAssignmentResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ListBaselineAssignments(ctx, state.LandingZoneId.ValueString(), state.ResourceType.ValueString(), state.AssignmentId.ValueString(), "")
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if len(data.BaselineAssignments) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	assignment := data.BaselineAssignments[0]

	state.AssignmentId = types.StringValue(assignment.Id)
	state.ResourceType = types.StringValue(assignment.Type)

	if assignment.Status.IsSet() && assignment.Status.Get() != nil {
		state.Status = types.StringValue(*assignment.Status.Get())
	} else {
		state.Status = types.StringNull()
	}

	if assignment.AccountCount.IsSet() && assignment.AccountCount.Get() != nil {
		state.AccountCount = types.Int64Value(int64(*assignment.AccountCount.Get()))
	} else {
		state.AccountCount = types.Int64Null()
	}

	if assignment.AccountAssignedCount.IsSet() && assignment.AccountAssignedCount.Get() != nil {
		state.AccountAssignedCount = types.Int64Value(int64(*assignment.AccountAssignedCount.Get()))
	} else {
		state.AccountAssignedCount = types.Int64Null()
	}

	if assignment.OuCount.IsSet() && assignment.OuCount.Get() != nil {
		state.OuCount = types.Int64Value(int64(*assignment.OuCount.Get()))
	} else {
		state.OuCount = types.Int64Null()
	}

	if assignment.OuAssignedCount.IsSet() && assignment.OuAssignedCount.Get() != nil {
		state.OuAssignedCount = types.Int64Value(int64(*assignment.OuAssignedCount.Get()))
	} else {
		state.OuAssignedCount = types.Int64Null()
	}

	if assignment.DetectiveGuardrailStatus.IsSet() && assignment.DetectiveGuardrailStatus.Get() != nil {
		state.DetectiveGuardrailStatus = types.StringValue(*assignment.DetectiveGuardrailStatus.Get())
	} else {
		state.DetectiveGuardrailStatus = types.StringNull()
	}

	if assignment.DetectiveGuardrailType.IsSet() && assignment.DetectiveGuardrailType.Get() != nil {
		state.DetectiveGuardrailType = types.StringValue(*assignment.DetectiveGuardrailType.Get())
	} else {
		state.DetectiveGuardrailType = types.StringNull()
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *baselineAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state cloudcontrol.BaselineAssignmentResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var fullPlan cloudcontrol.BaselineAssignmentResource
	diags = req.Plan.Get(ctx, &fullPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := cloudcontrol.BaselineAssignmentUpdateResource{
		LandingZoneId: fullPlan.LandingZoneId,
		AgreeYn:       fullPlan.AgreeYn,
	}

	assignmentId := state.AssignmentId.ValueString()

	updateResp, err := r.client.UpdateBaselineAssignment(ctx, assignmentId, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Baseline Assignment",
			"Could not update Baseline Assignment, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	state.JobId = types.StringValue(updateResp.JobId)
	state.ReregisterTrigger = fullPlan.ReregisterTrigger

	err = waitForBaselineAssignmentActiveStatus(ctx, r.client, state.LandingZoneId.ValueString(), assignmentId, state.ResourceType.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error waiting for baseline assignment update",
			"Error waiting for baseline assignment "+assignmentId+" to become active: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

func (r *baselineAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cloudcontrol.BaselineAssignmentResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ResourceType.ValueString() != "ACCOUNT" {
		resp.Diagnostics.AddWarning(
			"Baseline assignment not unassigned",
			"Only ACCOUNT resource type supports unassignment. "+
				"The resource has been removed from Terraform state, "+
				"but the baseline assignment still exists in the landing zone.",
		)
		return
	}

	_, err := r.client.DeleteBaselineAssignment(ctx, state.LandingZoneId.ValueString(), state.AssignmentId.ValueString(), state.ResourceType.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting Baseline Assignment",
			"Could not delete Baseline Assignment, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForBaselineAssignmentAbsent(ctx, r.client, state.LandingZoneId.ValueString(), state.AssignmentId.ValueString(), state.ResourceType.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error waiting for Baseline Assignment unregistration",
			"Error waiting for baseline assignment "+state.AssignmentId.ValueString()+" to fully unregister: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func waitForBaselineAssignmentAbsent(ctx context.Context, ccClient *cloudcontrol.Client, landingZoneId, assignmentId, resourceType string) error {
	return client.WaitForStatus(ctx, nil, []string{"PRESENT"}, []string{"ABSENT"}, func() (interface{}, string, error) {
		data, err := ccClient.ListBaselineAssignments(ctx, landingZoneId, resourceType, assignmentId, "")
		if err != nil {
			return nil, "", err
		}
		count := len(data.BaselineAssignments)
		statusSet := count > 0 &&
			data.BaselineAssignments[0].Status.IsSet() &&
			data.BaselineAssignments[0].Status.Get() != nil
		if count == 0 || !statusSet {
			return data, "ABSENT", nil
		}
		return data, "PRESENT", nil
	}, -1, -1, -1, -1)
}

func (r *baselineAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("assignment_id"), req, resp)
}
