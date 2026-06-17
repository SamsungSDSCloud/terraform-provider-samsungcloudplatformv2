package firewall

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &firewallFirewallDataSource{}
	_ datasource.DataSourceWithConfigure = &firewallFirewallDataSource{}
)

func NewFirewallFirewallDataSource() datasource.DataSource {
	return &firewallFirewallDataSource{}
}

type firewallFirewallDataSource struct {
	config  *scpsdk.Configuration
	client  *firewall.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *firewallFirewallDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_firewall"
}

// Schema defines the schema for the data source.
func (d *firewallFirewallDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves details of a firewall instance.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 68db67f78abd405da98a6056a8ee42af",
				Required: true,
			},
			common.ToSnakeCase("Firewall"): schema.SingleNestedAttribute{
				Description: "The firewall resource details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the resource.\n" +
							"  - example: 68db67f78abd405da98a6056a8ee42af",
						Computed: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example: 297615908b8e4ec69520a99a6777add3",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the resource.\n" +
							"  - example: fw-web-prod",
						Computed: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The identifier of the VPC that the resource belongs to.\n" +
							"  - example: 6a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "The name of the VPC that the resource belongs to.\n" +
							"  - example: vpc-prod-01",
						Computed: true,
					},
					common.ToSnakeCase("Loggable"): schema.BoolAttribute{
						Description: "The flag indicating whether firewall flow logs are stored.\n" +
							"  - example: True",
						Computed: true,
					},
					common.ToSnakeCase("FwResourceId"): schema.StringAttribute{
						Description: "The resource ID of the associated firewall service.\n" +
							"  - example: 01a41c6f439d4be78de081cc9d78a4f2",
						Computed: true,
					},
					common.ToSnakeCase("PreProductId"): schema.StringAttribute{
						Description: "The identifier of the firewall pre‑product associated with the service.\n" +
							"  - example: a637b0c278064513ab90ac9e91a92e03",
						Computed: true,
					},
					common.ToSnakeCase("ProductType"): schema.StringAttribute{
						Description: "The type of the firewall service.\n" +
							"  - example: IGW",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the resource.\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "The current status of the resource.\n" +
							"  - example: ENABLE",
						Computed: true,
					},
					common.ToSnakeCase("TotalRuleCount"): schema.Int32Attribute{
						Description: "Number of rules in the firewall.\n" +
							"  - example: 5",
						Computed: true,
					},
					common.ToSnakeCase("FlavorName"): schema.StringAttribute{
						Description: "Firewall size.\n" +
							"  - example: MEDIUM",
						Computed: true,
					},
					common.ToSnakeCase("FlavorRuleQuota"): schema.Int32Attribute{
						Description: "Firewall rule quota based on firewall size.\n" +
							"  - example: 200",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
							"  - example: 2025-01-15T10:30:00Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user ID that created the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
							"  - example: 2025-06-01T14:22:00Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user ID that modified the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *firewallFirewallDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *firewallFirewallDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state firewall.FirewallDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var defaultSize types.Int32
	var defaultPage types.Int32
	var defaultSort types.String
	var defaultName types.String
	var defaultVpcName types.String
	var defaultProductType types.List
	var defaultState types.List

	ids, err := GetFirewallList(d.clients, defaultPage, defaultSize, defaultSort, defaultName, defaultVpcName, defaultProductType, defaultState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Firewall",
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
			data, err := d.client.GetFirewall(state.Id.ValueString()) // client 를 호출한다.
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Reading Firewall",
					"Could not read Firewall ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
				)
				return
			}

			firewallElement := data.Firewall

			firewallModel := firewall.Firewall{
				Id:              types.StringValue(firewallElement.Id),
				AccountId:       types.StringValue(firewallElement.AccountId),
				Name:            types.StringValue(firewallElement.Name),
				VpcId:           types.StringPointerValue(firewallElement.VpcId.Get()),
				VpcName:         types.StringPointerValue(firewallElement.VpcName.Get()),
				Loggable:        types.BoolValue(firewallElement.Loggable),
				FwResourceId:    types.StringValue(firewallElement.FwResourceId),
				PreProductId:    types.StringPointerValue(firewallElement.PreProductId),
				ProductType:     types.StringValue(string(firewallElement.ProductType)),
				State:           types.StringValue(string(firewallElement.State)),
				Status:          types.StringValue(string(firewallElement.Status)),
				TotalRuleCount:  types.Int32Value(*firewallElement.TotalRuleCount),
				FlavorName:      types.StringPointerValue(firewallElement.FlavorName),
				FlavorRuleQuota: types.Int32Value(*firewallElement.FlavorRuleQuota),
				CreatedAt:       types.StringValue(firewallElement.CreatedAt.Format(time.RFC3339)),
				CreatedBy:       types.StringValue(firewallElement.CreatedBy),
				ModifiedAt:      types.StringValue(firewallElement.ModifiedAt.Format(time.RFC3339)),
				ModifiedBy:      types.StringValue(firewallElement.ModifiedBy),
			}
			firewallObjectValue, _ := types.ObjectValueFrom(ctx, firewallModel.AttributeTypes(), firewallModel)
			state.Firewall = firewallObjectValue
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
