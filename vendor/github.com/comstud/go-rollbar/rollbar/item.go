package rollbar

import (
	"errors"
	"fmt"
)

// Item struct
type Item struct {
	ID                      uint64   `json:"id"`
	Project_id              uint64   `json:"project_id"`
	Counter                 uint64   `json:"counter"`
	Environment             string   `json:"environment"`
	Platform                string   `json:"platform"`
	Framework               string   `json:"framework"`
	Hash                    string   `json:"hash"`
	Title                   string   `json:"title"`
	FirstOccurrenceId       uint64   `json:"first_occurrence_id"`
	FirstOccurenceTimestamp JSONTime `json:"first_occurrence_timestamp"`
	ActivatingOccurrenceId  uint64   `json:"activating_occurrence_id"`
	LastActivatedTimestamp  JSONTime `json:"last_activated_timestamp"`
	LastResolvedTimestamp   JSONTime `json:"last_resolved_timestamp"`
	LastMutedTimestamp      JSONTime `json:"last_muted_timestamp"`
	LastOccurrenceId        uint64   `json:"last_occurrence_id"`
	LastOccurenceTimestamp  JSONTime `json:"last_occurrence_timestamp"`
	TotalOccurrences        uint64   `json:"total_occurrences"`
	LastModifiedBy          uint64   `json:"last_modified_by"`
	Status                  string   `json:"status"`
	Level                   string   `json:"level"`
	// No idea what this is yet
	IntegrationsData interface{} `json:"integrations_data"`
}

// String representation in pretty JSON form
func (self *Item) String() string {
	return self.AsPrettyJSON()
}

// Item as json string
func (self *Item) AsJSON() string {
	return asJSON(self)
}

// Item as pretty json
func (self *Item) AsPrettyJSON() string {
	return asPrettyJSON(self)
}

// Container holding full response info for an item
type ItemResponse struct {
	BaseAPIResponse
	*Item `json:"result"`
}

// Get a single item by its id (id is NOT the same as the counter)
func (self *client) GetItem(id uint64) (*ItemResponse, error) {
	item_resp := &ItemResponse{}

	err := self.httpGet(
		fmt.Sprintf("/item/%d", id),
		nil,
		&item_resp,
	)

	if err != nil {
		return nil, err
	}

	return item_resp, nil
}

// Get a single item by its counter
func (self *client) GetItemByCounter(counter uint64) (*ItemResponse, error) {
	item_resp := &ItemResponse{}

	err := self.httpGet(
		fmt.Sprintf("/item_by_counter/%d", counter),
		nil,
		&item_resp,
	)

	if err != nil {
		return nil, err
	}

	return item_resp, nil
}

// Update an item's status by its id
func (self *client) SetItemStatus(id uint64, status string) error {
	item_update := map[string]interface{}{
		"status": status,
	}

	update_resp := &BaseAPIResponse{}

	err := self.httpPatch(
		fmt.Sprintf("/item/%d", id),
		&item_update,
		&update_resp,
	)

	if err != nil {
		return err
	}

	if update_resp.Err != 0 {
		return errors.New(update_resp.Message)
	}

	return nil
}

// Update an item's status by its counter
func (self *client) SetItemStatusByCounter(counter uint64, status string) error {
	item_resp, err := self.GetItemByCounter(counter)
	if err != nil || item_resp.Err != 0 {
		if err == nil {
			err = errors.New(item_resp.Message)
		}
		return fmt.Errorf("Error getting item ID: %s", err.Error())
	}

	return self.SetItemStatus(item_resp.ID, status)
}
