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
		Description: "Retrieves a list of images.\n\n" +
			"**GPU Image:**\n" +
			"- For GPU Server, use images with `scp_image_type` of `gpu_standard` or `gpu_custom`.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("ScpImageType"): schema.StringAttribute{
				Description: "SCP image type.\n" +
					"  - example: standard\n" +
					"  - Available values: standard, custom, gpu_standard, gpu_custom",
				MarkdownDescription: "SCP image type.\n" +
					"  - example: standard\n" +
					"  - Available values: standard, custom, gpu_standard, gpu_custom",
				Optional: true,
			},
			common.ToSnakeCase("ScpOriginalImageType"): schema.StringAttribute{
				Description:         "SCP original image type.\n  - example: standard",
				MarkdownDescription: "SCP original image type.\n  - example: standard",
				Optional:            true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description:         "Image name.\n  - example: ubuntu-22.04",
				MarkdownDescription: "Image name.\n  - example: ubuntu-22.04",
				Optional:            true,
			},
			common.ToSnakeCase("OsDistro"): schema.StringAttribute{
				Description: "OS distribution.\n" +
					"  - example: ubuntu\n" +
					"  - Available values: alma, centos, rhel, rocky, ubuntu, windows, oracle",
				MarkdownDescription: "OS distribution.\n" +
					"  - example: ubuntu\n" +
					"  - Available values: alma, centos, rhel, rocky, ubuntu, windows, oracle",
				Optional: true,
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Description: "Image status.\n" +
					"  - example: active\n",
				MarkdownDescription: "Image status.\n" +
					"  - example: active\n",
				Optional: true,
			},
			common.ToSnakeCase("Visibility"): schema.StringAttribute{
				Description: "Image visibility.\n" +
					"  - example: private\n" +
					"  - Available values: shared, private",
				MarkdownDescription: "Image visibility.\n" +
					"  - example: private\n" +
					"  - Available values: , shared, private",
				Optional: true,
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "List of image IDs.",
				MarkdownDescription: "List of image IDs.",
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
