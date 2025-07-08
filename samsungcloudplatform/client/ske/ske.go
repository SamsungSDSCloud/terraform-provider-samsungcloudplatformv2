package ske

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/ske/1.0"
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

	req = req.Size(request.Size.ValueInt32())
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

	req = req.ClusterCreateRequest(scpske.ClusterCreateRequest{
		Name:                                  request.Name.ValueString(),
		KubernetesVersion:                     request.KubernetesVersion.ValueString(),
		VpcId:                                 request.VpcId.ValueString(),
		SubnetId:                              request.SubnetId.ValueString(),
		VolumeId:                              request.VolumeId.ValueString(),
		CloudLoggingEnabled:                   request.CloudLoggingEnabled.ValueBool(),
		SecurityGroupIdList:                   securityGroupIdList,
		PrivateEndpointAccessControlResources: convertPrivateEndpointAccessControlResources(request.PrivateEndpointAccessControlResources),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteCluster(ctx context.Context, clusterId string) (*scpske.AsyncResponse, error) {
	req := client.sdkClient.SkeV1ClustersApiAPI.DeleteCluster(ctx, clusterId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetCluster(ctx context.Context, clusterId string) (*scpske.ClusterShowResponse, int, error) {
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

func (client *Client) CreateNodepool(ctx context.Context, request NodepoolResource) (*scpske.NodepoolShowResponse, error) {
	req := client.sdkClient.SkeV1NodepoolsApiAPI.CreateNodepool(ctx)

	req = req.NodepoolCreateRequest(scpske.NodepoolCreateRequest{
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

func (client *Client) GetNodepool(ctx context.Context, nodepoolId string) (*scpske.NodepoolShowResponse, int, error) {
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
		sId := privateEndpointAccessControlResource.Id.ValueString()
		sName := privateEndpointAccessControlResource.Name.ValueString()
		sType := privateEndpointAccessControlResource.Type.ValueString()
		result[i] = scpske.PrivateEndpointAccessControlResource{
			Id:   &sId,
			Name: &sName,
			Type: &sType,
		}
	}
	return result
}

func convertLablels(lablels []Labels) []scpske.NodepoolLabel {
	result := make([]scpske.NodepoolLabel, len(lablels))
	for i, lablel := range lablels {
		sKey := lablel.Key.ValueString()
		sValue := lablel.Value.ValueString()
		result[i] = scpske.NodepoolLabel{
			Key:   sKey,
			Value: &sValue,
		}
	}
	return result
}

func convertTaints(taints []Taints) []scpske.NodepoolTaint {
	result := make([]scpske.NodepoolTaint, len(taints))

	for i, taint := range taints {
		sEffect := scpske.TaintEffectEnum(taint.Effect.ValueString())
		sKey := taint.Key.ValueString()
		sValue := taint.Value.ValueString()

		result[i] = scpske.NodepoolTaint{
			Effect: &sEffect,
			Key:    sKey,
			Value:  &sValue,
		}
	}
	return result
}

// List of TaintEffectEnum
