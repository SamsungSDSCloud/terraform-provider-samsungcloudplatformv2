package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.TypeName = req.ProviderTypeName + "_vpc_private_nat"
}

// Schema defines the schema for the data source.
func (d *vpcPrivateNatDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of Private NAT.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("PrivateNatId"): schema.StringAttribute{
				Description: "The identifier of the private NAT.\n" +
					"  - example : 12f56e27070248a6a240a497e43fbe18 \n",
				Optional: true,
			},
			common.ToSnakeCase("PrivateNat"): schema.SingleNestedAttribute{
				Description: "Private NAT details",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the private NAT.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("Cidr"): schema.StringAttribute{
						Description: "The IP address range of the network in CIDR notation.\n" +
							"  - example : 192.167.0.0/24",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created in ISO 8601 format. \n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the resource.\n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource. \n" +
							"  - example : PrivateNat Description",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the private NAT.\n" +
							"  - example : 12f56e27070248a6a240a497e43fbe18",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified in ISO 8601 format. \n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By \n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the private NAT.\n" +
							"  - example : PrivateNatName",
						Computed: true,
					},
					common.ToSnakeCase("ServiceResourceId"): schema.StringAttribute{
						Description: "The identifier of the connected service resource.\n" +
							"  - example : 3f342bf9a557405b997c2cf48c89cbc2",
						Computed: true,
					},
					common.ToSnakeCase("ServiceResourceName"): schema.StringAttribute{
						Description: "Private NAT connected Service Resource Name \n" +
							"  - example : PrivateNatName",
						Computed: true,
					},
					common.ToSnakeCase("ServiceType"): schema.StringAttribute{
						Description: "The type of the connected service.\n" +
							"  - example : DIRECT_CONNECT | TRANSIT_GATEWAY",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the private NAT.\n" +
							"  - example : ACTIVE",
						Computed: true,
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

	d.client = inst.Client.VpcV1Dot2
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

	data, err := d.client.GetPrivateNat(ctx, state.PrivateNatId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading private nat",
			"Could not read private nat, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	privateNatModel := vpc.PrivateNat{
		Id:                  types.StringValue(data.PrivateNat.Id),
		Name:                types.StringValue(data.PrivateNat.Name),
		ServiceResourceId:   types.StringValue(data.PrivateNat.ServiceResourceId),
		ServiceResourceName: types.StringValue(data.PrivateNat.ServiceResourceName),
		ServiceType:         types.StringValue(string(data.PrivateNat.ServiceType)),
		Cidr:                types.StringValue(data.PrivateNat.Cidr),
		State:               types.StringValue(string(data.PrivateNat.State)),
		Description:         types.StringPointerValue(data.PrivateNat.Description.Get()),
		AccountId:           types.StringValue(data.PrivateNat.AccountId),
		CreatedAt:           types.StringValue(data.PrivateNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(data.PrivateNat.CreatedBy),
		ModifiedAt:          types.StringValue(data.PrivateNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(data.PrivateNat.ModifiedBy),
	}
	privateNatObjectValue, _ := types.ObjectValueFrom(ctx, privateNatModel.AttributeTypes(), privateNatModel)
	state.PrivateNat = privateNatObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
