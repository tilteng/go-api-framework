package app_context

import (
	"errors"

	"github.com/comstud/go-rollbar/rollbar"
)

var errNotImpl = errors.New("Not implemented")

type NOOPRollbarClient struct {
	apiBaseURL string
}

func (self *NOOPRollbarClient) APIBaseURL() string {
	return self.apiBaseURL
}

func (self *NOOPRollbarClient) SetAPIBaseURL(url string) rollbar.Client {
	self.apiBaseURL = url
	return self
}

func (self *NOOPRollbarClient) GetItem(id uint64) (*rollbar.ItemResponse, error) {
	return nil, errNotImpl
}

func (self *NOOPRollbarClient) GetItemByCounter(counter uint64) (*rollbar.ItemResponse, error) {
	return nil, errNotImpl
}

func (self *NOOPRollbarClient) SetItemStatus(id uint64, status string) error {
	return errNotImpl
}

func (self *NOOPRollbarClient) SetItemStatusByCounter(counter uint64, status string) error {
	return errNotImpl
}

func (self *NOOPRollbarClient) GetItemOccurrences(item_id uint64) (*rollbar.OccurrencesResponse, error) {
	return nil, errNotImpl
}

func (self *NOOPRollbarClient) GetItemOccurrencesWithPage(item_id uint64, page uint64) (*rollbar.OccurrencesResponse, error) {
	return nil, errNotImpl
}

func (self *NOOPRollbarClient) GetOccurrence(item_id uint64) (*rollbar.OccurrenceResponse, error) {
	return nil, errNotImpl
}

func (self *NOOPRollbarClient) GetOccurrences() (*rollbar.OccurrencesResponse, error) {
	return nil, errNotImpl
}

func (self *NOOPRollbarClient) GetOccurrencesWithPage(page uint64) (*rollbar.OccurrencesResponse, error) {
	return nil, errNotImpl
}

func (self *NOOPRollbarClient) NewMessageNotification(level rollbar.NotificationLevel, message string, custom rollbar.CustomInfo) *rollbar.MessageNotification {
	return &rollbar.MessageNotification{}
}

func (self *NOOPRollbarClient) NewTraceNotification(level rollbar.NotificationLevel, message string, custom rollbar.CustomInfo) *rollbar.TraceNotification {
	return &rollbar.TraceNotification{}
}

func (self *NOOPRollbarClient) NewTraceChainNotification(level rollbar.NotificationLevel, message string, custom rollbar.CustomInfo) *rollbar.TraceChainNotification {
	return &rollbar.TraceChainNotification{}
}

func (self *NOOPRollbarClient) NewCrashReportNotification(level rollbar.NotificationLevel, message string, custom rollbar.CustomInfo) *rollbar.CrashReportNotification {
	return &rollbar.CrashReportNotification{}
}

func (self *NOOPRollbarClient) SendNotification(notif rollbar.Notification) (*rollbar.NotificationResponse, error) {
	res := &rollbar.NotificationResponse{Err: 0}
	res.Result.UUID = "fake-uuid"
	return res, nil
}

func (self *NOOPRollbarClient) Options() *rollbar.ClientOptions {
	return &rollbar.ClientOptions{}
}
