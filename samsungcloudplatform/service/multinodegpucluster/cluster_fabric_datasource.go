package multinodegpucluster

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	multinodegpuclusterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/multinodegpucluster"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ClusterFabricDataSource{}
	_ datasource.DataSourceWithConfigure = &ClusterFabricDataSource{}
)

func NewClusterFabricDataSource() datasource.DataSource {
	return &ClusterFabricDataSource{}
}

type ClusterFabricDataSource struct {
	config  *scpsdk.Configuration
	client  *multinodegpuclusterClient.Client
	clients *client.SCPClient
}

func (multinodegpuclusterDS *ClusterFabricDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multinodegpucluster_cluster_fabric"
}

func (multinodegpuclusterDS *ClusterFabricDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (multinodegpuclusterDS *ClusterFabricDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ClusterFabricDataSourceSchema()
}

func ClusterFabricDataSourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Show Cluster Fabric",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
				MarkdownDescription: "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
			},
			"cluster_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Cluster Fabric name\n  - example: cluster001",
				MarkdownDescription: "Cluster Fabric name\n  - example: cluster001",
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
			"description": schema.StringAttribute{
				Computed:            true,
				Description:         "Description\n  - example: cluster fabric description",
				MarkdownDescription: "Description\n  - example: cluster fabric description",
			},
			"gpu_node_details": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "GPU Node details",
				MarkdownDescription: "GPU Node details",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"gpu_node_id": schema.StringAttribute{
							Computed:            true,
							Description:         "GPU Node ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
							MarkdownDescription: "GPU Node ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
						},
						"gpu_node_name": schema.StringAttribute{
							Computed:            true,
							Description:         "GPU Node name\n  - example: gpu-node-001",
							MarkdownDescription: "GPU Node name\n  - example: gpu-node-001",
						},
						"policy_ip": schema.StringAttribute{
							Computed:            true,
							Description:         "Policy IP\n  - example: 192.168.0.1",
							MarkdownDescription: "Policy IP\n  - example: 192.168.0.1",
						},
						"product_type_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Product type ID\n  - example: f90e8ef54cc2451b825608e9f95f7bcb",
							MarkdownDescription: "Product type ID\n  - example: f90e8ef54cc2451b825608e9f95f7bcb",
						},
						"server_type": schema.StringAttribute{
							Computed:            true,
							Description:         "Server type\n  - example: s1v8m32_metal",
							MarkdownDescription: "Server type\n  - example: s1v8m32_metal",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "GPU Node state\n  - example: RUNNING",
							MarkdownDescription: "GPU Node state\n  - example: RUNNING",
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster Fabric ID",
				MarkdownDescription: "Cluster Fabric ID",
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
			"node_pool_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Node Pool ID\n  - example: POOL001-krw1a",
				MarkdownDescription: "Node Pool ID\n  - example: POOL001-krw1a",
			},
			"pirp_id": schema.StringAttribute{
				Computed:            true,
				Description:         "PIRP ID\n  - example: PIRP-001",
				MarkdownDescription: "PIRP ID\n  - example: PIRP-001",
			},
			"product_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Product ID\n  - example: f90e8ef54cc2451b825608e9f95f7bcb",
				MarkdownDescription: "Product ID\n  - example: f90e8ef54cc2451b825608e9f95f7bcb",
			},
			"region_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Region ID\n  - example: kr-west1",
				MarkdownDescription: "Region ID\n  - example: kr-west1",
			},
			"server_type": schema.StringAttribute{
				Computed:            true,
				Description:         "Server type\n  - example: s1v8m32_metal",
				MarkdownDescription: "Server type\n  - example: s1v8m32_metal",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				Description:         "Cluster Fabric state\n  - example: ACTIVE",
				MarkdownDescription: "Cluster Fabric state\n  - example: ACTIVE",
			},
			"used_server_count": schema.Int64Attribute{
				Computed:            true,
				Description:         "Used server count\n  - example: 2",
				MarkdownDescription: "Used server count\n  - example: 2",
			},
		},
	}
}

func (multinodegpuclusterDS *ClusterFabricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state multinodegpuclusterClient.ClusterFabricDataSource

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := multinodegpuclusterDS.client.GetClusterFabric(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Cluster Fabric",
			err.Error()+ERROR_EXPLAIN+detail,
		)
		return
	}

	state.AccountId = types.StringValue(data.AccountId)
	state.ClusterName = types.StringValue(data.ClusterName)
	state.CreatedBy = types.StringValue(data.CreatedBy)
	state.CreatedAt = types.StringValue(data.CreatedAt.String())
	state.Description = types.StringValue(data.Description)
	state.NodePoolId = types.StringValue(data.NodePoolId)
	state.ModifiedBy = types.StringValue(data.ModifiedBy)
	state.ModifiedAt = types.StringValue(data.ModifiedAt.String())
	state.ProductId = types.StringValue(data.ProductId)
	state.RegionId = types.StringValue(data.RegionId)
	state.State = types.StringValue(data.State)

	if data.PirpId.IsSet() {
		state.PirpId = types.StringPointerValue(data.PirpId.Get())
	}
	if data.ServerType.IsSet() {
		state.ServerType = types.StringPointerValue(data.ServerType.Get())
	}
	if data.UsedServerCount.IsSet() {
		state.UsedServerCount = types.Int64Value(int64(*data.UsedServerCount.Get()))
	}

	gpuNodeDetails := make([]multinodegpuclusterClient.GpuNodeDetailsValue, 0)
	for _, detail := range data.GpuNodeDetails {
		gpuNodeDetail := multinodegpuclusterClient.GpuNodeDetailsValue{
			GpuNodeId:     types.StringValue(detail.GpuNodeId),
			GpuNodeName:   types.StringValue(detail.GpuNodeName),
			PolicyIp:      types.StringValue(detail.PolicyIp),
			ProductTypeId: types.StringValue(detail.ProductTypeId),
			State:         types.StringValue(detail.State),
		}
		if detail.ServerType.IsSet() {
			gpuNodeDetail.ServerType = types.StringPointerValue(detail.ServerType.Get())
		}
		gpuNodeDetails = append(gpuNodeDetails, gpuNodeDetail)
	}

	state.GpuNodeDetails, diags = types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: multinodegpuclusterClient.GpuNodeDetailsValue{}.AttributeTypes(),
	}, gpuNodeDetails)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
