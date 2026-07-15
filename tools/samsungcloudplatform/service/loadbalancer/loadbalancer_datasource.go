package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/loadbalancer" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadbalancerLoadbalancerDataSource{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLoadbalancerDataSource{}
)

// NewResourceManagerResourceGroupDataSources is a helper function to simplify the provider implementation.
func NewLoadbalancerLoadbalancerDataSource() datasource.DataSource {
	return &loadbalancerLoadbalancerDataSource{}
}

// resourceManagerResourceGroupDataSources is the data source implementation.
type loadbalancerLoadbalancerDataSource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLoadbalancerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_loadbalancer" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 복수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (d *loadbalancerLoadbalancerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Loadbalancer.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the LoadBalancer.\n" +
					"  - example : 46c681018e33453085ca7c8db54e0076\n",
				Optional: true,
			},
			common.ToSnakeCase("loadbalancer"): schema.SingleNestedAttribute{
				Description: "Details of the LoadBalancer.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
							"  - example : 2024-01-01T00:00:00Z\n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.\n" +
							"  - example : LoadBalancer for web traffic\n" +
							"  - maxLength : 255\n",
						Optional: true,
					},
					common.ToSnakeCase("FirewallId"): schema.StringAttribute{
						Description: "The firewall ID associated with the LoadBalancer.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("HealthCheckIp"): schema.ListAttribute{
						Description: "The list of health check IP addresses.\n" +
							"  - example : [\"192.168.1.1\", \"192.168.1.2\"]\n",
						ElementType: types.StringType,
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("LayerType"): schema.StringAttribute{
						Description: "The layer type of the Load Balancer.\n" +
							"  - example : L7\n" +
							"  - pattern : L4 | L7\n",
						Optional: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
							"  - example : 2024-01-01T00:00:00Z\n",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that last modified the resource.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the LoadBalancer.\n" +
							"  - example : LoadBalancer01\n" +
							"  - minLength : 1\n" +
							"  - maxLength : 63\n" +
							"  - pattern : ^[a-zA-Z0-9._-]+$\n",
						Optional: true,
					},
					common.ToSnakeCase("PublicNatEnabled"): schema.BoolAttribute{
						Description: "Whether public NAT is enabled.\n" +
							"  - example : true\n",
						Optional: true,
					},
					common.ToSnakeCase("ServiceIp"): schema.StringAttribute{
						Description: "The service IP address of the LoadBalancer.\n" +
							"  - example : 192.168.1.100\n",
						Optional: true,
					},
					common.ToSnakeCase("SourceNatIp"): schema.StringAttribute{
						Description: "The source NAT IP address.\n" +
							"  - example : 192.168.1.101\n",
						Optional: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the Load Balancer.\n" +
							"  - example : ACTIVE\n" +
							"  - pattern : CREATING | ACTIVE | DELETING | ERROR\n",
						Optional: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "The subnet ID where the LoadBalancer is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The VPC ID where the LoadBalancer is located.\n" +
							"  - example : 46c681018e33453085ca7c8db54e0076\n",
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadbalancerLoadbalancerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *loadbalancerLoadbalancerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state loadbalancer.LoadbalancerDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLoadbalancer(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show Loadbalancer",
			err.Error(),
		)
		return
	}

	var loadbalancerState = loadbalancer.LoadbalancerDetail{
		AccountId:        types.StringValue(data.Loadbalancer.AccountId),
		CreatedAt:        types.StringValue(data.Loadbalancer.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.Loadbalancer.CreatedBy),
		Description:      virtualserverutil.ToNullableStringValue(data.Loadbalancer.Description.Get()),
		FirewallId:       virtualserverutil.ToNullableStringValue(data.Loadbalancer.FirewallId.Get()),
		HealthCheckIp:    ToStringList(data.Loadbalancer.HealthCheckIp),
		Id:               types.StringValue(data.Loadbalancer.Id),
		LayerType:        types.StringValue(data.Loadbalancer.LayerType),
		ModifiedAt:       types.StringValue(data.Loadbalancer.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.Loadbalancer.ModifiedBy),
		Name:             types.StringValue(data.Loadbalancer.Name),
		PublicNatEnabled: common.ToNullableBoolValue(data.Loadbalancer.PublicNatEnabled.Get()),
		ServiceIp:        virtualserverutil.ToNullableStringValue(data.Loadbalancer.ServiceIp.Get()),
		SourceNatIp:      virtualserverutil.ToNullableStringValue(data.Loadbalancer.SourceNatIp.Get()),
		State:            types.StringValue(data.Loadbalancer.State),
		SubnetId:         types.StringValue(data.Loadbalancer.SubnetId),
		VpcId:            types.StringValue(data.Loadbalancer.VpcId),
	}

	loadbalancerObjectValue, _ := types.ObjectValueFrom(ctx, loadbalancerState.AttributeTypes(), loadbalancerState)
	state.LoadbalancerDetail = loadbalancerObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ToStringList(ips []string) basetypes.ListValue {
	if ips == nil {
		return types.ListNull(types.StringType)
	}

	elements := make([]attr.Value, len(ips))
	for i, v := range ips {
		elements[i] = types.StringValue(v)
	}
	return types.ListValueMust(types.StringType, elements)
}
