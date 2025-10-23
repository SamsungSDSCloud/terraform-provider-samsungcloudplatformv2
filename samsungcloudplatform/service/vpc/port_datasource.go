package vpc

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
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
		Description: "list of subnet.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "Limit \n" +
					"  - example : 10 \n" +
					"  - maximum : 10000 \n" +
					"  - minimum : 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Marker"): schema.StringAttribute{
				Description: "Marker \n" +
					"  - example : 607e0938521643b5b4b266f343fae693 \n" +
					"  - maxLength : 64 \n" +
					"  - minLength : 1",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Port ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Port Name \n" +
					"  - example : portName",
				Optional: true,
			},
			common.ToSnakeCase("SubnetName"): schema.StringAttribute{
				Description: "Subnet Name \n" +
					"  - example : subnetName",
				Optional: true,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "Subnet ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
				Description: "Attached Resource ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("FixedIpAddress"): schema.StringAttribute{
				Description: "Fixed IP Address \n" +
					"  - example : 172.24.4.2",
				Optional: true,
			},
			common.ToSnakeCase("MacAddress"): schema.StringAttribute{
				Description: "MAC Address \n" +
					"  - example : fa:16:3e:f7:32:c0",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State \n" +
					"  - example : CREATING | ACTIVE | DELETING | ERROR",
				Optional: true,
			},
			common.ToSnakeCase("Ports"): schema.ListNestedAttribute{
				Description: "A list of port.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Computed:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Name",
							Computed:    true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "AccountId",
							Computed:    true,
						},
						common.ToSnakeCase("SubnetId"): schema.StringAttribute{
							Description: "SubnetId",
							Computed:    true,
						},
						common.ToSnakeCase("SubnetName"): schema.StringAttribute{
							Description: "SubnetName",
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
						common.ToSnakeCase("FixedIpAddress"): schema.StringAttribute{
							Description: "FixedIpAddress",
							Computed:    true,
						},
						common.ToSnakeCase("MacAddress"): schema.StringAttribute{
							Description: "MacAddress",
							Computed:    true,
						},
						common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
							Description: "AttachedResourceId",
							Computed:    true,
						},
						common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
							Description: "AttachedResourceType",
							Computed:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "CreatedAt",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "ModifiedAt",
							Computed:    true,
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
