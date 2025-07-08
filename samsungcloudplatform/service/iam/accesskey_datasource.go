package iam

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
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
	_ datasource.DataSource              = &iamAccessKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &iamAccessKeyDataSource{}
)

// NewIamAccessKeyDataSource is a helper function to simplify the provider implementation.
func NewIamAccessKeyDataSource() datasource.DataSource {
	return &iamAccessKeyDataSource{}
}

// iamAccessKeyDataSource is the data source implementation.
type iamAccessKeyDataSource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *iamAccessKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_access_keys"
}

// Schema defines the schema for the data source.
func (d *iamAccessKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "list of access key.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Limit"): schema.Int32Attribute{
				Description: "Limit (between 1 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Marker"): schema.StringAttribute{
				Description: "Marker (between 1 and 64 characters)",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Account Id",
				Optional:    true,
			},
			common.ToSnakeCase("AccessKeys"): schema.ListNestedAttribute{
				Description: "A list of access key.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("AccessKey"): schema.StringAttribute{
							Description: "AccessKey",
							Computed:    true,
						},
						common.ToSnakeCase("AccessKeyType"): schema.StringAttribute{
							Description: "AccessKeyType",
							Computed:    true,
						},
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
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Description: "Description",
							Computed:    true,
						},
						common.ToSnakeCase("ExpirationTimestamp"): schema.StringAttribute{
							Description: "ExpirationTimestamp",
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
						common.ToSnakeCase("ParentAccessKeyId"): schema.StringAttribute{
							Description: "ParentAccessKeyId",
							Computed:    true,
						},
						common.ToSnakeCase("SecretKey"): schema.StringAttribute{
							Description: "SecretKey",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *iamAccessKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Iam
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *iamAccessKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.AccessKeyDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetAccessKeyList(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Access Keys",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, accessKey := range data.AccessKeys {
		accessKeyState := iam.AccessKey{
			AccessKey:           types.StringValue(accessKey.AccessKey),
			AccessKeyType:       types.StringValue(string(accessKey.AccessKeyType)),
			AccountId:           types.StringValue(accessKey.AccountId),
			CreatedAt:           types.StringValue(accessKey.CreatedAt.Format(time.RFC3339)),
			CreatedBy:           types.StringValue(accessKey.CreatedBy),
			Description:         types.StringPointerValue(accessKey.Description.Get()),
			ExpirationTimestamp: types.StringValue(accessKey.ExpirationTimestamp.Format(time.RFC3339)),
			Id:                  types.StringValue(accessKey.Id),
			ModifiedAt:          types.StringValue(accessKey.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:          types.StringValue(accessKey.ModifiedBy),
			ParentAccessKeyId:   types.StringPointerValue(accessKey.ParentAccessKeyId.Get()),
			SecretKey:           types.StringValue(accessKey.SecretKey),
		}

		state.AccessKeys = append(state.AccessKeys, accessKeyState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
