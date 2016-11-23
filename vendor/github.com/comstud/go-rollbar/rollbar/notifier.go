package rollbar

import "time"

type NotificationLevel string

const (
	LV_CRITICAL NotificationLevel = NotificationLevel("critical")
	LV_ERROR    NotificationLevel = NotificationLevel("error")
	LV_WARNING  NotificationLevel = NotificationLevel("warning")
	LV_INFO     NotificationLevel = NotificationLevel("info")
	LV_DEBUG    NotificationLevel = NotificationLevel("debug")
)

type CustomInfo map[string]interface{}

type Notification interface {
	GetEnvironment() string
	GetLevel() NotificationLevel
}

// NotificationResponse contains the API response for posting an item
type NotificationResponse struct {
	Err    int `json:"err"`
	Result struct {
		UUID string `json:"uuid"`
	} `json:"result"`
}

func (self *client) fillBaseNotification(base *baseNotification, level NotificationLevel, title string, custom CustomInfo) {
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

func (self *client) NewMessageNotification(level NotificationLevel, message string, custom CustomInfo) *MessageNotification {
	notif := &MessageNotification{}
	self.fillBaseNotification(&notif.baseNotification, level, message, custom)
	notif.notifierMessageBody.Message.Body = message

	return notif
}

func (self *client) NewTraceNotification(level NotificationLevel, message string, custom CustomInfo) *TraceNotification {
	notif := &TraceNotification{}

	self.fillBaseNotification(&notif.baseNotification, level, message, custom)
	return notif
}

func (self *client) NewTraceChainNotification(level NotificationLevel, message string, custom CustomInfo) *TraceChainNotification {
	notif := &TraceChainNotification{}

	self.fillBaseNotification(&notif.baseNotification, level, message, custom)
	return notif
}

func (self *client) NewCrashReportNotification(level NotificationLevel, message string, custom CustomInfo) *CrashReportNotification {
	notif := &CrashReportNotification{}

	self.fillBaseNotification(&notif.baseNotification, level, message, custom)
	return notif
}

func (self *client) SendNotification(notif Notification) (*NotificationResponse, error) {
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
