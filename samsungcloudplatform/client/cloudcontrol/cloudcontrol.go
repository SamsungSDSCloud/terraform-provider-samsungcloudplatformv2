package cloudcontrol

import (
	"context"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	cloudcontrol "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/cloudcontrol/1.2"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *cloudcontrol.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: cloudcontrol.NewAPIClient(config),
	}
}

// CreateLandingZone creates a new landing zone
func (client *Client) CreateLandingZone(ctx context.Context, request LandingZoneResource) (*cloudcontrol.LandingZoneCreateResponse, error) {
	req := client.sdkClient.CloudcontrolV1LandingZonesAPIsAPI.CreateLandingZone(ctx)

	sdkRequest, err := request.ToSdkType()
	if err != nil {
		return nil, err
	}

	req = req.LandingZoneCreateRequestV1Dot1(*sdkRequest)

	resp, _, err := req.Execute()
	return resp, err
}

// GetLandingZone retrieves a landing zone by ID
func (client *Client) GetLandingZone(ctx context.Context, landingZoneId string) (*cloudcontrol.LandingZoneShowResponseV1Dot1, error) {
	req := client.sdkClient.CloudcontrolV1LandingZonesAPIsAPI.ShowLandingZone(ctx, landingZoneId)

	resp, _, err := req.Execute()
	return resp, err
}

// UpdateLandingZone updates an existing landing zone
func (client *Client) UpdateLandingZone(ctx context.Context, landingZoneId string, request LandingZoneUpdateResource) (*cloudcontrol.LandingZoneCreateResponse, error) {
	req := client.sdkClient.CloudcontrolV1LandingZonesAPIsAPI.SetLandingZone(ctx, landingZoneId)

	sdkRequest, err := request.ToSdkType()
	if err != nil {
		return nil, err
	}

	req = req.LandingZoneUpdateRequest(*sdkRequest)

	resp, _, err := req.Execute()
	return resp, err
}

// DeleteLandingZone deletes a landing zone by ID
func (client *Client) DeleteLandingZone(ctx context.Context, landingZoneId string) (*cloudcontrol.LandingZoneCreateResponse, error) {
	req := client.sdkClient.CloudcontrolV1LandingZonesAPIsAPI.DeleteLandingZone(ctx, landingZoneId)

	resp, _, err := req.Execute()
	return resp, err
}

// CreateAccountFactoryAccount creates a new account factory account
func (client *Client) CreateAccountFactoryAccount(ctx context.Context, request AccountFactoryResource) (*cloudcontrol.AccountFactoryCreateResponse, error) {
	req := client.sdkClient.CloudcontrolV1AccountFactoryAPIsAPI.CreateAccountFactoryAccount(ctx)

	sdkRequest, err := request.ToSdkType()
	if err != nil {
		return nil, err
	}

	req = req.AccountFactoryCreateRequest(*sdkRequest)

	resp, _, err := req.Execute()
	return resp, err
}

// CreateBaselineAssignment creates a new baseline assignment
func (client *Client) CreateBaselineAssignment(ctx context.Context, assignmentId string, request BaselineAssignmentResource) (*cloudcontrol.BaselineAssignmentAddResponse, error) {
	req := client.sdkClient.CloudcontrolV1BaselineAssignmentsAPIsAPI.AddBaselineAssignment(ctx, assignmentId)

	sdkRequest, err := request.ToSdkType()
	if err != nil {
		return nil, err
	}

	req = req.BaselineAssignmentAddRequest(*sdkRequest)

	resp, _, err := req.Execute()
	return resp, err
}

// ListBaselineAssignments retrieves a list of baseline assignments
func (client *Client) ListBaselineAssignments(ctx context.Context, landingZoneId, resourceType, assignmentId, status string) (*cloudcontrol.BaselineAssignmentListResponseV1Dot1, error) {
	req := client.sdkClient.CloudcontrolV1BaselineAssignmentsAPIsAPI.ListBaselineAssignments(ctx)

	if landingZoneId != "" {
		req = req.LandingZoneId(landingZoneId)
	}
	if resourceType != "" {
		req = req.ResourceType(cloudcontrol.BaselineAssignmentResourceTypeEnum(resourceType))
	}
	if assignmentId != "" {
		req = req.AssignmentId(assignmentId)
	}
	if status != "" {
		req = req.Status(cloudcontrol.BaselineAssignmentStateEnum(status))
	}

	resp, _, err := req.Execute()
	return resp, err
}

// UpdateBaselineAssignment updates an existing baseline assignment
func (client *Client) UpdateBaselineAssignment(ctx context.Context, assignmentId string, request BaselineAssignmentUpdateResource) (*cloudcontrol.BaselineAssignmentAddResponse, error) {
	req := client.sdkClient.CloudcontrolV1BaselineAssignmentsAPIsAPI.UpdateBaselineAssignment(ctx, assignmentId)

	sdkRequest, err := request.ToSdkType()
	if err != nil {
		return nil, err
	}

	req = req.BaselineAssignmentUpdateRequest(*sdkRequest)

	resp, _, err := req.Execute()
	return resp, err
}

// DeleteBaselineAssignment deletes a baseline assignment
func (client *Client) DeleteBaselineAssignment(ctx context.Context, landingZoneId, assignmentId, resourceType string) (*cloudcontrol.BaselineAssignmentRemoveResponse, error) {
	removeReq := cloudcontrol.NewBaselineAssignmentRemoveRequest(
		landingZoneId,
		cloudcontrol.BaselineAssignmentTypeEnum(resourceType),
	)

	req := client.sdkClient.CloudcontrolV1BaselineAssignmentsAPIsAPI.RemoveBaselineAssignment(ctx, assignmentId)
	req = req.BaselineAssignmentRemoveRequest(*removeReq)

	resp, _, err := req.Execute()
	return resp, err
}

// ListGuardrails retrieves a list of guardrails
func (client *Client) ListGuardrails(ctx context.Context, landingZoneId, name, excludeUnitId, guidance, serviceName, status, sort string, size, page int32) (*cloudcontrol.GuardrailPageResponse, error) {
	req := client.sdkClient.CloudcontrolV1GuardrailsAPIsAPI.ListGuardrails(ctx)
	if landingZoneId != "" {
		req = req.LandingZoneId(landingZoneId)
	}
	if name != "" {
		req = req.Name(name)
	}
	if excludeUnitId != "" {
		req = req.ExcludeUnitId(excludeUnitId)
	}
	if guidance != "" {
		req = req.Guidance(guidance)
	}
	if serviceName != "" {
		req = req.ServiceName(serviceName)
	}
	if status != "" {
		req = req.Status(status)
	}
	if sort != "" {
		req = req.Sort(sort)
	}
	if size > 0 {
		req = req.Size(size)
	}
	if page > 0 {
		req = req.Page(page)
	}
	resp, _, err := req.Execute()
	return resp, err
}

// GetGuardrail retrieves a single guardrail by ID
func (client *Client) GetGuardrail(ctx context.Context, guardrailId, landingZoneId string) (*cloudcontrol.Guardrail, error) {
	req := client.sdkClient.CloudcontrolV1GuardrailsAPIsAPI.ShowGuardrail(ctx, guardrailId)

	if landingZoneId != "" {
		req = req.LandingZoneId(landingZoneId)
	}

	resp, _, err := req.Execute()
	if err != nil {
		return nil, err
	}

	guardrail := resp.Guardrail
	return &guardrail, nil
}

// EnableGuardrailBindings enables guardrail bindings for specified guardrails and units
func (client *Client) EnableGuardrailBindings(ctx context.Context, guardrailIds []string, unitIds []string, landingZoneId string) (*cloudcontrol.GuardrailAssignmentResponse, error) {
	req := cloudcontrol.NewGuardrailEnableRequest(guardrailIds, unitIds)

	if landingZoneId != "" {
		req.SetLandingZoneId(landingZoneId)
	}

	resp, _, err := client.sdkClient.CloudcontrolV1AssignmentsAPIsAPI.EnableGuardrailBindings(ctx).GuardrailEnableRequest(*req).Execute()
	return resp, err
}

// DisableGuardrailBindings disables guardrail bindings for specified guardrails and units
func (client *Client) DisableGuardrailBindings(ctx context.Context, guardrailIds []string, unitIds []string, landingZoneId string) (*cloudcontrol.GuardrailAssignmentResponse, error) {
	req := cloudcontrol.NewGuardrailDisableRequest(guardrailIds, unitIds)

	if landingZoneId != "" {
		req.SetLandingZoneId(landingZoneId)
	}

	resp, _, err := client.sdkClient.CloudcontrolV1AssignmentsAPIsAPI.DisableGuardrailBindings(ctx).GuardrailDisableRequest(*req).Execute()
	return resp, err
}

// ListGuardrailsForTarget lists guardrails for a target (unit/account)
func (client *Client) ListGuardrailsForTarget(ctx context.Context, unitId, landingZoneId, name, sort string, size, page int32) (*cloudcontrol.ListGuardrailsForTargetResponse, error) {
	req := client.sdkClient.CloudcontrolV1AssignmentsAPIsAPI.ListGuardrailsForTarget(ctx).TargetId(unitId)

	if landingZoneId != "" {
		req = req.LandingZoneId(landingZoneId)
	}
	if name != "" {
		req = req.Name(name)
	}
	if sort != "" {
		req = req.Sort(sort)
	}
	if size > 0 {
		req = req.Size(size)
	}
	if page > 0 {
		req = req.Page(page)
	}
	resp, _, err := req.Execute()
	return resp, err
}

// ListTargetsForGuardrail lists targets (accounts/OUs) for a guardrail
func (client *Client) ListTargetsForGuardrail(ctx context.Context, guardrailId, targetType, landingZoneId, name, sort string, size, page int32) (*cloudcontrol.ListTargetsForGuardrailResponse, error) {
	req := client.sdkClient.CloudcontrolV1AssignmentsAPIsAPI.ListTargetsForGuardrail(ctx).GuardrailId(guardrailId).TargetType(cloudcontrol.GuardrailAssignmentTargetType(targetType))
	if landingZoneId != "" {
		req = req.LandingZoneId(landingZoneId)
	}
	if name != "" {
		req = req.Name(name)
	}
	if sort != "" {
		req = req.Sort(sort)
	}
	if size > 0 {
		req = req.Size(size)
	}
	if page > 0 {
		req = req.Page(page)
	}
	resp, _, err := req.Execute()
	return resp, err
}