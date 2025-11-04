package vpn

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vpnVpnTunnelDataSources{}
	_ datasource.DataSourceWithConfigure = &vpnVpnTunnelDataSources{}
)

// NewVpnVpnTunnelDataSources is a helper function to simplify the provider implementation.
func NewVpnVpnTunnelDataSources() datasource.DataSource {
	return &vpnVpnTunnelDataSources{}
}

// vpnVpnTunnelDataSources is the data source implementation.
type vpnVpnTunnelDataSources struct {
	config  *scpsdk.Configuration
	client  *vpn.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpnVpnTunnelDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_vpn_tunnels"
}

func (d *vpnVpnTunnelDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of vpn tunnel",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size",
				Optional:    true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("VpnGatewayId"): schema.StringAttribute{
				Description: "VpnGatewayId",
				Optional:    true,
			},
			common.ToSnakeCase("VpnGatewayName"): schema.StringAttribute{
				Description: "VpnGatewayName",
				Optional:    true,
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Vpn Tunnel Id List",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpnVpnTunnelDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vpnVpnTunnelDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpn.VpnTunnelData1d1Source

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := getVpnTunnelList(d.clients, state.Page, state.Size, state.Sort, state.Name, state.VpnGatewayId, state.VpnGatewayName)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to Read Vpn Tunnels v1.1. Error: %s, Config: %+v", err.Error(), state)
		resp.Diagnostics.AddError(
			"VPN Tunnel Read Error",
			errorMessage,
		)
	}

	state.Ids = ids

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getVpnTunnelList(clients *client.SCPClient, page types.Int32, size types.Int32, sort types.String, name types.String,
	vpnGatewayId types.String, vpnGatewayName types.String) ([]types.String, error) {

	ctx := context.Background()

	data, err := clients.Vpn.GetVpnTunnelList(ctx, page, size, sort, name, vpnGatewayId, vpnGatewayName)
	if err != nil {
		return nil, err
	}

	contents := data.VpnTunnels

	var ids []types.String

	// Map response body to model
	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
