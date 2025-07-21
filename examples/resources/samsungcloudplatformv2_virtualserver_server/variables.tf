variable "name" {
  type = string
  default = "terraform_server"
}

variable "state" {
  type = string
  default = "ACTIVE"
}

variable "image_id" {
  type = string
  default = "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
}

variable "server_type_id" {
  type = string
  default = "s1v1m2"
}

variable "keypair_name" {
  type = string
  default = "terraform_keypair"
}

variable "lock" {
  type = bool
  default = false
}

variable "user_data" {
  type = string
  default = ""
}

variable "boot_volume" {
  type = object({
    size = number
    type = string
    delete_on_termination = bool
  })
  default = {
      size = 16
      type = "SSD"
      delete_on_termination = true
    }
}

variable "server_group_id" {
  type = string
  default = "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
}

variable "security_groups" {
  type = list(string)
  default = []
}

variable "networks_interface_1_subnet_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "extra_volumes_volume_1_size" {
  type = number
  default = 16
}

variable "extra_volumes_volume_1_type" {
  type = string
  default = "SSD"
}

variable "extra_volumes_volume_1_delete_on_termination" {
  type = bool
  default = true
}
