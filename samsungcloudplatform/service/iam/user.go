package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/importstate"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
const CreatedAt = "Created At"

var (
	_ resource.Resource                = &iamUserResource{}
	_ resource.ResourceWithConfigure   = &iamUserResource{}
	_ resource.ResourceWithImportState = &iamUserResource{}
)

// NewIamUserResource is a helper function to simplify the provider implementation.
func NewIamUserResource() resource.Resource {
	return &iamUserResource{}
}

// iamUserResource is the data source implementation.
type iamUserResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *iamUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_user"
}

func (r *iamUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM User.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Optional: true,
				Description: "Account ID to create the user in.\n" +
					"  - example : '123456789012'",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Description: "Human-readable description of the user.\n" +
					"  - example : 'My user description'",
			},
			"group_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of group IDs to add the user to.\n" +
					"  - example : ['grp-1234567890abcdef']",
			},
			"policy_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of policy IDs to attach to the user.\n" +
					"  - example : ['pol-1234567890abcdef']",
			},
			"password": schema.StringAttribute{
				Optional: true,
				Description: "Password for the user.\n" +
					"  - example : 'ExamplePassword123!'",
			},
			"tags": tag.ResourceSchema(),
			"temporary_password": schema.BoolAttribute{
				Optional: true,
				Description: "Whether the password is temporary and needs to be changed.\n" +
					"  - example : true",
			},
			"user_name": schema.StringAttribute{
				Optional: true,
				Description: "Unique username for the user.\n" +
					"  - example : 'john.doe'\n" +
					"  - maxLength: 64",
			},
			"password_reuse_count": schema.Int32Attribute{
				Optional: true,
				Description: "Number of previous passwords that cannot be reused.\n" +
					"  - example : 3\n" +
					"  - min: 0, max: 10",
			},
			"user_id": schema.StringAttribute{
				Computed: true,
				Description: "Unique identifier of the user.\n" +
					"  - example : 'usr-1234567890abcdef'",
			},
			"user": schema.SingleNestedAttribute{
				Description: "A detail of User.\n" +
					"  - example : '{account_id: 123456789012, company_name: Samsung SDS, created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, ...}'",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed: true,
						Description: "Account ID associated with the user.\n" +
							"  - example : '123456789012'",
					},
					"company_name": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Description: "Company name of the user.\n" +
							"  - example : 'Samsung SDS'",
					},
					"console_url": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Description: "URL to access the console.\n" +
							"  - example : 'https://console.example.com'",
					},
					"created_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the user was created.\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"created_by": schema.StringAttribute{
						Computed: true,
						Description: "User who created the user.\n" +
							"  - example : 'user@example.com'",
					},
					"description": schema.StringAttribute{
						Computed: true,
						Description: "Human-readable description of the user.\n" +
							"  - example : 'My user description'",
					},
					"dst_offset": schema.StringAttribute{
						Computed: true,
						Description: "Daylight saving time offset.\n" +
							"  - example : '+09:00'",
					},
					"email": schema.StringAttribute{
						Computed: true,
						Description: "Email address.\n" +
							"  - example : 'user@example.com'",
					},
					"email_authenticated": schema.BoolAttribute{
						Computed: true,
						Description: "Whether email is authenticated.\n" +
							"  - example : true",
					},
					"first_name": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Description: "First name of the user.\n" +
							"  - example : 'John'",
					},
					"id": schema.StringAttribute{
						Computed: true,
						Description: "Unique identifier of the user.\n" +
							"  - example : 'usr-1234567890abcdef'",
					},
					"last_login_at": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Description: "Timestamp when the user last logged in.\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"last_name": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Description: "Last name of the user.\n" +
							"  - example : 'Doe'",
					},
					"last_password_update_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the password was last updated.\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the user was last modified.\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
						Description: "User who last modified the user.\n" +
							"  - example : 'user@example.com'",
					},
					"name": schema.StringAttribute{
						Computed: true,
						Description: "Display name of the user.\n" +
							"  - example : 'John Doe'",
					},
					"password": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Description: "Password for the user account.\n" +
							"  - example : '********'",
					},
					"password_reuse_count": schema.Int64Attribute{
						Computed: true,
						Description: "Number of previous passwords that cannot be reused.\n" +
							"  - example : 3",
					},
					"phone_authenticated": schema.BoolAttribute{
						Computed: true,
						Description: "Whether the phone number has been authenticated.\n" +
							"  - example : true",
					},
					"policies": schema.ListNestedAttribute{
						Optional: true,
						Description: "List of policies attached to the role.\n" +
							"  - example : '[{account_id: 123456789012, created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, default_version_id: pol-1234567890abcdef, ...}]'",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"account_id": schema.StringAttribute{
									Optional: true,
									Description: "Account ID associated with the policy.\n" +
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
									Description: "Email address of the policy creator.\n" +
										"  - example : 'creator@example.com'",
								},
								"creator_name": schema.StringAttribute{
									Computed: true,
									Description: "Name of the policy creator.\n" +
										"  - example : 'Creator Name'",
								},
								"default_version_id": schema.StringAttribute{
									Computed: true,
									Description: "Default version ID of the policy.\n" +
										"  - example : 'pol-1234567890abcdef'",
								},
								"description": schema.StringAttribute{
									Computed: true,
									Description: "Description of the policy.\n" +
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
									Description: "Email address of the user who last modified the policy.\n" +
										"  - example : 'modifier@example.com'",
								},
								"modifier_name": schema.StringAttribute{
									Computed: true,
									Description: "Name of the user who last modified the policy.\n" +
										"  - example : 'Modifier Name'",
								},
								"policy_category": schema.StringAttribute{
									Computed: true,
									Description: "Category of the policy.\n" +
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
										"  - example : '[{created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, id: ver-1234567890abcdef, modified_at: 2024-05-17T00:23:17Z, modified_by: ef50cdc207f05f6fb8f20219f229ed1f, ...}]'",
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
																	Description: "List of actions allowed by this statement (e.g., iam:CreateRole, iam:ListUsers).\n" +
																		"  - example : ['iam:CreateRole']",
																	ElementType: types.StringType,
																},
																"not_action": schema.ListAttribute{
																	Optional: true,
																	Description: "List of actions that are explicitly excluded from this statement.\n" +
																		"  - example : ['iam:DeleteRole']",
																	ElementType: types.StringType,
																},
																"effect": schema.StringAttribute{
																	Computed: true,
																	Description: "Effect of the statement - either Allow or Deny.\n" +
																		"  - example : 'Allow'",
																},
																"resource": schema.ListAttribute{
																	Optional: true,
																	Description: "List of resources (ARNs or wildcards) that the statement applies to.\n" +
																		"  - example : ['*']",
																	ElementType: types.StringType,
																},
																"sid": schema.StringAttribute{
																	Computed: true,
																	Description: "Statement ID (SID) - unique identifier for this policy statement.\n" +
																		"  - example : 'Stmt1'",
																},
																"condition": schema.MapAttribute{
																	ElementType: types.MapType{
																		ElemType: types.ListType{
																			ElemType: types.StringType,
																		},
																	},
																	Optional: true,
																	Description: "Conditions that must be met for the policy statement to take effect.\n" +
																		"  - example : {'StringEquals': {'aws:PrincipalTag/department': ['IT']}}",
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
									Description: "Service Resource Name (SRN) - Unique identifier for the user in the SCP system.\n" +
										"  - example : 'srn:e:::::iam:policy/policy-12345678'",
								},
								"state": schema.StringAttribute{
									Computed: true,
									Description: "Current state of the user (e.g., ACTIVE, INACTIVE).\n" +
										"  - example : 'ACTIVE'",
								},
							},
						},
					},
					"timezone": schema.StringAttribute{
						Computed: true,
						Description: "User's timezone setting.\n" +
							"  - example : 'Asia/Seoul'",
					},
					"type": schema.StringAttribute{
						Computed: true,
						Description: "Type of user.\n" +
							"  - example : 'IAM'",
					},
					"tz_id": schema.StringAttribute{
						Computed: true,
						Description: "Timezone ID.\n" +
							"  - example : 'Asia/Seoul'",
					},
					"user_name": schema.StringAttribute{
						Computed: true,
						Description: "Unique username of the user.\n" +
							"  - example : 'john.doe'",
					},
					"utc_offset": schema.StringAttribute{
						Computed: true,
						Description: "User's UTC offset from UTC time.\n" +
							"  - example : '+09:00'",
					},
					"access_keys": schema.ListNestedAttribute{
						Computed: true,
						Description: "List of access keys associated with the user.\n" +
							"  - example : '[{access_key: ak-example-access-key-id, created_at: 2024-05-17T00:23:17Z, expiration_timestamp: 9999-12-31T23:59:59Z, id: 12345678-1234-1234-1234-1234567890ab, is_enabled: true}]'",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"access_key": schema.StringAttribute{
									Computed: true,
									Description: "The access key string value.\n" +
										"  - example : 'ak-example-access-key-id'",
								},
								"created_at": schema.StringAttribute{
									Computed: true,
									Description: "Timestamp when the access key was created.\n" +
										"  - example : '2024-01-01T00:00:00Z'",
								},
								"expiration_timestamp": schema.StringAttribute{
									Computed: true,
									Description: "Timestamp when the access key expires.\n" +
										"  - example : '2024-01-02T00:00:00Z'",
								},
								"id": schema.StringAttribute{
									Computed: true,
									Description: "Unique identifier of the access key.\n" +
										"  - example : '12345678-1234-1234-1234-1234567890ab'",
								},
								"is_enabled": schema.BoolAttribute{
									Computed: true,
									Description: "Whether the access key is enabled/active.\n" +
										"  - example : true",
								},
							},
						},
					},
					"groups": schema.ListNestedAttribute{
						Computed: true,
						Description: "List of groups the user belongs to.\n" +
							"  - example : '[{id: grp-1234567890abcdef, name: MyGroup}]'",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed: true,
									Description: "Unique identifier of the group.\n" +
										"  - example : 'grp-1234567890abcdef'",
								},
								"name": schema.StringAttribute{
									Computed: true,
									Description: "Display name of the group.\n" +
										"  - example : 'MyGroup'",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *iamUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importstate.ImportState(ctx, req, resp,
		path.Root("account_id"),
		path.Root("user_id"),
	)
}

func (r *iamUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.UserResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateUser(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	// user nil check
	userCompanyName := data.CompanyName.Get()
	if userCompanyName == nil {
		emptyStr := ""
		userCompanyName = &emptyStr
	}

	userFirstName := data.FirstName.Get()
	if userFirstName == nil {
		emptyStr := ""
		userFirstName = &emptyStr
	}

	userLastName := data.LastName.Get()
	if userLastName == nil {
		emptyStr := ""
		userLastName = &emptyStr
	}

	userLastLoginAt := data.LastLoginAt.Get()
	if userLastLoginAt == nil {
		emptyTime := time.Time{}
		userLastLoginAt = &emptyTime
	}

	userLastPasswordUpdateAt := data.LastPasswordUpdateAt.Get()
	if userLastPasswordUpdateAt == nil {
		emptyTime := time.Time{}
		userLastPasswordUpdateAt = &emptyTime
	}

	plan.UserId = types.StringValue(data.Id)

	// empty list
	accessKeyInfos := make([]iam.AccessKeyV1Dot4, 0)
	groupInfos := make([]iam.GroupInfo, 0)

	userState := iam.User{
		AccountId:            types.StringValue(*data.AccountId.Get()),
		CompanyName:          types.StringValue(*userCompanyName),
		ConsoleUrl:           types.StringValue(data.ConsoleUrl),
		CreatedAt:            types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(data.CreatedBy),
		Description:          types.StringValue(*data.Description.Get()),
		DstOffset:            types.StringValue(*data.DstOffset.Get()),
		Email:                types.StringValue(*data.Email.Get()),
		EmailAuthenticated:   types.BoolValue(data.EmailAuthenticated),
		FirstName:            types.StringValue(*userFirstName),
		Id:                   types.StringValue(data.Id),
		LastLoginAt:          types.StringValue(userLastLoginAt.Format(time.RFC3339)),
		LastName:             types.StringValue(*userLastName),
		LastPasswordUpdateAt: types.StringValue(userLastPasswordUpdateAt.Format(time.RFC3339)),
		ModifiedAt:           types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(data.ModifiedBy),
		Name:                 types.StringValue(*data.Name.Get()),
		Password:             types.StringValue(data.Password),
		PasswordReuseCount:   types.Int32Value(data.PasswordReuseCount),
		PhoneAuthenticated:   types.BoolValue(data.PhoneAuthenticated),
		Policies:             policies,
		Timezone:             types.StringValue(*data.Timezone.Get()),
		Type:                 types.StringValue(data.Type),
		TzId:                 types.StringValue(*data.TzId.Get()),
		UserName:             types.StringValue(*data.UserName.Get()),
		UtcOffset:            types.StringValue(*data.UtcOffset.Get()),
		AccessKeys:           accessKeyInfos,
		Groups:               groupInfos,
	}

	userObjectValue, diags := types.ObjectValueFrom(ctx, userState.AttributeTypes(), userState)
	plan.User = userObjectValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.UserResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetUser(ctx, state.AccountId.ValueString(), state.UserId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Show User",
			err.Error(),
		)
		return
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	// user nil check
	userCompanyName := data.CompanyName.Get()
	if userCompanyName == nil {
		emptyStr := ""
		userCompanyName = &emptyStr
	}

	userFirstName := data.FirstName.Get()
	if userFirstName == nil {
		emptyStr := ""
		userFirstName = &emptyStr
	}

	userLastName := data.LastName.Get()
	if userLastName == nil {
		emptyStr := ""
		userLastName = &emptyStr
	}

	userLastLoginAt := data.LastLoginAt.Get()
	if userLastLoginAt == nil {
		emptyTime := time.Time{}
		userLastLoginAt = &emptyTime
	}

	userLastPasswordUpdateAt := data.LastPasswordUpdateAt.Get()
	if userLastPasswordUpdateAt == nil {
		emptyTime := time.Time{}
		userLastPasswordUpdateAt = &emptyTime
	}

	// empty list
	accessKeyInfos := make([]iam.AccessKeyV1Dot4, 0)
	groupInfos := make([]iam.GroupInfo, 0)

	userState := iam.User{
		AccountId:            types.StringValue(*data.AccountId.Get()),
		CompanyName:          types.StringValue(*userCompanyName),
		CreatedAt:            types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(data.CreatedBy),
		Description:          types.StringValue(*data.Description.Get()),
		DstOffset:            types.StringValue(*data.DstOffset.Get()),
		Email:                types.StringValue(*data.Email.Get()),
		EmailAuthenticated:   types.BoolValue(data.EmailAuthenticated),
		FirstName:            types.StringValue(*userFirstName),
		Id:                   types.StringValue(data.Id),
		LastLoginAt:          types.StringValue(userLastLoginAt.Format(time.RFC3339)),
		LastName:             types.StringValue(*userLastName),
		LastPasswordUpdateAt: types.StringValue(userLastPasswordUpdateAt.Format(time.RFC3339)),
		ModifiedAt:           types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(data.ModifiedBy),
		Name:                 types.StringValue(*data.Name.Get()),
		PasswordReuseCount:   types.Int32Value(data.PasswordReuseCount),
		PhoneAuthenticated:   types.BoolValue(data.PhoneAuthenticated),
		Policies:             policies,
		Timezone:             types.StringValue(*data.Timezone.Get()),
		Type:                 types.StringValue(data.Type),
		TzId:                 types.StringValue(*data.TzId.Get()),
		UserName:             types.StringValue(*data.UserName.Get()),
		UtcOffset:            types.StringValue(*data.UtcOffset.Get()),
		AccessKeys:           accessKeyInfos,
		Groups:               groupInfos,
	}

	userObjectValue, diags := types.ObjectValueFrom(ctx, userState.AttributeTypes(), userState)
	state.User = userObjectValue
	state.UserId = types.StringValue(data.Id)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iam.UserResource
	var state iam.UserResource

	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)

	_, err := r.client.UpdateUser(ctx, state.AccountId.ValueString(), state.UserId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating User",
			"Could not update User, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetUser(ctx, state.AccountId.ValueString(), state.UserId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Unable to Read User",
			"Could not read User ID "+state.UserId.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	// user nil check
	userCompanyName := data.CompanyName.Get()
	if userCompanyName == nil {
		emptyStr := ""
		userCompanyName = &emptyStr
	}

	userFirstName := data.FirstName.Get()
	if userFirstName == nil {
		emptyStr := ""
		userFirstName = &emptyStr
	}

	userLastName := data.LastName.Get()
	if userLastName == nil {
		emptyStr := ""
		userLastName = &emptyStr
	}

	userLastLoginAt := data.LastLoginAt.Get()
	if userLastLoginAt == nil {
		emptyTime := time.Time{}
		userLastLoginAt = &emptyTime
	}

	userLastPasswordUpdateAt := data.LastPasswordUpdateAt.Get()
	if userLastPasswordUpdateAt == nil {
		emptyTime := time.Time{}
		userLastPasswordUpdateAt = &emptyTime
	}

	// empty list
	accessKeyInfos := make([]iam.AccessKeyV1Dot4, 0)
	groupInfos := make([]iam.GroupInfo, 0)

	userState := iam.User{
		AccountId:            types.StringValue(*data.AccountId.Get()),
		CompanyName:          types.StringValue(*userCompanyName),
		CreatedAt:            types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(data.CreatedBy),
		Description:          types.StringValue(*data.Description.Get()),
		DstOffset:            types.StringValue(*data.DstOffset.Get()),
		Email:                types.StringValue(*data.Email.Get()),
		EmailAuthenticated:   types.BoolValue(data.EmailAuthenticated),
		FirstName:            types.StringValue(*userFirstName),
		Id:                   types.StringValue(data.Id),
		LastLoginAt:          types.StringValue(userLastLoginAt.Format(time.RFC3339)),
		LastName:             types.StringValue(*userLastName),
		LastPasswordUpdateAt: types.StringValue(userLastPasswordUpdateAt.Format(time.RFC3339)),
		ModifiedAt:           types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(data.ModifiedBy),
		Name:                 types.StringValue(*data.Name.Get()),
		PasswordReuseCount:   types.Int32Value(data.PasswordReuseCount),
		PhoneAuthenticated:   types.BoolValue(data.PhoneAuthenticated),
		Policies:             policies,
		Timezone:             types.StringValue(*data.Timezone.Get()),
		Type:                 types.StringValue(data.Type),
		TzId:                 types.StringValue(*data.TzId.Get()),
		UserName:             types.StringValue(*data.UserName.Get()),
		UtcOffset:            types.StringValue(*data.UtcOffset.Get()),
		AccessKeys:           accessKeyInfos,
		Groups:               groupInfos,
	}

	userObjectValue, diags := types.ObjectValueFrom(ctx, userState.AttributeTypes(), userState)
	plan.User = userObjectValue
	plan.UserId = state.UserId

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam.UserResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(ctx, state.AccountId.ValueString(), state.UserId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam User",
			"Could not delete User, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
