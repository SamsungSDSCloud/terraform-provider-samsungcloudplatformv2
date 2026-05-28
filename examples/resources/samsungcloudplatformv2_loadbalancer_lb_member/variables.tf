variable "lb_server_group_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S LB_SERVER_GROUP_ID"
}

variable "lb_member" {
  type = object({
    name          = string
    member_ip     = string
    member_port   = string
    member_weight = string
    member_state  = string
    object_type   = string
    object_id     = string
  })
  default = {
    member_ip     = "1.1.1.213"
    member_port   = "65534"
    member_state  = "ENABLE"
    member_weight = "1"
    name          = "skubu11"
    object_id     = "ENTER YOUR RESOURCE'S OBJECT_ID"
    object_type   = "VM"
  }
}

variable "lb_member_modify" {
  type = object({
    member_port  = string
    member_state = string
  })
  default = {
    member_port  = "65530"
    member_state = "DISABLE"
  }
}


