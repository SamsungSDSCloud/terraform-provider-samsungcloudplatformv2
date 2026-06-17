package multinodegpucluster

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	multinodegpuclusterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/multinodegpucluster"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GpunodeDataSources{}
	_ datasource.DataSourceWithConfigure = &GpunodeDataSources{}
)

// NewcDataSources is a helper function to simplify the provider implementation.
func NewGpunodeDataSources() datasource.DataSource {
	return &GpunodeDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type GpunodeDataSources struct {
	config  *scpsdk.Configuration
	client  *multinodegpuclusterClient.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (multinodegpuclusterDS *GpunodeDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multinodegpucluster_gpunodes"
}

// Configure adds the provider configured client to the data source.
func (multinodegpuclusterDS *GpunodeDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expect *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	multinodegpuclusterDS.client = inst.Client.Mngc
	multinodegpuclusterDS.clients = inst.Client
}

// Schema defines the schema for the data source.
func (multinodegpuclusterDS *GpunodeDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GpuNodesDataSourcesSchema()
}

func GpuNodesDataSourcesSchema() schema.Schema {
	return schema.Schema{
		Description: "List of GPU Node",
		Attributes: map[string]schema.Attribute{
			"cluster_fabric_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Cluster Fabric ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
				MarkdownDescription: "Cluster Fabric ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
			},
			"cluster_fabric_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Cluster Fabric name\n  - example: cluster001",
				MarkdownDescription: "Cluster Fabric name\n  - example: cluster001",
			},
			"gpu_node_name": schema.StringAttribute{
				Optional:            true,
				Description:         "GPU Node name\n  - example: gpu-node-001",
				MarkdownDescription: "GPU Node name\n  - example: gpu-node-001",
			},
			"ids": schema.ListAttribute{
				Computed:            true,
				Description:         "GPU Node ID List\n  - example: ['aaaa8b745fa04852aad2aaa1f9907f2b','bbba8b745fa04852aad2aaa1f9907f2b']",
				MarkdownDescription: "GPU Node ID List\n  - example: ['aaaa8b745fa04852aad2aaa1f9907f2b','bbba8b745fa04852aad2aaa1f9907f2b']",
				ElementType:         types.StringType,
			},
			"ip": schema.StringAttribute{
				Optional:            true,
				Description:         "IP\n  - example: 192.168.0.1",
				MarkdownDescription: "IP\n  - example: 192.168.0.1",
			},
			"state": schema.StringAttribute{
				Optional:            true,
				Description:         "GPU Node state\n  - example: RUNNING",
				MarkdownDescription: "GPU Node state\n  - example: RUNNING",
			},
			"vpc_id": schema.StringAttribute{
				Optional:            true,
				Description:         "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
				MarkdownDescription: "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (multinodegpuclusterDS *GpunodeDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state multinodegpuclusterClient.GpuNodeList

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := multinodegpuclusterDS.client.GetGpuNodeList(ctx, state.GpuNodeName, state.State, state.Ip, state.VpcId, state.ClusterFabricName, state.ClusterFabricId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Resource Group",
			err.Error(),
		)
	}

	var ids []types.String

	// Map response body to id
	for _, content := range data.GpuNodes {
		ids = append(ids, types.StringValue(content.Id))
	}

	state.Ids = ids

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
