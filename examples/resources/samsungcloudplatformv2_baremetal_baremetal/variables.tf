variable "image_id" {
  type    = string
  default = ""
}

variable "init_script" {
  type    = string
  default = ""
}

variable "lock_enabled" {
  type    = bool
  default = false
}

variable "os_user_id" {
  type    = string
  default = ""
}

variable "os_user_password" {
  type    = string
  default = ""
}

variable "placement_group_name" {
  type    = string
  default = ""
}

variable "region_id" {
  type    = string
  default = ""
}

variable "subnet_id" {
  type    = string
  default = ""
}

variable "use_placement_group" {
  type    = bool
  default = false
}

variable "vpc_id" {
  type    = string
  default = ""
}

variable "server_details" {
  type = list(object({
    bare_metal_local_subnet_id         = optional(string)
    bare_metal_local_subnet_ip_address = optional(string)
    bare_metal_server_name             = string
    ip_address                         = optional(string)
    nat_enabled                        = bool
    public_ip_address_id               = string
    server_type_id                     = string
    use_hyper_threading                = optional(bool, false)
    state                              = optional(string)
  }))
  default = [{
    bare_metal_local_subnet_id         = null
    bare_metal_local_subnet_ip_address = null
    bare_metal_server_name             = ""
    ip_address                         = null
    nat_enabled                        = false
    public_ip_address_id               = ""
    server_type_id                     = ""
    state                              = null
    use_hyper_threading                = null
  }]
}

variable "tags" {
  type    = map(string)
  default = null
}

variable "create_timeouts" {
  type    = string
  default = ""
}

variable "delete_timeouts" {
  type    = string
  default = ""
}

