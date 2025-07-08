package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &skeNodepoolDataSources{}
	_ datasource.DataSourceWithConfigure = &skeNodepoolDataSources{}
)

func NewSkeNodepoolDataSources() datasource.DataSource {
	return &skeNodepoolDataSources{}
}

type skeNodepoolDataSources struct {
	config *scpsdk.Configuration
	client *ske.Client
}

// Metadata returns the data source type name.
func (d *skeNodepoolDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_nodepools"
}

// Schema defines the schema for the data source.
func (d *skeNodepoolDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of nodepool.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("ClusterId"): schema.StringAttribute{
				Description: "ClusterId",
				Required:    true,
			},
			common.ToSnakeCase("Nodepools"): schema.ListNestedAttribute{
				Description: "A list of nodepool.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
						common.ToSnakeCase("AutoRecoveryEnabled"): schema.BoolAttribute{
							Description: "AutoRecoveryEnabled",
							Computed:    true,
						},
						common.ToSnakeCase("AutoScaleEnabled"): schema.BoolAttribute{
							Description: "AutoScaleEnabled",
							Computed:    true,
						},
						common.ToSnakeCase("CurrentNodeCount"): schema.Int32Attribute{
							Description: "Name",
							Computed:    true,
						},
						common.ToSnakeCase("DesiredNodeCount"): schema.Int32Attribute{
							Description: "AccountId",
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
						common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
							Description: "KubernetesVersion",
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
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *skeNodepoolDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
}

// Read refreshes the Terraform state with the latest data.
func (d *skeNodepoolDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ske.NodepoolDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetNodePoolList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Nodepools",
			err.Error(),
		)
		return
	}

	contents := data.Nodepools

	// Map response body to model
	for _, nodepool := range contents {
		nodepoolState := ske.Nodepool{
			Id:                  types.StringValue(nodepool.Id),
			Name:                types.StringValue(nodepool.Name),
			AccountId:           types.StringValue(nodepool.AccountId),
			AutoRecoveryEnabled: types.BoolValue(nodepool.AutoRecoveryEnabled),
			AutoScaleEnabled:    types.BoolValue(nodepool.AutoScaleEnabled),
			CurrentNodeCount:    types.Int32Value(nodepool.CurrentNodeCount),
			DesiredNodeCount:    types.Int32Value(nodepool.DesiredNodeCount),
			Image: ske.Image{
				CustomImageName: types.StringPointerValue(nodepool.Image.CustomImageName.Get()),
				Os:              types.StringValue(nodepool.Image.Os),
				OsVersion:       types.StringValue(nodepool.Image.OsVersion),
			},
			KubernetesVersion: types.StringValue(nodepool.KubernetesVersion),
			ServerType: ske.ServerType{
				Description: types.StringValue(nodepool.ServerType.Description),
				Id:          types.StringValue(nodepool.ServerType.Id),
			},
			Status: types.StringValue(nodepool.Status),
			VolumeType: ske.VolumeTypeSummary{
				Id:   types.StringValue(nodepool.VolumeType.Id),
				Name: types.StringValue(nodepool.VolumeType.Name),
			},
		}

		state.Nodepools = append(state.Nodepools, nodepoolState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
