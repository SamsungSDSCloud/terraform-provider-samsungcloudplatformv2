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
	_ datasource.DataSource              = &loadbalancerLbMemberDataSources{}
	_ datasource.DataSourceWithConfigure = &loadbalancerLbMemberDataSources{}
)

// NewLoadBalancerLbMemberDataSources is a helper function to simplify the provider implementation.
func NewLoadbalancerLbMemberDataSources() datasource.DataSource {
	return &loadbalancerLbMemberDataSources{}
}

// loadbalancerLbMemberDataSources is the data source implementation.
type loadbalancerLbMemberDataSources struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *loadbalancerLbMemberDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_members"
}

// Schema defines the schema for the data source.
func (d *loadbalancerLbMemberDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Get List of Lb Members.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size",
				Optional:    true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("MemberIp"): schema.StringAttribute{
				Description: "MemberIp",
				Optional:    true,
			},
			common.ToSnakeCase("MemberPort"): schema.Int32Attribute{
				Description: "MemberPort",
				Optional:    true,
			},
			common.ToSnakeCase("LbServerGroupId"): schema.StringAttribute{
				Description: "LbServerGroupId",
				Optional:    true,
			},
			common.ToSnakeCase("LbMembers"): schema.ListNestedAttribute{
				Description: "A list of Lb Members.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Optional:    true,
						},
						common.ToSnakeCase("LbServerGroupId"): schema.StringAttribute{
							Description: "LbServerGroupId",
							Optional:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Name",
							Optional:    true,
						},
						common.ToSnakeCase("MemberIp"): schema.StringAttribute{
							Description: "MemberIp",
							Optional:    true,
						},
						common.ToSnakeCase("MemberPort"): schema.Int32Attribute{
							Description: "MemberPort",
							Optional:    true,
						},
						common.ToSnakeCase("MemberState"): schema.StringAttribute{
							Description: "MemberState",
							Optional:    true,
						},
						common.ToSnakeCase("MemberWeight"): schema.Int32Attribute{
							Description: "MemberWeight",
							Optional:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State",
							Optional:    true,
						},
						common.ToSnakeCase("HealthState"): schema.StringAttribute{
							Description: "HealthState",
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
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *loadbalancerLbMemberDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *loadbalancerLbMemberDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state loadbalancer.LbMemberDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetLbMemberList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Lb Members",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, lbMember := range data.Members {
		lbMemberState := loadbalancer.LbMember{
			Id:              types.StringValue(lbMember.Id),
			LbServerGroupId: types.StringValue(lbMember.LbServerGroupId),
			Name:            types.StringValue(lbMember.Name),
			MemberIp:        types.StringValue(lbMember.MemberIp),
			MemberPort:      types.Int32Value(lbMember.MemberPort),
			MemberState:     types.StringValue(lbMember.MemberState),
			MemberWeight:    types.Int32Value(lbMember.MemberWeight),
			ObjectType:      types.StringValue(string(lbMember.ObjectType)),
			ObjectId:        virtualserverutil.ToNullableStringValue(lbMember.ObjectId.Get()),
			HealthState:     types.StringValue(lbMember.HealthState),
			State:           types.StringValue(string(lbMember.State)),
			CreatedAt:       types.StringValue(lbMember.CreatedAt.Format(time.RFC3339)),
			CreatedBy:       types.StringValue(lbMember.CreatedBy),
			ModifiedAt:      types.StringValue(lbMember.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:      types.StringValue(lbMember.ModifiedBy),
		}

		state.LbMembers = append(state.LbMembers, lbMemberState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
