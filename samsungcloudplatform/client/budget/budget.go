package budget

import (
	"context"
	"math"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	budget "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/budget/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func convertAmountToFloat64(amount types.Number) float64 {
    if amount.IsNull() || amount.IsUnknown() {
        return 0.0
    }
    bigFloat := amount.ValueBigFloat()
    if bigFloat == nil {
        return 0.0
    }
    f, _ := bigFloat.Float64()
    return f
}

// float64 → budget.Amount 변환
func convertToBudgetAmount(value float64) budget.Amount {
    f32 := float32(value)
    return budget.Amount{
        Float32: &f32,
    }
}

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

func (client *Client) CreateAccountBudget(ctx context.Context, request BudgetResource) (*budget.BudgetAccountShowResponseV1dot1, error) {
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
		Amount:        convertToBudgetAmount(convertAmountToFloat64(request.Amount)),
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
	_, err := req.Execute()
	return err
}

func (client *Client) GetAccountBudgetList(ctx context.Context) (*budget.BudgetAccountPageResponseV1dot1, error) {
	req := client.sdkClient.BudgetV1AccountBudgetsAPIsAPI.ListAccountBudgets(ctx)
	req = req.Size(math.MaxInt32)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) SetAccountBudget(ctx context.Context, budgetId string, request BudgetResource) (*budget.BudgetAccountShowResponseV1dot1, error) {
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
		Amount:       convertToBudgetAmount(convertAmountToFloat64(request.Amount)),
		Name:          request.Name.ValueString(),
		Notifications: *budget.NewNullableNotificationSettingNew(convertNotifications),
		Prevention:    *budget.NewNullablePreventionSettingNew(convertPrevention),
		StartMonth:    request.StartMonth.ValueString(),
		Unit:          request.Unit.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetAccountBudget(ctx context.Context, budgetId string) (*budget.BudgetAccountShowResponseV1dot1, error) {
	req := client.sdkClient.BudgetV1AccountBudgetsAPIsAPI.ShowAccountBudget(ctx, budgetId)
	resp, _, err := req.Execute()
	return resp, err
}
