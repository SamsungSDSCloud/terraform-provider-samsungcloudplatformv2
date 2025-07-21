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
    description           = "aa"
    https_redirection     = false
    insert_client_ip      = null
    loadbalancer_id       = "8a463aa4b1dc4f279c3f53b94dc45e74"
    name                  = "terraform-http"
    persistence           = null
    protocol              = "HTTP"
    response_timeout      = 60
    server_group_id       = null
    service_port          = 34123
    session_duration_time = 120
    ssl_certificate = null
    url_handler = null
    url_redirection = [
      {
        "url_pattern" : "/url",
        "redirect_url_pattern" : "/redirect"
      },
      {
        "url_pattern" : "/url2",
        "redirect_url_pattern" : "/redirect2"
      }
    ]
    x_forwarded_proto = false
    x_forwarded_port  = false
    x_forwarded_for   = false
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
    description           = "description info"
    https_redirection     = null
    insert_client_ip      = null
    loadbalancer_id       = "8a463aa4b1dc4f279c3f53b94dc45e74"
    name                  = "terraform-https"
    persistence           = "source-ip"
    protocol              = "HTTPS"
    response_timeout      = 60
    server_group_id       = null
    service_port          = 34124
    session_duration_time = 120
    ssl_certificate = {
      "client_cert_id": "8a463aa4b1dc4f279c3f53b94dc45e74",
      "client_cert_level": "HIGH",
      "server_cert_id": "8a463aa4b1dc4f279c3f53b94dc45e74",
      "server_cert_level": "HIGH"
    }
    url_handler = [
      {
        "url_pattern": "/",
        "server_group_id": null
      },
      {
        "url_pattern": "/url",
        "server_group_id": null
      }
    ]
    url_redirection = null
    x_forwarded_proto = false
    x_forwarded_port  = false
    x_forwarded_for   = false
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
    description           = "description info"
    https_redirection     = null
    insert_client_ip      = null
    loadbalancer_id       = "8a463aa4b1dc4f279c3f53b94dc45e74"
    name                  = "terraform-udp"
    persistence           = null
    protocol              = "UDP"
    response_timeout      = null
    server_group_id       = null
    service_port          = 34124
    session_duration_time = 120
    ssl_certificate = null
    url_handler = null
    url_redirection = null
    x_forwarded_proto = null
    x_forwarded_port  = null
    x_forwarded_for   = null
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
    description           = "description info"
    https_redirection     = null
    insert_client_ip      = null
    loadbalancer_id       = "8a463aa4b1dc4f279c3f53b94dc45e74"
    name                  = "terraform-tcp"
    persistence           = "source-ip"
    protocol              = "TCP"
    response_timeout      = null
    server_group_id       = null
    service_port          = 34125
    session_duration_time = 120
    ssl_certificate = null
    url_handler = null
    url_redirection = null
    x_forwarded_proto = null
    x_forwarded_port  = null
    x_forwarded_for   = null
  }
}

