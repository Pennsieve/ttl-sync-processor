package models

import (
	"encoding/json"
	"fmt"
)

// PennsieveRecordID is the internal ID of the record in Pennsieve. Not usually seen by user, but needed for API calls
type PennsieveRecordID string

type ModelChanges struct {
	// The ID of the model. Can be empty or missing if the model does not exist.
	// In this case, Create below should be non-nil
	ID string `json:"id,omitempty"`
	// If Create is non-nil, the model should be created
	Create *ModelPropsCreate `json:"create,omitempty"`
	// Records describes the changes to the records of this model type
	Records RecordChanges `json:"records"`
}

type ModelPropsCreate struct {
	Model      ModelCreate      `json:"model"`
	Properties PropertiesCreate `json:"properties"`
}

// ModelCreate can be used as a payload for POST /models/datasets/<dataset id>/concepts to create a model
type ModelCreate struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Locked      bool   `json:"locked"`
}

// PropertiesCreate can be uses as a payload for PUT /models/datasets/<dataset id>/concepts/<model id>/properties to add properties to a model
type PropertiesCreate []PropertyCreate

type PropertyCreate struct {
	DisplayName  string          `json:"displayName"`
	Name         string          `json:"name"`
	DataType     json.RawMessage `json:"dataType"`
	ConceptTitle bool            `json:"conceptTitle"`
	Default      bool            `json:"default"`
	Required     bool            `json:"required"`
	IsEnum       bool            `json:"isEnum"`
	IsMultiValue bool            `json:"isMultiValue"`
	Value        string          `json:"value"`
	Locked       bool            `json:"locked"`
	Description  string          `json:"description"`
}

func (pc *PropertyCreate) SetDataType(dataType any) error {
	bytes, err := json.Marshal(dataType)
	if err != nil {
		return fmt.Errorf("error marshalling data type: %w", err)
	}
	pc.DataType = bytes
	return nil
}

type RecordChanges struct {
	// If DeleteAll is true, delete all records for the model. Model.ID should be non-empty in this case.
	DeleteAll bool `json:"delete_all"`
	// A list of RecordIDs to delete
	Delete []PennsieveRecordID `json:"delete"`
	// Create are records that should be created
	Create []RecordCreate `json:"create"`
	// Update are records that should be updated
	Update []RecordUpdate `json:"update"`
}

// RecordCreate can be used as a payload for POST /models/datasets/<dataset id>/concepts/<model id>/instances
type RecordCreate RecordValues

type RecordValue struct {
	Value any    `json:"value"`
	Name  string `json:"name"`
}

type RecordValues struct {
	Values []RecordValue `json:"values"`
}

// RecordUpdate can be used as a payload for PUT /models/datasets/<dataset id>/concepts/<model id>/instances/<record id> to update values in record
// Include both changed and unchanged values
// PennsieveID is omitted from json serialization so that this struct's serialized form can be used as a payload without including the record's id
// which belongs in the URL, not payload
type RecordUpdate struct {
	PennsieveID PennsieveRecordID `json:"-"`
	RecordValues
}
