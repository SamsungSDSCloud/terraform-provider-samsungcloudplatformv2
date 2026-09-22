package scr

import (
	"context"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/filter"
	scr11 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/scr/1.1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const reasonPrefix = "\nReason: "

// ---- ACL Resource ----

type AclResourceType struct {
	basetypes.ObjectType
}

type AclResource struct {
	ResourceId   types.String `tfsdk:"resource_id"`
	ResourceName types.String `tfsdk:"resource_name"`
	ResourceType types.String `tfsdk:"resource_type"`
	ResourceIps  types.List   `tfsdk:"resource_ips"`
	state        attr.ValueState
}

func (v AclResource) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"resource_id":   types.StringType,
		"resource_name": types.StringType,
		"resource_type": types.StringType,
		"resource_ips": types.ListType{
			ElemType: types.StringType,
		},
	}
}

// ---- Private ACL ----

type PrivateAclType struct {
	basetypes.ObjectType
}

type PrivateAcl struct {
	Enabled    types.Bool     `tfsdk:"enabled"`
	Resources  types.List     `tfsdk:"resources"`
	state      attr.ValueState
}

func (v PrivateAcl) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"resources": types.ListType{
			ElemType: AclResourceType{},
		},
	}
}

// ---- Public ACL ----

type PublicAclType struct {
	basetypes.ObjectType
}

type PublicAcl struct {
	Enabled    types.Bool     `tfsdk:"enabled"`
	Resources  types.List     `tfsdk:"resources"`
	state      attr.ValueState
}

func (v PublicAcl) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"resources": types.ListType{
			ElemType: AclResourceType{},
		},
	}
}

// ---- Registry Resource ----

type ContainerRegistryResource struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	State                   types.String `tfsdk:"state"`
	PrivateDomain           types.String `tfsdk:"private_domain"`
	PublicDomain            types.String `tfsdk:"public_domain"`
	BucketId                types.String `tfsdk:"bucket_id"`
	BucketName              types.String `tfsdk:"bucket_name"`
	BucketUsage             types.String `tfsdk:"bucket_usage"`
	PublicVisibleEnabled    types.Bool   `tfsdk:"public_visible_enabled"`
	PublicEndpointEnabled   types.Bool   `tfsdk:"public_endpoint_enabled"`
	PrivateAclEnabled       types.Bool   `tfsdk:"private_acl_enabled"`
	PrivateAclResources     types.List   `tfsdk:"private_acl_resources"`
	PublicAclEnabled        types.Bool   `tfsdk:"public_acl_enabled"`
	PublicAclResources      types.List   `tfsdk:"public_acl_resources"`
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	ModifiedAt              types.String `tfsdk:"modified_at"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
}

// ---- Registry Data Source ----

type ContainerRegistryDataSource struct {
	Id                      types.String    `tfsdk:"id"`
	Name                    types.String    `tfsdk:"name"`
	State                   types.String    `tfsdk:"state"`
	PrivateDomain           types.String    `tfsdk:"private_domain"`
	PublicDomain            types.String    `tfsdk:"public_domain"`
	BucketId                types.String    `tfsdk:"bucket_id"`
	BucketName              types.String    `tfsdk:"bucket_name"`
	BucketUsage             types.String    `tfsdk:"bucket_usage"`
	PublicVisibleEnabled    types.Bool      `tfsdk:"public_visible_enabled"`
	PublicEndpointEnabled   types.Bool      `tfsdk:"public_endpoint_enabled"`
	PrivateAclEnabled       types.Bool      `tfsdk:"private_acl_enabled"`
	PrivateAclResources     types.List      `tfsdk:"private_acl_resources"`
	PublicAclEnabled        types.Bool      `tfsdk:"public_acl_enabled"`
	PublicAclResources      types.List      `tfsdk:"public_acl_resources"`
	CreatedAt               types.String    `tfsdk:"created_at"`
	CreatedBy               types.String    `tfsdk:"created_by"`
	ModifiedAt              types.String    `tfsdk:"modified_at"`
	ModifiedBy              types.String    `tfsdk:"modified_by"`
	Filter                  []filter.Filter `tfsdk:"filter"`
}

// ---- Registries Data Source ----

type ContainerRegistryDataSources struct {
	Name        types.String    `tfsdk:"name"`
	Ids         []types.String  `tfsdk:"ids"`
	Filter      []filter.Filter `tfsdk:"filter"`
}

// ---- Image Data Source ----

type ImageDataSource struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	State              types.String `tfsdk:"state"`
	RegistryId         types.String `tfsdk:"registry_id"`
	RepositoryId       types.String `tfsdk:"repository_id"`
	PrivateEndpointUrl types.String `tfsdk:"private_endpoint_url"`
	PublicEndpointUrl  types.String `tfsdk:"public_endpoint_url"`
	PullCount          types.Int32  `tfsdk:"pull_count"`
	CreatedAt          types.String `tfsdk:"created_at"`
	ModifiedAt         types.String `tfsdk:"modified_at"`
}

// ---- Images Data Source ----

type ImageDataSources struct {
	RepositoryId types.String    `tfsdk:"repository_id"`
	Name         types.String    `tfsdk:"name"`
	Sort         types.String    `tfsdk:"sort"`
	Page         types.Int32     `tfsdk:"page"`
	Size         types.Int32     `tfsdk:"size"`
	Ids          []types.String  `tfsdk:"ids"`
	Filter       []filter.Filter `tfsdk:"filter"`
}

// ---- Lifecycle Policy ----

type LifecyclePolicyType struct {
	basetypes.ObjectType
}

type LifecyclePolicy struct {
	LifecyclePolicyEnabled  types.Bool   `tfsdk:"lifecycle_policy_enabled"`
	OutdatedRuleDuration    types.Int32  `tfsdk:"outdated_rule_duration"`
	OutdatedRuleEnabled     types.Bool   `tfsdk:"outdated_rule_enabled"`
	OutdatedRuleTagExpression types.String `tfsdk:"outdated_rule_tag_expression"`
	UntaggedRuleDuration    types.Int32  `tfsdk:"untagged_rule_duration"`
	UntaggedRuleEnabled     types.Bool   `tfsdk:"untagged_rule_enabled"`
	state                   attr.ValueState
}

func (v LifecyclePolicy) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"lifecycle_policy_enabled":   types.BoolType,
		"outdated_rule_duration":     types.Int32Type,
		"outdated_rule_enabled":      types.BoolType,
		"outdated_rule_tag_expression": types.StringType,
		"untagged_rule_duration":     types.Int32Type,
		"untagged_rule_enabled":      types.BoolType,
	}
}

// ---- Lock Policy ----

type LockPolicyType struct {
	basetypes.ObjectType
}

type LockPolicyModel struct {
	Locked types.Bool `tfsdk:"locked"`
	state  attr.ValueState
}

func (v LockPolicyModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"locked": types.BoolType,
	}
}

var LockPolicyModelType = basetypes.ObjectType{AttrTypes: LockPolicyModel{}.AttributeTypes()}

// ---- Pull Policy ----

type PullPolicyType struct {
	basetypes.ObjectType
}

type PullPolicyModel struct {
	CriticalLimit             types.Int32  `tfsdk:"critical_limit"`
	HighLimit                 types.Int32  `tfsdk:"high_limit"`
	UnmodifiedExcepted        types.Bool   `tfsdk:"unmodified_excepted"`
	UnscannedImagePullPrevented types.Bool `tfsdk:"unscanned_image_pull_prevented"`
	VulnerableImagePullPrevented types.Bool `tfsdk:"vulnerable_image_pull_prevented"`
	state                     attr.ValueState
}

func (v PullPolicyModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"critical_limit":                 types.Int32Type,
		"high_limit":                     types.Int32Type,
		"unmodified_excepted":            types.BoolType,
		"unscanned_image_pull_prevented": types.BoolType,
		"vulnerable_image_pull_prevented": types.BoolType,
	}
}

// ---- Scan Policy ----

type ScanPolicyType struct {
	basetypes.ObjectType
}

type ScanPolicyModel struct {
	AutoScanEnabled      types.Bool   `tfsdk:"auto_scan_enabled"`
	FixedVersionExcepted types.Bool   `tfsdk:"fixed_version_excepted"`
	LanguageExcepted     types.Bool   `tfsdk:"language_excepted"`
	ScanPolicyEnabled    types.Bool   `tfsdk:"scan_policy_enabled"`
	SecretExcepted       types.Bool   `tfsdk:"secret_excepted"`
	SeverityLimit        types.String `tfsdk:"severity_limit"`
	state               attr.ValueState
}

func (v ScanPolicyModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"auto_scan_enabled":      types.BoolType,
		"fixed_version_excepted": types.BoolType,
		"language_excepted":      types.BoolType,
		"scan_policy_enabled":    types.BoolType,
		"secret_excepted":        types.BoolType,
		"severity_limit":         types.StringType,
	}
}

// ---- Repository Resource ----

type RepositoryResource struct {
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	RegistryId            types.String `tfsdk:"registry_id"`
	State                 types.String `tfsdk:"state"`
	PrivateEndpointUrl    types.String `tfsdk:"private_endpoint_url"`
	PublicEndpointUrl     types.String `tfsdk:"public_endpoint_url"`
	PullPolicy            types.List   `tfsdk:"pull_policy"`
	ScanPolicy            types.List   `tfsdk:"scan_policy"`
	LifecyclePolicy       types.List   `tfsdk:"lifecycle_policy"`
	LockPolicy            types.List   `tfsdk:"lock_policy"`
	Tags                  types.Map    `tfsdk:"tags"`
	CreatedAt             types.String `tfsdk:"created_at"`
	CreatedBy             types.String `tfsdk:"created_by"`
	ModifiedAt            types.String `tfsdk:"modified_at"`
	ModifiedBy            types.String `tfsdk:"modified_by"`
}

// ---- Repository Data Source ----

type RepositoryDataSource struct {
	Id                    types.String    `tfsdk:"id"`
	Name                  types.String    `tfsdk:"name"`
	Description           types.String    `tfsdk:"description"`
	RegistryId            types.String    `tfsdk:"registry_id"`
	State                 types.String    `tfsdk:"state"`
	PrivateEndpointUrl    types.String    `tfsdk:"private_endpoint_url"`
	PublicEndpointUrl     types.String    `tfsdk:"public_endpoint_url"`
	PullPolicy            types.List      `tfsdk:"pull_policy"`
	ScanPolicy            types.List      `tfsdk:"scan_policy"`
	LifecyclePolicy       types.List      `tfsdk:"lifecycle_policy"`
	LockPolicy            types.List      `tfsdk:"lock_policy"`
	Tags                  types.Map       `tfsdk:"tags"`
	CreatedAt             types.String    `tfsdk:"created_at"`
	CreatedBy             types.String    `tfsdk:"created_by"`
	ModifiedAt            types.String    `tfsdk:"modified_at"`
	ModifiedBy            types.String    `tfsdk:"modified_by"`
	Filter                []filter.Filter `tfsdk:"filter"`
}

// ---- Repositories Data Source ----

type RepositoryDataSources struct {
	RegistryId types.String    `tfsdk:"registry_id"`
	Name       types.String    `tfsdk:"name"`
	Sort       types.String    `tfsdk:"sort"`
	Page       types.Int32     `tfsdk:"page"`
	Size       types.Int32     `tfsdk:"size"`
	Ids        []types.String  `tfsdk:"ids"`
	Filter     []filter.Filter `tfsdk:"filter"`
}

// ---- Image Resource ----

type ImageResourceModel struct {
	Id                types.String   `tfsdk:"id"`
	Name              types.String   `tfsdk:"name"`
	Description       types.String   `tfsdk:"description"`
	State             types.String   `tfsdk:"state"`
	RegistryId        types.String   `tfsdk:"registry_id"`
	RepositoryId      types.String   `tfsdk:"repository_id"`
	PullCount         types.Int32    `tfsdk:"pull_count"`
	PullPolicy        types.List `tfsdk:"pull_policy"`
	ScanPolicy        types.List `tfsdk:"scan_policy"`
	LifecyclePolicy   types.List `tfsdk:"lifecycle_policy"`
	LockPolicy        types.List `tfsdk:"lock_policy"`
	CreatedAt         types.String   `tfsdk:"created_at"`
	CreatedBy         types.String   `tfsdk:"created_by"`
	ModifiedAt        types.String   `tfsdk:"modified_at"`
	ModifiedBy        types.String   `tfsdk:"modified_by"`
}

// ---- Tag Data Source ----

type TagDataSource struct {
	Id                types.String    `tfsdk:"id"`
	ImageId           types.String    `tfsdk:"image_id"`
	RegistryId        types.String    `tfsdk:"registry_id"`
	RepositoryId      types.String    `tfsdk:"repository_id"`
	HashDigest        types.String    `tfsdk:"hash_digest"`
	Manifest          types.String    `tfsdk:"manifest"`
	ManifestMediaType types.String    `tfsdk:"manifest_media_type"`
	State             types.String    `tfsdk:"state"`
	ReferenceTags     types.List      `tfsdk:"reference_tags"`
	LockPolicy        types.List      `tfsdk:"lock_policy"`
	CreatedAt         types.String    `tfsdk:"created_at"`
	ModifiedAt        types.String    `tfsdk:"modified_at"`
	Filter            []filter.Filter `tfsdk:"filter"`
}

// ---- Tags Data Source ----

type TagsDataSource struct {
	ImageId       types.String    `tfsdk:"image_id"`
	Sort          types.String    `tfsdk:"sort"`
	Page          types.Int32     `tfsdk:"page"`
	Size          types.Int32     `tfsdk:"size"`
	ReferenceTags types.String    `tfsdk:"reference_tags"`
	TotalCount    types.Int32     `tfsdk:"total_count"`
	Tags          types.List      `tfsdk:"tags"`
	Filter        []filter.Filter `tfsdk:"filter"`
}

type TagListModel struct {
	Id                 types.String `tfsdk:"id"`
	HashDigest         types.String `tfsdk:"hash_digest"`
	ReferenceTags      types.List   `tfsdk:"reference_tags"`
	LockPolicy         types.List   `tfsdk:"lock_policy"`
	LastScannedAt      types.String `tfsdk:"last_scanned_at"`
	ModifiedAt         types.String `tfsdk:"modified_at"`
	PrivateEndpointUrl types.String `tfsdk:"private_endpoint_url"`
	PublicEndpointUrl  types.String `tfsdk:"public_endpoint_url"`
	ReScanNeeded       types.Bool   `tfsdk:"re_scan_needed"`
	ReferencedBy       types.String `tfsdk:"referenced_by"`
	ScanState          types.String `tfsdk:"scan_state"`
	ScanSummary        types.Object `tfsdk:"scan_summary"`
	Size               types.Int32  `tfsdk:"size"`
	State              types.String `tfsdk:"state"`
}

func (v TagListModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                   types.StringType,
		"hash_digest":          types.StringType,
		"reference_tags":       types.ListType{ElemType: types.StringType},
		"lock_policy":          types.ListType{ElemType: LockPolicyModelType},
		"last_scanned_at":      types.StringType,
		"modified_at":          types.StringType,
		"private_endpoint_url": types.StringType,
		"public_endpoint_url":  types.StringType,
		"re_scan_needed":       types.BoolType,
		"referenced_by":        types.StringType,
		"scan_state":           types.StringType,
		"scan_summary":         ScanSummaryModelType,
		"size":                 types.Int32Type,
		"state":                types.StringType,
	}
}

var TagListModelType = basetypes.ObjectType{AttrTypes: TagListModel{}.AttributeTypes()}

// ---- Tag Packages Data Source ----

type TagPackagesDataSource struct {
	TagsId         types.String    `tfsdk:"tags_id"`
	OsLanguage     types.String    `tfsdk:"os_language"`
	PackageName    types.String    `tfsdk:"package_name"`
	Sort           types.String    `tfsdk:"sort"`
	Page           types.Int32     `tfsdk:"page"`
	Size           types.Int32     `tfsdk:"size"`
	FilteredCount  types.Int32     `tfsdk:"filtered_count"`
	LastScannedAt  types.String    `tfsdk:"last_scanned_at"`
	ReleaseVersion types.String    `tfsdk:"release_version"`
	Packages       types.List      `tfsdk:"packages"`
	ScanSummary    types.Object    `tfsdk:"scan_summary"`
	Filter         []filter.Filter `tfsdk:"filter"`
}

type PackageReportModel struct {
	Category     types.String `tfsdk:"category"`
	OsLanguage   types.String `tfsdk:"os_language"`
	PackageName  types.String `tfsdk:"package_name"`
	ScanSummary  types.Object `tfsdk:"scan_summary"`
	Type         types.String `tfsdk:"type"`
	Version      types.String `tfsdk:"version"`
}

func (v PackageReportModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"category":     types.StringType,
		"os_language":   types.StringType,
		"package_name":  types.StringType,
		"scan_summary":  ScanSummaryModelType,
		"type":          types.StringType,
		"version":       types.StringType,
	}
}

var PackageReportModelType = basetypes.ObjectType{AttrTypes: PackageReportModel{}.AttributeTypes()}

// ---- Tag Secrets Data Source ----

type TagSecretsDataSource struct {
	TagsId         types.String    `tfsdk:"tags_id"`
	Sort           types.String    `tfsdk:"sort"`
	Page           types.Int32     `tfsdk:"page"`
	Size           types.Int32     `tfsdk:"size"`
	FileName       types.String    `tfsdk:"file_name"`
	FilteredCount  types.Int32     `tfsdk:"filtered_count"`
	LastScannedAt  types.String    `tfsdk:"last_scanned_at"`
	ReleaseVersion types.String    `tfsdk:"release_version"`
	Secrets        types.List      `tfsdk:"secrets"`
	SecretSummary  types.Object    `tfsdk:"secret_summary"`
	Filter         []filter.Filter `tfsdk:"filter"`
}

type SecretReportModel struct {
	Category           types.String `tfsdk:"category"`
	FileName           types.String `tfsdk:"file_name"`
	Match              types.String `tfsdk:"match"`
	RuleId             types.String `tfsdk:"rule_id"`
	Severity           types.String `tfsdk:"severity"`
	StartLine          types.Int32  `tfsdk:"start_line"`
	Target             types.String `tfsdk:"target"`
	Title              types.String `tfsdk:"title"`
	VulnerabilityClass types.String `tfsdk:"vulnerability_class"`
}

func (v SecretReportModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"category":            types.StringType,
		"file_name":           types.StringType,
		"match":               types.StringType,
		"rule_id":             types.StringType,
		"severity":            types.StringType,
		"start_line":          types.Int32Type,
		"target":              types.StringType,
		"title":               types.StringType,
		"vulnerability_class": types.StringType,
	}
}

var SecretReportModelType = basetypes.ObjectType{AttrTypes: SecretReportModel{}.AttributeTypes()}

type SecretSummaryModel struct {
	Critical    types.Int32 `tfsdk:"critical"`
	High        types.Int32 `tfsdk:"high"`
	Low         types.Int32 `tfsdk:"low"`
	Medium      types.Int32 `tfsdk:"medium"`
	TotalSecret types.Int32 `tfsdk:"total_secret"`
	Unknown     types.Int32 `tfsdk:"unknown"`
}

func (v SecretSummaryModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"critical":     types.Int32Type,
		"high":         types.Int32Type,
		"low":          types.Int32Type,
		"medium":       types.Int32Type,
		"total_secret": types.Int32Type,
		"unknown":      types.Int32Type,
	}
}

var SecretSummaryModelType = basetypes.ObjectType{AttrTypes: SecretSummaryModel{}.AttributeTypes()}

// ---- Tag Vulnerabilities Data Source ----

type TagVulnerabilitiesDataSource struct {
	TagsId                 types.String    `tfsdk:"tags_id"`
	OsLanguage             types.String    `tfsdk:"os_language"`
	PackageName            types.String    `tfsdk:"package_name"`
	Sort                   types.String    `tfsdk:"sort"`
	Page                   types.Int32     `tfsdk:"page"`
	Size                   types.Int32     `tfsdk:"size"`
	UpdateVersionAvailable types.Bool      `tfsdk:"update_version_available"`
	Severity               types.String    `tfsdk:"severity"`
	Category               types.String    `tfsdk:"category"`
	FilteredCount          types.Int32     `tfsdk:"filtered_count"`
	LastScannedAt          types.String    `tfsdk:"last_scanned_at"`
	ReleaseVersion         types.String    `tfsdk:"release_version"`
	Vulnerabilities        types.List      `tfsdk:"vulnerabilities"`
	ScanSummary            types.Object    `tfsdk:"scan_summary"`
	Filter                 []filter.Filter `tfsdk:"filter"`
}

type CveModel struct {
	Category       types.String `tfsdk:"category"`
	CurrentVersion types.String `tfsdk:"current_version"`
	CveCode        types.String `tfsdk:"cve_code"`
	Description    types.String `tfsdk:"description"`
	Links          types.String `tfsdk:"links"`
	OsLanguage     types.String `tfsdk:"os_language"`
	PackageName    types.String `tfsdk:"package_name"`
	Severity       types.String `tfsdk:"severity"`
	UpdateVersion  types.String `tfsdk:"update_version"`
	Vectors        types.Object `tfsdk:"vectors"`
}

func (v CveModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"category":        types.StringType,
		"current_version": types.StringType,
		"cve_code":        types.StringType,
		"description":     types.StringType,
		"links":           types.StringType,
		"os_language":     types.StringType,
		"package_name":    types.StringType,
		"severity":        types.StringType,
		"update_version":  types.StringType,
		"vectors":         VectorsModelType,
	}
}

var CveModelType = basetypes.ObjectType{AttrTypes: CveModel{}.AttributeTypes()}

type VectorsModel struct {
	AttackComplexity   types.String  `tfsdk:"attack_complexity"`
	AttackVector       types.String  `tfsdk:"attack_vector"`
	Availability       types.String  `tfsdk:"availability"`
	BaseSeverity       types.String  `tfsdk:"base_severity"`
	Confidentiality    types.String  `tfsdk:"confidentiality"`
	Cvss               types.Float64 `tfsdk:"cvss"`
	Integrity          types.String  `tfsdk:"integrity"`
	PrivilegesRequired types.String  `tfsdk:"privileges_required"`
	Scope              types.String  `tfsdk:"scope"`
	UserInteraction    types.String  `tfsdk:"user_interaction"`
}

func (v VectorsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"attack_complexity":   types.StringType,
		"attack_vector":       types.StringType,
		"availability":        types.StringType,
		"base_severity":       types.StringType,
		"confidentiality":     types.StringType,
		"cvss":                types.Float64Type,
		"integrity":           types.StringType,
		"privileges_required": types.StringType,
		"scope":               types.StringType,
		"user_interaction":    types.StringType,
	}
}

var VectorsModelType = basetypes.ObjectType{AttrTypes: VectorsModel{}.AttributeTypes()}

// ---- Scan Summary (shared) ----

type ScanSummaryModel struct {
	Critical           types.Int32 `tfsdk:"critical"`
	High               types.Int32 `tfsdk:"high"`
	Low                types.Int32 `tfsdk:"low"`
	Medium             types.Int32 `tfsdk:"medium"`
	Negligible         types.Int32 `tfsdk:"negligible"`
	TotalVulnerability types.Int32 `tfsdk:"total_vulnerability"`
	Unknown            types.Int32 `tfsdk:"unknown"`
}

func (v ScanSummaryModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"critical":            types.Int32Type,
		"high":                types.Int32Type,
		"low":                 types.Int32Type,
		"medium":              types.Int32Type,
		"negligible":          types.Int32Type,
		"total_vulnerability": types.Int32Type,
		"unknown":             types.Int32Type,
	}
}

var ScanSummaryModelType = basetypes.ObjectType{AttrTypes: ScanSummaryModel{}.AttributeTypes()}

// ---- Shared flatten helpers ----

func flattenScanSummary(ctx context.Context, summary scr11.ScanSummary) (types.Object, diag.Diagnostics) {
	model := ScanSummaryModel{
		Critical:           types.Int32Value(summary.Critical),
		High:               types.Int32Value(summary.High),
		Low:                types.Int32Value(summary.Low),
		Medium:             types.Int32Value(summary.Medium),
		Negligible:         types.Int32Value(summary.Negligible),
		TotalVulnerability: types.Int32Value(summary.TotalVulnerability),
		Unknown:            types.Int32Value(summary.Unknown),
	}
	return types.ObjectValueFrom(ctx, ScanSummaryModelType.AttrTypes, model)
}

func scanSummaryDataSourceSchemaAttrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"critical":            schema.Int32Attribute{Computed: true},
		"high":                schema.Int32Attribute{Computed: true},
		"low":                 schema.Int32Attribute{Computed: true},
		"medium":              schema.Int32Attribute{Computed: true},
		"negligible":          schema.Int32Attribute{Computed: true},
		"total_vulnerability": schema.Int32Attribute{Computed: true},
		"unknown":             schema.Int32Attribute{Computed: true},
	}
}
