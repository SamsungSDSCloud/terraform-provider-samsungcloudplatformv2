package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
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
		Description: "List of security group rule",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page \n" +
					"  - example : 0",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size \n" +
					"  - example : 20",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id \n" +
					"  - example : e09b390420d247e3b6699b2de1b44316",
				Optional: true,
			},
			common.ToSnakeCase("SecurityGroupId"): schema.StringAttribute{
				Description: "SecurityGroupId \n" +
					"  - example : 75a49bb9896a40a49c53e0d8b5a57c06",
				Required: true,
			},
			common.ToSnakeCase("RemoteIpPrefix"): schema.StringAttribute{
				Description: "RemoteIpPrefix \n" +
					"  - example : 10.10.10.10/32",
				Optional: true,
			},
			common.ToSnakeCase("RemoteGroupId"): schema.StringAttribute{
				Description: "RemoteGroupId \n" +
					"  - example : 54aa099f05224ddca958b08945dfb26a",
				Optional: true,
			},
			common.ToSnakeCase("Direction"): schema.StringAttribute{
				Description: "Direction \n" +
					"  - example : ingress",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ingress", "egress"),
				},
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description \n" +
					"  - example : securityGroupRuleDescription",
				Optional: true,
			},
			common.ToSnakeCase("Service"): schema.StringAttribute{
				Description: "Service \n" +
					"  - example : TCP 80",
				Optional: true,
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Security group rule Id List",
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
