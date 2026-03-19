package multinodegpucluster

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	multinodegpuclusterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/multinodegpucluster"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GpunodeDataSource{}
	_ datasource.DataSourceWithConfigure = &GpunodeDataSource{}
)

// NewDataSource is a helper function to simplify the provider implementation.
func NewGpunodeDataSource() datasource.DataSource {
	return &GpunodeDataSource{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type GpunodeDataSource struct {
	config  *scpsdk.Configuration
	client  *multinodegpuclusterClient.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (multinodegpuclusterDS *GpunodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multinodegpucluster_gpunode"
}

// Configure adds the provider configured client to the data source.
func (multinodegpuclusterDS *GpunodeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (multinodegpuclusterDS *GpunodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GpuNodeDataSourceSchema()
}

func GpuNodeDataSourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Show GPU Node",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
				MarkdownDescription: "Account ID\n  - example: f5c8e56a4d9b49a8bd89e14758a32d53",
			},
			"cluster_fabric_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Cluster Fabric ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
				MarkdownDescription: "Cluster Fabric ID\n  - example: 20c507a036c447cdb3b19468d8ea62ac",
			},
			"cluster_fabric_name": schema.StringAttribute{
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
			"gpu_node_name": schema.StringAttribute{
				Computed:            true,
				Description:         "GPU Node name\n  - example: gpu-node-001",
				MarkdownDescription: "GPU Node name\n  - example: gpu-node-001",
			},
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "GPU Node ID",
				MarkdownDescription: "GPU Node ID",
			},
			"image_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Image ID\n  - example: IMAGE-7XFMaJpLsapKvskFMjCtmm",
				MarkdownDescription: "Image ID\n  - example: IMAGE-7XFMaJpLsapKvskFMjCtmm",
			},
			"image_version": schema.StringAttribute{
				Computed:            true,
				Description:         "Image version\n  - example: RHEL 8.7 for BM",
				MarkdownDescription: "Image version\n  - example: RHEL 8.7 for BM",
			},
			"init_script": schema.StringAttribute{
				Computed:            true,
				Description:         "Init script\n  - example: init script",
				MarkdownDescription: "Init script\n  - example: init script",
			},
			"lock_enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Use Lock\n  - example: true",
				MarkdownDescription: "Use Lock\n  - example: true",
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
			"network_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Network ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
				MarkdownDescription: "Network ID\n  - example: ab313c43291e4b678f4bacffe10768ae",
			},
			"node_pool_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Node Pool ID\n  - example: POOL001-krw1a",
				MarkdownDescription: "Node Pool ID\n  - example: POOL001-krw1a",
			},
			"os_type": schema.StringAttribute{
				Computed:            true,
				Description:         "OS type\n  - example: WINDOWS",
				MarkdownDescription: "OS type\n  - example: WINDOWS",
			},
			"pfs_ip": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "PFS IP list\n  - example: [10.252.128.2 10.252.128.3]",
				MarkdownDescription: "PFS IP list\n  - example: [10.252.128.2 10.252.128.3]",
			},
			"policy_ip": schema.StringAttribute{
				Computed:            true,
				Description:         "Policy IP\n  - example: 192.168.0.1",
				MarkdownDescription: "Policy IP\n  - example: 192.168.0.1",
			},
			"policy_nat": schema.StringAttribute{
				Computed:            true,
				Description:         "Policy NAT\n  - example: 192.168.0.1",
				MarkdownDescription: "Policy NAT\n  - example: 192.168.0.1",
			},
			"policy_use_nat": schema.BoolAttribute{
				Computed:            true,
				Description:         "Policy use NAT\n  - example: true",
				MarkdownDescription: "Policy use NAT\n  - example: true",
			},
			"product_type_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Product type ID\n  - example: f90e8ef54cc2451b825608e9f95f7bcb",
				MarkdownDescription: "Product type ID\n  - example: f90e8ef54cc2451b825608e9f95f7bcb",
			},
			"region_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Region ID\n  - example: kr-west1",
				MarkdownDescription: "Region ID\n  - example: kr-west1",
			},
			"root_account": schema.StringAttribute{
				Computed:            true,
				Description:         "Root Account\n  - example: rootaccount",
				MarkdownDescription: "Root Account\n  - example: rootaccount",
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
			"time_zone": schema.StringAttribute{
				Computed:            true,
				Description:         "Time Zone\n  - example: Asia/Seoul",
				MarkdownDescription: "Time Zone\n  - example: Asia/Seoul",
			},
			"vpc_id": schema.StringAttribute{
				Computed:            true,
				Description:         "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
				MarkdownDescription: "VPC ID\n  - example: e58348b1bc9148e5af86500fd4ef99ca",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (multinodegpuclusterDS *GpunodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state multinodegpuclusterClient.GpuNodeDataSource

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := multinodegpuclusterDS.client.GetGpuNode(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read GPU Node",
			err.Error(),
		)
		return
	}

	// // Map response body to state
	state.AccountId = types.StringValue(data.AccountId)
	state.RegionId = types.StringValue(data.RegionId)

	state.ImageId = types.StringValue(data.ImageId)
	state.ImageVersion = types.StringValue(data.ImageVersion)
	state.OsType = types.StringValue(data.OsType)
	state.ProductTypeId = types.StringValue(data.ProductTypeId)
	state.RootAccount = types.StringValue(data.RootAccount)
	state.ServerType = types.StringValue(data.ServerType)
	state.State = types.StringValue(data.State)
	state.TimeZone = types.StringValue(data.TimeZone)

	state.GpuNodeName = types.StringValue(data.GpuNodeName)
	state.ClusterFabricName = types.StringValue(data.ClusterFabricName)
	if data.ClusterFabricId.IsSet() {
		state.ClusterFabricId = types.StringPointerValue(data.ClusterFabricId.Get())
	}
	state.NodePoolId = types.StringValue(data.NodePoolId)

	state.VpcId = types.StringValue(data.VpcId)
	state.NetworkId = types.StringValue(data.NetworkId)
	state.PolicyIp = types.StringValue(data.PolicyIp)
	if data.PolicyNat.IsSet() {
		state.PolicyNat = types.StringPointerValue(data.PolicyNat.Get())
	}
	state.PolicyUseNat = types.BoolPointerValue(data.PolicyUseNat)
	state.LockEnabled = types.BoolPointerValue(data.LockEnabled)
	state.InitScript = types.StringValue(data.InitScript)

	state.PfsIp, _ = types.ListValueFrom(ctx, types.StringType, data.PfsIp)

	// metadata info
	state.CreatedBy = types.StringValue(data.CreatedBy)
	state.CreatedAt = types.StringValue(data.CreatedAt.String())
	state.ModifiedBy = types.StringValue(data.ModifiedBy)
	state.ModifiedAt = types.StringValue(data.ModifiedAt.String())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
