package converter

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/ske"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/ske/1.4"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

func NodepoolResponseToNodepoolModel(data *scpske.NodepoolShowResponseV1Dot4) ske.Nodepool {
	nodepoolElement := data.GetNodepool()
	nodepoolModel := ske.Nodepool{
		Id:                  types.StringValue(nodepoolElement.Id),
		Name:                types.StringValue(nodepoolElement.Name),
		AccountId:           types.StringValue(nodepoolElement.AccountId),
		AutoRecoveryEnabled: types.BoolValue(nodepoolElement.AutoRecoveryEnabled),
		AutoScaleEnabled:    types.BoolValue(nodepoolElement.AutoScaleEnabled),
		Cluster: ske.IdMapType{
			Id: types.StringValue(nodepoolElement.Cluster.Id),
		},
		CurrentNodeCount: types.Int32Value(nodepoolElement.CurrentNodeCount),
		DesiredNodeCount: types.Int32Value(nodepoolElement.DesiredNodeCount),
		Image: ske.Image{
			CustomImageName: types.StringPointerValue(nodepoolElement.Image.CustomImageName.Get()),
			Os:              types.StringValue(nodepoolElement.Image.Os),
			OsVersion:       types.StringValue(nodepoolElement.Image.OsVersion),
			ScpGpuDriver:    types.StringPointerValue(nodepoolElement.Image.ScpGpuDriver.Get()),
		},
		Keypair: ske.NameMapType{
			Name: types.StringValue(nodepoolElement.Keypair.Name),
		},
		KubernetesVersion: types.StringValue(nodepoolElement.KubernetesVersion),
		MaxNodeCount:      types.Int32Value(nodepoolElement.MaxNodeCount),
		MinNodeCount:      types.Int32Value(nodepoolElement.MinNodeCount),
		ServerType: ske.ServerType{
			Description: types.StringValue(nodepoolElement.ServerType.Description),
			Id:          types.StringValue(nodepoolElement.ServerType.Id),
		},
		Status: types.StringValue(nodepoolElement.Status),
		VolumeType: ske.VolumeType{
			Encrypt: types.BoolValue(nodepoolElement.VolumeType.Encrypt),
			Id:      types.StringValue(nodepoolElement.VolumeType.Id),
			Name:    types.StringValue(nodepoolElement.VolumeType.Name),
		},
		VolumeSize:          types.Int32Value(nodepoolElement.VolumeSize),
		Labels:              MakeNodepoolLabelsModel(nodepoolElement.Labels),
		Taints:              MakeNodepoolTaintsModel(nodepoolElement.Taints),
		ServerGroupId:       types.StringPointerValue(nodepoolElement.ServerGroupId.Get()),
		AdvancedSettings:    MakeNodepoolAdvancedSettingsModel(nodepoolElement.AdvancedSettings),
		LinkedResources:     MakeLinkedResourcesModel(nodepoolElement.LinkedResources),
		VolumeMaxIops:       nodepoolElement.VolumeMaxIops.Get(),
		VolumeMaxThroughput: nodepoolElement.VolumeMaxThroughput.Get(),
		CreatedAt:           types.StringValue(nodepoolElement.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(nodepoolElement.CreatedBy),
		ModifiedAt:          types.StringValue(nodepoolElement.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(nodepoolElement.ModifiedBy),
	}
	return nodepoolModel
}

func MakeNodepoolLabelsModel(labels []scpske.NodepoolLabel) []ske.Label {
	var modelLabels []ske.Label
	for _, label := range labels {
		modelLabels = append(modelLabels, makeNodepoolLabelModel(&label))
	}
	return modelLabels
}

func MakeNodepoolTaintsModel(taints []scpske.NodepoolTaint) []ske.Taint {
	var modelTaints []ske.Taint
	for _, taint := range taints {
		modelTaints = append(modelTaints, makeNodepoolTaintModel(&taint))
	}
	return modelTaints
}

func MakeNodepoolAdvancedSettingsModel(advancedSettings scpske.NullableNodepoolAdvancedSettings) *ske.AdvancedSettings {
	value := advancedSettings.Get()
	if value != nil {
		return &ske.AdvancedSettings{
			AllowedUnsafeSysctls: types.StringPointerValue(value.AllowedUnsafeSysctls),
			ContainerLogMaxFiles: types.Int32Value(value.ContainerLogMaxFiles),
			ContainerLogMaxSize:  types.Int32Value(value.ContainerLogMaxSize),
			ImageGcHighThreshold: types.Int32Value(value.ImageGcHighThreshold),
			ImageGcLowThreshold:  types.Int32Value(value.ImageGcLowThreshold),
			MaxPods:              types.Int32Value(value.MaxPods),
			PodMaxPids:           types.Int32Value(value.PodMaxPids),
		}
	}
	return nil
}

func MakeLinkedResourcesModel(linkedResources []scpske.LinkedResource) []ske.LinkedResource {
	var modelLinkedResource []ske.LinkedResource
	for _, linkedResource := range linkedResources {
		modelLinkedResource = append(modelLinkedResource, makeLinkedResourceModel(&linkedResource))
	}
	return modelLinkedResource
}

func makeLinkedResourceModel(linkedResource *scpske.LinkedResource) ske.LinkedResource {
	return ske.LinkedResource{
		Id:   types.StringValue(linkedResource.Id),
		Name: types.StringValue(linkedResource.Name),
		Type: types.StringValue(linkedResource.Type),
	}
}

func makeNodepoolLabelModel(nodepoolLabels *scpske.NodepoolLabel) ske.Label {
	return ske.Label{
		Key:   types.StringValue(nodepoolLabels.GetKey()),
		Value: types.StringValue(nodepoolLabels.GetValue()),
	}
}

func makeNodepoolTaintModel(nodepoolTaints *scpske.NodepoolTaint) ske.Taint {
	return ske.Taint{
		Effect: types.StringValue(string(nodepoolTaints.GetEffect())),
		Key:    types.StringValue(nodepoolTaints.GetKey()),
		Value:  types.StringValue(nodepoolTaints.GetValue()),
	}
}
