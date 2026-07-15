provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_internet_gateway" "internetgateway" {
  vpc_id = var.igw_vpc_id
  type = var.igw_type
  description = var.igw_description
}
