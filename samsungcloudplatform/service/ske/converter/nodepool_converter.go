package converter

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/ske"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/ske/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
