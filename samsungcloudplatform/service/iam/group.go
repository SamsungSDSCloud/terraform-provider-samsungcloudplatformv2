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
	scpsdkiam "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/iam/1.4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &iamGroupResource{}
	_ resource.ResourceWithConfigure   = &iamGroupResource{}
	_ resource.ResourceWithImportState = &iamGroupResource{}
)

// NewIamGroupResource is a helper function to simplify the provider implementation.
func NewIamGroupResource() resource.Resource {
	return &iamGroupResource{}
}

// iamGroupResource is the data source implementation.
type iamGroupResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *iamGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group"
}

func (r *iamGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "Unique identifier of the group.\n" +
					"  - example : 'grp-1234567890abcdef'",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
				Description: "Display name of the group.\n" +
					"  - example : 'MyGroup'\n" +
					"  - maxLength: 24\n" +
					"  - minLength: 3",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Description: "Human-readable description of the group.\n" +
					"  - example : 'My group description'",
			},
			"tags": tag.ResourceSchema(),
			"policy_ids": schema.ListAttribute{
				Optional: true,
				Description: "List of policy IDs to attach to the group.\n" +
					"  - example : ['pol-1234567890abcdef']",
				ElementType: types.StringType,
			},
			"user_ids": schema.ListAttribute{
				Optional: true,
				Description: "List of user IDs to add as members of the group.\n" +
					"  - example : ['usr-1234567890abcdef']",
				ElementType: types.StringType,
			},
			"group": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Detailed information about the group.\n" +
					"  - example : '{created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, description: Descriptions for group, domain_name: scp, id: f39c460fade34fecb05ede8f904b24b7, name: TestGroup, type: USER_DEFINED, ...}'",
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
						Computed: true,
						Description: "List of members in the group.\n" +
							"  - example : '[{created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, user_id: f39c460fade34fecb05ede8f904b24b7, user_name: -, ...}]'",
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
						Description: "List of policies attached to the group.\n" +
							"  - example : '[{account_id: 123456789012, created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, ...}]'",
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
									Computed: true,
									Description: "List of versions of the policy.\n" +
										"  - example : '[{created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, id: v-1234567890abcdef, modified_at: 2024-05-17T00:23:17Z, modified_by: ef50cdc207f05f6fb8f20219f229ed1f, ...}]'",
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
												Computed: true,
												Description: "The policy document containing the permission definitions/\n" +
													"  - example : '{statement: [{action: [iam:CreateRole], effect: Allow, resource: [...], ...}]}'",
												Attributes: map[string]schema.Attribute{
													"statement": schema.ListNestedAttribute{
														Computed: true,
														Description: "List of policy statements defining the permissions.\n" +
															"  - example : '[{action: [iam:CreateRole], effect: Allow, resource: [srn:e::123456789012:::iam:role/12345678], ...}]'",
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
																	Description: "Condition for the policy statement.\n" +
																		"  - example : {'StringEquals': {'scp:PrincipalTag/department': ['engineering']}}",
																},
																"principal": schema.SingleNestedAttribute{
																	Computed: true,
																	Description: "Principal that is allowed or denied access.\n" +
																		"  - example : '{principal_string: srn:e::123456789012:::iam:user/12345678, principal_map: {SCP: [srn:e::123456789012:::iam:user/12345678]}}'",
																	Attributes: map[string]schema.Attribute{
																		"principal_string": schema.StringAttribute{
																			Computed: true,
																			Description: "Principal string identifier.\n" +
																				"  - example : 'srn:e::123456789012:::iam:user/12345678'",
																		},
																		"principal_map": schema.MapAttribute{
																			Computed: true,
																			Description: "Map of principal identifiers.\n" +
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
	}
}

func (r *iamGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ImportState imports the resource by ID.
func (r *iamGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *iamGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan iam.GroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateGroup(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating group",
			"Could not create group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// group
	group := data.Group

	// group members
	members := getGroupMembers(group.Members)

	// policies
	policies, hasError := getPolicies(ctx, group.Policies)
	if hasError {
		return
	}

	plan.Id = types.StringValue(group.Id)

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
	groupObjectValue, diags := types.ObjectValueFrom(ctx, groupState.AttributeTypes(), groupState)
	plan.Group = groupObjectValue
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *iamGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state iam.GroupResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Group
	_, err := r.client.UpdateGroup(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Group",
			"Could not update Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Unable to Read Group",
			"Could not read group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// group
	group := data.Group

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
	groupObjectValue, diags := types.ObjectValueFrom(ctx, groupState.AttributeTypes(), groupState)
	state.Group = groupObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state iam.GroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing iam group
	err := r.client.DeleteGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam group",
			"Could not delete Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *iamGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.GroupResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetGroup(ctx, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Show Group",
			err.Error()+"\nReason: "+detail,
		)
		return
	}

	// group
	group := data.Group

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

	groupObjectValue, diags := types.ObjectValueFrom(ctx, groupState.AttributeTypes(), groupState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Group = groupObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getGroupMembers(_members []scpsdkiam.GroupMember) []iam.Member {
	var members []iam.Member

	for _, member := range _members {
		var creatorLastLoginAt *string
		var userLastLoginAt *string

		if member.CreatorLastLoginAt.Get() != nil {
			t := member.CreatorLastLoginAt.Get().Format(time.RFC3339)
			creatorLastLoginAt = &t
		}
		if member.UserLastLoginAt.Get() != nil {
			t := member.UserLastLoginAt.Get().Format(time.RFC3339)
			userLastLoginAt = &t
		}

		groupNames := make([]types.String, 0, len(member.GroupNames))
		for _, groupName := range member.GroupNames {
			groupNames = append(groupNames, types.StringValue(groupName))
		}

		memberState := iam.Member{
			CreatedAt:          types.StringValue(member.CreatedAt.Format(time.RFC3339)),
			CreatedBy:          types.StringValue(member.CreatedBy),
			CreatorCreatedAt:   types.StringValue(member.CreatorCreatedAt.Format(time.RFC3339)),
			CreatorEmail:       types.StringPointerValue(member.CreatorEmail),
			CreatorLastLoginAt: types.StringPointerValue(creatorLastLoginAt),
			CreatorName:        types.StringPointerValue(member.CreatorName),
			GroupNames:         groupNames,
			UserCreatedAt:      types.StringValue(member.UserCreatedAt.Format(time.RFC3339)),
			UserEmail:          types.StringPointerValue(member.UserEmail),
			UserId:             types.StringValue(member.UserId),
			UserLastLoginAt:    types.StringPointerValue(userLastLoginAt),
			UserName:           types.StringPointerValue(member.UserName),
		}

		members = append(members, memberState)
	}
	return members
}

func getPolicies(ctx context.Context, _policies interface{}) ([]iam.Policy, bool) {
	var policies []iam.Policy

	switch v := _policies.(type) {
	case []scpsdkiam.Policy:
		policies = _getPolicies(ctx, v, policies)
	}

	return policies, false
}

func _getPolicies(ctx context.Context, _policies []scpsdkiam.Policy, policies []iam.Policy) []iam.Policy {
	for _, policy := range _policies {

		var policyVersions []iam.PolicyVersion
		//policy versions
		for _, policyVersion := range policy.PolicyVersions {

			var statements []iam.Statement
			for _, _statement := range policyVersion.PolicyDocument.Statement {

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

			policyDocument := iam.PolicyDocument{
				Version:   types.StringValue(policyVersion.PolicyDocument.Version),
				Statement: statements,
			}

			policyVersionState := iam.PolicyVersion{
				CreatedAt:         types.StringValue(policyVersion.CreatedAt.Format(time.RFC3339)),
				CreatedBy:         types.StringValue(policyVersion.CreatedBy),
				Id:                types.StringValue(*policyVersion.Id),
				ModifiedAt:        types.StringValue(policyVersion.ModifiedAt.Format(time.RFC3339)),
				ModifiedBy:        types.StringValue(policyVersion.ModifiedBy),
				PolicyDocument:    policyDocument,
				PolicyId:          types.StringValue(*policyVersion.PolicyId),
				PolicyVersionName: types.StringValue(*policyVersion.PolicyVersionName),
			}
			policyVersions = append(policyVersions, policyVersionState)

		}

		policyState := iam.Policy{
			AccountId:        types.StringPointerValue(policy.AccountId.Get()),
			CreatedAt:        types.StringValue(policy.CreatedAt.Format(time.RFC3339)),
			CreatedBy:        types.StringValue(policy.CreatedBy),
			CreatorEmail:     types.StringPointerValue(policy.CreatorEmail.Get()),
			CreatorName:      types.StringPointerValue(policy.CreatorName.Get()),
			DefaultVersionId: types.StringValue(*policy.DefaultVersionId),
			Description:      types.StringPointerValue(policy.Description.Get()),
			DomainName:       types.StringValue(policy.DomainName),
			Id:               types.StringValue(*policy.Id),
			ModifiedAt:       types.StringValue(policy.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:       types.StringValue(policy.ModifiedBy),
			ModifierEmail:    types.StringPointerValue(policy.ModifierEmail.Get()),
			ModifierName:     types.StringPointerValue(policy.ModifierName.Get()),
			PolicyCategory:   types.StringValue(string(*policy.PolicyCategory)),
			PolicyName:       types.StringValue(*policy.PolicyName),
			PolicyType:       types.StringValue(string(*policy.PolicyType)),
			PolicyVersions:   policyVersions,
			ResourceType:     types.StringPointerValue(policy.ResourceType.Get()),
			ServiceName:      types.StringPointerValue(policy.ServiceName.Get()),
			ServiceType:      types.StringPointerValue(policy.ServiceType.Get()),
			Srn:              types.StringValue(policy.Srn),
			State:            types.StringValue(string(*policy.State)),
		}

		policies = append(policies, policyState)
	}
	return policies
}
