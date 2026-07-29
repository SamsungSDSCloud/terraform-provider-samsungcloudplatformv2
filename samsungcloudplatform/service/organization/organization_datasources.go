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
	_ datasource.DataSource              = &organizationDataSources{}
	_ datasource.DataSourceWithConfigure = &organizationDataSources{}
)

// NewOrganizationDataSources is a helper function to simplify the provider implementation.
func NewOrganizationDataSources() datasource.DataSource {
	return &organizationDataSources{}
}

type organizationDataSources struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *organizationDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_organizations"
}

// Configure adds the provider configured client to the data source.
func (d *organizationDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Organizations",
		Attributes: map[string]schema.Attribute{
			"organizations": schema.ListNestedAttribute{
				Computed: true,
				Description: "Organization List. \n" +
					"  - example : '[{id: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5, name: My Organization, ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the organization. \n" +
								"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
						},
						"name": schema.StringAttribute{
							Computed: true,
							Description: "Organization Name. \n" +
								"  - example : 'My Organization' \n",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the organization was created. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who created the organization. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the organization was modified. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who modified the organization. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"master_account_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the master account that manages the organization. \n" +
								"  - example : '9f8e7d6c5b4a3z2y1x0w9v8u7t6s5r' \n",
						},
						"delegation_account_id": schema.StringAttribute{
							Computed: true,
							Description: "Delegation Account. \n" +
								"  - example : '124e7d6c5b4a3z2y1x0w9v8u7t6s5r' \n",
						},
						"root_unit_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the root organizational unit. \n" +
								"  - example : 'r-abcd7d6c5b4a3z2y1x0w9v8u7t6s5r' \n",
						},
						"use_scp_yn": schema.BoolAttribute{
							Computed: true,
							Description: "Control Policy Usage YN. \n" +
								"  - example : true \n",
						},
					},
				},
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
				Description: "Total organization count. \n" +
					"  - example : 1 \n",
			},
			"page": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page number. \n" +
					"  - example : 0 \n",
			},
			"size": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page size. \n" +
					"  - example : 20 \n",
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort criteria. \n" +
					"  - example : 'created_at:desc' \n",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort criteria from response. \n" +
					"  - example : ['created_at:desc'] \n",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Description: "Filter by organization name. \n" +
					"  - example : 'My Organization' \n",
			},
			"master_account_id": schema.StringAttribute{
				Optional: true,
				Description: "Filter by master account ID. \n" +
					"  - example : '9f8e7d6c5b4a3z2y1x0w9v8u7t6s5r' \n",
			},
			"organization_id": schema.StringAttribute{
				Optional: true,
				Description: "Filter by organization ID. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
			},
		},
	}
}

func (d *organizationDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationDataSourcesData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiRequest := organization.OrganizationDataSource{
		Size:            data.Size,
		Page:            data.Page,
		Sort:            data.Sort,
		Name:            data.Name,
		MasterAccountId: data.MasterAccountId,
		OrganizationId:  data.OrganizationId,
	}

	result, err := d.client.GetOrganizationList(ctx, apiRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Organizations",
			err.Error(),
		)
		return
	}

	data.Organizations = []organization.OrganizationSummaryItem{}
	for _, org := range result.Organizations {
		data.Organizations = append(data.Organizations, organization.OrganizationSummaryItem{
			Id:                  types.StringValue(org.Id),
			Name:                types.StringValue(org.Name),
			CreatedAt:           types.StringValue(org.CreatedAt.Format(time.RFC3339)),
			CreatedBy:           types.StringValue(org.CreatedBy),
			ModifiedAt:          types.StringValue(org.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:          types.StringValue(org.ModifiedBy),
			MasterAccountId:     types.StringValue(org.MasterAccountId),
			DelegationAccountId: types.StringPointerValue(org.DelegationAccountId.Get()),
			RootUnitId:          types.StringValue(org.RootUnitId),
			UseScpYn:            types.BoolValue(org.UseScpYn),
		})
	}

	data.TotalCount = types.Int64Value(int64(result.Count))
	data.Page = types.Int64Value(int64(result.Page))
	data.Size = types.Int64Value(int64(result.Size))
	if result.Sort != nil {
		sortVals, _ := types.ListValueFrom(ctx, types.StringType, result.Sort)
		data.SortResult = sortVals
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type OrganizationDataSourcesData struct {
	Organizations   []organization.OrganizationSummaryItem `tfsdk:"organizations"`
	TotalCount      types.Int64                            `tfsdk:"total_count"`
	Page            types.Int64                            `tfsdk:"page"`
	Size            types.Int64                            `tfsdk:"size"`
	Sort            types.String                           `tfsdk:"sort"`
	SortResult      types.List                             `tfsdk:"sort_result"`
	Name            types.String                           `tfsdk:"name"`
	MasterAccountId types.String                           `tfsdk:"master_account_id"`
	OrganizationId  types.String                           `tfsdk:"organization_id"`
}
