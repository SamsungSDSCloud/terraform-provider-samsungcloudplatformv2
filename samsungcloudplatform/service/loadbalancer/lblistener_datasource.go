package loadbalancer

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loadbalancer" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	loadbalancerutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/loadbalancer"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadbalancerLbListenerDataSource{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLbListenerDataSource{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewLoadbalancerLbListenerDataSource() datasource.DataSource {
	return &loadbalancerLbListenerDataSource{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type loadbalancerLbListenerDataSource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLbListenerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_listener" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *loadbalancerLbListenerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Retrieve details of a specific LB Listener.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the LB Listener.\n" +
					"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
				Optional: true,
			},
			common.ToSnakeCase("LbListener"): schema.SingleNestedAttribute{
				Description: "Details of the LB Listener.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier.\n" +
							"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z\n",
						Optional: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z\n",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
						Optional: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : LB Listener for web traffic\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
					common.ToSnakeCase("InsertClientIp"): schema.BoolAttribute{
						Description: "Whether to insert client IP in the header using Proxy Protocol v1.\n" +
							"  - example : true\n",
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
					common.ToSnakeCase("Persistence"): schema.StringAttribute{
						Description: "Session persistence configuration.\n" +
							"  - example : source-ip\n" +
							"  - pattern : source-ip | cookie\n",
						Optional: true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "The protocol used for the listener.\n" +
							"  - example : HTTP\n" +
							"  - pattern : TCP | UDP | HTTP | HTTPS | TLS | TCP_PROXY\n",
						Optional: true,
					},
					common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
						Description: "The ID of the server group associated with the listener.\n" +
							"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
						Optional: true,
					},
					common.ToSnakeCase("ServerGroupName"): schema.StringAttribute{
						Description: "The server group name for the listener.\n" +
							"  - example : ServerGroup01\n",
						Optional: true,
					},
					common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
						Description: "The service port number for the listener.\n" +
							"  - example : 80\n" +
							"  - minimum : 1\n" +
							"  - maximum : 65535\n",
						Optional: true,
					},
					common.ToSnakeCase("ResponseTimeout"): schema.Int32Attribute{
						Description: "The response timeout in seconds. Only for L7 protocols (HTTP/HTTPS).\n" +
							"  - example : 30\n" +
							"  - minimum : 1\n" +
							"  - maximum : 120\n",
						Optional: true,
					},
					common.ToSnakeCase("SessionDurationTime"): schema.Int32Attribute{
						Description: "The session duration time in seconds.\n" +
							"  - example : 3600\n" +
							"  - minimum : 1\n" +
							"  - maximum : 3600\n",
						Optional: true,
					},
					common.ToSnakeCase("SslCertificate"): schema.SingleNestedAttribute{
						Description: "SSL certificate configuration for the listener.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ClientCertId"): schema.StringAttribute{
								Description: "The client certificate ID.\n" +
									"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
								Optional: true,
							},
							common.ToSnakeCase("ClientCertLevel"): schema.StringAttribute{
								Description: "The client certificate validation level.\n" +
									"  - example : NORMAL\n" +
									"  - pattern : LOW | NORMAL | HIGH\n",
								Optional: true,
							},
							common.ToSnakeCase("ServerCertLevel"): schema.StringAttribute{
								Description: "The server certificate validation level.\n" +
									"  - example : NORMAL\n" +
									"  - pattern : LOW | NORMAL | HIGH\n",
								Optional: true,
							},
						},
					},
					common.ToSnakeCase("SniCertificate"): schema.ListNestedAttribute{
						Description: "SNI certificate configuration for multiple domains.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("SniCertId"): schema.StringAttribute{
									Description: "The SNI certificate ID.\n" +
										"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
									Optional: true,
								},
								common.ToSnakeCase("DomainName"): schema.StringAttribute{
									Description: "The domain name for SNI certificate.\n" +
										"  - example : example.com\n" +
										"  - minLength : 1\n" +
										"  - maxLength : 63\n",
									Optional: true,
								},
								common.ToSnakeCase("NotAfterDt"): schema.StringAttribute{
									Description: "The expiration date and time of the certificate.\n" +
										"  - example : 2024-12-31T23:59:59Z\n",
									Optional: true,
								},
							},
						},
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the LB Listener.\n" +
							"  - example : ACTIVE\n" +
							"  - pattern : CREATING | ACTIVE | DELETING | ERROR\n",
						Computed: true,
					},
					common.ToSnakeCase("UrlHandler"): schema.ListNestedAttribute{
						Description: "URL handler configuration for routing. Only for L7 protocols (HTTP/HTTPS).",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("UrlPattern"): schema.StringAttribute{
									Description: "The URL pattern for routing.\n" +
										"  - example : /api/v1\n" +
										"  - minLength : 1\n" +
										"  - maxLength : 63\n",
									Optional: true,
								},
								common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
									Description: "The ID of the server group to route traffic to.\n" +
										"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
									Optional: true,
								},
								common.ToSnakeCase("Seq"): schema.Int32Attribute{
									Description: "The sequence number for routing priority. 0 is reserved for default rule.\n" +
										"  - example : 1\n",
									Optional: true,
								},
							},
						},
					},
					common.ToSnakeCase("HttpsRedirection"): schema.SingleNestedAttribute{
						Description: "HTTPS redirection configuration. Only for HTTP protocol listeners.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "The protocol to redirect to.\n" +
									"  - example : HTTPS\n",
								Optional: true,
							},
							common.ToSnakeCase("Port"): schema.StringAttribute{
								Description: "The port number to redirect to.\n" +
									"  - example : 443\n" +
									"  - minimum : 1\n" +
									"  - maximum : 65534\n",
								Optional: true,
							},
							common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
								Description: "The HTTP response code for redirection.\n" +
									"  - example : 301\n",
								Optional: true,
							},
						},
					},
					common.ToSnakeCase("IdleTimeout"): schema.Int32Attribute{
						Description: "The idle timeout in seconds. Only for L7 protocols (HTTP/HTTPS).\n" +
							"  - example : 60\n" +
							"  - minimum : 1\n" +
							"  - maximum : 120\n",
						Optional: true,
					},
					common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
						Description: "The LoadBalancer ID associated with the listener.\n" +
							"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadbalancerLbListenerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *loadbalancerLbListenerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state loadbalancer.LbListenerDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLbListener(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show LbListener",
			err.Error(),
		)
		return
	}

	lbListenerState, _ := loadbalancerutil.ConvertResponse(data)

	lbListenerObjectValue, _ := types.ObjectValueFrom(ctx, lbListenerState.AttributeTypes(), lbListenerState)
	state.LbListenerDetail = lbListenerObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
