package scr

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scrImageDataSource{}
	_ datasource.DataSourceWithConfigure = &scrImageDataSource{}
)

func NewScrImageDataSource() datasource.DataSource {
	return &scrImageDataSource{}
}

type scrImageDataSource struct {
	config  *scpsdk.Configuration
	client  *scr.Client
	clients *client.SCPClient
}

func (d *scrImageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scr_image"
}

func (d *scrImageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Container Registry Image.",
		MarkdownDescription: "Get details of a Docker image in a Container Registry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Image ID",
				MarkdownDescription: "The unique identifier of the image.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Image name",
				MarkdownDescription: "The name of the image (digest).",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "Image state",
				MarkdownDescription: "The current state of the image.",
				Computed:            true,
			},
			"registry_id": schema.StringAttribute{
				Description:         "Registry ID",
				MarkdownDescription: "The ID of the registry containing the image.",
				Computed:            true,
			},
			"repository_id": schema.StringAttribute{
				Description:         "Repository ID",
				MarkdownDescription: "The ID of the repository containing the image.",
				Computed:            true,
			},
			"private_endpoint_url": schema.StringAttribute{
				Description:         "Private endpoint URL",
				MarkdownDescription: "The private endpoint URL for pulling the image.",
				Computed:            true,
			},
			"public_endpoint_url": schema.StringAttribute{
				Description:         "Public endpoint URL",
				MarkdownDescription: "The public endpoint URL for pulling the image.",
				Computed:            true,
			},
			"pull_count": schema.Int32Attribute{
				Description:         "Pull count",
				MarkdownDescription: "The number of times the image has been pulled.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				Description:         "Created at",
				MarkdownDescription: "The time the image was created.",
				Computed:            true,
			},
			"modified_at": schema.StringAttribute{
				Description:         "Modified at",
				MarkdownDescription: "The time the image was last modified.",
				Computed:            true,
			},
		},
	}
}

func (d *scrImageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *scrImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ImageDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	showResp, err := d.client.ShowImage(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Image",
			err.Error(),
		)
		return
	}

	image := showResp.Image
	state = ImageDataSource{
		Id:                 types.StringValue(image.Id),
		Name:               types.StringValue(image.Name),
		State:              types.StringValue(image.State),
		RegistryId:         types.StringValue(image.RegistryId),
		RepositoryId:       types.StringValue(image.RepositoryId),
		PrivateEndpointUrl: types.StringNull(),
		PublicEndpointUrl:  types.StringNull(),
		PullCount:          types.Int32Value(image.PullCount),
		CreatedAt:          types.StringValue(image.CreatedAt.Format(time.RFC3339)),
		ModifiedAt:         types.StringValue(image.ModifiedAt.Format(time.RFC3339)),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
