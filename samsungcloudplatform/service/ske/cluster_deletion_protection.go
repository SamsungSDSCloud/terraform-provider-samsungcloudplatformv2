package ske

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &skeClusterDeletionProtectionResource{}
	_ resource.ResourceWithConfigure   = &skeClusterDeletionProtectionResource{}
	_ resource.ResourceWithImportState = &skeClusterDeletionProtectionResource{}
)

func NewSkeClusterDeletionProtectionResource() resource.Resource {
	return &skeClusterDeletionProtectionResource{}
}

type skeClusterDeletionProtectionResource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

func (r *skeClusterDeletionProtectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_cluster_deletion_protection"
}

func (r *skeClusterDeletionProtectionResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages deletion protection for an SKE cluster.",
		MarkdownDescription: "Manages deletion protection for an SKE cluster. When enabled, the cluster cannot be deleted until protection is disabled.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the resource.\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				MarkdownDescription: "Identifier of the resource.\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Description:         "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				MarkdownDescription: "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"deletion_protection_enabled": schema.BoolAttribute{
				Description:         "Deletion protection flag. When set to true, cluster deletion is prevented.\n  - example: true",
				MarkdownDescription: "Deletion protection flag. When set to true, cluster deletion is prevented.\n  - example: true",
				Required:            true,
			},
			"resource_id": schema.StringAttribute{
				Description:         "Resource ID returned by the API, or cluster_id as a fallback when the API does not provide one\n  - example: 70a599e031e749b7b260868f441e862b",
				MarkdownDescription: "Resource ID returned by the API, or cluster_id as a fallback when the API does not provide one\n  - example: 70a599e031e749b7b260868f441e862b",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *skeClusterDeletionProtectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Ske
	r.clients = inst.Client
}

// resolveResourceId returns the API-provided resource_id, falling back to cluster_id
func resolveResourceId(apiResourceId string, clusterId types.String) types.String {
	if apiResourceId == "" {
		return clusterId
	}
	return types.StringValue(apiResourceId)
}

func (r *skeClusterDeletionProtectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan skeClusterDeletionProtectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, _, err := r.client.UpdateClusterDeletionProtection(ctx, plan.ClusterId.ValueString(), ske.ClusterResource{
		DeletionProtectionEnabled: types.BoolValue(plan.DeletionProtectionEnabled.ValueBool()),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SKE Cluster Deletion Protection",
			"Could not set deletion protection, unexpected error: "+err.Error(),
		)
		return
	}

	// Read back to ensure state matches server
	deletionProtection, httpStatus, err := r.client.GetClusterDeletionProtection(ctx, plan.ClusterId.ValueString())
	if err != nil {
		if httpStatus == http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Error Creating SKE Cluster Deletion Protection",
				"Cluster not found after creation: "+err.Error(),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Creating SKE Cluster Deletion Protection",
			"Could not read back deletion protection: "+err.Error(),
		)
		return
	}

	plan.Id = plan.ClusterId
	plan.ResourceId = resolveResourceId(result.ResourceId, plan.ClusterId)
	plan.DeletionProtectionEnabled = types.BoolValue(deletionProtection)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *skeClusterDeletionProtectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state skeClusterDeletionProtectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deletionProtection, httpStatus, err := r.client.GetClusterDeletionProtection(ctx, state.ClusterId.ValueString())
	if err != nil {
		if httpStatus == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SKE Cluster Deletion Protection",
			"Could not read deletion protection for cluster ID "+state.ClusterId.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Id = state.ClusterId
	state.DeletionProtectionEnabled = types.BoolValue(deletionProtection)

	// API GET doesn't return resource_id — set to cluster_id only if state value is missing
	if state.ResourceId.IsNull() || state.ResourceId.ValueString() == "" {
		state.ResourceId = state.ClusterId
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *skeClusterDeletionProtectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan skeClusterDeletionProtectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, _, err := r.client.UpdateClusterDeletionProtection(ctx, plan.ClusterId.ValueString(), ske.ClusterResource{
		DeletionProtectionEnabled: types.BoolValue(plan.DeletionProtectionEnabled.ValueBool()),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SKE Cluster Deletion Protection",
			"Could not update deletion protection, unexpected error: "+err.Error(),
		)
		return
	}

	// Read back to ensure state matches server
	deletionProtection, httpStatus, err := r.client.GetClusterDeletionProtection(ctx, plan.ClusterId.ValueString())
	if err != nil {
		if httpStatus == http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Error Updating SKE Cluster Deletion Protection",
				"Cluster not found after update: "+err.Error(),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Updating SKE Cluster Deletion Protection",
			"Could not read back deletion protection: "+err.Error(),
		)
		return
	}

	plan.Id = plan.ClusterId
	plan.ResourceId = resolveResourceId(result.ResourceId, plan.ClusterId)
	plan.DeletionProtectionEnabled = types.BoolValue(deletionProtection)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *skeClusterDeletionProtectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state skeClusterDeletionProtectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Disable deletion protection on delete (skip if already disabled)
	if state.DeletionProtectionEnabled.ValueBool() {
		_, httpStatus, err := r.client.UpdateClusterDeletionProtection(ctx, state.ClusterId.ValueString(), ske.ClusterResource{
			DeletionProtectionEnabled: types.BoolValue(false),
		})
		if err != nil && httpStatus != http.StatusNotFound {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Deleting SKE Cluster Deletion Protection",
				"Could not disable deletion protection, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}
}

func (r *skeClusterDeletionProtectionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

type skeClusterDeletionProtectionResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	ClusterId          types.String `tfsdk:"cluster_id"`
	DeletionProtectionEnabled types.Bool   `tfsdk:"deletion_protection_enabled"`
	ResourceId         types.String `tfsdk:"resource_id"`
}
