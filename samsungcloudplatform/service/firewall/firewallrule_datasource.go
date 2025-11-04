package firewall

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &firewallFirewallRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &firewallFirewallRuleDataSource{}
)

// NewFirewallFirewallRuleDataSource is a helper function to simplify the provider implementation.
func NewFirewallFirewallRuleDataSource() datasource.DataSource {
	return &firewallFirewallRuleDataSource{}
}

// firewallFirewallRuleDataSource is the data source implementation.
type firewallFirewallRuleDataSource struct {
	config  *scpsdk.Configuration
	client  *firewall.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *firewallFirewallRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_firewall_rule"
}

// Schema defines the schema for the data source.
func (d *firewallFirewallRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Firewall rule",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Firewall Rule ID \n" +
					"  - example : fd5aa79d25a84f52b02a954c125aa8d9",
				Required: true,
			},
			common.ToSnakeCase("FirewallId"): schema.StringAttribute{
				Description: "Firewall Id \n" +
					"  - example : 68db67f78abd405da98a6056a8ee42af",
				Required: true,
			},
			common.ToSnakeCase("FirewallRule"): schema.SingleNestedAttribute{
				Description: "Firewall rule",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("FirewallId"): schema.StringAttribute{
						Description: "FirewallId",
						Computed:    true,
					},
					common.ToSnakeCase("Sequence"): schema.Int32Attribute{
						Description: "Sequence",
						Computed:    true,
					},
					common.ToSnakeCase("SourceInterface"): schema.StringAttribute{
						Description: "SourceInterface",
						Computed:    true,
					},
					common.ToSnakeCase("SourceAddress"): schema.ListAttribute{
						Description: "SourceAddress",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("DestinationInterface"): schema.StringAttribute{
						Description: "DestinationInterface",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationAddress"): schema.ListAttribute{
						Description: "DestinationAddress",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("Service"): schema.ListNestedAttribute{
						Description: "Service",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("ServiceType"): schema.StringAttribute{
									Description: "ServiceType",
									Computed:    true,
								},
								common.ToSnakeCase("ServiceValue"): schema.StringAttribute{
									Description: "ServiceValue",
									Computed:    true,
								},
							},
						},
					},
					common.ToSnakeCase("Action"): schema.StringAttribute{
						Description: "Action",
						Computed:    true,
					},
					common.ToSnakeCase("Direction"): schema.StringAttribute{
						Description: "Direction",
						Computed:    true,
					},
					common.ToSnakeCase("VendorRuleId"): schema.StringAttribute{
						Description: "VendorRuleId",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
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
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *firewallFirewallRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Firewall
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *firewallFirewallRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state firewall.FirewallRuleDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var defaultSize types.Int32
	var defaultPage types.Int32
	var defaultSort types.String
	var defaultSrcIp types.String
	var defaultDstIp types.String
	var defaultDescription types.String
	var defaultState types.List
	var defaultStatus types.String
	var defaultFetchAll types.Bool

	ids, err := GetFirewallRuleList(d.clients, defaultPage, defaultSize, defaultSort, state.FirewallId,
		defaultSrcIp, defaultDstIp, defaultDescription, defaultState, defaultStatus, defaultFetchAll)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read firewall rule",
			err.Error(),
		)
	}

	if len(ids) > 0 {
		exist := false
		for _, v := range ids {
			if v == state.Id {
				exist = true
				break
			}
		}

		if exist {
			data, err := d.client.GetFirewallRule(state.Id.ValueString()) // client 를 호출한다.
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Reading Firewall Rule",
					"Could not read Firewall Rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
				)
				return
			}

			firewallRuleElement := data.FirewallRule

			sourceAddresses := make([]string, 0, len(firewallRuleElement.SourceAddress))
			for _, address := range firewallRuleElement.SourceAddress {
				sourceAddresses = append(sourceAddresses, address)
			}
			destinationAddresses := make([]string, 0, len(firewallRuleElement.DestinationAddress))
			for _, address := range firewallRuleElement.DestinationAddress {
				destinationAddresses = append(destinationAddresses, address)
			}
			services := make([]firewall.FirewallPort, 0, len(firewallRuleElement.Service))
			for _, service := range firewallRuleElement.Service {
				services = append(services, firewall.FirewallPort{
					ServiceType:  types.StringValue(string(service.ServiceType)),
					ServiceValue: types.StringValue(*service.ServiceValue),
				})
			}

			firewallRuleModel := firewall.FirewallRule{
				Id:                   types.StringValue(firewallRuleElement.Id),
				Name:                 types.StringPointerValue(firewallRuleElement.Name.Get()),
				FirewallId:           types.StringValue(firewallRuleElement.FirewallId),
				Sequence:             types.Int32Value(firewallRuleElement.Sequence),
				SourceInterface:      types.StringValue(firewallRuleElement.SourceInterface),
				SourceAddress:        sourceAddresses,
				DestinationInterface: types.StringValue(firewallRuleElement.DestinationInterface),
				DestinationAddress:   destinationAddresses,
				Service:              services,
				Action:               types.StringValue(string(firewallRuleElement.Action)),
				Direction:            types.StringValue(string(firewallRuleElement.Direction)),
				VendorRuleId:         types.StringValue(firewallRuleElement.VendorRuleId),
				Description:          types.StringPointerValue(firewallRuleElement.Description.Get()),
				State:                types.StringValue(string(firewallRuleElement.State)),
				Status:               types.StringValue(string(firewallRuleElement.Status)),
				CreatedAt:            types.StringValue(firewallRuleElement.CreatedAt.Format(time.RFC3339)),
				CreatedBy:            types.StringValue(firewallRuleElement.CreatedBy),
				ModifiedAt:           types.StringValue(firewallRuleElement.ModifiedAt.Format(time.RFC3339)),
				ModifiedBy:           types.StringValue(firewallRuleElement.ModifiedBy),
			}

			firewallRuleObjectValue, _ := types.ObjectValueFrom(ctx, firewallRuleModel.AttributeTypes(), firewallRuleModel)
			state.FirewallRule = firewallRuleObjectValue
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
