package configinspection

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/configinspection"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &configInspectionDiagnosisResource{}
	_ resource.ResourceWithConfigure   = &configInspectionDiagnosisResource{}
	_ resource.ResourceWithImportState = &configInspectionDiagnosisResource{}
)

// NewConfigInspectionDiagnosisResource is a helper function to simplify the provider implementation.
func NewConfigInspectionDiagnosisResource() resource.Resource {
	return &configInspectionDiagnosisResource{}
}

// configInspectionDiagnosisResource is the resource implementation.
type configInspectionDiagnosisResource struct {
	config  *scpsdk.Configuration //lint:ignore U1000 Ignore unused
	client  *configinspection.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name
func (r *configInspectionDiagnosisResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configinspection"
}

// Schema defines the resource schema
func (r *configInspectionDiagnosisResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Config inspection resource.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Account Identifier.\n" +
					"  - example : '0e3dffc50eb247a1adxxxxxxxxxxxxxx'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("AuthKeyRequest"): schema.SingleNestedAttribute{
				Description: "Auth key request",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
						Description: "Id of diagnosis.\n" +
							"  - example : 'DIA-943731CB8E3045C289xxxxxxxxxxxxxx'",
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					common.ToSnakeCase("AuthKeyCreatedAt"): schema.StringAttribute{
						Description: "Created date of authkey.\n" +
							"  - example : '2022-01-01 12:00:00'",
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					common.ToSnakeCase("AuthKeyExpiredAt"): schema.StringAttribute{
						Description: "Expired date of authkey.\n" +
							"  - example : '2023-01-01 12:00:00'",
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					common.ToSnakeCase("AuthKeyId"): schema.StringAttribute{
						Description: "Id of auth key.\n" +
							"  - example : '9b72a9856e494exxxxxxxxxxxxxxxxxx'",
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			common.ToSnakeCase("CspType"): schema.StringAttribute{
				Description: "Type of cloud service provider.\n" +
					"  - example : 'SCP'\n" +
					"  - enum : SCP | AWS | Azure",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("DiagnosisAccountId"): schema.StringAttribute{
				Description: "Account Id of diagnosis.\n" +
					"  - example : '0e3dffc50eb247a1adxxxxxxxxxxxxxx'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("DiagnosisCheckType"): schema.StringAttribute{
				Description: "Check type of diagnosis.\n" +
					"  - example : 'BP'\n" +
					"  - enum : BP | SSI",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
				Description: "Id of diagnosis.\n" +
					"  - example : 'DIA-943731CB8E3045C289xxxxxxxxxxxxxx'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("DiagnosisName"): schema.StringAttribute{
				Description: "Name of diagnosis.\n" +
					"  - example : 'Sample Diagnosis Name'\n" +
					"  - pattern : `^[a-zA-Z0-9-_]+$`",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("DiagnosisType"): schema.StringAttribute{
				Description: "Config inspection type.\n" +
					"  - example : 'Console'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("PlanType"): schema.StringAttribute{
				Description: "Billing plan for the inspection.\n" +
					"  - example : 'STANDARD'\n" +
					"  - enum : STANDARD | MONTHLY",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("ScheduleRequest"): schema.SingleNestedAttribute{
				Description: "Schedule request.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
						Description: "Id of diagnosis.\n" +
							"  - example : 'DIA-943731CB8E3045C289xxxxxxxxxxxxxx'",
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					common.ToSnakeCase("DiagnosisStartTimePattern"): schema.StringAttribute{
						Description: "Start time (5-minute increments, 00 to 23 hours, 00 to 55 minutes).\n" +
							"  - example : '08:00'",
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					common.ToSnakeCase("FrequencyType"): schema.StringAttribute{
						Description: "Schedule type (monthly, weekly, daily).\n" +
							"  - example : 'MONTH'",
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					common.ToSnakeCase("FrequencyValue"): schema.StringAttribute{
						Description: "Schedule value (01~31, MONDAY~SUNDAY, everyDay).\n" +
							"  - example : 1",
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					common.ToSnakeCase("UseDiagnosisCheckTypeBp"): schema.StringAttribute{
						Description: "Checklist Best Practice Use.\n" +
							"  - example : 'y'",
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					common.ToSnakeCase("UseDiagnosisCheckTypeSsi"): schema.StringAttribute{
						Description: "Checklist SSI usage.\n" +
							"  - example : 'y'",
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			common.ToSnakeCase("Tags"): tag.ResourceSchema(),

			// Output
			common.ToSnakeCase("Result"): schema.BoolAttribute{
				Description: "Result of diagnosis request (true, false).\n" +
					"  - example : true",
				Computed: true,
			},
			common.ToSnakeCase("CreatedDiagnosisId"): schema.StringAttribute{
				Description: "Id of created diagnosis.\n" +
					"  - example : 'DIA-943731CB8E3045C289xxxxxxxxxxxxxx'",
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *configInspectionDiagnosisResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	r.client = inst.Client.ConfigInspection
	r.clients = inst.Client
}

// Create creates a new resource
func (r *configInspectionDiagnosisResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state configinspection.ConfigInspectionDiagnosisResource

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.CreateConfigInspectionObject(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to create config inspection",
			fmt.Sprintf("An error occurred while creating config inspection: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	if res != nil {
		state.Result = types.BoolValue(res.Result)
		state.CreatedDiagnosisId = types.StringValue(res.DiagnosisId)
	} else {
		resp.Diagnostics.AddError(
			"Failed to create config inspection",
			"An error occurred while creating config inspection. Empty response",
		)
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads the resource state
func (r *configInspectionDiagnosisResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state configinspection.ConfigInspectionDiagnosisResource

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetConfigInspectionObjectDetail(ctx, state.CreatedDiagnosisId.ValueString())
	if err != nil {
		// Check if the error indicates the resource was not found
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to read config inspection detail",
			fmt.Sprintf("An error occurred while reading config inspection detail: %s. Details: %s", err.Error(), detail),
		)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			"Error reading config inspection",
			"An error occurred while reading config inspection. Empty response",
		)
		return
	}
	// Terraform resource only manage ID (used for deletion)
	state.CreatedDiagnosisId = types.StringValue(res.SummaryResponses.DiagnosisId)

	// Update input fields from API response for drift detection
	state.AccountId = types.StringPointerValue(res.AuthKeyResponses.UserId)
	state.DiagnosisAccountId = types.StringValue(res.SummaryResponses.DiagnosisAccountId)
	state.CspType = types.StringValue(res.SummaryResponses.CspType)
	state.DiagnosisCheckType = types.StringValue(res.SummaryResponses.DiagnosisCheckType)
	state.DiagnosisId = types.StringValue(res.SummaryResponses.DiagnosisId)
	state.DiagnosisName = types.StringValue(res.SummaryResponses.DiagnosisName)
	state.DiagnosisType = types.StringValue(res.SummaryResponses.DiagnosisType)
	state.PlanType = types.StringValue(res.SummaryResponses.PlanType)
	state.ScheduleRequest = &configinspection.DiagnosisScheduleRequest{
		DiagnosisId:               types.StringPointerValue(res.ScheduleResponse.DiagnosisId),
		DiagnosisStartTimePattern: types.StringPointerValue(res.ScheduleResponse.DiagnosisStartTimePattern),
		FrequencyType:             types.StringPointerValue(res.ScheduleResponse.FrequencyType),
		FrequencyValue:            types.StringPointerValue(res.ScheduleResponse.FrequencyValue),
		UseDiagnosisCheckTypeBp:   types.StringPointerValue(res.ScheduleResponse.UseDiagnosisCheckTypeBp),
		UseDiagnosisCheckTypeSsi:  types.StringPointerValue(res.ScheduleResponse.UseDiagnosisCheckTypeSsi),
	}
	state.AuthKeyRequest = &configinspection.AuthKeyRequest{
		AuthKeyId:        types.StringPointerValue(res.AuthKeyResponses.AuthKeyId),
		AuthKeyCreatedAt: types.StringPointerValue(res.AuthKeyResponses.AuthKeyCreatedAt),
		AuthKeyExpiredAt: types.StringPointerValue(res.AuthKeyResponses.AuthKeyExpiredAt),
	}

	// Save updated data into Terraform state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource
func (r *configInspectionDiagnosisResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Config inspection does not support update operations
	// This is a no-op implementation
	resp.Diagnostics.AddWarning(
		"Update not supported",
		"Config inspection resources do not support update operations. The resource will not be updated.",
	)
}

// Delete deletes the resource
func (r *configInspectionDiagnosisResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state configinspection.ConfigInspectionDiagnosisResource

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteConfigInspectionObject(ctx, state.CreatedDiagnosisId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to delete config inspection",
			fmt.Sprintf("An error occurred while deleting config inspection: %s. Details: %s", err.Error(), detail),
		)
		return
	}
}

// ImportState imports an existing resource into Terraform state using its ID.
func (r *configInspectionDiagnosisResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("created_diagnosis_id"), req, resp)
}
