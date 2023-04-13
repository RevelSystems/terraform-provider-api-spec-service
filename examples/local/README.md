# Test it locally

To be able to test the provider locally, you will need to execute these steps:

1. Clone this repository.
2. In the root directory, run `make install` to build the provider and save it inside the plugins folder, where terraform will be able to grab it.
3. Go to `examples/local` folder.
4. Run `terraform init` to setup the `terraform-provider-api-spec-service`.
5. Get an M2M token and put it in the `main.tf` file, `m2m_token` parameter (M2M token can also be set through an environment variable - `M2M_TOKEN`).
6. (Optional) Run `terraform plan` to see if the provider was sucesfully installed.
7. Run `terraform apply` to upload OpenAPI specification, that is being defined in the `oas_file_path` parameter of the `oas_document` resource (in the example, it will try to upload `api-doc-valid.json` that is inside the `examples/local` folder).
8. (Optional) You can switch the API Spec Service environment by providing `environment` parameter to the provider. Options: `dev (default), qa, prod`.
9. When you are done, run `terraform destroy` to remove the uploaded OpenAPI specification.