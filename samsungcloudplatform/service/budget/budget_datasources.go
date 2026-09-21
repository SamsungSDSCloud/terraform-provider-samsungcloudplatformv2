package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/budget"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	budgetsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/budget/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &budgetBudgetDataSources{}
	_ datasource.DataSourceWithConfigure = &budgetBudgetDataSources{}
)

// NewBudgetDataSources is a helper function to simplify the provider implementation.
func NewBudgetBudgetDataSources() datasource.DataSource {
	return &budgetBudgetDataSources{}
}

// budgetDataSource is the data source implementation.
type budgetBudgetDataSources struct {
	config  *scpsdk.Configuration
	client  *budget.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *budgetBudgetDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_budget_budgets"
}

type BudgetListItem struct {
	Amount       types.String `tfsdk:"amount"`
	CreatedAt    types.String `tfsdk:"created_at"`
	CreatedBy    types.String `tfsdk:"created_by"`
	Id           types.String `tfsdk:"id"`
	ModifiedAt   types.String `tfsdk:"modified_at"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
	Name         types.String `tfsdk:"name"`
	StartMonth   types.String `tfsdk:"start_month"`
	Type         types.String `tfsdk:"type"`
	Unit         types.String `tfsdk:"unit"`
	IsCostLinked types.Bool   `tfsdk:"is_cost_linked"`
}

type BudgetListDataSourceModel struct {
	Id          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	Filter  []filter.Filter `tfsdk:"filter"`
	Ids     []types.String  `tfsdk:"ids"`
	Budgets []BudgetListItem `tfsdk:"budgets"`
	BudgetCount   types.Int64     `tfsdk:"budget_count"`
	Page    types.Int64     `tfsdk:"page"`
	Size    types.Int64     `tfsdk:"size"`
	Sort    []types.String  `tfsdk:"sort"`
}

// Schema defines the schema for the data source.
func (d *budgetBudgetDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "list of account budget",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description:         "Unique ID of the budget.",
				MarkdownDescription: "The unique ID of the budget.\n\nExample: `bud-1234567890abcdef`",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description:         "Budget name",
				MarkdownDescription: "The name of the budget.\n\nExample: `ex_month_budget`",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 20),
				},
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "Budget ID List",
				MarkdownDescription: "The id list of budget.\n\nExample: `[\"bud-1234567890abcdef\"]`",
				Computed:            true,
			},
			// start add
			"budgets": schema.ListNestedAttribute{
				Computed: true,
				Description: "Budget list",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"amount":          schema.StringAttribute{Computed: true},
						"created_at":      schema.StringAttribute{Computed: true},
						"created_by":      schema.StringAttribute{Computed: true},
						"id":              schema.StringAttribute{Computed: true},
						"modified_at":     schema.StringAttribute{Computed: true},
						"modified_by":     schema.StringAttribute{Computed: true},
						"name":            schema.StringAttribute{Computed: true},
						"start_month":     schema.StringAttribute{Computed: true},
						"type":            schema.StringAttribute{Computed: true},
						"unit":            schema.StringAttribute{Computed: true},
						"is_cost_linked":  schema.BoolAttribute{Computed: true},
					},
				},
			},
			"budget_count": schema.Int64Attribute{
				Computed: true,
				Description: "Total count of budgets",
			},
			"page": schema.Int64Attribute{
				Computed: true,
				Description: "Current page number",
			},
			"size": schema.Int64Attribute{
				Computed: true,
				Description: "Page size",
			},
			"sort": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Sort criteria",
			},
		},
		// --end add
		Blocks: map[string]schema.Block{ // 필터는 Block 으로 정의한다.
			"filter": filter.DataSourceSchema(), // 필터 스키마는 공통으로 제공되는 함수를 이용하여 정의한다.
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *budgetBudgetDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func stringPtrToTypesString(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}

func timePtrToTypesString(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(time.RFC3339))
}

func nullableBoolToTypesBool(nullable budgetsdk.NullableBool) types.Bool {
	if !nullable.IsSet() {
		return types.BoolNull()
	}

	ptr := nullable.Get() // *bool
	if ptr == nil {
		return types.BoolNull()
	}

	return types.BoolValue(*ptr)
}

func (d *budgetBudgetDataSources) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var state BudgetListDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetAccountBudgetList(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Account Budget",
			err.Error(),
		)
		return
	}

	contents := data.Budgets

	if len(state.Filter) > 0 {
		indices, err := filter.GetFilterIndices(contents, state.Filter)
		if err != nil {
			resp.Diagnostics.AddError("Filter Error", err.Error())
			return
		}

		filteredContents := contents[:0]
		for i, resource := range contents {
			if common.Contains(indices, i) {
				filteredContents = append(filteredContents, resource)
			}
		}
		contents = filteredContents
	}

	ids := make([]types.String, 0, len(contents))
	budgets := make([]BudgetListItem, 0, len(contents))

	for _, content := range contents {
		id := content.Id

		ids = append(ids, types.StringValue(id))

		budgets = append(budgets, BudgetListItem{
			Amount:       types.StringValue(content.Amount),
			CreatedAt:    timePtrToTypesString(content.CreatedAt),
			CreatedBy:    stringPtrToTypesString(content.CreatedBy),
			Id:           types.StringValue(id),
			ModifiedAt:   timePtrToTypesString(content.ModifiedAt),
			ModifiedBy:   stringPtrToTypesString(content.ModifiedBy),
			Name:         types.StringValue(content.Name),
			StartMonth:   types.StringValue(content.StartMonth),
			Type:         types.StringValue(content.Type),
			Unit:         types.StringValue(content.Unit),
			IsCostLinked: nullableBoolToTypesBool(content.IsCostLinked),
		})
	}

	state.Ids = ids
	state.Budgets = budgets
	state.BudgetCount = types.Int64Value(int64(data.Count))
	state.Page = types.Int64Value(int64(data.Page))
	state.Size = types.Int64Value(int64(data.Size))

	sortValues := make([]types.String, 0, len(data.Sort))
	for _, sortItem := range data.Sort {
		sortValues = append(sortValues, types.StringValue(sortItem))
	}
	state.Sort = sortValues

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}