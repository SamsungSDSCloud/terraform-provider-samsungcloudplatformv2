package securitygroup

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/securitygroupv1d1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
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
	client  *securitygroupv1d1.Client
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
			common.ToSnakeCase("SecurityGroupRules"): schema.ListNestedAttribute{
				Description: "List of Security Group Rules.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the Security Group Rule.",
							Computed:    true,
						},
						common.ToSnakeCase("SecurityGroupId"): schema.StringAttribute{
							Description: "The identifier of the security group that the rule belongs to.",
							Computed:    true,
						},
						common.ToSnakeCase("Ethertype"): schema.StringAttribute{
							Description: "The layer 3 protocol name.",
							Computed:    true,
						},
						common.ToSnakeCase("Protocol"): schema.StringAttribute{
							Description: "The layer 4 protocol name, or NULL for ALL.",
							Computed:    true,
						},
						common.ToSnakeCase("PortRangeMin"): schema.Int32Attribute{
							Description: "The minimum port number in the range the rule applies to.",
							Computed:    true,
						},
						common.ToSnakeCase("PortRangeMax"): schema.Int32Attribute{
							Description: "The maximum port number in the range the rule applies to.",
							Computed:    true,
						},
						common.ToSnakeCase("RemoteIpPrefix"): schema.StringAttribute{
							Description: "The remote IP address range in CIDR notation.",
							Computed:    true,
						},
						common.ToSnakeCase("RemoteGroupId"): schema.StringAttribute{
							Description: "The identifier of the remote security group.",
							Computed:    true,
						},
						common.ToSnakeCase("RemoteGroupName"): schema.StringAttribute{
							Description: "The name of the remote security group.",
							Computed:    true,
						},
						common.ToSnakeCase("RemoteAddressGroupId"): schema.StringAttribute{
							Description: "The identifier of the remote address group.",
							Computed:    true,
						},
						common.ToSnakeCase("RemoteAddressGroupName"): schema.StringAttribute{
							Description: "The name of the remote address group.",
							Computed:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "A brief explanation or note about this rule.",
							Computed:    true,
						},
						common.ToSnakeCase("Direction"): schema.StringAttribute{
							Description: "The direction of the traffic the rule applies to.",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the rule was created.",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user ID that created the rule.",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the rule was last modified.",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user ID that modified the rule.",
							Computed:    true,
						},
					},
				},
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

	d.client = inst.Client.SecurityGroupV1d1
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *securityGroupRuleDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state securitygroupv1d1.SecurityGroupRuleDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, securityGroupRules, err := GetSecurityGroupRules(d.clients, state.Page, state.Size, state.Sort, state.Id, state.SecurityGroupId, state.RemoteIpPrefix, state.RemoteGroupId, state.Description, state.Direction, state.Service)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SecurityGroupRules",
			err.Error(),
		)
		return
	}

	state.Ids = ids
	state.SecurityGroupRules = securityGroupRules

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetSecurityGroupRules(clients *client.SCPClient, page types.Int32, size types.Int32, sort types.String, id types.String, securityGroupId types.String, remoteIpPrefix types.String, remoteGroupId types.String, description types.String, direction types.String, service types.String) ([]types.String, []securitygroupv1d1.SecurityGroupRule, error) {

	data, err := clients.SecurityGroupV1d1.GetSecurityGroupRuleList(page, size, sort, id, securityGroupId, remoteIpPrefix, remoteGroupId, description, direction, service)
	if err != nil {
		return nil, nil, err
	}

	contents := data.SecurityGroupRules

	var ids []types.String
	var rules []securitygroupv1d1.SecurityGroupRule

	// Map response body to model
	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))

		rule := securitygroupv1d1.SecurityGroupRule{
			Id:              types.StringValue(content.Id),
			SecurityGroupId: types.StringValue(content.SecurityGroupId),
			CreatedAt:       types.StringValue(content.CreatedAt.Format(time.RFC3339)),
			CreatedBy:       types.StringValue(content.CreatedBy),
			ModifiedAt:      types.StringValue(content.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:      types.StringValue(content.ModifiedBy),
			Description:     types.StringPointerValue(content.Description.Get()),
			Direction:       types.StringValue(string(content.Direction)),
		}

		if content.Ethertype.Get() != nil {
			rule.Ethertype = types.StringValue(*content.Ethertype.Get())
		}
		if content.Protocol.Get() != nil {
			rule.Protocol = types.StringValue(*content.Protocol.Get())
		}
		if content.PortRangeMin.Get() != nil {
			rule.PortRangeMin = types.Int32Value(*content.PortRangeMin.Get())
		}
		if content.PortRangeMax.Get() != nil {
			rule.PortRangeMax = types.Int32Value(*content.PortRangeMax.Get())
		}
		if content.RemoteIpPrefix.Get() != nil {
			rule.RemoteIpPrefix = types.StringValue(*content.RemoteIpPrefix.Get())
		}
		if content.RemoteGroupId.Get() != nil {
			rule.RemoteGroupId = types.StringValue(*content.RemoteGroupId.Get())
		}
		if content.RemoteGroupName.Get() != nil {
			rule.RemoteGroupName = types.StringValue(*content.RemoteGroupName.Get())
		}
		if content.RemoteAddressGroupId.Get() != nil {
			rule.RemoteAddressGroupId = types.StringValue(*content.RemoteAddressGroupId.Get())
		}
		if content.RemoteAddressGroupName.Get() != nil {
			rule.RemoteAddressGroupName = types.StringValue(*content.RemoteAddressGroupName.Get())
		}

		rules = append(rules, rule)
	}

	return ids, rules, nil
}
