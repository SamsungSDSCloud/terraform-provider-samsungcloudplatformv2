package securitygroup

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	securitygroupv1d1 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/securitygroupv1d1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &addressGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &addressGroupDataSource{}
)

// NewAddressGroupDataSource is a helper function to simplify the provider implementation.
func NewAddressGroupDataSource() datasource.DataSource {
	return &addressGroupDataSource{}
}

// addressGroupDataSource is the data source implementation.
type addressGroupDataSource struct {
	config  *scpsdk.Configuration
	client  *securitygroupv1d1.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *addressGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_address_group"
}

// Schema defines the schema for the data source.
func (r *addressGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Address Group Data Source",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the Address Group.\n" +
					"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
				Required: true,
			},
			common.ToSnakeCase("AddressGroup"): schema.SingleNestedAttribute{
				Description: "Address Group information.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The account ID associated with the resource.",
						Computed:    true,
					},
					common.ToSnakeCase("AddressCount"): schema.Int32Attribute{
						Description: "Number of addresses in the Address Group.",
						Computed:    true,
					},
					common.ToSnakeCase("AddressLimit"): schema.Int32Attribute{
						Description: "Maximum number of addresses allowed in the Address Group.",
						Computed:    true,
					},
					common.ToSnakeCase("Addresses"): schema.SetAttribute{
						Description: "A list of IP addresses or CIDR ranges in the Address Group.",
						ElementType: types.StringType,
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created in ISO 8601 format.",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user ID that created the resource.",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the resource.",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified in ISO 8601 format.",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user ID that modified the resource.",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the Address Group.",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the resource.",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *addressGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = inst.Client.SecurityGroupV1d1
	r.clients = inst.Client
}

// Read refreshes the data source.
func (r *addressGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config securitygroupv1d1.AddressGroupDataSource
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetAddressGroup(ctx, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading address group",
			"Could not read address group ID "+config.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError(
			"Error Reading address group",
			"Empty response from API",
		)
		return
	}

	agModel := securitygroupv1d1.AddressGroup{
		AccountId:    types.StringValue(data.AddressGroup.AccountId),
		AddressCount: types.Int32Value(data.AddressGroup.GetAddressCount()),
		AddressLimit: types.Int32Value(data.AddressGroup.GetAddressLimit()),
		CreatedAt:    types.StringValue(data.AddressGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:    types.StringValue(data.AddressGroup.CreatedBy),
		Description:  types.StringPointerValue(data.AddressGroup.Description.Get()),
		Id:           types.StringValue(data.AddressGroup.Id),
		ModifiedAt:   types.StringValue(data.AddressGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:   types.StringValue(data.AddressGroup.ModifiedBy),
		Name:         types.StringValue(data.AddressGroup.GetName()),
		State:        types.StringValue(data.AddressGroup.GetState()),
	}

	addressLst, dia := types.SetValueFrom(ctx, types.StringType, data.AddressGroup.Addresses)
	resp.Diagnostics.Append(dia...)
	if resp.Diagnostics.HasError() {
		return
	}
	agModel.Addresses = addressLst

	agObjectValue, d := types.ObjectValueFrom(ctx, agModel.AttributeTypes(), agModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.AddressGroup = agObjectValue

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
