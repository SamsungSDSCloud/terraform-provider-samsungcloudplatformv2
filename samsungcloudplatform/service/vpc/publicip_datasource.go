package vpc

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
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
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPublicipDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_publicips"
}

// Schema defines the schema for the data source.
func (d *vpcPublicipDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of publicip.",
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
			common.ToSnakeCase("IpAddress"): schema.StringAttribute{
				Description: "IP Address \n" +
					"  - example : 172.24.4.2",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State \n" +
					"  - example : RESERVED | ATTACHED",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
				Description: "Attached Resource Type \n" +
					"  - example : VM | LB | BM | NAT_GW",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
				Description: "Attached Resource Name \n" +
					"  - example : VirtualServerName",
				Optional: true,
			},
			common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
				Description: "Attached Resource ID \n" +
					"  - example : 4fa0e75475df40fbb1ab4e74bd60ff37",
				Optional: true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Type \n" +
					"  - example : IGW | GGW | SIGW",
				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "VPC ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Publicips"): schema.ListNestedAttribute{
				Description: "A list of publicip.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Computed:    true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "Type",
							Computed:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
							Computed:    true,
						},
						common.ToSnakeCase("IpAddress"): schema.StringAttribute{
							Description: "IpAddress",
							Computed:    true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "AccountId",
							Computed:    true,
						},
						common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
							Description: "AttachedResourceType",
							Computed:    true,
						},
						common.ToSnakeCase("AttachedResourceName"): schema.StringAttribute{
							Description: "AttachedResourceName",
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

	d.client = inst.Client.Vpc
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *vpcPublicipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.PublicipDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPublicipList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading public IP",
			"Could not read public IP, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	for _, publicip := range data.Publicips {
		publicipState := vpc.Publicip{
			IpAddress:            types.StringValue(publicip.IpAddress),
			CreatedAt:            types.StringValue(publicip.CreatedAt.Format(time.RFC3339)),
			CreatedBy:            types.StringValue(publicip.CreatedBy),
			Description:          types.StringPointerValue(publicip.Description.Get()),
			AttachedResourceName: types.StringPointerValue(publicip.AttachedResourceName.Get()),
			AttachedResourceId:   types.StringPointerValue(publicip.AttachedResourceId.Get()),
			Id:                   types.StringValue(publicip.Id),
			ModifiedAt:           types.StringValue(publicip.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:           types.StringValue(publicip.ModifiedBy),
			Type:                 types.StringValue(string(publicip.Type)),
			AccountId:            types.StringValue(publicip.AccountId),
			State:                types.StringValue(string(publicip.State)),
		}
		attachedResourceType := publicip.AttachedResourceType.Get()
		if attachedResourceType != nil {
			attachedResourceTypeStr := string(*attachedResourceType)
			publicipState.AttachedResourceType = types.StringPointerValue(&attachedResourceTypeStr)
		} else {
			publicipState.AttachedResourceType = types.StringPointerValue(nil)
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
