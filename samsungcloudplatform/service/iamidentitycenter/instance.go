package iamidentitycenter

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	iamidentitycenterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/iamidentitycenter"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/iam-identity-center/1.6"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &iamIdentityCenterInstanceResource{}
	_ resource.ResourceWithConfigure   = &iamIdentityCenterInstanceResource{}
	_ resource.ResourceWithImportState = &iamIdentityCenterInstanceResource{}
)

// NewIamIdentityCenterInstanceResource is a helper function to simplify the provider implementation.
func NewIamIdentityCenterInstanceResource() resource.Resource {
	return &iamIdentityCenterInstanceResource{}
}

// iamIdentityCenterInstanceResource is the data source implementation.
type iamIdentityCenterInstanceResource struct {
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *iamIdentityCenterInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_instance"
}

func (r *iamIdentityCenterInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM Identity Center Instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{},
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Instance Name\n  - example: My Instance\n  - maxLength: 128",
				MarkdownDescription: "Instance Name\n  - example: My Instance\n  - maxLength: 128",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Instance Description\n  - example: My Instance Description",
				MarkdownDescription: "Instance Description\n  - example: My Instance Description",
			},
			"identity_store_type": schema.StringAttribute{
				Required: true,
				Description: "- enum: [\"IDENTITY_CENTER_DIRECTORY\",\"ACTIVE_DIRECTORY\",\"OPENLDAP\",\"DIRECTORY_389DS\"]\n" +
					"  - example: IDENTITY_CENTER_DIRECTORY",
				MarkdownDescription: "- enum: [\"IDENTITY_CENTER_DIRECTORY\",\"ACTIVE_DIRECTORY\",\"OPENLDAP\",\"DIRECTORY_389DS\"]\n" +
					"  - example: IDENTITY_CENTER_DIRECTORY",
			},
			"identity_store_config": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Description: "Identity Store Config\n" +
					"  - example: map[bind_credential:password bind_dn:CN=Administrator,CN=Users,DC=test,DC=add connection_url:ldaps://test.subdomain.company.com:123]",
				MarkdownDescription: "Identity Store Config\n" +
					"  - example: map[bind_credential:password bind_dn:CN=Administrator,CN=Users,DC=test,DC=add connection_url:ldaps://test.subdomain.company.com:123]",
				Attributes: map[string]schema.Attribute{
					"bind_credential": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Bind Password\n  - maxLength: 128",
					},
					"bind_dn": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Bind DN\n  - maxLength: 255",
					},
					"connection_url": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "LDAP Server URL\n  - maxLength: 255",
					},
					"rdn_ldap_attribute": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "RDN LDAP Attribute\n  - maxLength: 64",
					},
					"user_object_classes": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "User Object Classes\n  - maxLength: 100",
					},
					"username_ldap_attribute": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Username LDAP Attribute\n  - maxLength: 64",
					},
					"users_dn": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Users DN\n  - maxLength: 255",
					},
				},
			},
			"self_managed_password": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Self Managed Password",
			},
			"instance": schema.SingleNestedAttribute{
				Computed: true,
				Description: "A detail of Instance.\n" +
					"  - example: '{account_id: 1234ab567cd64769e8f9g490hi304891, created_at: 2024-05-17T00:23:17Z, ...}'",
				MarkdownDescription: "A detail of Instance.\n" +
					"  - example: '{account_id: 1234ab567cd64769e8f9g490hi304891, created_at: 2024-05-17T00:23:17Z, ...}'",
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed: true,
						Description: "Account ID\n" +
							"  - example: 1234ab567cd64769e8f9g490hi304891",
					},
					"created_at": schema.StringAttribute{
						Computed: true,
						Description: "Created At\n" +
							"  - example: 2024-05-17T00:23:17Z",
					},
					"created_by": schema.StringAttribute{
						Computed: true,
						Description: "Created By\n" +
							"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
					},
					"creator_name": schema.StringAttribute{
						Computed: true,
						Description: "Creator Name\n" +
							"  - example: John Doe",
					},
					"delegated_account_id": schema.StringAttribute{
						Computed: true,
						Description: "Delegated Account ID\n" +
							"  - example: acc-2353",
					},
					"delegated_at": schema.StringAttribute{
						Computed: true,
						Description: "Delegation Date\n" +
							"  - example: 2023-10-15 14:30:00",
					},
					"delegated_by": schema.StringAttribute{
						Computed: true,
						Description: "Delegated By\n" +
							"  - example: sef2w3c8c29a449dbfa8681f8f1d78e2",
					},
					"description": schema.StringAttribute{
						Computed: true,
						Description: "Instance Description\n" +
							"  - example: My Instance Description",
					},
					"id": schema.StringAttribute{
						Computed: true,
						Description: "Instance ID\n" +
							"  - example: ssoins-12345",
					},
					"identity_store_id": schema.StringAttribute{
						Computed: true,
						Description: "Identity Store ID\n" +
							"  - example: d-12345",
					},
					"identity_store_type": schema.StringAttribute{
						Computed:    true,
						Description: "- enum: [\"IDENTITY_CENTER_DIRECTORY\",\"ACTIVE_DIRECTORY\",\"OPENLDAP\",\"DIRECTORY_389DS\"]",
					},
					"last_sync_at": schema.StringAttribute{
						Computed: true,
						Description: "Last Sync Date\n" +
							"  - example: 2025-08-06 10:30:00",
					},
					"last_sync_status": schema.StringAttribute{
						Computed:    true,
						Description: "- enum: [\"CREATED\",\"RUNNING\",\"FINISHED\",\"FAILED\",\"NONE\",\"ABORTED\"]",
					},
					"master_account_id": schema.StringAttribute{
						Computed: true,
						Description: "Master Account ID\n" +
							"  - example: acc-12345",
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
						Description: "Modified At\n" +
							"  - example: 2024-05-17T00:23:17Z",
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
						Description: "Modified By\n" +
							"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
					},
					"modifier_name": schema.StringAttribute{
						Computed: true,
						Description: "Modifier Name\n" +
							"  - example: Smith",
					},
					"name": schema.StringAttribute{
						Computed: true,
						Description: "Instance Name\n" +
							"  - example: My Instance",
					},
					"organization_id": schema.StringAttribute{
						Computed: true,
						Description: "Organization ID\n" +
							"  - example: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5",
					},
					"port": schema.Int64Attribute{
						Computed: true,
						Description: "Port\n" +
							"  - example: 443",
					},
					"public_ip": schema.StringAttribute{
						Computed:    true,
						Description: "Public IP",
					},
					"region": schema.StringAttribute{
						Computed: true,
						Description: "Region\n" +
							"  - example: kr-west1",
					},
					"resource_name": schema.StringAttribute{
						Computed: true,
						Description: "Resource Name\n" +
							"  - example: resource-name-1",
					},
					"resource_type": schema.StringAttribute{
						Computed: true,
						Description: "Resource Type [instance|permission-set]\n" +
							"  - example: virtual-server",
					},
					"resource_type_display_name": schema.StringAttribute{
						Computed: true,
						Description: "Resource Type Display Name\n" +
							"  - example: Virtual Server",
					},
					"self_managed_password": schema.BoolAttribute{
						Computed:    true,
						Description: "Self Managed Password",
					},
					"service": schema.StringAttribute{
						Computed: true,
						Description: "Service\n" +
							"  - example: virtualserver",
					},
					"service_name": schema.StringAttribute{
						Computed: true,
						Description: "Service Name\n" +
							"  - example: Virtual Server",
					},
					"srn": schema.StringAttribute{
						Computed: true,
						Description: "Instance SRN\n" +
							"  - example: srn:identity-center::prj-01234:kr-west-1::instance/ins-01234",
					},
					"sso_start_url": schema.StringAttribute{
						Computed: true,
						Description: "SSO Start URL\n" +
							"  - example: https://login.e.samsungsdscloud.com/d-12345",
					},
					"state": schema.StringAttribute{
						Computed:    true,
						Description: "- enum: [\"ACTIVE\",\"INACTIVE\"]",
					},
				},
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Region",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *iamIdentityCenterInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.clients = inst.Client
}

func waitForInstanceStatus(ctx context.Context, iamClient *iamidentitycenterClient.Client, instanceId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, _, err := iamClient.GetInstance(ctx, instanceId)
		if err != nil {
			return nil, "", err
		}

		var stateStr string
		if info.Instance.State != nil {
			stateStr = string(*info.Instance.State)
		} else {
			stateStr = "UNKNOWN"
		}

		return info, stateStr, nil
	}, -1, -1, -1, -1)
}

func waitForInstanceDeletion(ctx context.Context, iamClient *iamidentitycenterClient.Client, instanceId string) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for instance %s to be deleted", instanceId)
		case <-ticker.C:
			_, httpResp, err := iamClient.GetInstance(ctx, instanceId)
			if err != nil {
				if httpResp != nil && httpResp.StatusCode == 404 {
					return nil
				}
				return err
			}
		}
	}
}

func (r *iamIdentityCenterInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan instanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	createReq := iamidentitycenterClient.InstanceCreateRequest{
		Name:              plan.Name.ValueString(),
		IdentityStoreType: plan.IdentityStoreType.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createReq.Description = plan.Description.ValueString()
	}

	if !plan.SelfManagedPassword.IsNull() && !plan.SelfManagedPassword.IsUnknown() {
		v := plan.SelfManagedPassword.ValueBool()
		createReq.SelfManagedPassword = &v
	}

	// Handle identity store config
	if !plan.IdentityStoreConfig.IsNull() && !plan.IdentityStoreConfig.IsUnknown() {
		var config planmodel
		diags := plan.IdentityStoreConfig.As(ctx, &config, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		createReq.SetIdentityStoreConfig(iamidentitycenterClient.LdapConfigOptions{
			BindCredential:        stringToPointer(config.BindCredential.ValueString()),
			BindDn:                stringToPointer(config.BindDn.ValueString()),
			ConnectionUrl:         stringToPointer(config.ConnectionUrl.ValueString()),
			RdnLdapAttribute:      stringToPointer(config.RdnLdapAttribute.ValueString()),
			UserObjectClasses:     stringToPointer(config.UserObjectClasses.ValueString()),
			UsernameLdapAttribute: stringToPointer(config.UsernameLdapAttribute.ValueString()),
			UsersDn:               stringToPointer(config.UsersDn.ValueString()),
		})
	}

	result, err := r.clients.IamIdentityCenter.CreateInstance(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create IAM Identity Center Instance",
			err.Error(),
		)
		return
	}

	instanceId := result.Instance.Id

	instanceObj, _ := types.ObjectValueFrom(ctx, instanceDetailModel{}.AttributeTypes(), instanceDetailModel{})

	state := instanceResourceModel{
		ID:                  types.StringValue(instanceId),
		Name:                plan.Name,
		IdentityStoreType:   plan.IdentityStoreType,
		IdentityStoreConfig: plan.IdentityStoreConfig,
		SelfManagedPassword: plan.SelfManagedPassword,
		Region:              plan.Region,
		Instance:            instanceObj,
	}
	if !plan.Description.IsNull() {
		state.Description = plan.Description
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := waitForInstanceStatus(ctx, r.clients.IamIdentityCenter, instanceId, []string{}, []string{"ACTIVE", "INACTIVE"}); err != nil {
		resp.Diagnostics.AddWarning(
			"Warning: Instance created but status wait failed",
			"Instance was created but failed to wait for status: "+err.Error()+". State may be stale.",
		)
	}

	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

// Read refreshes the Terraform state with the latest data.
func (r *iamIdentityCenterInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state instanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := instanceIdFromModel(&state)

	result, httpResp, err := r.clients.IamIdentityCenter.GetInstance(ctx, instanceId)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Instance",
			err.Error(),
		)
		return
	}

	var newState instanceResourceModel
	diags = convertResponseToModel(ctx, result, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Keep the ID from the state; preserve config-driven name/identity_store_type
	// from the previous state, but fall back to the server values when the state
	// lacks them (e.g. right after import), so a partial PATCH update is possible.
	newState.ID = state.ID
	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		newState.Name = state.Name
	}
	if !state.IdentityStoreType.IsNull() && !state.IdentityStoreType.IsUnknown() {
		newState.IdentityStoreType = state.IdentityStoreType
	}

	if !state.IdentityStoreConfig.IsNull() && !state.IdentityStoreConfig.IsUnknown() {
		newState.IdentityStoreConfig = state.IdentityStoreConfig
	}

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

func instanceIdFromModel(m *instanceResourceModel) string {
	instanceId := m.ID.ValueString()
	if instanceId == "" {
		if !m.Instance.IsNull() && !m.Instance.IsUnknown() {
			if idAttr, ok := m.Instance.Attributes()["id"]; ok && !idAttr.IsNull() && !idAttr.IsUnknown() {
				instanceId = idAttr.(types.String).ValueString()
			}
		}
	}
	return instanceId
}

func resolveInstanceId(state, plan *instanceResourceModel) string {
	instanceId := instanceIdFromModel(state)
	if instanceId == "" {
		instanceId = instanceIdFromModel(plan)
	}
	return instanceId
}

func setIdentityStoreConfig(ctx context.Context, plan *instanceResourceModel, updateReq *iamidentitycenterClient.InstanceUpdateRequest) diag.Diagnostics {
	var config planmodel
	diags := plan.IdentityStoreConfig.As(ctx, &config, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return diags
	}
	updateReq.SetIdentityStoreConfig(iamidentitycenterClient.LdapConfigOptions{
		BindCredential:        stringToPointer(config.BindCredential.ValueString()),
		BindDn:                stringToPointer(config.BindDn.ValueString()),
		ConnectionUrl:         stringToPointer(config.ConnectionUrl.ValueString()),
		RdnLdapAttribute:      stringToPointer(config.RdnLdapAttribute.ValueString()),
		UserObjectClasses:     stringToPointer(config.UserObjectClasses.ValueString()),
		UsernameLdapAttribute: stringToPointer(config.UsernameLdapAttribute.ValueString()),
		UsersDn:               stringToPointer(config.UsersDn.ValueString()),
	})
	return nil
}

func mergeStateFromPlan(newState, plan *instanceResourceModel) {
	if !plan.ID.IsUnknown() {
		newState.ID = plan.ID
	}
	if !plan.Name.IsUnknown() {
		newState.Name = plan.Name
	}
	if !plan.IdentityStoreType.IsUnknown() {
		newState.IdentityStoreType = plan.IdentityStoreType
	}
	if !plan.Description.IsUnknown() {
		newState.Description = plan.Description
	}
	if !plan.IdentityStoreConfig.IsNull() && !plan.IdentityStoreConfig.IsUnknown() {
		newState.IdentityStoreConfig = plan.IdentityStoreConfig
	}
	if !plan.SelfManagedPassword.IsNull() && !plan.SelfManagedPassword.IsUnknown() {
		newState.SelfManagedPassword = plan.SelfManagedPassword
	}
	if !plan.Region.IsUnknown() {
		newState.Region = plan.Region
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *iamIdentityCenterInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan instanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get state to retrieve the instance ID from previous state
	var state instanceResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := resolveInstanceId(&state, &plan)

	updateReq := iamidentitycenterClient.InstanceUpdateRequest{}
	// SetInstance is a PATCH. Sending identity_store_type unconditionally makes
	// the backend 400 for LDAP-family stores when identity_store_config is
	// omitted, so store-type fields are sent only when the value changed.
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updateReq.Description = plan.Description.ValueString()
	}
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updateReq.Name = plan.Name.ValueString()
	}
	if !plan.IdentityStoreType.IsNull() && !plan.IdentityStoreType.IsUnknown() &&
		(state.IdentityStoreType.IsNull() || state.IdentityStoreType.IsUnknown() ||
			plan.IdentityStoreType.ValueString() != state.IdentityStoreType.ValueString()) {
		updateReq.IdentityStoreType = plan.IdentityStoreType.ValueString()
	}
	if !plan.SelfManagedPassword.IsNull() && !plan.SelfManagedPassword.IsUnknown() &&
		(state.SelfManagedPassword.IsNull() || state.SelfManagedPassword.IsUnknown() ||
			plan.SelfManagedPassword.ValueBool() != state.SelfManagedPassword.ValueBool()) {
		v := plan.SelfManagedPassword.ValueBool()
		updateReq.SelfManagedPassword = &v
	}
	if !plan.IdentityStoreConfig.IsNull() && !plan.IdentityStoreConfig.IsUnknown() &&
		!identityStoreConfigEqual(state.IdentityStoreConfig, plan.IdentityStoreConfig) {
		diags = setIdentityStoreConfig(ctx, &plan, &updateReq)
		resp.Diagnostics.Append(diags...)
	}

	result, err := r.clients.IamIdentityCenter.UpdateInstance(ctx, instanceId, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update IAM Identity Center Instance",
			err.Error(),
		)
		return
	}

	var newState instanceResourceModel
	diags = convertResponseToModel(ctx, result, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeStateFromPlan(&newState, &plan)

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}

// identityStoreConfigEqual reports whether two identity_store_config object
// values are semantically equal, so that IDC instance PATCH updates only send
// the store configuration when it actually changed.
func identityStoreConfigEqual(stateObj, planObj types.Object) bool {
	if stateObj.IsNull() || stateObj.IsUnknown() {
		return planObj.IsNull() || planObj.IsUnknown()
	}
	if planObj.IsNull() || planObj.IsUnknown() {
		return false
	}
	var s, p map[string]interface{}
	if err := stateObj.As(context.TODO(), &s, basetypes.ObjectAsOptions{}); err != nil {
		return false
	}
	if err := planObj.As(context.TODO(), &p, basetypes.ObjectAsOptions{}); err != nil {
		return false
	}
	return reflect.DeepEqual(s, p)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *iamIdentityCenterInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state instanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := instanceIdFromModel(&state)

	// Delete the instance
	httpResp, err := r.clients.IamIdentityCenter.DeleteInstance(ctx, instanceId)
	if err != nil {
		// If 409 Conflict, instance might have dependent resources or deletion is restricted
		// In such cases, we log a warning but consider the deletion successful from Terraform's perspective
		// because Terraform state should be cleared regardless
		if httpResp != nil && httpResp.StatusCode == 409 {
			resp.Diagnostics.AddWarning(
				"IAM Identity Center Instance deletion restricted",
				"Instance deletion returned 409 Conflict. This may indicate dependent resources exist or deletion is restricted. Instance will be removed from Terraform state.",
			)
			// Clear the state anyway - Terraform will no longer manage this resource
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to delete IAM Identity Center Instance",
			err.Error(),
		)
		return
	}

	// Wait for the instance to be fully deleted (until GetInstance returns 404)
	if err := waitForInstanceDeletion(ctx, r.clients.IamIdentityCenter, instanceId); err != nil {
		resp.Diagnostics.AddWarning(
			"Warning: Instance deleted but failed to verify deletion",
			err.Error(),
		)
	}
}

// ImportState imports the resource into Terraform state.
func (r *iamIdentityCenterInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper types and functions

type instanceResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	IdentityStoreType   types.String `tfsdk:"identity_store_type"`
	IdentityStoreConfig types.Object `tfsdk:"identity_store_config"`
	SelfManagedPassword types.Bool   `tfsdk:"self_managed_password"`
	Instance            types.Object `tfsdk:"instance"`
	Region              types.String `tfsdk:"region"`
}

type instanceDetailModel struct {
	AccountId               types.String `tfsdk:"account_id"`
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	CreatorName             types.String `tfsdk:"creator_name"`
	DelegatedAccountId      types.String `tfsdk:"delegated_account_id"`
	DelegatedAt             types.String `tfsdk:"delegated_at"`
	DelegatedBy             types.String `tfsdk:"delegated_by"`
	Description             types.String `tfsdk:"description"`
	ID                      types.String `tfsdk:"id"`
	IdentityStoreId         types.String `tfsdk:"identity_store_id"`
	IdentityStoreType       types.String `tfsdk:"identity_store_type"`
	LastSyncAt              types.String `tfsdk:"last_sync_at"`
	LastSyncStatus          types.String `tfsdk:"last_sync_status"`
	MasterAccountId         types.String `tfsdk:"master_account_id"`
	ModifiedAt              types.String `tfsdk:"modified_at"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	ModifierName            types.String `tfsdk:"modifier_name"`
	Name                    types.String `tfsdk:"name"`
	OrganizationId          types.String `tfsdk:"organization_id"`
	Port                    types.Int64  `tfsdk:"port"`
	PublicIp                types.String `tfsdk:"public_ip"`
	Region                  types.String `tfsdk:"region"`
	ResourceName            types.String `tfsdk:"resource_name"`
	ResourceType            types.String `tfsdk:"resource_type"`
	ResourceTypeDisplayName types.String `tfsdk:"resource_type_display_name"`
	SelfManagedPassword     types.Bool   `tfsdk:"self_managed_password"`
	Service                 types.String `tfsdk:"service"`
	ServiceName             types.String `tfsdk:"service_name"`
	Srn                     types.String `tfsdk:"srn"`
	SsoStartUrl             types.String `tfsdk:"sso_start_url"`
	State                   types.String `tfsdk:"state"`
}

func (i instanceDetailModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":                 types.StringType,
		"created_at":                 types.StringType,
		"created_by":                 types.StringType,
		"creator_name":               types.StringType,
		"delegated_account_id":       types.StringType,
		"delegated_at":               types.StringType,
		"delegated_by":               types.StringType,
		"description":                types.StringType,
		"id":                         types.StringType,
		"identity_store_id":          types.StringType,
		"identity_store_type":        types.StringType,
		"last_sync_at":               types.StringType,
		"last_sync_status":           types.StringType,
		"master_account_id":          types.StringType,
		"modified_at":                types.StringType,
		"modified_by":                types.StringType,
		"modifier_name":              types.StringType,
		"name":                       types.StringType,
		"organization_id":            types.StringType,
		"port":                       types.Int64Type,
		"public_ip":                  types.StringType,
		"region":                     types.StringType,
		"resource_name":              types.StringType,
		"resource_type":              types.StringType,
		"resource_type_display_name": types.StringType,
		"self_managed_password":      types.BoolType,
		"service":                    types.StringType,
		"service_name":               types.StringType,
		"srn":                        types.StringType,
		"sso_start_url":              types.StringType,
		"state":                      types.StringType,
	}
}

type planmodel struct {
	BindCredential        types.String `tfsdk:"bind_credential"`
	BindDn                types.String `tfsdk:"bind_dn"`
	ConnectionUrl         types.String `tfsdk:"connection_url"`
	RdnLdapAttribute      types.String `tfsdk:"rdn_ldap_attribute"`
	UserObjectClasses     types.String `tfsdk:"user_object_classes"`
	UsernameLdapAttribute types.String `tfsdk:"username_ldap_attribute"`
	UsersDn               types.String `tfsdk:"users_dn"`
}

func convertResponseToModel(ctx context.Context, result interface{}, model *instanceResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Handle different response types
	var instanceDetail *sdk.InstanceDetailAddResourceV1Dot6
	switch r := result.(type) {
	case *sdk.InstanceShowResponseAddResourceV1Dot6:
		instanceDetail = &r.Instance
	case *sdk.InstanceShowResponseV1Dot5:
		// For create/update responses, convert V1Dot5 to AddResourceV1Dot6 format
		detail := r.Instance
		instanceDetail = &sdk.InstanceDetailAddResourceV1Dot6{
			CreatedAt:           detail.CreatedAt,
			CreatedBy:           detail.CreatedBy,
			DelegatedAccountId:  detail.DelegatedAccountId,
			DelegatedAt:         detail.DelegatedAt,
			DelegatedBy:         detail.DelegatedBy,
			Description:         detail.Description,
			Id:                  detail.Id,
			IdentityStoreConfig: detail.IdentityStoreConfig,
			IdentityStoreId:     detail.IdentityStoreId,
			IdentityStoreType:   detail.IdentityStoreType,
			LastSyncAt:          detail.LastSyncAt,
			LastSyncStatus:      detail.LastSyncStatus,
			MasterAccountId:     detail.MasterAccountId,
			ModifiedAt:          detail.ModifiedAt,
			ModifiedBy:          detail.ModifiedBy,
			Name:                detail.Name,
			OrganizationId:      detail.OrganizationId,
			Port:                detail.Port,
			PublicIp:            detail.PublicIp,
			Region:              detail.Region,
			SelfManagedPassword: detail.SelfManagedPassword,
			Srn:                 detail.Srn,
			SsoStartUrl:         detail.SsoStartUrl,
			State:               detail.State,
		}
	default:
		diags.AddError("Unsupported response type", fmt.Sprintf("Got type %T", result))
		return diags
	}

	// Convert instance detail to model
	instanceModel := instanceDetailModel{
		CreatedAt:         types.StringValue(instanceDetail.CreatedAt.Format(time.RFC3339)),
		CreatedBy:         types.StringValue(instanceDetail.CreatedBy),
		Description:       toStringValue(instanceDetail.Description),
		ID:                types.StringValue(instanceDetail.Id),
		IdentityStoreId:   types.StringValue(instanceDetail.IdentityStoreId),
		IdentityStoreType: toStringValueFromEnum(instanceDetail.IdentityStoreType),
		MasterAccountId:   types.StringValue(instanceDetail.MasterAccountId),
		ModifiedAt:        types.StringValue(instanceDetail.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:        types.StringValue(instanceDetail.ModifiedBy),
		Name:              types.StringValue(instanceDetail.Name),
		OrganizationId:    types.StringValue(instanceDetail.OrganizationId),
		Region:            types.StringValue(instanceDetail.Region),
		Srn:               types.StringValue(instanceDetail.Srn),
		SsoStartUrl:       types.StringValue(instanceDetail.SsoStartUrl),
		State:             toStringValueFromEnum(instanceDetail.State),
	}

	// Optional fields
	instanceModel.AccountId = toStringValue(instanceDetail.AccountId)
	instanceModel.CreatorName = toStringValue(instanceDetail.CreatorName)
	instanceModel.DelegatedAccountId = toStringValue(instanceDetail.DelegatedAccountId)
	instanceModel.DelegatedAt = toTimeStringValue(instanceDetail.DelegatedAt)
	instanceModel.DelegatedBy = toStringValue(instanceDetail.DelegatedBy)
	instanceModel.LastSyncAt = toStringValue(instanceDetail.LastSyncAt)
	instanceModel.LastSyncStatus = toStringValueFromEnum(instanceDetail.LastSyncStatus)
	instanceModel.ModifierName = toStringValue(instanceDetail.ModifierName)
	instanceModel.Port = toInt64Value(instanceDetail.Port)
	instanceModel.PublicIp = toStringValue(instanceDetail.PublicIp)
	instanceModel.ResourceName = toStringValue(instanceDetail.ResourceName)
	instanceModel.ResourceType = toStringValue(instanceDetail.ResourceType)
	instanceModel.ResourceTypeDisplayName = toStringValue(instanceDetail.ResourceTypeDisplayName)
	instanceModel.SelfManagedPassword = toBoolValue(instanceDetail.SelfManagedPassword)
	instanceModel.Service = toStringValue(instanceDetail.Service)
	instanceModel.ServiceName = toStringValue(instanceDetail.ServiceName)

	instanceObj, objDiags := types.ObjectValueFrom(ctx, instanceModel.AttributeTypes(), instanceModel)
	diags.Append(objDiags...)
	model.Instance = instanceObj
	model.ID = types.StringValue(instanceDetail.Id)
	model.Name = types.StringValue(instanceDetail.Name)
	model.Region = types.StringValue(instanceDetail.Region)
	model.Description = toStringValue(instanceDetail.Description)
	model.IdentityStoreType = toStringValueFromEnum(instanceDetail.IdentityStoreType)
	model.SelfManagedPassword = toBoolValue(instanceDetail.SelfManagedPassword)

	// Convert identity store config
	identityStoreConfigType := planmodel{
		BindCredential:        types.StringNull(),
		BindDn:                types.StringNull(),
		ConnectionUrl:         types.StringNull(),
		RdnLdapAttribute:      types.StringNull(),
		UserObjectClasses:     types.StringNull(),
		UsernameLdapAttribute: types.StringNull(),
		UsersDn:               types.StringNull(),
	}

	if instanceDetail.IdentityStoreConfig.IsSet() {
		config := instanceDetail.IdentityStoreConfig.Get()

		if config != nil {
			identityStoreConfigType.BindCredential = toStringValueFromPtr(config.BindCredential)
			identityStoreConfigType.BindDn = toStringValueFromPtr(config.BindDn)
			identityStoreConfigType.ConnectionUrl = toStringValueFromPtr(config.ConnectionUrl)
			identityStoreConfigType.RdnLdapAttribute = toStringValueFromPtr(config.RdnLdapAttribute)
			identityStoreConfigType.UserObjectClasses = toStringValueFromPtr(config.UserObjectClasses)
			identityStoreConfigType.UsernameLdapAttribute = toStringValueFromPtr(config.UsernameLdapAttribute)
			identityStoreConfigType.UsersDn = toStringValueFromPtr(config.UsersDn)
		}
	}

	identityStoreConfigObj, cfgDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"bind_credential":         types.StringType,
		"bind_dn":                 types.StringType,
		"connection_url":          types.StringType,
		"rdn_ldap_attribute":      types.StringType,
		"user_object_classes":     types.StringType,
		"username_ldap_attribute": types.StringType,
		"users_dn":                types.StringType,
	}, identityStoreConfigType)
	diags.Append(cfgDiags...)
	model.IdentityStoreConfig = identityStoreConfigObj

	return diags
}

func stringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func toStringValueFromPtr(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

// toStringValue handles sdk.NullableString type
func toStringValue(ns sdk.NullableString) types.String {
	if !ns.IsSet() {
		return types.StringNull()
	}
	v := ns.Get()
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

func toStringValueFromEnum(v any) types.String {
	if v == nil {
		return types.StringNull()
	}
	switch val := v.(type) {
	case *string:
		if val == nil {
			return types.StringNull()
		}
		return types.StringValue(*val)
	case string:
		return types.StringValue(val)
	case *sdk.InstanceStateEnum:
		if val == nil {
			return types.StringNull()
		}
		return types.StringValue(string(*val))
	case *sdk.InstanceStoreTypeEnum:
		if val == nil {
			return types.StringNull()
		}
		return types.StringValue(string(*val))
	case sdk.NullableAsyncStatusEnum:
		if val.IsSet() && val.Get() != nil {
			return types.StringValue(string(*val.Get()))
		}
		return types.StringNull()
	default:
		return types.StringValue(fmt.Sprintf("%v", v))
	}
}

func toInt64ValueFromPtr(v *int32) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
}

func toTimeStringValue(ns sdk.NullableTime) types.String {
	if !ns.IsSet() {
		return types.StringNull()
	}
	if ns.Get() == nil {
		return types.StringNull()
	}
	return types.StringValue(ns.Get().Format(time.RFC3339))
}

func toInt64Value(ns sdk.NullableInt32) types.Int64 {
	if !ns.IsSet() {
		return types.Int64Null()
	}
	v := ns.Get()
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
}

func toBoolValue(ns sdk.NullableBool) types.Bool {
	if !ns.IsSet() {
		return types.BoolNull()
	}
	v := ns.Get()
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*v)
}
