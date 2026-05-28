variable "id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ID"
}

variable "name" {
  type    = string
  default = ""
}

variable "requester_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S REQUESTER_VPC_ID"
}

variable "requester_vpc_name" {
  type    = string
  default = ""
}

variable "approver_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S APPROVER_VPC_ID"
}

variable "approver_vpc_name" {
  type    = string
  default = ""
}

variable "account_type" {
  type    = string
  default = "SAME"
}

variable "size" {
  type    = number
  default = 10
}

variable "page" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = "created_at:desc"
}

variable "state" {
  type    = string
  default = "ACTIVE"
}


