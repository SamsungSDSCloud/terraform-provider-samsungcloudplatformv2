package acceptance_test

import (
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFullResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "~> 0.9.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFullTemplate(),
			},
		},
	})
}

func testAccFullTemplate() string {
	return `
locals {
  name_prefix    = "tfm"
  environment    = "t01"
  common_tags = {
    Project     = "terraformtest"
    Environment = local.environment
    Owner       = "terra"
  }
}

resource "samsungcloudplatformv2_vpc_vpc" "my_vpc" {
  name        = "${local.name_prefix}-vpc-${local.environment}"
  cidr        = "192.168.0.0/16"
  description = "Vpc generated from Terraform"
  tags        = local.common_tags
}

resource "samsungcloudplatformv2_vpc_subnet" "lb_subnet" {
  name        = "${local.name_prefix}PUBSUB${local.environment}"
  vpc_id      = samsungcloudplatformv2_vpc_vpc.my_vpc.id
  type        = "GENERAL"
  cidr        = "192.168.0.0/24"
  dns_nameservers    = ["8.8.8.8"]
  description = "Loadbalancer Subnet"
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

resource "samsungcloudplatformv2_vpc_subnet" "db_subnet" {
  name        = "${local.name_prefix}DBSUB${local.environment}"
  vpc_id      = samsungcloudplatformv2_vpc_vpc.my_vpc.id
  type        = "GENERAL"
  cidr        = "192.168.100.0/24"
  dns_nameservers    = ["8.8.8.8"]
  description = "DB Subnet"
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

resource "samsungcloudplatformv2_security_group_security_group" "lb_sg" {
  name        = "${local.name_prefix}-lb-SG-${local.environment}"
  description = "SecurityGroup generated from terraform"
  loggable    = false
  tags        = local.common_tags
}

resource "samsungcloudplatformv2_security_group_security_group" "bastion_sg" {
  name        = "${local.name_prefix}-bastion-SG-${local.environment}"
  description = "SecurityGroup generated from terraform"
  loggable    = false
  tags        = local.common_tags
}

resource "samsungcloudplatformv2_security_group_security_group" "k8s_sg" {
  name        = "${local.name_prefix}-k8s-SG-${local.environment}"
  description = "SecurityGroup generated from terraform"
  loggable    = false
  tags        = local.common_tags
}

resource "samsungcloudplatformv2_security_group_security_group_rule" "my_sg_rule_lb_http" {
  security_group_id = samsungcloudplatformv2_security_group_security_group.lb_sg.id
  ethertype         = "IPv4"
  protocol          = "TCP"
  direction         = "ingress"
  description       = "SecurityGroup Rule generated from Terraform"
  remote_ip_prefix  = "0.0.0.0/0"
  port_range_min    = 80
  port_range_max    = 80
}

resource "samsungcloudplatformv2_security_group_security_group_rule" "my_sg_rule_lb_https" {
  security_group_id = samsungcloudplatformv2_security_group_security_group.lb_sg.id
  ethertype         = "IPv4"
  protocol          = "TCP"
  direction         = "ingress"
  description       = "SecurityGroup Rule generated from Terraform"
  remote_ip_prefix  = "0.0.0.0/0"
  port_range_min    = 443
  port_range_max    = 443
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

resource "samsungcloudplatformv2_security_group_security_group_rule" "my_sg_rule_bastion_out_http" {
  security_group_id = samsungcloudplatformv2_security_group_security_group.bastion_sg.id
  ethertype         = "IPv4"
  protocol          = "TCP"
  direction         = "egress"
  description       = "SecurityGroup Rule generated from Terraform"
  remote_ip_prefix  = "0.0.0.0/0"
  port_range_min    = 80
  port_range_max    = 80
}

resource "samsungcloudplatformv2_security_group_security_group_rule" "my_sg_rule_bastion_out_https" {
  security_group_id = samsungcloudplatformv2_security_group_security_group.bastion_sg.id
  ethertype         = "IPv4"
  protocol          = "TCP"
  direction         = "egress"
  description       = "SecurityGroup Rule generated from Terraform"
  remote_ip_prefix  = "0.0.0.0/0"
  port_range_min    = 443
  port_range_max    = 443
}

resource "samsungcloudplatformv2_security_group_security_group_rule" "allow_bastion_mariadb" {
  security_group_id = samsungcloudplatformv2_security_group_security_group.bastion_sg.id
  ethertype         = "IPv4"
  protocol          = "TCP"
  direction         = "egress"
  description       = "SecurityGroup Rule generated from Terraform"
  remote_ip_prefix  = "192.168.100.0/24"
  port_range_min    = 2866
  port_range_max    = 2866
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

resource "samsungcloudplatformv2_security_group_security_group_rule" "allow_ssh_internal" {
  security_group_id = samsungcloudplatformv2_security_group_security_group.bastion_sg.id
  ethertype         = "IPv4"
  protocol          = "TCP"
  direction         = "ingress"
  description       = "SecurityGroup Rule generated from Terraform"
  remote_ip_prefix  = "192.168.0.1"
  port_range_min    = 22
  port_range_max    = 22
}

resource "time_sleep" "publicip_destroy_delay" {
  depends_on = [samsungcloudplatformv2_vpc_publicip.natgw_publicip]
  destroy_duration = "20s"
}

resource "samsungcloudplatformv2_vpc_publicip" "natgw_publicip" {
  depends_on = [samsungcloudplatformv2_vpc_internet_gateway.my_igw]
  description = "NAT Gateway Public ip"
  type        = "IGW"
  tags        = local.common_tags
}

resource "samsungcloudplatformv2_vpc_nat_gateway" "natgateway" {
  depends_on  = [time_sleep.publicip_destroy_delay]
  subnet_id   = samsungcloudplatformv2_vpc_subnet.k8s_subnet.id
  publicip_id = samsungcloudplatformv2_vpc_publicip.natgw_publicip.id
  description = "NAT Gateway generated from Terraform"
  tags        = local.common_tags
}
`
}
