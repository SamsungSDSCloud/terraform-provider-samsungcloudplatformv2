package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                = &organizationMembershipResource{}
	_ resource.ResourceWithConfigure   = &organizationMembershipResource{}
	_ resource.ResourceWithImportState = &organizationMembershipResource{}
)

func NewOrganizationMembershipResource() resource.Resource {
	return &organizationMembershipResource{}
}

type organizationMembershipResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *organizationMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_membership"
}

func (r *organizationMembershipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Organization Membership - Leave an organization",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Required: true,
			},
			"success": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Success Info. \n" +
					"  - example : '{success_id: b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0, success_name: My Organization}' \n",
				Attributes: map[string]schema.Attribute{
					"success_id": schema.StringAttribute{
						Description: "ID of the account that was left. \n" +
							"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
						Computed: true,
					},
					"success_name": schema.StringAttribute{
						Description: "Name of the organization that was left. \n" +
							"  - example : 'My Organization' \n",
						Computed: true,
					},
				},
			},
			"failed": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Failed Info. \n" +
					"  - example : '{error_code: Organization.AccountNotRemovable, failed_id: b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0, failed_name: My Organization, failed_caused: Payment method registration required.}' \n",
				Attributes: map[string]schema.Attribute{
					"failed_id": schema.StringAttribute{
						Description: "ID of the account that failed to leave. \n" +
							"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
						Computed: true,
					},
					"failed_name": schema.StringAttribute{
						Description: "Name of the organization that failed to leave. \n" +
							"  - example : 'My Organization' \n",
						Computed: true,
					},
					"failed_caused": schema.StringAttribute{
						Description: "Reason for leaving failure. \n" +
							"  - example : 'Payment method registration required.' \n",
						Computed: true,
					},
					"error_code": schema.StringAttribute{
						Description: "Error code returned when the operation fails. \n" +
							"  - example : 'Organization.AccountNotRemovable' \n",
						Computed: true,
					},
					"response": schema.MapAttribute{
						ElementType: types.StringType,
						Description: "Response details. \n" +
							"  - example : {} \n",
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *organizationMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *organizationMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var stateData organizationMembershipStateData
	diags := req.Plan.Get(ctx, &stateData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.LeaveOrganization(ctx, stateData.OrganizationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error leaving organization",
			err.Error(),
		)
		return
	}

	var successVal types.Object
	if result.HasSuccess() {
		successObj := result.GetSuccess()
		var successDiags diag.Diagnostics
		successVal, successDiags = types.ObjectValue(
			map[string]attr.Type{
				"success_id":   basetypes.StringType{},
				"success_name": basetypes.StringType{},
			},
			map[string]attr.Value{
				"success_id":   types.StringValue(successObj.GetSuccessId()),
				"success_name": types.StringValue(successObj.GetSuccessName()),
			},
		)
		resp.Diagnostics.Append(successDiags...)
	} else {
		successVal = types.ObjectNull(map[string]attr.Type{
			"success_id":   basetypes.StringType{},
			"success_name": basetypes.StringType{},
		})
	}

	var failedVal types.Object
	if result.HasFailed() {
		failedObj := result.GetFailed()
		responseMap, responseDiags := types.MapValue(types.StringType, map[string]attr.Value{})
		resp.Diagnostics.Append(responseDiags...)
		var failedDiags diag.Diagnostics
		failedVal, failedDiags = types.ObjectValue(
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
				"response":      responseMap,
			},
		)
		resp.Diagnostics.Append(failedDiags...)
	} else {
		failedVal = types.ObjectNull(map[string]attr.Type{
			"failed_id":     basetypes.StringType{},
			"failed_name":   basetypes.StringType{},
			"failed_caused": basetypes.StringType{},
			"error_code":    basetypes.StringType{},
			"response":      types.MapType{ElemType: types.StringType},
		})
	}

	state := organizationMembershipStateData{
		OrganizationId: stateData.OrganizationId,
		Success:        successVal,
		Failed:         failedVal,
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *organizationMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organizationMembershipStateData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *organizationMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// left empty.
}

func (r *organizationMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// left empty. The deletion happens in Create (leave organization action)
}

type OrganizationMembershipSuccessStateData struct {
	SuccessId   types.String `tfsdk:"success_id"`
	SuccessName types.String `tfsdk:"success_name"`
}

func (o OrganizationMembershipSuccessStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"success_id":   basetypes.StringType{},
		"success_name": basetypes.StringType{},
	}
}

type OrganizationMembershipFailedStateData struct {
	FailedId     types.String      `tfsdk:"failed_id"`
	FailedName   types.String      `tfsdk:"failed_name"`
	FailedCaused types.String      `tfsdk:"failed_caused"`
	ErrorCode    types.String      `tfsdk:"error_code"`
	Response     map[string]string `tfsdk:"response"`
}

func (o OrganizationMembershipFailedStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"failed_id":     basetypes.StringType{},
		"failed_name":   basetypes.StringType{},
		"failed_caused": basetypes.StringType{},
		"error_code":    basetypes.StringType{},
		"response":      types.MapType{ElemType: types.StringType},
	}
}

type organizationMembershipStateData struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	Success        types.Object `tfsdk:"success"`
	Failed         types.Object `tfsdk:"failed"`
}

func (o organizationMembershipStateData) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id": basetypes.StringType{},
		"success":         types.ObjectType{},
		"failed":          types.ObjectType{},
	}
}

func (r *organizationMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("organization_id"), req, resp)
}
