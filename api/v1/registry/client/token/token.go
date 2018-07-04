package token

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ivanilves/lstags/api/v1/registry/client/token/basic"
	"github.com/ivanilves/lstags/api/v1/registry/client/token/bearer"
	"github.com/ivanilves/lstags/api/v1/registry/client/token/none"
)

// Token is an abstraction for aggregated token-related information we get from authentication services
type Token interface {
	Method() string
	String() string
	ExpiresIn() int
}

type authHeader string

func extractAuthHeader(hh []string) (authHeader, error) {
	if len(hh) == 0 {
		return "None realm=none", nil
	}

	h := hh[0]

	if len(strings.SplitN(h, " ", 2)) != 2 {
		return "", errors.New("Unexpected 'Www-Authenticate' header: " + h)
	}

	return authHeader(h), nil
}

func getAuthMethod(h authHeader) string {
	return strings.SplitN(string(h), " ", 2)[0]
}

func getAuthParams(h authHeader) map[string]string {
	params := make(map[string]string)

	paramString := strings.SplitN(string(h), " ", 2)[1]

	for _, keyValueString := range strings.Split(paramString, ",") {
		kv := strings.Split(keyValueString, "=")
		if len(kv) == 2 {
			params[kv[0]] = strings.Trim(kv[1], "\"")
		}
	}

	return params
}

// New creates a new instance of Token in two steps:
// * detects authentication type ("Bearer", "Basic" or "None")
// * delegates actual authentication to the type-specific implementation
func New(url, username, password string) (Token, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	authHeader, err := extractAuthHeader(resp.Header["Www-Authenticate"])
	if err != nil {
		return nil, err
	}

	method := getAuthMethod(authHeader)
	params := getAuthParams(authHeader)

	switch method {
	case "None":
		return none.RequestToken()
	case "Basic":
		return basic.RequestToken(url, username, password)
	case "Bearer":
		return bearer.RequestToken(username, password, params)
	default:
		return nil, errors.New("Unknown authentication method: " + method)
	}
}
