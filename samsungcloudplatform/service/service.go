package service

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/backup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/baremetal"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/baremetalblockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/budget"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/cachestore"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/certificatemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/configinspection"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/directconnect"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/epas"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/eventstreams"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/gslb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/loggingaudit"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/mariadb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/mysql"
	network_logging "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/network-logging"
	billing "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/plannedcompute"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/postgresql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/quota"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/searchengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/sqlserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/vertica"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/vpn"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var ResourceConstructors = []func() resource.Resource{
	// IAM
	iam.NewIamAccessKeyResource,
	iam.NewIamGroupResource,
	iam.NewIamGroupMemberResource,
	iam.NewIamGroupPolicyBindingsResource,
	iam.NewIamPolicyResource,
	iam.NewIamRoleResource,
	iam.NewIamRolePolicyBindingsResource,
	iam.NewIamUserResource,
	iam.NewIamUserPolicyBindingsResource,

	resourcemanager.NewResourceManagerResourceGroupResource,
	ske.NewSkeClusterResource,
	ske.NewSkeNodepoolResource,
	vpc.NewVpcVpcResource,
	vpc.NewVpcSubnetResource,
	vpc.NewVpcPublicipResource,
	vpc.NewVpcPortResource,
	vpc.NewVpcNatGatewayResource,
	vpc.NewVpcInternetGatewayResource,
	vpn.NewVpnVpnGatewayResource,
	vpn.NewVpnVpnTunnelResource,
	vpc.NewVpcVpcEndpointResource,
	vpc.NewVpcPrivateNatResource,
	vpc.NewVpcPrivateNatIpResource,
	vpc.NewVpcPeeringApprovalResource,
	vpc.NewVpcVpcPeeringRuleResource,
	vpc.NewVpcTgwResource,
	vpc.NewVpcTgwRuleResource,
	vpc.NewVpcTgwVpcConnectionResource,
	network_logging.NewNetworkLoggingNetworkLoggingStorageResource,
	directconnect.NewDirectConnectDirectConnectResource,

	directconnect.NewDirectConnectRoutingRuleResource,
	securitygroup.NewSecurityGroupResource,
	securitygroup.NewSecurityGroupRuleResource,

	firewall.NewFirewallFirewallRuleResource,
	billing.NewBillingPlannedComputeResource,

	// Compute
	virtualserver.NewVirtualServerKeypairResource,
	virtualserver.NewVirtualServerVolumeResource,
	virtualserver.NewVirtualServerServerResource,
	virtualserver.NewVirtualServerImageResource,
	backup.NewBackupBackupResource,
	virtualserver.NewVirtualServerServerGroupResource,
	baremetal.NewBaremetalBaremetalResource,

	// Block storage(BM)
	baremetalblockstorage.NewBaremetalBlockStorageVolumeResource,

	// Storage
	filestorage.NewFileStorageVolumeResource,
	filestorage.NewFileStorageSnapshotScheduleResource,
	filestorage.NewFileStorageReplicationResource,

	// Database
	mysql.NewMysqlClusterResource,
	mariadb.NewMariadbClusterResource,
	postgresql.NewPostgresqlClusterResource,

	epas.NewEpasClusterResource,
	sqlserver.NewSqlserverClusterResource,
	cachestore.NewCachestoreClusterResource,
	searchengine.NewSearchengineClusterResource,
	eventstreams.NewEventstreamsClusterResource,
	vertica.NewVerticaClusterResource,

	// LoadBalancer
	loadbalancer.NewLoadbalancerLoadbalancerPublicNatIpResource,
	loadbalancer.NewLoadBalancerLoadBalancerResource,
	loadbalancer.NewLoadBalancerLbServerGroupResource,
	loadbalancer.NewLoadBalancerLbHealthCheckResource,

	// Monitoring
	cloudmonitoring.NewCloudMonitoringEventPolicyResource,

	// Gslb
	gslb.NewGslbGslbResource,

	// Dns
	dns.NewDnsPrivateDnsResource,
	dns.NewDnsPublicDomainNameResource,
	dns.NewDnsHostedZoneResource,
	dns.NewDnsRecordResource,

	// Loggingaudit
	loggingaudit.NewLoggingauditTrailResource,
	loadbalancer.NewLoadBalancerLbMemberResource,
	loadbalancer.NewLoadBalancerListenerResource,

	//peering
	vpc.NewVpcPeeringResource,

	//budget
	budget.NewBudgetBudgetResource,

	configinspection.NewConfigInspectionDiagnosisResource,

	//certificatemanager
	certificatemanager.NewCertificateManagerResource,
	certificatemanager.NewCertificateManagerSelfSignResource,
}

var DataSourceConstructors = []func() datasource.DataSource{
	// IAM
	iam.NewIamAccessKeyDataSource,
	iam.NewIamGroupDataSource,
	iam.NewIamGroupDataSources,
	iam.NewIamGroupMemberDataSources,
	iam.NewIamGroupPolicyBindingDataSources,
	iam.NewIamPolicyDataSource,
	iam.NewIamPolicyDataSources,
	iam.NewIamRoleDataSource,
	iam.NewIamRoleDataSources,
	iam.NewIamRolePolicyBindingDataSources,
	iam.NewIamUserDataSource,
	iam.NewIamUserDataSources,
	iam.NewIamUserPolicyBindingDataSources,

	resourcemanager.NewResourceManagerTagDataSource,
	resourcemanager.NewResourceManagerResourceTagDataSource,
	resourcemanager.NewResourceManagerResourceGroupDataSource,
	resourcemanager.NewResourceManagerResourceGroupDataSources,

	quota.NewQuotaAccountQuotaDataSource,
	quota.NewQuotaAccountQuotaDataSources,

	ske.NewSkeClusterDataSource,
	ske.NewSkeClusterDataSources,
	ske.NewSkeClusterKubeconfigDataSource,
	ske.NewSkeClusterUserKubeconfigDataSource,
	ske.NewSkeKubernetesVersionDataSources,
	ske.NewSkeNodepoolDataSource,
	ske.NewSkeNodepoolDataSources,
	ske.NewSkeNodepoolnodeDataSources,

	vpc.NewVpcVpcDataSource,
	vpc.NewVpcSubnetDataSource,
	vpc.NewVpcPublicipDataSource,
	vpc.NewVpcPortDataSource,
	vpn.NewVpnVpnGatewayDataSource,
	vpn.NewVpnVpnGatewayDataSources,
	vpn.NewVpnVpnTunnelDataSource,
	vpn.NewVpnVpnTunnelDataSources,
	vpc.NewVpcNatGatewayDataSource,
	vpc.NewVpcInternetGatewayDataSource,
	vpc.NewVpcVpcEndpointDataSource,
	vpc.NewVpcVpcPeeringRuleDataSource,
	vpc.NewVpcPrivateNatDataSource,
	vpc.NewVpcPrivateNatIpDataSource,
	vpc.NewTransitGatewayDataSources,
	vpc.NewTransitGatewayDataSource,
	vpc.NewTransitGatewayRoutingRuleDataSources,
	vpc.NewTransitGatewayVpcConnectionDataSources,
	network_logging.NewNetworkLoggingNetworkLoggingStorageDataSource,
	network_logging.NewNetworkLoggingNetworkLoggingConfigurationDataSource,
	directconnect.NewDirectConnectDirectConnectDataSource,
	billing.NewBillingPlannedComputeDataSource,

	directconnect.NewNetworkDirectConnectRoutingRuleDataSource,
	securitygroup.NewSecurityGroupDataSource,
	securitygroup.NewSecurityGroupRuleDataSource,

	securitygroup.NewSecurityGroupDataSources,
	securitygroup.NewSecurityGroupRuleDataSources,

	firewall.NewFirewallFirewallDataSource,
	firewall.NewFirewallFirewallDataSources,
	firewall.NewFirewallFirewallRuleDataSource,
	firewall.NewFirewallFirewallRuleDataSources,

	// Compute
	virtualserver.NewVirtualServerVolumeDataSource,
	virtualserver.NewVirtualServerVolumeDataSources,
	virtualserver.NewVirtualServerKeypairDataSource,
	virtualserver.NewVirtualServerKeypairDataSources,
	virtualserver.NewVirtualServerServerDataSource,
	virtualserver.NewVirtualServerServerDataSources,
	virtualserver.NewVirtualServerImageDataSource,
	virtualserver.NewVirtualServerImageDataSources,
	virtualserver.NewVirtualServerServerGroupDataSource,
	virtualserver.NewVirtualServerServerGroupDataSources,
	backup.NewBackupBackupDataSource,
	backup.NewBackupBackupDataSources,
	baremetal.NewBaremetalBaremetalDataSources,
	baremetal.NewBaremetalBaremetalDataSource,

	// Storage
	filestorage.NewFileStorageVolumeDataSources,
	filestorage.NewFileStorageSnapshotScheduleDataSource,
	filestorage.NewFileStorageVolumeDataSource,
	filestorage.NewFileStorageReplicationDataSource,
	filestorage.NewFileStorageReplicationDataSources,

	// Database
	mysql.NewMysqlClusterDataSource,
	mysql.NewMysqlClusterDataSources,
	mysql.NewMysqlEngineVersionDataSources,
	mariadb.NewMariadbClusterDataSource,
	mariadb.NewMariadbClusterDataSources,
	mariadb.NewMariadbEngineVersionDataSources,
	postgresql.NewPostgresqlClusterDataSource,
	postgresql.NewPostgresqlClusterDataSources,
	postgresql.NewPostgresqlEngineVersionDataSources,
	epas.NewEpasClusterDataSource,
	epas.NewEpasClusterDataSources,
	epas.NewEpasEngineVersionDataSources,
	sqlserver.NewSqlserverClusterDataSource,
	sqlserver.NewSqlserverClusterDataSources,
	sqlserver.NewSqlserverlEngineVersionDataSources,
	cachestore.NewCachestoreClusterDataSource,
	cachestore.NewCachestoreClusterDataSources,
	cachestore.NewCachestoreEngineVersionDataSources,
	searchengine.NewSearchengineClusterDataSource,
	searchengine.NewSearchengineClusterDataSources,
	searchengine.NewSearchengineEngineVersionDataSources,
	vertica.NewVerticaClusterDataSource,
	vertica.NewVerticaClusterDataSources,
	vertica.NewVerticaEngineVersionDataSources,
	eventstreams.NewEventstreamsClusterDataSource,
	eventstreams.NewEventstreamsClusterDataSources,
	eventstreams.NewEventstreamsEngineVersionDataSources,

	// LoadBalancer
	loadbalancer.NewLoadbalancerLoadbalancerDataSources,
	loadbalancer.NewLoadbalancerLoadbalancerDataSource,
	loadbalancer.NewLoadBalancerLbServerGroupDataSources,
	loadbalancer.NewLoadbalancerLbServerGroupDataSource,
	loadbalancer.NewLoadbalancerLbMemberDataSource,
	loadbalancer.NewLoadbalancerLbMemberDataSources,
	loadbalancer.NewLoadbalancerLbHealthCheckDataSource,
	loadbalancer.NewLoadbalancerLbHealthCheckDataSources,
	loadbalancer.NewLoadbalancerLbListenerDataSources,
	loadbalancer.NewLoadbalancerLbListenerDataSource,

	// Monitoring
	cloudmonitoring.NewCloudMonitoringEventDataSource,
	cloudmonitoring.NewCloudMonitoringEventDataSources,
	cloudmonitoring.NewCloudMonitoringEventAccountDataSources,
	cloudmonitoring.NewCloudMonitoringEventNotificationStateDataSources,
	cloudmonitoring.NewCloudMonitoringEventPolicyDataSource,
	cloudmonitoring.NewCloudMonitoringEventPolicyDataSources,
	cloudmonitoring.NewCloudMonitoringEventPolicyHistoryDataSources,
	cloudmonitoring.NewCloudMonitoringEventPolicyNotificationDataSources,
	cloudmonitoring.NewCloudMonitoringMetricDataSources,
	cloudmonitoring.NewCloudMonitoringProductTypeSources,
	cloudmonitoring.NewCloudMonitoringAccountProductSources,
	cloudmonitoring.NewCloudMonitoringAccountMemberSources,
	cloudmonitoring.NewCloudMonitoringAddressBookSources,
	cloudmonitoring.NewCloudMonitoringAddressMemberBookSources,
	cloudmonitoring.NewCloudMonitoringMetricPerfDataDataSources,

	// Gslb
	gslb.NewGslbGslbDataSources,
	gslb.NewGslbGslbDataSource,
	gslb.NewGslbGslbResourceDataSources,
	gslb.NewGslbGslbRRCDataSources,
	gslb.NewgslbGslbRRCUpdateDataSource,

	// Dns
	dns.NewDnsPrivateDnsDataSources,
	dns.NewDnsPrivateDnsDataSource,
	dns.NewDnsPublicDomainNameDataSources,
	dns.NewDnsPublicDomainNameDataSource,
	dns.NewDnsHostedZoneDataSources,
	dns.NewDnsHostedZoneDataSource,
	dns.NewDnsRecordDataSources,
	dns.NewDnsRecordDataSource,

	// Loggingaudit
	loggingaudit.NewLoggingauditTrailDataSource,
	loadbalancer.NewLoadbalancerLbCertificateDataSources,
	loadbalancer.NewLoadbalancerLbCertificateDataSource,

	// config inspections
	configinspection.NewConfigInspectionConfigInspectionDataSource,
	configinspection.NewConfigInspectionConfigInspectionDataSources,
	configinspection.NewConfigInspectionDiagnosisDataSources,
	configinspection.NewConfigInspectionDiagnosisDataSource,
	configinspection.NewConfigInspectionDiagnosisRequestDataSource,

	// vpc peering
	vpc.NewVpcVpcPeeringsDataSource,
	vpc.NewVpcVpcPeeringIdDataSource,

	// certificatemanager
	certificatemanager.NewCertificateManagerDataSource,
	certificatemanager.NewCertificateManagerDetailDataSource,

	// budget
	budget.NewBudgetBudgetDataSource,
	budget.NewBudgetBudgetDataSources,
}

var EphemeralResourceConstructors = []func() ephemeral.EphemeralResource{}
