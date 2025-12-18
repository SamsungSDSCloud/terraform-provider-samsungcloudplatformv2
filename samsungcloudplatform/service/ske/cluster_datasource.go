package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/ske"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/ske/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
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
		Description: "show cluster",
		Attributes: map[string]schema.Attribute{
			"cluster": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
						MarkdownDescription: "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
					},
					"cloud_logging_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Cloud Logging Enabled\n  - example: true",
						MarkdownDescription: "Cloud Logging Enabled\n  - example: true",
					},
					"cluster_namespace": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Namespace\n  - example: sample-cluster-12345",
						MarkdownDescription: "Cluster Namespace\n  - example: sample-cluster-12345",
					},
					"created_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
					},
					"created_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						MarkdownDescription: "ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
					},
					"kubernetes_version": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Version\n  - example: v1.29.8",
						MarkdownDescription: "Cluster Version\n  - example: v1.29.8",
					},
					"managed_security_group": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Managed Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "Managed Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "Managed Security Group Name\n  - example: sample-name",
								MarkdownDescription: "Managed Security Group Name\n  - example: sample-name",
							},
						},
						Computed:            true,
						Description:         "Managed Security Group",
						MarkdownDescription: "Managed Security Group",
					},
					"max_node_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Cluster Max Node Count\n  - example: 5",
						MarkdownDescription: "Cluster Max Node Count\n  - example: 5",
					},
					"modified_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
					},
					"modified_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Name\n  - example: sample-cluster",
						MarkdownDescription: "Cluster Name\n  - example: sample-cluster",
					},
					"node_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Cluster Node Count\n  - example: 5",
						MarkdownDescription: "Cluster Node Count\n  - example: 5",
					},
					"private_endpoint_access_control_resources": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Private Endpoint Access Control Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
									MarkdownDescription: "Private Endpoint Access Control Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "Private Endpoint Access Control Resource Name\n  - example: sample-name",
									MarkdownDescription: "Private Endpoint Access Control Resource Name\n  - example: sample-name",
								},
								"type": schema.StringAttribute{
									Computed:            true,
									Description:         "Private Endpoint Access Control Resource Type\n  - example: vm",
									MarkdownDescription: "Private Endpoint Access Control Resource Type\n  - example: vm",
								},
							},
						},
						Computed:            true,
						Description:         "Private Endpoint Access Control Resources",
						MarkdownDescription: "Private Endpoint Access Control Resources",
					},
					"private_endpoint_url": schema.StringAttribute{
						Computed:            true,
						Description:         "Private Kubeconfig Download Yn\n  - example: N",
						MarkdownDescription: "Private Kubeconfig Download Yn\n  - example: N",
					},
					"private_kubeconfig_download_yn": schema.StringAttribute{
						Computed:            true,
						Description:         "Private Endpoint URL\n  - example: https://sample-cluster.ske.private.kr-west1.samsungsdscloud.com:6443",
						MarkdownDescription: "Private Endpoint URL\n  - example: https://sample-cluster.ske.private.kr-west1.samsungsdscloud.com:6443",
					},
					"public_endpoint_access_control_ip": schema.StringAttribute{
						Computed:            true,
						Description:         "Public Endpoint Access Control IP\n  - example: 192.168.0.0",
						MarkdownDescription: "Public Endpoint Access Control IP\n  - example: 192.168.0.0",
					},
					"public_endpoint_url": schema.StringAttribute{
						Computed:            true,
						Description:         "Public Endpoint URL\n  - example: https://sample-cluster.ske.kr-west1.samsungsdscloud.com:6443",
						MarkdownDescription: "Public Endpoint URL\n  - example: https://sample-cluster.ske.kr-west1.samsungsdscloud.com:6443",
					},
					"public_kubeconfig_download_yn": schema.StringAttribute{
						Computed:            true,
						Description:         "Public Kubeconfig Download Yn\n  - example: N",
						MarkdownDescription: "Public Kubeconfig Download Yn\n  - example: N",
					},
					"security_group_list": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
									MarkdownDescription: "Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "Security Group Name\n  - example: sample-name",
									MarkdownDescription: "Security Group Name\n  - example: sample-name",
								},
							},
						},
						Computed:            true,
						Description:         "Connected Security Group List",
						MarkdownDescription: "Connected Security Group List",
					},
					"service_watch_logging_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Service Watch Enabled\n  - example: true",
						MarkdownDescription: "Service Watch Enabled\n  - example: true",
					},
					"status": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Status\n  - example: RUNNING",
						MarkdownDescription: "Cluster Status\n  - example: RUNNING",
					},
					"subnet": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Subnet ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "Subnet ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "Subnet Name\n  - example: sample-name",
								MarkdownDescription: "Subnet Name\n  - example: sample-name",
							},
						},
						Computed:            true,
						Description:         "Subnet of Cluster",
						MarkdownDescription: "Subnet of Cluster",
					},
					"volume": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Volume ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "Volume ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "Volume Name\n  - example: sample-name",
								MarkdownDescription: "Volume Name\n  - example: sample-name",
							},
						},
						Computed:            true,
						Description:         "Connected File Storage",
						MarkdownDescription: "Connected File Storage",
					},
					"vpc": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "VPC ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "VPC ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "VPC Name\n  - example: sample-name",
								MarkdownDescription: "VPC Name\n  - example: sample-name",
							},
						},
						Computed:            true,
						Description:         "VPC of Cluster",
						MarkdownDescription: "VPC of Cluster",
					},
				},
				Computed: true,
			},
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				MarkdownDescription: "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-f]{32}$"), ""),
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
		CloudLoggingEnabled:                   types.BoolValue(cluster.CloudLoggingEnabled),
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
		ServiceWatchLoggingEnabled:            types.BoolValue(cluster.ServiceWatchLoggingEnabled),
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
