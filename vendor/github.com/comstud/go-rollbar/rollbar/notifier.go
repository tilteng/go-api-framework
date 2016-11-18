package rollbar

import "time"

type notificationLevel string

const (
	LV_CRITICAL notificationLevel = notificationLevel("critical")
	LV_ERROR    notificationLevel = notificationLevel("error")
	LV_WARNING  notificationLevel = notificationLevel("warning")
	LV_INFO     notificationLevel = notificationLevel("info")
	LV_DEBUG    notificationLevel = notificationLevel("debug")
)

type CustomInfo map[string]interface{}

type Notification interface {
}

// NotificationResponse contains the API response for posting an item
type NotificationResponse struct {
	Err    int `json:"err"`
	Result struct {
		UUID string `json:"uuid"`
	} `json:"result"`
}

func (self *Client) sendNotification(notif interface{}) (*NotificationResponse, error) {
	notif_resp := &NotificationResponse{}
	err := self.httpPost(
		"/item/",
		map[string]interface{}{
			"access_token": self.accessToken,
			"data":         notif,
		},
		&notif_resp,
	)
	if err != nil {
		return nil, err
	}
	return notif_resp, nil
}

func (self *Client) fillBaseNotification(base *baseNotification, level notificationLevel, title string, custom CustomInfo) {
	base.client = self
	base.Environment = self.Environment
	base.Level = level
	base.Title = title
	base.Language = self.Language
	base.Timestamp = time.Now().Unix()
	base.Platform = self.Platform
	base.Framework = self.Framework
	base.Server = &self.NotifierServer
	base.CodeVersion = self.NotifierServer.CodeVersion
	base.Custom = custom
	if len(self.notifierName) != 0 {
		base.Notifier = &NotifierLibrary{
			Name:    self.notifierName,
			Version: self.notifierVersion,
		}
	}
}

func (self *Client) NewMessageNotification(level notificationLevel, message string, custom CustomInfo) *MessageNotification {
	notif := &MessageNotification{}
	self.fillBaseNotification(&notif.baseNotification, level, message, custom)
	notif.notifierMessageBody.Message.Body = message

	return notif
}

func (self *Client) NewTraceNotification(level notificationLevel, message string, custom CustomInfo) *TraceNotification {
	notif := &TraceNotification{}

	self.fillBaseNotification(&notif.baseNotification, level, message, custom)
	return notif
}

func (self *Client) NewTraceChainNotification(level notificationLevel, message string, custom CustomInfo) *TraceChainNotification {
	notif := &TraceChainNotification{}

	self.fillBaseNotification(&notif.baseNotification, level, message, custom)
	return notif
}

func (self *Client) NewCrashReportNotification(level notificationLevel, message string, custom CustomInfo) *CrashReportNotification {
	notif := &CrashReportNotification{}

	self.fillBaseNotification(&notif.baseNotification, level, message, custom)
	return notif
}
