package iam_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupMemberResourceTest(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupMemberWithUserAndGroup("samsungcloudplatformv2_iam_user.user_a.user_id"),
			},
			{
				Config: testAccGroupMemberWithUserAndGroup("samsungcloudplatformv2_iam_user.user_b.user_id"),
			},
			{
				Config: testAccGroupMemberRemove(),
			},
			{
				Config: testAccGroupRemoveForGroupMember(),
			},
		},
	})
}

func getIamUserA() string {
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_user" "user_a" {
				account_id = data.samsungcloudplatformv2_iam_access_keys.access_keys.access_keys[0].account_id
  				description = "test-acc iam user a desc"
				password = "U2NvcmVTY3AhMjM="
				temporary_password = true
				user_name = "test-acc-user-a"
			}`)
}

func getIamUserB() string {
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_user" "user_b" {
				account_id = data.samsungcloudplatformv2_iam_access_keys.access_keys.access_keys[0].account_id
  				description = "test-acc iam user b desc"
				password = "U2NvcmVTY3AhMjM="
				temporary_password = true
				user_name = "test-acc-user-b"
			}`)
}

func getGroupForGroupMember() string {
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_group" "group" {
				name = "test-acc-group"
  				description = "test-acc group desc"
				tags = {
					"test-acc-key" = "test-acc-value"
				}
				policy_ids = []
			}
	`)
}

func testAccGroupMemberWithUserAndGroup(userId string) string {
	return fmt.Sprintf(`
			%s
			%s
			%s
			%s

			resource "samsungcloudplatformv2_iam_group_member" "group_member" {
				group_id = samsungcloudplatformv2_iam_group.group.id
  				user_id = %s
			}`, getAccountIdForGroupMember(), getIamUserA(), getIamUserB(), getGroupForGroupMember(), userId)
}

func testAccGroupMemberRemove() string {
	return fmt.Sprintf(`
			%s
			%s
			%s
			%s
	`, getAccountIdForGroupMember(), getIamUserA(), getIamUserB(), getGroupForGroupMember())
}

func testAccGroupRemoveForGroupMember() string {
	return fmt.Sprintf(`
			%s
			%s
			%s
	`, getAccountIdForGroupMember(), getIamUserA(), getIamUserB())
}

func getAccountIdForGroupMember() string {
	return fmt.Sprintf(`
			data "samsungcloudplatformv2_iam_access_keys" "access_keys" {
				limit = 1
			}
	`)
}
