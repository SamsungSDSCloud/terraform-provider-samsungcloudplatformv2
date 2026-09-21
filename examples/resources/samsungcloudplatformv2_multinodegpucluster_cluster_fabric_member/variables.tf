variable "before_cluster_fabric_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S BEFORE_CLUSTER_FABRIC_ID"
}

variable "after_cluster_fabric_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S AFTER_CLUSTER_FABRIC_ID"
}

variable "gpu_node_id_list" {
  type    = list(string)
  default = ["87e5175a74e1456dbaba949d88e14968"]
}



