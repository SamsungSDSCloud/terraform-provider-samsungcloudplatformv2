package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
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
				Optional: true,
				Description: "Filter by Organization ID. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort criteria (e.g., 'created_at:desc'). \n" +
					"  - example : 'created_at:desc' \n",
			},
			"account_id": schema.StringAttribute{
				Optional: true,
				Description: "Filter by Account ID. \n" +
					"  - example : 'b4d3f2h1j0l9n8p7r6t5v4x3z2y1w0' \n",
			},
			"account_name": schema.StringAttribute{
				Optional: true,
				Description: "Filter by Account Name. \n" +
					"  - example : 'score-account' \n",
			},
			"account_email": schema.StringAttribute{
				Optional: true,
				Description: "Filter by Account Email. \n" +
					"  - example : 'score@samsung.com' \n",
			},
			"state": schema.StringAttribute{
				Optional: true,
				Description: "Filter by Invitation State. \n" +
					"  - example : 'INVITING' \n" +
					"  - allowed_values : ['INVITING', 'REFUSED', 'INVITED', 'CANCELED', 'EXPIRED'] \n",
			},
			"login_id": schema.StringAttribute{
				Optional: true,
				Description: "Filter by Login ID. \n" +
					"  - example : 'log-archive@samsung.com' \n",
			},
			"organization_invitations": schema.ListNestedAttribute{
				Computed: true,
				Description: "Invitation List. \n" +
					"  - example : '[{id: 0a36e0746dbf4908acf0357829701381, organization_id: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5, ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the invitation. \n" +
								"  - example : '0a36e0746dbf4908acf0357829701381' \n",
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the organization. \n" +
								"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
						},
						"account_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the account. \n" +
								"  - example : 'score-account' \n",
						},
						"account_email": schema.StringAttribute{
							Computed: true,
							Description: "Email address of the account. \n" +
								"  - example : 'score@samsung.com' \n",
						},
						"login_id": schema.StringAttribute{
							Computed: true,
							Description: "Login ID of the account. \n" +
								"  - example : 'log-archive@samsung.com' \n",
						},
						"target_account_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the target account. \n" +
								"  - example : '338c2fc8c29a449dbfa8681f8f1d78e5' \n",
						},
						"state": schema.StringAttribute{
							Computed: true,
							Description: "State of the invitation. \n" +
								"  - example : 'INVITING' \n" +
								"  - allowed_values : ['INVITING', 'REFUSED', 'INVITED', 'CANCELED', 'EXPIRED'] \n",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the invitation was created. \n" +
								"  - example : '2025-01-01T00:00:00Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who created the invitation. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the invitation was last modified. \n" +
								"  - example : '2025-01-01T00:00:00Z' \n",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who last modified the invitation. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
					},
				},
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
				Description: "Total organization invitation count. \n" +
					"  - example : 10 \n",
			},
			"page": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page Number. \n" +
					"  - example : 0 \n",
			},
			"size": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page Size. \n" +
					"  - example : 20 \n",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort Criteria. \n" +
					"  - example : ['created_at:desc'] \n",
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
