package network_logging

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/networklogging"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &networkLoggingNetworkLoggingStorageDataSource{}
	_ datasource.DataSourceWithConfigure = &networkLoggingNetworkLoggingStorageDataSource{}
)

// NewNetworkLoggingNetworkLoggingStorageDataSource is a helper function to simplify the provider implementation.
func NewNetworkLoggingNetworkLoggingStorageDataSource() datasource.DataSource {
	return &networkLoggingNetworkLoggingStorageDataSource{}
}

// networkLoggingNetworkLoggingStorageDataSource is the data source implementation.
type networkLoggingNetworkLoggingStorageDataSource struct {
	config  *scpsdk.Configuration
	client  *networklogging.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *networkLoggingNetworkLoggingStorageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_logging_network_logging_storages"
}

// Schema defines the schema for the data source.
func (d *networkLoggingNetworkLoggingStorageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of network logging storages",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: " Number of items returned per page. \n" +
					"  - example : 10 \n" +
					"  - maximum : 10000 \n" +
					"  - minimum : 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Marker"): schema.StringAttribute{
				Description: "Pagination Start ID. \n" +
					"  - example : 607e0938521643b5b4b266f343fae693 \n" +
					"  - maxLength : 64 \n" +
					"  - minLength : 1",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria in the format 'field_name:asc' for ascending or 'field_name:desc' for descending order. \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("ResourceType"): schema.StringAttribute{
				Description: "Type of the Resource. \n" +
					"  - example : FIREWALL | SECURITY_GROUP | NAT",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("FIREWALL", "SECURITY_GROUP", "NAT"),
				},
			},
			common.ToSnakeCase("NetworkLoggingStorages"): schema.ListNestedAttribute{
				Description: "A List of network logging storages",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Identifier of the Network Logging Storage. \n" +
				                "  - example : 026ee708da3748a28fca4b8fed43d7ce",
							Computed:    true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "Identifier of the Account. \n" +
			    	            "  - example : 232a7dbfb3df46ae81dc11a59fc058b0",
							Computed:    true,
						},
						common.ToSnakeCase("ResourceType"): schema.StringAttribute{
							Description: "Type of the Resource. \n" +
					            "  - example : FIREWALL | SECURITY_GROUP | NAT",
							Computed:    true,
						},
						common.ToSnakeCase("BucketName"): schema.StringAttribute{
							Description: "Name of the Bucket. \n" +
				                "  - example : bucket_name",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				            Description: "The timestamp when the resource was created, in ISO 8601 format. \n" +
                   	            "  - example : 2024-05-17T00:23:17Z",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				            Description: "The user id that created the resource. \n" +
                  	            "  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
				            Description: "The timestamp when the resource was last modified, in ISO 8601 format. \n" +
                   	            "  - example : 2024-05-17T00:23:17Z",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
				            Description: "The user id that last modified the resource. \n" +
                                "  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *networkLoggingNetworkLoggingStorageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.NetworkLogging
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *networkLoggingNetworkLoggingStorageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state networklogging.NetworkLoggingStorageDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetNetworkLoggingStorageList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read network logging storages.",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, networkLoggingStorage := range data.NetworkLoggingStorages {
		networkLoggingStorageState := networklogging.NetworkLoggingStorageResource{
			Id:           types.StringValue(networkLoggingStorage.Id),
			AccountId:    types.StringValue(networkLoggingStorage.AccountId),
			ResourceType: types.StringValue(string(networkLoggingStorage.ResourceType)),
			BucketName:   types.StringValue(networkLoggingStorage.BucketName),
			CreatedAt:    types.StringValue(networkLoggingStorage.CreatedAt.Format(time.RFC3339)),
			CreatedBy:    types.StringValue(networkLoggingStorage.CreatedBy),
			ModifiedAt:   types.StringValue(networkLoggingStorage.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:   types.StringValue(networkLoggingStorage.ModifiedBy),
		}

		state.NetworkLoggingStorages = append(state.NetworkLoggingStorages, networkLoggingStorageState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
