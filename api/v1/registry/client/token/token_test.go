package token

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ivanilves/lstags/util/getenv"
)

var registry = getenv.String("DOCKER_REGISTRY_HOSTNAME", "registry.hub.docker.com")

func TestExtractAuthHeader(t *testing.T) {
	testCases := []struct {
		input    []string
		expected authHeader
		isValid  bool
	}{
		{nil, "None realm=none", true},
		{[]string{}, "None realm=none", true},
		{[]string{"Basic realm=artifactory,scope=global"}, "Basic realm=artifactory,scope=global", true},
		{[]string{"Basic realm=artifactory", "expires=10s"}, "Basic realm=artifactory", true},
		{[]string{"megustapollofrito"}, "", false},
	}

	assert := assert.New(t)

	for _, tc := range testCases {
		output, err := extractAuthHeader(tc.input)

		assert.Equal(tc.expected, output)

		if tc.isValid {
			assert.Nil(err)
		} else {
			assert.NotNil(err)
		}
	}
}

func TestGetAuthMethod(t *testing.T) {
	testCases := []struct {
		input    authHeader
		expected string
	}{
		{"None realm=none", "None"},
		{"Basic realm=artifactory,anonymous=no", "Basic"},
	}

	assert := assert.New(t)

	for _, tc := range testCases {
		assert.Equal(tc.expected, getAuthMethod(tc.input))
	}
}

func TestGetAuthParams(t *testing.T) {
	testCases := []struct {
		input    authHeader
		expected map[string]string
	}{
		{"None realm=none", map[string]string{"realm": "none"}},
		{"Basic realm=quay.io,anonymous=no", map[string]string{"realm": "quay.io", "anonymous": "no"}},
		{"Basic realm=127.0.0.1,bogustoken", map[string]string{"realm": "127.0.0.1"}},
		{"Basic realm=artifactory,x=y=z", map[string]string{"realm": "artifactory"}},
	}

	assert := assert.New(t)

	for _, tc := range testCases {
		assert.Equal(tc.expected, getAuthParams(tc.input))
	}
}

func TestDockerHub(t *testing.T) {
	url := "https://registry.hub.docker.com/v2"

	tk, _ := New(url, "", "")

	assert := assert.New(t)

	assert.Equal("Bearer", tk.Method())
	assert.Len(tk.String(), 1963)
}
