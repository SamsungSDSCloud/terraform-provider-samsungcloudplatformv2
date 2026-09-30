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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &smpAlertManagerResource{}
	_ resource.ResourceWithConfigure   = &smpAlertManagerResource{}
	_ resource.ResourceWithImportState = &smpAlertManagerResource{}
)

func NewSmpAlertManagerResource() resource.Resource {
	return &smpAlertManagerResource{}
}

type smpAlertManagerResource struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (r *smpAlertManagerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_alert_manager"
}

func (r *smpAlertManagerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP AlertManager Resource",
		Attributes: map[string]schema.Attribute{
			"config_data": schema.StringAttribute{
				Description: "Base64 encoded YAML alert manager configuration file.",
				Required:    true,
			},
			"notification_group_id": schema.StringAttribute{
				Description: "Notification group ID.",
				Required:    true,
			},
			"alert_manager": schema.SingleNestedAttribute{
				Description: "Alert manager details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id":                    schema.StringAttribute{Computed: true},
					"config_data":           schema.StringAttribute{Computed: true},
					"state":                 schema.StringAttribute{Computed: true},
					"workspace_id":          schema.StringAttribute{Computed: true},
					"notification_group_id": schema.StringAttribute{Computed: true},
					"created_at":            schema.StringAttribute{Computed: true},
					"created_by":            schema.StringAttribute{Computed: true},
					"modified_at":           schema.StringAttribute{Computed: true},
					"modified_by":           schema.StringAttribute{Computed: true},
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID.",
				Required:    true,
			},
		},
	}
}

func (r *smpAlertManagerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(ErrUnexpectedConfigure, fmt.Sprintf(ErrUnexpectedConfigureFmt, req.ProviderData))
		return
	}
	r.client = inst.Client.Smp
	r.clients = inst.Client
}

func (r *smpAlertManagerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan smp.AlertManagerResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateAlertManager(ctx, plan.WorkspaceId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrCreateAlertManager, fmt.Sprintf(ErrCreateAlertManagerFmt, err.Error(), detail))
		return
	}

	am := result.GetAlertManager()
	amObj, diags := types.ObjectValue(
		smp.AlertManager{}.AttributeTypes(),
		map[string]attr.Value{
			"id":                    types.StringValue(am.GetId()),
			"config_data":           types.StringValue(am.GetConfigData()),
			"state":                 types.StringValue(am.GetState()),
			"workspace_id":          types.StringValue(am.GetWorkspaceId()),
			"notification_group_id": types.StringValue(am.GetNotificationGroupId()),
			"created_at":            types.StringValue(am.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":            types.StringValue(am.GetCreatedBy()),
			"modified_at":           nullableTimeTypes(am.GetModifiedAtOk()),
			"modified_by":           nullableStringTypes(am.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.AlertManager = amObj

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *smpAlertManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state smp.AlertManagerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetAlertManager(ctx, state.WorkspaceId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadAlertManager, fmt.Sprintf(ErrReadAlertManagerFmt, state.WorkspaceId.ValueString(), err.Error(), detail))
		return
	}

	am := result.GetAlertManager()
	amObj, diags := types.ObjectValue(
		smp.AlertManager{}.AttributeTypes(),
		map[string]attr.Value{
			"id":                    types.StringValue(am.GetId()),
			"config_data":           types.StringValue(am.GetConfigData()),
			"state":                 types.StringValue(am.GetState()),
			"workspace_id":          types.StringValue(am.GetWorkspaceId()),
			"notification_group_id": types.StringValue(am.GetNotificationGroupId()),
			"created_at":            types.StringValue(am.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":            types.StringValue(am.GetCreatedBy()),
			"modified_at":           nullableTimeTypes(am.GetModifiedAtOk()),
			"modified_by":           nullableStringTypes(am.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.AlertManager = amObj

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *smpAlertManagerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan smp.AlertManagerResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateAlertManager(ctx, plan.WorkspaceId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrUpdateAlertManager, fmt.Sprintf(ErrUpdateAlertManagerFmt, err.Error(), detail))
		return
	}

	result, err := r.client.GetAlertManager(ctx, plan.WorkspaceId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadAlertManager, fmt.Sprintf(ErrReadAlertManagerFmt, plan.WorkspaceId.ValueString(), err.Error(), detail))
		return
	}

	am := result.GetAlertManager()
	amObj, diags := types.ObjectValue(
		smp.AlertManager{}.AttributeTypes(),
		map[string]attr.Value{
			"id":                    types.StringValue(am.GetId()),
			"config_data":           types.StringValue(am.GetConfigData()),
			"state":                 types.StringValue(am.GetState()),
			"workspace_id":          types.StringValue(am.GetWorkspaceId()),
			"notification_group_id": types.StringValue(am.GetNotificationGroupId()),
			"created_at":            types.StringValue(am.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":            types.StringValue(am.GetCreatedBy()),
			"modified_at":           nullableTimeTypes(am.GetModifiedAtOk()),
			"modified_by":           nullableStringTypes(am.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.AlertManager = amObj

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *smpAlertManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state smp.AlertManagerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAlertManager(ctx, state.WorkspaceId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrDeleteAlertManager, fmt.Sprintf(ErrDeleteAlertManagerFmt, err.Error(), detail))
		return
	}
}

func (r *smpAlertManagerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("workspace_id"), req, resp)
}
