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
	_ datasource.DataSource              = &iamUserDataSource{}
	_ datasource.DataSourceWithConfigure = &iamUserDataSource{}
)

// NewIamUserDataSource is a helper function to simplify the provider implementation.
func NewIamUserDataSource() datasource.DataSource {
	return &iamUserDataSource{}
}

type iamUserDataSource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *iamUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_user"
}

// Configure adds the provider configured client to the data source.
func (d *iamUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show IAM User",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Optional: true,
				Description: "Account ID to filter users.\n" +
					"  - example : '123456789012'",
			},
			"user_id": schema.StringAttribute{
				Optional: true,
				Description: "Unique identifier of the user to retrieve.\n" +
					"  - example : 'usr-1234567890abcdef'",
			},
			"user": schema.SingleNestedAttribute{
				Description: "Detailed information about the user.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed: true,
						Description: "Account ID that owns the user.\n" +
							"  - example : '123456789012'",
					},
					"company_name": schema.StringAttribute{
						Computed: true,
						Description: "Company name of the user.\n" +
							"  - example : 'Samsung SDS'",
					},
					"console_url": schema.StringAttribute{
						Computed: true,
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
						Description: "Unique identifier.\n" +
							"  - example : 'usr-1234567890abcdef'",
					},
					"last_login_at": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Description: "Last login timestamp.\n" +
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
						Description: "User name.\n" +
							"  - example : 'John Doe'",
					},
					"password": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Description: "User password (masked for security).\n" +
							"  - example : '********'",
						MarkdownDescription: "User password (masked for security).\n" +
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
						Optional:    true,
						Description: "Policies",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"account_id": schema.StringAttribute{
									Computed:            true,
									Description:         "Account ID that owns the policy.\n  - example : '123456789012'",
									MarkdownDescription: "Account ID that owns the policy.\n  - example : '123456789012'",
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
									Description: "Email of the user who created the policy.\n" +
										"  - example : 'user@example.com'",
								},
								"creator_name": schema.StringAttribute{
									Computed: true,
									Description: "Name of the user who created the policy.\n" +
										"  - example : 'John Doe'",
								},
								"default_version_id": schema.StringAttribute{
									Computed:            true,
									Description:         "Default version ID of the policy.\n  - example : 'pol-1234567890abcdef'",
									MarkdownDescription: "Default version ID of the policy.\n  - example : 'pol-1234567890abcdef'",
								},
								"description": schema.StringAttribute{
									Computed:            true,
									Description:         "Description of the policy.\n  - example : 'My policy description'",
									MarkdownDescription: "Description of the policy.\n  - example : 'My policy description'",
								},
								"domain_name": schema.StringAttribute{
									Computed:            true,
									Description:         "Domain name associated with the policy.\n  - example : 'scp'",
									MarkdownDescription: "Domain name associated with the policy.\n  - example : 'scp'",
								},
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Policy ID.\n  - example : 'pol-1234567890abcdef'",
									MarkdownDescription: "Policy ID.\n  - example : 'pol-1234567890abcdef'",
								},
								"modified_at": schema.StringAttribute{
									Computed:            true,
									Description:         "Timestamp when the policy was last modified.\n  - example : '2024-01-01T00:00:00Z'",
									MarkdownDescription: "Timestamp when the policy was last modified.\n  - example : '2024-01-01T00:00:00Z'",
								},
								"modified_by": schema.StringAttribute{
									Computed:            true,
									Description:         "User who last modified the policy.\n  - example : 'user@example.com'",
									MarkdownDescription: "User who last modified the policy.\n  - example : 'user@example.com'",
								},
								"modifier_email": schema.StringAttribute{
									Computed:            true,
									Description:         "Email of the user who last modified the policy.\n  - example : 'user@example.com'",
									MarkdownDescription: "Email of the user who last modified the policy.\n  - example : 'user@example.com'",
								},
								"modifier_name": schema.StringAttribute{
									Computed:            true,
									Description:         "Name of the user who last modified the policy.\n  - example : 'John Doe'",
									MarkdownDescription: "Name of the user who last modified the policy.\n  - example : 'John Doe'",
								},
								"policy_category": schema.StringAttribute{
									Computed:            true,
									Description:         "Category of the policy.\n  - example : 'IDENTITY_BASED'",
									MarkdownDescription: "Category of the policy.\n  - example : 'IDENTITY_BASED'",
								},
								"policy_name": schema.StringAttribute{
									Computed:            true,
									Description:         "Name of the policy.\n  - example : 'MyPolicy'",
									MarkdownDescription: "Name of the policy.\n  - example : 'MyPolicy'",
								},
								"policy_type": schema.StringAttribute{
									Computed:            true,
									Description:         "Type of the policy.\n  - example : 'USER_DEFINED'",
									MarkdownDescription: "Type of the policy.\n  - example : 'USER_DEFINED'",
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
											},
											"created_by": schema.StringAttribute{
												Computed: true,
												Description: "User who created the policy version.\n" +
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
												Description: "User who last modified this policy version.\n" +
													"  - example : 'user@example.com'",
											},
											"policy_document": schema.SingleNestedAttribute{
												Computed:            true,
												Description:         "Policy document - JSON policy content defining permissions.\n  - example : {'Version': '2012-10-17', 'Statement': [...]}",
												MarkdownDescription: "Policy document - JSON policy content defining permissions.\n  - example : {'Version': '2012-10-17', 'Statement': [...]}",
												Attributes: map[string]schema.Attribute{
													"statement": schema.ListNestedAttribute{
														Computed:            true,
														Description:         "Statement - list of permission statements in the policy.\n  - example : [{'Sid': 'Stmt1', 'Effect': 'Allow', 'Action': [...], 'Resource': '*'}]",
														MarkdownDescription: "Statement - list of permission statements in the policy.\n  - example : [{'Sid': 'Stmt1', 'Effect': 'Allow', 'Action': [...], 'Resource': '*'}]",
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{
																"action": schema.ListAttribute{
																	Computed:            true,
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
																	Description:         "Effect of the policy statement (Allow or Deny).\n  - example : 'Allow'",
																	MarkdownDescription: "Effect of the policy statement (Allow or Deny).\n  - example : 'Allow'",
																},
																"resource": schema.ListAttribute{
																	Computed:            true,
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
												Computed:            true,
												Description:         "ID of the policy.\n  - example : 'pol-1234567890abcdef'",
												MarkdownDescription: "ID of the policy.\n  - example : 'pol-1234567890abcdef'",
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
									Description:         "Type of resource.\n  - example : 'policy'",
									MarkdownDescription: "Type of resource.\n  - example : 'policy'",
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
									Computed: true,
									Description: "Service Resource Name (SRN).\n" +
										"  - example : 'srn:e:::::iam:policy/policy-12345678'",
								},
								"state": schema.StringAttribute{
									Computed: true,
									Description: "User state.\n" +
										"  - example : 'ACTIVE'",
								},
							},
						},
					},
					"timezone": schema.StringAttribute{
						Computed: true,
						Description: "Timezone of the user.\n" +
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
						Description: "Unique username.\n" +
							"  - example : 'john.doe'",
					},
					"utc_offset": schema.StringAttribute{
						Computed: true,
						Description: "User's UTC offset from UTC time.\n" +
							"  - example : '+09:00'",
					},
					"access_keys": schema.ListNestedAttribute{
						Computed:            true,
						Description:         "Access Keys",
						MarkdownDescription: "Access Keys",
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
						Computed:            true,
						Description:         "Groups",
						MarkdownDescription: "Groups",
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

func (d *iamUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.UserDataSourceDetail

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetUser(ctx, state.AccountId.ValueString(), state.UserId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read User",
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

	// mapped access key info
	var accessKeyInfos []iam.AccessKeyV1Dot4
	for _, accessKeyInfo := range data.AccessKeys {
		accessKeyInfos = append(accessKeyInfos, iam.AccessKeyV1Dot4{
			AccessKey:           types.StringValue(accessKeyInfo.AccessKey),
			CreatedAt:           types.StringValue(accessKeyInfo.CreatedAt.Format(time.RFC3339)),
			ExpirationTimestamp: types.StringValue(accessKeyInfo.ExpirationTimestamp.Format(time.RFC3339)),
			Id:                  types.StringValue(accessKeyInfo.Id),
			IsEnabled:           types.BoolValue(accessKeyInfo.IsEnabled),
		})
	}

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
	}

	userObjectValue, _ := types.ObjectValueFrom(ctx, userState.AttributeTypes(), userState)
	state.User = userObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
