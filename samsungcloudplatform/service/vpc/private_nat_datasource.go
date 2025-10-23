package vpc

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vpcPrivateNatDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcPrivateNatDataSource{}
)

// NewVpcPrivateNatDataSource is a helper function to simplify the provider implementation.
func NewVpcPrivateNatDataSource() datasource.DataSource {
	return &vpcPrivateNatDataSource{}
}

// vpcPrivateNatDataSource is the data source implementation.
type vpcPrivateNatDataSource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPrivateNatDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_private_nats"
}

// Schema defines the schema for the data source.
func (d *vpcPrivateNatDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of private nats.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size \n" +
					"  - example : 20 \n" +
					"  - minimum : 0",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page \n" +
					"  - example : 0 \n" +
					"  - minimum : 0",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Private NAT Id \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Private NAT Name \n" +
					"  - example : privateNatName",
				Optional: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "VPC Name \n" +
					"  - example : vpcName",

				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "VPC Id \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("DirectConnectName"): schema.StringAttribute{
				Description: "Direct Connect Name \n" +
					"  - example : directConnectName",
				Optional: true,
			},
			common.ToSnakeCase("DirectConnectId"): schema.StringAttribute{
				Description: "Direct Connect Id \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "Private NAT IP range \n" +
					"  - example: 192.168.10.0/24 \n",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State \n" +
					"  - example : CREATING | ACTIVE | EDITING | DELETING | ERROR",
				Optional: true,
			},
			common.ToSnakeCase("PrivateNats"): schema.ListNestedAttribute{
				Description: "A list of private nat.",
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
						common.ToSnakeCase("VpcId"): schema.StringAttribute{
							Description: "VpcId",
							Computed:    true,
						},
						common.ToSnakeCase("VpcName"): schema.StringAttribute{
							Description: "VpcName",
							Computed:    true,
						},
						common.ToSnakeCase("DirectConnectId"): schema.StringAttribute{
							Description: "DirectConnectId",
							Computed:    true,
						},
						common.ToSnakeCase("DirectConnectName"): schema.StringAttribute{
							Description: "DirectConnectName",
							Computed:    true,
						},
						common.ToSnakeCase("Cidr"): schema.StringAttribute{
							Description: "Cidr",
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State",
							Computed:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
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
func (d *vpcPrivateNatDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vpcPrivateNatDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.PrivateNatDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPrivateNatList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading private nat",
			"Could not read private nat, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	for _, privateNat := range data.PrivateNats {
		privateNatState := vpc.PrivateNat{
			Id:                types.StringValue(privateNat.Id),
			Name:              types.StringValue(privateNat.Name),
			VpcId:             types.StringValue(privateNat.VpcId),
			VpcName:           types.StringPointerValue(privateNat.VpcName.Get()),
			DirectConnectId:   types.StringValue(privateNat.DirectConnectId),
			DirectConnectName: types.StringPointerValue(privateNat.DirectConnectName.Get()),
			Cidr:              types.StringValue(privateNat.Cidr),
			State:             types.StringValue(string(privateNat.State)),
			Description:       types.StringPointerValue(privateNat.Description.Get()),
			CreatedAt:         types.StringValue(privateNat.CreatedAt.Format(time.RFC3339)),
			CreatedBy:         types.StringValue(privateNat.CreatedBy),
			ModifiedAt:        types.StringValue(privateNat.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:        types.StringValue(privateNat.ModifiedBy),
		}
		state.PrivateNats = append(state.PrivateNats, privateNatState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
