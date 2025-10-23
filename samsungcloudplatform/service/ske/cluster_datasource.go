package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/ske/1.1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &skeClusterDataSource{}
	_ datasource.DataSourceWithConfigure = &skeClusterDataSource{}
)

// skeClusterDataSource is the data source implementation.
type skeClusterDataSource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

func NewSkeClusterDataSource() datasource.DataSource {
	return &skeClusterDataSource{}
}

//// datasource.DataSource Interface Methods

func (d *skeClusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_cluster"
}

func (d *skeClusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "cluster.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("Cluster"): schema.SingleNestedAttribute{
				Description: "cluster",
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
							Attributes: d.makePrivateEndpointAccessControlResourceSchema(),
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
						Attributes:  d.makeExternalResourceSchema(),
					},
					common.ToSnakeCase("Subnet"): schema.SingleNestedAttribute{
						Description: "Subnet",
						Computed:    true,
						Attributes:  d.makeExternalResourceSchema(),
					},

					common.ToSnakeCase("Volume"): schema.SingleNestedAttribute{
						Description: "Volume",
						Computed:    true,
						Attributes:  d.makeExternalResourceSchema(),
					},
					common.ToSnakeCase("SecurityGroupList"): schema.ListNestedAttribute{
						Description: "SecurityGroupList",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: d.makeExternalResourceSchema(),
						},
					},
					common.ToSnakeCase("ManagedSecurityGroup"): schema.SingleNestedAttribute{
						Description: "ManagedSecurityGroup",
						Computed:    true,
						Attributes:  d.makeExternalResourceSchema(),
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

func (d *skeClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ske.ClusterDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Cluster
	data, _, err := d.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			"Could not read Cluster ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	cluster := data.Cluster

	var securityGroups []ske.ExternalResource
	for _, securityGroup := range cluster.SecurityGroupList {
		securityGroups = append(securityGroups, d.makeExternalResourceModel((*scpske.ExternalResource)(&securityGroup)))
	}

	var privateEndpointAccessControlResources []ske.PrivateEndpointAccessControlResource
	for _, privateEndpointAccessControlResource := range cluster.PrivateEndpointAccessControlResources {

		privateEndpointAccessControlResources = append(privateEndpointAccessControlResources, d.makePrivateEndpointAccessControlResourceModel((*scpske.PrivateEndpointAccessControlResource)(&privateEndpointAccessControlResource)))
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
		Vpc:                                   d.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Vpc.Get())),
		Subnet:                                d.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Subnet.Get())),
		Volume:                                d.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Volume.Get())),
		SecurityGroupList:                     securityGroups,
		ManagedSecurityGroup:                  d.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Vpc.Get())),
		CreatedAt:                             types.StringValue(cluster.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                             types.StringValue(cluster.CreatedBy),
		ModifiedAt:                            types.StringValue(cluster.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:                            types.StringValue(cluster.ModifiedBy),
		Status:                                types.StringValue(cluster.Status),
		ServiceWatchLoggingEnabled:            types.BoolPointerValue(cluster.ServiceWatchLoggingEnabled),
	}
	clusterObjectValue, _ := types.ObjectValueFrom(ctx, clusterModel.AttributeTypes(), clusterModel)
	println("clusterObjectValue:", clusterObjectValue.String())
	state.Cluster = clusterObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *skeClusterDataSource) makeExternalResourceModel(externalResource *scpske.ExternalResource) ske.ExternalResource {
	return ske.ExternalResource{
		Id:   types.StringValue(externalResource.GetId()),
		Name: types.StringValue(externalResource.GetName()),
	}
}

func (d *skeClusterDataSource) makePrivateEndpointAccessControlResourceModel(privateEndpointAccessControlResource *scpske.PrivateEndpointAccessControlResource) ske.PrivateEndpointAccessControlResource {
	return ske.PrivateEndpointAccessControlResource{
		Id:   types.StringValue(privateEndpointAccessControlResource.GetId()),
		Name: types.StringValue(privateEndpointAccessControlResource.GetName()),
		Type: types.StringValue(privateEndpointAccessControlResource.GetType()),
	}
}

//// datasource.DataSourceWithConfigure Interface Methods

func (d *skeClusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Ske
	d.clients = inst.Client
}

// private
func (d *skeClusterDataSource) makeExternalResourceSchema() map[string]schema.Attribute {
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

func (d *skeClusterDataSource) makePrivateEndpointAccessControlResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		common.ToSnakeCase("Id"): schema.StringAttribute{
			Description: "Private Endpoint Access ControlResource Id",
			Computed:    true,
		},
		common.ToSnakeCase("Name"): schema.StringAttribute{
			Description: "Private Endpoint Access ControlResource Name",
			Computed:    true,
		},
		common.ToSnakeCase("Type"): schema.StringAttribute{
			Description: "Private Endpoint Access ControlResource Type",
			Computed:    true,
		},
	}
}
