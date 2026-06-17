package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpc"
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
	_ datasource.DataSource              = &vpcPortDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcPortDataSource{}
)

// NewVpcPortDataSource is a helper function to simplify the provider implementation.
func NewVpcPortDataSource() datasource.DataSource {
	return &vpcPortDataSource{}
}

// vpcSubnetDataSource is the data source implementation.
type vpcPortDataSource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPortDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_ports"
}

// Schema defines the schema for the data source.
func (d *vpcPortDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of port.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "Number of items returned per page.\n" +
					"  - example : 10 \n" +
					"  - maximum : 10000 \n" +
					"  - minimum : 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Marker"): schema.StringAttribute{
				Description: "Pagination Start ID.\n" +
					"  - example : 607e0938521643b5b4b266f343fae693 \n" +
					"  - maxLength : 64 \n" +
					"  - minLength : 1",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order. \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the port.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the port.\n" +
					"  - example : portName",
				Optional: true,
			},
			common.ToSnakeCase("SubnetName"): schema.StringAttribute{
				Description: "The name of the subnet that the port belongs to\n" +
					"  - example : subnetName",
				Optional: true,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "The identifier of the subnet that the port belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
				Description: "The identifier of the resource that this port is attached to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("FixedIpAddress"): schema.StringAttribute{
				Description: "The fixed IP address assigned to the port. \n" +
					"  - example : 172.24.4.2",
				Optional: true,
			},
			common.ToSnakeCase("MacAddress"): schema.StringAttribute{
				Description: "The MAC address of the port.\n" +
					"  - example : fa:16:3e:f7:32:c0",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the port. \n" +
					"  - example : CREATING | ACTIVE | DELETING | ERROR",
				Optional: true,
			},
			common.ToSnakeCase("Ports"): schema.ListNestedAttribute{
				Description: "A list of port.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the port.\n" +
								"  - example : 023c57b1-4f11-4836-8933-8d085e061492",
							Computed: true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the port.\n" +
								"  - example : portName",
							Computed: true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The identifier of the account that owns the port.\n" +
								"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
							Computed: true,
						},
						common.ToSnakeCase("SubnetId"): schema.StringAttribute{
							Description: "The identifier of the subnet that the port belongs to.\n" +
								"  - example : 023c57b14f11483689338d085e061492",
							Computed: true,
						},
						common.ToSnakeCase("SubnetName"): schema.StringAttribute{
							Description: "The name of the subnet that the port belongs to.\n" +
								"  - example : subnetName",
							Computed: true,
						},
						common.ToSnakeCase("VpcId"): schema.StringAttribute{
							Description: "The identifier of the VPC that the port belongs to.\n" +
								"  - example : 071bc63b767444c9afaab1c972d302d5",
							Computed: true,
						},
						common.ToSnakeCase("VpcName"): schema.StringAttribute{
							Description: "The name of the VPC that the port belongs to.\n" +
								"  - example : vpcName",
							Computed: true,
						},
						common.ToSnakeCase("FixedIpAddress"): schema.StringAttribute{
							Description: "The fixed IP address assigned to the port.\n" +
								"  - example : 192.168.1.100",
							Computed: true,
						},
						common.ToSnakeCase("MacAddress"): schema.StringAttribute{
							Description: "The MAC address of the port.\n" +
								"  - example : fa:16:3e:00:00:01",
							Computed: true,
						},
						common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
							Description: "The identifier of the resource that this port is attached to.\n" +
								"  - example : 9387e8f51de04a03994de8a9c3524935",
							Computed: true,
						},
						common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
							Description: "The type of the resource that this port is attached to.\n" +
								"  - example : VM",
							Computed: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource.\n" +
								"  - example : Port Description",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the port.\n" +
								"  - example : ACTIVE",
							Computed: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcPortDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Vpc
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *vpcPortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.PortDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPortList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading port",
			"Could not read port, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	for _, port := range data.Ports {
		portState := vpc.Port{
			Id:                   types.StringValue(port.Id),
			Name:                 types.StringValue(port.Name),
			AccountId:            types.StringValue(port.AccountId),
			SubnetId:             types.StringValue(port.SubnetId),
			SubnetName:           types.StringValue(port.SubnetName),
			VpcId:                types.StringValue(port.VpcId),
			VpcName:              types.StringValue(port.VpcName),
			FixedIpAddress:       types.StringValue(port.FixedIpAddress),
			MacAddress:           types.StringValue(port.MacAddress),
			AttachedResourceId:   types.StringValue(port.AttachedResourceId),
			AttachedResourceType: types.StringValue(port.AttachedResourceType),
			Description:          types.StringValue(port.Description),
			State:                types.StringValue(port.State),
			CreatedAt:            types.StringValue(port.CreatedAt.Format(time.RFC3339)),
			ModifiedAt:           types.StringValue(port.ModifiedAt.Format(time.RFC3339)),
		}

		state.Ports = append(state.Ports, portState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
