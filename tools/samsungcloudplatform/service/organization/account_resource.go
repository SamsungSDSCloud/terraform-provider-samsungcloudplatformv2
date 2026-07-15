package organization

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	orgsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &accountResource{}
	_ resource.ResourceWithConfigure   = &accountResource{}
	_ resource.ResourceWithImportState = &accountResource{}
)

func NewAccountResource() resource.Resource {
	return &accountResource{}
}

type accountResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *accountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_account"
}

func (r *accountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages accounts within an organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Account ID (Terraform internal ID). \n" +
					"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
				Optional: true,
				Computed: true,
			},
			"account_id": schema.StringAttribute{
				Description: "Organization Account ID. \n" +
					"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
				Optional: true,
				Computed: true,
			},
			"account_ids": schema.ListAttribute{
				Description: "Account IDs for batch delete. \n" +
					"  - example : ['b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0'] \n",
				Optional:    true,
				ElementType: types.StringType,
			},
			"login_id": schema.StringAttribute{
				Description: "Login ID of the account. \n" +
					"  - example : 'log-archive@samsung.com' \n",
				Optional: true,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "Account Name. \n" +
					"  - example : 'score-account' \n",
				Optional: true,
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Optional: true,
				Computed: true,
			},
			"role_name": schema.StringAttribute{
				Description: "Name of the role created in the account. \n" +
					"  - example : 'OrganizationAccountAccessRole' \n",
				Optional: true,
				Computed: true,
			},
			"parent_unit_id": schema.StringAttribute{
				Description: "Parent Organization Unit ID (for update/move). \n" +
					"  - example : 'ou-fc8c29a138d78e24bf1fa86812fc8b' \n",
				Optional: true,
			},
			"lazy_policy": schema.BoolAttribute{
				Description: "Linked Policy Query YN. \n" +
					"  - example : true \n",
				Optional: true,
				Computed: true,
			},
			"account": schema.SingleNestedAttribute{
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
						Description: "Service Control Policies. \n" +
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
				Computed: true,
				Description: "Organization Account Detail Info. \n" +
					"  - example : '{id: 0a36e0746dbf4908acf0357829701381, name: score-account, ...}' \n",
			},
			"success": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"success_id": schema.StringAttribute{
						Computed: true,
						Description: "ID of the account successfully created. \n" +
							"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
					},
					"success_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the account successfully created. \n" +
							"  - example : 'score-account' \n",
					},
				},
				Computed: true,
				Description: "Success Info. \n" +
					"  - example : '{success_id: b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0, success_name: score-account}' \n",
			},
			"failed": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"error_code": schema.StringAttribute{
						Computed: true,
						Description: "Error code returned when the operation fails. \n" +
							"  - example : 'Organization.AccountNotRemovable' \n",
					},
					"failed_caused": schema.StringAttribute{
						Computed: true,
						Description: "Failure Reason. \n" +
							"  - example : 'Account b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0 is not found' \n",
					},
					"response": schema.StringAttribute{
						Computed: true,
						Description: "Quota Response. \n" +
							"  - example : {} \n",
					},
				},
				Computed: true,
				Description: "Failed Info. \n" +
					"  - example : '{error_code: Organization.AccountNotRemovable, failed_caused: Account ... is not found, ...}' \n",
			},
		},
	}
}

func (r *accountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.Organization
	r.clients = inst.Client
}

func (r *accountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organization.AccountResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := plan.OrganizationId.ValueString()
	loginId := plan.LoginId.ValueString()
	name := plan.Name.ValueString()
	roleName := plan.RoleName.ValueString()

	data, err := r.client.CreateAccount(ctx, loginId, name, orgId, roleName)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Organization Account",
			"Could not create Organization Account, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var accountValue types.Object
	successId := ""
	successName := ""
	failedErrorCode := ""
	failedCaused := ""
	failedResponse := ""

	// Use HasSuccess() and GetSuccess() for nullable type
	if data.HasSuccess() {
		successData := data.GetSuccess()
		successId = successData.SuccessId
		successName = successData.SuccessName

		accountResp, err := r.client.GetAccount(ctx, successId)
		if err == nil {
			accountValue, _ = r.buildAccountValue(ctx, &accountResp.Account)
		} else {
			accountValue = types.ObjectNull(organization.AccountValue{}.AttributeTypes(ctx))
		}
	} else if data.HasFailed() {
		failedData := data.GetFailed()
		failedErrorCode = failedData.ErrorCode
		failedCaused = failedData.FailedCaused
		if failedData.Response != nil {
			responseBytes, _ := json.Marshal(failedData.Response)
			failedResponse = string(responseBytes)
		}
		accountValue = types.ObjectNull(organization.AccountValue{}.AttributeTypes(ctx))
	} else {
		accountValue = types.ObjectNull(organization.AccountValue{}.AttributeTypes(ctx))
	}

	successValue, diags := types.ObjectValue(organization.SuccessValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"success_id":   types.StringValue(successId),
		"success_name": types.StringValue(successName),
	})
	resp.Diagnostics.Append(diags...)

	failedResponseValue := failedResponse
	if failedResponseValue == "" {
		failedResponseValue = "{}"
	}

	failedValue, diags := types.ObjectValue(organization.FailedValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"error_code":    types.StringValue(failedErrorCode),
		"failed_caused": types.StringValue(failedCaused),
		"response":      types.StringValue(failedResponseValue),
	})
	resp.Diagnostics.Append(diags...)

	plan.Account = accountValue
	plan.Success = successValue
	plan.Failed = failedValue

	if successId != "" {
		plan.Id = types.StringValue(successId)
		plan.AccountId = types.StringValue(successId)

		err = waitForAccountActiveStatus(ctx, r.client, successId)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error waiting for account creation",
				"Error waiting for account "+successId+" to become active: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *accountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organization.AccountResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountId := state.Id.ValueString()
	if accountId == "" {
		resp.Diagnostics.AddError(
			"Unable to Read Organization Account",
			"Account ID is empty",
		)
		return
	}

	orgId := state.OrganizationId.ValueString()
	if orgId == "" {
		orgId = "default"
	}

	data, err := r.client.GetAccount(ctx, accountId)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	accountValue, diags := r.buildAccountValue(ctx, &data.Account)
	resp.Diagnostics.Append(diags...)
	state.Account = accountValue
	state.AccountId = types.StringValue(accountId)
	state.LoginId = types.StringValue(data.Account.GetLoginId())
	state.Name = types.StringValue(data.Account.GetName())
	state.LazyPolicy = types.BoolValue(false)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *accountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan organization.AccountResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountId := plan.Id.ValueString()
	parentUnitId := plan.ParentUnitId.ValueString()
	orgId := plan.OrganizationId.ValueString()

	if parentUnitId == "" {
		resp.Diagnostics.AddError(
			"Unable to Update Organization Account",
			"Parent Unit ID is required for update/move",
		)
		return
	}

	_, err := r.client.MoveAccount(ctx, accountId, parentUnitId, orgId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error moving Organization Account",
			"Could not move Organization Account, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetAccount(ctx, accountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Organization Account after move",
			err.Error(),
		)
		return
	}

	accountValue, diags := r.buildAccountValue(ctx, &data.Account)
	resp.Diagnostics.Append(diags...)
	plan.Account = accountValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *accountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	accountId := req.ID
	if accountId == "" {
		resp.Diagnostics.AddError(
			"Invalid import identifier",
			"The import ID cannot be empty",
		)
		return
	}

	state := organization.AccountResource{
		Id:         types.StringValue(accountId),
		AccountId:  types.StringValue(accountId),
		LoginId:    types.StringValue(""),
		Name:       types.StringValue(""),
		RoleName:   types.StringValue(""),
		LazyPolicy: types.BoolValue(false),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *accountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organization.AccountResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountId := state.Id.ValueString()
	accountIds := state.AccountIds
	orgId := state.OrganizationId.ValueString()

	if !accountIds.IsNull() && !accountIds.IsUnknown() {
		ids := make([]string, 0, 10)
		for _, elem := range accountIds.Elements() {
			if id, ok := elem.(types.String); ok {
				ids = append(ids, id.ValueString())
			}
		}

		if len(ids) > 0 {
			if orgId == "" {
				orgId = "default"
			}
			_, err := r.client.RemoveAccounts(ctx, ids, orgId)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error deleting Organization Accounts",
					"Could not delete Organization Accounts, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}
	} else if accountId != "" {
		_, err := r.client.DeleteAccount(ctx, accountId)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error deleting Organization Account",
				"Could not delete Organization Account, unexpected error: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Unable to Delete Organization Account",
			"Account ID or Account IDs is required in state or config",
		)
		return
	}
}

func (r *accountResource) buildAccountValue(ctx context.Context, account *orgsdk.OrganizationAccountWithPolicy) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var controlPolicies []attr.Value
	for _, policy := range account.ControlPolicies {
		policyValue, policyDiags := types.ObjectValue(organization.ControlPoliciesValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"policy_id":   types.StringValue(policy.PolicyId),
			"policy_name": types.StringValue(policy.PolicyName),
		})
		diags.Append(policyDiags...)
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

func waitForAccountActiveStatus(ctx context.Context, orgClient *organization.Client, accountId string) error {
	return client.WaitForStatus(ctx, nil, []string{"PENDING"}, []string{"ACTIVE", "RUNNING"}, func() (interface{}, string, error) {
		accountResp, err := orgClient.GetAccount(ctx, accountId)
		if err != nil {
			return nil, "", err
		}
		return accountResp, string(accountResp.Account.State), nil
	}, -1, -1, -1, -1)
}
