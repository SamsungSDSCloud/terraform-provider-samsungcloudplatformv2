variable "lb_server_group" {
  type = object({
    name               = string
    vpc_id             = string
    subnet_id          = string
    protocol           = string
    lb_method          = string
    description        = string
    lb_health_check_id = string
    tags               = map(string)
  })
  default = {
    description        = "123456789012345678901234567890123456789012345678902"
    lb_health_check_id = "ENTER YOUR RESOURCE'S LB_HEALTH_CHECK_ID"
    lb_method          = "ROUND_ROBIN"
    name               = "test-terraform"
    protocol           = "TCP"
    subnet_id          = "ENTER YOUR RESOURCE'S SUBNET_ID"
    tags = {
      key  = "value"
      key1 = "value1"
    }
    vpc_id = "ENTER YOUR RESOURCE'S VPC_ID"
  }
}

variable "lb_server_group_modify" {
  type = object({
    lb_method          = string
    description        = string
    lb_health_check_id = string
  })
  default = {
    description        = "987654321"
    lb_health_check_id = "ENTER YOUR RESOURCE'S LB_HEALTH_CHECK_ID"
    lb_method          = "LEAST_CONNECTION"
  }
}


