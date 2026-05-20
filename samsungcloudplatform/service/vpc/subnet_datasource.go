package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vpcSubnetDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcSubnetDataSource{}
)

// NewVpcSubnetDataSource is a helper function to simplify the provider implementation.
func NewVpcSubnetDataSource() datasource.DataSource {
	return &vpcSubnetDataSource{}
}

// vpcSubnetDataSource is the data source implementation.
type vpcSubnetDataSource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcSubnetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_subnets"
}

// Schema defines the schema for the data source.
func (d *vpcSubnetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of subnet.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("cidr"): schema.StringAttribute{
				Optional:            true,
				Description:         "Subnet Cidr",
				MarkdownDescription: "Subnet Cidr",
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Subnet ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "Subnet ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Subnet Name \n" +
					"  - example : subnetName",
				MarkdownDescription: "Subnet Name \n" +
					"  - example : subnetName",
				Optional: true,
			},
			common.ToSnakeCase("page"): schema.Int32Attribute{
				Optional:            true,
				Description:         "page",
				MarkdownDescription: "page",
				Validators: []validator.Int32{
					int32validator.Between(0, 99999),
				},
			},
			common.ToSnakeCase("size"): schema.Int32Attribute{
				Optional:            true,
				Description:         "size",
				MarkdownDescription: "size",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : created_at:desc",
				MarkdownDescription: "Sort \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description:         "- enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\"]",
				MarkdownDescription: "- enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\"]",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"CREATING",
						"ACTIVE",
						"DELETING",
						"DELETED",
						"ERROR",
					),
				},
				Optional: true,
			},
			common.ToSnakeCase("Subnets"): schema.ListNestedAttribute{
				Description:         "A list of subnet.",
				MarkdownDescription: "A list of subnet.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("account_id"): schema.StringAttribute{
							Computed:            true,
							Description:         "Account ID\n  - example: f1e6c81a2b054582878cb9724dc2ce9f",
							MarkdownDescription: "Account ID\n  - example: f1e6c81a2b054582878cb9724dc2ce9f",
						},
						common.ToSnakeCase("cidr"): schema.StringAttribute{
							Computed:            true,
							Description:         "Subnet Cidr\n  - example: 192.167.1.0/24",
							MarkdownDescription: "Subnet Cidr\n  - example: 192.167.1.0/24",
						},
						common.ToSnakeCase("created_at"): schema.StringAttribute{
							Computed:            true,
							Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
							MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
						},
						common.ToSnakeCase("created_by"): schema.StringAttribute{
							Computed:            true,
							Description:         "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							MarkdownDescription: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
						common.ToSnakeCase("gateway_ip_address"): schema.StringAttribute{
							Computed:            true,
							Description:         "Gateway IP Address\n  - example: 192.167.1.1",
							MarkdownDescription: "Gateway IP Address\n  - example: 192.167.1.1",
						},
						common.ToSnakeCase("id"): schema.StringAttribute{
							Computed:            true,
							Description:         "Subnet Id\n  - example: 023c57b14f11483689338d085e061492",
							MarkdownDescription: "Subnet Id\n  - example: 023c57b14f11483689338d085e061492",
						},
						common.ToSnakeCase("modified_at"): schema.StringAttribute{
							Computed:            true,
							Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
							MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
						},
						common.ToSnakeCase("modified_by"): schema.StringAttribute{
							Computed:            true,
							Description:         "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							MarkdownDescription: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
						common.ToSnakeCase("name"): schema.StringAttribute{
							Computed:            true,
							Description:         "Subnet Name\n  - maxLength: 20\n  - minLength: 3\n  - pattern: `^[a-zA-Z0-9-]*$`\n  - example: subnetName",
							MarkdownDescription: "Subnet Name\n  - maxLength: 20\n  - minLength: 3\n  - pattern: `^[a-zA-Z0-9-]*$`\n  - example: subnetName",
						},
						common.ToSnakeCase("state"): schema.StringAttribute{
							Computed:            true,
							Description:         "- enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\"]",
							MarkdownDescription: "- enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\"]",
						},
						common.ToSnakeCase("type"): schema.StringAttribute{
							Computed:            true,
							Description:         "- enum: [\"GENERAL\",\"LOCAL\",\"VPC_ENDPOINT\"]",
							MarkdownDescription: "- enum: [\"GENERAL\",\"LOCAL\",\"VPC_ENDPOINT\"]",
						},
						common.ToSnakeCase("vpc_id"): schema.StringAttribute{
							Computed:            true,
							Description:         "VPC Id\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
							MarkdownDescription: "VPC Id\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
						},
						common.ToSnakeCase("vpc_name"): schema.StringAttribute{
							Computed:            true,
							Description:         "VPC Name\n  - example: vpcName",
							MarkdownDescription: "VPC Name\n  - example: vpcName",
						},
					},
				},
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Total count\n" +
					"  - Example : 20",
				MarkdownDescription: "Total count\n" +
					"  - Example : 20",
				Computed: true,
			},
			common.ToSnakeCase("Type"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Type \n" +
					"  - example : [\"LOCAL\", \"GENERAL\", \"VPC_ENDPOINT\"]",
				MarkdownDescription: "Type \n" +
					"  - example : [\"LOCAL\", \"GENERAL\", \"VPC_ENDPOINT\"]",
				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "VPC ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "VPC ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "VPC Name \n" +
					"  - example : vpcName",
				MarkdownDescription: "VPC Name \n" +
					"  - example : vpcName",
				Optional: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcSubnetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vpcSubnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.SubnetDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetSubnetList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading subnet",
			"Could not read subnet, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	for _, subnet := range data.Subnets {
		subnetState := vpc.Subnet{
			Id:               types.StringValue(subnet.Id),
			Name:             types.StringValue(subnet.Name),
			AccountId:        types.StringValue(subnet.AccountId),
			VpcId:            types.StringValue(subnet.VpcId),
			VpcName:          types.StringValue(subnet.VpcName),
			Type:             types.StringValue(string(subnet.Type)),
			Cidr:             types.StringValue(subnet.Cidr),
			GatewayIpAddress: types.StringPointerValue(subnet.GatewayIpAddress.Get()),
			State:            types.StringValue(string(subnet.State)),
			CreatedAt:        types.StringValue(subnet.CreatedAt.Format(time.RFC3339)),
			CreatedBy:        types.StringValue(subnet.CreatedBy),
			ModifiedAt:       types.StringValue(subnet.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:       types.StringValue(subnet.ModifiedBy),
		}

		state.Subnets = append(state.Subnets, subnetState)
	}

	state.TotalCount = types.Int32Value(data.Count)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
