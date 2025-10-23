package iam

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-iam"

type AccessKeyDataSource struct {
	Limit      types.Int32  `tfsdk:"limit"`
	Marker     types.String `tfsdk:"marker"`
	Sort       types.String `tfsdk:"sort"`
	AccountId  types.String `tfsdk:"account_id"`
	AccessKeys []AccessKey  `tfsdk:"access_keys"`
}

type AccessKeyResource struct {
	Id                types.String `tfsdk:"id"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	AccessKeyType     types.String `tfsdk:"access_key_type"`
	AccountId         types.String `tfsdk:"account_id"`
	Description       types.String `tfsdk:"description"`
	Duration          types.String `tfsdk:"duration"`
	ParentAccessKeyId types.String `tfsdk:"parent_access_key_id"`
	Passcode          types.String `tfsdk:"passcode"`
	AccessKey         types.Object `tfsdk:"access_key"`
	IsEnabled         types.Bool   `tfsdk:"is_enabled"`
}

type AccessKey struct {
	AccessKey           types.String `tfsdk:"access_key"`
	AccessKeyType       types.String `tfsdk:"access_key_type"`
	AccountId           types.String `tfsdk:"account_id"`
	CreatedAt           types.String `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	Description         types.String `tfsdk:"description"`
	ExpirationTimestamp types.String `tfsdk:"expiration_timestamp"`
	Id                  types.String `tfsdk:"id"`
	IsEnabled           types.Bool   `tfsdk:"is_enabled"`
	ModifiedAt          types.String `tfsdk:"modified_at"`
	ModifiedBy          types.String `tfsdk:"modified_by"`
	ParentAccessKeyId   types.String `tfsdk:"parent_access_key_id"`
	SecretKey           types.String `tfsdk:"secret_key"`
}

func (m AccessKey) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_key":           types.StringType,
		"access_key_type":      types.StringType,
		"account_id":           types.StringType,
		"created_at":           types.StringType,
		"created_by":           types.StringType,
		"description":          types.StringType,
		"expiration_timestamp": types.StringType,
		"id":                   types.StringType,
		"is_enabled":           types.BoolType,
		"modified_at":          types.StringType,
		"modified_by":          types.StringType,
		"parent_access_key_id": types.StringType,
		"secret_key":           types.StringType,
	}
}

type GroupDataSource struct {
	Size   types.Int32  `tfsdk:"size"`
	Page   types.Int32  `tfsdk:"page"`
	Sort   types.String `tfsdk:"sort"`
	Name   types.String `tfsdk:"name"`
	Groups []Group      `tfsdk:"groups"`
}

type GroupDataSourceDetail struct {
	Id    types.String `tfsdk:"id"`
	Group types.Object `tfsdk:"group"`
}

type GroupResource struct {
	Id          types.String   `tfsdk:"id"`
	Description types.String   `tfsdk:"description"`
	Name        types.String   `tfsdk:"name"`
	Tags        types.Map      `tfsdk:"tags"`
	PolicyIds   []types.String `tfsdk:"policy_ids"`
	UserIds     []types.String `tfsdk:"user_ids"`
	Group       types.Object   `tfsdk:"group"`
}

type GroupMembersDataResource struct {
	GroupId      types.String `tfsdk:"group_id"`
	Size         types.Int32  `tfsdk:"size"`
	Page         types.Int32  `tfsdk:"page"`
	Sort         types.String `tfsdk:"sort"`
	UserName     types.String `tfsdk:"user_name"`
	UserEmail    types.String `tfsdk:"user_email"`
	CreatorName  types.String `tfsdk:"creator_name"`
	CreatorEmail types.String `tfsdk:"creator_email"`
	GroupMembers []Member     `tfsdk:"group_members"`
}

type GroupMemberResource struct {
	GroupId     types.String `tfsdk:"group_id"`
	UserId      types.String `tfsdk:"user_id"`
	GroupMember types.Object `tfsdk:"group_member"`
}

type GroupPolicyBindingsDataResource struct {
	GroupId             types.String `tfsdk:"group_id"`
	Size                types.Int32  `tfsdk:"size"`
	Page                types.Int32  `tfsdk:"page"`
	Sort                types.String `tfsdk:"sort"`
	PolicyId            types.String `tfsdk:"policy_id"`
	PolicyName          types.String `tfsdk:"policy_name"`
	PolicyType          types.String `tfsdk:"policy_type"`
	GroupPolicyBindings []Policy     `tfsdk:"group_policy_bindings"`
}

type GroupPolicyBindingsResource struct {
	GroupId             types.String   `tfsdk:"group_id"`
	PolicyIds           []types.String `tfsdk:"policy_ids"`
	GroupPolicyBindings types.List     `tfsdk:"group_policy_bindings"`
}

type Group struct {
	CreatedAt     types.String `tfsdk:"created_at"`
	CreatedBy     types.String `tfsdk:"created_by"`
	CreatorEmail  types.String `tfsdk:"creator_email"`
	CreatorName   types.String `tfsdk:"creator_name"`
	Description   types.String `tfsdk:"description"`
	DomainName    types.String `tfsdk:"domain_name"`
	Id            types.String `tfsdk:"id"`
	Members       []Member     `tfsdk:"members"`
	Policies      []Policy     `tfsdk:"policies"`
	ModifiedAt    types.String `tfsdk:"modified_at"`
	ModifiedBy    types.String `tfsdk:"modified_by"`
	ModifierEmail types.String `tfsdk:"modifier_email"`
	ModifierName  types.String `tfsdk:"modifier_name"`
	Name          types.String `tfsdk:"name"`
	ResourceType  types.String `tfsdk:"resource_type"`
	ServiceName   types.String `tfsdk:"service_name"`
	ServiceType   types.String `tfsdk:"service_type"`
	Srn           types.String `tfsdk:"srn"`
	GroupType     types.String `tfsdk:"type"`
	state         attr.ValueState
}

func (v Group) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_at":    types.StringType,
		"created_by":    types.StringType,
		"creator_email": types.StringType,
		"creator_name":  types.StringType,
		"description":   types.StringType,
		"domain_name":   types.StringType,
		"id":            types.StringType,
		"members": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"created_at":            types.StringType,
					"created_by":            types.StringType,
					"creator_created_at":    types.StringType,
					"creator_email":         types.StringType,
					"creator_last_login_at": types.StringType,
					"creator_name":          types.StringType,
					"group_names": types.ListType{
						ElemType: types.StringType,
					},
					"user_created_at":    types.StringType,
					"user_email":         types.StringType,
					"user_id":            types.StringType,
					"user_last_login_at": types.StringType,
					"user_name":          types.StringType,
				},
			},
		},
		"policies": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: getPolicyAttribute(),
			},
		},
		"modified_at":    types.StringType,
		"modified_by":    types.StringType,
		"modifier_email": types.StringType,
		"modifier_name":  types.StringType,
		"name":           types.StringType,
		"resource_type":  types.StringType,
		"service_name":   types.StringType,
		"service_type":   types.StringType,
		"srn":            types.StringType,
		"type":           types.StringType,
	}
}

func getPolicyAttribute() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":         types.StringType,
		"created_at":         types.StringType,
		"created_by":         types.StringType,
		"creator_email":      types.StringType,
		"creator_name":       types.StringType,
		"default_version_id": types.StringType,
		"description":        types.StringType,
		"domain_name":        types.StringType,
		"id":                 types.StringType,
		"modified_at":        types.StringType,
		"modified_by":        types.StringType,
		"modifier_email":     types.StringType,
		"modifier_name":      types.StringType,
		"policy_category":    types.StringType,
		"policy_name":        types.StringType,
		"policy_type":        types.StringType,
		"policy_versions": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"created_at":  types.StringType,
					"created_by":  types.StringType,
					"id":          types.StringType,
					"modified_at": types.StringType,
					"modified_by": types.StringType,
					"policy_document": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"statement": types.ListType{
								ElemType: types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"action": types.ListType{
											ElemType: types.StringType,
										},
										"effect": types.StringType,
										"not_action": types.ListType{
											ElemType: types.StringType,
										},
										"condition": types.MapType{
											ElemType: types.MapType{
												ElemType: types.ListType{
													ElemType: types.StringType,
												},
											},
										},
										"principal": types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"principal_string": types.StringType,
												"principal_map": types.MapType{
													ElemType: types.ListType{
														ElemType: types.StringType,
													},
												},
											},
										},
										"resource": types.ListType{
											ElemType: types.StringType,
										},
										"sid": types.StringType,
									},
								},
							},
							"version": types.StringType,
						},
					},
					"policy_id":           types.StringType,
					"policy_version_name": types.StringType,
				},
			},
		},
		"resource_type": types.StringType,
		"service_name":  types.StringType,
		"service_type":  types.StringType,
		"srn":           types.StringType,
		"state":         types.StringType,
	}
}

type Member struct {
	CreatedAt          types.String   `tfsdk:"created_at"`
	CreatedBy          types.String   `tfsdk:"created_by"`
	CreatorCreatedAt   types.String   `tfsdk:"creator_created_at"`
	CreatorEmail       types.String   `tfsdk:"creator_email"`
	CreatorLastLoginAt types.String   `tfsdk:"creator_last_login_at"`
	CreatorName        types.String   `tfsdk:"creator_name"`
	GroupNames         []types.String `tfsdk:"group_names"`
	UserCreatedAt      types.String   `tfsdk:"user_created_at"`
	UserEmail          types.String   `tfsdk:"user_email"`
	UserId             types.String   `tfsdk:"user_id"`
	UserLastLoginAt    types.String   `tfsdk:"user_last_login_at"`
	UserName           types.String   `tfsdk:"user_name"`
}

func (v Member) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_at":            types.StringType,
		"created_by":            types.StringType,
		"creator_created_at":    types.StringType,
		"creator_email":         types.StringType,
		"creator_last_login_at": types.StringType,
		"creator_name":          types.StringType,
		"group_names": types.ListType{
			ElemType: types.StringType,
		},
		"user_created_at":    types.StringType,
		"user_email":         types.StringType,
		"user_id":            types.StringType,
		"user_last_login_at": types.StringType,
		"user_name":          types.StringType,
	}
}

type PolicyDatasource struct {
	Size       types.Int32  `tfsdk:"size"`
	Page       types.Int32  `tfsdk:"page"`
	Sort       types.String `tfsdk:"sort"`
	Id         types.String `tfsdk:"id"`
	PolicyName types.String `tfsdk:"policy_name"`
	PolicyType types.String `tfsdk:"policy_type"`
	Policies   []Policy     `tfsdk:"policies"`
}

type PolicyDatasourceDetail struct {
	Id     types.String `tfsdk:"id"`
	Policy types.Object `tfsdk:"policy"`
}

type PolicyResource struct {
	Id            types.String          `tfsdk:"id"`
	PolicyName    types.String          `tfsdk:"policy_name"`
	Description   types.String          `tfsdk:"description"`
	Tags          types.Map             `tfsdk:"tags"`
	PolicyVersion PolicyVersionResource `tfsdk:"policy_version"`
	Policy        types.Object          `tfsdk:"policy"`
}

type PolicyVersionResource struct {
	PolicyDocument IdentityPolicyDocument `tfsdk:"policy_document"`
}

type IdentityPolicyDocument struct {
	Version   types.String        `tfsdk:"version"`
	Statement []IdentityStatement `tfsdk:"statement"`
}

type IdentityStatement struct {
	Sid       types.String   `tfsdk:"sid"`
	Effect    types.String   `tfsdk:"effect"`
	Resource  []types.String `tfsdk:"resource"`
	Action    []types.String `tfsdk:"action"`
	NotAction []types.String `tfsdk:"not_action"`
	Condition types.Map      `tfsdk:"condition"`
}

type Policy struct {
	AccountId        types.String    `tfsdk:"account_id"`
	CreatedAt        types.String    `tfsdk:"created_at"`
	CreatedBy        types.String    `tfsdk:"created_by"`
	CreatorEmail     types.String    `tfsdk:"creator_email"`
	CreatorName      types.String    `tfsdk:"creator_name"`
	DefaultVersionId types.String    `tfsdk:"default_version_id"`
	Description      types.String    `tfsdk:"description"`
	DomainName       types.String    `tfsdk:"domain_name"`
	Id               types.String    `tfsdk:"id"`
	ModifiedAt       types.String    `tfsdk:"modified_at"`
	ModifiedBy       types.String    `tfsdk:"modified_by"`
	ModifierEmail    types.String    `tfsdk:"modifier_email"`
	ModifierName     types.String    `tfsdk:"modifier_name"`
	PolicyCategory   types.String    `tfsdk:"policy_category"`
	PolicyName       types.String    `tfsdk:"policy_name"`
	PolicyType       types.String    `tfsdk:"policy_type"`
	PolicyVersions   []PolicyVersion `tfsdk:"policy_versions"`
	ResourceType     types.String    `tfsdk:"resource_type"`
	ServiceName      types.String    `tfsdk:"service_name"`
	ServiceType      types.String    `tfsdk:"service_type"`
	Srn              types.String    `tfsdk:"srn"`
	State            types.String    `tfsdk:"state"`
}

func (v Policy) Attributes() map[string]attr.Type {
	return getPolicyAttribute()
}

type PolicyVersion struct {
	CreatedAt         types.String   `tfsdk:"created_at"`
	CreatedBy         types.String   `tfsdk:"created_by"`
	Id                types.String   `tfsdk:"id"`
	ModifiedAt        types.String   `tfsdk:"modified_at"`
	ModifiedBy        types.String   `tfsdk:"modified_by"`
	PolicyDocument    PolicyDocument `tfsdk:"policy_document"`
	PolicyId          types.String   `tfsdk:"policy_id"`
	PolicyVersionName types.String   `tfsdk:"policy_version_name"`
}

type PolicyDocument struct {
	Version   types.String `tfsdk:"version"`
	Statement []Statement  `tfsdk:"statement"`
}

type Statement struct {
	Sid       types.String   `tfsdk:"sid"`
	Effect    types.String   `tfsdk:"effect"`
	Resource  []types.String `tfsdk:"resource"`
	Action    []types.String `tfsdk:"action"`
	NotAction []types.String `tfsdk:"not_action"`
	Principal Principal      `tfsdk:"principal"`
	Condition types.Map      `tfsdk:"condition"`
}

type Principal struct {
	PrincipalString types.String `tfsdk:"principal_string"`
	PrincipalMap    types.Map    `tfsdk:"principal_map"`
}

type RoleDataSource struct {
	Size      types.Int32  `tfsdk:"size"`
	Page      types.Int32  `tfsdk:"page"`
	Sort      types.String `tfsdk:"sort"`
	Name      types.String `tfsdk:"name"`
	RoleType  types.String `tfsdk:"role_type"`
	AccountId types.String `tfsdk:"account_id"`
	Roles     []Role       `tfsdk:"roles"`
}

type RoleDataSourceDetail struct {
	Id   types.String `tfsdk:"id"`
	Role types.Object `tfsdk:"role"`
}

type RoleResource struct {
	Id                       types.String                     `tfsdk:"id"`
	AccountId                types.String                     `tfsdk:"account_id"`
	Name                     types.String                     `tfsdk:"name"`
	Description              types.String                     `tfsdk:"description"`
	MaxSessionDuration       types.Int32                      `tfsdk:"max_session_duration"`
	AssumeRolePolicyDocument *PolicyDocument                  `tfsdk:"assume_role_policy_document"`
	PolicyIds                []types.String                   `tfsdk:"policy_ids"`
	Principals               []CreateRoleTrustPolicyPrincipal `tfsdk:"principals"`
	Tags                     types.Map                        `tfsdk:"tags"`
	Role                     types.Object                     `tfsdk:"role"`
}

type CreateRoleTrustPolicyPrincipal struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

type Role struct {
	AccountId                types.String   `tfsdk:"account_id"`
	AssumeRolePolicyDocument PolicyDocument `tfsdk:"assume_role_policy_document"`
	CreatedAt                types.String   `tfsdk:"created_at"`
	CreatedBy                types.String   `tfsdk:"created_by"`
	CreatorEmail             types.String   `tfsdk:"creator_email"`
	CreatorName              types.String   `tfsdk:"creator_name"`
	Description              types.String   `tfsdk:"description"`
	Id                       types.String   `tfsdk:"id"`
	MaxSessionDuration       types.Int32    `tfsdk:"max_session_duration"`
	ModifiedAt               types.String   `tfsdk:"modified_at"`
	ModifiedBy               types.String   `tfsdk:"modified_by"`
	ModifierEmail            types.String   `tfsdk:"modifier_email"`
	ModifierName             types.String   `tfsdk:"modifier_name"`
	Name                     types.String   `tfsdk:"name"`
	Policies                 []Policy       `tfsdk:"policies"`
	Type                     types.String   `tfsdk:"type"`
}

func (v Role) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id": types.StringType,
		"assume_role_policy_document": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"statement": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"action": types.ListType{
								ElemType: types.StringType,
							},
							"effect": types.StringType,
							"not_action": types.ListType{
								ElemType: types.StringType,
							},
							"condition": types.MapType{
								ElemType: types.MapType{
									ElemType: types.ListType{
										ElemType: types.StringType,
									},
								},
							},
							"resource": types.ListType{
								ElemType: types.StringType,
							},
							"principal": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"principal_string": types.StringType,
									"principal_map": types.MapType{
										ElemType: types.ListType{
											ElemType: types.StringType,
										},
									},
								},
							},
							"sid": types.StringType,
						},
					},
				},
				"version": types.StringType,
			},
		},
		"created_at":           types.StringType,
		"created_by":           types.StringType,
		"creator_email":        types.StringType,
		"creator_name":         types.StringType,
		"description":          types.StringType,
		"id":                   types.StringType,
		"max_session_duration": types.Int32Type,
		"modified_at":          types.StringType,
		"modified_by":          types.StringType,
		"modifier_email":       types.StringType,
		"modifier_name":        types.StringType,
		"name":                 types.StringType,
		"policies": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: getPolicyAttribute(),
			},
		},
		"type": types.StringType,
	}
}

type RolePolicyBindingsDataSource struct {
	RoleId             types.String `tfsdk:"role_id"`
	Size               types.Int32  `tfsdk:"size"`
	Page               types.Int32  `tfsdk:"page"`
	Sort               types.String `tfsdk:"sort"`
	PolicyId           types.String `tfsdk:"policy_id"`
	PolicyName         types.String `tfsdk:"policy_name"`
	PolicyType         types.String `tfsdk:"policy_type"`
	RolePolicyBindings []Policy     `tfsdk:"role_policy_bindings"`
}

type RolePolicyBindingsResource struct {
	RoleId             types.String   `tfsdk:"role_id"`
	PolicyIds          []types.String `tfsdk:"policy_ids"`
	RolePolicyBindings types.List     `tfsdk:"role_policy_bindings"`
}

type UserDataSource struct {
	Size      types.Int32  `tfsdk:"size"`
	Page      types.Int32  `tfsdk:"page"`
	Sort      types.String `tfsdk:"sort"`
	Email     types.String `tfsdk:"email"`
	UserName  types.String `tfsdk:"user_name"`
	Type      types.String `tfsdk:"type"`
	AccountId types.String `tfsdk:"account_id"`
	Users     []User       `tfsdk:"users"`
}

type UserDataSourceDetail struct {
	AccountId types.String `tfsdk:"account_id"`
	UserId    types.String `tfsdk:"user_id"`
	User      types.Object `tfsdk:"user"`
}

type UserResource struct {
	AccountId          types.String   `tfsdk:"account_id"`
	UserId             types.String   `tfsdk:"user_id"`
	Description        types.String   `tfsdk:"description"`
	GroupIds           []types.String `tfsdk:"group_ids"`
	Password           types.String   `tfsdk:"password"`
	PolicyIds          []types.String `tfsdk:"policy_ids"`
	Tags               types.Map      `tfsdk:"tags"`
	TemporaryPassword  types.Bool     `tfsdk:"temporary_password"`
	UserName           types.String   `tfsdk:"user_name"`
	PasswordReuseCount types.Int32    `tfsdk:"password_reuse_count"`
	User               types.Object   `tfsdk:"user"`
}

type User struct {
	AccountId            types.String `tfsdk:"account_id"`
	CompanyName          types.String `tfsdk:"company_name"`
	ConsoleUrl           types.String `tfsdk:"console_url"`
	CreatedAt            types.String `tfsdk:"created_at"`
	CreatedBy            types.String `tfsdk:"created_by"`
	Description          types.String `tfsdk:"description"`
	DstOffset            types.String `tfsdk:"dst_offset"`
	Email                types.String `tfsdk:"email"`
	EmailAuthenticated   types.Bool   `tfsdk:"email_authenticated"`
	FirstName            types.String `tfsdk:"first_name"`
	Id                   types.String `tfsdk:"id"`
	LastLoginAt          types.String `tfsdk:"last_login_at"`
	LastName             types.String `tfsdk:"last_name"`
	LastPasswordUpdateAt types.String `tfsdk:"last_password_update_at"`
	ModifiedAt           types.String `tfsdk:"modified_at"`
	ModifiedBy           types.String `tfsdk:"modified_by"`
	Name                 types.String `tfsdk:"name"`
	Password             types.String `tfsdk:"password"`
	PasswordReuseCount   types.Int32  `tfsdk:"password_reuse_count"`
	PhoneAuthenticated   types.Bool   `tfsdk:"phone_authenticated"`
	Policies             []Policy     `tfsdk:"policies"`
	Timezone             types.String `tfsdk:"timezone"`
	Type                 types.String `tfsdk:"type"`
	TzId                 types.String `tfsdk:"tz_id"`
	UserName             types.String `tfsdk:"user_name"`
	UtcOffset            types.String `tfsdk:"utc_offset"`
}

func (v User) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":              types.StringType,
		"company_name":            types.StringType,
		"console_url":             types.StringType,
		"created_at":              types.StringType,
		"created_by":              types.StringType,
		"description":             types.StringType,
		"dst_offset":              types.StringType,
		"email":                   types.StringType,
		"email_authenticated":     types.BoolType,
		"first_name":              types.StringType,
		"id":                      types.StringType,
		"last_login_at":           types.StringType,
		"last_name":               types.StringType,
		"last_password_update_at": types.StringType,
		"modified_at":             types.StringType,
		"modified_by":             types.StringType,
		"name":                    types.StringType,
		"password":                types.StringType,
		"password_reuse_count":    types.Int32Type,
		"phone_authenticated":     types.BoolType,
		"policies": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: getPolicyAttribute(),
			},
		},
		"timezone":   types.StringType,
		"type":       types.StringType,
		"tz_id":      types.StringType,
		"user_name":  types.StringType,
		"utc_offset": types.StringType,
	}
}

type UserPolicyBindingsDataSource struct {
	UserId             types.String `tfsdk:"user_id"`
	Size               types.Int32  `tfsdk:"size"`
	Page               types.Int32  `tfsdk:"page"`
	Sort               types.String `tfsdk:"sort"`
	PolicyId           types.String `tfsdk:"policy_id"`
	PolicyName         types.String `tfsdk:"policy_name"`
	PolicyType         types.String `tfsdk:"policy_type"`
	UserPolicyBindings []Policy     `tfsdk:"user_policy_bindings"`
}

type UserPolicyBindingsResource struct {
	UserId             types.String   `tfsdk:"user_id"`
	PolicyIds          []types.String `tfsdk:"policy_ids"`
	UserPolicyBindings types.List     `tfsdk:"user_policy_bindings"`
}
