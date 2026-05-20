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
	_ datasource.DataSource              = &vpcNatGatewayDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcNatGatewayDataSource{}
)

// NewVpcNatGatewayDataSource is a helper function to simplify the provider implementation.
func NewVpcNatGatewayDataSource() datasource.DataSource {
	return &vpcNatGatewayDataSource{}
}

// vpcNatGatewayDataSource is the data source implementation.
type vpcNatGatewayDataSource struct {
	config  *scpsdk.Configuration
	client  *vpcv1d2.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcNatGatewayDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_nat_gateways"
}

// Schema defines the schema for the data source.
func (d *vpcNatGatewayDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of NAT Gateways.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size \n" +
					"  - example : 20",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page \n" +
					"  - example : 0",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "NAT Gateway Name \n" +
					"  - example : NatGatewayName",
				Optional: true,
			},
			common.ToSnakeCase("NatGatewayIpAddress"): schema.StringAttribute{
				Description: "NAT Gateway IP Address \n" +
					"  - example : 192.167.0.5",
				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "VPC ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "VPC Name \n" +
					"  - example : vpcName",
				Optional: true,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "Subnet ID \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Optional: true,
			},
			common.ToSnakeCase("SubnetName"): schema.StringAttribute{
				Description: "Subnet Name \n" +
					"  - example : subnetName",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "NAT Gateway State \n" +
					"  - example : CREATING | ACTIVE | DELETING | DELETED | ERROR",
				Optional: true,
			},
			common.ToSnakeCase("NatGateways"): schema.ListNestedAttribute{
				Description: "A list of NAT Gateways.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "NAT Gateway ID",
							Computed:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "NAT Gateway Name",
							Computed:    true,
						},
						common.ToSnakeCase("NatGatewayIpAddress"): schema.StringAttribute{
							Description: "NAT Gateway IP Address",
							Computed:    true,
						},
						common.ToSnakeCase("VpcId"): schema.StringAttribute{
							Description: "VPC ID",
							Computed:    true,
						},
						common.ToSnakeCase("VpcName"): schema.StringAttribute{
							Description: "VPC Name",
							Computed:    true,
						},
						common.ToSnakeCase("SubnetId"): schema.StringAttribute{
							Description: "Subnet ID",
							Computed:    true,
						},
						common.ToSnakeCase("SubnetName"): schema.StringAttribute{
							Description: "Subnet Name",
							Computed:    true,
						},
						common.ToSnakeCase("SubnetCidr"): schema.StringAttribute{
							Description: "Subnet CIDR",
							Computed:    true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "Account ID",
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "NAT Gateway State",
							Computed:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "NAT Gateway Description",
							Computed:    true,
						},
						common.ToSnakeCase("PublicipId"): schema.StringAttribute{
							Description: "PublicIP ID",
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
					},
				},
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Count",
				Computed:    true,
			},
			common.ToSnakeCase("SortFinal"): schema.ListAttribute{
				Description: "List of sort condition \n" +
					"  - example : [\"created_at:desc\"]",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcNatGatewayDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vpcNatGatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpcv1d2.NatGatewayDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.ListNatGateways(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading NAT gateway",
			"Could not read NAT gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	for _, natgateway := range data.NatGateways {
		natgatewayState := vpcv1d2.NatGateway{
			Id:                  types.StringValue(natgateway.Id),
			Name:                types.StringValue(natgateway.Name),
			NatGatewayIpAddress: types.StringValue(natgateway.NatGatewayIpAddress),
			VpcId:               types.StringValue(natgateway.VpcId),
			VpcName:             types.StringValue(natgateway.VpcName),
			SubnetId:            types.StringValue(natgateway.SubnetId),
			SubnetName:          types.StringValue(natgateway.SubnetName),
			SubnetCidr:          types.StringValue(natgateway.SubnetCidr),
			AccountId:           types.StringValue(natgateway.AccountId),
			State:               types.StringValue(string(natgateway.State)),
			Description:         types.StringPointerValue(natgateway.Description.Get()),
			PublicipId:          types.StringPointerValue(natgateway.PublicipId.Get()),
			CreatedAt:           types.StringValue(natgateway.CreatedAt.Format(time.RFC3339)),
			CreatedBy:           types.StringValue(natgateway.CreatedBy),
			ModifiedAt:          types.StringValue(natgateway.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:          types.StringValue(natgateway.ModifiedBy),
		}
		state.NatGateways = append(state.NatGateways, natgatewayState)
	}

	state.TotalCount = types.Int32Value(data.Count)
	state.Page = types.Int32Value(data.Page)
	state.Size = types.Int32Value(data.Size)
	for _, sortVal := range data.Sort {
		state.SortFinal = append(state.SortFinal, types.StringValue(sortVal))
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
