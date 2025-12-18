variable "budget_name" {
  type    = string
  default = ""
}

variable "budget_amount" {
  type    = number
  default = 0
}

variable "budget_unit" {
  type    = string
  default = ""
}

variable "budget_start_month" {
  type    = string
  default = ""
}

variable "budget_notifications" {
  type    = map
  default = null
}

variable "budget_prevention" {
  type    = map
  default = null
}



