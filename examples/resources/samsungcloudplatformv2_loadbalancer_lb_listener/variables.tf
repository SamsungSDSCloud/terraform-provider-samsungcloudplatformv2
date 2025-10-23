variable "lb_listener1" {
  type = object({
    description           = string
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
      server_cert_level = string
    })
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
      seq             = number
    }))
    https_redirection = object({
      protocol      = string
      port          = string
      response_code = string
    })
    url_redirection   = string
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
    routing_action    = string
    condition_type    = string
  })
  default = {
    condition_type = ""
    description    = ""
    https_redirection = {
      port          = ""
      protocol      = ""
      response_code = ""
    }
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    routing_action        = ""
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_level = ""
    }
    url_handler = [{
      seq             = 0
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection   = ""
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}

variable "lb_listener2" {
  type = object({
    description           = string
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
      server_cert_level = string
    })
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
      seq             = number
    }))
    https_redirection = object({
      protocol      = string
      port          = string
      response_code = string
    })
    url_redirection   = string
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
    routing_action    = string
    condition_type    = string
  })
  default = {
    condition_type = ""
    description    = ""
    https_redirection = {
      port          = ""
      protocol      = ""
      response_code = ""
    }
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    routing_action        = ""
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_level = ""
    }
    url_handler = [{
      seq             = 0
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection   = ""
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}

variable "lb_listener3" {
  type = object({
    description           = string
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
      server_cert_level = string
    })
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
      seq             = number
    }))
    https_redirection = object({
      protocol      = string
      port          = string
      response_code = string
    })
    url_redirection   = string
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
    routing_action    = string
    condition_type    = string
  })
  default = {
    condition_type = ""
    description    = ""
    https_redirection = {
      port          = ""
      protocol      = ""
      response_code = ""
    }
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    routing_action        = ""
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_level = ""
    }
    url_handler = [{
      seq             = 0
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection   = ""
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}

variable "lb_listener_https" {
  type = object({
    description           = string
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
      server_cert_level = string
    })
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
      seq             = number
    }))
    https_redirection = object({
      protocol      = string
      port          = string
      response_code = string
    })
    url_redirection   = string
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
    routing_action    = string
    condition_type    = string
  })
  default = {
    condition_type = ""
    description    = ""
    https_redirection = {
      port          = ""
      protocol      = ""
      response_code = ""
    }
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    routing_action        = ""
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_level = ""
    }
    url_handler = [{
      seq             = 0
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection   = ""
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}

variable "lb_listener_udp" {
  type = object({
    description           = string
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
      server_cert_level = string
    })
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
      seq             = number
    }))
    https_redirection = object({
      protocol      = string
      port          = string
      response_code = string
    })
    url_redirection   = string
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
    routing_action    = string
    condition_type    = string
  })
  default = {
    condition_type = ""
    description    = ""
    https_redirection = {
      port          = ""
      protocol      = ""
      response_code = ""
    }
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    routing_action        = ""
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_level = ""
    }
    url_handler = [{
      seq             = 0
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection   = ""
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}

variable "lb_listener_tcp" {
  type = object({
    description           = string
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
      server_cert_level = string
    })
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
      seq             = number
    }))
    https_redirection = object({
      protocol      = string
      port          = string
      response_code = string
    })
    url_redirection   = string
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
    routing_action    = string
    condition_type    = string
  })
  default = {
    condition_type = ""
    description    = ""
    https_redirection = {
      port          = ""
      protocol      = ""
      response_code = ""
    }
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    routing_action        = ""
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_level = ""
    }
    url_handler = [{
      seq             = 0
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection   = ""
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}

variable "lb_listener_tls" {
  type = object({
    description           = string
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
      server_cert_level = string
    })
    sni_certificate = list(object({
      sni_cert_id = string
      domain_name = string
    }))
    url_handler = list(object({
      url_pattern     = string
      server_group_id = string
      seq             = number
    }))
    https_redirection = object({
      protocol      = string
      port          = string
      response_code = string
    })
    url_redirection   = string
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
    routing_action    = string
    condition_type    = string
  })
  default = {
    condition_type = ""
    description    = ""
    https_redirection = {
      port          = ""
      protocol      = ""
      response_code = ""
    }
    insert_client_ip      = false
    loadbalancer_id       = ""
    name                  = ""
    persistence           = ""
    protocol              = ""
    response_timeout      = 0
    routing_action        = ""
    server_group_id       = ""
    service_port          = 0
    session_duration_time = 0
    sni_certificate = [{
      domain_name = ""
      sni_cert_id = ""
    }]
    ssl_certificate = {
      client_cert_id    = ""
      client_cert_level = ""
      server_cert_level = ""
    }
    url_handler = [{
      seq             = 0
      server_group_id = ""
      url_pattern     = ""
    }]
    url_redirection   = ""
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}

