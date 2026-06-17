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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iamPolicyDataSources{}
	_ datasource.DataSourceWithConfigure = &iamPolicyDataSources{}
)

func NewIamPolicyDataSources() datasource.DataSource { return &iamPolicyDataSources{} }

type iamPolicyDataSources struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (d *iamPolicyDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_policies"
}

func (d *iamPolicyDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamPolicyDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show IAM Policies",
		Attributes: map[string]schema.Attribute{
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
					"  - example : 'createdAt,desc'",
			},
			"id": schema.StringAttribute{
				Optional: true,
				Description: "Filter by policy ID.\n" +
					"  - example : 'pol-1234567890abcdef'",
			},
			"policy_name": schema.StringAttribute{
				Optional: true,
				Description: "Filter by policy name.\n" +
					"  - example : 'MyPolicy'",
			},
			"policy_type": schema.StringAttribute{
				Optional: true,
				Description: "Type of the policy (e.g., IDENTITY_BASED, RESOURCE_BASED).\n" +
					"  - example : 'IDENTITY_BASED'",
				MarkdownDescription: "Type of the policy (e.g., IDENTITY_BASED, RESOURCE_BASED).\n" +
					"  - example : 'IDENTITY_BASED'",
			},
			"policies": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "A list of policy",
				MarkdownDescription: "A list of policy",
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
								"  - example : 'pol-1234567890abcdef'",
						},
						"description": schema.StringAttribute{
							Computed: true,
							Description: "Human-readable description of the policy.\n" +
								"  - example : 'My policy description'",
						},
						"domain_name": schema.StringAttribute{
							Computed: true,
							Description: "Domain name associated with the policy.\n" +
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
							Description: "Category of the policy.\n" +
								"  - example : 'IDENTITY_BASED'",
						},
						"policy_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the policy.\n" +
								"  - example : 'PolicyName'",
						},
						"policy_type": schema.StringAttribute{
							Computed: true,
							Description: "Type of the policy.\n" +
								"  - example : 'USER_DEFINED'",
							MarkdownDescription: "Type of the policy.\n" +
								"  - example : 'USER_DEFINED'",
						},
						"policy_versions": schema.ListNestedAttribute{
							Optional:            true,
							Description:         "Policy Versions",
							MarkdownDescription: "Policy Versions",
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
															Optional: true,
															Description: "List of actions allowed by this statement.\n" +
																"  - example : ['iam:CreateRole']",
															ElementType: types.StringType,
														},
														"not_action": schema.ListAttribute{
															Optional:    true,
															Description: "List of actions that are not allowed by this statement.\n  - example : ['iam:DeleteRole']",
															ElementType: types.StringType,
														},
														"effect": schema.StringAttribute{
															Computed: true,
															Description: "Effect of the statement (Allow or Deny).\n" +
																"  - example : 'Allow'",
														},
														"resource": schema.ListAttribute{
															Computed:            true,
															Description:         "Resources that the policy statement applies to.\n  - example : ['*']",
															MarkdownDescription: "Resources that the policy statement applies to.\n  - example : ['*']",
															ElementType:         types.StringType,
														},
														"sid": schema.StringAttribute{
															Computed:            true,
															Description:         "Statement ID (SID) - unique identifier for the policy statement.\n  - example : 'Sid1'",
															MarkdownDescription: "Statement ID (SID) - unique identifier for the policy statement.\n  - example : 'Sid1'",
														},
														"condition": schema.MapAttribute{
															Description: "Condition for the policy statement. Specifies constraints on when the policy applies. For example, use with IpAddress or DateGreaterThan.\n" +
																"  - example : {\"aws:PrincipalTag/department\": [\"engineering\"]}",
															MarkdownDescription: "Condition for the policy statement. Specifies constraints on when the policy applies. For example, use with IpAddress or DateGreaterThan.\n" +
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
		},
	}
}

func (d *iamPolicyDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.PolicyDatasource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetPolicies(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Policies",
			err.Error(),
		)
		return
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	state.Policies = policies

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
