package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadbalancerLbHealthCheckDataSources{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLbHealthCheckDataSources{}
)

// NewLoadBalancerLbHealthCheckDataSources is a helper function to simplify the provider implementation.
func NewLoadbalancerLbHealthCheckDataSources() datasource.DataSource {
	return &loadbalancerLbHealthCheckDataSources{}
}

// loadbalancerLbHealthCheckDataSources is the data source implementation.
type loadbalancerLbHealthCheckDataSources struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLbHealthCheckDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_health_checks"
}

// Schema defines the schema for the data source.
func (d *loadbalancerLbHealthCheckDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Get List of Lb Health Checks.",
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
			common.ToSnakeCase("Protocol"): schema.ListAttribute{
				Description: "Protocol",
				Optional:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "SubnetId",
				Optional:    true,
			},
			common.ToSnakeCase("LbHealthChecks"): schema.ListNestedAttribute{
				Description: "A list of Lb Health Checks.",
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
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State",
							Optional:    true,
						},
						common.ToSnakeCase("LbServerGroupCount"): schema.Int32Attribute{
							Description: "LbServerGroupCount",
							Optional:    true,
						},
						common.ToSnakeCase("HealthCheckType"): schema.StringAttribute{
							Description: "HealthCheckType",
							Optional:    true,
						},
						common.ToSnakeCase("Protocol"): schema.StringAttribute{
							Description: "Protocol",
							Optional:    true,
						},
						common.ToSnakeCase("SubnetId"): schema.StringAttribute{
							Description: "SubnetId",
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
func (d *loadbalancerLbHealthCheckDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *loadbalancerLbHealthCheckDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state loadbalancer.LbHealthCheckDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLbHealthCheckList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read LbHealthChecks",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, lbHealthCheck := range data.LbHealthChecks {
		lbHealthCheckState := loadbalancer.LbHealthCheck{
			Id:                 types.StringValue(lbHealthCheck.Id),
			Name:               types.StringValue(lbHealthCheck.Name),
			Protocol:           types.StringValue(string(lbHealthCheck.Protocol)),
			State:              types.StringValue(lbHealthCheck.State),
			SubnetId:           types.StringValue(lbHealthCheck.SubnetId),
			LbServerGroupCount: types.Int32Value(lbHealthCheck.LbServerGroupCount),
			CreatedAt:          types.StringValue(lbHealthCheck.CreatedAt.Format(time.RFC3339)),
			CreatedBy:          types.StringValue(lbHealthCheck.CreatedBy),
			ModifiedAt:         types.StringValue(lbHealthCheck.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:         types.StringValue(lbHealthCheck.ModifiedBy),
			HealthCheckType:    types.StringValue(string(lbHealthCheck.HealthCheckType)),
		}

		state.LbHealthChecks = append(state.LbHealthChecks, lbHealthCheckState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
