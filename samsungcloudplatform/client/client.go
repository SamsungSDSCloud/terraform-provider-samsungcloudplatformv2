package client

import (
	"fmt"
	"net/http"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/backup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/baremetal"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/baremetalblockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/billing"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/budget"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cachestore"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/certificatemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cloudcontrol"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/configinspection"
	dc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/directconnect"
	dcv1d1 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/directconnectv1d1"
	dcv1d2 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/directconnectv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/epas"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/eventstreams"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/firewall"
	firewallv1d1 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/firewallv1d1"
	firewallv1d2 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/firewallv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/gslb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/iamidentitycenter"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/loadbalancerv1d4"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/loggingaudit"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/mariadb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/multinodegpucluster"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/mysql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/networklogging"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/parallelfilestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/postgresql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/quota"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/scr"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/searchengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/securitygroup"
	securitygroupv1d1 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/securitygroupv1d1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/sqlserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vertica"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1d3"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/config"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
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
	VpcV1     *vpcv1.Client
	VpcV1Dot2 *vpcv1d2.Client
	VpcV1Dot3 *vpcv1d3.Client

	// DirectConnect
	DirectConnect     *dc.Client
	DirectConnectV1d1 *dcv1d1.Client
	DirectConnectV1d2 *dcv1d2.Client

	// Firewall
	Firewall     *firewall.Client
	FirewallV1d1 *firewallv1d1.Client
	FirewallV1d2 *firewallv1d2.Client

	// VPN
	Vpn *vpn.Client

	// NetworkLogging
	NetworkLogging *networklogging.Client

	// SecurityGroup
	SecurityGroup     *securitygroup.Client
	SecurityGroupV1d1 *securitygroupv1d1.Client

	// Kubernetes
	Ske *ske.Client

	// Compute
	VirtualServer *virtualserver.Client
	Backup        *backup.Client
	Baremetal     *baremetal.Client

	// Storage
	BaremetalBlockStorage *baremetalblockstorage.Client
	FileStorage           *filestorage.Client
	ParallelFileStorage   *parallelfilestorage.Client

	// Database
	Mysql        *mysql.Client
	Mariadb      *mariadb.Client
	Postgresql   *postgresql.Client
	Epas         *epas.Client
	Sqlserver    *sqlserver.Client
	Cachestore   *cachestore.Client
	Searchengine *searchengine.Client
	Scr          *scr.Client
	Eventstreams *eventstreams.Client
	Vertica      *vertica.Client

	// Platform
	Iam               *iam.Client
	IamIdentityCenter *iamidentitycenter.Client
	ResourceManager   *resourcemanager.Client
	Billing           *billing.Client
	Budget            *budget.Client
	LoggingAudit      *loggingaudit.Client
	Quota             *quota.Client
	Organization      *organization.Client
	CloudControl      *cloudcontrol.Client

	// LoadBalancer
	LoadBalancer     *loadbalancer.Client
	LoadBalancerV1d4 *loadbalancerv1d4.Client
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

	// Multi-node GPU Cluster
	Mngc *multinodegpucluster.Client
	// ServiceWatch
	ServiceWatch *servicewatch.Client
}

var AllowSDKDefaultVersion = map[string][]string{
	// VPC
	vpc.ServiceType: {"v1.1", "v1.2", "v1.3"},

	// CertificateManager V1
	certificatemanager.ServiceType: {"v1.1"},

	// DirectConnect
	dc.ServiceType: {"v1.0", "v1.1", "v1.2"},

	// Firewall
	firewall.ServiceType: {"v1.0", "v1.1", "v1.2"},

	// VPN
	vpn.ServiceType: {"v1.1"},

	// NetworkLogging
	networklogging.ServiceType: {"v1.0"},

	// SecurityGroup
	securitygroup.ServiceType: {"v1.0", "v1.1"},

	// Kubernetes
	ske.ServiceType: {"v1.6"},

	// Compute
	virtualserver.ServiceType: {"v1.5"},
	backup.ServiceType:        {"v1.4"},
	baremetal.ServiceType:     {"v1.2"},

	// Storage
	baremetalblockstorage.ServiceType: {"v1.4"},
	filestorage.ServiceType:           {"v1.2"},
	parallelfilestorage.ServiceType:   {"v1.1"},

	// Database
	mysql.ServiceType:        {"v1.3"},
	mariadb.ServiceType:      {"v1.3"},
	postgresql.ServiceType:   {"v1.3"},
	epas.ServiceType:         {"v1.3"},
	sqlserver.ServiceType:    {"v1.2"},
	cachestore.ServiceType:   {"v1.2"},
	searchengine.ServiceType: {"v1.2"},
	scr.ServiceType:          {"v1.1"},
	eventstreams.ServiceType: {"v1.2"},
	vertica.ServiceType:      {"v1.2"},

	// Platform
	iam.ServiceType:               {"v1.4"},
	iamidentitycenter.ServiceType: {"v1.6"},
	resourcemanager.ServiceType:   {"v1.0"},
	billing.ServiceType:           {"v1.0"},
	budget.ServiceType:            {"v1.1"},
	loggingaudit.ServiceType:      {"v1.1"},
	quota.ServiceType:             {"v1.5"},
	organization.ServiceType:      {"v1.3"},
	cloudcontrol.ServiceType:      {"v1.2"},

	// LoadBalancer
	loadbalancer.ServiceType: {"v1.3", "v1.4"},

	// Monitoring
	cloudmonitoring.ServiceType: {"v1.0"},

	// Gslb
	gslb.ServiceType: {"v1.2"},

	// Dns
	dns.ServiceType: {"v1.4"},

	// ConfigInspection
	configinspection.ServiceType: {"v1.1"},

	// Multi-node GPU Cluster
	multinodegpucluster.ServiceType: {"v1.3"},
	// ServiceWatch
	servicewatch.ServiceType: {"v1.5"},

	// Misc.

}

func NewDefaultConfig(config *config.ProviderConfig, serviceType string) *scpsdk.Configuration {
	tlsConfig, err := common.CreateTlsConfig()
	if err != nil {
		fmt.Println(
			"Failed to build TLS config",
			"Error details: "+err.Error(),
		)
		return nil
	}

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

		VpcV1:     vpcv1.NewClient(NewDefaultConfig(providerConfig, vpcv1.ServiceType)),
		VpcV1Dot2: vpcv1d2.NewClient(NewDefaultConfig(providerConfig, vpcv1d2.ServiceType)),
		VpcV1Dot3: vpcv1d3.NewClient(NewDefaultConfig(providerConfig, vpcv1d3.ServiceType)),

		// DirectConnect
		DirectConnect:     dc.NewClient(NewDefaultConfig(providerConfig, dc.ServiceType)),
		DirectConnectV1d1: dcv1d1.NewClient(NewDefaultConfig(providerConfig, dcv1d1.ServiceType)),
		DirectConnectV1d2: dcv1d2.NewClient(NewDefaultConfig(providerConfig, dcv1d2.ServiceType)),

		// Firewall
		Firewall:     firewall.NewClient(NewDefaultConfig(providerConfig, firewall.ServiceType)),
		FirewallV1d1: firewallv1d1.NewClient(NewDefaultConfig(providerConfig, firewallv1d1.ServiceType)),
		FirewallV1d2: firewallv1d2.NewClient(NewDefaultConfig(providerConfig, firewallv1d2.ServiceType)),

		// VPN
		Vpn: vpn.NewClient(NewDefaultConfig(providerConfig, vpn.ServiceType)),

		// NetworkLogging
		NetworkLogging: networklogging.NewClient(NewDefaultConfig(providerConfig, networklogging.ServiceType)),

		// SecurityGroup
		SecurityGroup:     securitygroup.NewClient(NewDefaultConfig(providerConfig, securitygroup.ServiceType)),
		SecurityGroupV1d1: securitygroupv1d1.NewClient(NewDefaultConfig(providerConfig, securitygroupv1d1.ServiceType)),

		// Kubernetes
		Ske: ske.NewClient(NewDefaultConfig(providerConfig, ske.ServiceType)),
		// Compute
		VirtualServer: virtualserver.NewClient(NewDefaultConfig(providerConfig, virtualserver.ServiceType)),
		Backup:        backup.NewClient(NewDefaultConfig(providerConfig, backup.ServiceType)),
		Baremetal:     baremetal.NewClient(NewDefaultConfig(providerConfig, baremetal.ServiceType)),

		// Storage
		BaremetalBlockStorage: baremetalblockstorage.NewClient(NewDefaultConfig(providerConfig, baremetalblockstorage.ServiceType)),
		FileStorage:           filestorage.NewClient(NewDefaultConfig(providerConfig, filestorage.ServiceType)),
		ParallelFileStorage:   parallelfilestorage.NewClient(NewDefaultConfig(providerConfig, parallelfilestorage.ServiceType)),

		// Database
		Mysql:        mysql.NewClient(NewDefaultConfig(providerConfig, mysql.ServiceType)),
		Mariadb:      mariadb.NewClient(NewDefaultConfig(providerConfig, mariadb.ServiceType)),
		Postgresql:   postgresql.NewClient(NewDefaultConfig(providerConfig, postgresql.ServiceType)),
		Epas:         epas.NewClient(NewDefaultConfig(providerConfig, epas.ServiceType)),
		Sqlserver:    sqlserver.NewClient(NewDefaultConfig(providerConfig, sqlserver.ServiceType)),
		Cachestore:   cachestore.NewClient(NewDefaultConfig(providerConfig, cachestore.ServiceType)),
		Searchengine: searchengine.NewClient(NewDefaultConfig(providerConfig, searchengine.ServiceType)),
		Scr:          scr.NewClient(NewDefaultConfig(providerConfig, scr.ServiceType)),
		Eventstreams: eventstreams.NewClient(NewDefaultConfig(providerConfig, eventstreams.ServiceType)),
		Vertica:      vertica.NewClient(NewDefaultConfig(providerConfig, vertica.ServiceType)),

		// Platform
		Iam:               iam.NewClient(NewDefaultConfig(providerConfig, iam.ServiceType)),
		IamIdentityCenter: iamidentitycenter.NewClient(NewDefaultConfig(providerConfig, iamidentitycenter.ServiceType)),
		ResourceManager:   resourcemanager.NewClient(NewDefaultConfig(providerConfig, resourcemanager.ServiceType)),
		Billing:           billing.NewClient((NewDefaultConfig(providerConfig, billing.ServiceType))),
		Budget:            budget.NewClient((NewDefaultConfig(providerConfig, budget.ServiceType))),
		Quota:             quota.NewClient(NewDefaultConfig(providerConfig, quota.ServiceType)),
		Organization:      organization.NewClient(NewDefaultConfig(providerConfig, organization.ServiceType)),
		CloudControl:      cloudcontrol.NewClient(NewDefaultConfig(providerConfig, cloudcontrol.ServiceType)),

		// LoadBalancer
		LoadBalancer:     loadbalancer.NewClient(NewDefaultConfig(providerConfig, loadbalancer.ServiceType)),
		LoadBalancerV1d4: loadbalancerv1d4.NewClient(NewDefaultConfig(providerConfig, loadbalancer.ServiceType)),
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

		// Multi-node GPU Cluster
		Mngc: multinodegpucluster.NewClient(NewDefaultConfig(providerConfig, multinodegpucluster.ServiceType)),
		// ServiceWatch
		ServiceWatch: servicewatch.NewClient(NewDefaultConfig(providerConfig, servicewatch.ServiceType)),
	}

	return client, nil
}
