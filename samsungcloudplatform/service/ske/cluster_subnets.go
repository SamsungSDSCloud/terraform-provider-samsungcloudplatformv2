package ske

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &skeClusterSubnetsResource{}
	_ resource.ResourceWithConfigure   = &skeClusterSubnetsResource{}
	_ resource.ResourceWithImportState = &skeClusterSubnetsResource{}
)

func NewSkeClusterSubnetsResource() resource.Resource {
	return &skeClusterSubnetsResource{}
}

type skeClusterSubnetsResource struct {
	config  *ske.Client
	clients *client.SCPClient
}

func (r *skeClusterSubnetsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_cluster_subnets"
}

func (r *skeClusterSubnetsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Cluster Subnets",
		MarkdownDescription: "Manages additional subnet ID list of a SKE cluster via PUT /v1/clusters/{cluster_id}/subnets (REPLACE semantics).",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster ID",
				MarkdownDescription: "The ID of the cluster whose additional subnets are managed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"additional_subnet_id_list": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Description:         "Additional subnet ID list (max 2)",
				MarkdownDescription: "List of additional subnet IDs to attach to the cluster. REPLACES the entire list on update. Max 2 items.",
				Validators: []validator.Set{
					setvalidator.SizeAtMost(2),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID (= cluster_id)",
				MarkdownDescription: "The ID of the resource. Equals the cluster_id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *skeClusterSubnetsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.config = inst.Client.Ske
	r.clients = inst.Client
}

func (r *skeClusterSubnetsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ske.ClusterSubnetsResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	additionalSubnetIdList := setToStringSlice(plan.AdditionalSubnetIdList)

	_, err := r.config.SetClusterSubnets(ctx, plan.ClusterId.ValueString(), additionalSubnetIdList)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Setting Cluster Subnets",
			"Could not set cluster subnets, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForClusterStatus(ctx, r.config, plan.ClusterId.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Setting Cluster Subnets",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return
	}

	plan.Id = plan.ClusterId
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

func (r *skeClusterSubnetsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ske.ClusterSubnetsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, httpStatus, err := r.config.GetClusterSubnets(ctx, state.ClusterId.ValueString())
	if err != nil {
		if httpStatus == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster Subnets",
			"Could not read cluster subnets for cluster ID "+state.ClusterId.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	state.AdditionalSubnetIdList = toTypesStringSet(list)
	state.Id = state.ClusterId

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *skeClusterSubnetsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ske.ClusterSubnetsResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ske.ClusterSubnetsResource
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	planSet, diagsSet := types.SetValueFrom(ctx, types.StringType, plan.AdditionalSubnetIdList)
	resp.Diagnostics.Append(diagsSet...)
	stateSet, diagsSet := types.SetValueFrom(ctx, types.StringType, state.AdditionalSubnetIdList)
	resp.Diagnostics.Append(diagsSet...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planSet.Equal(stateSet) {
		_, err := r.config.SetClusterSubnets(ctx, plan.ClusterId.ValueString(), setToStringSlice(plan.AdditionalSubnetIdList))
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Updating Cluster Subnets",
				"Could not update cluster subnets, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}

		err = waitForClusterStatus(ctx, r.config, plan.ClusterId.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Cluster Subnets",
				"Error waiting for cluster to become running: "+err.Error(),
			)
			return
		}
	}

	plan.Id = plan.ClusterId
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

func (r *skeClusterSubnetsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ske.ClusterSubnetsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.config.SetClusterSubnets(ctx, state.ClusterId.ValueString(), []string{})
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Clearing Cluster Subnets",
			"Could not clear cluster subnets, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForClusterStatus(ctx, r.config, state.ClusterId.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Clearing Cluster Subnets",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return
	}
}

func (r *skeClusterSubnetsResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// ---- helpers ----

func setToStringSlice(values types.Set) []string {
	if values.IsNull() || values.IsUnknown() {
		return []string{}
	}
	result := make([]string, 0, len(values.Elements()))
	for _, v := range values.Elements() {
		result = append(result, v.(types.String).ValueString())
	}
	return result
}

func toTypesStringSet(values []string) types.Set {
	if values == nil {
		return types.SetNull(types.StringType)
	}
	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elems = append(elems, types.StringValue(v))
	}
	result, _ := types.SetValue(types.StringType, elems)
	return result
}
