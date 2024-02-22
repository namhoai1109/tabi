package model

import "fmt"

type ErrorDetail struct {
	Field       string `json:"field"`
	Location    string `json:"location"`
	Issue       string `json:"issue"`
	Description string `json:"description"`
}

type ErrorLinks struct {
	Href    string `json:"href"`
	Rel     string `json:"rel"`
	EncType string `json:"encType"`
}

type ErrorResponse struct {
	Name    string        `json:"name"`
	Message string        `json:"message"`
	DebugID string        `json:"debug_id"`
	Details []ErrorDetail `json:"details"`
	Links   []ErrorLinks  `json:"links"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Name: %s \n Message: %s \n DebugID: %s \n Details: %+v", e.Name, e.Message, e.DebugID, e.Details)
}
