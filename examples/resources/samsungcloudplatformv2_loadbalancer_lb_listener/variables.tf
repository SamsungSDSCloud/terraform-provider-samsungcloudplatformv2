variable "lb_listener1" {
  type = object({
    description           = string
    https_redirection     = bool
    insert_client_ip      = bool
    loadbalancer_id       = string
    name                  = string
    persistence           = string
    protocol              = string
    response_timeout      = number
    server_group_id       = string
    service_port          = number
    session_duration_time = number
    ssl_certificate = object({
      client_cert_id    = string
      client_cert_level = string
      server_cert_id    = string
      server_cert_level = string
    })
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
    }))
    url_redirection = list(object({
      url_pattern          = string
      redirect_url_pattern = string
    }))
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
  })
  default = {
    description           = ""
    https_redirection     = false
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_id    = ""
      server_cert_level = ""
    }
    url_handler = [{
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection = [{
      redirect_url_pattern = ""
      url_pattern          = ""
    }]
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}

variable "lb_listener2" {
  type = object({
    description           = string
    https_redirection     = bool
    insert_client_ip      = bool
    loadbalancer_id       = string
    name                  = string
    persistence           = string
    protocol              = string
    response_timeout      = number
    server_group_id       = string
    service_port          = number
    session_duration_time = number
    ssl_certificate = object({
      client_cert_id    = string
      client_cert_level = string
      server_cert_id    = string
      server_cert_level = string
    })
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
    }))
    url_redirection = list(object({
      url_pattern          = string
      redirect_url_pattern = string
    }))
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
  })
  default = {
    description           = ""
    https_redirection     = false
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_id    = ""
      server_cert_level = ""
    }
    url_handler = [{
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection = [{
      redirect_url_pattern = ""
      url_pattern          = ""
    }]
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}

variable "lb_listener3" {
  type = object({
    description           = string
    https_redirection     = bool
    insert_client_ip      = bool
    loadbalancer_id       = string
    name                  = string
    persistence           = string
    protocol              = string
    response_timeout      = number
    server_group_id       = string
    service_port          = number
    session_duration_time = number
    ssl_certificate = object({
      client_cert_id    = string
      client_cert_level = string
      server_cert_id    = string
      server_cert_level = string
    })
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
    }))
    url_redirection = list(object({
      url_pattern          = string
      redirect_url_pattern = string
    }))
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
  })
  default = {
    description           = ""
    https_redirection     = false
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_id    = ""
      server_cert_level = ""
    }
    url_handler = [{
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection = [{
      redirect_url_pattern = ""
      url_pattern          = ""
    }]
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}

variable "lb_listener4" {
  type = object({
    description           = string
    https_redirection     = bool
    insert_client_ip      = bool
    loadbalancer_id       = string
    name                  = string
    persistence           = string
    protocol              = string
    response_timeout      = number
    server_group_id       = string
    service_port          = number
    session_duration_time = number
    ssl_certificate = object({
      client_cert_id    = string
      client_cert_level = string
      server_cert_id    = string
      server_cert_level = string
    })
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
    }))
    url_redirection = list(object({
      url_pattern          = string
      redirect_url_pattern = string
    }))
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
  })
  default = {
    description           = ""
    https_redirection     = false
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_id    = ""
      server_cert_level = ""
    }
    url_handler = [{
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection = [{
      redirect_url_pattern = ""
      url_pattern          = ""
    }]
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}



