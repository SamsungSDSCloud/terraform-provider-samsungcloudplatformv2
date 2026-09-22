package multinodegpucluster

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	multinodegpuclusterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/multinodegpucluster"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ClusterFabricDataSources{}
	_ datasource.DataSourceWithConfigure = &ClusterFabricDataSources{}
)

func NewClusterFabricDataSources() datasource.DataSource {
	return &ClusterFabricDataSources{}
}

type ClusterFabricDataSources struct {
	config  *scpsdk.Configuration
	client  *multinodegpuclusterClient.Client
	clients *client.SCPClient
}

func (multinodegpuclusterDS *ClusterFabricDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multinodegpucluster_cluster_fabrics"
}

func (multinodegpuclusterDS *ClusterFabricDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (multinodegpuclusterDS *ClusterFabricDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ClusterFabricsDataSourcesSchema()
}

func ClusterFabricsDataSourcesSchema() schema.Schema {
	return schema.Schema{
		Description: "List of Cluster Fabric",
		Attributes: map[string]schema.Attribute{
			"cluster_fabric_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Cluster Fabric name\n  - example: cluster001",
				MarkdownDescription: "Cluster Fabric name\n  - example: cluster001",
			},
			"ids": schema.ListAttribute{
				Computed:            true,
				Description:         "Cluster Fabric ID List",
				MarkdownDescription: "Cluster Fabric ID List",
				ElementType:         types.StringType,
			},
			"node_pool_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Node Pool ID\n  - example: POOL001-krw1a",
				MarkdownDescription: "Node Pool ID\n  - example: POOL001-krw1a",
			},
			"state": schema.StringAttribute{
				Optional:            true,
				Description:         "Cluster Fabric state\n  - example: ACTIVE",
				MarkdownDescription: "Cluster Fabric state\n  - example: ACTIVE",
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (multinodegpuclusterDS *ClusterFabricDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state multinodegpuclusterClient.ClusterFabricList

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := multinodegpuclusterDS.client.GetClusterFabricList(ctx, state.ClusterFabricName, state.State, state.NodePoolId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Cluster Fabrics",
			err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	contents := data.ClusterFabrics

	if len(state.Filter) > 0 {
		indices, err := filter.GetFilterIndices(contents, state.Filter)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to filter Cluster Fabrics",
				err.Error(),
			)
			return
		}

		filtered := contents[:0]
		for i, content := range contents {
			if common.Contains(indices, i) {
				filtered = append(filtered, content)
			}
		}
		contents = filtered
	}

	var ids []types.String

	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	state.Ids = ids

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
