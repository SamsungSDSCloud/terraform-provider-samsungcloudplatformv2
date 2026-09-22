package scr

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scrContainerRegistryDataSources{}
	_ datasource.DataSourceWithConfigure = &scrContainerRegistryDataSources{}
)

func NewScrContainerRegistryDataSources() datasource.DataSource {
	return &scrContainerRegistryDataSources{}
}

type scrContainerRegistryDataSources struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrContainerRegistryDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_container_registries"
}

func (d *scrContainerRegistryDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List of Container Registries.",
		MarkdownDescription: "Get a list of Container Registry IDs.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "Registry name",
				MarkdownDescription: "Filter by registry name.",
				Optional:            true,
			},
			"ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "Registry ID list",
				MarkdownDescription: "List of registry IDs.",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *scrContainerRegistryDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Scr
	d.clients = inst.Client
}

func (d *scrContainerRegistryDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ContainerRegistryDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := ""
	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		name = state.Name.ValueString()
	}

	listResp, err := d.client.ListRegistries(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Registries",
			err.Error(),
		)
		return
	}

	contents := listResp.Registries

	if len(state.Filter) > 0 {
		filteredContents := contents[:0]
		indices, err := filter.GetFilterIndices(contents, state.Filter)
		if err != nil {
			resp.Diagnostics.AddError("Filter Error", err.Error())
			return
		}

		for i, registry := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, registry)
			}
		}
		contents = filteredContents
	}

	var ids []types.String
	for _, registry := range contents {
		ids = append(ids, types.StringValue(registry.Id))
	}

	state.Ids = ids

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
