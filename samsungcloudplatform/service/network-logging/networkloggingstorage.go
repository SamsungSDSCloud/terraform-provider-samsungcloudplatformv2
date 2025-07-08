package network_logging

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/networklogging"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &networkLoggingNetworkLoggingStorageResource{}
	_ resource.ResourceWithConfigure = &networkLoggingNetworkLoggingStorageResource{}
)

// NewNetworkLoggingNetworkLoggingStorageResource is a helper function to simplify the provider implementation.
func NewNetworkLoggingNetworkLoggingStorageResource() resource.Resource {
	return &networkLoggingNetworkLoggingStorageResource{}
}

// networkLoggingNetworkLoggingStorageResource is the data source implementation.
type networkLoggingNetworkLoggingStorageResource struct {
	config  *scpsdk.Configuration
	client  *networklogging.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *networkLoggingNetworkLoggingStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_logging_network_logging_storage"
}

// Schema defines the schema for the data source.
func (r *networkLoggingNetworkLoggingStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network logging storage",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "AccountId",
				Computed:    true,
			},
			common.ToSnakeCase("ResourceType"): schema.StringAttribute{
				Description: "ResourceType \n" +
					"  - example : FIREWALL | SECURITY_GROUP | NAT",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("FIREWALL", "SECURITY_GROUP", "NAT"),
				},
			},
			common.ToSnakeCase("BucketName"): schema.StringAttribute{
				Description: "BucketName",
				Required:    true,
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
	}
}

// Configure adds the provider configured client to the data source.
func (r *networkLoggingNetworkLoggingStorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.NetworkLogging
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *networkLoggingNetworkLoggingStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan networklogging.NetworkLoggingStorageResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new network logging storage
	data, err := r.client.CreateNetworkLoggingStorage(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating network logging storage",
			"Could not create network logging storage, unexpected error: "+err.Error(),
		)
		return
	}

	networkLoggingStorage := data.NetworkLoggingStorage

	plan.Id = types.StringValue(networkLoggingStorage.Id)
	plan.AccountId = types.StringValue(networkLoggingStorage.AccountId)
	plan.ResourceType = types.StringValue(string(networkLoggingStorage.ResourceType))
	plan.BucketName = types.StringValue(networkLoggingStorage.BucketName)
	plan.CreatedAt = types.StringValue(networkLoggingStorage.CreatedAt.Format(time.RFC3339))
	plan.CreatedBy = types.StringValue(networkLoggingStorage.CreatedBy)
	plan.ModifiedAt = types.StringValue(networkLoggingStorage.ModifiedAt.Format(time.RFC3339))
	plan.ModifiedBy = types.StringValue(networkLoggingStorage.ModifiedBy)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *networkLoggingNetworkLoggingStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkLoggingNetworkLoggingStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkLoggingNetworkLoggingStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state networklogging.NetworkLoggingStorageResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing network logging storage
	err := r.client.DeleteNetworkLoggingStorage(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting network logging storage",
			"Could not delete network logging storage, unexpected error: "+err.Error(),
		)
		return
	}
}
