package iam

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &iamGroupPolicyBindingDataSources{}
	_ datasource.DataSourceWithConfigure = &iamGroupPolicyBindingDataSources{}
)

func NewIamGroupPolicyBindingDataSources() datasource.DataSource {
	return &iamGroupPolicyBindingDataSources{}
}

type iamGroupPolicyBindingDataSources struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (d *iamGroupPolicyBindingDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group_policy_bindings"
}

func (d *iamGroupPolicyBindingDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamGroupPolicyBindingDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Group Policy Bindings",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Optional: true,
				Description: "Group ID to filter policy bindings.\n" +
					"  - example : 'group-12345678'",
			},
			"size": schema.Int32Attribute{
				Optional: true,
				Description: "Size (between 1 and 10000)\n" +
					"  - example : 100",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			"page": schema.Int32Attribute{
				Optional: true,
				Description: "Page (between 0 and 10000)\n" +
					"  - example : 0",
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort order for results.\n" +
					"  - example : 'created_at,desc'",
			},
			"policy_id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the policy.\n  - example : 'policy-12345678'",
				MarkdownDescription: "ID of the policy.\n  - example : 'policy-12345678'",
			},
			"policy_version_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the policy version.\n  - example : 'POLICY_VERSION_1'",
				MarkdownDescription: "Name of the policy version.\n  - example : 'POLICY_VERSION_1'",
			},
			"policy_name": schema.StringAttribute{
				Optional: true,
				Description: "Filter by policy name.\n" +
					"  - example : 'MyPolicy'",
			},
			"policy_type": schema.StringAttribute{
				Optional: true,
				Description: "Type of the policy (e.g., USER_DEFINED, SYSTEM_DEFINED).\n" +
					"  - example : 'USER_DEFINED'",
				MarkdownDescription: "Type of the policy (e.g., USER_DEFINED, SYSTEM_DEFINED).\n" +
					"  - example : 'USER_DEFINED'",
			},
			"group_policy_bindings": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Group Policy Bindings",
				MarkdownDescription: "Group Policy Bindings",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Account ID that owns this policy binding.\n  - example : '123456789012'",
							MarkdownDescription: "Account ID that owns this policy binding.\n  - example : '123456789012'",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Timestamp when the policy binding was created.\n  - example : '2024-01-01T00:00:00Z'",
							MarkdownDescription: "Timestamp when the policy binding was created.\n  - example : '2024-01-01T00:00:00Z'",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							Description:         "User who created the policy binding.\n  - example : 'user@example.com'",
							MarkdownDescription: "User who created the policy binding.\n  - example : 'user@example.com'",
						},
						"creator_email": schema.StringAttribute{
							Computed:            true,
							Description:         "Email of the user who created the policy binding.\n  - example : 'user@example.com'",
							MarkdownDescription: "Email of the user who created the policy binding.\n  - example : 'user@example.com'",
						},
						"creator_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Name of the user who created the policy binding.\n  - example : 'John Doe'",
							MarkdownDescription: "Name of the user who created the policy binding.\n  - example : 'John Doe'",
						},
						"default_version_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Default version ID of the policy.\n  - example : 'pol-1234567890abcdef'",
							MarkdownDescription: "Default version ID of the policy.\n  - example : 'pol-1234567890abcdef'",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "Human-readable description of the policy binding.\n  - example : 'My policy description'",
							MarkdownDescription: "Human-readable description of the policy binding.\n  - example : 'My policy description'",
						},
						"domain_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Domain name associated with the policy binding.\n  - example : 'scp'",
							MarkdownDescription: "Domain name associated with the policy binding.\n  - example : 'scp'",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Unique identifier of the policy binding.\n  - example : 'pol-1234567890abcdef'",
							MarkdownDescription: "Unique identifier of the policy binding.\n  - example : 'pol-1234567890abcdef'",
						},
						"modified_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Timestamp when the policy binding was last modified.\n  - example : '2024-01-01T00:00:00Z'",
							MarkdownDescription: "Timestamp when the policy binding was last modified.\n  - example : '2024-01-01T00:00:00Z'",
						},
						"modified_by": schema.StringAttribute{
							Computed:            true,
							Description:         "User who last modified the policy binding.\n  - example : 'user@example.com'",
							MarkdownDescription: "User who last modified the policy binding.\n  - example : 'user@example.com'",
						},
						"modifier_email": schema.StringAttribute{
							Computed:            true,
							Description:         "Email of the user who last modified the policy binding.\n  - example : 'user@example.com'",
							MarkdownDescription: "Email of the user who last modified the policy binding.\n  - example : 'user@example.com'",
						},
						"modifier_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Name of the user who last modified the policy binding.\n  - example : 'John Doe'",
							MarkdownDescription: "Name of the user who last modified the policy binding.\n  - example : 'John Doe'",
						},
						"policy_category": schema.StringAttribute{
							Computed:            true,
							Description:         "Category of the policy.\n  - example : 'IDENTITY_BASED' | 'RESOURCE_BASED'",
							MarkdownDescription: "Category of the policy.\n  - example : 'IDENTITY_BASED' | 'RESOURCE_BASED'",
						},
						"policy_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Name of the policy.\n  - example : 'MyPolicy'",
							MarkdownDescription: "Name of the policy.\n  - example : 'MyPolicy'",
						},
						"policy_type": schema.StringAttribute{
							Computed:            true,
							Description:         "Type of the policy.\n  - example : 'USER_DEFINED' | 'SYSTEM_MANAGED'",
							MarkdownDescription: "Type of the policy.\n  - example : 'USER_DEFINED' | 'SYSTEM_MANAGED'",
						},
						"policy_versions": schema.ListNestedAttribute{
							Optional:            true,
							Description:         "Policy Versions",
							MarkdownDescription: "Policy Versions",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"created_at": schema.StringAttribute{
										Computed:            true,
										Description:         "Timestamp when the policy version was created.\n  - example : '2024-01-01T00:00:00Z'",
										MarkdownDescription: "Timestamp when the policy version was created.\n  - example : '2024-01-01T00:00:00Z'",
									},
									"created_by": schema.StringAttribute{
										Computed:            true,
										Description:         "User who created the policy version.\n  - example : 'user@example.com'",
										MarkdownDescription: "User who created the policy version.\n  - example : 'user@example.com'",
									},
									"id": schema.StringAttribute{
										Computed:            true,
										Description:         "Policy version ID - unique identifier for this policy version.\n  - example : 'pol-1234567890abcdef'",
										MarkdownDescription: "Policy version ID - unique identifier for this policy version.\n  - example : 'pol-1234567890abcdef'",
									},
									"modified_at": schema.StringAttribute{
										Computed:            true,
										Description:         "Timestamp when the policy version was last modified.\n  - example : '2024-01-01T00:00:00Z'",
										MarkdownDescription: "Timestamp when the policy version was last modified.\n  - example : '2024-01-01T00:00:00Z'",
									},
									"modified_by": schema.StringAttribute{
										Computed:            true,
										Description:         "User who last modified the policy version.\n  - example : 'user@example.com'",
										MarkdownDescription: "User who last modified the policy version.\n  - example : 'user@example.com'",
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
															Optional:            true,
															Description:         "Actions allowed or denied by the policy statement.\n  - example : ['iam:CreateRole']",
															MarkdownDescription: "Actions allowed or denied by the policy statement.\n  - example : ['iam:CreateRole']",
															ElementType:         types.StringType,
														},
														"not_action": schema.ListAttribute{
															Optional:            true,
															Description:         "Actions that are excluded from the policy statement.\n  - example : ['iam:DeleteRole']",
															MarkdownDescription: "Actions that are excluded from the policy statement.\n  - example : ['iam:DeleteRole']",
															ElementType:         types.StringType,
														},
														"effect": schema.StringAttribute{
															Computed:            true,
															Description:         "Effect of the policy statement (Allow or Deny).\n  - example : 'Allow' | 'Deny'",
															MarkdownDescription: "Effect of the policy statement (Allow or Deny).\n  - example : 'Allow' | 'Deny'",
														},
														"resource": schema.ListAttribute{
															Optional:            true,
															Description:         "Resources that the policy statement applies to.\n  - example : ['*']",
															MarkdownDescription: "Resources that the policy statement applies to.\n  - example : ['*']",
															ElementType:         types.StringType,
														},
														"sid": schema.StringAttribute{
															Computed:            true,
															Description:         "Statement ID (SID) - unique identifier for the policy statement.\n  - example : 'Stmt1'",
															MarkdownDescription: "Statement ID (SID) - unique identifier for the policy statement.\n  - example : 'Stmt1'",
														},
														"condition": schema.MapAttribute{
															ElementType: types.MapType{
																ElemType: types.ListType{
																	ElemType: types.StringType,
																},
															},
															Optional:            true,
															Description:         "Condition for the policy statement - specifies when the policy effect takes effect.\n  - example : {'StringEquals': {'scp:PrincipalTag/department': ['engineering']}}",
															MarkdownDescription: "Condition for the policy statement - specifies when the policy effect takes effect.\n  - example : {'StringEquals': {'scp:PrincipalTag/department': ['engineering']}}",
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
												Computed:            true,
												Description:         "Policy Version\n  - example : '2024-07-01'",
												MarkdownDescription: "Policy Version\n  - example : '2024-07-01'",
											},
										},
									},

									"policy_id": schema.StringAttribute{
										Computed:            true,
										Description:         "ID of the policy associated with this binding.\n  - example : 'policy-12345678'",
										MarkdownDescription: "ID of the policy associated with this binding.\n  - example : 'policy-12345678'",
									},
									"policy_version_name": schema.StringAttribute{
										Computed:            true,
										Description:         "Name of the policy version.\n  - example : 'POLICY_VERSION_1'",
										MarkdownDescription: "Name of the policy version.\n  - example : 'POLICY_VERSION_1'",
									},
								},
							},
						},
						"resource_type": schema.StringAttribute{
							Computed:            true,
							Description:         "Type of resource the policy applies to.\n  - example : 'policy'",
							MarkdownDescription: "Type of resource the policy applies to.\n  - example : 'policy'",
						},
						"service_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Name of the service.\n  - example : 'Identity Access Management'",
							MarkdownDescription: "Name of the service.\n  - example : 'Identity Access Management'",
						},
						"service_type": schema.StringAttribute{
							Computed:            true,
							Description:         "Type of service.\n  - example : 'iam'",
							MarkdownDescription: "Type of service.\n  - example : 'iam'",
						},
						"srn": schema.StringAttribute{
							Computed:            true,
							Description:         "Service Resource Name (SRN).\n  - example : 'srn:e:::::iam:policy/policy-12345678'",
							MarkdownDescription: "Service Resource Name (SRN).\n  - example : 'srn:e:::::iam:policy/policy-12345678'",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "State of the policy binding.\n  - example : 'ACTIVE'",
							MarkdownDescription: "State of the policy binding.\n  - example : 'ACTIVE'",
						},
					},
				},
			},
		},
	}
}

func (d *iamGroupPolicyBindingDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.GroupPolicyBindingsDataResource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetGroupPolicyBindings(ctx, state.GroupId.ValueString(), state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Group Policies",
			err.Error(),
		)
		return
	}

	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	state.GroupPolicyBindings = policies

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
