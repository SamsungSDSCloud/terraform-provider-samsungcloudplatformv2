package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"log"
	"reflect"
	"regexp"
	"strings"
)

const (
	NetworkProductGroup      string = "NETWORKING"
	VpcProductName           string = "VPC Traffic"
	DirectConnectProductName string = "Direct Connect"
	PublicIpProductName      string = "Reserved IP"
	SecurityGroupProductName string = "Security Group"
	FileStorageProductName   string = "File Storage(New)"

	StorageProductGroup           string = "STORAGE"
	ContainerProductGroup         string = "CONTAINER"
	KubernetesEngineVmProductName string = "Kubernetes Engine VM"

	ContractProductType string = "CONTRACT_DISCOUNT"

	CreatingState    string = "CREATING"
	ReservedState    string = "RESERVED"
	ActiveState      string = "ACTIVE"
	InActiveState    string = "INACTIVE"
	DeployingState   string = "DEPLOYING"
	DeletedState     string = "DELETED"
	TerminatingState string = "TERMINATING"
	RunningState     string = "RUNNING"
	AvailableState   string = "AVAILABLE"
	UnavailableState string = "UNAVAILABLE"
	UnknownState     string = "UNKNOWN"
	ErrorState       string = "ERROR"
	EditingState     string = "EDITING"
	StartingState    string = "STARTING"
	StoppingState    string = "STOPPING"
	StoppedState     string = "STOPPED"
	RestartingState  string = "RESTARTING"
	SoftDeletedState string = "SOFT_DELETED"
	UpgradingState   string = "UPGRADING"
	InUseState       string = "IN_USE"
	DeletingState    string = "DELETING"
	TerminatedState  string = "TERMINATED"

	VpcPublicIpPurpose            string = "NAT"
	VpcPublicIpNetworkServiceType string = "VPC"

	ServicedGroupCompute       string = "COMPUTE"
	ServicedForVirtualServer   string = "Virtual Server"
	ServicedForGpuServer       string = "GPU Server"
	ServicedForBaremetalServer string = "Baremetal Server"
	ServicedGroupDatabase      string = "DATABASE"
	ServicedForPostgresql      string = "PostgreSQL"
	ServicedForMariadb         string = "Mariadb"
	ServicedForMySql           string = "MySql"
	ServicedForEpas            string = "EPAS"
	ServicedForSqlServer       string = "Microsoft SQL Server"
	ServicedForTibero          string = "Tibero"

	ProductTypeDisk string = "DISK"

	// Product & Product Group state

	ProductActiveState    string = "ACTIVE"
	ProductAvailableState string = "AVAILABLE"

	// Product & Product Group key

	ProductDefaultDisk      string = "DEFAULT_DISK"
	ProductDisk             string = "DISK"
	ProductScale            string = "SCALE"
	ProductIP               string = "IP"
	ProductContractDiscount string = "CONTRACT_DISCOUNT"
	ProductCloudType        string = "CLOUD_TYPE"
	ProductMonitoringTool   string = "MONITORING_TOOL"
	ProductOS               string = "OS"
	ProductPGLevel1         string = "PG_LEVEL1"
	ProductPGLevel2         string = "PG_LEVEL2"
	ProductServiceLevel     string = "SERVICE_LEVEL"

	OsTypeWindows string = "WINDOWS"

	LinuxAdminAccount   string = "root"
	WindowsAdminAccount string = "Administrator"

	DeploymentEnvironmentDev string = "DEV"
	DeploymentEnvironmentPrd string = "PRD"

	BlockStorageTypeOS      string = "OS"
	BlockStorageTypeData    string = "DATA"
	BlockStorageTypeArchive string = "ARCHIVE"
)

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func EncodeBase64(input string) string {
	base64Encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return base64Encoded
}

func HasField(obj interface{}, fieldName string) bool {
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return false
	}

	field := val.FieldByName(fieldName)
	return field.IsValid()
}

func ToMap(in any) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &inInterface)

	m := map[string]interface{}{}

	for field, val := range inInterface {
		field = ToSnakeCase(field)

		log.Println("KV Pair: ", field, val)

		var typeOfValue []interface{} //TODO:

		if reflect.TypeOf(val) == reflect.TypeOf(typeOfValue) {
			m[field] = ConvertStructToMaps(val.([]interface{}))
		} else {
			m[field] = val
		}

	}
	return m
}

func ConvertStructToMaps[T any](contents []T) []map[string]interface{} {
	var contentMaps []map[string]interface{}

	for _, content := range contents {
		contentMaps = append(contentMaps, ToMap(content))
	}

	return contentMaps
}

func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("unsupported type %T", v)
	}
}

func Contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func GetIdFromSrn(srn string) string {
	parts := strings.Split(srn, "/")
	id := parts[len(parts)-1]
	return id
}

func NullTagCheck(tagsMap types.Map, plan_tags types.Map) types.Map {
	if len(tagsMap.Elements()) == 0 && plan_tags.IsNull() {
		return types.MapNull(types.StringType)
	}
	return tagsMap
}

func ToNullableBoolValue(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*v)
}

func ToNullableInt32Value(v *int32) types.Int32 {
	if v == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*v)
}
