package gslb

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/gslb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
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
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size",
				Optional:    true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("GslbId"): schema.StringAttribute{
				Description: "GslbId",
				Optional:    true,
			},
			common.ToSnakeCase("GslbResources"): schema.ListNestedAttribute{
				Description: "A list of gslb resource.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "created at",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "created by",
							Computed:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
							Optional:    true,
						},
						common.ToSnakeCase("Destination"): schema.StringAttribute{
							Description: "Destination",
							Optional:    true,
						},
						common.ToSnakeCase("Disabled"): schema.BoolAttribute{
							Description: "Disabled",
							Optional:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Optional:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "modified at",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "modified by",
							Computed:    true,
						},
						common.ToSnakeCase("Region"): schema.StringAttribute{
							Description: "Region",
							Optional:    true,
						},
						common.ToSnakeCase("Weight"): schema.Int32Attribute{
							Description: "Weight",
							Optional:    true,
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

	for _, gslbResouece := range data.GslbResources {
		gslbResourceState := gslb.GslbResourceDetail{
			CreatedAt:   types.StringValue(gslbResouece.CreatedAt.Format(time.RFC3339)),
			CreatedBy:   types.StringValue(gslbResouece.CreatedBy),
			Description: virtualserverutil.ToNullableStringValue(gslbResouece.Description.Get()),
			Destination: types.StringValue(gslbResouece.Destination),
			Disabled:    common.ToNullableBoolValue(gslbResouece.Disabled.Get()),
			Id:          types.StringValue(gslbResouece.Id),
			ModifiedAt:  types.StringValue(gslbResouece.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:  types.StringValue(gslbResouece.ModifiedBy),
			Region:      types.StringValue(gslbResouece.Region),
			Weight:      common.ToNullableInt32Value(gslbResouece.Weight.Get()),
		}

		state.GslbResources = append(state.GslbResources, gslbResourceState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}
