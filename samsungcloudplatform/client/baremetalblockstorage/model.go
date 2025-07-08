package baremetalblockstorage

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-baremetal-blockstorage"

type VolumeResource struct {
	Id          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	DiskType    types.String   `tfsdk:"disk_type"`
	SizeGb      types.Int32    `tfsdk:"size_gb"`
	Attachments []Attachment   `tfsdk:"attachments"`
	Region      types.String   `tfsdk:"region"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
	Tags        types.Map      `tfsdk:"tags"`
}

type Attachment struct {
	ObjectId   types.String `tfsdk:"object_id"`
	ObjectType types.String `tfsdk:"object_type"`
}
