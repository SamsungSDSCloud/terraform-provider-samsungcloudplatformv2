package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc_peerings"
}

// Schema defines the schema for the data source.
func (d *vpcPeeringsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of vpc peering.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size\n" +
					"  - Example: 20",
				Optional: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page\n" +
					"  - Example: 0",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort\n" +
					"  - Example: created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "VPC Peering Name",
				Optional:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State\n" +
					"  - Enum: CREATING | ACTIVE | DELETING | DELETED | ERROR | EDITING | CREATING_REQUESTING | REJECTED | CANCELED | DELETING_REQUESTING",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "VPC Peering ID",
				Optional:    true,
			},
			common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
				Description: "Requester VPC ID",
				Optional:    true,
			},
			common.ToSnakeCase("RequesterVpcName"): schema.StringAttribute{
				Description: "Requester VPC Name",
				Optional:    true,
			},
			common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
				Description: "Approver VPC ID",
				Optional:    true,
			},
			common.ToSnakeCase("ApproverVpcName"): schema.StringAttribute{
				Description: "Approver VPC Name",
				Optional:    true,
			},
			common.ToSnakeCase("AccountType"): schema.StringAttribute{
				Description: "Account Type\n" +
					"  - Enum: SAME | DIFFERENT",
				Optional: true,
			},
			common.ToSnakeCase("VpcPeerings"): schema.ListNestedAttribute{
				Description: "Vpc Peerings",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountType"): schema.StringAttribute{
							Description: "Account Type\n" +
								"  - Enum: SAME | DIFFERENT",
							Computed: true,
						},
						common.ToSnakeCase("ApproverVpcAccountId"): schema.StringAttribute{
							Description: "Approver VPC Account ID",
							Computed:    true,
						},
						common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
							Description: "Approver VPC ID",
							Computed:    true,
						},
						common.ToSnakeCase("ApproverVpcName"): schema.StringAttribute{
							Description: "Approver VPC Name",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "Created At\n" +
								"  - Example: 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "Created By\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "VPC Peering Description",
							Computed:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "VPC Peering ID",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "Modified At\n" +
								"  - Example: 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "Modified By\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "VPC Peering Name\n" +
								"  - Minimum length: 3\n" +
								"  - Maximum length: 20\n" +
								"  - Pattern: ^[a-zA-Z0-9-]*$",
							Computed: true,
						},
						common.ToSnakeCase("RequesterVpcAccountId"): schema.StringAttribute{
							Description: "Requester VPC Account ID",
							Computed:    true,
						},
						common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
							Description: "Requester VPC ID",
							Computed:    true,
						},
						common.ToSnakeCase("RequesterVpcName"): schema.StringAttribute{
							Description: "Requester VPC Name",
							Computed:    true,
						},
						common.ToSnakeCase("DeleteRequesterAccountId"): schema.StringAttribute{
							Description: "Requester VPC Account ID",
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State\n" +
								"  - Enum: CREATING | ACTIVE | DELETING | DELETED | ERROR | EDITING | CREATING_REQUESTING | REJECTED | CANCELED | DELETING_REQUESTING",
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
			"Error reading vpc peering",
			"Could not read vpc peering, unexpected error: "+err.Error()+"\nReason: "+detail,
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
