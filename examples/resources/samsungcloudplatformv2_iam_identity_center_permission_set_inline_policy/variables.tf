variable "instance_id" {
  type        = string
  description = "The ID of the Identity Center Instance"
  default     = "ENTER YOUR RESOURCE'S INSTANCE_ID"
}

variable "permission_set_id" {
  type        = string
  description = "The ID of the IAM Identity Center Permission Set"
  default     = "ENTER YOUR RESOURCE'S PERMISSION_SET_ID"
}

variable "policy_document" {
  type        = any
  description = "Inline Policy JSON Document"
  default = {
    Statement = [{
      Action   = ["s3:*"]
      Effect   = "Allow"
      Resource = ["*"]
    }]
    Version = "2012-10-17"
  }
}



