package virtualserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualServerServerGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualServerServerGroupDataSource{}
)

func NewVirtualServerServerGroupDataSource() datasource.DataSource {
	return &virtualServerServerGroupDataSource{}
}

type virtualServerServerGroupDataSource struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (d *virtualServerServerGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_server_group"
}

func (d *virtualServerServerGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of server groups.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "ID",
				Optional:    true,
			},
			common.ToSnakeCase("ServerGroup"): schema.SingleNestedAttribute{
				Description: "Server Group.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "ID",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("Policy"): schema.StringAttribute{
						Description: "Policy",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID",
						Computed:    true,
					},
					common.ToSnakeCase("UserId"): schema.StringAttribute{
						Description: "User ID",
						Computed:    true,
					},
					common.ToSnakeCase("Members"): schema.ListAttribute{
						Description: "Members",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *virtualServerServerGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.VirtualServer
	d.clients = inst.Client

}

func (d *virtualServerServerGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualserver.ServerGroupDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetServerGroups(d.clients, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Server Group",
			err.Error(),
		)
	}

	if len(ids) > 0 || !state.Id.IsNull() {
		id := virtualserverutil.SetResourceIdentifier(state.Id, ids, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		serverGroup, err := d.client.GetServerGroup(ctx, id.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Reading Server Group",
				"Could not read Server Group ID "+id.ValueString()+": "+err.Error()+"\nReason: "+detail,
			)
			return
		}

		members := make([]attr.Value, len(serverGroup.Members))
		for i, member := range serverGroup.Members {
			members[i] = types.StringValue(member)
		}

		serverGroupModel := virtualserver.ServerGroup{
			Id:        types.StringValue(serverGroup.Id),
			Name:      types.StringValue(serverGroup.Name),
			Policy:    types.StringValue(serverGroup.Policy),
			AccountId: types.StringValue(serverGroup.AccountId),
			UserId:    types.StringValue(serverGroup.UserId),
			Members:   types.ListValueMust(types.StringType, members),
		}
		serverGroupObjectValue, _ := types.ObjectValueFrom(ctx, serverGroupModel.AttributeTypes(), serverGroupModel)
		state.ServerGroup = serverGroupObjectValue
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
