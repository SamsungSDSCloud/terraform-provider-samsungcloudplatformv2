package ske

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/ske/1.1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
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

func (client *Client) GetClusterList(ctx context.Context, request ClusterDataSourceIds) (*scpske.ClusterListResponse, error) {
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
	//todo: will do later
	//req = req.Status(request.Status)
	//req = req.KubernetesVersion(request.KubernetesVersion.ValueString())

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateCluster(ctx context.Context, request ClusterResource) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.CreateCluster(ctx)

	var securityGroupIdList []string
	for _, securityGroupId := range request.SecurityGroupIdList {
		securityGroupIdList = append(securityGroupIdList, securityGroupId.ValueString())
	}

	tags := convertTags(request.Tags.Elements())

	req = req.ClusterCreateRequestV1Dot1(scpske.ClusterCreateRequestV1Dot1{
		Name:                                  request.Name.ValueString(),
		KubernetesVersion:                     request.KubernetesVersion.ValueString(),
		VpcId:                                 request.VpcId.ValueString(),
		SubnetId:                              request.SubnetId.ValueString(),
		VolumeId:                              request.VolumeId.ValueString(),
		CloudLoggingEnabled:                   request.CloudLoggingEnabled.ValueBool(),
		SecurityGroupIdList:                   securityGroupIdList,
		PrivateEndpointAccessControlResources: convertPrivateEndpointAccessControlResources(request.PrivateEndpointAccessControlResources),
		PublicEndpointAccessControlIp:         *scpske.NewNullableString(request.PublicEndpointAccessControlIp.ValueStringPointer()),
		ServiceWatchLoggingEnabled:            request.ServiceWatchLoggingEnabled.ValueBool(), // v1.1
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

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*scpske.ClusterShowResponseV1Dot1, int, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.ShowCluster(ctx, clusterId)

	resp, httpResponse, err := req.Execute()
	if httpResponse == nil {
		return nil, 0, err
	}
	return resp, httpResponse.StatusCode, err
}

func (client *Client) UpdateClusterLogging(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterSetResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterLogging(ctx, clusterId)

	req = req.ClusterLoggingSetRequest(scpske.ClusterLoggingSetRequest{
		CloudLoggingEnabled: request.CloudLoggingEnabled.ValueBool(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpgradeCluster(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterSetResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterUpgrade(ctx, clusterId)

	req = req.ClusterUpgradeSetRequest(scpske.ClusterUpgradeSetRequest{
		KubernetesVersion: request.KubernetesVersion.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateClusterSecurityGroups(ctx context.Context, clusterId string, request ClusterResource) (*scpske.ClusterShowResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.SetClusterSecurityGroups(ctx, clusterId)

	var securityGroupIdList []string
	for _, securityGroupId := range request.SecurityGroupIdList {
		securityGroupIdList = append(securityGroupIdList, securityGroupId.ValueString())
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
		PrivateEndpointAccessControlResources: convertPrivateEndpointAccessControlResources(request.PrivateEndpointAccessControlResources),
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

func (client *Client) GetNodePoolList(ctx context.Context, request NodepoolDataSources) (*scpske.NodepoolListResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.ListNodepools(ctx, request.ClusterId.ValueString())

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateNodepool(ctx context.Context, request NodepoolResource) (*scpske.NodepoolShowResponseV1Dot1, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.CreateNodepool(ctx)

	req = req.NodepoolCreateRequestV1Dot1(scpske.NodepoolCreateRequestV1Dot1{
		Name:              request.Name.ValueString(),
		ClusterId:         request.ClusterId.ValueString(),
		CustomImageId:     *scpske.NewNullableString(request.CustomImageId.ValueStringPointer()),
		DesiredNodeCount:  *scpske.NewNullableInt32(request.DesiredNodeCount.ValueInt32Pointer()),
		ImageOs:           request.ImageOs.ValueString(),
		ImageOsVersion:    request.ImageOsVersion.ValueString(),
		Labels:            convertLablels(request.Labels),
		Taints:            convertTaints(request.Taints),
		IsAutoRecovery:    request.IsAutoRecovery.ValueBool(),
		IsAutoScale:       request.IsAutoScale.ValueBool(),
		KeypairName:       request.KeypairName.ValueString(),
		KubernetesVersion: request.KubernetesVersion.ValueString(),
		MaxNodeCount:      *scpske.NewNullableInt32(request.MaxNodeCount.ValueInt32Pointer()),
		MinNodeCount:      *scpske.NewNullableInt32(request.MinNodeCount.ValueInt32Pointer()),
		ServerTypeId:      request.ServerTypeId.ValueString(),
		VolumeTypeName:    request.VolumeTypeName.ValueString(),
		VolumeSize:        request.VolumeSize.ValueInt32(),
		ServerGroupId:     *scpske.NewNullableString(request.ServerGroupId.ValueStringPointer()), // v1.1
		AdvancedSettings:  convertAdvancedSettings(request.AdvancedSettings),                     // v1.1
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

func (client *Client) UpgradeNodepool(ctx context.Context, nodepoolId string) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.SetNodepoolUpgrade(ctx, nodepoolId)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteNodepool(ctx context.Context, nodepoolId string) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.DeleteNodepool(ctx, nodepoolId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetNodepool(ctx context.Context, nodepoolId string) (*scpske.NodepoolShowResponseV1Dot1, int, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.ShowNodepool(ctx, nodepoolId)

	resp, httpResponse, err := req.Execute()
	return resp, httpResponse.StatusCode, err
}

func (client *Client) CheckNodepoolList(ctx context.Context, clusterId string) (*scpske.NodepoolListResponse, error) {
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
