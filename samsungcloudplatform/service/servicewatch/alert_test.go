package servicewatch_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccAlertAddress = "samsungcloudplatformv2_servicewatch_alert.alert"

// testAccAlertRecipients reads the notification recipients to exercise from the
// environment, since valid ids are account specific. The recipient block is left out
// of the configuration when SCP_TF_TEST_RECIPIENT_IDS is unset, so the rest of the
// alert test still runs.
//
//	SCP_TF_TEST_RECIPIENT_IDS   comma separated recipient ids
//	SCP_TF_TEST_RECIPIENT_TYPE  USER (default) or GROUP
func testAccAlertRecipients() ([]string, string) {
	raw := os.Getenv("SCP_TF_TEST_RECIPIENT_IDS")
	if raw == "" {
		return nil, ""
	}

	ids := make([]string, 0)
	for _, id := range strings.Split(raw, ",") {
		if id = strings.TrimSpace(id); id != "" {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil, ""
	}

	recipientType := os.Getenv("SCP_TF_TEST_RECIPIENT_TYPE")
	if recipientType == "" {
		recipientType = "USER"
	}
	return ids, recipientType
}

// testAccAlertRecipientBlock renders the recipient attributes, or nothing when unset.
func testAccAlertRecipientBlock() string {
	ids, recipientType := testAccAlertRecipients()
	if len(ids) == 0 {
		return ""
	}
	return fmt.Sprintf("recipient_ids = %s\n\t\t\t\t  recipient_type = %q",
		hclStringList(ids), recipientType)
}

func testAccAlertRecipientChecks() []resource.TestCheckFunc {
	ids, recipientType := testAccAlertRecipients()
	if len(ids) == 0 {
		return nil
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(testAccAlertAddress, "recipient_type", recipientType),
		resource.TestCheckResourceAttr(testAccAlertAddress, "recipient_ids.#", strconv.Itoa(len(ids))),
	}
	for i, id := range ids {
		checks = append(checks,
			resource.TestCheckResourceAttr(testAccAlertAddress, fmt.Sprintf("recipient_ids.%d", i), id))
	}
	return checks
}

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
				Check: resource.ComposeAggregateTestCheckFunc(testAccAlertRecipientChecks()...),
			},
			{
				// step2. Alert Update (Description, Evaluation Method)
				Config: testAccAlertUpdate(alertName, "test-acc alert modified", "Y",
					map[string]string{
						"test-acc-key": "test-acc-value"}),
				Check: resource.ComposeAggregateTestCheckFunc(testAccAlertRecipientChecks()...),
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

// hclStringList renders a Go slice as an HCL list literal.
func hclStringList(items []string) string {
	quoted := make([]string, 0, len(items))
	for _, item := range items {
		quoted = append(quoted, fmt.Sprintf("%q", item))
	}
	return "[" + strings.Join(quoted, ", ") + "]"
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
				  %s
				  tags = %s
			}
	`, name, description, testAccAlertRecipientBlock(), tagsJson)
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
				  %s
				  tags = %s
			}
	`, name, description, activatedYn, testAccAlertRecipientBlock(), tagsJson)
}
