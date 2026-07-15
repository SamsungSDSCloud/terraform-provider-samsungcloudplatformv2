package vpc

import (
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
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
				Description: "The number of items per page.\n" +
					"  - example : 3",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination.\n" +
					"  - example : 1",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for decending order.\n" +
					"  - example : created_at:asc ",
				Optional: true,
			},
			common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
				Description: "The identifier of the transit gateway that the resource belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the resource belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "The name of the VPC that the resource belongs to.\n" +
					"  - example : vpcName",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the resource.\n" +
					"  - example : ACTIVE",
				Optional: true,
			},
			common.ToSnakeCase("TransitGatewayVpcConnections"): schema.ListNestedAttribute{
				Description: "A list of TransitGateway VpcConnection.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The identifier of the account that owns the connection.\n" +
								"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
							Computed: true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
                            Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
                                "  - example : 2024-05-17T00:23:17Z",
                            Computed: true,
                        },
                        common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
                            Description: "The user id that created the resource.\n" +
                                "  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
                            Computed: true,
                        },
                        common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
                            Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
                                "  - example : 2024-05-17T00:23:17Z",
                            Computed: true,
                        },
                        common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
                            Description: "The user id that modified the resource.\n" +
                                "  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
                            Computed: true,
                        },
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the connection.\n" +
								"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed: true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the connection.\n" +
								"  - enum: CREATING, ACTIVE, DELETING, DELETED, ERROR, EDITING\n" +
								"  - example : ACTIVE",
							Computed: true,
						},
						common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
							Description: "The identifier of the transit gateway that the connection belongs to.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("VpcId"): schema.StringAttribute{
							Description: "The identifier of the VPC that the connection belongs to.\n" +
								"  - example : 7df8abb4912e4709b1cb237daccca7a8",
							Computed: true,
						},
						common.ToSnakeCase("VpcName"): schema.StringAttribute{
							Description: "The name of the VPC that the connection belongs to.\n" +
								"  - example : vpcName",
							Computed: true,
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
