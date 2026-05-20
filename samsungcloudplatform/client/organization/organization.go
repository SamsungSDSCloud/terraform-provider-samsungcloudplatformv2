package organization

import (
	"context"
	"fmt"
	"net/http"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	organization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *organization.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: organization.NewAPIClient(config),
	}
}

func (client *Client) GetOrganizationList(ctx context.Context, request OrganizationDataSource) (*organization.OrganizationPageResponse, error) {
	req := client.sdkClient.OrganizationV1OrganizationsAPIsAPI.ListOrganizations(ctx)
	if !request.Size.IsNull() {
		req = req.Size(int32(request.Size.ValueInt64()))
	} else {
		req = req.Size(20) // 기본값
	}
	if !request.Page.IsNull() {
		req = req.Page(int32(request.Page.ValueInt64()))
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.MasterAccountId.IsNull() {
		req = req.MasterAccountId(request.MasterAccountId.ValueString())
	}
	if !request.OrganizationId.IsNull() {
		req = req.OrganizationId(request.OrganizationId.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetOrganization(ctx context.Context, organizationId string) (*organization.OrganizationShowResponse, error) {
	req := client.sdkClient.OrganizationV1OrganizationsAPIsAPI.ShowOrganization(ctx, organizationId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateOrganization(ctx context.Context, request OrganizationResource) (*organization.OrganizationShowResponse, error) {
	req := client.sdkClient.OrganizationV1OrganizationsAPIsAPI.CreateOrganization(ctx)

	req = req.OrganizationCreateRequest(organization.OrganizationCreateRequest{
		Name: request.Name.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateOrganization(ctx context.Context, organizationId string, request OrganizationResource) (*organization.OrganizationShowResponse, error) {
	req := client.sdkClient.OrganizationV1OrganizationsAPIsAPI.SetOrganization(ctx, organizationId)

	var delegationAccountId organization.NullableString
	if !request.DelegationAccountId.IsNull() && !request.DelegationAccountId.IsUnknown() {
		delegationAccountId = *organization.NewNullableString(request.DelegationAccountId.ValueStringPointer())
	}

	var name organization.NullableString
	if !request.Name.IsNull() && !request.Name.IsUnknown() {
		name = *organization.NewNullableString(request.Name.ValueStringPointer())
	}

	var useScpYn organization.NullableBool
	if !request.UseScpYn.IsNull() && !request.UseScpYn.IsUnknown() {
		useScpYn = *organization.NewNullableBool(request.UseScpYn.ValueBoolPointer())
	}

	req = req.OrganizationSetRequest(organization.OrganizationSetRequest{
		DelegationAccountId: delegationAccountId,
		Name:                name,
		OrganizationId:      *organization.NewNullableString(&organizationId),
		UseScpYn:            useScpYn,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteOrganization(ctx context.Context, organizationId string) (*organization.OrganizationDeleteResponse, error) {
	req := client.sdkClient.OrganizationV1OrganizationsAPIsAPI.DeleteOrganization(ctx, organizationId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetOrganizationUnit(ctx context.Context, unitId string, organizationId string) (*organization.OrganizationUnitShowResponse, error) {
	req := client.sdkClient.OrganizationV1OrganizationUnitsAPIsAPI.ShowOrganizationUnit(ctx, unitId)

	if organizationId != "" {
		req = req.OrganizationId(organizationId)
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetOrganizationUnits(ctx context.Context, parentUnitId string, organizationId string, name string, excludePolicyId string) (*organization.OrganizationUnitListResponseV1dot1, error) {
	req := client.sdkClient.OrganizationV1OrganizationUnitsAPIsAPI.ListOrganizationUnits(ctx)
	req = req.ParentUnitId(parentUnitId)

	if organizationId != "" {
		req = req.OrganizationId(organizationId)
	}
	if name != "" {
		req = req.Name(name)
	}
	if excludePolicyId != "" {
		req = req.ExcludePolicyId(excludePolicyId)
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetOrganizationUnitParents(ctx context.Context, unitId string, organizationId string) (*organization.OrganizationUnitParentsResponse, error) {
	req := client.sdkClient.OrganizationV1OrganizationUnitsAPIsAPI.ListParents(ctx, unitId)

	if organizationId != "" {
		req = req.OrganizationId(organizationId)
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateOrganizationUnit(ctx context.Context, request OrganizationUnitResource) (*organization.OrganizationUnitCreateResponse, error) {
	req := client.sdkClient.OrganizationV1OrganizationUnitsAPIsAPI.CreateOrganizationUnit(ctx)

	var description organization.NullableString
	if !request.Description.IsNull() && !request.Description.IsUnknown() {
		description = *organization.NewNullableString(request.Description.ValueStringPointer())
	}

	var organizationId organization.NullableString
	if !request.OrganizationId.IsNull() && !request.OrganizationId.IsUnknown() {
		organizationId = *organization.NewNullableString(request.OrganizationId.ValueStringPointer())
	}

	var parentUnitId organization.NullableString
	if !request.ParentUnitId.IsNull() && !request.ParentUnitId.IsUnknown() {
		parentUnitId = *organization.NewNullableString(request.ParentUnitId.ValueStringPointer())
	}

	var policyIds []string
	for _, policyId := range request.PolicyIds {
		policyIds = append(policyIds, policyId.ValueString())
	}

	req = req.OrganizationUnitCreateRequest(organization.OrganizationUnitCreateRequest{
		Description:    description,
		Name:           request.Name.ValueString(),
		OrganizationId: organizationId,
		ParentUnitId:   parentUnitId,
		PolicyIds:      policyIds,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateInvitation(ctx context.Context, request InvitationResource) (*organization.InvitationCreateResponse, error) {
	req := client.sdkClient.OrganizationV1InvitationsAPIsAPI.CreateInvitation(ctx)

	var organizationId organization.NullableString
	if !request.OrganizationId.IsNull() && !request.OrganizationId.IsUnknown() {
		organizationId = *organization.NewNullableString(request.OrganizationId.ValueStringPointer())
	}

	targetLoginIds := make([]string, 0, len(request.TargetLoginIds.Elements()))
	for _, elem := range request.TargetLoginIds.Elements() {
		targetLoginIds = append(targetLoginIds, elem.(types.String).ValueString())
	}

	req = req.InvitationCreateRequest(organization.InvitationCreateRequest{
		OrganizationId: organizationId,
		TargetLoginIds: targetLoginIds,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CancelInvitations(ctx context.Context, request InvitationCancelRequest) (*organization.InvitationCancelResponse, error) {
	req := client.sdkClient.OrganizationV1InvitationsAPIsAPI.CancelInvitations(ctx)

	var organizationId string
	if !request.OrganizationId.IsNull() && !request.OrganizationId.IsUnknown() {
		organizationId = request.OrganizationId.ValueString()
	}

	var ids []string
	if !request.Ids.IsNull() && !request.Ids.IsUnknown() {
		elems := request.Ids.Elements()
		ids = make([]string, 0, len(elems))
		for _, e := range elems {
			if sv, ok := e.(types.String); ok {
				ids = append(ids, sv.ValueString())
			}
		}
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no invitation IDs to cancel")
	}

	cancelReq := organization.InvitationCancelRequest{
		Ids: ids,
	}
	if organizationId != "" {
		cancelReq.SetOrganizationId(organizationId)
	}

	req = req.InvitationCancelRequest(cancelReq)

	resp, rawResp, err := req.Execute()
	if err != nil {
		errMsg := fmt.Sprintf("CancelInvitations error: %v (orgId: %q, ids: %v)", err, organizationId, ids)
		if sdkErr, ok := err.(*scpsdk.GenericOpenAPIError); ok {
			errMsg += fmt.Sprintf(", SDK Error Body: %s", string(sdkErr.Body()))
		} else if rawResp != nil {
			errMsg += fmt.Sprintf(", StatusCode: %d", rawResp.StatusCode)
		}
		return nil, fmt.Errorf(errMsg)
	}
	return resp, err
}

func (client *Client) AcceptInvitation(ctx context.Context, invitationId string) (*organization.InvitationAcceptResponse, error) {
	req := client.sdkClient.OrganizationV1InvitationsAPIsAPI.AcceptInvitation(ctx, invitationId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeclineInvitation(ctx context.Context, invitationId string) (*organization.InvitationShowResponse, error) {
	req := client.sdkClient.OrganizationV1InvitationsAPIsAPI.DeclineInvitation(ctx, invitationId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetAccountInvitations(ctx context.Context) (*organization.AccountInvitationResponse, error) {
	req := client.sdkClient.OrganizationV1InvitationsAPIsAPI.ListAccountInvitations(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) GetOrganizationInvitations(ctx context.Context, request OrganizationInvitationsDataSource) (*organization.OrganizationInvitationPageResponse, error) {
	req := c.sdkClient.OrganizationV1InvitationsAPIsAPI.ListOrganizationInvitations(ctx)

	if !request.OrganizationId.IsNull() {
		req = req.OrganizationId(request.OrganizationId.ValueString())
	}
	if !request.Size.IsNull() {
		req = req.Size(int32(request.Size.ValueInt64()))
	}
	if !request.Page.IsNull() {
		req = req.Page(int32(request.Page.ValueInt64()))
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.AccountId.IsNull() {
		req = req.AccountId(request.AccountId.ValueString())
	}
	if !request.AccountName.IsNull() {
		req = req.AccountName(request.AccountName.ValueString())
	}
	if !request.AccountEmail.IsNull() {
		req = req.AccountEmail(request.AccountEmail.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(organization.InvitationState(request.State.ValueString()))
	}
	if !request.LoginId.IsNull() {
		req = req.LoginId(request.LoginId.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateOrganizationUnit(ctx context.Context, unitId string, request OrganizationUnitResource) (*organization.OrganizationUnitShowResponse, error) {
	req := client.sdkClient.OrganizationV1OrganizationUnitsAPIsAPI.SetOrganizationUnit(ctx, unitId)

	var description organization.NullableString
	if !request.Description.IsNull() && !request.Description.IsUnknown() {
		description = *organization.NewNullableString(request.Description.ValueStringPointer())
	}

	var name organization.NullableString
	if !request.Name.IsNull() && !request.Name.IsUnknown() {
		name = *organization.NewNullableString(request.Name.ValueStringPointer())
	}

	req = req.OrganizationUnitSetRequest(organization.OrganizationUnitSetRequest{
		Description: description,
		Name:        name,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteOrganizationUnit(ctx context.Context, unitId string, organizationId string) (*organization.OrganizationUnitDeleteResponse, error) {
	req := client.sdkClient.OrganizationV1OrganizationUnitsAPIsAPI.DeleteOrganizationUnits(ctx)

	orgId := organizationId
	req = req.OrganizationUnitDeleteRequest(organization.OrganizationUnitDeleteRequest{
		Ids:            []string{unitId},
		OrganizationId: *organization.NewNullableString(&orgId),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPoliciesForTarget(ctx context.Context, targetId string, organizationId string) (*organization.ListPoliciesForTargetResponse, error) {
	req := client.sdkClient.OrganizationV1AssignmentsAPIsAPI.ListPoliciesForTarget(ctx)
	req = req.TargetId(targetId)

	if organizationId != "" {
		req = req.OrganizationId(organizationId)
	}

	resp, _, err := req.Execute()
	return resp, err
}

// =====================
// Organization Account Methods
// =====================

func (client *Client) GetAccount(ctx context.Context, accountId string) (*organization.AccountShowResponse, error) {
	req := client.sdkClient.OrganizationV1AccountsAPIsAPI.ShowAccount(ctx, accountId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) ListAccounts(ctx context.Context, request AccountListDataSourceRequest) (*organization.AccountPageResponse, error) {
	req := client.sdkClient.OrganizationV1AccountsAPIsAPI.ListAccounts(ctx)

	if !request.OrganizationId.IsNull() {
		req = req.OrganizationId(request.OrganizationId.ValueString())
	}
	if !request.Size.IsNull() {
		req = req.Size(int32(request.Size.ValueInt64()))
	} else {
		req = req.Size(100) // 기본값
	}
	if !request.Page.IsNull() {
		req = req.Page(int32(request.Page.ValueInt64()))
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.Email.IsNull() {
		req = req.Email(request.Email.ValueString())
	}
	if !request.LoginId.IsNull() {
		req = req.LoginId(request.LoginId.ValueString())
	}
	if !request.JoinedStartDate.IsNull() {
		req = req.JoinedStartDate(request.JoinedStartDate.ValueString())
	}
	if !request.JoinedEndDate.IsNull() {
		req = req.JoinedEndDate(request.JoinedEndDate.ValueString())
	}
	if !request.JoinedMethod.IsNull() {
		req = req.JoinedMethod(organization.JoinedMethodEnum(request.JoinedMethod.ValueString()))
	}
	if !request.ExcludePolicyId.IsNull() {
		req = req.ExcludePolicyId(request.ExcludePolicyId.ValueString())
	}
	if !request.ParentUnitId.IsNull() {
		req = req.ParentUnitId(request.ParentUnitId.ValueString())
	}
	if !request.ParentUnitName.IsNull() {
		req = req.ParentUnitName(request.ParentUnitName.ValueString())
	}
	if !request.Type.IsNull() {
		req = req.Type_(organization.OrganizationAccountTypeEnum(request.Type.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateAccount(ctx context.Context, loginId, name, organizationId, roleName string) (*organization.AccountCreateResponse, error) {
	req := client.sdkClient.OrganizationV1AccountsAPIsAPI.CreateAccount(ctx)

	req = req.OrganizationUserCreateRequest(organization.OrganizationUserCreateRequest{
		LoginId:        loginId,
		Name:           name,
		OrganizationId: organizationId,
		RoleName:       roleName,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) MoveAccount(ctx context.Context, accountId, parentUnitId, organizationId string) (*organization.AccountMoveResponse, error) {
	return client.MoveAccounts(ctx, []string{accountId}, parentUnitId, organizationId)
}

func (client *Client) MoveAccounts(ctx context.Context, accountIds []string, parentUnitId, organizationId string) (*organization.AccountMoveResponse, error) {
	req := client.sdkClient.OrganizationV1AccountsAPIsAPI.MoveAccount(ctx)

	var orgId organization.NullableString
	if organizationId != "" {
		orgId = *organization.NewNullableString(&organizationId)
	}

	req = req.MoveAccountRequest(organization.MoveAccountRequest{
		ParentUnitId:     parentUnitId,
		TargetAccountIds: accountIds,
		OrganizationId:   orgId,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteAccount(ctx context.Context, accountId string) (*http.Response, error) {
	req := client.sdkClient.OrganizationV1AccountsAPIsAPI.DeleteAccount(ctx, accountId)

	resp, err := req.Execute()
	return resp, err
}

func (client *Client) RemoveAccounts(ctx context.Context, targetAccountIds []string, organizationId string) (*organization.AccountRemoveResponse, error) {
	req := client.sdkClient.OrganizationV1AccountsAPIsAPI.RemoveAccounts(ctx)

	var orgId organization.NullableString
	if organizationId != "" {
		orgId = *organization.NewNullableString(&organizationId)
	}

	req = req.RemoveAccountsRequest(organization.RemoveAccountsRequest{
		TargetAccountIds: targetAccountIds,
		OrganizationId:   orgId,
	})

	resp, _, err := req.Execute()
	return resp, err
}

// =====================
// Delegation Policy Methods
// =====================

func (client *Client) CreateDelegationPolicy(ctx context.Context, request DelegationPolicyResource) (*organization.DelegationPolicyShowResponse, error) {
	req := client.sdkClient.OrganizationV1DelegationPolicyAPIsAPI.CreateDelegationPolicy(ctx)

	var organizationId organization.NullableString
	if !request.OrganizationId.IsNull() && !request.OrganizationId.IsUnknown() {
		organizationId = *organization.NewNullableString(request.OrganizationId.ValueStringPointer())
	}

	document, err := request.Document.toSdkType()
	if err != nil {
		return nil, fmt.Errorf("failed to convert document: %w", err)
	}

	req = req.DelegationPolicyCreateRequest(organization.DelegationPolicyCreateRequest{
		Document:       *document,
		OrganizationId: organizationId,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetDelegationPolicy(ctx context.Context, organizationId string) (*organization.DelegationPolicyShowResponse, error) {
	req := client.sdkClient.OrganizationV1DelegationPolicyAPIsAPI.ShowDelegationPolicy(ctx)

	if organizationId != "" {
		req = req.OrganizationId(organizationId)
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateDelegationPolicy(ctx context.Context, request DelegationPolicyResource) (*organization.DelegationPolicyShowResponse, error) {
	req := client.sdkClient.OrganizationV1DelegationPolicyAPIsAPI.SetDelegationPolicy(ctx)

	var organizationId organization.NullableString
	if !request.OrganizationId.IsNull() && !request.OrganizationId.IsUnknown() {
		organizationId = *organization.NewNullableString(request.OrganizationId.ValueStringPointer())
	}

	document, err := request.Document.toSdkType()
	if err != nil {
		return nil, fmt.Errorf("failed to convert document: %w", err)
	}

	req = req.DelegationPolicySetRequest(organization.DelegationPolicySetRequest{
		Document:       *document,
		OrganizationId: organizationId,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteDelegationPolicy(ctx context.Context, organizationId string) (*http.Response, error) {
	req := client.sdkClient.OrganizationV1DelegationPolicyAPIsAPI.DeleteDelegationPolicy(ctx)

	if organizationId != "" {
		req = req.OrganizationId(organizationId)
	}

	resp, err := req.Execute()
	return resp, err
}

func (client *Client) GetServiceControlPolicy(ctx context.Context, policyId string, organizationId string) (*organization.ServiceControlPolicyShowResponse, error) {
	req := client.sdkClient.OrganizationV1ServiceControlPoliciesAPIsAPI.ShowServiceControlPolicy(ctx, policyId)

	if organizationId != "" {
		req = req.OrganizationId(organizationId)
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) ListServiceControlPolicies(ctx context.Context, request ServiceControlPolicyListDataSourceRequest) (*organization.ServiceControlPolicyPageResponse, error) {
	req := c.sdkClient.OrganizationV1ServiceControlPoliciesAPIsAPI.ListServiceControlPolicies(ctx)

	if !request.OrganizationId.IsNull() {
		req = req.OrganizationId(request.OrganizationId.ValueString())
	}
	if !request.Name.IsNull() {
		name := request.Name.ValueString()
		req = req.Name(organization.Name{String: &name})
	}
	if !request.Type.IsNull() {
		req = req.Type_(organization.PolicyTypeEnum(request.Type.ValueString()))
	}
	if !request.ExcludeTargetId.IsNull() {
		req = req.ExcludeTargetId(request.ExcludeTargetId.ValueString())
	}
	if !request.Size.IsNull() {
		req = req.Size(int32(request.Size.ValueInt64()))
	}
	if !request.Page.IsNull() {
		req = req.Page(int32(request.Page.ValueInt64()))
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateServiceControlPolicy(ctx context.Context, request ServiceControlPolicyResource) (*organization.ServiceControlPolicyShowResponse, error) {
	req := client.sdkClient.OrganizationV1ServiceControlPoliciesAPIsAPI.CreateServiceControlPolicy(ctx)

	var description organization.NullableString
	if !request.Description.IsNull() && !request.Description.IsUnknown() {
		description = *organization.NewNullableString(request.Description.ValueStringPointer())
	}

	document, err := documentObjectToSdkType(request.Document)
	if err != nil {
		return nil, fmt.Errorf("failed to convert document: %w", err)
	}

	var organizationId string
	if !request.OrganizationId.IsNull() && !request.OrganizationId.IsUnknown() {
		organizationId = request.OrganizationId.ValueString()
	}

	var policyType *organization.PolicyTypeEnum
	if !request.Type.IsNull() && !request.Type.IsUnknown() {
		pt := organization.PolicyTypeEnum(request.Type.ValueString())
		policyType = &pt
	}

	req = req.ServiceControlPolicyCreateRequest(organization.ServiceControlPolicyCreateRequest{
		Description:    description,
		Document:       *document,
		Name:           request.Name.ValueString(),
		OrganizationId: organizationId,
		Type:           policyType,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateServiceControlPolicy(ctx context.Context, policyId string, organizationId string, request ServiceControlPolicyResource) (*organization.ServiceControlPolicyShowResponse, error) {
	req := client.sdkClient.OrganizationV1ServiceControlPoliciesAPIsAPI.SetServiceControlPolicy(ctx, policyId)

	var description organization.NullableString
	if !request.Description.IsNull() && !request.Description.IsUnknown() {
		description = *organization.NewNullableString(request.Description.ValueStringPointer())
	}

	var name organization.NullableString
	if !request.Name.IsNull() && !request.Name.IsUnknown() {
		name = *organization.NewNullableString(request.Name.ValueStringPointer())
	}

	document, err := documentObjectToSdkType(request.Document)
	if err != nil {
		return nil, fmt.Errorf("failed to convert document: %w", err)
	}

	req = req.ServiceControlPolicySetRequest(organization.ServiceControlPolicySetRequest{
		Description:    description,
		Document:       *organization.NewNullableServiceControlPolicyDocument(document),
		Name:           name,
		OrganizationId: organizationId,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteServiceControlPolicies(ctx context.Context, organizationId string, policyIds []string) (*organization.ServiceControlPolicyDeleteResponse, error) {
	req := client.sdkClient.OrganizationV1ServiceControlPoliciesAPIsAPI.DeleteServiceControlPolicies(ctx)

	req = req.ServiceControlPolicyDeleteRequest(organization.ServiceControlPolicyDeleteRequest{
		Ids:            policyIds,
		OrganizationId: organizationId,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (o ServiceControlPolicyDocumentValue) toSdkType() (*organization.ServiceControlPolicyDocument, error) {
	doc := organization.NewServiceControlPolicyDocument()

	var statements []organization.ServiceControlPolicyStatement
	if !o.Statement.IsNull() && !o.Statement.IsUnknown() {
		statementElems := o.Statement.Elements()
		statements = make([]organization.ServiceControlPolicyStatement, 0, len(statementElems))
		for _, stmtElem := range statementElems {
			stmtObj, ok := stmtElem.(types.Object)
			if !ok {
				continue
			}

			stmtAttrs := stmtObj.Attributes()
			stmt := organization.NewServiceControlPolicyStatement()

			if sidAttr, ok := stmtAttrs["sid"]; ok && !sidAttr.IsNull() && !sidAttr.IsUnknown() {
				stmt.SetSid(sidAttr.(types.String).ValueString())
			}

			if effectAttr, ok := stmtAttrs["effect"]; ok && !effectAttr.IsNull() && !effectAttr.IsUnknown() {
				stmt.SetEffect(effectAttr.(types.String).ValueString())
			}

			if actionAttr, ok := stmtAttrs["action"]; ok && !actionAttr.IsNull() && !actionAttr.IsUnknown() {
				actionList, ok := actionAttr.(types.List)
				if ok {
					actionElems := actionList.Elements()
					actions := make([]string, 0, len(actionElems))
					for _, a := range actionElems {
						if s, ok := a.(types.String); ok {
							actions = append(actions, s.ValueString())
						}
					}
					stmt.SetAction(actions)
				}
			}

			if notActionAttr, ok := stmtAttrs["not_action"]; ok && !notActionAttr.IsNull() && !notActionAttr.IsUnknown() {
				notActionList, ok := notActionAttr.(types.List)
				if ok {
					notActionElems := notActionList.Elements()
					notActions := make([]string, 0, len(notActionElems))
					for _, a := range notActionElems {
						if s, ok := a.(types.String); ok {
							notActions = append(notActions, s.ValueString())
						}
					}
					stmt.SetNotAction(notActions)
				}
			}

			if resourceAttr, ok := stmtAttrs["resource"]; ok && !resourceAttr.IsNull() && !resourceAttr.IsUnknown() {
				resourceList, ok := resourceAttr.(types.List)
				if ok {
					resourceElems := resourceList.Elements()
					resources := make([]string, 0, len(resourceElems))
					for _, r := range resourceElems {
						if s, ok := r.(types.String); ok {
							resources = append(resources, s.ValueString())
						}
					}
					stmt.SetResource(resources)
				}
			}

			statements = append(statements, *stmt)
		}
		doc.SetStatement(statements)
	}

	if !o.Version.IsNull() && !o.Version.IsUnknown() {
		doc.SetVersion(o.Version.ValueString())
	}

	return doc, nil
}

func documentObjectToSdkType(docObj types.Object) (*organization.ServiceControlPolicyDocument, error) {
	doc := organization.NewServiceControlPolicyDocument()

	attrs := docObj.Attributes()

	if stmtAttr, ok := attrs["statement"]; ok && !stmtAttr.IsNull() && !stmtAttr.IsUnknown() {
		statementList, ok := stmtAttr.(types.List)
		if ok {
			statementElems := statementList.Elements()
			statements := make([]organization.ServiceControlPolicyStatement, 0, len(statementElems))
			for _, stmtElem := range statementElems {
				stmtObj, ok := stmtElem.(types.Object)
				if !ok {
					continue
				}

				stmtAttrs := stmtObj.Attributes()
				stmt := organization.NewServiceControlPolicyStatement()

				if sidAttr, ok := stmtAttrs["sid"]; ok && !sidAttr.IsNull() && !sidAttr.IsUnknown() {
					stmt.SetSid(sidAttr.(types.String).ValueString())
				}

				if effectAttr, ok := stmtAttrs["effect"]; ok && !effectAttr.IsNull() && !effectAttr.IsUnknown() {
					stmt.SetEffect(effectAttr.(types.String).ValueString())
				}

				if actionAttr, ok := stmtAttrs["action"]; ok && !actionAttr.IsNull() && !actionAttr.IsUnknown() {
					actionList, ok := actionAttr.(types.List)
					if ok {
						actionElems := actionList.Elements()
						actions := make([]string, 0, len(actionElems))
						for _, a := range actionElems {
							if s, ok := a.(types.String); ok {
								actions = append(actions, s.ValueString())
							}
						}
						stmt.SetAction(actions)
					}
				}

				if notActionAttr, ok := stmtAttrs["not_action"]; ok && !notActionAttr.IsNull() && !notActionAttr.IsUnknown() {
					notActionList, ok := notActionAttr.(types.List)
					if ok {
						notActionElems := notActionList.Elements()
						notActions := make([]string, 0, len(notActionElems))
						for _, a := range notActionElems {
							if s, ok := a.(types.String); ok {
								notActions = append(notActions, s.ValueString())
							}
						}
						stmt.SetNotAction(notActions)
					}
				}

				if resourceAttr, ok := stmtAttrs["resource"]; ok && !resourceAttr.IsNull() && !resourceAttr.IsUnknown() {
					resourceList, ok := resourceAttr.(types.List)
					if ok {
						resourceElems := resourceList.Elements()
						resources := make([]string, 0, len(resourceElems))
						for _, r := range resourceElems {
							if s, ok := r.(types.String); ok {
								resources = append(resources, s.ValueString())
							}
						}
						stmt.SetResource(resources)
					}
				}

				statements = append(statements, *stmt)
			}
			doc.SetStatement(statements)
		}
	}

	if versionAttr, ok := attrs["version"]; ok && !versionAttr.IsNull() && !versionAttr.IsUnknown() {
		doc.SetVersion(versionAttr.(types.String).ValueString())
	}

	return doc, nil
}
