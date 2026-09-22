package database

// Block storage group role_type value sets shared by the database services.
//
// The DBaaS SDK narrows role_type by the position it appears in, so a single
// allow-list is not correct for every attribute:
//
//   - cluster creation uses the per-product create enum
//     (BlockStorageGroupRoleType / OsDataBlockStorageGroupRoleType /
//     OsDataBackupBlockStorageGroupRoleType)
//   - adding a block storage to an existing instance group uses
//     ExtraBlockStorageGroupRoleType, which excludes OS
//
// Keep these in sync with library/<service>/<version>/model_*block_storage_group_role_type.go
// in the terraform-sdk-samsungcloudplatformv2 checkout.
var (
	// BSRoleTypesFull mirrors BlockStorageGroupRoleType, the create-time enum for
	// the relational products (mysql, mariadb, postgresql, epas).
	BSRoleTypesFull = []string{"OS", "DATA", "ARCHIVE", "TEMP", "BACKUP"}

	// BSRoleTypesOsData mirrors OsDataBlockStorageGroupRoleType, the create-time
	// enum for cachestore, eventstreams, sqlserver and vertica.
	BSRoleTypesOsData = []string{"OS", "DATA"}

	// BSRoleTypesOsDataBackup mirrors OsDataBackupBlockStorageGroupRoleType, the
	// create-time enum for searchengine.
	BSRoleTypesOsDataBackup = []string{"OS", "DATA", "BACKUP"}

	// BSRoleTypesExtra mirrors ExtraBlockStorageGroupRoleType, used when adding a
	// block storage to an existing instance group. OS is deliberately absent: the
	// OS disk is created with the instance and cannot be added later.
	BSRoleTypesExtra = []string{"DATA", "ARCHIVE", "TEMP", "BACKUP"}

	// BSRoleTypesExtraSqlserver is the add-time set for sqlserver, which does NOT
	// share ExtraBlockStorageGroupRoleType with the other products.
	//
	// sqlserver has its own SqlserverExtraBlockStorageGroupRoleType with DATA as the
	// only member. The SDK now generates that enum type (it briefly came through as a
	// bare string while the spec carried an empty enum block — see error.md A-1), but
	// a Terraform-supplied role_type is still a runtime string, so converting it to the
	// SDK type does not validate it. This set is what the provider checks against.
	//
	// DATA is confirmed against the DBaaS API source (osscp/scp-dbaas):
	// scp_dbaas/common/enum/role_type.py declares
	// SqlserverExtraBlockStorageGroupRoleType with DATA as its only member, and
	// scp_dbaas/api/models/sqlserver/sqlserver_block_storage_model.py notes
	// "SQL Server 의 추가 블록 스토리지 역할 유형은 DATA 만 가능 (TEMP, ARCHIVE, BACKUP X)".
	// The empty enum is a spec-generation defect, not an API definition error.
	// See error.md A-1.
	BSRoleTypesExtraSqlserver = []string{"DATA"}
)

// ContainsRoleType reports whether roleType is a member of allowed.
func ContainsRoleType(allowed []string, roleType string) bool {
	for _, v := range allowed {
		if roleType == v {
			return true
		}
	}
	return false
}

// IsExtraBlockStorageRoleType reports whether roleType is accepted by the
// AddBlockStorages API for the products that share
// ExtraBlockStorageGroupRoleType. sqlserver does NOT: use
// BSRoleTypesExtraSqlserver with ContainsRoleType for that product.
func IsExtraBlockStorageRoleType(roleType string) bool {
	return ContainsRoleType(BSRoleTypesExtra, roleType)
}
