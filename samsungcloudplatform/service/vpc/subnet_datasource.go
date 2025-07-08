package vpc

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
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
				Description: "Subnet ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Subnet Name \n" +
					"  - example : subnetName",
				Optional: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "VPC Name \n" +
					"  - example : vpcName",
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
			common.ToSnakeCase("Type"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Type \n" +
					"  - example : [\"LOCAL\", \"GENERAL\", \"VPC_ENDPOINT\"]",
				Optional: true,
			},
			common.ToSnakeCase("Subnets"): schema.ListNestedAttribute{
				Description: "A list of subnet.",
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
						common.ToSnakeCase("VpcID"): schema.StringAttribute{
							Description: "VpcID",
							Computed:    true,
						},
						common.ToSnakeCase("VpcName"): schema.StringAttribute{
							Description: "VpcName",
							Computed:    true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "Type",
							Computed:    true,
						},
						common.ToSnakeCase("Cidr"): schema.StringAttribute{
							Description: "Cidr",
							Computed:    true,
						},
						common.ToSnakeCase("GatewayIpAddress"): schema.StringAttribute{
							Description: "GatewayIpAddress",
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
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "CreatedBy",
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
					},
				},
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

	d.client = inst.Client.Vpc
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

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
