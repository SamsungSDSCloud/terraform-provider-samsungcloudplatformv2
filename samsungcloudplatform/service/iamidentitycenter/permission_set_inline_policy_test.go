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

func TestAccIamIdentityCenterPermissionSetInlinePolicyResourceTest(t *testing.T) {
	testPermissionSetName := fmt.Sprintf("test-acc-psin-%d", time.Now().UnixNano())
	testPolicyDocument := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`
	existingInstanceId := getExistingInstanceId()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionSetInlinePolicyCreate(existingInstanceId, testPermissionSetName, testPolicyDocument),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_iam_identity_center_permission_set_inline_policy.policy", "permission_set_id"),
				),
			},
		},
	})
}

func testAccPermissionSetInlinePolicyCreate(instanceId string, permissionSetName string, policyDocument string) string {
	return fmt.Sprintf(`
	resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
		instance_id      = "%s"
		name             = "%s"
		description      = "test permission set for inline policy"
		session_duration = 3600
	}

	resource "samsungcloudplatformv2_iam_identity_center_permission_set_inline_policy" "policy" {
		instance_id       = "%s"
		permission_set_id = samsungcloudplatformv2_iam_identity_center_permission_set.permission_set.id
		policy_document   = <<EOT
%s
EOT
	}`, instanceId, permissionSetName, instanceId, policyDocument)
}
