package gslb

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/gslb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
		Description: "The GSLB Regional Resource List.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Optional: true,
				Description: "size\n" +
					"  - Example: 20",
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Optional: true,
				Description: "page\n" +
					"  - Example: 0",
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Optional: true,
				Description: "sort\n" +
					"  - Example: created_at:desc",
			},
			common.ToSnakeCase("Region"): schema.StringAttribute{
				Optional:    true,
				Description: "The GSLB Resource Region.",
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Optional:    true,
				Description: "The GSLB Resource Status.",
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Optional:    true,
				Description: "The Name of the gslb.",
			},

			// Output
			common.ToSnakeCase("RegionalGslbs"): schema.ListNestedAttribute{
				Description: "GSLB Routing Control List.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Algorithm"): schema.StringAttribute{
							Description: "The GSLB Algorithm.",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "Created At\n" +
								"  - Example: 2024-05-17T00:23:17Z\n" +
								"  - Format: date-time",
							Computed: true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "Created By\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
							Computed:    true,
						},
						common.ToSnakeCase("EnvUsage"): schema.StringAttribute{
							Description: "The GSLB Environment Usage.",
							Computed:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "ID\n" +
								"  - Example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
							Computed: true,
						},
						common.ToSnakeCase("LinkedRegionalResourceCount"): schema.Int32Attribute{
							Description: "The GSLB Linked Resource Count Per Region.",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "Modified At\n" +
								"  - Example: 2024-05-17T00:23:17Z\n" +
								"  - Format: date-time",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "Modified By\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The Name of the gslb.",
							Computed:    true,
						},
						common.ToSnakeCase("Region"): schema.StringAttribute{
							Description: "The GSLB Resource Region.",
							Computed:    true,
						},
						common.ToSnakeCase("Status"): schema.StringAttribute{
							Description: "The GSLB Resource Status.",
							Computed:    true,
						},
					},
				},
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Computed: true,
				Description: "Count\n" +
					"  - Example: 20",
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
