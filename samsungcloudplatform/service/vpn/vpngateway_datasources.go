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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vpnVpnGatewayDataSources{}
	_ datasource.DataSourceWithConfigure = &vpnVpnGatewayDataSources{}
)

// NewVpnVpnGatewayDataSources is a helper function to simplify the provider implementation.
func NewVpnVpnGatewayDataSources() datasource.DataSource {
	return &vpnVpnGatewayDataSources{}
}

// vpnVpnGatewayDataSources is the data source implementation.
type vpnVpnGatewayDataSources struct {
	config  *scpsdk.Configuration
	client  *vpn.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpnVpnGatewayDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_vpn_gateways"
}

func (d *vpnVpnGatewayDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of vpn gateway",
		Attributes: map[string]schema.Attribute{
		common.ToSnakeCase("Size"): schema.Int32Attribute{
			Description: "The number of items per page.\n" +
				"  - example: 20\n" +
				"  - constraints: min: 1",
			Optional:    true,
		},
		common.ToSnakeCase("Sort"): schema.StringAttribute{
			Description: "The sorting criteria.\n" +
				"  - example: created_at:desc\n" +
				"  - valid: field_name:asc or field_name:desc.",
			Optional:    true,
		},
		common.ToSnakeCase("Page"): schema.Int32Attribute{
			Description: "The page number for pagination.\n" +
				"  - example: 1\n" +
				"  - constraints: min: 1",
			Optional:    true,
		},
		common.ToSnakeCase("Name"): schema.StringAttribute{
			Description: "The name of the resource.\n" +
				"  - example: vpnGWProd\n" +
				"  - valid: English letters and numbers only\n" +
				"  - constraints: minLength: 1, maxLength: 20",
			Optional:    true,
		},
		common.ToSnakeCase("IpAddress"): schema.StringAttribute{
			Description: "The IP address assigned to the resource.\n" +
				"  - example: 10.0.0.0/24",
			Optional:    true,
		},
		common.ToSnakeCase("VpcId"): schema.StringAttribute{
			Description: "The identifier of the VPC that the resource belongs to.\n" +
				"  - example: f32265726b694b32920aa3b111f4c715",
			Optional:    true,
		},
		common.ToSnakeCase("VpcName"): schema.StringAttribute{
			Description: "The name of the VPC that the resource belongs to.\n" +
				"  - example: vpcProd01",
			Optional:    true,
		},
		common.ToSnakeCase("Ids"): schema.ListAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Vpn gateway Id List.\n" +
				"  - example: [8e83f42d823941d7a4883f0f99101ef9]",
		},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpnVpnGatewayDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vpnVpnGatewayDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpn.VpnGatewayDataSourceIds

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetVpnGatewayList(d.clients, state.Size, state.Page, state.Sort, state.Name, state.IpAddress, state.VpcName, state.VpcId) // client 를 호출한다.
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to Read Vpn Gateways. Error: %s, Config: %+v", err.Error(), state)
		resp.Diagnostics.AddError(
			"VPN Gateway Read Error",
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

func GetVpnGatewayList(clients *client.SCPClient, size types.Int32, page types.Int32, sort types.String, name types.String,
	ipAddress types.String, vpcName types.String, vpcId types.String) ([]types.String, error) {

	ctx := context.Background()

	data, err := clients.Vpn.GetVpnGatewayList(ctx, size, page, sort, name, ipAddress, vpcName, vpcId)
	if err != nil {
		return nil, err
	}

	contents := data.VpnGateways

	var ids []types.String

	// Map response body to model
	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
