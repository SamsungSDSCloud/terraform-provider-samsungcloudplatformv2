variable "name" {
  type    = string
  default = ""
}

variable "capacity_tb" {
  type    = number
  default = 100
}

variable "zone" {
  type    = string
  default = "kr-west1-a"
}

variable "access_rules" {
  type = list(object({
    object_type = string,
    object_id   = string
  }))
  default = []
}

variable "tags" {
  type = map(string)
  default = {
    test_terraform = "test_terraform_value"
  }
}



