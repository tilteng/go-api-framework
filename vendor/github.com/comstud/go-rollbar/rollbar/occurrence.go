package rollbar

import (
	"errors"
	"fmt"
	"net/url"
)

type OccurrenceMetadata struct {
	CustomerTimestamp JSONTime    `json:"customer_timestamp"`
	TimestampMS       int64       `json:"timestamp_ms"`
	APIServerHostname string      `json:"api_server_hostname"`
	Debug             interface{} `json:"debug"`
}

type OccurrenceNotifier struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type OccurrenceBody struct {
	Message map[string]interface{} `json:"message"`
}

type OccurrenceData struct {
	Environment string             `json:"environment"`
	Timestamp   JSONTime           `json:"timestamp"`
	CodeVersion string             `json:"code_version"`
	Platform    string             `json:"platform"`
	Level       string             `json:"level"`
	Notifier    OccurrenceNotifier `json:"notifier"`
	Context     string             `json:"context"`
	Title       string             `json:"title"`
	Body        OccurrenceBody     `json:"body"`
	Metadata    OccurrenceMetadata `json:"metadata"`
	Framework   string             `json:"framework"`
	UUID        string             `json:"uuid"`
}

// Occurrence object as returned from API
type Occurrence struct {
	ID         uint64         `json:"id"`
	Project_id uint64         `json:"project_id"`
	Timestamp  JSONTime       `json:"timestamp"`
	Version    uint           `json:"version"`
	Data       OccurrenceData `json:"data"`
	Billable   uint           `json:"billable"`
}

// String representation of an occurrence (pretty json)
func (self *Occurrence) String() string {
	return self.AsPrettyJSON()
}

// Occurrence as a json string
func (self *Occurrence) AsJSON() string {
	return asJSON(self)
}

// Occurrence as a pretty json string
func (self *Occurrence) AsPrettyJSON() string {
	return asPrettyJSON(self)
}

// Full API response for a single occurrence
type OccurrenceResponse struct {
	BaseAPIResponse
	*Occurrence `json:"result"`
}

// Container for multiple occurrences
type OccurrencesResult struct {
	Occurrences []*Occurrence `json:"instances"`
}

// Full API response for multiple occurrences
type OccurrencesResponse struct {
	rollbar *client
	BaseAPIResponse
	Page               uint64
	*OccurrencesResult `json:"result"`
}

// String representation of occurrences response (pretty json)
func (self *OccurrencesResponse) String() string {
	return self.AsPrettyJSON()
}

// Occurrences response as a json string
func (self *OccurrencesResponse) AsJSON() string {
	s := "["
	for i, occur := range self.Occurrences {
		if i != 0 {
			s += ","
		}
		s += asJSON(occur)
	}
	return s + "]"
}

// Occurrences response as a pretty json string
func (self *OccurrencesResponse) AsPrettyJSON() string {
	s := "["
	for i, occur := range self.Occurrences {
		if i != 0 {
			s += ","
		}
		s += asPrettyJSON(occur)
	}
	return s + "]"
}

// Does an occurrences response have more pages?
func (self *OccurrencesResponse) HasMorePages() bool {
	return self.IsSuccess() && len(self.Occurrences) > 0
}

// Get the next page of currences
func (self *OccurrencesResponse) GetNextPage() (*OccurrencesResponse, error) {
	if !self.HasMorePages() {
		return self, nil
	}
	resp := &OccurrencesResponse{
		rollbar: self.rollbar,
		Page:    self.Page + 1,
	}
	return self.rollbar.getOccurrences(resp)
}

func (self *client) getOccurrences(resp *OccurrencesResponse) (*OccurrencesResponse, error) {
	query := url.Values{
		"page": []string{fmt.Sprintf("%d", resp.Page)},
	}

	err := self.httpGet("/instances", query, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (self *client) getItemOccurrences(item_id uint64, resp *OccurrencesResponse) (*OccurrencesResponse, error) {
	query := url.Values{
		"page": []string{fmt.Sprintf("%d", resp.Page)},
	}

	err := self.httpGet(fmt.Sprintf("/item/%d/instances", item_id), query, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Get an occurrence by its id (id is NOT the same as the counter)
func (self *client) GetOccurrence(id uint64) (*OccurrenceResponse, error) {
	occur_resp := &OccurrenceResponse{}

	err := self.httpGet(
		fmt.Sprintf("/instance/%d", id),
		nil,
		&occur_resp,
	)
	if err != nil {
		return nil, err
	}

	return occur_resp, nil
}

// Get first page of all occurrences
func (self *client) GetOccurrences() (*OccurrencesResponse, error) {
	resp := &OccurrencesResponse{
		rollbar: self,
		Page:    1,
	}
	return self.getOccurrences(resp)
}

// Get a specific page of all occurrences
func (self *client) GetOccurrencesWithPage(page uint64) (*OccurrencesResponse, error) {
	if page == 0 {
		return nil, errors.New("Page must be greater than 0")
	}
	resp := &OccurrencesResponse{
		rollbar: self,
		Page:    page,
	}
	return self.getOccurrences(resp)
}

// Get first page of occurrences for an item (by item id -- NOT the counter)
func (self *client) GetItemOccurrences(item_id uint64) (*OccurrencesResponse, error) {
	resp := &OccurrencesResponse{
		rollbar: self,
		Page:    1,
	}
	return self.getItemOccurrences(item_id, resp)
}

// Get a specific page of occurrences for an item (by item id -- NOT the counter)
func (self *client) GetItemOccurrencesWithPage(item_id uint64, page uint64) (*OccurrencesResponse, error) {
	if page == 0 {
		return nil, errors.New("Page must be greater than 0")
	}
	resp := &OccurrencesResponse{
		rollbar: self,
		Page:    page,
	}
	return self.getItemOccurrences(item_id, resp)
}
