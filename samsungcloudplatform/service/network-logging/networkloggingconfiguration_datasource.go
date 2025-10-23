package network_logging

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/networklogging"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
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
	_ datasource.DataSource              = &networkLoggingNetworkLoggingConfigurationDataSource{}
	_ datasource.DataSourceWithConfigure = &networkLoggingNetworkLoggingConfigurationDataSource{}
)

// NewNetworkLoggingNetworkLoggingConfigurationDataSource is a helper function to simplify the provider implementation.
func NewNetworkLoggingNetworkLoggingConfigurationDataSource() datasource.DataSource {
	return &networkLoggingNetworkLoggingConfigurationDataSource{}
}

// networkLoggingNetworkLoggingConfigurationDataSource is the data source implementation.
type networkLoggingNetworkLoggingConfigurationDataSource struct {
	config  *scpsdk.Configuration
	client  *networklogging.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *networkLoggingNetworkLoggingConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_logging_network_logging_configurations"
}

// Schema defines the schema for the data source.
func (d *networkLoggingNetworkLoggingConfigurationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of network logging configurations",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "Limit \n" +
					"  - example : 10 \n" +
					"  - maximum : 10000 \n" +
					"  - minimum : 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Marker"): schema.StringAttribute{
				Description: "Marker \n" +
					"  - example : 607e0938521643b5b4b266f343fae693 \n" +
					"  - maxLength : 64 \n" +
					"  - minLength : 1",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("ResourceId"): schema.StringAttribute{
				Description: "ResourceId \n" +
					"  - example : xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Optional: true,
			},
			common.ToSnakeCase("ResourceType"): schema.StringAttribute{
				Description: "ResourceType \n" +
					"  - example : FIREWALL | SECURITY_GROUP | NAT",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("FIREWALL", "SECURITY_GROUP", "NAT"),
				},
			},
			common.ToSnakeCase("ResourceName"): schema.StringAttribute{
				Description: "ResourceName \n" +
					"  - example : FW_IGW_xxxxxx",
				Optional: true,
			},
			common.ToSnakeCase("NetworkLoggingConfigurations"): schema.ListNestedAttribute{
				Description: "A List of network logging configurations",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Computed:    true,
						},
						common.ToSnakeCase("AccountId"): schema.StringAttribute{
							Description: "AccountId",
							Computed:    true,
						},
						common.ToSnakeCase("ResourceId"): schema.StringAttribute{
							Description: "ResourceId",
							Computed:    true,
						},
						common.ToSnakeCase("ResourceType"): schema.StringAttribute{
							Description: "ResourceType",
							Computed:    true,
						},
						common.ToSnakeCase("ResourceName"): schema.StringAttribute{
							Description: "ResourceName",
							Computed:    true,
						},
						common.ToSnakeCase("BucketName"): schema.StringAttribute{
							Description: "BucketName",
							Computed:    true,
						},
						common.ToSnakeCase("SecurityGroupLogId"): schema.StringAttribute{
							Description: "SecurityGroupLogId",
							Computed:    true,
						},
						common.ToSnakeCase("UpInterface"): schema.StringAttribute{
							Description: "UpInterface",
							Computed:    true,
						},
						common.ToSnakeCase("DownInterface"): schema.StringAttribute{
							Description: "DownInterface",
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
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "ModifiedAt",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "ModifiedBy",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *networkLoggingNetworkLoggingConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *networkLoggingNetworkLoggingConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state networklogging.NetworkLoggingConfigurationDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetNetworkLoggingConfigurationList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read network logging configurations.",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, networkLoggingConfiguration := range data.NetworkLoggingConfigurations {
		networkLoggingConfigurationState := networklogging.NetworkLoggingConfiguration{
			Id:                 types.StringValue(networkLoggingConfiguration.Id),
			AccountId:          types.StringValue(networkLoggingConfiguration.AccountId),
			ResourceId:         types.StringValue(networkLoggingConfiguration.ResourceId),
			ResourceType:       types.StringValue(string(networkLoggingConfiguration.ResourceType)),
			ResourceName:       types.StringValue(networkLoggingConfiguration.ResourceName),
			BucketName:         types.StringValue(networkLoggingConfiguration.BucketName),
			SecurityGroupLogId: types.StringPointerValue(networkLoggingConfiguration.SecurityGroupLogId.Get()),
			UpInterface:        types.StringPointerValue(networkLoggingConfiguration.UpInterface.Get()),
			DownInterface:      types.StringPointerValue(networkLoggingConfiguration.DownInterface.Get()),
			CreatedAt:          types.StringValue(networkLoggingConfiguration.CreatedAt.Format(time.RFC3339)),
			CreatedBy:          types.StringValue(networkLoggingConfiguration.CreatedBy),
			ModifiedAt:         types.StringValue(networkLoggingConfiguration.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:         types.StringValue(networkLoggingConfiguration.ModifiedBy),
		}

		state.NetworkLoggingConfigurations = append(state.NetworkLoggingConfigurations, networkLoggingConfigurationState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
