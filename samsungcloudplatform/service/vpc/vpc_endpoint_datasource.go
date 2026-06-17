package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpc"
	vpcV1Dot2 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
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
	_ datasource.DataSource              = &vpcVpcEndpointDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcVpcEndpointDataSource{}
)

// NewVpcVpcEndpointDataSource is a helper function to simplify the provider implementation.
func NewVpcVpcEndpointDataSource() datasource.DataSource {
	return &vpcVpcEndpointDataSource{}
}

// vpcNatGatewayDataSource is the data source implementation.
type vpcVpcEndpointDataSource struct {
	config    *scpsdk.Configuration
	client    *vpc.Client
	client1d2 *vpcV1Dot2.Client
	clients   *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcVpcEndpointDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc_endpoints"
}

// Schema defines the schema for the data source.
func (d *vpcVpcEndpointDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of vpcendpoints.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page. \n" +
					"  - example : 20 \n" +
					"  - minimum : 0",
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination.\n" +
					"  - example : 0 \n" +
					"  - minimum : 0",
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order.\n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the endpoint.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the endpoint.\n" +
					"  - example : vpcName",
				Optional: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "The name of the VPC that the endpoint belongs to. \n" +
					"  - example : vpcName",

				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the endpoint belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "The identifier of the subnet that the endpoint belongs to.\n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Optional: true,
			},
			common.ToSnakeCase("ResourceType"): schema.StringAttribute{
				Description: "The type of the target resource.(File Storage : FS, Object Storage : OBS, Container Registry : SCR, DNS : DNS)\n" +
					"  - example : FS | OBS | SCR | DNS",
				Optional: true,
			},
			common.ToSnakeCase("ResourceKey"): schema.StringAttribute{
				Description: "The key identifying the target resource of the endpoint.\n" +
					"  - example(case: SCR/DNS) : 07c5364702384471b650147321b52173 \n" +
					"  - example(case: FS/OBS) : 1.1.1.1",
				Optional: true,
			},
			common.ToSnakeCase("EndpointIpAddress"): schema.StringAttribute{
				Description: "The IP address of the endpoint.\n" +
					"  - example : 1.1.1.1",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the endpoint.\n" +
					"  - example : CREATING | ACTIVE | EDITING | DELETING | ERROR",
				Optional: true,
			},

			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Count \n" +
					"  - example : 20 \n",
				Computed: true,
			},
			common.ToSnakeCase("VpcEndpoints"): schema.ListNestedAttribute{
				Description: "A list of endpoints.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the endpoint.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the endpoint.\n" +
								"  - example : endpointName",
							Computed: true,
						},
						common.ToSnakeCase("VpcId"): schema.StringAttribute{
							Description: "The identifier of the VPC that the endpoint belongs to.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("VpcName"): schema.StringAttribute{
							Description: "The name of the VPC that the endpoint belongs to.\n" +
								"  - example : vpcName",
							Computed: true,
						},
						common.ToSnakeCase("SubnetId"): schema.StringAttribute{
							Description: "The identifier of the subnet that the endpoint belongs to.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("SubnetName"): schema.StringAttribute{
							Description: "The name of the subnet that the endpoint belongs to.\n" +
								"  - example : resourceName",
							Computed: true,
						},
						common.ToSnakeCase("EndpointIpAddress"): schema.StringAttribute{
							Description: "The IP address of the endpoint.\n" +
								"  - example : 1.1.1.1",
							Computed: true,
						},
						common.ToSnakeCase("ResourceType"): schema.StringAttribute{
							Description: "The type of the target resource.(File Storage : FS, Object Storage : OBS, Container Registry : SCR, DNS : DNS)\n" +
								"  - example : FS | OBS | SCR | DNS",
							Computed: true,
						},
						common.ToSnakeCase("ResourceKey"): schema.StringAttribute{
							Description: "The key identifying the target resource of the endpoint.\n" +
								"  - example : 07c5364702384471b650147321b52173",
							Computed: true,
						},
						common.ToSnakeCase("ResourceInfo"): schema.StringAttribute{
							Description: "The information about the target resource of the endpoint.\n" +
								"  - example : x.samsungsdscloud.com(Registry)",
							Computed: true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The identifier of the account that owns the endpoint.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the endpoint.\n" +
								"  - example : ACTIVE",
							Computed: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource.\n" +
								"  - example : resourceDescription",
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
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user id that modified the resource.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcVpcEndpointDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client1d2 = inst.Client.VpcV1Dot2 // For VPC endpoint list v1.2
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *vpcVpcEndpointDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpcV1Dot2.VpcEndpointDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client1d2.GetVpcEndpointList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading vpc endpoint",
			"Could not read vpc endpoint, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	state.TotalCount = types.Int32Value(data.Count)
	state.Page = types.Int32Value(data.Page)
	state.Size = types.Int32Value(data.Size)

	// Map response body to model
	for _, vpcendpoint := range data.VpcEndpoints {
		vpcendpointState := vpcV1Dot2.VpcEndpoint{
			Id:                types.StringValue(vpcendpoint.Id),
			Name:              types.StringValue(vpcendpoint.Name),
			VpcId:             types.StringValue(vpcendpoint.VpcId),
			VpcName:           types.StringValue(vpcendpoint.VpcName),
			SubnetId:          types.StringValue(vpcendpoint.SubnetId),
			SubnetName:        types.StringValue(vpcendpoint.SubnetName),
			EndpointIpAddress: types.StringValue(vpcendpoint.EndpointIpAddress),
			ResourceType:      types.StringValue(string(vpcendpoint.ResourceType)),
			ResourceKey:       types.StringValue(vpcendpoint.ResourceKey),
			ResourceInfo:      types.StringValue(vpcendpoint.ResourceKey),
			AccountId:         types.StringValue(vpcendpoint.AccountId),
			State:             types.StringValue(string(vpcendpoint.State)),
			Description:       types.StringPointerValue(vpcendpoint.Description.Get()),
			CreatedAt:         types.StringValue(vpcendpoint.CreatedAt.Format(time.RFC3339)),
			CreatedBy:         types.StringValue(vpcendpoint.CreatedBy),
			ModifiedAt:        types.StringValue(vpcendpoint.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:        types.StringValue(vpcendpoint.ModifiedBy),
		}
		state.VpcEndpoints = append(state.VpcEndpoints, vpcendpointState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
