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

func TestAccIamIdentityCenterPermissionSetCustomPolicyResourceTest(t *testing.T) {
	testPermissionSetName := fmt.Sprintf("test-acc-pscp-%d", time.Now().UnixNano())
	testCustomPolicyName := fmt.Sprintf("test-acc-custom-policy-%d", time.Now().UnixNano())
	existingInstanceId := getExistingInstanceId()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionSetCustomPolicyCreate(existingInstanceId, testPermissionSetName, testCustomPolicyName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_iam_identity_center_permission_set_custom_policy.policy", "permission_set_id"),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_permission_set_custom_policy.policy", "name", testCustomPolicyName),
				),
			},
		},
	})
}

func testAccPermissionSetCustomPolicyCreate(instanceId string, permissionSetName string, customPolicyName string) string {
	return fmt.Sprintf(`
	resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
		instance_id      = "%s"
		name             = "%s"
		description      = "test permission set for custom policy"
		session_duration = 3600
	}

	resource "samsungcloudplatformv2_iam_identity_center_permission_set_custom_policy" "policy" {
		instance_id       = "%s"
		permission_set_id = samsungcloudplatformv2_iam_identity_center_permission_set.permission_set.id
		name              = "%s"
	}`, instanceId, permissionSetName, instanceId, customPolicyName)
}
