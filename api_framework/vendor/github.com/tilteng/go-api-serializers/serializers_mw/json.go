package serializers_mw

import (
	"encoding/json"
	"io"
)

var _jsonSerializer = jsonSerializer{}

type jsonSerializer struct{}

func (self jsonSerializer) GetMimeType() string {
	return "application/json"
}

func (self jsonSerializer) Serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (self jsonSerializer) SerializeToWriter(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (self jsonSerializer) Deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (self jsonSerializer) DeserializeFromReader(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
