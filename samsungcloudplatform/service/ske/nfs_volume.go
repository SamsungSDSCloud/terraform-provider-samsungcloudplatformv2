package ske

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpske16 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/ske/1.6"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &skeClusterNfsVolumeResource{}
	_ resource.ResourceWithConfigure   = &skeClusterNfsVolumeResource{}
	_ resource.ResourceWithImportState = &skeClusterNfsVolumeResource{}
)

// NewSkeClusterNfsVolumeResource is a helper function to simplify the provider implementation.
func NewSkeClusterNfsVolumeResource() resource.Resource {
	return &skeClusterNfsVolumeResource{}
}

// skeClusterNfsVolumeResource is the resource implementation.
type skeClusterNfsVolumeResource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *skeClusterNfsVolumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_cluster_nfs_volume"
}

// Schema defines the schema for the resource.
func (r *skeClusterNfsVolumeResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attaches or detaches an NFS volume to/from an SKE cluster.\n" +
			"  - API: PUT /v1/clusters/{cluster_id}/nfs-volume\n" +
			"  - Note: Mutually exclusive with the `nfs_volume_id` argument of the `samsungcloudplatformv2_ske_cluster` resource for the same cluster. Use only one of them to manage the NFS volume of a cluster.",
		MarkdownDescription: "Attaches or detaches an NFS volume to/from an SKE cluster.\n" +
			"  - API: PUT /v1/clusters/{cluster_id}/nfs-volume\n" +
			"  - Note: Mutually exclusive with the `nfs_volume_id` argument of the `samsungcloudplatformv2_ske_cluster` resource for the same cluster. Use only one of them to manage the NFS volume of a cluster.",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				MarkdownDescription: "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nfs_volume_id": schema.StringAttribute{
				Required:            true,
				Description:         "NFS Volume ID to attach to the cluster\n  - example: bfdbabf2-04d9-4e8b-a205-020f8e6da438",
				MarkdownDescription: "NFS Volume ID to attach to the cluster\n  - example: bfdbabf2-04d9-4e8b-a205-020f8e6da438",
			},
			"id": schema.StringAttribute{
				Description: "Identifier of the resource (cluster_id).\n - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						MarkdownDescription: "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Name\n  - example: sample-cluster",
						MarkdownDescription: "Cluster Name\n  - example: sample-cluster",
					},
					"status": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Status\n  - pattern: RUNNING|CREATING|UPDATING|DELETING\n  - example: RUNNING",
						MarkdownDescription: "Cluster Status\n  - pattern: RUNNING|CREATING|UPDATING|DELETING\n  - example: RUNNING",
					},
					"nfs_volume_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Attached NFS Volume ID\n  - example: bfdbabf2-04d9-4e8b-a205-020f8e6da438",
						MarkdownDescription: "Attached NFS Volume ID\n  - example: bfdbabf2-04d9-4e8b-a205-020f8e6da438",
					},
				},
				Computed:    true,
				Description: "Cluster information after NFS volume operation",
				MarkdownDescription: "Cluster information after NFS volume operation",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *skeClusterNfsVolumeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *skeClusterNfsVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ske.ClusterNfsVolumeResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := plan.ClusterId.ValueString()
	nfsVolumeId := plan.NfsVolumeId.ValueString()

	// Set NFS volume on the cluster
	_, err := r.client.SetClusterNfsVolume(ctx, clusterId, nfsVolumeId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Setting NFS Volume",
			"Could not set NFS volume on cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(clusterId)

	// Wait for cluster to return to RUNNING status
	err = waitForClusterStatus(ctx, r.client, clusterId, []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Setting NFS Volume",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return
	}

	// Read back the cluster state
	clusterData, _, err := r.client.GetCluster(ctx, clusterId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			"Could not read cluster after NFS volume operation: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	clusterObj, diags := r.makeClusterNfsVolumeModel(ctx, (*scpske16.ClusterV1Dot6)(&clusterData.Cluster))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Cluster = clusterObj

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *skeClusterNfsVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ske.ClusterNfsVolumeResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.ClusterId.ValueString()

	// Get refreshed cluster value
	data, httpStatus, err := r.client.GetCluster(ctx, clusterId)
	if err != nil {
		if httpStatus == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster NFS Volume",
			"Could not read cluster ID "+clusterId+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	clusterObj, diags := r.makeClusterNfsVolumeModel(ctx, (*scpske16.ClusterV1Dot6)(&data.Cluster))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.Cluster.NfsVolumeId.Get() != nil {
		state.NfsVolumeId = types.StringValue(*data.Cluster.NfsVolumeId.Get())
	} else {
		state.NfsVolumeId = types.StringNull()
	}

	state.Id = types.StringValue(clusterId)
	state.Cluster = clusterObj

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *skeClusterNfsVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ske.ClusterNfsVolumeResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := plan.ClusterId.ValueString()
	nfsVolumeId := plan.NfsVolumeId.ValueString()

	// Update NFS volume on the cluster
	_, err := r.client.SetClusterNfsVolume(ctx, clusterId, nfsVolumeId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating NFS Volume",
			"Could not update NFS volume on cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(clusterId)

	// Wait for cluster to return to RUNNING status
	err = waitForClusterStatus(ctx, r.client, clusterId, []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating NFS Volume",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return
	}

	// Read back the cluster state
	clusterData, _, err := r.client.GetCluster(ctx, clusterId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			"Could not read cluster after NFS volume update: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	clusterObj, diags := r.makeClusterNfsVolumeModel(ctx, (*scpske16.ClusterV1Dot6)(&clusterData.Cluster))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Cluster = clusterObj

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the NFS volume from the cluster.
func (r *skeClusterNfsVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ske.ClusterNfsVolumeResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.ClusterId.ValueString()

	// Detach NFS volume from the cluster (set to null)
	_, httpStatus, err := r.client.UnsetClusterNfsVolume(ctx, clusterId)
	if err != nil {
		if httpStatus == http.StatusNotFound {
			// Cluster already deleted; nothing to detach.
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Detaching NFS Volume",
			"Could not detach NFS volume from cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Wait for cluster to return to RUNNING status
	err = waitForClusterStatus(ctx, r.client, clusterId, []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Detaching NFS Volume",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return
	}
}

func (r *skeClusterNfsVolumeResource) makeClusterNfsVolumeModel(ctx context.Context, cluster *scpske16.ClusterV1Dot6) (types.Object, diag.Diagnostics) {
	nfsVolumeId := ""
	if cluster.NfsVolumeId.Get() != nil {
		nfsVolumeId = *cluster.NfsVolumeId.Get()
	}

	clusterModel := struct {
		Id          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Status      types.String `tfsdk:"status"`
		NfsVolumeId types.String `tfsdk:"nfs_volume_id"`
	}{
		Id:          types.StringValue(cluster.Id),
		Name:        types.StringValue(cluster.Name),
		Status:      types.StringValue(cluster.Status),
		NfsVolumeId: types.StringValue(nfsVolumeId),
	}

	attributeTypes := map[string]attr.Type{
		"id":            types.StringType,
		"name":          types.StringType,
		"status":        types.StringType,
		"nfs_volume_id": types.StringType,
	}

	obj, diags := types.ObjectValueFrom(ctx, attributeTypes, clusterModel)
	return obj, diags
}

func (r *skeClusterNfsVolumeResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}
