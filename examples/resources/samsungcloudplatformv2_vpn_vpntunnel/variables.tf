variable "description" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
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
    dpd_retry_interval           = 0
    ike_version                  = 0
    peer_gateway_ip              = ""
    phase1_diffie_hellman_groups = [0]
    phase1_encryptions           = [""]
    phase1_life_time             = 0
    pre_shared_key               = ""
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
    perfect_forward_secrecy      = ""
    phase2_diffie_hellman_groups = [0]
    phase2_encryptions           = [""]
    phase2_life_time             = 0
    remote_subnets               = [""]
  }
}

variable "tags" {
  type    = map(string)
  default = null
}

variable "vpn_gateway_id" {
  type    = string
  default = ""
}

