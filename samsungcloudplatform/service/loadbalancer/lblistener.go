package loadbalancer

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	loadbalancerutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/loadbalancer"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &loadbalancerLbListenerResource{}
	_ resource.ResourceWithConfigure = &loadbalancerLbListenerResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewLoadBalancerListenerResource() resource.Resource {
	return &loadbalancerLbListenerResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type loadbalancerLbListenerResource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLbListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_listener" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *loadbalancerLbListenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "LB Listener resource for managing listener configurations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.\n" +
					"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
							"  - example : listener-web-01\n" +
							"  - minLength : 3\n" +
							"  - maxLength : 63\n" +
							"  - pattern : ^[a-zA-Z0-9-_]*$\n",
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
							"  - example : TCP\n" +
							"  - pattern : TCP | UDP | HTTP | HTTPS | TLS | TCP_PROXY\n",
						Optional: true,
					},
					common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
						Description: "The ID of the server group associated with the listener. Required for TCP, UDP, and TLS protocols.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
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
							"  - maximum : 65534\n",
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
						Description: "The session duration time in seconds. Required for L4 protocols (TCP/UDP).\n" +
							"  - example : 60\n" +
							"  - L7(HTTP/HTTPS) : minimum 1, maximum 120\n" +
							"  - L4(TCP/TLS) : minimum 60, maximum 3600 (60-second increments)\n" +
							"  - UDP : minimum 60, maximum 180 (60-second increments)\n",
						Optional: true,
					},
					common.ToSnakeCase("SslCertificate"): schema.SingleNestedAttribute{
						Description: "SSL certificate configuration for the listener.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ClientCertId"): schema.StringAttribute{
								Description: "The client certificate ID.\n" +
									"  - example : 46c681018e33453085ca7c8db54e0076\n",
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
										"  - example : 46c681018e33453085ca7c8db54e0076\n",
									Optional: true,
								},
								common.ToSnakeCase("DomainName"): schema.StringAttribute{
									Description: "The domain name for SNI certificate. Must be unique within the listener.\n" +
										"  - example : example.com\n" +
										"  - minLength : 1\n" +
										"  - maxLength : 63\n" +
										"  - pattern : ^[a-zA-Z0-9](?:[a-zA-Z0-9.-]{0,61}[a-zA-Z0-9])?$\n",
									Optional: true,
								},
								common.ToSnakeCase("NotAfterDt"): schema.StringAttribute{
									Description: "The expiration date and time of the certificate. Read-only.\n" +
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
									Description: "The ID of the server group to route traffic to when the URL pattern matches.\n" +
										"  - example : 46c681018e33453085ca7c8db54e0076\n",
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
					common.ToSnakeCase("UrlRedirection"): schema.StringAttribute{
						Description: "URL redirection configuration. Must start with http:// or https://.\n" +
							"  - example : https://example.com\n" +
							"  - maxLength : 255\n" +
							"  - pattern (HTTP) : ^http://[A-Za-z0-9:/.\\-?=&#]{0,248}$\n" +
							"  - pattern (HTTPS) : ^https://[A-Za-z0-9:/.\\-?=&#]{0,247}$\n",
						Optional: true,
					},
					common.ToSnakeCase("XForwardedFor"): schema.BoolAttribute{
						Description: "X-Forwarded-For header configuration.\n" +
							"  - example : true\n",
						Optional: true,
					},
					common.ToSnakeCase("XForwardedPort"): schema.BoolAttribute{
						Description: "X-Forwarded-Port header configuration.\n" +
							"  - example : true\n",
						Optional: true,
					},
					common.ToSnakeCase("XForwardedProto"): schema.BoolAttribute{
						Description: "X-Forwarded-Proto header configuration.\n" +
							"  - example : true\n",
						Optional: true,
					},
					common.ToSnakeCase("RoutingAction"): schema.StringAttribute{
						Description: "The routing action type. 'LB_SERVER_GROUP' for URL handler routing, 'URL_REDIRECT' for HTTPS/URL redirection.\n" +
							"  - example : LB_SERVER_GROUP\n",
						Optional: true,
					},
					common.ToSnakeCase("ConditionType"): schema.StringAttribute{
						Description: "The condition type for routing. 'URL_PATH' or 'HOST_HEADER' for URL handler.\n" +
							"  - example : URL_PATH\n",
						Optional: true,
					},
					common.ToSnakeCase("IdleTimeout"): schema.Int32Attribute{
						Description: "The idle timeout in seconds. Only applicable for L7 protocols (HTTP, HTTPS).\n" +
							"  - example : 60\n" +
							"  - minimum : 60\n" +
							"  - maximum : 3600\n",
						Optional: true,
					},
					common.ToSnakeCase("HstsMaxAge"): schema.Int32Attribute{
						Description: "HSTS max age in seconds.\n" +
							"  - example : 31536000\n",
						Optional: true,
					},
				},
			},
			common.ToSnakeCase("LbListenerCreate"): schema.SingleNestedAttribute{
				Description: "Parameters for creating a new LB Listener.",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
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
					common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
						Description: "The LoadBalancer ID associated with the listener.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LB Listener.\n" +
							"  - example : Listener01\n" +
							"  - minLength : 1\n" +
							"  - maxLength : 63\n",
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
						Validators: []validator.String{
							stringvalidator.OneOf("HTTP", "HTTPS", "TCP", "UDP", "TLS", "TCP_PROXY"),
						},
					},
					common.ToSnakeCase("ResponseTimeout"): schema.Int32Attribute{
						Description: "The response timeout in seconds. Only for L7 protocols (HTTP/HTTPS).\n" +
							"  - example : 30\n" +
							"  - minimum : 1\n" +
							"  - maximum : 120\n",
						Optional: true,
					},
					common.ToSnakeCase("ServerGroupId"): schema.StringAttribute{
						Description: "The ID of the server group associated with the listener. Required for TCP, UDP, and TLS protocols.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
						Description: "The service port number for the listener.\n" +
							"  - example : 80\n" +
							"  - minimum : 1\n" +
							"  - maximum : 65534\n",
						Optional: true,
					},
					common.ToSnakeCase("SessionDurationTime"): schema.Int32Attribute{
						Description: "The session duration time in seconds. Required for L4 protocols (TCP/UDP).\n" +
							"  - example : 60\n" +
							"  - L7(HTTP/HTTPS) : minimum 1, maximum 120\n" +
							"  - L4(TCP/TLS) : minimum 60, maximum 3600 (60-second increments)\n" +
							"  - UDP : minimum 60, maximum 180 (60-second increments)\n",
						Optional: true,
					},

					common.ToSnakeCase("SslCertificate"): schema.SingleNestedAttribute{
						Description: "SSL certificate configuration for the listener.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ClientCertId"): schema.StringAttribute{
								Description: "The client certificate ID.\n" +
									"  - example : 46c681018e33453085ca7c8db54e0076\n",
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
										"  - example : 46c681018e33453085ca7c8db54e0076\n",
									Optional: true,
								},
								common.ToSnakeCase("DomainName"): schema.StringAttribute{
									Description: "The domain name for SNI certificate. Must be unique within the listener.\n" +
										"  - example : example.com\n" +
										"  - minLength : 1\n" +
										"  - maxLength : 63\n" +
										"  - pattern : ^[a-zA-Z0-9](?:[a-zA-Z0-9.-]{0,61}[a-zA-Z0-9])?$\n",
									Optional: true,
								},
							},
						},
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
									Description: "The ID of the server group to route traffic to when the URL pattern matches.\n" +
										"  - example : 46c681018e33453085ca7c8db54e0076\n",
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
					common.ToSnakeCase("UrlRedirection"): schema.StringAttribute{
						Description: "URL redirection configuration. Must start with http:// or https://.\n" +
							"  - example : https://example.com\n" +
							"  - maxLength : 255\n" +
							"  - pattern (HTTP) : ^http://[A-Za-z0-9:/.\\-?=&#]{0,248}$\n" +
							"  - pattern (HTTPS) : ^https://[A-Za-z0-9:/.\\-?=&#]{0,247}$\n",
						Optional: true,
					},
					common.ToSnakeCase("XForwardedFor"): schema.BoolAttribute{
						Description: "X-Forwarded-For header configuration.\n" +
							"  - example : true\n",
						Optional: true,
					},
					common.ToSnakeCase("XForwardedPort"): schema.BoolAttribute{
						Description: "X-Forwarded-Port header configuration.\n" +
							"  - example : true\n",
						Optional: true,
					},
					common.ToSnakeCase("XForwardedProto"): schema.BoolAttribute{
						Description: "X-Forwarded-Proto header configuration.\n" +
							"  - example : true\n",
						Optional: true,
					},
					common.ToSnakeCase("RoutingAction"): schema.StringAttribute{
						Description: "The routing action type. 'LB_SERVER_GROUP' for URL handler routing, 'URL_REDIRECT' for HTTPS/URL redirection.\n" +
							"  - example : LB_SERVER_GROUP\n",
						Optional: true,
					},
					common.ToSnakeCase("ConditionType"): schema.StringAttribute{
						Description: "The condition type for routing. 'URL_PATH' or 'HOST_HEADER' for URL handler.\n" +
							"  - example : URL_PATH\n",
						Optional: true,
					},
					common.ToSnakeCase("IdleTimeout"): schema.Int32Attribute{
						Description: "The idle timeout in seconds. Only applicable for L7 protocols (HTTP, HTTPS).\n" +
							"  - example : 60\n" +
							"  - minimum : 60\n" +
							"  - maximum : 3600\n",
						Optional: true,
					},
					common.ToSnakeCase("HstsMaxAge"): schema.Int32Attribute{
						Description: "HSTS max age in seconds.\n" +
							"  - example : 31536000\n",
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loadbalancerLbListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.LoadBalancer
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *loadbalancerLbListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan loadbalancer.LbListenerResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateLbListenerCreate(ctx, r.client, plan.LbListenerCreate); err != nil {
		resp.Diagnostics.AddError("Error creating LB Listener", err.Error())
		return
	}

	// Create new Lb Listener
	data, err := r.client.CreateLbListener(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Listener",
			"Could not create Lb Listener, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.Listener.Id)

	// Map response body to schema and populate Computed attribute values
	lbListenerModel := loadbalancerutil.ConvertResponse(data)
	lbListenerOjbectValue, diags := types.ObjectValueFrom(ctx, lbListenerModel.AttributeTypes(), lbListenerModel)
	plan.LbListener = lbListenerOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	// Wait for active state
	err = waitForLBListenerStatus(ctx, r.client, data.Listener.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Lb Listener",
			"Error waiting for Lb Listener to become active: "+err.Error(),
		)
		return
	}

	// use read function to fetch latest data
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)

	resp.State = readResp.State
}

// Read refreshes the Terraform state with the latest data.
func (r *loadbalancerLbListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state loadbalancer.LbListenerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from LB Listener
	data, err := r.client.GetLbListener(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Lb Listener",
			"Could not read Lb Listener, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	lbListenerModel := loadbalancerutil.ConvertResponse(data)

	lbListenerObjectValue, diags := types.ObjectValueFrom(ctx, lbListenerModel.AttributeTypes(), lbListenerModel)
	state.LbListener = lbListenerObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loadbalancerLbListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state loadbalancer.LbListenerResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateLbListenerCreate(ctx, r.client, state.LbListenerCreate); err != nil {
		resp.Diagnostics.AddError("Error updating LB Listener", err.Error())
		return
	}

	// Update existing order
	_, err := r.client.UpdateLbListener(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Listener",
			"Could not create Lb Listener, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetFirewallRule as UpdateFirewallRule items are not populated.
	data, err := r.client.GetLbListener(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Listener",
			"Could not update Lb Listener, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	lbListenerModel := loadbalancerutil.ConvertResponse(data)

	lbListenerObjectValue, diags := types.ObjectValueFrom(ctx, lbListenerModel.AttributeTypes(), lbListenerModel)
	state.LbListener = lbListenerObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *loadbalancerLbListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancer.LbListenerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing LB Server Group
	err := r.client.DeleteLbListener(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting LB Listener",
			"Could not delete lb listener, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func validateServerGroup(ctx context.Context, lbClient *loadbalancer.Client, serverGroupId string) error {
	serverGroup, err := lbClient.GetLbServerGroup(ctx, serverGroupId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		return fmt.Errorf("LBServer Group ID is not valid, unexpected error: %s\nReason: %s", err.Error(), detail)
	}

	if serverGroup.LbServerGroup.State != "ACTIVE" {
		return fmt.Errorf("LB Server Group with ID %s is not in ACTIVE state", serverGroupId)
	}
	return nil
}

func validateLbListenerCreate(ctx context.Context, lbClient *loadbalancer.Client, create loadbalancer.LbListenerCreate) error {
	protocol := strings.ToUpper(create.Protocol.ValueString())
	layer4Protocols := map[string]bool{"TCP": true, "UDP": true, "TLS": true, "TCP_PROXY": true}
	if layer4Protocols[protocol] && len(create.UrlHandler) > 0 {
		return fmt.Errorf("URL Handler is not supported for %s protocol", protocol)
	}

	if !create.ServerGroupId.IsNull() {
		if err := validateServerGroup(ctx, lbClient, create.ServerGroupId.ValueString()); err != nil {
			return err
		}
	}

	for _, v := range create.UrlHandler {
		if !v.Seq.IsNull() && v.Seq.ValueInt32() == 0 {
			if v.UrlPattern.ValueString() != "default" {
				return fmt.Errorf("URL Handler with seq 0 must have url_pattern set to 'default'")
			}
		}

		if !v.ServerGroupId.IsNull() {
			if err := validateServerGroup(ctx, lbClient, v.ServerGroupId.ValueString()); err != nil {
				return err
			}
		}
	}
	return nil
}

func waitForLBListenerStatus(ctx context.Context, loadbalancerClient *loadbalancer.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := loadbalancerClient.GetLbListener(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.Listener.State), nil
	}, -1, -1, -1, -1)
}
