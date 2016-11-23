package rollbar

import "encoding/json"

// base notifcation object to post as 'data'. This base
type baseNotification struct {
	// Required

	Environment string `json:"environment"`

	// Optional

	Level       NotificationLevel `json:"level,omitempty"`
	Timestamp   int64             `json:"timestamp,omitempty"`
	CodeVersion string            `json:"code_verison,omitempty"`
	Platform    string            `json:"platform,omitempty"`
	Language    string            `json:"language,omitempty"`
	Framework   string            `json:"framework,omitempty"`
	Context     string            `json:"context,omitempty"`
	Request     *NotifierRequest  `json:"request,omitempty"`
	Person      *NotifierPerson   `json:"person,omitempty"`
	Server      *NotifierServer   `json:"server,omitempty"`
	Client      *NotifierClient   `json:"client,omitempty"`

	// Custom data
	Custom CustomInfo `json:"custom,omitempty"`

	Fingerprint string `json:"fingerprint,omitempty"`

	// Title -- max 255 characters
	Title string `json:"title,omitempty"`

	// Up to 36 characters.
	UUID string `json:"uuid,omitempty"`

	// Optional info that describes the library used to send event
	Notifier *NotifierLibrary `json:"notifier,omitempty"`
}

func (self *baseNotification) GetEnvironment() string {
	return self.Environment
}

func (self *baseNotification) GetLevel() NotificationLevel {
	return self.Level
}

// Optional data about client making the request
type NotifierClient struct {
	Javascript *NotifierJavascriptClient `json:"javascript,omitempty"`

	// Additional arbitrary data
	Custom CustomInfo `json:"-"`
}

func (self *NotifierClient) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	if self.Javascript != nil {
		obj["javascript"] = self.Javascript
	}
	if self.Custom != nil {
		for k, v := range self.Custom {
			obj[k] = v
		}
	}
	return json.Marshal(obj)
}

// Javascript info to use in NotifierClient
type NotifierJavascriptClient struct {
	Browser             string `json:"browser,omitempty"`
	CodeVersion         string `json:"code_version,omitempty"`
	SourceMapEnabled    *bool  `json:"source_map_enabled,omitempty"`
	GuessUncaughtFrames *bool  `json:"guess_uncaught_frames,omitempty"`

	// Additional arbitrary data
	Custom CustomInfo `json:"-"`
}

func (self *NotifierJavascriptClient) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	if len(self.Browser) != 0 {
		obj["browser"] = self.Browser
	}
	if len(self.CodeVersion) != 0 {
		obj["code_version"] = self.CodeVersion
	}
	if self.SourceMapEnabled != nil {
		obj["source_map_enabled"] = *self.SourceMapEnabled
	}
	if self.GuessUncaughtFrames != nil {
		obj["guess_uncaught_frames"] = *self.GuessUncaughtFrames
	}
	if self.Custom != nil {
		for k, v := range self.Custom {
			obj[k] = v
		}
	}
	return json.Marshal(obj)
}

// Optional user affected by event. Rollbar indexes by ID, username, and
// email. ID is unique. Most recent username,email used for an ID will
// replace older data for the ID.
type NotifierPerson struct {
	// Required. Unique ID to represent user in your system
	ID       string `json:"id"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`

	// Additional data
	Custom CustomInfo `json:"-"`
}

func (self *NotifierPerson) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	if len(self.ID) != 0 {
		obj["id"] = self.ID
	}
	if len(self.Username) != 0 {
		obj["username"] = self.Username
	}
	if len(self.Email) != 0 {
		obj["email"] = self.Email
	}
	if self.Custom != nil {
		for k, v := range self.Custom {
			obj[k] = v
		}
	}
	return json.Marshal(obj)
}

// Optional data about the request event occurred in. Can be any arbitrary
// key/value. Methods on NotifierRequest() exist for keys that rollbar
// understands.
type NotifierRequest struct {
	// Full URL where event occurred
	URL string `json:"url,omitempty"`

	// Request method.. e.g, "POST"
	Method string `json:"method,omitempty"`

	// Headers -- formatted like HTTP
	Headers map[string]string `json:"headers,omitempty"`

	// Any routing parameters (i.e. for use with Rails Routes)
	Params map[string]string `json:"params,omitempty"`

	// GET query string params
	GETParams map[string]string `json:"GET,omitempty"`

	// Raw query string
	QueryString string `json:"query_string,omitempty"`

	// POST body params
	POSTParams map[string]interface{} `json:"POST,omitempty"`

	// Raw POST body
	Body string `json:"body,omitempty"`

	// User's IP address as string
	UserIP string `json:"user_ip,omitempty"`

	// Additional data
	Custom CustomInfo `json:"-"`
}

func (self *NotifierRequest) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	if len(self.URL) != 0 {
		obj["url"] = self.URL
	}
	if len(self.Method) != 0 {
		obj["method"] = self.Method
	}
	if len(self.Headers) != 0 {
		obj["headers"] = self.Headers
	}
	if len(self.Params) != 0 {
		obj["params"] = self.Params
	}
	if len(self.GETParams) != 0 {
		obj["GET"] = self.GETParams
	}
	if len(self.QueryString) != 0 {
		obj["query_string"] = self.QueryString
	}
	if len(self.POSTParams) != 0 {
		obj["POST"] = self.POSTParams
	}
	if len(self.Body) != 0 {
		obj["body"] = self.Body
	}
	if len(self.UserIP) != 0 {
		obj["user_ip"] = self.UserIP
	}
	if self.Custom != nil {
		for k, v := range self.Custom {
			obj[k] = v
		}
	}
	return json.Marshal(obj)
}

// Optional data about the server
type NotifierServer struct {
	// Server hostname (will be indexed)
	Host string `json:"host"`

	// Path to application code root, not including final slash
	Root string `json:"root,omitempty"`

	// Code branch name
	Branch string `json:"branch,omitempty"`

	// Optional code version
	CodeVersion string `json:"code_version,omitempty"`

	// Additional arbitrary data
	Custom CustomInfo `json:"-"`
}

func (self *NotifierServer) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	if len(self.Host) != 0 {
		obj["host"] = self.Host
	}
	if len(self.Root) != 0 {
		obj["root"] = self.Root
	}
	if len(self.Branch) != 0 {
		obj["branch"] = self.Branch
	}
	if len(self.CodeVersion) != 0 {
		obj["code_version"] = self.CodeVersion
	}
	if self.Custom != nil {
		for k, v := range self.Custom {
			obj[k] = v
		}
	}
	return json.Marshal(obj)
}

// Optional data about the notifier library
type NotifierLibrary struct {
	// Optional name describing notifier (this) library
	Name string `json:"name,omitempty"`
	// Optional version of notifier library
	Version string `json:"version,omitempty"`
}
