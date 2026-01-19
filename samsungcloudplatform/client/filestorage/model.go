package filestorage

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-filestorage" // 해당 서비스의 서비스 타입(keystone 에 등록된 service type)을 추가한다.

type VolumeDataSourceIds struct {
	Offset   types.Int32    `tfsdk:"offset"`
	Limit    types.Int32    `tfsdk:"limit"`
	Sort     types.String   `tfsdk:"sort"`
	Name     types.String   `tfsdk:"name"`
	TypeName types.String   `tfsdk:"type_name"`
	Ids      []types.String `tfsdk:"ids"`
}

type VolumeDataSource struct {
	AccountId               types.String `tfsdk:"account_id"`
	CreatedAt               types.String `tfsdk:"created_at"`
	EncryptionEnabled       types.Bool   `tfsdk:"encryption_enabled"`
	EndpointPath            types.String `tfsdk:"endpoint_path"`
	FileUnitRecoveryEnabled types.Bool   `tfsdk:"file_unit_recovery_enabled"`
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Path                    types.String `tfsdk:"path"`
	Protocol                types.String `tfsdk:"protocol"`
	Purpose                 types.String `tfsdk:"purpose"`
	State                   types.String `tfsdk:"state"`
	TypeId                  types.String `tfsdk:"type_id"`
	TypeName                types.String `tfsdk:"type_name"`
	Usage                   types.Int64  `tfsdk:"usage"`
}

type VolumeResource struct {
	AccountId               types.String         `tfsdk:"account_id"`
	CifsPassword            types.String         `tfsdk:"cifs_password"`
	CreatedAt               types.String         `tfsdk:"created_at"`
	EncryptionEnabled       types.Bool           `tfsdk:"encryption_enabled"`
	EndpointPath            types.String         `tfsdk:"endpoint_path"`
	FileUnitRecoveryEnabled types.Bool           `tfsdk:"file_unit_recovery_enabled"`
	Id                      types.String         `tfsdk:"id"`
	Name                    types.String         `tfsdk:"name"`
	NameUuid                types.String         `tfsdk:"name_uuid"`
	Path                    types.String         `tfsdk:"path"`
	Protocol                types.String         `tfsdk:"protocol"`
	Purpose                 types.String         `tfsdk:"purpose"`
	State                   types.String         `tfsdk:"state"`
	TypeId                  types.String         `tfsdk:"type_id"`
	TypeName                types.String         `tfsdk:"type_name"`
	Usage                   types.Int64          `tfsdk:"usage"`
	Tags                    types.Map            `tfsdk:"tags"`
	AccessRules             []AccessRuleResource `tfsdk:"access_rules"`
}

type AccessRuleResource struct {
	ObjectId   types.String `tfsdk:"object_id"`
	ObjectType types.String `tfsdk:"object_type"`
}

// -------------------- SnapshotSchedule -------------------- //

type SnapshotScheduleDataSource struct {
	VolumeId               types.String       `tfsdk:"volume_id"`
	SnapshotPolicyEnabled  types.Bool         `tfsdk:"snapshot_policy_enabled"`
	SnapshotRetentionCount types.Int32        `tfsdk:"snapshot_retention_count"`
	SnapshotSchedules      []SnapshotSchedule `tfsdk:"snapshot_schedules"`
}

type SnapshotScheduleResource struct {
	VolumeId               types.String     `tfsdk:"volume_id"`
	SnapshotPolicyEnabled  types.Bool       `tfsdk:"snapshot_policy_enabled"`
	SnapshotRetentionCount types.Int32      `tfsdk:"snapshot_retention_count"`
	SnapshotSchedule       SnapshotSchedule `tfsdk:"snapshot_schedule"`
}

type SnapshotSchedule struct {
	DayOfWeek types.String `tfsdk:"day_of_week"`
	Frequency types.String `tfsdk:"frequency"`
	Hour      types.String `tfsdk:"hour"`
	Id        types.String `tfsdk:"id"`
}

func (m SnapshotSchedule) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"day_of_week": types.StringType,
		"frequency":   types.StringType,
		"hour":        types.StringType,
		"id":          types.StringType,
	}
}

// -------------------- Replication -------------------- //
type ReplicationResources struct {
	VolumeId       types.String   `tfsdk:"volume_id"`
	ReplicationIds []types.String `tfsdk:"ids"`
}

type ReplicationResource struct {
	CifsPassword                 types.String `tfsdk:"cifs_password"`
	Name                         types.String `tfsdk:"name"`
	Region                       types.String `tfsdk:"region"`
	ReplicationFrequency         types.String `tfsdk:"replication_frequency"`
	VolumeId                     types.String `tfsdk:"volume_id"`
	ReplicationId                types.String `tfsdk:"replication_id"`
	ReplicationPolicy            types.String `tfsdk:"replication_policy"`
	ReplicationStatus            types.String `tfsdk:"replication_status"`
	ReplicationVolumeAccessLevel types.String `tfsdk:"replication_volume_access_level"`
	ReplicationVolumeId          types.String `tfsdk:"replication_volume_id"`
	ReplicationVolumeName        types.String `tfsdk:"replication_volume_name"`
	ReplicationVolumeRegion      types.String `tfsdk:"replication_volume_region"`
	SourceVolumeAccessLevel      types.String `tfsdk:"source_volume_access_level"`
	SourceVolumeId               types.String `tfsdk:"source_volume_id"`
	SourceVolumeName             types.String `tfsdk:"source_volume_name"`
	SourceVolumeRegion           types.String `tfsdk:"source_volume_region"`
	ReplicationUpdateType        types.String `tfsdk:"replication_update_type"`
	ReplicationType              types.String `tfsdk:"replication_type"`
	BackupRetentionCount         types.Int32  `tfsdk:"backup_retention_count"`
}

func (m ReplicationResource) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cifs_password":                   types.StringType,
		"name":                            types.StringType,
		"region":                          types.StringType,
		"replication_frequency":           types.StringType,
		"volume_id":                       types.StringType,
		"replication_id":                  types.StringType,
		"replication_policy":              types.StringType,
		"replication_status":              types.StringType,
		"replication_volume_access_level": types.StringType,
		"replication_volume_id":           types.StringType,
		"replication_volume_name":         types.StringType,
		"replication_volume_region":       types.StringType,
		"source_volume_access_level":      types.StringType,
		"source_volume_id":                types.StringType,
		"source_volume_name":              types.StringType,
		"source_volume_region":            types.StringType,
		"replication_update_type":         types.StringType,
	}
}

type VolumeReplicationPolicy struct {
	BackupRetentionCount  types.Int32  `json:"backup_retention_count"`
	ReplicationFrequency  types.String `tfsdk:"replication_frequency"`
	ReplicationPolicy     types.String `tfsdk:"replication_policy"`
	ReplicationUpdateType types.String `tfsdk:"replication_update_type"`
}

func (m VolumeReplicationPolicy) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"replication_frequency":   types.StringType,
		"replication_policy":      types.StringType,
		"replication_update_type": types.StringType,
	}
}

type ReplicationDataSource struct {
	Id          types.String `tfsdk:"id"`
	VolumeId    types.String `tfsdk:"volume_id"`
	Replication types.Object `tfsdk:"replication"`
}

type Replication struct {
	ReplicationFrequency         types.String `tfsdk:"replication_frequency"`
	ReplicationId                types.String `tfsdk:"replication_id"`
	ReplicationPolicy            types.String `tfsdk:"replication_policy"`
	ReplicationStatus            types.String `tfsdk:"replication_status"`
	ReplicationVolumeAccessLevel types.String `tfsdk:"replication_volume_access_level"`
	ReplicationVolumeId          types.String `tfsdk:"replication_volume_id"`
	ReplicationVolumeName        types.String `tfsdk:"replication_volume_name"`
	ReplicationVolumeRegion      types.String `tfsdk:"replication_volume_region"`
	SourceVolumeAccessLevel      types.String `tfsdk:"source_volume_access_level"`
	SourceVolumeId               types.String `tfsdk:"source_volume_id"`
	SourceVolumeName             types.String `tfsdk:"source_volume_name"`
	SourceVolumeRegion           types.String `tfsdk:"source_volume_region"`
}

func (m Replication) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"replication_frequency":           types.StringType,
		"replication_id":                  types.StringType,
		"replication_policy":              types.StringType,
		"replication_status":              types.StringType,
		"replication_volume_access_level": types.StringType,
		"replication_volume_id":           types.StringType,
		"replication_volume_name":         types.StringType,
		"replication_volume_region":       types.StringType,
		"source_volume_access_level":      types.StringType,
		"source_volume_id":                types.StringType,
		"source_volume_name":              types.StringType,
		"source_volume_region":            types.StringType,
	}
}
