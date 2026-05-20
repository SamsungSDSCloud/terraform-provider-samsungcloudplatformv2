package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &VpcSubnetVipDataSource{}
	_ datasource.DataSourceWithConfigure = &VpcSubnetVipDataSource{}
)

// NewVpcSubnetVipDataSource is a helper function to simplify the provider implementation.
func NewVpcSubnetVipDataSource() datasource.DataSource {
	return &VpcSubnetVipDataSource{}
}

// VpcSubnetVipDataSource is the data source implementation.
type VpcSubnetVipDataSource struct {
	config  *scpsdk.Configuration
	client  *vpcv1d2.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *VpcSubnetVipDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_subnet_vip"
}

// Schema defines the schema for the data source.
// Schema defines the schema for the data source.
func (d *VpcSubnetVipDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get detail Subnet Vip",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "Subnet ID \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
			},
			common.ToSnakeCase("VipId"): schema.StringAttribute{
				Description: "Subnet Vip ID \n" +
					"  - example : 0466a9448d9a4411a86055939e451c8f",
				Required: true,
			},

			// Output

			common.ToSnakeCase("SubnetVip"): schema.SingleNestedAttribute{
				Description: "Subnet vip detail",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Subnet Vip Id",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created At",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created By",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified At",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State \n" +
							"  - enum : CREATING, ACTIVE, DELETING, DELETED, ERROR",
						Computed: true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "Subnet ID",
						Computed:    true,
					},
					common.ToSnakeCase("VipPortId"): schema.StringAttribute{
						Description: "Vip Port Id",
						Computed:    true,
					},
					common.ToSnakeCase("VirtualIpAddress"): schema.StringAttribute{
						Description: "Virtual IP Address \n" +
							"  - example : 192.168.20.6",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("ConnectedPorts"): schema.ListNestedAttribute{
						Description: "Connected Ports",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "Connected Port Id",
									Computed:    true,
								},
								common.ToSnakeCase("PortId"): schema.StringAttribute{
									Description: "Port ID",
									Computed:    true,
								},
								common.ToSnakeCase("PortName"): schema.StringAttribute{
									Description: "Port Name",
									Computed:    true,
								},
								common.ToSnakeCase("PortIpAddress"): schema.StringAttribute{
									Description: "Port IP Address",
									Computed:    true,
								},
								common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
									Description: "Connected resource ID",
									Computed:    true,
								},
								common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
									Description: "Connected resource Type",
									Computed:    true,
								},
							},
						},
					},
					common.ToSnakeCase("StaticNat"): schema.SingleNestedAttribute{
						Description: "Static NAT Info",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ExternalIpAddress"): schema.StringAttribute{
								Description: "Static Nat External Ip Address",
								Computed:    true,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "Static Nat Id",
								Computed:    true,
							},
							common.ToSnakeCase("PublicipId"): schema.StringAttribute{
								Description: "Publicip ID",
								Computed:    true,
							},
							common.ToSnakeCase("State"): schema.StringAttribute{
								Description: "Static Nat State",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *VpcSubnetVipDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.VpcV1Dot2
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *VpcSubnetVipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpcv1d2.SubnetVipDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.ShowSubnetVip(ctx, state.SubnetId.ValueString(), state.VipId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading subnet",
			"Could not read subnet, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	// Map API response to state
	state.SubnetVip = &vpcv1d2.VpcSubnetVipDetail{
		Id:               types.StringValue(data.SubnetVip.Id),
		CreatedAt:        types.StringValue(data.SubnetVip.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.SubnetVip.CreatedBy),
		ModifiedAt:       types.StringValue(data.SubnetVip.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.SubnetVip.ModifiedBy),
		State:            types.StringValue((string)(*data.SubnetVip.State.Ptr())),
		SubnetId:         types.StringValue(data.SubnetVip.SubnetId),
		VipPortId:        types.StringValue(data.SubnetVip.VipPortId),
		VirtualIpAddress: types.StringValue(data.SubnetVip.VirtualIpAddress),
	}
	if data.SubnetVip.Description.IsSet() {
		if desc := data.SubnetVip.Description.Get(); desc != nil {
			state.SubnetVip.Description = types.StringValue(*desc)
		}
	}

	state.SubnetVip.ConnectedPorts = []vpcv1d2.ConnectedPortInfo{}
	if data.SubnetVip.ConnectedPorts != nil {
		for _, port := range data.SubnetVip.ConnectedPorts {
			state.SubnetVip.ConnectedPorts = append(state.SubnetVip.ConnectedPorts, vpcv1d2.ConnectedPortInfo{
				Id:                   types.StringValue(port.Id),
				PortId:               types.StringValue(port.PortId),
				PortName:             types.StringValue(port.PortName),
				PortIpAddress:        types.StringValue(port.PortIpAddress),
				AttachedResourceId:   types.StringValue(port.AttachedResourceId),
				AttachedResourceType: types.StringValue(port.AttachedResourceType),
			})
		}
	}

	if data.SubnetVip.StaticNat.IsSet() {
		resultStaticNat := data.SubnetVip.StaticNat.Get()
		if resultStaticNat != nil {
			state.SubnetVip.StaticNat = &vpcv1d2.StaticNatSummary{
				ExternalIpAddress: types.StringValue(resultStaticNat.ExternalIpAddress),
				Id:                types.StringValue(resultStaticNat.Id),
				PublicipId:        types.StringValue(resultStaticNat.PublicipId),
				State:             types.StringValue(resultStaticNat.State),
			}
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
