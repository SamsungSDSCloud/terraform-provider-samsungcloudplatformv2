package ske

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &skeClusterDeletionProtectionDataSource{}
	_ datasource.DataSourceWithConfigure = &skeClusterDeletionProtectionDataSource{}
)

func NewSkeClusterDeletionProtectionDataSource() datasource.DataSource {
	return &skeClusterDeletionProtectionDataSource{}
}

type skeClusterDeletionProtectionDataSource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

func (d *skeClusterDeletionProtectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_cluster_deletion_protection"
}

func (d *skeClusterDeletionProtectionDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Gets the deletion protection status of an SKE cluster.",
		MarkdownDescription: "Gets the deletion protection status of an SKE cluster.",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Description:         "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				MarkdownDescription: "Cluster ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				Required:            true,
			},
			"deletion_protection_enabled": schema.BoolAttribute{
				Description:         "Deletion protection flag. True when cluster deletion is prevented.\n  - example: true",
				MarkdownDescription: "Deletion protection flag. True when cluster deletion is prevented.\n  - example: true",
				Computed:            true,
			},
		},
	}
}

func (d *skeClusterDeletionProtectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *skeClusterDeletionProtectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config skeClusterDeletionProtectionDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deletionProtection, _, err := d.client.GetClusterDeletionProtection(ctx, config.ClusterId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SKE Cluster Deletion Protection",
			"Could not read deletion protection for cluster ID "+config.ClusterId.ValueString()+": "+err.Error(),
		)
		return
	}

	config.DeletionProtectionEnabled = types.BoolValue(deletionProtection)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

type skeClusterDeletionProtectionDataSourceModel struct {
	ClusterId          types.String `tfsdk:"cluster_id"`
	DeletionProtectionEnabled types.Bool   `tfsdk:"deletion_protection_enabled"`
}
