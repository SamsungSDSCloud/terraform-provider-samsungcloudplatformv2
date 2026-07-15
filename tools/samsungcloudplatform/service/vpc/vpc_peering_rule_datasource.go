package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
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
	client  *vpcv1d2.Client
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
				Required: true,
				Description: "The identifier of the VPC peering.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Description: "The number of items per page. \n" +
					"  - Example: 20",
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Description: "The page number for pagination. \n" +
					"  - Example: 0",
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Optional: true,
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for decending order. \n" +
					"  - Example: created_at:desc",
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Optional: true,
				Description: "The unique identifier of the VPC peering rule.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
			},
			common.ToSnakeCase("SourceVpcId"): schema.StringAttribute{
				Optional: true,
				Description: "The identifier of the source VPC.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
			},
			common.ToSnakeCase("SourceVpcType"): schema.StringAttribute{
				Optional: true,
				Description: "The type of the source VPC.\n" +
                    "  - Example : REQUESTER_VPC | APPROVER_VPC\n" +
                    "  - Reference : VpcPeeringRuleDestinationVpcType",
			},
			common.ToSnakeCase("DestinationVpcId"): schema.StringAttribute{
				Optional: true,
				Description: "The identifier of the destination VPC.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
			},
			common.ToSnakeCase("DestinationVpcType"): schema.StringAttribute{
				Optional: true,
				Description: "The type of the destination VPC.\n" +
                "  - Example : REQUESTER_VPC | APPROVER_VPC\n" +
                "  - Reference : VpcPeeringRuleDestinationVpcType",
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Optional: true,
				Description: "The destination IP address range in CIDR notation.\n" +
					"  - example : 10.10.10.0/24",
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Optional: true,
				Description: "The current lifecycle state of the VPC peering rule.\n" +
					"  - example : ACTIVE",
			},

			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
                Description: "The total number of VPC peering rules.\n" +
                    "  - example : 2",
                Computed: true,
            },
			common.ToSnakeCase("SortFinal"): schema.ListAttribute{
				Description: "List of sort condition \n" +
					"  - example : [\"created_at:desc\"]",
				ElementType: types.StringType,
				Computed:    true,
			},
			common.ToSnakeCase("VpcPeeringRules"): schema.ListNestedAttribute{
				Description: "List of VPC peering rules",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Computed: true,
							Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
								"  - Example: 2024-05-17T00:23:17Z",
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Computed: true,
							Description: "The user id that created the resource.\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
						common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
							Computed: true,
							Description: "The destination IP address range in CIDR notation.\n" +
								"  - example : 10.10.10.0/24",
						},
						common.ToSnakeCase("DestinationVpcId"): schema.StringAttribute{
							Computed: true,
							Description: "The identifier of the destination VPC.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						},
						common.ToSnakeCase("DestinationVpcName"): schema.StringAttribute{
							Computed: true,
							Description: "The name of the destination VPC.\n" +
								"  - example : resourceName",
						},
						common.ToSnakeCase("DestinationVpcType"): schema.StringAttribute{
							Computed: true,
							Description: "The type of the destination VPC.\n" +
								"  - Example: REQUESTER_VPC | APPROVER_VPC\n" +
								"  - Reference: VpcPeeringRuleDestinationVpcType",
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Computed: true,
							Description: "The unique identifier of the VPC peering rule.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Computed: true,
							Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
								"  - Example: 2024-05-17T00:23:17Z",
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Computed: true,
							Description: "The user id that modified the resource.\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
						common.ToSnakeCase("SourceVpcId"): schema.StringAttribute{
							Computed: true,
							Description: "The identifier of the source VPC.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						},
						common.ToSnakeCase("SourceVpcName"): schema.StringAttribute{
							Computed: true,
							Description: "The name of the source VPC.\n" +
								"  - example : resourceName",
						},
						common.ToSnakeCase("SourceVpcType"): schema.StringAttribute{
							Computed: true,
							Description: "The type of the source VPC.\n" +
								"  - Example: REQUESTER_VPC | APPROVER_VPC\n" +
								"  - Reference: VpcPeeringRuleDestinationVpcType",
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Computed: true,
							Description: "The current lifecycle state of the VPC peering rule.\n" +
								"  - Example : CREATING | ACTIVE | DELETING | DELETED | ERROR\n" +
								"  - Reference: RoutingRuleState",
						},
						common.ToSnakeCase("VpcPeeringId"): schema.StringAttribute{
							Computed: true,
							Description: "The identifier of the VPC peering.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
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

	d.client = inst.Client.VpcV1Dot2
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *vpcVpcPeeringRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpcv1d2.VpcPeeringRuleDataSource

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
		vpcpeeringState := vpcv1d2.VpcPeeringRule{
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

	state.TotalCount = types.Int32Value(int32(data.Count))
	state.Page = types.Int32Value(data.Page)
	state.Size = types.Int32Value(data.Size)
	for _, sortVal := range data.Sort {
		state.SortFinal = append(state.SortFinal, types.StringValue(sortVal))
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
