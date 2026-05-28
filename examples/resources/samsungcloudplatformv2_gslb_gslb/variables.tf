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
    algorithm   = "round_robin"
    description = "bbb"
    env_usage   = "PUBLIC"
    health_check = {
      health_check_interval      = 10
      health_check_probe_timeout = 6
      health_check_user_id       = "ENTER YOUR RESOURCE'S HEALTH_CHECK_USER_ID"
      health_check_user_password = "ENTER YOUR RESOURCE'S HEALTH_CHECK_USER_PASSWORD"
      protocol                   = "http"
      receive_string             = "asdsad11"
      send_string                = "asdsad2"
      service_port               = 40
      timeout                    = 30
    }
    name = "terraform.gslb.dev2.samsungsdscloud.com"
    resources = [{
      description = "string"
      destination = "1.1.1.1"
      region      = "KR-WEST-1"
      weight      = 40
      }, {
      description = "string"
      destination = "3.3.3.3"
      region      = "KR-WEST-2"
      weight      = 50
    }]
  }
}

variable "tag" {
  type = object({
    test_terraform_tag_key = string
  })
  default = {
    test_terraform_tag_key = "test_terraform_tag_value"
  }
}



