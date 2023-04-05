package provider

import "net/http"

type Config struct {
	M2MToken    string
	Environment string
}

var environments = map[string]Environment{
	"dev": {
		Host: "http://api-spec-service.dev.int.revelup.io",
	},
	"qa": {
		Host: "http://api-spec-service.qa.int.revelup.io",
	},
	"prod": {
		Host: "http://api-spec-service.int.revelup.io",
	},
}

func (c *Config) Client() (*Client, error) {
	httpClient := &http.Client{}

	client := &Client{
		M2MToken:    c.M2MToken,
		Environment: environments[c.Environment],
		HTTPClient:  httpClient,
	}

	return client, nil
}
