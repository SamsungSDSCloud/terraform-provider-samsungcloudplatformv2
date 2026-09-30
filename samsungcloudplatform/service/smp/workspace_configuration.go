package smp

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/smp"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &smpWorkspaceConfigurationResource{}
	_ resource.ResourceWithConfigure   = &smpWorkspaceConfigurationResource{}
	_ resource.ResourceWithImportState = &smpWorkspaceConfigurationResource{}
)

func NewSmpWorkspaceConfigurationResource() resource.Resource {
	return &smpWorkspaceConfigurationResource{}
}

type smpWorkspaceConfigurationResource struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

// --------------------------------------------------------------------------
// METADATA
// --------------------------------------------------------------------------
func (r *smpWorkspaceConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_workspace_configuration"
}

// --------------------------------------------------------------------------
// SCHEMA
// --------------------------------------------------------------------------
func (r *smpWorkspaceConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP WorkspaceConfiguration Resource",
		Attributes: map[string]schema.Attribute{
			// ----- Computed block that returns the whole configuration -----
			"workspace_configuration": schema.SingleNestedAttribute{
				Description: "Workspace configuration details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id":               schema.StringAttribute{Computed: true},
					"state":            schema.StringAttribute{Computed: true},
					"workspace_id":     schema.StringAttribute{Computed: true},
					"retention_period": schema.Int64Attribute{Computed: true},
					"created_at":       schema.StringAttribute{Computed: true},
					"created_by":       schema.StringAttribute{Computed: true},
					"modified_at":      schema.StringAttribute{Computed: true},
					"modified_by":      schema.StringAttribute{Computed: true},
				},
			},

			// ----- Input attributes -------------------------------------------------
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"retention_period": schema.Int64Attribute{
				Description: "Retention period in days.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// --------------------------------------------------------------------------
// CONFIGURE
// --------------------------------------------------------------------------
func (r *smpWorkspaceConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(ErrUnexpectedConfigure,
			fmt.Sprintf(ErrUnexpectedConfigureFmt, req.ProviderData))
		return
	}
	r.client = inst.Client.Smp
	r.clients = inst.Client
}

// --------------------------------------------------------------------------
// CREATE
// --------------------------------------------------------------------------
func (r *smpWorkspaceConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan smp.WorkspaceConfigurationResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The SDK does not expose a dedicated “Create” endpoint – we use Update (Set) instead.
	result, err := r.client.UpdateWorkspaceConfiguration(ctx, plan.WorkspaceId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrCreateWorkspaceConfiguration,
			fmt.Sprintf(ErrCreateWorkspaceConfigurationFmt, err.Error(), detail))
		return
	}

	wc := result.GetWorkspaceConfiguration()
	wcObj, diags := types.ObjectValue(
		smp.WorkspaceConfiguration{}.AttributeTypes(),
		map[string]attr.Value{
			"id":               types.StringValue(wc.GetId()),
			"state":            types.StringValue(wc.GetState()),
			"workspace_id":     types.StringValue(wc.GetWorkspaceId()),
			"retention_period": types.Int64Value(int64(wc.GetRetentionPeriod())),
			"created_at":       types.StringValue(wc.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":       types.StringValue(wc.GetCreatedBy()),
			"modified_at":      nullableTimeTypes(wc.GetModifiedAtOk()),
			"modified_by":      nullableStringTypes(wc.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.WorkspaceConfiguration = wcObj
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// --------------------------------------------------------------------------
// READ
// --------------------------------------------------------------------------
func (r *smpWorkspaceConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state smp.WorkspaceConfigurationResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetWorkspaceConfiguration(ctx, state.WorkspaceId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadWorkspaceConfiguration,
			fmt.Sprintf(ErrReadWorkspaceConfigurationFmt,
				state.WorkspaceId.ValueString(),
				err.Error(),
				detail))
		return
	}

	wc := result.GetWorkspaceConfiguration()
	wcObj, diags := types.ObjectValue(
		smp.WorkspaceConfiguration{}.AttributeTypes(),
		map[string]attr.Value{
			"id":               types.StringValue(wc.GetId()),
			"state":            types.StringValue(wc.GetState()),
			"workspace_id":     types.StringValue(wc.GetWorkspaceId()),
			"retention_period": types.Int64Value(int64(wc.GetRetentionPeriod())),
			"created_at":       types.StringValue(wc.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":       types.StringValue(wc.GetCreatedBy()),
			"modified_at":      nullableTimeTypes(wc.GetModifiedAtOk()),
			"modified_by":      nullableStringTypes(wc.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.WorkspaceConfiguration = wcObj
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// --------------------------------------------------------------------------
// UPDATE
// --------------------------------------------------------------------------
func (r *smpWorkspaceConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan smp.WorkspaceConfigurationResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateWorkspaceConfiguration(ctx, plan.WorkspaceId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrUpdateWorkspaceConfiguration,
			fmt.Sprintf(ErrUpdateWorkspaceConfigurationFmt, err.Error(), detail))
		return
	}

	result, err := r.client.GetWorkspaceConfiguration(ctx, plan.WorkspaceId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadWorkspaceConfiguration,
			fmt.Sprintf(ErrReadWorkspaceConfigurationFmt,
				plan.WorkspaceId.ValueString(),
				err.Error(),
				detail))
		return
	}

	wc := result.GetWorkspaceConfiguration()
	wcObj, diags := types.ObjectValue(
		smp.WorkspaceConfiguration{}.AttributeTypes(),
		map[string]attr.Value{
			"id":               types.StringValue(wc.GetId()),
			"state":            types.StringValue(wc.GetState()),
			"workspace_id":     types.StringValue(wc.GetWorkspaceId()),
			"retention_period": types.Int64Value(int64(wc.GetRetentionPeriod())),
			"created_at":       types.StringValue(wc.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":       types.StringValue(wc.GetCreatedBy()),
			"modified_at":      nullableTimeTypes(wc.GetModifiedAtOk()),
			"modified_by":      nullableStringTypes(wc.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.WorkspaceConfiguration = wcObj
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// --------------------------------------------------------------------------
// DELETE
// --------------------------------------------------------------------------
func (r *smpWorkspaceConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state smp.WorkspaceConfigurationResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The SDK does not expose a delete endpoint; we just call the no‑op wrapper.
	err := r.client.DeleteWorkspaceConfiguration(ctx, state.WorkspaceId.ValueString())
	if err != nil {
	  detail := client.GetDetailFromError(err)
	  resp.Diagnostics.AddError(ErrDeleteWorkspaceConfiguration,
	      fmt.Sprintf(ErrDeleteWorkspaceConfigurationFmt, err.Error(), detail))
	  return
   }
}

// --------------------------------------------------------------------------
// IMPORT
// --------------------------------------------------------------------------
func (r *smpWorkspaceConfigurationResource) ImportState(ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("workspace_id"), req, resp)
}
