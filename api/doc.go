// The server API is generated using /api/v1.yaml definition.

//go:generate -command oapi-codegen go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
//go:generate oapi-codegen --generate std-http-server,strict-server,types --package api -o server.gen.go v1.yaml
//go:generate oapi-codegen --generate client --package api -o server_client.gen.go v1.yaml

package api

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)
