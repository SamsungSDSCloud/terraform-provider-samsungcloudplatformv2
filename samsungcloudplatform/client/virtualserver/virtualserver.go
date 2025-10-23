package virtualserver

import (
	"context"
	virtualservercommon "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	scpvirtualserver "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/virtualserver/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math"
	"net/http"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpvirtualserver.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpvirtualserver.NewAPIClient(config),
	}
}

// -------------------- Volume -------------------- //
func (client *Client) GetVolumeList() (*scpvirtualserver.VolumeListResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.VirtualserverV1VolumesApiAPI.ListVolumes(ctx)
	resp, _, err := req.Execute()
	return resp, err
}
func (client *Client) GetVolumeListWithParam(Name types.String, State types.String, Bootable types.Bool) (*scpvirtualserver.VolumeListResponse, error) {
	ctx := context.Background()
	req := client.sdkClient.VirtualserverV1VolumesApiAPI.ListVolumes(ctx)
	if !Name.IsNull() {
		req = req.Name(Name.ValueString())
	}
	if !State.IsNull() {
		req = req.State(State.ValueString())
	}
	if !Bootable.IsNull() {
		req = req.Bootable(Bootable.ValueBool())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateVolume(ctx context.Context, request VolumeResource) (*scpvirtualserver.VolumeShowResponse, error) {
	req := client.sdkClient.VirtualserverV1VolumesApiAPI.CreateVolume(ctx)

	//Tags
	var TagsObject []scpvirtualserver.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := scpvirtualserver.Tag{
			Key:   k,
			Value: *scpvirtualserver.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.VolumeCreateRequest(scpvirtualserver.VolumeCreateRequest{
		Name:       request.Name.ValueString(),
		Size:       request.Size.ValueInt32(),
		VolumeType: *scpvirtualserver.NewNullableString(request.VolumeType.ValueStringPointer()),
		Tags:       TagsObject,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVolume(ctx context.Context, volumeId string) (*scpvirtualserver.VolumeShowResponse, error) {
	req := client.sdkClient.VirtualserverV1VolumesApiAPI.ShowVolume(ctx, volumeId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetDefaultVolumeType(ctx context.Context) (*scpvirtualserver.VolumeTypeDetailResponse, error) {
	req := client.sdkClient.VirtualserverV1VolumeTypeApiAPI.ShowDefaultVolumeType(ctx)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVolume(ctx context.Context, volumeId string, request VolumeResource) (*scpvirtualserver.VolumeShowResponse, error) {
	req := client.sdkClient.VirtualserverV1VolumesApiAPI.UpdateVolume(ctx, volumeId)

	req = req.VolumeUpdateRequest(scpvirtualserver.VolumeUpdateRequest{
		Name: *scpvirtualserver.NewNullableString(request.Name.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) ExtendVolume(ctx context.Context, volumeId string, request VolumeResource) (*scpvirtualserver.VolumeShowResponse, error) {
	req := client.sdkClient.VirtualserverV1VolumesApiAPI.ExtendVolume(ctx, volumeId)

	req = req.VolumeExtendRequest(scpvirtualserver.VolumeExtendRequest{
		Size: request.Size.ValueInt32(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteVolume(ctx context.Context, volumeId string) error {
	req := client.sdkClient.VirtualserverV1VolumesApiAPI.DeleteVolume(ctx, volumeId)

	_, err := req.Execute()
	return err
}

func (client *Client) AttachVolume(ctx context.Context, volumeId string, serverId string) (*scpvirtualserver.VolumeServerAttachResponse, error) {
	req := client.sdkClient.VirtualserverV1VolumeServersApiAPI.AttachVolumeToVirtualServer(ctx, volumeId)

	req = req.VolumeServerAttachRequest(scpvirtualserver.VolumeServerAttachRequest{
		ServerId: serverId,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DetachVolume(ctx context.Context, volumeId string, serverId string) error {
	req := client.sdkClient.VirtualserverV1VolumeServersApiAPI.DetachVolumeFromVirtualServer(ctx, volumeId, serverId)

	_, err := req.Execute()
	return err
}

// -------------------- Keypair -------------------- //

func (client *Client) GetKeypairList() (*scpvirtualserver.KeypairListResponse, error) {
	ctx := context.Background()

	req := client.sdkClient.VirtualserverV1KeypairsAPI.ListKeypairs(ctx)
	req = req.Limit(math.MaxInt32)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateKeypair(ctx context.Context, request KeypairResource) (*scpvirtualserver.KeypairCreateResponse, error) {
	req := client.sdkClient.VirtualserverV1KeypairsAPI.CreateKeypair(ctx)

	if request.PublicKey.ValueString() == "" {
		request.PublicKey = types.StringNull()
	}

	//Tags
	var TagsObject []scpvirtualserver.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := scpvirtualserver.Tag{
			Key:   k,
			Value: *scpvirtualserver.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.KeypairCreateRequest(scpvirtualserver.KeypairCreateRequest{
		Name:      request.Name.ValueString(),
		PublicKey: *scpvirtualserver.NewNullableString(request.PublicKey.ValueStringPointer()),
		Tags:      TagsObject,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetKeypair(ctx context.Context, keypairName string) (*scpvirtualserver.KeypairShowResponse, error) {
	req := client.sdkClient.VirtualserverV1KeypairsAPI.ShowKeypair(ctx, keypairName)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteKeypair(ctx context.Context, keypairName string) error {
	req := client.sdkClient.VirtualserverV1KeypairsAPI.DeleteKeypair(ctx, keypairName)
	_, err := req.Execute()
	return err
}

// -------------------- Server -------------------- //

func (client *Client) GetServerList(Name types.String, Ip types.String, State types.String,
	ProductCategory types.String, ProductOffering types.String, VpcId types.String, ServerTypeId types.String,
	AutoScalingGroupId types.String) (*scpvirtualserver.ServerListResponse, error) {
	ctx := context.Background()

	req := client.sdkClient.VirtualserverV1ServersAPI.ListVirtualServers(ctx)
	req = req.Limit(math.MaxInt32)

	if !Name.IsNull() {
		req = req.Name(Name.ValueString())
	}
	if !Ip.IsNull() {
		req = req.Ip(Ip.ValueString())
	}
	if !State.IsNull() {
		req = req.State(State.ValueString())
	}
	if !ProductCategory.IsNull() {
		req = req.ProductCategory(scpvirtualserver.ServerProductCategory(ProductCategory.ValueString()))
	}
	if !ProductOffering.IsNull() {
		req = req.ProductOffering(scpvirtualserver.ProductOffering{
			ServerProductOffering:        scpvirtualserver.ServerProductOffering(ProductOffering.String()).Ptr(),
			ArrayOfServerProductOffering: nil,
		})
	}
	if !VpcId.IsNull() {
		req = req.VpcId(VpcId.ValueString())
	}
	if !ServerTypeId.IsNull() {
		req = req.ServerTypeId(ServerTypeId.ValueString())
	}
	if !AutoScalingGroupId.IsNull() {
		req = req.AutoScalingGroupId(AutoScalingGroupId.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetServerInterfaceList(ctx context.Context, serverId string) (*scpvirtualserver.ServerInterfaceListResponse, error) {
	req := client.sdkClient.VirtualserverV1ServersAPI.ListServerInterfaces(ctx, serverId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetServerSecurityGroupList(ctx context.Context, serverId string) (*scpvirtualserver.ServerSecurityGroupListResponse, error) {
	req := client.sdkClient.VirtualserverV1ServerActionAPI.ListVirtualServerSecurityGroups(ctx, serverId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetServerVolumeList(ctx context.Context, serverId string) (*scpvirtualserver.ServerVolumesResponse, error) {
	req := client.sdkClient.VirtualserverV1ServersAPI.ListServerVolumes(ctx, serverId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateServer(ctx context.Context, request ServerResource) (*scpvirtualserver.ServerCreateResponse, error) {
	req := client.sdkClient.VirtualserverV1ServersAPI.CreateVirtualServer(ctx)

	//MaxCount
	var defaultMaxCount int32 = 1
	var maxCount = *scpvirtualserver.NewNullableInt32(&defaultMaxCount)

	//Metadata
	metadataMapValue := make(map[string]interface{})
	for k, v := range request.Metadata.Elements() {
		metadataMapValue[k] = v.String()
	}

	//Networks
	var networks []scpvirtualserver.Network
	var networkMap map[string]ServerResourceNetwork
	request.Networks.ElementsAs(ctx, &networkMap, false)

	for _, network := range networkMap {
		networkState := &scpvirtualserver.Network{
			FixedIp:    *scpvirtualserver.NewNullableString(network.FixedIp.ValueStringPointer()),
			PortId:     *scpvirtualserver.NewNullableString(network.PortId.ValueStringPointer()),
			PublicIpId: *scpvirtualserver.NewNullableString(network.PublicIpId.ValueStringPointer()),
			SubnetId:   *scpvirtualserver.NewNullableString(network.SubnetId.ValueStringPointer()),
		}
		virtualservercommon.UnsetNilFields(networkState)
		networks = append(networks, *networkState)
	}

	//ProductCategory
	var productCategory scpvirtualserver.NullableServerProductCategory
	if !request.ProductCategory.IsNull() {
		serverProductCategory, _ := scpvirtualserver.NewServerProductCategoryFromValue(request.ProductCategory.ValueString())
		productCategory.Set(serverProductCategory)
	}

	//ProductOffering
	var productOffering scpvirtualserver.NullableServerProductOffering
	if !request.ProductOffering.IsNull() {
		serverProductOffering, _ := scpvirtualserver.NewServerProductOfferingFromValue(request.ProductOffering.ValueString())
		productOffering.Set(serverProductOffering)
	}

	//SecurityGroups
	var securityGroups []string
	for _, securityGroup := range request.SecurityGroups.Elements() {
		securityGroups = append(securityGroups, securityGroup.(types.String).ValueString())
	}

	//Tags
	var TagsObject []scpvirtualserver.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := scpvirtualserver.Tag{
			Key:   k,
			Value: *scpvirtualserver.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	//Volumes
	var volumes []scpvirtualserver.Volume
	imageSourceType, _ := scpvirtualserver.NewVolumeSourceTypeFromValue("image")
	blankSourceType, _ := scpvirtualserver.NewVolumeSourceTypeFromValue("blank")

	defaultVolumeType, _ := client.GetDefaultVolumeType(ctx)
	volumeType := request.BootVolume.Type.ValueStringPointer()
	if request.BootVolume.Type.IsNull() || request.BootVolume.Type.IsUnknown() {
		volumeType = defaultVolumeType.Name.Get()
	}
	deleteOnTermination := request.BootVolume.DeleteOnTermination.ValueBoolPointer()
	if request.BootVolume.DeleteOnTermination.IsNull() || request.BootVolume.DeleteOnTermination.IsUnknown() {
		deleteOnTermination = nil
	}

	volumeState := scpvirtualserver.Volume{
		BootIndex:           0,
		DeleteOnTermination: *scpvirtualserver.NewNullableBool(deleteOnTermination),
		Size:                request.BootVolume.Size.ValueInt32(),
		SourceType:          *scpvirtualserver.NewNullableVolumeSourceType(imageSourceType),
		Type:                *scpvirtualserver.NewNullableString(volumeType),
	}
	volumes = append(volumes, volumeState)

	var extraVolumeMap map[string]ServerResourceVolume
	request.ExtraVolumes.ElementsAs(ctx, &extraVolumeMap, false)
	for _, volume := range extraVolumeMap {
		volumeType := volume.Type.ValueStringPointer()
		if volume.Type.IsNull() || volume.Type.IsUnknown() {
			volumeType = defaultVolumeType.Name.Get()
		}
		deleteOnTermination := volume.DeleteOnTermination.ValueBoolPointer()
		if volume.DeleteOnTermination.IsNull() || volume.DeleteOnTermination.IsUnknown() {
			deleteOnTermination = nil
		}
		volumeState := scpvirtualserver.Volume{
			BootIndex:           -1,
			DeleteOnTermination: *scpvirtualserver.NewNullableBool(deleteOnTermination),
			Size:                volume.Size.ValueInt32(),
			SourceType:          *scpvirtualserver.NewNullableVolumeSourceType(blankSourceType),
			Type:                *scpvirtualserver.NewNullableString(volumeType),
		}
		volumes = append(volumes, volumeState)
	}

	reqState := &scpvirtualserver.ServerCreateRequest{
		ImageId:         request.ImageId.ValueString(),
		KeypairName:     request.KeypairName.ValueString(),
		Lock:            *scpvirtualserver.NewNullableBool(request.Lock.ValueBoolPointer()),
		MaxCount:        maxCount,
		Metadata:        metadataMapValue,
		Name:            request.Name.ValueString(),
		Networks:        networks,
		ProductCategory: productCategory,
		ProductOffering: productOffering,
		SecurityGroups:  securityGroups,
		ServerGroupId:   *scpvirtualserver.NewNullableString(request.ServerGroupId.ValueStringPointer()),
		ServerTypeId:    request.ServerTypeId.ValueString(),
		Tags:            TagsObject,
		UserData:        *scpvirtualserver.NewNullableString(request.UserData.ValueStringPointer()),
		Volumes:         volumes,
	}

	virtualservercommon.UnsetNilFields(reqState)

	req = req.ServerCreateRequest(*reqState)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateServerInterface(ctx context.Context, serverId string, request ServerResourceNetwork) (*scpvirtualserver.InterfaceResponse, error) {
	req := client.sdkClient.VirtualserverV1ServersAPI.CreateServerInterface(ctx, serverId)

	// Fixed IPs
	var fixedIps []scpvirtualserver.InterfaceAttachmentFixedIp
	if !request.FixedIp.IsNull() && !request.FixedIp.IsUnknown() {
		fixedIpState := scpvirtualserver.InterfaceAttachmentFixedIp{
			IpAddress: request.FixedIp.ValueString(),
		}
		fixedIps = append(fixedIps, fixedIpState)
	}

	var portId, subnetId *string
	if !request.PortId.IsNull() && !request.PortId.IsUnknown() {
		portId = request.PortId.ValueStringPointer()
	} else {
		portId = nil
	}

	if !request.SubnetId.IsNull() && !request.SubnetId.IsUnknown() {
		subnetId = request.SubnetId.ValueStringPointer()
	} else {
		subnetId = nil
	}

	reqState := &scpvirtualserver.ServerInterfaceCreateRequest{
		FixedIps: fixedIps,
		PortId:   *scpvirtualserver.NewNullableString(portId),
		SubnetId: *scpvirtualserver.NewNullableString(subnetId),
	}
	req = req.ServerInterfaceCreateRequest(*reqState)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateServerInterfaceNat(ctx context.Context, serverId string, portId string, publicIpId string) (*http.Response, error) {
	req := client.sdkClient.VirtualserverV1ServersAPI.CreateServerInterfaceNat(ctx, serverId, portId)

	reqState := &scpvirtualserver.ServerStaticNatCreateRequest{
		PublicipId: publicIpId,
	}
	req = req.ServerStaticNatCreateRequest(*reqState)
	resp, err := req.Execute()
	return resp, err
}

func (client *Client) GetServer(ctx context.Context, serverId string) (*scpvirtualserver.ServerShowResponse, error) {
	req := client.sdkClient.VirtualserverV1ServersAPI.ShowVirtualServer(ctx, serverId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetServerInterface(ctx context.Context, serverId string, portId string) (*scpvirtualserver.ServerInterfaceListResponse, error) {
	req := client.sdkClient.VirtualserverV1ServersAPI.ShowServerInterface(ctx, serverId, portId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateServer(ctx context.Context, request ServerResource) (*scpvirtualserver.ServerShowResponse, error) {
	req := client.sdkClient.VirtualserverV1ServersAPI.UpdateVirtualServer(ctx, request.Id.ValueString())

	reqState := &scpvirtualserver.ServerUpdateRequest{
		Name: request.Name.ValueString(),
	}

	req = req.ServerUpdateRequest(*reqState)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateServerLock(ctx context.Context, serverId string) error {
	req := client.sdkClient.VirtualserverV1ServerActionAPI.LockVirtualServer(ctx, serverId)
	_, err := req.Execute()
	return err
}

func (client *Client) UpdateServerUnlock(ctx context.Context, serverId string) error {
	req := client.sdkClient.VirtualserverV1ServerActionAPI.UnlockVirtualServer(ctx, serverId)
	_, err := req.Execute()
	return err
}

func (client *Client) UpdateServerSecurityGroupAttach(ctx context.Context, serverId string, securityGroupId string) error {
	req := client.sdkClient.VirtualserverV1ServerActionAPI.AttachVirtualServerSecurityGroup(ctx, serverId)
	reqState := &scpvirtualserver.ServerSecurityGroupActionRequestBody{
		SecurityGroupId: securityGroupId,
	}
	req = req.ServerSecurityGroupActionRequestBody(*reqState)
	_, err := req.Execute()
	return err
}

func (client *Client) UpdateServerSecurityGroupDetach(ctx context.Context, serverId string, securityGroupId string) error {
	req := client.sdkClient.VirtualserverV1ServerActionAPI.DetachVirtualServerSecurityGroup(ctx, serverId, securityGroupId)
	_, err := req.Execute()
	return err
}

func (client *Client) UpdateServerType(ctx context.Context, request ServerResource) error {
	req := client.sdkClient.VirtualserverV1ServerActionAPI.SetVirtualServerType(ctx, request.Id.ValueString())
	reqState := &scpvirtualserver.ServerSetServerTypeRequestBody{
		ServerType: request.ServerTypeId.ValueString(),
	}
	req = req.ServerSetServerTypeRequestBody(*reqState)
	_, err := req.Execute()
	return err
}

func (client *Client) UpdateServerVolume(ctx context.Context, serverId string, volumeId string, deleteOnTermination bool) error {
	req := client.sdkClient.VirtualserverV1ServersAPI.UpdateServerVolume(ctx, serverId, volumeId)
	reqState := &scpvirtualserver.ServerVolumesUpdateRequest{
		DeleteOnTermination: *scpvirtualserver.NewNullableBool(&deleteOnTermination),
		VolumeId:            volumeId,
	}
	req = req.ServerVolumesUpdateRequest(*reqState)
	_, err := req.Execute()
	return err
}

func (client *Client) StopServer(ctx context.Context, serverId string) error {
	req := client.sdkClient.VirtualserverV1ServerActionAPI.StopVirtualServer(ctx, serverId)
	_, err := req.Execute()
	return err
}

func (client *Client) StartServer(ctx context.Context, serverId string) error {
	req := client.sdkClient.VirtualserverV1ServerActionAPI.StartVirtualServer(ctx, serverId)
	_, err := req.Execute()
	return err
}

func (client *Client) DeleteServer(ctx context.Context, serverId string) error {
	req := client.sdkClient.VirtualserverV1ServersAPI.DeleteVirtualServer(ctx, serverId)
	_, err := req.Execute()
	return err
}

func (client *Client) DeleteServerInterface(ctx context.Context, serverId string, portId string) error {
	req := client.sdkClient.VirtualserverV1ServersAPI.DeleteServerInterface(ctx, serverId, portId)
	_, err := req.Execute()
	return err
}

func (client *Client) DeleteServerInterfaceNat(ctx context.Context, serverId string, portId string, natId string) error {
	req := client.sdkClient.VirtualserverV1ServersAPI.DeleteServerInterfaceNat(ctx, serverId, portId, natId)
	_, err := req.Execute()
	return err
}

// -------------------- Image -------------------- //

func (client *Client) GetImageList(ScpImageType types.String, ScpOriginalImageType types.String, Name types.String,
	OsDistro types.String, Status types.String, Visibility types.String) (*scpvirtualserver.ImageListResponse, error) {
	ctx := context.Background()

	req := client.sdkClient.VirtualserverV1VirtualserverV1ImagesAPI.ListImages(ctx)
	req = req.Limit(1000) // image api 최대 제약

	if !ScpImageType.IsNull() {
		req = req.ScpImageType(ScpImageType.ValueString())
	}
	if !ScpOriginalImageType.IsNull() {
		req = req.ScpOriginalImageType(ScpOriginalImageType.ValueString())
	}
	if !Name.IsNull() {
		req = req.Name(Name.ValueString())
	}
	if !OsDistro.IsNull() {
		req = req.OsDistro(OsDistro.ValueString())
	}
	if !Status.IsNull() {
		req = req.Status(Status.ValueString())
	}
	if !Visibility.IsNull() {
		req = req.Visibility(scpvirtualserver.ImageVisibilityEnum(Visibility.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateImage(ctx context.Context, request ImageResource) (*scpvirtualserver.ImageShowResponse, error) {
	req := client.sdkClient.VirtualserverV1VirtualserverV1ImagesAPI.CreateImage(ctx)

	//Tags
	var TagsObject []scpvirtualserver.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := scpvirtualserver.Tag{
			Key:   k,
			Value: *scpvirtualserver.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	reqState := &scpvirtualserver.ImageCreateRequest{
		ContainerFormat: request.ContainerFormat.ValueStringPointer(),
		DiskFormat:      request.DiskFormat.ValueStringPointer(),
		MinDisk:         request.MinDisk.ValueInt32Pointer(),
		MinRam:          request.MinRam.ValueInt32Pointer(),
		Name:            request.Name.ValueString(),
		OsDistro:        scpvirtualserver.ImageOsDistroEnum(request.OsDistro.ValueString()),
		Protected:       request.Protected.ValueBoolPointer(),
		Url:             request.Url.ValueString(),
		Visibility:      request.Visibility.ValueStringPointer(),
		Tags:            TagsObject,
	}

	req = req.ImageCreateRequest(*reqState)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateImageFromServer(ctx context.Context, request ImageResource) (*scpvirtualserver.ServerCreateImageResponse, error) {
	req := client.sdkClient.VirtualserverV1ServerActionAPI.CreateVirtualServerCustomImage(ctx, request.InstanceId.ValueString())

	reqState := &scpvirtualserver.ServerCreateImageRequestBody{
		ImageName: request.Name.ValueString(),
	}

	req = req.ServerCreateImageRequestBody(*reqState)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetImage(ctx context.Context, imageId string) (*scpvirtualserver.ImageShowResponse, error) {
	req := client.sdkClient.VirtualserverV1VirtualserverV1ImagesAPI.ShowImage(ctx, imageId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteImage(ctx context.Context, imageId string) error {
	req := client.sdkClient.VirtualserverV1VirtualserverV1ImagesAPI.DeleteImage(ctx, imageId)
	_, err := req.Execute()
	return err
}

func (client *Client) UpdateImage(ctx context.Context, imageId string, request ImageResource) (*scpvirtualserver.ImageShowResponse, error) {
	req := client.sdkClient.VirtualserverV1VirtualserverV1ImagesAPI.UpdateImage(ctx, imageId)

	reqState := &scpvirtualserver.ImageSetRequest{
		MinDisk:    *scpvirtualserver.NewNullableInt32(request.MinDisk.ValueInt32Pointer()),
		MinRam:     *scpvirtualserver.NewNullableInt32(request.MinRam.ValueInt32Pointer()),
		Protected:  *scpvirtualserver.NewNullableBool(request.Protected.ValueBoolPointer()),
		Visibility: *scpvirtualserver.NewNullableString(request.Visibility.ValueStringPointer()),
	}

	req = req.ImageSetRequest(*reqState)

	resp, _, err := req.Execute()
	return resp, err
}

// -------------------- Server Group -------------------- //

func (client *Client) GetServerGroupList() (*scpvirtualserver.ServerGroupListResponse, error) {
	ctx := context.Background()

	req := client.sdkClient.VirtualserverV1ServerGroupsAPI.ListServerGroups(ctx)
	req = req.Limit(math.MaxInt32)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateServerGroup(ctx context.Context, request ServerGroupResource) (*scpvirtualserver.ServerGroup, error) {
	req := client.sdkClient.VirtualserverV1ServerGroupsAPI.CreateServerGroup(ctx)

	//Tags
	var TagsObject []scpvirtualserver.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := scpvirtualserver.Tag{
			Key:   k,
			Value: *scpvirtualserver.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObject = append(TagsObject, tagObject)
	}

	req = req.ServerGroupCreateRequest(scpvirtualserver.ServerGroupCreateRequest{
		Name:   request.Name.ValueString(),
		Policy: request.Policy.ValueString(),
		Tags:   TagsObject,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetServerGroup(ctx context.Context, serverGroupId string) (*scpvirtualserver.ServerGroup, error) {
	req := client.sdkClient.VirtualserverV1ServerGroupsAPI.ShowServerGroup(ctx, serverGroupId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteServerGroup(ctx context.Context, serverGroupId string) error {
	req := client.sdkClient.VirtualserverV1ServerGroupsAPI.DeleteServerGroup(ctx, serverGroupId)
	_, err := req.Execute()
	return err
}
