package smp

import (
	"context"
	"math"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	smp "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/smp/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *smp.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: smp.NewAPIClient(config),
	}
}

// ----------------------------
// Workspace API
// ----------------------------

func (client *Client) GetWorkspaceList(ctx context.Context, request WorkspaceDataSources) (*smp.WorkspacePageResponse, error) {
	req := client.sdkClient.SmpV1WorkspacesAPIsAPI.ListWorkspaces(ctx)
	if !request.Size.IsNull() && !request.Size.IsUnknown() {
		req = req.Size(int32(request.Size.ValueInt64()))
	} else {
		req = req.Size(math.MaxInt32)
	}
	if !request.Page.IsNull() && !request.Page.IsUnknown() {
		req = req.Page(int32(request.Page.ValueInt64()))
	}
	if !request.Id.IsNull() && request.Id.ValueString() != "" {
		req = req.Id(request.Id.ValueString())
	}
	if !request.IdLike.IsNull() && request.IdLike.ValueString() != "" {
		if request.Id.IsNull() || request.Id.ValueString() == "" {
			req = req.Id(request.IdLike.ValueString())
		}
	}
	if !request.Name.IsNull() && request.Name.ValueString() != "" {
		req = req.Name(request.Name.ValueString())
	}
	if !request.NameLike.IsNull() && request.NameLike.ValueString() != "" {
		if request.Name.IsNull() || request.Name.ValueString() == "" {
			req = req.Name(request.NameLike.ValueString())
		}
	}
	if !request.SortBy.IsNull() && request.SortBy.ValueString() != "" {
		sort := request.SortBy.ValueString()
		if !request.SortOrder.IsNull() && request.SortOrder.ValueString() != "" {
			sort = sort + ":" + request.SortOrder.ValueString()
		}
		req = req.Sort(sort)
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetWorkspace(ctx context.Context, workspaceId string) (*smp.WorkspaceShowResponse, error) {
	req := client.sdkClient.SmpV1WorkspacesAPIsAPI.ShowWorkspace(ctx, workspaceId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateWorkspace(ctx context.Context, request WorkspaceResource) (*smp.WorkspaceShowResponse, error) {
	req := client.sdkClient.SmpV1WorkspacesAPIsAPI.CreateWorkspace(ctx)
	createReq := smp.WorkspaceCreateRequest{
		Name:              request.Name.ValueString(),
		PrivateAclEnabled: true,
	}
	if !request.PrivateAclEnabled.IsNull() && !request.PrivateAclEnabled.IsUnknown() {
		createReq.PrivateAclEnabled = request.PrivateAclEnabled.ValueBool()
	}
	if !request.PrivateAclResources.IsNull() && !request.PrivateAclResources.IsUnknown() {
		var resources []smp.PrivateAclResource
		for _, v := range request.PrivateAclResources.Elements() {
			obj := v.(types.Object)
			raw := obj.Attributes()
			resourceId := raw["resource_id"].(types.String).ValueString()
			resourceType := raw["resource_type"].(types.String).ValueString()
			res := smp.PrivateAclResource{
				ResourceId:   resourceId,
				ResourceType: resourceType,
			}
			resourceName := raw["resource_name"].(types.String)
			if !resourceName.IsNull() && !resourceName.IsUnknown() {
				res.SetResourceName(resourceName.ValueString())
			}
			resources = append(resources, res)
		}
		if len(resources) > 0 {
			createReq.PrivateAclResources = resources
		}
	}
	req = req.WorkspaceCreateRequest(createReq)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateWorkspace(ctx context.Context, workspaceId string, request WorkspaceResource) (*smp.WorkspaceShowResponse, error) {
	req := client.sdkClient.SmpV1WorkspacesAPIsAPI.SetWorkspace(ctx, workspaceId)
	req = req.WorkspaceSetRequest(smp.WorkspaceSetRequest{
		Name: request.Name.ValueString(),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteWorkspace(ctx context.Context, workspaceId string) error {
	req := client.sdkClient.SmpV1WorkspacesAPIsAPI.DeleteWorkspace(ctx, workspaceId)
	_, err := req.Execute()
	return err
}

func (client *Client) SetWorkspacePrivateAclResources(ctx context.Context, workspaceId string, request PrivateAclResourcesSetRequest) (*smp.PrivateAclResourcesSetResponse, error) {
	req := client.sdkClient.SmpV1WorkspacesAPIsAPI.SetWorkspacePrivateAclResources(ctx, workspaceId)
	sdkReq := smp.PrivateAclResourcesSetRequest{}
	if !request.PrivateAclEnabled.IsNull() && !request.PrivateAclEnabled.IsUnknown() {
		sdkReq.SetPrivateAclEnabled(request.PrivateAclEnabled.ValueBool())
	}
	if !request.PrivateAclResources.IsNull() && !request.PrivateAclResources.IsUnknown() {
		var resources []smp.PrivateAclResource
		for _, v := range request.PrivateAclResources.Elements() {
			obj := v.(types.Object)
			raw := obj.Attributes()
			resourceId := raw["resource_id"].(types.String).ValueString()
			resourceType := raw["resource_type"].(types.String).ValueString()
			res := smp.PrivateAclResource{
				ResourceId:   resourceId,
				ResourceType: resourceType,
			}
			resourceName := raw["resource_name"].(types.String)
			if !resourceName.IsNull() && !resourceName.IsUnknown() {
				res.SetResourceName(resourceName.ValueString())
			}
			resources = append(resources, res)
		}
		sdkReq.PrivateAclResources = resources
	}
	req = req.PrivateAclResourcesSetRequest(sdkReq)
	resp, _, err := req.Execute()
	return resp, err
}

// ----------------------------
// RuleNamespace API (all require workspaceId)
// ----------------------------

func (client *Client) GetRuleNamespaceList(ctx context.Context, workspaceId string, request RuleNamespaceDataSources) (*smp.RuleNamespacePageResponse, error) {
	req := client.sdkClient.SmpV1RuleNamespacesAPIsAPI.ListRuleNamespaces(ctx, workspaceId)
	if !request.Size.IsNull() && !request.Size.IsUnknown() {
		req = req.Size(int32(request.Size.ValueInt64()))
	} else {
		req = req.Size(math.MaxInt32)
	}
	if !request.Page.IsNull() && !request.Page.IsUnknown() {
		req = req.Page(int32(request.Page.ValueInt64()))
	}
	if !request.Id.IsNull() && request.Id.ValueString() != "" {
		req = req.Id(request.Id.ValueString())
	}
	if !request.IdLike.IsNull() && request.IdLike.ValueString() != "" {
		if request.Id.IsNull() || request.Id.ValueString() == "" {
			req = req.Id(request.IdLike.ValueString())
		}
	}
	if !request.Name.IsNull() && request.Name.ValueString() != "" {
		req = req.Name(request.Name.ValueString())
	}
	if !request.NameLike.IsNull() && request.NameLike.ValueString() != "" {
		if request.Name.IsNull() || request.Name.ValueString() == "" {
			req = req.Name(request.NameLike.ValueString())
		}
	}
	if !request.SortBy.IsNull() && request.SortBy.ValueString() != "" {
		sort := request.SortBy.ValueString()
		if !request.SortOrder.IsNull() && request.SortOrder.ValueString() != "" {
			sort = sort + ":" + request.SortOrder.ValueString()
		}
		req = req.Sort(sort)
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetRuleNamespace(ctx context.Context, workspaceId string, ruleNamespaceId string) (*smp.RuleNamespaceShowResponse, error) {
	req := client.sdkClient.SmpV1RuleNamespacesAPIsAPI.ShowRuleNamespace(ctx, workspaceId, ruleNamespaceId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateRuleNamespace(ctx context.Context, workspaceId string, request RuleNamespaceResource) (*smp.RuleNamespaceShowResponse, error) {
	req := client.sdkClient.SmpV1RuleNamespacesAPIsAPI.CreateRuleNamespace(ctx, workspaceId)
	req = req.RuleNamespaceCreateRequest(smp.RuleNamespaceCreateRequest{
		Name:       request.Name.ValueString(),
		ConfigData: request.ConfigData.ValueString(),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateRuleNamespace(ctx context.Context, workspaceId string, ruleNamespaceId string, request RuleNamespaceResource) (*smp.RuleNamespaceShowResponse, error) {
	req := client.sdkClient.SmpV1RuleNamespacesAPIsAPI.SetRuleNamespace(ctx, workspaceId, ruleNamespaceId)
	req = req.RuleNamespaceSetRequest(smp.RuleNamespaceSetRequest{
		ConfigData: request.ConfigData.ValueString(),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteRuleNamespace(ctx context.Context, workspaceId string, ruleNamespaceId string) error {
	req := client.sdkClient.SmpV1RuleNamespacesAPIsAPI.DeleteRuleNamespace(ctx, workspaceId, ruleNamespaceId)
	_, err := req.Execute()
	return err
}

// ----------------------------
// AlertManager API (all require workspaceId)
// ----------------------------

func (client *Client) GetAlertManager(ctx context.Context, workspaceId string) (*smp.AlertManagerShowResponse, error) {
	req := client.sdkClient.SmpV1AlertManagerAPIsAPI.ShowAlertManager(ctx, workspaceId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateAlertManager(ctx context.Context, workspaceId string, request AlertManagerResource) (*smp.AlertManagerShowResponse, error) {
	req := client.sdkClient.SmpV1AlertManagerAPIsAPI.CreateAlertManager(ctx, workspaceId)
	req = req.AlertManagerCreateRequest(smp.AlertManagerCreateRequest{
		ConfigData:          request.ConfigData.ValueString(),
		NotificationGroupId: request.NotificationGroupId.ValueString(),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateAlertManager(ctx context.Context, workspaceId string, request AlertManagerResource) (*smp.AlertManagerShowResponse, error) {
	req := client.sdkClient.SmpV1AlertManagerAPIsAPI.SetAlertManager(ctx, workspaceId)
	req = req.AlertManagerCreateRequest(smp.AlertManagerCreateRequest{
		ConfigData:          request.ConfigData.ValueString(),
		NotificationGroupId: request.NotificationGroupId.ValueString(),
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteAlertManager(ctx context.Context, workspaceId string) error {
	req := client.sdkClient.SmpV1AlertManagerAPIsAPI.DeleteAlertManager(ctx, workspaceId)
	_, err := req.Execute()
	return err
}

// ----------------------------
// WorkspaceConfiguration API (Show + Set only — no Create/Delete in SDK)
// ----------------------------

func (client *Client) GetWorkspaceConfiguration(ctx context.Context, workspaceId string) (*smp.WorkspaceConfigurationShowResponse, error) {
	req := client.sdkClient.SmpV1WorkspaceConfigurationAPIsAPI.ShowWorkspaceConfiguration(ctx, workspaceId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateWorkspaceConfiguration(ctx context.Context, request WorkspaceConfigurationResource) (*smp.WorkspaceConfigurationShowResponse, error) {
	// SDK has no Create endpoint; use Set as create
	workspaceId := request.WorkspaceId.ValueString()
	req := client.sdkClient.SmpV1WorkspaceConfigurationAPIsAPI.SetWorkspaceConfiguration(ctx, workspaceId)
	retentionPeriod := int32(0)
	if !request.RetentionPeriod.IsNull() && !request.RetentionPeriod.IsUnknown() {
		retentionPeriod = int32(request.RetentionPeriod.ValueInt64())
	}
	req = req.WorkspaceConfigurationCreateRequest(smp.WorkspaceConfigurationCreateRequest{
		RetentionPeriod: retentionPeriod,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateWorkspaceConfiguration(ctx context.Context, workspaceId string, request WorkspaceConfigurationResource) (*smp.WorkspaceConfigurationShowResponse, error) {
	req := client.sdkClient.SmpV1WorkspaceConfigurationAPIsAPI.SetWorkspaceConfiguration(ctx, workspaceId)
	retentionPeriod := int32(0)
	if !request.RetentionPeriod.IsNull() && !request.RetentionPeriod.IsUnknown() {
		retentionPeriod = int32(request.RetentionPeriod.ValueInt64())
	}
	req = req.WorkspaceConfigurationCreateRequest(smp.WorkspaceConfigurationCreateRequest{
		RetentionPeriod: retentionPeriod,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteWorkspaceConfiguration(ctx context.Context, workspaceId string) error {
	// SDK has no Delete endpoint for workspace configuration.
	// Return nil — the resource cannot be deleted via API.
	return nil
}

// ----------------------------
// NotificationGroup API
// ----------------------------

func (client *Client) GetNotificationGroupList(ctx context.Context, request NotificationGroupDataSources) (*smp.NotificationGroupPageResponse, error) {
	req := client.sdkClient.SmpV1NotificationGroupsAPIsAPI.ListNotificationGroups(ctx)
	if !request.Size.IsNull() && !request.Size.IsUnknown() {
		req = req.Size(int32(request.Size.ValueInt64()))
	} else {
		req = req.Size(math.MaxInt32)
	}
	if !request.Page.IsNull() && !request.Page.IsUnknown() {
		req = req.Page(int32(request.Page.ValueInt64()))
	}
	if !request.Id.IsNull() && request.Id.ValueString() != "" {
		req = req.Id(request.Id.ValueString())
	}
	if !request.IdLike.IsNull() && request.IdLike.ValueString() != "" {
		if request.Id.IsNull() || request.Id.ValueString() == "" {
			req = req.Id(request.IdLike.ValueString())
		}
	}
	if !request.Name.IsNull() && request.Name.ValueString() != "" {
		req = req.Name(request.Name.ValueString())
	}
	if !request.NameLike.IsNull() && request.NameLike.ValueString() != "" {
		if request.Name.IsNull() || request.Name.ValueString() == "" {
			req = req.Name(request.NameLike.ValueString())
		}
	}
	if !request.SortBy.IsNull() && request.SortBy.ValueString() != "" {
		sort := request.SortBy.ValueString()
		if !request.SortOrder.IsNull() && request.SortOrder.ValueString() != "" {
			sort = sort + ":" + request.SortOrder.ValueString()
		}
		req = req.Sort(sort)
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetNotificationGroup(ctx context.Context, notificationGroupId string) (*smp.NotificationGroupShowResponse, error) {
	req := client.sdkClient.SmpV1NotificationGroupsAPIsAPI.ShowNotificationGroup(ctx, notificationGroupId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateNotificationGroup(ctx context.Context, request NotificationGroupResource) (*smp.NotificationGroupResponse, error) {
	req := client.sdkClient.SmpV1NotificationGroupsAPIsAPI.CreateNotificationGroup(ctx)
	createReq := smp.NotificationGroupCreateRequest{
		Name: request.Name.ValueString(),
	}
	if !request.Description.IsNull() {
		createReq.SetDescription(request.Description.ValueString())
	}
	if !request.RecipientUserIds.IsNull() && !request.RecipientUserIds.IsUnknown() {
		var ids []string
		for _, v := range request.RecipientUserIds.Elements() {
			ids = append(ids, v.(types.String).ValueString())
		}
		createReq.RecipientUserIds = ids
	}
	req = req.NotificationGroupCreateRequest(createReq)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateNotificationGroup(ctx context.Context, notificationGroupId string, request NotificationGroupResource) (*smp.NotificationGroupResponse, error) {
	req := client.sdkClient.SmpV1NotificationGroupsAPIsAPI.SetNotificationGroup(ctx, notificationGroupId)
	setReq := smp.NotificationGroupSetRequest{
		Name: request.Name.ValueString(),
	}
	if !request.Description.IsNull() {
		setReq.SetDescription(request.Description.ValueString())
	}
	req = req.NotificationGroupSetRequest(setReq)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateNotificationGroupRecipients(ctx context.Context, notificationGroupId string, request NotificationGroupResource) error {
	req := client.sdkClient.SmpV1NotificationGroupsAPIsAPI.SetNotificationGroupRecipient(ctx, notificationGroupId)
	var ids []string
	for _, v := range request.RecipientUserIds.Elements() {
		ids = append(ids, v.(types.String).ValueString())
	}
	recipientReq := smp.NotificationGroupRecipientRequest{
		RecipientUserIds: ids,
	}
	req = req.NotificationGroupRecipientRequest(recipientReq)
	_, _, err := req.Execute()
	return err
}

func (client *Client) DeleteNotificationGroup(ctx context.Context, notificationGroupId string) error {
	req := client.sdkClient.SmpV1NotificationGroupsAPIsAPI.DeleteNotificationGroup(ctx, notificationGroupId)
	_, err := req.Execute()
	return err
}
