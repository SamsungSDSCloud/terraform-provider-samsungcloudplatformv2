package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
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
				Optional: true,
				Description: "The IP address range of the network in CIDR notation.\n" +
					"  - example : 192.168.0.0/24 \n" +
					"  - maxMask : /28\n" +
					"  - minMask : /16",
				MarkdownDescription: "The IP address range of the network in CIDR notation.\n" +
					"  - example : 192.168.0.0/24 \n" +
					"  - maxMask : /28\n" +
					"  - minMask : /16",
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the subnet.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "The unique identifier of the subnet.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the subnet.\n" +
					"  - example : subnetName",
				MarkdownDescription: "The name of the subnet.\n" +
					"  - example : subnetName",
				Optional: true,
			},
			common.ToSnakeCase("page"): schema.Int32Attribute{
				Optional: true,
				Description: "The page number for pagination.\n" +
					"  - example : 0 ",
				MarkdownDescription: "The page number for pagination.\n" +
					"  - example : 0 ",
				Validators: []validator.Int32{
					int32validator.Between(0, 99999),
				},
			},
			common.ToSnakeCase("size"): schema.Int32Attribute{
				Optional: true,
				Description: "The number of items per page.\n" +
					"  - example : 20 ",
				MarkdownDescription: "The number of items per page.\n" +
					"  - example : 20 ",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order.\n" +
					"  - example : created_at:desc",
				MarkdownDescription: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order.\n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the subnet." +
					"  - enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\"]\n" +
                    "  - example : ACTIVE",
				MarkdownDescription: "The current lifecycle state of the subnet." +
					"  - enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\"]\n" +
                    "  - example : ACTIVE",
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
							Computed: true,
							Description: "The identifier of the account that owns the subnet.\n" +
								"  - example: f1e6c81a2b054582878cb9724dc2ce9f",
							MarkdownDescription: "The identifier of the account that owns the subnet.\n" +
								"  - example: f1e6c81a2b054582878cb9724dc2ce9f",
						},
						common.ToSnakeCase("cidr"): schema.StringAttribute{
							Computed: true,
							Description: "The IP address range of the network in CIDR notation.\n" +
								"  - example: 192.168.0.0/24",
							MarkdownDescription: "The IP address range of the network in CIDR notation.\n" +
								"  - example: 192.168.0.0/24",
						},
						common.ToSnakeCase("created_at"): schema.StringAttribute{
							Computed: true,
							Description: "The timestamp when the subnet was created in ISO 8601 format.\n" +
								"  - example: 2024-05-17T00:23:17Z",
							MarkdownDescription: "The timestamp when the subnet was created in ISO 8601 format.\n" +
								"  - example: 2024-05-17T00:23:17Z",
						},
						common.ToSnakeCase("created_by"): schema.StringAttribute{
							Computed: true,
							Description: "The user id that created the subnet.\n" +
								"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							MarkdownDescription: "The user id that created the subnet.\n" +
								"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
						common.ToSnakeCase("gateway_ip_address"): schema.StringAttribute{
							Computed: true,
							Description: "The gateway IP address of the subnet.\n" +
								"  - example: 192.168.0.1",
							MarkdownDescription: "The gateway IP address of the subnet.\n" +
								"  - example: 192.168.0.1",
						},
						common.ToSnakeCase("id"): schema.StringAttribute{
							Computed: true,
							Description: "The unique identifier of the subnet.\n" +
								"  - example: 023c57b14f11483689338d085e061492",
							MarkdownDescription: "The unique identifier of the subnet.\n" +
								"  - example: 023c57b14f11483689338d085e061492",
						},
						common.ToSnakeCase("modified_at"): schema.StringAttribute{
							Computed: true,
							Description: "The timestamp when the subnet was last modified in ISO 8601 format.\n" +
								"  - example: 2024-05-17T00:23:17Z",
							MarkdownDescription: "The timestamp when the subnet was last modified in ISO 8601 format.\n" +
								"  - example: 2024-05-17T00:23:17Z",
						},
						common.ToSnakeCase("modified_by"): schema.StringAttribute{
							Computed: true,
							Description: "The user id that modified the subnet.\n" +
								"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							MarkdownDescription: "The user id that modified the subnet.\n" +
								"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
						common.ToSnakeCase("name"): schema.StringAttribute{
							Computed: true,
							Description: "The name of the subnet.\n" +
								"  - maxLength: 20\n" +
								"  - minLength: 3\n" +
								"  - pattern: `^[a-zA-Z0-9-]*$`\n" +
								"  - example: subnetName",
							MarkdownDescription: "The name of the subnet.\n" +
								"  - maxLength: 20\n" +
								"  - minLength: 3\n" +
								"  - pattern: `^[a-zA-Z0-9-]*$`\n" +
								"  - example: subnetName",
						},
						common.ToSnakeCase("state"): schema.StringAttribute{
							Computed: true,
							Description: "The current lifecycle state of the subnet." +
								"  - enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\"]\n" +
                                "  - example : ACTIVE",
							MarkdownDescription: "The current lifecycle state of the subnet." +
								"  - enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\"]\n" +
                                "  - example : ACTIVE",
						},
						common.ToSnakeCase("type"): schema.StringAttribute{
							Computed: true,
							Description: "The type of the subnet.\n" +
								"  - enum: [\"GENERAL\",\"LOCAL\",\"VPC_ENDPOINT\"]\n" +
                                "  - example : GENERAL",
							MarkdownDescription: "The type of the subnet.\n" +
								"  - enum: [\"GENERAL\",\"LOCAL\",\"VPC_ENDPOINT\"]\n" +
                                "  - example : GENERAL",
						},
						common.ToSnakeCase("vpc_id"): schema.StringAttribute{
							Computed: true,
							Description: "The identifier of the VPC that the subnet belongs to.\n" +
								"  - example: 7df8abb4912e4709b1cb237daccca7a8",
							MarkdownDescription: "The identifier of the VPC that the subnet belongs to.\n" +
								"  - example: 7df8abb4912e4709b1cb237daccca7a8",
						},
						common.ToSnakeCase("vpc_name"): schema.StringAttribute{
							Computed: true,
							Description: "The name of the VPC that the subnet belongs to.\n" +
								"  - example: vpcName",
							MarkdownDescription: "The name of the VPC that the subnet belongs to.\n" +
								"  - example: vpcName",
						},
					},
				},
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Total count of the Subnet.\n" +
					"  - Example : 20",
				MarkdownDescription: "Total count\n" +
					"  - Example : 20",
				Computed: true,
			},
			common.ToSnakeCase("Type"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The type of the subnet.\n" +
					"  - example : [\"LOCAL\", \"GENERAL\", \"VPC_ENDPOINT\"]",
				MarkdownDescription: "Type \n" +
					"  - example : [\"LOCAL\", \"GENERAL\", \"VPC_ENDPOINT\"]",
				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the subnet belongs to. \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "The identifier of the VPC that the subnet belongs to. \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "The name of the VPC that the subnet belongs to.\n" +
					"  - example : vpcName",
				MarkdownDescription: "The name of the VPC that the subnet belongs to.\n" +
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
