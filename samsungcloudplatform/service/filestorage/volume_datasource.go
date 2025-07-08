package filestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &fileStorageVolumeDataSource{}
	_ datasource.DataSourceWithConfigure = &fileStorageVolumeDataSource{}
)

func NewFileStorageVolumeDataSource() datasource.DataSource {
	return &fileStorageVolumeDataSource{}
}

type fileStorageVolumeDataSource struct {
	config  *scpsdk.Configuration
	client  *filestorage.Client
	clients *client.SCPClient
}

func (d *fileStorageVolumeDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_filestorage_volume"
}

func (d *fileStorageVolumeDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Show Volume",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "ID",
				Optional:    true,
			},
			common.ToSnakeCase("Volume"): schema.SingleNestedAttribute{
				Description: "Volume",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID \n" +
							"  - example : 'rwww523320dfvwbbefefsdvwdadsfa24c' \n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created At \n" +
							"  - example : '2024-07-30T04:54:33.219373' \n",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Account ID \n" +
							"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Volume Name \n" +
							"  - example : 'my_volume' \n",
						Computed: true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "Protocol \n" +
							"  - example : 'NFS' \n",
						Computed: true,
					},
					common.ToSnakeCase("Purpose"): schema.StringAttribute{
						Description: "Volume Purpose \n" +
							"  - example : 'none' \n",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "Volume State \n" +
							"  - example : 'available' \n",
						Computed: true,
					},
					common.ToSnakeCase("TypeId"): schema.StringAttribute{
						Description: "Volume Type ID \n" +
							"  - example : 'jef22f67-ee83-4gg2-2ab6-3lf774ekfjdu' \n",
						Computed: true,
					},
					common.ToSnakeCase("TypeName"): schema.StringAttribute{
						Description: "Volume Type Name \n" +
							"  - example : 'HDD' \n",
						Computed: true,
					},
				},
			},
		},
	}
}

func (d *fileStorageVolumeDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	inst, ok := request.ProviderData.(client.Instance)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Plase report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	d.client = inst.Client.FileStorage
	d.clients = inst.Client
}

func (d *fileStorageVolumeDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var state filestorage.VolumeDataSource
	diags := request.Config.Get(ctx, &state)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	volume, err := d.client.GetVolume(ctx, state.Id.ValueString())

	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError("Error Reading Volume",
			"Could not read Volume Id "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail)
	}
	volumeModel := filestorage.Volume{
		AccountId: types.StringValue(volume.AccountId),
		CreatedAt: types.StringValue(volume.CreatedAt.Format(time.RFC3339)),
		Id:        types.StringValue(volume.Id),
		Name:      types.StringValue(volume.Name),
		Protocol:  types.StringValue(volume.Protocol),
		Purpose:   types.StringValue(volume.Purpose),
		State:     types.StringValue(volume.State),
		TypeId:    types.StringValue(volume.TypeId),
		TypeName:  types.StringValue(volume.TypeName),
	}

	volumeObjectValue, _ := types.ObjectValueFrom(ctx, volumeModel.AttributeTypes(), volumeModel)
	state.Volume = volumeObjectValue

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}
