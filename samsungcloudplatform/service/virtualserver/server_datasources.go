package virtualserver

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualServerServerDataSources{}
	_ datasource.DataSourceWithConfigure = &virtualServerServerDataSources{}
)

func NewVirtualServerServerDataSources() datasource.DataSource {
	return &virtualServerServerDataSources{}
}

type virtualServerServerDataSources struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (d *virtualServerServerDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_servers"
}

func (d *virtualServerServerDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a list of virtual servers.",
		MarkdownDescription: "Retrieves a list of virtual server IDs matching the specified criteria.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description:         "Server name to filter by.\n  - example: my-server\n  - minLength: 1\n  - maxLength: 255",
				MarkdownDescription: "Server name to filter by.\n  - example: my-server\n  - minLength: 1\n  - maxLength: 255",
				Optional:            true,
			},
			common.ToSnakeCase("Ip"): schema.StringAttribute{
				Description:         "IP address to filter servers.\n  - example: 192.168.1.100",
				MarkdownDescription: "IP address to filter servers.\n  - example: 192.168.1.100",
				Optional:            true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description:         "Server state to filter.\n  - Available values: ACTIVE, SHUTOFF, ERROR",
				MarkdownDescription: "Server state to filter.\n  - Available values: ACTIVE, SHUTOFF, ERROR",
				Optional:            true,
			},
			common.ToSnakeCase("ProductCategory"): schema.StringAttribute{
				Description:         "Product category.\n  - Available values: compute, container",
				MarkdownDescription: "Product category.\n  - Available values: compute, container",
				Optional:            true,
			},
			common.ToSnakeCase("ProductOffering"): schema.StringAttribute{
				Description:         "Product offering.\n  - Available values: virtual_server, gpu_server, k8s_vm, k8s_gpu_vm",
				MarkdownDescription: "Product offering.\n  - Available values: virtual_server, gpu_server, k8s_vm, k8s_gpu_vm\n  - note: Use gpu_server for GPU instances",
				Optional:            true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description:         "VPC ID.\n  - example: cc976b621087484ea5fd527f4b78708b",
				MarkdownDescription: "VPC ID.\n  - example: cc976b621087484ea5fd527f4b78708b",
				Optional:            true,
			},
			common.ToSnakeCase("ServerTypeId"): schema.StringAttribute{
				Description:         "Server type ID.\n  - example: s1v1m2",
				MarkdownDescription: "Server type ID.\n  - example: s1v1m2",
				Optional:            true,
			},
			common.ToSnakeCase("AutoScalingGroupId"): schema.StringAttribute{
				Description:         "Auto Scaling Group ID.\n  - example: 52613bd852b04b39adcb15a8364d856d",
				MarkdownDescription: "Auto Scaling Group ID.\n  - example: 52613bd852b04b39adcb15a8364d856d",
				Optional:            true,
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "List of retrieved server IDs.",
				MarkdownDescription: "List of retrieved server IDs.",
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *virtualServerServerDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.VirtualServer
	d.clients = inst.Client
}

func (d *virtualServerServerDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualserver.ServerDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetServers(d.clients, state.Name, state.Ip, state.State, state.ProductCategory, state.ProductOffering,
		state.VpcId, state.ServerTypeId, state.AutoScalingGroupId, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	state.Ids = ids

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetServers(clients *client.SCPClient, name types.String, ip types.String, state types.String,
	productCategory types.String, productOffering types.String, vpcId types.String, serverTypeId types.String,
	autoScalingGroupId types.String, filters []filter.Filter) ([]types.String, error) {
	data, err := clients.VirtualServer.GetServerList(name, ip, state, productCategory, productOffering, vpcId, serverTypeId, autoScalingGroupId)
	if err != nil {
		return nil, err
	}

	contents := data.Servers
	filteredContents := data.Servers

	if len(filters) > 0 {
		filteredContents = filteredContents[:0]
		indices, err := filter.GetFilterIndices(contents, filters)
		if err != nil {
			return nil, err
		}

		for i, resource := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, resource)
			}
		}
		contents = filteredContents
	}

	var ids []types.String

	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
