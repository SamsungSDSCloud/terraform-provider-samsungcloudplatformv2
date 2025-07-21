variable "vpn_gateway_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "name" {
  type = string
  default = "vpntunnelName"
}

variable "phase1" {
  type = object({
    ike_version           = number
    diffie_hellman_groups = list(number)
    encryptions           = list(string)
    dpd_retry_interval    = number
    life_time             = number
    peer_gateway_ip       = string
    pre_shared_key        = string
  })
  default = {
    ike_version           = 1
    diffie_hellman_groups = [29, 30]
    encryptions           = ["des-sha1"]
    dpd_retry_interval    = 120
    life_time             = 3600
    peer_gateway_ip       = "203.0.113.2"
    pre_shared_key        = "shared012"
  }
}

variable "phase2" {
  type = object({
    diffie_hellman_groups = list(number)
    encryptions           = list(string)
    life_time             = number
    perfect_forward_secrecy = string
    remote_subnet         = string
  })
  default = {
    diffie_hellman_groups = [ 20, 21 ]
    encryptions           = ["des-sha1"]
    life_time             = 43200
    perfect_forward_secrecy = "ENABLE"
    remote_subnet         = "30.10.0.0/24"
  }
}

variable "description" {
  type = string
  default = "description info"
}

variable "tags" {
  type    = map(string)
  default = {
    "vpn_tag_key" = "vpn_tag_value"
  }
}