package configinspection

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/configinspection"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &configInspectionDiagnosisResource{}
	_ resource.ResourceWithConfigure = &configInspectionDiagnosisResource{}
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
				Description: "Account Id\n" +
					"  - Example: 0e3dffc50eb247a1adf4f2e5c82c4f99",
				Required: true,
			},
			common.ToSnakeCase("AuthKeyRequest"): schema.SingleNestedAttribute{
				Description: "Auth key request",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
						Description: "Diagnosis ID\n" +
							"  - Example: DIA-943731CB8E3045C289BAECAEC3532097",
						Optional: true,
					},
					common.ToSnakeCase("AuthKeyCreatedAt"): schema.StringAttribute{
						Description: "Auth key created at\n" +
							"  - Example: 2022-01-01 12:00:00",
						Optional: true,
					},
					common.ToSnakeCase("AuthKeyExpiredAt"): schema.StringAttribute{
						Description: "Auth key expired at\n" +
							"  - Example: 2023-01-01 12:00:00",
						Optional: true,
					},
					common.ToSnakeCase("AuthKeyId"): schema.StringAttribute{
						Description: "Auth key ID\n" +
							"  - Example: 9b72a9856e494e67afc69atd3631fe38",
						Required: true,
					},
				},
			},
			common.ToSnakeCase("CspType"): schema.StringAttribute{
				Description: "Type of cloud service provider\n" +
					"  - Example: SCP",
				Required: true,
			},
			common.ToSnakeCase("DiagnosisAccountId"): schema.StringAttribute{
				Description: "Id of diagnosis\n" +
					"  - Example: 0e3dffc50eb247a1adf4f2e5c82c4f99",
				Required: true,
			},
			common.ToSnakeCase("DiagnosisCheckType"): schema.StringAttribute{
				Description: "Check type of diagnosis\n" +
					"  - Example: BP",
				Required: true,
			},
			common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
				Description: "Id of diagnosis\n" +
					"  - Example: DIA-943731CB8E3045C289BAECAEC3532097",
				Required: true,
			},
			common.ToSnakeCase("DiagnosisName"): schema.StringAttribute{
				Description: "Name of diagnosis\n" +
					"  - Example: Sample Diagnosis Name",
				Required: true,
			},
			common.ToSnakeCase("DiagnosisType"): schema.StringAttribute{
				Description: "Diagnosis Type\n" +
					"  - Example: Console",
				Required: true,
			},
			common.ToSnakeCase("PlanType"): schema.StringAttribute{
				Description: "Plan\n" +
					"  - Example: STANDARD",
				Required: true,
			},
			common.ToSnakeCase("ScheduleRequest"): schema.SingleNestedAttribute{
				Description: "Schedule request",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("DiagnosisId"): schema.StringAttribute{
						Description: "Diagnosis ID\n" +
							"  - Example: DIA-943731CB8E3045C289BAECAEC3532097",
						Required: true,
					},
					common.ToSnakeCase("DiagnosisStartTimePattern"): schema.StringAttribute{
						Description: "Diagnosis start time pattern\n" +
							"  - Example: 08:00",
						Required: true,
					},
					common.ToSnakeCase("FrequencyType"): schema.StringAttribute{
						Description: "Frequency type\n" +
							"  - Example: MONTH",
						Required: true,
					},
					common.ToSnakeCase("FrequencyValue"): schema.StringAttribute{
						Description: "Frequency value\n" +
							"  - Example:1",
						Required: true,
					},
					common.ToSnakeCase("UseDiagnosisCheckTypeBp"): schema.StringAttribute{
						Description: "Use diagnosis check type BP\n" +
							"  - Example: y",
						Required: true,
					},
					common.ToSnakeCase("UseDiagnosisCheckTypeSsi"): schema.StringAttribute{
						Description: "Use diagnosis check type SSI\n" +
							"  - Example: y",
						Required: true,
					},
				},
			},
			common.ToSnakeCase("Tags"): tag.ResourceSchema(),

			// Output
			common.ToSnakeCase("Result"): schema.BoolAttribute{
				Description: "Result",
				Computed:    true,
			},
			common.ToSnakeCase("CreatedDiagnosisId"): schema.StringAttribute{
				Description: "Id of created diagnosis",
				Computed:    true,
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

	// Wait for sdk changes returned diagnosis ID
	state.Result = types.BoolValue(res.Result)
	state.CreatedDiagnosisId = types.StringValue(res.DiagnosisId)

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

	_, err := r.client.GetConfigInspectionObjectDetail(ctx, state.CreatedDiagnosisId.ValueString())
	if err != nil {
		// Check if the error indicates the resource was not found
		if err.Error() == "404 Not Found" {
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
