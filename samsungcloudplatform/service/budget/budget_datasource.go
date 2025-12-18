package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/budget"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
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
				Description: "ID",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name (between 1 and 64 characters)",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			common.ToSnakeCase("Budget"): schema.SingleNestedAttribute{
				Description: "Account budget.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("BudgetId"): schema.StringAttribute{
						Description: "BudgetId",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("Amount"): schema.Int32Attribute{
						Description: "Amount",
						Computed:    true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "Type",
						Computed:    true,
					},
					common.ToSnakeCase("Unit"): schema.StringAttribute{
						Description: "Unit",
						Computed:    true,
					},
					common.ToSnakeCase("StartMonth"): schema.StringAttribute{
						Description: "StartMonth",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
				},
			},
			common.ToSnakeCase("Notifications"): schema.SingleNestedAttribute{
				Description: "Notifications",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("IsUseNotification"): schema.BoolAttribute{
						Description: "IsUseNotification",
						Optional:    true,
					},
					common.ToSnakeCase("NotificationSendPeriod"): schema.StringAttribute{
						Description: "NotificationSendPeriod",
						Optional:    true,
					},
					common.ToSnakeCase("Receivers"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "Receivers",
						Optional:    true,
					},
					common.ToSnakeCase("Thresholds"): schema.ListAttribute{
						ElementType: types.Int32Type,
						Description: "Thresholds",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("Prevention"): schema.SingleNestedAttribute{
				Description: "Prevention",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("IsUsePrevention"): schema.BoolAttribute{
						Description: "IsUsePrevention",
						Optional:    true,
					},
					common.ToSnakeCase("Receivers"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "Receivers",
						Optional:    true,
					},
					common.ToSnakeCase("Threshold"): schema.Int32Attribute{
						Description: "Threshold",
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

	budgetData := data.Budget
	budgetModel := budget.Budget{
		BudgetId:   types.StringValue(budgetData.Id),
		Name:       types.StringValue(budgetData.Name),
		Amount:     types.Int32Value(budgetData.Amount),
		BudgetType: types.StringValue(budgetData.Type),
		Unit:       types.StringValue(budgetData.Unit),
		CreatedAt:  types.StringValue(budgetData.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringPointerValue(budgetData.CreatedBy),
		ModifiedAt: types.StringValue(budgetData.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringPointerValue(budgetData.ModifiedBy),
	}
	budgetObjectValue, _ := types.ObjectValueFrom(ctx, budgetModel.AttributeTypes(), budgetModel)
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

	var notificationReceiversListValue, _ = types.ListValueFrom(ctx, basetypes.StringType{}, notificationsReceivers)
	var notificationThresholdsListValue, _ = types.ListValueFrom(ctx, basetypes.Int32Type{}, notificationsThresholds)

	notificationsModel := budget.Notifications{
		IsUseNotification:      types.BoolValue(*notificationsData.IsUseNotification.Get()),
		NotificationSendPeriod: types.StringValue(*notificationsData.NotificationSendPeriod.Get()),
		Receivers:              notificationReceiversListValue,
		Thresholds:             notificationThresholdsListValue,
	}
	notificationsObjectValue, _ := types.ObjectValueFrom(ctx, notificationsModel.AttributeTypes(), notificationsModel)
	state.Notifications = notificationsObjectValue

	preventionData := data.Prevention

	var preventionReceivers []string
	for _, receiver := range preventionData.Receivers {
		preventionReceivers = append(preventionReceivers, receiver)
	}

	var preventionReceiversListValue, _ = types.ListValueFrom(ctx, basetypes.StringType{}, preventionReceivers)

	preventionModel := budget.Prevention{
		IsUsePrevention: types.BoolValue(*preventionData.IsUsePrevention.Get()),
		Receivers:       preventionReceiversListValue,
		Threshold:       common.ToNullableInt32Value(preventionData.Threshold.Get()),
	}

	preventionObjectValue, _ := types.ObjectValueFrom(ctx, preventionModel.AttributeTypes(), preventionModel)
	state.Prevention = preventionObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
