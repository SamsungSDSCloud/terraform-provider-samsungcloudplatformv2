package parallelfilestorage

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/client"
	scpparallelfilestorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/library/parallelfilestorage/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpparallelfilestorage.APIClient // 서비스의 client 를 구조체에 추가한다.
}

func NewClient(config *scpsdk.Configuration) *Client { // client 생성 함수를 추가한다.
	return &Client{
		Config:    config,
		sdkClient: scpparallelfilestorage.NewAPIClient(config),
	}
}

// Parallel FileStorage
func (client *Client) GetVolumeList(ctx context.Context, request VolumeDataSourceIds) (*scpparallelfilestorage.VolumeListResponseV1Dot1, error) {
	req := client.sdkClient.ParallelFilestorageV1VolumesAPIsAPI.ListVolumes(ctx)
	if !request.Offset.IsNull() {
		req = req.Offset(request.Offset.ValueInt32())
	}
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
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

func (client *Client) CreateVolume(ctx context.Context, request VolumeResource) (*scpparallelfilestorage.VolumeV1Dot1, error) {
	req := client.sdkClient.ParallelFilestorageV1VolumesAPIsAPI.CreateVolume(ctx)

	if request.Id.ValueString() == "" {
		request.Id = types.StringNull()
	}

	// Tags
	var TagsObjects []scpparallelfilestorage.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := scpparallelfilestorage.Tag{
			Key:   k,
			Value: *scpparallelfilestorage.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObjects = append(TagsObjects, tagObject)
	}

	req = req.VolumeCreateRequestV1Dot1(scpparallelfilestorage.VolumeCreateRequestV1Dot1{
		Name:         request.Name.ValueString(),
		CapacityTb:   request.CapacityTb.ValueInt32(),
		Tags:         TagsObjects,
		Zone:  		  request.Zone.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVolume(ctx context.Context, id string) (*scpparallelfilestorage.VolumeShowResponseV1Dot1, error) {
	req := client.sdkClient.ParallelFilestorageV1VolumesAPIsAPI.ShowVolume(ctx, id)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVolumeCapacityTb(ctx context.Context, id string, request VolumeResource) error {
	req := client.sdkClient.ParallelFilestorageV1VolumesAPIsAPI.SetVolumeCapacity(ctx, id)

	req = req.VolumeCapacityRequest(scpparallelfilestorage.VolumeCapacityRequest{
		CapacityTb: request.CapacityTb.ValueInt32(),
	})

	_, err := req.Execute()
	return err
}

func (client *Client) DeleteVolume(ctx context.Context, id string) error {
	req := client.sdkClient.ParallelFilestorageV1VolumesAPIsAPI.DeleteVolume(ctx, id)

	_, err := req.Execute()
	return err
}

// AccessRules
func (client *Client) GetVolumeAccessRules(ctx context.Context, id string) (*scpparallelfilestorage.AccessRuleListResponse, error) {
	req := client.sdkClient.ParallelFilestorageV1AccessRulesAPIsAPI.ListAccessRule(ctx, id)
	resp, _, err := req.Execute()
	return resp, err
}

func toSDKAccessRule(in AccessRuleResource) scpparallelfilestorage.AccessRule {
	return scpparallelfilestorage.AccessRule{
		ObjectId:   in.ObjectId.ValueString(),
		ObjectType: in.ObjectType.ValueString(),
	}
}

func (client *Client) UpdateVolumeAccessRule(ctx context.Context, id string, request []AccessRuleResource, action string) error {
	req := client.sdkClient.ParallelFilestorageV1AccessRulesAPIsAPI.SetAccessRule(ctx, id)
	
	sdkRules := make([]scpparallelfilestorage.AccessRule, len(request))

	var addAccessRules []scpparallelfilestorage.AccessRule
	var removeAccessRules []scpparallelfilestorage.AccessRule

	for i, r := range request {
		sdkRules[i] = toSDKAccessRule(r)
	}

	if action == "add" {
		addAccessRules = sdkRules
		removeAccessRules = []scpparallelfilestorage.AccessRule{}
	} else {
		removeAccessRules = sdkRules
		addAccessRules = []scpparallelfilestorage.AccessRule{}
	}
	req = req.AccessRuleUpdateRequest(scpparallelfilestorage.AccessRuleUpdateRequest{
		AddAccessRules:    addAccessRules,
		RemoveAccessRules: removeAccessRules,
	})
	_, _, err := req.Execute()
	return err
}
