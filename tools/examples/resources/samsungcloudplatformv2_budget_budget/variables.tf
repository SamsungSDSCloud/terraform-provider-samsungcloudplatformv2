variable "budget_name" {
  type    = string
  default = "test_budget"
}

variable "budget_amount" {
  type    = number
  default = 5000000
}

variable "budget_unit" {
  type    = string
  default = "MONTHLY"
}

variable "budget_start_month" {
  type    = string
  default = "2025-12"
}

variable "budget_notifications" {
  type = map
  default = {
    is_use_notification = false
  }
}

variable "budget_prevention" {
  type = map
  default = {
    is_use_prevention = false
  }
}




