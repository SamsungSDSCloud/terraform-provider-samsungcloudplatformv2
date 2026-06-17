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
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadBalancerLbHealthCheckDataSource{}
	_ datasource.DataSourceWithConfigure = &loadBalancerLbHealthCheckDataSource{}
)

// NewLoadBalancerResourceGroupDataSource is a helper function to simplify the provider implementation.
func NewLoadbalancerLbHealthCheckDataSource() datasource.DataSource {
	return &loadBalancerLbHealthCheckDataSource{}
}

// loadBalancerLbHealthCheckDataSource is the data source implementation.
type loadBalancerLbHealthCheckDataSource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadBalancerLbHealthCheckDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_health_check"
}

// Schema defines the schema for the data source.
func (d *loadBalancerLbHealthCheckDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details of a specific LB Health Check.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the LB Health Check.",
				Optional:    true,
			},
			common.ToSnakeCase("LbHealthCheck"): schema.SingleNestedAttribute{
				Description: "Details of the LB Health Check.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource (max 255 characters). This helps identify the purpose or usage of the resource.",
						Optional:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The VPC ID where the resource is located.",
						Optional:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The subnet ID where the resource is located.",
						Optional:    true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "The protocol used for the listener (e.g., TCP, HTTP, HTTPS).",
						Optional:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the Health Check (CREATING, ACTIVE, DELETING, ERROR).",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LB Health Check (1-63 characters, alphanumeric with spaces, hyphens, underscores, and dots allowed).",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckPort"): schema.Int32Attribute{
						Description: "The port number used for health checks (1-65534).",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
						Description: "The interval between health checks in seconds (1-180).",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckTimeout"): schema.Int32Attribute{
						Description: "The timeout for health check responses in seconds (1-180). Must be less than or equal to the interval.",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckCount"): schema.Int32Attribute{
						Description: "The number of consecutive health check failures before marking as unhealthy (1-10).",
						Optional:    true,
					},
					common.ToSnakeCase("HttpMethod"): schema.StringAttribute{
						Description: "The HTTP method used for health checks (GET, POST).",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckUrl"): schema.StringAttribute{
						Description: "The URL path for HTTP health checks (1-50 characters, must start with '/').",
						Optional:    true,
					},
					common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
						Description: "The expected HTTP response code for health checks (200-599).",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckType"): schema.StringAttribute{
						Description: "The type of health check (DEFAULT, CUSTOM).",
						Optional:    true,
					},
					common.ToSnakeCase("RequestData"): schema.StringAttribute{
						Description: "The request data sent during health checks (max 255 characters).",
						Optional:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadBalancerLbHealthCheckDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *loadBalancerLbHealthCheckDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state loadbalancer.LbHealthCheckDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from routing rule
	data, err := d.client.GetLbHealthCheck(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading lb health check",
			"Could not read lb health check ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var lbHealthCheckState = loadbalancer.LbHealthCheckDetail{
		Name:                types.StringValue(data.LbHealthCheck.Name),
		VpcId:               loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.VpcId.Get()),
		SubnetId:            loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.SubnetId.Get()),
		Protocol:            loadbalancerutil.ToNullableStringValue((*string)(data.LbHealthCheck.Protocol)),
		HealthCheckPort:     loadbalancerutil.ToNullableInt32Value(data.LbHealthCheck.HealthCheckPort.Get()),
		HealthCheckInterval: types.Int32Value(*data.LbHealthCheck.HealthCheckInterval),
		HealthCheckTimeout:  types.Int32Value(*data.LbHealthCheck.HealthCheckTimeout),
		HealthCheckCount:    types.Int32Value(*data.LbHealthCheck.HealthCheckCount),
		HttpMethod:          loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.HttpMethod.Get()),
		HealthCheckUrl:      loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.HealthCheckUrl.Get()),
		ResponseCode:        loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.ResponseCode.Get()),
		RequestData:         loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.RequestData.Get()),
		HealthCheckType:     types.StringValue(string(data.LbHealthCheck.HealthCheckType)),
		Description:         loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.Description.Get()),
		State:               types.StringValue(data.LbHealthCheck.State),
		AccountId:           loadbalancerutil.ToNullableStringValue(data.LbHealthCheck.AccountId.Get()),
		ModifiedBy:          types.StringValue(data.LbHealthCheck.ModifiedBy),
		ModifiedAt:          types.StringValue(data.LbHealthCheck.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(data.LbHealthCheck.CreatedBy),
		CreatedAt:           types.StringValue(data.LbHealthCheck.CreatedAt.Format(time.RFC3339)),
	}
	lbHealthCheckObjectValue, _ := types.ObjectValueFrom(ctx, lbHealthCheckState.AttributeTypes(), lbHealthCheckState)
	state.LbHealthCheckDetail = lbHealthCheckObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
