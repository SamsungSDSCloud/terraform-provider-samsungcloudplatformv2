package database

// InstanceGroupResponse mirrors cachestore.RedisInstanceGroupResponse using only Go primitive types.
type InstanceGroupResponse struct {
	BlockStorageGroups []BlockStorageGroupResponse
	Id                 string
	Instances          []InstanceResponse
	RoleType           string
	ServerTypeName     string
}

// BlockStorageGroupResponse mirrors cachestore.BlockStorageGroupResponse using only Go primitive types.
type BlockStorageGroupResponse struct {
	Id         string
	Name       string
	RoleType   string
	SizeGb     int32
	VolumeType string
}

// InstanceResponse mirrors cachestore.RedisInstanceResponse using only Go primitive types.
type InstanceResponse struct {
	Name             string
	PublicIpAddress  string
	PublicIpId       string
	RoleType         string
	ServiceIpAddress string
	ServiceState     string
}
