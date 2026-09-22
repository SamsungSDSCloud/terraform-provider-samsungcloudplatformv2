package iamidentitycenter_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	iamidentitycentercommon "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/iamidentitycenter"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// getExistingInstanceId looks up an existing IAM Identity Center instance via the API.
// If found, returns its ID so the test can reuse it.
// If not found, returns "" and the test will create one dynamically.
func getExistingInstanceId() string {
	region := os.Getenv("SCP_REGION")
	if region == "" {
		region = "kr-west1"
	}
	client, err := iamidentitycentercommon.SharedClientForRegion(region)
	if err != nil {
		return ""
	}
	instances, err := client.Client.IamIdentityCenter.ListInstances(context.Background(), 0, 0, "")
	if err != nil {
		return ""
	}
	allInstances := instances.GetInstances()
	if len(allInstances) > 0 {
		return allInstances[0].Id
	}
	return ""
}

func TestAccIamIdentityCenterInstanceResourceTest(t *testing.T) {
	if getExistingInstanceId() != "" {
		t.Skip("An IAM Identity Center instance already exists. Only one instance is allowed per account. Delete the existing instance manually or via sweeper and re-run this test.")
	}

	testInstanceName := fmt.Sprintf("test-acc-%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceCreate(
					testInstanceName,
					"test acc IAM Identity Center Instance description",
					"IDENTITY_CENTER_DIRECTORY",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_instance.instance", "name", testInstanceName),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_instance.instance", "description", "test acc IAM Identity Center Instance description"),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_instance.instance", "identity_store_type", "IDENTITY_CENTER_DIRECTORY"),
				),
			},
			{
				Config: testAccInstanceUpdate(
					testInstanceName,
					"test acc IAM Identity Center Instance updated description",
					"IDENTITY_CENTER_DIRECTORY",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_instance.instance", "name", testInstanceName),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_instance.instance", "description", "test acc IAM Identity Center Instance updated description"),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_instance.instance", "identity_store_type", "IDENTITY_CENTER_DIRECTORY"),
				),
			},
			{
				Config: testAccInstanceDelete(),
			},
		},
	})
}

func testAccInstanceCreate(name string, description string, identityStoreType string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name = "%s"
			description = "%s"
			identity_store_type = "%s"
		}`, name, description, identityStoreType)
}

func testAccInstanceUpdate(name string, description string, identityStoreType string) string {
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name = "%s"
			description = "%s"
			identity_store_type = "%s"
		}`, name, description, identityStoreType)
}

func testAccInstanceDelete() string {
	// Return a valid HCL config that does not define any managed resources.
	// Terraform detects the instance is in state but absent from config and destroys it.
	return `
	locals {
		done = "true"
	}
	`
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_iam_identity_center_instance", &resource.Sweeper{
		Name: "samsungcloudplatformv2_iam_identity_center_instance",
		F:    sweepInstance,
	})
}

func sweepInstance(region string) error {
	scpClient, err := iamidentitycentercommon.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	instances, err := scpClient.Client.IamIdentityCenter.ListInstances(context.Background(), 0, 0, "")
	if err != nil {
		return fmt.Errorf("error listing IAM Identity Center instances: %v", err)
	}

	for _, instance := range instances.GetInstances() {
		// Get instance detail to check the name
		instanceDetail, _, err := scpClient.Client.IamIdentityCenter.GetInstance(context.Background(), instance.Id)
		if err != nil {
			fmt.Printf("warning: failed to get instance %s detail: %v (continuing...)\n", instance.Id, err)
			continue
		}

		// Only delete instances created by tests (name starts with "test-acc-")
		if !strings.HasPrefix(instanceDetail.Instance.Name, "test-acc-") {
			fmt.Printf("skipping instance %s (not a test instance)\n", instance.Id)
			continue
		}

		// Retry delete with delay for instances with dependent resources
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			_, err = scpClient.Client.IamIdentityCenter.DeleteInstance(context.Background(), instance.Id)
			if err == nil {
				fmt.Printf("deleted IAM Identity Center instance: %s\n", instance.Id)
				break
			}
			if i < maxRetries-1 && (strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "Conflict")) {
				fmt.Printf("retrying delete for instance %s (attempt %d/%d)\n", instance.Id, i+1, maxRetries)
				time.Sleep(5 * time.Second)
			} else {
				fmt.Printf("warning: failed to delete instance %s: %v (continuing...)\n", instance.Id, err)
				break
			}
		}
	}

	return nil
}
