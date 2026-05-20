package organization

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	orgsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &accountResource{}
	_ resource.ResourceWithConfigure = &accountResource{}
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
		Description: "Organization Account",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Account ID",
				Optional:    true,
				Computed:    true,
			},
			"account_id": schema.StringAttribute{
				Description: "Account ID (Alias)",
				Optional:    true,
				Computed:    true,
			},
			"account_ids": schema.ListAttribute{
				Description: "Account IDs for batch delete",
				Optional:    true,
				ElementType: types.StringType,
			},
			"login_id": schema.StringAttribute{
				Description: "Login ID",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Account Name",
				Optional:    true,
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "Organization ID",
				Optional:    true,
				Computed:    true,
			},
			"role_name": schema.StringAttribute{
				Description: "Role Name",
				Optional:    true,
				Computed:    true,
			},
			"parent_unit_id": schema.StringAttribute{
				Description: "Parent Organization Unit ID (for update/move)",
				Optional:    true,
			},
			"lazy_policy": schema.BoolAttribute{
				Description: "Linked Policy Query YN",
				Optional:    true,
				Computed:    true,
			},
			"account": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"control_policies": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"policy_id": schema.StringAttribute{
									Computed: true,
								},
								"policy_name": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"created_at": schema.StringAttribute{
						Computed: true,
					},
					"created_by": schema.StringAttribute{
						Computed: true,
					},
					"creator_name": schema.StringAttribute{
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
					"modifier_name": schema.StringAttribute{
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
					"srn": schema.StringAttribute{
						Computed: true,
					},
					"state": schema.StringAttribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
					},
				},
				Computed: true,
			},
			"success": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"success_id": schema.StringAttribute{
						Computed: true,
					},
					"success_name": schema.StringAttribute{
						Computed: true,
					},
				},
				Computed: true,
			},
			"failed": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"error_code": schema.StringAttribute{
						Computed: true,
					},
					"failed_caused": schema.StringAttribute{
						Computed: true,
					},
					"response": schema.StringAttribute{
						Computed: true,
					},
				},
				Computed: true,
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
			accountValue = r.buildAccountValue(ctx, &accountResp.Account)
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
		resp.Diagnostics.AddError(
			"Unable to Read Organization Account",
			err.Error(),
		)
		return
	}

	accountValue := r.buildAccountValue(ctx, &data.Account)
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

	accountValue := r.buildAccountValue(ctx, &data.Account)
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

func (r *accountResource) buildAccountValue(ctx context.Context, account *orgsdk.OrganizationAccountWithPolicy) types.Object {
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
