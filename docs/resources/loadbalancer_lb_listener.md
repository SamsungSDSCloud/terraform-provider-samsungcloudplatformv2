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


locals {
  lb_listeners = {
    "listener1" = { lb_listener_create = var.lb_listener1 },
    "listener2" = { lb_listener_create = var.lb_listener2 },
    "listener3" = { lb_listener_create = var.lb_listener3 },
    "listener4" = { lb_listener_create = var.lb_listener4 },
  }
}

resource "samsungcloudplatformv2_loadbalancer_lb_listener" "lblistener" {
  for_each           = local.lb_listeners
  lb_listener_create = each.value.lb_listener_create
}

# output "lb_listener" {
#   value = samsungcloudplatformv2_loadbalancer_lb_listener.lblistener
# }

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

- `description` (String) Description
- `https_redirection` (Boolean) HttpsRedirection
- `insert_client_ip` (Boolean) InsertClientIp
- `loadbalancer_id` (String) LoadbalancerId
- `name` (String) Name
- `persistence` (String) Persistence
- `protocol` (String) Protocol
- `response_timeout` (Number) ResponseTimeout
- `server_group_id` (String) ServerGroupId
- `service_port` (Number) ServicePort
- `session_duration_time` (Number) SessionDurationTime
- `ssl_certificate` (Attributes) SslCertificate (see [below for nested schema](#nestedatt--lb_listener_create--ssl_certificate))
- `url_handler` (Attributes List) UrlHandler (see [below for nested schema](#nestedatt--lb_listener_create--url_handler))
- `url_redirection` (Attributes List) UrlRedirection (see [below for nested schema](#nestedatt--lb_listener_create--url_redirection))
- `x_forwarded_for` (Boolean) XForwardedFor
- `x_forwarded_port` (Boolean) XForwardedPort
- `x_forwarded_proto` (Boolean) XForwardedProto

<a id="nestedatt--lb_listener_create--ssl_certificate"></a>
### Nested Schema for `lb_listener_create.ssl_certificate`

Optional:

- `client_cert_id` (String) ClientCertId
- `client_cert_level` (String) ClientCertLevel
- `server_cert_id` (String) ServerCertId
- `server_cert_level` (String) ServerCertLevel


<a id="nestedatt--lb_listener_create--url_handler"></a>
### Nested Schema for `lb_listener_create.url_handler`

Optional:

- `server_group_id` (String) ServerGroupId
- `url_pattern` (String) UrlPattern


<a id="nestedatt--lb_listener_create--url_redirection"></a>
### Nested Schema for `lb_listener_create.url_redirection`

Optional:

- `redirect_url_pattern` (String) RedirectUrlPattern
- `url_pattern` (String) UrlPattern



<a id="nestedatt--lb_listener"></a>
### Nested Schema for `lb_listener`

Optional:

- `created_at` (String) created at
- `created_by` (String) created by
- `description` (String) Description
- `https_redirection` (Boolean) HttpsRedirection
- `insert_client_ip` (Boolean) InsertClientIp
- `modified_at` (String) modified at
- `modified_by` (String) modified by
- `name` (String) Name
- `persistence` (String) Persistence
- `protocol` (String) Protocol
- `response_timeout` (Number) ResponseTimeout
- `server_group_id` (String) ServerGroupId
- `server_group_name` (String) ServerGroupName
- `service_port` (Number) ServicePort
- `session_duration_time` (Number) SessionDurationTime
- `ssl_certificate` (Attributes) SslCertificate (see [below for nested schema](#nestedatt--lb_listener--ssl_certificate))
- `url_handler` (Attributes List) UrlHandler (see [below for nested schema](#nestedatt--lb_listener--url_handler))
- `url_redirection` (Attributes List) UrlRedirection (see [below for nested schema](#nestedatt--lb_listener--url_redirection))
- `x_forwarded_for` (Boolean) XForwardedFor
- `x_forwarded_port` (Boolean) XForwardedPort
- `x_forwarded_proto` (Boolean) XForwardedProto

Read-Only:

- `id` (String) id
- `state` (String) State

<a id="nestedatt--lb_listener--ssl_certificate"></a>
### Nested Schema for `lb_listener.ssl_certificate`

Optional:

- `client_cert_id` (String) ClientCertId
- `client_cert_level` (String) ClientCertLevel
- `server_cert_id` (String) ServerCertId
- `server_cert_level` (String) ServerCertLevel


<a id="nestedatt--lb_listener--url_handler"></a>
### Nested Schema for `lb_listener.url_handler`

Optional:

- `server_group_id` (String) ServerGroupId
- `url_pattern` (String) UrlPattern


<a id="nestedatt--lb_listener--url_redirection"></a>
### Nested Schema for `lb_listener.url_redirection`

Optional:

- `redirect_url_pattern` (String) RedirectUrlPattern
- `url_pattern` (String) UrlPattern