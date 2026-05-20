package vpc

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.TypeName = req.ProviderTypeName + "_vpc_transit_gateways"
}

func (d *tgwDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of Transit Gateway",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("firewall_connection_state"): schema.StringAttribute{
				Description:         "- enum: [ATTACHING, ACTIVE, DETACHING, DELETED, INACTIVE, ERROR]",
				MarkdownDescription: "- enum: [ATTACHING, ACTIVE, DETACHING, DELETED, INACTIVE, ERROR]",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"ATTACHING",
						"ACTIVE",
						"DETACHING",
						"DELETED",
						"INACTIVE",
						"ERROR",
					),
				},
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Transit Gateway Name \n" +
					"  - example : TransitGatewayName",
				MarkdownDescription: "Transit Gateway Name \n" +
					"  - example : TransitGatewayName",
				Optional: true,
			},
			common.ToSnakeCase("page"): schema.Int32Attribute{
				Optional:            true,
				Description:         "page",
				MarkdownDescription: "page",
				Validators: []validator.Int32{
					int32validator.Between(0, 99999),
				},
			},
			common.ToSnakeCase("size"): schema.Int32Attribute{
				Optional:            true,
				Description:         "Size (between 1 and 10000)",
				MarkdownDescription: "Size (between 1 and 10000)",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : created_at:desc",
				MarkdownDescription: "Sort \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description:         "- enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\", \"EDITTING\"]",
				MarkdownDescription: "- enum: [\"CREATING\",\"ACTIVE\",\"DELETING\",\"DELETED\",\"ERROR\", \"EDITTING\"]",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"CREATING",
						"ACTIVE",
						"DELETING",
						"DELETED",
						"ERROR",
						"EDITTING",
					),
				},
				Optional: true,
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
								"  - example : Tgw description\n",
							Computed: true,
						},
						common.ToSnakeCase("firewall_connection_state"): schema.StringAttribute{
							Description: "firewall connection state",
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
							Description: "Name\n" +
								"  - example : Tgw name\n",
							Computed: true,
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
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "Total count\n" +
					"  - Example : 20",
				MarkdownDescription: "Total count\n" +
					"  - Example : 20",
				Computed: true,
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

	d.client = inst.Client.VpcV1Dot2
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
		tgwState := vpc.MapToTgw(d)
		state.Tgws = append(state.Tgws, tgwState)
	}

	state.TotalCount = types.Int32Value(data.Count)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
