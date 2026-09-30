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
	_ resource.Resource                = &smpRuleNamespaceResource{}
	_ resource.ResourceWithConfigure   = &smpRuleNamespaceResource{}
	_ resource.ResourceWithImportState = &smpRuleNamespaceResource{}
)

func NewSmpRuleNamespaceResource() resource.Resource {
	return &smpRuleNamespaceResource{}
}

type smpRuleNamespaceResource struct {
	config  *scpsdk.Configuration
	client  *smp.Client
	clients *client.SCPClient
}

func (r *smpRuleNamespaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smp_rule_namespace"
}

func (r *smpRuleNamespaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SMP RuleNamespace Resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Rule namespace name.",
				Required:    true,
			},
			"config_data": schema.StringAttribute{
				Description: "Base64 encoded YAML rules file.",
				Required:    true,
			},
			"rule_namespace": schema.SingleNestedAttribute{
				Description: "Rule namespace details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id":           schema.StringAttribute{Computed: true},
					"name":         schema.StringAttribute{Computed: true},
					"config_data":  schema.StringAttribute{Computed: true},
					"state":        schema.StringAttribute{Computed: true},
					"workspace_id": schema.StringAttribute{Computed: true},
					"created_at":   schema.StringAttribute{Computed: true},
					"created_by":   schema.StringAttribute{Computed: true},
					"modified_at":  schema.StringAttribute{Computed: true},
					"modified_by":  schema.StringAttribute{Computed: true},
				},
			},
			"rule_namespace_id": schema.StringAttribute{
				Description: "Rule namespace ID.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID.",
				Required:    true,
			},
		},
	}
}

func (r *smpRuleNamespaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *smpRuleNamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan smp.RuleNamespaceResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateRuleNamespace(ctx, plan.WorkspaceId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrCreateRuleNamespace, fmt.Sprintf(ErrCreateRuleNamespaceFmt, err.Error(), detail))
		return
	}

	rn := result.GetRuleNamespace()
	rnObj, diags := types.ObjectValue(
		smp.RuleNamespace{}.AttributeTypes(),
		map[string]attr.Value{
			"id":           types.StringValue(rn.GetId()),
			"name":         types.StringValue(rn.GetName()),
			"config_data":  types.StringValue(rn.GetConfigData()),
			"state":        types.StringValue(rn.GetState()),
			"workspace_id": types.StringValue(rn.GetWorkspaceId()),
			"created_at":   types.StringValue(rn.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":   types.StringValue(rn.GetCreatedBy()),
			"modified_at":  nullableTimeTypes(rn.GetModifiedAtOk()),
			"modified_by":  nullableStringTypes(rn.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.RuleNamespaceId = types.StringValue(rn.GetId())
	plan.RuleNamespace = rnObj

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *smpRuleNamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state smp.RuleNamespaceResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetRuleNamespace(ctx, state.WorkspaceId.ValueString(), state.RuleNamespaceId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadRuleNamespace, fmt.Sprintf(ErrReadRuleNamespaceFmt, state.RuleNamespaceId.ValueString(), err.Error(), detail))
		return
	}

	rn := result.GetRuleNamespace()
	rnObj, diags := types.ObjectValue(
		smp.RuleNamespace{}.AttributeTypes(),
		map[string]attr.Value{
			"id":           types.StringValue(rn.GetId()),
			"name":         types.StringValue(rn.GetName()),
			"config_data":  types.StringValue(rn.GetConfigData()),
			"state":        types.StringValue(rn.GetState()),
			"workspace_id": types.StringValue(rn.GetWorkspaceId()),
			"created_at":   types.StringValue(rn.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":   types.StringValue(rn.GetCreatedBy()),
			"modified_at":  nullableTimeTypes(rn.GetModifiedAtOk()),
			"modified_by":  nullableStringTypes(rn.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.RuleNamespaceId = types.StringValue(rn.GetId())
	state.RuleNamespace = rnObj

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *smpRuleNamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan smp.RuleNamespaceResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateRuleNamespace(ctx, plan.WorkspaceId.ValueString(), plan.RuleNamespaceId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrUpdateRuleNamespace, fmt.Sprintf(ErrUpdateRuleNamespaceFmt, err.Error(), detail))
		return
	}

	result, err := r.client.GetRuleNamespace(ctx, plan.WorkspaceId.ValueString(), plan.RuleNamespaceId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrReadRuleNamespace, fmt.Sprintf(ErrReadRuleNamespaceFmt, plan.RuleNamespaceId.ValueString(), err.Error(), detail))
		return
	}

	rn := result.GetRuleNamespace()
	rnObj, diags := types.ObjectValue(
		smp.RuleNamespace{}.AttributeTypes(),
		map[string]attr.Value{
			"id":           types.StringValue(rn.GetId()),
			"name":         types.StringValue(rn.GetName()),
			"config_data":  types.StringValue(rn.GetConfigData()),
			"state":        types.StringValue(rn.GetState()),
			"workspace_id": types.StringValue(rn.GetWorkspaceId()),
			"created_at":   types.StringValue(rn.GetCreatedAt().Format(TimeFormatDisplay)),
			"created_by":   types.StringValue(rn.GetCreatedBy()),
			"modified_at":  nullableTimeTypes(rn.GetModifiedAtOk()),
			"modified_by":  nullableStringTypes(rn.GetModifiedByOk()),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.RuleNamespaceId = types.StringValue(rn.GetId())
	plan.RuleNamespace = rnObj

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *smpRuleNamespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state smp.RuleNamespaceResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRuleNamespace(ctx, state.WorkspaceId.ValueString(), state.RuleNamespaceId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(ErrDeleteRuleNamespace, fmt.Sprintf(ErrDeleteRuleNamespaceFmt, err.Error(), detail))
		return
	}
}

func (r *smpRuleNamespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("rule_namespace_id"), req, resp)
}
