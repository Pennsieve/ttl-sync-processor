package models

type ModelChanges struct {
	// The ID of the model. Can be empty or missing if the model does not exist.
	// In this case, Create below should be non-nil
	ID string `json:"id,omitempty"`
	// If Create is non-nil, the model should be created
	Create *ModelCreate `json:"create,omitempty"`
	// Records describes the changes to the records of this model type
	Records RecordChanges `json:"records"`
}

type ModelCreate struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type RecordChanges struct {
	// If DeleteAll is true, delete all records for the model. Model.ID should be non-empty in this case.
	DeleteAll bool `json:"delete_all"`
	// A list of RecordIDs to delete
	Delete []string `json:"delete"`
	// Create are records that should be created
	Create []RecordCreate `json:"create"`
	// Update are records that should be updated
	Update []RecordUpdate `json:"update"`
}

type RecordCreate struct{}

type RecordUpdate struct{}
