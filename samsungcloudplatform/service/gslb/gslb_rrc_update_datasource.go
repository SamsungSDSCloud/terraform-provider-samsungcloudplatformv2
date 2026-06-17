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
	_ datasource.DataSource              = &gslbGslbRRCUpdateDataSource{}
	_ datasource.DataSourceWithConfigure = &gslbGslbRRCUpdateDataSource{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewgslbGslbRRCUpdateDataSource() datasource.DataSource {
	return &gslbGslbRRCUpdateDataSource{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type gslbGslbRRCUpdateDataSource struct {
	config  *scpsdk.Configuration
	client  *gslb.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *gslbGslbRRCUpdateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gslb_gslb_rrc_update" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *gslbGslbRRCUpdateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Updates GSLB regional routing control settings.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("GslbId"): schema.StringAttribute{
				Description: "The unique identifier of the GSLB.\n" +
					"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				Required: true,
			},
			common.ToSnakeCase("Region"): schema.StringAttribute{
				Description: "The region where the GSLB resource is located.\n" +
					"  - example : kr-west1",
				Required: true,
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Description: "Whether to use the GSLB resource. (e.g., ENABLE, DISABLE)\n" +
					"  - example : ENABLE",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ENABLE", "DISABLE"),
				},
			},

			// Output
			common.ToSnakeCase("RegionalGslb"): schema.SingleNestedAttribute{
				Description: "Details of the updated GSLB regional routing control.",
				Computed:    true,
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
	}
}

// Configure adds the provider configured client to the data source.
func (d *gslbGslbRRCUpdateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *gslbGslbRRCUpdateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state gslb.GslbRegionalRoutingControlUpdateDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.UpdateGslbRegionalRoutingControl(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Gslb Regional Routing Control.",
			err.Error(),
		)
		return
	}

	regionalGslbModelData := gslb.GslbRoutingControl{
		Algorithm:                   types.StringValue(data.RegionalGslb.Algorithm),
		CreatedAt:                   types.StringValue(data.RegionalGslb.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                   types.StringValue(data.RegionalGslb.CreatedBy),
		Description:                 virtualserverutil.ToNullableStringValue(data.RegionalGslb.Description.Get()),
		EnvUsage:                    types.StringValue(data.RegionalGslb.EnvUsage),
		Id:                          types.StringValue(data.RegionalGslb.Id),
		LinkedRegionalResourceCount: types.Int32Value(data.RegionalGslb.LinkedRegionalResourceCount),
		ModifiedAt:                  types.StringValue(data.RegionalGslb.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:                  types.StringValue(data.RegionalGslb.ModifiedBy),
		Name:                        types.StringValue(data.RegionalGslb.Name),
		Region:                      types.StringValue(data.RegionalGslb.Region),
		Status:                      types.StringValue(data.RegionalGslb.Status),
	}

	gslbObjectValue, _ := types.ObjectValueFrom(ctx, regionalGslbModelData.AttributeTypes(), regionalGslbModelData)

	state.RegionalGslb = gslbObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
