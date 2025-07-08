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
				Description: "Id",
				Required:    true,
			},
			common.ToSnakeCase("VpnGateway"): schema.SingleNestedAttribute{
				Description: "VPN Gateway",
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
					common.ToSnakeCase("IpAddress"): schema.StringAttribute{
						Description: "IpAddress",
						Computed:    true,
					},
					common.ToSnakeCase("IpId"): schema.StringAttribute{
						Description: "IpId",
						Computed:    true,
					},
					common.ToSnakeCase("IpType"): schema.StringAttribute{
						Description: "IpType",
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
