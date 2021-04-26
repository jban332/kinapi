package openapi3

import (
	"context"
	"errors"
	"fmt"

	"github.com/getkin/kin-openapi/jsoninfo"
)

// Swagger is the root of an OpenAPI v3 document
type Swagger struct {
	ExtensionProps
	OpenAPI      string               `json:"openapi" yaml:"openapi"` // Required
	Components   Components           `json:"components,omitempty" yaml:"components,omitempty"`
	Info         *Info                `json:"info" yaml:"info"`   // Required
	Paths        Paths                `json:"paths" yaml:"paths"` // Required
	Security     SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
	Servers      Servers              `json:"servers,omitempty" yaml:"servers,omitempty"`
	Tags         Tags                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDocs        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

func (doc *Swagger) MarshalJSON() ([]byte, error) {
	return jsoninfo.MarshalStrictStruct(doc)
}

func (doc *Swagger) UnmarshalJSON(data []byte) error {
	return jsoninfo.UnmarshalStrictStruct(data, doc)
}

func (doc *Swagger) AddOperation(path string, method string, operation *Operation) {
	paths := doc.Paths
	if paths == nil {
		paths = make(Paths)
		doc.Paths = paths
	}
	pathItem := paths[path]
	if pathItem == nil {
		pathItem = &PathItem{}
		paths[path] = pathItem
	}
	pathItem.SetOperation(method, operation)
}

func (doc *Swagger) AddServer(server *Server) {
	doc.Servers = append(doc.Servers, server)
}

func (doc *Swagger) Validate(c context.Context) error {
	if doc.OpenAPI == "" {
		return errors.New("value of openapi must be a non-empty string")
	}

	// NOTE: only mention info/components/paths/... key in this func's errors.

	{
		wrap := func(e error) error { return fmt.Errorf("invalid components: %v", e) }
		if err := doc.Components.Validate(c); err != nil {
			return wrap(err)
		}
	}

	{
		wrap := func(e error) error { return fmt.Errorf("invalid info: %v", e) }
		if v := doc.Info; v != nil {
			if err := v.Validate(c); err != nil {
				return wrap(err)
			}
		} else {
			return wrap(errors.New("must be an object"))
		}
	}

	{
		wrap := func(e error) error { return fmt.Errorf("invalid paths: %v", e) }
		if v := doc.Paths; v != nil {
			if err := v.Validate(c); err != nil {
				return wrap(err)
			}
		} else {
			return wrap(errors.New("must be an object"))
		}
	}

	{
		wrap := func(e error) error { return fmt.Errorf("invalid security: %v", e) }
		if v := doc.Security; v != nil {
			if err := v.Validate(c); err != nil {
				return wrap(err)
			}
		}
	}

	{
		wrap := func(e error) error { return fmt.Errorf("invalid servers: %v", e) }
		if v := doc.Servers; v != nil {
			if err := v.Validate(c); err != nil {
				return wrap(err)
			}
		}
	}

	return nil
}