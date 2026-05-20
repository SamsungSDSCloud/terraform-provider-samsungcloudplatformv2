package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Nested object types for success_ids
type successIdItem struct {
	SuccessId   string `tfsdk:"success_id"`
	SuccessName string `tfsdk:"success_name"`
}

type failedIdItem struct {
	ErrorCode    string            `tfsdk:"error_code"`
	FailedCaused string            `tfsdk:"failed_caused"`
	FailedId     string            `tfsdk:"failed_id"`
	FailedName   string            `tfsdk:"failed_name"`
	Response     map[string]string `tfsdk:"response"`
}

var (
	_ resource.Resource              = &accountRemoveResource{}
	_ resource.ResourceWithConfigure = &accountRemoveResource{}
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
				Description: "Organization ID",
				Required:    true,
			},
			"account_id": schema.StringAttribute{
				Description: "Account ID to remove (single)",
				Optional:    true,
			},
			"target_account_ids": schema.ListAttribute{
				Description: "Account IDs to remove (multiple)",
				ElementType: types.StringType,
				Optional:    true,
			},
			"success_ids": schema.ListNestedAttribute{
				Description: "Successfully removed accounts",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"success_id": schema.StringAttribute{
							Computed: true,
						},
						"success_name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"failed_ids": schema.ListNestedAttribute{
				Description: "Failed accounts",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Computed: true,
						},
						"failed_caused": schema.StringAttribute{
							Computed: true,
						},
						"failed_id": schema.StringAttribute{
							Computed: true,
						},
						"failed_name": schema.StringAttribute{
							Computed: true,
						},
						"response": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
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

func (r *accountRemoveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planData accountRemoveStateData
	diags := req.Plan.Get(ctx, &planData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := planData.OrganizationId.ValueString()

	accountIds := make([]string, 0, 10)

	singleAccountId := planData.AccountId.ValueString()
	if singleAccountId != "" {
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
		resp.Diagnostics.AddError(
			"Error removing Organization Accounts",
			"Either account_id or target_account_ids is required",
		)
		return
	}

	state := accountRemoveStateData{
		OrganizationId: planData.OrganizationId,
		AccountId:      planData.AccountId,
	}

	accountIdStrings := make([]string, 0, len(accountIds))
	accountIdStrings = append(accountIdStrings, accountIds...)
	targetAccountIdsValue, _ := types.ListValueFrom(ctx, types.StringType, accountIdStrings)
	state.TargetAccountIds = targetAccountIdsValue

	result, err := r.client.RemoveAccounts(ctx, accountIds, orgId)

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddWarning(
			"Error removing Organization Accounts",
			"Could not remove Organization Accounts: "+err.Error()+"\nReason: "+detail,
		)
		state.SuccessIds = types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
			"success_id":   basetypes.StringType{},
			"success_name": basetypes.StringType{},
		}}, []attr.Value{})
		state.FailedIds = types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
			"error_code":    basetypes.StringType{},
			"failed_caused": basetypes.StringType{},
			"failed_id":     basetypes.StringType{},
			"failed_name":   basetypes.StringType{},
			"response":      types.MapType{ElemType: types.StringType},
		}}, []attr.Value{})
		diags := resp.State.Set(ctx, state)
		resp.Diagnostics.Append(diags...)
		return
	}

	if result != nil && len(result.SuccessIds) > 0 {
		successIdAttrTypes := map[string]attr.Type{
			"success_id":   basetypes.StringType{},
			"success_name": basetypes.StringType{},
		}

		successItems := make([]attr.Value, 0, len(result.SuccessIds))
		for _, s := range result.SuccessIds {
			obj, _ := types.ObjectValue(successIdAttrTypes, map[string]attr.Value{
				"success_id":   types.StringValue(s.SuccessId),
				"success_name": types.StringValue(s.SuccessName),
			})
			successItems = append(successItems, obj)
		}
		state.SuccessIds = types.ListValueMust(types.ObjectType{AttrTypes: successIdAttrTypes}, successItems)
	} else {
		state.SuccessIds = types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
			"success_id":   basetypes.StringType{},
			"success_name": basetypes.StringType{},
		}}, []attr.Value{})
	}

	failedIdAttrTypes := map[string]attr.Type{
		"error_code":    basetypes.StringType{},
		"failed_caused": basetypes.StringType{},
		"failed_id":     basetypes.StringType{},
		"failed_name":   basetypes.StringType{},
		"response":      types.MapType{ElemType: types.StringType},
	}

	if result != nil && len(result.FailedIds) > 0 {
		failedItems := make([]attr.Value, 0, len(result.FailedIds))
		var failedMsg string
		for _, failed := range result.FailedIds {
			var responseValue types.Map
			if failed.Response == nil || len(failed.Response) == 0 {
				responseValue = types.MapValueMust(types.StringType, map[string]attr.Value{})
			} else {
				responseMap := map[string]attr.Value{}
				for k, v := range failed.Response {
					responseMap[k] = types.StringValue(fmt.Sprintf("%v", v))
				}
				responseValue, _ = types.MapValue(types.StringType, responseMap)
			}
			obj, _ := types.ObjectValue(failedIdAttrTypes, map[string]attr.Value{
				"error_code":    types.StringValue(failed.ErrorCode),
				"failed_caused": types.StringValue(failed.FailedCaused),
				"failed_id":     types.StringValue(failed.FailedId),
				"failed_name":   types.StringValue(failed.FailedName),
				"response":      responseValue,
			})
			failedItems = append(failedItems, obj)
			if failed.FailedId != "" {
				failedMsg += fmt.Sprintf("\n- %s: %s (%s)", failed.FailedId, failed.FailedName, failed.FailedCaused)
			}
		}
		if failedMsg != "" {
			resp.Diagnostics.AddWarning(
				"Some accounts failed to remove",
				"Failed to remove the following accounts:"+failedMsg,
			)
		}
		state.FailedIds = types.ListValueMust(types.ObjectType{AttrTypes: failedIdAttrTypes}, failedItems)
	} else {
		state.FailedIds = types.ListValueMust(types.ObjectType{AttrTypes: failedIdAttrTypes}, []attr.Value{})
	}

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

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *accountRemoveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *accountRemoveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
