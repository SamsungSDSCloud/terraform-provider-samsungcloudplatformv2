package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &invitationAcceptResource{}
	_ resource.ResourceWithConfigure = &invitationAcceptResource{}
)

func NewInvitationAcceptResource() resource.Resource {
	return &invitationAcceptResource{}
}

type invitationAcceptResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *invitationAcceptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_invitation_accept"
}

func (r *invitationAcceptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Invitation Accept - Use a separate provider alias with target account credentials",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Invitation ID",
				Required:    true,
			},
			"master_account_email": schema.StringAttribute{
				Description: "Master Account Email",
				Computed:    true,
			},
			"success_id": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"success_id": schema.StringAttribute{
						Description: "Success ID (Organization ID)",
						Computed:    true,
					},
					"success_name": schema.StringAttribute{
						Description: "Success Name (Organization Name)",
						Computed:    true,
					},
				},
			},
			"failed_id": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"failed_id": schema.StringAttribute{
						Description: "Failed ID",
						Computed:    true,
					},
					"failed_name": schema.StringAttribute{
						Description: "Failed Name",
						Computed:    true,
					},
					"failed_caused": schema.StringAttribute{
						Description: "Failed Caused",
						Computed:    true,
					},
					"error_code": schema.StringAttribute{
						Description: "Error Code",
						Computed:    true,
					},
					"response": schema.MapAttribute{
						ElementType: types.StringType,
						Description: "Response",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *invitationAcceptResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *invitationAcceptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var stateData invitationAcceptStateData
	diags := req.Plan.Get(ctx, &stateData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.AcceptInvitation(ctx, stateData.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error accepting invitation",
			err.Error(),
		)
		return
	}

	masterAccountEmail := result.GetMasterAccountEmail()

	var successIdVal types.Object
	if result.HasSuccessId() {
		successObj := result.GetSuccessId()
		successIdVal, _ = types.ObjectValue(
			map[string]attr.Type{
				"success_id":   basetypes.StringType{},
				"success_name": basetypes.StringType{},
			},
			map[string]attr.Value{
				"success_id":   types.StringValue(successObj.GetSuccessId()),
				"success_name": types.StringValue(successObj.GetSuccessName()),
			},
		)
	}

	var failedIdVal types.Object
	if result.HasFailedId() {
		failedObj := result.GetFailedId()
		failedIdVal, _ = types.ObjectValue(
			map[string]attr.Type{
				"failed_id":     basetypes.StringType{},
				"failed_name":   basetypes.StringType{},
				"failed_caused": basetypes.StringType{},
				"error_code":    basetypes.StringType{},
				"response":      types.MapType{ElemType: types.StringType},
			},
			map[string]attr.Value{
				"failed_id":     types.StringValue(failedObj.GetFailedId()),
				"failed_name":   types.StringValue(failedObj.GetFailedName()),
				"failed_caused": types.StringValue(failedObj.GetFailedCaused()),
				"error_code":    types.StringValue(failedObj.GetErrorCode()),
				"response":      types.MapValueMust(types.StringType, map[string]attr.Value{}),
			},
		)
	}

	state := invitationAcceptStateData{
		Id:                 stateData.Id,
		MasterAccountEmail: types.StringValue(masterAccountEmail),
		SuccessId:          successIdVal,
		FailedId:           failedIdVal,
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *invitationAcceptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *invitationAcceptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *invitationAcceptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

type InvitationAcceptSuccessStateData struct {
	SuccessId   types.String `tfsdk:"success_id"`
	SuccessName types.String `tfsdk:"success_name"`
}

func (o InvitationAcceptSuccessStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"success_id":   basetypes.StringType{},
		"success_name": basetypes.StringType{},
	}
}

type InvitationAcceptFailedStateData struct {
	FailedId     types.String      `tfsdk:"failed_id"`
	FailedName   types.String      `tfsdk:"failed_name"`
	FailedCaused types.String      `tfsdk:"failed_caused"`
	ErrorCode    types.String      `tfsdk:"error_code"`
	Response     map[string]string `tfsdk:"response"`
}

func (o InvitationAcceptFailedStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"failed_id":     basetypes.StringType{},
		"failed_name":   basetypes.StringType{},
		"failed_caused": basetypes.StringType{},
		"error_code":    basetypes.StringType{},
		"response":      types.MapType{ElemType: types.StringType},
	}
}

type invitationAcceptStateData struct {
	Id                 types.String `tfsdk:"id"`
	MasterAccountEmail types.String `tfsdk:"master_account_email"`
	SuccessId          types.Object `tfsdk:"success_id"`
	FailedId           types.Object `tfsdk:"failed_id"`
}

func (o invitationAcceptStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                   basetypes.StringType{},
		"master_account_email": basetypes.StringType{},
		"success_id":           types.ObjectType{},
		"failed_id":            types.ObjectType{},
	}
}
