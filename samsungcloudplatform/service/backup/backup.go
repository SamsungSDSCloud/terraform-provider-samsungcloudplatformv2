package backup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/backup"
	common "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	backuputil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/backup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/region"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpbackup "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/backup/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource              = &backupBackupResource{}
	_ resource.ResourceWithConfigure = &backupBackupResource{}
)

func NewBackupBackupResource() resource.Resource {
	return &backupBackupResource{}
}

type backupBackupResource struct {
	config  *scpsdk.Configuration
	client  *backup.Client
	clients *client.SCPClient
}

func (r *backupBackupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_backup"
}

func (r *backupBackupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Backup",
		Attributes: map[string]schema.Attribute{
			"region": region.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Backup name \n" +
					"  - example: 'terraformtestbackup01'",
				Required: true,
			},
			common.ToSnakeCase("PolicyCategory"): schema.StringAttribute{
				Description: "Backup policy category \n" +
					"  - example: 'AGENTLESS'",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("AGENTLESS"),
				},
			},
			common.ToSnakeCase("PolicyType"): schema.StringAttribute{
				Description: "Backup policy type \n" +
					"  - example: 'VM_IMAGE'",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("VM_IMAGE"),
				},
			},
			common.ToSnakeCase("ServerUuid"): schema.StringAttribute{
				Description: "Backup server UUID \n" +
					"  - example: 'a16687f2-3abc-4f40-bb5d-ee79ea21249d'",
				Required: true,
			},
			common.ToSnakeCase("ServerCategory"): schema.StringAttribute{
				Description: "Backup server category \n" +
					"  - example: 'VIRTUAL_SERVER'",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("VIRTUAL_SERVER", "GPU_SERVER"),
				},
			},
			common.ToSnakeCase("EncryptEnabled"): schema.BoolAttribute{
				Description: "Whether to use Encryption \n" +
					"  - example: true",
				Required: true,
				Validators: []validator.Bool{
					boolvalidator.Equals(true),
				},
			},
			common.ToSnakeCase("RetentionPeriod"): schema.StringAttribute{
				Description: "Backup retention period \n" +
					"  - example: 'MONTH_1' \n" +
					"  - pattern: WEEK_2 / MONTH_1 / MONTH_3 / MONTH_6 / YEAR_1",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("WEEK_2", "MONTH_1", "MONTH_3", "MONTH_6", "YEAR_1"),
				},
			},
			common.ToSnakeCase("Schedules"): schema.ListNestedAttribute{
				Description: "Backup Schedules",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "Schedule type \n" +
								"  - example: 'FULL' \n" +
								"  - pattern: FULL, INCREMENTAL",
							Required: true,
						},
						common.ToSnakeCase("Frequency"): schema.StringAttribute{
							Description: "Schedule frequency type \n" +
								"  - example: 'DAILY' \n" +
								"  - pattern: MONTHLY / WEEKLY / DAILY",
							Required: true,
						},
						common.ToSnakeCase("StartTime"): schema.StringAttribute{
							Description: "Backup schedule start time \n" +
								"  - example: '11:00:00'",
							Required: true,
						},
						common.ToSnakeCase("StartDay"): schema.StringAttribute{
							Description: "Backup schedule start day \n" +
								"  - example: 'MON' \n" +
								"  - pattern: MON / TUE / WED / THU / FRI / SAT / SUN",
							Optional: true,
						},
						common.ToSnakeCase("StartWeek"): schema.StringAttribute{
							Description: "Backup schedule start week \n" +
								"  - example: 'WEEK_2' \n" +
								"  - pattern: WEEK_1 / WEEK_2 / WEEK_3 / WEEK_4 / WEEK_LAST",
							Optional: true,
						},
					},
				},
			},
			"tags": tag.ResourceSchema(),
		},
	}
}

func (r *backupBackupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Backup
	r.clients = inst.Client
}

func (r *backupBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan backup.BackupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Region.IsNull() {
		r.client.Config.Region = plan.Region.ValueString()
		r.clients.Iam.Config.Region = plan.Region.ValueString()
	}

	data, err := r.client.CreateBackup(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating backup",
			"Could not create backup, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	getData, err := r.client.GetBackup(ctx, data.Resource.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Backup",
			"Could not create Backup, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, "backup", "backup", data.Resource.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Backup Tags",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	state, err := r.MapGetResponseToState(ctx, getData, plan, tagsMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Server",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *backupBackupResource) MapGetResponseToState(ctx context.Context, getData *scpbackup.BackupDetailResponse1Dot1, plan backup.BackupResource, tagsMap types.Map) (backup.BackupResource, error) {
	getSchedules, err := r.client.GetScheduleList(ctx, getData.Id)
	if err != nil {
		return backup.BackupResource{}, err
	}

	var backupSchedules []backup.Schedule
	if len(plan.Schedules) == 0 {
		for _, backupSchedule := range getSchedules.Contents {
			backupSchedules = append(backupSchedules, backup.Schedule{
				Frequency: types.StringValue(string(backupSchedule.Frequency)),
				StartDay:  types.StringPointerValue(backupSchedule.StartDay.Get()),
				StartTime: types.StringPointerValue(backupSchedule.StartTime.Get()),
				StartWeek: types.StringPointerValue(backupSchedule.StartWeek.Get()),
				Type:      types.StringValue(backupSchedule.Type),
			})
		}
	} else {
		for _, planSchedule := range plan.Schedules {
			Frequency := planSchedule.Frequency
			StartDay := planSchedule.StartDay
			StartTime := planSchedule.StartTime
			StartWeek := planSchedule.StartWeek
			Type := planSchedule.Type

			for _, backupSchedule := range getSchedules.Contents {
				if Frequency == types.StringValue(string(backupSchedule.Frequency)) &&
					StartDay == types.StringPointerValue(backupSchedule.StartDay.Get()) &&
					StartTime == types.StringPointerValue(backupSchedule.StartTime.Get()) &&
					StartWeek == types.StringPointerValue(backupSchedule.StartWeek.Get()) &&
					Type == types.StringValue(backupSchedule.Type) {
					backupSchedules = append(backupSchedules, backup.Schedule{
						Frequency: types.StringValue(string(backupSchedule.Frequency)),
						StartDay:  types.StringPointerValue(backupSchedule.StartDay.Get()),
						StartTime: types.StringPointerValue(backupSchedule.StartTime.Get()),
						StartWeek: types.StringPointerValue(backupSchedule.StartWeek.Get()),
						Type:      types.StringValue(backupSchedule.Type),
					})
					break
				}
			}
		}
	}

	return backup.BackupResource{
		Id:              types.StringValue(getData.Id),
		Name:            types.StringValue(getData.Name),
		PolicyCategory:  types.StringValue(getData.PolicyCategory),
		PolicyType:      types.StringValue(getData.PolicyType),
		ServerUuid:      types.StringValue(getData.ServerUuid),
		ServerCategory:  types.StringValue(getData.ServerCategory),
		EncryptEnabled:  types.BoolPointerValue(getData.EncryptEnabled.Get()),
		RetentionPeriod: types.StringValue(getData.RetentionPeriod),
		Schedules:       backupSchedules,
		Tags:            tagsMap,
		Region:          plan.Region,
	}, nil
}

func (r *backupBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state backup.BackupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Region.IsNull() {
		r.client.Config.Region = state.Region.ValueString()
		r.clients.Iam.Config.Region = state.Region.ValueString()
	}

	getData, err := r.client.GetBackup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Backup",
			"Could not read Backup ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, "backup", "backup", state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Backup Tags",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, state.Tags)

	newState, err := r.MapGetResponseToState(ctx, getData, state, tagsMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Server",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *backupBackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	handlers := []*backup.UpdateHandler{
		{
			Fields:  []string{"RetentionPeriod"},
			Handler: r.handlerUpdateRetentionPeriod,
		},
		{
			Fields:  []string{"Schedules"},
			Handler: r.handlerUpdateSchedule,
		},
		{
			Fields:  []string{"Tags"},
			Handler: r.handlerUpdateTag,
		},
	}

	var plan backup.BackupResource
	var state backup.BackupResource
	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var settableFileds []string
	for attrName, attribute := range req.Plan.Schema.GetAttributes() {
		if attribute.IsRequired() || attribute.IsOptional() {
			settableFileds = append(settableFileds, backuputil.SnakeToPascal(attrName))
		}
	}

	changeFields, err := backuputil.GetChangedFields(plan, state, settableFileds)
	if err != nil {
		return
	}

	immuntableFields := []string{"Name", "PolicyCategory", "PolicyType", "ServerUuid", "ServerCategory", "EncryptEnabled"}

	if backuputil.IsOverlapFields(immuntableFields, changeFields) {
		resp.Diagnostics.AddError(
			"Error Updating Backup",
			"Immutable fields cannot be modified: "+strings.Join(immuntableFields, ", "),
		)
		return
	}

	if !state.Region.IsNull() {
		r.client.Config.Region = state.Region.ValueString()
		r.clients.Iam.Config.Region = state.Region.ValueString()
	}

	for _, h := range handlers {
		if backuputil.IsOverlapFields(h.Fields, changeFields) {
			if err := h.Handler(ctx, req, resp); err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Backup",
					"Could not update backup, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}
	}

	data, err := r.client.GetBackup(ctx, plan.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading backup",
			"Could not read backup ID "+plan.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, "backup", "backup", plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Backup Tags",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	newState, err := r.MapGetResponseToState(ctx, data, plan, tagsMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Server",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *backupBackupResource) handlerUpdateRetentionPeriod(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan backup.BackupResource
	req.Plan.Get(ctx, &plan)

	_, err := r.client.UpdateBackupRetentionPeriod(ctx, plan.Id.ValueString(), plan)
	if err != nil {
		return err
	}
	return nil
}

func (r *backupBackupResource) handlerUpdateSchedule(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan backup.BackupResource
	req.Plan.Get(ctx, &plan)

	_, err := r.client.UpdateSchedule(ctx, plan.Id.ValueString(), plan)
	if err != nil {
		return err
	}
	return nil
}

func (r *backupBackupResource) handlerUpdateTag(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan backup.BackupResource
	req.Plan.Get(ctx, &plan)

	_, err := tag.UpdateTags(r.clients, "backup", "backup", plan.Id.ValueString(), plan.Tags.Elements())
	if err != nil {
		return err
	}
	return nil
}

func (r *backupBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state backup.BackupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Region.IsNull() {
		r.client.Config.Region = state.Region.ValueString()
		r.clients.Iam.Config.Region = state.Region.ValueString()
	}

	tag.UpdateTags(r.clients, "backup", "backup", state.Id.ValueString(), make(map[string]attr.Value))

	err := r.client.DeleteBackup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting backup",
			"Could not delete backup, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
