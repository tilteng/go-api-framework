package serializers_mw

import (
	"context"
	"io"

	"github.com/tilteng/go-api-router/api_router"
)

type Serializer interface {
	GetMimeType() string
	Deserialize([]byte, interface{}) error
	DeserializeFromReader(io.Reader, interface{}) error
	Serialize(interface{}) ([]byte, error)
	SerializeToWriter(io.Writer, interface{}) error
}

var serializers = map[string]Serializer{
	"application/json": _jsonSerializer,
}

type RequestContext interface {
	WriteSerializedResponse(context.Context, interface{}) error
	ReadDeserializedBody(context.Context, interface{}) error
}

type requestContext struct {
	rctx         *api_router.RequestContext
	deserializer Serializer
	serializer   Serializer
}

func (self *requestContext) WriteSerializedResponse(_ context.Context, v interface{}) error {
	self.rctx.WriteStatusHeader()
	return self.serializer.SerializeToWriter(
		self.rctx.ResponseWriter(),
		v,
	)
}

func (self *requestContext) ReadDeserializedBody(_ context.Context, v interface{}) error {
	return self.deserializer.DeserializeFromReader(
		self.rctx.Body(),
		v,
	)
}
