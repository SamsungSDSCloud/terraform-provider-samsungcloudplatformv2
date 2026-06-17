package iam

import (
	"context"
	"fmt"
	"time"

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
	_ datasource.DataSource              = &iamGroupDataSources{}
	_ datasource.DataSourceWithConfigure = &iamGroupDataSources{}
)

// NewIamGroupDataSources is a helper function to simplify the provider implementation.
func NewIamGroupDataSources() datasource.DataSource {
	return &iamGroupDataSources{}
}

type iamGroupDataSources struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *iamGroupDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_groups"
}

// Configure adds the provider configured client to the data source.
func (d *iamGroupDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamGroupDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show IAM Groups",
		Attributes: map[string]schema.Attribute{
			"size": schema.Int32Attribute{
				Optional: true,
				Description: "Number of results to return per page.\n" +
					"  - example : 100\n" +
					"  - min: 1, max: 10000",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			"page": schema.Int32Attribute{
				Optional: true,
				Description: "Page number to retrieve.\n" +
					"  - example : 0\n" +
					"  - min: 0, max: 10000",
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort order for results (e.g., 'createdAt,desc').\n" +
					"  - example : 'createdAt,desc'",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Description: "Filter groups by name.\n" +
					"  - example : 'MyGroup'",
			},
			"groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of groups matching the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the group was created.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who created the group.\n" +
								"  - example : 'user@example.com'",
						},
						"creator_email": schema.StringAttribute{
							Computed: true,
							Description: "Email of the user who created the group.\n" +
								"  - example : 'user@example.com'",
						},
						"creator_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the user who created the group.\n" +
								"  - example : 'John Doe'",
						},
						"description": schema.StringAttribute{
							Computed: true,
							Description: "Human-readable description of the group.\n" +
								"  - example : 'My group description'",
						},
						"domain_name": schema.StringAttribute{
							Computed: true,
							Description: "Domain name associated with the group.\n" +
								"  - example : 'scp'",
						},
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the group.\n" +
								"  - example : 'grp-1234567890abcdef'",
						},
						"members": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of members in the group.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"created_at": schema.StringAttribute{
										Computed: true,
										Description: "Timestamp when the member was added.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
									},
									"created_by": schema.StringAttribute{
										Computed: true,
										Description: "User who added the member.\n" +
											"  - example : 'user@example.com'",
									},
									"creator_created_at": schema.StringAttribute{
										Computed: true,
										Description: "Timestamp when the creator was created.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
									},
									"creator_email": schema.StringAttribute{
										Computed: true,
										Description: "Email of the creator.\n" +
											"  - example : 'user@example.com'",
									},
									"creator_last_login_at": schema.StringAttribute{
										Computed: true,
										Description: "Timestamp of the creator's last login.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
									},
									"creator_name": schema.StringAttribute{
										Computed: true,
										Description: "Name of the creator.\n" +
											"  - example : 'John Doe'",
									},
									"group_names": schema.ListAttribute{
										ElementType: types.StringType,
										Computed:    true,
										Description: "Names of the groups the user belongs to.\n" +
											"  - example : ['MyGroup']",
									},
									"user_created_at": schema.StringAttribute{
										Computed: true,
										Description: "Timestamp when the user was created.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
									},
									"user_email": schema.StringAttribute{
										Computed: true,
										Description: "Email of the user.\n" +
											"  - example : 'member@example.com'",
									},
									"user_id": schema.StringAttribute{
										Computed: true,
										Description: "Unique identifier of the user.\n" +
											"  - example : 'usr-1234567890abcdef'",
									},
									"user_last_login_at": schema.StringAttribute{
										Computed: true,
										Description: "Timestamp of the user's last login.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
									},
									"user_name": schema.StringAttribute{
										Computed: true,
										Description: "Name of the user.\n" +
											"  - example : 'Jane Doe'",
									},
								},
							},
						},
						"policies": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_id": schema.StringAttribute{
										Computed: true,
										Description: "Account ID that owns the policy.\n" +
											"  - example : '123456789012'",
									},
									"created_at": schema.StringAttribute{
										Computed: true,
										Description: "Timestamp when the policy was created.\n" +
											"  - example : '2024-01-01T00:00:00Z'",
									},
									"created_by": schema.StringAttribute{
										Computed: true,
										Description: "User who created the policy.\n" +
											"  - example : 'user@example.com'",
									},
									"creator_email": schema.StringAttribute{
										Computed: true,
										Description: "Email of the policy creator.\n" +
											"  - example : 'user@example.com'",
									},
									"creator_name": schema.StringAttribute{
										Computed: true,
										Description: "Name of the policy creator.\n" +
											"  - example : 'John Doe'",
									},
									"default_version_id": schema.StringAttribute{
										Computed: true,
										Description: "Default version ID of the policy.\n" +
											"  - example : 'v1'",
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
											"  - example : 'IDENTITY_BASED' | 'RESOURCE_BASED' | 'SESSION'",
									},
									"policy_name": schema.StringAttribute{
										Computed: true,
										Description: "Name of the policy.\n" +
											"  - example : 'MyPolicy'",
									},
									"policy_type": schema.StringAttribute{
										Computed: true,
										Description: "Type of the policy.\n" +
											"  - example : 'SYSTEM_MANAGED' | 'USER_DEFINED' | 'INLINE'",
									},
									"policy_versions": schema.ListNestedAttribute{
										Computed:    true,
										Description: "List of versions of the policy.",
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
														"  - example : 'v-1234567890abcdef'",
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
													Computed:    true,
													Description: "The policy document containing the permission definitions.",
													Attributes: map[string]schema.Attribute{
														"statement": schema.ListNestedAttribute{
															Computed:    true,
															Description: "List of policy statements defining the permissions.",
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
																			"  - example : ['iam:DeleteRole']",
																		ElementType: types.StringType,
																	},
																	"effect": schema.StringAttribute{
																		Computed: true,
																		Description: "Effect of the statement (allow or deny).\n" +
																			"  - example : 'Allow'",
																	},
																	"resource": schema.ListAttribute{
																		Computed: true,
																		Description: "List of resources the statement applies to.\n" +
																			"  - example : ['srn:e::123456789012:::iam:role/12345678']",
																		ElementType: types.StringType,
																	},
																	"sid": schema.StringAttribute{
																		Computed: true,
																		Description: "Statement ID for the statement.\n" +
																			"  - example : 'Stmt1'",
																	},
																	"condition": schema.MapAttribute{
																		ElementType: types.MapType{
																			ElemType: types.ListType{
																				ElemType: types.StringType,
																			},
																		},
																		Computed: true,
																		Description: "Condition for the statement.\n" +
																			"  - example : {'StringEquals': {'scp:PrincipalTag/department': ['engineering']}}",
																	},
																	"principal": schema.SingleNestedAttribute{
																		Computed:    true,
																		Description: "Principal that is allowed or denied access.",
																		Attributes: map[string]schema.Attribute{
																			"principal_string": schema.StringAttribute{
																				Computed: true,
																				Description: "String representation of the principal.\n" +
																					"  - example : 'srn:e::123456789012:::iam:user/12345678'",
																			},
																			"principal_map": schema.MapAttribute{
																				Computed: true,
																				Description: "Map of principal attributes.\n" +
																					"  - example : {'SCP': ['srn:e::123456789012:::iam:user/12345678']}",
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
															Description: "Policy document version.\n" +
																"  - example : '2024-07-01'",
														},
													},
												},
												"policy_id": schema.StringAttribute{
													Computed: true,
													Description: "ID of the policy this version belongs to.\n" +
														"  - example : 'pol-1234567890abcdef'",
												},
												"policy_version_name": schema.StringAttribute{
													Computed: true,
													Description: "Name of the policy version.\n" +
														"  - example : 'v1'",
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
										Description: "Name of the service the policy applies to.\n" +
											"  - example : 'compute'",
									},
									"service_type": schema.StringAttribute{
										Computed: true,
										Description: "Type of service the policy applies to.\n" +
											"  - example : 'EC2'",
									},
									"srn": schema.StringAttribute{
										Computed: true,
										Description: "Samsung Resource Name (SRN) of the policy.\n" +
											"  - example : 'srn:cloud:iam::123456789012:policy/my-policy'",
									},
									"state": schema.StringAttribute{
										Computed: true,
										Description: "State of the policy.\n" +
											"  - example : 'ACTIVE' | 'INACTIVE' | 'DELETED'",
									},
								},
							},
						},

						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the group was last modified.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who last modified the group.\n" +
								"  - example : 'user@example.com'",
						},
						"modifier_email": schema.StringAttribute{
							Computed: true,
							Description: "Email of the user who last modified the group.\n" +
								"  - example : 'user@example.com'",
						},
						"modifier_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the user who last modified the group.\n" +
								"  - example : 'John Doe'",
						},
						"name": schema.StringAttribute{
							Computed: true,
							Description: "Display name of the group.\n" +
								"  - example : 'MyGroup'",
						},
						"resource_type": schema.StringAttribute{
							Computed: true,
							Description: "Type of resource the group applies to.\n" +
								"  - example : 'group'",
						},
						"service_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the service the group applies to.\n" +
								"  - example : 'iam'",
						},
						"service_type": schema.StringAttribute{
							Computed: true,
							Description: "Type of service the group applies to.\n" +
								"  - example : 'IAM'",
						},
						"srn": schema.StringAttribute{
							Computed: true,
							Description: "Samsung Resource Name (SRN) of the group.\n" +
								"  - example : 'srn:cloud:iam::123456789012:group/my-group'",
						},
						"type": schema.StringAttribute{
							Computed: true,
							Description: "Type of group.\n" +
								"  - example : 'USER_DEFINED' | 'DEFAULT'",
						},
					},
				},
			},
		},
	}
}

func (d *iamGroupDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.GroupDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetGroups(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Groups",
			err.Error(),
		)
		return
	}

	for _, group := range data.Groups {

		// group members
		members := getGroupMembers(group.Members)

		// policies
		policies, hasError := getPolicies(ctx, group.Policies)
		if hasError {
			return
		}

		groupState := iam.Group{
			CreatedAt:     types.StringValue(group.CreatedAt.Format(time.RFC3339)),
			CreatedBy:     types.StringValue(group.CreatedBy),
			CreatorEmail:  types.StringPointerValue(group.CreatorEmail),
			CreatorName:   types.StringPointerValue(group.CreatorName),
			Description:   types.StringPointerValue(group.Description.Get()),
			DomainName:    types.StringValue(group.DomainName),
			Id:            types.StringValue(group.Id),
			Members:       members,
			Policies:      policies,
			ModifiedAt:    types.StringValue(group.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:    types.StringValue(group.ModifiedBy),
			ModifierEmail: types.StringPointerValue(group.ModifierEmail),
			ModifierName:  types.StringPointerValue(group.ModifierName),
			Name:          types.StringValue(group.Name),
			ResourceType:  types.StringPointerValue(group.ResourceType.Get()),
			ServiceName:   types.StringPointerValue(group.ServiceName.Get()),
			ServiceType:   types.StringPointerValue(group.ServiceType.Get()),
			Srn:           types.StringPointerValue(group.Srn.Get()),
			GroupType:     types.StringValue(group.Type),
		}

		state.Groups = append(state.Groups, groupState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
