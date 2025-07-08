package filestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &fileStorageSnapshotScheduleResource{}
	_ resource.ResourceWithConfigure = &fileStorageSnapshotScheduleResource{}
)

func NewFileStorageSnapshotScheduleResource() resource.Resource {
	return &fileStorageSnapshotScheduleResource{}
}

type fileStorageSnapshotScheduleResource struct {
	config  *scpsdk.Configuration
	client  *filestorage.Client
	clients *client.SCPClient
}

func (r *fileStorageSnapshotScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filestorage_snapshot_schedule"
}

func (r *fileStorageSnapshotScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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
				Optional: true,
			},
			common.ToSnakeCase("SnapshotSchedule"): schema.SingleNestedAttribute{
				Description: "Snapshot Schedule",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("DayOfWeek"): schema.StringAttribute{
						Description: "Day Of Week \n" +
							"  - example : 'MON' \n" +
							"  - pattern: '^(SUN|MON|TUE|WED|THU|FRI|SAT)$' \n",
						Optional: true,
					},
					common.ToSnakeCase("Frequency"): schema.StringAttribute{
						Description: "Frequency \n" +
							"  - example : 'DAILY' \n" +
							"  - pattern: '^(WEEKLY|DAILY)$' \n",
						Optional: true,
					},
					common.ToSnakeCase("Hour"): schema.StringAttribute{
						Description: "Hour \n" +
							"  - example : '0' \n" +
							"  - maximum : 23 \n" +
							"  - minimum : 0  \n" +
							"  - pattern: '^([0-9]|1[0-9]|2[0-3])$' \n",
						Optional: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{Description: "ID", Computed: true},
				},
			},
		},
	}
}

func (r *fileStorageSnapshotScheduleResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	inst, ok := request.ProviderData.(client.Instance)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	r.client = inst.Client.FileStorage
	r.clients = inst.Client
}

func (r *fileStorageSnapshotScheduleResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	// Values from plan
	var plan filestorage.SnapshotScheduleResource
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Exist Check Schedule
	getData, err := r.client.GetSnapshotScheduleList(ctx, plan.VolumeId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Reading Snapshot Schedule",
			"Could not read Snapshot Schedule"+getData.GetVolumeId()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if getData.GetSnapshotPolicyEnabled() {
		var SnapshotSheduleId = ""
		var snapshotScheduleState = filestorage.SnapshotSchedule{
			DayOfWeek: types.StringPointerValue(nil),
			Frequency: types.StringValue(""),
			Hour:      types.StringValue(""),
			Id:        types.StringPointerValue(nil),
		}

		for _, snapshotScheduleList := range getData.SnapshotSchedule {
			snapshotScheduleState = filestorage.SnapshotSchedule{
				DayOfWeek: types.StringPointerValue(snapshotScheduleList.DayOfWeek.Get()),
				Frequency: types.StringValue(snapshotScheduleList.Frequency),
				Hour:      types.StringValue(snapshotScheduleList.Hour),
				Id:        types.StringPointerValue(snapshotScheduleList.Id.Get()),
			}
		}

		SnapshotSheduleId = snapshotScheduleState.Id.ValueString()

		// Update Schedule
		_, err := r.client.UpdateSnapshotSchedule(ctx, SnapshotSheduleId, plan)
		if err != nil {
			detail := client.GetDetailFromError(err)
			response.Diagnostics.AddError(
				"Error updating SnapshotSchedule",
				"Could not update SnapshotSchedule, unexpected error: \nReason: "+detail,
			)
			return
		}
	} else {
		// Create Schedule
		_, err := r.client.CreateSnapshotSchedule(ctx, plan)
		if err != nil {
			detail := client.GetDetailFromError(err)
			response.Diagnostics.AddError(
				"Error creating SnapshotSchedule",
				"Could not create SnapshotSchedule, unexpected error: \nReason: "+detail,
			)
			return
		}
	}

	// Read Schedule
	getSchedule, err := r.client.GetSnapshotScheduleList(ctx, plan.VolumeId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Reading Snapshot Schedule",
			"Could not read Snapshot Schedule "+plan.VolumeId.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	snapshotSchedule := getSchedule
	// Map response body to schema and populate Computed attribute values
	plan.VolumeId = types.StringValue(snapshotSchedule.VolumeId)
	plan.SnapshotRetentionCount = types.Int32PointerValue(snapshotSchedule.SnapshotRetentionCount)
	plan.SnapshotPolicyEnabled = types.BoolPointerValue(snapshotSchedule.SnapshotPolicyEnabled.Get())

	var snapshotScheduleState = filestorage.SnapshotSchedule{
		DayOfWeek: types.StringPointerValue(nil),
		Frequency: types.StringValue(""),
		Hour:      types.StringValue(""),
		Id:        types.StringPointerValue(nil),
	}

	for _, snapshotScheduleList := range snapshotSchedule.SnapshotSchedule {
		snapshotScheduleState = filestorage.SnapshotSchedule{
			DayOfWeek: types.StringPointerValue(snapshotScheduleList.DayOfWeek.Get()),
			Frequency: types.StringValue(snapshotScheduleList.Frequency),
			Hour:      types.StringValue(snapshotScheduleList.Hour),
			Id:        types.StringPointerValue(snapshotScheduleList.Id.Get()),
		}

	}

	plan.SnapshotSchedule = snapshotScheduleState

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

}

func (r *fileStorageSnapshotScheduleResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state filestorage.SnapshotScheduleResource

	diags := request.State.Get(ctx, &state)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetSnapshotScheduleList(ctx, state.VolumeId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Reading Snapshot Schedule",
			"Could not read Volume ID "+state.VolumeId.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	state.SnapshotRetentionCount = types.Int32PointerValue(data.SnapshotRetentionCount)
	state.SnapshotPolicyEnabled = types.BoolPointerValue(data.SnapshotPolicyEnabled.Get())

	var snapshotScheduleState = filestorage.SnapshotSchedule{
		DayOfWeek: types.StringPointerValue(nil),
		Frequency: types.StringValue(""),
		Hour:      types.StringValue(""),
		Id:        types.StringPointerValue(nil),
	}

	for _, snapshotScheduleList := range data.SnapshotSchedule {
		snapshotScheduleState = filestorage.SnapshotSchedule{
			DayOfWeek: types.StringPointerValue(snapshotScheduleList.DayOfWeek.Get()),
			Frequency: types.StringValue(snapshotScheduleList.Frequency),
			Hour:      types.StringValue(snapshotScheduleList.Hour),
			Id:        types.StringPointerValue(snapshotScheduleList.Id.Get()),
		}

	}

	state.SnapshotSchedule = snapshotScheduleState
	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

}

func (r *fileStorageSnapshotScheduleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// Values from plan
	var plan filestorage.SnapshotScheduleResource
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Exist Check Schedule
	getData, err := r.client.GetSnapshotScheduleList(ctx, plan.VolumeId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Reading Snapshot Schedule",
			"Could not read Snapshot Schedule"+getData.GetVolumeId()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var SnapshotSheduleId = ""

	var snapshotScheduleState = filestorage.SnapshotSchedule{
		DayOfWeek: types.StringPointerValue(nil),
		Frequency: types.StringValue(""),
		Hour:      types.StringValue(""),
		Id:        types.StringPointerValue(nil),
	}

	for _, snapshotScheduleList := range getData.SnapshotSchedule {
		snapshotScheduleState = filestorage.SnapshotSchedule{
			DayOfWeek: types.StringPointerValue(snapshotScheduleList.DayOfWeek.Get()),
			Frequency: types.StringValue(snapshotScheduleList.Frequency),
			Hour:      types.StringValue(snapshotScheduleList.Hour),
			Id:        types.StringPointerValue(snapshotScheduleList.Id.Get()),
		}
	}

	SnapshotSheduleId = snapshotScheduleState.Id.ValueString()

	// Update Schedule
	_, err = r.client.UpdateSnapshotSchedule(ctx, SnapshotSheduleId, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error updating SnapshotSchedule",
			"Could not update SnapshotSchedule, unexpected error: \nReason: "+detail,
		)
		return
	}

	// Read Schedule
	getSchedule, err := r.client.GetSnapshotScheduleList(ctx, plan.VolumeId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Reading Snapshot Schedule",
			"Could not read Snapshot Schedule "+plan.VolumeId.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	snapshotSchedule := getSchedule
	// Map response body to schema and populate Computed attribute values
	plan.VolumeId = types.StringValue(snapshotSchedule.VolumeId)
	plan.SnapshotRetentionCount = types.Int32PointerValue(snapshotSchedule.SnapshotRetentionCount)
	plan.SnapshotPolicyEnabled = types.BoolPointerValue(snapshotSchedule.SnapshotPolicyEnabled.Get())

	for _, snapshotScheduleList := range snapshotSchedule.SnapshotSchedule {
		snapshotScheduleState = filestorage.SnapshotSchedule{
			DayOfWeek: types.StringPointerValue(snapshotScheduleList.DayOfWeek.Get()),
			Frequency: types.StringValue(snapshotScheduleList.Frequency),
			Hour:      types.StringValue(snapshotScheduleList.Hour),
			Id:        types.StringPointerValue(snapshotScheduleList.Id.Get()),
		}

	}

	plan.SnapshotSchedule = snapshotScheduleState

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}
}

func (r *fileStorageSnapshotScheduleResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// Values from plan
	var state filestorage.SnapshotScheduleResource
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Exist Check Schedule
	getData, err := r.client.GetSnapshotScheduleList(ctx, state.VolumeId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Reading Snapshot Schedule",
			"Could not read Snapshot Schedule"+getData.GetVolumeId()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var SnapshotSheduleId = ""

	var snapshotScheduleState = filestorage.SnapshotSchedule{
		DayOfWeek: types.StringPointerValue(nil),
		Frequency: types.StringValue(""),
		Hour:      types.StringValue(""),
		Id:        types.StringPointerValue(nil),
	}

	for _, snapshotScheduleList := range getData.SnapshotSchedule {
		snapshotScheduleState = filestorage.SnapshotSchedule{
			DayOfWeek: types.StringPointerValue(snapshotScheduleList.DayOfWeek.Get()),
			Frequency: types.StringValue(snapshotScheduleList.Frequency),
			Hour:      types.StringValue(snapshotScheduleList.Hour),
			Id:        types.StringPointerValue(snapshotScheduleList.Id.Get()),
		}
	}

	SnapshotSheduleId = snapshotScheduleState.Id.ValueString()

	// Update Schedule
	err = r.client.DeleteSnapshotSchedule(ctx, SnapshotSheduleId, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Deleting SnapshotSchedule",
			"Could not delete SnapshotSchedule, unexpected error: \nReason: "+detail,
		)
		return
	}
}
