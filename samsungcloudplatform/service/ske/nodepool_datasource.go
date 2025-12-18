package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/ske/converter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/ske/1.1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &skeNodepoolDataSource{}
	_ datasource.DataSourceWithConfigure = &skeNodepoolDataSource{}
)

// skeNodepoolDataSource is the data source implementation.
type skeNodepoolDataSource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

func NewSkeNodepoolDataSource() datasource.DataSource {
	return &skeNodepoolDataSource{}
}

//// datasource.DataSource Interface Methods

func (d *skeNodepoolDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_nodepool"
}

func (d *skeNodepoolDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "show nodepool",
		Attributes: map[string]schema.Attribute{
			"nodepool": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
						MarkdownDescription: "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
					},
					"advanced_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"allowed_unsafe_sysctls": schema.StringAttribute{
								Computed:            true,
								Description:         "Node Pool Allowed unsafe sysctls\n  - example: kernel.msg*,net.ipv4.route.min_pmtu",
								MarkdownDescription: "Node Pool Allowed unsafe sysctls\n  - example: kernel.msg*,net.ipv4.route.min_pmtu",
							},
							"container_log_max_files": schema.Int64Attribute{
								Computed:            true,
								Description:         "Node Pool container log max files\n  - maximum: 10\n  - minimum: 2\n  - example: 5",
								MarkdownDescription: "Node Pool container log max files\n  - maximum: 10\n  - minimum: 2\n  - example: 5",
							},
							"container_log_max_size": schema.Int64Attribute{
								Computed:            true,
								Description:         "Node Pool container log max size\n  - maximum: 100\n  - minimum: 10\n  - example: 10",
								MarkdownDescription: "Node Pool container log max size\n  - maximum: 100\n  - minimum: 10\n  - example: 10",
							},
							"image_gc_high_threshold": schema.Int64Attribute{
								Computed:            true,
								Description:         "Node Pool image GC high threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 85",
								MarkdownDescription: "Node Pool image GC high threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 85",
							},
							"image_gc_low_threshold": schema.Int64Attribute{
								Computed:            true,
								Description:         "Node Pool image GC low threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 80",
								MarkdownDescription: "Node Pool image GC low threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 80",
							},
							"max_pods": schema.Int64Attribute{
								Computed:            true,
								Description:         "Node Pool max pod number\n  - maximum: 250\n  - minimum: 10\n  - example: 110",
								MarkdownDescription: "Node Pool max pod number\n  - maximum: 250\n  - minimum: 10\n  - example: 110",
							},
							"pod_max_pids": schema.Int64Attribute{
								Computed:            true,
								Description:         "Node Pool Pod Max pids constraint\n  - maximum: 4.194304e+06\n  - minimum: 1024\n  - example: 4096",
								MarkdownDescription: "Node Pool Pod Max pids constraint\n  - maximum: 4.194304e+06\n  - minimum: 1024\n  - example: 4096",
							},
						},
						Computed:            true,
						Description:         "Node Pool Advanced Settings",
						MarkdownDescription: "Node Pool Advanced Settings",
					},
					"auto_recovery_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Is Auto Recovery\n  - example: true",
						MarkdownDescription: "Is Auto Recovery\n  - example: true",
					},
					"auto_scale_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Is Auto Scale\n  - example: true",
						MarkdownDescription: "Is Auto Scale\n  - example: true",
					},
					"cluster": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
								MarkdownDescription: "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
							},
						},
						Computed:            true,
						Description:         "Cluster",
						MarkdownDescription: "Cluster",
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
					"current_node_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Current Node Count\n  - example: 1",
						MarkdownDescription: "Current Node Count\n  - example: 1",
					},
					"desired_node_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Desired Node Count\n  - example: 2",
						MarkdownDescription: "Desired Node Count\n  - example: 2",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Nodepool ID\n  - example: bdfda539-bd2e-4a5c-9021-ec6d52d1ca79",
						MarkdownDescription: "Nodepool ID\n  - example: bdfda539-bd2e-4a5c-9021-ec6d52d1ca79",
					},
					"image": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"custom_image_name": schema.StringAttribute{
								Computed:            true,
								Description:         "Custom Image Name\n  - example: custom-image",
								MarkdownDescription: "Custom Image Name\n  - example: custom-image",
							},
							"os": schema.StringAttribute{
								Computed:            true,
								Description:         "Image OS\n  - example: ubuntu",
								MarkdownDescription: "Image OS\n  - example: ubuntu",
							},
							"os_version": schema.StringAttribute{
								Computed:            true,
								Description:         "Image OS Version\n  - example: 22.04",
								MarkdownDescription: "Image OS Version\n  - example: 22.04",
							},
						},
						Computed:            true,
						Description:         "Image",
						MarkdownDescription: "Image",
					},
					"keypair": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "Keypair Name\n  - example: test_keypair",
								MarkdownDescription: "Keypair Name\n  - example: test_keypair",
							},
						},
						Computed:            true,
						Description:         "Keypair Name",
						MarkdownDescription: "Keypair Name",
					},
					"kubernetes_version": schema.StringAttribute{
						Computed:            true,
						Description:         "Kubernetes Version\n  - example: v1.29.8",
						MarkdownDescription: "Kubernetes Version\n  - example: v1.29.8",
					},
					"labels": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Computed:            true,
									Description:         "Node Pool Label Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
									MarkdownDescription: "Node Pool Label Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
								},
								"value": schema.StringAttribute{
									Computed:            true,
									Description:         "Node Pool Label Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
									MarkdownDescription: "Node Pool Label Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
								},
							},
						},
						Computed:            true,
						Description:         "Node Pool Labels",
						MarkdownDescription: "Node Pool Labels",
					},
					"max_node_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Max Node Count\n  - example: 5",
						MarkdownDescription: "Max Node Count\n  - example: 5",
					},
					"min_node_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Min Node Count\n  - example: 1",
						MarkdownDescription: "Min Node Count\n  - example: 1",
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
						Description:         "Nodepool Name\n  - example: sample-nodepool",
						MarkdownDescription: "Nodepool Name\n  - example: sample-nodepool",
					},
					"server_group_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Server Group ID\n  - example: 2b8d33d5-4de5-40a5-a34c-7e30204133xc",
						MarkdownDescription: "Server Group ID\n  - example: 2b8d33d5-4de5-40a5-a34c-7e30204133xc",
					},
					"server_type": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"description": schema.StringAttribute{
								Computed:            true,
								Description:         "Server Type Description\n  - example: Standard",
								MarkdownDescription: "Server Type Description\n  - example: Standard",
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Server Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
								MarkdownDescription: "Server Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
							},
						},
						Computed:            true,
						Description:         "Server Type",
						MarkdownDescription: "Server Type",
					},
					"status": schema.StringAttribute{
						Computed:            true,
						Description:         "Node Pool Status\n  - example: Running",
						MarkdownDescription: "Node Pool Status\n  - example: Running",
					},
					"taints": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"effect": schema.StringAttribute{
									Computed:            true,
									Description:         "- enum: [\"NoSchedule\",\"NoExecute\",\"PreferNoSchedule\"]",
									MarkdownDescription: "- enum: [\"NoSchedule\",\"NoExecute\",\"PreferNoSchedule\"]",
								},
								"key": schema.StringAttribute{
									Computed:            true,
									Description:         "Node Pool Taint Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
									MarkdownDescription: "Node Pool Taint Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
								},
								"value": schema.StringAttribute{
									Computed:            true,
									Description:         "Node Pool Taint Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
									MarkdownDescription: "Node Pool Taint Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
								},
							},
						},
						Computed:            true,
						Description:         "Node Pool Taints",
						MarkdownDescription: "Node Pool Taints",
					},
					"volume_size": schema.Int64Attribute{
						Computed:            true,
						Description:         "Volume Size\n  - example: 104",
						MarkdownDescription: "Volume Size\n  - example: 104",
					},
					"volume_type": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"encrypt": schema.BoolAttribute{
								Computed:            true,
								Description:         "Volume Type Encrypt\n  - example: true",
								MarkdownDescription: "Volume Type Encrypt\n  - example: true",
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Volume Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
								MarkdownDescription: "Volume Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "Volume Type Name\n  - example: SSD",
								MarkdownDescription: "Volume Type Name\n  - example: SSD",
							},
						},
						Computed:            true,
						Description:         "Volume Type",
						MarkdownDescription: "Volume Type",
					},
				},
				Computed: true,
			},
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "Nodepool ID\n  - example: bdfda539-bd2e-4a5c-9021-ec6d52d1ca79",
				MarkdownDescription: "Nodepool ID\n  - example: bdfda539-bd2e-4a5c-9021-ec6d52d1ca79",
			},
		},
	}
}

func (d *skeNodepoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ske.NodepoolDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Nodepool
	data, _, err := d.client.GetNodepool(ctx, state.Id.ValueString())

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Nodepool",
			"Could not read Nodepool ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	nodepoolModel := ske.Nodepool{
		Id:                  types.StringValue(data.Nodepool.Id),
		Name:                types.StringValue(data.Nodepool.Name),
		AccountId:           types.StringValue(data.Nodepool.AccountId),
		AutoRecoveryEnabled: types.BoolValue(data.Nodepool.AutoRecoveryEnabled),
		AutoScaleEnabled:    types.BoolValue(data.Nodepool.AutoScaleEnabled),
		Cluster: ske.IdMapType{
			Id: types.StringValue(data.Nodepool.Cluster.Id),
		},
		CurrentNodeCount: types.Int32Value(data.Nodepool.CurrentNodeCount),
		DesiredNodeCount: types.Int32Value(data.Nodepool.DesiredNodeCount),
		Image: ske.Image{
			CustomImageName: types.StringPointerValue(data.Nodepool.Image.CustomImageName.Get()),
			Os:              types.StringValue(data.Nodepool.Image.Os),
			OsVersion:       types.StringValue(data.Nodepool.Image.OsVersion),
		},
		Keypair: ske.NameMapType{
			Name: types.StringValue(data.Nodepool.Keypair.Name),
		},
		KubernetesVersion: types.StringValue(data.Nodepool.KubernetesVersion),
		Labels:            converter.MakeNodepoolLabelsModel(data.Nodepool.Labels),
		Taints:            converter.MakeNodepoolTaintsModel(data.Nodepool.Taints),
		MaxNodeCount:      types.Int32Value(data.Nodepool.MaxNodeCount),
		MinNodeCount:      types.Int32Value(data.Nodepool.MinNodeCount),
		ServerType: ske.ServerType{
			Description: types.StringValue(data.Nodepool.ServerType.Description),
			Id:          types.StringValue(data.Nodepool.ServerType.Id),
		},
		Status: types.StringValue(data.Nodepool.Status),
		VolumeType: ske.VolumeType{
			Encrypt: types.BoolValue(data.Nodepool.VolumeType.Encrypt),
			Id:      types.StringValue(data.Nodepool.VolumeType.Id),
			Name:    types.StringValue(data.Nodepool.VolumeType.Name),
		},
		VolumeSize:       types.Int32Value(data.Nodepool.VolumeSize),
		CreatedAt:        types.StringValue(data.Nodepool.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.Nodepool.CreatedBy),
		ModifiedAt:       types.StringValue(data.Nodepool.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.Nodepool.ModifiedBy),
		ServerGroupId:    types.StringPointerValue(data.Nodepool.ServerGroupId.Get()),
		AdvancedSettings: converter.MakeNodepoolAdvancedSettingsModel(data.Nodepool.AdvancedSettings),
	}
	nodepoolObjectValue, diags := types.ObjectValueFrom(ctx, nodepoolModel.AttributeTypes(), nodepoolModel)
	state.Nodepool = nodepoolObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *skeNodepoolDataSource) makeExternalResourceModel(externalResource *scpske.ExternalResource) ske.ExternalResource {
	return ske.ExternalResource{
		Id:   types.StringValue(externalResource.GetId()),
		Name: types.StringValue(externalResource.GetName()),
	}
}

//// datasource.DataSourceWithConfigure Interface Methods

func (d *skeNodepoolDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
