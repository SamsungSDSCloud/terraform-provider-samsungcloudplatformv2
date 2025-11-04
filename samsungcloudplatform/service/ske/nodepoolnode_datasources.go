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
	_ datasource.DataSource              = &skeNodepoolnodeDataSources{}
	_ datasource.DataSourceWithConfigure = &skeNodepoolnodeDataSources{}
)

func NewSkeNodepoolnodeDataSources() datasource.DataSource {
	return &skeNodepoolnodeDataSources{}
}

type skeNodepoolnodeDataSources struct {
	config *scpsdk.Configuration
	client *ske.Client
}

// Metadata returns the data source type name.
func (d *skeNodepoolnodeDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_nodepoolnodes"
}

// Schema defines the schema for the data source.
func (d *skeNodepoolnodeDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of nodepoolnode.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("NodepoolId"): schema.StringAttribute{
				Description: "NodepoolId",
				Required:    true,
			},
			common.ToSnakeCase("Nodes"): schema.ListNestedAttribute{
				Description: "A list of node",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Name",
							Computed:    true,
						},
						common.ToSnakeCase("KubernetesVersion"): schema.StringAttribute{
							Description: "KubernetesVersion",
							Computed:    true,
						},
						common.ToSnakeCase("Status"): schema.StringAttribute{
							Description: "Status",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *skeNodepoolnodeDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *skeNodepoolnodeDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ske.NodepoolnodeDataSources

	diags := req.Config.Get(ctx, &state)
	fmt.Println(diags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetNodepoolNodeList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Nodepools",
			err.Error(),
		)
		return
	}

	contents := data.Nodes

	// Map response body to model
	for _, node := range contents {
		nodeState := ske.NodeInNodepool{
			Name:              types.StringPointerValue(node.Name),
			KubernetesVersion: types.StringPointerValue(node.KubernetesVersion),
			Status:            types.StringPointerValue(node.Status),
		}

		state.Nodes = append(state.Nodes, nodeState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
