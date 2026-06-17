package virtualserver

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/filter"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &virtualServerVolumeDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualServerVolumeDataSource{}
)

// NewComputeVolumeDataSource is a helper function to simplify the provider implementation.
func NewVirtualServerVolumeDataSource() datasource.DataSource {
	return &virtualServerVolumeDataSource{}
}

// virtualServerVolumeDataSource is the data source implementation.
type virtualServerVolumeDataSource struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *virtualServerVolumeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_volume"
}

// Schema defines the schema for the data source.
func (d *virtualServerVolumeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves volume information.",
		MarkdownDescription: "Retrieves block storage volume information for virtual servers.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Volume ID.\n" +
					"  - example: 550e8400-e29b-41d4-a716-446655440000",
				MarkdownDescription: "Volume ID.\n" +
					"  - example: 550e8400-e29b-41d4-a716-446655440000",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Volume name.\n" +
					"  - example: my-volume\n" +
					"  - minLength: 1\n" +
					"  - maxLength: 255",
				MarkdownDescription: "Volume name.\n" +
					"  - example: my-volume\n" +
					"  - minLength: 1\n" +
					"  - maxLength: 255",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "Volume state.\n" +
					"  - Available values: available, reserved, attaching, detaching, in-use, awaiting-transfer, error, etc.",
				MarkdownDescription: "Volume state.\n" +
					"  - Available values: available, reserved, attaching, detaching, in-use, awaiting-transfer, error, etc.",
				Optional: true,
			},
			common.ToSnakeCase("Bootable"): schema.BoolAttribute{
				Description:         "Bootable flag.",
				MarkdownDescription: "Whether the volume is bootable.\n  - example: false",
				Optional:            true,
			},
			common.ToSnakeCase("Volume"): schema.SingleNestedAttribute{
				Description:         "Volume details.",
				MarkdownDescription: "Volume details including size, type, state, and attached servers.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description:         "Volume ID.",
						MarkdownDescription: "Volume ID.",
						Computed:            true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description:         "Volume name.",
						MarkdownDescription: "Volume name.",
						Computed:            true,
					},
					common.ToSnakeCase("Size"): schema.Int32Attribute{
						Description:         "Volume size (GB).",
						MarkdownDescription: "Volume size in GB.",
						Computed:            true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description:         "Volume state.",
						MarkdownDescription: "Volume state.",
						Computed:            true,
					},
					common.ToSnakeCase("UserId"): schema.StringAttribute{
						Description:         "User ID.",
						MarkdownDescription: "User ID.",
						Computed:            true,
					},
					common.ToSnakeCase("VolumeType"): schema.StringAttribute{
						Description:         "Volume type.",
						MarkdownDescription: "Volume type.",
						Computed:            true,
					},
					common.ToSnakeCase("Encrypted"): schema.BoolAttribute{
						Description:         "Encryption flag.",
						MarkdownDescription: "Whether the volume is encrypted.",
						Computed:            true,
					},
					common.ToSnakeCase("Bootable"): schema.BoolAttribute{
						Description:         "Bootable flag.",
						MarkdownDescription: "Whether the volume is bootable.",
						Computed:            true,
					},
					common.ToSnakeCase("Multiattach"): schema.BoolAttribute{
						Description:         "Multi-attach flag.",
						MarkdownDescription: "Whether the volume can be attached to multiple servers.",
						Computed:            true,
					},
					common.ToSnakeCase("Servers"): schema.ListNestedAttribute{
						Description:         "List of attached servers.",
						MarkdownDescription: "List of attached servers.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description:         "Server ID.",
									MarkdownDescription: "Server ID.",
									Computed:            true,
								},
							},
						},
					},
					common.ToSnakeCase("MaxIops"): schema.Int32Attribute{
						Description:         "Maximum IOPS per second.",
						MarkdownDescription: "Maximum IOPS per second.",
						Computed:            true,
					},
					common.ToSnakeCase("MaxThroughput"): schema.Int32Attribute{
						Description:         "Maximum throughput per second (MB/s).",
						MarkdownDescription: "Maximum throughput per second (MB/s).",
						Computed:            true,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *virtualServerVolumeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

// Read refreshes the Terraform state with the latest data.
func (d *virtualServerVolumeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualserver.VolumeDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetVolumes(d.clients, state.Name, state.State, state.Bootable, state.Filter)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Volumes",
			err.Error(),
		)
		return
	}
	if len(ids) > 0 || !state.Id.IsNull() {
		id := virtualserverutil.SetResourceIdentifier(state.Id, ids, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		volume, err := d.client.GetVolume(ctx, id.ValueString())

		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Reading Server",
				"Could not read Server ID "+id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}
		var servers []virtualserver.VolumeServer
		for _, server := range volume.Servers {
			serverState := virtualserver.VolumeServer{
				Id: types.StringValue(server.Id),
			}
			servers = append(servers, serverState)
		}
		volumeModel := virtualserver.Volume{
			Id:            types.StringValue(volume.Id),
			Name:          types.StringPointerValue(volume.Name.Get()),
			Size:          types.Int32Value(volume.Size),
			State:         types.StringValue(string(volume.State)),
			UserId:        types.StringValue(volume.UserId),
			VolumeType:    types.StringValue(volume.VolumeType),
			Encrypted:     types.BoolValue(volume.Encrypted),
			Bootable:      types.BoolValue(volume.Bootable),
			Multiattach:   types.BoolValue(volume.Multiattach),
			Servers:       servers,
			MaxIops:       types.Int32PointerValue(volume.MaxIops.Get()),
			MaxThroughput: types.Int32PointerValue(volume.MaxThroughput.Get()),
		}
		volumeObjectValue, _ := types.ObjectValueFrom(ctx, volumeModel.AttributeTypes(), volumeModel)
		state.Volume = volumeObjectValue
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
