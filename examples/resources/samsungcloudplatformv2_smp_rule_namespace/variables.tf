variable "rule_namespace_name" {
  description = "Rule namespace name."
  type        = string
  default     = "tf-test-rule-namespace"
}

variable "config_data" {
  description = "Base64 encoded YAML rules file."
  type        = string
  default     = "Z3JvdXBzOgogIC0gbmFtZTogZXhhbXBsZQogICAgcnVsZXM6CiAgICAgIC0gYWxlcnQ6IEhpZ2hDUFUKICAgICAgICBleHByOiBjcHVfdXNhZ2UgPiA4MAo="
}

variable "workspace_id" {
  description = "Workspace ID."
  type        = string
  default     = "ENTER YOUR RESOURCE'S WORKSPACE_ID"
}



