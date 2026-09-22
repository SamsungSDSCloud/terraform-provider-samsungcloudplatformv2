package iamidentitycenter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIamIdentityCenterPermissionSetManagedPolicyResourceTest(t *testing.T) {
	testPermissionSetName := fmt.Sprintf("test-acc-psmp-%d", time.Now().UnixNano())
	testManagedPolicyId := "37f2e31ff86b415698d7e8eeafab445d"
	testManagedPolicyName := "AmazonS3ReadOnlyAccess"
	existingInstanceId := getExistingInstanceId()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionSetManagedPolicyCreate(existingInstanceId, testPermissionSetName, testManagedPolicyId, testManagedPolicyName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_iam_identity_center_permission_set_managed_policy.policy", "permission_set_id"),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_permission_set_managed_policy.policy", "managed_policy_id", testManagedPolicyId),
				),
			},
		},
	})
}

func testAccPermissionSetManagedPolicyCreate(instanceId string, permissionSetName string, managedPolicyId string, managedPolicyName string) string {
	return fmt.Sprintf(`
	resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
		instance_id      = "%s"
		name             = "%s"
		description      = "test permission set for managed policy"
		session_duration = 3600
	}

	resource "samsungcloudplatformv2_iam_identity_center_permission_set_managed_policy" "policy" {
		instance_id         = "%s"
		permission_set_id   = samsungcloudplatformv2_iam_identity_center_permission_set.permission_set.id
		managed_policy_id   = "%s"
		managed_policy_name = "%s"
	}`, instanceId, permissionSetName, instanceId, managedPolicyId, managedPolicyName)
}
