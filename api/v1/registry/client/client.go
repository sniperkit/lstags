// Package client provides Docker registry client API
package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ivanilves/lstags/api/v1/registry/client/token"
)

// DefaultConcurrentRequests will be used if no explicit ConcurrentRequests configured
var DefaultConcurrentRequests = 32

// DefaultRetryDelay will be used if no explicit RetryDelay configured
var DefaultRetryDelay = 30 * time.Second

// MaxConcurrentRequests is a hard limit for simultaneous registry requests
const MaxConcurrentRequests = 256

// RegistryClient is an abstraction to wrap logic of working with Docker registry
// incl. connection, authentification, authorization, obtaining information etc...
type RegistryClient struct {
	registry string
	// Config has general configuration of the registry client instance
	Config Config
	// Token is an authentication token obtained after registry login
	Token token.Token
}

// Config has configuration parameters for RegistryClient creation
type Config struct {
	// ConcurrentRequests defines how much requests to registry we could run concurrently
	ConcurrentRequests int
	// WaitBetween defines how much we will wait between batches of requests
	WaitBetween time.Duration
	// RetryRequests defines how much retries we will do to the failed HTTP request
	RetryRequests int
	// RetryDelay defines how much we will wait between failed HTTP request and retry
	RetryDelay time.Duration
	// TraceRequests sets if we will print out registry HTTP request traces
	TraceRequests bool
	// IsInsecure sets if we want to communicate registry over plain HTTP instead of HTTPS
	IsInsecure bool
}

// New creates and validates new RegistryClient instance
func New(registry string, config Config) (*RegistryClient, error) {
	if config.ConcurrentRequests == 0 {
		config.ConcurrentRequests = DefaultConcurrentRequests
	}

	if config.RetryDelay == 0 {
		config.RetryDelay = DefaultRetryDelay
	}

	if config.ConcurrentRequests > MaxConcurrentRequests {
		err := fmt.Errorf(
			"Could not run more than %d concurrent requests (%d configured)",
			MaxConcurrentRequests,
			config.ConcurrentRequests,
		)

		return nil, err
	}

	return &RegistryClient{registry: registry, Config: config}, nil
}

func (cli *RegistryClient) webScheme() string {
	if cli.Config.IsInsecure {
		return "http://"
	}

	return "https://"
}

// URL formats a valid URL for the V2 registry
func (cli *RegistryClient) URL() string {
	return cli.webScheme() + cli.registry + "/v2/"
}

// Ping checks basic connectivity to the registry
func (cli *RegistryClient) Ping() error {
	resp, err := http.Get(cli.URL())
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 401 {
		return fmt.Errorf("Unexpected status: %s", resp.Status)
	}

	return nil
}

// Login logs in to the registry (return error, if failed)
func (cli *RegistryClient) Login(username, password string) error {
	tk, err := token.New(cli.URL(), username, password)
	if err != nil {
		return err
	}

	cli.Token = tk

	return nil
}

// IsLoggedIn indicates if we are logged in to registry or not
func (cli *RegistryClient) IsLoggedIn() bool {
	return cli.Token != nil
}

// GetRepositories() {}

// GetTags(repository string) {}
