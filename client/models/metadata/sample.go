package metadata

import (
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
)

const SampleModelName = "sample"
const SampleDisplayName = "Sample"

// Keys should match the json struct tag

const SampleIDKey = "id"
const PrimaryKeyKey = "primary_key"

type Sample struct {
	ID         changesetmodels.ExternalInstanceID `json:"id"`
	PrimaryKey string                             `json:"primary_key"`
}

func (s Sample) ExternalID() changesetmodels.ExternalInstanceID {
	return changesetmodels.ExternalInstanceID(s.ID)
}

type SavedSample struct {
	PennsieveID changesetmodels.PennsieveInstanceID `json:"-"`
	Sample
}

func (ss SavedSample) GetPennsieveID() changesetmodels.PennsieveInstanceID {
	return ss.PennsieveID
}
