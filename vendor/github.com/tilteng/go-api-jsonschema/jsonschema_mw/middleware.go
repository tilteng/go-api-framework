package jsonschema_mw

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tilteng/go-logger/logger"
	"github.com/xeipuuv/gojsonschema"
)

type JSONSchemaResult struct {
	errors []*JSONSchemaResultError
}

func (self *JSONSchemaResult) Errors() []*JSONSchemaResultError {
	return self.errors
}

type JSONSchemaResultError struct {
	internalError string
	resultError   gojsonschema.ResultError
}

func (self *JSONSchemaResultError) String() string {
	if self.internalError != "" {
		return self.internalError
	}
	return self.resultError.String()
}

type ErrorHandler func(context.Context, *JSONSchemaResult) bool

func (self ErrorHandler) Error(ctx context.Context, result *JSONSchemaResult) bool {
	return self(ctx, result)
}

type JSONSchema struct {
	schema     *gojsonschema.Schema
	jsonString string
}

func (self *JSONSchema) GetSchema() *gojsonschema.Schema {
	return self.schema
}

func (self *JSONSchema) GetJSONString() string {
	return self.jsonString
}

type JSONSchemaMiddleware struct {
	jsonSchemas    map[string]*JSONSchema
	logger         logger.CtxLogger
	errorHandler   ErrorHandler
	linkPathPrefix string
}

func (self *JSONSchemaMiddleware) GetSchema(name string) *JSONSchema {
	schema, _ := self.jsonSchemas[name]
	return schema
}

func (self *JSONSchemaMiddleware) GetSchemas() map[string]*JSONSchema {
	return self.jsonSchemas
}

func (self *JSONSchemaMiddleware) LoadFromPath(ctx context.Context, base_path string) error {
	return filepath.Walk(base_path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Error reading schema from %s: %s", path, err)
		}

		json_string := string(bytes)

		loader := gojsonschema.NewStringLoader(json_string)
		schema, err := gojsonschema.NewSchema(loader)
		if err != nil {
			return fmt.Errorf("Error loading schema from %s: %s", path, err)
		}

		name := filepath.Base(path)
		name = name[0 : len(name)-5]

		self.jsonSchemas[name] = &JSONSchema{
			schema:     schema,
			jsonString: json_string,
		}

		if self.logger != nil {
			self.logger.LogDebug(ctx, "Loaded schema "+name)
		}

		return nil
	})
}

func (self *JSONSchemaMiddleware) NewWrapper(schema *gojsonschema.Schema, linkpath string) *JSONSchemaWrapper {
	if linkpath != "" {
		if self.linkPathPrefix != "" {
			linkpath = self.linkPathPrefix + "/" + linkpath
		}
	}
	return &JSONSchemaWrapper{
		errorHandler: self.errorHandler,
		schema:       schema,
		linkPath:     linkpath,
	}
}

func (self *JSONSchemaMiddleware) NewWrapperFromSchemaName(ctx context.Context, name string) *JSONSchemaWrapper {
	schema := self.GetSchema(name)
	if schema == nil {
		panic(fmt.Errorf("Couldn't find json schema with name '%s'", name))
	}
	return self.NewWrapper(schema.GetSchema(), name)
}

func (self *JSONSchemaMiddleware) NewWrapperFromRouteOptions(ctx context.Context, opts ...interface{}) *JSONSchemaWrapper {
	var schema_name string
	for _, opt_map_i := range opts {
		opt_map, ok := opt_map_i.(map[string]string)
		if !ok {
			continue
		}
		schema_name, ok = opt_map["jsonschema"]
		if !ok {
			continue
		}
		return self.NewWrapperFromSchemaName(ctx, schema_name)
	}
	return nil
}

func (self *JSONSchemaMiddleware) SetLogger(logger logger.CtxLogger) *JSONSchemaMiddleware {
	self.logger = logger
	return self
}

func NewMiddleware(error_handler ErrorHandler) *JSONSchemaMiddleware {
	return &JSONSchemaMiddleware{
		errorHandler: error_handler,
		jsonSchemas:  map[string]*JSONSchema{},
	}
}

func NewMiddlewareWithLinkPathPrefix(error_handler ErrorHandler, link_path_prefix string) *JSONSchemaMiddleware {
	return &JSONSchemaMiddleware{
		errorHandler:   error_handler,
		jsonSchemas:    map[string]*JSONSchema{},
		linkPathPrefix: link_path_prefix,
	}
}
