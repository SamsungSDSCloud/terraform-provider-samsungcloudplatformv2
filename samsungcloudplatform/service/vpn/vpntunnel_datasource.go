package vpn

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
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
				Description: "Id",
				Required:    true,
			},
			common.ToSnakeCase("VpnTunnel"): schema.SingleNestedAttribute{
				Description: "VPN Tunnel",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId \n  - example: 0e3dffc50eb247a1adf4f2e5c82c4f99 ",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt \n - example : 2024-05-17T00:23:17Z",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy \n - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description \n - example : Example Description for VPN Tunnel",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id \n - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt \n - example : 2024-05-17T00:23:17Z",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy \n - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name \n - example: ExampleVpnTunnel1 ",
						Computed:    true,
					},
					common.ToSnakeCase("Phase1"): schema.SingleNestedAttribute{
						Description: "Phase1",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"dpd_retry_interval": schema.Int32Attribute{
								Description: "DpdRetryInterval \n - example: 60",
								Computed:    true,
							},
							"ike_version": schema.Int32Attribute{
								Description: "IkeVersion \n - example: 2",
								Computed:    true,
							},
							"life_time": schema.Int32Attribute{
								Description: "LifeTime \n - example: 86400 ",
								Computed:    true,
							},
							"peer_gateway_ip": schema.StringAttribute{
								Description: "PeerGatewayIp \n - example: 123.0.0.2",
								Computed:    true,
							},
							"diffie_hellman_groups": schema.ListAttribute{
								Description: "VPN Tunnel ISAKMP Diffie-Hellman Group 목록 \n - example : [\n   \"30\",\n    \"31\",\n   \"32\"\n  ]",
								Computed:    true,
								ElementType: types.Int32Type,
							},
							"encryptions": schema.ListAttribute{
								Description: "VPN Tunnel ISAKMP Proposal 목록 \n - example : [\n   \"null-md5\",\n    \"aes128gcm\",\n   \"chacha20poly1305\"\n  ]",
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
					common.ToSnakeCase("Phase2"): schema.SingleNestedAttribute{
						Description: "Phase2",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"life_time": schema.Int32Attribute{
								Description: "LifeTime \n - example: 86400 ",
								Computed:    true,
							},
							"perfect_forward_secrecy": schema.StringAttribute{
								Description: "PerfectForwardSecrecy \n - example: ENABLE",
								Computed:    true,
							},
							"remote_subnets": schema.ListAttribute{
								Description: "VPN Tunnel IPSec Remote Subnets \n - example : [\n   \"10.1.1.0/24\",\n    \"10.1.2.0/24\",\n   \"10.1.3.0/24\"\n  ]",
								Computed:    true,
								ElementType: types.StringType,
							},
							"diffie_hellman_groups": schema.ListAttribute{
								Description: "VPN Tunnel ISAKMP Diffie-Hellman Group 목록 \n - example : [\n   \"30\",\n    \"31\",\n   \"32\"\n  ]",
								Computed:    true,
								ElementType: types.Int32Type,
							},
							"encryptions": schema.ListAttribute{
								Description: "VPN Tunnel ISAKMP Proposal 목록 \n - example : [\n   \"null-md5\",\n    \"aes128gcm\",\n   \"chacha20poly1305\"\n  ]",
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State \n - example: ACTIVE",
						Computed:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status \n - example : DOWN",
						Computed:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId\n - example: ceb44ea5ecb34a49b16495f9a63b0718",
						Computed:    true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "VpcName \n - example: ExampleVPC1",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayId"): schema.StringAttribute{
						Description: "VpnGatewayId \n- example: b156740b6335468d8354eb9ef8eddf5a",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayIpAddress"): schema.StringAttribute{
						Description: "VpnGatewayIpAddress \n - example: 123.0.0.1",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayName"): schema.StringAttribute{
						Description: "VpnGatewayName \n - example: ExampleVpnGW1",
						Computed:    true,
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
