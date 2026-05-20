package vpc

import (
	"context"
	"fmt"

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
	_ datasource.DataSource              = &vpcVpcDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcVpcDataSource{}
)

// NewVpcVpcDataSource is a helper function to simplify the provider implementation.
func NewVpcVpcDataSource() datasource.DataSource {
	return &vpcVpcDataSource{}
}

// vpcVpcDataSource is the data source implementation.
type vpcVpcDataSource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcVpcDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpcs"
}

// Schema defines the schema for the data source.
func (d *vpcVpcDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of vpc.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "VPC CIDR \n" +
					"  - example : 192.167.0.0/18",
				MarkdownDescription: "VPC CIDR \n" +
					"  - example : 192.167.0.0/18",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "VPC ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "VPC ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "VPC Name \n" +
					"  - example : vpcName",
				MarkdownDescription: "VPC Name \n" +
					"  - example : vpcName",
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
				Description:         "- enum: [\"CREATING\",\"ACTIVE\",\"DELETED\",\"ERROR\"]",
				MarkdownDescription: "- enum: [\"CREATING\",\"ACTIVE\",\"DELETED\",\"ERROR\"]",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"CREATING",
						"ACTIVE",
						"DELETED",
						"ERROR",
					),
				},
			},
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Computed:            true,
				Description:         "count\n  - example: 20",
				MarkdownDescription: "count\n  - example: 20",
			},
			common.ToSnakeCase("Vpcs"): schema.ListNestedAttribute{
				Description: "A list of vpc.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description:         "Account ID\n  - example: f1e6c81a2b054582878cb9724dc2ce9f",
							MarkdownDescription: "Account ID\n  - example: f1e6c81a2b054582878cb9724dc2ce9f",
							Computed:            true,
						},
						common.ToSnakeCase("cidr_count"): schema.Int32Attribute{
							Computed:            true,
							Description:         "Cidr Count\n  - example: 20",
							MarkdownDescription: "Cidr Count\n  - example: 20",
						},
						common.ToSnakeCase("cidrs"): schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"cidr": schema.StringAttribute{
										Computed:            true,
										Description:         "VPC Cidr\n  - example: 192.167.0.0/18",
										MarkdownDescription: "VPC Cidr\n  - example: 192.167.0.0/18",
									},
									"created_at": schema.StringAttribute{
										Computed:            true,
										Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
										MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
									},
									"created_by": schema.StringAttribute{
										Computed:            true,
										Description:         "Created By\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
										MarkdownDescription: "Created By\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
									},
									"id": schema.StringAttribute{
										Computed:            true,
										Description:         "Cidr ID\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
										MarkdownDescription: "Cidr ID\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
									},
								},
							},
							Computed: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
							MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
							Computed:            true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description:         "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							MarkdownDescription: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed:            true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description:         "Description\n  - maxLength: 50\n  - example: vpcDescription",
							MarkdownDescription: "Description\n  - maxLength: 50\n  - example: vpcDescription",
							Computed:            true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description:         "VPC Id\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
							MarkdownDescription: "VPC Id\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
							Computed:            true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
							MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
							Computed:            true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description:         "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							MarkdownDescription: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed:            true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description:         "VPC Name\n  - maxLength: 20\n  - minLength: 3\n  - pattern: `^[a-zA-Z0-9-]*$`\n  - example: vpcName",
							MarkdownDescription: "VPC Name\n  - maxLength: 20\n  - minLength: 3\n  - pattern: `^[a-zA-Z0-9-]*$`\n  - example: vpcName",
							Computed:            true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description:         "- enum: [\"CREATING\",\"ACTIVE\",\"DELETED\",\"ERROR\"]",
							MarkdownDescription: "- enum: [\"CREATING\",\"ACTIVE\",\"DELETED\",\"ERROR\"]",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcVpcDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vpcVpcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.VpcDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetVpcList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading vpc",
			"Could not read vpc, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcs := make([]vpc.VpcDSValue, len(data.Vpcs))

	// Map response body to model
	for pos, vpcElement := range data.Vpcs {

		vpcs[pos] = vpc.ResponseToVpcDSValue(vpcElement)

	}

	state.Vpcs = vpcs

	state.TotalCount = types.Int32Value(data.Count)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
