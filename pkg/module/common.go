package module

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	xhttp "github.com/ibm-security-verify/verifyctl/pkg/util/http"
)

type VerifyError struct {
	MessageID          string `json:"messageId" yaml:"messageId"`
	MessageDescription string `json:"messageDescription" yaml:"messageDescription"`
}

func HandleCommonErrors(ctx context.Context, response *xhttp.Response, defaultError string) error {
	if response.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Login again.")
	}

	if response.StatusCode == http.StatusForbidden {
		return fmt.Errorf("You are not allowed to make this request. Check the client or application entitlements.")
	}

	if response.StatusCode == http.StatusBadRequest {
		var errorMessage VerifyError
		if err := json.Unmarshal(response.Body, &errorMessage); err != nil {
			return fmt.Errorf("bad request: %s", defaultError)
		}
		// If the expected fields are not populated, return the raw response body.
		if errorMessage.MessageID == "" && errorMessage.MessageDescription == "" {
			return fmt.Errorf("bad request: %s", string(response.Body))
		}
		return fmt.Errorf("%s %s", errorMessage.MessageID, errorMessage.MessageDescription)
	}

	if response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Resource not found")
	}

	return nil
}
