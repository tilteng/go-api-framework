package api_framework

import (
	"context"

	"github.com/comstud/go-rollbar/rollbar"
	"github.com/tilteng/go-api-jsonschema/jsonschema_mw"
	"github.com/tilteng/go-errors/errors"
)

var ErrInternalServerError = errors.ErrInternalServerError
var ErrJSONSchemaValidationFailed = errors.ErrJSONSchemaValidationFailed
var ErrRouteNotFound = errors.ErrRouteNotFound

// Called to format an error or errors. Pass to custom callback, if set.
func (self *Controller) formatErrors(ctx context.Context, errtype errors.ErrorType) interface{} {
	rctx := self.RequestContext(ctx)
	if self.options.ErrorFormatter != nil {
		return self.options.ErrorFormatter.FormatErrors(rctx, errtype)
	}
	return errtype.AsJSONAPIResponse()
}

// Called when a panic occurs. Pass to custom callback, if set.
func (self *Controller) handlePanic(ctx context.Context, v interface{}) {
	rctx := self.RequestContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			self.logger.LogErrorf(
				rctx,
				"Received panic while processing another panic: %+v",
				r,
			)
			rctx.SetStatus(500)
			self.WriteResponse(rctx, nil)
		}
	}()

	if self.options.PanicHandler != nil {
		self.options.PanicHandler.Panic(rctx, v)
		return
	}

	err_obj, ok := v.(*errors.Error)
	if !ok {
		err_obj = ErrInternalServerError.New(rctx, "")
		err_obj.SetInternal(v)
	}

	self.WriteResponse(rctx, err_obj)
}

// Called when a serializer erorr occurs. Pass to custom callback, if set.
func (self *Controller) handleSerializerError(ctx context.Context, err error) {
	rctx := self.RequestContext(ctx)

	if self.options.SerializerErrorHandler != nil {
		self.options.SerializerErrorHandler.Error(rctx, err)
		return
	}

	rctx.SetStatus(500)
	rctx.WriteResponseString(err.Error())
}

// Called when a json schema validation failure or error occurs.
// Pass to custom callback, if set.
func (self *Controller) handleJSONSchemaError(ctx context.Context, result *jsonschema_mw.JSONSchemaResult) bool {
	rctx := self.RequestContext(ctx)

	if self.options.JSONSchemaErrorHandler != nil {
		return self.options.JSONSchemaErrorHandler.Error(rctx, result)
	}

	json_errors := result.Errors()
	api_errors := make(errors.Errors, 0, len(json_errors))
	for _, json_err := range json_errors {
		err := ErrJSONSchemaValidationFailed.New(rctx, "")
		err.Details = json_err.String()
		api_errors.AddError(err)
	}

	self.WriteResponse(rctx, api_errors)

	return false
}

func NewErrorHandler(ctx context.Context, err errors.ErrorType) {
	rctx := RequestContextFromContext(ctx)
	if rctx == nil {
		return
	}

	status := err.GetStatus()
	if status < 500 {
		return
	}

	json, json_err := err.AsJSON()
	if json_err != nil {
		rctx.LogErrorf("Returning exception: %+v", err)
	} else {
		rctx.LogError("Returning exception: " + json)
	}

	if !rctx.RollbarEnabled() {
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				rctx.LogErrorf("Received error in rollbar goroutine: %+v",
					r,
				)
			}
		}()

		var notif rollbar.Notification

		custom_info := rollbar.CustomInfo{
			"error":    err,
			"route":    rctx.CurrentRoute().FullPath(),
			"trace_id": rctx.GetTraceID(),
			"span_id":  rctx.GetSpanID(),
		}

		title := err.GetInternalError()
		if len(title) == 0 {
			title = err.GetTitle()
		}

		trace := err.GetStackTrace()
		if len(trace) != 0 {
			tnotif := rctx.RollbarClient().NewTraceNotification(
				rollbar.LV_CRITICAL,
				title,
				custom_info,
			)

			tnotif.Trace.Exception = &rollbar.NotifierException{
				Class:       err.GetName(),
				Message:     err.GetTitle(),
				Description: err.GetDetails(),
			}

			frames := make([]*rollbar.NotifierFrame, len(trace), len(trace))
			for i, frame := range trace {
				frames[i] = &rollbar.NotifierFrame{
					Filename: frame.Filename,
					Method:   frame.Function,
					Line:     frame.LineNo,
				}
			}

			tnotif.Trace.Frames = frames
			notif = tnotif
		} else {
			notif = rctx.RollbarClient().NewMessageNotification(
				rollbar.LV_ERROR,
				title,
				custom_info,
			)
		}

		http_req := rctx.HTTPRequest()
		notif.SetRequest(&rollbar.NotifierRequest{
			URL:    http_req.URL.String(),
			Method: http_req.Method,
		})
		rctx.RollbarClient().SendNotification(notif)
	}()
}
