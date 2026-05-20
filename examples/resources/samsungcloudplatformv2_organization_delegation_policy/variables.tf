variable "organization_id" {
  type    = string
  default = ""
}

variable "document" {
  type = object({
    version = string
    statement = list(object({
      sid    = optional(string)
      effect = string
      action = list(string)
      principal = optional(object({
        scp = optional(list(string))
      }))
      resource = list(string)
    }))
  })
  default = {
    statement = [{
      action    = [""]
      effect    = ""
      principal = null
      resource  = [""]
      sid       = null
    }]
    version = ""
  }
}

