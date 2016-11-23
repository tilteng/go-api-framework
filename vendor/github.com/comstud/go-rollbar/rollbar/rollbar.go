package rollbar

import (
	"log"
	"net/http"
	"os"
	"runtime"
)

const (
	DEFAULT_API_BASE_URL         = "https://api.rollbar.com/api/1"
	DEFAULT_NOTIFIER_NAME        = "go-rollbar"
	DEFAULT_NOTIFIER_VERSION     = "0.0.1"
	DEFAULT_NOTIFIER_ENVIRONMENT = "development"
)

// Options that can be directory modified on the client
type ClientOptions struct {
	// Logger to use
	Logger *log.Logger

	// The following affect the notifier

	// Your environment. Defaults to "development"
	Environment string

	NotifierServer NotifierServer
	Platform       string
	Language       string
	Framework      string
}

// Client interface
type Client interface {
	APIBaseURL() string
	SetAPIBaseURL(base_url string) Client
	Options() *ClientOptions
	GetItem(id uint64) (*ItemResponse, error)
	GetItemByCounter(counter uint64) (*ItemResponse, error)
	SetItemStatus(id uint64, status string) error
	SetItemStatusByCounter(counter uint64, status string) error
	GetItemOccurrences(item_id uint64) (*OccurrencesResponse, error)
	GetItemOccurrencesWithPage(item_id uint64, page uint64) (*OccurrencesResponse, error)
	GetOccurrence(id uint64) (*OccurrenceResponse, error)
	GetOccurrences() (*OccurrencesResponse, error)
	GetOccurrencesWithPage(page uint64) (*OccurrencesResponse, error)
	NewMessageNotification(level NotificationLevel, message string, custom CustomInfo) *MessageNotification
	NewTraceNotification(level NotificationLevel, message string, custom CustomInfo) *TraceNotification
	NewTraceChainNotification(level NotificationLevel, message string, custom CustomInfo) *TraceChainNotification
	NewCrashReportNotification(level NotificationLevel, message string, custom CustomInfo) *CrashReportNotification
	SendNotification(notif Notification) (*NotificationResponse, error)
}

var DefaultClientOptions ClientOptions

func init() {
	hostname, _ := os.Hostname()

	DefaultClientOptions = ClientOptions{
		Environment: DEFAULT_NOTIFIER_ENVIRONMENT,
		NotifierServer: NotifierServer{
			Host: hostname,
		},
		Platform: runtime.GOOS,
		Language: "go",
		Logger:   log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds),
	}
}

// Create a new client with specified access token
func NewClient(access_token string) (Client, error) {
	return &client{
		httpClient:      &http.Client{},
		apiBaseURL:      DEFAULT_API_BASE_URL,
		accessToken:     access_token,
		notifierName:    DEFAULT_NOTIFIER_NAME,
		notifierVersion: DEFAULT_NOTIFIER_VERSION,
		ClientOptions:   DefaultClientOptions,
	}, nil
}

func NewNOOPClient() Client {
	return &noopClient{}
}
