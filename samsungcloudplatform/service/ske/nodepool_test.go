package ske_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSkeResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "~> 0.9.0",
			},
			"local": {
				Source:            "hashicorp/local",
				VersionConstraint: "~> 2.5.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "~> 3.5.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSkeTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.samsungcloudplatformv2_ske_nodepools.this", "nodepools.0.subnet_id"),
				),
			},
		},
	})
}

func testAccSkeTemplate() string {
	return fmt.Sprintf(`%s %s %s %s %s`,
		getLocals("k8s"), getNetwork(), getKeypair(), getCluster(), getNodepool())
}

func getLocals(environment string) string {
	return fmt.Sprintf(`
		locals {
		  name_prefix    = "zzz"
		  environment    = "%s"
		  common_tags = {
			Project     = "terraformtest"
			Environment = local.environment
			Owner       = "terra"
		  }
		}
	`, environment)
}

func getCluster() string {
	return `
		resource "samsungcloudplatformv2_filestorage_volume" "k8s_file_storage" {
		  name         = "${local.name_prefix}fs${local.environment}"
		  protocol     = "NFS"
		  type_name    = "HDD"
		  tags         = local.common_tags
		}
		
		resource "samsungcloudplatformv2_ske_cluster" "my_cluster" {
		  name                            = "${local.name_prefix}k8s${local.environment}"
		  kubernetes_version              = "v1.31.8"
		  vpc_id                          = samsungcloudplatformv2_vpc_vpc.my_vpc.id
		  default_subnet_id               = samsungcloudplatformv2_vpc_subnet.k8s_subnet.id
		  security_group_id_list          = [samsungcloudplatformv2_security_group_security_group.k8s_sg.id]
		  nfs_volume_id                   = samsungcloudplatformv2_filestorage_volume.k8s_file_storage.id
		  public_endpoint_access_control_ip = "192.168.0.1"
		  service_watch_logging_enabled       = false
		}
	`
}

func getNodepool() string {
	return `
		resource "samsungcloudplatformv2_ske_nodepool" "nodepool" {
		  name                 = "${local.name_prefix}nodepool${local.environment}"
		  cluster_id           = samsungcloudplatformv2_ske_cluster.my_cluster.id
		  image_os             = "ubuntu"
		  image_os_version     = "22.04"
		  kubernetes_version   = "v1.31.8"
		  server_type_id       = "s2v2m4"
		  volume_type_name     = "SSD"
		  volume_size          = "104"
		  keypair_name         = samsungcloudplatformv2_virtualserver_keypair.keypair.name
		  is_auto_recovery     = false
		  is_auto_scale        = true
		  min_node_count       = 1
		  max_node_count       = 3
		  desired_node_count   = 1
		}

		data "samsungcloudplatformv2_ske_nodepools" "this" {
		  cluster_id = samsungcloudplatformv2_ske_cluster.my_cluster.id
		}

		output "nodepool_subnet_id" {
		  value = data.samsungcloudplatformv2_ske_nodepools.this.nodepools[0].subnet_id
		}
	`
}

func getNetwork() string {
	return `
		resource "samsungcloudplatformv2_vpc_vpc" "my_vpc" {
		  name        = "${local.name_prefix}-vpc-${local.environment}"
		  cidr        = "192.168.0.0/16"
		  description = "Vpc generated from Terraform"
		  tags        = local.common_tags
		}
		
		resource "samsungcloudplatformv2_vpc_subnet" "k8s_subnet" {
		  name        = "${local.name_prefix}PRISUB${local.environment}"
		  vpc_id      = samsungcloudplatformv2_vpc_vpc.my_vpc.id
		  type        = "GENERAL"
		  cidr        = "192.168.50.0/24"
		  dns_nameservers    = ["8.8.8.8"]
		  description = "k8s Subnet"
		  tags        = local.common_tags
		}
		
		resource "samsungcloudplatformv2_vpc_internet_gateway" "my_igw" {
		  vpc_id            = samsungcloudplatformv2_vpc_vpc.my_vpc.id
		  type              = "IGW"
		  firewall_enabled  = true
		  firewall_loggable = false
		  description       = "Internet GW generated from Terraform"
		  tags        = local.common_tags
		}
		
		resource "samsungcloudplatformv2_firewall_firewall_rule" "my_igw_fwrule_systemupdate" {
		  depends_on = [samsungcloudplatformv2_vpc_internet_gateway.my_igw]
		  firewall_id = samsungcloudplatformv2_vpc_internet_gateway.my_igw.internet_gateway.firewall_id
		  firewall_rule_create = {
			action              = "ALLOW"
			description         = "Rule from terraform"
			destination_address = ["0.0.0.0/0"]
			direction           = "OUTBOUND"
			service = [
			  {
				service_type  = "TCP"
				service_value = "80"
			  },
			  {
				service_type  = "TCP"
				service_value = "443"
			  }
			]
			source_address = ["192.168.0.0/24", "192.168.50.0/24"]
			status         = "ENABLE"
		  }
		}
		
		resource "samsungcloudplatformv2_firewall_firewall_rule" "my_igw_fwrule_webservice" {
		  depends_on = [
			samsungcloudplatformv2_vpc_internet_gateway.my_igw,
			samsungcloudplatformv2_firewall_firewall_rule.my_igw_fwrule_systemupdate,
		  ]
		  firewall_id = samsungcloudplatformv2_vpc_internet_gateway.my_igw.internet_gateway.firewall_id
		  firewall_rule_create = {
			action              = "ALLOW"
			description         = "Rule from terraform"
			destination_address = ["192.168.0.0/24"]
			direction           = "INBOUND"
			service = [
			  {
				service_type  = "TCP"
				service_value = "80"
			  },
			  {
				service_type  = "TCP"
				service_value = "443"
			  }
			]
			source_address = ["0.0.0.0/0"]
			status         = "ENABLE"
		  }
		}
		
		resource "samsungcloudplatformv2_firewall_firewall_rule" "my_igw_fwrule_k8s" {
		  depends_on = [
			samsungcloudplatformv2_vpc_internet_gateway.my_igw,
			samsungcloudplatformv2_firewall_firewall_rule.my_igw_fwrule_webservice,
		  ]
		  firewall_id = samsungcloudplatformv2_vpc_internet_gateway.my_igw.internet_gateway.firewall_id
		  firewall_rule_create = {
			action              = "ALLOW"
			description         = "Rule from terraform"
			destination_address = ["192.168.50.0/24"]
			direction           = "INBOUND"
			service = [
			  {
				service_type  = "TCP"
				service_value = "80"
			  },
			  {
				service_type  = "TCP"
				service_value = "443"
			  },
			  {
				service_type  = "TCP"
				service_value = "6443"
			  }
			]
			source_address = ["0.0.0.0/0"]
			status         = "ENABLE"
		  }
		}
		
		resource "samsungcloudplatformv2_firewall_firewall_rule" "my_igw_fwrule_ssh" {
		  depends_on = [
			samsungcloudplatformv2_vpc_internet_gateway.my_igw,
			samsungcloudplatformv2_firewall_firewall_rule.my_igw_fwrule_k8s,
		  ]
		  firewall_id = samsungcloudplatformv2_vpc_internet_gateway.my_igw.internet_gateway.firewall_id
		  firewall_rule_create = {
			action              = "ALLOW"
			description         = "Rule from terraform"
			destination_address = ["192.168.0.0/24"]
			direction           = "INBOUND"
			service = [
			  {
				service_type  = "TCP"
				service_value = "22"
			  }
			]
			source_address = ["192.168.0.1"]
			status         = "ENABLE"
		  }
		}
		
		resource "samsungcloudplatformv2_security_group_security_group" "k8s_sg" {
		  name        = "${local.name_prefix}-k8s-SG-${local.environment}"
		  description = "SecurityGroup generated from terraform"
		  loggable    = false
		  tags        = local.common_tags
		}
		
		resource "samsungcloudplatformv2_security_group_security_group_rule" "my_sg_rule_kubectl" {
		  security_group_id = samsungcloudplatformv2_security_group_security_group.k8s_sg.id
		  ethertype         = "IPv4"
		  protocol          = "TCP"
		  direction         = "ingress"
		  description       = "SecurityGroup Rule generated from Terraform"
		  remote_ip_prefix  = "192.168.0.1"
		  port_range_min    = 6443
		  port_range_max    = 6443
		}
		
		resource "samsungcloudplatformv2_security_group_security_group_rule" "my_sg_rule_update_http" {
		  security_group_id = samsungcloudplatformv2_security_group_security_group.k8s_sg.id
		  ethertype         = "IPv4"
		  protocol          = "TCP"
		  direction         = "egress"
		  description       = "SecurityGroup Rule generated from Terraform"
		  remote_ip_prefix  = "0.0.0.0/0"
		  port_range_min    = 80
		  port_range_max    = 80
		}
		
		resource "samsungcloudplatformv2_security_group_security_group_rule" "my_sg_rule_update_https" {
		  security_group_id = samsungcloudplatformv2_security_group_security_group.k8s_sg.id
		  ethertype         = "IPv4"
		  protocol          = "TCP"
		  direction         = "egress"
		  description       = "SecurityGroup Rule generated from Terraform"
		  remote_ip_prefix  = "0.0.0.0/0"
		  port_range_min    = 443
		  port_range_max    = 443
		}
		
		resource "samsungcloudplatformv2_security_group_security_group_rule" "allow_k8s_mariadb" {
		  security_group_id = samsungcloudplatformv2_security_group_security_group.k8s_sg.id
		  ethertype         = "IPv4"
		  protocol          = "TCP"
		  direction         = "egress"
		  description       = "SecurityGroup Rule generated from Terraform"
		  remote_ip_prefix  = "192.168.100.0/24"
		  port_range_min    = 2866
		  port_range_max    = 2866
		}
		
		resource "time_sleep" "publicip_destroy_delay" {
		  depends_on = [samsungcloudplatformv2_vpc_publicip.natgw_publicip]
		  destroy_duration = "120s"
		}
		
		resource "samsungcloudplatformv2_vpc_nat_gateway" "natgateway" {
		  depends_on  = [time_sleep.publicip_destroy_delay]
		  subnet_id   = samsungcloudplatformv2_vpc_subnet.k8s_subnet.id
		  publicip_id = samsungcloudplatformv2_vpc_publicip.natgw_publicip.id
		  description = "NAT Gateway generated from Terraform"
		  tags        = local.common_tags
		}

		resource "samsungcloudplatformv2_vpc_publicip" "natgw_publicip" {
		  description = "NAT Gateway Public ip"
		  type        = "IGW"
		#  tags        = local.common_tags
		# PublicIP 가 IGW 보다 먼저(뒤에) 삭제되게 함
		  depends_on = [samsungcloudplatformv2_vpc_internet_gateway.my_igw]
		}
	`
}

func getKeypair() string {
	return `
		resource "samsungcloudplatformv2_virtualserver_keypair" "keypair" {
		  name = "${local.name_prefix}-keypair-${local.environment}"
		  tags = local.common_tags
		}
		
		resource "local_file" "my_keypair" {
		  content         = samsungcloudplatformv2_virtualserver_keypair.keypair.private_key
		  filename        = pathexpand("./${local.name_prefix}-keypair-${local.environment}.pem")
		  file_permission = "0600"
		}
	`
}
