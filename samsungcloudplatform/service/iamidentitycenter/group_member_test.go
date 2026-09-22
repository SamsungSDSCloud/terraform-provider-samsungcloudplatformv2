package iamidentitycenter_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	iamidentitycenterclient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/iamidentitycenter"
	iamidentitycentercommon "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/iamidentitycenter"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIamIdentityCenterGroupMemberResourceTest(t *testing.T) {
	testInstanceName := fmt.Sprintf("test-acc-gm-i-%d", time.Now().UnixNano())
	testGroupName := fmt.Sprintf("test-acc-gm-g-%d", time.Now().UnixNano())
	testUserId := fmt.Sprintf("test-acc-gm-u-%d", time.Now().UnixNano())

	existingInstanceId := getExistingInstanceId()
	existingGroupId := os.Getenv("SCP_TEST_GROUP_ID")
	testUserUuid := os.Getenv("SCP_TEST_USER_UUID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupMemberCreate(
					testInstanceName,
					testGroupName,
					testUserId,
					existingInstanceId,
					existingGroupId,
					testUserUuid,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_iam_identity_center_group_member.member", "group_id"),
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_iam_identity_center_group_member.member", "member_id"),
				),
			},
			{
				Config: testAccGroupMemberUpdate(
					testInstanceName,
					testGroupName,
					testUserId,
					existingInstanceId,
					existingGroupId,
					testUserUuid,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_iam_identity_center_group_member.member", "group_id"),
				),
			},
		},
	})
}

func testAccGroupMemberCreate(instanceName string, groupName string, userId string, existingInstanceId string, existingGroupId string, existingUserUuid string) string {
	// Case 1: All three existing resources provided — use them directly
	if existingInstanceId != "" && existingGroupId != "" && existingUserUuid != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_group_member" "member" {
			instance_id = "%s"
			group_id    = "%s"
			member_id   = "%s"
		}`, existingInstanceId, existingGroupId, existingUserUuid)
	}

	// Case 2: Existing instance provided but group/user need to be created dynamically
	if existingInstanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_group" "group" {
			instance_id = "%s"
			name        = "%s"
			description = "test group for group member"
		}

		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = "%s"
			user_id     = "%s"
			name        = "Test User"
			email       = "%s@example.com"
			password    = "P@ssw0rd!2024Secure"
		}

		resource "samsungcloudplatformv2_iam_identity_center_group_member" "member" {
			instance_id = "%s"
			group_id    = samsungcloudplatformv2_iam_identity_center_group.group.id
			member_id   = samsungcloudplatformv2_iam_identity_center_user.user.user_uuid
		}`, existingInstanceId, groupName, existingInstanceId, userId, userId, existingInstanceId)
	}

	// Case 3: No existing resources — create everything (may fail with 409 if instance already exists)
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name               = "%s"
			description        = "test instance for group member"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_group" "group" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			name        = "%s"
			description = "test group for group member"
		}

		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			user_id     = "%s"
			name        = "Test User"
			email       = "%s@example.com"
			password    = "P@ssw0rd!2024Secure"
		}

		resource "samsungcloudplatformv2_iam_identity_center_group_member" "member" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			group_id    = samsungcloudplatformv2_iam_identity_center_group.group.id
			member_id   = samsungcloudplatformv2_iam_identity_center_user.user.user_uuid
		}`, instanceName, groupName, userId, userId)
}

func testAccGroupMemberUpdate(instanceName string, groupName string, userId string, existingInstanceId string, existingGroupId string, existingUserUuid string) string {
	// Case 1: All three existing resources provided — use them directly
	if existingInstanceId != "" && existingGroupId != "" && existingUserUuid != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_group_member" "member" {
			instance_id = "%s"
			group_id    = "%s"
			member_id   = "%s"
		}`, existingInstanceId, existingGroupId, existingUserUuid)
	}

	// Case 2: Existing instance provided but group/user need to be created dynamically
	if existingInstanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_group" "group" {
			instance_id = "%s"
			name        = "%s"
			description = "test group for group member"
		}

		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = "%s"
			user_id     = "%s"
			name        = "Test User"
			email       = "%s@example.com"
			password    = "P@ssw0rd!2024Secure"
		}

		resource "samsungcloudplatformv2_iam_identity_center_group_member" "member" {
			instance_id = "%s"
			group_id    = samsungcloudplatformv2_iam_identity_center_group.group.id
			member_id   = samsungcloudplatformv2_iam_identity_center_user.user.user_uuid
		}`, existingInstanceId, groupName, existingInstanceId, userId, userId, existingInstanceId)
	}

	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name               = "%s"
			description        = "test instance for group member"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_group" "group" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			name        = "%s"
			description = "test group for group member"
		}

		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			user_id     = "%s"
			name        = "Test User"
			email       = "%s@example.com"
			password    = "P@ssw0rd!2024Secure"
		}

		resource "samsungcloudplatformv2_iam_identity_center_group_member" "member" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			group_id    = samsungcloudplatformv2_iam_identity_center_group.group.id
			member_id   = samsungcloudplatformv2_iam_identity_center_user.user.user_uuid
		}`, instanceName, groupName, userId, userId)
}

func TestAccIamIdentityCenterGroupMemberDataSourceTest(t *testing.T) {
	testInstanceName := fmt.Sprintf("test-acc-gm-ds-i-%d", time.Now().UnixNano())
	testGroupName := fmt.Sprintf("test-acc-gm-ds-g-%d", time.Now().UnixNano())

	existingInstanceId := getExistingInstanceId()
	existingGroupId := os.Getenv("SCP_TEST_GROUP_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupMemberDataSource(testInstanceName, testGroupName, existingInstanceId, existingGroupId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.samsungcloudplatformv2_iam_identity_center_group_member.member", "group_id"),
				),
			},
		},
	})
}

func testAccGroupMemberDataSource(instanceName string, groupName string, existingInstanceId string, existingGroupId string) string {
	// Case 1: Both instance and group provided — use directly
	if existingInstanceId != "" && existingGroupId != "" {
		return fmt.Sprintf(`
		data "samsungcloudplatformv2_iam_identity_center_group_member" "member" {
			instance_id = "%s"
			group_id    = "%s"
		}`, existingInstanceId, existingGroupId)
	}

	// Case 2: Only instance provided — create group dynamically
	if existingInstanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_group" "group" {
			instance_id = "%s"
			name        = "%s"
			description = "test group for group member data source"
		}

		data "samsungcloudplatformv2_iam_identity_center_group_member" "member" {
			instance_id = "%s"
			group_id    = samsungcloudplatformv2_iam_identity_center_group.group.id
		}`, existingInstanceId, groupName, existingInstanceId)
	}

	// Case 3: Nothing provided — create everything
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name               = "%s"
			description        = "test instance for group member data source"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_group" "group" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			name        = "%s"
			description = "test group for group member data source"
		}

		data "samsungcloudplatformv2_iam_identity_center_group_member" "member" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			group_id    = samsungcloudplatformv2_iam_identity_center_group.group.id
		}`, instanceName, groupName)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_iam_identity_center_group_member", &resource.Sweeper{
		Name: "samsungcloudplatformv2_iam_identity_center_group_member",
		F:    sweepGroupMember,
	})
}

func sweepGroupMember(region string) error {
	ctx := context.Background()
	scpClient, err := iamidentitycentercommon.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	instances, err := scpClient.Client.IamIdentityCenter.ListInstances(ctx, 0, 0, "")
	if err != nil {
		return fmt.Errorf("error listing instances: %v", err)
	}

	for _, inst := range instances.GetInstances() {
		instanceId := inst.Id

		sweepGroupMemberUsers(ctx, scpClient.Client.IamIdentityCenter, instanceId)
		if !sweepGroupMemberGroups(ctx, scpClient.Client.IamIdentityCenter, instanceId) {
			continue
		}
		sweepGroupMemberInstance(ctx, scpClient.Client.IamIdentityCenter, instanceId)
	}

	return nil
}

func sweepGroupMemberUsers(ctx context.Context, iamClient *iamidentitycenterclient.Client, instanceId string) {
	users, err := iamClient.ListUsers(ctx, instanceId, "", 100, 0, "", "", "")
	if err != nil {
		fmt.Printf("warning: error listing users for instance %s: %v\n", instanceId, err)
		return
	}

	for _, u := range users.GetUsers() {
		if strings.HasPrefix(u.UserId, "test-acc-gm-") {
			_, _, err = iamClient.DeleteUser(ctx, u.Id, instanceId)
			if err != nil {
				fmt.Printf("warning: failed to delete user %s: %v\n", u.Id, err)
			}
		}
	}
}

func sweepGroupMemberGroups(ctx context.Context, iamClient *iamidentitycenterclient.Client, instanceId string) bool {
	groups, err := iamClient.ListGroups(ctx, instanceId, "", 100, 0, "", "", "")
	if err != nil {
		fmt.Printf("warning: error listing groups for instance %s: %v\n", instanceId, err)
		return false
	}

	for _, g := range groups.GetGroups() {
		if strings.HasPrefix(g.Name, "test-acc-gm-") {
			_, _, err = iamClient.DeleteGroup(ctx, g.Id, instanceId)
			if err != nil {
				fmt.Printf("warning: failed to delete group %s: %v\n", g.Id, err)
			}
		}
	}

	return true
}

func sweepGroupMemberInstance(ctx context.Context, iamClient *iamidentitycenterclient.Client, instanceId string) {
	instanceDetail, _, err := iamClient.GetInstance(ctx, instanceId)
	if err != nil {
		return
	}

	if strings.HasPrefix(instanceDetail.Instance.Name, "test-acc-gm-") {
		_, err = iamClient.DeleteInstance(ctx, instanceId)
		if err != nil {
			fmt.Printf("warning: failed to delete instance %s: %v\n", instanceId, err)
		}
	}
}
