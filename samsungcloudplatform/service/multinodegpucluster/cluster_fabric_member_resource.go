package multinodegpucluster

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	multinodegpuclusterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/multinodegpucluster"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &ClusterFabricMemberResource{}
	_ resource.ResourceWithConfigure = &ClusterFabricMemberResource{}
)

func NewClusterFabricMemberResource() resource.Resource {
	return &ClusterFabricMemberResource{}
}

type ClusterFabricMemberResource struct {
	config  *scpsdk.Configuration
	client  *multinodegpuclusterClient.Client
	clients *client.SCPClient
}

func (multinodegpuclusterRS *ClusterFabricMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multinodegpucluster_cluster_fabric_member"
}

func (multinodegpuclusterRS *ClusterFabricMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expect *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	multinodegpuclusterRS.client = inst.Client.Mngc
	multinodegpuclusterRS.clients = inst.Client
}

func (multinodegpuclusterRS *ClusterFabricMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterFabricMemberResourceSchema()
}

func ClusterFabricMemberResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Cluster Fabric Member",
		Attributes: map[string]schema.Attribute{
			"after_cluster_fabric_id": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster Fabric ID that the GPU Nodes will belong to after the change\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
				MarkdownDescription: "Cluster Fabric ID that the GPU Nodes will belong to after the change\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"before_cluster_fabric_id": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster Fabric ID that the GPU Nodes belong to before the change\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
				MarkdownDescription: "Cluster Fabric ID that the GPU Nodes belong to before the change\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"gpu_node_id_list": schema.ListAttribute{
				Required:            true,
				ElementType:         types.StringType,
				Description:         "GPU Node ID List to move between Cluster Fabrics. Growing this list on an existing resource moves only the newly added GPU Nodes; nodes already moved by a prior apply are not resent (they are no longer members of before_cluster_fabric_id, so resending them is rejected by the API).",
				MarkdownDescription: "GPU Node ID List to move between Cluster Fabrics. Growing this list on an existing resource moves only the newly added GPU Nodes; nodes already moved by a prior apply are not resent (they are no longer members of before_cluster_fabric_id, so resending them is rejected by the API).",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Cluster Fabric ID that the GPU Nodes belong to after the change",
				MarkdownDescription: "Cluster Fabric ID that the GPU Nodes belong to after the change",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (multinodegpuclusterRS *ClusterFabricMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan multinodegpuclusterClient.ClusterFabricMember

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	gpuNodeIdList := make([]string, 0, len(plan.GpuNodeIdList))
	for _, gpuNodeId := range plan.GpuNodeIdList {
		gpuNodeIdList = append(gpuNodeIdList, gpuNodeId.ValueString())
	}

	_, err := multinodegpuclusterRS.client.ModifyClusterFabricMembers(ctx, plan.BeforeClusterFabricId.ValueString(), plan.AfterClusterFabricId.ValueString(), gpuNodeIdList)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error modifying Cluster Fabric members",
			"Could not modify Cluster Fabric members, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	plan.Id = plan.AfterClusterFabricId

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (multinodegpuclusterRS *ClusterFabricMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state multinodegpuclusterClient.ClusterFabricMember

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResponse, err := multinodegpuclusterRS.client.GetClusterFabric(ctx, state.AfterClusterFabricId.ValueString())
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster Fabric",
			"Could not read Cluster Fabric ID "+state.AfterClusterFabricId.ValueString()+": "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (multinodegpuclusterRS *ClusterFabricMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan multinodegpuclusterClient.ClusterFabricMember

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state multinodegpuclusterClient.ClusterFabricMember

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	alreadyMoved := make(map[string]bool, len(state.GpuNodeIdList))
	for _, gpuNodeId := range state.GpuNodeIdList {
		alreadyMoved[gpuNodeId.ValueString()] = true
	}

	newGpuNodeIdList := make([]string, 0, len(plan.GpuNodeIdList))
	for _, gpuNodeId := range plan.GpuNodeIdList {
		id := gpuNodeId.ValueString()
		if !alreadyMoved[id] {
			newGpuNodeIdList = append(newGpuNodeIdList, id)
		}
	}

	plan.Id = plan.AfterClusterFabricId

	if len(newGpuNodeIdList) == 0 {
		diags = resp.State.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		return
	}

	_, err := multinodegpuclusterRS.client.ModifyClusterFabricMembers(ctx, plan.BeforeClusterFabricId.ValueString(), plan.AfterClusterFabricId.ValueString(), newGpuNodeIdList)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error modifying Cluster Fabric members",
			"Could not modify Cluster Fabric members, unexpected error: "+err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (multinodegpuclusterRS *ClusterFabricMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete not supported",
		"Cluster Fabric member changes cannot be reverted. It is only removed from the Terraform state.",
	)
}
