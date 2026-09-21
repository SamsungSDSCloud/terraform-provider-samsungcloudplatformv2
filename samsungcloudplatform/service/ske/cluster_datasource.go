package ske

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/ske/1.6"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
					"additional_subnet_id_list": schema.ListAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						Description:         "Additional Subnet ID List",
						MarkdownDescription: "List of additional subnet IDs associated with the cluster.",
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
						Computed: true,
						Description: "Cluster Version\n" +
							"  - pattern: ^v[0-9]{1}\\.[0-9]{1,2}\\.[0-9]{1,2}$\n" +
							"  - pattern: v1.31.X|v1.32.X|v1.33.X|v1.34.X\n" +
							"  - example: v1.34.3\n" +
							"  - Use the samsungcloudplatformv2_ske_kubernetes_versions data source to query the SKE service for all Kubernetes versions supported. (ex v1.31.X|v1.32.X|v1.33.X|v1.34.X)",
						MarkdownDescription: "Cluster Version\n" +
							"  - pattern: ^v[0-9]{1}\\.[0-9]{1,2}\\.[0-9]{1,2}$\n" +
							"  - pattern: v1.31.X|v1.32.X|v1.33.X|v1.34.X\n" +
							"  - example: v1.34.3\n" +
							"  - Use the samsungcloudplatformv2_ske_kubernetes_versions data source to query the SKE service for all Kubernetes versions supported. (ex v1.31.X|v1.32.X|v1.33.X|v1.34.X)",
					},
					"linked_resources": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Linked Resource ID\n  - example: res-12345678",
									MarkdownDescription: "Linked Resource ID\n  - example: res-12345678",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "Linked Resource Name\n  - example: my-resource",
									MarkdownDescription: "Linked Resource Name\n  - example: my-resource",
								},
								"type": schema.StringAttribute{
									Computed:            true,
									Description:         "Linked Resource Type (fs/obs)\n  - pattern: fs|obs\n  - example: fs",
									MarkdownDescription: "Linked Resource Type (fs/obs)\n  - pattern: fs|obs\n  - example: fs",
								},
							},
						},
						Computed:            true,
						Description:         "List of linked resources associated with the cluster\n  - example: {id='res-12345678', name='my-resource', type='fs'}",
						MarkdownDescription: "List of linked resources associated with the cluster\n  - example: {id='res-12345678', name='my-resource', type='fs'}",
					},
					"managed_security_group": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Managed Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "Managed Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
						},
						Computed: true,
						Description: "Managed Security Group\n" +
							exampleNestedResourceID,
						MarkdownDescription: "Managed Security Group\n" +
							exampleNestedResourceID,
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
									Description:         "Private Endpoint Access Control Resource Type\n  - pattern: vm|bm|gpuvm|mngc|devops\n  - example: vm",
									MarkdownDescription: "Private Endpoint Access Control Resource Type\n  - pattern: vm|bm|gpuvm|mngc|devops\n  - example: vm",
								},
							},
						},
						Computed: true,
						Description: "Private Endpoint Access Control Resources\n" +
							"  - example: {id='2a9be312-5d4b-4bc8-b2ae-35100fa9241f', name='sample-name', type='vm'}",
						MarkdownDescription: "Private Endpoint Access Control Resources\n" +
							"  - example: {id='2a9be312-5d4b-4bc8-b2ae-35100fa9241f', name='sample-name', type='vm'}",
					},
					"private_endpoint_url": schema.StringAttribute{
						Computed:            true,
						Description:         "Private Endpoint URL\n  - example: https://sample-cluster.ske.private.kr-west1.samsungsdscloud.com:6443",
						MarkdownDescription: "Private Endpoint URL\n  - example: https://sample-cluster.ske.private.kr-west1.samsungsdscloud.com:6443",
					},
					"private_kubeconfig_download_yn": schema.StringAttribute{
						Computed:            true,
						Description:         "Private Kubeconfig Download Yn\n  - pattern: Y|N\n  - example: N",
						MarkdownDescription: "Private Kubeconfig Download Yn\n  - pattern: Y|N\n  - example: N",
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
						Description:         "Public Kubeconfig Download Yn\n  - pattern: Y|N\n  - example: N",
						MarkdownDescription: "Public Kubeconfig Download Yn\n  - pattern: Y|N\n  - example: N",
					},
					"security_group_list": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
									MarkdownDescription: "Security Group ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								},
							},
						},
						Computed: true,
						Description: "Connected Security Group List\n" +
							exampleNestedResourceID,
						MarkdownDescription: "Connected Security Group List\n" +
							exampleNestedResourceID,
					},
					"service_watch_logging_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Service Watch Enabled\n  - example: true",
						MarkdownDescription: "Service Watch Enabled\n  - example: true",
					},
					"status": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Status\n  - pattern: RUNNING|CREATING|UPDATING|DELETING\n  - example: RUNNING",
						MarkdownDescription: "Cluster Status\n  - pattern: RUNNING|CREATING|UPDATING|DELETING\n  - example: RUNNING",
					},
					"subnet": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Subnet ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "Subnet ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
						},
						Computed: true,
						Description: "Subnet of Cluster\n" +
							exampleNestedResourceID,
						MarkdownDescription: "Subnet of Cluster\n" +
							exampleNestedResourceID,
					},
					"volume": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "NFS Volume ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "NFS Volume ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
						},
						Computed: true,
						Description: "Connected File Storage\n" +
							exampleNestedResourceID,
						MarkdownDescription: "Connected File Storage\n" +
							exampleNestedResourceID,
					},
					"vpc": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "VPC ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "VPC ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
						},
						Computed: true,
						Description: "VPC of Cluster\n" +
							exampleNestedResourceID,
						MarkdownDescription: "VPC of Cluster\n" +
							exampleNestedResourceID,
					},
				},
				Computed:            true,
				Description:         "Cluster\n - example: https://registry.terraform.io/providers/SamsungSDSCloud/samsungcloudplatformv2/latest/docs/resources/ske_cluster#nested-schema-for-cluster",
				MarkdownDescription: "Cluster\n - example: https://registry.terraform.io/providers/SamsungSDSCloud/samsungcloudplatformv2/latest/docs/resources/ske_cluster#nested-schema-for-cluster",
			},
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				MarkdownDescription: "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-f]{32}$"), ""),
				},
			},
			"deletion_protection_enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Cluster Deletion Protection Enabled\n  - example: true",
				MarkdownDescription: "Cluster Deletion Protection Enabled\n  - example: true",
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

	// Get refreshed value from Cluster (v1.6)
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

	var securityGroups []ske.ExternalResourceId
	for _, sgId := range cluster.SecurityGroupIdList {
		securityGroups = append(securityGroups, ske.ExternalResourceId{
			Id: types.StringValue(sgId),
		})
	}

	var privateEndpointAccessControlResources []ske.PrivateEndpointAccessControlResource
	for _, privateEndpointAccessControlResource := range cluster.PrivateEndpointAccessControlResources {

		privateEndpointAccessControlResources = append(privateEndpointAccessControlResources, d.makePrivateEndpointAccessControlResourceModel((*scpske.PrivateEndpointAccessControlResource)(&privateEndpointAccessControlResource)))
	}

	clusterModel := ske.Cluster{
		Id:                                    types.StringValue(cluster.Id),
		Name:                                  types.StringValue(cluster.Name),
		AccountId:                             types.StringValue(cluster.AccountId),
		AdditionalSubnetIdList:                func() types.List { v, _ := types.ListValueFrom(ctx, types.StringType, func() []string { ids := make([]string, len(cluster.AdditionalSubnetIdList)); copy(ids, cluster.AdditionalSubnetIdList); return ids }()); return v }(),
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
		Vpc:                                   d.nullableStringToExternalResourceId(cluster.VpcId),
		Subnet:                                d.nullableStringToExternalResourceId(cluster.DefaultSubnetId),
		Volume:                                d.nullableStringToExternalResourceId(cluster.NfsVolumeId),
		SecurityGroupList:                     securityGroups,
		ManagedSecurityGroup:                  d.nullableStringToExternalResourceId(cluster.ManagedSecurityGroupId),
		CreatedAt:                             types.StringValue(cluster.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                             types.StringValue(cluster.CreatedBy),
		ModifiedAt:                            types.StringValue(cluster.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:                            types.StringValue(cluster.ModifiedBy),
		Status:                                types.StringValue(cluster.Status),
		ServiceWatchLoggingEnabled:            types.BoolValue(cluster.ServiceWatchLoggingEnabled),
		LinkedResources:                       convertLinkedResourcesFromSDK(cluster.LinkedResources),
	}
	clusterObjectValue, _ := types.ObjectValueFrom(ctx, clusterModel.AttributeTypes(), clusterModel)
	state.Cluster = clusterObjectValue
	state.DeletionProtectionEnabled = types.BoolValue(cluster.DeletionProtectionEnabled)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *skeClusterDataSource) nullableStringToExternalResourceId(ns scpske.NullableString) ske.ExternalResourceId {
	if ns.Get() == nil {
		return ske.ExternalResourceId{
			Id: types.StringNull(),
		}
	}
	return ske.ExternalResourceId{
		Id: types.StringValue(*ns.Get()),
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
