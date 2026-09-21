variable "organization_id" {
  type        = string
  description = "Organization ID"
  default     = "ENTER YOUR RESOURCE'S ORGANIZATION_ID"
}

variable "entity" {
  type        = string
  description = "Entity type for the binding"
  default     = "POLICY"
}

variable "target_ids" {
  type        = set(string)
  description = "Set of target IDs to bind policies to"
  default     = ["ou-2a4c34d21ff14161adec23c83974d267"]
}

variable "policy_ids" {
  type        = set(string)
  description = "Set of control policy IDs to attach to the target"
  default     = ["9bfdb307d88a452e8151aa6091170b88", "e0102d253a8c4c76a5c7d4aa8e0ea128"]
}



