package models

type Dataset struct {
	Models        []ModelChanges        `json:"models"`
	Relationships []RelationshipChanges `json:"relationships"`
}
