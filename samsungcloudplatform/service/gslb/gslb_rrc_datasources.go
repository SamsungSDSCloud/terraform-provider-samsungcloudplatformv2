package gslb

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/gslb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &gslbGslbRRCDataSources{}
	_ datasource.DataSourceWithConfigure = &gslbGslbRRCDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewGslbGslbRRCDataSources() datasource.DataSource {
	return &gslbGslbRRCDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type gslbGslbRRCDataSources struct {
	config  *scpsdk.Configuration
	client  *gslb.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *gslbGslbRRCDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gslb_gslb_rrcs" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *gslbGslbRRCDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Gets a list of GSLB regional routing controls.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Optional: true,
				Description: "The number of items per page.\n" +
					"  - example : 20",
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Optional: true,
				Description: "The page number for pagination.\n" +
					"  - example : 0",
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Optional: true,
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order.\n" +
					"  - example : createdAt:asc",
			},
			common.ToSnakeCase("Region"): schema.StringAttribute{
				Optional: true,
				Description: "The region to filter GSLB resources.\n" +
					"  - example : kr-west1",
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Optional: true,
				Description: "Whether to use the GSLB resource. (e.g., ENABLE, DISABLE)\n" +
					"  - example : ENABLE",
				Validators: []validator.String{
					stringvalidator.OneOf("ENABLE", "DISABLE"),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Optional: true,
				Description: "The name to filter GSLB resources.\n" +
					"  - example : example.gslb.e.samsungsdscloud.com",
			},

			// Output
			common.ToSnakeCase("RegionalGslbs"): schema.ListNestedAttribute{
				Description: "A list of GSLB regional routing controls.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Algorithm"): schema.StringAttribute{
							Description: "The load balancing algorithm for GSLB traffic distribution (e.g., ROUND_ROBIN, RATIO).\n" +
								"  - example : ROUND_ROBIN",
							Computed: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user id that created the resource.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
								"  - example : Example Description for GSLB",
							Computed: true,
						},
						common.ToSnakeCase("EnvUsage"): schema.StringAttribute{
							Description: "The environment usage type for the GSLB (e.g., PUBLIC).\n" +
								"  - example : PUBLIC",
							Computed: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the GSLB.\n" +
								"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
							Computed: true,
						},
						common.ToSnakeCase("LinkedRegionalResourceCount"): schema.Int32Attribute{
							Description: "The number of resources linked to this GSLB in the region.\n" +
								"  - example : 2",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user id that last modified the resource.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the GSLB.\n" +
								"  - example : example.gslb.e.samsungsdscloud.com",
							Computed: true,
						},
						common.ToSnakeCase("Region"): schema.StringAttribute{
							Description: "The region where the GSLB resource is located.\n" +
								"  - example : kr-west1",
							Computed: true,
						},
						common.ToSnakeCase("Status"): schema.StringAttribute{
							Description: "Whether to use the GSLB resource. (e.g., ENABLE, DISABLE)\n" +
								"  - example : ENABLE",
							Computed: true,
						},
					},
				},
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Computed: true,
				Description: "The total number of GSLB regional routing controls.\n" +
					"  - example : 2",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *gslbGslbRRCDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	d.client = inst.Client.Gslb
	d.clients = inst.Client
}

func (d *gslbGslbRRCDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state gslb.GslbRegionalRoutingControlListDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetGslbRegionalRoutingControlList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Gslb Regional Routing Control list.",
			err.Error(),
		)
		return
	}

	state.TotalCount = types.Int32Value(data.Count)
	state.Size = types.Int32Value(data.Size)
	state.Page = types.Int32Value(data.Page)
	for _, gslbRrc := range data.RegionalGslbs {
		gslbState := gslb.GslbRoutingControl{
			Algorithm:                   types.StringValue(gslbRrc.Algorithm),
			CreatedAt:                   types.StringValue(gslbRrc.CreatedAt.Format(time.RFC3339)),
			CreatedBy:                   types.StringValue(gslbRrc.CreatedBy),
			Description:                 virtualserverutil.ToNullableStringValue(gslbRrc.Description.Get()),
			EnvUsage:                    types.StringValue(gslbRrc.EnvUsage),
			Id:                          types.StringValue(gslbRrc.Id),
			LinkedRegionalResourceCount: types.Int32Value(gslbRrc.LinkedRegionalResourceCount),
			ModifiedAt:                  types.StringValue(gslbRrc.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:                  types.StringValue(gslbRrc.ModifiedBy),
			Name:                        types.StringValue(gslbRrc.Name),
			Region:                      types.StringValue(gslbRrc.Region),
			Status:                      types.StringValue(gslbRrc.Status),
		}

		state.RegionalGslbs = append(state.RegionalGslbs, gslbState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}
