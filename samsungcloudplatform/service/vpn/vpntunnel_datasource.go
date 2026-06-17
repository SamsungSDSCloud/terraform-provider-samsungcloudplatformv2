package vpn

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vpnVpnTunnelDataSource{}
	_ datasource.DataSourceWithConfigure = &vpnVpnTunnelDataSource{}
)

// NewVpnVpnTunnelDataSource is a helper function to simplify the provider implementation.
func NewVpnVpnTunnelDataSource() datasource.DataSource {
	return &vpnVpnTunnelDataSource{}
}

// vpnVpnTunnelDataSource is the data source implementation.
type vpnVpnTunnelDataSource struct {
	config  *scpsdk.Configuration
	client  *vpn.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpnVpnTunnelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_vpn_tunnel"
}

func (d *vpnVpnTunnelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPN Tunnel",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 6a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d",
				Required: true,
			},
			common.ToSnakeCase("VpnTunnel"): schema.SingleNestedAttribute{
				Description: "VPN Tunnel",
				Computed:    true,
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
					common.ToSnakeCase("Phase1"): schema.SingleNestedAttribute{
						Description: "The IKE phase 1 negotiation settings of the VPN tunnel.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"dpd_retry_interval": schema.Int32Attribute{
								Description: "The Dead Peer Detection retry interval in seconds.\n" +
									"  - example: 60",
								Computed: true,
							},
							"ike_version": schema.Int32Attribute{
								Description: "The IKE (Internet Key Exchange) protocol version.\n" +
									"  - example: 2",
								Computed: true,
							},
							"life_time": schema.Int32Attribute{
								Description: "The lifetime of the IKE phase 1 security association in seconds.\n" +
									"  - example: 86400",
								Computed: true,
							},
							"peer_gateway_ip": schema.StringAttribute{
								Description: "The IP address of the peer VPN gateway.\n" +
									"  - example: 123.0.0.2",
								Computed: true,
							},
							"diffie_hellman_groups": schema.ListAttribute{
								Description: "The list of Diffie-Hellman groups for IKE phase 1.\n" +
									"  - example: [30, 31, 32]",
								Computed:    true,
								ElementType: types.Int32Type,
							},
							"encryptions": schema.ListAttribute{
								Description: "The list of encryption algorithms for IKE phase 1.\n" +
									"  - example: [\"null-md5\", \"aes128gcm\", \"chacha20poly1305\"]",
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
					common.ToSnakeCase("Phase2"): schema.SingleNestedAttribute{
						Description: "The IKE phase 2 negotiation settings of the VPN tunnel.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"life_time": schema.Int32Attribute{
								Description: "The lifetime of the IKE phase 2 security association in seconds.\n" +
									"  - example: 86400",
								Computed: true,
							},
							"perfect_forward_secrecy": schema.StringAttribute{
								Description: "The Perfect Forward Secrecy setting for IKE phase 2.\n" +
									"  - example: ENABLE",
								Computed: true,
							},
							"remote_subnets": schema.ListAttribute{
								Description: "The list of remote subnets for IKE phase 2.\n" +
									"  - example: [\"10.1.1.0/24\", \"10.1.2.0/24\", \"10.1.3.0/24\"]",
								Computed:    true,
								ElementType: types.StringType,
							},
							"diffie_hellman_groups": schema.ListAttribute{
								Description: "The list of Diffie-Hellman groups for IKE phase 2.\n" +
									"  - example: [30, 31, 32]",
								Computed:    true,
								ElementType: types.Int32Type,
							},
							"encryptions": schema.ListAttribute{
								Description: "The list of encryption algorithms for IKE phase 2.\n" +
									"  - example: [\"null-md5\", \"aes128gcm\", \"chacha20poly1305\"]",
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the resource.\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "The current status of the vpn tunnel.\n" +
							"  - example: UP",
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
					common.ToSnakeCase("VpnGatewayId"): schema.StringAttribute{
						Description: "The identifier of the VPN gateway that the resource belongs to.\n" +
							"  - example: 01c543eb4b8d42a9a3502345d4025147",
						Computed: true,
					},
					common.ToSnakeCase("VpnGatewayIpAddress"): schema.StringAttribute{
						Description: "The IP address of the VPN gateway.\n" +
							"  - example: 10.0.0.0/24",
						Computed: true,
					},
					common.ToSnakeCase("VpnGatewayName"): schema.StringAttribute{
						Description: "The name of the VPN gateway that the resource belongs to.\n" +
							"  - example: vpnGWProd",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpnVpnTunnelDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vpnVpnTunnelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpn.VpnTunnelDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var defaultSize types.Int32
	var defaultPage types.Int32
	var defaultSort types.String
	var defaultName types.String
	var defaultVpnGatewayId types.String
	var defaultVpnGatewayName types.String

	ids, err := getVpnTunnelList(d.clients, defaultPage, defaultSize, defaultSort, defaultName, defaultVpnGatewayId, defaultVpnGatewayName)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to Read Vpn Tunnels. Error: %s, Config: %+v", err.Error(), state)
		resp.Diagnostics.AddError(
			"VPN Tunnel v1.1 Read Error",
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

		fmt.Println("get state", state)
		fmt.Println("GetVpnTunnel.getList", ids)

		if exist {
			data, err := d.client.GetVpnTunnel(ctx, state.Id.ValueString()) // client 를 호출한다.
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Reading Vpn Tunnel v1.1",
					"Could not read Vpn Tunnel ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
				)
				return
			}

			vpnTunnelElement := data.VpnTunnel
			vpnTunnelModel := vpn.VpnTunnel{
				AccountId:           types.StringValue(vpnTunnelElement.AccountId),
				CreatedAt:           types.StringValue(vpnTunnelElement.CreatedAt.Format(time.RFC3339)),
				CreatedBy:           types.StringValue(vpnTunnelElement.CreatedBy),
				Description:         types.StringPointerValue(vpnTunnelElement.Description.Get()),
				Id:                  types.StringValue(vpnTunnelElement.Id),
				ModifiedAt:          types.StringValue(vpnTunnelElement.ModifiedAt.Format(time.RFC3339)),
				ModifiedBy:          types.StringValue(vpnTunnelElement.ModifiedBy),
				Name:                types.StringValue(vpnTunnelElement.Name),
				Phase1:              mapPhase1Detail(vpnTunnelElement.Phase1),
				Phase2:              mapPhase2Detail(vpnTunnelElement.Phase2),
				State:               types.StringValue(string(vpnTunnelElement.State)),
				Status:              types.StringValue(string(vpnTunnelElement.Status)),
				VpcId:               types.StringValue(vpnTunnelElement.VpcId),
				VpcName:             types.StringValue(vpnTunnelElement.VpcName),
				VpnGatewayId:        types.StringValue(vpnTunnelElement.VpnGatewayId),
				VpnGatewayIpAddress: types.StringValue(vpnTunnelElement.VpnGatewayIpAddress),
				VpnGatewayName:      types.StringValue(vpnTunnelElement.VpnGatewayName),
			}

			vpnTunnelObjectValue, _ := types.ObjectValueFrom(ctx, vpnTunnelModel.AttributeTypes(), vpnTunnelModel)

			state.VpnTunnel = vpnTunnelObjectValue
			diags = resp.State.Set(ctx, &state)

		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
