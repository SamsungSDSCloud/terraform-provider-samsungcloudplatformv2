output "privatenat" {
  value = {
    private_nat: data.samsungcloudplatformv2_vpc_private_nat.my_private_nat.private_nat
  }
}
