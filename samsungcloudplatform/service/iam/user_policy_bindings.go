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
	_ resource.Resource                = &iamUserPolicyBindingsResource{}
	_ resource.ResourceWithConfigure   = &iamUserPolicyBindingsResource{}
	_ resource.ResourceWithImportState = &iamUserPolicyBindingsResource{}
)

func NewIamUserPolicyBindingsResource() resource.Resource {
	return &iamUserPolicyBindingsResource{}
}

type iamUserPolicyBindingsResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (r *iamUserPolicyBindingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_user_policy_bindings"
}

func (r *iamUserPolicyBindingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamUserPolicyBindingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM User Policy Bindings.",
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				Optional: true,
				Description: "ID of the user to bind policies to.\n" +
					"  - example : 'user-12345678'",
			},
			"policy_ids": schema.ListAttribute{
				Optional: true,
				Description: "List of policy IDs to bind to the user.\n" +
					"  - example : ['policy-12345678', 'policy-87654321']",
				ElementType: types.StringType,
			},
			"user_policy_bindings": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of policy bindings attached to the user.\n" +
					"  - example : '[{account_id: 123456789012, created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, default_version_id: pol-1234567890abcdef, description: My policy description, ...}]'",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed:    true,
							Description: "Account ID of the account that owns the user policy binding.\n  - example : '123456789012'",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the user policy binding was created.\n  - example : '2024-01-01T00:00:00Z'",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "User who created the user policy binding.\n  - example : 'user@example.com'",
						},
						"creator_email": schema.StringAttribute{
							Computed:    true,
							Description: "Email address of the user who created the user policy binding.\n  - example : 'user@example.com'",
						},
						"creator_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the user who created the user policy binding.\n  - example : 'John Doe'",
						},
						"default_version_id": schema.StringAttribute{
							Computed:    true,
							Description: "Default version ID of the policy.\n  - example : 'pol-1234567890abcdef'",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable description of the user policy binding.\n  - example : 'My policy description'",
						},
						"domain_name": schema.StringAttribute{
							Computed:    true,
							Description: "Domain name associated with the user policy binding.\n  - example : 'scp'",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier of the user policy binding.\n  - example : 'pol-1234567890abcdef'",
						},
						"modified_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the user policy binding was last modified.\n  - example : '2024-01-01T00:00:00Z'",
						},
						"modified_by": schema.StringAttribute{
							Computed:    true,
							Description: "User who last modified the user policy binding.\n  - example : 'user@example.com'",
						},
						"modifier_email": schema.StringAttribute{
							Computed:    true,
							Description: "Email address of the user who last modified the user policy binding.\n  - example : 'user@example.com'",
						},
						"modifier_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the user who last modified the user policy binding.\n  - example : 'John Doe'",
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
								"  - example : '[{created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, id: pol-1234567890abcdef, policy_id: pol-1234567890abcdef, policy_version_name: POLICY_VERSION_1, ...}]'",
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
															Description: "Principal - The entity (user, service, or account) that the policy statement applies to.\n" +
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

func (r *iamUserPolicyBindingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("user_id"), req, resp)
}

// No polling is needed because the AddUserPolicyBindings API call is synchronous
// and returns the created policy bindings directly in the response.
func (r *iamUserPolicyBindingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.UserPolicyBindingsResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data, err := r.client.AddUserPolicyBindings(ctx, plan.UserId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error add user policy bindings",
			"Could not add user policy bindings, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	userPolicyObjectValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: iam.Policy{}.Attributes(),
	}, policies)
	plan.UserPolicyBindings = userPolicyObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamUserPolicyBindingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iam.UserPolicyBindingsResource
	var state iam.UserPolicyBindingsResource

	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// detach
	err := r.client.RemoveUserPolicyBindings(ctx, state.UserId.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting user policy bindings",
			"Could not delete user policy bindings, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// attach
	data, err := r.client.AddUserPolicyBindings(ctx, plan.UserId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error add user policy bindings",
			"Could not add user policy bindings, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	userPolicyObjectValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: iam.Policy{}.Attributes(),
	}, policies)
	plan.UserPolicyBindings = userPolicyObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamUserPolicyBindingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam.UserPolicyBindingsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RemoveUserPolicyBindings(ctx, state.UserId.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting user policy bindings",
			"Could not delete user policy bindings, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *iamUserPolicyBindingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.UserPolicyBindingsResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetUserPolicyBindings(ctx, state.UserId.ValueString(), iam.UserPolicyBindingsDataSource{Size: basetypes.NewInt32Value(20)})
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to read user policy bindings",
			err.Error(),
		)
		return
	}

	getAllUserPolicies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	var userPolicy []iam.Policy
	for _, policy := range getAllUserPolicies {
		if slices.Contains(state.PolicyIds, policy.Id) {
			userPolicy = append(userPolicy, policy)
		}
	}

	userPolicyObjectValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: iam.Policy{}.Attributes(),
	}, userPolicy)
	state.UserPolicyBindings = userPolicyObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
