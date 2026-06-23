package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loadbalancer"
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
	_ datasource.DataSource              = &loadbalancerLbServerGroupDataSources{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLbServerGroupDataSources{}
)

// NewLoadBalancerLbServerGroupDataSources is a helper function to simplify the provider implementation.
func NewLoadBalancerLbServerGroupDataSources() datasource.DataSource {
	return &loadbalancerLbServerGroupDataSources{}
}

// loadbalancerLbServerGroupDataSources is the data source implementation.
type loadbalancerLbServerGroupDataSources struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLbServerGroupDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_server_groups"
}

// Schema defines the schema for the data source.
func (d *loadbalancerLbServerGroupDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Get List of Lb Server Groups.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example : 20\n",
				Optional: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number.\n" +
					"  - example : 0\n",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sort order.\n" +
					"  - example : name:asc\n",
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
			common.ToSnakeCase("Protocol"): schema.ListAttribute{
				Description: "The protocol used for the server group.\n" +
					"  - example : TCP\n" +
					"  - pattern : TCP | UDP\n",
				Optional:    true,
				ElementType: types.StringType,
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
			common.ToSnakeCase("LbServerGroups"): schema.ListNestedAttribute{
				Description: "List of LB Server Groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the LB Server Group.\n" +
								"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
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
						common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
							Description: "The LoadBalancer ID associated with the server group.\n" +
								"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
							Optional: true,
						},
						common.ToSnakeCase("LbName"): schema.StringAttribute{
							Description: "The name of the LoadBalancer.\n" +
								"  - example : LoadBalancer01\n",
							Optional: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current state of the LB Server Group.\n" +
								"  - example : ACTIVE\n" +
								"  - pattern : CREATING | ACTIVE | DELETING | ERROR | EDITING\n",
							Optional: true,
						},
						common.ToSnakeCase("VpcId"): schema.StringAttribute{
							Description: "The VPC ID where the resource is located.\n" +
								"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e\n",
							Optional: true,
						},
						common.ToSnakeCase("LbServerGroupMemberCount"): schema.Int32Attribute{
							Description: "The number of members in the LB Server Group.\n" +
								"  - example : 3\n",
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
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadbalancerLbServerGroupDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *loadbalancerLbServerGroupDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state loadbalancer.LbServerGroupDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLbServerGroupList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read LbServerGroups",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, lbservergroup := range data.LbServerGroups {
		lbservergroupState := loadbalancer.LbServerGroup{
			Id:                       types.StringValue(lbservergroup.Id),
			Name:                     types.StringValue(lbservergroup.Name),
			Protocol:                 types.StringValue(string(lbservergroup.Protocol)),
			LoadbalancerId:           types.StringPointerValue(lbservergroup.LoadbalancerId.Get()),
			LbName:                   virtualserverutil.ToNullableStringValue(lbservergroup.LbName.Get()),
			State:                    types.StringValue(lbservergroup.State),
			VpcId:                    types.StringValue(lbservergroup.VpcId),
			LbServerGroupMemberCount: types.Int32Value(lbservergroup.LbServerGroupMemberCount),
			CreatedAt:                types.StringValue(lbservergroup.CreatedAt.Format(time.RFC3339)),
			CreatedBy:                types.StringValue(lbservergroup.CreatedBy),
			ModifiedAt:               types.StringValue(lbservergroup.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:               types.StringValue(lbservergroup.ModifiedBy),
		}

		state.LbServerGroups = append(state.LbServerGroups, lbservergroupState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
