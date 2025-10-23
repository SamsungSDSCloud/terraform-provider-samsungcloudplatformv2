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
    description        = ""
    lb_health_check_id = ""
    lb_method          = ""
    name               = ""
    protocol           = ""
    subnet_id          = ""
    tags               = null
    vpc_id             = ""
  }
}

variable "lb_server_group_modify" {
  type = object({
    lb_method          = string
    description        = string
    lb_health_check_id = string
  })
  default = {
    description        = ""
    lb_health_check_id = ""
    lb_method          = ""
  }
}

