package resourcemanager

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/resourcemanager/1.0"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *resourcemanager.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: resourcemanager.NewAPIClient(config),
	}
}

// Tag

func (client *Client) GetTagList(ctx context.Context, request TagDataSource) (*resourcemanager.TagListResponse, error) {
	req := client.sdkClient.ResourcemanagerV1TagsAPIsAPI.ListTags(ctx)
	if !request.Key.IsNull() {
		req = req.Key(request.Key.ValueString())
	}
	if !request.Value.IsNull() {
		req = req.Value(request.Value.ValueString())
	}
	if !request.ResourceIdentifier.IsNull() {
		req = req.ResourceIdentifier(request.ResourceIdentifier.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

// Resource Tag (srn)

func (client *Client) GetResourceTagList(ctx context.Context, request ResourceTagDataSource) (*resourcemanager.TagShowListResponse, error) {
	req := client.sdkClient.ResourcemanagerV1TagsAPIsAPI.ListResourceTags(ctx, request.EncodedSrn.ValueString())
	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetResourceTags(srn string) (*resourcemanager.TagShowListResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.ResourcemanagerV1TagsAPIsAPI.ListResourceTags(ctx, srn)
	req = req.Size(math.MaxInt32)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) BulkUpdateResourceTags(srn string, tagElements map[string]attr.Value) (*resourcemanager.TagBaseResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.ResourcemanagerV1TagsAPIsAPI.UpdateResourceTags(ctx, srn)

	tagsObject := make([]resourcemanager.Tag, 0)

	for k, v := range tagElements {
		tagObject := resourcemanager.Tag{
			Key:   k,
			Value: *resourcemanager.NewNullableString(v.(types.String).ValueStringPointer()),
		}

		tagsObject = append(tagsObject, tagObject)
	}

	req = req.TagSetRequest(resourcemanager.TagSetRequest{
		Tags: tagsObject,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteResourceTags(ctx context.Context, srn string) error {
	req := client.sdkClient.ResourcemanagerV1TagsAPIsAPI.DeleteResourceTags(ctx, srn)

	_, err := req.Execute()
	return err
}

// Resource Tag (srn, key)

func (client *Client) GetResourceTag(ctx context.Context, srn string, key string) (*resourcemanager.TagShowResponse, error) {
	req := client.sdkClient.ResourcemanagerV1TagsAPIsAPI.ShowResourceTag(ctx, srn, key)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateResourceTag(ctx context.Context, srn string, key string, value string) (*resourcemanager.TagBaseResponse, error) {
	req := client.sdkClient.ResourcemanagerV1TagsAPIsAPI.UpdateResourceTagValue(ctx, srn, key)

	req = req.TagValue(resourcemanager.TagValue{
		Value: value,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteResourceTag(ctx context.Context, srn string, key string) error {
	req := client.sdkClient.ResourcemanagerV1TagsAPIsAPI.DeleteResourceTag(ctx, srn, key)

	_, err := req.Execute()
	return err
}

// Resource Group

func (client *Client) GetResourceGroupList(Id types.String, Name types.String) (*resourcemanager.ResourceGroupPageResponse, error) {
	ctx := context.Background()

	req := client.sdkClient.ResourcemanagerV1ResourceGroupsAPIsAPI.ListResourceGroups(ctx)
	req = req.Size(math.MaxInt32)
	if !Id.IsNull() {
		req = req.Sort(Id.ValueString())
	}
	if !Name.IsNull() {
		req = req.Sort(Name.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateResourceGroup(ctx context.Context, request ResourceGroupResource) (*resourcemanager.ResourceGroupCreateResponse, error) {
	req := client.sdkClient.ResourcemanagerV1ResourceGroupsAPIsAPI.CreateResourceGroup(ctx)

	var GroupDefinitionTagsObject []resourcemanager.Tag

	for k, v := range request.GroupDefinitionTags.Elements() {
		tagObject := resourcemanager.Tag{
			Key:   k,
			Value: *resourcemanager.NewNullableString(v.(types.String).ValueStringPointer()),
		}

		GroupDefinitionTagsObject = append(GroupDefinitionTagsObject, tagObject)
	}

	var resourceTypes []string
	request.ResourceTypes.ElementsAs(ctx, &resourceTypes, false)

	req = req.ResourceGroupCreateRequest(resourcemanager.ResourceGroupCreateRequest{
		Description:   *resourcemanager.NewNullableString(request.Description.ValueStringPointer()),
		Name:          request.Name.ValueString(),
		ResourceTypes: resourceTypes,
		Tags:          GroupDefinitionTagsObject,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetResourceGroup(ctx context.Context, resourceGroupId string) (*resourcemanager.ResourceGroupShowResponse, error) {
	req := client.sdkClient.ResourcemanagerV1ResourceGroupsAPIsAPI.ShowResourceGroup(ctx, resourceGroupId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateResourceGroup(ctx context.Context, resourceGroupId string, request ResourceGroupResource) (*resourcemanager.ResourceGroupCreateResponse, error) {
	req := client.sdkClient.ResourcemanagerV1ResourceGroupsAPIsAPI.SetResourceGroup(ctx, resourceGroupId)

	var GroupDefinitionTagsObject []resourcemanager.Tag

	for k, v := range request.GroupDefinitionTags.Elements() {
		tagObject := resourcemanager.Tag{
			Key:   k,
			Value: *resourcemanager.NewNullableString(v.(types.String).ValueStringPointer()),
		}

		GroupDefinitionTagsObject = append(GroupDefinitionTagsObject, tagObject)
	}

	var resourceTypes []string
	request.ResourceTypes.ElementsAs(ctx, &resourceTypes, false)

	req = req.ResourceGroupUpdateRequest(resourcemanager.ResourceGroupUpdateRequest{
		Description:   *resourcemanager.NewNullableString(request.Description.ValueStringPointer()),
		ResourceTypes: resourceTypes,
		Tags:          GroupDefinitionTagsObject,
	})
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteResourceGroup(ctx context.Context, resourceGroupId string) error {
	req := client.sdkClient.ResourcemanagerV1ResourceGroupsAPIsAPI.DeleteResourceGroup(ctx, resourceGroupId)

	_, err := req.Execute()
	return err
}
