package iam

import (
	"context"
	"fmt"
	"sort"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scpiam1d0 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/iam/1.0"
	scpiam1d1 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/iam/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Client struct {
	Config       *scpsdk.Configuration
	sdkClient1d0 *scpiam1d0.APIClient
	sdkClient1d1 *scpiam1d1.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:       config,
		sdkClient1d0: scpiam1d0.NewAPIClient(config),
		sdkClient1d1: scpiam1d1.NewAPIClient(config),
	}
}

func (client *Client) GetAccessKeyList(ctx context.Context, request AccessKeyDataSource) (*scpiam1d0.ListAccessKeyResponse, error) {
	req := client.sdkClient1d0.IamV1AccessKeysApiAPI.AccessKeyList(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.AccountId.IsNull() {
		req = req.AccountId(request.AccountId.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateAccessKey(ctx context.Context, request AccessKeyResource) (*scpiam1d0.AccessKeyResponse, error) {
	req := client.sdkClient1d0.IamV1AccessKeysApiAPI.AccessKeyCreate(ctx)

	req = req.AccessKeyCreateRequest(scpiam1d0.AccessKeyCreateRequest{
		AccessKeyType:     scpiam1d0.AccessKeyTypeCreateRequestEnum(request.AccessKeyType.ValueString()),
		AccountId:         request.AccountId.ValueStringPointer(),
		Description:       *scpiam1d0.NewNullableString(request.Description.ValueStringPointer()),
		Duration:          *scpiam1d0.NewNullableString(request.Duration.ValueStringPointer()),
		ParentAccessKeyId: *scpiam1d0.NewNullableString(request.ParentAccessKeyId.ValueStringPointer()),
		Passcode:          *scpiam1d0.NewNullableString(request.Passcode.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetAccessKey(ctx context.Context, accessKeyId string) (*scpiam1d0.AccessKeyResponse, error) {
	req := client.sdkClient1d0.IamV1AccessKeysApiAPI.AccessKeyShow(ctx, accessKeyId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateAccessKey(ctx context.Context, accessKeyId string, request AccessKeyResource) (*scpiam1d0.AccessKeyResponse, error) {
	req := client.sdkClient1d0.IamV1AccessKeysApiAPI.AccessKeySet(ctx, accessKeyId)

	req = req.AccessKeyUpdateRequest(scpiam1d0.AccessKeyUpdateRequest{
		IsEnabled: request.IsEnabled.ValueBoolPointer(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteAccessKey(ctx context.Context, accessKeyId string) error {
	req := client.sdkClient1d0.IamV1AccessKeysApiAPI.AccessKeyDelete(ctx, accessKeyId)

	_, err := req.Execute()
	return err
}

func (client *Client) GetEndpointList() (*scpiam1d0.ListEndpointsResponse, error) {
	ctx := context.Background()

	req := client.sdkClient1d0.IamV1EndpointsApiAPI.ListEndpoints(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

var regions []string

func (client *Client) GetRegionList() []string {
	if len(regions) == 0 {
		ctx := context.Background()

		req := client.sdkClient1d0.IamV1EndpointsApiAPI.ListEndpoints(ctx)

		resp, _, _ := req.Execute()

		regionMap := make(map[string]bool)
		var regions []string

		for _, endpoint := range resp.Endpoints {
			if !regionMap[endpoint.Region] {
				regionMap[endpoint.Region] = true
				regions = append(regions, endpoint.Region)
			}
		}

		sort.Slice(regions, func(i, j int) bool {
			return regions[i] < regions[j]
		})
	}

	return regions
}

func (client *Client) GetAccountId() (string, error) {
	ctx := context.Background()
	data, err := client.GetAccessKeyList(ctx, AccessKeyDataSource{})
	if err != nil {
		return "", err
	}

	if len(data.AccessKeys) == 0 {
		return "", fmt.Errorf("failed to find Account ID")
	}

	accessKey := data.AccessKeys[0]
	return accessKey.AccountId, nil
}

// / GROUP ///
func (client *Client) GetGroups(ctx context.Context, request GroupDataSource) (*scpiam1d0.GroupPageResponse, error) {
	req := client.sdkClient1d0.IamV1GroupsApiAPI.ListGroup(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetGroup(ctx context.Context, groupId string) (*scpiam1d0.GroupShowResponse, error) {
	req := client.sdkClient1d0.IamV1GroupsApiAPI.ShowGroup(ctx, groupId)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) CreateGroup(ctx context.Context, request GroupResource) (*scpiam1d0.GroupShowResponse, error) {
	req := client.sdkClient1d0.IamV1GroupsApiAPI.CreateGroup(ctx)

	//tag
	var TagsObject []map[string]string

	for k, v := range request.Tags.Elements() {
		tagObject := make(map[string]string)
		tagObject["key"] = k
		tagObject["value"] = v.(types.String).ValueString()

		TagsObject = append(TagsObject, tagObject)
	}

	//policy ids
	var policyIds []string
	for _, policyId := range request.PolicyIds {
		policyIds = append(policyIds, policyId.ValueString())
	}

	//user ids
	var userIds []string
	for _, userId := range request.UserIds {
		userIds = append(userIds, userId.ValueString())
	}

	req = req.GroupCreateRequest(scpiam1d0.GroupCreateRequest{
		Description: *scpiam1d0.NewNullableString(request.Description.ValueStringPointer()),
		Name:        request.Name.ValueString(),
		Tags:        TagsObject,
		PolicyIds:   policyIds,
		UserIds:     userIds,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateGroup(ctx context.Context, groupId string, request GroupResource) (*scpiam1d0.GroupShowResponse, error) {
	req := client.sdkClient1d0.IamV1GroupsApiAPI.SetGroup(ctx, groupId)

	req = req.GroupSetRequest(scpiam1d0.GroupSetRequest{
		Description: *scpiam1d0.NewNullableString(request.Description.ValueStringPointer()),
		Name:        request.Name.ValueString(),
	})

	resp, _, err := req.Execute()

	return resp, err
}

func (client *Client) DeleteGroup(ctx context.Context, groupId string) error {
	req := client.sdkClient1d0.IamV1GroupsApiAPI.DeleteGroup(ctx, groupId)

	_, err := req.Execute()
	return err
}

func (client *Client) GetGroupMembers(ctx context.Context, groupId string, request GroupMembersDataResource) (*scpiam1d0.GroupMemberPageResponse, error) {
	req := client.sdkClient1d0.IamV1GroupsApiAPI.ListGroupMember(ctx, groupId)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.UserName.IsNull() {
		req = req.UserName(request.UserName.ValueString())
	}
	if !request.UserEmail.IsNull() {
		req = req.UserName(request.UserEmail.ValueString())
	}
	if !request.CreatorName.IsNull() {
		req = req.UserName(request.CreatorName.ValueString())
	}
	if !request.CreatorEmail.IsNull() {
		req = req.UserName(request.CreatorEmail.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) AddGroupMember(ctx context.Context, groupId string, request GroupMemberResource) (*scpiam1d0.GroupMemberCreateResponse, error) {
	req := client.sdkClient1d0.IamV1GroupsApiAPI.AddGroupMember(ctx, groupId)

	if !request.UserId.IsNull() {
		req = req.GroupMemberCreateRequest(scpiam1d0.GroupMemberCreateRequest{UserId: request.UserId.ValueString()})
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) RemoveGroupMember(ctx context.Context, groupId string, request GroupMemberResource) error {
	req := client.sdkClient1d0.IamV1GroupsApiAPI.RemoveGroupMember(ctx, groupId, request.UserId.ValueString())

	_, err := req.Execute()
	return err
}

func (client *Client) GetGroupPolicyBindings(ctx context.Context, groupId string, request GroupPolicyBindingsDataResource) (*scpiam1d0.GroupPolicyPageResponse, error) {
	req := client.sdkClient1d0.IamV1GroupsApiAPI.ListGroupPolicyBinding(ctx, groupId)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.PolicyId.IsNull() {
		req = req.PolicyId(request.PolicyId.ValueString())
	}
	if !request.PolicyName.IsNull() {
		req = req.PolicyName(request.PolicyName.ValueString())
	}
	if !request.PolicyType.IsNull() {
		req = req.PolicyType(scpiam1d0.PolicyType{nil, request.PolicyType.ValueStringPointer()})
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) AddGroupPolicyBindings(ctx context.Context, groupId string, request GroupPolicyBindingsResource) (*scpiam1d0.GroupPolicyResponse, error) {
	req := client.sdkClient1d0.IamV1GroupsApiAPI.AddGroupPolicyBinding(ctx, groupId)

	var policyIds []string
	for _, policyId := range request.PolicyIds {
		policyIds = append(policyIds, policyId.ValueString())
	}

	req = req.GroupPolicyBindingRequest(scpiam1d0.GroupPolicyBindingRequest{PolicyIds: policyIds})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) RemoveGroupPolicyBindings(ctx context.Context, groupId string, request GroupPolicyBindingsResource) error {

	for _, policyId := range request.PolicyIds {
		req := client.sdkClient1d0.IamV1GroupsApiAPI.RemoveGroupPolicyBinding(ctx, groupId, policyId.ValueString())

		_, err := req.Execute()

		if err != nil {
			return err
		}
	}

	return nil
}

// / POLICY ///
func (client *Client) GetPolicies(ctx context.Context, request PolicyDatasource) (*scpiam1d0.PolicyPageResponse, error) {
	req := client.sdkClient1d0.IamV1PoliciesApiAPI.ListPolicy(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}
	if !request.PolicyName.IsNull() {
		req = req.PolicyName(request.PolicyName.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPolicy(ctx context.Context, policyId string) (*scpiam1d0.PolicyShowResponse, error) {
	req := client.sdkClient1d0.IamV1PoliciesApiAPI.ShowPolicy(ctx, policyId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreatePolicy(ctx context.Context, request PolicyResource) (*scpiam1d0.PolicyShowResponse, error) {
	req := client.sdkClient1d0.IamV1PoliciesApiAPI.CreatePolicy(ctx)

	//tag
	var TagsObject []map[string]string
	for k, v := range request.Tags.Elements() {
		tagObject := make(map[string]string)
		tagObject["key"] = k
		tagObject["value"] = v.(types.String).ValueString()

		TagsObject = append(TagsObject, tagObject)
	}

	//policy version
	var statements []scpiam1d0.Statement
	for _, _statement := range request.PolicyVersion.PolicyDocument.Statement {

		// resource
		var resources []string
		for _, _resource := range _statement.Resource {
			resources = append(resources, _resource.ValueString())
		}

		// action
		var actions []string
		for _, _action := range _statement.Action {
			actions = append(actions, _action.ValueString())
		}

		// not action
		var notActions []string
		for _, _notAction := range _statement.NotAction {
			notActions = append(notActions, _notAction.ValueString())
		}

		// condition
		condition := map[string]map[string][]interface{}{}

		for key, value := range _statement.Condition.Elements() {
			typeMap, ok := value.(basetypes.MapValue)
			if !ok || typeMap.IsUnknown() || typeMap.IsNull() {
				continue
			}

			innerMap := convertMapToGoInnerCondition(typeMap)

			condition[key] = innerMap
		}

		statement := scpiam1d0.Statement{
			Sid:       _statement.Sid.ValueStringPointer(),
			Effect:    _statement.Effect.ValueString(),
			Resource:  resources,
			Action:    actions,
			NotAction: notActions,
		}

		if len(condition) > 0 {
			statement.SetCondition(condition)
		} else {
			statement.SetCondition(nil)
		}

		statements = append(statements, statement)
	}

	policyVersion := scpiam1d0.PolicyVersionCreateRequest{
		PolicyDocument: scpiam1d0.IamPolicyDocument{
			Statement: statements,
			Version:   request.PolicyVersion.PolicyDocument.Version.ValueString(),
		},
	}

	req = req.PolicyCreateRequest(scpiam1d0.PolicyCreateRequest{
		Description:   *scpiam1d0.NewNullableString(request.Description.ValueStringPointer()),
		PolicyName:    request.PolicyName.ValueString(),
		Tags:          TagsObject,
		PolicyVersion: policyVersion,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePolicy(ctx context.Context, policyId string, request PolicyResource) (*scpiam1d0.PolicyShowResponse, error) {
	req := client.sdkClient1d0.IamV1PoliciesApiAPI.SetPolicy(ctx, policyId)

	//policy version
	var statements []scpiam1d0.Statement
	for _, _statement := range request.PolicyVersion.PolicyDocument.Statement {

		// resource
		var resources []string
		for _, _resource := range _statement.Resource {
			resources = append(resources, _resource.ValueString())
		}

		// action
		var actions []string
		for _, _action := range _statement.Action {
			actions = append(actions, _action.ValueString())
		}

		// not action
		var notActions []string
		for _, _notAction := range _statement.NotAction {
			notActions = append(notActions, _notAction.ValueString())
		}

		// condition
		condition := map[string]map[string][]interface{}{}

		for key, value := range _statement.Condition.Elements() {
			typeMap, ok := value.(basetypes.MapValue)
			if !ok || typeMap.IsUnknown() || typeMap.IsNull() {
				continue
			}

			innerMap := convertMapToGoInnerCondition(typeMap)

			condition[key] = innerMap
		}

		statement := scpiam1d0.Statement{
			Sid:       _statement.Sid.ValueStringPointer(),
			Effect:    _statement.Effect.ValueString(),
			Resource:  resources,
			Action:    actions,
			NotAction: notActions,
		}

		if len(condition) > 0 {
			statement.SetCondition(condition)
		} else {
			statement.SetCondition(nil)
		}

		statements = append(statements, statement)
	}

	//policy version
	policyVersion := scpiam1d0.PolicyVersionCreateRequest{
		PolicyDocument: scpiam1d0.IamPolicyDocument{
			Statement: statements,
			Version:   request.PolicyVersion.PolicyDocument.Version.ValueString(),
		},
	}

	req = req.PolicySetRequest(scpiam1d0.PolicySetRequest{
		Description:   *scpiam1d0.NewNullableString(request.Description.ValueStringPointer()),
		PolicyName:    *scpiam1d0.NewNullableString(request.PolicyName.ValueStringPointer()),
		PolicyVersion: *scpiam1d0.NewNullablePolicyVersionCreateRequest(&policyVersion),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeletePolicy(ctx context.Context, policyId string) error {
	req := client.sdkClient1d0.IamV1PoliciesApiAPI.DeletePolicy(ctx, policyId)

	_, err := req.Execute()
	return err
}

func convertMapToGoInnerCondition(m basetypes.MapValue) map[string][]interface{} {
	result := map[string][]interface{}{}

	if m.IsNull() || m.IsUnknown() {
		return result
	}

	for key, value := range m.Elements() {
		// value는 types.ListValue 타입이어야 함
		listVal, ok := value.(basetypes.ListValue)
		if !ok || listVal.IsUnknown() || listVal.IsNull() {
			continue
		}

		stringValues := []interface{}{}
		for _, v := range listVal.Elements() {
			s, ok := v.(basetypes.StringValue)
			if !ok || s.IsUnknown() || s.IsNull() {
				continue
			}
			stringValues = append(stringValues, s.ValueString())
		}

		result[key] = stringValues
	}

	return result
}

// / ROLE ///
func (client *Client) GetRoles(ctx context.Context, request RoleDataSource) (*scpiam1d0.RolePageResponse, error) {
	req := client.sdkClient1d0.IamV1RolesApiAPI.ListRole(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}

	if !request.RoleType.IsNull() {
		req = req.RoleType(request.RoleType.ValueString())
	}

	if !request.AccountId.IsNull() {
		req = req.AccountId(request.AccountId.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetRole(ctx context.Context, roleId string) (*scpiam1d0.RoleShowResponse, error) {
	req := client.sdkClient1d0.IamV1RolesApiAPI.ShowRole(ctx, roleId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateRole(ctx context.Context, request RoleResource) (*scpiam1d0.RoleShowResponse, error) {
	req := client.sdkClient1d0.IamV1RolesApiAPI.CreateRole(ctx)

	// tag
	var TagsObject []map[string]string

	for k, v := range request.Tags.Elements() {
		tagObject := make(map[string]string)
		tagObject["key"] = k
		tagObject["value"] = v.(types.String).ValueString()

		TagsObject = append(TagsObject, tagObject)
	}

	//policy ids
	var policyIds []string
	for _, policyId := range request.PolicyIds {
		policyIds = append(policyIds, policyId.ValueString())
	}

	//pricipals
	var roleTrustPolicyPrincipals []scpiam1d0.RoleTrustPolicyPrincipal
	for _, principal := range request.Principals {
		roleTrustPolicyPrincipals = append(roleTrustPolicyPrincipals, scpiam1d0.RoleTrustPolicyPrincipal{
			Type:  principal.Type.ValueString(),
			Value: principal.Value.ValueString(),
		})
	}

	// assume role policy
	var statements []scpiam1d0.Statement
	var version string
	if request.AssumeRolePolicyDocument != nil {
		version = request.AssumeRolePolicyDocument.Version.ValueString()

		for _, _statement := range request.AssumeRolePolicyDocument.Statement {

			// resource
			var resources []string
			for _, _resource := range _statement.Resource {
				resources = append(resources, _resource.ValueString())
			}

			// action
			var actions []string
			for _, _action := range _statement.Action {
				actions = append(actions, _action.ValueString())
			}

			// not action
			var notActions []string
			for _, _notAction := range _statement.NotAction {
				notActions = append(notActions, _notAction.ValueString())
			}

			// principal
			var principalData interface{}

			isPrincipalStringSet := !_statement.Principal.PrincipalString.IsNull()
			isPrincipalMapSet := !_statement.Principal.PrincipalMap.IsNull()

			if isPrincipalStringSet {
				principalData = _statement.Principal.PrincipalString.ValueStringPointer()
			} else if isPrincipalMapSet {
				var principalMap = map[string][]string{}
				for key, value := range _statement.Principal.PrincipalMap.Elements() {
					listVal, ok := value.(basetypes.ListValue)
					if !ok || listVal.IsUnknown() || listVal.IsNull() {
						continue
					}
					var principalItem []string
					for _, v := range listVal.Elements() {
						s, ok := v.(basetypes.StringValue)
						if !ok || s.IsUnknown() || s.IsNull() {
							continue
						}
						principalItem = append(principalItem, s.ValueString())
					}
					principalMap[key] = principalItem
				}
				principalData = principalMap
			} else {
				principalData = nil
			}

			principal := createNullablePrincipal(principalData)

			// condition
			condition := map[string]map[string][]interface{}{}
			for key, value := range _statement.Condition.Elements() {
				typeMap, ok := value.(basetypes.MapValue)
				if !ok || typeMap.IsUnknown() || typeMap.IsNull() {
					continue
				}
				innerMap := convertMapToGoInnerCondition(typeMap)
				condition[key] = innerMap
			}

			statement := scpiam1d0.Statement{
				Sid:       _statement.Sid.ValueStringPointer(),
				Effect:    _statement.Effect.ValueString(),
				Resource:  resources,
				Action:    actions,
				NotAction: notActions,
				Principal: *principal,
			}

			if len(condition) > 0 {
				statement.SetCondition(condition)
			} else {
				statement.SetCondition(nil)
			}

			statements = append(statements, statement)
		}
	}

	policyDocument := scpiam1d0.NewNullablePolicyDocument(&scpiam1d0.PolicyDocument{
		Statement: statements,
		Version:   version,
	})

	if policyDocument.Get().Statement == nil {
		policyDocument.Unset()
	}

	req = req.RoleCreateRequest(scpiam1d0.RoleCreateRequest{
		Description:              *scpiam1d0.NewNullableString(request.Description.ValueStringPointer()),
		Name:                     request.Name.ValueString(),
		MaxSessionDuration:       request.MaxSessionDuration.ValueInt32Pointer(),
		Tags:                     TagsObject,
		PolicyIds:                policyIds,
		Principals:               roleTrustPolicyPrincipals,
		AssumeRolePolicyDocument: *policyDocument,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func createNullablePrincipal(principalData interface{}) *scpiam1d0.NullablePrincipal {
	if principalData == nil {
		return scpiam1d0.NewNullablePrincipal(nil)
	}
	switch v := principalData.(type) {
	case string:
		principal := &scpiam1d0.Principal{
			String: &v,
		}
		return scpiam1d0.NewNullablePrincipal(principal)
	case map[string][]string:
		principal := &scpiam1d0.Principal{
			MapmapOfStringarrayOfString: &v,
		}
		return scpiam1d0.NewNullablePrincipal(principal)
	default:
		return scpiam1d0.NewNullablePrincipal(nil)
	}
}

func (client *Client) UpdateRole(ctx context.Context, roleId string, request RoleResource) (*scpiam1d0.RoleShowResponse, error) {
	req := client.sdkClient1d0.IamV1RolesApiAPI.SetRole(ctx, roleId)

	req = req.RoleSetRequest(scpiam1d0.RoleSetRequest{
		Description:        *scpiam1d0.NewNullableString(request.Description.ValueStringPointer()),
		MaxSessionDuration: *scpiam1d0.NewNullableInt32(request.MaxSessionDuration.ValueInt32Pointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteRole(ctx context.Context, roleId string) error {
	req := client.sdkClient1d0.IamV1RolesApiAPI.DeleteRole(ctx, roleId)
	_, err := req.Execute()
	return err
}

func (client *Client) GetRolePolicyBindings(ctx context.Context, roleId string, request RolePolicyBindingsDataSource) (*scpiam1d0.RolePolicyBindingPageResponse, error) {
	req := client.sdkClient1d0.IamV1RolesApiAPI.ListRolePolicyBindings(ctx, roleId)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.PolicyId.IsNull() {
		req = req.PolicyId(request.PolicyId.ValueString())
	}
	if !request.PolicyName.IsNull() {
		req = req.PolicyName(request.PolicyName.ValueString())
	}
	if !request.PolicyType.IsNull() {
		req = req.PolicyType(scpiam1d0.PolicyType3{nil, request.PolicyType.ValueStringPointer()})
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) AddRolePolicyBindings(ctx context.Context, roleId string, request RolePolicyBindingsResource) (*scpiam1d0.RolePolicyBindingResponse, error) {
	req := client.sdkClient1d0.IamV1RolesApiAPI.AddRolePolicyBindings(ctx, roleId)

	var policyIds []string
	for _, policyId := range request.PolicyIds {
		policyIds = append(policyIds, policyId.ValueString())
	}

	req = req.RolePolicyBindingRequest(scpiam1d0.RolePolicyBindingRequest{
		PolicyIds: policyIds,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) RemoveRolePolicyBindings(ctx context.Context, roleId string, request RolePolicyBindingsResource) error {
	for _, policyId := range request.PolicyIds {
		req := client.sdkClient1d0.IamV1RolesApiAPI.RemoveRolePolicyBinding(ctx, roleId, policyId.ValueString())

		_, err := req.Execute()

		if err != nil {
			return err
		}
	}

	return nil
}

// / USER ///
func (client *Client) GetUsers(ctx context.Context, accountId string, request UserDataSource) (*scpiam1d1.ListIAMUserResponse, error) {
	req := client.sdkClient1d1.IamV1AccountsApiAPI.ListIAMUser(ctx, accountId)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Email.IsNull() {
		req = req.Email(request.Email.ValueString())
	}
	if !request.Email.IsNull() {
		req = req.UserName(request.UserName.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetUser(ctx context.Context, accountId string, userId string) (*scpiam1d1.IAMUserDetailResponse, error) {
	req := client.sdkClient1d1.IamV1AccountsApiAPI.GetIAMUser(ctx, accountId, userId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateUser(ctx context.Context, request UserResource) (*scpiam1d1.IAMCreateUserResponse, error) {
	req := client.sdkClient1d1.IamV1AccountsApiAPI.CreateIAMUser(ctx, request.AccountId.ValueString())

	//group ids
	var groupIds []string
	for _, groupId := range request.GroupIds {
		groupIds = append(groupIds, groupId.ValueString())
	}

	//policy ids
	var policyIds []string
	for _, policyId := range request.PolicyIds {
		policyIds = append(policyIds, policyId.ValueString())
	}

	//tag
	var TagsObject []map[string]string

	for k, v := range request.Tags.Elements() {
		tagObject := make(map[string]string)
		tagObject["key"] = k
		tagObject["value"] = v.(types.String).ValueString()

		TagsObject = append(TagsObject, tagObject)
	}

	req = req.IAMUserCreateRequest(scpiam1d1.IAMUserCreateRequest{
		Description:       *scpiam1d1.NewNullableString(request.Description.ValueStringPointer()),
		UserName:          request.UserName.ValueString(),
		Password:          request.Password.ValueString(),
		TemporaryPassword: request.TemporaryPassword.ValueBool(),
		GroupIds:          groupIds,
		PolicyIds:         policyIds,
		Tags:              TagsObject,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateUser(ctx context.Context, accountId string, userId string, request UserResource) (*scpiam1d1.UserResponse, error) {
	req := client.sdkClient1d1.IamV1AccountsApiAPI.UpdateIAMUser(ctx, accountId, userId)

	req = req.IAMUserUpdateRequest(scpiam1d1.IAMUserUpdateRequest{
		Description:        *scpiam1d1.NewNullableString(request.Description.ValueStringPointer()),
		PasswordReuseCount: request.PasswordReuseCount.ValueInt32(),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteUser(ctx context.Context, accountId string, userId string) error {
	req := client.sdkClient1d1.IamV1AccountsApiAPI.DeleteIAMUser(ctx, accountId, userId)

	_, err := req.Execute()
	return err
}

func (client *Client) GetUserPolicyBindings(ctx context.Context, userId string, request UserPolicyBindingsDataSource) (*scpiam1d1.UserPolicyPageResponse, error) {
	req := client.sdkClient1d1.IamV1UsersApiAPI.ListUserPolicyBindings(ctx, userId)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.PolicyId.IsNull() {
		req = req.PolicyId(request.PolicyId.ValueString())
	}
	if !request.PolicyName.IsNull() {
		req = req.PolicyName(request.PolicyName.ValueString())
	}
	if !request.PolicyType.IsNull() {
		req = req.PolicyType(scpiam1d1.PolicyType4{nil, request.PolicyType.ValueStringPointer()})
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) AddUserPolicyBindings(ctx context.Context, userId string, request UserPolicyBindingsResource) (*scpiam1d1.UserPolicyResponse, error) {
	req := client.sdkClient1d1.IamV1UsersApiAPI.AddUserPolicyBinding(ctx, userId)

	var policyIds []string
	for _, policyId := range request.PolicyIds {
		policyIds = append(policyIds, policyId.ValueString())
	}

	req = req.UserPolicyRequest(scpiam1d1.UserPolicyRequest{
		PolicyIds: policyIds,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) RemoveUserPolicyBindings(ctx context.Context, userId string, request UserPolicyBindingsResource) error {
	for _, policyId := range request.PolicyIds {
		req := client.sdkClient1d1.IamV1UsersApiAPI.RemoveUserPolicyBinding(ctx, userId, policyId.ValueString())

		_, err := req.Execute()

		if err != nil {
			return err
		}
	}

	return nil
}
