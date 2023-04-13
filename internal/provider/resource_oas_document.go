package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
				Description: "Path to OpenAPI specification file",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			// Optional inputs
			"oas_title": {
				Description: "The title inside the OpenAPI specification",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"oas_version": {
				Description: "The version inside the OpenAPI specification",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
		},
	}
})

func resourceOASDocumentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	oasFilePath := d.Get("oas_file_path").(string)
	oasFileExtension := filepath.Ext(oasFilePath)

	// Check OAS file extension
	contentType, contentTypeFound := fileExtensionToContentType[oasFileExtension]
	if !contentTypeFound {
		return diag.Errorf("Invalid file extension. Only json and yaml (yml) are allowed")
	}

	// Open OAS file
	oasFile, err := os.Open(oasFilePath)
	if err != nil {
		return diag.FromErr(err)
	}
	defer oasFile.Close()

	resp, err := client.UploadOASDocument(oasFile, contentType)
	if err != nil {
		return diag.FromErr(err)
	}

	oas_title := strings.ToLower(resp.Title)
	oas_version := strings.ToLower(resp.Version)

	d.SetId(fmt.Sprintf("%s#%s", oas_title, oas_version))
	d.Set("oas_title", oas_title)
	d.Set("oas_version", oas_version)

	return resourceOASDocumentRead(ctx, d, meta)
}

func resourceOASDocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	oasTitle := d.Get("oas_title").(string)
	oasVersion := d.Get("oas_version").(string)

	_, err := client.GetOASDocument(oasTitle, oasVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOASDocumentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	oasTitle := d.Get("oas_title").(string)
	oasVersion := d.Get("oas_version").(string)

	err := client.DeleteOASDocument(oasTitle, oasVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
