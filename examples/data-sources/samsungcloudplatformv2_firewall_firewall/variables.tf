variable "id" {
  type = string
  default =  "68db67f78abd405da98a6056a8ee42af"
}

variable "product_type" {
  type = list(string)
  default =  ["LB", "IGW"]
}
