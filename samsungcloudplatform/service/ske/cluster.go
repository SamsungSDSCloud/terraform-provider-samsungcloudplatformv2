package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/region"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/ske/1.1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"strings"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &skeClusterResource{}
	_ resource.ResourceWithConfigure = &skeClusterResource{}
)

// NewSkeClusterResource is a helper function to simplify the provider implementation.
func NewSkeClusterResource() resource.Resource {
	return &skeClusterResource{}
}

// skeClusterResource is the data source implementation.
type skeClusterResource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *skeClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_cluster"
}

// Schema defines the schema for the data source.
func (r *skeClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "cluster",
		Attributes: map[string]schema.Attribute{
			"region": region.DataSourceSchema(),
			"tags":   tag.DataSourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the cluster",
				Computed:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Required:    true,
			},
			common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
				Description: "KubernetesVersion",
				Required:    true,
			},
			common.ToSnakeCase("VpcID"): schema.StringAttribute{
				Description: "VpcID",
				Required:    true,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "SubnetId",
				Required:    true,
			},
			common.ToSnakeCase("VolumeId"): schema.StringAttribute{
				Description: "VolumeId",
				Required:    true,
			},
			common.ToSnakeCase("CloudLoggingEnabled"): schema.BoolAttribute{
				Description: "CloudLoggingEnabled",
				Required:    true,
			},
			common.ToSnakeCase("SecurityGroupIdList"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "SecurityGroupIdList",
				Required:    true,
			},
			common.ToSnakeCase("VolumeId"): schema.StringAttribute{
				Description: "VolumeId",
				Required:    true,
			},
			common.ToSnakeCase("PrivateEndpointAccessControlResources"): schema.ListNestedAttribute{
				Description: "PrivateEndpointAccessControlResources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.makePrivateEndpointAccessControlResourceListSchema(),
				},
				Optional: true,
			},
			common.ToSnakeCase("PublicEndpointAccessControlIp"): schema.StringAttribute{
				Description: "PublicEndpointAccessControlIp",
				Optional:    true,
			},
			// v1.1
			common.ToSnakeCase("ServiceWatchLoggingEnabled"): schema.BoolAttribute{
				Description: "ServiceWatchLoggingEnabled",
				Required:    true,
				Optional:    false,
			},
			common.ToSnakeCase("Cluster"): schema.SingleNestedAttribute{
				Description: "Cluster",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("CloudLoggingEnabled"): schema.BoolAttribute{
						Description: "CloudLoggingEnabled",
						Computed:    true,
					},
					common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
						Description: "KubernetesVersion",
						Computed:    true,
					},
					common.ToSnakeCase("ClusterNamespace"): schema.StringAttribute{
						Description: "ClusterNamespace",
						Computed:    true,
					},
					common.ToSnakeCase("MaxNodeCount"): schema.Int32Attribute{
						Description: "MaxNodeCount",
						Computed:    true,
					},
					common.ToSnakeCase("NodeCount"): schema.Int32Attribute{
						Description: "NodeCount",
						Computed:    true,
					},
					common.ToSnakeCase("PrivateEndpointUrl"): schema.StringAttribute{
						Description: "PrivateEndpointUrl",
						Computed:    true,
					},
					common.ToSnakeCase("PrivateKubeconfigDownloadYn"): schema.StringAttribute{
						Description: "PrivateKubeconfigDownloadYn",
						Computed:    true,
					},
					common.ToSnakeCase("PrivateEndpointAccessControlResources"): schema.ListNestedAttribute{
						Description: "PrivateEndpointAccessControlResources",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: r.makePrivateEndpointAccessControlResourceSchema(),
						},
					},
					common.ToSnakeCase("PublicEndpointUrl"): schema.StringAttribute{
						Description: "PublicEndpointUrl",
						Computed:    true,
					},
					common.ToSnakeCase("PublicKubeconfigDownloadYn"): schema.StringAttribute{
						Description: "PublicKubeconfigDownloadYn",
						Computed:    true,
					},
					common.ToSnakeCase("PublicEndpointAccessControlIp"): schema.StringAttribute{
						Description: "PublicEndpointAccessControlIp",
						Computed:    true,
					},
					common.ToSnakeCase("Vpc"): schema.SingleNestedAttribute{
						Description: "Vpc",
						Computed:    true,
						Attributes:  r.makeExternalResourceSchema(),
					},
					common.ToSnakeCase("Subnet"): schema.SingleNestedAttribute{
						Description: "Subnet",
						Computed:    true,
						Attributes:  r.makeExternalResourceSchema(),
					},
					common.ToSnakeCase("Volume"): schema.SingleNestedAttribute{
						Description: "Volume",
						Computed:    true,
						Attributes:  r.makeExternalResourceSchema(),
					},
					common.ToSnakeCase("SecurityGroupList"): schema.ListNestedAttribute{
						Description: "SecurityGroupList",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: r.makeExternalResourceSchema(),
						},
					},
					common.ToSnakeCase("ManagedSecurityGroup"): schema.SingleNestedAttribute{
						Description: "ManagedSecurityGroup",
						Computed:    true,
						Attributes:  r.makeExternalResourceSchema(),
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Computed:    true,
					},
					// v1.1
					common.ToSnakeCase("ServiceWatchLoggingEnabled"): schema.BoolAttribute{
						Description: "ServiceWatchLoggingEnabled",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *skeClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.Ske
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *skeClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ske.ClusterResource
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new cluster
	data, err := r.client.CreateCluster(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Creating Cluster",
			"Could not create cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.ResourceId)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForClusterStatus(ctx, r.client, data.ResourceId, []string{"CREATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
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
func (r *skeClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ske.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from cluster
	data, _, err := r.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			"Could not read cluster ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	cluster := data.Cluster

	var securityGroups []ske.ExternalResource
	for _, securityGroup := range cluster.SecurityGroupList {
		securityGroups = append(securityGroups, r.makeExternalResourceModel((*scpske.ExternalResource)(&securityGroup)))
	}

	var privateEndpointAccessControlResources []ske.PrivateEndpointAccessControlResource
	for _, privateEndpointAccessControlResource := range cluster.PrivateEndpointAccessControlResources {
		privateEndpointAccessControlResources = append(privateEndpointAccessControlResources, r.makePrivateEndpointAccessControlResourceModel((*scpske.PrivateEndpointAccessControlResource)(&privateEndpointAccessControlResource)))
	}

	clusterModel := ske.Cluster{
		Id:                                    types.StringValue(cluster.Id),
		Name:                                  types.StringValue(cluster.Name),
		AccountId:                             types.StringValue(cluster.AccountId),
		CloudLoggingEnabled:                   types.BoolPointerValue(cluster.CloudLoggingEnabled),
		KubernetesVersion:                     types.StringValue(cluster.KubernetesVersion),
		ClusterNamespace:                      types.StringValue(cluster.ClusterNamespace),
		MaxNodeCount:                          types.Int32PointerValue(cluster.MaxNodeCount.Get()),
		NodeCount:                             types.Int32PointerValue(cluster.NodeCount.Get()),
		PrivateEndpointUrl:                    types.StringValue(cluster.PrivateEndpointUrl),
		PrivateKubeconfigDownloadYn:           types.StringValue(cluster.PrivateKubeconfigDownloadYn),
		PrivateEndpointAccessControlResources: privateEndpointAccessControlResources,
		PublicEndpointUrl:                     types.StringValue(cluster.GetPublicEndpointUrl()),
		PublicKubeconfigDownloadYn:            types.StringValue(cluster.PublicKubeconfigDownloadYn),
		PublicEndpointAccessControlIp:         types.StringValue(cluster.GetPublicEndpointAccessControlIp()),
		Vpc:                                   r.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Vpc.Get())),
		Subnet:                                r.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Subnet.Get())),
		Volume:                                r.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Volume.Get())),
		SecurityGroupList:                     securityGroups,
		ManagedSecurityGroup:                  r.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Vpc.Get())),
		CreatedAt:                             types.StringValue(cluster.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                             types.StringValue(cluster.CreatedBy),
		ModifiedAt:                            types.StringValue(cluster.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:                            types.StringValue(cluster.ModifiedBy),
		Status:                                types.StringValue(cluster.Status),
		ServiceWatchLoggingEnabled:            types.BoolValue(cluster.GetServiceWatchLoggingEnabled()), // v1.1
	}
	clusterObjectValue, _ := types.ObjectValueFrom(ctx, clusterModel.AttributeTypes(), clusterModel)
	state.Cluster = clusterObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *skeClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state ske.ClusterResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	req.State.Get(ctx, &state)

	err := r.syncCloudLoggingEnabled(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncSecurityGroupList(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncKubernetesVersion(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncPrivateEndpointAccessControlResources(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncPublicEndpointAccessControlIp(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncServiceWatchLoggingEnabled(ctx, state, &plan, resp)
	if err != nil {
		return
	}

	diags = resp.State.Set(ctx, &plan)
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

func (r *skeClusterResource) syncCloudLoggingEnabled(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if state.CloudLoggingEnabled.Equal(plan.CloudLoggingEnabled) {
		return nil
	}
	data, err := r.client.UpdateClusterLogging(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating CloudLoggingEnabled",
			"Could not update cloud logging enabled, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncSecurityGroupList(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	securityGroupIdListPlan, _ := types.ListValueFrom(ctx, types.StringType, plan.SecurityGroupIdList)
	securityGroupIdListState, _ := types.ListValueFrom(ctx, types.StringType, state.SecurityGroupIdList)
	if securityGroupIdListPlan.Equal(securityGroupIdListState) {
		return nil
	}
	data, err := r.client.UpdateClusterSecurityGroups(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating SecurityGroupList",
			"Could not update security group list, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	plan.Id = types.StringValue(data.Cluster.Id)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncKubernetesVersion(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.KubernetesVersion.Equal(state.KubernetesVersion) {
		return nil
	}
	data, err := r.client.UpgradeCluster(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating KubernetesVersion",
			"Could not update kubernetes version, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncPrivateEndpointAccessControlResources(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if reflect.DeepEqual(plan.PrivateEndpointAccessControlResources, state.PrivateEndpointAccessControlResources) {
		return nil
	}
	data, err := r.client.UpdatePrivateEndpointAccessControlResources(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating PrivateEndpointAccessControlResources",
			"Could not update cluster private endpoint access control resources, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncPublicEndpointAccessControlIp(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.PublicEndpointAccessControlIp.Equal(state.PublicEndpointAccessControlIp) {
		return nil
	}
	data, err := r.client.UpdatePublicEndpointAccessControlIps(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating PublicEndpointAccessControlIp",
			"Could not update public endpoint access control ip, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncServiceWatchLoggingEnabled(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.ServiceWatchLoggingEnabled.Equal(state.ServiceWatchLoggingEnabled) {
		return nil
	}
	data, err := r.client.UpdateServiceWatchLoggingEnabled(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating ServiceWatchLoggingEnabled",
			"Could not update service watch logging enabled, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *skeClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ske.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing cluster
	data, err := r.client.DeleteCluster(ctx, state.Id.ValueString())
	if err != nil && !strings.Contains(err.Error(), "404") {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Cluster",
			"Could not delete cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForClusterStatus(ctx, r.client, data.ResourceId, []string{}, []string{"DELETED"}, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Cluster",
			"Error waiting for cluster to become deleted: "+err.Error(),
		)
		return
	}
}

// private
func (r *skeClusterResource) makeExternalResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		common.ToSnakeCase("Id"): schema.StringAttribute{
			Description: "External Resource Id",
			Computed:    true,
		},
		common.ToSnakeCase("Name"): schema.StringAttribute{
			Description: "External Resource Id",
			Computed:    true,
		},
	}
}

func (r *skeClusterResource) makeExternalResourceModel(externalResource *scpske.ExternalResource) ske.ExternalResource {
	return ske.ExternalResource{
		Id:   types.StringValue(externalResource.GetId()),
		Name: types.StringValue(externalResource.GetName()),
	}
}

// private
func (r *skeClusterResource) makePrivateEndpointAccessControlResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		common.ToSnakeCase("Id"): schema.StringAttribute{
			Description: "External Resource Id",
			Computed:    true,
		},
		common.ToSnakeCase("Name"): schema.StringAttribute{
			Description: "External Resource Name",
			Computed:    true,
		},
		common.ToSnakeCase("Type"): schema.StringAttribute{
			Description: "External Resource Type",
			Computed:    true,
		},
	}
}

// private
func (r *skeClusterResource) makePrivateEndpointAccessControlResourceListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		common.ToSnakeCase("Id"): schema.StringAttribute{
			Description: "External Resource Id",
			Required:    true,
		},
		common.ToSnakeCase("Name"): schema.StringAttribute{
			Description: "External Resource Name",
			Required:    true,
		},
		common.ToSnakeCase("Type"): schema.StringAttribute{
			Description: "External Resource Type",
			Required:    true,
		},
	}
}

func (r *skeClusterResource) makePrivateEndpointAccessControlResourceModel(privateEndpointAccessControlResource *scpske.PrivateEndpointAccessControlResource) ske.PrivateEndpointAccessControlResource {
	return ske.PrivateEndpointAccessControlResource{
		Id:   types.StringValue(privateEndpointAccessControlResource.GetId()),
		Name: types.StringValue(privateEndpointAccessControlResource.GetName()),
		Type: types.StringValue(privateEndpointAccessControlResource.GetName()),
	}
}

//

func waitForClusterStatus(ctx context.Context, skeClient *ske.Client, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpStatus, err := skeClient.GetCluster(ctx, id)
		if httpStatus == 200 {
			return info, info.Cluster.Status, nil
		} else if httpStatus == 404 {
			if errorOnNotFound {
				return nil, "", fmt.Errorf("cluster with id=%s not found", id)
			}

			return info, "DELETED", nil
		} else if err != nil {
			return nil, "", err
		}

		return info, info.Cluster.Status, nil
	})
}
