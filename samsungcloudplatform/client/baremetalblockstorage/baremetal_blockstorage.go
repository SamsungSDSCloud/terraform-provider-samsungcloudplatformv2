package baremetalblockstorage

import (
	"context"
	"fmt"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	baremetalblockstorage1d2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/baremetal-blockstorage/1.2"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"strconv"
	"time"
)

type Client struct {
	Config       *scpsdk.Configuration
	sdkClient1d2 *baremetalblockstorage1d2.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:       config,
		sdkClient1d2: baremetalblockstorage1d2.NewAPIClient(config),
	}
}

func (client *Client) CreateBlockStorage(ctx context.Context, request VolumeResource) (*baremetalblockstorage1d2.AsyncResponse, error) {
	req := client.sdkClient1d2.BaremetalBlockstorageV1VolumeV1APIsAPI.CreateVolume(ctx)

	diskType, _ := baremetalblockstorage1d2.NewDiskTypeFromValue(request.DiskType.ValueString())

	attachments := client.getAttachmentListModelList(request.Attachments)

	tags := make([]baremetalblockstorage1d2.TagModel, 0)
	for k, v := range request.Tags.Elements() {
		tag := baremetalblockstorage1d2.TagModel{}

		key := baremetalblockstorage1d2.NullableString{}
		key.Set(&k)

		tag.Key = key

		if v != nil {
			value := baremetalblockstorage1d2.NullableString{}
			value.Set(v.(types.String).ValueStringPointer())
			tag.Value = value
		}
		tags = append(tags, tag)
	}

	requestBody := baremetalblockstorage1d2.VolumeCreateRequestV1Dot2{
		Name:        request.Name.ValueString(),
		DiskType:    *diskType,
		SizeGb:      request.SizeGb.ValueInt32(),
		Attachments: attachments,
		Tags:        tags,
	}

	if !request.QoS.IsNull() {
		attributes := request.QoS.Attributes()
		iops, _ := strconv.ParseInt(attributes["iops"].String(), 10, 32)
		throughput, _ := strconv.ParseInt(attributes["throughput"].String(), 10, 32)
		qos := baremetalblockstorage1d2.QoSModel{Iops: int32(iops), Throughput: int32(throughput)}
		requestBody.Qos = &qos
	}

	req = req.VolumeCreateRequestV1Dot2(requestBody)

	response, _, err := req.Execute()
	return response, err
}

func (client *Client) GetBlockStorage(ctx context.Context, blockStorageId string) (*baremetalblockstorage1d2.VolumeResponseV1Dot2, int, error) {
	req := client.sdkClient1d2.BaremetalBlockstorageV1VolumeV1APIsAPI.ShowVolume(ctx, blockStorageId)
	response, c, err := req.Execute()
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) AttachBlockStorages(ctx context.Context, blockStorageId string, attachObjectList []Attachment) (*baremetalblockstorage1d2.VolumeAttachmentResponse, int, error) {
	req := client.sdkClient1d2.BaremetalBlockstorageV1VolumeV1APIsAPI.CreateVolumeAttachments(ctx, blockStorageId)

	attachments := client.getAttachmentListModelList(attachObjectList)

	req = req.VolumeAttachmentRequest(baremetalblockstorage1d2.VolumeAttachmentRequest{Attachments: attachments})

	response, c, err := req.Execute()

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) DetachBlockStorages(ctx context.Context, blockStorageId string, detachObjectIdList []string) (*baremetalblockstorage1d2.VolumeAttachmentResponse, int, error) {
	req := client.sdkClient1d2.BaremetalBlockstorageV1VolumeV1APIsAPI.DeleteVolumeAttachments(ctx, blockStorageId)

	req = req.VolumeDetachRequest(baremetalblockstorage1d2.VolumeDetachRequest{
		Attachments: detachObjectIdList,
	})

	response, c, err := req.Execute()

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) UpdateBlockStorageQoS(ctx context.Context, blockStorageId string, iops int32, throughput int32) (*baremetalblockstorage1d2.SetVolumeQoSResponse, int, error) {
	req := client.sdkClient1d2.BaremetalBlockstorageV1VolumeV1APIsAPI.SetVolumeQos(ctx, blockStorageId)

	req = req.SetVolumeQoSRequest(baremetalblockstorage1d2.SetVolumeQoSRequest{
		Iops:       &iops,
		Throughput: &throughput,
	})

	response, c, err := req.Execute()
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) DeleteBlockStorage(ctx context.Context, blockStorageId string) (*baremetalblockstorage1d2.AsyncResponse, int, error) {
	req := client.sdkClient1d2.BaremetalBlockstorageV1VolumeV1APIsAPI.DeleteVolume(ctx, blockStorageId)

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

func (client *Client) getAttachmentListModelList(attachmentList []Attachment) []baremetalblockstorage1d2.AttachmentListModel {
	attachments := make([]baremetalblockstorage1d2.AttachmentListModel, 0)
	for _, attachment := range attachmentList {
		objectType, _ := baremetalblockstorage1d2.NewBlockStorageAttachmentObjectTypeFromValue(attachment.ObjectType.ValueString())
		attachments = append(attachments, baremetalblockstorage1d2.AttachmentListModel{
			ObjectType: objectType,
			ObjectId:   attachment.ObjectId.ValueStringPointer(),
		})
	}
	return attachments
}
