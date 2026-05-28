variable "image_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S IMAGE_ID"
}

variable "init_script" {
  type    = string
  default = "init script"
}

variable "lock_enabled" {
  type    = bool
  default = false
}

variable "os_user_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S OS_USER_ID"
}

variable "os_user_password" {
  type    = string
  default = "ENTER YOUR RESOURCE'S OS_USER_PASSWORD"
}

variable "placement_group_name" {
  type    = string
  default = "pg-group"
}

variable "region_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S REGION_ID"
}

variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
}

variable "use_placement_group" {
  type    = bool
  default = true
}

variable "vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VPC_ID"
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
    bare_metal_local_subnet_id         = "ENTER YOUR RESOURCE'S BARE_METAL_LOCAL_SUBNET_ID"
    bare_metal_local_subnet_ip_address = ""
    bare_metal_server_name             = "terraform-test01"
    ip_address                         = "192.168.1.3"
    nat_enabled                        = false
    public_ip_address_id               = "ENTER YOUR RESOURCE'S PUBLIC_IP_ADDRESS_ID"
    server_type_id                     = "ENTER YOUR RESOURCE'S SERVER_TYPE_ID"
    state                              = "RUNNING"
    use_hyper_threading                = true
  }]
}

variable "tags" {
  type = map(string)
  default = {
    test_tag_key  = "test_tag_value"
    test_tag_key2 = "test_tag_value2"
  }
}

variable "create_timeouts" {
  type    = string
  default = "60m"
}

variable "delete_timeouts" {
  type    = string
  default = "60m"
}


