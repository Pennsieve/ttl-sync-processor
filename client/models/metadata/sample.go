package metadata

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

const SampleModelName = "sample"
const SampleDisplayName = "Sample"

// Keys should match the json struct tag

const SampleIDKey = "id"

type Sample struct {
	ID string `json:"id"`
}

func (s Sample) GetID() string {
	return s.ID
}

type SavedSample struct {
	PennsieveID changesetmodels.PennsieveRecordID `json:"-"`
	Sample
}

func (ss SavedSample) GetPennsieveID() changesetmodels.PennsieveRecordID {
	return ss.PennsieveID
}
