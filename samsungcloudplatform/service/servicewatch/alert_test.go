package servicewatch_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlertResourceTest(t *testing.T) {
	alertName := fmt.Sprintf("test-acc-alert-%s", time.Now().Format("20060102_150405"))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// step1. Alert Create
				Config: testAccAlertCreate(alertName, "test-acc alert",
					map[string]string{
						"test-acc-key": "test-acc-value"}),
			},
			{
				// step2. Alert Update (Description, Evaluation Method)
				Config: testAccAlertUpdate(alertName, "test-acc alert modified", "Y",
					map[string]string{
						"test-acc-key": "test-acc-value"}),
			},
			{
				// step3. Alert Update (Activate)
				Config: testAccAlertUpdate(alertName, "test-acc alert modified", "N",
					map[string]string{
						"test-acc-key": "test-acc-value"}),
			},
			// step4. Delete Alert
		},
	})
}

func testAccAlertCreate(name string,
	description string,
	tags map[string]string) string {

	tagsJson, _ := json.Marshal(tags)
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_servicewatch_alert" "alert"{
				  name = "%s"
				  description = "%s"
				  type = "METRIC_ALERT"
				  level = "MIDDLE"
				  activated_yn = "Y"
				  namespace_name = "Virtual Server"
				  metric_name = "CPU Usage"
				  dimensions = [{key="resource_id", value="d5b49100-e3e3-4d10-b2e9-9da68aed7747"}]
				  period = 60
				  statistic = "AVG"
				  evaluation_count = 3
				  violation_count = 3
				  operator = "GTE"
				  threshold = 10
				  missing_data_option = "IGNORE"
				  tags = %s
			}
	`, name, description, tagsJson)
}

func testAccAlertUpdate(name string,
	description string,
	activatedYn string,
	tags map[string]string) string {

	tagsJson, _ := json.Marshal(tags)
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_servicewatch_alert" "alert"{
				  name = "%s"
				  description = "%s"
				  type = "METRIC_ALERT"
				  level = "MIDDLE"
				  activated_yn = "%s"
				  namespace_name = "Virtual Server"
				  metric_name = "CPU Usage"
				  dimensions = [{key="resource_id", value="d5b49100-e3e3-4d10-b2e9-9da68aed7747"}]
				  period = 60
				  statistic = "AVG"
				  evaluation_count = 3
				  violation_count = 3
				  operator = "RANGE"
				  upper_bound = 20
				  lower_bound = 10
				  missing_data_option = "IGNORE"
				  tags = %s
			}
	`, name, description, activatedYn, tagsJson)
}
