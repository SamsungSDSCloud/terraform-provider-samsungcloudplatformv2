package loggingaudit

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/loggingaudit/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *loggingaudit.APIClient // 서비스의 client 를 구조체에 추가한다.
}

func NewClient(config *scpsdk.Configuration) *Client { // client 생성 함수를 추가한다.
	return &Client{
		Config:    config,
		sdkClient: loggingaudit.NewAPIClient(config),
	}
}

// Trail
func (client *Client) GetTrailList(ctx context.Context, request TrailDataSource) (*loggingaudit.TrailListResponseV1dot1, error) {
	req := client.sdkClient.LoggingauditV1TrailsApiAPI.ListTrails(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.BucketName.IsNull() {
		req = req.BucketName(request.BucketName.String())
	}
	if !request.State.IsNull() {
		req = req.State(request.State.String())
	}
	if !request.ResourceType.IsNull() {
		req = req.ResourceType(request.ResourceType.String())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetTrail(ctx context.Context, trailId string) (*loggingaudit.TrailShowResponseV1dot1, error) {
	req := client.sdkClient.LoggingauditV1TrailsApiAPI.ShowTrail(ctx, trailId)

	resp, _, err := req.Execute()
	return resp, err
}

func extractString(s basetypes.StringValue) string {
	return s.String() // 또는 s.Value, 구조체에 따라
}

func (client *Client) CreateTrail(ctx context.Context, request TrailResource) (*loggingaudit.TrailShowResponseV1dot1, error) {
	req := client.sdkClient.LoggingauditV1TrailsApiAPI.CreateTrail(ctx)

	var tags []map[string]string

	for _, tag := range request.TagCreateRequests {
		tags = append(tags, map[string]string{
			"key":   tag.Key.ValueString(),
			"value": tag.Value.ValueString(),
		})
	}

	req = req.TrailCreateRequestV1dot1(loggingaudit.TrailCreateRequestV1dot1{
		AccountId:           request.AccountId.ValueString(),
		BucketName:          request.BucketName.ValueString(),
		BucketRegion:        request.BucketRegion.ValueString(),
		LogTypeTotalYn:      *loggingaudit.NewNullableString(request.LogTypeTotalYn.ValueStringPointer()),
		LogVerificationYn:   *loggingaudit.NewNullableString(request.LogVerificationYn.ValueStringPointer()),
		RegionNames:         ConvertStringListToInterfaceList(request.RegionNames),
		RegionTotalYn:       *loggingaudit.NewNullableString(request.RegionTotalYn.ValueStringPointer()),
		ResourceTypeTotalYn: *loggingaudit.NewNullableString(request.ResourceTypeTotalYn.ValueStringPointer()),
		TagCreateRequests:   tags,
		TargetLogTypes:      ConvertStringListToInterfaceList(request.TargetLogTypes),
		TargetResourceTypes: ConvertStringListToInterfaceList(request.TargetResourceTypes),
		TargetUsers:         ConvertStringListToInterfaceList(request.TargetUsers),
		TrailDescription:    *loggingaudit.NewNullableString(request.TrailDescription.ValueStringPointer()),
		TrailName:           request.TrailName.ValueString(),
		TrailSaveType:       request.TrailSaveType.ValueString(),
		UserTotalYn:         *loggingaudit.NewNullableString(request.UserTotalYn.ValueStringPointer()),
		OrganizationTrailYn: *loggingaudit.NewNullableString(request.OrganizationTrailYn.ValueStringPointer()),
		LogArchiveAccountId: *loggingaudit.NewNullableString(request.LogArchiveAccountId.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteTrailKey(ctx context.Context, trailId string) error {
	req := client.sdkClient.LoggingauditV1TrailsApiAPI.DeleteTrail(ctx, trailId)

	_, err := req.Execute()
	return err
}

func (client *Client) SetTrail(ctx context.Context, trailId string, request TrailResource) (*loggingaudit.TrailShowResponseV1dot1, error) {
	req := client.sdkClient.LoggingauditV1TrailsApiAPI.SetTrail(ctx, trailId)

	req = req.TrailSetRequestV1dot1(loggingaudit.TrailSetRequestV1dot1{
		LogTypeTotalYn:      *loggingaudit.NewNullableString(request.LogTypeTotalYn.ValueStringPointer()),
		LogVerificationYn:   *loggingaudit.NewNullableString(request.LogVerificationYn.ValueStringPointer()),
		RegionNames:         ConvertStringListToInterfaceList(request.RegionNames),
		RegionTotalYn:       *loggingaudit.NewNullableString(request.RegionTotalYn.ValueStringPointer()),
		ResourceTypeTotalYn: *loggingaudit.NewNullableString(request.ResourceTypeTotalYn.ValueStringPointer()),
		TargetLogTypes:      ConvertStringListToInterfaceList(request.TargetLogTypes),
		TargetResourceTypes: ConvertStringListToInterfaceList(request.TargetResourceTypes),
		TargetUsers:         ConvertStringListToInterfaceList(request.TargetUsers),
		TrailDescription:    *loggingaudit.NewNullableString(request.TrailDescription.ValueStringPointer()),
		TrailSaveType:       *loggingaudit.NewNullableString(request.TrailSaveType.ValueStringPointer()),
		UserTotalYn:         *loggingaudit.NewNullableString(request.UserTotalYn.ValueStringPointer()),
		OrganizationTrailYn: *loggingaudit.NewNullableString(request.OrganizationTrailYn.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) StartTrail(ctx context.Context, trailId string) error {
	req := client.sdkClient.LoggingauditV1TrailsApiAPI.StartTrail(ctx, trailId)

	_, _, err := req.Execute()
	return err
}

func (client *Client) StopTrail(ctx context.Context, trailId string) error {
	req := client.sdkClient.LoggingauditV1TrailsApiAPI.StopTrail(ctx, trailId)

	_, _, err := req.Execute()
	return err
}
