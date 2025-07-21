variable "account_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "contract_type" {
  type = string
  default = "03"
}

variable "os_type" {
  type = string
  default = "WINDOWS"
}

variable "server_type" {
  type = string
  default = "s1v16m128" 
}

variable "service_id" {
  type = string
  default = "VIRTUAL_SERVER"
}

variable "service_name" {
  type = string
  default = "Virtual Server"
}

variable "action" {
  type = string
  default = "EXTEND_APPLY"
  
}

variable "region" {
    type = string
    default = "kr-west1"
}
