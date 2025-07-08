package virtualserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualServerKeypairDataSources{}
	_ datasource.DataSourceWithConfigure = &virtualServerKeypairDataSources{}
)

func NewVirtualServerKeypairDataSources() datasource.DataSource {
	return &virtualServerKeypairDataSources{}
}

type virtualServerKeypairDataSources struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (d *virtualServerKeypairDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_keypairs"
}

func (d *virtualServerKeypairDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of keypair.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Names"): schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Keypair Name List",
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *virtualServerKeypairDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *virtualServerKeypairDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualserver.KeypairDataSourceNames

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	names, err := GetKeypairs(d.clients, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Keypairs",
			err.Error(),
		)
		return
	}

	state.Names = names

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetKeypairs(clients *client.SCPClient, filters []filter.Filter) ([]types.String, error) {
	data, err := clients.VirtualServer.GetKeypairList()
	if err != nil {
		return nil, err
	}

	contents := data.Keypairs
	filteredContents := data.Keypairs

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

	var names []types.String

	for _, content := range contents {
		names = append(names, types.StringValue(content.Name))
	}

	return names, nil
}
