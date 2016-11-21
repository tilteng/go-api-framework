package jsonschema_mw

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/tilteng/go-api-router/api_router"
	"github.com/xeipuuv/gojsonschema"
)

type JSONSchemaWrapper struct {
	errorHandler ErrorHandler
	linkPath     string
	schema       *gojsonschema.Schema
}

func (self *JSONSchemaWrapper) validateBody(ctx context.Context, rctx *api_router.RequestContext, body []byte) bool {
	our_result := &JSONSchemaResult{}

	loader := gojsonschema.NewStringLoader(string(body))
	resp, err := self.schema.Validate(loader)
	if err != nil {
		our_result.errors = []*JSONSchemaResultError{
			&JSONSchemaResultError{
				internalError: "Error validating body: " + err.Error(),
			},
		}
	} else if resp.Valid() {
		return true
	} else {
		json_errors := resp.Errors()
		our_result.errors = make(
			[]*JSONSchemaResultError, len(json_errors), len(json_errors),
		)
		for i, result := range json_errors {
			our_result.errors[i] = &JSONSchemaResultError{
				resultError: result,
			}
		}
	}

	if self.errorHandler != nil {
		return self.errorHandler.Error(ctx, our_result)
	}

	var str string
	for i, result := range our_result.errors {
		if i == 0 {
			str += result.String()
		} else {
			str += "," + result.String()
		}
	}
	rctx.SetStatus(400)
	rctx.WriteResponse(nil) // Force writing of status
	panic(str)
}

func (self *JSONSchemaWrapper) SetErrorHandler(error_handler ErrorHandler) *JSONSchemaWrapper {
	self.errorHandler = error_handler
	return self
}

func (self *JSONSchemaWrapper) Wrap(next api_router.RouteFn) api_router.RouteFn {
	return func(ctx context.Context) {
		rctx := api_router.RequestContextFromContext(ctx)
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
