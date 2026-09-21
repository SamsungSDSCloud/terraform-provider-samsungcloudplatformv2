package iamidentitycenter_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIamIdentityCenterAccountAssignmentResourceTest(t *testing.T) {
	testInstanceName := fmt.Sprintf("test-acc-aa-i-%d", time.Now().UnixNano())
	testPermissionSetName := fmt.Sprintf("test-acc-aa-ps-%d", time.Now().UnixNano())
	testUserId := fmt.Sprintf("test-acc-aa-u-%d", time.Now().UnixNano())

	targetAccountId := os.Getenv("SCP_TEST_TARGET_ACCOUNT_ID")
	if targetAccountId == "" {
		t.Skip("Set SCP_TEST_TARGET_ACCOUNT_ID to run this test")
	}

	existingInstanceId := getExistingInstanceId()

	resourceName := "samsungcloudplatformv2_iam_identity_center_account_assignment.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAccountAssignmentCreate(
					testInstanceName,
					testPermissionSetName,
					testUserId,
					targetAccountId,
					existingInstanceId,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "account_assignment_id"),
					resource.TestCheckResourceAttr(resourceName, "target_account_id", targetAccountId),
					resource.TestCheckResourceAttr(resourceName, "principal_type", "USER"),
				),
			},
			{
				Config: testAccAccountAssignmentDelete(
					testInstanceName,
					testPermissionSetName,
					testUserId,
					targetAccountId,
					existingInstanceId,
				),
			},
		},
	})
}

func testAccAccountAssignmentCreate(instanceName string, permissionSetName string, userId string, targetAccountId string, existingInstanceId string) string {
	if existingInstanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = "%s"
			user_id     = "%s"
			name        = "Test User"
			email       = "%s@example.com"
			password    = "P@ssw0rd!2024Secure"
		}

		resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
			instance_id      = "%s"
			name             = "%s"
			description      = "test permission set for account assignment"
			session_duration = 3600
		}

		resource "samsungcloudplatformv2_iam_identity_center_account_assignment" "example" {
			instance_id       = "%s"
			target_account_id = "%s"
			principal_id      = samsungcloudplatformv2_iam_identity_center_user.user.user_uuid
			principal_type    = "USER"
			permission_set_id = samsungcloudplatformv2_iam_identity_center_permission_set.permission_set.id
		}`, existingInstanceId, userId, userId, existingInstanceId, permissionSetName, existingInstanceId, targetAccountId)
	}

	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name                = "%s"
			description         = "test instance for account assignment"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			user_id     = "%s"
			name        = "Test User"
			email       = "%s@example.com"
			password    = "P@ssw0rd!2024Secure"
		}

		resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
			instance_id      = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			name             = "%s"
			description      = "test permission set for account assignment"
			session_duration = 3600
		}

		resource "samsungcloudplatformv2_iam_identity_center_account_assignment" "example" {
			instance_id       = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			target_account_id = "%s"
			principal_id      = samsungcloudplatformv2_iam_identity_center_user.user.user_uuid
			principal_type    = "USER"
			permission_set_id = samsungcloudplatformv2_iam_identity_center_permission_set.permission_set.id
		}`, instanceName, userId, userId, permissionSetName, targetAccountId)
}

func testAccAccountAssignmentDelete(instanceName string, permissionSetName string, userId string, targetAccountId string, existingInstanceId string) string {
	return testAccAccountAssignmentCreate(instanceName, permissionSetName, userId, targetAccountId, existingInstanceId)
}
