package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &iamAccessKeyResource{}
	_ resource.ResourceWithConfigure   = &iamAccessKeyResource{}
	_ resource.ResourceWithImportState = &iamAccessKeyResource{}
)

// NewIamAccessKeyResource is a helper function to simplify the provider implementation.
func NewIamAccessKeyResource() resource.Resource {
	return &iamAccessKeyResource{}
}

// iamAccessKeyResource is the data source implementation.
type iamAccessKeyResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *iamAccessKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_access_key"
}

// Schema defines the schema for the data source.
func (r *iamAccessKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM Access Key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the access key resource.\n" +
					"  - example : '12345678-1234-1234-1234-1234567890ab'",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the access key.\n" +
					"  - example : '2024-01-01 00:00:00'",
				Computed: true,
			},
			common.ToSnakeCase("AccessKeyType"): schema.StringAttribute{
				Description: "Type of access key determining its expiration policy.\n" +
					"  - example : 'PERMANENT' | 'TEMPORARY' | 'SECRET_VAULT_TEMPORARY'",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("PERMANENT", "TEMPORARY"),
				},
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Account ID that owns the access key.\n" +
					"  - example : '123456789012'",
				Computed: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Human-readable description of the access key.\n" +
					"  - example : 'My Access Key'",
				Optional: true,
			},
			common.ToSnakeCase("Duration"): schema.StringAttribute{
				Description: "Duration for temporary access key validity (in ISO 8601 duration format).\n" +
					"  - example : 'PT1H' (1 hour)",
				Optional: true,
			},
			common.ToSnakeCase("ParentAccessKeyId"): schema.StringAttribute{
				Description: "Parent access key ID if this is a derived key.\n" +
					"  - example : '12345678-1234-1234-1234-1234567890ab'",
				Optional: true,
			},
			common.ToSnakeCase("IsEnabled"): schema.BoolAttribute{
				Description: "Whether the access key is enabled/active.\n" +
					"  - example : true",
				Computed: true,
				Optional: true,
			},
			common.ToSnakeCase("Passcode"): schema.StringAttribute{
				Description: "Passcode required for access key creation.\n" +
					"  - example : '123456'",
				Optional: true,
			},
			common.ToSnakeCase("AccessKey"): schema.SingleNestedAttribute{
				Description: "Access key details.\n" +
					"  - example : '{access_key: b754b12b39da4ce29a40c5e324650bd0, access_key_type: PERMANENT, account_id: f39c460fade34fecb05ede8f904b24b7, ...}'",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccessKey"): schema.StringAttribute{
						Description: "The access key string value.\n" +
							"  - example : 'ak-example-access-key-id'",
						Computed: true,
					},
					common.ToSnakeCase("AccessKeyType"): schema.StringAttribute{
						Description: "Type of access key determining its expiration policy.\n" +
							"  - example : 'PERMANENT' | 'TEMPORARY' | 'SECRET_VAULT_TEMPORARY'",
						Computed: true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID that owns the access key.\n" +
							"  - example : '123456789012'",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Timestamp when the access key was created.\n" +
							"  - example : '2024-01-01T00:00:00Z'",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "User who created the access key.\n" +
							"  - example : 'user@example.com'",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Human-readable description of the access key.\n" +
							"  - example : 'My access key description'",
						Computed: true,
					},
					common.ToSnakeCase("ExpirationTimestamp"): schema.StringAttribute{
						Description: "Timestamp when the access key expires (for temporary keys).\n" +
							"  - example : '2024-01-02T00:00:00Z'",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Unique identifier of the access key.\n" +
							"  - example : 'ak-1234567890abcdef'",
						Computed: true,
					},
					common.ToSnakeCase("IsEnabled"): schema.BoolAttribute{
						Description: "Whether the access key is enabled/active.\n" +
							"  - example : true",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Timestamp when the access key was last modified.\n" +
							"  - example : '2024-01-01T00:00:00Z'",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "User who last modified the access key.\n" +
							"  - example : 'user@example.com'",
						Computed: true,
					},
					common.ToSnakeCase("ParentAccessKeyId"): schema.StringAttribute{
						Description: "Parent access key ID if this is a derived key.\n" +
							"  - example : '12345678-1234-1234-1234-1234567890ab'",
						Computed: true,
					},
					common.ToSnakeCase("SecretKey"): schema.StringAttribute{
						Description: "The secret key string value (only available at creation time).\n" +
							"  - example : 'sk-example-secret-key-value'",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *iamAccessKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Iam
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *iamAccessKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan iam.AccessKeyResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountId, err := r.clients.Iam.GetAccountId()
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading AccountId",
			"Could not read Account ID, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.AccountId = types.StringValue(accountId)

	// Create new access key
	data, err := r.client.CreateAccessKey(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating access key",
			"Could not create access key, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	accessKey := data.AccessKey
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(accessKey.Id)

	accessKeyModel := iam.AccessKey{
		AccessKey:           types.StringValue(accessKey.AccessKey),
		AccessKeyType:       types.StringValue(string(accessKey.AccessKeyType)),
		AccountId:           types.StringValue(accessKey.AccountId),
		CreatedAt:           types.StringValue(accessKey.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(accessKey.CreatedBy),
		Description:         types.StringPointerValue(accessKey.Description.Get()),
		ExpirationTimestamp: types.StringValue(accessKey.ExpirationTimestamp.Format(time.RFC3339)),
		Id:                  types.StringValue(accessKey.Id),
		IsEnabled:           types.BoolValue(accessKey.IsEnabled),
		ModifiedAt:          types.StringValue(accessKey.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(accessKey.ModifiedBy),
		ParentAccessKeyId:   types.StringPointerValue(accessKey.ParentAccessKeyId.Get()),
		SecretKey:           types.StringValue(accessKey.SecretKey),
	}
	accessKeyObjectValue, diags := types.ObjectValueFrom(ctx, accessKeyModel.AttributeTypes(), accessKeyModel)
	plan.AccessKey = accessKeyObjectValue
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.IsEnabled = types.BoolValue(accessKeyModel.IsEnabled.ValueBool())
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *iamAccessKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state iam.AccessKeyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from access key
	data, err := r.client.GetAccessKey(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading AccessKey",
			"Could not read AccessKey ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	accessKey := data.AccessKey

	accessKeyModel := iam.AccessKey{
		AccessKey:           types.StringValue(accessKey.AccessKey),
		AccessKeyType:       types.StringValue(string(accessKey.AccessKeyType)),
		AccountId:           types.StringValue(accessKey.AccountId),
		CreatedAt:           types.StringValue(accessKey.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(accessKey.CreatedBy),
		Description:         types.StringPointerValue(accessKey.Description.Get()),
		ExpirationTimestamp: types.StringValue(accessKey.ExpirationTimestamp.Format(time.RFC3339)),
		Id:                  types.StringValue(accessKey.Id),
		IsEnabled:           types.BoolValue(accessKey.IsEnabled),
		ModifiedAt:          types.StringValue(accessKey.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(accessKey.ModifiedBy),
		ParentAccessKeyId:   types.StringPointerValue(accessKey.ParentAccessKeyId.Get()),
		SecretKey:           types.StringValue(accessKey.SecretKey),
	}
	accessKeyObjectValue, diags := types.ObjectValueFrom(ctx, accessKeyModel.AttributeTypes(), accessKeyModel)
	state.AccessKey = accessKeyObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *iamAccessKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state iam.AccessKeyResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Account Key
	_, err := r.client.UpdateAccessKey(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Account Key",
			"Could not update Account Key, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetAccountKey as UpdateAccountKey items are not populated.
	data, err := r.client.GetAccessKey(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading account key",
			"Could not read account ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	accessKey := data.AccessKey

	accessKeyModel := iam.AccessKey{
		AccessKey:           types.StringValue(accessKey.AccessKey),
		AccessKeyType:       types.StringValue(string(accessKey.AccessKeyType)),
		AccountId:           types.StringValue(accessKey.AccountId),
		CreatedAt:           types.StringValue(accessKey.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(accessKey.CreatedBy),
		Description:         types.StringPointerValue(accessKey.Description.Get()),
		ExpirationTimestamp: types.StringValue(accessKey.ExpirationTimestamp.Format(time.RFC3339)),
		Id:                  types.StringValue(accessKey.Id),
		IsEnabled:           types.BoolValue(accessKey.IsEnabled),
		ModifiedAt:          types.StringValue(accessKey.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(accessKey.ModifiedBy),
		ParentAccessKeyId:   types.StringPointerValue(accessKey.ParentAccessKeyId.Get()),
		SecretKey:           types.StringValue(accessKey.SecretKey),
	}
	accessKeyObjectValue, diags := types.ObjectValueFrom(ctx, accessKeyModel.AttributeTypes(), accessKeyModel)
	state.AccessKey = accessKeyObjectValue
	state.AccountId = types.StringValue(accessKey.AccountId)
	state.IsEnabled = types.BoolValue(accessKeyModel.IsEnabled.ValueBool())
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState imports the resource by ID.
func (r *iamAccessKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *iamAccessKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state iam.AccessKeyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing iam
	err := r.client.DeleteAccessKey(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting iam",
			"Could not delete iam, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
