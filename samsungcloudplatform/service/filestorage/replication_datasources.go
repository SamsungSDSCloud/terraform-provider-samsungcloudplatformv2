package filestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &fileStorageReplicationDataSources{}
	_ datasource.DataSourceWithConfigure = &fileStorageReplicationDataSources{}
)

func NewFileStorageReplicationDataSources() datasource.DataSource {
	return &fileStorageReplicationDataSources{}
}

type fileStorageReplicationDataSources struct {
	config  *scpsdk.Configuration
	client  *filestorage.Client
	clients *client.SCPClient
}

func (d *fileStorageReplicationDataSources) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_filestorage_replications"
}

func (d *fileStorageReplicationDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Lists of Volume Replications.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("VolumeId"): schema.StringAttribute{
				Description: "Volume ID \n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
				Required: true,
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Volume Replication List",
			},
		},
	}
}

func (d *fileStorageReplicationDataSources) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

func (d *fileStorageReplicationDataSources) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var state filestorage.ReplicationResources

	diags := request.Config.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ids, err := GetVolumeReplications(d.clients, ctx, state)
	if err != nil {
		response.Diagnostics.AddError("Unable to Read Volume Replications", err.Error())
		return
	}

	state.ReplicationIds = ids

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

}

func GetVolumeReplications(client *client.SCPClient, ctx context.Context, state filestorage.ReplicationResources) ([]types.String, error) {
	data, err := client.FileStorage.GetVolumeReplicationList(ctx, state.VolumeId.ValueString())
	if err != nil {
		return nil, err
	}

	contents := data.Replications

	var ids []types.String

	for _, content := range contents {
		ids = append(ids, types.StringValue(content.ReplicationId))
	}

	return ids, nil
}
