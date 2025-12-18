package budget

import (
	"context"
	"math"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	budget "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/budget/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *budget.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: budget.NewAPIClient(config),
	}
}

func (client *Client) CreateAccountBudget(ctx context.Context, request BudgetResource) (*budget.BudgetAccountShowResponse, error) {
	req := client.sdkClient.BudgetV1AccountBudgetsAPIsAPI.CreateAccountBudget(ctx)

	var convertReceivers []string
	if request.Notifications.Receivers.IsNull() || request.Notifications.Receivers.IsUnknown() {
		convertReceivers = []string{}
	} else {
		for _, elem := range request.Notifications.Receivers.Elements() {
			strVal := elem.(types.String)
			convertReceivers = append(convertReceivers, strVal.ValueString())
		}
	}

	var convertThresholds []int32
	if request.Notifications.Thresholds.IsNull() || request.Notifications.Thresholds.IsUnknown() {
		convertThresholds = []int32{}
	} else {
		for _, elem := range request.Notifications.Thresholds.Elements() {
			strVal := elem.(types.Int32)
			convertThresholds = append(convertThresholds, strVal.ValueInt32())
		}
	}

	var notifications = request.Notifications
	var convertNotifications = &budget.NotificationSettingNew{}
	// TODO - Validation
	convertNotifications = &budget.NotificationSettingNew{
		IsUseNotification:      *budget.NewNullableBool(notifications.IsUseNotification.ValueBoolPointer()),
		NotificationSendPeriod: *budget.NewNullableString(notifications.NotificationSendPeriod.ValueStringPointer()),
		Receivers:              convertReceivers,
		Thresholds:             convertThresholds,
	}

	var convertPreventionReceivers []string
	if request.Prevention.Receivers.IsNull() || request.Prevention.Receivers.IsUnknown() {
		convertPreventionReceivers = []string{}
	} else {
		for _, elem := range request.Prevention.Receivers.Elements() {
			strVal := elem.(types.String)
			convertPreventionReceivers = append(convertPreventionReceivers, strVal.ValueString())
		}
	}

	var prevention = request.Prevention
	var convertPrevention = &budget.PreventionSettingNew{}
	// TODO - Validation
	convertPrevention = &budget.PreventionSettingNew{
		IsUsePrevention: *budget.NewNullableBool(prevention.IsUsePrevention.ValueBoolPointer()),
		Receivers:       convertPreventionReceivers,
		Threshold:       *budget.NewNullableInt32(prevention.Threshold.ValueInt32Pointer()),
	}

	req = req.BudgetCreateRequest(budget.BudgetCreateRequest{
		Amount:        request.Amount.ValueInt32(),
		Name:          request.Name.ValueString(),
		Notifications: *budget.NewNullableNotificationSettingNew(convertNotifications),
		Prevention:    *budget.NewNullablePreventionSettingNew(convertPrevention),
		StartMonth:    request.StartMonth.ValueString(),
		Unit:          request.Unit.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteAccountBudget(ctx context.Context, budgetId string) error {
	req := client.sdkClient.BudgetV1AccountBudgetsAPIsAPI.DeleteAccountBudget(ctx, budgetId)
	_, _, err := req.Execute()
	return err
}

func (client *Client) GetAccountBudgetList(ctx context.Context) (*budget.BudgetAccountPageResponse, error) {
	req := client.sdkClient.BudgetV1AccountBudgetsAPIsAPI.ListAccountBudgets(ctx)
	req = req.Size(math.MaxInt32)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) SetAccountBudget(ctx context.Context, budgetId string, request BudgetResource) (*budget.BudgetAccountShowResponse, error) {
	req := client.sdkClient.BudgetV1AccountBudgetsAPIsAPI.SetAccountBudget(ctx, budgetId)

	var convertReceivers []string
	if request.Notifications.Receivers.IsNull() || request.Notifications.Receivers.IsUnknown() {
		convertReceivers = []string{}
	} else {
		for _, elem := range request.Notifications.Receivers.Elements() {
			strVal := elem.(types.String)
			convertReceivers = append(convertReceivers, strVal.ValueString())
		}
	}

	var convertThresholds []int32
	if request.Notifications.Thresholds.IsNull() || request.Notifications.Thresholds.IsUnknown() {
		convertThresholds = []int32{}
	} else {
		for _, elem := range request.Notifications.Thresholds.Elements() {
			strVal := elem.(types.Int32)
			convertThresholds = append(convertThresholds, strVal.ValueInt32())
		}
	}

	var notifications = request.Notifications
	var convertNotifications = &budget.NotificationSettingNew{}
	// TODO - Validation
	convertNotifications = &budget.NotificationSettingNew{
		IsUseNotification:      *budget.NewNullableBool(notifications.IsUseNotification.ValueBoolPointer()),
		NotificationSendPeriod: *budget.NewNullableString(notifications.NotificationSendPeriod.ValueStringPointer()),
		Receivers:              convertReceivers,
		Thresholds:             convertThresholds,
	}

	var convertPreventionReceivers []string
	if request.Prevention.Receivers.IsNull() || request.Prevention.Receivers.IsUnknown() {
		convertPreventionReceivers = []string{}
	} else {
		for _, elem := range request.Prevention.Receivers.Elements() {
			strVal := elem.(types.String)
			convertPreventionReceivers = append(convertPreventionReceivers, strVal.ValueString())
		}
	}

	var prevention = request.Prevention
	var convertPrevention = &budget.PreventionSettingNew{}
	// TODO - Validation
	convertPrevention = &budget.PreventionSettingNew{
		IsUsePrevention: *budget.NewNullableBool(prevention.IsUsePrevention.ValueBoolPointer()),
		Receivers:       convertPreventionReceivers,
		Threshold:       *budget.NewNullableInt32(prevention.Threshold.ValueInt32Pointer()),
	}

	req = req.BudgetSetRequest(budget.BudgetSetRequest{
		Amount:        request.Amount.ValueInt32(),
		Name:          request.Name.ValueString(),
		Notifications: *budget.NewNullableNotificationSettingNew(convertNotifications),
		Prevention:    *budget.NewNullablePreventionSettingNew(convertPrevention),
		StartMonth:    request.StartMonth.ValueString(),
		Unit:          request.Unit.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetAccountBudget(ctx context.Context, budgetId string) (*budget.BudgetAccountShowResponse, error) {
	req := client.sdkClient.BudgetV1AccountBudgetsAPIsAPI.ShowAccountBudget(ctx, budgetId)
	resp, _, err := req.Execute()
	return resp, err
}
