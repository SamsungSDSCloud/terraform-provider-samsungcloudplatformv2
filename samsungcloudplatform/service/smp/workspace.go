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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                = &smpWorkspaceResource{}
	_ resource.ResourceWithConfigure   = &smpWorkspaceResource{}
	_ resource.ResourceWithImportState = &smpWorkspaceResource{}
	_ resource.ResourceWithModifyPlan  = &smpWorkspaceResource{}
)

func NewSmpWorkspaceResource() resource.Resource {
	return &smpWorkspaceResource{}
}

type smpWorkspaceResource struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

// --------------------------------------------------------------------------
// METADATA
// --------------------------------------------------------------------------
func (r *smpWorkspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_workspace"
}

// --------------------------------------------------------------------------
// SCHEMA
// --------------------------------------------------------------------------
func (r *smpWorkspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP Workspace Resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Workspace name.",
				Required:    true,
			},
			"private_acl_enabled": schema.BoolAttribute{
				Description: "Enable private ACL access.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"private_acl_resources": schema.ListNestedAttribute{
				Description: "List of private ACL resources.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_id": schema.StringAttribute{
							Description: "Virtual server.",
							Required:    true,
						},
						"resource_name": schema.StringAttribute{
							Description: "Resource name.",
							Optional:    true,
						},
						"resource_type": schema.StringAttribute{
							Description: "Resource type: virtual_server.",
							Required:    true,
						},
					},
				},
			},
			"workspace": schema.SingleNestedAttribute{
				Description: "Workspace details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id":                         schema.StringAttribute{Computed: true},
					"name":                       schema.StringAttribute{Computed: true},
					"state":                      schema.StringAttribute{Computed: true},
					"account_id":                 schema.StringAttribute{Computed: true},
					"created_at":                 schema.StringAttribute{Computed: true},
					"created_by":                 schema.StringAttribute{Computed: true},
					"modified_at":                schema.StringAttribute{Computed: true},
					"modified_by":                schema.StringAttribute{Computed: true},
					"private_acl_enabled":        schema.BoolAttribute{Computed: true},
					"private_acl_resources": schema.ListNestedAttribute{
						Computed:     true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"resource_id":   schema.StringAttribute{Computed: true},
								"resource_name": schema.StringAttribute{Computed: true},
								"resource_type": schema.StringAttribute{Computed: true},
							},
						},
					},
					"private_prometheus_endpoint": schema.StringAttribute{Computed: true},
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// --------------------------------------------------------------------------
// MODIFY PLAN
// --------------------------------------------------------------------------
// ModifyPlan merges stable fields from the prior state into the plan so the
// Computed mirror object (workspace) does not show "(known after apply)" on
// every update. Mutable fields are only marked unknown when the input
// actually changed.
func (r *smpWorkspaceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip plan modification when destroying the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	// Skip if there's no existing state (create)
	if req.State.Raw.IsNull() {
		return
	}

	var plan smp.WorkspaceResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state smp.WorkspaceResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Workspace.IsNull() || state.Workspace.IsUnknown() {
		return
	}

	var stateWs smp.Workspace
	resp.Diagnostics.Append(state.Workspace.As(ctx, &stateWs, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only mark mutable fields as unknown when the input actually changes
	changed := !plan.Name.Equal(state.Name) ||
		!plan.PrivateAclEnabled.Equal(state.PrivateAclEnabled) ||
		!plan.PrivateAclResources.Equal(state.PrivateAclResources)

	modelWs := smp.Workspace{
		Id:                        stateWs.Id,
		Name:                      plan.Name,
		State:                     stateWs.State,
		AccountId:                 stateWs.AccountId,
		CreatedAt:                 stateWs.CreatedAt,
		CreatedBy:                 stateWs.CreatedBy,
		PrivateAclEnabled:         plan.PrivateAclEnabled,
		PrivateAclResources:       plan.PrivateAclResources,
		PrivatePrometheusEndpoint: stateWs.PrivatePrometheusEndpoint,
	}

	if changed {
		modelWs.ModifiedAt = types.StringUnknown()
		modelWs.ModifiedBy = types.StringUnknown()
		modelWs.PrivatePrometheusEndpoint = types.StringUnknown()
	} else {
		modelWs.ModifiedAt = stateWs.ModifiedAt
		modelWs.ModifiedBy = stateWs.ModifiedBy
	}

	mergedObj, diags := types.ObjectValueFrom(ctx, modelWs.AttributeTypes(), modelWs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Workspace = mergedObj
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

// --------------------------------------------------------------------------
// CONFIGURE
// --------------------------------------------------------------------------
func (r *smpWorkspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *smpWorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan smp.WorkspaceResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateWorkspace(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrCreateWorkspace,
			fmt.Sprintf(ErrCreateWorkspaceFmt, err.Error(), detail))
		return
	}

	workspace := result.GetWorkspace()
	workspaceObj, diags := types.ObjectValue(
		smp.Workspace{}.AttributeTypes(),
		map[string]attr.Value{
			"id":                          types.StringValue(workspace.GetId()),
			"name":                        types.StringValue(workspace.GetName()),
			"state":                       types.StringValue(workspace.GetState()),
			"account_id":                  types.StringValue(workspace.GetAccountId()),
			"created_at":                  types.StringValue(workspace.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":                  types.StringValue(workspace.GetCreatedBy()),
			"modified_at":                 nullableTimeTypes(workspace.GetModifiedAtOk()),
			"modified_by":                 nullableStringTypes(workspace.GetModifiedByOk()),
			"private_acl_enabled":         nullableBoolTypes(workspace.GetPrivateAclEnabledOk()),
			"private_acl_resources":       privateAclResourcesToList(ctx, workspace.GetPrivateAclResources(), &resp.Diagnostics),
			"private_prometheus_endpoint": nullableStringTypes(workspace.GetPrivatePrometheusEndpointOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.WorkspaceId = types.StringValue(workspace.GetId())
	plan.Workspace = workspaceObj

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// --------------------------------------------------------------------------
// READ
// --------------------------------------------------------------------------
func (r *smpWorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state smp.WorkspaceResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetWorkspace(ctx, state.WorkspaceId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadWorkspace,
			fmt.Sprintf(ErrReadWorkspaceFmt,
				state.WorkspaceId.ValueString(),
				err.Error(),
				detail))
		return
	}

	workspace := result.GetWorkspace()
	workspaceObj, diags := types.ObjectValue(
		smp.Workspace{}.AttributeTypes(),
		map[string]attr.Value{
			"id":                          types.StringValue(workspace.GetId()),
			"name":                        types.StringValue(workspace.GetName()),
			"state":                       types.StringValue(workspace.GetState()),
			"account_id":                  types.StringValue(workspace.GetAccountId()),
			"created_at":                  types.StringValue(workspace.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":                  types.StringValue(workspace.GetCreatedBy()),
			"modified_at":                 nullableTimeTypes(workspace.GetModifiedAtOk()),
			"modified_by":                 nullableStringTypes(workspace.GetModifiedByOk()),
			"private_acl_enabled":         nullableBoolTypes(workspace.GetPrivateAclEnabledOk()),
			"private_acl_resources":       privateAclResourcesToList(ctx, workspace.GetPrivateAclResources(), &resp.Diagnostics),
			"private_prometheus_endpoint": nullableStringTypes(workspace.GetPrivatePrometheusEndpointOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.WorkspaceId = types.StringValue(workspace.GetId())
	state.Workspace = workspaceObj

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// --------------------------------------------------------------------------
// UPDATE
// --------------------------------------------------------------------------
func (r *smpWorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan smp.WorkspaceResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateWorkspace(ctx, plan.WorkspaceId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrUpdateWorkspace,
			fmt.Sprintf(ErrUpdateWorkspaceFmt, err.Error(), detail))
		return
	}

	// Update private ACL resources if they changed
	if !plan.PrivateAclResources.IsNull() && !plan.PrivateAclResources.IsUnknown() {
		_, err := r.client.SetWorkspacePrivateAclResources(ctx, plan.WorkspaceId.ValueString(), smp.PrivateAclResourcesSetRequest{
			PrivateAclEnabled:   plan.PrivateAclEnabled,
			PrivateAclResources: plan.PrivateAclResources,
		})
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError("Unable to update private ACL resources",
				fmt.Sprintf("Error updating private ACL resources: %s, Detail: %s", err.Error(), detail))
			return
		}
	}

	result, err := r.client.GetWorkspace(ctx, plan.WorkspaceId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadWorkspace,
			fmt.Sprintf(ErrReadWorkspaceFmt,
				plan.WorkspaceId.ValueString(),
				err.Error(),
				detail))
		return
	}

	workspace := result.GetWorkspace()
	workspaceObj, diags := types.ObjectValue(
		smp.Workspace{}.AttributeTypes(),
		map[string]attr.Value{
			"id":                          types.StringValue(workspace.GetId()),
			"name":                        types.StringValue(workspace.GetName()),
			"state":                       types.StringValue(workspace.GetState()),
			"account_id":                  types.StringValue(workspace.GetAccountId()),
			"created_at":                  types.StringValue(workspace.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":                  types.StringValue(workspace.GetCreatedBy()),
			"modified_at":                 nullableTimeTypes(workspace.GetModifiedAtOk()),
			"modified_by":                 nullableStringTypes(workspace.GetModifiedByOk()),
			"private_acl_enabled":         nullableBoolTypes(workspace.GetPrivateAclEnabledOk()),
			"private_acl_resources":       privateAclResourcesToList(ctx, workspace.GetPrivateAclResources(), &resp.Diagnostics),
			"private_prometheus_endpoint": nullableStringTypes(workspace.GetPrivatePrometheusEndpointOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.WorkspaceId = types.StringValue(workspace.GetId())
	plan.Workspace = workspaceObj

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// --------------------------------------------------------------------------
// DELETE
// --------------------------------------------------------------------------
func (r *smpWorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state smp.WorkspaceResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWorkspace(ctx, state.WorkspaceId.ValueString())
	if err != nil {
	  detail := client.GetDetailFromError(err)
	  resp.Diagnostics.AddError(ErrDeleteWorkspace,
	      fmt.Sprintf(ErrDeleteWorkspaceFmt, err.Error(), detail))
	  return
   }
}

// --------------------------------------------------------------------------
// IMPORT
// --------------------------------------------------------------------------
func (r *smpWorkspaceResource) ImportState(ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("workspace_id"), req, resp)
}
