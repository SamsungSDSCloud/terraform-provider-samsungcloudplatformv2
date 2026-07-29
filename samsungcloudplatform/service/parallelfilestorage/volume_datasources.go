package parallelfilestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client/parallelfilestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &parallelFileStorageVolumeDataSources{}
	_ datasource.DataSourceWithConfigure = &parallelFileStorageVolumeDataSources{}
)

func NewParallelFileStorageVolumeDataSources() datasource.DataSource {
	return &parallelFileStorageVolumeDataSources{}
}

type parallelFileStorageVolumeDataSources struct {
	config  *scpsdk.Configuration
	client  *parallelfilestorage.Client
	clients *client.SCPClient
}

func (d *parallelFileStorageVolumeDataSources) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_parallel_filestorage_volumes"
}

func (d *parallelFileStorageVolumeDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Lists of Volumes.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Offset"): schema.Int32Attribute{
				Description: "Offset \n" +
					"  - example : 0 \n" +
					"  - maximum: 10000  \n" +
					"  - minimum: 0  \n",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "Limit \n" +
					"  - example : 20 \n",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : 'created_at:asc' \n",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Volume Name \n" +
					"  - example : 'my_volume' \n" +
					"  - maxLength: 21  \n" +
					"  - minLength: 3  \n" +
					"  - pattern: `^[a-z]([a-z0-9_]){2,20}$` \n",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 21),
				},
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Volume ID List",
			},
		},
	}
}

func (d *parallelFileStorageVolumeDataSources) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

	d.client = inst.Client.ParallelFileStorage
	d.clients = inst.Client
}

func (d *parallelFileStorageVolumeDataSources) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var state parallelfilestorage.VolumeDataSourceIds

	diags := request.Config.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ids, err := GetVolumeIds(d.clients, ctx, state)
	if err != nil {
		response.Diagnostics.AddError("Unable to Read Volumes", err.Error())
		return
	}

	state.Ids = ids

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

}

func GetVolumeIds(client *client.SCPClient, ctx context.Context, state parallelfilestorage.VolumeDataSourceIds) ([]types.String, error) {
	data, err := client.ParallelFileStorage.GetVolumeList(ctx, state)
	if err != nil {
		return nil, err
	}

	contents := data.ParallelFilestorages

	var ids []types.String
	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
