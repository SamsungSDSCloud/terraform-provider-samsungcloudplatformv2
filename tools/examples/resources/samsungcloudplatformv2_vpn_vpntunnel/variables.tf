variable "description" {
  type    = string
  default = "update 0912test123"
}

variable "name" {
  type    = string
  default = "test06162026"
}

variable "phase1" {
  type = object({
    dpd_retry_interval           = number
    ike_version                  = number
    peer_gateway_ip              = string
    phase1_diffie_hellman_groups = list(number)
    phase1_encryptions           = list(string)
    phase1_life_time             = number
    pre_shared_key               = string
  })
  default = {
    dpd_retry_interval           = 60
    ike_version                  = 1
    peer_gateway_ip              = "192.169.1.6"
    phase1_diffie_hellman_groups = [20, 21]
    phase1_encryptions           = ["des-sha1"]
    phase1_life_time             = 86400
    pre_shared_key               = "test1234578"
  }
}

variable "phase2" {
  type = object({
    perfect_forward_secrecy      = string
    phase2_diffie_hellman_groups = list(number)
    phase2_encryptions           = list(string)
    phase2_life_time             = number
    remote_subnets               = list(string)
  })
  default = {
    perfect_forward_secrecy      = "ENABLE"
    phase2_diffie_hellman_groups = [20, 21]
    phase2_encryptions           = ["des-sha1"]
    phase2_life_time             = 43200
    remote_subnets               = ["5.6.0.0/22", "5.7.0.0/24"]
  }
}

variable "tags" {
  type = map(string)
  default = {
    key1 = "222 updated"
  }
}

variable "vpn_gateway_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPN_GATEWAY_ID"
}



