package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vpcVpcPeeringRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcVpcPeeringRuleDataSource{}
)

// NewVpcVpcPeeringRuleDataSource is a helper function to simplify the provider implementation.
func NewVpcVpcPeeringRuleDataSource() datasource.DataSource {
	return &vpcVpcPeeringRuleDataSource{}
}

// vpcNatGatewayDataSource is the data source implementation.
type vpcVpcPeeringRuleDataSource struct {
	_config *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcVpcPeeringRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc_peering_rules"
}

// Schema defines the schema for the data source.
func (d *vpcVpcPeeringRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of VPC peering rules",
		Attributes: map[string]schema.Attribute{
			// Input Params
			common.ToSnakeCase("VpcPeeringId"): schema.StringAttribute{
				Required:    true,
				Description: "VPC Peering ID",
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Optional: true,
				Description: "Size \n" +
					"  - Example: 20",
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Optional: true,
				Description: "Page \n" +
					"  - Example: 0",
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Optional: true,
				Description: "Sort \n" +
					"  - Example: created_at:desc",
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Optional:    true,
				Description: "VPC Peering Rule ID",
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Optional:    true,
				Description: "Name",
			},
			common.ToSnakeCase("SourceVpcId"): schema.StringAttribute{
				Optional:    true,
				Description: "Source VPC ID",
			},
			common.ToSnakeCase("SourceVpcType"): schema.StringAttribute{
				Optional:    true,
				Description: "Source VPC Type",
			},
			common.ToSnakeCase("DestinationVpcId"): schema.StringAttribute{
				Optional:    true,
				Description: "Destination VPC ID",
			},
			common.ToSnakeCase("DestinationVpcType"): schema.StringAttribute{
				Optional:    true,
				Description: "Destination VPC Type",
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Optional:    true,
				Description: "Destination CIDR",
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Optional:    true,
				Description: "State",
			},

			// Response
			common.ToSnakeCase("VpcPeeringRules"): schema.ListNestedAttribute{
				Description: "List of VPC peering rules",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Computed: true,
							Description: "Created At\n" +
								"  - Example: 2024-05-17T00:23:17Z",
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Computed: true,
							Description: "Created By\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
						common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
							Computed:    true,
							Description: "Destination CIDR",
						},
						common.ToSnakeCase("DestinationVpcId"): schema.StringAttribute{
							Computed:    true,
							Description: "Destination VPC ID",
						},
						common.ToSnakeCase("DestinationVpcName"): schema.StringAttribute{
							Computed:    true,
							Description: "Destination VPC Name",
						},
						common.ToSnakeCase("DestinationVpcType"): schema.StringAttribute{
							Computed: true,
							Description: "Destination VPC Type\n" +
								"  - Example: REQUESTER_VPC | APPROVER_VPC\n" +
								"  - Reference: VpcPeeringRuleDestinationVpcType",
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Computed:    true,
							Description: "VPC Peering Rule ID",
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Computed: true,
							Description: "Modified At\n" +
								"  - Example: 2024-05-17T00:23:17Z",
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Computed: true,
							Description: "Modified By\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
						common.ToSnakeCase("SourceVpcId"): schema.StringAttribute{
							Computed:    true,
							Description: "Source VPC ID",
						},
						common.ToSnakeCase("SourceVpcName"): schema.StringAttribute{
							Computed:    true,
							Description: "Source VPC Name",
						},
						common.ToSnakeCase("SourceVpcType"): schema.StringAttribute{
							Computed: true,
							Description: "Source VPC Type\n" +
								"  - Example: REQUESTER_VPC | APPROVER_VPC\n" +
								"  - Reference: VpcPeeringRuleDestinationVpcType",
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Computed: true,
							Description: "State\n" +
								"  - Example : CREATING | ACTIVE | DELETING | DELETED | ERROR\n" +
								"  - Reference: RoutingRuleState",
						},
						common.ToSnakeCase("VpcPeeringId"): schema.StringAttribute{
							Computed:    true,
							Description: "VPC Peering ID",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcVpcPeeringRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Vpc
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *vpcVpcPeeringRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.VpcPeeringRuleDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetVpcPeeringRuleList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to retrieve VPC peering rules",
			fmt.Sprintf("An error occurred while retrieving VPC peering rules: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	// Map response body to model
	for _, vpcPeering := range data.VpcPeeringRules {
		vpcpeeringState := vpc.VpcPeeringRule{
			Id:                 types.StringValue(vpcPeering.Id),
			CreatedAt:          types.StringValue(vpcPeering.CreatedAt.Format(time.RFC3339)),
			CreatedBy:          types.StringValue(vpcPeering.CreatedBy),
			DestinationCidr:    types.StringValue(vpcPeering.DestinationCidr),
			DestinationVpcId:   types.StringValue(vpcPeering.DestinationVpcId),
			DestinationVpcName: types.StringValue(vpcPeering.DestinationVpcName),
			DestinationVpcType: types.StringValue(string(*vpcPeering.DestinationVpcType.Ptr())),
			ModifiedAt:         types.StringValue(vpcPeering.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:         types.StringValue(vpcPeering.ModifiedBy),
			SourceVpcId:        types.StringValue(vpcPeering.SourceVpcId),
			SourceVpcName:      types.StringValue(vpcPeering.SourceVpcName),
			SourceVpcType:      types.StringValue(string(*vpcPeering.SourceVpcType.Ptr())),
			State:              types.StringValue(string(*vpcPeering.State.Ptr())),
			VpcPeeringId:       types.StringValue(vpcPeering.VpcPeeringId),
		}
		state.VpcPeeringRules = append(state.VpcPeeringRules, vpcpeeringState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
