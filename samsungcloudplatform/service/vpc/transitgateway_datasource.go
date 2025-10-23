package vpc

import (
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
)

var (
	_ datasource.DataSource              = &tgwDataSource{}
	_ datasource.DataSourceWithConfigure = &tgwDataSource{}
)

func NewTransitGatewayDataSource() datasource.DataSource {
	return &tgwDataSource{}
}

type tgwDataSource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

func (d *tgwDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transitgateway"
}

func (d *tgwDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of TGW.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "AccountId",
				Computed:    true,
			},
			common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
				Description: "Bandwidth",
				Computed:    true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description: "CreatedAt",
				Computed:    true,
			},
			common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				Description: "CreatedBy",
				Computed:    true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description",
				Computed:    true,
			},
			common.ToSnakeCase("FirewallIds"): schema.StringAttribute{
				Description: "FirewallIds",
				Computed:    true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Computed:    false,
				Required:    true,
			},
			common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
				Description: "ModifiedAt",
				Computed:    true,
			},
			common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
				Description: "ModifiedBy",
				Computed:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Computed:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State" +
					" - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING",
				Computed: true,
			},
			common.ToSnakeCase("UplinkEnabled"): schema.BoolAttribute{
				Description: "UplinkEnabled",
				Computed:    true,
			},
		},
	}
}

func (d *tgwDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Vpc
	d.clients = inst.Client
}

func (d *tgwDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.TgwDataSourceDetail

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tgwResp, err := d.client.GetTransitGatewayInfo(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read GetTransitGatewayInfo",
			err.Error(),
		)
		return
	}

	if tgwResp == nil {
		return
	}
	data := tgwResp.TransitGateway

	// Map response body to model
	state.Id = types.StringValue(data.Id)
	state.Description = types.StringPointerValue(data.Description.Get())
	state.Name = types.StringValue(data.Name)
	state.AccountId = types.StringValue(data.AccountId)
	state.Bandwidth = types.Int32PointerValue(data.Bandwidth.Get())
	state.CreatedAt = types.StringValue(data.CreatedAt.Format(time.RFC3339))
	state.CreatedBy = types.StringValue(data.CreatedBy)
	state.FirewallIds = types.StringPointerValue(data.FirewallIds.Get())
	state.ModifiedAt = types.StringValue(data.ModifiedAt.Format(time.RFC3339))
	state.ModifiedBy = types.StringValue(data.ModifiedBy)
	state.State = types.StringValue(string(data.State))
	state.UplinkEnabled = types.BoolPointerValue(data.UplinkEnabled)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
