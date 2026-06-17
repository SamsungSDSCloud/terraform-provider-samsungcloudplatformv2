package gslb

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/gslb" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &gslbGslbDataSource{}
	_ datasource.DataSourceWithConfigure = &gslbGslbDataSource{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewGslbGslbDataSource() datasource.DataSource {
	return &gslbGslbDataSource{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type gslbGslbDataSource struct {
	config  *scpsdk.Configuration
	client  *gslb.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *gslbGslbDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gslb_gslb" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *gslbGslbDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Shows details of a specific Global Server Load Balancer.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the GSLB to query.\n" +
					"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
				Optional: true,
			},
			common.ToSnakeCase("GslbDetail"): schema.SingleNestedAttribute{
				Description: "Details of the Global Server Load Balancer.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Algorithm"): schema.StringAttribute{
						Description: "The load balancing algorithm for GSLB traffic distribution (e.g., ROUND_ROBIN, RATIO).\n" +
							"  - example : ROUND_ROBIN",
						Optional: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z",
						Optional: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Optional: true,
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
					common.ToSnakeCase("HealthCheck"): schema.SingleNestedAttribute{
						Description: "Health check configuration for monitoring GSLB endpoint availability.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
								Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
									"  - example : 2024-05-17T00:23:17Z",
								Optional: true,
							},
							common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
								Description: "The user id that created the resource.\n" +
									"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
								Optional: true,
							},
							common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
								Description: "The GSLB Health Check Interval.\n" +
									"  - example : 30\n" +
									"  - Range: 5 to 299",
								Optional: true,
							},
							common.ToSnakeCase("HealthCheckProbeTimeout"): schema.Int32Attribute{
								Description: "The GSLB Health Check Probe Timeout.\n" +
									"  - example : 10\n" +
									"  - Range: 5 to 300",
								Optional: true,
							},
							common.ToSnakeCase("HealthCheckUserId"): schema.StringAttribute{
								Description: "The GSLB Health Check User Name.\n" +
									"  - example : healthcheck_user\n" +
									"  - Max length: 60",
								Optional: true,
							},
							common.ToSnakeCase("HealthCheckUserPassword"): schema.StringAttribute{
								Description: "The GSLB Health Check Password.\n" +
                              		"  - example : **********",
								Optional: true,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "The unique identifier of the health check configuration.\n" +
									"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
								Computed: true,
							},
							common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
								Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
									"  - example : 2024-05-17T00:23:17Z",
								Optional: true,
							},
							common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
								Description: "The user id that last modified the resource.\n" +
									"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
								Optional: true,
							},
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "The protocol used for health checks (e.g., ICMP, TCP, HTTP, HTTPS).\n" +
									"  - example : TCP",
								Optional: true,
							},
							common.ToSnakeCase("ReceiveString"): schema.StringAttribute{
								Description: "The GSLB Health Check Receive String.\n" +
									"  - example : HTTP/1.1 200 OK\n" +
									"  - Max length: 300",
								Optional: true,
							},
							common.ToSnakeCase("SendString"): schema.StringAttribute{
								Description: "The GSLB Health Check Send String. If no input is provided, it operates as a \"GET /\" request.\n" +
									"  - example : GET /",
								Optional: true,
							},
							common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
								Description: "The GSLB Health Check Service Port.\n" +
									"  - example : 80\n" +
									"  - Range: 1 to 65535",
								Optional: true,
							},
							common.ToSnakeCase("Timeout"): schema.Int32Attribute{
								Description: "The GSLB Health Check Timeout. It must be greater than the Interval.\n" +
									"  - example : 40\n" +
									"  - Range: 6 to 300",
								Optional: true,
							},
						},
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the GSLB.\n" +
							"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						Computed: true,
					},
					common.ToSnakeCase("LinkedResourceCount"): schema.Int32Attribute{
						Description: "The number of resources linked to this GSLB.\n" +
							"  - example : 2",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the GSLB.\n" +
							"  - example : example.gslb.e.samsungsdscloud.com",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the GSLB (e.g., ACTIVE, CREATING, EDITING, ERROR, DELETING).\n" +
							"  - example : ACTIVE",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *gslbGslbDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *gslbGslbDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state gslb.GslbDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetGslb(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show Gslb",
			err.Error(),
		)
		return
	}

	var healthCheck *gslb.HealthCheck
	var healthCheckFromData = data.Gslb.HealthCheck.Get()
	if healthCheckFromData != nil {
		healthCheck = &gslb.HealthCheck{
			CreatedAt:               types.StringValue(healthCheckFromData.CreatedAt.Format(time.RFC3339)),
			CreatedBy:               types.StringValue(healthCheckFromData.CreatedBy),
			HealthCheckInterval:     types.Int32Value(healthCheckFromData.GetHealthCheckInterval()),
			HealthCheckProbeTimeout: types.Int32Value(healthCheckFromData.GetHealthCheckProbeTimeout()),
			HealthCheckUserId:       types.StringValue(healthCheckFromData.GetHealthCheckUserId()),
			HealthCheckUserPassword: types.StringValue(healthCheckFromData.GetHealthCheckUserPassword()),
			Id:                      types.StringValue(healthCheckFromData.Id),
			ModifiedAt:              types.StringValue(healthCheckFromData.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:              types.StringValue(healthCheckFromData.ModifiedBy),
			Protocol:                types.StringValue(healthCheckFromData.Protocol),
			ReceiveString:           types.StringValue(healthCheckFromData.GetReceiveString()),
			SendString:              types.StringValue(healthCheckFromData.GetSendString()),
			ServicePort:             types.Int32Value(healthCheckFromData.GetServicePort()),
			Timeout:                 types.Int32Value(healthCheckFromData.GetTimeout()),
		}
	}
	gslbState := gslb.GslbDetail{
		Algorithm:           types.StringValue(data.Gslb.Algorithm),
		CreatedAt:           types.StringValue(data.Gslb.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(data.Gslb.CreatedBy),
		Description:         virtualserverutil.ToNullableStringValue(data.Gslb.Description.Get()),
		EnvUsage:            types.StringValue(data.Gslb.EnvUsage),
		HealthCheck:         healthCheck,
		Id:                  types.StringValue(data.Gslb.Id),
		LinkedResourceCount: types.Int32Value(data.Gslb.LinkedResourceCount),
		ModifiedAt:          types.StringValue(data.Gslb.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(data.Gslb.ModifiedBy),
		Name:                types.StringValue(data.Gslb.Name),
		State:               types.StringValue(data.Gslb.State),
	}

	gslbObjectValue, _ := types.ObjectValueFrom(ctx, gslbState.AttributeTypes(), gslbState)
	state.GslbDetail = gslbObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
