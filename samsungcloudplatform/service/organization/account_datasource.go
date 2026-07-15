package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &accountDataSource{}
	_ datasource.DataSourceWithConfigure = &accountDataSource{}
)

func NewAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

type accountDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (d *accountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_account"
}

func (d *accountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Organization Account",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Required: true,
				Description: "Organization Account ID. \n" +
					"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
			},
			"account": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Organization Account Detail Info. \n" +
					"  - example : '{id: 0a36e0746dbf4908acf0357829701381, name: score-account, ...}' \n",
				Attributes: map[string]schema.Attribute{
					"control_policies": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"policy_id": schema.StringAttribute{
									Computed: true,
									Description: "Control Policy ID. \n" +
										"  - example : 'f98e76d54c32b10a9z8y7x6w5v4u3' \n",
								},
								"policy_name": schema.StringAttribute{
									Computed: true,
									Description: "Control Policy Name. \n" +
										"  - example : 'test-policy-name' \n",
								},
							},
						},
						Computed: true,
						Description: "Service control policies. \n" +
							"  - example : '[{policy_id: f98e76d54c32b10a9z8y7x6w5v4u3, policy_name: test-policy-name}]' \n",
					},
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
					"creator_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the account creator. \n" +
							"  - example : 'John Doe na' \n",
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
					"modifier_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the account modifier. \n" +
							"  - example : 'Alice' \n",
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
					"srn": schema.StringAttribute{
						Computed: true,
						Description: "Samsung Resource Name (SRN) uniquely identifying this resource. \n" +
							"  - example : 'srn:dev2::1b8c29a138d78e24bf1fa86812fcaa:kr-west1::organizations/account/b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
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
	}
}

func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organization.AccountDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountId := state.AccountId.ValueString()
	if accountId == "" {
		resp.Diagnostics.AddError(
			"Unable to Read Organization Account",
			"Account ID is empty",
		)
		return
	}

	data, err := d.client.GetAccount(ctx, accountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Organization Account",
			err.Error(),
		)
		return
	}

	accountValue, accountDiags := d.buildAccountValue(ctx, &data.Account)
	resp.Diagnostics.Append(accountDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Account = accountValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *accountDataSource) buildAccountValue(ctx context.Context, account *sdkorganization.OrganizationAccountWithPolicy) (types.Object, diag.Diagnostics) {
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
