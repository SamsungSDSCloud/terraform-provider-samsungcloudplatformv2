package iam

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                = &iamRolePolicyBindingsResource{}
	_ resource.ResourceWithConfigure   = &iamRolePolicyBindingsResource{}
	_ resource.ResourceWithImportState = &iamRolePolicyBindingsResource{}
)

func NewIamRolePolicyBindingsResource() resource.Resource {
	return &iamRolePolicyBindingsResource{}
}

type iamRolePolicyBindingsResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (r *iamRolePolicyBindingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role_policy_bindings"
}

func (r *iamRolePolicyBindingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamRolePolicyBindingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Role Policy Bindings",
		Attributes: map[string]schema.Attribute{
			"role_id": schema.StringAttribute{
				Optional: true,
				Description: "ID of the role to bind policies to.\n" +
					"  - example : 'role-12345678'",
			},
			"policy_ids": schema.ListAttribute{
				Optional: true,
				Description: "List of policy IDs to bind to the role.\n" +
					"  - example : ['policy-12345678', 'policy-87654321']",
				ElementType: types.StringType,
			},
			"role_policy_bindings": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of policy bindings attached to the role.\n" +
					"  - example : '[{account_id: 123456789012, created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, creator_email: user@example.com, ...}]'",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed:    true,
							Description: "Account ID of the account that owns the role policy binding.\n  - example : '123456789012'",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the role policy binding was created.\n  - example : '2024-01-01T00:00:00Z'",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "User who created the role policy binding.\n  - example : 'user@example.com'",
						},
						"creator_email": schema.StringAttribute{
							Computed:    true,
							Description: "Email address of the user who created the role policy binding.\n  - example : 'user@example.com'",
						},
						"creator_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the user who created the role policy binding.\n  - example : 'John Doe'",
						},
						"default_version_id": schema.StringAttribute{
							Computed:    true,
							Description: "Default version ID of the policy.\n  - example : 'pol-1234567890abcdef'",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable description of the role policy binding.\n  - example : 'My policy description'",
						},
						"domain_name": schema.StringAttribute{
							Computed:    true,
							Description: "Domain name associated with the role policy binding.\n  - example : 'scp'",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier of the role policy binding.\n  - example : 'pol-1234567890abcdef'",
						},
						"modified_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the role policy binding was last modified.\n  - example : '2024-01-01T00:00:00Z'",
						},
						"modified_by": schema.StringAttribute{
							Computed:    true,
							Description: "User who last modified the role policy binding.\n  - example : 'user@example.com'",
						},
						"modifier_email": schema.StringAttribute{
							Computed:    true,
							Description: "Email address of the user who last modified the role policy binding.\n  - example : 'user@example.com'",
						},
						"modifier_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the user who last modified the role policy binding.\n  - example : 'John Doe'",
						},
						"policy_category": schema.StringAttribute{
							Computed:    true,
							Description: "Category of the policy (e.g., IDENTITY_BASED or RESOURCE_BASED).\n  - example : 'IDENTITY_BASED' | 'RESOURCE_BASED'",
						},
						"policy_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the policy.\n  - example : 'MyPolicy'",
						},
						"policy_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the policy.\n  - example : 'USER_DEFINED'",
						},
						"policy_versions": schema.ListNestedAttribute{
							Optional: true,
							Description: "List of versions associated with the policy.\n" +
								"  - example : '[{created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229dev1f, id: pol-1234567890abcdef, modified_at: 2024-05-17T00:23:17Z, modified_by: ef50cdc207f05f6fb8f20219f229ed1f, ...}]'",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"created_at": schema.StringAttribute{
										Computed:    true,
										Description: "Timestamp when the policy version was created.\n  - example : '2024-01-01T00:00:00Z'",
									},
									"created_by": schema.StringAttribute{
										Computed:    true,
										Description: "User who created the policy version.\n  - example : 'user@example.com'",
									},
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "Unique identifier of the policy version.\n  - example : 'pol-1234567890abcdef'",
									},
									"modified_at": schema.StringAttribute{
										Computed:    true,
										Description: "Timestamp when the policy version was last modified.\n  - example : '2024-01-01T00:00:00Z'",
									},
									"modified_by": schema.StringAttribute{
										Computed:    true,
										Description: "User who last modified the policy version.\n  - example : 'user@example.com'",
									},
									"policy_document": schema.SingleNestedAttribute{
										Computed: true,
										Description: "The policy document containing permission definitions for this policy version.\n" +
											"  - example : '{statement: [{action: [iam:CreateRole], effect: Allow, resource: [*], ...}], version: 2024-07-01}'",
										Attributes: map[string]schema.Attribute{
											"statement": schema.ListNestedAttribute{
												Computed: true,
												Description: "List of policy statements that define the permissions granted or denied.\n" +
													"  - example : '[{action: [iam:CreateRole], effect: Allow, resource: [*], sid: Sid1, ...}]'",
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"action": schema.ListAttribute{
															Optional:    true,
															Description: "Actions permitted by the policy statement (e.g., iam:CreateRole, iam:ListUsers).\n  - example : ['iam:CreateRole']",
															ElementType: types.StringType,
														},
														"not_action": schema.ListAttribute{
															Optional:    true,
															Description: "Actions explicitly excluded from the policy statement.\n  - example : ['iam:DeleteRole']",
															ElementType: types.StringType,
														},
														"effect": schema.StringAttribute{
															Computed:    true,
															Description: "Effect of the policy statement - either Allow or Deny.\n  - example : 'Allow'",
														},
														"resource": schema.ListAttribute{
															Optional:    true,
															Description: "Resources that the policy statement applies to (ARNs or wildcards).\n  - example : ['*']",
															ElementType: types.StringType,
														},
														"sid": schema.StringAttribute{
															Computed:    true,
															Description: "Statement ID - unique identifier for this policy statement.\n  - example : 'Sid1'",
														},
														"condition": schema.MapAttribute{
															ElementType: types.MapType{
																ElemType: types.ListType{
																	ElemType: types.StringType,
																},
															},
															Optional:    true,
															Description: "Conditions that must be met for the policy statement to take effect.\n  - example : {'StringEquals': {'aws:PrincipalTag/department': ['IT']}}",
														},
														"principal": schema.SingleNestedAttribute{
															Optional: true,
															Description: "The entity (user, service, or account) that the policy statement applies to.\n" +
																"  - example : '{principal_string: 123456789012, principal_map: {AWS: [arn:aws:iam::123456789012:root]}}'",
															Attributes: map[string]schema.Attribute{
																"principal_string": schema.StringAttribute{
																	Optional:    true,
																	Description: "Principal as a string value (e.g., AWS account ID or IAM user ARN).\n  - example : '123456789012'",
																},
																"principal_map": schema.MapAttribute{
																	Optional: true,
																	ElementType: types.ListType{
																		ElemType: types.StringType,
																	},
																	Description: "Principal as a map - supports multiple principal types (e.g., AWS, Federated, etc.).\n  - example : {'AWS': ['arn:aws:iam::123456789012:root']}",
																},
															},
														},
													},
												},
											},
											"version": schema.StringAttribute{
												Computed:    true,
												Description: "Policy Version\n  - example : '2024-07-01'",
											},
										},
									},

									"policy_id": schema.StringAttribute{
										Computed:    true,
										Description: "Unique identifier of the policy.\n  - example : 'pol-1234567890abcdef'",
									},
									"policy_version_name": schema.StringAttribute{
										Computed:    true,
										Description: "Name of the policy version.\n  - example : 'POLICY_VERSION_1'",
									},
								},
							},
						},
						"resource_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of resource the policy applies to.\n  - example : 'policy'",
						},
						"service_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the service the policy is associated with.\n  - example : 'Identity Access Management'",
						},
						"service_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of service the policy is associated with.\n  - example : 'iam'",
						},
						"srn": schema.StringAttribute{
							Computed:    true,
							Description: "Samsung Resource Name (SRN) - Unique identifier for the policy binding in the SCP system.\n  - example : 'srn:e:::::iam:policy/policy-12345678'",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "Current state of the policy binding (e.g., ACTIVE, INACTIVE).\n  - example : 'ACTIVE'",
						},
					},
				},
			},
		},
	}

}

func (r *iamRolePolicyBindingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("role_id"), req, resp)
}

// Create implements resource.Resource.
// No polling is needed because the AddRolePolicyBindings API call is synchronous
// and returns the created policy bindings directly in the response.
func (r *iamRolePolicyBindingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.RolePolicyBindingsResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data, err := r.client.AddRolePolicyBindings(ctx, plan.RoleId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error add role policy bindings",
			"Could not add role policy bindings, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	rolePolicyObjectValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: iam.Policy{}.Attributes(),
	}, policies)
	plan.RolePolicyBindings = rolePolicyObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamRolePolicyBindingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iam.RolePolicyBindingsResource
	var state iam.RolePolicyBindingsResource

	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// detach
	err := r.client.RemoveRolePolicyBindings(ctx, state.RoleId.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting role policy bindings",
			"Could not delete role policy bindings, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// attach
	data, err := r.client.AddRolePolicyBindings(ctx, plan.RoleId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error add role policy bindings",
			"Could not add role policy bindings, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	rolePolicyObjectValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: iam.Policy{}.Attributes(),
	}, policies)
	plan.RolePolicyBindings = rolePolicyObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamRolePolicyBindingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam.RolePolicyBindingsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RemoveRolePolicyBindings(ctx, state.RoleId.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting role policy bindings",
			"Could not delete role policy bindings, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *iamRolePolicyBindingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.RolePolicyBindingsResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetRolePolicyBindings(ctx, state.RoleId.ValueString(), iam.RolePolicyBindingsDataSource{Size: basetypes.NewInt32Value(20)})
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to read role policy bindings",
			err.Error(),
		)
		return
	}

	getAllRolePolicies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	var rolePolicy []iam.Policy
	for _, policy := range getAllRolePolicies {
		if slices.Contains(state.PolicyIds, policy.Id) {
			rolePolicy = append(rolePolicy, policy)
		}
	}

	rolePolicyObjectValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: iam.Policy{}.Attributes(),
	}, rolePolicy)
	state.RolePolicyBindings = rolePolicyObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
