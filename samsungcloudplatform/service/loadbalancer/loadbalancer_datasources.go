package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadbalancerLoadbalancerDataSources{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLoadbalancerDataSources{}
)

// NewLoadbalancerLoadbalancerDataSources is a helper function to simplify the provider implementation.
func NewLoadbalancerLoadbalancerDataSources() datasource.DataSource {
	return &loadbalancerLoadbalancerDataSources{}
}

// loadbalancerLoadbalancerDataSources is the data source implementation.
type loadbalancerLoadbalancerDataSources struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLoadbalancerDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_loadbalancers"
}

// Schema defines the schema for the data source.
func (d *loadbalancerLoadbalancerDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Get List of Loadbalancers.",
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
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("ServiceIp"): schema.StringAttribute{
				Description: "ServiceIp",
				Optional:    true,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "SubnetId",
				Optional:    true,
			},
			common.ToSnakeCase("Loadbalancers"): schema.ListNestedAttribute{
				Description: "A list of Loadbalancers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Optional:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Name",
							Optional:    true,
						},
						common.ToSnakeCase("PublicNatIp"): schema.StringAttribute{
							Description: "PublicNatIp",
							Optional:    true,
						},
						common.ToSnakeCase("ServiceIp"): schema.StringAttribute{
							Description: "ServiceIp",
							Optional:    true,
						},
						common.ToSnakeCase("SourceNatIp"): schema.StringAttribute{
							Description: "SourceNatIp",
							Optional:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State",
							Optional:    true,
						},
						common.ToSnakeCase("ListenerCount"): schema.Int32Attribute{
							Description: "ListenerCount",
							Optional:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "created at",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "created by",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "modified at",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "modified by",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadbalancerLoadbalancerDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.LoadBalancer
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *loadbalancerLoadbalancerDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state loadbalancer.LoadbalancerDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLoadbalancerList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Loadbalancers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, loadbalancerElement := range data.Loadbalancers {
		loadbalancerState := loadbalancer.Loadbalancer{
			Id:            types.StringValue(loadbalancerElement.Id),
			Name:          virtualserverutil.ToNullableStringValue(loadbalancerElement.Name.Get()),
			PublicNatIp:   virtualserverutil.ToNullableStringValue(loadbalancerElement.PublicNatIp.Get()),
			ServiceIp:     virtualserverutil.ToNullableStringValue(loadbalancerElement.ServiceIp.Get()),
			SourceNatIp:   virtualserverutil.ToNullableStringValue(loadbalancerElement.SourceNatIp.Get()),
			State:         virtualserverutil.ToNullableStringValue(loadbalancerElement.State.Get()),
			ListenerCount: ToNullableIntValue(loadbalancerElement.ListenerCount.Get()),
			CreatedAt:     types.StringValue(loadbalancerElement.CreatedAt.Format(time.RFC3339)),
			CreatedBy:     types.StringValue(loadbalancerElement.CreatedBy),
			ModifiedAt:    types.StringValue(loadbalancerElement.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:    types.StringValue(loadbalancerElement.ModifiedBy),
		}

		state.Loadbalancers = append(state.Loadbalancers, loadbalancerState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ToNullableIntValue(v *int32) types.Int32 {
	if v == nil {
		return types.Int32Value(0)
	}
	return types.Int32Value(*v)
}
