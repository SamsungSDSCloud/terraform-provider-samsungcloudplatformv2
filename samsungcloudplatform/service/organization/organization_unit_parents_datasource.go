package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &organizationUnitParentsDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationUnitParentsDataSource{}
)

// NewOrganizationUnitParentsDataSource is a helper function to simplify the provider implementation.
func NewOrganizationUnitParentsDataSource() datasource.DataSource {
	return &organizationUnitParentsDataSource{}
}

type organizationUnitParentsDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *organizationUnitParentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_unit_parents"
}

// Configure adds the provider configured client to the data source.
func (d *organizationUnitParentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationUnitParentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List Parent Organization Units",
		Attributes: map[string]schema.Attribute{
			"unit_id": schema.StringAttribute{
				Required: true,
				Description: "Organization Unit ID. \n" +
					"  - example : 'ou-c29a138f8f1d78e24dbfa8681fc2fc8' \n",
			},
			"parents": schema.ListNestedAttribute{
				Computed: true,
				Description: "Parent Organization Unit Info. \n" +
					"  - example : '[{id: ou-c29a138f8f1d78e24dbfa8681fc2fc8, name: score-organization-unit, ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Organization Unit ID. \n" +
								"  - example : 'ou-c29a138f8f1d78e24dbfa8681fc2fc8' \n",
						},
						"name": schema.StringAttribute{
							Computed: true,
							Description: "Organization Unit Name. \n" +
								"  - example : 'score-organization-unit' \n",
						},
						"depth": schema.Int64Attribute{
							Computed: true,
							Description: "Hierarchy (0~5). \n" +
								"  - example : 1 \n",
						},
						"type": schema.StringAttribute{
							Computed: true,
							Description: "Type (Root or Organization Unit). \n" +
								"  - example : 'OU' \n",
						},
						"parent_unit_id": schema.StringAttribute{
							Computed: true,
							Description: "Root or Parent Organization Unit ID. \n" +
								"  - example : 'ou-fc8c29a138d78e24bf1fa86812fc8b' \n",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the organization unit was created. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who created the organization unit. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the organization unit was modified. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who modified the organization unit. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"service_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the service to which the policy applies. \n" +
								"  - example : 'Organization' \n",
						},
						"description": schema.StringAttribute{
							Computed: true,
							Description: "Organization Unit Description. \n" +
								"  - example : 'Score Organization Unit' \n",
						},
					},
				},
			},
		},
	}
}

func (d *organizationUnitParentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationUnitParentsDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetOrganizationUnitParents(
		ctx,
		data.UnitId.ValueString(),
		"",
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Organization Unit Parents",
			err.Error(),
		)
		return
	}

	parents := make([]OrganizationUnitParentItem, 0)
	if result != nil && result.Parents != nil {
		for _, parent := range result.Parents {
			createdAt := ""
			if !parent.CreatedAt.IsZero() {
				createdAt = parent.CreatedAt.Format(time.RFC3339)
			}

			modifiedAt := ""
			if !parent.ModifiedAt.IsZero() {
				modifiedAt = parent.ModifiedAt.Format(time.RFC3339)
			}

			description := ""
			if parent.Description.IsSet() && !IsNil(parent.Description.Get()) {
				description = *parent.Description.Get()
			}

			item := OrganizationUnitParentItem{
				Id:           types.StringValue(parent.Id),
				Name:         types.StringValue(parent.Name),
				Depth:        types.Int64Value(int64(parent.Depth)),
				Type:         types.StringValue(parent.Type),
				ParentUnitId: types.StringValue(parent.GetParentUnitId()),
				CreatedAt:    types.StringValue(createdAt),
				CreatedBy:    types.StringValue(parent.CreatedBy),
				ModifiedAt:   types.StringValue(modifiedAt),
				ModifiedBy:   types.StringValue(parent.ModifiedBy),
				ServiceName:  types.StringValue(parent.ServiceName),
				Description:  types.StringValue(description),
			}
			parents = append(parents, item)
		}
	}

	data.Parents = parents

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type OrganizationUnitParentsDataSourceModel struct {
	UnitId  types.String                 `tfsdk:"unit_id"`
	Parents []OrganizationUnitParentItem `tfsdk:"parents"`
}

// OrganizationUnitParentItem is the item in the list
type OrganizationUnitParentItem struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Depth        types.Int64  `tfsdk:"depth"`
	Type         types.String `tfsdk:"type"`
	ParentUnitId types.String `tfsdk:"parent_unit_id"`
	CreatedAt    types.String `tfsdk:"created_at"`
	CreatedBy    types.String `tfsdk:"created_by"`
	ModifiedAt   types.String `tfsdk:"modified_at"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
	ServiceName  types.String `tfsdk:"service_name"`
	Description  types.String `tfsdk:"description"`
}
