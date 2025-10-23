package vpc

import (
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
)

var (
	_ datasource.DataSource              = &tgwDataSources{}
	_ datasource.DataSourceWithConfigure = &tgwDataSources{}
)

func NewTransitGatewayDataSources() datasource.DataSource {
	return &tgwDataSources{}
}

type tgwDataSources struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

func (d *tgwDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transitgateways"
}

func (d *tgwDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of TGW.",
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
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "State",
				Optional:    true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "id",
				Optional:    true,
			},
			common.ToSnakeCase("Tgws"): schema.ListNestedAttribute{
				Description: "A list of tgw.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
							Description: "Description\n" +
								"  - example : Tgw description\n" +
								"  - maxLength : 50\n" +
								"  - minLength : 1",
							Optional: true,
						},
						common.ToSnakeCase("FirewallIds"): schema.StringAttribute{
							Description: "FirewallIds",
							Optional:    true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Optional:    true,
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
							Description: "Name\n" +
								"  - example : Tgw name\n" +
								"  - pattern : ^[a-zA-Z0-9]*$\n" +
								"  - maxLength : 20\n" +
								"  - minLength : 3",
							Required: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "State" +
								" - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING",
							Computed: true,
						},
						common.ToSnakeCase("UplinkEnabled"): schema.BoolAttribute{
							Description: "UplinkEnabled" +
								"  - example : false\n",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *tgwDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tgwDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpc.TgwDataSource

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := d.client.GetTransitGatewayList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read GetTransitGatewayList",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, d := range data.TransitGateways {
		tgwState := vpc.Tgw{
			Id:            types.StringValue(d.Id),
			Description:   types.StringPointerValue(d.Description.Get()),
			Name:          types.StringValue(d.Name),
			AccountId:     types.StringValue(d.AccountId),
			Bandwidth:     types.Int32PointerValue(d.Bandwidth.Get()),
			CreatedAt:     types.StringValue(d.CreatedAt.Format(time.RFC3339)),
			CreatedBy:     types.StringValue(d.CreatedBy),
			FirewallIds:   types.StringPointerValue(d.FirewallIds.Get()),
			ModifiedAt:    types.StringValue(d.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:    types.StringValue(d.ModifiedBy),
			State:         types.StringValue(string(d.State)),
			UplinkEnabled: types.BoolValue(d.GetUplinkEnabled()),
		}
		state.Tgws = append(state.Tgws, tgwState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
