package client

import (
	"context"
	"encoding/json"
	"fmt"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"strings"
	"time"
)

const DefaultTimeout time.Duration = 120 * time.Minute

type Instance struct {
	Client *SCPClient
}

func IsTransientError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	if strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "temporary failure") ||
		strings.Contains(errStr, "network is unreachable") ||
		strings.Contains(errStr, "dial tcp") ||
		strings.Contains(errStr, "bad gateway") ||
		strings.Contains(errStr, "service unavailable") ||
		strings.Contains(errStr, "notresolvable") ||
		strings.Contains(errStr, "no route to host") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "server misbehaving") ||
		strings.Contains(errStr, "no address associated") ||
		strings.Contains(errStr, "host is unreachable") ||
		strings.Contains(errStr, "connection aborted") ||
		strings.Contains(errStr, "unexpected eof") ||
		strings.Contains(errStr, "failed to connect") ||
		strings.Contains(errStr, "use of closed network connection") ||
		strings.Contains(errStr, "too many open files") ||
		strings.Contains(errStr, "certificate verify failed") ||
		strings.Contains(errStr, "x509") {
		return true
	}

	return false
}

func WaitForStatus(ctx context.Context, client *SCPClient, pendingStates []string, targetStates []string, refreshFunc retry.StateRefreshFunc, timeout time.Duration, delay time.Duration, minTimeout time.Duration, maxConsecutiveErrors int) error {
	if timeout < 0 {
		timeout = DefaultTimeout
	}
	if delay < 0 {
		delay = 2 * time.Second
	}
	if minTimeout < 0 {
		minTimeout = 3 * time.Second
	}
	if maxConsecutiveErrors < 0 {
		maxConsecutiveErrors = 3
	}

	refreshWrapper := retryWithConsecutiveErrors(refreshFunc, maxConsecutiveErrors)

	stateConf := &retry.StateChangeConf{
		Pending:    pendingStates,
		Target:     targetStates,
		Refresh:    refreshWrapper,
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: minTimeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting : %s", err)
	}

	return nil
}

func retryWithConsecutiveErrors(refreshFunc retry.StateRefreshFunc, maxConsecutiveErrors int) retry.StateRefreshFunc {
	consecutiveErrors := 0

	return func() (interface{}, string, error) {
		result, state, err := refreshFunc()

		if err != nil {
			if IsTransientError(err) {
				consecutiveErrors++
				if consecutiveErrors >= maxConsecutiveErrors {
					return result, state, err
				}
				return result, state, nil
			}
			return result, state, err
		}

		consecutiveErrors = 0
		return result, state, nil
	}
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
