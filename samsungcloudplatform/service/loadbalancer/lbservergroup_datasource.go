package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loadbalancer" // client 를 import 한다.
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
		Description: "Retrieve details of a specific LB Server Group.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the LB Server Group.\n" +
					"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
				Optional: true,
			},
			common.ToSnakeCase("LbServerGroup"): schema.SingleNestedAttribute{
				Description: "Details of the LB Server Group.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
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
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : Server group for web servers\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
					common.ToSnakeCase("LbMethod"): schema.StringAttribute{
						Description: "The load balancing method.\n" +
							"  - example : ROUND_ROBIN\n" +
							"  - pattern : ROUND_ROBIN | LEAST_CONNECTION | IP_HASH | WEIGHTED_ROUND_ROBIN | WEIGHTED_LEAST_CONNECTION\n",
						Optional: true,
					},
					common.ToSnakeCase("LbName"): schema.StringAttribute{
						Description: "The name of the LoadBalancer.\n" +
							"  - example : LoadBalancer01\n",
						Optional: true,
					},
					common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
						Description: "The LoadBalancer ID associated with the server group.\n" +
							"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
						Optional: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the LB Server Group.\n" +
							"  - example : ACTIVE\n" +
							"  - pattern : CREATING | ACTIVE | DELETING | ERROR | EDITING\n",
						Optional: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LB Server Group.\n" +
							"  - example : ServerGroup01\n" +
							"  - minLength : 1\n" +
							"  - maxLength : 63\n" +
							"  - pattern : ^[a-zA-Z0-9._-]+$\n",
						Optional: true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "The protocol for the server group.\n" +
							"  - example : TCP\n" +
							"  - pattern : TCP | UDP\n",
						Optional: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The VPC ID where the resource is located.\n" +
							"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
						Optional: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The subnet ID where the resource is located.\n" +
							"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
						Optional: true,
					},
					common.ToSnakeCase("LbHealthCheckId"): schema.StringAttribute{
						Description: "The LB Health Check ID.\n" +
							"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
						Optional: true,
					},
				},
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
