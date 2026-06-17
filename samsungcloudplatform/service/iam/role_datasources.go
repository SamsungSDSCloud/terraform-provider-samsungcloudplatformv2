package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpsdkiam "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/iam/1.4"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iamRoleDataSources{}
	_ datasource.DataSourceWithConfigure = &iamRoleDataSources{}
)

func NewIamRoleDataSources() datasource.DataSource { return &iamRoleDataSources{} }

type iamRoleDataSources struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (d *iamRoleDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_roles"
}

func (d *iamRoleDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamRoleDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show IAM Roles",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)\n" +
					"  - example : 100",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page (between 0 and 10000)\n" +
					"  - example : 0",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort order for results.\n" +
					"  - example : 'created_at,desc'",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Filter by role name.\n" +
					"  - example : 'AdminRole'",
				Optional: true,
			},
			common.ToSnakeCase("RoleType"): schema.StringAttribute{
				Description: "Filter by role type (USER_DEFINED or SYSTEM_DEFINED).\n" +
					"  - example : 'USER_DEFINED'",
				Optional: true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Filter by account ID.\n" +
					"  - example : '123456789012'",
				Optional: true,
			},
			common.ToSnakeCase("Roles"): schema.ListNestedAttribute{
				Description: "List of roles matching the filter criteria.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed: true,
							Description: "Account ID that owns the role.\n" +
								"  - example : '123456789012'",
						},
						"assume_role_policy_document": schema.SingleNestedAttribute{
							Computed:            true,
							Description:         "Policy document that grants assumed role permissions.",
							MarkdownDescription: "Policy document that grants assumed role permissions.",
							Attributes: map[string]schema.Attribute{
								"statement": schema.ListNestedAttribute{
									Computed:            true,
									Description:         "List of policy statements defining the permissions.",
									MarkdownDescription: "List of policy statements defining the permissions.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.ListAttribute{
												Computed: true,
												Description: "List of actions allowed by this statement.\n" +
													"  - example : ['iam:AssumeRole']",
												ElementType: types.StringType,
											},
											"not_action": schema.ListAttribute{
												Optional: true,
												Description: "List of actions that are not allowed by this statement.\n" +
													"  - example : ['iam:AssumeRole']",
												ElementType: types.StringType,
											},
											"effect": schema.StringAttribute{
												Computed: true,
												Description: "Effect of the statement (Allow or Deny).\n" +
													"  - example : 'Allow'",
												MarkdownDescription: "Effect of the statement (Allow or Deny).\n" +
													"  - example : 'Allow'",
											},
											"resource": schema.ListAttribute{
												Computed: true,
												Description: "List of resources the statement applies to.\n" +
													"  - example : ['*']",
												ElementType: types.StringType,
											},
											"principal": schema.SingleNestedAttribute{
												Optional:            true,
												Description:         "Principal that is allowed or denied access.",
												MarkdownDescription: "Principal that is allowed or denied access.",
												Attributes: map[string]schema.Attribute{
													"principal_string": schema.StringAttribute{
														Optional: true,
														Description: "Principal as a string. Specifies the IAM user, role, or account that the policy applies to.\n" +
															"  - example : 'arn:aws:iam::123456789012:user/admin'",
														MarkdownDescription: "Principal as a string. Specifies the IAM user, role, or account that the policy applies to.\n" +
															"  - example : 'arn:aws:iam::123456789012:user/admin'",
													},
													"principal_map": schema.MapAttribute{
														Optional: true,
														Description: "Principal as a map. Specifies multiple principals using key-value pairs.\n" +
															"  - example : {\"AWS\": [\"arn:aws:iam::123456789012:root\"]}",
														MarkdownDescription: "Principal as a map. Specifies multiple principals using key-value pairs.\n" +
															"  - example : {\"AWS\": [\"arn:aws:iam::123456789012:root\"]}",
														ElementType: types.ListType{
															ElemType: types.StringType,
														},
													},
												},
											},
											"sid": schema.StringAttribute{
												Computed: true,
												Description: "Statement ID for the statement.\n" +
													"  - example : 'Stmt1'",
												MarkdownDescription: "Statement ID for the statement.\n" +
													"  - example : 'Stmt1'",
											},
											"condition": schema.MapAttribute{
												Description: "Condition for the policy statement. Specifies constraints on when the policy applies.\n" +
													"  - example : {\"aws:PrincipalTag/department\": [\"engineering\"]}",
												MarkdownDescription: "Condition for the policy statement. Specifies constraints on when the policy applies.\n" +
													"  - example : {\"aws:PrincipalTag/department\": [\"engineering\"]}",
												ElementType: types.MapType{
													ElemType: types.ListType{
														ElemType: types.StringType,
													},
												},
												Optional: true,
											},
										},
									},
								},
								"version": schema.StringAttribute{
									Computed: true,
									Description: "Policy document version.\n" +
										"  - example : '2024-07-01'",
									MarkdownDescription: "Policy document version.\n" +
										"  - example : '2024-07-01'",
								},
							},
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the role was created.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
							MarkdownDescription: "Timestamp when the role was created.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who created the role.\n" +
								"  - example : 'user@example.com'",
							MarkdownDescription: "User who created the role.\n" +
								"  - example : 'user@example.com'",
						},
						"creator_email": schema.StringAttribute{
							Computed: true,
							Description: "Email of the user who created the role.\n" +
								"  - example : 'user@example.com'",
							MarkdownDescription: "Email of the user who created the role.\n" +
								"  - example : 'user@example.com'",
						},
						"creator_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the user who created the role.\n" +
								"  - example : 'John Doe'",
							MarkdownDescription: "Name of the user who created the role.\n" +
								"  - example : 'John Doe'",
						},
						"description": schema.StringAttribute{
							Computed: true,
							Description: "Human-readable description of the role.\n" +
								"  - example : 'My role description'",
						},
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the role.\n" +
								"  - example : 'rol-1234567890abcdef'",
							MarkdownDescription: "Unique identifier of the role.\n" +
								"  - example : 'rol-1234567890abcdef'",
						},
						"max_session_duration": schema.Int64Attribute{
							Computed: true,
							Description: "Maximum duration in seconds that the assumed role can be active.\n" +
								"  - example : 3600",
							MarkdownDescription: "Maximum duration in seconds that the assumed role can be active.\n" +
								"  - example : 3600",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the role was last modified.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
							MarkdownDescription: "Timestamp when the role was last modified.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who last modified the role.\n" +
								"  - example : 'user@example.com'",
							MarkdownDescription: "User who last modified the role.\n" +
								"  - example : 'user@example.com'",
						},
						"modifier_email": schema.StringAttribute{
							Computed: true,
							Description: "Email of the user who last modified the role.\n" +
								"  - example : 'user@example.com'",
							MarkdownDescription: "Email of the user who last modified the role.\n" +
								"  - example : 'user@example.com'",
						},
						"modifier_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the user who last modified the role.\n" +
								"  - example : 'John Doe'",
							MarkdownDescription: "Name of the user who last modified the role.\n" +
								"  - example : 'John Doe'",
						},
						"name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the role.\n" +
								"  - example : 'MyRole'",
						},
						"policies": schema.ListNestedAttribute{
							Optional:            true,
							Description:         "List of policies attached to the role.",
							MarkdownDescription: "List of policies attached to the role.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_id": schema.StringAttribute{
										Optional: true,
										Description: "Account ID that owns the policy.\n" +
											"  - example : '123456789012'",
										MarkdownDescription: "Account ID that owns the policy.\n" +
											"  - example : '123456789012'",
									},
									"created_at": schema.StringAttribute{
										Computed: true,
										Description: "Timestamp when the policy was created.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
										MarkdownDescription: "Timestamp when the policy was created.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
									},
									"created_by": schema.StringAttribute{
										Computed: true,
										Description: "User who created the policy.\n" +
											"  - example : 'user@example.com'",
										MarkdownDescription: "User who created the policy.\n" +
											"  - example : 'user@example.com'",
									},
									"creator_email": schema.StringAttribute{
										Computed: true,
										Description: "Email of the user who created the policy.\n" +
											"  - example : 'user@example.com'",
										MarkdownDescription: "Email of the user who created the policy.\n" +
											"  - example : 'user@example.com'",
									},
									"creator_name": schema.StringAttribute{
										Computed: true,
										Description: "Name of the user who created the policy.\n" +
											"  - example : 'John Doe'",
										MarkdownDescription: "Name of the user who created the policy.\n" +
											"  - example : 'John Doe'",
									},
									"default_version_id": schema.StringAttribute{
										Computed: true,
										Description: "Default version ID of the policy.\n" +
											"  - example : 'pol-1234567890abcdef'",
										MarkdownDescription: "Default version ID of the policy.\n" +
											"  - example : 'pol-1234567890abcdef'",
									},
									"description": schema.StringAttribute{
										Computed: true,
										Description: "Human-readable description of the policy.\n" +
											"  - example : 'My policy description'",
										MarkdownDescription: "Human-readable description of the policy.\n" +
											"  - example : 'My policy description'",
									},
									"domain_name": schema.StringAttribute{
										Computed: true,
										Description: "Domain name associated with the policy.\n" +
											"  - example : 'scp'",
										MarkdownDescription: "Domain name associated with the policy.\n" +
											"  - example : 'scp'",
									},
									"id": schema.StringAttribute{
										Computed: true,
										Description: "Unique identifier of the policy.\n" +
											"  - example : 'pol-1234567890abcdef'",
										MarkdownDescription: "Unique identifier of the policy.\n" +
											"  - example : 'pol-1234567890abcdef'",
									},
									"modified_at": schema.StringAttribute{
										Computed: true,
										Description: "Timestamp when the policy was last modified.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
										MarkdownDescription: "Timestamp when the policy was last modified.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
									},
									"modified_by": schema.StringAttribute{
										Computed: true,
										Description: "User who last modified the policy.\n" +
											"  - example : 'user@example.com'",
										MarkdownDescription: "User who last modified the policy.\n" +
											"  - example : 'user@example.com'",
									},
									"modifier_email": schema.StringAttribute{
										Computed: true,
										Description: "Email of the user who last modified the policy.\n" +
											"  - example : 'user@example.com'",
										MarkdownDescription: "Email of the user who last modified the policy.\n" +
											"  - example : 'user@example.com'",
									},
									"modifier_name": schema.StringAttribute{
										Computed: true,
										Description: "Name of the user who last modified the policy.\n" +
											"  - example : 'John Doe'",
										MarkdownDescription: "Name of the user who last modified the policy.\n" +
											"  - example : 'John Doe'",
									},
									"policy_category": schema.StringAttribute{
										Computed: true,
										Description: "Category of the policy.\n" +
											"  - example : 'IDENTITY_BASED'",
										MarkdownDescription: "Category of the policy.\n" +
											"  - example : 'IDENTITY_BASED'",
									},
									"policy_name": schema.StringAttribute{
										Computed: true,
										Description: "Name of the policy.\n" +
											"  - example : 'MyPolicy'",
										MarkdownDescription: "Name of the policy.\n" +
											"  - example : 'MyPolicy'",
									},
									"policy_type": schema.StringAttribute{
										Computed: true,
										Description: "Type of the policy (USER_DEFINED or SYSTEM_MANAGED).\n" +
											"  - example : 'USER_DEFINED'",
										MarkdownDescription: "Type of the policy (USER_DEFINED or SYSTEM_MANAGED).\n" +
											"  - example : 'USER_DEFINED'",
									},
									"policy_versions": schema.ListNestedAttribute{
										Optional:            true,
										Description:         "List of policy versions.",
										MarkdownDescription: "List of policy versions.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"created_at": schema.StringAttribute{
													Computed: true,
													Description: "Timestamp when the policy version was created.\n" +
														"  - example : '2024-01-01T00:00:00Z'",
													MarkdownDescription: "Timestamp when the policy version was created.\n" +
														"  - example : '2024-01-01T00:00:00Z'",
												},
												"created_by": schema.StringAttribute{
													Computed: true,
													Description: "User who created the policy version.\n" +
														"  - example : 'user@example.com'",
													MarkdownDescription: "User who created the policy version.\n" +
														"  - example : 'user@example.com'",
												},
												"id": schema.StringAttribute{
													Computed: true,
													Description: "Unique identifier of the policy version.\n" +
														"  - example : 'pol-1234567890abcdef'",
													MarkdownDescription: "Unique identifier of the policy version.\n" +
														"  - example : 'pol-1234567890abcdef'",
												},
												"modified_at": schema.StringAttribute{
													Computed: true,
													Description: "Timestamp when the policy version was last modified.\n" +
														"  - example : '2024-01-01T00:00:00Z'",
													MarkdownDescription: "Timestamp when the policy version was last modified.\n" +
														"  - example : '2024-01-01T00:00:00Z'",
												},
												"modified_by": schema.StringAttribute{
													Computed: true,
													Description: "User who last modified the policy version.\n" +
														"  - example : 'user@example.com'",
													MarkdownDescription: "User who last modified the policy version.\n" +
														"  - example : 'user@example.com'",
												},
												"policy_document": schema.SingleNestedAttribute{
													Computed:            true,
													Description:         "The policy document containing the permission definitions.",
													MarkdownDescription: "The policy document containing the permission definitions.",
													Attributes: map[string]schema.Attribute{
														"statement": schema.ListNestedAttribute{
															Computed:            true,
															Description:         "List of policy statements defining the permissions.",
															MarkdownDescription: "List of policy statements defining the permissions.",
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{
																	"action": schema.ListAttribute{
																		Computed: true,
																		Description: "List of actions allowed by this statement.\n" +
																			"  - example : ['iam:CreateRole']",
																		ElementType: types.StringType,
																	},
																	"not_action": schema.ListAttribute{
																		Optional: true,
																		Description: "List of actions that are not allowed by this statement.\n" +
																			"  - example : ['iam:CreateRole']",
																		ElementType: types.StringType,
																	},
																	"effect": schema.StringAttribute{
																		Computed: true,
																		Description: "Effect of the statement (Allow or Deny).\n" +
																			"  - example : 'Allow'",
																		MarkdownDescription: "Effect of the statement (Allow or Deny).\n" +
																			"  - example : 'Allow'",
																	},
																	"resource": schema.ListAttribute{
																		Computed: true,
																		Description: "List of resources the statement applies to.\n" +
																			"  - example : ['*']",
																		ElementType: types.StringType,
																	},
																	"sid": schema.StringAttribute{
																		Computed: true,
																		Description: "Statement ID for the statement.\n" +
																			"  - example : 'Stmt1'",
																		MarkdownDescription: "Statement ID for the statement.\n" +
																			"  - example : 'Stmt1'",
																	},
																	"principal": schema.SingleNestedAttribute{
																		Optional:            true,
																		Description:         "Principal that is allowed or denied access.",
																		MarkdownDescription: "Principal that is allowed or denied access.",
																		Attributes: map[string]schema.Attribute{
																			"principal_string": schema.StringAttribute{
																				Optional: true,
																				Description: "Principal as a string. Specifies the IAM user, role, or account that the policy applies to.\n" +
																					"  - example : 'arn:aws:iam::123456789012:user/admin'",
																				MarkdownDescription: "Principal as a string. Specifies the IAM user, role, or account that the policy applies to.\n" +
																					"  - example : 'arn:aws:iam::123456789012:user/admin'",
																			},
																			"principal_map": schema.MapAttribute{
																				Optional: true,
																				Description: "Principal as a map. Specifies multiple principals using key-value pairs.\n" +
																					"  - example : {\"AWS\": [\"arn:aws:iam::123456789012:root\"]}",
																				MarkdownDescription: "Principal as a map. Specifies multiple principals using key-value pairs.\n" +
																					"  - example : {\"AWS\": [\"arn:aws:iam::123456789012:root\"]}",
																				ElementType: types.ListType{
																					ElemType: types.StringType,
																				},
																			},
																		},
																	},
																	"condition": schema.MapAttribute{
																		Description: "Condition for the policy statement. Specifies constraints on when the policy applies.\n" +
																			"  - example : {\"aws:PrincipalTag/department\": [\"engineering\"]}",
																		MarkdownDescription: "Condition for the policy statement. Specifies constraints on when the policy applies.\n" +
																			"  - example : {\"aws:PrincipalTag/department\": [\"engineering\"]}",
																		ElementType: types.MapType{
																			ElemType: types.ListType{
																				ElemType: types.StringType,
																			},
																		},
																		Optional: true,
																	},
																},
															},
														},
														"version": schema.StringAttribute{
															Computed: true,
															Description: "Policy document version.\n" +
																"  - example : '2024-07-01'",
															MarkdownDescription: "Policy document version.\n" +
																"  - example : '2024-07-01'",
														},
													},
												},
												"policy_id": schema.StringAttribute{
													Computed: true,
													Description: "ID of the policy this version belongs to.\n" +
														"  - example : 'pol-1234567890abcdef'",
													MarkdownDescription: "ID of the policy this version belongs to.\n" +
														"  - example : 'pol-1234567890abcdef'",
												},
												"policy_version_name": schema.StringAttribute{
													Computed: true,
													Description: "Name of the policy version.\n" +
														"  - example : 'POLICY_VERSION_1'",
													MarkdownDescription: "Name of the policy version.\n" +
														"  - example : 'POLICY_VERSION_1'",
												},
											},
										},
									},
									"resource_type": schema.StringAttribute{
										Computed: true,
										Description: "Type of resource the policy applies to.\n" +
											"  - example : 'policy'",
										MarkdownDescription: "Type of resource the policy applies to.\n" +
											"  - example : 'policy'",
									},
									"service_name": schema.StringAttribute{
										Computed: true,
										Description: "Name of the service the policy applies to.\n" +
											"  - example : 'Identity Access Management'",
										MarkdownDescription: "Name of the service the policy applies to.\n" +
											"  - example : 'Identity Access Management'",
									},
									"service_type": schema.StringAttribute{
										Computed: true,
										Description: "Type of service the policy applies to.\n" +
											"  - example : 'iam'",
										MarkdownDescription: "Type of service the policy applies to.\n" +
											"  - example : 'iam'",
									},
									"srn": schema.StringAttribute{
										Computed: true,
										Description: "Samsung Resource Name (SRN) of the policy.\n" +
											"  - example : 'srn:e:::::iam:policy/policy-12345678'",
										MarkdownDescription: "Samsung Resource Name (SRN) of the policy.\n" +
											"  - example : 'srn:e:::::iam:policy/policy-12345678'",
									},
									"state": schema.StringAttribute{
										Computed: true,
										Description: "State of the policy.\n" +
											"  - example : 'ACTIVE'",
										MarkdownDescription: "State of the policy.\n" +
											"  - example : 'ACTIVE'",
									},
								},
							},
						},
						"type": schema.StringAttribute{
							Computed: true,
							Description: "Type of role.\n" +
								"  - example : 'SERVICE'",
							MarkdownDescription: "Type of role.\n" +
								"  - example : 'SERVICE'",
						},
					},
				},
			},
		},
	}
}

func (d *iamRoleDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.RoleDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetRoles(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Roles",
			err.Error(),
		)
		return
	}

	// roles
	roles, hasError := getRoles(ctx, data.Roles)
	if hasError {
		return
	}

	state.Roles = roles

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getRoles(ctx context.Context, _roles []scpsdkiam.Role) ([]iam.Role, bool) {
	var roles []iam.Role
	for _, role := range _roles {

		var statements []iam.Statement
		for _, _statement := range role.AssumeRolePolicyDocument.Statement {

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
				Sid:       types.StringPointerValue(_statement.Sid),
				Effect:    types.StringValue(_statement.Effect),
				Resource:  resources,
				Action:    actions,
				NotAction: notActions,
				Principal: principal,
				Condition: condition,
			}

			statements = append(statements, statement)
		}

		assumeRolePolicyDocument := iam.PolicyDocument{
			Version:   types.StringValue(role.AssumeRolePolicyDocument.Version),
			Statement: statements,
		}

		// policies
		policies, _ := getPolicies(ctx, role.Policies)

		// role nil check
		roleAccountId := role.AccountId.Get()
		if roleAccountId == nil {
			emptyStr := ""
			roleAccountId = &emptyStr
		}

		roleDescription := role.Description.Get()
		if roleDescription == nil {
			emptyStr := ""
			roleDescription = &emptyStr
		}

		roleCreatorName := role.CreatorName.Get()
		if roleCreatorName == nil {
			emptyStr := ""
			roleCreatorName = &emptyStr
		}

		roleCreatorEmail := role.CreatorEmail.Get()
		if roleCreatorEmail == nil {
			emptyStr := ""
			roleCreatorEmail = &emptyStr
		}

		roleModifierName := role.ModifierName.Get()
		if roleModifierName == nil {
			emptyStr := ""
			roleModifierName = &emptyStr
		}

		roleModifierEmail := role.ModifierEmail.Get()
		if roleModifierEmail == nil {
			emptyStr := ""
			roleModifierEmail = &emptyStr
		}

		roleState := iam.Role{
			Id:                       types.StringValue(role.Id),
			Name:                     types.StringValue(role.Name),
			Type:                     types.StringValue(string(*role.Type)),
			AccountId:                types.StringValue(*roleAccountId),
			Description:              types.StringValue(*roleDescription),
			MaxSessionDuration:       types.Int32Value(role.MaxSessionDuration),
			CreatedAt:                types.StringValue(role.CreatedAt.Format(time.RFC3339)),
			CreatedBy:                types.StringValue(role.CreatedBy),
			ModifiedAt:               types.StringValue(role.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:               types.StringValue(role.ModifiedBy),
			CreatorName:              types.StringValue(*roleCreatorName),
			CreatorEmail:             types.StringValue(*roleCreatorEmail),
			ModifierName:             types.StringValue(*roleModifierName),
			ModifierEmail:            types.StringValue(*roleModifierEmail),
			Policies:                 policies,
			AssumeRolePolicyDocument: assumeRolePolicyDocument,
		}

		roles = append(roles, roleState)
	}

	return roles, false
}
