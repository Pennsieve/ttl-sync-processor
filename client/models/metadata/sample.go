package metadata

import (
	"fmt"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
)

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

const SampleSubjectLinkName = SubjectModelName
const SampleSubjectLinkDisplayName = SubjectDisplayName

type SampleSubjectLink struct {
	SampleID  string
	SubjectID string
	//id is not exported. just a cached value for GetID()
	id string
}

func (l SampleSubjectLink) GetID() string {
	if len(l.id) == 0 {
		l.id = fmt.Sprintf("%s:%s", l.SampleID, l.SubjectID)
	}
	return l.id
}

type SavedSample struct {
	PennsieveID changesetmodels.PennsieveInstanceID `json:"-"`
	Sample
}

func (ss SavedSample) GetPennsieveID() changesetmodels.PennsieveInstanceID {
	return ss.PennsieveID
}

type SavedSampleSubjectLink struct {
	PennsieveID changesetmodels.PennsieveInstanceID `json:"-"`
	SampleSubjectLink
}

func (sl SavedSampleSubjectLink) GetPennsieveID() changesetmodels.PennsieveInstanceID {
	return sl.PennsieveID
}
