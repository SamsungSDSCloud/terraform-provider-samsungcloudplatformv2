package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
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
				Description: "ClusterId\n  - example: 628a6c3f05454f2699da171a0c2f50b1",
				Required:    true,
			},
			common.ToSnakeCase("Nodepools"): schema.ListNestedAttribute{
				Description: "A list of nodepool.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id\n  - example: 0d9ef630-a557-48f3-b6fa-04d054834d11",
							Computed:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Name\n  - example: test-nodepool",
							Computed:    true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "AccountId\n  - example: 27bb070b564349f8a31cc60734cc36a5",
							Computed:    true,
						},
						common.ToSnakeCase("AutoRecoveryEnabled"): schema.BoolAttribute{
							Description: "AutoRecoveryEnabled\n  - example: false",
							Computed:    true,
						},
						common.ToSnakeCase("AutoScaleEnabled"): schema.BoolAttribute{
							Description: "AutoScaleEnabled\n  - example: false",
							Computed:    true,
						},
						common.ToSnakeCase("CurrentNodeCount"): schema.Int32Attribute{
							Description: "Name\n  - example: 1",
							Computed:    true,
						},
						common.ToSnakeCase("DesiredNodeCount"): schema.Int32Attribute{
							Description: "AccountId\n  - example: 1",
							Computed:    true,
						},
						common.ToSnakeCase("Image"): schema.SingleNestedAttribute{
							Description: "Image",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("CustomImageName"): schema.StringAttribute{
									Description: "CustomImageName\n  - example: null",
									Computed:    true,
								},
								common.ToSnakeCase("Os"): schema.StringAttribute{
									Description: "Os\n  - example: ubuntu",
									Computed:    true,
								},
								common.ToSnakeCase("OsVersion"): schema.StringAttribute{
									Description: "OsVersion\n  - example: 22.04",
									Computed:    true,
								},
							},
						},
						common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
							Description: "KubernetesVersion\n  - example: v1.31.8",
							Computed:    true,
						},
						common.ToSnakeCase("ServerType"): schema.SingleNestedAttribute{
							Description: "ServerType",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Description"): schema.StringAttribute{
									Description: "Description\n  - example: vCPU 1 | Memory 2G",
									Computed:    true,
								},
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "Id\n  - example: s1v1m2",
									Computed:    true,
								},
							},
						},
						common.ToSnakeCase("Status"): schema.StringAttribute{
							Description: "Status\n  - example: Running",
							Computed:    true,
						},
						common.ToSnakeCase("VolumeType"): schema.SingleNestedAttribute{
							Description: "VolumeType",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "Id\n  - example: a6d4a8a2-4db1-45bb-b85c-9f3a57b304c6",
									Computed:    true,
								},
								common.ToSnakeCase("Name"): schema.StringAttribute{
									Description: "Name\n  - example: SSD",
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
		nodepoolState := ske.NodepoolSummary{
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
