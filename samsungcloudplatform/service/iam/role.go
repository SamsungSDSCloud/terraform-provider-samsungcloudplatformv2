package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpiam1d0 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/iam/1.4"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &iamRoleResource{}
	_ resource.ResourceWithConfigure   = &iamRoleResource{}
	_ resource.ResourceWithImportState = &iamRoleResource{}
)

// NewIamRoleResource is a helper function to simplify the provider implementation.
func NewIamRoleResource() resource.Resource {
	return &iamRoleResource{}
}

// iamRoleResource is the data source implementation.
type iamRoleResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *iamRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role"
}

func (r *iamRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM Role.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Optional: true,
				Description: "Account ID to create the role in.\n" +
					"  - example : '123456789012'",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Description: "Human-readable description of the role.\n" +
					"  - example : 'My role description'",
			},
			"max_session_duration": schema.Int32Attribute{
				Optional: true,
				Description: "Maximum duration for a session using this role (in seconds).\n" +
					"  - example : 3600\n" +
					"  - min: 900, max: 43200",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Description: "Name of the role.\n" +
					"  - example : 'MyRole'\n" +
					"  - maxLength: 64\n" +
					"  - minLength: 1",
			},
			"policy_ids": schema.ListAttribute{
				Optional: true,
				Description: "List of policy IDs to attach to the role.\n" +
					"  - example : ['pol-1234567890abcdef']",
				ElementType: types.StringType,
			},
			"principals": schema.ListNestedAttribute{
				Optional: true,
				Description: "List of principals who can assume this role.\n" +
					"  - example : '[{type: Account, value: ec2.amazonaws.com}]'",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Optional: true,
							Description: "Type of principal.\n" +
								"  - example : 'Account' | 'SCP' | 'Federated' | 'Service'",
						},
						"value": schema.StringAttribute{
							Optional: true,
							Description: "Value identifier for the principal.\n" +
								"  - example : 'ec2.amazonaws.com'",
						},
					},
				},
			},
			"tags": tag.ResourceSchema(),
			"assume_role_policy_document": schema.SingleNestedAttribute{
				Optional: true,
				Description: "Policy document that grants an entity permission to assume the role.\n" +
					"  - example : '{statement: [{action: [sts:AssumeRole], effect: Allow, principal: {principal_map: {Account: [123456789012]}}, ...}], version: 2024-07-01}'",
				Attributes: map[string]schema.Attribute{
					"statement": schema.ListNestedAttribute{
						Optional: true,
						Description: "List of policy statements defining who can assume this role.\n" +
							"  - example : '[{action: [sts:AssumeRole], effect: Allow, principal: {principal_map: {Account: [123456789012]}}, sid: Sid1, ...}]'",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"action": schema.ListAttribute{
									Optional: true,
									Description: "List of actions allowed by this statement.\n" +
										"  - example : ['sts:AssumeRole']",
									ElementType: types.StringType,
								},
								"not_action": schema.ListAttribute{
									Optional: true,
									Description: "List of actions that are not allowed by this statement.\n" +
										"  - example : ['sts:AssumeRole']",
									ElementType: types.StringType,
								},
								"effect": schema.StringAttribute{
									Optional: true,
									Description: "Effect of the statement (allow or deny).\n" +
										"  - example : 'Allow'",
								},
								"resource": schema.ListAttribute{
									Optional: true,
									Description: "List of resources the statement applies to.\n" +
										"  - example : ['*']",
									ElementType: types.StringType,
								},
								"principal": schema.SingleNestedAttribute{
									Optional: true,
									Description: "Principal who can assume the role (e.g., AWS account, service, or federated user).\n" +
										"  - example : '{principal_string: arn:aws:iam::123456789012:root, principal_map: {Account: [123456789012]}}'",
									Attributes: map[string]schema.Attribute{
										"principal_string": schema.StringAttribute{
											Optional:    true,
											Description: "String representation of the principal (e.g., ARN or account ID).\n  - example : 'arn:aws:iam::123456789012:root'",
										},
										"principal_map": schema.MapAttribute{
											Optional:    true,
											Description: "Map of principal type to list of principal identifiers.\n  - example : {\"Account\": [\"123456789012\"]}",
											ElementType: types.ListType{
												ElemType: types.StringType,
											},
										},
									},
								},
								"sid": schema.StringAttribute{
									Optional: true,
									Description: "Statement ID for the statement.\n" +
										"  - example : 'Sid1'",
								},
								"condition": schema.MapAttribute{
									ElementType: types.MapType{
										ElemType: types.ListType{
											ElemType: types.StringType,
										},
									},
									Optional:    true,
									Description: "Condition for the policy statement that determines when the policy is in effect.\n  - example : {\"StringEquals\": {\"aws:PrincipalTag/department\": [\"finance\"]}}",
								},
							},
						},
					},
					"version": schema.StringAttribute{
						Optional: true,
						Description: "Policy document version.\n" +
							"  - example : '2024-07-01'",
					},
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				Description: "Unique identifier of the role.\n" +
					"  - example : 'role-1234567890abcdef'",
			},
			"role": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Detailed information about the role.\n" +
					"  - example : '{account_id: 123456789012, created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, description: My role description, id: role-1234567890abcdef, max_session_duration: 3600, ...}'",
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the role was created.\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"created_by": schema.StringAttribute{
						Computed: true,
						Description: "User who created the role.\n" +
							"  - example : 'user@example.com'",
					},
					"creator_email": schema.StringAttribute{
						Computed: true,
						Description: "Email of the user who created the role.\n" +
							"  - example : 'user@example.com'",
					},
					"creator_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the user who created the role.\n" +
							"  - example : 'John Doe'",
					},
					"description": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Description: "Human-readable description of the role.\n" +
							"  - example : 'My role description'",
					},
					"account_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Description: "Account ID that owns the role.\n" +
							"  - example : '123456789012'",
					},
					"id": schema.StringAttribute{
						Computed: true,
						Description: "Unique identifier of the role.\n" +
							"  - example : 'role-1234567890abcdef'",
					},
					"max_session_duration": schema.Int32Attribute{
						Computed: true,
						Description: "Maximum duration for a session using this role (in seconds).\n" +
							"  - example : 3600",
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the role was last modified.\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
						Description: "User who last modified the role.\n" +
							"  - example : 'user@example.com'",
					},
					"modifier_email": schema.StringAttribute{
						Computed: true,
						Description: "Email of the user who last modified the role.\n" +
							"  - example : 'user@example.com'",
					},
					"modifier_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the user who last modified the role.\n" +
							"  - example : 'John Doe'",
					},
					"name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the role.\n" +
							"  - example : 'AdminRole'",
					},
					"assume_role_policy_document": schema.SingleNestedAttribute{
						Optional: true,
						Description: "Policy document that grants an entity permission to assume the role.\n" +
							"  - example : '{statement: [{action: [sts:AssumeRole], effect: Allow, principal: {principal_map: {Account: [123456789012]}}, ...}], version: 2024-07-01}'",
						Attributes: map[string]schema.Attribute{
							"statement": schema.ListNestedAttribute{
								Optional: true,
								Computed: true,
								Description: "List of policy statements defining who can assume this role.\n" +
									"  - example : '[{action: [sts:AssumeRole], effect: Allow, principal: {principal_map: {Account: [123456789012]}}, sid: Sid1, ...}]'",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"action": schema.ListAttribute{
											Optional: true,
											Description: "List of actions allowed by this statement.\n" +
												"  - example : ['sts:AssumeRole']",
											ElementType: types.StringType,
										},
										"not_action": schema.ListAttribute{
											Optional: true,
											Description: "List of actions that are not allowed by this statement.\n" +
												"  - example : ['sts:AssumeRole']",
											ElementType: types.StringType,
										},
										"effect": schema.StringAttribute{
											Computed: true,
											Description: "Effect of the statement.\n" +
												"  - example : 'Allow'",
										},
										"resource": schema.ListAttribute{
											Optional: true,
											Description: "List of resources the statement applies to.\n" +
												"  - example : ['*']",
											ElementType: types.StringType,
										},
										"principal": schema.SingleNestedAttribute{
											Optional: true,
											Description: "Principal - The entity (user, service, or account) that can assume this role.\n" +
												"  - example : '{principal_string: arn:aws:iam::123456789012:root, principal_map: {Account: [123456789012]}}'",
											Attributes: map[string]schema.Attribute{
												"principal_string": schema.StringAttribute{
													Optional:    true,
													Description: "Principal as a string value (e.g., AWS account ID or IAM user ARN).\n  - example : 'arn:aws:iam::123456789012:root'",
												},
												"principal_map": schema.MapAttribute{
													Optional: true,
													Computed: true,
													ElementType: types.ListType{
														ElemType: types.StringType,
													},
													Description: "Principal as a map - supports multiple principal types (e.g., AWS, Federated, etc.).\n  - example : {'AWS': ['arn:aws:iam::123456789012:root']}",
												},
											},
										},
										"sid": schema.StringAttribute{
											Computed: true,
											Description: "Statement ID.\n" +
												"  - example : 'AllowAssumeRole'",
										},
										"condition": schema.MapAttribute{
											ElementType: types.MapType{
												ElemType: types.ListType{
													ElemType: types.StringType,
												},
											},
											Optional:    true,
											Description: "Condition for the policy statement.\n  - example : {\"StringEquals\": {\"aws:PrincipalTag/department\": [\"finance\"]}}",
										},
									},
								},
							},
							"version": schema.StringAttribute{
								Computed: true,
								Description: "Policy Version\n" +
									"  - example : '2024-07-01'",
							},
						},
					},
					"policies": schema.ListNestedAttribute{
						Optional: true,
						Description: "List of policies attached to the role.\n" +
							"  - example : '[{account_id: 123456789012, created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, id: pol-1234567890abcdef, description: My role description, ...}]'",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"account_id": schema.StringAttribute{
									Optional: true,
									Description: "Account ID that owns the role.\n" +
										"  - example : '123456789012'",
								},
								"created_at": schema.StringAttribute{
									Computed: true,
									Description: "Timestamp when the role was created.\n" +
										"  - example : '2024-01-01T00:00:00Z'",
								},
								"created_by": schema.StringAttribute{
									Computed: true,
									Description: "User who created the role.\n" +
										"  - example : 'user@example.com'",
								},
								"creator_email": schema.StringAttribute{
									Computed: true,
									Description: "Email of the user who created the role.\n" +
										"  - example : 'user@example.com'",
								},
								"creator_name": schema.StringAttribute{
									Computed: true,
									Description: "Name of the user who created the role.\n" +
										"  - example : 'John Doe'",
								},
								"default_version_id": schema.StringAttribute{
									Computed: true,
									Description: "Default version ID of the policy.\n" +
										"  - example : 'pol-1234567890abcde'",
								},
								"description": schema.StringAttribute{
									Computed: true,
									Description: "Human-readable description of the role.\n" +
										"  - example : 'My role description'",
								},
								"domain_name": schema.StringAttribute{
									Computed: true,
									Description: "Domain name associated with the role.\n" +
										"  - example : 'scp'",
								},
								"id": schema.StringAttribute{
									Computed: true,
									Description: "Unique identifier of the policy.\n" +
										"  - example : 'pol-1234567890abcdef'",
								},
								"modified_at": schema.StringAttribute{
									Computed: true,
									Description: "Timestamp when the policy was last modified.\n" +
										"  - example : '2024-01-01T00:00:00Z'",
								},
								"modified_by": schema.StringAttribute{
									Computed: true,
									Description: "User who last modified the policy.\n" +
										"  - example : 'user@example.com'",
								},
								"modifier_email": schema.StringAttribute{
									Computed: true,
									Description: "Email of the user who last modified the policy.\n" +
										"  - example : 'user@example.com'",
								},
								"modifier_name": schema.StringAttribute{
									Computed: true,
									Description: "Name of the user who last modified the policy.\n" +
										"  - example : 'John Doe'",
								},
								"policy_category": schema.StringAttribute{
									Computed: true,
									Description: "Category of the policy (e.g., IDENTITY_BASED or RESOURCE_BASED).\n" +
										"  - example : 'IDENTITY_BASED'",
								},
								"policy_name": schema.StringAttribute{
									Computed: true,
									Description: "Name of the policy.\n" +
										"  - example : 'MyPolicy'",
								},
								"policy_type": schema.StringAttribute{
									Computed: true,
									Description: "Type of the policy.\n" +
										"  - example : 'USER_DEFINED'",
								},
								"policy_versions": schema.ListNestedAttribute{
									Optional: true,
									Description: "List of versions associated with the policy.\n" +
										"  - example : '[{created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, id: ver-1234567890abcdef, policy_id: pol-1234567890abcdef, policy_version_name: POLICY_VERSION_1, ...}]'",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"created_at": schema.StringAttribute{
												Computed: true,
												Description: "Timestamp when the policy version was created.\n" +
													"  - example : '2024-01-01T00:00:00Z'",
											},
											"created_by": schema.StringAttribute{
												Computed: true,
												Description: "User who created the policy version.\n" +
													"  - example : 'user@example.com'",
											},
											"id": schema.StringAttribute{
												Computed: true,
												Description: "Unique identifier of the policy version.\n" +
													"  - example : 'ver-1234567890abcdef'",
											},
											"modified_at": schema.StringAttribute{
												Computed: true,
												Description: "Timestamp when the policy version was last modified.\n" +
													"  - example : '2024-01-01T00:00:00Z'",
											},
											"modified_by": schema.StringAttribute{
												Computed: true,
												Description: "User who last modified the policy version.\n" +
													"  - example : 'user@example.com'",
											},
											"policy_document": schema.SingleNestedAttribute{
												Computed: true,
												Description: "The policy document containing permission definitions for this policy version.\n" +
													"  - example : '{statement: [{action: [iam:CreateRole], effect: Allow, resource: [*], ...}], version: 2024-07-01}'",
												Attributes: map[string]schema.Attribute{
													"statement": schema.ListNestedAttribute{
														Computed: true,
														Description: "List of policy statements that define the permissions granted or denied.\n" +
															"  - example : '[{action: [iam:CreateRole], effect: Allow, resource: [*], sid: Stmt1, ...}]'",
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{
																"action": schema.ListAttribute{
																	Optional: true,
																	Description: "List of actions allowed.\n" +
																		"  - example : ['iam:CreateRole']",
																	ElementType: types.StringType,
																},
																"not_action": schema.ListAttribute{
																	Optional: true,
																	Description: "List of actions that are not allowed.\n" +
																		"  - example : ['iam:CreateRole']",
																	ElementType: types.StringType,
																},
																"effect": schema.StringAttribute{
																	Computed: true,
																	Description: "Effect of the statement.\n" +
																		"  - example : 'Allow'",
																},
																"resource": schema.ListAttribute{
																	Optional: true,
																	Description: "List of resources.\n" +
																		"  - example : ['*']",
																	ElementType: types.StringType,
																},
																"sid": schema.StringAttribute{
																	Computed: true,
																	Description: "Statement ID.\n" +
																		"  - example : 'Stmt1'",
																},
																"condition": schema.MapAttribute{
																	ElementType: types.MapType{
																		ElemType: types.ListType{
																			ElemType: types.StringType,
																		},
																	},
																	Optional:    true,
																	Description: "Condition for the policy statement.\n  - example : {\"StringEquals\": {\"aws:PrincipalTag/department\": [\"finance\"]}}",
																},
																"principal": schema.SingleNestedAttribute{
																	Optional: true,
																	Description: "Principal who can access the resource.\n" +
																		"  - example : '{principal_string: arn:aws:iam::123456789012:root, principal_map: {Account: [123456789012]}}'",
																	Attributes: map[string]schema.Attribute{
																		"principal_string": schema.StringAttribute{
																			Optional:    true,
																			Description: "String representation of the principal.\n  - example : 'arn:aws:iam::123456789012:root'",
																		},
																		"principal_map": schema.MapAttribute{
																			Optional:    true,
																			Description: "Map of principal type to list of principal identifiers.\n  - example : {\"Account\": [\"123456789012\"]}",
																			ElementType: types.ListType{
																				ElemType: types.StringType,
																			},
																		},
																	},
																},
															},
														},
													},
													"version": schema.StringAttribute{
														Computed: true,
														Description: "Policy Version\n" +
															"  - example : '2024-07-01'",
													},
												},
											},
											"policy_id": schema.StringAttribute{
												Computed: true,
												Description: "Unique identifier of the policy.\n" +
													"  - example : 'pol-1234567890abcdef'",
											},
											"policy_version_name": schema.StringAttribute{
												Computed: true,
												Description: "Name of the policy version.\n" +
													"  - example : 'POLICY_VERSION_1'",
											},
										},
									},
								},
								"resource_type": schema.StringAttribute{
									Computed: true,
									Description: "Type of resource the policy applies to.\n" +
										"  - example : 'policy'",
								},
								"service_name": schema.StringAttribute{
									Computed: true,
									Description: "Name of the service the policy is associated with.\n" +
										"  - example : 'Identity Access Management'",
								},
								"service_type": schema.StringAttribute{
									Computed: true,
									Description: "Type of service the policy is associated with.\n" +
										"  - example : 'iam'",
								},
								"srn": schema.StringAttribute{
									Computed: true,
									Description: "Service Resource Name (SRN) - Unique identifier for the role in the SCP system.\n" +
										"  - example : 'srn:e:::::iam:policy/policy-12345678'",
								},
								"state": schema.StringAttribute{
									Computed: true,
									Description: "Current state of the role (e.g., ACTIVE, INACTIVE).\n" +
										"  - example : 'ACTIVE'",
								},
							},
						},
					},
					"type": schema.StringAttribute{
						Computed: true,
						Description: "Type of role.\n" +
							"  - example : 'USER_DEFINED' | 'SERVICE' | 'SERVICE_LINKED'",
					},
				},
			},
		},
	}
}

func (r *iamRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *iamRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.RoleResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateRole(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating role",
			"Could not create role, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// No polling needed - CreateRole API returns the complete role object directly
	// with no pending status field. The role is immediately readable.

	// assume role policy document
	var assumeRolePolicyDocument iam.PolicyDocument

	var statements []iam.Statement
	for _, _statement := range data.Role.AssumeRolePolicyDocument.Statement {
		// resource
		resources := make([]types.String, 0, len(_statement.Resource))
		for _, _resource := range _statement.Resource {
			resources = append(resources, types.StringValue(_resource))
		}

		// action
		actions := make([]types.String, 0, len(_statement.Action))
		for _, _action := range _statement.Action {
			actions = append(actions, types.StringValue(_action))
		}

		// not action
		notActions := make([]types.String, 0, len(_statement.NotAction))
		for _, _notAction := range _statement.NotAction {
			notActions = append(notActions, types.StringValue(_notAction))
		}

		// principal
		principal, _ := convertPrincipal(ctx, _statement.Principal)

		// condition
		condition, _ := convertCondition(ctx, _statement.Condition)

		statement := iam.Statement{
			Sid:       types.StringValue(*_statement.Sid),
			Effect:    types.StringValue(_statement.Effect),
			Resource:  resources,
			Action:    actions,
			NotAction: notActions,
			Principal: principal,
			Condition: condition,
		}

		statements = append(statements, statement)
	}

	assumeRolePolicyDocument = iam.PolicyDocument{
		Version:   types.StringValue(data.Role.AssumeRolePolicyDocument.Version),
		Statement: statements,
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Role.Policies)
	if hasError {
		return
	}

	plan.Id = types.StringValue(data.Role.Id)

	roleState := iam.Role{
		AccountId:                types.StringValue(*data.Role.AccountId.Get()),
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		CreatedAt:                types.StringValue(data.Role.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(data.Role.CreatedBy),
		CreatorEmail:             types.StringValue(*data.Role.CreatorEmail.Get()),
		CreatorName:              types.StringValue(*data.Role.CreatorName.Get()),
		Description:              types.StringPointerValue(data.Role.Description.Get()),
		Id:                       types.StringValue(data.Role.Id),
		MaxSessionDuration:       types.Int32Value(data.Role.MaxSessionDuration),
		ModifiedAt:               types.StringValue(data.Role.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(data.Role.ModifiedBy),
		ModifierEmail:            types.StringValue(*data.Role.ModifierEmail.Get()),
		ModifierName:             types.StringValue(*data.Role.ModifierName.Get()),
		Name:                     types.StringValue(data.Role.Name),
		Policies:                 policies,
		Type:                     types.StringValue(string(*data.Role.Type)),
	}

	roleObjectValue, diags := types.ObjectValueFrom(ctx, roleState.AttributeTypes(), roleState)
	plan.Role = roleObjectValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.RoleResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetRole(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Show Role",
			err.Error(),
		)
		return
	}

	// assume role policy document
	var assumeRolePolicyDocument iam.PolicyDocument

	var statements []iam.Statement
	for _, _statement := range data.Role.AssumeRolePolicyDocument.Statement {
		// resource
		resources := make([]types.String, 0, len(_statement.Resource))
		for _, _resource := range _statement.Resource {
			resources = append(resources, types.StringValue(_resource))
		}

		// action
		actions := make([]types.String, 0, len(_statement.Action))
		for _, _action := range _statement.Action {
			actions = append(actions, types.StringValue(_action))
		}

		// not action
		notActions := make([]types.String, 0, len(_statement.NotAction))
		for _, _notAction := range _statement.NotAction {
			notActions = append(notActions, types.StringValue(_notAction))
		}

		// principal
		principal, _ := convertPrincipal(ctx, _statement.Principal)

		// condition
		condition, _ := convertCondition(ctx, _statement.Condition)

		statement := iam.Statement{
			Sid:       types.StringValue(*_statement.Sid),
			Effect:    types.StringValue(_statement.Effect),
			Resource:  resources,
			Action:    actions,
			NotAction: notActions,
			Principal: principal,
			Condition: condition,
		}

		statements = append(statements, statement)
	}

	assumeRolePolicyDocument = iam.PolicyDocument{
		Version:   types.StringValue(data.Role.AssumeRolePolicyDocument.Version),
		Statement: statements,
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Role.Policies)
	if hasError {
		return
	}

	roleState := iam.Role{
		AccountId:                types.StringValue(*data.Role.AccountId.Get()),
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		CreatedAt:                types.StringValue(data.Role.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(data.Role.CreatedBy),
		CreatorEmail:             types.StringValue(*data.Role.CreatorEmail.Get()),
		CreatorName:              types.StringValue(*data.Role.CreatorName.Get()),
		Description:              types.StringPointerValue(data.Role.Description.Get()),
		Id:                       types.StringValue(data.Role.Id),
		MaxSessionDuration:       types.Int32Value(data.Role.MaxSessionDuration),
		ModifiedAt:               types.StringValue(data.Role.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(data.Role.ModifiedBy),
		ModifierEmail:            types.StringValue(*data.Role.ModifierEmail.Get()),
		ModifierName:             types.StringValue(*data.Role.ModifierName.Get()),
		Name:                     types.StringValue(data.Role.Name),
		Policies:                 policies,
		Type:                     types.StringValue(string(*data.Role.Type)),
	}

	roleObjectValue, diags := types.ObjectValueFrom(ctx, roleState.AttributeTypes(), roleState)
	state.Role = roleObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iam.RoleResource
	var state iam.RoleResource

	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateRole(ctx, state.Id.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Role",
			"Could not update Role, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetRole(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Unable to Read Role",
			"Could not read role ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// assume role policy document
	var assumeRolePolicyDocument iam.PolicyDocument

	var statements []iam.Statement
	for _, _statement := range data.Role.AssumeRolePolicyDocument.Statement {
		// resource
		resources := make([]types.String, 0, len(_statement.Resource))
		for _, _resource := range _statement.Resource {
			resources = append(resources, types.StringValue(_resource))
		}

		// action
		actions := make([]types.String, 0, len(_statement.Action))
		for _, _action := range _statement.Action {
			actions = append(actions, types.StringValue(_action))
		}

		// not action
		notActions := make([]types.String, 0, len(_statement.NotAction))
		for _, _notAction := range _statement.NotAction {
			notActions = append(notActions, types.StringValue(_notAction))
		}

		// principal
		principal, _ := convertPrincipal(ctx, _statement.Principal)

		// condition
		condition, _ := convertCondition(ctx, _statement.Condition)

		statement := iam.Statement{
			Sid:       types.StringValue(*_statement.Sid),
			Effect:    types.StringValue(_statement.Effect),
			Resource:  resources,
			Action:    actions,
			NotAction: notActions,
			Principal: principal,
			Condition: condition,
		}

		statements = append(statements, statement)
	}

	assumeRolePolicyDocument = iam.PolicyDocument{
		Version:   types.StringValue(data.Role.AssumeRolePolicyDocument.Version),
		Statement: statements,
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Role.Policies)
	if hasError {
		return
	}

	roleState := iam.Role{
		AccountId:                types.StringValue(*data.Role.AccountId.Get()),
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		CreatedAt:                types.StringValue(data.Role.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(data.Role.CreatedBy),
		CreatorEmail:             types.StringValue(*data.Role.CreatorEmail.Get()),
		CreatorName:              types.StringValue(*data.Role.CreatorName.Get()),
		Description:              types.StringPointerValue(data.Role.Description.Get()),
		Id:                       types.StringValue(data.Role.Id),
		MaxSessionDuration:       types.Int32Value(data.Role.MaxSessionDuration),
		ModifiedAt:               types.StringValue(data.Role.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(data.Role.ModifiedBy),
		ModifierEmail:            types.StringValue(*data.Role.ModifierEmail.Get()),
		ModifierName:             types.StringValue(*data.Role.ModifierName.Get()),
		Name:                     types.StringValue(data.Role.Name),
		Policies:                 policies,
		Type:                     types.StringValue(string(*data.Role.Type)),
	}

	roleObjectValue, diags := types.ObjectValueFrom(ctx, roleState.AttributeTypes(), roleState)

	plan.Role = roleObjectValue
	plan.Id = state.Id

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam.RoleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRole(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam Role",
			"Could not delete Role, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func convertPrincipal(ctx context.Context, principal interface{}) (iam.Principal, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch v := principal.(type) {

	case scpiam1d0.NullablePrincipal:

		if v.Get() == nil {
			_principal := iam.Principal{
				PrincipalString: types.StringNull(),
				PrincipalMap:    types.MapNull(types.ListType{ElemType: types.StringType}),
			}

			return _principal, diags
		}

		if v.Get().String != nil {
			_principal := iam.Principal{
				PrincipalString: types.StringValue(*v.Get().String),
				PrincipalMap:    types.MapNull(types.ListType{ElemType: types.StringType}),
			}
			return _principal, diags
		}

		if v.Get().MapmapOfStringarrayOfString != nil {
			tempMap := make(map[string]types.List, len(*v.Get().MapmapOfStringarrayOfString))
			for key, val := range *v.Get().MapmapOfStringarrayOfString {
				listVal, listDiags := types.ListValueFrom(ctx, types.StringType, val)
				diags.Append(listDiags...)
				tempMap[key] = listVal
			}

			mapVal, mapDiags := types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, tempMap)
			diags.Append(mapDiags...)

			if diags.HasError() {
				_principal := iam.Principal{
					PrincipalString: types.StringNull(),
					PrincipalMap:    types.MapNull(types.ListType{ElemType: types.StringType}),
				}
				return _principal, diags
			}

			_principal := iam.Principal{
				PrincipalString: types.StringNull(),
				PrincipalMap:    mapVal,
			}
			return _principal, diags
		}

	default:
		diags.AddError(
			"Error converting principal",
			fmt.Sprintf("Unsupported principal type: %T", v),
		)
		_principal := iam.Principal{
			PrincipalString: types.StringNull(),
			PrincipalMap:    types.MapNull(types.ListType{ElemType: types.StringType}),
		}
		return _principal, diags
	}

	_principal := iam.Principal{
		PrincipalString: types.StringNull(),
		PrincipalMap:    types.MapNull(types.ListType{ElemType: types.StringType}),
	}

	return _principal, diags
}
