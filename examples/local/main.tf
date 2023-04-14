terraform {
  required_providers {
    api-spec-service = {
      source = "registry.terraform.io/RevelSystems/api-spec-service"
    }
  }
}

provider "api-spec-service" {
  m2m_token = "<M2M_TOKEN>"
  // Instead of m2m_token, you can provide client_id and client_secret to allow the provider to get the token
  client_id     = "<CLIENT_ID>"
  client_secret = "<CLIENT_SECRET>"
}

resource "oas_document" "testoas" {
  provider      = api-spec-service
  oas_file_path = "api-doc-valid.json"
}
