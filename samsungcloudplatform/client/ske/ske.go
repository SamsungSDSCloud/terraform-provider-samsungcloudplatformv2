package ske

import (
	"context"
	"fmt"
	"io"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/ske/1.6"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpske.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpske.NewAPIClient(config),
	}
}

//------------ Cluster -------------------//

func (client *Client) GetClusterList(ctx context.Context, request ClusterDataSourceIds) (*scpske.ClusterListResponseV1Dot6, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.ListClusters(ctx)
	if request.Size != nil {
		req = req.Size(*request.Size)
	}
	if request.Page != nil {
		req = req.Page(*request.Page)
	}
	req = req.Sort(request.Sort.ValueString())
	req = req.Name(request.Name.ValueString())
	req = req.SubnetId(request.SubnetId.ValueString())
	if len(request.Status) > 0 {
		statuses := make([]*string, len(request.Status))
		for i, s := range request.Status {
			v := s.ValueString()
			statuses[i] = &v
		}
		req = req.Status(scpske.Status{ArrayOfPtrString: &statuses})
	}
	if len(request.KubernetesVersion) > 0 {
		versions := make([]string, len(request.KubernetesVersion))
		for i, v := range request.KubernetesVersion {
			versions[i] = v.ValueString()
		}
		req = req.KubernetesVersion(scpske.KubernetesVersion{ArrayOfString: &versions})
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.CreateCluster(ctx)

	securityGroupIdList := []string{}
	for _, securityGroupId := range stringListToSlice(request.SecurityGroupIdList) {
		securityGroupIdList = append(securityGroupIdList, securityGroupId)
	}

	additionalSubnetIdList := []string{}
	for _, additionalSubnetId := range stringListToSlice(request.AdditionalSubnetIdList) {
		additionalSubnetIdList = append(additionalSubnetIdList, additionalSubnetId)
	}

	tags := convertTags(request.Tags.Elements())

	req = req.ClusterCreateRequestV1Dot6(scpske.ClusterCreateRequestV1Dot6{
		Name:                                  request.Name.ValueString(),
		KubernetesVersion:                     request.KubernetesVersion.ValueString(),
		VpcId:                                 request.VpcId.ValueString(),
		DefaultSubnetId:                       request.DefaultSubnetId.ValueString(),
		AdditionalSubnetIdList:                additionalSubnetIdList,
		NfsVolumeId:                           *scpske.NewNullableString(request.NfsVolumeId.ValueStringPointer()),
		DeletionProtectionEnabled:             request.DeletionProtectionEnabled.ValueBoolPointer(),
		LinkedResources:                       convertLinkedResources(LinkedResourcesFromList(ctx, request.LinkedResources)),
		SecurityGroupIdList:                   securityGroupIdList,
		PrivateEndpointAccessControlResources: convertPrivateEndpointAccessControlResources(privateEndpointAccessControlResourcesFromList(ctx, request.PrivateEndpointAccessControlResources)),
		PublicEndpointAccessControlIp:         *scpske.NewNullableString(request.PublicEndpointAccessControlIp.ValueStringPointer()),
		ServiceWatchLoggingEnabled:            request.ServiceWatchLoggingEnabled.ValueBool(),
		Tags:                                  tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteCluster(ctx context.Context, clusterId string) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.DeleteCluster(ctx, clusterId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*scpske.ClusterShowResponseV1Dot6, int, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.ShowCluster(ctx, clusterId)

	resp, httpResponse, err := req.Execute()
	if httpResponse == nil {
		return nil, 0, err
	}
	return resp, httpResponse.StatusCode, err
}

func (client *Client) UpgradeCluster(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterSetResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterUpgrade(ctx, clusterId)

	req = req.ClusterUpgradeSetRequest(scpske.ClusterUpgradeSetRequest{
		KubernetesVersion: request.KubernetesVersion.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateClusterSecurityGroups(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterShowResponseV1Dot5, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterSecurityGroups(ctx, clusterId)

	securityGroupIdList := []string{}
	for _, securityGroupId := range stringListToSlice(request.SecurityGroupIdList) {
		securityGroupIdList = append(securityGroupIdList, securityGroupId)
	}

	req = req.ClusterSecurityGroupsSetRequest(scpske.ClusterSecurityGroupsSetRequest{
		SecurityGroupIdList: securityGroupIdList,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePrivateEndpointAccessControlResources(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterSetResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterPrivateAccessControl(ctx, clusterId)

	req = req.ClusterPrivateAccessControlSetRequest(scpske.ClusterPrivateAccessControlSetRequest{
		PrivateEndpointAccessControlResources: convertPrivateEndpointAccessControlResources(privateEndpointAccessControlResourcesFromList(ctx, request.PrivateEndpointAccessControlResources)),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePublicEndpointAccessControlIps(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterSetResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterPublicAccessControl(ctx, clusterId)

	req = req.ClusterPublicAccessControlSetRequest(scpske.ClusterPublicAccessControlSetRequest{
		PublicEndpointAccessControlIp: request.PublicEndpointAccessControlIp.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateServiceWatchLoggingEnabled(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterSetResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterServiceWatchLogging(ctx, clusterId)

	req = req.ClusterServiceWatchLoggingSetRequest(scpske.ClusterServiceWatchLoggingSetRequest{
		ServiceWatchLoggingEnabled: request.ServiceWatchLoggingEnabled.ValueBool(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetKubeConfig(ctx context.Context, clusterId string, kubeconfig_type string) (string, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.CreateClusterKubeconfig(ctx, clusterId)

	req = req.KubeconfigType(scpske.ClusterKubeconfigType(kubeconfig_type))

	resp, err := req.Execute()
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	byteArray, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(byteArray), err
}

func (client *Client) GetUserKubeConfig(ctx context.Context, clusterId string, kubeconfigType string) (string, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.ShowClusterUserKubeconfig(ctx, clusterId)
	req = req.KubeconfigType(scpske.ClusterKubeconfigType(kubeconfigType))

	resp, err := req.Execute()
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	byteArray, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(byteArray), err
}

//------------ Kubernetes Version -------------------//

func (client *Client) GetKubernetesVersionList(ctx context.Context) (*scpske.KubernetesVersionListResponse, error) {
	req := client.sdkClient.SkeV1KubernetesVersionsApiAPI.ListKubernetesVersions(ctx)
	resp, _, err := req.Execute()
	return resp, err
}

//------------ Nodepool -------------------//

func (client *Client) GetNodePoolList(ctx context.Context, request NodepoolDataSources) (*scpske.NodepoolListResponseV1Dot6, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.ListNodepools(ctx, request.ClusterId.ValueString())

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateNodepool(ctx context.Context, request NodepoolResource) (*scpske.NodepoolShowResponseV1Dot5, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.CreateNodepool(ctx)

	req = req.NodepoolCreateRequestV1Dot5(scpske.NodepoolCreateRequestV1Dot5{
		Name:                request.Name.ValueString(),
		ClusterId:           request.ClusterId.ValueString(),
		CustomImageId:       *scpske.NewNullableString(request.CustomImageId.ValueStringPointer()),
		DesiredNodeCount:    *scpske.NewNullableInt32(request.DesiredNodeCount.ValueInt32Pointer()),
		ImageOs:             request.ImageOs.ValueString(),
		ImageOsVersion:      request.ImageOsVersion.ValueString(),
		Labels:              convertLablels(request.Labels),
		Taints:              convertTaints(request.Taints),
		IsAutoRecovery:      request.IsAutoRecovery.ValueBool(),
		IsAutoScale:         request.IsAutoScale.ValueBool(),
		KeypairName:         request.KeypairName.ValueString(),
		KubernetesVersion:   request.KubernetesVersion.ValueString(),
		MaxNodeCount:        *scpske.NewNullableInt32(request.MaxNodeCount.ValueInt32Pointer()),
		MinNodeCount:        *scpske.NewNullableInt32(request.MinNodeCount.ValueInt32Pointer()),
		ServerTypeId:        request.ServerTypeId.ValueString(),
		VolumeTypeName:      request.VolumeTypeName.ValueString(),
		VolumeSize:          request.VolumeSize.ValueInt32(),
		ServerGroupId:       *scpske.NewNullableString(request.ServerGroupId.ValueStringPointer()),     // v1.1
		AdvancedSettings:    convertAdvancedSettings(request.AdvancedSettings),                         // v1.1
		LinkedResources:     convertLinkedResources(request.LinkedResources),                           // v1.3
		VolumeMaxIops:       *scpske.NewNullableInt32(request.VolumeMaxIops.ValueInt32Pointer()),       // v1.4
		VolumeMaxThroughput: *scpske.NewNullableInt32(request.VolumeMaxThroughput.ValueInt32Pointer()), // v1.4
		ScpGpuDriver:        *scpske.NewNullableString(request.ScpGpuDriver.ValueStringPointer()),      // v1.4
		PreferredIps:        *scpske.NewNullableString(request.PreferredIps.ValueStringPointer()),      // v1.5
		SubnetId:            *scpske.NewNullableString(request.SubnetId.ValueStringPointer()),          // v1.5
		Zone:                *scpske.NewNullableString(request.Zone.ValueStringPointer()),              // v1.5
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateNodepool(ctx context.Context, nodepoolId string, request NodepoolResource) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.SetNodepool(ctx, nodepoolId)

	req = req.NodepoolUpdateRequest(scpske.NodepoolUpdateRequest{
		DesiredNodeCount: *scpske.NewNullableInt32(request.DesiredNodeCount.ValueInt32Pointer()),
		IsAutoRecovery:   *scpske.NewNullableBool(request.IsAutoRecovery.ValueBoolPointer()),
		IsAutoScale:      *scpske.NewNullableBool(request.IsAutoScale.ValueBoolPointer()),
		MaxNodeCount:     *scpske.NewNullableInt32(request.MaxNodeCount.ValueInt32Pointer()),
		MinNodeCount:     *scpske.NewNullableInt32(request.MinNodeCount.ValueInt32Pointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateNodepoolLabels(ctx context.Context, nodepoolId string, request NodepoolResource) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.SetNodepoolLabels(ctx, nodepoolId)

	req = req.NodepoolLabelsSetRequest(scpske.NodepoolLabelsSetRequest{
		Labels: convertLablels(request.Labels),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateNodepoolTaints(ctx context.Context, nodepoolId string, request NodepoolResource) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.SetNodepoolTaints(ctx, nodepoolId)

	req = req.NodepoolTaintsSetRequest(scpske.NodepoolTaintsSetRequest{
		Taints: convertTaints(request.Taints),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateNodepoolLinkedResources(ctx context.Context, nodepoolId string, request NodepoolResource) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.SetNodepoolLinkedResources(ctx, nodepoolId)

	req = req.NodepoolLinkedResourcesSetRequest(scpske.NodepoolLinkedResourcesSetRequest{
		LinkedResources: convertLinkedResources(request.LinkedResources),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpgradeNodepool(ctx context.Context, request NodepoolResource) (*scpske.AsyncResponse, error) {
	reqV1Dot4 := scpske.NewNodepoolUpgradeSetRequestV1Dot4WithDefaults()
	reqV1Dot4.SetOsVersion(request.ImageOsVersion.ValueString())
	if request.ScpGpuDriver.ValueStringPointer() != nil {
		reqV1Dot4.SetScpGpuDriver(request.ScpGpuDriver.ValueString())
	}
	req := client.sdkClient.SkeV1NodepoolsApiAPI.SetNodepoolUpgrade(ctx, request.Id.ValueString())
	req = req.NodepoolUpgradeSetRequestV1Dot4(*reqV1Dot4)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteNodepool(ctx context.Context, nodepoolId string) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.DeleteNodepool(ctx, nodepoolId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetNodepool(ctx context.Context, nodepoolId string) (*scpske.NodepoolShowResponseV1Dot5, int, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.ShowNodepool(ctx, nodepoolId)

	resp, httpResponse, err := req.Execute()
	return resp, httpResponse.StatusCode, err
}

func (client *Client) CheckNodepoolList(ctx context.Context, clusterId string) (*scpske.NodepoolListResponseV1Dot6, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.ListNodepools(ctx, clusterId)

	resp, _, err := req.Execute()
	return resp, err
}

// ------------ Nodepoolnode-------------------//

func (client *Client) GetNodepoolNodeList(ctx context.Context, request NodepoolnodeDataSources) (*scpske.NodeListInNodepoolResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.ListNodepoolNodes(ctx, request.NodepoolId.ValueString())

	resp, _, err := req.Execute()
	return resp, err
}

//------------ Cluster Linked Resources (v1.6) -------------------//

func (client *Client) UpdateClusterLinkedResources(ctx context.Context, clusterId string, linkedResources []LinkedResource) (*scpske.ClusterSetResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterLinkedResources(ctx, clusterId)

	req = req.ClusterLinkedResourcesSetRequest(scpske.ClusterLinkedResourcesSetRequest{
		LinkedResources: convertLinkedResources(linkedResources),
	})

	resp, _, err := req.Execute()

	return resp, err
}

func convertPrivateEndpointAccessControlResources(privateEndpointAccessControlResources []PrivateEndpointAccessControlResource) []scpske.PrivateEndpointAccessControlResource {
	result := make([]scpske.PrivateEndpointAccessControlResource, len(privateEndpointAccessControlResources))
	for i, privateEndpointAccessControlResource := range privateEndpointAccessControlResources {
		result[i] = scpske.PrivateEndpointAccessControlResource{
			Id:   privateEndpointAccessControlResource.Id.ValueString(),
			Name: privateEndpointAccessControlResource.Name.ValueString(),
			Type: privateEndpointAccessControlResource.Type.ValueString(),
		}
	}
	return result
}

func convertLablels(lablels []Label) []scpske.NodepoolLabel {
	result := make([]scpske.NodepoolLabel, len(lablels))
	for i, label := range lablels {
		sKey := label.Key.ValueString()
		sValue := label.Value.ValueString()
		result[i] = scpske.NodepoolLabel{
			Key:   sKey,
			Value: &sValue,
		}
	}
	return result
}

func convertTaints(taints []Taint) []scpske.NodepoolTaint {
	result := make([]scpske.NodepoolTaint, len(taints))

	for i, taint := range taints {
		result[i] = scpske.NodepoolTaint{
			Effect: scpske.TaintEffectEnum(taint.Effect.ValueString()),
			Key:    taint.Key.ValueString(),
			Value:  taint.Value.ValueStringPointer(),
		}
	}
	return result
}

// List of TaintEffectEnum
func convertAdvancedSettings(advancedSettings *AdvancedSettings) scpske.NullableNodepoolAdvancedSettings {
	result := scpske.NullableNodepoolAdvancedSettings{}
	if advancedSettings != nil {
		value := scpske.NodepoolAdvancedSettings{
			AllowedUnsafeSysctls: advancedSettings.AllowedUnsafeSysctls.ValueStringPointer(),
			ContainerLogMaxFiles: advancedSettings.ContainerLogMaxFiles.ValueInt32(),
			ContainerLogMaxSize:  advancedSettings.ContainerLogMaxSize.ValueInt32(),
			ImageGcHighThreshold: advancedSettings.ImageGcHighThreshold.ValueInt32(),
			ImageGcLowThreshold:  advancedSettings.ImageGcLowThreshold.ValueInt32(),
			MaxPods:              advancedSettings.MaxPods.ValueInt32(),
			PodMaxPids:           advancedSettings.PodMaxPids.ValueInt32(),
		}
		result.Set(&value)
	} else {
		result.Unset()
	}
	return result
}

func convertTags(elements map[string]attr.Value) []scpske.Tag {
	var tags []scpske.Tag
	for k, v := range elements {
		tagObject := scpske.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}

func convertLinkedResources(linkedResources []LinkedResource) []scpske.LinkedResource {
	result := make([]scpske.LinkedResource, len(linkedResources))

	for i, linkedResource := range linkedResources {
		result[i] = scpske.LinkedResource{
			Id:   linkedResource.Id.ValueString(),
			Name: linkedResource.Name.ValueString(),
			Type: linkedResource.Type.ValueString(),
		}
	}
	return result
}

// LinkedResourcesFromList extracts []LinkedResource from a types.List (handles null/unknown).
func LinkedResourcesFromList(ctx context.Context, list types.List) []LinkedResource {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var result []LinkedResource
	_ = list.ElementsAs(ctx, &result, false)
	return result
}

// privateEndpointAccessControlResourcesFromList extracts []PrivateEndpointAccessControlResource from a types.List (handles null/unknown).
func privateEndpointAccessControlResourcesFromList(ctx context.Context, list types.List) []PrivateEndpointAccessControlResource {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var result []PrivateEndpointAccessControlResource
	_ = list.ElementsAs(ctx, &result, false)
	return result
}

// stringListToSlice extracts []string from a types.List of strings (handles null/unknown).
func stringListToSlice(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var result []string
	_ = list.ElementsAs(context.Background(), &result, false)
	return result
}

func (client *Client) SetClusterNfsVolume(ctx context.Context, clusterId string, nfsVolumeId string) (*scpske.ClusterSetResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterNfsVolume(ctx, clusterId)
	req = req.ClusterNfsVolumeSetRequest(scpske.ClusterNfsVolumeSetRequest{
		NfsVolumeId: *scpske.NewNullableString(&nfsVolumeId),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UnsetClusterNfsVolume(ctx context.Context, clusterId string) (*scpske.ClusterSetResponse, int, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterNfsVolume(ctx, clusterId)
	nullNfs := scpske.NewNullableString(nil)
	req = req.ClusterNfsVolumeSetRequest(scpske.ClusterNfsVolumeSetRequest{
		NfsVolumeId: *nullNfs,
	})

	resp, httpResp, err := req.Execute()
	httpStatus := 0
	if httpResp != nil {
		httpStatus = httpResp.StatusCode
	}
	return resp, httpStatus, err
}

//------------ Nodepool Image -------------------//

func (client *Client) GetNodepoolImageList(ctx context.Context, request NodepoolImageDataSources) (*scpske.NodepoolImageListResponseV1Dot6, error) {
	req := client.sdkClient.SkeV1ImagesAPIAPI.ListImages(ctx)
	req = req.ScpOriginalImageType(request.ScpOriginalImageType.ValueString())
	req = req.Size(request.Size.ValueInt32())
	req = req.Page(request.Page.ValueInt32())
	req = req.Sort(request.Sort.ValueString())
	req = req.KubernetesVersion(request.KubernetesVersion.ValueString())
	req = req.Os(request.Os.ValueString())

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) SetNodepoolPreferredIps(ctx context.Context, nodepoolId string, preferredIps string) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.SetNodepoolPreferredIps(ctx, nodepoolId)

	req = req.NodepoolPreferredIpsSetRequest(scpske.NodepoolPreferredIpsSetRequest{
		PreferredIps: preferredIps,
	})

	resp, _, err := req.Execute()
	return resp, err
}

//------------ Deletion Protection -------------------//

func (client *Client) GetClusterDeletionProtection(ctx context.Context, clusterId string) (bool, int, error) {
	data, httpStatus, err := client.GetCluster(ctx, clusterId)
	if err != nil {
		return false, httpStatus, fmt.Errorf("failed to get cluster %s: %w", clusterId, err)
	}

	return data.Cluster.DeletionProtectionEnabled, httpStatus, nil
}

//------------ V1.6 Cluster Update Methods -------------------//

func (client *Client) UpdateClusterSubnets(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterSetResponse, error) {
	additionalSubnetIdList := []string{}
	for _, id := range stringListToSlice(request.AdditionalSubnetIdList) {
		additionalSubnetIdList = append(additionalSubnetIdList, id)
	}

	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterSubnets(ctx, clusterId)
	req = req.ClusterSubnetsSetRequest(scpske.ClusterSubnetsSetRequest{
		AdditionalSubnetIdList: additionalSubnetIdList,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) SetClusterSubnets(ctx context.Context, clusterId string, additionalSubnetIdList []string) (*scpske.ClusterSetResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterSubnets(ctx, clusterId)
	req = req.ClusterSubnetsSetRequest(scpske.ClusterSubnetsSetRequest{
		AdditionalSubnetIdList: additionalSubnetIdList,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetClusterSubnets(ctx context.Context, clusterId string) ([]string, int, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.ShowCluster(ctx, clusterId)
	resp, httpResponse, err := req.Execute()
	if httpResponse == nil {
		return nil, 0, err
	}
	if resp == nil {
		return nil, httpResponse.StatusCode, err
	}
	return resp.Cluster.AdditionalSubnetIdList, httpResponse.StatusCode, err
}

func (client *Client) UpdateClusterNfsVolume(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterSetResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterNfsVolume(ctx, clusterId)
	req = req.ClusterNfsVolumeSetRequest(scpske.ClusterNfsVolumeSetRequest{
		NfsVolumeId: *scpske.NewNullableString(request.NfsVolumeId.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateClusterDeletionProtection(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterSetResponse, int, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterDeletionProtection(ctx, clusterId)
	req = req.ClusterDeletionProtectionSetRequest(scpske.ClusterDeletionProtectionSetRequest{
		DeletionProtectionEnabled: request.DeletionProtectionEnabled.ValueBool(),
	})

	resp, httpResponse, err := req.Execute()
	if httpResponse == nil {
		return nil, 0, err
	}
	return resp, httpResponse.StatusCode, err
}
