package virtualserver

// Service Name
const (
	ServiceNameVirtualServer = "virtualserver"
	ServiceNameGpuServer     = "gpuserver"
	ServiceNameVpc           = "vpc"
)

// Resource Type
const (
	ResourceTypeImage         = "image"
	ResourceTypeKeypair       = "keypair"
	ResourceTypeVirtualServer = "virtual-server"
	ResourceTypeGpuServer     = "gpu-server"
	ResourceTypeVolume        = "volume"
	ResourceTypePort          = "port"
	ResourceTypeServerGroup   = "server-group"
)

// Product Offering
const (
	ProductOfferingGpuServer    = "gpu_server"
	ProductOfferingK8sGpuServer = "k8s_gpu_vm"
)

// SCP Image Type
const (
	ScpImageTypeGpuCustom = "gpu_custom"
)
