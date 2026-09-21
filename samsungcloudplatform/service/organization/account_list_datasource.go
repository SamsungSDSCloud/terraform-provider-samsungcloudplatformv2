package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/organization/1.3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &accountListDataSource{}
	_ datasource.DataSourceWithConfigure = &accountListDataSource{}
)

func NewAccountListDataSource() datasource.DataSource {
	return &accountListDataSource{}
}

type accountListDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (d *accountListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_accounts"
}

func (d *accountListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Organization
	d.clients = inst.Client
}

func (d *accountListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List Organization Accounts",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Optional: true,
				Description: "Filter by Organization ID. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
			},
			"size": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page size. \n" +
					"  - example : 20 \n",
			},
			"page": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page number. \n" +
					"  - example : 0 \n",
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort criteria. \n" +
					"  - example : 'created_at:desc' \n",
			},
			"id": schema.StringAttribute{
				Optional: true,
				Description: "Filter by Account ID. \n" +
					"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Description: "Filter by Account Name. \n" +
					"  - example : 'score-account' \n",
			},
			"email": schema.StringAttribute{
				Optional: true,
				Description: "Filter by Account Email. \n" +
					"  - example : 'score@samsung.com' \n",
			},
			"login_id": schema.StringAttribute{
				Optional: true,
				Description: "Filter by Login ID. \n" +
					"  - example : 'log-archive@samsung.com' \n",
			},
			"joined_start_date": schema.StringAttribute{
				Optional: true,
				Description: "Filter by joined start date. \n" +
					"  - example : '2026-04-11T12:12:12.123Z' \n",
			},
			"joined_end_date": schema.StringAttribute{
				Optional: true,
				Description: "Filter by joined end date. \n" +
					"  - example : '2026-04-11T12:12:12.123Z' \n",
			},
			"joined_method": schema.StringAttribute{
				Optional: true,
				Description: "Filter by joined method. \n" +
					"  - example : 'INVITED' \n" +
					"  - allowed_values : ['INVITED', 'CREATED', 'CLOUD_CONTROL'] \n",
			},
			"exclude_policy_id": schema.StringAttribute{
				Optional: true,
				Description: "Policy ID to Exclude. \n" +
					"  - example : '238c2fc8c29a449dbfa8681f8f1d78e2' \n",
			},
			"parent_unit_id": schema.StringAttribute{
				Optional: true,
				Description: "Filter by parent unit ID. \n" +
					"  - example : 'ou-fc8c29a138d78e24bf1fa86812fc8b' \n",
			},
			"parent_unit_name": schema.StringAttribute{
				Optional: true,
				Description: "Filter by parent unit name. \n" +
					"  - example : 'parent-unit-name' \n",
			},
			"type": schema.StringAttribute{
				Optional: true,
				Description: "Filter by account type. \n" +
					"  - example : 'MEMBER' \n" +
					"  - allowed_values : ['MANAGEMENT', 'DELEGATION', 'MEMBER', 'NONE', 'LOG_ARCHIVE', 'AUDIT'] \n",
			},
			"accounts": schema.ListNestedAttribute{
				Computed: true,
				Description: "Organization Account List. \n" +
					"  - example : '[{id: 0a36e0746dbf4908acf0357829701381, name: score-account, ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the account was created. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who created the account. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"email": schema.StringAttribute{
							Computed: true,
							Description: "Account Email. \n" +
								"  - example : 'score@samsung.com' \n",
						},
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the account in the organization. \n" +
								"  - example : '0a36e0746dbf4908acf0357829701381' \n",
						},
						"joined_method": schema.StringAttribute{
							Computed: true,
							Description: "Method by which the account joined the organization. \n" +
								"  - example : 'INVITED' \n" +
								"  - allowed_values : ['INVITED', 'CREATED', 'CLOUD_CONTROL'] \n",
						},
						"joined_time": schema.StringAttribute{
							Computed: true,
							Description: "Joined Datetime. \n" +
								"  - example : '2024-05-17T12:34:56.789Z' \n",
						},
						"login_id": schema.StringAttribute{
							Computed: true,
							Description: "Login ID of the account. \n" +
								"  - example : 'log-archive@samsung.com' \n",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the account was modified. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who modified the account. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"name": schema.StringAttribute{
							Computed: true,
							Description: "Account Name. \n" +
								"  - example : 'score-account' \n",
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the organization. \n" +
								"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
						},
						"parent_unit_id": schema.StringAttribute{
							Computed: true,
							Description: "Root or Parent Organization Unit ID. \n" +
								"  - example : 'ou-fc8c29a138d78e24bf1fa86812fc8b' \n",
						},
						"parent_unit_name": schema.StringAttribute{
							Computed: true,
							Description: "Root or Parent Organization Unit Name. \n" +
								"  - example : 'parent-unit-name' \n",
						},
						"state": schema.StringAttribute{
							Computed: true,
							Description: "Account Status. \n" +
								"  - example : 'ACTIVE' \n" +
								"  - allowed_values : ['ACTIVE', 'SUSPENDED', 'PENDING_CLOSURE'] \n",
						},
						"type": schema.StringAttribute{
							Computed: true,
							Description: "Account Type. \n" +
								"  - example : 'MEMBER' \n" +
								"  - allowed_values : ['MANAGEMENT', 'DELEGATION', 'MEMBER', 'NONE', 'LOG_ARCHIVE', 'AUDIT'] \n",
						},
					},
				},
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
				Description: "Total account count. \n" +
					"  - example : 5 \n",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort criteria from response. \n" +
					"  - example : ['created_at:desc'] \n",
			},
		},
	}
}

func (d *accountListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organization.AccountListDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiRequest := organization.AccountListDataSourceRequest{
		OrganizationId:  state.OrganizationId,
		Size:            state.Size,
		Page:            state.Page,
		Sort:            state.Sort,
		Id:              state.Id,
		Name:            state.Name,
		Email:           state.Email,
		LoginId:         state.LoginId,
		JoinedStartDate: state.JoinedStartDate,
		JoinedEndDate:   state.JoinedEndDate,
		JoinedMethod:    state.JoinedMethod,
		ExcludePolicyId: state.ExcludePolicyId,
		ParentUnitId:    state.ParentUnitId,
		ParentUnitName:  state.ParentUnitName,
		Type:            state.Type,
	}

	data, err := d.client.ListAccounts(ctx, apiRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to List Organization Accounts",
			err.Error(),
		)
		return
	}

	accounts := make([]attr.Value, 0)
	for _, account := range data.Accounts {
		accountValue := d.buildAccountValueFromSummary(ctx, &account)
		accounts = append(accounts, accountValue)
	}

	accountsList, accountsListDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.AccountListValue{}.AttributeTypes(ctx)},
		accounts,
	)
	resp.Diagnostics.Append(accountsListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sortListVal, diags := types.ListValueFrom(ctx, types.StringType, data.Sort)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Accounts = accountsList
	state.TotalCount = types.Int64Value(int64(data.Count))
	state.Page = types.Int64Value(int64(data.Page))
	state.Size = types.Int64Value(int64(data.Size))
	state.SortResult = sortListVal

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *accountListDataSource) buildAccountValueFromSummary(ctx context.Context, account *sdkorganization.AccountSummary) types.Object {
	var email, loginId, name, parentUnitName string
	if account.Email.IsSet() {
		email = *account.Email.Get()
	}
	if account.LoginId.IsSet() {
		loginId = *account.LoginId.Get()
	}
	if account.Name.IsSet() {
		name = *account.Name.Get()
	}
	if account.ParentUnitName.IsSet() && account.ParentUnitName.Get() != nil {
		parentUnitName = *account.ParentUnitName.Get()
	}

	accountValue, _ := types.ObjectValue(organization.AccountListValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"created_at":       types.StringValue(account.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
		"created_by":       types.StringValue(account.CreatedBy),
		"email":            types.StringValue(email),
		"id":               types.StringValue(account.Id),
		"joined_method":    types.StringValue(string(account.JoinedMethod)),
		"joined_time":      types.StringValue(account.JoinedTime.Format("2006-01-02T15:04:05.000Z")),
		"login_id":         types.StringValue(loginId),
		"modified_at":      types.StringValue(account.ModifiedAt.Format("2006-01-02T15:04:05.000Z")),
		"modified_by":      types.StringValue(account.ModifiedBy),
		"name":             types.StringValue(name),
		"organization_id":  types.StringValue(account.OrganizationId),
		"parent_unit_id":   types.StringValue(account.ParentUnitId),
		"parent_unit_name": types.StringValue(parentUnitName),
		"state":            types.StringValue(string(account.State)),
		"type":             types.StringValue(string(account.Type)),
	})

	return accountValue
}

func (d *accountListDataSource) buildAccountValue(ctx context.Context, account *sdkorganization.OrganizationAccountWithPolicy) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var controlPolicies []attr.Value
	for _, policy := range account.ControlPolicies {
		policyValue, policyDiags := types.ObjectValue(organization.ControlPoliciesValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"policy_id":   types.StringValue(policy.PolicyId),
			"policy_name": types.StringValue(policy.PolicyName),
		})
		diags.Append(policyDiags...)
		if diags.HasError() {
			return types.ObjectNull(organization.AccountValue{}.AttributeTypes(ctx)), diags
		}
		controlPolicies = append(controlPolicies, policyValue)
	}

	controlPoliciesList, listDiags := types.ListValue(types.ObjectType{AttrTypes: organization.ControlPoliciesValue{}.AttributeTypes(ctx)}, controlPolicies)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.AccountValue{}.AttributeTypes(ctx)), diags
	}

	var email, loginId, name, srn, parentUnitName string
	if account.Email.IsSet() {
		email = *account.Email.Get()
	}
	if account.LoginId.IsSet() {
		loginId = *account.LoginId.Get()
	}
	if account.Name.IsSet() {
		name = *account.Name.Get()
	}
	if account.Srn.IsSet() {
		srn = *account.Srn.Get()
	}
	if account.ParentUnitName.IsSet() && account.ParentUnitName.Get() != nil {
		parentUnitName = *account.ParentUnitName.Get()
	}

	accountValue, accountDiags := types.ObjectValue(organization.AccountValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"control_policies": controlPoliciesList,
		"created_at":       types.StringValue(account.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
		"created_by":       types.StringValue(account.CreatedBy),
		"creator_name":     types.StringValue(account.GetCreatorName()),
		"email":            types.StringValue(email),
		"id":               types.StringValue(account.Id),
		"joined_method":    types.StringValue(string(account.JoinedMethod)),
		"joined_time":      types.StringValue(account.JoinedTime.Format("2006-01-02T15:04:05.000Z")),
		"login_id":         types.StringValue(loginId),
		"modified_at":      types.StringValue(account.ModifiedAt.Format("2006-01-02T15:04:05.000Z")),
		"modified_by":      types.StringValue(account.ModifiedBy),
		"modifier_name":    types.StringValue(account.GetModifierName()),
		"name":             types.StringValue(name),
		"organization_id":  types.StringValue(account.OrganizationId),
		"parent_unit_id":   types.StringValue(account.ParentUnitId),
		"parent_unit_name": types.StringValue(parentUnitName),
		"srn":              types.StringValue(srn),
		"state":            types.StringValue(string(account.State)),
		"type":             types.StringValue(string(account.Type)),
	})
	diags.Append(accountDiags...)
	return accountValue, diags
}
