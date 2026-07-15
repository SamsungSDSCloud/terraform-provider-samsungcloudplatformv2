provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_subnet_vips" "subnet_vip_list" {
  subnet_id          = var.subnet_id
  size               = var.size
  page               = var.page
  sort               = var.sort
  virtual_ip_address = var.virtual_ip_address
  public_ip_address  = var.public_ip_address
}
