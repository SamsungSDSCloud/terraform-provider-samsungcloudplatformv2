package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/budget"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &budgetBudgetDataSource{}
	_ datasource.DataSourceWithConfigure = &budgetBudgetDataSource{}
)

// NewBudgetBudgetDataSource is a helper function to simplify the provider implementation.
func NewBudgetBudgetDataSource() datasource.DataSource {
	return &budgetBudgetDataSource{}
}

// budgetBudgetDataSource is the data source implementation.
type budgetBudgetDataSource struct {
	config  *scpsdk.Configuration
	client  *budget.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *budgetBudgetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_budget_budget"
}

// Schema defines the schema for the data source.
func (d *budgetBudgetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Account budget.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "ID of the budget to filter by.",
				MarkdownDescription: "ID of the budget to filter by.\n\nExample: `aaaaaa123456789abcede`",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Budget name.",
				MarkdownDescription: "Name of the budget to filter by.\n\nAllowed length: 20 korean character, english character, digit, space, +=,.@-_.\n\nExample: `ex_month_budget`",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 20),
				},
			},
			common.ToSnakeCase("Budget"): schema.SingleNestedAttribute{
				Description: "Account budget details.",
				MarkdownDescription: "Details of the account budget.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("BudgetId"): schema.StringAttribute{
						Description: "Unique ID of the budget.",
						MarkdownDescription: "The unique ID of the budget.\n\nExample: `ab1234567890abcdef`",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description:         "Budget name",
                        MarkdownDescription: "The name of the budget.\n\nExample: `ex_month_budget`",
						Computed:    true,
					},
					common.ToSnakeCase("Amount"): schema.Int32Attribute{
						Description:         "Budget amount",
                        MarkdownDescription: "The budget amount in the specified currency unit.\n\nExample: `1000000`",
						Computed:    true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description:         "Budget type",
						MarkdownDescription: "The type of the budget. Default type is  `COST`\n\nAllowed values: `COST`\n\nExample: `COST`",
						Computed:    true,
					},
					common.ToSnakeCase("Unit"): schema.StringAttribute{
                        Description:         "Budget management unit",
                        MarkdownDescription: "The budget manage unit.\n\nAllowed values: `MONTHLY` | `OVERALL`\n\nExample: `MONTHLY`",
						Computed:    true,
					},
					common.ToSnakeCase("StartMonth"): schema.StringAttribute{
                        Description:         "Budget start month",
                        MarkdownDescription: "The month when the budget period starts.\n\nExample: `2024-01`",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
                        Description:         "Created datetime",
                        MarkdownDescription: "The datetime when the budget was created.\n\nExample: `2024-01-15T00:00:00`",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
                        Description:         "Created user",
                        MarkdownDescription: "The user who created the budget.\n\nExample: `user@example.com`",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
                        Description:         "Modified datetime",
                        MarkdownDescription: "The datetime when the budget was last modified.\n\nExample: `2024-01-15T00:00:00`",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
                        Description:         "Modified user",
                        MarkdownDescription: "The user who last modified the budget.\n\nExample: `user@example.com`",
						Computed:    true,
					},
				},
			},
			common.ToSnakeCase("Notifications"): schema.SingleNestedAttribute{
                Description:         "Notification settings for the budget",
                MarkdownDescription: "Settings for budget alerts when usage exceeds defined thresholds.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("IsUseNotification"): schema.BoolAttribute{
                        Description:         "Notification use state",
                        MarkdownDescription: "Whether to enable budget notifications.\n\nAllowed values:`true`|`false`\n\nExample: `true`",
						Optional:    true,
					},
					common.ToSnakeCase("NotificationSendPeriod"): schema.StringAttribute{
                        Description:         "Notification send period",
                        MarkdownDescription: "When to send budget notifications.\n\nAllowed values: `FIRST` | `DAILY` | `NONE`\n\nExample: `FIRST`",
						Optional:    true,
					},
					common.ToSnakeCase("Receivers"): schema.ListAttribute{
						ElementType: types.StringType,
                        Description:         "List of notification recipient email addresses",
                        MarkdownDescription: "Email addresses to receive budget notifications.\n\nExample: `[\"user@example.com\"]`",
						Optional:    true,
					},
					common.ToSnakeCase("Thresholds"): schema.ListAttribute{
						ElementType: types.Int32Type,
                        Description:         "List of threshold percentages for notifications",
                        MarkdownDescription: "Percentage thresholds at which to send notifications.\n\nAllowed values: `70` | `80` | `90` | `100`\n\nExample: `[80, 100]`",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("Prevention"): schema.SingleNestedAttribute{
                Description:         "Auto generation prevention settings for the budget",
                MarkdownDescription: "Settings to prevent new resource generation when budget usage exceeds the threshold.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("IsUsePrevention"): schema.BoolAttribute{
                        Description:         "Auto Generation prevent use state",
                        MarkdownDescription: "Whether to prevent new resource creation when budget threshold is exceeded.\n\nAllowed values:`true`|`false`\n\nExample: `true`",
						Optional:    true,
					},
					common.ToSnakeCase("Receivers"): schema.ListAttribute{
						ElementType: types.StringType,
                        Description:         "List of notification recipient email addresses",
                        MarkdownDescription: "Email addresses to notify when budget prevention is triggered.\n\nExample: `[\"admin@example.com\"]`",
						Optional:    true,
					},
					common.ToSnakeCase("Threshold"): schema.Int32Attribute{
                        Description:         "Prevention threshold",
                        MarkdownDescription: "Budget threshold percentage to trigger prevention.\n\nAllowed values: One of the numbers `70` | `80` | `90` | `100`\n\nExample: `90`",
						Optional:    true,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{ // 필터는 Block 으로 정의한다.
			"filter": filter.DataSourceSchema(), // 필터 스키마는 공통으로 제공되는 함수를 이용하여 정의한다.
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *budgetBudgetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	d.client = inst.Client.Budget
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *budgetBudgetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state budget.BudgetDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetAccountBudget(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Account Budget",
			err.Error(),
		)
		return
	}

    if data == nil {
        resp.Diagnostics.AddError("Error reading budget", "empty response from API")
        return
    }

    budgetData := data.Budget
	budgetModel := budget.Budget{
		BudgetId:   types.StringValue(budgetData.Id),
		Name:       types.StringValue(budgetData.Name),
		Amount:     convertStringToInt32(budgetData.Amount),
		BudgetType: types.StringValue(budgetData.Type),
		Unit:       types.StringValue(budgetData.Unit),
		CreatedAt:  types.StringValue(budgetData.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringPointerValue(budgetData.CreatedBy),
		ModifiedAt: types.StringValue(budgetData.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringPointerValue(budgetData.ModifiedBy),
		StartMonth: types.StringValue(budgetData.StartMonth), // 추가: StartMonth 매핑
	}
    budgetObjectValue, objDiags := types.ObjectValueFrom(ctx, budgetModel.AttributeTypes(), budgetModel)
    if objDiags.HasError() {
        resp.Diagnostics.Append(objDiags...)
        return
    }

    state.Budget = budgetObjectValue

	notificationsData := data.Notifications

	var notificationsReceivers []string
	for _, receiver := range notificationsData.Receivers {
		notificationsReceivers = append(notificationsReceivers, receiver)
	}

	var notificationsThresholds []int32
	for _, threshold := range notificationsData.Thresholds {
		notificationsThresholds = append(notificationsThresholds, threshold)
	}

	var notificationReceiversListValue,d1 = types.ListValueFrom(ctx, basetypes.StringType{}, notificationsReceivers)
    if d1.HasError() {
        resp.Diagnostics.Append(d1...)
        return
    }
	var notificationThresholdsListValue, d2 = types.ListValueFrom(ctx, basetypes.Int32Type{}, notificationsThresholds)
    if d2.HasError() {
        resp.Diagnostics.Append(d2...)
        return
    }

	notificationsModel := budget.Notifications{
		IsUseNotification:      types.BoolValue(*notificationsData.IsUseNotification.Get()),
		NotificationSendPeriod: types.StringValue(*notificationsData.NotificationSendPeriod.Get()),
		Receivers:              notificationReceiversListValue,
		Thresholds:             notificationThresholdsListValue,
	}
	notificationsObjectValue, d3 := types.ObjectValueFrom(ctx, notificationsModel.AttributeTypes(), notificationsModel)
    if d3.HasError() {
        resp.Diagnostics.Append(d3...)
        return
    }

	state.Notifications = notificationsObjectValue

	preventionData := data.Prevention

	var preventionReceivers []string
	for _, receiver := range preventionData.Receivers {
		preventionReceivers = append(preventionReceivers, receiver)
	}

	var preventionReceiversListValue, d4 = types.ListValueFrom(ctx, basetypes.StringType{}, preventionReceivers)
    if d4.HasError() {
        resp.Diagnostics.Append(d4...)
        return
    }

	preventionModel := budget.Prevention{
		IsUsePrevention: types.BoolValue(*preventionData.IsUsePrevention.Get()),
		Receivers:       preventionReceiversListValue,
		Threshold:       common.ToNullableInt32Value(preventionData.Threshold.Get()),
	}

	preventionObjectValue, d5 := types.ObjectValueFrom(ctx, preventionModel.AttributeTypes(), preventionModel)
    if d5.HasError() {
        resp.Diagnostics.Append(d5...)
        return
    }
	state.Prevention = preventionObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
