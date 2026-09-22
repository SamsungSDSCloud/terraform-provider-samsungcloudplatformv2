provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_internet_gateway" "internetgateway" {
  vpc_id            = var.igw_vpc_id
  type              = var.igw_type
  description       = var.igw_description
  loggable          = var.loggable
  firewall_enabled  = var.firewall_enabled
  firewall_loggable = var.firewall_loggable
  tags              = var.tags
}
