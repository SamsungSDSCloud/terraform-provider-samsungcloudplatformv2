package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
				Optional:    true,
				Description: "Filter by Organization ID",
			},
			"size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Page size",
			},
			"page": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Page number",
			},
			"sort": schema.StringAttribute{
				Optional:    true,
				Description: "Sort criteria (e.g., 'created_at:desc')",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by Account ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by Account Name",
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by Account Email",
			},
			"login_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by Login ID",
			},
			"joined_start_date": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by joined start date (e.g., '2026-04-11T12:12:12.123Z')",
			},
			"joined_end_date": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by joined end date (e.g., '2026-04-11T12:12:12.123Z')",
			},
			"joined_method": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by joined method",
			},
			"exclude_policy_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by exclude policy ID",
			},
			"parent_unit_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by parent unit ID",
			},
			"parent_unit_name": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by parent unit name",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by account type",
			},
			"accounts": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"created_by": schema.StringAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"joined_method": schema.StringAttribute{
							Computed: true,
						},
						"joined_time": schema.StringAttribute{
							Computed: true,
						},
						"login_id": schema.StringAttribute{
							Computed: true,
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
						},
						"parent_unit_id": schema.StringAttribute{
							Computed: true,
						},
						"parent_unit_name": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort criteria from response",
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

	accountsList := types.ListValueMust(
		types.ObjectType{AttrTypes: organization.AccountListValue{}.AttributeTypes(ctx)},
		accounts,
	)

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

func (d *accountListDataSource) buildAccountValue(ctx context.Context, account *sdkorganization.OrganizationAccountWithPolicy) types.Object {
	var controlPolicies []attr.Value
	for _, policy := range account.ControlPolicies {
		policyValue, _ := types.ObjectValue(organization.ControlPoliciesValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"policy_id":   types.StringValue(policy.PolicyId),
			"policy_name": types.StringValue(policy.PolicyName),
		})
		controlPolicies = append(controlPolicies, policyValue)
	}

	controlPoliciesList := types.ListValueMust(
		types.ObjectType{AttrTypes: organization.ControlPoliciesValue{}.AttributeTypes(ctx)},
		controlPolicies,
	)

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

	accountValue, _ := types.ObjectValue(organization.AccountValue{}.AttributeTypes(ctx), map[string]attr.Value{
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

	return accountValue
}
