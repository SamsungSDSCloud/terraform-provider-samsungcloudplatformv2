package firewall

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	firewall "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/firewallv1d1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &firewallFirewallDataSources{}
	_ datasource.DataSourceWithConfigure = &firewallFirewallDataSources{}
)

func NewFirewallFirewallDataSources() datasource.DataSource {
	return &firewallFirewallDataSources{}
}

type firewallFirewallDataSources struct {
	config  *scpsdk.Configuration
	client  *firewall.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *firewallFirewallDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_firewalls"
}

// Schema defines the schema for the data source.

func (d *firewallFirewallDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of firewall",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination.\n" +
					"  - example: 1\n" +
					"  - constraints: min: 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example: 20\n" +
					"  - constraints: min: 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria.\n" +
					"  - example: created_at:desc\n" +
					"  - valid: field_name:asc or field_name:desc",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the resource.\n" +
					"  - example: fw-web-prod\n" +
					"  - constraints: maxLength: 255",
				Optional: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "The name of the VPC that the resource belongs to.\n" +
					"  - example: vpc-prod-01",
				Optional: true,
			},
			common.ToSnakeCase("ProductType"): schema.ListAttribute{
				Description: "The type of the firewall service.\n" +
					"  - example: IGW\n" +
					"  - valid: IGW, GGW, DGW, LB, SIGW, TGW_IGW, TGW_GGW, TGW_DGW, TGW_SIGW, TGW_BM",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("IGW", "GGW", "DGW", "LB", "SIGW", "TGW_IGW", "TGW_GGW", "TGW_DGW", "TGW_SIGW", "TGW_BM"),
					),
				},
			},
			common.ToSnakeCase("State"): schema.ListAttribute{
				Description: "The current state of the resource.\n" +
					"  - example: ACTIVE\n" +
					"  - valid: CREATING, ACTIVE, EDITING, DELETING, ERROR, DEPLOYING",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("CREATING", "ACTIVE", "EDITING", "DELETING", "ERROR", "DEPLOYING"),
					),
				},
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "The total number of Firewall resources.\n" +
					"  - example : 2",
				Computed: true,
			},
			common.ToSnakeCase("Firewalls"): schema.ListNestedAttribute{
				Description: "A list of firewalls.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the firewall.\n" +
								"  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
							Computed: true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The account ID that owns the firewall.\n" +
								"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the firewall.\n" +
								"  - example: fw-web-prod",
							Computed: true,
						},
						common.ToSnakeCase("VpcId"): schema.StringAttribute{
							Description: "The VPC ID that the firewall belongs to.\n" +
								"  - example: vpc-1234567890abcdef",
							Computed: true,
						},
						common.ToSnakeCase("VpcName"): schema.StringAttribute{
							Description: "The VPC name that the firewall belongs to.\n" +
								"  - example: vpc-prod-01",
							Computed: true,
						},
						common.ToSnakeCase("Loggable"): schema.BoolAttribute{
							Description: "Whether logging is enabled for the firewall.\n" +
								"  - example: true",
							Computed: true,
						},
						common.ToSnakeCase("PreProductId"): schema.StringAttribute{
							Description: "The pre-product ID associated with the firewall.\n" +
								"  - example: igw-1234567890abcdef",
							Computed: true,
						},
						common.ToSnakeCase("ProductType"): schema.StringAttribute{
							Description: "The type of the firewall service.\n" +
								"  - example: IGW\n" +
								"  - valid: IGW, GGW, DGW, LB, SIGW, TGW_IGW, TGW_GGW, TGW_DGW, TGW_SIGW, TGW_BM",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current state of the firewall.\n" +
								"  - example: ACTIVE\n" +
								"  - valid: CREATING, ACTIVE, EDITING, DELETING, ERROR, DEPLOYING",
							Computed: true,
						},
						common.ToSnakeCase("Status"): schema.StringAttribute{
							Description: "The status of the firewall.\n" +
								"  - example: ENABLED",
							Computed: true,
						},
						common.ToSnakeCase("TotalRuleCount"): schema.Int32Attribute{
							Description: "The total number of rules in the firewall.\n" +
								"  - example: 10",
							Computed: true,
						},
						common.ToSnakeCase("FlavorName"): schema.StringAttribute{
							Description: "The flavor name of the firewall.\n" +
								"  - example: FIREWALL.STANDARD",
							Computed: true,
						},
						common.ToSnakeCase("FlavorRuleQuota"): schema.Int32Attribute{
							Description: "The maximum number of rules allowed by the firewall flavor.\n" +
								"  - example: 100",
							Computed: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the firewall was created, in ISO 8601 format.\n" +
								"  - example: 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user ID that created the firewall.\n" +
								"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the firewall was last modified, in ISO 8601 format.\n" +
								"  - example: 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user ID that last modified the firewall.\n" +
								"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("ZoneResources"): schema.ListNestedAttribute{
							Description: "A list of zone resources associated with the firewall.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("AllocateState"): schema.StringAttribute{
										Description: "The allocation state of the zone resource.\n" +
											"  - example: ALLOCATED",
										Computed: true,
									},
									common.ToSnakeCase("FwResourceId"): schema.StringAttribute{
										Description: "The resource ID of the firewall in the zone.\n" +
											"  - example: fw-zone-1234567890abcdef",
										Computed: true,
									},
									common.ToSnakeCase("Zone"): schema.StringAttribute{
										Description: "The zone name where the firewall is deployed.\n" +
											"  - example: zone-1",
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *firewallFirewallDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.FirewallV1d1
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *firewallFirewallDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state firewall.FirewallDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataResp, err := d.client.GetFirewallList(state.Page, state.Size, state.Sort, state.Name, state.VpcName, state.ProductType, state.State)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read firewalls.",
			err.Error(),
		)
		return
	}

	for _, fwRes := range dataResp.Firewalls {
		var zoneResources []firewall.ZoneResources
		if len(fwRes.ZoneResources) > 0 {
			for _, v := range fwRes.ZoneResources {
				zoneResource := firewall.ZoneResources{
					AllocateState: types.StringValue(string(v.AllocateState)),
					FwResourceId:  types.StringValue(v.FwResourceId),
					Zone:          types.StringValue(v.Zone),
				}
				zoneResources = append(zoneResources, zoneResource)
			}
		}

		firewallModel := firewall.Firewall{
			Id:              types.StringValue(fwRes.Id),
			AccountId:       types.StringValue(fwRes.AccountId),
			Name:            types.StringValue(fwRes.Name),
			VpcId:           types.StringPointerValue(fwRes.VpcId.Get()),
			VpcName:         types.StringPointerValue(fwRes.VpcName.Get()),
			Loggable:        types.BoolValue(fwRes.Loggable),
			PreProductId:    types.StringPointerValue(fwRes.PreProductId),
			ProductType:     types.StringValue(string(fwRes.ProductType)),
			State:           types.StringValue(string(fwRes.State)),
			Status:          types.StringValue(string(fwRes.Status)),
			TotalRuleCount:  types.Int32PointerValue(fwRes.TotalRuleCount),
			FlavorName:      types.StringPointerValue(fwRes.FlavorName),
			FlavorRuleQuota: types.Int32PointerValue(fwRes.FlavorRuleQuota),
			CreatedAt:       types.StringValue(fwRes.CreatedAt.Format(time.RFC3339)),
			CreatedBy:       types.StringValue(fwRes.CreatedBy),
			ModifiedAt:      types.StringValue(fwRes.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:      types.StringValue(fwRes.ModifiedBy),
			ZoneResources:   zoneResources,
		}

		state.Firewalls = append(state.Firewalls, firewallModel)
	}

	state.TotalCount = types.Int32Value(dataResp.Count)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
