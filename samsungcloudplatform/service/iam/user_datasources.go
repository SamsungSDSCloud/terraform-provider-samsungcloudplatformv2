package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scpsdkiam "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/iam/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iamUserDataSources{}
	_ datasource.DataSourceWithConfigure = &iamUserDataSources{}
)

func NewIamUserDataSources() datasource.DataSource { return &iamUserDataSources{} }

type iamUserDataSources struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (d *iamUserDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_users"
}

func (d *iamUserDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamUserDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Users.",
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
			common.ToSnakeCase("Email"): schema.StringAttribute{
				Description: "Email",
				Optional:    true,
			},
			common.ToSnakeCase("UserName"): schema.StringAttribute{
				Description: "User Name",
				Optional:    true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "User Type",
				Optional:    true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Account ID",
				Optional:    true,
			},
			common.ToSnakeCase("Users"): schema.ListNestedAttribute{
				Description: "A list of user.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Account ID",
							MarkdownDescription: "Account ID",
						},
						"company_name": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "Company Name",
							MarkdownDescription: "Company Name",
						},
						"console_url": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "Console URL",
							MarkdownDescription: "Console URL",
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
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "Description",
							MarkdownDescription: "Description",
						},
						"dst_offset": schema.StringAttribute{
							Computed:            true,
							Description:         "Dst Offset",
							MarkdownDescription: "Dst Offset",
						},
						"email": schema.StringAttribute{
							Computed:            true,
							Description:         "Email",
							MarkdownDescription: "Email",
						},
						"email_authenticated": schema.BoolAttribute{
							Computed:            true,
							Description:         "Email Authenticated",
							MarkdownDescription: "Email Authenticated",
						},
						"first_name": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "First Name",
							MarkdownDescription: "First Name",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "ID",
							MarkdownDescription: "ID",
						},
						"last_login_at": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "Last Login At",
							MarkdownDescription: "Last Login At",
						},
						"last_name": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "Last Name",
							MarkdownDescription: "Last Name",
						},
						"last_password_update_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Last Password Update At",
							MarkdownDescription: "Last Password Update At",
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
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Name",
							MarkdownDescription: "Name",
						},
						"password": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "Password",
							MarkdownDescription: "Password",
						},
						"password_reuse_count": schema.Int64Attribute{
							Computed:            true,
							Description:         "Password Reuse Count",
							MarkdownDescription: "Password Reuse Count",
						},
						"phone_authenticated": schema.BoolAttribute{
							Computed:            true,
							Description:         "Phone Authenticated",
							MarkdownDescription: "Phone Authenticated",
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
						"timezone": schema.StringAttribute{
							Computed:            true,
							Description:         "Timezone",
							MarkdownDescription: "Timezone",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "Type",
							MarkdownDescription: "Type",
						},
						"tz_id": schema.StringAttribute{
							Computed:            true,
							Description:         "TZ ID",
							MarkdownDescription: "TZ ID",
						},
						"user_name": schema.StringAttribute{
							Computed:            true,
							Description:         "User Name",
							MarkdownDescription: "User Name",
						},
						"utc_offset": schema.StringAttribute{
							Computed:            true,
							Description:         "UTC Offset",
							MarkdownDescription: "UTC Offset",
						},
					},
				},
			},
		},
	}
}

func (d *iamUserDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.UserDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetUsers(ctx, state.AccountId.ValueString(), state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Users",
			err.Error(),
		)
		return
	}

	// users
	users, hasError := getUsers(ctx, data.Users)
	if hasError {
		return
	}

	state.Users = users

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getUsers(ctx context.Context, _users []scpsdkiam.IAMUserResponse) ([]iam.User, bool) {
	var users []iam.User
	for _, user := range _users {

		// policies
		policies, _ := getPolicies(ctx, user.Policies)

		// user nil check
		userCompanyName := user.CompanyName.Get()
		if userCompanyName == nil {
			emptyStr := ""
			userCompanyName = &emptyStr
		}

		userFirstName := user.FirstName.Get()
		if userFirstName == nil {
			emptyStr := ""
			userFirstName = &emptyStr
		}

		userLastName := user.LastName.Get()
		if userLastName == nil {
			emptyStr := ""
			userLastName = &emptyStr
		}

		userLastLoginAt := user.LastLoginAt.Get()
		if userLastLoginAt == nil {
			emptyTime := time.Time{}
			userLastLoginAt = &emptyTime
		}
		userLastPasswordUpdateAt := user.LastPasswordUpdateAt.Get()
		if userLastPasswordUpdateAt == nil {
			emptyTime := time.Time{}
			userLastPasswordUpdateAt = &emptyTime
		}

		userState := iam.User{
			AccountId:            types.StringValue(*user.AccountId.Get()),
			CompanyName:          types.StringValue(*userCompanyName),
			CreatedAt:            types.StringValue(user.CreatedAt.Format(time.RFC3339)),
			CreatedBy:            types.StringValue(user.CreatedBy),
			Description:          types.StringValue(*user.Description.Get()),
			DstOffset:            types.StringValue(*user.DstOffset.Get()),
			Email:                types.StringValue(*user.Email.Get()),
			EmailAuthenticated:   types.BoolValue(user.EmailAuthenticated),
			FirstName:            types.StringValue(*userFirstName),
			Id:                   types.StringValue(user.Id),
			LastLoginAt:          types.StringValue(userLastLoginAt.Format(time.RFC3339)),
			LastName:             types.StringValue(*userLastName),
			LastPasswordUpdateAt: types.StringValue(userLastPasswordUpdateAt.Format(time.RFC3339)),
			ModifiedAt:           types.StringValue(user.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:           types.StringValue(user.ModifiedBy),
			Name:                 types.StringValue(*user.Name.Get()),
			PasswordReuseCount:   types.Int32Value(user.PasswordReuseCount),
			PhoneAuthenticated:   types.BoolValue(user.PhoneAuthenticated),
			Policies:             policies,
			Timezone:             types.StringValue(*user.Timezone.Get()),
			Type:                 types.StringValue(user.Type),
			TzId:                 types.StringValue(*user.TzId.Get()),
			UserName:             types.StringValue(*user.UserName.Get()),
			UtcOffset:            types.StringValue(*user.UtcOffset.Get()),
		}

		users = append(users, userState)
	}
	return users, false
}
