package organization

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/organization"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var successIdAttrTypes = map[string]attr.Type{
	"success_id":   basetypes.StringType{},
	"success_name": basetypes.StringType{},
}

var failedIdAttrTypes = map[string]attr.Type{
	"error_code":    basetypes.StringType{},
	"failed_caused": basetypes.StringType{},
	"failed_id":     basetypes.StringType{},
	"failed_name":   basetypes.StringType{},
	"response":      types.MapType{ElemType: types.StringType},
}

var (
	_ resource.Resource                = &accountRemoveResource{}
	_ resource.ResourceWithConfigure   = &accountRemoveResource{}
	_ resource.ResourceWithImportState = &accountRemoveResource{}
)

func NewAccountRemoveResource() resource.Resource {
	return &accountRemoveResource{}
}

type accountRemoveResource struct {
	config  interface{}
	client  *organization.Client
	clients *client.SCPClient
}

func (r *accountRemoveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_account_remove"
}

func (r *accountRemoveResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Remove Organization Accounts - Delete accounts by IDs",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Required: true,
			},
			"account_id": schema.StringAttribute{
				Description: "Account ID to remove (single). \n" +
					"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
				Optional: true,
			},
			"target_account_ids": schema.ListAttribute{
				Description: "Account IDs to remove (multiple). \n" +
					"  - example : ['b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0'] \n",
				ElementType: types.StringType,
				Optional:    true,
			},
			"success_ids": schema.ListNestedAttribute{
				Description: "Successfully removed accounts. \n" +
					"  - example : '[{success_id: b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0, success_name: score-account}]' \n",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"success_id": schema.StringAttribute{
							Computed: true,
							Description: "Organization Account ID. \n" +
								"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
						},
						"success_name": schema.StringAttribute{
							Computed: true,
							Description: "Account Name. \n" +
								"  - example : 'score-account' \n",
						},
					},
				},
			},
			"failed_ids": schema.ListNestedAttribute{
				Description: "Failed accounts. \n" +
					"  - example : '[{failed_id: 1a5xz23h1j0l9n8p7r6t5v4x3z2y1w0, error_code: Organization.AccountNotRemovable, ...}]' \n",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
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
						"failed_id": schema.StringAttribute{
							Computed: true,
							Description: "Organization Account ID. \n" +
								"  - example : '1a5xz23h1j0l9n8p7r6t5v4x3z2y1w0' \n",
						},
						"failed_name": schema.StringAttribute{
							Computed: true,
							Description: "Account Name. \n" +
								"  - example : 'score-account' \n",
						},
						"response": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Quota Response. \n" +
								"  - example : {} \n",
						},
					},
				},
			},
		},
	}
}

func (r *accountRemoveResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.Instance, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	r.client = inst.Client.Organization
	r.clients = inst.Client
}

func buildResponseMapValue(response map[string]interface{}) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(response) == 0 {
		return types.MapValue(types.StringType, map[string]attr.Value{})
	}
	responseMap := map[string]attr.Value{}
	for k, v := range response {
		responseMap[k] = types.StringValue(fmt.Sprintf("%v", v))
	}
	result, diags := types.MapValue(types.StringType, responseMap)
	return result, diags
}

func (r *accountRemoveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planData accountRemoveStateData
	diags := req.Plan.Get(ctx, &planData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := planData.OrganizationId.ValueString()
	accountIds := make([]string, 0, 10)

	if singleAccountId := planData.AccountId.ValueString(); singleAccountId != "" {
		accountIds = append(accountIds, singleAccountId)
	}

	if !planData.TargetAccountIds.IsNull() && !planData.TargetAccountIds.IsUnknown() {
		for _, elem := range planData.TargetAccountIds.Elements() {
			if id, ok := elem.(types.String); ok {
				accountIds = append(accountIds, id.ValueString())
			}
		}
	}

	if len(accountIds) == 0 {
		resp.Diagnostics.AddError("Error removing Organization Accounts", "Either account_id or target_account_ids is required")
		return
	}

	state := accountRemoveStateData{
		OrganizationId: planData.OrganizationId,
		AccountId:      planData.AccountId,
	}
	targetAccountIdsValue, targetDiags := types.ListValueFrom(ctx, types.StringType, accountIds)
	resp.Diagnostics.Append(targetDiags...)
	state.TargetAccountIds = targetAccountIdsValue

	result, err := r.client.RemoveAccounts(ctx, accountIds, orgId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddWarning("Error removing Organization Accounts", "Could not remove Organization Accounts: "+err.Error()+"\nReason: "+detail)
		successIds, successIdsDiags := types.ListValue(types.ObjectType{AttrTypes: successIdAttrTypes}, []attr.Value{})
		resp.Diagnostics.Append(successIdsDiags...)
		failedIds, failedIdsDiags := types.ListValue(types.ObjectType{AttrTypes: failedIdAttrTypes}, []attr.Value{})
		resp.Diagnostics.Append(failedIdsDiags...)
		state.SuccessIds = successIds
		state.FailedIds = failedIds
		diags := resp.State.Set(ctx, state)
		resp.Diagnostics.Append(diags...)
		return
	}

	for _, accountId := range accountIds {
		err = waitForAccountRemoveStatus(ctx, r.client, accountId)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error waiting for account removal",
				"Error waiting for account "+accountId+" to be removed: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	successItems := make([]attr.Value, 0)
	if result != nil {
		for _, s := range result.SuccessIds {
			obj, objDiags := types.ObjectValue(successIdAttrTypes, map[string]attr.Value{
				"success_id":   types.StringValue(s.SuccessId),
				"success_name": types.StringValue(s.SuccessName),
			})
			resp.Diagnostics.Append(objDiags...)
			successItems = append(successItems, obj)
		}
	}
	successIds, successIdsDiags := types.ListValue(types.ObjectType{AttrTypes: successIdAttrTypes}, successItems)
	resp.Diagnostics.Append(successIdsDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.SuccessIds = successIds

	failedItems := make([]attr.Value, 0)
	var failedMsg string
	if result != nil {
		for _, failed := range result.FailedIds {
			responseMap, responseDiags := buildResponseMapValue(failed.Response)
			resp.Diagnostics.Append(responseDiags...)

			obj, objDiags := types.ObjectValue(failedIdAttrTypes, map[string]attr.Value{
				"error_code":    types.StringValue(failed.ErrorCode),
				"failed_caused": types.StringValue(failed.FailedCaused),
				"failed_id":     types.StringValue(failed.FailedId),
				"failed_name":   types.StringValue(failed.FailedName),
				"response":      responseMap,
			})
			resp.Diagnostics.Append(objDiags...)

			failedItems = append(failedItems, obj)
			if failed.FailedId != "" {
				failedMsg += fmt.Sprintf("\n- %s: %s (%s)", failed.FailedId, failed.FailedName, failed.FailedCaused)
			}
		}
	}
	if failedMsg != "" {
		resp.Diagnostics.AddWarning("Some accounts failed to remove", "Failed to remove the following accounts:"+failedMsg)
	}
	failedIds, failedIdsDiags := types.ListValue(types.ObjectType{AttrTypes: failedIdAttrTypes}, failedItems)
	resp.Diagnostics.Append(failedIdsDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.FailedIds = failedIds

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *accountRemoveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accountRemoveStateData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountIds := make([]string, 0, 10)
	if singleAccountId := state.AccountId.ValueString(); singleAccountId != "" {
		accountIds = append(accountIds, singleAccountId)
	}
	if !state.TargetAccountIds.IsNull() && !state.TargetAccountIds.IsUnknown() {
		for _, elem := range state.TargetAccountIds.Elements() {
			if id, ok := elem.(types.String); ok {
				accountIds = append(accountIds, id.ValueString())
			}
		}
	}

	for _, accountId := range accountIds {
		_, err := r.client.GetAccount(ctx, accountId)
		if err != nil {
			if !strings.Contains(err.Error(), "404") {
				resp.Diagnostics.AddError(
					"Failed to read account",
					fmt.Sprintf("accountId: %s, err: %s", accountId, err.Error()),
				)
				return
			}
			resp.State.RemoveResource(ctx)
			return
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *accountRemoveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// left empty.
}

func (r *accountRemoveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// left empty.
}

type accountRemoveStateData struct {
	OrganizationId   types.String `tfsdk:"organization_id"`
	AccountId        types.String `tfsdk:"account_id"`
	TargetAccountIds types.List   `tfsdk:"target_account_ids"`
	SuccessIds       types.List   `tfsdk:"success_ids"`
	FailedIds        types.List   `tfsdk:"failed_ids"`
}

func (o accountRemoveStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	successIdAttrTypes := map[string]attr.Type{
		"success_id":   basetypes.StringType{},
		"success_name": basetypes.StringType{},
	}
	failedIdAttrTypes := map[string]attr.Type{
		"error_code":    basetypes.StringType{},
		"failed_caused": basetypes.StringType{},
		"failed_id":     basetypes.StringType{},
		"failed_name":   basetypes.StringType{},
		"response":      types.MapType{ElemType: types.StringType},
	}

	return map[string]attr.Type{
		"organization_id": basetypes.StringType{},
		"account_id":      basetypes.StringType{},
		"target_account_ids": types.ListType{
			ElemType: types.StringType,
		},
		"success_ids": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: successIdAttrTypes,
			},
		},
		"failed_ids": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: failedIdAttrTypes,
			},
		},
	}
}

func (r *accountRemoveResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) < 1 {
		resp.Diagnostics.AddError(
			"Invalid import format",
			"Expected format: organization_id:account_id or organization_id:account_id1,account_id2,...",
		)
		return
	}

	orgId := parts[0]

	var accountId string
	var accountIds []string

	if len(parts) > 1 {
		if strings.Contains(parts[1], ",") {
			accountIds = strings.Split(parts[1], ",")
		} else {
			accountId = parts[1]
			if len(parts) > 2 {
				accountIds = parts[2:]
			}
		}
	}

	state := accountRemoveStateData{
		OrganizationId: types.StringValue(orgId),
		AccountId:      types.StringValue(accountId),
	}

	if len(accountIds) > 0 {
		targetAccountIdsValue, targetDiags := types.ListValueFrom(ctx, types.StringType, accountIds)
		resp.Diagnostics.Append(targetDiags...)
		state.TargetAccountIds = targetAccountIdsValue
	}

	successIds, successIdsDiags := types.ListValue(types.ObjectType{AttrTypes: successIdAttrTypes}, []attr.Value{})
	resp.Diagnostics.Append(successIdsDiags...)
	failedIds, failedIdsDiags := types.ListValue(types.ObjectType{AttrTypes: failedIdAttrTypes}, []attr.Value{})
	resp.Diagnostics.Append(failedIdsDiags...)
	state.SuccessIds = successIds
	state.FailedIds = failedIds

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func waitForAccountRemoveStatus(ctx context.Context, orgClient *organization.Client, accountId string) error {
	return client.WaitForStatus(ctx, nil, []string{"EXISTS"}, []string{"REMOVED"}, func() (interface{}, string, error) {
		accountResp, err := orgClient.GetAccount(ctx, accountId)
		if err != nil {
			return "REMOVED", "REMOVED", nil
		}
		return accountResp, "EXISTS", nil
	}, -1, -1, -1, -1)
}
