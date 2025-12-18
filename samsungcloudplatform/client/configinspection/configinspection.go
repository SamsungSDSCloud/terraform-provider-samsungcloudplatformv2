package configinspection

import (
	"context"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	configinspection "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/configinspection/1.1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *configinspection.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: configinspection.NewAPIClient(config),
	}
}

func (client *Client) GetConfigInspectionList(ctx context.Context, request ConfigInspectionDataSources) (*configinspection.ConfigInspectionListResponse, error) {
	req := client.sdkClient.ConfiginspectionV1ConfigInspectionApiAPI.ConfigInspectionList(ctx)
	if !request.WithCount.IsNull() {
		req = req.WithCount(request.WithCount.ValueString())
	}

	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}

	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
	}

	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}

	if !request.IsMine.IsNull() {
		req = req.IsMine(request.IsMine.ValueBool())
	}

	if !request.DiagnosisID.IsNull() {
		req = req.DiagnosisId(request.DiagnosisID.ValueString())
	}

	if !request.DiagnosisName.IsNull() {
		req = req.DiagnosisName(request.DiagnosisName.ValueString())
	}

	if !request.CSPType.IsNull() {
		req = req.CspType(request.CSPType.ValueString())
	}

	if !request.DiagnosisAccountID.IsNull() {
		req = req.DiagnosisAccountId(request.DiagnosisAccountID.ValueString())
	}

	if !request.StartDate.IsNull() {
		req = req.StartDate(request.StartDate.ValueString())
	}

	if !request.EndDate.IsNull() {
		req = req.EndDate(request.EndDate.ValueString())
	}

	// Handle RecentDiagnosisState slice
	if len(request.RecentDiagnosisState) > 0 {
		var states []*string
		for _, state := range request.RecentDiagnosisState {
			if !state.IsNull() {
				states = append(states, state.ValueStringPointer())
			}
		}
		if len(states) > 0 {
			req = req.RecentDiagnosisState(states)
		}
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetConfigInspectionObjectDetail(ctx context.Context, diagnosisId string) (*configinspection.DiagnosisObjectDetailResponse, error) {
	resp, _, err := client.sdkClient.ConfiginspectionV1ConfigInspectionApiAPI.GetDiagnosisObjectDetail(ctx, diagnosisId).Execute()
	return resp, err
}

func (client *Client) CreateConfigInspectionObject(ctx context.Context, request ConfigInspectionDiagnosisResource) (*configinspection.ConfigInspectionCreateResponse, error) {
	req := client.sdkClient.ConfiginspectionV1ConfigInspectionApiAPI.CreateDiagnosisObject(ctx)

	tags := convertToTags(request.Tags.Elements())

	authKeyRequest := configinspection.AuthKeyRequest{}
	if request.AuthKeyRequest != nil {
		authKeyRequest = configinspection.AuthKeyRequest{
			AuthKeyId: request.AuthKeyRequest.AuthKeyId.ValueString(),
		}
		authKeyRequest.AuthKeyCreatedAt.Set(request.AuthKeyRequest.AuthKeyCreatedAt.ValueStringPointer())
		authKeyRequest.AuthKeyExpiredAt.Set(request.AuthKeyRequest.AuthKeyExpiredAt.ValueStringPointer())
		authKeyRequest.DiagnosisId.Set(request.AuthKeyRequest.DiagnosisId.ValueStringPointer())

	}

	scheduleRequest := configinspection.DiagnosisScheduleRequest{}
	if request.ScheduleRequest != nil {
		scheduleRequest = configinspection.DiagnosisScheduleRequest{
			DiagnosisId:               request.ScheduleRequest.DiagnosisId.ValueString(),
			DiagnosisStartTimePattern: request.ScheduleRequest.DiagnosisStartTimePattern.ValueString(),
			FrequencyType:             request.ScheduleRequest.FrequencyType.ValueString(),
			FrequencyValue:            request.ScheduleRequest.FrequencyValue.ValueString(),
			UseDiagnosisCheckTypeBp:   request.ScheduleRequest.UseDiagnosisCheckTypeBp.ValueString(),
			UseDiagnosisCheckTypeSsi:  request.ScheduleRequest.UseDiagnosisCheckTypeSsi.ValueString(),
		}
	}

	req = req.DiagnosisObjectRequest(configinspection.DiagnosisObjectRequest{
		AccountId:          request.AccountId.ValueString(),
		AuthKeyRequest:     authKeyRequest,
		CspType:            request.CspType.ValueString(),
		DiagnosisAccountId: request.DiagnosisAccountId.ValueString(),
		DiagnosisCheckType: request.DiagnosisCheckType.ValueString(),
		DiagnosisId:        request.DiagnosisId.ValueString(),
		DiagnosisName:      request.DiagnosisName.ValueString(),
		DiagnosisType:      request.DiagnosisType.ValueString(),
		PlanType:           request.PlanType.ValueString(),
		ScheduleRequest:    scheduleRequest,
		Tags:               tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteConfigInspectionObject(ctx context.Context, diagnosisId string) (*configinspection.TerminateResponse, error) {
	resp, _, err := client.sdkClient.ConfiginspectionV1ConfigInspectionApiAPI.TerminateDiagnosisObject(ctx, diagnosisId).Execute()
	return resp, err
}

func (client *Client) GetConfigInspectionDiagnosisResultList(ctx context.Context, request ConfigInspectionDiagnosisResultListDataSources) (*configinspection.DiagnosisResultListResponse, error) {
	req := client.sdkClient.ConfiginspectionV1ConfigInspectionApiAPI.GetDiagnosisResultList(ctx)

	if !request.WithCount.IsNull() {
		req = req.WithCount(request.WithCount.ValueString())
	}

	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}

	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
	}

	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}

	if !request.DiagnosisID.IsNull() {
		req = req.DiagnosisId(request.DiagnosisID.ValueString())
	}

	if !request.DiagnosisName.IsNull() {
		req = req.DiagnosisName(request.DiagnosisName.ValueString())
	}

	if !request.CSPType.IsNull() {
		req = req.CspType(request.CSPType.ValueString())
	}

	if !request.StartDate.IsNull() {
		req = req.StartDate(request.StartDate.ValueString())
	}

	if !request.EndDate.IsNull() {
		req = req.EndDate(request.EndDate.ValueString())
	}

	if !request.AccountId.IsNull() {
		req = req.AccountId(request.AccountId.ValueString())
	}

	if !request.DiagnosisState.IsNull() {
		req = req.DiagnosisState(request.DiagnosisState.ValueString())
	}

	if !request.UserId.IsNull() {
		req = req.UserId(request.UserId.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetConfigInspectionDiagnosisResultDetail(ctx context.Context, request ConfigInspectionDiagnosisResultDetailDataSource) (*configinspection.DiagnosisResultDetailResponse, error) {
	req := client.sdkClient.ConfiginspectionV1ConfigInspectionApiAPI.GetDiagnosisResultDetail(ctx)

	if !request.DiagnosisId.IsNull() {
		req = req.DiagnosisId(request.DiagnosisId.ValueString())
	}

	if !request.DiagnosisRequestSequence.IsNull() {
		req = req.DiagnosisRequestSequence(request.DiagnosisRequestSequence.ValueString())
	}

	if !request.WithCount.IsNull() {
		req = req.WithCount(request.WithCount.ValueString())
	}

	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}

	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
	}

	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) RequestNewConfigInspectionDiagnosis(ctx context.Context, request configinspection.DiagnosisRequest) (*configinspection.CheckResponse, error) {
	req := client.sdkClient.ConfiginspectionV1ConfigInspectionApiAPI.DiagnosisRequest(ctx)

	req = req.DiagnosisRequest(request)

	resp, _, err := req.Execute()
	return resp, err
}

func convertToTags(elements map[string]attr.Value) []configinspection.Tag {
	var tags []configinspection.Tag
	for k, v := range elements {
		tagObject := configinspection.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}
