package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	_ "github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	_ "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &securityGroupRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &securityGroupRuleDataSource{}
)

// NewSecurityGroupDataSource is a helper function to simplify the provider implementation.
func NewSecurityGroupRuleDataSource() datasource.DataSource {
	return &securityGroupRuleDataSource{}
}

// securityGroupDataSource is the data source implementation.
type securityGroupRuleDataSource struct {
	config  *scpsdk.Configuration
	client  *securitygroup.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *securityGroupRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_security_group_rule"
}

// Schema defines the schema for the data source.
func (d *securityGroupRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.

	resp.Schema = schema.Schema{
		Description: "Retrieves details of a security group rule.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
				Required: true,
			},
			common.ToSnakeCase("SecurityGroupRule"): schema.SingleNestedAttribute{
				Description: "The security group rule resource details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("SecurityGroupId"): schema.StringAttribute{
						Description: "The identifier of the security group that the resource belongs to.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("ethertype"): schema.StringAttribute{
						Description: "The Ethernet protocol type the rule applies to.\n" +
							"  - example: IPv4",
						Computed: true,
					},
					common.ToSnakeCase("protocol"): schema.StringAttribute{
						Description: "The network protocol the rule applies to.\n" +
							"  - example: TCP",
						Computed: true,
					},
					common.ToSnakeCase("portRangeMin"): schema.Int32Attribute{
						Description: "The minimum port number of the rule's port range.\n" +
							"  - example: 5",
						Computed: true,
					},
					common.ToSnakeCase("portRangeMax"): schema.Int32Attribute{
						Description: "The maximum port number of the rule's port range.\n" +
							"  - example: 10",
						Computed: true,
					},
					common.ToSnakeCase("RemoteIpPrefix"): schema.StringAttribute{
						Description: "The remote IP address range the rule applies to in CIDR notation.\n" +
							"  - example: 10.0.0.0/24",
						Computed: true,
					},
					common.ToSnakeCase("RemoteGroupId"): schema.StringAttribute{
						Description: "The identifier of the remote security group the rule applies to.\n" +
							"  - example: ce5a565f-20fa-48f7-b06d-be0f03d2b50c",
						Computed: true,
					},
					common.ToSnakeCase("RemoteGroupName"): schema.StringAttribute{
						Description: "The name of the remote security group the rule applies to.\n" +
							"  - example: sg-db-prod",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.\n" +
							"  - example: Security group for web tier",
						Computed: true,
					},
					common.ToSnakeCase("Direction"): schema.StringAttribute{
						Description: "The direction of the traffic the rule applies to.\n" +
							"  - example: ingress",
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
func (d *securityGroupRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *securityGroupRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state securitygroup.SecurityGroupRuleDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data, err := d.client.GetSecurityGroupRule(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Security Group Rule",
			"Could not read Security Group Rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	securityGroupRuleElement := data.SecurityGroupRule
	// Map response body to model
	securityGroupRuleModel := securitygroup.SecurityGroupRule{
		Id:              types.StringValue(securityGroupRuleElement.Id),
		SecurityGroupId: types.StringValue(securityGroupRuleElement.SecurityGroupId),
		Ethertype:       types.StringPointerValue(securityGroupRuleElement.Ethertype.Get()),
		Protocol:        types.StringPointerValue(securityGroupRuleElement.Protocol.Get()),
		PortRangeMin:    types.Int32PointerValue(securityGroupRuleElement.PortRangeMin.Get()),
		PortRangeMax:    types.Int32PointerValue(securityGroupRuleElement.PortRangeMax.Get()),
		RemoteIpPrefix:  types.StringPointerValue(securityGroupRuleElement.RemoteIpPrefix.Get()),
		RemoteGroupId:   types.StringPointerValue(securityGroupRuleElement.RemoteGroupId.Get()),
		RemoteGroupName: types.StringPointerValue(securityGroupRuleElement.RemoteGroupName.Get()),
		Description:     types.StringPointerValue(securityGroupRuleElement.Description.Get()),
		Direction:       types.StringValue(string(securityGroupRuleElement.Direction)),
		CreatedAt:       types.StringValue(securityGroupRuleElement.CreatedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(securityGroupRuleElement.CreatedBy),
		ModifiedAt:      types.StringValue(securityGroupRuleElement.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:      types.StringValue(securityGroupRuleElement.ModifiedBy),
	}
	securityGroupObjectValue, _ := types.ObjectValueFrom(ctx, securityGroupRuleModel.AttributeTypes(), securityGroupRuleModel)
	state.SecurityGroupRule = securityGroupObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
