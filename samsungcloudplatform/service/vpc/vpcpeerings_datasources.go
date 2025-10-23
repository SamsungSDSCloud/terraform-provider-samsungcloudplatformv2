package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/vpcv1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vpcPeeringsDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcPeeringsDataSource{}
)

// VpcPeeringsDataSource is a helper function to simplify the provider implementation.
func NewVpcVpcPeeringsDataSource() datasource.DataSource {
	return &vpcPeeringsDataSource{}
}

// VpcPeeringsDataSource is the data source implementation.
type vpcPeeringsDataSource struct {
	config  *scpsdk.Configuration
	client  *vpcv1.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPeeringsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_peerings"
}

// Schema defines the schema for the data source.
func (d *vpcPeeringsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of vpc peering.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State",
				Optional:    true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "id",
				Optional:    true,
			},
			common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
				Description: "RequesterVpcId",
				Optional:    true,
			},
			common.ToSnakeCase("RequesterVpcName"): schema.StringAttribute{
				Description: "RequesterVpcName",
				Optional:    true,
			},
			common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
				Description: "ApproverVpcId",
				Optional:    true,
			},
			common.ToSnakeCase("ApproverVpcName"): schema.StringAttribute{
				Description: "ApproverVpcName",
				Optional:    true,
			},
			common.ToSnakeCase("AccountType"): schema.StringAttribute{
				Description: "AccountType",
				Optional:    true,
			},
			common.ToSnakeCase("VpcPeerings"): schema.ListNestedAttribute{
				Description: "A list vpc peering.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountType"): schema.StringAttribute{
							Description: "AccountType",
							Computed:    true,
						},
						common.ToSnakeCase("ApproverVpcAccountId"): schema.StringAttribute{
							Description: "ApproverVpcAccountId",
							Computed:    true,
						},
						common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
							Description: "ApproverVpcId",
							Computed:    true,
						},
						common.ToSnakeCase("ApproverVpcName"): schema.StringAttribute{
							Description: "ApproverVpcName",
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
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description\n" +
								"  - example : Tgw description\n" +
								"  - maxLength : 50\n" +
								"  - minLength : 1",
							Optional: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Optional:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "ModifiedAt",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "ModifiedBy",
							Computed:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Name\n" +
								"  - example : Tgw name\n" +
								"  - pattern : ^[a-zA-Z0-9]*$\n" +
								"  - maxLength : 20\n" +
								"  - minLength : 3",
							Required: true,
						},
						common.ToSnakeCase("RequesterVpcAccountId"): schema.StringAttribute{
							Description: "RequesterVpcAccountId",
							Computed:    true,
						},
						common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
							Description: "RequesterVpcId",
							Computed:    true,
						},
						common.ToSnakeCase("RequesterVpcName"): schema.StringAttribute{
							Description: "RequesterVpcName",
							Computed:    true,
						},
						common.ToSnakeCase("DeleteRequesterAccountId"): schema.StringAttribute{
							Description: "DeleteRequesterAccountId",
							Optional:    true,
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State" +
								" - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcPeeringsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.VpcV1
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *vpcPeeringsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpcv1.VpcPeeringListDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetListVpcPeering(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading internet gateway",
			"Could not read internet gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	for _, igw := range data.VpcPeerings {

		vps := vpcv1.VpcPeering{
			Id:                       types.StringValue(igw.Id),
			Name:                     types.StringValue(igw.Name),
			AccountType:              types.StringValue(string(igw.AccountType)),
			ApproverVpcAccountId:     types.StringValue(igw.ApproverVpcAccountId),
			ApproverVpcId:            types.StringValue(igw.ApproverVpcId),
			ApproverVpcName:          types.StringValue(igw.ApproverVpcName),
			Description:              types.StringPointerValue(igw.Description.Get()),
			RequesterVpcAccountId:    types.StringValue(igw.RequesterVpcAccountId),
			RequesterVpcId:           types.StringValue(igw.RequesterVpcId),
			RequesterVpcName:         types.StringValue(igw.RequesterVpcName),
			DeleteRequesterAccountId: types.StringPointerValue(igw.DeleteRequesterAccountId.Get()),
			CreatedAt:                types.StringValue(igw.CreatedAt.Format(time.RFC3339)),
			CreatedBy:                types.StringValue(igw.CreatedBy),
			ModifiedAt:               types.StringValue(igw.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:               types.StringValue(igw.ModifiedBy),
			State:                    types.StringValue(string(igw.State)),
		}

		state.VpcPeerings = append(state.VpcPeerings, vps)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
