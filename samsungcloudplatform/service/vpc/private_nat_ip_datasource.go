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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vpcPrivateNatIpDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcPrivateNatIpDataSource{}
)

// NewVpcPrivateNatIpDataSource is a helper function to simplify the provider implementation.
func NewVpcPrivateNatIpDataSource() datasource.DataSource {
	return &vpcPrivateNatIpDataSource{}
}

// vpcPrivateNatIpDataSource is the data source implementation.
type vpcPrivateNatIpDataSource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPrivateNatIpDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_private_nat_ips"
}

// Schema defines the schema for the data source.
func (d *vpcPrivateNatIpDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of private Nat ips.",
		Attributes: map[string]schema.Attribute{
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
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the private Nat ip.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("PrivateNatId"): schema.StringAttribute{
				Description: "The identifier of the private NAT that the private Nat ip belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("IpAddress"): schema.StringAttribute{
				Description: "The IP address assigned to the private Nat.\n" +
					"  - example : 192.168.10.0",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
				Description: "The name of the resource that this private Nat ip is attached to.\n" +
					"  - example : attachedResourceName",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
				Description: "The type of the resource that this private Nat ip is attached to.\n" +
					"  - example : VM | BM",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
				Description: "The identifier of the resource that this private Nat ip is attached to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the private Nat ip.\n " +
					"  - example : ATTACHED | RESERVED",
				Optional: true,
			},
			common.ToSnakeCase("PrivateNatIps"): schema.ListNestedAttribute{
				Description: "A list of private nat ip.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the private Nat ip.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("IpAddress"): schema.StringAttribute{
							Description: "The IP address assigned to the private Nat.\n" +
								"  - example : 192.167.0.5",
							Computed: true,
						},
						common.ToSnakeCase("PrivateNatId"): schema.StringAttribute{
							Description: "The identifier of the private NAT that the private Nat ip belongs to.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("PrivateNatName"): schema.StringAttribute{
							Description: "The name of the private NAT that the private NAT ip belongs to.\n" +
								"  - example : privatenatName",
							Computed: true,
						},
						common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
							Description: "The name of the resource that this private Nat ip is attached to.\n" +
								"  - example : resourceName",
							Computed: true,
						},
						common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
							Description: "The type of the resource that this private Nat ip is attached to.\n" +
								"  - example : VM | BM",
							Computed: true,
						},
						common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
							Description: "The identifier of the resource that this private Nat ip is attached to.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the private Nat ip.\n" +
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
func (d *vpcPrivateNatIpDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vpcPrivateNatIpDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.PrivateNatIpDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPrivateNatIpList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading private nat ip",
			"Could not read private nat ip, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	for _, privateNatIp := range data.PrivateNatIps {
		ptr := privateNatIp.AttachedResourceType.Get()
		var attachedResourceType *string
		if ptr != nil {
			s := string(*ptr)
			attachedResourceType = &s
		} else {
			attachedResourceType = nil
		}

		privateNatIpState := vpc.PrivateNatIp{
			Id:                   types.StringValue(privateNatIp.Id),
			IpAddress:            types.StringValue(privateNatIp.IpAddress),
			PrivateNatId:         types.StringValue(privateNatIp.PrivateNatId),
			AttachedResourceName: types.StringPointerValue(privateNatIp.AttachedResourceName.Get()),
			AttachedResourceType: types.StringPointerValue(attachedResourceType),
			AttachedResourceId:   types.StringPointerValue(privateNatIp.AttachedResourceId.Get()),
			State:                types.StringValue(string(privateNatIp.State)),
			CreatedAt:            types.StringValue(privateNatIp.CreatedAt.Format(time.RFC3339)),
			CreatedBy:            types.StringValue(privateNatIp.CreatedBy),
			ModifiedAt:           types.StringValue(privateNatIp.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:           types.StringValue(privateNatIp.ModifiedBy),
		}
		state.PrivateNatIps = append(state.PrivateNatIps, privateNatIpState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
