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
	_ datasource.DataSource              = &virtualServerImageDataSources{}
	_ datasource.DataSourceWithConfigure = &virtualServerImageDataSources{}
)

func NewVirtualServerImageDataSources() datasource.DataSource {
	return &virtualServerImageDataSources{}
}

type virtualServerImageDataSources struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (d *virtualServerImageDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_images"
}

func (d *virtualServerImageDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of images.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("ScpImageType"): schema.StringAttribute{
				Description: "SCP Image type",
				Optional:    true,
			},
			common.ToSnakeCase("ScpOriginalImageType"): schema.StringAttribute{
				Description: "SCP Original Image type",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("OsDistro"): schema.StringAttribute{
				Description: "OS Distro",
				Optional:    true,
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Description: "Status",
				Optional:    true,
			},
			common.ToSnakeCase("Visibility"): schema.StringAttribute{
				Description: "Visibility",
				Optional:    true,
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Image ID List",
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *virtualServerImageDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *virtualServerImageDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualserver.ImageDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetImages(d.clients, state.ScpImageType, state.ScpOriginalImageType, state.Name, state.OsDistro, state.Status, state.Visibility, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Images",
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

func GetImages(clients *client.SCPClient, scpImageType types.String, scpOriginalImageType types.String, name types.String,
	osDistro types.String, status types.String, visibility types.String, filters []filter.Filter) ([]types.String, error) {
	data, err := clients.VirtualServer.GetImageList(scpImageType, scpOriginalImageType, name, osDistro, status, visibility)
	if err != nil {
		return nil, err
	}

	contents := data.Images
	filteredContents := data.Images

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
