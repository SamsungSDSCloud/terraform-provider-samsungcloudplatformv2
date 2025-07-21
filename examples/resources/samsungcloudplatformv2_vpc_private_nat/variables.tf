variable "private_nat_name" {
  type = string
  default = "sample"
}

variable "private_nat_direct_connect_id" {
  type = string
  default = "sample-direct-connect-id"
}

variable "private_nat_cidr" {
  type = string
  default = "sample-cidr"
}

variable "private_nat_description" {
  type = string
  default = "sample-description"
}

variable "private_nat_tags" {
  type    = map(string)
  default = {
    "tf_key1" = "tf_val1"
  }
}
