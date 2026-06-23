package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loadbalancer" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadbalancerLbListenerDataSources{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLbListenerDataSources{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewLoadbalancerLbListenerDataSources() datasource.DataSource {
	return &loadbalancerLbListenerDataSources{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type loadbalancerLbListenerDataSources struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLbListenerDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_listeners" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *loadbalancerLbListenerDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "List all LB Listeners.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example : 20\n",
				Optional: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number.\n" +
					"  - example : 0\n",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sort order.\n" +
					"  - example : name:asc\n",
				Optional: true,
			},
			common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
				Description: "The LoadBalancer ID associated with the listener.\n" +
					"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current state of the LB Listener.\n" +
					"  - example : ACTIVE\n" +
					"  - pattern : CREATING | ACTIVE | DELETING | ERROR\n",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the LB Listener.\n" +
					"  - example : Listener01\n" +
					"  - minLength : 1\n" +
					"  - maxLength : 63\n" +
					"  - pattern : ^[a-zA-Z0-9._-]+$\n",
				Optional: true,
			},
			common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
				Description: "The service port number for the listener.\n" +
					"  - example : 80\n" +
					"  - minimum : 1\n" +
					"  - maximum : 65535\n",
				Optional: true,
			},
			common.ToSnakeCase("LbListeners"): schema.ListNestedAttribute{
				Description: "List of LB Listeners.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the LB Listener.\n" +
								"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
							Optional: true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the LB Listener.\n" +
								"  - example : Listener01\n" +
								"  - minLength : 1\n" +
								"  - maxLength : 63\n" +
								"  - pattern : ^[a-zA-Z0-9._-]+$\n",
							Optional: true,
						},
						common.ToSnakeCase("Protocol"): schema.StringAttribute{
							Description: "The protocol used for the listener.\n" +
								"  - example : HTTP\n" +
								"  - pattern : TCP | UDP | HTTP | HTTPS | TLS | TCP_PROXY\n",
							Optional: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current state of the LB Listener.\n" +
								"  - example : ACTIVE\n" +
								"  - pattern : CREATING | ACTIVE | DELETING | ERROR\n",
							Optional: true,
						},
						common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
							Description: "The service port number for the listener.\n" +
								"  - example : 80\n" +
								"  - minimum : 1\n" +
								"  - maximum : 65535\n",
							Optional: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z\n",
							Computed: true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user id that created the resource.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z\n",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user id that last modified the resource.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadbalancerLbListenerDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *loadbalancerLbListenerDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state loadbalancer.LbListenerDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLbListenerList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read LbListeners",
			err.Error(),
		)
		return
	}

	for _, lblistener := range data.Listeners {
		lblistenerState := loadbalancer.LbListener{
			Id:          types.StringValue(lblistener.Id),
			Name:        types.StringValue(lblistener.Name),
			Protocol:    types.StringValue(string(lblistener.Protocol)),
			State:       types.StringValue(*lblistener.State.Get()),
			ServicePort: types.Int32Value(lblistener.ServicePort),
			CreatedAt:   types.StringValue(lblistener.CreatedAt.Format(time.RFC3339)),
			CreatedBy:   types.StringValue(lblistener.CreatedBy),
			ModifiedAt:  types.StringValue(lblistener.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:  types.StringValue(lblistener.ModifiedBy),
		}

		state.LbListeners = append(state.LbListeners, lblistenerState)

		// Set state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}
