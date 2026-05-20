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
				Required:            true,
				Description:         "Organization Account ID",
				MarkdownDescription: "Organization Account ID",
			},
			"account": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Organization Account Detail Info",
				MarkdownDescription: "Organization Account Detail Info",
				Attributes: map[string]schema.Attribute{
					"control_policies": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"policy_id": schema.StringAttribute{
									Computed:            true,
									Description:         "Control Policy ID",
									MarkdownDescription: "Control Policy ID",
								},
								"policy_name": schema.StringAttribute{
									Computed:            true,
									Description:         "Control Policy Name",
									MarkdownDescription: "Control Policy Name",
								},
							},
						},
						Computed: true,
					},
					"created_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Created At",
						MarkdownDescription: "Created At",
					},
					"created_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Created By",
						MarkdownDescription: "Created By",
					},
					"creator_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Creator Name",
						MarkdownDescription: "Creator Name",
					},
					"email": schema.StringAttribute{
						Computed:            true,
						Description:         "Account Email",
						MarkdownDescription: "Account Email",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "ID",
						MarkdownDescription: "ID",
					},
					"joined_method": schema.StringAttribute{
						Computed:            true,
						Description:         "Joined Method",
						MarkdownDescription: "Joined Method",
					},
					"joined_time": schema.StringAttribute{
						Computed:            true,
						Description:         "Joined Datetime",
						MarkdownDescription: "Joined Datetime",
					},
					"login_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Login ID",
						MarkdownDescription: "Login ID",
					},
					"modified_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified At",
						MarkdownDescription: "Modified At",
					},
					"modified_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified By",
						MarkdownDescription: "Modified By",
					},
					"modifier_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Modifier Name",
						MarkdownDescription: "Modifier Name",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Account Name",
						MarkdownDescription: "Account Name",
					},
					"organization_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Organization ID",
						MarkdownDescription: "Organization ID",
					},
					"parent_unit_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Root or Parent Organization Unit ID",
						MarkdownDescription: "Root or Parent Organization Unit ID",
					},
					"parent_unit_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Root or Parent Organization Unit Name",
						MarkdownDescription: "Root or Parent Organization Unit Name",
					},
					"srn": schema.StringAttribute{
						Computed:            true,
						Description:         "SRN",
						MarkdownDescription: "SRN",
					},
					"state": schema.StringAttribute{
						Computed:            true,
						Description:         "Account State",
						MarkdownDescription: "Account State",
					},
					"type": schema.StringAttribute{
						Computed:            true,
						Description:         "Account Type",
						MarkdownDescription: "Account Type",
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

	accountValue := d.buildAccountValue(ctx, &data.Account)
	state.Account = accountValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *accountDataSource) buildAccountValue(ctx context.Context, account *sdkorganization.OrganizationAccountWithPolicy) types.Object {
	var controlPolicies []attr.Value
	for _, policy := range account.ControlPolicies {
		policyValue, _ := types.ObjectValue(organization.ControlPoliciesValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"policy_id":   types.StringValue(policy.PolicyId),
			"policy_name": types.StringValue(policy.PolicyName),
		})
		controlPolicies = append(controlPolicies, policyValue)
	}

	controlPoliciesList := types.ListValueMust(types.ObjectType{AttrTypes: organization.ControlPoliciesValue{}.AttributeTypes(ctx)}, controlPolicies)

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
