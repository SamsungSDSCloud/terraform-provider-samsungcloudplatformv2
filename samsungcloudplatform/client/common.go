package client

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	StateCreating = "CREATING"
	StateActive   = "ACTIVE"
	StateEditing  = "EDITING"
	StateDeleting = "DELETING"
	StateDeleted  = "DELETED"
	StateError    = "ERROR"
)

const DefaultTimeout time.Duration = 120 * time.Minute

const DefaultWaiterTimeout time.Duration = 30 * time.Minute

// WaiterFunc is a lifecycle-specific waiter that receives a refresh function.
type WaiterFunc func(context.Context, retry.StateRefreshFunc) error

// WaiterFuncWithStates is a waiter that accepts custom pending and target states.
type WaiterFuncWithStates func(context.Context, []string, []string, retry.StateRefreshFunc) error

// WaitForResource waits for pendingStates → targetStates with default timeout.
// If StateDeleted is in targetStates, 404 errors are automatically mapped to StateDeleted.
func WaitForResource(ctx context.Context, pendingStates []string, targetStates []string, refreshFunc retry.StateRefreshFunc) error {
	// Wrap to catch 404 errors when waiting for deletion
	if slices.Contains(targetStates, StateDeleted) {
		originalRefresh := refreshFunc
		refreshFunc = func() (interface{}, string, error) {
			result, state, err := originalRefresh()
			if err != nil && strings.HasPrefix(err.Error(), "404") {
				return struct{}{}, StateDeleted, nil
			}
			return result, state, err
		}
	}
	return WaitForStatus(ctx, nil, pendingStates, targetStates, refreshFunc, DefaultWaiterTimeout, -1, -1, -1)
}


// WaitForResourceCreated waits for CREATING → ACTIVE.
func WaitForResourceCreated(ctx context.Context, refreshFunc retry.StateRefreshFunc) error {
	return WaitForResource(ctx, []string{StateCreating}, []string{StateActive}, refreshFunc)
}

// WaitForResourceUpdated waits for EDITING → ACTIVE.
func WaitForResourceUpdated(ctx context.Context, refreshFunc retry.StateRefreshFunc) error {
	return WaitForResource(ctx, []string{StateEditing}, []string{StateActive}, refreshFunc)
}

// WaitForResourceDeleted waits for DELETING → gone (404).
func WaitForResourceDeleted(ctx context.Context, refreshFunc retry.StateRefreshFunc) error {
	return WaitForResource(ctx, []string{StateDeleting}, []string{StateDeleted}, refreshFunc)
}

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
				if s, ok := d.(string); ok {
					details = append(details, s)
				}
			}
		}
	}

	return strings.Join(details, ", ")
}
