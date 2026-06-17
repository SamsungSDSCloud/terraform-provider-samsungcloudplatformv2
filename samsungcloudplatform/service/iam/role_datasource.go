package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iamRoleDataSource{}
	_ datasource.DataSourceWithConfigure = &iamRoleDataSource{}
)

// NewIamRoleDataSource is a helper function to simplify the provider implementation.
func NewIamRoleDataSource() datasource.DataSource {
	return &iamRoleDataSource{}
}

type iamRoleDataSource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *iamRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role"
}

// Configure adds the provider configured client to the data source.
func (d *iamRoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamRoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show IAM Role",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Description: "Unique identifier of the role to retrieve.\n" +
					"  - example : 'rol-1234567890abcdef'",
			},
			"role": schema.SingleNestedAttribute{
				Description: "Detailed information about the role.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed: true,
						Description: "Account ID that owns the role.\n" +
							"  - example : '123456789012'",
					},
					"assume_role_policy_document": schema.SingleNestedAttribute{
						Computed:            true,
						Description:         "Policy document that grants an entity permission to assume the role.",
						MarkdownDescription: "Assume Role Policy Document",
						Attributes: map[string]schema.Attribute{
							"statement": schema.ListNestedAttribute{
								Computed:            true,
								Description:         "Statement - list of permission statements in the policy.\n  - example : [{'Sid': 'Stmt1', 'Effect': 'Allow', 'Action': [...], 'Resource': '*'}]",
								MarkdownDescription: "Statement - list of permission statements in the policy.\n  - example : [{'Sid': 'Stmt1', 'Effect': 'Allow', 'Action': [...], 'Resource': '*'}]",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"action": schema.ListAttribute{
											Computed: true,
											Description: "List of actions allowed by this statement.\n" +
												"  - example : ['iam:AssumeRole']",
											ElementType: types.StringType,
										},
										"not_action": schema.ListAttribute{
											Computed: true,
											Description: "List of actions that are not allowed by this statement.\n" +
												"  - example : ['iam:AssumeRole']",
											ElementType: types.StringType,
										},
										"effect": schema.StringAttribute{
											Computed: true,
											Description: "Effect of the statement (Allow or Deny).\n" +
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
											Description:         "Principal - the entity (user, group, or service) that the policy statement applies to.\n  - example : {'Service': ['ec2.amazonaws.com']}",
											MarkdownDescription: "Principal - the entity (user, group, or service) that the policy statement applies to.\n  - example : {'Service': ['ec2.amazonaws.com']}",
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
											Computed:            true,
											Description:         "Statement ID (SID) - unique identifier for the policy statement.\n  - example : 'Stmt1'",
											MarkdownDescription: "Statement ID (SID) - unique identifier for the policy statement.\n  - example : 'Stmt1'",
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
								Description: "Policy Version.\n" +
									"  - example : '2024-07-01'",
							},
						},
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
					"description": schema.StringAttribute{
						Computed: true,
						Description: "Human-readable description of the role.\n" +
							"  - example : 'My role description'",
					},
					"id": schema.StringAttribute{
						Computed: true,
						Description: "Unique identifier of the role.\n" +
							"  - example : 'rol-1234567890abcdef'",
					},
					"max_session_duration": schema.Int64Attribute{
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
							"  - example : 'MyRole'",
					},
					"policies": schema.ListNestedAttribute{
						Optional:            true,
						Description:         "List of policies attached to the role.",
						MarkdownDescription: "List of policies attached to the role.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"account_id": schema.StringAttribute{
									Computed: true,
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
										"  - example : 'v1'",
									MarkdownDescription: "Default version ID of the policy.\n" +
										"  - example : 'v1'",
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
									Description: "Domain name associated with the role.\n" +
										"  - example : 'scp'",
									MarkdownDescription: "Domain name associated with the role.\n" +
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
										"  - example : 'scp'",
									MarkdownDescription: "Category of the policy.\n" +
										"  - example : 'scp'",
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
									Description: "Type of the policy (IDENTITY_BASED or RESOURCE_BASED).\n" +
										"  - example : 'IDENTITY_BASED'",
									MarkdownDescription: "Type of the policy (IDENTITY_BASED or RESOURCE_BASED).\n" +
										"  - example : 'IDENTITY_BASED'",
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
												Description:         "Policy Document",
												MarkdownDescription: "Policy Document",
												Attributes: map[string]schema.Attribute{
													"statement": schema.ListNestedAttribute{
														Computed:            true,
														Description:         "Statement - list of permission statements in the policy.\n  - example : [{'Sid': 'Stmt1', 'Effect': 'Allow', 'Action': [...], 'Resource': '*'}]",
														MarkdownDescription: "Statement - list of permission statements in the policy.\n  - example : [{'Sid': 'Stmt1', 'Effect': 'Allow', 'Action': [...], 'Resource': '*'}]",
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{
																"action": schema.ListAttribute{
																	Computed: true,
																	Description: "List of actions allowed by this statement.\n" +
																		"  - example : ['iam:CreateRole']",
																	ElementType: types.StringType,
																},
																"not_action": schema.ListAttribute{
																	Computed: true,
																	Description: "List of actions that are not allowed by this statement.\n" +
																		"  - example : ['iam:CreateRole']",
																	ElementType: types.StringType,
																},
																"effect": schema.StringAttribute{
																	Computed: true,
																	Description: "Effect of the statement (Allow or Deny).\n" +
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
																"principal": schema.SingleNestedAttribute{
																	Optional:            true,
																	Description:         "Principal - the entity (user, group, or service) that the policy statement applies to.\n  - example : {'Service': ['ec2.amazonaws.com']}",
																	MarkdownDescription: "Principal - the entity (user, group, or service) that the policy statement applies to.\n  - example : {'Service': ['ec2.amazonaws.com']}",
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
															},
														},
													},
													"version": schema.StringAttribute{
														Computed: true,
														Description: "Policy Version.\n" +
															"  - example : '2024-07-01'",
													},
												},
											},
											"policy_id": schema.StringAttribute{
												Computed: true,
												Description: "ID of the policy.\n" +
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
									Description: "Name of the service.\n" +
										"  - example : 'Identity Access Management'",
								},
								"service_type": schema.StringAttribute{
									Computed: true,
									Description: "Type of service.\n" +
										"  - example : 'iam'",
								},
								"srn": schema.StringAttribute{
									Computed: true,
									Description: "Service Resource Name (SRN).\n" +
										"  - example : 'srn:e:::::iam:policy/policy-12345678'",
								},
								"state": schema.StringAttribute{
									Computed: true,
									Description: "Current state of the policy.\n" +
										"  - example : 'ACTIVE'",
								},
							},
						},
					},
					"type": schema.StringAttribute{
						Computed: true,
						Description: "Type of role.\n" +
							"  - example : 'SERVICE'",
					},
				},
			},
		},
	}
}

func (d *iamRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.RoleDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetRole(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Role",
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

	// role nil check
	roleAccountId := data.Role.AccountId.Get()
	if roleAccountId == nil {
		emptyStr := ""
		roleAccountId = &emptyStr
	}

	roleDescription := data.Role.Description.Get()
	if roleDescription == nil {
		emptyStr := ""
		roleDescription = &emptyStr
	}

	roleCreatorName := data.Role.CreatorName.Get()
	if roleCreatorName == nil {
		emptyStr := ""
		roleCreatorName = &emptyStr
	}

	roleCreatorEmail := data.Role.CreatorEmail.Get()
	if roleCreatorEmail == nil {
		emptyStr := ""
		roleCreatorEmail = &emptyStr
	}

	roleModifierName := data.Role.ModifierName.Get()
	if roleModifierName == nil {
		emptyStr := ""
		roleModifierName = &emptyStr
	}

	roleModifierEmail := data.Role.ModifierEmail.Get()
	if roleModifierEmail == nil {
		emptyStr := ""
		roleModifierEmail = &emptyStr
	}

	roleState := iam.Role{
		AccountId:                types.StringValue(*roleAccountId),
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		CreatedAt:                types.StringValue(data.Role.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(data.Role.CreatedBy),
		CreatorEmail:             types.StringValue(*roleCreatorEmail),
		CreatorName:              types.StringValue(*roleCreatorName),
		Description:              types.StringValue(*roleDescription),
		Id:                       types.StringValue(data.Role.Id),
		MaxSessionDuration:       types.Int32Value(data.Role.MaxSessionDuration),
		ModifiedAt:               types.StringValue(data.Role.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(data.Role.ModifiedBy),
		ModifierEmail:            types.StringValue(*roleModifierEmail),
		ModifierName:             types.StringValue(*roleModifierName),
		Name:                     types.StringValue(data.Role.Name),
		Policies:                 policies,
		Type:                     types.StringValue(string(*data.Role.Type)),
	}

	roleObjectValue, _ := types.ObjectValueFrom(ctx, roleState.AttributeTypes(), roleState)
	state.Role = roleObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
