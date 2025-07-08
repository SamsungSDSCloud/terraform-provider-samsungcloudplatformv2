package client

import (
	"context"
	"encoding/json"
	"fmt"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strings"
	"time"
)

const DefaultTimeout time.Duration = 120 * time.Minute

type Instance struct {
	Client *SCPClient
}

func WaitForStatus(ctx context.Context, client *SCPClient, pendingStates []string, targetStates []string, refreshFunc resource.StateRefreshFunc) error {
	stateConf := &resource.StateChangeConf{
		Pending:    pendingStates,
		Target:     targetStates,
		Refresh:    refreshFunc,
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting : %s", err)
	}

	return nil
}

func GetDetailFromError(err error) string {
	var data map[string]interface{}
	body := err.(*scpsdk.GenericOpenAPIError).Body()
	err = json.Unmarshal(body, &data)

	var details []string
	errors := data["errors"]

	for _, err := range errors.([]interface{}) {
		errorMap := err.(map[string]interface{})
		detail := errorMap["detail"]
		switch detail.(type) {
		case string:
			details = append(details, detail.(string))
		case []interface{}:
			for _, d := range detail.([]interface{}) {
				details = append(details, d.(string))
			}
		}
	}

	return strings.Join(details, ", ")
}
