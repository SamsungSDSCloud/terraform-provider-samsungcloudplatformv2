variable "lb_server_group_id" {
  type    = string
  default = ""
}

variable "lb_member" {
  type = object({
    name          = string
    member_ip     = string
    member_port   = string
    member_weight = string
    object_type   = string
    object_id     = string
  })
  default = {
    member_ip     = ""
    member_port   = ""
    member_weight = ""
    name          = ""
    object_id     = ""
    object_type   = ""
  }
}

variable "lb_member_modify" {
  type = object({
    member_port   = string
    member_weight = string
    #     member_state= string
  })
  default = {
    member_port   = ""
    member_weight = ""
  }
}

