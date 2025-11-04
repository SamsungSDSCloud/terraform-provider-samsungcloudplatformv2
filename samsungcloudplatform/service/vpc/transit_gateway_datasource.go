package vpc

import (
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
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
	resp.TypeName = req.ProviderTypeName + "_vpc_transit_gateway"
}

func (d *tgwDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Transit Gateway",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Id",
				Required:    true,
			},
			common.ToSnakeCase("TransitGateway"): schema.SingleNestedAttribute{
				Description: "Transit Gateway",
				Computed:    true,
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
						Computed:    true,
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
							" - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING \n" +
							" - example : CREATING \n",
						Computed: true,
					},
					common.ToSnakeCase("UplinkEnabled"): schema.BoolAttribute{
						Description: "UplinkEnabled",
						Computed:    true,
					},
				},
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

	transitGateway := vpc.Tgw{
		Id:            types.StringValue(data.Id),
		Description:   types.StringPointerValue(data.Description.Get()),
		Name:          types.StringValue(data.Name),
		AccountId:     types.StringValue(data.AccountId),
		Bandwidth:     types.Int32PointerValue(data.Bandwidth.Get()),
		CreatedAt:     types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:     types.StringValue(data.CreatedBy),
		FirewallIds:   types.StringPointerValue(data.FirewallIds.Get()),
		ModifiedAt:    types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:    types.StringValue(data.ModifiedBy),
		State:         types.StringValue(string(data.State)),
		UplinkEnabled: types.BoolPointerValue(data.UplinkEnabled),
	}

	tgwObjectValue, _ := types.ObjectValueFrom(ctx, transitGateway.AttributeTypes(), transitGateway)
	state.TransitGateway = tgwObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
