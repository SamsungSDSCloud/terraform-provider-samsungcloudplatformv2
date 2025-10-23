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
		Description: "nodepool",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("Nodepool"): schema.SingleNestedAttribute{
				Description: "Nodepool",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "id",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "name",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "account id",
						Computed:    true,
					},
					common.ToSnakeCase("AutoRecoveryEnabled"): schema.BoolAttribute{
						Description: "AutoRecoveryEnabled",
						Computed:    true,
					},
					common.ToSnakeCase("AutoScaleEnabled"): schema.BoolAttribute{
						Description: "AutoScaleEnabled",
						Computed:    true,
					},
					common.ToSnakeCase("Cluster"): schema.SingleNestedAttribute{
						Description: "Cluster",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "Id",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("CurrentNodeCount"): schema.Int32Attribute{
						Description: "CurrentNodeCount",
						Computed:    true,
					},
					common.ToSnakeCase("DesiredNodeCount"): schema.Int32Attribute{
						Description: "DesiredNodeCount",
						Computed:    true,
					},
					common.ToSnakeCase("Image"): schema.SingleNestedAttribute{
						Description: "Image",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("CustomImageName"): schema.StringAttribute{
								Description: "CustomImageName",
								Computed:    true,
							},
							common.ToSnakeCase("Os"): schema.StringAttribute{
								Description: "Os",
								Computed:    true,
							},
							common.ToSnakeCase("OsVersion"): schema.StringAttribute{
								Description: "OsVersion",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("Keypair"): schema.SingleNestedAttribute{
						Description: "Keypair",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Name"): schema.StringAttribute{
								Description: "Name",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
						Description: "KubernetesVersion",
						Computed:    true,
					},
					common.ToSnakeCase("Labels"): schema.ListNestedAttribute{
						Description: "Labels",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: d.makeNodepoolLabelsSchema(),
						},
					},
					common.ToSnakeCase("Taints"): schema.ListNestedAttribute{
						Description: "Taints",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: d.makeNodepoolTaintsSchema(),
						},
					},
					common.ToSnakeCase("MaxNodeCount"): schema.Int32Attribute{
						Description: "MaxNodeCount",
						Computed:    true,
					},
					common.ToSnakeCase("MinNodeCount"): schema.Int32Attribute{
						Description: "MinNodeCount",
						Computed:    true,
					},
					common.ToSnakeCase("ServerType"): schema.SingleNestedAttribute{
						Description: "ServerType",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Description"): schema.StringAttribute{
								Description: "Description",
								Computed:    true,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "Id",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Computed:    true,
					},
					common.ToSnakeCase("VolumeType"): schema.SingleNestedAttribute{
						Description: "VolumeType",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Encrypt"): schema.BoolAttribute{
								Description: "Encrypt",
								Computed:    true,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "Id",
								Computed:    true,
							},
							common.ToSnakeCase("Name"): schema.StringAttribute{
								Description: "Name",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("VolumeSize"): schema.Int32Attribute{
						Description: "VolumeSize",
						Computed:    true,
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
					common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
						Description: "ServerGroupId",
						Computed:    true,
					},
					common.ToSnakeCase("AdvancedSettings"): schema.SingleNestedAttribute{
						Description: "AdvancedSettings",
						Computed:    true,
						Required:    false,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("AllowedUnsafeSysctls"): schema.StringAttribute{
								Description: "AllowedUnsafeSysctls",
								Computed:    true,
							},
							common.ToSnakeCase("ContainerLogMaxFiles"): schema.Int32Attribute{
								Description: "ContainerLogMaxFiles",
								Computed:    true,
							},
							common.ToSnakeCase("ContainerLogMaxSize"): schema.Int32Attribute{
								Description: "ContainerLogMaxSize",
								Computed:    true,
							},
							common.ToSnakeCase("ImageGcHighThreshold"): schema.Int32Attribute{
								Description: "ImageGcHighThreshold",
								Computed:    true,
							},
							common.ToSnakeCase("ImageGcLowThreshold"): schema.Int32Attribute{
								Description: "ImageGcLowThreshold",
								Computed:    true,
							},
							common.ToSnakeCase("MaxPods"): schema.Int32Attribute{
								Description: "MaxPods",
								Computed:    true,
							},
							common.ToSnakeCase("PodMaxPids"): schema.Int32Attribute{
								Description: "PodMaxPids",
								Computed:    true,
							},
						},
					},
				},
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

	var labels []ske.Labels
	for _, label := range data.Nodepool.Labels {
		labels = append(labels, d.makeNodepoolLabelsModel((*scpske.NodepoolLabel)(&label)))
	}

	var taints []ske.Taints
	for _, taint := range data.Nodepool.Taints {
		taints = append(taints, d.makeNodepoolTaintsModel((*scpske.NodepoolTaint)(&taint)))
	}

	nodepoolModel := ske.NodepoolDetail{
		Id:                  types.StringPointerValue(data.Nodepool.Id),
		Name:                types.StringPointerValue(data.Nodepool.Name),
		AccountId:           types.StringPointerValue(data.Nodepool.AccountId),
		AutoRecoveryEnabled: types.BoolPointerValue(data.Nodepool.AutoRecoveryEnabled),
		AutoScaleEnabled:    types.BoolPointerValue(data.Nodepool.AutoScaleEnabled),
		Cluster: ske.IdMapType{
			Id: types.StringPointerValue(data.Nodepool.Cluster.Id),
		},
		CurrentNodeCount: types.Int32PointerValue(data.Nodepool.CurrentNodeCount),
		DesiredNodeCount: types.Int32PointerValue(data.Nodepool.DesiredNodeCount),
		Image: ske.Image{
			CustomImageName: types.StringPointerValue(data.Nodepool.Image.CustomImageName.Get()),
			Os:              types.StringPointerValue(data.Nodepool.Image.Os),
			OsVersion:       types.StringPointerValue(data.Nodepool.Image.OsVersion),
		},
		Keypair: ske.NameMapType{
			Name: types.StringPointerValue(data.Nodepool.Keypair.Name),
		},
		KubernetesVersion: types.StringPointerValue(data.Nodepool.KubernetesVersion),
		Labels:            labels,
		Taints:            taints,
		MaxNodeCount:      types.Int32PointerValue(data.Nodepool.MaxNodeCount),
		MinNodeCount:      types.Int32PointerValue(data.Nodepool.MinNodeCount),
		ServerType: ske.ServerType{
			Description: types.StringPointerValue(data.Nodepool.ServerType.Description),
			Id:          types.StringPointerValue(data.Nodepool.ServerType.Id),
		},
		Status: types.StringPointerValue(data.Nodepool.Status),
		VolumeType: ske.VolumeType{
			Encrypt: types.BoolPointerValue(data.Nodepool.VolumeType.Encrypt),
			Id:      types.StringPointerValue(data.Nodepool.VolumeType.Id),
			Name:    types.StringPointerValue(data.Nodepool.VolumeType.Name),
		},
		VolumeSize:       types.Int32PointerValue(data.Nodepool.VolumeSize),
		CreatedAt:        types.StringValue(data.Nodepool.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.Nodepool.CreatedBy),
		ModifiedAt:       types.StringValue(data.Nodepool.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.Nodepool.ModifiedBy),
		ServerGroupId:    types.StringPointerValue(data.Nodepool.ServerGroupId.Get()),
		AdvancedSettings: d.makeNodepoolAdvancedSettingsModel(data.Nodepool.AdvancedSettings),
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
func (d *skeNodepoolDataSource) makeNodepoolLabelsModel(nodepoolLabels *scpske.NodepoolLabel) ske.Labels {
	return ske.Labels{
		Key:   types.StringValue(nodepoolLabels.GetKey()),
		Value: types.StringValue(nodepoolLabels.GetValue()),
	}
}

func (d *skeNodepoolDataSource) makeNodepoolTaintsModel(nodepoolTaints *scpske.NodepoolTaint) ske.Taints {
	return ske.Taints{
		Effect: types.StringValue(string(nodepoolTaints.GetEffect())),
		Key:    types.StringValue(nodepoolTaints.GetKey()),
		Value:  types.StringValue(nodepoolTaints.GetValue()),
	}
}

func (d *skeNodepoolDataSource) makeNodepoolAdvancedSettingsModel(advancedSettings scpske.NullableNodepoolAdvancedSettings) *ske.AdvancedSettings {
	value := advancedSettings.Get()
	if value != nil {
		return &ske.AdvancedSettings{
			AllowedUnsafeSysctls: types.StringValue(value.AllowedUnsafeSysctls),
			ContainerLogMaxFiles: types.Int32Value(value.ContainerLogMaxFiles),
			ContainerLogMaxSize:  types.Int32Value(value.ContainerLogMaxSize),
			ImageGcHighThreshold: types.Int32Value(value.ImageGcHighThreshold),
			ImageGcLowThreshold:  types.Int32Value(value.ImageGcLowThreshold),
			MaxPods:              types.Int32Value(value.MaxPods),
			PodMaxPids:           types.Int32Value(value.PodMaxPids),
		}
	}
	return nil
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

// private
func (d *skeNodepoolDataSource) makeExternalResourceSchema() map[string]schema.Attribute {
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

func (d *skeNodepoolDataSource) makeNodepoolLabelsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		common.ToSnakeCase("Key"): schema.StringAttribute{
			Description: "Key",
			Computed:    true,
		},
		common.ToSnakeCase("Value"): schema.StringAttribute{
			Description: "Value",
			Computed:    true,
		},
	}
}

func (d *skeNodepoolDataSource) makeNodepoolTaintsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		common.ToSnakeCase("Effect"): schema.StringAttribute{
			Description: "Effect",
			Computed:    true,
		},
		common.ToSnakeCase("Key"): schema.StringAttribute{
			Description: "Key",
			Computed:    true,
		},
		common.ToSnakeCase("Value"): schema.StringAttribute{
			Description: "Value",
			Computed:    true,
		},
	}
}
