package rollbar

import "errors"

var errNotImpl = errors.New("Not implemented")

type noopClient struct {
	apiBaseURL string
}

func (self *noopClient) APIBaseURL() string {
	return self.apiBaseURL
}

func (self *noopClient) SetAPIBaseURL(url string) Client {
	self.apiBaseURL = url
	return self
}

func (self *noopClient) GetItem(id uint64) (*ItemResponse, error) {
	return nil, errNotImpl
}

func (self *noopClient) GetItemByCounter(counter uint64) (*ItemResponse, error) {
	return nil, errNotImpl
}

func (self *noopClient) SetItemStatus(id uint64, status string) error {
	return errNotImpl
}

func (self *noopClient) SetItemStatusByCounter(counter uint64, status string) error {
	return errNotImpl
}

func (self *noopClient) GetItemOccurrences(item_id uint64) (*OccurrencesResponse, error) {
	return nil, errNotImpl
}

func (self *noopClient) GetItemOccurrencesWithPage(item_id uint64, page uint64) (*OccurrencesResponse, error) {
	return nil, errNotImpl
}

func (self *noopClient) GetOccurrence(item_id uint64) (*OccurrenceResponse, error) {
	return nil, errNotImpl
}

func (self *noopClient) GetOccurrences() (*OccurrencesResponse, error) {
	return nil, errNotImpl
}

func (self *noopClient) GetOccurrencesWithPage(page uint64) (*OccurrencesResponse, error) {
	return nil, errNotImpl
}

func (self *noopClient) NewMessageNotification(level NotificationLevel, message string, custom CustomInfo) *MessageNotification {
	return NewMessageNotification(level, message, custom)
}

func (self *noopClient) NewTraceNotification(level NotificationLevel, message string, custom CustomInfo) *TraceNotification {
	return NewTraceNotification(level, message, custom)
}

func (self *noopClient) NewTraceChainNotification(level NotificationLevel, message string, custom CustomInfo) *TraceChainNotification {
	return NewTraceChainNotification(level, message, custom)
}

func (self *noopClient) NewCrashReportNotification(level NotificationLevel, message string, custom CustomInfo) *CrashReportNotification {
	return NewCrashReportNotification(level, message, custom)
}

func (self *noopClient) SendNotification(notif Notification) (*NotificationResponse, error) {
	res := &NotificationResponse{Err: 0}
	res.Result.UUID = "fake-uuid"
	return res, nil
}

func (self *noopClient) Options() *ClientOptions {
	return &ClientOptions{}
}
