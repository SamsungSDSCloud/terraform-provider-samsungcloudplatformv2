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
    algorithm   = "ROUND_ROBIN"
    description = "ExampleDescriptionForGSLB"
    env_usage   = "PUBLIC"
    health_check = {
      health_check_interval      = 10
      health_check_probe_timeout = 6
      health_check_user_id       = "ENTER YOUR RESOURCE'S HEALTH_CHECK_USER_ID"
      health_check_user_password = "ENTER YOUR RESOURCE'S HEALTH_CHECK_USER_PASSWORD"
      protocol                   = "HTTP"
      receive_string             = "ExampleReceiveString1"
      send_string                = "ExampleSendString1"
      service_port               = 80
      timeout                    = 30
    }
    name = "example.gslb.e.samsungsdscloud.com"
    resources = [{
      description = "ExampleResource1"
      destination = "1.1.1.1"
      region      = "KR-WEST-1"
      weight      = 40
      }, {
      description = "ExampleResource2"
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



