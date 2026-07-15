package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &tgwRoutingRuleDataSources{}
	_ datasource.DataSourceWithConfigure = &tgwRoutingRuleDataSources{}
)

func NewTransitGatewayRoutingRuleDataSources() datasource.DataSource {
	return &tgwRoutingRuleDataSources{}
}

type tgwRoutingRuleDataSources struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

func (d *tgwRoutingRuleDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_transit_gateway_rules"
}

func (d *tgwRoutingRuleDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of transit gateway rule.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
				Description: "The identifier of the transit gateway that the rule belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("page"): schema.Int32Attribute{
				Optional: true,
				Description: "The page number for pagination.\n" +
					"  - example : 2",
				MarkdownDescription: "The page number for pagination.\n" +
                    "  - example : 2",
				Validators: []validator.Int32{
					int32validator.Between(0, 99999),
				},
			},			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for decending order.\n" +
					"  - example : created_at:desc",
				MarkdownDescription: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for decending order.\n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the rule.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("TgwConnectionVpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the transit gateway connection belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("TgwConnectionVpcName"): schema.StringAttribute{
				Description: "The name of the VPC that the transit gateway connection belongs to.\n" +
					"  - example : vpcName",
				Optional: true,
			},
			common.ToSnakeCase("SourceType"): schema.StringAttribute{
				Description: "The type of the source\n" +
					"  - enum:VPC, TGW\n" +
					"  - example:VPC",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"VPC",
						"TGW",
					),
				},
			},
			common.ToSnakeCase("DestinationType"): schema.StringAttribute{
				Description: "The type of the destination.\n" +
					"  - example : VPC | TGW",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"VPC",
						"TGW",
					),
				},
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "The destination IP address range in CIDR notation.\n" +
					"  - example : 10.10.10.0/24",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the rule.\n" +
					"  - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING\n" +
					"  - example:ACTIVE",
				MarkdownDescription: "The current lifecycle state of the rule.\n" +
					"  - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING\n" +
					"  - example:ACTIVE",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"CREATING",
						"ACTIVE",
						"DELETED",
						"ERROR",
						"DELETING",
						"EDITING",
					),
				},
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Computed:            true,
				Description:         "count\n  - example: 20",
				MarkdownDescription: "count\n  - example: 20",
			},
			common.ToSnakeCase("rule_type"): schema.StringAttribute{
				Description: "The type of the rule.\n" +
					"  - enum: TGW_VPC, TGW_UPLINK\n" +
					"  - example:TGW_VPC",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"TGW_VPC",
						"TGW_UPLINK",
					),
				},
			},
			common.ToSnakeCase("RoutingRules"): schema.ListNestedAttribute{
				Description: "Routing rules",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The identifier of the account that owns the rule.\n" +
								"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
							Computed: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
                            Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
                                "  - example : 2024-05-17T00:23:17Z",
                            Computed: true,
                        },
                        common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
                            Description: "The user id that created the resource.\n" +
                                "  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
                            Computed: true,
                        },
                        common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
                            Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
                                "  - example : 2024-05-17T00:23:17Z",
                            Computed: true,
                        },
                        common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
                            Description: "The user id that modified the resource.\n" +
                                "  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
                            Computed: true,
                        },
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource.\n" +
								"  - example : description",
							Computed: true,
						},
						common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
							Description: "The destination IP address range in CIDR notation.\n" +
								"  - example : 10.10.10.0/24",
							Computed: true,
						},
						common.ToSnakeCase("DestinationResourceId"): schema.StringAttribute{
							Description: "The identifier of the destination resource.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("DestinationResourceName"): schema.StringAttribute{
							Description: "The name of the destination resource.\n" +
								"  - example : resourcename",
							Computed: true,
						},
						common.ToSnakeCase("DestinationType"): schema.StringAttribute{
							Description: "The type of the destination.\n" +
								"  - example : VPC | TGW",
							Computed: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the rule.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("SourceResourceId"): schema.StringAttribute{
							Description: "The identifier of the source resource.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("SourceResourceName"): schema.StringAttribute{
							Description: "The name of the source resource.\n" +
								"  - example : sourceName",
							Computed: true,
						},
						common.ToSnakeCase("SourceType"): schema.StringAttribute{
							Description: "The type of the source.\n" +
								"  - example : VPC | TGW",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the rule.\n" +
								"  - example : ACTIVE",
							Computed: true,
						},
						common.ToSnakeCase("TgwConnectionVpcId"): schema.StringAttribute{
							Description: "The identifier of the VPC that the transit gateway connection belongs to.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("TgwConnectionVpcName"): schema.StringAttribute{
							Description: "The name of the VPC that the transit gateway connection belongs to.\n" +
								"  - example : vpcName",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *tgwRoutingRuleDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tgwRoutingRuleDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.TransitGatewayRuleDataSource

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := d.client.GetTGWRuleList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading routing rule",
			"Could not read routing rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	for _, routingRule := range data.TransitGatewayRules {
		routingRuleState := vpc.RoutingRule{
			AccountId:               types.StringValue(routingRule.AccountId),
			CreatedAt:               types.StringValue(routingRule.CreatedAt.Format(time.RFC3339)),
			CreatedBy:               types.StringValue(routingRule.CreatedBy),
			Description:             types.StringValue(routingRule.Description),
			DestinationCidr:         types.StringValue(routingRule.DestinationCidr),
			DestinationResourceId:   types.StringPointerValue(routingRule.DestinationResourceId.Get()),
			DestinationResourceName: types.StringPointerValue(routingRule.DestinationResourceName.Get()),
			DestinationType:         types.StringValue(string(routingRule.DestinationType)),
			Id:                      types.StringValue(routingRule.Id),
			ModifiedAt:              types.StringValue(routingRule.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:              types.StringValue(routingRule.ModifiedBy),
			SourceResourceId:        types.StringPointerValue(routingRule.SourceResourceId.Get()),
			SourceResourceName:      types.StringPointerValue(routingRule.SourceResourceName.Get()),
			SourceType:              types.StringValue(string(routingRule.SourceType)),
			State:                   types.StringValue(string(routingRule.State)),
			TgwConnectionVpcId:      types.StringPointerValue(routingRule.TgwConnectionVpcId.Get()),
			TgwConnectionVpcName:    types.StringPointerValue(routingRule.TgwConnectionVpcName.Get()),
		}

		state.RoutingRules = append(state.RoutingRules, routingRuleState)
	}

	state.TotalCount = types.Int32Value(data.Count)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
