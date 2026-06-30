package budget

import (
	"context"
	"fmt"
	"time"
	"strings"
	"strconv"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/budget" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &budgetBudgetResource{}
	_ resource.ResourceWithConfigure = &budgetBudgetResource{}
	// [BUDGET-FIX-01] ImportState 인터페이스 추가
	_ resource.ResourceWithImportState = &budgetBudgetResource{}
)

func convertStringToInt32(value string) types.Int32 {
    if value == "" {
        return types.Int32Null()
    }
    v, _ := strconv.Atoi(value)
    return types.Int32Value(int32(v))
}

// NewBudgetBudgetResource is a helper function to simplify the provider implementation.
func NewBudgetBudgetResource() resource.Resource {
	return &budgetBudgetResource{}
}

// ============================================================================
// [BUDGET-FIX-01] ImportState 구현 - 기존 리소스 import 기능 추가
// ============================================================================
func (r *budgetBudgetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// ID 기반 단일 리소스 import 시 표준 패턴
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// budgetDataSource is the data source implementation.
type budgetBudgetResource struct {
	config  *scpsdk.Configuration
	client  *budget.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *budgetBudgetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_budget_budget"
}

// Schema defines the schema for the data source.
func (r *budgetBudgetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = BudgetDataSourceSchema()
}

// Configure adds the provider configured client to the data source.
func (r *budgetBudgetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Budget
	r.clients = inst.Client
}

// Create creates the resource and sets Terraform state.
func (r *budgetBudgetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan budget.BudgetResource
	diags := req.Plan.Get(ctx, &plan)

	if len(diags) > 0 {
		for i, diag := range diags {
			fmt.Printf("  [%d] Severity: %s, Summary: %s, Detail: %s\n",
				i,
				diag.Severity(),
				diag.Summary(),
				diag.Detail())
		}
	} else {
		fmt.Println("No diagnostics found.")
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateAccountBudget(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating budget",
			"Could not create budget, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
    if data == nil {
        resp.Diagnostics.AddError(
            "Error creating budget",
            "Empty response from API",
        )
        return
    }

	budgetData := data.Budget
	plan.Id = types.StringValue(budgetData.Id)
    // [BUDGET-FIX-02] 상태 폴링/재시도 로직 추가
    // 리소스 생성이 완료될 때까지 폴링
    // ID가 유효한 경우에만 폴링
	if !plan.Id.IsNull() && plan.Id.ValueString() != "" {
		// 생성 후 리소스가 준비될 때까지 대기
		err = r.waitForBudgetReady(ctx, plan.Id.ValueString(), 60*time.Second)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for budget",
				"Budget was created but failed to become ready: "+err.Error(),
			)
			return
		}

		// 최신 상태 다시 조회
		// 1. 변수명을 showData로 변경하여 타입 충돌 방지
		showData, err := r.client.GetAccountBudget(ctx, plan.Id.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error reading budget after creation",
				// 2. plan.Id를 문자열로 변환 (ValueString 사용)
				"Could not budget ID "+plan.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}
		// 3. 변수명 동기화
		budgetData = showData.Budget
	}

	budgetModel := budget.Budget{
		Amount:     convertStringToInt32(budgetData.Amount),
		CreatedAt:  types.StringValue(budgetData.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringPointerValue(budgetData.CreatedBy),
		BudgetId:   types.StringValue(budgetData.Id),
		ModifiedAt: types.StringValue(budgetData.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringPointerValue(budgetData.ModifiedBy),
		Name:       types.StringValue(budgetData.Name),
		StartMonth: types.StringValue(budgetData.StartMonth),
		Unit:       types.StringValue(budgetData.Unit),
	}
	budgetObjectValue, diags := types.ObjectValueFrom(ctx, budgetModel.AttributeTypes(), budgetModel)
    // [BUDGET-FIX-04] diagnostics 폐기 방지 - diags가 있으면 Append
    if diags.HasError() {
        resp.Diagnostics.Append(diags...)
        return
    }

	plan.Budget = budgetObjectValue
	plan.Id = types.StringValue(budgetData.Id)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *budgetBudgetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state budget.BudgetResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.SetAccountBudget(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Budget",
			"Could not update Budget, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetResourceGroup as UpdateResourceGroup items are not populated.
	data, err := r.client.GetAccountBudget(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Budget",
			"Could not read Budget ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
    if data == nil {
        resp.Diagnostics.AddError(
            "Error reading budget",
            "Empty response from API",
        )
        return
    }

	budgetData := data.Budget
	budgetModel := budget.Budget{
		Amount:     convertStringToInt32(budgetData.Amount),
		CreatedAt:  types.StringValue(budgetData.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringPointerValue(budgetData.CreatedBy),
		BudgetId:   types.StringValue(budgetData.Id),
		ModifiedAt: types.StringValue(budgetData.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringPointerValue(budgetData.ModifiedBy),
		Name:       types.StringValue(budgetData.Name),
		StartMonth: types.StringValue(budgetData.StartMonth),
		Unit:       types.StringValue(budgetData.Unit),
	}
	budgetObjectValue, diags := types.ObjectValueFrom(ctx, budgetModel.AttributeTypes(), budgetModel)
    // [BUDGET-FIX-04] diagnostics 폐기 방지
    if diags.HasError() {
        resp.Diagnostics.Append(diags...)
        return
    }

	state.Budget = budgetObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and sets the updated Terraform state on success.
func (r *budgetBudgetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state budget.BudgetResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteAccountBudget(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Budget",
			"Could not delete Budget, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

// Read read the resource and sets the updated Terraform state on success.
func (r *budgetBudgetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state budget.BudgetResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetAccountBudget(ctx, state.Id.ValueString())
	if err != nil {
		// BUDGET-FIX-03] NotFound(404) 시 RemoveResource 호출
		// 외부에서 삭제된 리소스인 경우 Terraform 상태에서 제거
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Account Budget",
			"Could not read Budget ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

    if data == nil {
        resp.Diagnostics.AddError("Error reading budget", "empty response from API")
        return
    }

	budgetElement := data.Budget

	budgetModel := budget.Budget{
		BudgetId:   types.StringValue(budgetElement.Id),
		Name:       types.StringValue(budgetElement.Name),
		Amount:     convertStringToInt32(budgetElement.Amount),
		BudgetType: types.StringValue(budgetElement.Type),
		Unit:       types.StringValue(budgetElement.Unit),
		CreatedAt:  types.StringValue(budgetElement.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringPointerValue(budgetElement.CreatedBy),
		ModifiedAt: types.StringValue(budgetElement.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringPointerValue(budgetElement.ModifiedBy),
		StartMonth: types.StringValue(budgetElement.StartMonth), // 추가: StartMonth 매핑
	}

	budgetObjectValue, diags := types.ObjectValueFrom(ctx, budgetModel.AttributeTypes(), budgetModel)
	if diags.HasError() {
       resp.Diagnostics.Append(diags...)
       return
    }

	state.Budget = budgetObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func BudgetDataSourceSchema() schema.Schema {
    return schema.Schema{
        Attributes: map[string]schema.Attribute{
            common.ToSnakeCase("Id"): schema.StringAttribute{
                Description:        "Identifier of the resource.",
                MarkdownDescription: "The unique identifier of the budget resource.\n\nExample: `bud-1234567890abcdef`",
                Computed:   true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            common.ToSnakeCase("LastUpdated"): schema.StringAttribute{
                Description:        "Timestamp of the last Terraform update of the Budget",
                MarkdownDescription: "The timestamp when the budget was last updated by Terraform.\n\nExample: `Monday, 02 Jan 2006 15:04:05 -0700`",
                Computed:   true,
            },
            common.ToSnakeCase("Budget"): schema.SingleNestedAttribute{
                Attributes: map[string]schema.Attribute{
                    common.ToSnakeCase("Amount"): schema.Int32Attribute{
                        Computed:            true,
                        Description:         "Budget amount",
                        MarkdownDescription: "The budget amount in the specified currency unit.\n\nExample: `1000000`",
                    },
                    common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
                        Computed:            true,
                        Description:         "Created datetime",
                        MarkdownDescription: "The datetime when the budget was created.\n\nExample: `2024-01-15T00:00:00`",
                    },
                    common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
                        Computed:            true,
                        Description:         "Created user",
                        MarkdownDescription: "The user who created the budget.\n\nExample: `user@example.com`",
                    },
                    common.ToSnakeCase("BudgetId"): schema.StringAttribute{
                        Computed:            true,
                        Description:         "Budget id",
                        MarkdownDescription: "The unique ID of the budget.\n\nExample: `bud-1234567890abcdef`",
                    },
                    common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
                        Computed:            true,
                        Description:         "Modified datetime",
                        MarkdownDescription: "The datetime when the budget was last modified.\n\nExample: `2024-01-15T00:00:00`",
                    },
                    common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
                        Computed:            true,
                        Description:         "Modified user",
                        MarkdownDescription: "The user who last modified the budget.\n\nExample: `user@example.com`",
                    },
                    common.ToSnakeCase("Name"): schema.StringAttribute{
                        Computed:            true,
                        Description:         "Budget name",
                        MarkdownDescription: "The name of the budget.\n\nExample: `Monthly Service Budget`",
                    },
                    common.ToSnakeCase("StartMonth"): schema.StringAttribute{
                        Computed:            true,
                        Description:         "Budget start month",
                        MarkdownDescription: "The month when the budget period starts.\n\nExample: `2024-01`",
                    },
                    common.ToSnakeCase("Type"): schema.StringAttribute{
                        Computed:            true,
                        Description:         "Budget type",
                        MarkdownDescription: "The type of the budget. Default type is  `COST`\n\nAllowed values: `COST`\n\nExample: `COST`",
                    },
                    common.ToSnakeCase("Unit"): schema.StringAttribute{
                        Computed:            true,
                        Description:         "Budget management unit",
                        MarkdownDescription: "The budget manage unit.\n\nAllowed values: `MONTHLY` | `OVERALL`\n\nExample: `MONTHLY`",
                    },
                },
                Computed: true,
            },
            common.ToSnakeCase("Name"): schema.StringAttribute{
                Description:        "Budget name ",
                MarkdownDescription: "The name of the budget.\n\nAllowed length: 20\n\nAllowed values: English, Korean, numbers and special characters such `as +=, .@, and -`\n\nExample: `bud-1234567890abcdef`",
                Optional:    true,
                Validators: []validator.String{
                    stringvalidator.LengthBetween(1, 64),
                },
            },
            common.ToSnakeCase("Amount"): schema.Int32Attribute{
                Description:        "Budget amount to set",
                MarkdownDescription: "The budget amount to set. Must be a positive numeric value.\n\nExample: `1000000`",
                Optional:    true,
            },
            common.ToSnakeCase("StartMonth"): schema.StringAttribute{
                Description:        "Start month of the budget period",
                MarkdownDescription: "The start month of the budget in YYYY-MM format.\n\nExample: `2024-01`",
                Optional:    true,
            },
            common.ToSnakeCase("Unit"): schema.StringAttribute{
                Description:        "Budget management unit",
                MarkdownDescription: "The budget manage unit.\n\nAllowed values: `MONTHLY` | `OVERALL`\n\nExample: `MONTHLY`",
                Optional:    true,
            },
            common.ToSnakeCase("Notifications"): schema.SingleNestedAttribute{
                Attributes: map[string]schema.Attribute{
                    common.ToSnakeCase("IsUseNotification"): schema.BoolAttribute{
                        Optional:            true,
                        Description:         "Notification use state",
                        MarkdownDescription: "Whether to enable budget notifications.\n\nAllowed values:`true`|`false`\n\nExample: `true`",
                    },
                    common.ToSnakeCase("NotificationSendPeriod"): schema.StringAttribute{
                        Optional:            true,
                        Description:         "Notification send period",
                        MarkdownDescription: "When to send budget notifications.\n\nAllowed values: `FIRST` | `DAILY` | `NONE`\n\nExample: `FIRST`",
                    },
                    common.ToSnakeCase("Receivers"): schema.ListAttribute{
                        ElementType:         types.StringType,
                        Optional:            true,
                        Description:         "List of notification recipient email addresses",
                        MarkdownDescription: "Email addresses to receive budget notifications.\n\nExample: `[\"user@example.com\"]`",
                    },
                    common.ToSnakeCase("Thresholds"): schema.ListAttribute{
                        ElementType:         types.Int32Type,
                        Optional:            true,
                        Description:         "List of threshold percentages for notifications",
                        MarkdownDescription: "Percentage thresholds at which to send notifications.\n\nAllowed values: `70` | `80` | `90` | `100`\n\nExample: `[80, 100]`",
                    },
                },
                Optional:            true,
                Description:         "Notification settings for the budget",
                MarkdownDescription: "Settings for budget alerts when usage exceeds defined thresholds.",
            },
            common.ToSnakeCase("Prevention"): schema.SingleNestedAttribute{
                Attributes: map[string]schema.Attribute{
                    common.ToSnakeCase("IsUsePrevention"): schema.BoolAttribute{
                        Optional:            true,
                        Description:         "Auto Generation prevent use state",
                        MarkdownDescription: "Whether to prevent new resource creation when budget threshold is exceeded.\n\nAllowed values:`true`|`false`\n\nExample: `true`",
                    },
                    common.ToSnakeCase("Receivers"): schema.ListAttribute{
                        ElementType:         types.StringType,
                        Optional:            true,
                        Description:         "List of notification recipient email addresses",
                        MarkdownDescription: "Email addresses to notify when budget prevention is triggered.\n\nExample: `[\"admin@example.com\"]`",
                    },
                    common.ToSnakeCase("Threshold"): schema.Int32Attribute{
                        Optional:            true,
                        Description:         "Prevention threshold",
                        MarkdownDescription: "Budget threshold percentage to trigger prevention.\n\nAllowed values: One of the numbers `70` | `80` | `90` | `100`\n\nExample: `90`",
                    },
                },
                Optional:            true,
                Description:         "Auto generation prevention settings for the budget",
                MarkdownDescription: "Settings to prevent new resource generation when budget usage exceeds the threshold.",
            },
        },
    }
}

// [BUDGET-FIX-02] 상태 폴링/재시도 로직을 위한 waitForBudgetReady 추가
func (r *budgetBudgetResource) waitForBudgetReady(ctx context.Context, id string, timeout time.Duration) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return fmt.Errorf("timeout waiting for budget %s to be ready", id)
		case <-ticker.C:
			data, err := r.client.GetAccountBudget(ctx, id)
			if err != nil {
				// 404 Not Found - 생성 중일 수 있음, 계속 대기
				if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
					continue
				}
				// 일시적 오류는 무시하고 계속 폴링 (네트워크 문제 등)
				if client.IsTransientError(err) {
					continue
				}
				// 그 외의 영구적인 오류는 반환
				return err
			}

			// 상태 확인 - 실제 API 응답에 맞게 수정 필요
			budget := data.Budget
			if budget.Id != "" {
				return nil
			}
		}
	}
}


