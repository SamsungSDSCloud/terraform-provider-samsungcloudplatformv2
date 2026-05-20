package organization

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &invitationResource{}
	_ resource.ResourceWithConfigure = &invitationResource{}
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
	resp.TypeName = req.ProviderTypeName + "_invitation"
}

func (r *invitationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Invitation",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "Organization ID",
				Optional:    true,
			},
			"target_login_ids": schema.ListAttribute{
				Description: "Target Login IDs",
				ElementType: types.StringType,
				Optional:    true,
			},
			"success_ids": schema.ListAttribute{
				Description: "Success IDs",
				ElementType: types.StringType,
				Computed:    true,
			},
			"failed_ids": schema.ListNestedAttribute{
				Description: "Failed IDs",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Computed: true,
						},
						"failed_caused": schema.StringAttribute{
							Computed: true,
						},
						"failed_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"response": schema.StringAttribute{
							Computed: true,
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

		var failedIdList types.List
		if f.FailedId != "" {
			var failedIdSlice []string
			if err := json.Unmarshal([]byte(f.FailedId), &failedIdSlice); err == nil {
				failedIdList, _ = types.ListValueFrom(ctx, types.StringType, failedIdSlice)
			} else {
				failedIdList = types.ListValueMust(types.StringType, []attr.Value{types.StringValue(f.FailedId)})
			}
		} else {
			failedIdList = types.ListNull(types.StringType)
		}

		failedIds = append(failedIds, types.ObjectValueMust(
			organization.InvitationCreateFailCausedValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"error_code":    types.StringValue(f.ErrorCode),
				"failed_caused": types.StringValue(f.FailedCaused),
				"failed_ids":    failedIdList,
				"response":      types.StringValue(respMap),
			},
		))
	}

	successIdsList, _ := types.ListValueFrom(ctx, types.StringType, result.SuccessIds)

	emptyFailedIdsAttrTypes := organization.InvitationCreateFailCausedValue{}.AttributeTypes(ctx)
	var failedIdsObjList types.List
	objectType := types.ObjectType{AttrTypes: emptyFailedIdsAttrTypes}

	if len(failedIds) > 0 {
		failedIdsObjList = types.ListValueMust(objectType, failedIds)
	} else {
		failedIdsObjList = types.ListValueMust(objectType, []attr.Value{})
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

func (r *invitationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	successIdsList, _ := types.ListValueFrom(ctx, types.StringType, result.SuccessIds)

	emptyFailedIdsAttrTypes := organization.InvitationCreateFailCausedValue{}.AttributeTypes(ctx)
	var failedIdsObjList types.List
	objectType := types.ObjectType{AttrTypes: emptyFailedIdsAttrTypes}

	failedIds := make([]attr.Value, 0, len(result.FailedIds))
	for _, f := range result.FailedIds {
		respMap := ""
		if f.Response != nil {
			respMap = fmt.Sprintf("%v", f.Response)
		}

		var failedIdList types.List
		if f.FailedId != "" {
			var failedIdSlice []string
			if err := json.Unmarshal([]byte(f.FailedId), &failedIdSlice); err == nil {
				failedIdList, _ = types.ListValueFrom(ctx, types.StringType, failedIdSlice)
			} else {
				failedIdList = types.ListValueMust(types.StringType, []attr.Value{types.StringValue(f.FailedId)})
			}
		} else {
			failedIdList = types.ListNull(types.StringType)
		}

		failedIds = append(failedIds, types.ObjectValueMust(
			organization.InvitationCreateFailCausedValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"error_code":    types.StringValue(f.ErrorCode),
				"failed_caused": types.StringValue(f.FailedCaused),
				"failed_ids":    failedIdList,
				"response":      types.StringValue(respMap),
			},
		))
	}

	if len(failedIds) > 0 {
		failedIdsObjList = types.ListValueMust(objectType, failedIds)
	} else {
		failedIdsObjList = types.ListValueMust(objectType, []attr.Value{})
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
