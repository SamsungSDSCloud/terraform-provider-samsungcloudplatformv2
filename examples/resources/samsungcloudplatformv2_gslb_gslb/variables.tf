variable "gslb" {
  type = object({
    algorithm   = string
    description = string
    env_usage   = string
    health_check = object({
      health_check_interval      = number
      health_check_probe_timeout = number
      health_check_user_id       = string
      health_check_user_password = string
      protocol                   = string
      receive_string             = string
      send_string                = string
      service_port               = number
      timeout                    = number
    })
    name = string
    resources = list(object({
      description = string
      destination = string
      disabled    = bool
      region      = string
      weight      = number
    }))
  })
  default = {
    algorithm   = ""
    description = ""
    env_usage   = ""
    health_check = {
      health_check_interval      = 0
      health_check_probe_timeout = 0
      health_check_user_id       = ""
      health_check_user_password = ""
      protocol                   = ""
      receive_string             = ""
      send_string                = ""
      service_port               = 0
      timeout                    = 0
    }
    name = ""
    resources = [{
      description = ""
      destination = ""
      disabled    = false
      region      = ""
      weight      = 0
    }]
  }
}

variable "tag" {
  type = object({
    test_terraform_tag_key = string
  })
  default = {
    test_terraform_tag_key = ""
  }
}


