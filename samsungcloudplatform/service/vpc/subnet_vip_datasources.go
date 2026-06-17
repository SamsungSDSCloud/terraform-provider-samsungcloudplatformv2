package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &VpcSubnetVipDataSources{}
	_ datasource.DataSourceWithConfigure = &VpcSubnetVipDataSources{}
)

// NewVpcSubnetVipDataSource is a helper function to simplify the provider implementation.
func NewVpcSubnetVipDataSources() datasource.DataSource {
	return &VpcSubnetVipDataSources{}
}

// VpcSubnetVipDataSources is the data source implementation.
type VpcSubnetVipDataSources struct {
	config  *scpsdk.Configuration
	client  *vpcv1d2.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *VpcSubnetVipDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_subnet_vips"
}

// Schema defines the schema for the data source.
func (d *VpcSubnetVipDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A list of Subnet Vip",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "The identifier of the subnet that the subnet vip belongs to.\n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example : 20 \n" +
					"  - minimum : 0",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination.\n" +
					"  - example : 0 \n" +
					"  - minimum : 0",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for decending order.\n" +
					"  - example : created_at:desc",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("VirtualIpAddress"): schema.StringAttribute{
				Description: "The virtual IP address assigned to the resource.\n" +
					"  - example : 192.168.20.6",
				Optional: true,
			},
			common.ToSnakeCase("PublicIpAddress"): schema.StringAttribute{
				Description: "The ip address of the Public IP Address \n" +
					"  - example : 100.112.9.84",
				Optional: true,
			},

			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "The count of Subnet Vips.\n" +
					"  - example : 2",
				Computed: true,
			},
			common.ToSnakeCase("SubnetVips"): schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of Subnet Vips",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("ConnectedPortCount"): schema.Int32Attribute{
							Description: "The count of connected ports\n" +
								"  - example : 2",
							Computed: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the subnet vip was created in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user id that created the subnet vip.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the subnet vip.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the subnet vip was last modified in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user id that modified the subnet vip.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the subnet vip.\n" +
								"  - enum : CREATING, ACTIVE, DELETING, DELETED, ERROR\n" +
                                "  - example : ACTIVE",
							Computed: true,
						},
						common.ToSnakeCase("StaticNat"): schema.SingleNestedAttribute{
							Description: "Static NAT Info",
							Computed:    true,
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("ExternalIpAddress"): schema.StringAttribute{
									Description: "Static Nat External Ip Address\n" +
										"  - example : 192.168.0.1",
									Computed: true,
								},
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "The unique identifier of the static nat.\n" +
										"  - example : 7df8abb4912e4709b1cb237daccca7a8",
									Computed: true,
								},
								common.ToSnakeCase("PublicipId"): schema.StringAttribute{
									Description: "The identifier of the public IP address.\n" +
										"  - example : 7df8abb4912e4709b1cb237daccca7a8",
									Computed: true,
								},
								common.ToSnakeCase("State"): schema.StringAttribute{
									Description: "The current lifecycle state of the static nat.\n" +
										"  - example : ACTIVE",
									Computed: true,
								},
							},
						},
						common.ToSnakeCase("VirtualIpAddress"): schema.StringAttribute{
							Description: "The virtual IP address assigned to the subnet vip.\n" +
								"  - example : 192.168.0.1",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *VpcSubnetVipDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *VpcSubnetVipDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpcv1d2.SubnetVipDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.ListSubnetVips(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading subnet",
			"Could not read subnet, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	state.TotalCount = types.Int32Value(data.Count)
	state.Page = types.Int32Value(data.Page)
	state.Size = types.Int32Value(data.Size)
	if len(data.Sort) > 0 {
		state.Sort = types.StringValue(data.Sort[0])
	}

	state.SubnetVips = []vpcv1d2.SubnetVipSummary{}
	if data.SubnetVips != nil {
		for _, subnetVipRaw := range data.SubnetVips {
			subnetVip := vpcv1d2.SubnetVipSummary{
				Id:                 types.StringValue(subnetVipRaw.Id),
				CreatedAt:          types.StringValue(subnetVipRaw.CreatedAt.Format(time.RFC3339)),
				CreatedBy:          types.StringValue(subnetVipRaw.CreatedBy),
				ModifiedAt:         types.StringValue(subnetVipRaw.ModifiedAt.Format(time.RFC3339)),
				ModifiedBy:         types.StringValue(subnetVipRaw.ModifiedBy),
				State:              types.StringValue((string)(*subnetVipRaw.State.Ptr())),
				ConnectedPortCount: types.Int32PointerValue(subnetVipRaw.ConnectedPortCount),
				VirtualIpAddress:   types.StringValue(subnetVipRaw.VirtualIpAddress),
			}
			if subnetVipRaw.StaticNat.IsSet() {
				subnetVipRawStaticNat := subnetVipRaw.StaticNat.Get()
				if subnetVipRawStaticNat != nil {
					subnetVip.StaticNat = &vpcv1d2.StaticNatSummary{
						ExternalIpAddress: types.StringValue(subnetVipRawStaticNat.ExternalIpAddress),
						Id:                types.StringValue(subnetVipRawStaticNat.Id),
						PublicipId:        types.StringValue(subnetVipRawStaticNat.PublicipId),
						State:             types.StringValue(subnetVipRawStaticNat.State),
					}
				}
			}

			state.SubnetVips = append(state.SubnetVips, subnetVip)
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
