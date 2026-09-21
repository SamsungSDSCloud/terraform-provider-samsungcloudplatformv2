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


# resource "samsungcloudplatformv2_loadbalancer_lb_listener" "lblistener" {
#  lb_listener_create = var.lb_listener_https
# }

locals {
  lb_listener_full = {
    # --- Base fields (only in lb_listener_https) ---
    loadbalancer_id   = var.lb_listener_https.loadbalancer_id
    name              = var.lb_listener_https.name
    protocol          = var.lb_listener_https.protocol
    service_port      = var.lb_listener_https.service_port
    routing_action    = var.lb_listener_https.routing_action
    tags              = var.lb_listener_https.tags

    # --- Resource fields
    description           = var.lb_listener_resource.description
    insert_client_ip      = var.lb_listener_resource.insert_client_ip
    persistence           = var.lb_listener_resource.persistence
    response_timeout      = var.lb_listener_resource.response_timeout
    session_duration_time = var.lb_listener_resource.session_duration_time
    idle_timeout          = var.lb_listener_resource.idle_timeout
    support_http2         = var.lb_listener_resource.support_http2
    x_forwarded_for       = var.lb_listener_resource.x_forwarded_for
    x_forwarded_port      = var.lb_listener_resource.x_forwarded_port
    x_forwarded_proto     = var.lb_listener_resource.x_forwarded_proto
    hsts_config           = var.lb_listener_resource.hsts_config

    # --- Rule fields
    condition_type    = var.lb_listener_rule.condition_type
    server_group_id   = var.lb_listener_rule.server_group_id
    url_redirection   = var.lb_listener_rule.url_redirection
    https_redirection = var.lb_listener_rule.https_redirection
    url_handler       = var.lb_listener_rule.url_handler

    # --- Certificate fields
    ssl_certificate = var.lb_listener_certificate.ssl_certificate
    sni_certificate = var.lb_listener_certificate.sni_certificate
  }
}

resource "samsungcloudplatformv2_loadbalancer_lb_listener" "lblistener" {
  lb_listener_create = local.lb_listener_full
}