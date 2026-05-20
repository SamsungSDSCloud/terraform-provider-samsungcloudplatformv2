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
	_ datasource.DataSource              = &vpcInternetGatewayDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcInternetGatewayDataSource{}
)

// NewVpcInternetGatewayDataSource is a helper function to simplify the provider implementation.
func NewVpcInternetGatewayDataSource() datasource.DataSource {
	return &vpcInternetGatewayDataSource{}
}

// vpcInternetGatewayDataSource is the data source implementation.
type vpcInternetGatewayDataSource struct {
	config  *scpsdk.Configuration
	client  *vpcv1d2.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcInternetGatewayDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_internet_gateways"
}

// Schema defines the schema for the data source.
func (d *vpcInternetGatewayDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of internet gateways.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size \n" +
					"  - example : 20 \n" +
					"  - minimum : 0",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page \n" +
					"  - example : 0 \n" +
					"  - minimum : 0",
				Optional: true,
				Computed: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Internet Gateway ID \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Internet Gateway Name \n" +
					"  - example : IGW_VPCname",
				Optional: true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Internet Gateway Type \n" +
					"  - example : IGW | GGW | SIGW",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State \n" +
					"  - example : CREATING | ACTIVE | DELETING | ERROR",
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
			common.ToSnakeCase("InternetGateways"): schema.ListNestedAttribute{
				Description: "A list of internet gateways.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Internet Gateway ID",
							Computed:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Internet Gateway Name",
							Computed:    true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "Account ID",
							Computed:    true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "Internet Gateway Type",
							Computed:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
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
						common.ToSnakeCase("Loggable"): schema.BoolAttribute{
							Description: "NAT Loggable",
							Computed:    true,
						},
						common.ToSnakeCase("FirewallId"): schema.StringAttribute{
							Description: "Firewall ID",
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
							Description: "State",
							Computed:    true,
						},
					},
				},
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Total count",
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
func (d *vpcInternetGatewayDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vpcInternetGatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpcv1d2.InternetGatewayDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.ListInternetGateways(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading internet gateway",
			"Could not read internet gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	for _, igw := range data.InternetGateways {
		igwState := vpcv1d2.InternetGateway{
			Id:          types.StringValue(igw.Id),
			Name:        types.StringValue(igw.Name),
			AccountId:   types.StringValue(igw.AccountId),
			Type:        types.StringValue(string(igw.Type)),
			Description: types.StringPointerValue(igw.Description.Get()),
			VpcId:       types.StringValue(igw.VpcId),
			VpcName:     types.StringValue(igw.VpcName),
			Loggable:    types.BoolValue(igw.GetLoggable()),
			FirewallId:  types.StringValue(igw.GetFirewallId()),
			CreatedAt:   types.StringValue(igw.CreatedAt.Format(time.RFC3339)),
			CreatedBy:   types.StringValue(igw.CreatedBy),
			ModifiedAt:  types.StringValue(igw.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:  types.StringValue(igw.ModifiedBy),
			State:       types.StringValue(string(igw.State)),
		}

		state.InternetGateways = append(state.InternetGateways, igwState)
	}

	state.TotalCount = types.Int32Value(int32(data.Count))
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
