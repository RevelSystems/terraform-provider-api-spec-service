package provider

import "net/http"

type Config struct {
	M2MToken    string
	Environment string
	HTTPClient  *http.Client
}
