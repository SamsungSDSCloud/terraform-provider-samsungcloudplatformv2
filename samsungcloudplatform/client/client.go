package client

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/backup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/baremetal"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/baremetalblockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/billing"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/budget"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/cachestore"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/certificatemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/configinspection"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/directconnect"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/epas"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/eventstreams"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/gslb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loggingaudit"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/mariadb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/mysql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/networklogging"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/postgresql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/quota"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/searchengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/sqlserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vertica"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/config"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"net/http"
)

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SCPClient struct {
	// CertificateManager
	CertificateManager *certificatemanager.Client

	// VPC
	Vpc *vpc.Client

	// VPC
	VpcV1 *vpcv1.Client

	// DirectConnect
	DirectConnect *directconnect.Client

	// Firewall
	Firewall *firewall.Client

	// VPN
	Vpn *vpn.Client

	// NetworkLogging
	NetworkLogging *networklogging.Client

	// SecurityGroup
	SecurityGroup *securitygroup.Client

	// Kubernetes
	Ske *ske.Client

	// Compute
	VirtualServer *virtualserver.Client
	Backup        *backup.Client
	Baremetal     *baremetal.Client

	// Storage
	BaremetalBlockStorage *baremetalblockstorage.Client
	FileStorage           *filestorage.Client

	// Database
	Mysql        *mysql.Client
	Mariadb      *mariadb.Client
	Postgresql   *postgresql.Client
	Epas         *epas.Client
	Sqlserver    *sqlserver.Client
	Cachestore   *cachestore.Client
	Searchengine *searchengine.Client
	Eventstreams *eventstreams.Client
	Vertica      *vertica.Client

	// Platform
	Iam             *iam.Client
	ResourceManager *resourcemanager.Client
	Billing         *billing.Client
	Budget          *budget.Client
	LoggingAudit    *loggingaudit.Client
	Quota           *quota.Client

	// LoadBalancer
	LoadBalancer *loadbalancer.Client

	// Monitoring
	CloudMonitoring *cloudmonitoring.Client

	// Gslb
	Gslb *gslb.Client

	// Dns
	Dns *dns.Client

	// Misc.

	// Config
	config *config.ProviderConfig

	// Security
	ConfigInspection *configinspection.Client
}

var AllowSDKDefaultVersion = map[string][]string{
	// VPC
	vpc.ServiceType: {"v1.1"},

	// CertificateManager V1
	certificatemanager.ServiceType: {"v1.1"},

	// DirectConnect
	directconnect.ServiceType: {"v1.0"},

	// Firewall
	firewall.ServiceType: {"v1.0"},

	// VPN
	vpn.ServiceType: {"v1.1"},

	// NetworkLogging
	networklogging.ServiceType: {"v1.0"},

	// SecurityGroup
	securitygroup.ServiceType: {"v1.0"},

	// Kubernetes
	ske.ServiceType: {"v1.1"},

	// Compute
	virtualserver.ServiceType: {"v1.1"},
	backup.ServiceType:        {"v1.1"},
	baremetal.ServiceType:     {"v1.1"},

	// Storage
	baremetalblockstorage.ServiceType: {"v1.2"},
	filestorage.ServiceType:           {"v1.1"},

	// Database
	mysql.ServiceType:        {"v1.0"},
	mariadb.ServiceType:      {"v1.0"},
	postgresql.ServiceType:   {"v1.0"},
	epas.ServiceType:         {"v1.0"},
	sqlserver.ServiceType:    {"v1.0"},
	cachestore.ServiceType:   {"v1.0"},
	searchengine.ServiceType: {"v1.0"},
	eventstreams.ServiceType: {"v1.0"},
	vertica.ServiceType:      {"v1.0"},

	// Platform
	iam.ServiceType:             {"v1.2"},
	resourcemanager.ServiceType: {"v1.0"},
	billing.ServiceType:         {"v1.0"},
	budget.ServiceType:          {"v1.0"},
	loggingaudit.ServiceType:    {"v1.1"},
	quota.ServiceType:           {"v1.2"},

	// LoadBalancer
	loadbalancer.ServiceType: {"v1.2"},

	// Monitoring
	cloudmonitoring.ServiceType: {"v1.0"},

	// Gslb
	gslb.ServiceType: {"v1.1"},

	// Dns
	dns.ServiceType: {"v1.3"},

	// ConfigInspection
	configinspection.ServiceType: {"v1.1"},

	// Misc.

}

func NewDefaultConfig(config *config.ProviderConfig, serviceType string) *scpsdk.Configuration {
	tlsConfig, _ := common.CreateTlsConfig()

	cfg := &scpsdk.Configuration{
		AuthUrl:         config.AuthUrl.ValueString(),
		ServiceType:     serviceType,
		AllowSDKVersion: AllowSDKDefaultVersion[serviceType],
		AccountId:       "",
		DefaultRegion:   config.DefaultRegion.ValueString(),
		Endpoint:        config.EndpointOverride.ValueString(),
		Credentials: &scpsdk.Credentials{
			AccessKey: config.AccessKey.ValueString(),
			SecretKey: config.SecretKey.ValueString(),
			AuthToken: config.AuthToken.ValueString(),
		},
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
				Proxy:           http.ProxyFromEnvironment,
			},
			//Timeout: DefaultTimeout, // Default timeout
		},

		DefaultHeader: make(map[string]string),
		UserAgent:     "scpclient/0.0.1",
	}

	return cfg
}

func NewSCPClient(providerConfig *config.ProviderConfig) (*SCPClient, error) {
	client := &SCPClient{
		// CertificateManager
		CertificateManager: certificatemanager.NewClient(NewDefaultConfig(providerConfig, certificatemanager.ServiceType)),

		// VPC
		Vpc: vpc.NewClient(NewDefaultConfig(providerConfig, vpc.ServiceType)),

		VpcV1: vpcv1.NewClient(NewDefaultConfig(providerConfig, vpcv1.ServiceType)),

		// DirectConnect
		DirectConnect: directconnect.NewClient(NewDefaultConfig(providerConfig, directconnect.ServiceType)),

		// Firewall
		Firewall: firewall.NewClient(NewDefaultConfig(providerConfig, firewall.ServiceType)),

		// VPN
		Vpn: vpn.NewClient(NewDefaultConfig(providerConfig, vpn.ServiceType)),

		// NetworkLogging
		NetworkLogging: networklogging.NewClient(NewDefaultConfig(providerConfig, networklogging.ServiceType)),

		// SecurityGroup
		SecurityGroup: securitygroup.NewClient(NewDefaultConfig(providerConfig, securitygroup.ServiceType)),

		// Kubernetes
		Ske: ske.NewClient(NewDefaultConfig(providerConfig, ske.ServiceType)),
		// Compute
		VirtualServer: virtualserver.NewClient(NewDefaultConfig(providerConfig, virtualserver.ServiceType)),
		Backup:        backup.NewClient(NewDefaultConfig(providerConfig, backup.ServiceType)),
		Baremetal:     baremetal.NewClient(NewDefaultConfig(providerConfig, baremetal.ServiceType)),

		// Storage
		BaremetalBlockStorage: baremetalblockstorage.NewClient(NewDefaultConfig(providerConfig, baremetalblockstorage.ServiceType)),
		FileStorage:           filestorage.NewClient(NewDefaultConfig(providerConfig, filestorage.ServiceType)),

		// Database
		Mysql:        mysql.NewClient(NewDefaultConfig(providerConfig, mysql.ServiceType)),
		Mariadb:      mariadb.NewClient(NewDefaultConfig(providerConfig, mariadb.ServiceType)),
		Postgresql:   postgresql.NewClient(NewDefaultConfig(providerConfig, postgresql.ServiceType)),
		Epas:         epas.NewClient(NewDefaultConfig(providerConfig, epas.ServiceType)),
		Sqlserver:    sqlserver.NewClient(NewDefaultConfig(providerConfig, sqlserver.ServiceType)),
		Cachestore:   cachestore.NewClient(NewDefaultConfig(providerConfig, cachestore.ServiceType)),
		Searchengine: searchengine.NewClient(NewDefaultConfig(providerConfig, searchengine.ServiceType)),
		Eventstreams: eventstreams.NewClient(NewDefaultConfig(providerConfig, eventstreams.ServiceType)),
		Vertica:      vertica.NewClient(NewDefaultConfig(providerConfig, vertica.ServiceType)),

		// Platform
		Iam:             iam.NewClient(NewDefaultConfig(providerConfig, iam.ServiceType)),
		ResourceManager: resourcemanager.NewClient(NewDefaultConfig(providerConfig, resourcemanager.ServiceType)),
		Billing:         billing.NewClient((NewDefaultConfig(providerConfig, billing.ServiceType))),
		Budget:          budget.NewClient((NewDefaultConfig(providerConfig, budget.ServiceType))),
		Quota:           quota.NewClient(NewDefaultConfig(providerConfig, quota.ServiceType)),

		// LoadBalancer
		LoadBalancer: loadbalancer.NewClient(NewDefaultConfig(providerConfig, loadbalancer.ServiceType)),

		// Monitoring
		CloudMonitoring: cloudmonitoring.NewClient(NewDefaultConfig(providerConfig, cloudmonitoring.ServiceType)),

		// Gslb
		Gslb: gslb.NewClient(NewDefaultConfig(providerConfig, gslb.ServiceType)),

		// Dns
		Dns: dns.NewClient(NewDefaultConfig(providerConfig, dns.ServiceType)),

		// LoggingAudit
		LoggingAudit: loggingaudit.NewClient(NewDefaultConfig(providerConfig, loggingaudit.ServiceType)),

		// Misc.

		// Config
		config: providerConfig,

		// Security
		ConfigInspection: configinspection.NewClient(NewDefaultConfig(providerConfig, configinspection.ServiceType)),
	}

	return client, nil
}
