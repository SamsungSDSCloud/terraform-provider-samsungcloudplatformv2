variable "approver_vpc_account_id" {
  type    = string
  default = ""
}

variable "approver_vpc_id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "requester_vpc_id" {
  type    = string
  default = ""
}

variable "tags" {
  type    = map(string)
  default = null
}

variable "description" {
  default = null
}

