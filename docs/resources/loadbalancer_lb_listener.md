---
page_title: "samsungcloudplatformv2_loadbalancer_lb_listener Resource - samsungcloudplatformv2"
subcategory: LB Listener
description: |-
  Lb Listener.
---

# samsungcloudplatformv2_loadbalancer_lb_listener (Resource)

Lb Listener.

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}


# locals {
#   lb_listeners = {
#     "listener1" = { lb_listener_create = var.lb_listener1 },
#     "listener2" = { lb_listener_create = var.lb_listener2 },
#     "listener3" = { lb_listener_create = var.lb_listener3 },
#     "listener4" = { lb_listener_create = var.lb_listener4 },
#   }
# }
#
# resource "samsungcloudplatformv2_loadbalancer_lb_listener" "lblistener" {
#   for_each           = local.lb_listeners
#   lb_listener_create = each.value.lb_listener_create
# }


resource "samsungcloudplatformv2_loadbalancer_lb_listener" "lblistener" {
  lb_listener_create = var.lb_listener_https
}


output "lb_listener" {
  value = samsungcloudplatformv2_loadbalancer_lb_listener.lblistener
}

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
    idle_timeout      = number
    hsts_max_age      = number
  })
  default = {
    condition_type        = "URL_PATH"
    description           = "aa2"
    hsts_max_age          = null
    https_redirection     = null
    idle_timeout          = null
    insert_client_ip      = null
    loadbalancer_id       = "ENTER YOUR RESOURCE'S LOADBALANCER_ID"
    name                  = "terraform-http-url-handler"
    persistence           = "source-ip"
    protocol              = "HTTP"
    response_timeout      = 60
    routing_action        = "LB_SERVER_GROUP"
    server_group_id       = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
    service_port          = 34124
    session_duration_time = 120
    ssl_certificate       = null
    url_handler = [{
      seq             = 0
      server_group_id = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
      url_pattern     = "default"
      }, {
      seq             = 1
      server_group_id = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
      url_pattern     = "/"
    }]
    url_redirection   = null
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
    idle_timeout      = number
    hsts_max_age      = number
  })
  default = {
    condition_type = "PROTOCOL_PORT"
    description    = "aa"
    hsts_max_age   = null
    https_redirection = {
      port          = "443"
      protocol      = "HTTPS"
      response_code = "301"
    }
    idle_timeout          = null
    insert_client_ip      = null
    loadbalancer_id       = "ENTER YOUR RESOURCE'S LOADBALANCER_ID"
    name                  = "terraform-http-protocol_port"
    persistence           = null
    protocol              = "HTTP"
    response_timeout      = 60
    routing_action        = "URL_REDIRECT"
    server_group_id       = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
    service_port          = 34123
    session_duration_time = 120
    ssl_certificate       = null
    url_handler           = null
    url_redirection       = null
    x_forwarded_for       = false
    x_forwarded_port      = false
    x_forwarded_proto     = false
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
    idle_timeout      = number
    hsts_max_age      = number
  })
  default = {
    condition_type        = "HOST_HEADER"
    description           = "aa2"
    hsts_max_age          = null
    https_redirection     = null
    idle_timeout          = 600
    insert_client_ip      = null
    loadbalancer_id       = "ENTER YOUR RESOURCE'S LOADBALANCER_ID"
    name                  = "terraform-http-redirect"
    persistence           = null
    protocol              = "HTTP"
    response_timeout      = null
    routing_action        = "URL_REDIRECT"
    server_group_id       = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
    service_port          = 34124
    session_duration_time = null
    ssl_certificate       = null
    url_handler           = null
    url_redirection       = "http://redirect.com"
    x_forwarded_for       = false
    x_forwarded_port      = false
    x_forwarded_proto     = false
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
    idle_timeout      = number
    hsts_max_age      = number
  })
  default = {
    condition_type        = "URL_PATH"
    description           = "aa2"
    hsts_max_age          = null
    https_redirection     = null
    idle_timeout          = 240
    insert_client_ip      = null
    loadbalancer_id       = "ENTER YOUR RESOURCE'S LOADBALANCER_ID"
    name                  = "terraform-https"
    persistence           = "source-ip"
    protocol              = "HTTPS"
    response_timeout      = null
    routing_action        = "LB_SERVER_GROUP"
    server_group_id       = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
    service_port          = 34124
    session_duration_time = null
    ssl_certificate = {
      client_cert_id    = "ENTER YOUR RESOURCE'S CLIENT_CERT_ID"
      client_cert_level = "HIGH"
      server_cert_level = null
    }
    url_handler = [{
      seq             = 0
      server_group_id = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
      url_pattern     = "default"
      }, {
      seq             = 1
      server_group_id = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
      url_pattern     = "/"
    }]
    url_redirection   = null
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
    idle_timeout      = number
    hsts_max_age      = number
  })
  default = {
    condition_type        = null
    description           = "aa2"
    hsts_max_age          = null
    https_redirection     = null
    idle_timeout          = null
    insert_client_ip      = null
    loadbalancer_id       = "ENTER YOUR RESOURCE'S LOADBALANCER_ID"
    name                  = "terraform-udp"
    persistence           = null
    protocol              = "UDP"
    response_timeout      = null
    routing_action        = "LB_SERVER_GROUP"
    server_group_id       = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
    service_port          = 34124
    session_duration_time = 120
    ssl_certificate       = null
    url_handler           = null
    url_redirection       = null
    x_forwarded_for       = null
    x_forwarded_port      = null
    x_forwarded_proto     = null
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
    idle_timeout      = number
    hsts_max_age      = number
  })
  default = {
    condition_type        = null
    description           = "aa2"
    hsts_max_age          = null
    https_redirection     = null
    idle_timeout          = null
    insert_client_ip      = null
    loadbalancer_id       = "ENTER YOUR RESOURCE'S LOADBALANCER_ID"
    name                  = "terraform-tcp"
    persistence           = "source-ip"
    protocol              = "TCP"
    response_timeout      = null
    routing_action        = "LB_SERVER_GROUP"
    server_group_id       = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
    service_port          = 34125
    session_duration_time = 120
    ssl_certificate       = null
    url_handler           = null
    url_redirection       = null
    x_forwarded_for       = null
    x_forwarded_port      = null
    x_forwarded_proto     = null
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
    idle_timeout      = number
    hsts_max_age      = number
  })
  default = {
    condition_type        = null
    description           = "this listener is made from terraform"
    hsts_max_age          = null
    https_redirection     = null
    idle_timeout          = null
    insert_client_ip      = null
    loadbalancer_id       = "ENTER YOUR RESOURCE'S LOADBALANCER_ID"
    name                  = "terraform-tls"
    persistence           = "source-ip"
    protocol              = "TLS"
    response_timeout      = null
    routing_action        = "LB_SERVER_GROUP"
    server_group_id       = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
    service_port          = 34122
    session_duration_time = 120
    sni_certificate       = null
    ssl_certificate = {
      client_cert_id    = "ENTER YOUR RESOURCE'S CLIENT_CERT_ID"
      client_cert_level = "HIGH"
      server_cert_level = "HIGH"
    }
    url_handler       = null
    url_redirection   = null
    x_forwarded_for   = null
    x_forwarded_port  = null
    x_forwarded_proto = null
  }
}

variable "update_sni_certificate" {
  type = object({
    sni_certificate = list(object({
      sni_cert_id = string
      domain_name = string
    }))
    x_forwarded_proto = bool
    x_forwarded_port  = bool
    x_forwarded_for   = bool
  })
  default = {
    sni_certificate = [{
      domain_name = "test"
      sni_cert_id = "ENTER YOUR RESOURCE'S SNI_CERT_ID"
    }]
    x_forwarded_for   = false
    x_forwarded_port  = false
    x_forwarded_proto = false
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `lb_listener_create` (Attributes) Create Lb Listener. (see [below for nested schema](#nestedatt--lb_listener_create))

### Read-Only

- `id` (String) Identifier of the resource.
- `lb_listener` (Attributes) A detail of Lb Listener. (see [below for nested schema](#nestedatt--lb_listener))

<a id="nestedatt--lb_listener_create"></a>
### Nested Schema for `lb_listener_create`

Optional:

- `condition_type` (String) ConditionType
- `description` (String) Description
- `hsts_max_age` (Number) HstsMaxAge
- `https_redirection` (Attributes) HttpsRedirection (see [below for nested schema](#nestedatt--lb_listener_create--https_redirection))
- `idle_timeout` (Number) IdleTimeout
- `insert_client_ip` (Boolean) InsertClientIp
- `loadbalancer_id` (String) LoadbalancerId
- `name` (String) Name
- `persistence` (String) Persistence
- `protocol` (String) Protocol
- `response_timeout` (Number) ResponseTimeout
- `routing_action` (String) RoutingAction
- `server_group_id` (String) ServerGroupId
- `service_port` (Number) ServicePort
- `session_duration_time` (Number) SessionDurationTime
- `sni_certificate` (Attributes List) SniCertificate (see [below for nested schema](#nestedatt--lb_listener_create--sni_certificate))
- `ssl_certificate` (Attributes) SslCertificate (see [below for nested schema](#nestedatt--lb_listener_create--ssl_certificate))
- `url_handler` (Attributes List) UrlHandler (see [below for nested schema](#nestedatt--lb_listener_create--url_handler))
- `url_redirection` (String) UrlRedirection
- `x_forwarded_for` (Boolean) XForwardedFor
- `x_forwarded_port` (Boolean) XForwardedPort
- `x_forwarded_proto` (Boolean) XForwardedProto

<a id="nestedatt--lb_listener_create--https_redirection"></a>
### Nested Schema for `lb_listener_create.https_redirection`

Optional:

- `port` (String) Port
- `protocol` (String) Protocol
- `response_code` (String) ResponseCode


<a id="nestedatt--lb_listener_create--sni_certificate"></a>
### Nested Schema for `lb_listener_create.sni_certificate`

Optional:

- `domain_name` (String) DomainName
- `sni_cert_id` (String) SniCertId


<a id="nestedatt--lb_listener_create--ssl_certificate"></a>
### Nested Schema for `lb_listener_create.ssl_certificate`

Optional:

- `client_cert_id` (String) ClientCertId
- `client_cert_level` (String) ClientCertLevel
- `server_cert_level` (String) ServerCertLevel


<a id="nestedatt--lb_listener_create--url_handler"></a>
### Nested Schema for `lb_listener_create.url_handler`

Optional:

- `seq` (Number) Seq
- `server_group_id` (String) ServerGroupId
- `url_pattern` (String) UrlPattern



<a id="nestedatt--lb_listener"></a>
### Nested Schema for `lb_listener`

Optional:

- `condition_type` (String) ConditionType
- `created_at` (String) created at
- `created_by` (String) created by
- `description` (String) Description
- `hsts_max_age` (Number) HstsMaxAge
- `https_redirection` (Attributes) HttpsRedirection (see [below for nested schema](#nestedatt--lb_listener--https_redirection))
- `idle_timeout` (Number) IdleTimeout
- `insert_client_ip` (Boolean) InsertClientIp
- `modified_at` (String) modified at
- `modified_by` (String) modified by
- `name` (String) Name
- `persistence` (String) Persistence
- `protocol` (String) Protocol
- `response_timeout` (Number) ResponseTimeout
- `routing_action` (String) RoutingAction
- `server_group_id` (String) ServerGroupId
- `server_group_name` (String) ServerGroupName
- `service_port` (Number) ServicePort
- `session_duration_time` (Number) SessionDurationTime
- `sni_certificate` (Attributes List) SniCertificate (see [below for nested schema](#nestedatt--lb_listener--sni_certificate))
- `ssl_certificate` (Attributes) SslCertificate (see [below for nested schema](#nestedatt--lb_listener--ssl_certificate))
- `url_handler` (Attributes List) UrlHandler (see [below for nested schema](#nestedatt--lb_listener--url_handler))
- `url_redirection` (String) UrlRedirection
- `x_forwarded_for` (Boolean) XForwardedFor
- `x_forwarded_port` (Boolean) XForwardedPort
- `x_forwarded_proto` (Boolean) XForwardedProto

Read-Only:

- `id` (String) id
- `state` (String) State

<a id="nestedatt--lb_listener--https_redirection"></a>
### Nested Schema for `lb_listener.https_redirection`

Optional:

- `port` (String) Port
- `protocol` (String) Protocol
- `response_code` (String) ResponseCode


<a id="nestedatt--lb_listener--sni_certificate"></a>
### Nested Schema for `lb_listener.sni_certificate`

Optional:

- `domain_name` (String) DomainName
- `not_after_dt` (String) NotAfterDt
- `sni_cert_id` (String) SniCertId


<a id="nestedatt--lb_listener--ssl_certificate"></a>
### Nested Schema for `lb_listener.ssl_certificate`

Optional:

- `client_cert_id` (String) ClientCertId
- `client_cert_level` (String) ClientCertLevel
- `server_cert_level` (String) ServerCertLevel


<a id="nestedatt--lb_listener--url_handler"></a>
### Nested Schema for `lb_listener.url_handler`

Optional:

- `seq` (Number) Seq
- `server_group_id` (String) ServerGroupId
- `url_pattern` (String) UrlPattern