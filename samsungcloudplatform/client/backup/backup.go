package backup

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpbackup "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/backup/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math"
	"strings"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpbackup.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpbackup.NewAPIClient(config),
	}
}

func (client *Client) GetBackupList(Name types.String, ServerName types.String) (*scpbackup.BackupListResponse, error) {
	ctx := context.Background()

	req := client.sdkClient.BackupV1BackupsApiAPI.ListBackups(ctx)
	req = req.Size(math.MaxInt32)

	if !Name.IsNull() {
		req = req.Name(Name.ValueString())
	}
	if !ServerName.IsNull() {
		req = req.ServerName(ServerName.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetScheduleList(ctx context.Context, backupId string) (*scpbackup.BackupScheduleListResponse, error) {
	req := client.sdkClient.BackupV1BackupSchedulesApiAPI.ListBackupSchedules(ctx, backupId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateBackup(ctx context.Context, request BackupResource) (*scpbackup.SyncResponse, error) {
	req := client.sdkClient.BackupV1BackupsApiAPI.CreateBackup(ctx)

	// schedule
	var convertedSchedule []scpbackup.BackupScheduleCreateRequest
	for _, backupSchedule := range request.Schedules {
		var convertedStartDay scpbackup.NullableBackupScheduleDay
		if !backupSchedule.StartDay.IsNull() {
			convertedStartDay.Set(scpbackup.BackupScheduleDay(backupSchedule.StartDay.ValueString()).Ptr())
		} else {
			convertedStartDay.Unset()
		}

		var convertedStartWeek scpbackup.NullableBackupScheduleWeek
		if !backupSchedule.StartWeek.IsNull() {
			convertedStartWeek.Set(scpbackup.BackupScheduleWeek(backupSchedule.StartWeek.ValueString()).Ptr())
		} else {
			convertedStartWeek.Unset()
		}

		convertedSchedule = append(convertedSchedule, scpbackup.BackupScheduleCreateRequest{
			Frequency: scpbackup.BackupScheduleFrequency(backupSchedule.Frequency.ValueString()),
			StartDay:  convertedStartDay,
			StartTime: backupSchedule.StartTime.ValueString(),
			StartWeek: convertedStartWeek,
			Type:      scpbackup.BackupScheduleType(backupSchedule.Type.ValueString()),
		})
	}

	// tags
	var TagsObject []scpbackup.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := scpbackup.Tag{
			Key:   k,
			Value: strings.ReplaceAll(v.String(), "\"", ""),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.BackupCreateRequest(scpbackup.BackupCreateRequest{
		Name:            request.Name.ValueString(),
		PolicyCategory:  scpbackup.BackupPolicyCategory(request.PolicyCategory.ValueString()),
		PolicyType:      scpbackup.BackupPolicyType(request.PolicyType.ValueString()),
		ServerUuid:      request.ServerUuid.ValueString(),
		ServerCategory:  scpbackup.ServerCategory(request.ServerCategory.ValueString()),
		EncryptEnabled:  request.EncryptEnabled.ValueBoolPointer(),
		RetentionPeriod: scpbackup.BackupRetentionPeriod(request.RetentionPeriod.ValueString()).Ptr(),
		Schedules:       convertedSchedule,
		Tags:            TagsObject,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetBackup(ctx context.Context, backupId string) (*scpbackup.BackupDetailResponse, error) {
	req := client.sdkClient.BackupV1BackupsApiAPI.ShowBackup(ctx, backupId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateBackupRetentionPeriod(ctx context.Context, backupId string, request BackupResource) (*scpbackup.SyncResponse, error) {
	req := client.sdkClient.BackupV1BackupsApiAPI.UpdateRetentionPeriod(ctx, backupId)

	req = req.RetentionPeriodUpdateBody(scpbackup.RetentionPeriodUpdateBody{
		RetentionPeriod: scpbackup.BackupRetentionPeriod(request.RetentionPeriod.ValueString()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateSchedule(ctx context.Context, backupId string, request BackupResource) (*scpbackup.SyncResponse, error) {
	req := client.sdkClient.BackupV1BackupSchedulesApiAPI.SetBackupSchedules(ctx, backupId)

	// schedule
	var convertedSchedule []scpbackup.BackupScheduleCreateRequest
	for _, backupSchedule := range request.Schedules {

		var convertedStartDay scpbackup.NullableBackupScheduleDay
		if !backupSchedule.StartDay.IsNull() {
			convertedStartDay.Set(scpbackup.BackupScheduleDay(backupSchedule.StartDay.ValueString()).Ptr())
		} else {
			convertedStartDay.Unset()
		}

		var convertedStartWeek scpbackup.NullableBackupScheduleWeek
		if !backupSchedule.StartWeek.IsNull() {
			convertedStartWeek.Set(scpbackup.BackupScheduleWeek(backupSchedule.StartWeek.ValueString()).Ptr())
		} else {
			convertedStartWeek.Unset()
		}

		convertedSchedule = append(convertedSchedule, scpbackup.BackupScheduleCreateRequest{
			Frequency: scpbackup.BackupScheduleFrequency(backupSchedule.Frequency.ValueString()),
			StartDay:  convertedStartDay,
			StartTime: backupSchedule.StartTime.ValueString(),
			StartWeek: convertedStartWeek,
			Type:      scpbackup.BackupScheduleType(backupSchedule.Type.ValueString()),
		})
	}

	req = req.ModifyBackupSchedulesRequest(scpbackup.ModifyBackupSchedulesRequest{
		Schedules: convertedSchedule,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteBackup(ctx context.Context, backupId string) error {
	req := client.sdkClient.BackupV1BackupsApiAPI.DeleteBackup(ctx, backupId)

	_, _, err := req.Execute()
	return err
}
