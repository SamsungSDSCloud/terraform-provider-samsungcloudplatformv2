package resourcemanager_test

import (
	"encoding/json"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccResourceGroupResourceTest(t *testing.T) {
	resourceTypes := []string{"iam:policy"}
	tags := map[string]string{
		"key":   "test-acc-key",
		"value": "test-acc-value",
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupCreate(
					"test-acc-resource_group",
					"test acc resource_group description",
					resourceTypes,
					tags),
			},
			{
				Config: testAccResourceGroupUpdate(
					"test-acc-resource_group",
					"test acc resource_group update description",
					resourceTypes,
					tags),
			},
			{
				Config: testAccResourceGroupDelete(
					"test-acc-resource_group",
					"test acc resource_group update description",
					resourceTypes,
					tags),
			},
		},
	})
}

func testAccResourceGroupCreate(name string, description string, resourceTypes []string, tags map[string]string) string {
	resourceTypesJson, _ := json.Marshal(resourceTypes)
	tagsJson, _ := json.Marshal(tags)

	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_resourcemanager_resource_group" "resource_group"{
				name = "%s"
				description = "%s"
				resource_types = %s
				tags = %s
			}
	`, name, description, resourceTypesJson, tagsJson)
}

func testAccResourceGroupUpdate(name string, description string, resourceTypes []string, tags map[string]string) string {
	resourceTypesJson, _ := json.Marshal(resourceTypes)
	tagsJson, _ := json.Marshal(tags)

	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_resourcemanager_resource_group" "resource_group"{
				name = "%s"
				description = "%s"
				resource_types = %s
				tags = %s
			}
	`, name, description, resourceTypesJson, tagsJson)
}

func testAccResourceGroupDelete(name string, description string, resourceTypes []string, tags map[string]string) string {
	resourceTypesJson, _ := json.Marshal(resourceTypes)
	tagsJson, _ := json.Marshal(tags)

	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_resourcemanager_resource_group" "resource_group"{
				name = "%s"
				description = "%s"
				resource_types = %s
				tags = %s
			}
	`, name, description, resourceTypesJson, tagsJson)
}
