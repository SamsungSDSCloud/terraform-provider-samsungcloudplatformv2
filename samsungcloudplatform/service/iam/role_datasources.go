package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpsdkiam "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/iam/1.2"
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
		Description: "Show Roles.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size (between 1 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page (between 0 and 10000)",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("RoleType"): schema.StringAttribute{
				Description: "Role Type",
				Optional:    true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Account ID",
				Optional:    true,
			},
			common.ToSnakeCase("Roles"): schema.ListNestedAttribute{
				Description: "A list of role.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Account ID",
							MarkdownDescription: "Account ID",
						},
						"assume_role_policy_document": schema.SingleNestedAttribute{
							Computed:            true,
							Description:         "Assume Role Policy Document",
							MarkdownDescription: "Assume Role Policy Document",
							Attributes: map[string]schema.Attribute{
								"statement": schema.ListNestedAttribute{
									Computed:            true,
									Description:         "Statement",
									MarkdownDescription: "Statement",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.ListAttribute{
												Optional:            true,
												Description:         "Action",
												MarkdownDescription: "Action",
												ElementType:         types.StringType,
											},
											"not_action": schema.ListAttribute{
												Optional:            true,
												Description:         "Not Action",
												MarkdownDescription: "Not Action",
												ElementType:         types.StringType,
											},
											"effect": schema.StringAttribute{
												Computed:            true,
												Description:         "Effect",
												MarkdownDescription: "Effect",
											},
											"resource": schema.ListAttribute{
												Optional:            true,
												Description:         "Resource",
												MarkdownDescription: "Resource",
												ElementType:         types.StringType,
											},
											"principal": schema.SingleNestedAttribute{
												Optional:            true,
												Description:         "Principal",
												MarkdownDescription: "Principal",
												Attributes: map[string]schema.Attribute{
													"principal_string": schema.StringAttribute{
														Optional: true,
													},
													"principal_map": schema.MapAttribute{
														Optional: true,
														ElementType: types.ListType{
															ElemType: types.StringType,
														},
													},
												},
											},
											"sid": schema.StringAttribute{
												Computed:            true,
												Description:         "SID",
												MarkdownDescription: "SID",
											},
											"condition": schema.MapAttribute{
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
									Computed:            true,
									Description:         "Policy Version",
									MarkdownDescription: "Policy Version",
								},
							},
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Created At",
							MarkdownDescription: "Created At",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							Description:         "Created By",
							MarkdownDescription: "Created By",
						},
						"creator_email": schema.StringAttribute{
							Computed:            true,
							Description:         "Creator Email",
							MarkdownDescription: "Creator Email",
						},
						"creator_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Creator Name",
							MarkdownDescription: "Creator Name",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "Description",
							MarkdownDescription: "Description",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "ID",
							MarkdownDescription: "ID",
						},
						"max_session_duration": schema.Int64Attribute{
							Computed:            true,
							Description:         "Max Session Duration",
							MarkdownDescription: "Max Session Duration",
						},
						"modified_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Modified At",
							MarkdownDescription: "Modified At",
						},
						"modified_by": schema.StringAttribute{
							Computed:            true,
							Description:         "Modified By",
							MarkdownDescription: "Modified By",
						},
						"modifier_email": schema.StringAttribute{
							Computed:            true,
							Description:         "Modifier Email",
							MarkdownDescription: "Modifier Email",
						},
						"modifier_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Modifier Name",
							MarkdownDescription: "Modifier Name",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Name",
							MarkdownDescription: "Name",
						},
						"policies": schema.ListNestedAttribute{
							Optional:            true,
							Description:         "Policies",
							MarkdownDescription: "Policies",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_id": schema.StringAttribute{
										Optional:            true,
										Description:         "Account ID",
										MarkdownDescription: "Account ID",
									},
									"created_at": schema.StringAttribute{
										Computed:            true,
										Description:         "Created At",
										MarkdownDescription: "Created At",
									},
									"created_by": schema.StringAttribute{
										Computed:            true,
										Description:         "Created By",
										MarkdownDescription: "Created By",
									},
									"creator_email": schema.StringAttribute{
										Computed:            true,
										Description:         "Creator Email",
										MarkdownDescription: "Creator Email",
									},
									"creator_name": schema.StringAttribute{
										Computed:            true,
										Description:         "Creator Name",
										MarkdownDescription: "Creator Name",
									},
									"default_version_id": schema.StringAttribute{
										Computed:            true,
										Description:         "Default Version ID",
										MarkdownDescription: "Default Version ID",
									},
									"description": schema.StringAttribute{
										Computed:            true,
										Description:         "Description",
										MarkdownDescription: "Description",
									},
									"domain_name": schema.StringAttribute{
										Computed:            true,
										Description:         "Domain Name",
										MarkdownDescription: "Domain Name",
									},
									"id": schema.StringAttribute{
										Computed:            true,
										Description:         "ID",
										MarkdownDescription: "ID",
									},
									"modified_at": schema.StringAttribute{
										Computed:            true,
										Description:         "Modified At",
										MarkdownDescription: "Modified At",
									},
									"modified_by": schema.StringAttribute{
										Computed:            true,
										Description:         "Modified By",
										MarkdownDescription: "Modified By",
									},
									"modifier_email": schema.StringAttribute{
										Computed:            true,
										Description:         "Modifier Email",
										MarkdownDescription: "Modifier Email",
									},
									"modifier_name": schema.StringAttribute{
										Computed:            true,
										Description:         "Modifier Name",
										MarkdownDescription: "Modifier Name",
									},
									"policy_category": schema.StringAttribute{
										Computed:            true,
										Description:         "Policy Category",
										MarkdownDescription: "Policy Category",
									},
									"policy_name": schema.StringAttribute{
										Computed:            true,
										Description:         "Policy Name",
										MarkdownDescription: "Policy Name",
									},
									"policy_type": schema.StringAttribute{
										Computed:            true,
										Description:         "Policy Type",
										MarkdownDescription: "Policy Type",
									},
									"policy_versions": schema.ListNestedAttribute{
										Optional:            true,
										Description:         "Policy Versions",
										MarkdownDescription: "Policy Versions",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"created_at": schema.StringAttribute{
													Computed:            true,
													Description:         "Created At",
													MarkdownDescription: "Created At",
												},
												"created_by": schema.StringAttribute{
													Computed:            true,
													Description:         "Created By",
													MarkdownDescription: "Created By",
												},
												"id": schema.StringAttribute{
													Computed:            true,
													Description:         "ID",
													MarkdownDescription: "ID",
												},
												"modified_at": schema.StringAttribute{
													Computed:            true,
													Description:         "Modified At",
													MarkdownDescription: "Modified At",
												},
												"modified_by": schema.StringAttribute{
													Computed:            true,
													Description:         "Modified By",
													MarkdownDescription: "Modified By",
												},
												"policy_document": schema.SingleNestedAttribute{
													Computed:            true,
													Description:         "Policy Document",
													MarkdownDescription: "Policy Document",
													Attributes: map[string]schema.Attribute{
														"statement": schema.ListNestedAttribute{
															Computed:            true,
															Description:         "Statement",
															MarkdownDescription: "Statement",
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{
																	"action": schema.ListAttribute{
																		Optional:            true,
																		Description:         "Action",
																		MarkdownDescription: "Action",
																		ElementType:         types.StringType,
																	},
																	"not_action": schema.ListAttribute{
																		Optional:            true,
																		Description:         "Not Action",
																		MarkdownDescription: "Not Action",
																		ElementType:         types.StringType,
																	},
																	"effect": schema.StringAttribute{
																		Computed:            true,
																		Description:         "Effect",
																		MarkdownDescription: "Effect",
																	},
																	"resource": schema.ListAttribute{
																		Optional:            true,
																		Description:         "Resource",
																		MarkdownDescription: "Resource",
																		ElementType:         types.StringType,
																	},
																	"sid": schema.StringAttribute{
																		Computed:            true,
																		Description:         "SID",
																		MarkdownDescription: "SID",
																	},
																	"principal": schema.SingleNestedAttribute{
																		Optional:            true,
																		Description:         "Principal",
																		MarkdownDescription: "Principal",
																		Attributes: map[string]schema.Attribute{
																			"principal_string": schema.StringAttribute{
																				Optional: true,
																			},
																			"principal_map": schema.MapAttribute{
																				Optional: true,
																				ElementType: types.ListType{
																					ElemType: types.StringType,
																				},
																			},
																		},
																	},
																	"condition": schema.MapAttribute{
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
															Computed:            true,
															Description:         "Policy Version",
															MarkdownDescription: "Policy Version",
														},
													},
												},
												"policy_id": schema.StringAttribute{
													Computed:            true,
													Description:         "Policy ID",
													MarkdownDescription: "Policy ID",
												},
												"policy_version_name": schema.StringAttribute{
													Computed:            true,
													Description:         "Policy Version Name",
													MarkdownDescription: "Policy Version Name",
												},
											},
										},
									},
									"resource_type": schema.StringAttribute{
										Computed:            true,
										Description:         "Resource Type",
										MarkdownDescription: "Resource Type",
									},
									"service_name": schema.StringAttribute{
										Computed:            true,
										Description:         "Service Name",
										MarkdownDescription: "Service Name",
									},
									"service_type": schema.StringAttribute{
										Computed:            true,
										Description:         "Service Type",
										MarkdownDescription: "Service Type",
									},
									"srn": schema.StringAttribute{
										Computed:            true,
										Description:         "SRN",
										MarkdownDescription: "SRN",
									},
									"state": schema.StringAttribute{
										Computed:            true,
										Description:         "State",
										MarkdownDescription: "State",
									},
								},
							},
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "Type",
							MarkdownDescription: "Type",
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
