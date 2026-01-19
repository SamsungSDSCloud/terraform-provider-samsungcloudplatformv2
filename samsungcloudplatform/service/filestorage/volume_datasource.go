package filestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/filestorage"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
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

func (d *fileStorageVolumeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = VolumeDataSourceSchema()
}
func VolumeDataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Computed: true,
				Description: "Account ID \n" +
					"  - example : 'rwww523320dfvwbbefefsdvwdadsfa24c' \n",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				Description: "Created At \n" +
					"  - example : '2024-07-30T04:54:33.219373' \n",
			},
			"encryption_enabled": schema.BoolAttribute{
				Computed: true,
				Description: "Volume Encryption Enabled \n" +
					"  - example : 'true'",
			},
			"endpoint_path": schema.StringAttribute{
				Computed: true,
				Description: "Volume Endpoint Path \n" +
					"  - example : 'xxx.xx.xxx.xxx'",
			},
			"file_unit_recovery_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Required: true,
				Description: "ID \n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Volume Name \n" +
					"  - example : 'my_volume' \n",
			},
			"path": schema.StringAttribute{
				Computed: true,
				Description: "Volume Mount Path \n" +
					"  - example : 'xxx.xx.xxx.xxx'",
			},
			"protocol": schema.StringAttribute{
				Computed: true,
				Description: "Protocol \n" +
					"  - example : 'NFS' \n",
			},
			"purpose": schema.StringAttribute{
				Computed: true,
				Description: "Volume Purpose \n" +
					"  - example : 'none' \n",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "Volume State \n" +
					"  - example : 'available' \n",
			},
			"type_id": schema.StringAttribute{
				Computed: true,
				Description: "Volume Type ID \n" +
					"  - example : 'jef22f67-ee83-4gg2-2ab6-3lf774ekfjdu' \n",
			},
			"type_name": schema.StringAttribute{
				Computed: true,
				Description: "Volume Type Name \n" +
					"  - example : 'HDD' \n",
			},
			"usage": schema.Int64Attribute{
				Computed: true,
				Description: "Volume Usage \n" +
					"  - example : '100000' \n",
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

	if volume.AccountId != "" {
		state.AccountId = types.StringValue(volume.AccountId)
	}
	if !volume.CreatedAt.IsZero() {
		state.CreatedAt = types.StringValue(volume.CreatedAt.Format(time.RFC3339))
	}
	if volume.Id != "" {
		state.Id = types.StringValue(volume.Id)
	}
	if volume.Name != "" {
		state.Name = types.StringValue(volume.Name)
	}
	if volume.Protocol != "" {
		state.Protocol = types.StringValue(volume.Protocol)
	}
	if volume.Purpose != "" {
		state.Purpose = types.StringValue(volume.Purpose)
	}
	if volume.State != "" {
		state.State = types.StringValue(volume.State)
	}
	if volume.TypeId != "" {
		state.TypeId = types.StringValue(volume.TypeId)
	}
	if volume.TypeName != "" {
		state.TypeName = types.StringValue(volume.TypeName)
	}

	state.EncryptionEnabled = types.BoolValue(volume.EncryptionEnabled)

	if volume.EndpointPath.Get() != nil {
		state.EndpointPath = types.StringValue(*volume.EndpointPath.Get())
	}

	state.FileUnitRecoveryEnabled = types.BoolValue(*volume.FileUnitRecoveryEnabled.Get())

	if volume.Path.Get() != nil {
		state.Path = types.StringValue(*volume.Path.Get())
	}
	if volume.Usage.Get() != nil {
		state.Usage = types.Int64Value(*volume.Usage.Get())
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}
