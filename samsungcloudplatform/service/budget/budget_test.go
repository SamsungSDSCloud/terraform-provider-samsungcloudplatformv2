package budget_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceType = "samsungcloudplatformv2_budget_budget"
	resourceName = "samsungcloudplatformv2_budget_budget.budget" // type.name
)

func TestAccBudgetResourceTest(t *testing.T) {
	budgetName := fmt.Sprintf("tf-acc-b-%s", acctest.RandString(11))
	currentStartMonth := time.Now().Format("2006-01")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},

		Steps: []resource.TestStep{
			{
				// 1. Budget Create
				Config: testAccBudgetConfig(budgetName, 100000, currentStartMonth, "MONTHLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", budgetName),
					resource.TestCheckResourceAttr(resourceName, "amount", "100000"),
					resource.TestCheckResourceAttr(resourceName, "notifications.is_use_notification", "false"),
				),
			},
			{
				// 2. Update: Amount & Notifications
				Config: testAccBudgetNotificationConfig(budgetName, 250000, currentStartMonth, "MONTHLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "amount", "250000"),
					resource.TestCheckResourceAttr(resourceName, "notifications.is_use_notification", "true"),
					resource.TestCheckResourceAttr(resourceName, "notifications.thresholds.0", "80"),
				),
			},
			{
				// 3. Update: Prevention
				Config: testAccBudgetPreventionConfig(budgetName, 250000, currentStartMonth, "MONTHLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "prevention.is_use_prevention", "true"),
					resource.TestCheckResourceAttr(resourceName, "prevention.threshold", "80"),
				),
			},
		},
	})
}

// common resource generation functions
func testAccBudgetBaseConfig(name string, amount int, startMonth string, unit string) string {
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_budget_budget" "budget" {
				name        = "%s"
				amount      = %d
				start_month = "%s"
				unit        = "%s"
			`, name, amount, startMonth, unit)
}

func testAccBudgetConfig(name string, amount int, startMonth string, unit string) string {
	return testAccBudgetBaseConfig(name, amount, startMonth, unit) + `
				notifications = {
					is_use_notification = false
				}
				prevention = {
					is_use_prevention = false
				}
			}`
}

func testAccBudgetNotificationConfig(name string, amount int, startMonth string, unit string) string {
	return testAccBudgetBaseConfig(name, amount, startMonth, unit) + `
				notifications = {
					is_use_notification = true
					thresholds = [80]
					notification_send_period      = "FIRST"
					receivers    = ["test@example.com"]
				}
				prevention = {
					is_use_prevention = false
				}
			}`
}

func testAccBudgetPreventionConfig(name string, amount int, startMonth string, unit string) string {
	return testAccBudgetBaseConfig(name, amount, startMonth, unit) + `
				notifications = {
					is_use_notification = false
				}
				prevention = {
					is_use_prevention = true
					threshold    = 80
					receivers = ["test@example.com"]
				}
			}`
}
