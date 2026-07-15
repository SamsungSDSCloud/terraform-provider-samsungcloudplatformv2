package organization

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &invitationResource{}
	_ resource.ResourceWithConfigure   = &invitationResource{}
	_ resource.ResourceWithImportState = &invitationResource{}
)

func NewInvitationResource() resource.Resource {
	return &invitationResource{}
}

type invitationResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *invitationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_invitation"
}

func (r *invitationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages invitations for accounts to join the organization",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Optional: true,
			},
			"target_login_ids": schema.ListAttribute{
				Description: "Invitation Receiver Account Email List. \n" +
					"  - example : ['score-1@samsung.com'] \n",
				ElementType: types.StringType,
				Optional:    true,
			},
			"success_ids": schema.ListAttribute{
				Description: "Organization Invitation Success Login ID List. \n" +
					"  - example : ['log-archive@samsung.com'] \n",
				ElementType: types.StringType,
				Computed:    true,
			},
			"failed_ids": schema.ListNestedAttribute{
				Description: "Organization Invitation Fail Login ID List. \n" +
					"  - example : '[{error_code: Invitation.AlreadySentError, failed_id: [score-1@samsung.com], ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Computed: true,
							Description: "Error code returned when the operation fails. \n" +
								"  - example : 'Invitation.AlreadySentError' \n",
						},
						"failed_caused": schema.StringAttribute{
							Computed: true,
							Description: "Failure Reason. \n" +
								"  - example : 'Account b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0 is not found' \n",
						},
						"failed_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "List of IDs that failed to be processed. \n" +
								"  - example : ['user1@samsung.com'] \n",
						},
						"response": schema.StringAttribute{
							Computed: true,
							Description: "Quota Response. \n" +
								"  - example : {} \n",
						},
					},
				},
				Computed: true,
				Optional: true,
			},
		},
	}
}

func (r *invitationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *invitationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planData organization.InvitationResource
	diags := req.Plan.Get(ctx, &planData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := planData.OrganizationId.ValueString()
	if orgId == "" {
		orgId = "default"
	}
	planData.OrganizationId = types.StringValue(orgId)

	result, err := r.client.CreateInvitation(ctx, planData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating invitation",
			err.Error(),
		)
		return
	}

	successIds := make([]attr.Value, 0, len(result.SuccessIds))
	for _, id := range result.SuccessIds {
		successIds = append(successIds, types.StringValue(id))
	}

	failedIds := make([]attr.Value, 0, len(result.FailedIds))
	for _, f := range result.FailedIds {
		respMap := ""
		if f.Response != nil {
			respMap = fmt.Sprintf("%v", f.Response)
		}

		failedIdList, hasErr := buildFailedIdList(ctx, f.FailedId, &resp.Diagnostics)
		if hasErr {
			return
		}

		failedIdObj, failedIdObjDiags := types.ObjectValue(
			organization.InvitationCreateFailCausedValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"error_code":    types.StringValue(f.ErrorCode),
				"failed_caused": types.StringValue(f.FailedCaused),
				"failed_ids":    failedIdList,
				"response":      types.StringValue(respMap),
			},
		)
		resp.Diagnostics.Append(failedIdObjDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		failedIds = append(failedIds, failedIdObj)
	}

	successIdsList, successIdsDiags := types.ListValueFrom(ctx, types.StringType, result.SuccessIds)
	resp.Diagnostics.Append(successIdsDiags...)

	emptyFailedIdsAttrTypes := organization.InvitationCreateFailCausedValue{}.AttributeTypes(ctx)

	objectType := types.ObjectType{AttrTypes: emptyFailedIdsAttrTypes}
	failedIdsObjList, failedIdsListDiags := types.ListValue(objectType, failedIds)
	resp.Diagnostics.Append(failedIdsListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := organization.InvitationResource{
		OrganizationId: planData.OrganizationId,
		TargetLoginIds: planData.TargetLoginIds,
		SuccessIds:     successIdsList,
		FailedIds:      failedIdsObjList,
	}

	for _, id := range result.SuccessIds {
		err = waitForInvitationReady(ctx, r.client, id, orgId)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error waiting for invitation creation",
				"Error waiting for invitation "+id+" to be ready: "+err.Error()+"\nReason: "+detail,
			)
			return
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organization.InvitationResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	successIds := make([]string, 0, 10)
	for _, elem := range state.SuccessIds.Elements() {
		if id, ok := elem.(types.String); ok {
			successIds = append(successIds, id.ValueString())
		}
	}

	if len(successIds) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	request := organization.OrganizationInvitationsDataSource{
		OrganizationId: state.OrganizationId,
	}
	result, err := r.client.GetOrganizationInvitations(ctx, request)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	existingIds := make(map[string]bool)
	for _, inv := range result.GetOrganizationInvitations() {
		existingIds[inv.Id] = true
	}

	for _, id := range successIds {
		if !existingIds[id] {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData organization.InvitationResource
	diags := req.Plan.Get(ctx, &planData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var stateData organization.InvitationResource
	diags = req.State.Get(ctx, &stateData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if stateData.SuccessIds.IsNull() == false && stateData.SuccessIds.IsUnknown() == false && len(stateData.SuccessIds.Elements()) > 0 {
		cancelReq := organization.InvitationCancelRequest{
			OrganizationId: stateData.OrganizationId,
			Ids:            stateData.SuccessIds,
		}
		_, err := r.client.CancelInvitations(ctx, cancelReq)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error canceling old invitation",
				err.Error(),
			)
			return
		}
	}
	result, err := r.client.CreateInvitation(ctx, planData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating invitation",
			err.Error(),
		)
		return
	}
	successIdsList, successIdsDiags := types.ListValueFrom(ctx, types.StringType, result.SuccessIds)
	resp.Diagnostics.Append(successIdsDiags...)
	emptyFailedIdsAttrTypes := organization.InvitationCreateFailCausedValue{}.AttributeTypes(ctx)
	objectType := types.ObjectType{AttrTypes: emptyFailedIdsAttrTypes}
	failedIds := make([]attr.Value, 0, len(result.FailedIds))
	for _, f := range result.FailedIds {
		respMap := ""
		if f.Response != nil {
			respMap = fmt.Sprintf("%v", f.Response)
		}
		failedIdList, hasErr := buildFailedIdList(ctx, f.FailedId, &resp.Diagnostics)
		if hasErr {
			return
		}
		failedIdObj, failedIdObjDiags := types.ObjectValue(
			organization.InvitationCreateFailCausedValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"error_code":    types.StringValue(f.ErrorCode),
				"failed_caused": types.StringValue(f.FailedCaused),
				"failed_ids":    failedIdList,
				"response":      types.StringValue(respMap),
			},
		)
		resp.Diagnostics.Append(failedIdObjDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		failedIds = append(failedIds, failedIdObj)
	}
	failedIdsObjList, failedIdsListDiags := types.ListValue(objectType, failedIds)
	resp.Diagnostics.Append(failedIdsListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state := organization.InvitationResource{
		OrganizationId: planData.OrganizationId,
		TargetLoginIds: planData.TargetLoginIds,
		SuccessIds:     successIdsList,
		FailedIds:      failedIdsObjList,
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organization.InvitationResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cancelReq := organization.InvitationCancelRequest{
		OrganizationId: state.OrganizationId,
		Ids:            state.SuccessIds,
	}
	_, err := r.client.CancelInvitations(ctx, cancelReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error canceling invitation",
			err.Error(),
		)
		return
	}
}

func (r *invitationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) < 1 {
		resp.Diagnostics.AddError(
			"Invalid import format",
			"Expected format: organization_id:login_id1,login_id2,...",
		)
		return
	}

	orgId := parts[0]
	var targetLoginIds []string
	if len(parts) > 1 {
		targetLoginIds = strings.Split(parts[1], ",")
	}

	targetLoginIdsList, targetLoginIdsDiags := types.ListValueFrom(ctx, types.StringType, targetLoginIds)
	resp.Diagnostics.Append(targetLoginIdsDiags...)

	emptyFailedIdsAttrTypes := organization.InvitationCreateFailCausedValue{}.AttributeTypes(ctx)
	failedIdsObjList, failedIdsDiags := types.ListValue(types.ObjectType{AttrTypes: emptyFailedIdsAttrTypes}, []attr.Value{})
	resp.Diagnostics.Append(failedIdsDiags...)

	successIds, successIdsDiags := types.ListValue(types.StringType, []attr.Value{})
	resp.Diagnostics.Append(successIdsDiags...)

	state := organization.InvitationResource{
		OrganizationId: types.StringValue(orgId),
		TargetLoginIds: targetLoginIdsList,
		SuccessIds:     successIds,
		FailedIds:      failedIdsObjList,
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func waitForInvitationReady(ctx context.Context, orgClient *organization.Client, loginId string, orgId string) error {
	return client.WaitForStatus(ctx, nil, []string{}, []string{"INVITING"}, func() (interface{}, string, error) {
		request := organization.OrganizationInvitationsDataSource{
			OrganizationId: types.StringValue(orgId),
		}
		result, err := orgClient.GetOrganizationInvitations(ctx, request)
		if err != nil {
			return nil, "", err
		}
		for _, inv := range result.GetOrganizationInvitations() {
			if inv.LoginId.Get() != nil && *inv.LoginId.Get() == loginId {
				return inv, string(inv.State), nil
			}
		}
		return "INVITING", "INVITING", nil
	}, -1, -1, -1, -1)
}

func buildFailedIdList(ctx context.Context, failedId string, diags *diag.Diagnostics) (types.List, bool) {
	if failedId == "" {
		return types.ListNull(types.StringType), false
	}
	var failedIdSlice []string
	if err := json.Unmarshal([]byte(failedId), &failedIdSlice); err == nil {
		list, d := types.ListValueFrom(ctx, types.StringType, failedIdSlice)
		diags.Append(d...)
		return list, diags.HasError()
	}
	list, d := types.ListValue(types.StringType, []attr.Value{types.StringValue(failedId)})
	diags.Append(d...)
	return list, diags.HasError()
}
