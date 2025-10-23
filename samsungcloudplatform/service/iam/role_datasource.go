package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
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
		Description: "Show Role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Description:         "Role ID",
				MarkdownDescription: "Role ID",
			},
			"role": schema.SingleNestedAttribute{
				Description: "A detail of Role.",
				Computed:    true,
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
																"condition": schema.MapAttribute{
																	ElementType: types.MapType{
																		ElemType: types.ListType{
																			ElemType: types.StringType,
																		},
																	},
																	Optional: true,
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
