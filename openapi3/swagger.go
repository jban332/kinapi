package openapi3

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/jsoninfo"
)

type Swagger struct {
	ExtensionProps
	OpenAPI      string               `json:"openapi"` // Required
	Info         *Info                `json:"info"`    // Required
	Servers      Servers              `json:"servers,omitempty"`
	Paths        Paths                `json:"paths"` // Required
	Components   Components           `json:"components,omitempty"`
	Security     SecurityRequirements `json:"security,omitempty"`
	ExternalDocs *ExternalDocs        `json:"externalDocs,omitempty"`
}

func (swagger *Swagger) MarshalJSON() ([]byte, error) {
	return jsoninfo.MarshalStrictStruct(swagger)
}

func (swagger *Swagger) UnmarshalJSON(data []byte) error {
	return jsoninfo.UnmarshalStrictStruct(data, swagger)
}

func (swagger *Swagger) AddOperation(path string, method string, operation *Operation) {
	paths := swagger.Paths
	if paths == nil {
		paths = make(Paths)
		swagger.Paths = paths
	}
	pathItem := paths[path]
	if pathItem == nil {
		pathItem = &PathItem{}
		paths[path] = pathItem
	}
	pathItem.SetOperation(method, operation)
}

func (swagger *Swagger) AddServer(server *Server) {
	swagger.Servers = append(swagger.Servers, server)
}

func (swagger *Swagger) Validate(c context.Context) error {
	if swagger.OpenAPI == "" {
		return fmt.Errorf("Variable 'openapi' must be a non-empty JSON string")
	}
	if err := swagger.Components.Validate(c); err != nil {
		return fmt.Errorf("Error when validating Components: %s", err.Error())
	}
	if v := swagger.Security; v != nil {
		if err := v.Validate(c); err != nil {
			return fmt.Errorf("Error when validating Security: %s", err.Error())
		}
	}
	if v := swagger.Servers; v != nil {
		if err := v.Validate(c); err != nil {
			return fmt.Errorf("Error when validating Servers: %s", err.Error())
		}
	}
	if v := swagger.Paths; v != nil {
		if err := v.Validate(c); err != nil {
			return fmt.Errorf("Error when validating Paths: %s", err.Error())
		}
	} else {
		return fmt.Errorf("Variable 'paths' must be a JSON object")
	}
	if v := swagger.Info; v != nil {
		if err := v.Validate(c); err != nil {
			return fmt.Errorf("Error when validating Info: %s", err.Error())
		}
	} else {
		return fmt.Errorf("Variable 'info' must be a JSON object")
	}
	return nil
}
