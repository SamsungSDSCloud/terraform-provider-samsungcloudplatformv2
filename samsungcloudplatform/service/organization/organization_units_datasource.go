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
	_ datasource.DataSource              = &organizationUnitsDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationUnitsDataSource{}
)

// NewOrganizationUnitsDataSource is a helper function to simplify the provider implementation.
func NewOrganizationUnitsDataSource() datasource.DataSource {
	return &organizationUnitsDataSource{}
}

type organizationUnitsDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *organizationUnitsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_units"
}

// Configure adds the provider configured client to the data source.
func (d *organizationUnitsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationUnitsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List Organization Units",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID",
				MarkdownDescription: "ID",
			},
			"parent_unit_id": schema.StringAttribute{
				Required:            true,
				Description:         "Parent Organization Unit ID",
				MarkdownDescription: "Parent Organization Unit ID",
			},
			"organization_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Organization ID",
				MarkdownDescription: "Organization ID",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Organization Unit Name",
				MarkdownDescription: "Organization Unit Name",
			},
			"exclude_policy_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Policy ID to Exclude",
				MarkdownDescription: "Policy ID to Exclude",
			},
			"organization_units": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "A list of Organization Units",
				MarkdownDescription: "A list of Organization Units",
				NestedObject: schema.NestedAttributeObject{
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
						"email": schema.StringAttribute{
							Computed:            true,
							Description:         "Email",
							MarkdownDescription: "Email",
						},
						"depth": schema.Int64Attribute{
							Computed:            true,
							Description:         "Hierarchy (0~5)",
							MarkdownDescription: "Hierarchy (0~5)",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "Type (Root or Organization Unit)",
							MarkdownDescription: "Type (Root or Organization Unit)",
						},
						"organization_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Organization ID",
							MarkdownDescription: "Organization ID",
						},
						"parent_unit_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Parent Organization Unit ID",
							MarkdownDescription: "Parent Organization Unit ID",
						},
						"parent_unit_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Parent Organization Unit Name",
							MarkdownDescription: "Parent Organization Unit Name",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Created At",
							MarkdownDescription: "Created At",
						},
						"joined_method": schema.StringAttribute{
							Computed:            true,
							Description:         "Joined Method",
							MarkdownDescription: "Joined Method",
						},
						"joined_time": schema.StringAttribute{
							Computed:            true,
							Description:         "Joined Time",
							MarkdownDescription: "Joined Time",
						},
						"login_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Login ID",
							MarkdownDescription: "Login ID",
						},
					},
				},
			},
		},
	}
}

func (d *organizationUnitsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state OrganizationUnitsDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentUnitId := state.ParentUnitId.ValueString()
	if parentUnitId == "" {
		parentUnitId = "root"
	}

	result, err := d.client.GetOrganizationUnits(
		ctx,
		parentUnitId,
		state.OrganizationId.ValueString(),
		state.Name.ValueString(),
		state.ExcludePolicyId.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Organization Units",
			err.Error(),
		)
		return
	}

	state.Id = types.StringValue("organization_units")
	state.ParentUnitId = types.StringValue(parentUnitId)

	organizationUnits := make([]OrganizationUnitItem, 0)
	if result != nil && result.OrganizationUnits != nil {
		for _, unit := range result.OrganizationUnits {
			createdAt := unit.GetCreatedAt()
			joinedTime := unit.GetJoinedTime()

			item := OrganizationUnitItem{
				Id:             types.StringValue(unit.Id),
				Name:           types.StringValue(unit.GetName()),
				Email:          types.StringValue(unit.GetEmail()),
				Depth:          types.Int64Value(int64(unit.GetDepth())),
				Type:           types.StringValue(unit.Type),
				OrganizationId: types.StringValue(unit.OrganizationId),
				ParentUnitId:   types.StringValue(unit.GetParentUnitId()),
				ParentUnitName: types.StringValue(unit.GetParentUnitName()),
				CreatedAt:      types.StringValue(createdAt.Format(time.RFC3339)),
				JoinedMethod:   types.StringValue(string(unit.GetJoinedMethod())),
				JoinedTime:     types.StringValue(joinedTime.Format(time.RFC3339)),
				LoginId:        types.StringValue(unit.GetLoginId()),
			}
			organizationUnits = append(organizationUnits, item)
		}
	}

	state.OrganizationUnits = organizationUnits

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type OrganizationUnitsDataSource struct {
	Id                types.String           `tfsdk:"id"`
	ParentUnitId      types.String           `tfsdk:"parent_unit_id"`
	OrganizationId    types.String           `tfsdk:"organization_id"`
	Name              types.String           `tfsdk:"name"`
	ExcludePolicyId   types.String           `tfsdk:"exclude_policy_id"`
	OrganizationUnits []OrganizationUnitItem `tfsdk:"organization_units"`
}

// OrganizationUnitItem is the item in the list
type OrganizationUnitItem struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Email          types.String `tfsdk:"email"`
	Depth          types.Int64  `tfsdk:"depth"`
	Type           types.String `tfsdk:"type"`
	OrganizationId types.String `tfsdk:"organization_id"`
	ParentUnitId   types.String `tfsdk:"parent_unit_id"`
	ParentUnitName types.String `tfsdk:"parent_unit_name"`
	CreatedAt      types.String `tfsdk:"created_at"`
	JoinedMethod   types.String `tfsdk:"joined_method"`
	JoinedTime     types.String `tfsdk:"joined_time"`
	LoginId        types.String `tfsdk:"login_id"`
}

// IsNil checks if the pointer is nil
func IsNil(ptr interface{}) bool {
	if ptr == nil {
		return true
	}
	switch v := ptr.(type) {
	case *string:
		return v == nil
	case *int32:
		return v == nil
	case *time.Time:
		return v == nil
	default:
		return false
	}
}
