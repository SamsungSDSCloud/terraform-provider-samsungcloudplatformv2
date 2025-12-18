package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loadbalancer" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	loadbalancerutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/loadbalancer"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
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
		Description: "Show Lb Health Check.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("LbHealthCheck"): schema.SingleNestedAttribute{
				Description: "A detail of Lb Health Check.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
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
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId",
						Optional:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "SubnetId",
						Optional:    true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "Protocol",
						Optional:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckPort"): schema.Int32Attribute{
						Description: "HealthCheckPort",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
						Description: "HealthCheckInterval",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckTimeout"): schema.Int32Attribute{
						Description: "HealthCheckTimeout",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckCount"): schema.Int32Attribute{
						Description: "HealthCheckCount",
						Optional:    true,
					},
					common.ToSnakeCase("HttpMethod"): schema.StringAttribute{
						Description: "HttpMethod",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckUrl"): schema.StringAttribute{
						Description: "HealthCheckUrl",
						Optional:    true,
					},
					common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
						Description: "ResponseCode",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckType"): schema.StringAttribute{
						Description: "HealthCheckType",
						Optional:    true,
					},
					common.ToSnakeCase("RequestData"): schema.StringAttribute{
						Description: "RequestData",
						Optional:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
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
