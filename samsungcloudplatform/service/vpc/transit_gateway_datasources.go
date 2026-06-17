package vpc

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
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
				Description: "The current lifecycle state of the firewall connection. \n" +
					"  - enum: [ATTACHING, ACTIVE, DETACHING, DELETED, INACTIVE, ERROR]\n" +
					"  - example: ACTIVE",
				MarkdownDescription: "The current lifecycle state of the firewall connection. \n" +
					"  - enum: [ATTACHING, ACTIVE, DETACHING, DELETED, INACTIVE, ERROR]\n" +
					"  - example: ACTIVE",
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
				Description: "The name of the transit gateway. \n" +
					"  - example : TransitGatewayName",
				MarkdownDescription: "The name of the transit gateway. \n" +
					"  - example : TransitGatewayName",
				Optional: true,
			},
			common.ToSnakeCase("page"): schema.Int32Attribute{
				Optional: true,
				Description: "The page number for pagination.\n" +
					"  - example : 0",
				MarkdownDescription: "The page number for pagination.\n" +
					"  - example : 0",
				Validators: []validator.Int32{
					int32validator.Between(0, 99999),
				},
			},
			common.ToSnakeCase("size"): schema.Int32Attribute{
				Optional: true,
				Description: "The number of items per page.\n" +
					"  - example : 20",
				MarkdownDescription: "The number of items per page.\n" +
					"  - example : 20",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for decending order. \n" +
					"  - example : created_at:desc",
				MarkdownDescription: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for decending order. \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the transit gateway. \n" +
					"  - enum: [ATTACHING, ACTIVE, DETACHING, DELETED, INACTIVE, ERROR]\n" +
					"  - example:ACTIVE",
				MarkdownDescription: "The current lifecycle state of the transit gateway. \n" +
					"  - enum: [ATTACHING, ACTIVE, DETACHING, DELETED, INACTIVE, ERROR]\n" +
					"  - example:ACTIVE",
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
				Description: "The unique identifier of the transit gateway.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("Tgws"): schema.ListNestedAttribute{
				Description: "A list of tgw.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The identifier of the account that owns the transit gateway.\n" +
								"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
							Computed: true,
						},
						common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
							Description: "The bandwidth capacity of the connection.\n" +
								"  - example : 1",
							Computed: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the transit gateway was created in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user id that created the transit gateway.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Enter a brief explanation or note about this transit gateway. This help identify the purpose or usage of the resource.\n" +
								"  - example : Tgw description",
							Computed: true,
						},
						common.ToSnakeCase("firewall_connection_state"): schema.StringAttribute{
							Description: "The current lifecycle state of the firewall connection. \n" +
								"  - example : INACTIVE",
							Computed: true,
						},
						common.ToSnakeCase("FirewallIds"): schema.StringAttribute{
							Description: "List of firewall IDs\n" +
								"  - example : bbb93aca123f4bb2b2c0f206f4a86b2b",
							Computed: true,
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the transit gateway.\n" +
								"  - example : fe860e0af0c04dcd8182b84f907f31f4",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the transit gateway was last modified in ISO 8601 format.\n" +
								"  - example : 2024-05-17T00:23:17Z",
							Computed: true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user id that modified the transit gateway.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the transit gateway.\n" +
								"  - example : Tgw name",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the transit gateway. \n" +
								"  - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING\n" +
								"  - example:ACTIVE",
							Computed: true,
						},
						common.ToSnakeCase("UplinkEnabled"): schema.BoolAttribute{
							Description: "Whether the uplink is enabled.\n" +
								"  - example : false",
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
