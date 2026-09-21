package ske

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/ske/1.6"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &skeClusterLinkedResourcesResource{}
	_ resource.ResourceWithConfigure   = &skeClusterLinkedResourcesResource{}
	_ resource.ResourceWithImportState = &skeClusterLinkedResourcesResource{}
)

// skeClusterLinkedResourcesResource is the resource implementation.
type skeClusterLinkedResourcesResource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

// ClusterLinkedResourcesResourceModel maps the schema to a Go struct.
type ClusterLinkedResourcesResourceModel struct {
	Id              types.String         `tfsdk:"id"`
	ClusterId       types.String         `tfsdk:"cluster_id"`
	LinkedResources []ske.LinkedResource `tfsdk:"linked_resources"`
}

// NewSkeClusterLinkedResourcesResource is a helper function to simplify the provider implementation.
func NewSkeClusterLinkedResourcesResource() resource.Resource {
	return &skeClusterLinkedResourcesResource{}
}

// Metadata returns the resource type name.
func (r *skeClusterLinkedResourcesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_cluster_linked_resources"
}

// Schema defines the schema for the resource.
func (r *skeClusterLinkedResourcesResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterLinkedResourcesResourceSchema()
}

func ClusterLinkedResourcesResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages linked resources for an SKE cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource. Uses cluster_id as the resource ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
				MarkdownDescription: "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"linked_resources": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource ID\n  - example: res-12345678",
							MarkdownDescription: "Linked Resource ID\n  - example: res-12345678",
						},
						"name": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource Name\n  - example: my-resource",
							MarkdownDescription: "Linked Resource Name\n  - example: my-resource",
						},
						"type": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource Type (fs/obs)\n  - pattern: fs|obs\n  - example: fs",
							MarkdownDescription: "Linked Resource Type (fs/obs)\n  - pattern: fs|obs\n  - example: fs",
							Validators: []validator.String{
								stringvalidator.OneOf("fs", "obs"),
							},
						},
					},
				},
				Required:            true,
				Description:         "List of linked resources associated with the cluster\n  - example: {id='res-12345678', name='my-resource', type='fs'}",
				MarkdownDescription: "List of linked resources associated with the cluster\n  - example: {id='res-12345678', name='my-resource', type='fs'}",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure sets up the resource with the provider configuration.
func (r *skeClusterLinkedResourcesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create sets the linked resources on a cluster.
func (r *skeClusterLinkedResourcesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClusterLinkedResourcesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateClusterLinkedResources(ctx, plan.ClusterId.ValueString(), plan.LinkedResources)

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Setting Cluster Linked Resources",
			"Could not set linked resources, unexpected error: "+err.Error()+"\nReason: "+detail,
		)

		return
	}

	plan.Id = types.StringValue(plan.ClusterId.ValueString())

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

// Read refreshes the Terraform state with the latest data.
func (r *skeClusterLinkedResourcesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClusterLinkedResourcesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, httpStatus, err := r.client.GetCluster(ctx, state.ClusterId.ValueString())

	if err != nil {
		if httpStatus == 404 {
			resp.State.RemoveResource(ctx)
			return
		}

		detail := client.GetDetailFromError(err)

		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			"Could not read cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)

		return
	}

	if data != nil {
		cluster := data.Cluster
		state.ClusterId = types.StringValue(cluster.Id)
		state.Id = types.StringValue(cluster.Id)
		state.LinkedResources = convertLinkedResourcesFromSDK(cluster.LinkedResources)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
}

// Update updates the linked resources on a cluster.
func (r *skeClusterLinkedResourcesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ClusterLinkedResourcesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateClusterLinkedResources(ctx, plan.ClusterId.ValueString(), plan.LinkedResources)

	if err != nil {
		detail := client.GetDetailFromError(err)

		resp.Diagnostics.AddError(
			"Error Updating Cluster Linked Resources",
			"Could not update linked resources, unexpected error: "+err.Error()+"\nReason: "+detail,
		)

		return
	}

	plan.Id = types.StringValue(plan.ClusterId.ValueString())

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

// Delete removes all linked resources from a cluster by setting an empty list.
func (r *skeClusterLinkedResourcesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClusterLinkedResourcesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateClusterLinkedResources(ctx, state.ClusterId.ValueString(), []ske.LinkedResource{})

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Cluster Linked Resources",
			"Could not delete linked resources, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

// ImportState imports an existing cluster's linked resources.
func (r *skeClusterLinkedResourcesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// convertLinkedResourcesFromSDK converts SDK LinkedResource models to client LinkedResource models.
func convertLinkedResourcesFromSDK(sdkResources []scpske.LinkedResource) []ske.LinkedResource {
	if sdkResources == nil {
		return nil
	}
	result := make([]ske.LinkedResource, len(sdkResources))

	for i, lr := range sdkResources {
		result[i] = ske.LinkedResource{
			Id:   types.StringValue(lr.Id),
			Name: types.StringValue(lr.Name),
			Type: types.StringValue(lr.Type),
		}
	}

	return result
}
