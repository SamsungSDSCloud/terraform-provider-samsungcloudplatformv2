package vpn

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
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
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
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
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("Phase1"): schema.SingleNestedAttribute{
						Description: "Phase1",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"diffie_hellman_groups": schema.ListAttribute{
								Description: "DiffieHellmanGroups",
								Computed:    true,
								ElementType: types.Int32Type,
							},
							"dpd_retry_interval": schema.Int32Attribute{
								Description: "DpdRetryInterval",
								Computed:    true,
							},
							"encryptions": schema.ListAttribute{
								Description: "Encryptions",
								Computed:    true,
								ElementType: types.StringType,
							},
							"ike_version": schema.Int32Attribute{
								Description: "IkeVersion",
								Computed:    true,
							},
							"life_time": schema.Int32Attribute{
								Description: "LifeTime",
								Computed:    true,
							},
							"peer_gateway_ip": schema.StringAttribute{
								Description: "PeerGatewayIp",
								Computed:    true,
							},
							"pre_shared_key": schema.StringAttribute{
								Description: "PreSharedKey",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("Phase2"): schema.SingleNestedAttribute{
						Description: "Phase2",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"diffie_hellman_groups": schema.ListAttribute{
								Description: "DiffieHellmanGroups",
								Computed:    true,
								ElementType: types.Int32Type,
							},
							"encryptions": schema.ListAttribute{
								Description: "Encryptions",
								Computed:    true,
								ElementType: types.StringType,
							},
							"life_time": schema.Int32Attribute{
								Description: "LifeTime",
								Computed:    true,
							},
							"perfect_forward_secrecy": schema.StringAttribute{
								Description: "PerfectForwardSecrecy",
								Computed:    true,
							},
							"remote_subnet": schema.StringAttribute{
								Description: "RemoteSubnet",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId",
						Computed:    true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "VpcName",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayId"): schema.StringAttribute{
						Description: "VpnGatewayId",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayIpAddress"): schema.StringAttribute{
						Description: "VpnGatewayIpAddress",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayName"): schema.StringAttribute{
						Description: "VpnGatewayName",
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
	var defaultPeerGatewayIp types.String
	var defaultRemoteSubnet types.String

	ids, err := GetVpnTunnelList(d.clients, defaultPage, defaultSize, defaultSort, defaultName,
		defaultVpnGatewayId, defaultVpnGatewayName, defaultPeerGatewayIp, defaultRemoteSubnet)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to Read Vpn Tunnels. Error: %s, Config: %+v", err.Error(), state)
		resp.Diagnostics.AddError(
			"VPN Tunnel Read Error",
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
			data, err := d.client.GetVpnTunnel(ctx, state.Id.ValueString()) // client 를 호출한다.
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Reading Vpn Tunnel",
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
				VpcId:               types.StringValue(vpnTunnelElement.VpcId),
				VpcName:             types.StringValue(vpnTunnelElement.VpcName),
				VpnGatewayId:        types.StringValue(vpnTunnelElement.VpnGatewayId),
				VpnGatewayIpAddress: types.StringValue(vpnTunnelElement.VpnGatewayIpAddress),
				VpnGatewayName:      types.StringValue(vpnTunnelElement.VpnGatewayName),
			}
			vpnTunnelObjectValue, _ := types.ObjectValueFrom(ctx, vpnTunnelModel.AttributeTypes(), vpnTunnelModel)
			state.VpnTunnel = vpnTunnelObjectValue
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
