variable "lb_server_group" {
  type = object({
    name = string
    vpc_id = string
    subnet_id= string
    protocol= string
    lb_method= string
    description= string
    lb_health_check_id= string
    tags=map(string)
  })
  default = {
    name = "terraform"
    vpc_id = "8a463aa4b1dc4f279c3f53b94dc45e74"
    subnet_id= "8a463aa4b1dc4f279c3f53b94dc45e74"
    protocol= "TCP"
    lb_method= "ROUND_ROBIN"
    description= "description info"
    lb_health_check_id="8a463aa4b1dc4f279c3f53b94dc45e74"
    tags={
      "key" : "value"
      "key1" : "value1"
    }
  }
}

variable "lb_server_group_modify" {
  type = object({
    lb_method= string
    description= string
    lb_health_check_id= string
  })
  default = {
    lb_method= "LEAST_CONNECTION"
    description= "description info"
    lb_health_check_id="8a463aa4b1dc4f279c3f53b94dc45e74"
  }
}