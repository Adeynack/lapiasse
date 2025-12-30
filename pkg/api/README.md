# API Package

## Source of Truth

The API is described as an _OpenAPI Specification_ in `lapiasse.oas.yaml`. From the _OAS_, _Go_ code is generated into `lapiasse.gen.go`.

Here is the flow, in order:

| File                         | Format       | Function                                                                                                        |
| ---------------------------- | ------------ | --------------------------------------------------------------------------------------------------------------- |
| `lapiasse.oas.yaml`          | OpenAPI YAML | The _OpenAPI Specification_ (source of truth).                                                                  |
| `lapiasse.oapi-codegen.yaml` | YAML         | Configuration of the _OAPI Codegen_ tool, that generates _Go_ code from the _OAS_.                              |
| `lapiasse.gen.go`            | Go           | Resulting type-safe _Go_ code – interfaces and structs – that are then implemented in the `controller` package. |
