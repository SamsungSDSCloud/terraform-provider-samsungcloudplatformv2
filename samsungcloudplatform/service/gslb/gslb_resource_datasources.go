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
		Description: "Gets a list of resources linked to a GSLB.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example : 20",
				Optional: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination.\n" +
					"  - example : 0",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order.\n" +
					"  - example : createdAt:asc",
				Optional: true,
			},
			common.ToSnakeCase("GslbId"): schema.StringAttribute{
				Description: "The unique identifier of the GSLB.\n" +
					"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				Required: true,
			},

			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "The total number of GSLB resources.\n" +
					"  - example : 2",
				Computed: true,
			},
			common.ToSnakeCase("GslbResources"): schema.ListNestedAttribute{
				Description: "A list of resources linked to the GSLB.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
								"  - example : Example Description for GSLB Resource",
							Computed: true,
						},
						common.ToSnakeCase("Destination"): schema.StringAttribute{
							Description: "The destination endpoint for the GSLB resource.\n" +
								"  - example : 192.168.1.100",
							Computed: true,
						},
						common.ToSnakeCase("HealthCheckStatus"): schema.StringAttribute{
							Description: "The health check status of the GSLB resource (e.g., CONNECTED, DISCONNECTED).\n" +
								"  - example : CONNECTED",
							Computed: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the GSLB resource.\n" +
								"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
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
						common.ToSnakeCase("Weight"): schema.Int32Attribute{
							Description: "The weight for load balancing distribution (0-100).\n" +
								"  - example : 50",
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
