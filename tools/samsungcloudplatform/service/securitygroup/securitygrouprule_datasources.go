package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	_ "github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	_ "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &securityGroupRuleDataSources{}
	_ datasource.DataSourceWithConfigure = &securityGroupRuleDataSources{}
)

// NewSecurityGroupDataSource is a helper function to simplify the provider implementation.
func NewSecurityGroupRuleDataSources() datasource.DataSource {
	return &securityGroupRuleDataSources{}
}

// securityGroupDataSource is the data source implementation.
type securityGroupRuleDataSources struct {
	config  *scpsdk.Configuration
	client  *securitygroup.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *securityGroupRuleDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_security_group_rules"
}

// Schema defines the schema for the data source.
func (d *securityGroupRuleDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.

		resp.Schema = schema.Schema{
			Description: "List of security group rules",
			Attributes: map[string]schema.Attribute{
				common.ToSnakeCase("Page"): schema.Int32Attribute{
					Description: "The page number for pagination.\n" +
						"  - example: 1\n" +
						"  - constraints: min: 1",
					Optional: true,
					Validators: []validator.Int32{
						int32validator.AtLeast(0),
					},
				},
				common.ToSnakeCase("Size"): schema.Int32Attribute{
					Description: "The number of items per page.\n" +
						"  - example: 20\n" +
						"  - constraints: min: 1",
					Optional: true,
					Validators: []validator.Int32{
						int32validator.AtLeast(0),
					},
				},
				common.ToSnakeCase("Sort"): schema.StringAttribute{
					Description: "The sorting criteria.\n" +
						"  - example: created_at:desc\n" +
						"  - valid: field_name:asc or field_name:desc",
					Optional: true,
				},
				common.ToSnakeCase("Id"): schema.StringAttribute{
					Description: "The unique identifier of the resource.\n" +
						"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
					Optional: true,
				},
				common.ToSnakeCase("SecurityGroupId"): schema.StringAttribute{
					Description: "The identifier of the security group that the resource belongs to.\n" +
						"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
					Required: true,
				},
				common.ToSnakeCase("RemoteIpPrefix"): schema.StringAttribute{
					Description: "The remote IP address range the rule applies to in CIDR notation.\n" +
						"  - example: 10.0.0.0/24\n" +
						"  - valid: IPv4 CIDR",
					Optional: true,
				},
				common.ToSnakeCase("RemoteGroupId"): schema.StringAttribute{
					Description: "The identifier of the remote security group the rule applies to.\n" +
						"  - example: ce5a565f-20fa-48f7-b06d-be0f03d2b50c",
					Optional: true,
				},
				common.ToSnakeCase("Direction"): schema.StringAttribute{
					Description: "The direction of the traffic the rule applies to.\n" +
						"  - example: ingress\n" +
						"  - valid: ingress, egress",
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("ingress", "egress"),
					},
				},
				common.ToSnakeCase("Description"): schema.StringAttribute{
					Description: "A brief explanation or note about this resource.\n" +
						"  - example: Security group for web tier\n" +
						"  - constraints: maxLength: 255",
					Optional: true,
				},
				common.ToSnakeCase("Service"): schema.StringAttribute{
					Description: "The service ports the rule applies to.\n" +
						"  - example: TCP 80",
					Optional: true,
				},
				common.ToSnakeCase("Ids"): schema.ListAttribute{
					ElementType: types.StringType,
					Computed:    true,
					Description: "The list of Security Group Rule identifiers.\n" +
						"  - example: [sgr-aaa111, sgr-bbb222, sgr-ccc333]",
				},
			},
		}
}

// Configure adds the provider configured client to the data source.
func (d *securityGroupRuleDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *securityGroupRuleDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state securitygroup.SecurityGroupRuleDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetSecurityGroupRules(d.clients, state.Page, state.Size, state.Sort, state.Id, state.SecurityGroupId, state.RemoteIpPrefix, state.RemoteGroupId, state.Description, state.Direction, state.Service)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SecurityGroupRules",
			err.Error(),
		)
		return
	}

	state.Ids = ids // 상태 데이터에 ID 리스트를 추가한다.

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetSecurityGroupRules(clients *client.SCPClient, page types.Int32, size types.Int32, sort types.String, id types.String, securityGroupId types.String, remoteIpPrefix types.String, remoteGroupId types.String, description types.String, direction types.String, service types.String) ([]types.String, error) {

	data, err := clients.SecurityGroup.GetSecurityGroupRuleList(page, size, sort, id, securityGroupId, remoteIpPrefix, remoteGroupId, description, direction, service)
	if err != nil {
		return nil, err
	}

	contents := data.SecurityGroupRules

	var ids []types.String

	// Map response body to model
	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
