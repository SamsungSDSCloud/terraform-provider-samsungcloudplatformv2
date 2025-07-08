package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &loadBalancerLbMemberDataSource{}
	_ datasource.DataSourceWithConfigure = &loadBalancerLbMemberDataSource{}
)

// NewLoadBalancerResourceGroupDataSource is a helper function to simplify the provider implementation.
func NewLoadbalancerLbMemberDataSource() datasource.DataSource {
	return &loadBalancerLbMemberDataSource{}
}

// loadBalancerLbMemberDataSource is the data source implementation.
type loadBalancerLbMemberDataSource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadBalancerLbMemberDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_member"
}

// Schema defines the schema for the data source.
func (d *loadBalancerLbMemberDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Lb Member.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Optional:    true,
			},
			common.ToSnakeCase("LbServerGroupId"): schema.StringAttribute{
				Description: "LbServerGroupId",
				Optional:    true,
			},
			common.ToSnakeCase("LbMember"): schema.SingleNestedAttribute{
				Description: "A detail of Lb Member.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
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
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Optional:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "SubnetId",
						Optional:    true,
					},
					common.ToSnakeCase("Uuid"): schema.StringAttribute{
						Description: "Uuid",
						Optional:    true,
					},
					common.ToSnakeCase("ObjectId"): schema.StringAttribute{
						Description: "ObjectId",
						Optional:    true,
					},
					common.ToSnakeCase("ObjectType"): schema.StringAttribute{
						Description: "ObjectType",
						Optional:    true,
					},
					common.ToSnakeCase("MemberWeight"): schema.Int32Attribute{
						Description: "MemberWeight",
						Optional:    true,
					},
					common.ToSnakeCase("MemberState"): schema.StringAttribute{
						Description: "MemberState",
						Optional:    true,
					},
					common.ToSnakeCase("MemberPort"): schema.Int32Attribute{
						Description: "MemberPort",
						Optional:    true,
					},
					common.ToSnakeCase("MemberIp"): schema.StringAttribute{
						Description: "MemberIp",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("LbServerGroupId"): schema.StringAttribute{
						Description: "LbServerGroupId",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadBalancerLbMemberDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.LoadBalancer
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *loadBalancerLbMemberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state loadbalancer.LbMemberDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from routing rule
	data, err := d.client.GetLbMember(ctx, state.LbServerGroupId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading lb member",
			"Could not read lb member ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	lbMember := data.Member.Get()

	var lbMemberState = loadbalancer.LbMemberDetail{
		Name:            types.StringValue(lbMember.Name),
		LbServerGroupId: types.StringValue(lbMember.LbServerGroupId),
		State:           types.StringValue(string(lbMember.State)),
		MemberIp:        types.StringValue(lbMember.MemberIp),
		MemberPort:      types.Int32Value(lbMember.MemberPort),
		MemberState:     types.StringValue(lbMember.MemberState),
		MemberWeight:    types.Int32Value(lbMember.MemberWeight),
		ObjectType:      types.StringValue(string(lbMember.ObjectType)),
		ObjectId:        virtualserverutil.ToNullableStringValue(lbMember.ObjectId.Get()),
		SubnetId:        types.StringValue(lbMember.SubnetId),
		Uuid:            types.StringValue(lbMember.Uuid),
		ModifiedBy:      types.StringValue(lbMember.ModifiedBy),
		ModifiedAt:      types.StringValue(lbMember.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(lbMember.CreatedBy),
		CreatedAt:       types.StringValue(lbMember.CreatedAt.Format(time.RFC3339)),
	}
	lbMemberObjectValue, _ := types.ObjectValueFrom(ctx, lbMemberState.AttributeTypes(), lbMemberState)
	state.LbMemberDetail = lbMemberObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
