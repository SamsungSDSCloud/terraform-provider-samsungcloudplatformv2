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
		Description: "list of private nat ips.",
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
				Description: "Private NAT IP ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("PrivateNatId"): schema.StringAttribute{
				Description: "Private NAT ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("IpAddress"): schema.StringAttribute{
				Description: "IP address \n" +
					"  - example : 192.168.10.0",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
				Description: "Attached Resource Name \n" +
					"  - example : attachedResourceName",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
				Description: "Attached Resource Type \n" +
					"  - example : VM | BM",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
				Description: "Attached Resource Id \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State \n" +
					"  - example : ATTACHED | RESERVED",
				Optional: true,
			},
			common.ToSnakeCase("PrivateNatIps"): schema.ListNestedAttribute{
				Description: "A list of private nat ip.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Computed:    true,
						},
						common.ToSnakeCase("IpAddress"): schema.StringAttribute{
							Description: "IpAddress",
							Computed:    true,
						},
						common.ToSnakeCase("PrivateNatId"): schema.StringAttribute{
							Description: "PrivateNatId",
							Computed:    true,
						},
						common.ToSnakeCase("PrivateNatName"): schema.StringAttribute{
							Description: "PrivateNatName",
							Computed:    true,
						},
						common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
							Description: "AttachedResourceName",
							Computed:    true,
						},
						common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
							Description: "AttachedResourceType",
							Computed:    true,
						},
						common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
							Description: "AttachedResourceId",
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
