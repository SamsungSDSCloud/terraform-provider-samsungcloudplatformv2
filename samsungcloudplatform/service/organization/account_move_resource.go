package organization

import (
	"context"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &accountMoveResource{}
	_ resource.ResourceWithConfigure = &accountMoveResource{}
)

func NewAccountMoveResource() resource.Resource {
	return &accountMoveResource{}
}

type accountMoveResource struct {
	config *client.SCPClient
	client *organization.Client
}

func (r *accountMoveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_account_move"
}

func (r *accountMoveResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Move Organization Accounts to different OU",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "Organization ID",
				Required:    true,
			},
			"parent_unit_id": schema.StringAttribute{
				Description: "Target Parent Organization Unit ID",
				Required:    true,
			},
			"target_account_ids": schema.ListAttribute{
				Description: "Account IDs to move",
				Required:    true,
				ElementType: types.StringType,
			},
			"success_ids": schema.ListNestedAttribute{
				Description: "Successfully moved account information",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"success_id": schema.StringAttribute{
							Description: "Success ID",
							Computed:    true,
						},
						"success_name": schema.StringAttribute{
							Description: "Success Name",
							Computed:    true,
						},
						"target_id": schema.StringAttribute{
							Description: "Target ID",
							Computed:    true,
						},
						"target_name": schema.StringAttribute{
							Description: "Target Name",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *accountMoveResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
}

func (r *accountMoveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organization.AccountMoveResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := plan.OrganizationId.ValueString()
	parentUnitId := plan.ParentUnitId.ValueString()

	targetAccountIds := make([]string, 0, 10)
	for _, elem := range plan.TargetAccountIds.Elements() {
		if id, ok := elem.(types.String); ok {
			targetAccountIds = append(targetAccountIds, id.ValueString())
		}
	}

	moveResp, err := r.client.MoveAccounts(ctx, targetAccountIds, parentUnitId, orgId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error moving Organization Accounts",
			"Could not move Organization Accounts, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var successIdsList []types.Object
	if moveResp != nil && moveResp.SuccessIds != nil {
		for _, s := range moveResp.SuccessIds {
			successObj, diags := types.ObjectValue(
				organization.MoveAccountSuccessValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"success_id":   types.StringValue(s.SuccessId),
					"success_name": types.StringValue(s.SuccessName),
					"target_id":    types.StringValue(s.TargetId),
					"target_name":  types.StringValue(s.TargetName),
				},
			)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			successIdsList = append(successIdsList, successObj)
		}
	}

	var successIdsValue types.List
	if len(successIdsList) > 0 {
		successIdsValue, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: organization.MoveAccountSuccessValue{}.AttributeTypes(ctx)}, successIdsList)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		successIdsValue = types.ListNull(types.ObjectType{AttrTypes: organization.MoveAccountSuccessValue{}.AttributeTypes(ctx)})
	}

	stateData := organization.AccountMoveResource{
		OrganizationId:   types.StringValue(orgId),
		ParentUnitId:     types.StringValue(parentUnitId),
		TargetAccountIds: plan.TargetAccountIds,
		SuccessIds:       successIdsValue,
	}

	diags = resp.State.Set(ctx, &stateData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *accountMoveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *accountMoveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *accountMoveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

type SCPClient struct {
}
