package vpn

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vpnVpnGatewayDataSource{}
	_ datasource.DataSourceWithConfigure = &vpnVpnGatewayDataSource{}
)

// NewVpnVpnGatewayDataSource is a helper function to simplify the provider implementation.
func NewVpnVpnGatewayDataSource() datasource.DataSource {
	return &vpnVpnGatewayDataSource{}
}

// vpnVpnGatewayDataSource is the data source implementation.
type vpnVpnGatewayDataSource struct {
	config  *scpsdk.Configuration
	client  *vpn.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpnVpnGatewayDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_vpn_gateway"
}

func (d *vpnVpnGatewayDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPN Gateway",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 6a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d",
				Required: true,
			},
			common.ToSnakeCase("VpnGateway"): schema.SingleNestedAttribute{
				Description: "The identifier of the VPN gateway that the resource belongs to.\n" +
					"  - example: 01c543eb4b8d42a9a3502345d4025147",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example: 297615908b8e4ec69520a99a6777add3",
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
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.\n" +
							"  - example: VPN test",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the resource.\n" +
							"  - example: 6a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("IpAddress"): schema.StringAttribute{
						Description: "The IP address assigned to the resource.\n" +
							"  - example: 10.0.0.0/24",
						Computed: true,
					},
					common.ToSnakeCase("IpId"): schema.StringAttribute{
						Description: "The identifier of the IP address assigned to the resource.\n" +
							"  - example: bd07e102fe574edf8a1748957c45bdbf",
						Computed: true,
					},
					common.ToSnakeCase("IpType"): schema.StringAttribute{
						Description: "The type of the IP address assigned to the resource.\n" +
							"  - example: PUBLIC",
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
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the resource.\n" +
							"  - example: vpnGWProd",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the resource.\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "The identifier of the VPC that the resource belongs to.\n" +
							"  - example: f32265726b694b32920aa3b111f4c715",
						Computed: true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "The name of the VPC that the resource belongs to.\n" +
							"  - example: vpcProd01",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpnVpnGatewayDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Vpn
	d.clients = inst.Client
}

func (d *vpnVpnGatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpn.VpnGatewayDataSource

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var defaultSize types.Int32
	var defaultPage types.Int32
	var defaultSort types.String
	var defaultName types.String
	var defaultIpAddress types.String
	var defaultVpcName types.String
	var defaultVpcId types.String

	ids, err := GetVpnGatewayList(d.clients, defaultSize, defaultPage, defaultSort, defaultName, defaultIpAddress, defaultVpcName, defaultVpcId) // client 를 호출한다.
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to Read Vpn Gateway. Error: %s, Config: %+v", err.Error(), state)
		resp.Diagnostics.AddError(
			"VPN Gateway Read Error",
			errorMessage,
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
			data, err := d.client.GetVpnGateway(ctx, state.Id.ValueString()) // client 를 호출한다.
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Reading Vpn Gateway",
					"Could not read Vpn Gateway ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
				)
				return
			}

			vpnGatewayElement := data.VpnGateway

			vpnGatewayModel := vpn.VpnGateway{
				AccountId:   types.StringValue(vpnGatewayElement.AccountId),
				CreatedAt:   types.StringValue(vpnGatewayElement.CreatedAt.Format(time.RFC3339)),
				CreatedBy:   types.StringValue(vpnGatewayElement.CreatedBy),
				Description: types.StringPointerValue(vpnGatewayElement.Description.Get()),
				Id:          types.StringValue(vpnGatewayElement.Id),
				IpAddress:   types.StringValue(vpnGatewayElement.IpAddress),
				IpId:        types.StringValue(vpnGatewayElement.IpId),
				IpType:      types.StringValue(vpnGatewayElement.IpType),
				ModifiedAt:  types.StringValue(vpnGatewayElement.ModifiedAt.Format(time.RFC3339)),
				ModifiedBy:  types.StringValue(vpnGatewayElement.ModifiedBy),
				Name:        types.StringValue(vpnGatewayElement.Name),
				State:       types.StringValue(string(vpnGatewayElement.State)),
				VpcId:       types.StringValue(vpnGatewayElement.VpcId),
				VpcName:     types.StringValue(vpnGatewayElement.VpcName),
			}
			vpnGatewayObjectValue, _ := types.ObjectValueFrom(ctx, vpnGatewayModel.AttributeTypes(), vpnGatewayModel)
			state.VpnGateway = vpnGatewayObjectValue
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
