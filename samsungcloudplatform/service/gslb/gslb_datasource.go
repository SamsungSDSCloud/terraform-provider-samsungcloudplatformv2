package gslb

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/gslb" // client 를 import 한다.
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
		Description: "Show Gslb.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("GslbDetail"): schema.SingleNestedAttribute{
				Description: "A detail of Gslb.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Algorithm"): schema.StringAttribute{
						Description: "Algorithm",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("EnvUsage"): schema.StringAttribute{
						Description: "EnvUsage",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheck"): schema.SingleNestedAttribute{
						Description: "HealthCheck",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
								Description: "created at",
								Optional:    true,
							},
							common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
								Description: "created by",
								Optional:    true,
							},
							common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
								Description: "HealthCheckInterval",
								Optional:    true,
							},
							common.ToSnakeCase("HealthCheckProbeTimeout"): schema.Int32Attribute{
								Description: "HealthCheckProbeTimeout",
								Optional:    true,
							},
							common.ToSnakeCase("HealthCheckUserId"): schema.StringAttribute{
								Description: "HealthCheckUserId",
								Optional:    true,
							},
							common.ToSnakeCase("HealthCheckUserPassword"): schema.StringAttribute{
								Description: "HealthCheckUserPassword",
								Optional:    true,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "id",
								Computed:    true,
							},
							common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
								Description: "modified at",
								Optional:    true,
							},
							common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
								Description: "modified by",
								Optional:    true,
							},
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "Protocol",
								Optional:    true,
							},
							common.ToSnakeCase("ReceiveString"): schema.StringAttribute{
								Description: "ReceiveString",
								Optional:    true,
							},
							common.ToSnakeCase("SendString"): schema.StringAttribute{
								Description: "SendString",
								Optional:    true,
							},
							common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
								Description: "ServicePort",
								Optional:    true,
							},
							common.ToSnakeCase("Timeout"): schema.Int32Attribute{
								Description: "Timeout",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "id",
						Computed:    true,
					},
					common.ToSnakeCase("LinkedResourceCount"): schema.Int32Attribute{
						Description: "LinkedResourceCount",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
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
