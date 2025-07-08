variable "vpn_gateway_id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
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
    diffie_hellman_groups = [0]
    dpd_retry_interval    = 0
    encryptions           = [""]
    ike_version           = 0
    life_time             = 0
    peer_gateway_ip       = ""
    pre_shared_key        = ""
  }
}

variable "phase2" {
  type = object({
    diffie_hellman_groups   = list(number)
    encryptions             = list(string)
    life_time               = number
    perfect_forward_secrecy = string
    remote_subnet           = string
  })
  default = {
    diffie_hellman_groups   = [0]
    encryptions             = [""]
    life_time               = 0
    perfect_forward_secrecy = ""
    remote_subnet           = ""
  }
}

variable "description" {
  type    = string
  default = ""
}

variable "tags" {
  type    = map(string)
  default = null
}

