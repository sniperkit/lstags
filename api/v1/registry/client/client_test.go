package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"time"

	"github.com/ivanilves/lstags/util/getenv"
)

var registry = getenv.String("DOCKER_REGISTRY_HOSTNAME", "registry.hub.docker.com")

func TestNewWithEmptyConfig(t *testing.T) {
	assert := assert.New(t)

	c, _ := New(registry, Config{})

	assert.Equal(DefaultConcurrentRequests, c.Config.ConcurrentRequests)
	assert.Equal(DefaultRetryDelay, c.Config.RetryDelay)
}

func TestNewWithInvalidConfig(t *testing.T) {
	assert := assert.New(t)

	_, err := New(registry, Config{ConcurrentRequests: 9000})

	assert.NotNil(err)
}

func TestNewWithDefinedConfig(t *testing.T) {
	concurrentRequests := 77
	retryDelay := 5 * time.Second

	assert := assert.New(t)

	c, _ := New(
		registry,
		Config{
			ConcurrentRequests: concurrentRequests,
			RetryDelay:         retryDelay,
		},
	)

	assert.Equal(concurrentRequests, c.Config.ConcurrentRequests)
	assert.Equal(retryDelay, c.Config.RetryDelay)
}

func TestWebSchemeDefault(t *testing.T) {
	assert := assert.New(t)

	c, _ := New(registry, Config{})

	assert.Equal("https://", c.webScheme())
}

func TestWebSchemeInsecure(t *testing.T) {
	assert := assert.New(t)

	c, _ := New(registry, Config{IsInsecure: true})

	assert.Equal("http://", c.webScheme())
}

func TestPing(t *testing.T) {
	assert := assert.New(t)

	c, _ := New(registry, Config{})

	assert.Nil(c.Ping())
}

func TestPingNoHost(t *testing.T) {
	assert := assert.New(t)

	c, _ := New("i.do.not.exist.sorry", Config{})

	assert.NotNil(c.Ping())
}

func TestPingNotARegistry(t *testing.T) {
	assert := assert.New(t)

	c, _ := New("www.google.com", Config{})

	assert.NotNil(c.Ping())
}

func TestLoginAnonymous(t *testing.T) {
	assert := assert.New(t)

	cli, _ := New(registry, Config{})

	assert.False(cli.IsLoggedIn())

	assert.Nil(cli.Login("", ""))

	assert.True(cli.IsLoggedIn())
}

func TestLoginAuthenticated(t *testing.T) {
	username := getenv.String("DOCKER_REGISTRY_USERNAME", "")
	password := getenv.String("DOCKER_REGISTRY_PASSWORD", "")

	if username == "" && password == "" {
		t.Skip("DOCKER_REGISTRY_USERNAME and DOCKER_REGISTRY_PASSWORD env variables not defined")
	}

	assert := assert.New(t)

	cli, _ := New(registry, Config{})

	assert.False(cli.IsLoggedIn())

	assert.Nil(cli.Login(username, password))

	assert.True(cli.IsLoggedIn())
}

func TestLoginInvalid(t *testing.T) {
	assert := assert.New(t)

	cli, _ := New(registry, Config{})

	assert.False(cli.IsLoggedIn())

	assert.NotNil(cli.Login("notavaliduser", "omgnotavalidpassword"))

	assert.False(cli.IsLoggedIn())
}
