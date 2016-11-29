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
	SetEnvironment(env string)
	GetLevel() NotificationLevel
	SetLevel(level NotificationLevel)
	GetTimestamp() time.Time
	SetTimestamp(t time.Time) Notification
	GetCodeVersion() string
	SetCodeVersion(vers string) Notification
	GetPlatform() string
	SetPlatform(platform string) Notification
	GetLanguage() string
	SetLanguage(lang string) Notification
	GetFramework() string
	SetFramework(framework string) Notification
	GetContext() string
	SetContext(context string) Notification
	GetRequest() *NotifierRequest
	SetRequest(req *NotifierRequest) Notification
	GetPerson() *NotifierPerson
	SetPerson(person *NotifierPerson) Notification
	GetServer() *NotifierServer
	SetServer(server *NotifierServer) Notification
	GetClient() *NotifierClient
	SetClient(client *NotifierClient) Notification
	GetCustom() CustomInfo
	SetCustom(custom CustomInfo) Notification
	GetFingerprint() string
	SetFingerprint(fingerprint string) Notification
	GetTitle() string
	SetTitle(title string) Notification
	GetUUID() string
	SetUUID(uuid string) Notification
	GetNotifier() *NotifierLibrary
	SetNotifier(notifier *NotifierLibrary) Notification
}

// NotificationResponse contains the API response for posting an item
type NotificationResponse struct {
	Err    int `json:"err"`
	Result struct {
		UUID string `json:"uuid"`
	} `json:"result"`
}

func (self *client) fillBaseNotification(base *baseNotification) {
	base.Environment = self.Environment
	base.Language = self.Language
	base.Platform = self.Platform
	base.Framework = self.Framework
	base.Server = &self.NotifierServer
	base.CodeVersion = self.NotifierServer.CodeVersion
	if len(self.notifierName) != 0 {
		base.Notifier = &NotifierLibrary{
			Name:    self.notifierName,
			Version: self.notifierVersion,
		}
	}
}

func (self *client) NewMessageNotification(level NotificationLevel, message string, custom CustomInfo) *MessageNotification {
	notif := NewMessageNotification(level, message, custom)
	self.fillBaseNotification(&notif.baseNotification)
	notif.notifierMessageBody.Message.Body = message
	return notif
}

func (self *client) NewTraceNotification(level NotificationLevel, message string, custom CustomInfo) *TraceNotification {
	notif := NewTraceNotification(level, message, custom)
	self.fillBaseNotification(&notif.baseNotification)
	return notif
}

func (self *client) NewTraceChainNotification(level NotificationLevel, message string, custom CustomInfo) *TraceChainNotification {
	notif := NewTraceChainNotification(level, message, custom)
	self.fillBaseNotification(&notif.baseNotification)
	return notif
}

func (self *client) NewCrashReportNotification(level NotificationLevel, message string, custom CustomInfo) *CrashReportNotification {
	notif := NewCrashReportNotification(level, message, custom)
	self.fillBaseNotification(&notif.baseNotification)
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

func NewMessageNotification(level NotificationLevel, title string, custom CustomInfo) *MessageNotification {
	notif := &MessageNotification{}
	notif.self = notif
	notif.Timestamp = time.Now().Unix()
	notif.Level = level
	notif.Title = title
	notif.Custom = custom
	return notif
}

func NewTraceNotification(level NotificationLevel, title string, custom CustomInfo) *TraceNotification {
	notif := &TraceNotification{}
	notif.self = notif
	notif.Timestamp = time.Now().Unix()
	notif.Level = level
	notif.Title = title
	notif.Custom = custom
	return notif
}

func NewTraceChainNotification(level NotificationLevel, title string, custom CustomInfo) *TraceChainNotification {
	notif := &TraceChainNotification{}
	notif.self = notif
	notif.Timestamp = time.Now().Unix()
	notif.Level = level
	notif.Title = title
	notif.Custom = custom
	return notif
}

func NewCrashReportNotification(level NotificationLevel, title string, custom CustomInfo) *CrashReportNotification {
	notif := &CrashReportNotification{}
	notif.self = notif
	notif.Timestamp = time.Now().Unix()
	notif.Level = level
	notif.Title = title
	notif.Custom = custom
	return notif
}
