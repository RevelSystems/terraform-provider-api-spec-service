terraform {
  required_providers {
    api-spec-service = {
      source = "registry.terraform.io/RevelSystems/api-spec-service"
    }
  }
}

provider "api-spec-service" {
  m2m_token   = "<M2M-token>"
  environment = "dev"
}

resource "oas_document" "testoas" {
  provider      = api-spec-service
  oas_file_path = "api-doc-valid.json"
}
