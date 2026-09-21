output "lbMembers" {
	value = {
		count: data.samsungcloudplatformv2_loadbalancer_lb_members.lbmembers.total_count
		lbmembers: data.samsungcloudplatformv2_loadbalancer_lb_members.lbmembers
	}
}