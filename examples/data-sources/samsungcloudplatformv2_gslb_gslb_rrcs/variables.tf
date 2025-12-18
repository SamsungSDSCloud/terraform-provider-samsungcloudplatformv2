variable "size" {
  description = "size"
  type        = number
  default     = 0
}

variable "page" {
  description = "page"
  type        = number
  default     = 0
}

variable "sort" {
  description = "sort"
  type        = string
  default     = ""
}

variable "region" {
  description = "The GSLB Resource Region."
  type        = string
  default     = ""
}

variable "status" {
  description = "The GSLB Resource Status."
  type        = string
  default     = ""
}

variable "name" {
  description = "The Name of the gslb."
  type        = string
  default     = ""
}

