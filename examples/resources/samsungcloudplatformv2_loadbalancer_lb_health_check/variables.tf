variable "lb_health_check_http" {
  type = object({
    name                  = string
    vpc_id                = string
    subnet_id             = string
    protocol              = string
    health_check_port     = string
    health_check_interval = string
    health_check_timeout  = string
    health_check_count    = string
    http_method           = string
    health_check_url      = string
    response_code         = string
    request_data          = string
    description           = string
    tags                  = map(string)
  })
  default = {
    description           = ""
    health_check_count    = ""
    health_check_interval = ""
    health_check_port     = ""
    health_check_timeout  = ""
    health_check_url      = ""
    http_method           = ""
    name                  = ""
    protocol              = ""
    request_data          = ""
    response_code         = ""
    subnet_id             = ""
    tags                  = null
    vpc_id                = ""
  }
}

variable "lb_health_check_https" {
  type = object({
    name                  = string
    vpc_id                = string
    subnet_id             = string
    protocol              = string
    health_check_port     = string
    health_check_interval = string
    health_check_timeout  = string
    health_check_count    = string
    http_method           = string
    health_check_url      = string
    response_code         = string
    request_data          = string
    description           = string
    tags                  = map(string)
  })
  default = {
    description           = ""
    health_check_count    = ""
    health_check_interval = ""
    health_check_port     = ""
    health_check_timeout  = ""
    health_check_url      = ""
    http_method           = ""
    name                  = ""
    protocol              = ""
    request_data          = ""
    response_code         = ""
    subnet_id             = ""
    tags                  = null
    vpc_id                = ""
  }
}

variable "lb_health_check_tcp" {
  type = object({
    name                  = string
    vpc_id                = string
    subnet_id             = string
    protocol              = string
    health_check_port     = string
    health_check_interval = string
    health_check_timeout  = string
    health_check_count    = string
    description           = string
    tags                  = map(string)
  })
  default = {
    description           = ""
    health_check_count    = ""
    health_check_interval = ""
    health_check_port     = ""
    health_check_timeout  = ""
    name                  = ""
    protocol              = ""
    subnet_id             = ""
    tags                  = null
    vpc_id                = ""
  }
}

variable "lb_health_check_modify" {
  type = object({
    protocol              = string
    health_check_port     = string
    health_check_interval = string
    health_check_timeout  = string
    health_check_count    = string
    http_method           = string
    health_check_url      = string
    response_code         = string
    request_data          = string
    description           = string
  })
  default = {
    description           = ""
    health_check_count    = ""
    health_check_interval = ""
    health_check_port     = ""
    health_check_timeout  = ""
    health_check_url      = ""
    http_method           = ""
    protocol              = ""
    request_data          = ""
    response_code         = ""
  }
}

