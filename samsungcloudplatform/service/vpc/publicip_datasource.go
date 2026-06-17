package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
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
	_ datasource.DataSource              = &vpcPublicipDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcPublicipDataSource{}
)

// NewVpcPublicipDataSource is a helper function to simplify the provider implementation.
func NewVpcPublicipDataSource() datasource.DataSource {
	return &vpcPublicipDataSource{}
}

// vpcPublicipDataSource is the data source implementation.
type vpcPublicipDataSource struct {
	config  *scpsdk.Configuration
	client  *vpcV1Dot2.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPublicipDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_publicips"
}

// Schema defines the schema for the data source.
func (d *vpcPublicipDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of PublicIPs.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
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
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for decending order.\n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("IpAddress"): schema.StringAttribute{
				Description: "The IP address assigned to the resource.\n" +
					"  - example : 192.167.0.5",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the public ip.\n" +
					"  - example : RESERVED | ATTACHED | DELETED",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
				Description: "The type of the resource that this public ip is attached to.\n" +
					"  - example : VM | ALB | LB | BM | DB | NAT_GW | GPU_NODE | VPN | GPU_SERVER | EPAS | POSTGRESQL | MARIADB | SQLSERVER | CACHESTORE | SCALABLEDB | EVENTSTREAMS | SEARCHENGINE | VERTICA | SUBNET | MYSQL",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
				Description: "The name of the resource that this public ip is attached to.\n" +
					"  - example : Attached NAT Gateway Name",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
				Description: "The identifier of the resource that this public ip is attached to.\n" +
					"  - example : 37e6db41f5124184a43251a63124cdc9",
				Optional: true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "The type of the public ip.\n" +
					"  - example : IGW | GGW | SIGW",
				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the public ip belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},

			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "The total number of publicIPs.\n" +
                    "  - example : 2",
				Computed: true,
			},
			common.ToSnakeCase("Publicips"): schema.ListNestedAttribute{
				Description: "A list of public IPs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the public ip.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("IpAddress"): schema.StringAttribute{
							Description: "The IP address assigned to the resource.\n" +
								"  - example : 192.167.0.5",
							Computed: true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The identifier of the account that owns the public ip.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "The type of the public ip.\n" +
                                "  - example : IGW | GGW | SIGW",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the public ip.\n" +
								"  - example : RESERVED | ATTACHED | DELETED",
							Computed: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource.\n" +
								"  - example : resourceDescription",
							Computed: true,
						},
						common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
							Description: "The type of the resource that this public ip is attached to.\n" +
								"  - example : VM | BM",
							Computed: true,
						},
						common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
							Description: "The identifier of the resource that this public ip is attached to.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
							Description: "The name of the resource that this public ip is attached to.\n" +
								"  - example : resourceName",
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
func (d *vpcPublicipDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vpcPublicipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpcV1Dot2.PublicipDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.ListPublicips(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading public IP list",
			"Could not read public IP list, unexpected error: "+err.Error()+"\nReason: "+detail,
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
	for _, publicip := range data.Publicips {
		publicipState := vpcV1Dot2.PublicIp{
			IpAddress:   types.StringValue(publicip.IpAddress),
			CreatedAt:   types.StringValue(publicip.CreatedAt.Format(time.RFC3339)),
			CreatedBy:   types.StringValue(publicip.CreatedBy),
			Description: types.StringPointerValue(publicip.Description.Get()),
			Id:          types.StringValue(publicip.Id),
			ModifiedAt:  types.StringValue(publicip.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:  types.StringValue(publicip.ModifiedBy),
			Type:        types.StringValue(string(publicip.Type)),
			AccountId:   types.StringValue(publicip.AccountId),
			State:       types.StringValue(string(publicip.State)),
		}

		// Handle nullable AttachedResourceType
		if publicip.AttachedResourceType.Get() != nil {
			attachedResourceType := string(*publicip.AttachedResourceType.Get())
			publicipState.AttachedResourceType = types.StringValue(attachedResourceType)
		}

		// Handle nullable AttachedResourceId
		if publicip.AttachedResourceId.Get() != nil {
			publicipState.AttachedResourceId = types.StringValue(*publicip.AttachedResourceId.Get())
		}

		// Handle nullable AttachedResourceName
		if publicip.AttachedResourceName.Get() != nil {
			publicipState.AttachedResourceName = types.StringValue(*publicip.AttachedResourceName.Get())
		}

		state.Publicips = append(state.Publicips, publicipState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
