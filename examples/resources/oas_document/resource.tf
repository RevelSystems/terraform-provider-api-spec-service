resource "oas_document" "testoas_json" {
  provider      = api-spec-service
  oas_file_path = "api-doc-valid.json"
}

resource "oas_document" "testoas_yaml" {
  provider      = api-spec-service
  oas_file_path = "api-doc-valid.yaml"
}
