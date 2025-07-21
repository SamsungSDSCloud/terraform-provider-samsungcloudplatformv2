variable "lb_health_check_http" {
  type = object({
    name = string
    vpc_id = string
    subnet_id= string
    protocol= string
    health_check_port= string
    health_check_interval= string
    health_check_timeout= string
    health_check_count= string
    http_method= string
    health_check_url= string
    response_code= string
    request_data= string
    description= string
    tags=map(string)
  })
  default = {
    name = "terraform-http"
    vpc_id = "8a463aa4b1dc4f279c3f53b94dc45e74"
    subnet_id= "8a463aa4b1dc4f279c3f53b94dc45e74"
    protocol= "HTTP"
    health_check_port= "100"
    health_check_interval= "30"
    health_check_timeout= "20"
    health_check_count= "5"
    http_method= "POST"
    health_check_url= "/terraform"
    response_code= "200"
    request_data= "username=name&password=123456789"
    description= "description info"
    tags={
      "key" : "value"
      "key1" : "value1"
    }
  }
}

variable "lb_health_check_tcp" {
  type = object({
    name = string
    vpc_id = string
    subnet_id= string
    protocol= string
    health_check_port= string
    health_check_interval= string
    health_check_timeout= string
    health_check_count= string
    description= string
    tags=map(string)
  })
  default = {
    name = "terraform-tcp"
    vpc_id = "8a463aa4b1dc4f279c3f53b94dc45e74"
    subnet_id= "8a463aa4b1dc4f279c3f53b94dc45e74"
    protocol= "TCP"
    health_check_port= "100"
    health_check_interval= "30"
    health_check_timeout= "20"
    health_check_count= "5"
    description= "description info"
    tags={
      "key" : "value"
      "key1" : "value1"
    }
  }
}

variable "lb_health_check_modify" {
  type = object({
    protocol= string
    health_check_port= string
    health_check_interval= string
    health_check_timeout= string
    health_check_count= string
    http_method= string
    health_check_url= string
    response_code= string
    request_data= string
    description= string
  })
  default = {
    protocol= "TCP"
    health_check_port= "200"
    health_check_interval= "20"
    health_check_timeout= "10"
    health_check_count= "7"
    http_method= null
    health_check_url= null
    response_code= null
    request_data= null
    description= "description info"
  }
}