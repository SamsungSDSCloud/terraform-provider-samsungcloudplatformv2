variable "approver_vpc_account_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S APPROVER_VPC_ACCOUNT_ID"
}

variable "approver_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S APPROVER_VPC_ID"
}

variable "name" {
  type    = string
  default = "nam-test-vpc"
}

variable "requester_vpc_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S REQUESTER_VPC_ID"
}

variable "tags" {
  type = map(string)
  default = {
    test_tag_key  = "test_tag_value"
    test_tag_key2 = "test_tag_value2"
  }
}

variable "description" {
  default = "description aaa 1111"
}


