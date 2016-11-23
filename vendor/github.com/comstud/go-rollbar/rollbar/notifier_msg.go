package rollbar

import "encoding/json"

type MessageNotification struct {
	baseNotification
	notifierMessageBody `json:"body"`
}

type notifierMessageBody struct {
	Message NotifierMessage `json:"message,omitempty"`
}

// Object to use in NotificationData
type NotifierMessage struct {
	// Primary message text (required)
	Body string `json:"body"`

	// Arbitrary custom data
	Custom CustomInfo `json:"-"`
}

func (self *NotifierMessage) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{
		"body": self.Body,
	}
	if self.Custom != nil {
		for k, v := range self.Custom {
			obj[k] = v
		}
	}
	return json.Marshal(obj)
}
