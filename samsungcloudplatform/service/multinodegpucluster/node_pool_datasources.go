package multinodegpucluster

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	multinodegpuclusterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/multinodegpucluster"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                     = &NodePoolDataSources{}
	_ datasource.DataSourceWithConfigure        = &NodePoolDataSources{}
	_ datasource.DataSourceWithConfigValidators = &NodePoolDataSources{}
)

func NewNodePoolDataSources() datasource.DataSource {
	return &NodePoolDataSources{}
}

type NodePoolDataSources struct {
	config  *scpsdk.Configuration
	client  *multinodegpuclusterClient.Client
	clients *client.SCPClient
}

func (multinodegpuclusterDS *NodePoolDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multinodegpucluster_node_pools"
}

func (multinodegpuclusterDS *NodePoolDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (multinodegpuclusterDS *NodePoolDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = NodePoolsDataSourcesSchema()
}

func (multinodegpuclusterDS *NodePoolDataSources) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("cluster_fabric_id"),
			path.MatchRoot("subnet_id"),
		),
		datasourcevalidator.RequiredTogether(
			path.MatchRoot("subnet_id"),
			path.MatchRoot("zone"),
		),
	}
}

func NodePoolsDataSourcesSchema() schema.Schema {
	return schema.Schema{
		Description: "List of Node Pool",
		Attributes: map[string]schema.Attribute{
			"cluster_fabric_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Cluster Fabric ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
				MarkdownDescription: "Cluster Fabric ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
			},
			"ids": schema.ListAttribute{
				Computed:            true,
				Description:         "Node Pool ID List",
				MarkdownDescription: "Node Pool ID List",
				ElementType:         types.StringType,
			},
			"node_pool_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Node Pool ID\n  - example: POOL001-krw1a",
				MarkdownDescription: "Node Pool ID\n  - example: POOL001-krw1a",
			},
			"subnet_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Subnet ID (required with zone when cluster_fabric_id is empty)\n  - example: ab313c43291e4b678f4bacffe10768ae",
				MarkdownDescription: "Subnet ID (required with zone when cluster_fabric_id is empty)\n  - example: ab313c43291e4b678f4bacffe10768ae",
			},
			"zone": schema.StringAttribute{
				Optional:            true,
				Description:         "Availability Zone (required when subnet_id is provided)\n  - example: kr-west1-a",
				MarkdownDescription: "Availability Zone (required when subnet_id is provided)\n  - example: kr-west1-a",
			},
		},
	}
}

func (multinodegpuclusterDS *NodePoolDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state multinodegpuclusterClient.NodePoolList

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := multinodegpuclusterDS.client.GetNodePoolList(ctx, state.SubnetId, state.ClusterFabricId, state.NodePoolId, state.Zone)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Node Pools",
			err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	var ids []types.String

	for _, content := range data.NodePools {
		ids = append(ids, types.StringValue(content.NodePoolId))
	}

	state.Ids = ids

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
