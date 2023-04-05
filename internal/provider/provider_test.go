package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"api-spec-service": func() (*schema.Provider, error) {
		return New()(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New()().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = New()()
}

func testPreCheck(t *testing.T) {
	if v := os.Getenv("M2M_TOKEN"); v == "" {
		t.Fatal("M2M_TOKEN must be set for acceptance tests")
	}
}
