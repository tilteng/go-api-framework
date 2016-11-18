package rollbar

import (
	"encoding/json"
	"time"
)

type BaseAPIResponse struct {
	Err     int    `json:"err"`
	Message string `json:"message"`
}

func (self *BaseAPIResponse) IsError() bool {
	return self.Err != 0
}

func (self *BaseAPIResponse) IsSuccess() bool {
	return self.Err == 0
}

type JSONTime struct{ time.Time }

func (self JSONTime) MarshallJSON() ([]byte, error) {
	return json.Marshal(self.Unix())
}

func (self *JSONTime) UnmarshalJSON(buf []byte) error {
	var epoch int64
	err := json.Unmarshal(buf, &epoch)
	if err == nil {
		*self = JSONTime{time.Unix(epoch, 0)}
	}
	return err
}

func asJSON(v interface{}) string {
	stuff, _ := json.Marshal(v)
	return string(stuff)
}

func asPrettyJSON(v interface{}) string {
	stuff, _ := json.MarshalIndent(v, "", "  ")
	return string(stuff)
}
