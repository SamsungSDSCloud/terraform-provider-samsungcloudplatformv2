package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
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
				Description: "The number of items per page.\n" +
					"  - Example: 20",
				Optional: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination.\n" +
					"  - Example: 0",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for decending order.\n" +
					"  - Example: created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the peering.\n" +
					"  - example : resourceName",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the resource.\n" +
					"  - Enum: CREATING | ACTIVE | DELETING | DELETED | ERROR | EDITING | CREATING_REQUESTING | REJECTED | CANCELED | DELETING_REQUESTING\n" +
					"  - example:ACTIVE",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the peering.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
				Description: "The identifier of the requester VPC.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("RequesterVpcName"): schema.StringAttribute{
				Description: "The name of the requester VPC.\n" +
					"  - example : vpcName",
				Optional: true,
			},
			common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
				Description: "The identifier of the approver VPC.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("ApproverVpcName"): schema.StringAttribute{
				Description: "The name of the approver VPC.\n" +
					"  - example : vpcName",
				Optional: true,
			},
			common.ToSnakeCase("AccountType"): schema.StringAttribute{
				Description: "The type of account.\n" +
					"  - Enum: SAME | DIFFERENT\n" +
					"  - example:SAME",
				Optional: true,
			},
			common.ToSnakeCase("VpcPeerings"): schema.ListNestedAttribute{
				Description: "List of Vpc Peering",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountType"): schema.StringAttribute{
							Description: "The type of account.\n" +
								"  - Enum: SAME | DIFFERENT\n" +
								"  - example:SAME",
							Computed: true,
						},
						common.ToSnakeCase("ApproverVpcAccountId"): schema.StringAttribute{
							Description: "The identifier of the account that the approver VPC belongs to.\n" +
								"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
							Computed: true,
						},
						common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
							Description: "The identifier of the approver VPC.\n" +
								"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
							Computed: true,
						},
						common.ToSnakeCase("ApproverVpcName"): schema.StringAttribute{
							Description: "The name of the approver VPC.\n" +
								"  - example : vpcName",
							Computed: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
								"  - Example: 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user id that created the resource.\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Enter a brief explanation or note about this resource. This help identify the purpose or usage of the resource.\n" +
								"  - example : resourceDescription",
							Computed: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the peering.\n" +
								"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
								"  - Example: 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user id that modified the resource.\n" +
								"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the peering.\n" +
								"  - example : peering name\n" +
								"  - Minimum length: 3\n" +
								"  - Maximum length: 20\n" +
								"  - Pattern: ^[a-zA-Z0-9-]*$",
							Computed: true,
						},
						common.ToSnakeCase("RequesterVpcAccountId"): schema.StringAttribute{
							Description: "The identifier of the account that the requester VPC belongs to.\n" +
								"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
							Computed: true,
						},
						common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
							Description: "The identifier of the requester VPC.\n" +
								"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
							Computed: true,
						},
						common.ToSnakeCase("RequesterVpcName"): schema.StringAttribute{
							Description: "The name of the requester VPC.\n" +
								"  - example : resourceName",
							Computed: true,
						},
						common.ToSnakeCase("DeleteRequesterAccountId"): schema.StringAttribute{
							Description: "The identifier of account that the deletion requester belongs to.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the peering.\n" +
								"  - Enum: CREATING | ACTIVE | DELETING | DELETED | ERROR | EDITING | CREATING_REQUESTING | REJECTED | CANCELED | DELETING_REQUESTING\n" +
								"  - example:ACTIVE",
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
