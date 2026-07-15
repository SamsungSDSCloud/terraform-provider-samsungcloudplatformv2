package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vpcPrivateNatDataSources{}
	_ datasource.DataSourceWithConfigure = &vpcPrivateNatDataSources{}
)

// NewVpcPrivateNatDataSources is a helper function to simplify the provider implementation.
func NewVpcPrivateNatDataSources() datasource.DataSource {
	return &vpcPrivateNatDataSources{}
}

// vpcPrivateNatDataSources is the data source implementation.
type vpcPrivateNatDataSources struct {
	config  *scpsdk.Configuration
	client  *vpcv1d2.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPrivateNatDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_private_nats"
}

// Schema defines the schema for the data source.
func (d *vpcPrivateNatDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of private nat.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page. \n" +
					"  - example : 20 \n" +
					"  - minimum : 0",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination. \n" +
					"  - example : 0 \n" +
					"  - minimum : 0",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order. \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the private NAT.\n" +
					"  - example : PrivateNatName",
				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the private NAT belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "The IP address range of the network in CIDR notation. \n" +
					"  - example: 192.167.0.0/24 \n",
				Optional: true,
			},
			common.ToSnakeCase("ServiceResourceId"): schema.StringAttribute{
				Description: "The identifier of the connected service resource.\n" +
					"  - example : 3f342bf9a557405b997c2cf48c89cbc2",
				Optional: true,
			},
			common.ToSnakeCase("ServiceType"): schema.StringAttribute{
				Description: "The type of the connected service.\n" +
					"  - example : DIRECT_CONNECT | TRANSIT_GATEWAY",
				Optional: true,
			},
			common.ToSnakeCase("ServiceResourceName"): schema.StringAttribute{
				Description: "Private NAT connected Service Resource Name \n" +
					"  - example : Service Resource Name",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the private NAT.\n" +
					"  - example : CREATING | ACTIVE | DELETING | DELETED | ERROR",
				Optional: true,
			},

			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "The total number of private nat.\n" +
					"  - example : 2",
				Computed: true,
			},
			common.ToSnakeCase("PrivateNats"): schema.ListNestedAttribute{
				Description: "A list of private nat.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
							Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
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
							Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user id that modified the resource. \n" +
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
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcPrivateNatDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vpcPrivateNatDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpcv1d2.PrivateNatDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.ListPrivateNats(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading private nat",
			"Could not read private nat, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	state.TotalCount = types.Int32Value(data.Count)
	state.Page = types.Int32Value(data.Page)
	state.Size = types.Int32Value(data.Size)
	if len(data.Sort) > 0 {
		state.Sort = types.StringValue(data.Sort[0])
	}

	// Map response body to model
	state.PrivateNats = []vpcv1d2.PrivateNat{}
	for _, privateNat := range data.PrivateNats {
		privateNatState := vpcv1d2.PrivateNat{
			AccountId:           types.StringValue(privateNat.AccountId),
			Cidr:                types.StringValue(privateNat.Cidr),
			CreatedAt:           types.StringValue(privateNat.CreatedAt.Format(time.RFC3339)),
			CreatedBy:           types.StringValue(privateNat.CreatedBy),
			Description:         types.StringPointerValue(privateNat.Description.Get()),
			Id:                  types.StringValue(privateNat.Id),
			ModifiedAt:          types.StringValue(privateNat.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:          types.StringValue(privateNat.ModifiedBy),
			Name:                types.StringValue(privateNat.Name),
			ServiceResourceId:   types.StringValue(privateNat.ServiceResourceId),
			ServiceResourceName: types.StringValue(privateNat.ServiceResourceName),
			ServiceType:         types.StringValue(string(privateNat.ServiceType)),
			State:               types.StringValue(string(privateNat.State)),
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
