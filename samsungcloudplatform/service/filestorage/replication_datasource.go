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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &fileStorageReplicationDataSource{}
	_ datasource.DataSourceWithConfigure = &fileStorageReplicationDataSource{}
)

func NewFileStorageReplicationDataSource() datasource.DataSource {
	return &fileStorageReplicationDataSource{}
}

type fileStorageReplicationDataSource struct {
	config  *scpsdk.Configuration
	client  *filestorage.Client
	clients *client.SCPClient
}

func (f *fileStorageReplicationDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

	f.client = inst.Client.FileStorage
	f.clients = inst.Client
}

func (f *fileStorageReplicationDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_filestorage_replication"
}

func (f *fileStorageReplicationDataSource) Schema(_ context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Show Replication",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID",
				Required:    true,
			},
			"volume_id": schema.StringAttribute{
				Description: "Volume ID \n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
				Required: true,
			},
			common.ToSnakeCase("Replication"): schema.SingleNestedAttribute{
				Description: "Replication.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"replication_frequency": schema.StringAttribute{
						Description: "Replication Frequency \n" +
							"  - example : '5min' \n",
						Computed: true,
					},
					"replication_id": schema.StringAttribute{
						Description: "Replication ID \n" +
							"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
						Computed: true,
					},
					"replication_policy": schema.StringAttribute{
						Description: "Replication Policy \n" +
							"  - example : 'use' \n",
						Computed: true,
					},
					"replication_status": schema.StringAttribute{
						Description: "Replication Status \n" +
							"  - example : 'creating' \n",
						Computed: true,
					},
					"replication_volume_access_level": schema.StringAttribute{
						Description: "Target Access Level \n" +
							"  - example : 'ro' \n",
						Computed: true,
					},
					"replication_volume_id": schema.StringAttribute{
						Description: "Target Volume ID \n" +
							"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
						Computed: true,
					},
					"replication_volume_name": schema.StringAttribute{
						Description: "Target Volume Name \n" +
							"  - example : 'my_volume' \n",
						Computed: true,
					},
					"replication_volume_region": schema.StringAttribute{
						Description: "Target Volume Region \n" +
							"  - example : 'kr-west1' \n",
						Computed: true,
					},
					"source_volume_access_level": schema.StringAttribute{
						Description: "Source Access Level \n" +
							"  - example : 'ro' \n",
						Computed: true,
					},
					"source_volume_id": schema.StringAttribute{
						Description: "Source Volume ID \n" +
							"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
						Computed: true,
					},
					"source_volume_name": schema.StringAttribute{
						Description: "Source Volume Name \n" +
							"  - example : 'my_volume' \n",
						Computed: true,
					},
					"source_volume_region": schema.StringAttribute{
						Description: "Source Volume Region \n" +
							"  - example : 'kr-west1' \n",
						Computed: true,
					},
				},
			},
		},
	}
}

func (f *fileStorageReplicationDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var state filestorage.ReplicationDataSource

	diags := request.Config.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	data, err := f.client.GetVolumeReplication(ctx, state.Id.ValueString(), state.VolumeId.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Unable to Read Replication Policy", err.Error())
		return
	}

	var replicationState = filestorage.Replication{
		ReplicationFrequency:         types.StringValue(data.ReplicationFrequency),
		ReplicationId:                types.StringValue(data.ReplicationId),
		ReplicationPolicy:            types.StringValue(data.ReplicationPolicy),
		ReplicationStatus:            types.StringValue(data.ReplicationStatus),
		ReplicationVolumeAccessLevel: types.StringValue(data.ReplicationVolumeAccessLevel),
		ReplicationVolumeId:          types.StringValue(data.ReplicationVolumeId),
		ReplicationVolumeName:        types.StringValue(data.ReplicationVolumeName),
		ReplicationVolumeRegion:      types.StringValue(data.ReplicationVolumeRegion),
		SourceVolumeAccessLevel:      types.StringValue(data.SourceVolumeAccessLevel),
		SourceVolumeId:               types.StringValue(data.SourceVolumeId),
		SourceVolumeName:             types.StringValue(data.SourceVolumeName),
		SourceVolumeRegion:           types.StringValue(data.SourceVolumeRegion),
	}
	replicationObjectValue, _ := types.ObjectValueFrom(ctx, replicationState.AttributeTypes(), replicationState)

	state.Replication = replicationObjectValue

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}
}
