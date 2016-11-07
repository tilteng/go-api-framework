package serializers

import (
	"io"

	controller "github.com/comstud/go-api-controller"
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

type requestContext struct {
	rctx         *controller.RequestContext
	deserializer Serializer
	serializer   Serializer
}

func (self *requestContext) WriteSerializedResponse(v interface{}) error {
	self.rctx.WriteStatusHeader()
	return self.serializer.SerializeToWriter(
		self.rctx.ResponseWriter(),
		v,
	)
}

func (self *requestContext) DeserializedBody(v interface{}) error {
	return self.deserializer.DeserializeFromReader(
		self.rctx.Body(),
		v,
	)
}
