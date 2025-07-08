package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ListAccountMember
// [GET] /v1/cloudmonitorings/product/v1/accounts/members
type cloudMonitoringAccountMemberDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

var (
	_ datasource.DataSource              = &cloudMonitoringAccountMemberDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringAccountMemberDataSources{}
)

func NewCloudMonitoringAccountMemberSources() datasource.DataSource {
	return &cloudMonitoringAccountMemberDataSources{}
}

func (d *cloudMonitoringAccountMemberDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_accountmembers"
}

func (d *cloudMonitoringAccountMemberDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Schema of cloudMonitoringAccountMemberDataSources.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("AccountMembers"): schema.ListNestedAttribute{
				Description: "AccountMembers",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CompanyName"): schema.StringAttribute{
							Description: "CompanyName",
							Optional:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "CreatedBy",
							Optional:    true,
						},
						common.ToSnakeCase("CreatedByEmail"): schema.StringAttribute{
							Description: "CreatedByEmail",
							Optional:    true,
						},
						common.ToSnakeCase("CreatedByName"): schema.StringAttribute{
							Description: "CreatedByName",
							Optional:    true,
						},
						common.ToSnakeCase("CreatedDt"): schema.StringAttribute{
							Description: "CreatedDt",
							Optional:    true,
						},
						common.ToSnakeCase("Email"): schema.StringAttribute{
							Description: "Email",
							Optional:    true,
						},
						common.ToSnakeCase("LastAccessDt"): schema.StringAttribute{
							Description: "LastAccessDt",
							Optional:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "ModifiedBy",
							Optional:    true,
						},
						common.ToSnakeCase("ModifiedByEmail"): schema.StringAttribute{
							Description: "ModifiedByEmail",
							Optional:    true,
						},
						common.ToSnakeCase("ModifiedByName"): schema.StringAttribute{
							Description: "ModifiedByName",
							Optional:    true,
						},
						common.ToSnakeCase("ModifiedDt"): schema.StringAttribute{
							Description: "ModifiedDt",
							Optional:    true,
						},
						common.ToSnakeCase("OrganizationId"): schema.StringAttribute{
							Description: "OrganizationId",
							Optional:    true,
						},
						common.ToSnakeCase("PositionName"): schema.StringAttribute{
							Description: "PositionName",
							Optional:    true,
						},
						common.ToSnakeCase("UserGroupCount"): schema.StringAttribute{
							Description: "UserGroupCount",
							Optional:    true,
						},
						common.ToSnakeCase("UserId"): schema.StringAttribute{
							Description: "UserId",
							Optional:    true,
						},
						common.ToSnakeCase("UserName"): schema.StringAttribute{
							Description: "UserName",
							Optional:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *cloudMonitoringAccountMemberDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.CloudMonitoring
	d.clients = inst.Client
}

func (d *cloudMonitoringAccountMemberDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudmonitoring.AccountMemberDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.GetAccountMembers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, element := range res.Contents {
		resState := cloudmonitoring.AccountMember{
			CompanyName:     types.StringValue(element.GetCompanyName()),
			CreatedBy:       types.StringValue(element.GetCreatedBy()),
			CreatedByEmail:  types.StringValue(element.GetCreatedByEmail()),
			CreatedByName:   types.StringValue(element.GetCreatedByName()),
			CreatedDt:       types.StringValue(element.GetCreatedDt()),
			Email:           types.StringValue(element.GetEmail()),
			LastAccessDt:    types.StringValue(element.GetLastAccessDt()),
			ModifiedBy:      types.StringValue(element.GetModifiedBy()),
			ModifiedByEmail: types.StringValue(element.GetModifiedByEmail()),
			ModifiedByName:  types.StringValue(element.GetModifiedByName()),
			ModifiedDt:      types.StringValue(element.GetModifiedDt()),
			OrganizationId:  types.StringValue(element.GetOrganizationId()),
			PositionName:    types.StringValue(element.GetPositionName()),
			UserGroupCount:  types.StringValue(element.GetUserGroupCount()),
			UserId:          types.StringValue(element.GetUserId()),
			UserName:        types.StringValue(element.GetUserName()),
		}

		state.AccountMembers = append(state.AccountMembers, resState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
