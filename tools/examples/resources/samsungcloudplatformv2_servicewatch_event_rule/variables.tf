variable "event_ids" {
  type    = list(string)
  default = ["ENTER YOUR RESOURCE'S EVENT_IDS"]
}

variable "event_rule_name" {
  type    = string
  default = "tf_event_rule_test"
}

variable "recipient_ids" {
  type    = list(string)
  default = []
}

variable "service_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SERVICE_ID"
}

variable "resource_type_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S RESOURCE_TYPE_ID"
}

variable "srn_list" {
  type    = list(string)
  default = []
}

variable "description" {
  type    = string
  default = "terraform test"
}

variable "tags" {
  type = map(string)
  default = {
    test_terraform_tag_key = "test_terraform_tag_value"
  }
}


