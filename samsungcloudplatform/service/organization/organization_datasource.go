package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &organizationDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationDataSource{}
)

// NewOrganizationDataSource is a helper function to simplify the provider implementation.
func NewOrganizationDataSource() datasource.DataSource {
	return &organizationDataSource{}
}

type organizationDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *organizationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Configure adds the provider configured client to the data source.
func (d *organizationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Organization ID",
				MarkdownDescription: "Organization ID",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Organization Name",
				MarkdownDescription: "Organization Name",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Created At",
				MarkdownDescription: "Created At",
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Created By",
				MarkdownDescription: "Created By",
			},
			"master_account_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Master Account ID",
				MarkdownDescription: "Master Account ID",
			},
			"master_account_email": schema.StringAttribute{
				Computed:            true,
				Description:         "Master Account Email",
				MarkdownDescription: "Master Account Email",
			},
			"delegation_account_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Delegation Account ID",
				MarkdownDescription: "Delegation Account ID",
			},
			"root_unit_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Root Unit ID",
				MarkdownDescription: "Root Unit ID",
			},
			"srn": schema.StringAttribute{
				Computed:            true,
				Description:         "SRN",
				MarkdownDescription: "SRN",
			},
			"use_scp_yn": schema.BoolAttribute{
				Computed:            true,
				Description:         "Use SCP",
				MarkdownDescription: "Use SCP",
			},
			"creator_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Creator Name",
				MarkdownDescription: "Creator Name",
			},
			"modified_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified At",
				MarkdownDescription: "Modified At",
			},
			"modified_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified By",
				MarkdownDescription: "Modified By",
			},
			"modifier_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Modifier Name",
				MarkdownDescription: "Modifier Name",
			},
		},
	}
}

func (d *organizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data organization.OrganizationResource

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := data.Id.ValueString()

	result, err := d.client.GetOrganization(ctx, organizationId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Organization",
			err.Error(),
		)
		return
	}

	org := result.Organization

	data.Id = types.StringValue(org.Id)
	data.Name = types.StringValue(org.Name)
	data.CreatedAt = types.StringValue(org.CreatedAt.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(org.CreatedBy)
	data.MasterAccountId = types.StringValue(org.MasterAccountId)
	data.MasterAccountEmail = types.StringValue(org.MasterAccountEmail)
	data.DelegationAccountId = types.StringPointerValue(org.DelegationAccountId.Get())
	data.RootUnitId = types.StringValue(org.RootUnitId)
	data.Srn = types.StringValue(org.Srn)
	data.UseScpYn = types.BoolValue(org.UseScpYn)
	data.CreatorName = types.StringValue(org.GetCreatorName())
	data.ModifiedAt = types.StringValue(org.ModifiedAt.Format(time.RFC3339))
	data.ModifiedBy = types.StringValue(org.ModifiedBy)
	data.ModifierName = types.StringValue(org.GetModifierName())

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
