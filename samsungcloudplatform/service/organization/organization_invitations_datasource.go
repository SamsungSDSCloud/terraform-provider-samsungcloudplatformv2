package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ datasource.DataSource              = &organizationInvitationsDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationInvitationsDataSource{}
)

func NewOrganizationInvitationsDataSource() datasource.DataSource {
	return &organizationInvitationsDataSource{}
}

type organizationInvitationsDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (d *organizationInvitationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_invitations"
}

func (d *organizationInvitationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationInvitationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List Organization Invitations",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by Organization ID",
			},
			"sort": schema.StringAttribute{
				Optional:    true,
				Description: "Sort criteria (e.g., 'created_at:desc')",
			},
			"account_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by Account ID",
			},
			"account_name": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by Account Name",
			},
			"account_email": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by Account Email",
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by Invitation State",
			},
			"login_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by Login ID",
			},
			"organization_invitations": schema.ListAttribute{
				Computed:    true,
				ElementType: types.ObjectType{AttrTypes: organization.GetInvitationItemAttributeTypes()},
				Description: "A list of Organization Invitations",
			},
			"total_count": schema.Int64Attribute{
				Computed:            true,
				Description:         "Total Count",
				MarkdownDescription: "Total Count",
			},
			"page": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Page Number",
				MarkdownDescription: "Page Number",
			},
			"size": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Page Size",
				MarkdownDescription: "Page Size",
			},
			"sort_result": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Sort Criteria",
				MarkdownDescription: "Sort Criteria",
			},
		},
	}
}

func (d *organizationInvitationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var configData OrganizationInvitationsDataSource
	diags := req.Config.Get(ctx, &configData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiRequest := organization.OrganizationInvitationsDataSource{
		OrganizationId: configData.OrganizationId,
		Size:           configData.Size,
		Page:           configData.Page,
		Sort:           configData.Sort,
		AccountId:      configData.AccountId,
		AccountName:    configData.AccountName,
		AccountEmail:   configData.AccountEmail,
		State:          configData.State,
		LoginId:        configData.LoginId,
	}

	result, err := d.client.GetOrganizationInvitations(ctx, apiRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Organization Invitations",
			err.Error(),
		)
		return
	}

	invitationAttrTypes := organization.GetInvitationItemAttributeTypes()
	invitationObjects := []organization.InvitationItem{}
	for _, inv := range result.GetOrganizationInvitations() {
		invitationObjects = append(invitationObjects, organization.InvitationItem{
			Id:              types.StringValue(inv.Id),
			OrganizationId:  types.StringValue(inv.OrganizationId),
			AccountName:     types.StringValue(inv.GetAccountName()),
			AccountEmail:    types.StringValue(inv.GetAccountEmail()),
			LoginId:         types.StringValue(inv.GetLoginId()),
			TargetAccountId: types.StringValue(inv.TargetAccountId),
			State:           types.StringValue(string(inv.State)),
			CreatedAt:       types.StringValue(inv.CreatedAt.Format(time.RFC3339)),
			CreatedBy:       types.StringValue(inv.CreatedBy),
			ModifiedAt:      types.StringValue(inv.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:      types.StringValue(inv.ModifiedBy),
		})
	}

	invitationsListVal, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: invitationAttrTypes}, invitationObjects)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sortListVal, diags := types.ListValueFrom(ctx, types.StringType, result.GetSort())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	configData.OrganizationInvitations = invitationsListVal
	configData.TotalCount = types.Int64Value(int64(result.GetCount()))
	configData.Page = types.Int64Value(int64(result.GetPage()))
	configData.Size = types.Int64Value(int64(result.GetSize()))
	configData.SortResult = sortListVal

	diags = resp.State.Set(ctx, &configData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type OrganizationInvitationsDataSource struct {
	OrganizationInvitations types.List   `tfsdk:"organization_invitations"`
	TotalCount              types.Int64  `tfsdk:"total_count"`
	Page                    types.Int64  `tfsdk:"page"`
	Size                    types.Int64  `tfsdk:"size"`
	SortResult              types.List   `tfsdk:"sort_result"`
	OrganizationId          types.String `tfsdk:"organization_id"`
	Sort                    types.String `tfsdk:"sort"`
	AccountId               types.String `tfsdk:"account_id"`
	AccountName             types.String `tfsdk:"account_name"`
	AccountEmail            types.String `tfsdk:"account_email"`
	State                   types.String `tfsdk:"state"`
	LoginId                 types.String `tfsdk:"login_id"`
}

func (o OrganizationInvitationsDataSource) AttributeTypes(ctx context.Context) map[string]attr.Type {
	invitationAttrTypes := organization.GetInvitationItemAttributeTypes()
	invitationObjType := types.ObjectType{AttrTypes: invitationAttrTypes}
	return map[string]attr.Type{
		"organization_invitations": types.ListType{ElemType: invitationObjType},
		"total_count":              basetypes.Int64Type{},
		"page":                     basetypes.Int64Type{},
		"size":                     basetypes.Int64Type{},
		"sort_result":              types.ListType{ElemType: types.StringType},
		"organization_id":          basetypes.StringType{},
		"sort":                     basetypes.StringType{},
		"account_id":               basetypes.StringType{},
		"account_name":             basetypes.StringType{},
		"account_email":            basetypes.StringType{},
		"state":                    basetypes.StringType{},
		"login_id":                 basetypes.StringType{},
	}
}
