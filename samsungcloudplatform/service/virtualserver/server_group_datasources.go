package virtualserver

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualServerServerGroupDataSources{}
	_ datasource.DataSourceWithConfigure = &virtualServerServerGroupDataSources{}
)

func NewVirtualServerServerGroupDataSources() datasource.DataSource {
	return &virtualServerServerGroupDataSources{}
}

type virtualServerServerGroupDataSources struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (d *virtualServerServerGroupDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_server_groups"
}

func (d *virtualServerServerGroupDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of server groups.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Server Group ID List",
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *virtualServerServerGroupDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *virtualServerServerGroupDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualserver.ServerGroupDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetServerGroups(d.clients, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Server Groups",
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

func GetServerGroups(clients *client.SCPClient, filters []filter.Filter) ([]types.String, error) {
	data, err := clients.VirtualServer.GetServerGroupList()
	if err != nil {
		return nil, err
	}

	contents := data.ServerGroups
	filteredContents := data.ServerGroups

	if len(filters) > 0 {
		filteredContents = filteredContents[:0]
		indices := filter.GetFilterIndices(contents, filters) // 필터링된 컨텐츠의 Index 정보 리턴

		for i, resource := range contents {
			if common.Contains(indices, i) { // Index 정보 기준으로 필터링 진행
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
