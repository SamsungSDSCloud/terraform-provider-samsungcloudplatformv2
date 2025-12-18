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
	_ datasource.DataSource              = &gslbGslbResourceDataSources{}
	_ datasource.DataSourceWithConfigure = &gslbGslbResourceDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewGslbGslbResourceDataSources() datasource.DataSource {
	return &gslbGslbResourceDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type gslbGslbResourceDataSources struct {
	config  *scpsdk.Configuration
	client  *gslb.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *gslbGslbResourceDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gslb_gslb_resources" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *gslbGslbResourceDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Get List of Gslb Resource.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size\n" +
					"  - Example: 20",
				Optional: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page\n" +
					"  - Example: 0",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort\n" +
					"  - Example: created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("GslbId"): schema.StringAttribute{
				Description: "GslbId",
				Required:    true,
			},

			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Count\n" +
					"  - Example: 20",
				Computed: true,
			},
			common.ToSnakeCase("GslbResources"): schema.ListNestedAttribute{
				Description: "A list of gslb resource.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "Created at\n" +
								"  - Example: 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "Created by\n" +
								"  - Example: 90dfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description\n" +
								"  - Max length:50",
							Computed: true,
						},
						common.ToSnakeCase("Destination"): schema.StringAttribute{
							Description: "The GSLB Resource Destination.",
							Computed:    true,
						},
						common.ToSnakeCase("HealthCheckStatus"): schema.StringAttribute{
							Description: "The GSLB Resource Health Check Status." +
								"  - Example: CONNECTED | DISCONNECTED",
							Computed: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "ID\n" +
								"  - Example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "Modified at\n" +
								"  - Example: 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "Modified by\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Region"): schema.StringAttribute{
							Description: "The GSLB Resource Region.",
							Computed:    true,
						},
						common.ToSnakeCase("Status"): schema.StringAttribute{
							Description: "The GSLB Resource Status.\n" +
								"  - Example: ENABLE | DISABLE",
							Computed: true,
						},
						common.ToSnakeCase("Weight"): schema.Int32Attribute{
							Description: "The GSLB Resource Weight.\n" +
								"  - Min: 0\n" +
								"  - Max: 100\n" +
								"  - Default: 0",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *gslbGslbResourceDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *gslbGslbResourceDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state gslb.GslbResourceDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetGslbResourceList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Gslb Resources",
			err.Error(),
		)
		return
	}

	for _, gslbResource := range data.GslbResources {
		gslbResourceState := gslb.GslbResourceDetail{
			CreatedAt:   types.StringValue(gslbResource.CreatedAt.Format(time.RFC3339)),
			CreatedBy:   types.StringValue(gslbResource.CreatedBy),
			Description: virtualserverutil.ToNullableStringValue(gslbResource.Description.Get()),
			Destination: types.StringValue(gslbResource.Destination),
			Id:          types.StringValue(gslbResource.Id),
			ModifiedAt:  types.StringValue(gslbResource.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:  types.StringValue(gslbResource.ModifiedBy),
			Region:      types.StringValue(gslbResource.Region),
			Weight:      common.ToNullableInt32Value(gslbResource.Weight.Get()),
		}

		if gslbResource.HealthCheckStatus.IsValid() {
			gslbResourceState.HealthCheckStatus = types.StringValue(string(gslbResource.HealthCheckStatus))
		}

		if gslbResource.Status.IsValid() {
			gslbResourceState.Status = types.StringValue(string(gslbResource.Status))
		}

		state.GslbResources = append(state.GslbResources, gslbResourceState)
		state.TotalCount = types.Int32Value(data.Count)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}
