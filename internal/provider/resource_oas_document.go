package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var environments = map[string]map[string]string{
	"dev": {
		"host": "http://api-spec-service.dev.int.revelup.io",
	},
	"qa": {
		"host": "http://api-spec-service.qa.int.revelup.io",
	},
	"prod": {
		"host": "http://api-spec-service.int.revelup.io",
	},
}

var fileExtensionToContentType = map[string]string{
	".json": "application/json",
	".yaml": "application/yaml",
	".yml":  "application/yaml",
}

var _ = registerResource("oas_document", func() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOASDocumentCreate,
		ReadContext:   resourceOASDocumentRead,
		DeleteContext: resourceOASDocumentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"oas_file_path": {
				Description: "Path to OAS file",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
})

func resourceOASDocumentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(Config)
	env := environments[config.Environment]

	oasFilePath := d.Get("oas_file_path").(string)
	oasFileExtension := filepath.Ext(oasFilePath)

	// Check OAS file extension
	contentType := fileExtensionToContentType[oasFileExtension]
	if contentType == "" {
		return diag.Errorf("Invalid file extension. Only json and yaml (yml) are allowed")
	}

	// Open OAS file
	oasFile, err := os.Open(oasFilePath)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get stats about the OAS file
	oasFileInfo, err := oasFile.Stat()
	if err != nil {
		return diag.FromErr(err)
	}

	// Create a new multipart form
	multiBuf := &bytes.Buffer{}
	multiWriter := multipart.NewWriter(multiBuf)

	// Add the file to the multipart form
	part, err := multiWriter.CreatePart(map[string][]string{
		"Content-Disposition": {"form-data; name=\"document\"; filename=\"" + oasFileInfo.Name() + "\""},
		"Content-Type":        {contentType},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = io.Copy(part, oasFile)
	if err != nil {
		return diag.FromErr(err)
	}

	// Close the multipart writer
	err = multiWriter.Close()
	if err != nil {
		return diag.FromErr(err)
	}

	// Create a new HTTP request with the multipart form as the body
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/documents", env["host"]), multiBuf)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.M2MToken))
	req.Header.Set("Content-Type", multiWriter.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	resp, err := config.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode > 299 {
		return diag.FromErr(fmt.Errorf(string(respBody)))
	}

	// TODO: construct id from title and version
	d.SetId("test-id")

	return nil
}

func resourceOASDocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceOASDocumentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
