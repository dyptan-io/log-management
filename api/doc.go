// The server API is generated using /api/v1.yaml definition.

//go:generate -command oapi-codegen go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen
//go:generate oapi-codegen --generate chi-server,types --package api -o server.gen.go v1.yaml
//go:generate oapi-codegen --generate client --package api -o server_client.gen.go v1.yaml

package api
