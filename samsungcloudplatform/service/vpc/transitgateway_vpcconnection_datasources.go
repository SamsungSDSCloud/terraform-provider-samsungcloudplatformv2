package vpc

import (
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
)

var (
	_ datasource.DataSource              = &tgwVpcConnectionDataSources{}
	_ datasource.DataSourceWithConfigure = &tgwVpcConnectionDataSources{}
)

func NewTransitGatewayVpcConnectionDataSources() datasource.DataSource {
	return &tgwVpcConnectionDataSources{}
}

type tgwVpcConnectionDataSources struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

func (d *tgwVpcConnectionDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_transit_gateway_vpc_connections"
}

func (d *tgwVpcConnectionDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of TGW VPC.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
				Description: "Transit Gateway ID",
				Required:    true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "id",
				Optional:    true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "vpc id",
				Optional:    true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "Vpc Name",
				Optional:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State",
				Optional:    true,
			},
			common.ToSnakeCase("TransitGatewayVpcConnections"): schema.ListNestedAttribute{
				Description: "A list of TransitGateway VpcConnection.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "AccountId",
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
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State" +
								" - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING",
							Computed: true,
						},
						common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
							Description: "Transit Gateway Id",
							Computed:    true,
						},
						common.ToSnakeCase("VpcId"): schema.StringAttribute{
							Description: "vpc id",
							Computed:    true,
						},
						common.ToSnakeCase("VpcName"): schema.StringAttribute{
							Description: "Vpc Name",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *tgwVpcConnectionDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tgwVpcConnectionDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.TgwVpcConnectionDataSource

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := d.client.GetTgwVpcConnectionList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read GetTgw VpcConnectionList",
			err.Error(),
		)
		return
	}

	fmt.Println("GetTgwVpcConnectionList.data", data)

	// Map response body to model
	for _, d := range data.TransitGatewayVpcConnections {
		tgwState := vpc.TgwVpcConnection{
			AccountId:        types.StringValue(d.AccountId),
			CreatedAt:        types.StringValue(d.CreatedAt.Format(time.RFC3339)),
			CreatedBy:        types.StringValue(d.CreatedBy),
			Id:               types.StringValue(d.Id),
			ModifiedAt:       types.StringValue(d.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:       types.StringValue(d.ModifiedBy),
			State:            types.StringValue(string(d.State)),
			TransitGatewayId: types.StringValue(d.TransitGatewayId),
			VpcId:            types.StringValue(d.VpcId),
			VpcName:          types.StringValue(d.VpcName),
		}
		state.TgwVpcConnections = append(state.TgwVpcConnections, tgwState)
	}

	fmt.Println("GetTgwVpcConnectionList.state", state)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
