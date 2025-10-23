package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
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
		Description: "Security group",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id \n" +
					"  - example : f09708a755e24fceb4e15f7f5c82b0c1",
				Required: true,
			},
			common.ToSnakeCase("SecurityGroup"): schema.SingleNestedAttribute{
				Description: "Security group",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("Loggable"): schema.BoolAttribute{
						Description: "loggable",
						Computed:    true,
					},
					common.ToSnakeCase("RuleCount"): schema.Int32Attribute{
						Description: "rule count",
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
