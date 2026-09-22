package iamidentitycenter

import (
	"context"
	"fmt"
	"net/http"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/iam-identity-center/1.6"
)

const ServiceType = "scp-iam-identity-center"

// Client is the IAM Identity Center API client
type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *sdk.APIClient
}

// NewClient creates a new IAM Identity Center client
func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: sdk.NewAPIClient(config),
	}
}

// identityStoreTypeFromString maps a documented identity store type name to its
// SDK enum value, returning an error for any unknown value instead of silently
// defaulting to a fixed type.
func identityStoreTypeFromString(value string) (sdk.InstanceStoreTypeEnum, error) {
	switch value {
	case "IDENTITY_CENTER_DIRECTORY":
		return sdk.INSTANCESTORETYPEENUM_IDENTITY_CENTER_DIRECTORY, nil
	case "ACTIVE_DIRECTORY":
		return sdk.INSTANCESTORETYPEENUM_ACTIVE_DIRECTORY, nil
	case "OPENLDAP":
		return sdk.INSTANCESTORETYPEENUM_OPENLDAP, nil
	case "DIRECTORY_389DS":
		return sdk.INSTANCESTORETYPEENUM_DIRECTORY_389_DS, nil
	default:
		return "", fmt.Errorf("unsupported identity_store_type %q: must be one of IDENTITY_CENTER_DIRECTORY, ACTIVE_DIRECTORY, OPENLDAP, DIRECTORY_389DS", value)
	}
}

// CreateInstance creates a new IAM Identity Center instance
func (client *Client) CreateInstance(ctx context.Context, request InstanceCreateRequest) (*sdk.InstanceShowResponseV1Dot5, error) {
	req := client.sdkClient.IamIdentityCenterV1InstancesAPIsAPI.CreateInstance(ctx)

	// Build LDAP config if provided
	var ldapConfig *sdk.LdapConfigCreateRequest
	if request.HasIdentityStoreConfig() {
		ldapConfig = &sdk.LdapConfigCreateRequest{
			BindCredential:        request.BindCredential,
			BindDn:                request.BindDn,
			ConnectionUrl:         request.ConnectionUrl,
			RdnLdapAttribute:      request.RdnLdapAttribute,
			UserObjectClasses:     request.UserObjectClasses,
			UsernameLdapAttribute: request.UsernameLdapAttribute,
			UsersDn:               request.UsersDn,
		}
	}

	identityStoreType, err := identityStoreTypeFromString(request.IdentityStoreType)
	if err != nil {
		return nil, err
	}

	createReq := sdk.InstanceCreateRequestV1Dot5{
		Name:              request.Name,
		IdentityStoreType: &identityStoreType,
	}

	if request.Description != "" {
		createReq.Description = *sdk.NewNullableString(&request.Description)
	}

	if request.SelfManagedPassword != nil {
		createReq.SelfManagedPassword = *sdk.NewNullableBool(request.SelfManagedPassword)
	}

	if ldapConfig != nil {
		createReq.IdentityStoreConfig = *sdk.NewNullableLdapConfigCreateRequest(ldapConfig)
	}

	req = req.InstanceCreateRequestV1Dot5(createReq)

	resp, _, err := req.Execute()
	if err != nil {
		if apiErr, ok := err.(*scpsdk.GenericOpenAPIError); ok {
			body := string(apiErr.Body())
			if body != "" {
				return resp, fmt.Errorf("%s: %s", apiErr.Error(), body)
			}
		}
	}
	return resp, err
}

// GetInstance retrieves an IAM Identity Center instance by ID
func (client *Client) GetInstance(ctx context.Context, instanceId string) (*sdk.InstanceShowResponseAddResourceV1Dot6, *http.Response, error) {
	req := client.sdkClient.IamIdentityCenterV1InstancesAPIsAPI.ShowInstance(ctx, instanceId)

	resp, httpResp, err := req.Execute()
	return resp, httpResp, err
}

// ListInstances lists all IAM Identity Center instances
func (client *Client) ListInstances(ctx context.Context, size int32, page int32, sort string) (*sdk.InstancePageResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1InstancesAPIsAPI.ListInstances(ctx)

	if size > 0 {
		req = req.Size(size)
	}
	if page > 0 {
		req = req.Page(page)
	}
	if sort != "" {
		req = req.Sort(sort)
	}

	resp, _, err := req.Execute()
	return resp, err
}

// ListGroups lists all IAM Identity Center groups
func (client *Client) ListGroups(ctx context.Context, instanceId string, name string, size int32, page int32, sort string, excludedUserUuid string, excludedAccountId string) (*sdk.GroupPageResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1GroupsAPIsAPI.ListGroups(ctx)

	if instanceId != "" {
		req = req.InstanceId(instanceId)
	}
	if name != "" {
		req = req.Name(name)
	}
	if size > 0 {
		req = req.Size(size)
	}
	if page > 0 {
		req = req.Page(page)
	}
	if sort != "" {
		req = req.Sort(sort)
	}
	if excludedUserUuid != "" {
		req = req.ExcludedUserUuid(excludedUserUuid)
	}
	if excludedAccountId != "" {
		req = req.ExcludedAccountId(excludedAccountId)
	}

	resp, _, err := req.Execute()
	return resp, err
}

// UpdateInstance updates an IAM Identity Center instance
func (client *Client) UpdateInstance(ctx context.Context, instanceId string, request InstanceUpdateRequest) (*sdk.InstanceShowResponseV1Dot5, error) {
	req := client.sdkClient.IamIdentityCenterV1InstancesAPIsAPI.SetInstance(ctx, instanceId)

	updateReq := sdk.InstanceSetRequest1Dot5{}

	if request.Description != "" {
		updateReq.Description = *sdk.NewNullableString(&request.Description)
	}

	if request.Name != "" {
		updateReq.Name = *sdk.NewNullableString(&request.Name)
	}

	if request.IdentityStoreType != "" {
		storeType, err := identityStoreTypeFromString(request.IdentityStoreType)
		if err != nil {
			return nil, err
		}
		updateReq.IdentityStoreType = *sdk.NewNullableInstanceStoreTypeEnum(&storeType)
	}

	if request.SelfManagedPassword != nil {
		updateReq.SelfManagedPassword = *sdk.NewNullableBool(request.SelfManagedPassword)
	}

	if request.HasIdentityStoreConfig() {
		ldapConfig := sdk.LdapConfig{
			BindCredential:        request.BindCredential,
			BindDn:                request.BindDn,
			ConnectionUrl:         request.ConnectionUrl,
			RdnLdapAttribute:      request.RdnLdapAttribute,
			UserObjectClasses:     request.UserObjectClasses,
			UsernameLdapAttribute: request.UsernameLdapAttribute,
			UsersDn:               request.UsersDn,
		}
		updateReq.IdentityStoreConfig = *sdk.NewNullableLdapConfig(&ldapConfig)
	}

	req = req.InstanceSetRequest1Dot5(updateReq)

	resp, _, err := req.Execute()
	return resp, err
}

// DeleteInstance deletes an IAM Identity Center instance
func (client *Client) DeleteInstance(ctx context.Context, instanceId string) (*http.Response, error) {
	req := client.sdkClient.IamIdentityCenterV1InstancesAPIsAPI.DeleteInstance(ctx, instanceId)

	return req.Execute()
}

// InstanceCreateRequest represents the request for creating an instance
type InstanceCreateRequest struct {
	Name                   string
	Description            string
	IdentityStoreType      string
	SelfManagedPassword    *bool
	identityStoreConfigSet bool
	BindCredential         *string
	BindDn                 *string
	ConnectionUrl          *string
	RdnLdapAttribute       *string
	UserObjectClasses      *string
	UsernameLdapAttribute  *string
	UsersDn                *string
}

// HasIdentityStoreConfig returns whether the identity store config has been set
func (r *InstanceCreateRequest) HasIdentityStoreConfig() bool {
	return r.identityStoreConfigSet
}

// SetIdentityStoreConfig sets the identity store config
func (r *InstanceCreateRequest) SetIdentityStoreConfig(config LdapConfigOptions) {
	r.identityStoreConfigSet = true
	r.BindCredential = config.BindCredential
	r.BindDn = config.BindDn
	r.ConnectionUrl = config.ConnectionUrl
	r.RdnLdapAttribute = config.RdnLdapAttribute
	r.UserObjectClasses = config.UserObjectClasses
	r.UsernameLdapAttribute = config.UsernameLdapAttribute
	r.UsersDn = config.UsersDn
}

// LdapConfigOptions represents the LDAP configuration options
type LdapConfigOptions struct {
	BindCredential        *string
	BindDn                *string
	ConnectionUrl         *string
	RdnLdapAttribute      *string
	UserObjectClasses     *string
	UsernameLdapAttribute *string
	UsersDn               *string
}

// InstanceUpdateRequest represents the request for updating an instance
type InstanceUpdateRequest struct {
	Description            string
	Name                   string
	IdentityStoreType      string
	SelfManagedPassword    *bool
	identityStoreConfigSet bool
	BindCredential         *string
	BindDn                 *string
	ConnectionUrl          *string
	RdnLdapAttribute       *string
	UserObjectClasses      *string
	UsernameLdapAttribute  *string
	UsersDn                *string
}

// HasIdentityStoreConfig returns whether the identity store config has been set
func (r *InstanceUpdateRequest) HasIdentityStoreConfig() bool {
	return r.identityStoreConfigSet
}

// SetIdentityStoreConfig sets the identity store config
func (r *InstanceUpdateRequest) SetIdentityStoreConfig(config LdapConfigOptions) {
	r.identityStoreConfigSet = true
	r.BindCredential = config.BindCredential
	r.BindDn = config.BindDn
	r.ConnectionUrl = config.ConnectionUrl
	r.RdnLdapAttribute = config.RdnLdapAttribute
	r.UserObjectClasses = config.UserObjectClasses
	r.UsernameLdapAttribute = config.UsernameLdapAttribute
	r.UsersDn = config.UsersDn
}

// CreateGroup creates a new IAM Identity Center group
func (client *Client) CreateGroup(ctx context.Context, request GroupCreateRequest) (*sdk.GroupShowResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1GroupsAPIsAPI.CreateGroup(ctx)

	createReq := sdk.GroupCreateRequest{
		InstanceId: request.InstanceId,
		Name:       request.Name,
	}

	if request.Description != "" {
		createReq.Description = *sdk.NewNullableString(&request.Description)
	}

	if len(request.UserUuids) > 0 {
		createReq.UserUuids = request.UserUuids
	}

	req = req.GroupCreateRequest(createReq)

	resp, _, err := req.Execute()
	return resp, err
}

// GetGroup retrieves an IAM Identity Center group by ID
func (client *Client) GetGroup(ctx context.Context, groupId string, instanceId string) (*sdk.GroupShowResponse, *http.Response, error) {
	req := client.sdkClient.IamIdentityCenterV1GroupsAPIsAPI.ShowGroup(ctx, groupId)

	if instanceId != "" {
		req = req.InstanceId(instanceId)
	}

	result, r, err := req.Execute()
	return result, r, err
}

// DeleteGroup deletes an IAM Identity Center group
func (client *Client) DeleteGroup(ctx context.Context, groupId string, instanceId string) (*sdk.GroupRemoveResponse, *http.Response, error) {
	req := client.sdkClient.IamIdentityCenterV1GroupsAPIsAPI.DeleteGroup(ctx, groupId)

	instanceBaseReq := sdk.NewInstanceBaseRequest(instanceId)
	req = req.InstanceBaseRequest(*instanceBaseReq)

	result, r, err := req.Execute()
	return result, r, err
}

func (client *Client) UpdateGroup(ctx context.Context, groupId string, request GroupUpdateRequest) (*sdk.GroupSetResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1GroupsAPIsAPI.SetGroup(ctx, groupId)

	updateReq := sdk.GroupSetRequest{
		InstanceId: request.InstanceId,
	}

	if request.Name != "" {
		updateReq.Name = *sdk.NewNullableString(&request.Name)
	}

	if request.Description != "" {
		updateReq.Description = *sdk.NewNullableString(&request.Description)
	}

	req = req.GroupSetRequest(updateReq)

	resp, _, err := req.Execute()
	return resp, err
}

type GroupUpdateRequest struct {
	InstanceId  string
	Name        string
	Description string
}

type GroupCreateRequest struct {
	InstanceId  string
	Name        string
	Description string
	UserUuids   []string
}

// AddUsersToGroup adds users to an IAM Identity Center group
func (client *Client) AddUsersToGroup(ctx context.Context, groupId string, request sdk.GroupUsersRequest) (*sdk.GroupShowResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1GroupsAPIsAPI.CreateBulkGroupUsers(ctx, groupId)
	req = req.GroupUsersRequest(request)
	resp, _, err := req.Execute()
	return resp, err
}

// RemoveUsersFromGroup removes users from an IAM Identity Center group
func (client *Client) RemoveUsersFromGroup(ctx context.Context, groupId string, request sdk.GroupUsersRequest) error {
	req := client.sdkClient.IamIdentityCenterV1GroupsAPIsAPI.DeleteBulkGroupUsers(ctx, groupId)
	req = req.GroupUsersRequest(request)
	_, err := req.Execute()
	return err
}

// ListGroupUsers lists users of an IAM Identity Center group
func (client *Client) ListGroupUsers(ctx context.Context, groupId string, instanceId string, userId string, size int32, page int32) (*sdk.GroupUsersPageResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1GroupsAPIsAPI.ListGroupUsers(ctx, groupId)

	if instanceId != "" {
		req = req.InstanceId(instanceId)
	}
	if userId != "" {
		req = req.UserId(userId)
	}
	if size > 0 {
		req = req.Size(size)
	}
	if page > 0 {
		req = req.Page(page)
	}

	resp, _, err := req.Execute()
	return resp, err
}

// CreateUser creates a new IAM Identity Center user
func (client *Client) CreateUser(ctx context.Context, request sdk.UserCreateRequest) (*sdk.UserResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1UsersAPIsAPI.CreateUser(ctx)

	req = req.UserCreateRequest(request)

	resp, _, err := req.Execute()
	return resp, err
}

// GetUser retrieves an IAM Identity Center user by ID
func (client *Client) GetUser(ctx context.Context, instanceId string, userId string) (*sdk.UserResponse, *http.Response, error) {
	req := client.sdkClient.IamIdentityCenterV1UsersAPIsAPI.ShowUser(ctx, userId)

	if instanceId != "" {
		req = req.InstanceId(instanceId)
	}

	resp, httpResp, err := req.Execute()
	return resp, httpResp, err
}

// ListUsers lists all IAM Identity Center users
func (client *Client) ListUsers(ctx context.Context, instanceId string, userId string, size int32, page int32, sort string, excludedGroupId string, excludedAccountId string) (*sdk.UserPageResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1UsersAPIsAPI.ListUsers(ctx)

	if instanceId != "" {
		req = req.InstanceId(instanceId)
	}
	if userId != "" {
		req = req.UserId(userId)
	}
	if size > 0 {
		req = req.Size(size)
	}
	if page > 0 {
		req = req.Page(page)
	}
	if sort != "" {
		req = req.Sort(sort)
	}
	if excludedGroupId != "" {
		req = req.ExcludedGroupId(excludedGroupId)
	}
	if excludedAccountId != "" {
		req = req.ExcludedAccountId(excludedAccountId)
	}

	resp, _, err := req.Execute()
	return resp, err
}

// UpdateUser updates an IAM Identity Center user
func (client *Client) UpdateUser(ctx context.Context, userId string, request sdk.UserSetRequest) (*sdk.UserResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1UsersAPIsAPI.SetUser(ctx, userId)

	req = req.UserSetRequest(request)

	resp, _, err := req.Execute()
	return resp, err
}

// DeleteUser deletes an IAM Identity Center user
func (client *Client) DeleteUser(ctx context.Context, userId string, instanceId string) (*sdk.UserRemoveResponse, *http.Response, error) {
	req := client.sdkClient.IamIdentityCenterV1UsersAPIsAPI.DeleteUser(ctx, userId)
	req = req.InstanceBaseRequest(sdk.InstanceBaseRequest{
		InstanceId: instanceId,
	})

	return req.Execute()
}

// CreatePermissionSet creates a new IAM Identity Center permission set
func (client *Client) CreatePermissionSet(ctx context.Context, request PermissionSetCreateRequest) (*sdk.PermissionSetDetailsV1Dot2, error) {
	req := client.sdkClient.IamIdentityCenterV1PermissionSetsAPIsAPI.CreatePermissionSet(ctx)

	createReq := sdk.PermissionSetCreateRequestV1Dot2{
		InstanceId:      request.InstanceId,
		Name:            request.Name,
		SessionDuration: request.SessionDuration,
	}

	if request.Description != "" {
		createReq.Description = *sdk.NewNullableString(&request.Description)
	}

	if len(request.CustomPolicies) > 0 {
		createReq.CustomPolicies = request.CustomPolicies
	}

	if len(request.ManagedPolicies) > 0 {
		createReq.ManagedPolicies = request.ManagedPolicies
	}

	if len(request.InlinePolicies) > 0 {
		createReq.InlinePolicies = request.InlinePolicies
	}

	req = req.PermissionSetCreateRequestV1Dot2(createReq)

	resp, _, err := req.Execute()
	if err != nil {
		if apiErr, ok := err.(*scpsdk.GenericOpenAPIError); ok {
			body := string(apiErr.Body())
			if body != "" {
				return resp, fmt.Errorf("%s: %s", apiErr.Error(), body)
			}
		}
	}
	return resp, err
}

// GetPermissionSet retrieves an IAM Identity Center permission set by ID
func (client *Client) GetPermissionSet(ctx context.Context, permissionSetId string, instanceId string) (*sdk.ShowPermissionSetAddResourceV1Dot2, *http.Response, error) {
	req := client.sdkClient.IamIdentityCenterV1PermissionSetsAPIsAPI.ShowPermissionSet(ctx, permissionSetId)

	if instanceId != "" {
		req = req.InstanceId(instanceId)
	}

	result, r, err := req.Execute()
	return result, r, err
}

// DeletePermissionSet deletes an IAM Identity Center permission set
func (client *Client) DeletePermissionSet(ctx context.Context, permissionSetId string, instanceId string) error {
	req := client.sdkClient.IamIdentityCenterV1PermissionSetsAPIsAPI.DeletePermissionSet(ctx, permissionSetId)
	req = req.InstanceId(instanceId)

	_, err := req.Execute()
	return err
}

// UpdatePermissionSet updates an IAM Identity Center permission set
func (client *Client) UpdatePermissionSet(ctx context.Context, permissionSetId string, request PermissionSetUpdateRequest) (*sdk.PermissionSetDetailsV1Dot2, error) {
	req := client.sdkClient.IamIdentityCenterV1PermissionSetsAPIsAPI.SetPermissionSet(ctx, permissionSetId)

	updateReq := sdk.PermissionSetSetRequest{}

	if request.InstanceId != "" {
		updateReq.InstanceId = request.InstanceId
	}

	if request.Description != "" {
		updateReq.Description = *sdk.NewNullableString(&request.Description)
	}

	if request.SessionDuration != 0 {
		updateReq.SessionDuration = *sdk.NewNullableInt32(&request.SessionDuration)
	}

	req = req.PermissionSetSetRequest(updateReq)

	resp, _, err := req.Execute()
	return resp, err
}

// ListPermissionSets lists all IAM Identity Center permission sets
func (client *Client) ListPermissionSets(ctx context.Context, instanceId string, name string, size int32, page int32, sort string) (*sdk.PermissionSetDetailsPageResponseV1Dot2, error) {
	req := client.sdkClient.IamIdentityCenterV1PermissionSetsAPIsAPI.ListPermissionSets(ctx)

	if instanceId != "" {
		req = req.InstanceId(instanceId)
	}
	if name != "" {
		req = req.Name(name)
	}
	if size > 0 {
		req = req.Size(size)
	}
	if page > 0 {
		req = req.Page(page)
	}
	if sort != "" {
		req = req.Sort(sort)
	}

	resp, _, err := req.Execute()
	return resp, err
}

type PermissionSetUpdateRequest struct {
	InstanceId      string
	Name            string
	Description     string
	SessionDuration int32
}

type PermissionSetCreateRequest struct {
	InstanceId      string
	Name            string
	Description     string
	SessionDuration int32
	CustomPolicies  []string
	InlinePolicies  []map[string]interface{}
	ManagedPolicies []sdk.ManagedPolicy
}

type PoliciesSetRequest struct {
	InstanceId      string
	CustomPolicies  []string
	InlinePolicies  []map[string]interface{}
	ManagedPolicies []map[string]interface{}
}

type PoliciesDeleteRequest struct {
	InstanceId string
	PolicyIds  []string
}

// GetPolicies returns the policies of an IAM Identity Center permission set
func (client *Client) GetPolicies(ctx context.Context, permissionSetId string, instanceId string, category string, name string, size int32, page int32) (*sdk.PermissionSetPoliciesPageResponseV1Dot2, error) {
	req := client.sdkClient.IamIdentityCenterV1PermissionSetsAPIsAPI.ListPermissionSetPolicies(ctx, permissionSetId)

	req = req.InstanceId(instanceId)

	if size > 0 {
		req = req.Size(size)
	}
	if page > 0 {
		req = req.Page(page)
	}
	if category != "" {
		var cat sdk.PermissionSetPolicyTypeEnumV1Dot2
		switch category {
		case "MANAGED_POLICY":
			cat = sdk.PERMISSIONSETPOLICYTYPEENUMV1DOT2_MANAGED_POLICY
		case "CUSTOM_POLICY":
			cat = sdk.PERMISSIONSETPOLICYTYPEENUMV1DOT2_CUSTOM_POLICY
		case "INLINE_POLICY":
			cat = sdk.PERMISSIONSETPOLICYTYPEENUMV1DOT2_INLINE_POLICY
		}
		req = req.Category(cat)
	}
	if name != "" {
		req = req.Name(name)
	}

	resp, _, err := req.Execute()
	return resp, err
}

// SetPolicies sets the policies of an IAM Identity Center permission set
func (client *Client) SetPolicies(ctx context.Context, permissionSetId string, request PoliciesSetRequest) (*sdk.PermissionSetDetailsV1Dot2, error) {
	req := client.sdkClient.IamIdentityCenterV1PermissionSetsAPIsAPI.SetPermissionSetPolicies(ctx, permissionSetId)

	setReq := sdk.PermissionSetPoliciesSetRequestV1Dot2{
		InstanceId: request.InstanceId,
	}

	if len(request.CustomPolicies) > 0 {
		setReq.CustomPolicies = request.CustomPolicies
	}

	if len(request.InlinePolicies) > 0 {
		setReq.InlinePolicies = request.InlinePolicies
	}

	if len(request.ManagedPolicies) > 0 {
		setReq.ManagedPolicies = request.ManagedPolicies
	}

	req = req.PermissionSetPoliciesSetRequestV1Dot2(setReq)

	resp, _, err := req.Execute()
	if err != nil {
		if apiErr, ok := err.(*scpsdk.GenericOpenAPIError); ok {
			body := string(apiErr.Body())
			if body != "" {
				return resp, fmt.Errorf("%s: %s", apiErr.Error(), body)
			}
		}
	}
	return resp, err
}

// DeletePolicies deletes the policies of an IAM Identity Center permission set
func (client *Client) DeletePolicies(ctx context.Context, permissionSetId string, request PoliciesDeleteRequest) error {
	req := client.sdkClient.IamIdentityCenterV1PermissionSetsAPIsAPI.DeletePermissionSetPolicies(ctx, permissionSetId)

	deleteReq := sdk.PermissionSetPoliciesDeleteRequest{
		InstanceId: request.InstanceId,
		PolicyIds:  request.PolicyIds,
	}

	req = req.PermissionSetPoliciesDeleteRequest(deleteReq)

	_, err := req.Execute()
	return err
}

// CreateAccountAssignment creates IAM Identity Center account assignments
func (client *Client) CreateAccountAssignment(ctx context.Context, request AccountAssignmentCreateRequest) (*sdk.AccountAssignmentsCreateResponseV1Dot2, error) {
	req := client.sdkClient.IamIdentityCenterV1AccountAssignmentsAPIsAPI.CreateAccountAssignment(ctx)

	createReq := sdk.AccountAssignmentsCreateRequestV1Dot6{
		InstanceId: request.InstanceId,
	}

	if len(request.Accounts) > 0 {
		createReq.Accounts = request.Accounts
	}

	if len(request.PermissionSets) > 0 {
		createReq.PermissionSets = request.PermissionSets
	}

	if len(request.Principals) > 0 {
		createReq.Principals = request.Principals
	}

	req = req.AccountAssignmentsCreateRequestV1Dot6(createReq)

	resp, _, err := req.Execute()
	if err != nil {
		if apiErr, ok := err.(*scpsdk.GenericOpenAPIError); ok {
			body := string(apiErr.Body())
			if body != "" {
				return resp, fmt.Errorf("%s: %s", apiErr.Error(), body)
			}
		}
	}
	return resp, err
}

// GetAccountAssignment retrieves an IAM Identity Center account assignment by ID
func (client *Client) GetAccountAssignment(ctx context.Context, instanceId string, targetAccountId string) (*sdk.AccountAssignmentsPageResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1AccountAssignmentsAPIsAPI.ListAccountAssignments(ctx)

	req = req.InstanceId(instanceId)
	req = req.TargetAccountId(targetAccountId)

	resp, _, err := req.Execute()
	return resp, err
}

// ListAccountAssignments lists IAM Identity Center account assignments with pagination
func (client *Client) ListAccountAssignments(ctx context.Context, instanceId string, targetAccountId string, size int32, page int32, sort string, targetAccountName string, targetAccountEmail string, permissionSetId string, principalName string, roleSrn string) (*sdk.AccountAssignmentsPageResponse, error) {
	req := client.sdkClient.IamIdentityCenterV1AccountAssignmentsAPIsAPI.ListAccountAssignments(ctx)

	req = req.InstanceId(instanceId)
	if targetAccountId != "" {
		req = req.TargetAccountId(targetAccountId)
	}
	if size > 0 {
		req = req.Size(size)
	}
	if page > 0 {
		req = req.Page(page)
	}
	if sort != "" {
		req = req.Sort(sort)
	}
	if targetAccountName != "" {
		req = req.TargetAccountName(targetAccountName)
	}
	if targetAccountEmail != "" {
		req = req.TargetAccountEmail(targetAccountEmail)
	}
	if permissionSetId != "" {
		req = req.PermissionSetId(permissionSetId)
	}
	if principalName != "" {
		req = req.PrincipalName(principalName)
	}
	if roleSrn != "" {
		req = req.RoleSrn(roleSrn)
	}

	resp, _, err := req.Execute()
	return resp, err
}

// ListAllAccountAssignments lists all account assignments of a target account,
// following pagination until every result is collected.
func (client *Client) ListAllAccountAssignments(ctx context.Context, instanceId string, targetAccountId string) (*sdk.AccountAssignmentsPageResponse, error) {
	const pageSize int32 = 100

	var all []sdk.AccountAssignment
	var totalCount int32
	page := int32(0)

	for {
		req := client.sdkClient.IamIdentityCenterV1AccountAssignmentsAPIsAPI.ListAccountAssignments(ctx).
			InstanceId(instanceId).
			TargetAccountId(targetAccountId).
			Size(pageSize).
			Page(page)

		resp, _, err := req.Execute()
		if err != nil {
			return nil, err
		}
		if resp == nil {
			break
		}

		totalCount = resp.Count
		if resp.AccountAssignments != nil {
			all = append(all, resp.AccountAssignments...)
		}

		if totalCount > 0 && int32(len(all)) >= totalCount {
			break
		}
		if len(resp.AccountAssignments) == 0 {
			break
		}
		if page > 1000 {
			break
		}
		page++
	}

	return &sdk.AccountAssignmentsPageResponse{
		AccountAssignments: all,
		Count:              totalCount,
	}, nil
}

// DeleteAccountAssignment deletes an IAM Identity Center account assignment
func (client *Client) DeleteAccountAssignment(ctx context.Context, accountAssignmentId string, request AccountAssignmentDeleteRequest) error {
	req := client.sdkClient.IamIdentityCenterV1AccountAssignmentsAPIsAPI.DeleteAccountAssignment(ctx, accountAssignmentId)

	deleteReq := sdk.AccountAssignmentDeleteRequest{
		InstanceId:      request.InstanceId,
		TargetAccountId: request.TargetAccountId,
	}

	if request.OnlyUserRemoved {
		deleteReq.OnlyUserRemoved = *sdk.NewNullableBool(&request.OnlyUserRemoved)
	}

	req = req.AccountAssignmentDeleteRequest(deleteReq)

	_, err := req.Execute()
	return err
}

type AccountAssignmentCreateRequest struct {
	InstanceId     string
	Accounts       []sdk.AccountV1Dot6
	PermissionSets []sdk.PermissionSetInfoV1Dot6
	Principals     []sdk.Principal
}

type AccountAssignmentDeleteRequest struct {
	InstanceId      string
	TargetAccountId string
	OnlyUserRemoved bool
}
