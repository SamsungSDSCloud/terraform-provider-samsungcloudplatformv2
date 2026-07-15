output "private_nat_output" {
  value = {
    private_nat: samsungcloudplatformv2_vpc_private_nat.privatenat.private_nat,
  }
}
