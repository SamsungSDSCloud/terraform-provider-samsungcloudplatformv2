package directconnect

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/directconnect"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
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
	client  *directconnect.Client
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
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "Limit \n" +
					"  - example : 10 \n" +
					"  - maximum : 10000 \n" +
					"  - minimum : 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Marker"): schema.StringAttribute{
				Description: "Marker \n" +
					"  - example : 607e0938521643b5b4b266f343fae693 \n" +
					"  - maxLength : 64 \n" +
					"  - minLength : 1",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("DirectConnectId"): schema.StringAttribute{
				Description: "Direct Connect ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Routing Rule ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("DestinationType"): schema.StringAttribute{
				Description: "Destination Type \n" +
					"  - example : ON-PREM | VPC",
				Optional: true,
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "Destination CIDR \n" +
					"  - example : 10.10.10.0/24",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State \n" +
					"  - example : CREATING | ACTIVE | DELETING | ERROR",
				Optional: true,
			},
			common.ToSnakeCase("RoutingRules"): schema.ListNestedAttribute{
				Description: "A list of routing rule.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Computed:    true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "AccountId",
							Computed:    true,
						},
						common.ToSnakeCase("OwnerId"): schema.StringAttribute{
							Description: "OwnerId",
							Computed:    true,
						},
						common.ToSnakeCase("OwnerType"): schema.StringAttribute{
							Description: "OwnerType",
							Computed:    true,
						},
						common.ToSnakeCase("DestinationType"): schema.StringAttribute{
							Description: "DestinationType",
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
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
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
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "ModifiedAt",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "ModifiedBy",
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State",
							Computed:    true,
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

	d.client = inst.Client.DirectConnect
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *networkDirectConnectRoutingRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state directconnect.RoutingRuleDataSource

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
	for _, routingRule := range data.RoutingRules {
		routingRuleState := directconnect.RoutingRule{
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
