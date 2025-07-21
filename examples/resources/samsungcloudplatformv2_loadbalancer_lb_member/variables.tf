variable "lb_server_group_id"{
  type = string
  default = "22"
}

variable "lb_member" {
  type = object({
    name = string
    member_ip = string
    member_port= string
    member_weight= string
    object_type= string
    object_id= string
  })
  default = {
      name: "terraform",
      member_ip: "1.1.1.213",
      member_port: "65534",
      member_weight: "10",
      object_type: "VM",
      object_id: "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
    }
}

variable "lb_member_modify" {
  type = object({
    member_port= string
    member_weight= string
    member_state= string
  })
  default = {
      member_port: "65533",
      member_weight: "100",
      member_state: "DISABLE"
  }
}