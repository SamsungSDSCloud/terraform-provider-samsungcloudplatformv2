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
				Description: "The unique identifier of the LB Listener.",
				Optional:    true,
			},
			common.ToSnakeCase("LbListener"): schema.SingleNestedAttribute{
				Description: "Details of the LB Listener.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier.",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource (max 255 characters). This helps identify the purpose or usage of the resource.",
						Optional:    true,
					},
					common.ToSnakeCase("InsertClientIp"): schema.BoolAttribute{
						Description: "Whether to insert client IP in the header using Proxy Protocol v1.",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LB Listener (1-63 characters, alphanumeric with spaces, hyphens, underscores, and dots allowed).",
						Optional:    true,
					},
					common.ToSnakeCase("Persistence"): schema.StringAttribute{
						Description: "Session persistence configuration (e.g., 'source-ip', 'cookie').",
						Optional:    true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "The protocol used for the listener (TCP, UDP, HTTP, HTTPS, TLS, TCP_PROXY).",
						Optional:    true,
					},
					common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
						Description: "The ID of the server group associated with the listener. Required for TCP, UDP, and TLS protocols. This field is optional.",
						Optional:    true,
					},
					common.ToSnakeCase("ServerGroupName"): schema.StringAttribute{
						Description: "The server group name for the listener.",
						Optional:    true,
					},
					common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
						Description: "The service port number for the listener.",
						Optional:    true,
					},
					common.ToSnakeCase("ResponseTimeout"): schema.Int32Attribute{
						Description: "The response timeout in seconds (1-120). Only for L7 protocols (HTTP/HTTPS). Cannot be used with idle_timeout.",
						Optional:    true,
					},
					common.ToSnakeCase("SessionDurationTime"): schema.Int32Attribute{
						Description: "The session duration time in seconds. Required for L4 protocols (TCP/UDP). L7: 1-120, TCP/TLS: 60-3600 (60-second increments), UDP: 60-180 (60-second increments). Cannot be used with idle_timeout for L7.",
						Optional:    true,
					},
					common.ToSnakeCase("SslCertificate"): schema.SingleNestedAttribute{
						Description: "SSL certificate configuration for the listener.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ClientCertId"): schema.StringAttribute{
								Description: "The client certificate ID.",
								Optional:    true,
							},
							common.ToSnakeCase("ClientCertLevel"): schema.StringAttribute{
								Description: "The client certificate validation level (LOW, NORMAL, HIGH).",
								Optional:    true,
							},
							common.ToSnakeCase("ServerCertLevel"): schema.StringAttribute{
								Description: "The server certificate validation level (LOW, NORMAL, HIGH).",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("SniCertificate"): schema.ListNestedAttribute{
						Description: "SNI certificate configuration for multiple domains.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("SniCertId"): schema.StringAttribute{
									Description: "The SNI certificate ID.",
									Optional:    true,
								},
								common.ToSnakeCase("DomainName"): schema.StringAttribute{
									Description: "The domain name for SNI certificate (1-63 characters, alphanumeric with dots and hyphens, must start and end with alphanumeric). Must be unique within the listener.",
									Optional:    true,
								},
								common.ToSnakeCase("NotAfterDt"): schema.StringAttribute{
									Description: "The expiration date and time of the certificate (e.g., '2024-12-31T23:59:59Z'). Read-only.",
									Optional:    true,
								},
							},
						},
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the LB Listener (CREATING, ACTIVE, DELETING, ERROR).",
						Computed:    true,
					},
					common.ToSnakeCase("UrlHandler"): schema.ListNestedAttribute{
						Description: "URL handler configuration for routing (max 20 entries). Only for L7 protocols (HTTP/HTTPS). Requires at least one default rule (seq=0, url_pattern='default').",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("UrlPattern"): schema.StringAttribute{
									Description: "The URL pattern for routing (1-63 characters). For URL_PATH condition: alphanumeric with /_-., for HOST_HEADER condition: alphanumeric with .-. Example: '/api/v1' or 'example.com'.",
									Optional:    true,
								},
								common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
									Description: "The ID of the server group to route traffic to when the URL pattern matches. Required in url_handler.",
									Optional:    true,
								},
								common.ToSnakeCase("Seq"): schema.Int32Attribute{
									Description: "The sequence number for routing priority. 0 is reserved for default rule. Example: 1, 2, 3...",
									Optional:    true,
								},
							},
						},
					},
					common.ToSnakeCase("HttpsRedirection"): schema.SingleNestedAttribute{
						Description: "HTTPS redirection configuration. Only for HTTP protocol listeners.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "The protocol to redirect to. Must be 'HTTPS'.",
								Optional:    true,
							},
							common.ToSnakeCase("Port"): schema.StringAttribute{
								Description: "The port number to redirect to (1-65534).",
								Optional:    true,
							},
							common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
								Description: "The HTTP response code for redirection. Must be '301'.",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("UrlRedirection"): schema.StringAttribute{
						Description: "URL redirection configuration (max 8 entries).",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedFor"): schema.BoolAttribute{
						Description: "X-Forwarded-For header configuration.",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedPort"): schema.BoolAttribute{
						Description: "X-Forwarded-Port header configuration.",
						Optional:    true,
					},
					common.ToSnakeCase("XForwardedProto"): schema.BoolAttribute{
						Description: "X-Forwarded-Proto header configuration.",
						Optional:    true,
					},
					common.ToSnakeCase("RoutingAction"): schema.StringAttribute{
						Description: "The routing action type. 'LB_SERVER_GROUP' for URL handler routing, 'URL_REDIRECT' for HTTPS/URL redirection. Required for L4 protocols. Set only during creation.",
						Optional:    true,
					},
					common.ToSnakeCase("ConditionType"): schema.StringAttribute{
						Description: "The condition type for routing. 'URL_PATH' or 'HOST_HEADER' for URL handler, 'PROTOCOL_PORT' for HTTPS redirection. Only for L7 protocols. Cannot be modified when changing url_handler.",
						Optional:    true,
					},
					common.ToSnakeCase("IdleTimeout"): schema.Int32Attribute{
						Description: "The idle timeout in seconds (60-3600, in 60-second increments). Only applicable for L7 protocols (HTTP, HTTPS).",
						Optional:    true,
					},
					common.ToSnakeCase("HstsMaxAge"): schema.Int32Attribute{
						Description: "HSTS max age in seconds.",
						Optional:    true,
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

	lbListenerState := loadbalancerutil.ConvertResponse(data)

	lbListenerObjectValue, _ := types.ObjectValueFrom(ctx, lbListenerState.AttributeTypes(), lbListenerState)
	state.LbListenerDetail = lbListenerObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
