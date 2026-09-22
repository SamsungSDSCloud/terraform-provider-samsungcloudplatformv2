package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &delegationAccountsDataSource{}
	_ datasource.DataSourceWithConfigure = &delegationAccountsDataSource{}
)

func NewDelegationAccountsDataSource() datasource.DataSource {
	return &delegationAccountsDataSource{}
}

type delegationAccountsDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (d *delegationAccountsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_delegation_accounts"
}

func (d *delegationAccountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *delegationAccountsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get Delegation Account list",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
			},
			"service_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Service type of the delegated account. \n" +
					"  - allowed_values : ['identity-center', 'resource-optimizer', 'infrastructure-builder'] \n" +
					"  - example : 'identity-center' \n",
			},
			"account_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Delegation Account ID. \n" +
					"  - example : '0a7ca6b5693c45a28736e84346c6d6c5' \n",
			},
			"size": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page size. \n" +
					"  - example : '20' \n",
			},
			"page": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page number. \n" +
					"  - example : '0' \n",
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Sort criteria. \n" +
					"  - example : 'createdAt,desc' \n",
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
				Description: "Total number of delegation accounts matching the filter. \n" +
					"  - example : '10' \n",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort criteria. \n" +
					"  - example : '[\"created_at:asc\"]' \n",
			},
			"delegation_accounts": schema.ListNestedAttribute{
				Computed: true,
				Description: "Delegation Account list. \n" +
					"  - example : '[{account_id: b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0, created_at: 2025-01-01T00:00:00.000Z, created_by: c23fb561c689455993874fa5d5ed4a2f, id: 0a36e0746dbf4908acf0357829701381, organization_id: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5, service_type: identity-center}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed: true,
							Description: "Delegation Account ID. \n" +
								"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Created timestamp. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "Creator ID. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the delegation record. \n" +
								"  - example : '0a36e0746dbf4908acf0357829701381' \n",
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the organization. \n" +
								"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
						},
						"service_type": schema.StringAttribute{
							Computed: true,
							Description: "Service type of the delegated account. \n" +
								"  - allowed_values : ['identity-center', 'resource-optimizer', 'infrastructure-builder'] \n" +
								"  - example : 'identity-center' \n",
						},
					},
				},
			},
		},
	}
}

func (d *delegationAccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organization.DelegationAccountsDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := state.OrganizationId.ValueString()
	serviceType := state.ServiceType.ValueString()
	accountId := state.AccountId.ValueString()

	var size int32
	if !state.Size.IsNull() && !state.Size.IsUnknown() {
		size = int32(state.Size.ValueInt64())
	}

	var page int32
	if !state.Page.IsNull() && !state.Page.IsUnknown() {
		page = int32(state.Page.ValueInt64())
	}

	var sortStr string
	if !state.Sort.IsNull() && !state.Sort.IsUnknown() {
		sortStr = state.Sort.ValueString()
	}

	data, err := d.client.ListDelegationAccounts(ctx, orgId, serviceType, size, page, sortStr, accountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Delegation Accounts",
			err.Error(),
		)
		return
	}

	state.TotalCount = types.Int64Value(int64(data.Count))
	state.Page = types.Int64Value(int64(page))
	state.Size = types.Int64Value(int64(size))

	if data.Sort != nil && len(data.Sort) > 0 {
		sortListVal, sortDiags := types.ListValueFrom(ctx, types.StringType, data.Sort)
		resp.Diagnostics.Append(sortDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.SortResult = sortListVal
	}

	var accounts []attr.Value
	for _, item := range data.DelegationAccounts {
		accountValue, accDiags := types.ObjectValue(organization.DelegationAccountValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"account_id":      types.StringValue(item.AccountId),
			"created_at":      types.StringValue(item.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
			"created_by":      types.StringValue(item.CreatedBy),
			"id":              types.StringValue(item.Id),
			"organization_id": types.StringValue(item.OrganizationId),
			"service_type":    types.StringValue(item.ServiceType),
		})
		resp.Diagnostics.Append(accDiags...)
		accounts = append(accounts, accountValue)
	}

	accountsList, listDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.DelegationAccountValue{}.AttributeTypes(ctx)},
		accounts,
	)
	resp.Diagnostics.Append(listDiags...)
	state.DelegationAccounts = accountsList

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
