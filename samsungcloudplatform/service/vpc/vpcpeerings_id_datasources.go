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
				Description: "The unique identifier of the peering.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("VpcPeering"): schema.SingleNestedAttribute{
				Description: "A detail of VpcPeering.",
				Computed:    true,
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
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
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
						Description: "The timestamp when the resource was last modified in ISO 8601 format.\n " +
							"  - Example: 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that modified the resource.\n " +
							"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the resource.\n" +
							"  - Minimum length: 3\n" +
							"  - Maximum length: 20\n" +
							"  - Pattern: ^[a-zA-Z0-9-]*$",
						Computed: true,
					},
					common.ToSnakeCase("RequesterVpcAccountId"): schema.StringAttribute{
						Description: "The identifier of the account that the requester VPC belongs to.\n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9fn",
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
