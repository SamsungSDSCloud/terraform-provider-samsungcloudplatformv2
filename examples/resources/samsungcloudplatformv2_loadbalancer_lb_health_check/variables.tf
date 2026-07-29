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
    description           = "123456789012345678901234567890123456789012345678902"
    health_check_count    = "5"
    health_check_interval = "30"
    health_check_port     = "100"
    health_check_timeout  = "20"
    health_check_url      = "/test-terraform"
    http_method           = "POST"
    name                  = "test-terraform-http"
    protocol              = "HTTP"
    request_data          = "username=test&password=123456789"
    response_code         = "200"
    subnet_id             = "ENTER YOUR RESOURCE'S SUBNET_ID"
    tags = {
      key  = "value"
      key1 = "value1"
    }
    vpc_id = "ENTER YOUR RESOURCE'S VPC_ID"
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
    description           = "123456789012345678901234567890123456789012345678902"
    health_check_count    = "5"
    health_check_interval = "30"
    health_check_port     = null
    health_check_timeout  = "20"
    health_check_url      = "/test-terraform"
    http_method           = "POST"
    name                  = "test-terraform-http"
    protocol              = "HTTPS"
    request_data          = "username=test&password=123456789"
    response_code         = "200"
    subnet_id             = "ENTER YOUR RESOURCE'S SUBNET_ID"
    tags = {
      key  = "value"
      key1 = "value1"
    }
    vpc_id = "ENTER YOUR RESOURCE'S VPC_ID"
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
    description           = "123456789012345678901234567890123456789012345678902"
    health_check_count    = "5"
    health_check_interval = "30"
    health_check_port     = "100"
    health_check_timeout  = "20"
    name                  = "test-terraform-tcp"
    protocol              = "TCP"
    subnet_id             = "ENTER YOUR RESOURCE'S SUBNET_ID"
    tags = {
      key  = "value"
      key1 = "value1"
    }
    vpc_id = "ENTER YOUR RESOURCE'S VPC_ID"
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
    description           = "this is a test"
    health_check_count    = "7"
    health_check_interval = "20"
    health_check_port     = "200"
    health_check_timeout  = "10"
    health_check_url      = null
    http_method           = null
    protocol              = "TCP"
    request_data          = null
    response_code         = null
  }
}


