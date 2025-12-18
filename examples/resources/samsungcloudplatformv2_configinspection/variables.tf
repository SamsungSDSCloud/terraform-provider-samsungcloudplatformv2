variable "account_id" {
  description = "account Id"
  type        = string
  default     = ""
}

variable "csp_type" {
  description = "Type of cloud service provider"
  type        = string
  default     = ""
}

variable "diagnosis_account_id" {
  description = "Id of diagnosis"
  type        = string
  default     = ""
}

variable "diagnosis_check_type" {
  description = "Check type of diagnosis"
  type        = string
  default     = ""
}

variable "diagnosis_id" {
  description = "Id of diagnosis"
  type        = string
  default     = ""
}

variable "diagnosis_name" {
  description = "Name of diagnosis"
  type        = string
  default     = ""
}

variable "diagnosis_type" {
  description = "diagnosis Type"
  type        = string
  default     = ""
}

variable "plan_type" {
  description = "plan Type"
  type        = string
  default     = ""
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
    auth_key_created_at = null
    auth_key_expired_at = null
    auth_key_id         = ""
    diagnosis_id        = ""
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
    diagnosis_id                 = ""
    diagnosis_start_time_pattern = ""
    frequency_type               = ""
    frequency_value              = ""
    use_diagnosis_check_type_bp  = ""
    use_diagnosis_check_type_ssi = ""
  }
}

variable "tags" {
  description = "Tags for the resource"
  type        = map(string)
  default     = null
}


