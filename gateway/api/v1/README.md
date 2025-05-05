# API Code Generation

This directory contains OpenAPI specification and configuration for the Gateway API.

## Prerequisites

Install the oapi-codegen tool:

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

## Files

- `openapi.yaml` - OpenAPI 3.0 specification
- `oapi-codegen.yaml` - Configuration for code generation
- Generated code is located in `../../internal/controller/http/v1/handler/models.gen.go` and `../../internal/controller/http/v1/swagger.gen.go`

## Generating Code

To generate/update the API code, run from the root directory:

```bash
make generate-api
```

This command will update the generated code with:

- Chi server implementation
- Request/response models
- Validation logic

## Configuration

The `oapi-models-codegen.yaml` configures code generation with:

```yaml
package: v1
generate:
  models: true # Generates data models
output: ../../internal/controller/http/v1/handler/models.gen.go # Output file location
```

The `oapi-swagger-codegen.yaml` configures code generation with:

```yaml
package: v1
generate:
  embedded-spec: true # Embeds OpenAPI spec in generated code
output: ../../internal/controller/http/v1/swagger.gen.go # Output file location
```
