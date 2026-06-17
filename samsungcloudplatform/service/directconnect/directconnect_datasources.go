package directconnect

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	directconnectv1d1 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/directconnectv1d1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &directConnectDirectConnectDataSource{}
	_ datasource.DataSourceWithConfigure = &directConnectDirectConnectDataSource{}
)

// NewDirectConnectDirectConnectDataSource is a helper function to simplify the provider implementation.
func NewDirectConnectDirectConnectDataSource() datasource.DataSource {
	return &directConnectDirectConnectDataSource{}
}

// directConnectDirectConnectDataSource1 is the data source implementation.
type directConnectDirectConnectDataSource struct {
	config  *scpsdk.Configuration
	client  *directconnectv1d1.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *directConnectDirectConnectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directconnect_direct_connects"
}

// Schema defines the schema for the data source.
func (d *directConnectDirectConnectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of direct connect.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page. \n" +
					"  - example : 20 \n" +
					"  - minimum : 0",
				Optional: true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "The page number for pagination. \n" +
					"  - example : 0 \n" +
					"  - minimum : 0",
				Optional: true,
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order. \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the direct connect. \n" +
					"  - example : fe860e0af0c04dcd8182b84f907f31f4",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "The name of the direct connect. \n" +
					"  - example : directConnectName",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the direct connect. \n" +
					"  - example : CREATING | ACTIVE | EDITING | DELETING | ERROR",
				Optional: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the direct connect belongs to. \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Optional: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "The name of the VPC that the direct connect belongs to. \n" +
					"  - example : vpcName",
				Optional: true,
			},

			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Description: "The total number of Direct Connect.\n" +
                    "  - example : 5",
				Computed:    true,
			},
			// Output
			common.ToSnakeCase("SortFinal"): schema.ListAttribute{
				Description: "List of sort condition \n" +
					"  - example : [\"created_at:desc\"]",
				ElementType: types.StringType,
				Computed:    true,
			},
			common.ToSnakeCase("DirectConnects"): schema.ListNestedAttribute{
				Description: "A list of direct connect.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Identifier of the direct connect.\n" +
								"  - example : fe860e0af0c04dcd8182b84f907f31f4",
							Computed:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the direct connect.\n" +
								"  - example : directConnectName",
							Computed:    true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "The identifier of the account that owns the direct connect.\n" +
								"  - example : 27bb070b564349f8a31cc60734cc36a5",
							Computed:    true,
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Enter a brief explanation or note about this direct connect. This help identify the purpose or usage of the resource.\n" +
								"  - example : Direct Connect description\n" +
								"  - maxLength : 50\n" +
								"  - minLength : 1",
							Computed:    true,
						},
						common.ToSnakeCase("VpcId"): schema.StringAttribute{
							Description: "The identifier of the VPC that the direct connect belongs to.\n" +
								"  - example : 023c57b14f11483689338d085e061492",
							Computed:    true,
						},
						common.ToSnakeCase("VpcName"): schema.StringAttribute{
							Description: "The name of the VPC that the direct connect belongs to.\n" +
								"  - example : vpc-prod-01",
							Computed:    true,
						},
						common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
							Description: "The bandwidth capacity(1Gpbs, 10Gpbs, 20Gpbs or 40Gpbs) of the connection.\n" +
								"  - example : 1 | 10 | 20 | 40",
							Computed:    true,
						},
						common.ToSnakeCase("FirewallId"): schema.StringAttribute{
							Description: "The identifier of the firewall associated with the direct connect.\n" +
								"  - example : 68db67f78abd405da98a6056a8ee42af",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
                                "  - example : 2024-05-17T00:23:17Z",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "The user id that created the resource.\n" +
                                "  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						    Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
                                "  - example : 2024-05-17T00:23:17Z",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "The user id that last modified the resource.\n" +
                                "  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed:    true,
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Description: "The current lifecycle state of the direct connect.\n" +
								"  - example : CREATING | ACTIVE | EDITING | DELETING | ERROR",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *directConnectDirectConnectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.DirectConnectV1d1
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *directConnectDirectConnectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state directconnectv1d1.DirectConnectDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetDirectConnectList(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading direct connect",
			"Could not read direct connect, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to model
	state.TotalCount = types.Int32Value(data.Count)
	state.Page = types.Int32Value(data.Page)
	state.Size = types.Int32Value(data.Size)
	for _, sortVal := range data.Sort {
		state.SortFinal = append(state.SortFinal, types.StringValue(sortVal))
	}

	for _, dcon := range data.DirectConnects {
		dconState := directconnectv1d1.DirectConnect{
			Id:         types.StringValue(dcon.Id),
			Name:       types.StringValue(dcon.Name),
			AccountId:  types.StringValue(dcon.AccountId),
			VpcId:      types.StringValue(dcon.VpcId),
			VpcName:    types.StringValue(dcon.VpcName),
			Bandwidth:  types.Int32Value(dcon.Bandwidth),
			CreatedAt:  types.StringValue(dcon.CreatedAt.Format(time.RFC3339)),
			CreatedBy:  types.StringValue(dcon.CreatedBy),
			ModifiedAt: types.StringValue(dcon.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy: types.StringValue(dcon.ModifiedBy),
			State:      types.StringValue(string(dcon.State)),
		}
		if dcon.Description.IsSet() {
			if val := dcon.Description.Get(); val != nil {
				dconState.Description = types.StringValue(*val)
			}
		}

		if dcon.FirewallId.IsSet() {
			if val := dcon.FirewallId.Get(); val != nil {
				dconState.FirewallId = types.StringValue(*val)
			}
		}

		state.DirectConnects = append(state.DirectConnects, dconState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
