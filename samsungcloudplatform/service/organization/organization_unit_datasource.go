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
	_ datasource.DataSource              = &organizationUnitDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationUnitDataSource{}
)

// NewOrganizationUnitDataSource is a helper function to simplify the provider implementation.
func NewOrganizationUnitDataSource() datasource.DataSource {
	return &organizationUnitDataSource{}
}

type organizationUnitDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *organizationUnitDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_unit"
}

// Configure adds the provider configured client to the data source.
func (d *organizationUnitDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationUnitDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Organization Unit",
		Attributes: map[string]schema.Attribute{
			"unit_id": schema.StringAttribute{
				Required:            true,
				Description:         "Organization Unit ID",
				MarkdownDescription: "Organization Unit ID",
			},
			"organization_unit": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Organization Unit Info",
				MarkdownDescription: "Organization Unit Info",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Organization Unit ID",
						MarkdownDescription: "Organization Unit ID",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Organization Unit Name",
						MarkdownDescription: "Organization Unit Name",
					},
					"description": schema.StringAttribute{
						Computed:            true,
						Description:         "Organization Unit Description",
						MarkdownDescription: "Organization Unit Description",
					},
					"parent_unit_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Parent Organization Unit ID",
						MarkdownDescription: "Parent Organization Unit ID",
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
					"depth": schema.Int64Attribute{
						Computed:            true,
						Description:         "Hierarchy (0~5)",
						MarkdownDescription: "Hierarchy (0~5)",
					},
					"service_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Service Name",
						MarkdownDescription: "Service Name",
					},
					"srn": schema.StringAttribute{
						Computed:            true,
						Description:         "SRN",
						MarkdownDescription: "SRN",
					},
					"type": schema.StringAttribute{
						Computed:            true,
						Description:         "Type (Root or Organization Unit)",
						MarkdownDescription: "Type (Root or Organization Unit)",
					},
				},
			},
		},
	}
}

func (d *organizationUnitDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationUnitDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	unitId := data.UnitId.ValueString()

	result, err := d.client.GetOrganizationUnit(ctx, unitId, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Organization Unit",
			err.Error(),
		)
		return
	}

	unit := result.OrganizationUnit

	data.OrganizationUnit = &OrganizationUnitInfo{
		Id:           types.StringValue(unit.Id),
		Name:         types.StringValue(unit.Name),
		Description:  types.StringValue(unit.GetDescription()),
		ParentUnitId: types.StringValue(unit.GetParentUnitId()),
		CreatedAt:    types.StringValue(unit.CreatedAt.Format(time.RFC3339)),
		CreatedBy:    types.StringValue(unit.CreatedBy),
		CreatorName:  types.StringValue(unit.GetCreatorName()),
		ModifiedAt:   types.StringValue(unit.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:   types.StringValue(unit.ModifiedBy),
		ModifierName: types.StringValue(unit.GetModifierName()),
		Depth:        types.Int64Value(int64(unit.Depth)),
		ServiceName:  types.StringValue(unit.ServiceName),
		Srn:          types.StringValue(unit.GetSrn()),
		Type:         types.StringValue(unit.Type),
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type OrganizationUnitDataSourceModel struct {
	UnitId           types.String          `tfsdk:"unit_id"`
	OrganizationUnit *OrganizationUnitInfo `tfsdk:"organization_unit"`
}

// OrganizationUnitInfo contains the organization unit details
type OrganizationUnitInfo struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	ParentUnitId types.String `tfsdk:"parent_unit_id"`
	CreatedAt    types.String `tfsdk:"created_at"`
	CreatedBy    types.String `tfsdk:"created_by"`
	CreatorName  types.String `tfsdk:"creator_name"`
	ModifiedAt   types.String `tfsdk:"modified_at"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
	ModifierName types.String `tfsdk:"modifier_name"`
	Depth        types.Int64  `tfsdk:"depth"`
	ServiceName  types.String `tfsdk:"service_name"`
	Srn          types.String `tfsdk:"srn"`
	Type         types.String `tfsdk:"type"`
}
