package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetAPIKey extracts an API Key from
// the headers of an HTTP request
// Example:
// Authorization ApiKey {insert apiKey here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("missing authorization header")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("malformed authorization header provided")
	}

	if vals[0] != "ApiKey" {
		return "", fmt.Errorf("malformed authorization header key provided: %s", vals[0])
	}

	key := vals[1]

	if len(key) < 64 {
		return "", errors.New("malformed api_key provided")
	}

	return key, nil
}
