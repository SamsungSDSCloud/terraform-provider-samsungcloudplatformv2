variable "account_id" {
  description = "account Id"
  type        = string
  default     = "ENTER YOUR RESOURCE'S ACCOUNT_ID"
}

variable "csp_type" {
  description = "Type of cloud service provider"
  type        = string
  default     = "SCP"
}

variable "diagnosis_account_id" {
  description = "Id of diagnosis"
  type        = string
  default     = "ENTER YOUR RESOURCE'S DIAGNOSIS_ACCOUNT_ID"
}

variable "diagnosis_check_type" {
  description = "Check type of diagnosis"
  type        = string
  default     = "BP"
}

variable "diagnosis_id" {
  description = "Id of diagnosis"
  type        = string
  default     = "ENTER YOUR RESOURCE'S DIAGNOSIS_ID"
}

variable "diagnosis_name" {
  description = "Name of diagnosis"
  type        = string
  default     = "gdcv-test-123"
}

variable "diagnosis_type" {
  description = "diagnosis Type"
  type        = string
  default     = "Console"
}

variable "plan_type" {
  description = "plan Type"
  type        = string
  default     = "STANDARD"
}

variable "auth_key_request" {
  description = "Auth key request"
  type = object({
    diagnosis_id        = string
    auth_key_created_at = optional(string)
    auth_key_expired_at = optional(string)
    auth_key_id         = string
  })
  default = {
    auth_key_id  = "ENTER YOUR RESOURCE'S AUTH_KEY_ID"
    diagnosis_id = "ENTER YOUR RESOURCE'S DIAGNOSIS_ID"
  }
}

variable "schedule_request" {
  description = "Schedule request"
  type = object({
    diagnosis_id                 = string
    diagnosis_start_time_pattern = string
    frequency_type               = string
    frequency_value              = string
    use_diagnosis_check_type_bp  = string
    use_diagnosis_check_type_ssi = string
  })
  default = {
    diagnosis_id                 = "ENTER YOUR RESOURCE'S DIAGNOSIS_ID"
    diagnosis_start_time_pattern = "00:00"
    frequency_type               = "week"
    frequency_value              = "monday"
    use_diagnosis_check_type_bp  = "y"
    use_diagnosis_check_type_ssi = "n"
  }
}

variable "tags" {
  description = "Tags for the resource"
  type        = map(string)
  default = {
    tag1 = "tag1"
  }
}



