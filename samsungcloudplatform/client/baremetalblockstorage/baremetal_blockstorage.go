package baremetalblockstorage

import (
	"context"
	"fmt"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	baremetalblockstorage1d0 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/baremetal-blockstorage/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"time"
)

type Client struct {
	Config       *scpsdk.Configuration
	sdkClient1d0 *baremetalblockstorage1d0.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:       config,
		sdkClient1d0: baremetalblockstorage1d0.NewAPIClient(config),
	}
}

func (client *Client) CreateBlockStorage(ctx context.Context, request VolumeResource) (*baremetalblockstorage1d0.AsyncResponse, error) {
	req := client.sdkClient1d0.BaremetalBlockstorageV1VolumeV1APIsAPI.CreateVolume(ctx)

	diskType, _ := baremetalblockstorage1d0.NewDiskTypeFromValue(request.DiskType.ValueString())

	attachments := client.getAttachmentListModelList(request.Attachments)

	tags := make([]baremetalblockstorage1d0.TagModel, 0)
	for k, v := range request.Tags.Elements() {
		tag := baremetalblockstorage1d0.TagModel{}

		key := baremetalblockstorage1d0.NullableString{}
		key.Set(&k)

		tag.Key = key

		if v != nil {
			value := baremetalblockstorage1d0.NullableString{}
			value.Set(v.(types.String).ValueStringPointer())
			tag.Value = value
		}
		tags = append(tags, tag)
	}

	req = req.VolumeCreateRequest(baremetalblockstorage1d0.VolumeCreateRequest{
		Name:        request.Name.ValueString(),
		DiskType:    *diskType,
		SizeGb:      request.SizeGb.ValueInt32(),
		Attachments: attachments,
		Tags:        tags,
	})

	response, _, err := req.Execute()
	return response, err
}

func (client *Client) GetBlockStorage(ctx context.Context, blockStorageId string) (*baremetalblockstorage1d0.VolumeResponse, int, error) {
	req := client.sdkClient1d0.BaremetalBlockstorageV1VolumeV1APIsAPI.ShowVolume(ctx, blockStorageId)
	response, c, err := req.Execute()
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) AttachBlockStorages(ctx context.Context, blockStorageId string, attachObjectList []Attachment) (*baremetalblockstorage1d0.VolumeResponse, int, error) {
	req := client.sdkClient1d0.BaremetalBlockstorageV1VolumeV1APIsAPI.CreateVolumeAttachments(ctx, blockStorageId)

	attachments := client.getAttachmentListModelList(attachObjectList)

	req = req.VolumeAttachmentRequest(baremetalblockstorage1d0.VolumeAttachmentRequest{Attachments: attachments})

	response, c, err := req.Execute()

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) DetachBlockStorages(ctx context.Context, blockStorageId string, detachObjectIdList []string) (*baremetalblockstorage1d0.VolumeResponse, int, error) {
	req := client.sdkClient1d0.BaremetalBlockstorageV1VolumeV1APIsAPI.DeleteVolumeAttachments(ctx, blockStorageId)

	req = req.VolumeDetachRequest(baremetalblockstorage1d0.VolumeDetachRequest{
		Attachments: detachObjectIdList,
	})

	response, c, err := req.Execute()

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) DeleteBlockStorage(ctx context.Context, blockStorageId string) (*baremetalblockstorage1d0.AsyncResponse, int, error) {
	req := client.sdkClient1d0.BaremetalBlockstorageV1VolumeV1APIsAPI.DeleteVolume(ctx, blockStorageId)

	response, c, err := req.Execute()

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) WaitForStatus(ctx context.Context, pendingStates []string, targetStates []string, timeout time.Duration,
	refreshFunc retry.StateRefreshFunc) error {
	stateConf := &retry.StateChangeConf{
		Pending:    pendingStates,
		Target:     targetStates,
		Refresh:    refreshFunc,
		Timeout:    timeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting : %s", err)
	}

	return nil
}

func (client *Client) getAttachmentListModelList(attachmentList []Attachment) []baremetalblockstorage1d0.AttachmentListModel {
	attachments := make([]baremetalblockstorage1d0.AttachmentListModel, 0)
	for _, attachment := range attachmentList {
		objectType, _ := baremetalblockstorage1d0.NewBlockStorageAttachmentObjectTypeFromValue(attachment.ObjectType.ValueString())
		attachments = append(attachments, baremetalblockstorage1d0.AttachmentListModel{
			ObjectType: objectType,
			ObjectId:   attachment.ObjectId.ValueStringPointer(),
		})
	}
	return attachments
}
