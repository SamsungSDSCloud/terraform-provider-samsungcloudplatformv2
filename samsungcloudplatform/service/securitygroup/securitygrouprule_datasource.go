package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
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
		Description: "Security group rule",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id \n" +
					"  - example : e09b390420d247e3b6699b2de1b44316",
				Required: true,
			},
			common.ToSnakeCase("SecurityGroupRule"): schema.SingleNestedAttribute{
				Description: "Security group rule",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("SecurityGroupId"): schema.StringAttribute{
						Description: "SecurityGroupId",
						Computed:    true,
					},
					common.ToSnakeCase("ethertype"): schema.StringAttribute{
						Description: "ethertype",
						Computed:    true,
					},
					common.ToSnakeCase("protocol"): schema.StringAttribute{
						Description: "protocol",
						Computed:    true,
					},
					common.ToSnakeCase("portRangeMin"): schema.Int32Attribute{
						Description: "portRangeMin",
						Computed:    true,
					},
					common.ToSnakeCase("portRangeMax"): schema.Int32Attribute{
						Description: "portRangeMax",
						Computed:    true,
					},
					common.ToSnakeCase("RemoteIpPrefix"): schema.StringAttribute{
						Description: "RemoteIpPrefix",
						Computed:    true,
					},
					common.ToSnakeCase("RemoteGroupId"): schema.StringAttribute{
						Description: "RemoteGroupId",
						Computed:    true,
					},
					common.ToSnakeCase("RemoteGroupName"): schema.StringAttribute{
						Description: "RemoteGroupName",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("Direction"): schema.StringAttribute{
						Description: "Direction",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Computed:    true,
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
