package vpc

import (
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
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
		Description: "List of vpc transit gateway rule.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
				Description: "transit gateway id",
				Required:    true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "id",
				Optional:    true,
			},
			common.ToSnakeCase("TgwConnectionVpcId"): schema.StringAttribute{
				Description: "tgw connection vpc id",
				Optional:    true,
			},
			common.ToSnakeCase("TgwConnectionVpcName"): schema.StringAttribute{
				Description: "tgw connection vpc name",
				Optional:    true,
			},
			common.ToSnakeCase("SourceType"): schema.StringAttribute{
				Description: "source type" +
					" - enum(VPC, TGW)",
				Optional: true,
			},
			common.ToSnakeCase("DestinationType"): schema.StringAttribute{
				Description: "destination type" +
					" - enum(VPC, TGW)",
				Optional: true,
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "destination cidr",
				Optional:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State" +
					" - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING",
				Optional: true,
			},
			common.ToSnakeCase("RoutingRules"): schema.ListNestedAttribute{
				Description: "transit gateway rules",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "AccountId",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "CreatedAt",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "CreatedBy",
							Computed:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
							Computed:    true,
						},
						common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
							Description: "DestinationCidr",
							Computed:    true,
						},
						common.ToSnakeCase("DestinationResourceId"): schema.StringAttribute{
							Description: "DestinationResourceId",
							Computed:    true,
						},
						common.ToSnakeCase("DestinationResourceName"): schema.StringAttribute{
							Description: "DestinationResourceName",
							Computed:    true,
						},
						common.ToSnakeCase("DestinationType"): schema.StringAttribute{
							Description: "DestinationType",
							Computed:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "id",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "ModifiedAt",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "ModifiedBy",
							Computed:    true,
						},
						common.ToSnakeCase("SourceResourceId"): schema.StringAttribute{
							Description: "SourceResourceId",
							Computed:    true,
						},
						common.ToSnakeCase("SourceResourceName"): schema.StringAttribute{
							Description: "SourceResourceName",
							Computed:    true,
						},
						common.ToSnakeCase("SourceType"): schema.StringAttribute{
							Description: "SourceType",
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State",
							Computed:    true,
						},
						common.ToSnakeCase("TgwConnectionVpcId"): schema.StringAttribute{
							Description: "TgwConnectionVpcId",
							Computed:    true,
						},
						common.ToSnakeCase("TgwConnectionVpcName"): schema.StringAttribute{
							Description: "TgwConnectionVpcName",
							Computed:    true,
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

	d.client = inst.Client.Vpc
	d.clients = inst.Client
}

func (d *tgwRoutingRuleDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.RoutingRuleDataSource

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := d.client.GetRoutingRuleList(ctx, state)
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

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
