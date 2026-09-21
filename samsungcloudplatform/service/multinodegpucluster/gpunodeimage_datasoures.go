package multinodegpucluster

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	multinodegpuclusterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/multinodegpucluster"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GpunodeImageDataSources{}
	_ datasource.DataSourceWithConfigure = &GpunodeImageDataSources{}
)

// NewGpunodeImageDataSources is a helper function to simplify the provider implementation.
func NewGpunodeImageDataSources() datasource.DataSource {
	return &GpunodeImageDataSources{}
}

// GpunodeImageDataSources is the data source implementation.
type GpunodeImageDataSources struct {
	config  *scpsdk.Configuration
	client  *multinodegpuclusterClient.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *GpunodeImageDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multinodegpucluster_gpunode_images"
}

// Configure adds the provider configured client to the data source.
func (d *GpunodeImageDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Mngc
	d.clients = inst.Client
}

// Schema defines the schema for the data source.
func (d *GpunodeImageDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GpuNodeImagesDataSourcesSchema()
}

func GpuNodeImagesDataSourcesSchema() schema.Schema {
	return schema.Schema{
		Description: "List of GPU Node Image",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Region ID\n  - example: kr-west1",
				MarkdownDescription: "Region ID\n  - example: kr-west1",
			},
			"images": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "GPU Node Image List",
				MarkdownDescription: "GPU Node Image List",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Image ID\n  - example: 10a599e031e749b7b260868f441e862b",
							MarkdownDescription: "Image ID\n  - example: 10a599e031e749b7b260868f441e862b",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Image name\n  - example: Ubuntu 24.04 GPU (B300/580.105.08) for BM",
							MarkdownDescription: "Image name\n  - example: Ubuntu 24.04 GPU (B300/580.105.08) for BM",
						},
						"os_distro": schema.StringAttribute{
							Computed:            true,
							Description:         "OS Distro\n  - example: UBUNTU",
							MarkdownDescription: "OS Distro\n  - example: UBUNTU",
						},
						"scp_os_version": schema.StringAttribute{
							Computed:            true,
							Description:         "SCP OS version\n  - example: 24.04",
							MarkdownDescription: "SCP OS version\n  - example: 24.04",
						},
						"scp_image_type": schema.StringAttribute{
							Computed:            true,
							Description:         "SCP image type\n  - example: STANDARD",
							MarkdownDescription: "SCP image type\n  - example: STANDARD",
						},
						"priority": schema.StringAttribute{
							Computed:            true,
							Description:         "GPU Node Image priority\n  - example: 1",
							MarkdownDescription: "GPU Node Image priority\n  - example: 1",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Created at\n  - example: 2024-05-01T00:00:00Z",
							MarkdownDescription: "Created at\n  - example: 2024-05-01T00:00:00Z",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *GpunodeImageDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state multinodegpuclusterClient.GpuNodeImageList

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetImageList(ctx, state.RegionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read GPU Node Images",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, image := range data.Images {
		imageState := multinodegpuclusterClient.GpuNodeImage{
			CreatedAt:    types.StringValue(image.CreatedAt.Format(time.RFC3339)),
			Id:           types.StringValue(image.Id),
			Name:         types.StringValue(image.Name),
			OsDistro:     types.StringValue(image.OsDistro),
			Priority:     types.StringValue(image.Priority),
			ScpImageType: types.StringValue(image.ScpImageType),
			ScpOsVersion: types.StringValue(image.ScpOsVersion),
		}
		state.Images = append(state.Images, imageState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
