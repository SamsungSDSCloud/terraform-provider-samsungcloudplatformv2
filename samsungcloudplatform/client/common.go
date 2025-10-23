package client

import (
	"context"
	"encoding/json"
	"fmt"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
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

	// Check if the error is of type *scpsdk.GenericOpenAPIError
	if genericErr, ok := err.(*scpsdk.GenericOpenAPIError); ok {
		body := genericErr.Body()
		err := json.Unmarshal(body, &data)
		if err != nil {
			return "Error parsing error body: " + err.Error()
		}
	} else {
		// If the error is not of type *scpsdk.GenericOpenAPIError, return a generic error message
		return "Unknown error: " + err.Error()
	}

	var details []string
	errors, ok := data["errors"].([]interface{})
	if !ok {
		return "Invalid error data"
	}

	for _, err := range errors {
		errorMap, ok := err.(map[string]interface{})
		if !ok {
			continue
		}
		detail, ok := errorMap["detail"]
		if !ok {
			continue
		}
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
