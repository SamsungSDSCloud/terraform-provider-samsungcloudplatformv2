---
page_title: "samsungcloudplatformv2_loadbalancer_lb_listener Resource - samsungcloudplatformv2"
subcategory: LB Listener
description: |-
  LB Listener resource for managing listener configurations.
---

# samsungcloudplatformv2_loadbalancer_lb_listener (Resource)

LB Listener resource for managing listener configurations.

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

- `lb_listener_create` (Attributes) Parameters for creating a new LB Listener. (see [below for nested schema](#nestedatt--lb_listener_create))

### Read-Only

- `id` (String) Identifier of the resource.
  - example: YOUR RESOURCE'S ID
- `lb_listener` (Attributes) Details of the LB Listener. (see [below for nested schema](#nestedatt--lb_listener))

<a id="nestedatt--lb_listener_create"></a>
### Nested Schema for `lb_listener_create`

Optional:

- `condition_type` (String) The condition type for routing. 'URL_PATH' or 'HOST_HEADER' for URL handler.
  - example : URL_PATH
- `description` (String) Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.
  - example : LB Listener for web traffic
  - maxLength : 255
- `hsts_max_age` (Number) HSTS max age in seconds.
  - example : 31536000
- `https_redirection` (Attributes) HTTPS redirection configuration. Only for HTTP protocol listeners. (see [below for nested schema](#nestedatt--lb_listener_create--https_redirection))
- `idle_timeout` (Number) The idle timeout in seconds. Only applicable for L7 protocols (HTTP, HTTPS).
  - example : 60
  - minimum : 60
  - maximum : 3600
- `insert_client_ip` (Boolean) Whether to insert client IP in the header using Proxy Protocol v1.
  - example : true
- `loadbalancer_id` (String) The LoadBalancer ID associated with the listener.
  - example: YOUR RESOURCE'S LOADBALANCER_ID
- `name` (String) The name of the LB Listener.
  - example : Listener01
  - minLength : 1
  - maxLength : 63
- `persistence` (String) Session persistence configuration.
  - example : source-ip
  - pattern : source-ip | cookie
- `protocol` (String) The protocol used for the listener.
  - example : HTTP
  - pattern : TCP | UDP | HTTP | HTTPS | TLS | TCP_PROXY
- `response_timeout` (Number) The response timeout in seconds. Only for L7 protocols (HTTP/HTTPS).
  - example : 30
  - minimum : 1
  - maximum : 120
- `routing_action` (String) The routing action type. 'LB_SERVER_GROUP' for URL handler routing, 'URL_REDIRECT' for HTTPS/URL redirection.
  - example : LB_SERVER_GROUP
- `server_group_id` (String) The ID of the server group associated with the listener. Required for TCP, UDP, and TLS protocols.
  - example: YOUR RESOURCE'S SERVER_GROUP_ID
- `service_port` (Number) The service port number for the listener.
  - example : 80
  - minimum : 1
  - maximum : 65534
- `session_duration_time` (Number) The session duration time in seconds. Required for L4 protocols (TCP/UDP).
  - example : 60
  - L7(HTTP/HTTPS) : minimum 1, maximum 120
  - L4(TCP/TLS) : minimum 60, maximum 3600 (60-second increments)
  - UDP : minimum 60, maximum 180 (60-second increments)
- `sni_certificate` (Attributes List) SNI certificate configuration for multiple domains. (see [below for nested schema](#nestedatt--lb_listener_create--sni_certificate))
- `ssl_certificate` (Attributes) SSL certificate configuration for the listener. (see [below for nested schema](#nestedatt--lb_listener_create--ssl_certificate))
- `url_handler` (Attributes List) URL handler configuration for routing. Only for L7 protocols (HTTP/HTTPS). (see [below for nested schema](#nestedatt--lb_listener_create--url_handler))
- `url_redirection` (String) URL redirection configuration. Must start with http:// or https://.
  - example : https://example.com
  - maxLength : 255
  - pattern (HTTP) : ^http://[A-Za-z0-9:/.\-?=&#]{0,248}$
  - pattern (HTTPS) : ^https://[A-Za-z0-9:/.\-?=&#]{0,247}$
- `x_forwarded_for` (Boolean) X-Forwarded-For header configuration.
  - example : true
- `x_forwarded_port` (Boolean) X-Forwarded-Port header configuration.
  - example : true
- `x_forwarded_proto` (Boolean) X-Forwarded-Proto header configuration.
  - example : true

<a id="nestedatt--lb_listener_create--https_redirection"></a>
### Nested Schema for `lb_listener_create.https_redirection`

Optional:

- `port` (String) The port number to redirect to.
  - example : 443
  - minimum : 1
  - maximum : 65534
- `protocol` (String) The protocol to redirect to.
  - example : HTTPS
- `response_code` (String) The HTTP response code for redirection.
  - example : 301


<a id="nestedatt--lb_listener_create--sni_certificate"></a>
### Nested Schema for `lb_listener_create.sni_certificate`

Optional:

- `domain_name` (String) The domain name for SNI certificate. Must be unique within the listener.
  - example : example.com
  - minLength : 1
  - maxLength : 63
  - pattern : ^[a-zA-Z0-9](?:[a-zA-Z0-9.-]{0,61}[a-zA-Z0-9])?$
- `sni_cert_id` (String) The SNI certificate ID.
  - example: YOUR RESOURCE'S SNI_CERT_ID


<a id="nestedatt--lb_listener_create--ssl_certificate"></a>
### Nested Schema for `lb_listener_create.ssl_certificate`

Optional:

- `client_cert_id` (String) The client certificate ID.
  - example: YOUR RESOURCE'S CLIENT_CERT_ID
- `client_cert_level` (String) The client certificate validation level.
  - example : NORMAL
  - pattern : LOW | NORMAL | HIGH
- `server_cert_level` (String) The server certificate validation level.
  - example : NORMAL
  - pattern : LOW | NORMAL | HIGH


<a id="nestedatt--lb_listener_create--url_handler"></a>
### Nested Schema for `lb_listener_create.url_handler`

Optional:

- `seq` (Number) The sequence number for routing priority. 0 is reserved for default rule.
  - example : 1
- `server_group_id` (String) The ID of the server group to route traffic to when the URL pattern matches.
  - example: YOUR RESOURCE'S SERVER_GROUP_ID
- `url_pattern` (String) The URL pattern for routing.
  - example : /api/v1
  - minLength : 1
  - maxLength : 63



<a id="nestedatt--lb_listener"></a>
### Nested Schema for `lb_listener`

Optional:

- `condition_type` (String) The condition type for routing. 'URL_PATH' or 'HOST_HEADER' for URL handler.
  - example : URL_PATH
- `created_at` (String) The timestamp when the resource was created, in ISO 8601 format.
  - example : 2024-05-17T00:23:17Z
- `created_by` (String) The user id that created the resource.
  - example: YOUR RESOURCE'S CREATED_BY
- `description` (String) Enter a brief explanation or note about this resource. This helps identify the purpose or usage of the resource.
  - example : LB Listener for web traffic
  - maxLength : 255
- `hsts_max_age` (Number) HSTS max age in seconds.
  - example : 31536000
- `https_redirection` (Attributes) HTTPS redirection configuration. Only for HTTP protocol listeners. (see [below for nested schema](#nestedatt--lb_listener--https_redirection))
- `idle_timeout` (Number) The idle timeout in seconds. Only applicable for L7 protocols (HTTP, HTTPS).
  - example : 60
  - minimum : 60
  - maximum : 3600
- `insert_client_ip` (Boolean) Whether to insert client IP in the header using Proxy Protocol v1.
  - example : true
- `modified_at` (String) The timestamp when the resource was last modified, in ISO 8601 format.
  - example : 2024-05-17T00:23:17Z
- `modified_by` (String) The user id that last modified the resource.
  - example: YOUR RESOURCE'S MODIFIED_BY
- `name` (String) The name of the LB Listener.
  - example : listener-web-01
  - minLength : 3
  - maxLength : 63
  - pattern : ^[a-zA-Z0-9-_]*$
- `persistence` (String) Session persistence configuration.
  - example : source-ip
  - pattern : source-ip | cookie
- `protocol` (String) The protocol used for the listener.
  - example : TCP
  - pattern : TCP | UDP | HTTP | HTTPS | TLS | TCP_PROXY
- `response_timeout` (Number) The response timeout in seconds. Only for L7 protocols (HTTP/HTTPS).
  - example : 30
  - minimum : 1
  - maximum : 120
- `routing_action` (String) The routing action type. 'LB_SERVER_GROUP' for URL handler routing, 'URL_REDIRECT' for HTTPS/URL redirection.
  - example : LB_SERVER_GROUP
- `server_group_id` (String) The ID of the server group associated with the listener. Required for TCP, UDP, and TLS protocols.
  - example: YOUR RESOURCE'S SERVER_GROUP_ID
- `server_group_name` (String) The server group name for the listener.
  - example : ServerGroup01
- `service_port` (Number) The service port number for the listener.
  - example : 80
  - minimum : 1
  - maximum : 65534
- `session_duration_time` (Number) The session duration time in seconds. Required for L4 protocols (TCP/UDP).
  - example : 60
  - L7(HTTP/HTTPS) : minimum 1, maximum 120
  - L4(TCP/TLS) : minimum 60, maximum 3600 (60-second increments)
  - UDP : minimum 60, maximum 180 (60-second increments)
- `sni_certificate` (Attributes List) SNI certificate configuration for multiple domains. (see [below for nested schema](#nestedatt--lb_listener--sni_certificate))
- `ssl_certificate` (Attributes) SSL certificate configuration for the listener. (see [below for nested schema](#nestedatt--lb_listener--ssl_certificate))
- `url_handler` (Attributes List) URL handler configuration for routing. Only for L7 protocols (HTTP/HTTPS). (see [below for nested schema](#nestedatt--lb_listener--url_handler))
- `url_redirection` (String) URL redirection configuration. Must start with http:// or https://.
  - example : https://example.com
  - maxLength : 255
  - pattern (HTTP) : ^http://[A-Za-z0-9:/.\-?=&#]{0,248}$
  - pattern (HTTPS) : ^https://[A-Za-z0-9:/.\-?=&#]{0,247}$
- `x_forwarded_for` (Boolean) X-Forwarded-For header configuration.
  - example : true
- `x_forwarded_port` (Boolean) X-Forwarded-Port header configuration.
  - example : true
- `x_forwarded_proto` (Boolean) X-Forwarded-Proto header configuration.
  - example : true

Read-Only:

- `id` (String) The unique identifier.
  - example: YOUR RESOURCE'S ID
- `state` (String) The current state of the LB Listener.
  - example : ACTIVE
  - pattern : CREATING | ACTIVE | DELETING | ERROR

<a id="nestedatt--lb_listener--https_redirection"></a>
### Nested Schema for `lb_listener.https_redirection`

Optional:

- `port` (String) The port number to redirect to.
  - example : 443
  - minimum : 1
  - maximum : 65534
- `protocol` (String) The protocol to redirect to.
  - example : HTTPS
- `response_code` (String) The HTTP response code for redirection.
  - example : 301


<a id="nestedatt--lb_listener--sni_certificate"></a>
### Nested Schema for `lb_listener.sni_certificate`

Optional:

- `domain_name` (String) The domain name for SNI certificate. Must be unique within the listener.
  - example : example.com
  - minLength : 1
  - maxLength : 63
  - pattern : ^[a-zA-Z0-9](?:[a-zA-Z0-9.-]{0,61}[a-zA-Z0-9])?$
- `not_after_dt` (String) The expiration date and time of the certificate. Read-only.
  - example : 2024-12-31T23:59:59Z
- `sni_cert_id` (String) The SNI certificate ID.
  - example: YOUR RESOURCE'S SNI_CERT_ID


<a id="nestedatt--lb_listener--ssl_certificate"></a>
### Nested Schema for `lb_listener.ssl_certificate`

Optional:

- `client_cert_id` (String) The client certificate ID.
  - example: YOUR RESOURCE'S CLIENT_CERT_ID
- `client_cert_level` (String) The client certificate validation level.
  - example : NORMAL
  - pattern : LOW | NORMAL | HIGH
- `server_cert_level` (String) The server certificate validation level.
  - example : NORMAL
  - pattern : LOW | NORMAL | HIGH


<a id="nestedatt--lb_listener--url_handler"></a>
### Nested Schema for `lb_listener.url_handler`

Optional:

- `seq` (Number) The sequence number for routing priority. 0 is reserved for default rule.
  - example : 1
- `server_group_id` (String) The ID of the server group to route traffic to when the URL pattern matches.
  - example: YOUR RESOURCE'S SERVER_GROUP_ID
- `url_pattern` (String) The URL pattern for routing.
  - example : /api/v1
  - minLength : 1
  - maxLength : 63