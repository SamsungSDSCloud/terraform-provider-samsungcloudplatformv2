package baremetalblockstorage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	baremetalblockstorage1d4 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/baremetal-blockstorage/1.4"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type Client struct {
	Config       *scpsdk.Configuration
	sdkClient1d4 *baremetalblockstorage1d4.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:       config,
		sdkClient1d4: baremetalblockstorage1d4.NewAPIClient(config),
	}
}

func (client *Client) CreateBlockStorage(ctx context.Context, request VolumeResource) (*baremetalblockstorage1d4.AsyncResponse, error) {
	req := client.sdkClient1d4.BaremetalBlockstorageV1VolumeV1APIsAPI.CreateVolume(ctx)

	diskType, _ := baremetalblockstorage1d4.NewDiskTypeFromValue(request.DiskType.ValueString())

	attachments := client.getAttachmentListModelList(request.Attachments)

	tags := make([]baremetalblockstorage1d4.TagModel, 0)
	for k, v := range request.Tags.Elements() {
		tag := baremetalblockstorage1d4.TagModel{}

		key := baremetalblockstorage1d4.NullableString{}
		key.Set(&k)

		tag.Key = key

		if v != nil {
			value := baremetalblockstorage1d4.NullableString{}
			value.Set(v.(types.String).ValueStringPointer())
			tag.Value = value
		}
		tags = append(tags, tag)
	}

	requestBody := baremetalblockstorage1d4.VolumeCreateRequestV1Dot4{
		Zone:        request.Zone.ValueString(),
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
		qos := baremetalblockstorage1d4.QoSModel{Iops: int32(iops), Throughput: int32(throughput)}
		requestBody.Qos = &qos
	}

	req = req.VolumeCreateRequestV1Dot4(requestBody)

	response, _, err := req.Execute()
	return response, err
}

func (client *Client) GetBlockStorage(ctx context.Context, blockStorageId string) (*baremetalblockstorage1d4.VolumeResponseV1Dot4, int, error) {
	req := client.sdkClient1d4.BaremetalBlockstorageV1VolumeV1APIsAPI.ShowVolume(ctx, blockStorageId)
	response, c, err := req.Execute()
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) AttachBlockStorages(ctx context.Context, blockStorageId string, attachObjectList []Attachment) (*baremetalblockstorage1d4.VolumeAttachmentResponse, int, error) {
	req := client.sdkClient1d4.BaremetalBlockstorageV1VolumeV1APIsAPI.CreateVolumeAttachments(ctx, blockStorageId)

	attachments := client.getAttachmentListModelList(attachObjectList)

	req = req.VolumeAttachmentRequest(baremetalblockstorage1d4.VolumeAttachmentRequest{Attachments: attachments})

	response, c, err := req.Execute()

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) DetachBlockStorages(ctx context.Context, blockStorageId string, detachObjectIdList []string) (*baremetalblockstorage1d4.VolumeAttachmentResponse, int, error) {
	req := client.sdkClient1d4.BaremetalBlockstorageV1VolumeV1APIsAPI.DeleteVolumeAttachments(ctx, blockStorageId)

	req = req.VolumeDetachRequest(baremetalblockstorage1d4.VolumeDetachRequest{
		Attachments: detachObjectIdList,
	})

	response, c, err := req.Execute()

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) UpdateBlockStorageQoS(ctx context.Context, blockStorageId string, iops int32, throughput int32) (*baremetalblockstorage1d4.SetVolumeQoSResponse, int, error) {
	req := client.sdkClient1d4.BaremetalBlockstorageV1VolumeV1APIsAPI.SetVolumeQos(ctx, blockStorageId)

	req = req.SetVolumeQoSRequest(baremetalblockstorage1d4.SetVolumeQoSRequest{
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

func (client *Client) DeleteBlockStorage(ctx context.Context, blockStorageId string) (*baremetalblockstorage1d4.AsyncResponse, int, error) {
	req := client.sdkClient1d4.BaremetalBlockstorageV1VolumeV1APIsAPI.DeleteVolume(ctx, blockStorageId)

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

func (client *Client) getAttachmentListModelList(attachmentList []Attachment) []baremetalblockstorage1d4.AttachmentListModel {
	attachments := make([]baremetalblockstorage1d4.AttachmentListModel, 0)
	for _, attachment := range attachmentList {
		objectType, _ := baremetalblockstorage1d4.NewBlockStorageAttachmentObjectTypeFromValue(attachment.ObjectType.ValueString())
		attachments = append(attachments, baremetalblockstorage1d4.AttachmentListModel{
			ObjectType: objectType,
			ObjectId:   attachment.ObjectId.ValueStringPointer(),
		})
	}
	return attachments
}
