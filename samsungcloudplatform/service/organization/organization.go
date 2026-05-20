package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &organizationResource{}
	_ resource.ResourceWithConfigure = &organizationResource{}
)

// NewOrganizationResource is a helper function to simplify the provider implementation.
func NewOrganizationResource() resource.Resource {
	return &organizationResource{}
}

// organizationResource is the resource implementation.
type organizationResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *organizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema defines the schema for the resource.
func (r *organizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Organization ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Organization Name",
				Required:    true,
			},
			"delegation_account_id": schema.StringAttribute{
				Description: "Delegation Account ID",
				Optional:    true,
				Computed:    true,
			},
			"use_scp_yn": schema.BoolAttribute{
				Description: "Use SCP",
				Optional:    true,
				Computed:    true,
			},
			// Computed fields
			"created_at": schema.StringAttribute{
				Description: "Created At",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "Created By",
				Computed:    true,
			},
			"creator_name": schema.StringAttribute{
				Description: "Creator Name",
				Computed:    true,
			},
			"modified_at": schema.StringAttribute{
				Description: "Modified At",
				Computed:    true,
			},
			"modified_by": schema.StringAttribute{
				Description: "Modified By",
				Computed:    true,
			},
			"modifier_name": schema.StringAttribute{
				Description: "Modifier Name",
				Computed:    true,
			},
			"master_account_id": schema.StringAttribute{
				Description: "Master Account ID",
				Computed:    true,
			},
			"master_account_email": schema.StringAttribute{
				Description: "Master Account Email",
				Computed:    true,
			},
			"root_unit_id": schema.StringAttribute{
				Description: "Root Unit ID",
				Computed:    true,
			},
			"srn": schema.StringAttribute{
				Description: "SRN",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *organizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.Organization
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *organizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan organization.OrganizationResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Organization
	data, err := r.client.CreateOrganization(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Organization",
			"Could not create Organization, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	org := data.Organization

	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(org.Id)
	plan.CreatedAt = types.StringValue(org.CreatedAt.Format(time.RFC3339))
	plan.CreatedBy = types.StringValue(org.CreatedBy)
	plan.CreatorName = types.StringValue(org.GetCreatorName())
	plan.ModifiedAt = types.StringValue(org.ModifiedAt.Format(time.RFC3339))
	plan.ModifiedBy = types.StringValue(org.ModifiedBy)
	plan.ModifierName = types.StringValue(org.GetModifierName())
	plan.MasterAccountId = types.StringValue(org.MasterAccountId)
	plan.MasterAccountEmail = types.StringValue(org.MasterAccountEmail)
	plan.RootUnitId = types.StringValue(org.RootUnitId)
	plan.Srn = types.StringValue(org.Srn)
	delegationAccountId := org.DelegationAccountId.Get()
	if delegationAccountId == nil || *delegationAccountId == "" {
		plan.DelegationAccountId = types.StringNull()
	} else {
		plan.DelegationAccountId = types.StringValue(*delegationAccountId)
	}
	plan.UseScpYn = types.BoolValue(org.UseScpYn)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads the resource state from Terraform state and updates the resource data.
func (r *organizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state organization.OrganizationResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed Organization value from SCP
	data, err := r.client.GetOrganization(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Organization",
			err.Error(),
		)
		return
	}

	org := data.Organization

	// Overwrite items to refresh state
	state.Id = types.StringValue(org.Id)
	state.Name = types.StringValue(org.Name)
	state.CreatedAt = types.StringValue(org.CreatedAt.Format(time.RFC3339))
	state.CreatedBy = types.StringValue(org.CreatedBy)
	state.CreatorName = types.StringValue(org.GetCreatorName())
	state.ModifiedAt = types.StringValue(org.ModifiedAt.Format(time.RFC3339))
	state.ModifiedBy = types.StringValue(org.ModifiedBy)
	state.ModifierName = types.StringValue(org.GetModifierName())
	state.MasterAccountId = types.StringValue(org.MasterAccountId)
	state.MasterAccountEmail = types.StringValue(org.MasterAccountEmail)
	state.RootUnitId = types.StringValue(org.RootUnitId)
	state.Srn = types.StringValue(org.Srn)
	state.DelegationAccountId = types.StringPointerValue(org.DelegationAccountId.Get())
	state.UseScpYn = types.BoolValue(org.UseScpYn)

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *organizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan organization.OrganizationResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Organization
	_, err := r.client.UpdateOrganization(ctx, plan.Id.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Organization",
			"Could not update Organization, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Read updated Organization
	data, err := r.client.GetOrganization(ctx, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Unable to Read Organization",
			"Could not read Organization ID "+plan.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	org := data.Organization

	// Update plan with new values
	plan.Id = types.StringValue(org.Id)
	plan.Name = types.StringValue(org.Name)
	plan.CreatedAt = types.StringValue(org.CreatedAt.Format(time.RFC3339))
	plan.CreatedBy = types.StringValue(org.CreatedBy)
	plan.CreatorName = types.StringValue(org.GetCreatorName())
	plan.ModifiedAt = types.StringValue(org.ModifiedAt.Format(time.RFC3339))
	plan.ModifiedBy = types.StringValue(org.ModifiedBy)
	plan.ModifierName = types.StringValue(org.GetModifierName())
	plan.MasterAccountId = types.StringValue(org.MasterAccountId)
	plan.MasterAccountEmail = types.StringValue(org.MasterAccountEmail)
	plan.RootUnitId = types.StringValue(org.RootUnitId)
	plan.Srn = types.StringValue(org.Srn)

	plan.UseScpYn = types.BoolValue(org.UseScpYn)

	delegationAccountId := org.DelegationAccountId.Get()
	if delegationAccountId == nil || *delegationAccountId == "" {
		plan.DelegationAccountId = types.StringNull()
	} else {
		plan.DelegationAccountId = types.StringValue(*delegationAccountId)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *organizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state organization.OrganizationResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Organization
	_, err := r.client.DeleteOrganization(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting Organization",
			"Could not delete Organization, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
