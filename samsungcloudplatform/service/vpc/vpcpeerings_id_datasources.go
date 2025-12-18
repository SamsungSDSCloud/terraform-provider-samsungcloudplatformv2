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
	_ datasource.DataSource              = &vpcPeeringIdDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcPeeringIdDataSource{}
)

// VpcPeeringsDataSource is a helper function to simplify the provider implementation.
func NewVpcVpcPeeringIdDataSource() datasource.DataSource {
	return &vpcPeeringIdDataSource{}
}

// VpcPeeringsDataSource is the data source implementation.
type vpcPeeringIdDataSource struct {
	config  *scpsdk.Configuration
	client  *vpcv1.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPeeringIdDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc_peering"
}

// Schema defines the schema for the data source.
func (d *vpcPeeringIdDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Vpc peering",

		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Required:    true,
			},
			common.ToSnakeCase("VpcPeering"): schema.SingleNestedAttribute{
				Description: "A detail of VpcPeering.",
				Computed:    true,
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
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcPeeringIdDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vpcPeeringIdDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpcv1.VpcPeeringDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetVpcPeering(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading vpc peering detail",
			"Could not read vpc peering detail, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	vpc := vpcv1.VpcPeering{
		Id:                       types.StringValue(data.VpcPeering.Id),
		Name:                     types.StringValue(data.VpcPeering.Name),
		AccountType:              types.StringValue(string(data.VpcPeering.AccountType)),
		ApproverVpcAccountId:     types.StringValue(data.VpcPeering.ApproverVpcAccountId),
		ApproverVpcId:            types.StringValue(data.VpcPeering.ApproverVpcId),
		ApproverVpcName:          types.StringValue(data.VpcPeering.ApproverVpcName),
		Description:              types.StringPointerValue(data.VpcPeering.Description.Get()),
		RequesterVpcAccountId:    types.StringValue(data.VpcPeering.RequesterVpcAccountId),
		RequesterVpcId:           types.StringValue(data.VpcPeering.RequesterVpcId),
		RequesterVpcName:         types.StringValue(data.VpcPeering.RequesterVpcName),
		DeleteRequesterAccountId: types.StringPointerValue(data.VpcPeering.DeleteRequesterAccountId.Get()),
		CreatedAt:                types.StringValue(data.VpcPeering.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(data.VpcPeering.CreatedBy),
		ModifiedAt:               types.StringValue(data.VpcPeering.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(data.VpcPeering.ModifiedBy),
		State:                    types.StringValue(string(data.VpcPeering.State)),
	}

	vpcObjectValue, _ := types.ObjectValueFrom(ctx, vpc.AttributeTypes(), vpc)
	fmt.Printf("plan create-=----------------------\n\n\n\n\n\n\n")
	fmt.Printf("plan create-=----------------------%+v\n\n\n\n\n\n\n", vpc)
	state.VpcPeering = vpcObjectValue
	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
