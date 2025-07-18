---
page_title: "samsungcloudplatformv2_vpn_vpn_tunnel Data Source - samsungcloudplatformv2"
subcategory: VPN Tunnel
description: |-
  VPN Tunnel
---

# samsungcloudplatformv2_vpn_vpn_tunnel (Data Source)

VPN Tunnel

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpn_vpn_tunnel" "vpntunnel" {
  id = var.id
}


output "vpn_tunnel" {
  value = data.samsungcloudplatformv2_vpn_vpn_tunnel.vpntunnel.vpn_tunnel
}


variable "id" {
  type    = string
  default = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Id

### Read-Only

- `vpn_tunnel` (Attributes) VPN Tunnel (see [below for nested schema](#nestedatt--vpn_tunnel))

<a id="nestedatt--vpn_tunnel"></a>
### Nested Schema for `vpn_tunnel`

Read-Only:

- `account_id` (String) AccountId
- `created_at` (String) CreatedAt
- `created_by` (String) CreatedBy
- `description` (String) Description
- `id` (String) Id
- `modified_at` (String) ModifiedAt
- `modified_by` (String) ModifiedBy
- `name` (String) Name
- `phase1` (Attributes) Phase1 (see [below for nested schema](#nestedatt--vpn_tunnel--phase1))
- `phase2` (Attributes) Phase2 (see [below for nested schema](#nestedatt--vpn_tunnel--phase2))
- `state` (String) State
- `vpc_id` (String) VpcId
- `vpc_name` (String) VpcName
- `vpn_gateway_id` (String) VpnGatewayId
- `vpn_gateway_ip_address` (String) VpnGatewayIpAddress
- `vpn_gateway_name` (String) VpnGatewayName

<a id="nestedatt--vpn_tunnel--phase1"></a>
### Nested Schema for `vpn_tunnel.phase1`

Read-Only:

- `diffie_hellman_groups` (List of Number) DiffieHellmanGroups
- `dpd_retry_interval` (Number) DpdRetryInterval
- `encryptions` (List of String) Encryptions
- `ike_version` (Number) IkeVersion
- `life_time` (Number) LifeTime
- `peer_gateway_ip` (String) PeerGatewayIp
- `pre_shared_key` (String) PreSharedKey


<a id="nestedatt--vpn_tunnel--phase2"></a>
### Nested Schema for `vpn_tunnel.phase2`

Read-Only:

- `diffie_hellman_groups` (List of Number) DiffieHellmanGroups
- `encryptions` (List of String) Encryptions
- `life_time` (Number) LifeTime
- `perfect_forward_secrecy` (String) PerfectForwardSecrecy
- `remote_subnet` (String) RemoteSubnet