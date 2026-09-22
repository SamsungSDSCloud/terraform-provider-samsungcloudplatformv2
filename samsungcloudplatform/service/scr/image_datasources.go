package scr

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scrImageDataSources{}
	_ datasource.DataSourceWithConfigure = &scrImageDataSources{}
)

func NewScrImageDataSources() datasource.DataSource {
	return &scrImageDataSources{}
}

type scrImageDataSources struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrImageDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_images"
}

func (d *scrImageDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List of Container Registry Images.",
		MarkdownDescription: "Get a list of Docker image IDs in a Container Registry repository.",
		Attributes: map[string]schema.Attribute{
			"repository_id": schema.StringAttribute{
				Description:         "Repository ID",
				MarkdownDescription: "The ID of the repository to list images from.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Image name filter",
				MarkdownDescription: "Filter images by name.",
				Optional:            true,
			},
			"sort": schema.StringAttribute{
				Description:         "Sort",
				MarkdownDescription: "Sort order, e.g. `name:asc`.",
				Optional:            true,
			},
			"page": schema.Int32Attribute{
				Description:         "Page number",
				MarkdownDescription: "Page number (0-based).",
				Optional:            true,
			},
			"size": schema.Int32Attribute{
				Description:         "Page size",
				MarkdownDescription: "Number of items per page.",
				Optional:            true,
			},
			"ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "Image ID list",
				MarkdownDescription: "List of image IDs.",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *scrImageDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *scrImageDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ImageDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := ""
	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		name = state.Name.ValueString()
	}
	sort := ""
	if !state.Sort.IsNull() && !state.Sort.IsUnknown() {
		sort = state.Sort.ValueString()
	}
	var page, size int32
	if !state.Page.IsNull() && !state.Page.IsUnknown() {
		page = state.Page.ValueInt32()
	}
	if !state.Size.IsNull() && !state.Size.IsUnknown() {
		size = state.Size.ValueInt32()
	}

	listResp, err := d.client.ListImages(ctx, state.RepositoryId.ValueString(), name, sort, page, size)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Images",
			err.Error(),
		)
		return
	}

	contents := listResp.Images
	filteredContents := contents

	if len(state.Filter) > 0 {
		filteredContents = filteredContents[:0]
		indices, err := filter.GetFilterIndices(contents, state.Filter)
		if err != nil {
			resp.Diagnostics.AddError("Filter Error", err.Error())
			return
		}

		for i, image := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, image)
			}
		}
		contents = filteredContents
	}

	var ids []types.String
	for _, image := range contents {
		ids = append(ids, types.StringValue(image.Id))
	}

	state.Ids = ids

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
