package filestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &fileStorageSnapshotScheduleDataSource{}
	_ datasource.DataSourceWithConfigure = &fileStorageSnapshotScheduleDataSource{}
)

func NewFileStorageSnapshotScheduleDataSource() datasource.DataSource {
	return &fileStorageSnapshotScheduleDataSource{}
}

type fileStorageSnapshotScheduleDataSource struct {
	config  *scpsdk.Configuration
	client  *filestorage.Client
	clients *client.SCPClient
}

func (d *fileStorageSnapshotScheduleDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_filestorage_snapshot_schedule"
}

func (d *fileStorageSnapshotScheduleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Lists of SnapshotSchedules.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("VolumeId"): schema.StringAttribute{
				Description: "Volume ID \n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
				Required: true,
			},
			common.ToSnakeCase("SnapshotPolicyEnabled"): schema.BoolAttribute{
				Description: "Snapshot Policy Enabled \n" +
					"  - example : 'true' \n",
				Computed: true,
			},
			common.ToSnakeCase("SnapshotRetentionCount"): schema.Int32Attribute{
				Description: "Snapshot Retention Count (If not entered, 10 will be applied.) \n" +
					"  - example : 1 \n" +
					"  - maximum : 128 \n" +
					"  - minimum : 1  \n",
				Computed: true,
			},
			common.ToSnakeCase("SnapshotSchedules"): schema.ListNestedAttribute{
				Description: "List of SnapshotSchedules.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("DayOfWeek"): schema.StringAttribute{
							Description: "Day Of Week \n" +
								"  - example : 'MON' \n" +
								"  - pattern: '^(SUN|MON|TUE|WED|THU|FRI|SAT)$' \n",
							Computed: true,
						},
						common.ToSnakeCase("Frequency"): schema.StringAttribute{
							Description: "Frequency \n" +
								"  - example : 'DAILY' \n" +
								"  - pattern: '^(WEEKLY|DAILY)$' \n",
							Computed: true,
						},
						common.ToSnakeCase("Hour"): schema.StringAttribute{
							Description: "Hour \n" +
								"  - example : '0' \n" +
								"  - maximum : 23 \n" +
								"  - minimum : 0  \n" +
								"  - pattern: '^([0-9]|1[0-9]|2[0-3])$' \n",
							Computed: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{Description: "ID", Computed: true},
					},
				},
			},
		},
	}
}

func (d *fileStorageSnapshotScheduleDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

func (d *fileStorageSnapshotScheduleDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var state filestorage.SnapshotScheduleDataSource

	diags := request.Config.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetSnapshotScheduleList(ctx, state.VolumeId.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Unable to Read SnapshotSchedules", err.Error())
		return
	}

	state.SnapshotRetentionCount = types.Int32PointerValue(data.SnapshotRetentionCount)
	state.SnapshotPolicyEnabled = types.BoolPointerValue(data.SnapshotPolicyEnabled.Get())

	for _, snapshotScheduleList := range data.SnapshotSchedule {
		var snapshotScheduleState = filestorage.SnapshotSchedule{
			DayOfWeek: types.StringPointerValue(snapshotScheduleList.DayOfWeek.Get()),
			Frequency: types.StringValue(snapshotScheduleList.Frequency),
			Hour:      types.StringValue(snapshotScheduleList.Hour),
			Id:        types.StringPointerValue(snapshotScheduleList.Id.Get()),
		}
		state.SnapshotSchedules = append(state.SnapshotSchedules, snapshotScheduleState)
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

}
