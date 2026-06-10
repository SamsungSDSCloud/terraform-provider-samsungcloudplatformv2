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
			"local": {
				Source:            "hashicorp/local",
				VersionConstraint: "~> 2.5.0",
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "~> 0.9.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "~> 3.5.0",
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

# Random password for database (14 chars: upper, lower, digit, special)
resource "random_password" "db_password" {
  length           = 14
  upper            = true
  lower            = true
  numeric          = true
  special          = true
  override_special = "!@#"

  min_upper   = 1
  min_lower   = 1
  min_numeric = 1
  min_special = 1
}

data "samsungcloudplatformv2_virtualserver_image" "ubuntu_2404" {
  name = "Ubuntu 24.04"
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

# Delay to ensure resources are deleted before PublicIP
# Create: IGW → PublicIP → time_sleep → Server/NAT
# Destroy: Server/NAT → (wait 20s) → time_sleep → PublicIP
resource "time_sleep" "publicip_destroy_delay" {
  depends_on = [samsungcloudplatformv2_vpc_publicip.natgw_publicip]
  destroy_duration = "20s"
}

# Delay for bastion server (VM attached to bastion_publicip)
resource "time_sleep" "bastion_destroy_delay" {
  depends_on = [samsungcloudplatformv2_vpc_publicip.bastion_publicip]
  destroy_duration = "20s"
}

resource "samsungcloudplatformv2_vpc_publicip" "bastion_publicip" {
  depends_on = [samsungcloudplatformv2_vpc_internet_gateway.my_igw]
  description = "Bastion Public ip"
  type        = "IGW"
  tags        = local.common_tags
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

resource "samsungcloudplatformv2_virtualserver_keypair" "keypair" {
  name = "${local.name_prefix}-keypair-${local.environment}"
  tags        = local.common_tags
}

resource "local_file" "my_keypair" {
  content         = samsungcloudplatformv2_virtualserver_keypair.keypair.private_key
  filename        = "/tmp/${local.name_prefix}-keypair-${local.environment}.pem"
  file_permission = "0600"
}


resource "samsungcloudplatformv2_virtualserver_server" "bastion" {
  depends_on = [time_sleep.bastion_destroy_delay]

  name        = "${local.name_prefix}-bastion-${local.environment}"
  state     = "ACTIVE"
  security_groups = [samsungcloudplatformv2_security_group_security_group.bastion_sg.id]
  image_id = data.samsungcloudplatformv2_virtualserver_image.ubuntu_2404.image.id
  server_type_id  = "s2v1m2"

  boot_volume = {
    size = 48
    type = "SSD"
  }

  networks = {
    interface_1 = {
      subnet_id  = samsungcloudplatformv2_vpc_subnet.lb_subnet.id
      public_ip_id = samsungcloudplatformv2_vpc_publicip.bastion_publicip.id
    }
  }

  keypair_name = samsungcloudplatformv2_virtualserver_keypair.keypair.name

  tags = merge(
    local.common_tags,
    {
      Name    = "${local.name_prefix}-bastion-${local.environment}"
      Role    = "Bastion Host"
      Purpose = "Secure Gateway"
    }
  )
}

resource "samsungcloudplatformv2_filestorage_volume" "k8s_file_storage" {
  name         = "${local.name_prefix}fs_${local.environment}"
  protocol     = "NFS"
  type_name    = "HDD"
  access_rules = []
  tags        = local.common_tags

  lifecycle {
    ignore_changes = [access_rules]
  }
}

resource "samsungcloudplatformv2_ske_cluster" "my_cluster" {
  name                            = "${local.name_prefix}k8s${local.environment}"
  kubernetes_version              = "v1.31.8"
  vpc_id                          = samsungcloudplatformv2_vpc_vpc.my_vpc.id
  subnet_id                       = samsungcloudplatformv2_vpc_subnet.k8s_subnet.id
  security_group_id_list          = [samsungcloudplatformv2_security_group_security_group.k8s_sg.id]
  volume_id                       = samsungcloudplatformv2_filestorage_volume.k8s_file_storage.id
  public_endpoint_access_control_ip = "192.168.0.1"
  cloud_logging_enabled               = false
  service_watch_logging_enabled       = false
}

resource "samsungcloudplatformv2_ske_nodepool" "nodepool" {
  name                 = "${local.name_prefix}-nodepool-${local.environment}"
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

# MariaDB Cluster
resource "samsungcloudplatformv2_mariadb_cluster" "my_mariadb" {
  name                     = "${local.name_prefix}mdb"
  subnet_id                = samsungcloudplatformv2_vpc_subnet.db_subnet.id
  dbaas_engine_version_id  = "1cd2c28ba72447daaaf7e4d7e8dd720b"
  instance_name_prefix     = "${local.name_prefix}-mariadb"
  ha_enabled               = false
  timezone                 = "Asia/Seoul"
  nat_enabled              = false
  service_state            = "RUNNING"

  allowable_ip_addresses = ["192.168.0.0/24", "192.168.50.0/24"]

  init_config_option = {
    audit_enabled          = false
    backup_option          = {}
    database_character_set = "utf8"
    database_name          = "mdb"
    database_port          = 2866
    database_user_name     = "dbadmin"
    database_user_password = random_password.db_password.result
  }

  instance_groups = [
    {
      role_type        = "ACTIVE"
      server_type_name = "db2v1m2"
      instances = [
        {
          role_type = "ACTIVE"
        }
      ]
      block_storage_groups = [
        {
          "role_type"   = "OS"
          "size_gb"     = 104
          "volume_type" = "SSD"
        }
      ]
    }
  ]

  maintenance_option = {
    default = {
      period_hour = null
      starting_day_of_week = null
      starting_time = null
      use_maintenance_option = false
    }
  }
  tags               = local.common_tags
}
`
}
