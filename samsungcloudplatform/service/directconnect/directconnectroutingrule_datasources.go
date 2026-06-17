package directconnect

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	directconnectv1d1 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/directconnectv1d1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &networkDirectConnectRoutingRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &networkDirectConnectRoutingRuleDataSource{}
)

// NewNetworkDirectConnectRoutingRuleDataSource is a helper function to simplify the provider implementation.
func NewNetworkDirectConnectRoutingRuleDataSource() datasource.DataSource {
	return &networkDirectConnectRoutingRuleDataSource{}
}

// networkRoutingRuleDataSource is the data source implementation.
type networkDirectConnectRoutingRuleDataSource struct {
	config  *scpsdk.Configuration
	client  *directconnectv1d1.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *networkDirectConnectRoutingRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directconnect_routing_rules"
}

// Schema defines the schema for the data source.
func (d *networkDirectConnectRoutingRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of routing rule.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page. \n" +
					"  - example : 20 \n" +
					"  - minimum : 0",
				Optional: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination. \n" +
					"  - example : 0 \n" +
					"  - minimum : 0",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order. \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("DirectConnectId"): schema.StringAttribute{
				Description: "The identifier of the direct connect.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the routing rule.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("DestinationType"): schema.StringAttribute{
				Description: "The type of the routing destination. In the VPC, the Direct Connect direction is ON_PREMISE, in the opposite direction—from Direct Connect toward the VPC—the direction is VPC.\n" +
					"  -  example : ON-PREMISE | VPC",
				Optional: true,
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "The destination IP address range in CIDR notation. \n" +
					"  - example : 10.10.10.0/24",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the routing rule. \n" +
					"  - example : CREATING | ACTIVE | DELETING | ERROR",
				Optional: true,
			},

			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "The total number of Direct Connect routing rule.\n" +
					"  - example : 5",
				Computed: true,
			},
			common.ToSnakeCase("SortFinal"): schema.ListAttribute{
				Description: "List of sort condition \n" +
					"  - example : [\"created_at:desc\"]",
				ElementType: types.StringType,
				Computed:    true,
			},
			common.ToSnakeCase("RoutingRules"): schema.ListNestedAttribute{
				Description: "A list of routing rule.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the routing rule.\n" +
								"  - example : fe860e0af0c04dcd8182b84f907f31f4",
							Computed: true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The identifier of the account that owns the direct connect.\n" +
								"  - example : 27bb070b564349f8a31cc60734cc36a5",
							Computed: true,
						},
						common.ToSnakeCase("OwnerId"): schema.StringAttribute{
							Description: "The identifier of the routing rule owner.\n" +
								"  - example : 0fdd87aab8cb46f59b7c1f81ed03fb3e",
							Computed: true,
						},
						common.ToSnakeCase("OwnerType"): schema.StringAttribute{
							Description: "The type of the routing rule owner.\n" +
								"  - example : DIRECT_CONNECT",
							Computed: true,
						},
						common.ToSnakeCase("DestinationType"): schema.StringAttribute{
							Description: "The type of the routing destination.In the VPC, the Direct Connect direction is ON_PREMISE, in the opposite direction—from Direct Connect toward the VPC—the direction is VPC.\n" +
								"  -  example : ON-PREMISE | VPC",
							Computed: true,
						},
						common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
							Description: "The destination IP address range in CIDR notation.\n" +
								"  - example : 10.10.10.0/24",
							Computed: true,
						},
						common.ToSnakeCase("DestinationResourceId"): schema.StringAttribute{
							Description: "The identifier of the destination resource.When the Destination Type is VPC, provide the VpcId.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("DestinationResourceName"): schema.StringAttribute{
							Description: "The name of the destination resource.When the Destination Type is VPC, provide the Vpc name.\n" +
								"  - example : Resource Name",
							Computed: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Enter a brief explanation or note about this routing rule. This help identify the purpose or usage of the resource.\n" +
								"  - example : Routing Rule description\n" +
								"  - maxLength : 50\n" +
								"  - minLength : 1",
							Computed: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user id that created the resource.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user id that last modified the resource.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the routing rule.\n" +
								"  - example : CREATING | ACTIVE | EDITING | DELETING | ERROR",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *networkDirectConnectRoutingRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.DirectConnectV1d1
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *networkDirectConnectRoutingRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state directconnectv1d1.RoutingRuleDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetRoutingRuleList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading routing rule",
			"Could not read routing rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	state.TotalCount = types.Int32Value(data.Count)
	state.Page = types.Int32Value(data.Page)
	state.Size = types.Int32Value(data.Size)
	for _, sortVal := range data.Sort {
		state.SortFinal = append(state.SortFinal, types.StringValue(sortVal))
	}

	for _, routingRule := range data.RoutingRules {
		routingRuleState := directconnectv1d1.RoutingRule{
			Id:                      types.StringValue(routingRule.Id),
			AccountId:               types.StringValue(routingRule.AccountId),
			OwnerId:                 types.StringValue(routingRule.OwnerId),
			OwnerType:               types.StringValue(string(routingRule.OwnerType)),
			DestinationType:         types.StringValue(string(routingRule.DestinationType)),
			DestinationCidr:         types.StringValue(routingRule.DestinationCidr),
			DestinationResourceId:   types.StringPointerValue(routingRule.DestinationResourceId.Get()),
			DestinationResourceName: types.StringPointerValue(routingRule.DestinationResourceName.Get()),
			Description:             types.StringValue(routingRule.Description),
			CreatedAt:               types.StringValue(routingRule.CreatedAt.Format(time.RFC3339)),
			CreatedBy:               types.StringValue(routingRule.CreatedBy),
			ModifiedAt:              types.StringValue(routingRule.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:              types.StringValue(routingRule.ModifiedBy),
			State:                   types.StringValue(string(routingRule.State)),
		}

		state.RoutingRules = append(state.RoutingRules, routingRuleState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
