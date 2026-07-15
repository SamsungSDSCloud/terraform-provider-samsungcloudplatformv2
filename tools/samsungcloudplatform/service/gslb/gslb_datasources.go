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
	_ datasource.DataSource              = &gslbGslbDataSources{}
	_ datasource.DataSourceWithConfigure = &gslbGslbDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewGslbGslbDataSources() datasource.DataSource {
	return &gslbGslbDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type gslbGslbDataSources struct {
	config  *scpsdk.Configuration
	client  *gslb.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *gslbGslbDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gslb_gslbs" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *gslbGslbDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Gets a list of Global Server Load Balancers.",
		Attributes: map[string]schema.Attribute{
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
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The state to filter GSLBs (e.g., ACTIVE, CREATING, EDITING, ERROR, DELETING).\n" +
					"  - example : ACTIVE",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "CREATING", "EDITING", "ERROR", "DELETING"),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name to filter GSLBs.\n" +
					"  - example : example.gslb.e.samsungsdscloud.com",
				Optional: true,
			},
			common.ToSnakeCase("Gslbs"): schema.ListNestedAttribute{
				Description: "A list of Global Server Load Balancers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Algorithm"): schema.StringAttribute{
							Description: "The load balancing algorithm for GSLB traffic distribution (e.g., ROUND_ROBIN, RATIO).\n" +
								"  - example : ROUND_ROBIN",
							Optional: true,
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
							Optional: true,
						},
						common.ToSnakeCase("EnvUsage"): schema.StringAttribute{
							Description: "The environment usage type for the GSLB (e.g., PUBLIC).\n" +
								"  - example : PUBLIC",
							Optional: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the GSLB.\n" +
								"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
							Optional: true,
						},
						common.ToSnakeCase("LinkedResourceCount"): schema.Int32Attribute{
							Description: "The number of resources linked to this GSLB.\n" +
								"  - example : 2",
							Optional: true,
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
							Optional: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current state of the GSLB (e.g., ACTIVE, CREATING, EDITING, ERROR, DELETING).\n" +
								"  - example : ACTIVE",
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *gslbGslbDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *gslbGslbDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state gslb.GslbDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetGslbList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Gslbs",
			err.Error(),
		)
		return
	}

	for _, gslbRes := range data.Gslbs {
		gslbState := gslb.Gslb{
			Algorithm:           types.StringValue(gslbRes.Algorithm),
			CreatedAt:           types.StringValue(gslbRes.CreatedAt.Format(time.RFC3339)),
			CreatedBy:           types.StringValue(gslbRes.CreatedBy),
			Description:         virtualserverutil.ToNullableStringValue(gslbRes.Description.Get()),
			EnvUsage:            types.StringValue(gslbRes.EnvUsage),
			Id:                  types.StringValue(gslbRes.Id),
			LinkedResourceCount: types.Int32Value(gslbRes.LinkedResourceCount),
			ModifiedAt:          types.StringValue(gslbRes.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:          types.StringValue(gslbRes.ModifiedBy),
			Name:                types.StringValue(gslbRes.Name),
			State:               types.StringValue(gslbRes.State),
		}

		state.Gslbs = append(state.Gslbs, gslbState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}
