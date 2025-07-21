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
    algorithm   = "round_robin"
    description = "description info"
    env_usage   = "PUBLIC"
    health_check = {
      health_check_interval      = 5
      timeout                    = 7
      health_check_probe_timeout = 6
      health_check_user_id       = null
      health_check_user_password = null
      protocol                   = "http"
      receive_string             = "receivestring"
      send_string                = "sendstring"
      service_port               = 40
    }
    name = "terraform.gslb.s.samsungsdscloud.com"
    resources = [
      {
        description = "string"
        destination = "10.10.10.10"
        disabled    = true
        region      = "KR-WEST-1"
        weight      = 40
      },
      {
        "description" : "string",
        "destination" : "20.20.20.20",
        "disabled" : true,
        "region" : "KR-WEST-2",
        "weight" : 50
      }
    ]
  }
}

variable "tag" {
  type = object({
    terraform_tag_key = string
  })
  default = {
    terraform_tag_key = "terraform_tag_value"
  }
}
