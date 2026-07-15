package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &networkSecurityGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &networkSecurityGroupDataSource{}
)

// NewSecurityGroupDataSource is a helper function to simplify the provider implementation.
func NewSecurityGroupDataSource() datasource.DataSource {
	return &networkSecurityGroupDataSource{}
}

// securityGroupDataSource is the data source implementation.
type networkSecurityGroupDataSource struct {
	config  *scpsdk.Configuration
	client  *securitygroup.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *networkSecurityGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_security_group"
}

// Schema defines the schema for the data source.
func (d *networkSecurityGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
				Description: "The security group resource details.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
				Required: true,
			},
			common.ToSnakeCase("SecurityGroup"): schema.SingleNestedAttribute{
		Description: "Retrieves details of a security group.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.\n" +
							"  - example: 297615908b8e4ec69520a99a6777add3",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the Security Group.\n" +
							"  - example: sg-web-prod",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.\n" +
							"  - example: Security group for web tier",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the resource.\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("Loggable"): schema.BoolAttribute{
						Description: "Enable flow log for the security group.\n" +
							"  - example: true",
						Computed: true,
					},
					common.ToSnakeCase("RuleCount"): schema.Int32Attribute{
						Description: "Number of rules in the Security Group.\n" +
							"  - example: 5",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
							"  - example: 2025-01-15T10:30:00Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user ID that created the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
							"  - example: 2025-06-01T14:22:00Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user ID that modified the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *networkSecurityGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.SecurityGroup
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *networkSecurityGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state securitygroup.SecurityGroupDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var defaultSize types.Int32
	var defaultPage types.Int32
	var defaultSort types.String
	var defaultName types.String

	ids, err := GetSecurityGroups(d.clients, defaultPage, defaultSize, defaultSort, defaultName, state.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Security Group",
			"Could not read Security Group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if len(ids) > 0 {
		id := ids[0]

		data, err := d.client.GetSecurityGroup(ctx, id.ValueString()) // client 를 호출한다.
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Reading Security Group",
				"Could not read Security Group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}

		securityGroupElement := data.SecurityGroup

		securityGroupModel := securitygroup.SecurityGroup{
			Id:          types.StringValue(securityGroupElement.Id),
			Name:        types.StringValue(securityGroupElement.Name),
			AccountId:   types.StringValue(securityGroupElement.AccountId),
			Description: types.StringPointerValue(securityGroupElement.Description.Get()),
			State:       types.StringValue(securityGroupElement.State),
			Loggable:    types.BoolValue(securityGroupElement.Loggable),
			RuleCount:   types.Int32PointerValue(securityGroupElement.RuleCount),
			CreatedAt:   types.StringValue(securityGroupElement.CreatedAt.Format(time.RFC3339)),
			CreatedBy:   types.StringValue(securityGroupElement.CreatedBy),
			ModifiedAt:  types.StringValue(securityGroupElement.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:  types.StringValue(securityGroupElement.ModifiedBy),
		}
		securityGroupObjectValue, _ := types.ObjectValueFrom(ctx, securityGroupModel.AttributeTypes(), securityGroupModel)
		state.SecurityGroup = securityGroupObjectValue
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
