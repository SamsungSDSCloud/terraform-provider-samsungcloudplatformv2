package client

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/backup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/baremetalblockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/billing"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/cachestore"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/directconnect"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/epas"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/eventstreams"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/gslb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/loggingaudit"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/mariadb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/mysql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/networklogging"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/postgresql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/quota"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/searchengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/sqlserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/vertica"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/config"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"io/ioutil"
	"net/http"
	"os"
)

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SCPClient struct {
	// VPC
	Vpc *vpc.Client

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
}

func createTlsConfig() (*tls.Config, error) {
	certPath := os.Getenv("SSL_CERT_FILE")
	var certPool *x509.CertPool
	var err error

	if certPath == "" {
		certPool, err = x509.SystemCertPool()
	} else {
		crt, err := ioutil.ReadFile(certPath)
		if err != nil {
			return nil, err
		}
		certPool = x509.NewCertPool()
		certPool.AppendCertsFromPEM(crt)
	}

	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs: certPool,
	}, nil
}

func NewDefaultConfig(config *config.ProviderConfig, serviceType string) *scpsdk.Configuration {
	tlsConfig, _ := createTlsConfig()

	cfg := &scpsdk.Configuration{
		AuthUrl:       config.AuthUrl.ValueString(),
		ServiceType:   serviceType,
		AccountId:     "",
		DefaultRegion: config.DefaultRegion.ValueString(),
		Endpoint:      config.EndpointOverride.ValueString(),
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
		// VPC
		Vpc: vpc.NewClient(NewDefaultConfig(providerConfig, vpc.ServiceType)),

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
	}

	return client, nil
}
