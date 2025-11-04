package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loadbalancer" // client 를 import 한다.
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
	_ datasource.DataSource              = &loadBalancerLbServerGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &loadBalancerLbServerGroupDataSource{}
)

// NewLoadBalancerResourceGroupDataSource is a helper function to simplify the provider implementation.
func NewLoadbalancerLbServerGroupDataSource() datasource.DataSource {
	return &loadBalancerLbServerGroupDataSource{}
}

// loadBalancerLbServerGroupDataSource is the data source implementation.
type loadBalancerLbServerGroupDataSource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadBalancerLbServerGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_server_group"
}

// Schema defines the schema for the data source.
func (d *loadBalancerLbServerGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Lb Server Group.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("LbServerGroup"): schema.SingleNestedAttribute{
				Description: "A detail of Lb Server Group.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{common.ToSnakeCase("AccountId"): schema.StringAttribute{
					Description: "AccountId",
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
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("LbMethod"): schema.StringAttribute{
						Description: "LbMethod",
						Optional:    true,
					},
					common.ToSnakeCase("LbName"): schema.StringAttribute{
						Description: "LbName",
						Optional:    true,
					},
					common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
						Description: "LoadbalancerId",
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
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "Protocol",
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
					common.ToSnakeCase("LbHealthCheckId"): schema.StringAttribute{
						Description: "LbHealthCheckId",
						Optional:    true,
					}},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadBalancerLbServerGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *loadBalancerLbServerGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state loadbalancer.LbServerGroupDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from routing rule
	data, err := d.client.GetLbServerGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading lb server group",
			"Could not read lb server group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var lbServerGroupState = loadbalancer.LbServerGroupDetail{
		Name:            types.StringValue(data.LbServerGroup.Name),
		Protocol:        types.StringValue(string(data.LbServerGroup.Protocol)),
		LoadbalancerId:  types.StringPointerValue(data.LbServerGroup.LoadbalancerId.Get()),
		LbName:          virtualserverutil.ToNullableStringValue(data.LbServerGroup.LbName.Get()),
		LbMethod:        types.StringValue(string(data.LbServerGroup.LbMethod)),
		LbHealthCheckId: virtualserverutil.ToNullableStringValue(data.LbServerGroup.LbHealthCheckId.Get()),
		State:           types.StringValue(data.LbServerGroup.State),
		VpcId:           types.StringValue(data.LbServerGroup.VpcId),
		SubnetId:        types.StringValue(data.LbServerGroup.SubnetId),
		AccountId:       types.StringValue(data.LbServerGroup.AccountId),
		Description:     virtualserverutil.ToNullableStringValue(data.LbServerGroup.Description.Get()),
		ModifiedBy:      types.StringValue(data.LbServerGroup.ModifiedBy),
		ModifiedAt:      types.StringValue(data.LbServerGroup.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(data.LbServerGroup.CreatedBy),
		CreatedAt:       types.StringValue(data.LbServerGroup.CreatedAt.Format(time.RFC3339)),
	}
	lbServerGroupObjectValue, _ := types.ObjectValueFrom(ctx, lbServerGroupState.AttributeTypes(), lbServerGroupState)
	state.LbServerGroupDetail = lbServerGroupObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
