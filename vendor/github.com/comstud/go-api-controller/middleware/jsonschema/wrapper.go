package jsonschema

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/comstud/go-api-controller"
	"github.com/xeipuuv/gojsonschema"
)

type JSONSchemaWrapper struct {
	errorHandler ErrorHandler
	linkPath     string
	schema       *gojsonschema.Schema
}

func (self *JSONSchemaWrapper) validateBody(ctx context.Context, rctx *controller.RequestContext, body []byte) bool {
	loader := gojsonschema.NewStringLoader(string(body))
	result, err := self.schema.Validate(loader)
	if err != nil {
		panic(fmt.Sprintf("Error validating body: %s", err.Error()))
	}
	if !result.Valid() {
		if self.errorHandler != nil {
			return self.errorHandler.Error(ctx, result)
		}
		var str string
		for _, desc := range result.Errors() {
			if str != "" {
				str += ","
			}
			str += fmt.Sprintf("%s", desc)
		}
		rctx.SetStatus(400)
		rctx.WriteResponse(nil) // Force writing of status
		panic(str)
	}
	return true
}

func (self *JSONSchemaWrapper) SetErrorHandler(error_handler ErrorHandler) *JSONSchemaWrapper {
	self.errorHandler = error_handler
	return self
}

func (self *JSONSchemaWrapper) Wrap(next controller.ControllerFn) controller.ControllerFn {
	return func(ctx context.Context) {
		rctx := controller.RequestContextFromContext(ctx)
		if self.linkPath != "" {
			rctx.SetResponseHeader(
				"Link",
				fmt.Sprintf(`<%s>; rel="describedBy"`,
					self.linkPath,
				),
			)
		}
		body := rctx.Body()
		buf, err := ioutil.ReadAll(body)
		if err != nil {
			panic(fmt.Sprintf("Error reading body: %s", err))
		}
		defer body.Close()
		rctx.SetBody(ioutil.NopCloser(bytes.NewBuffer(buf)))
		if self.validateBody(ctx, rctx, buf) {
			next(ctx)
		}
	}
}
