variable "gslb" {
  type = object({
    algorithm   = string
    description = optional(string)
    env_usage   = string
    health_check = optional(object({
      health_check_interval      = optional(number)
      health_check_probe_timeout = optional(number)
      health_check_user_id       = optional(string)
      health_check_user_password = optional(string)
      protocol                   = string
      receive_string             = optional(string)
      send_string                = optional(string)
      service_port               = optional(number)
      timeout                    = optional(number)
    }))
    name = string
    resources = list(object({
      description = optional(string)
      destination = optional(string)
      region      = optional(string)
      weight      = optional(number)
    }))
  })
  default = {
    algorithm    = ""
    description  = null
    env_usage    = ""
    health_check = null
    name         = ""
    resources = [{
      description = null
      destination = null
      region      = null
      weight      = null
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


