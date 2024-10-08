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

type SampleSubject struct {
	SampleID  string
	SubjectID string
}

func (l SampleSubject) GetID() string {
	return fmt.Sprintf("%s:%s", l.SampleID, l.SubjectID)
}

type SavedSample struct {
	PennsieveID changesetmodels.PennsieveInstanceID `json:"-"`
	Sample
}

func (ss SavedSample) GetPennsieveID() changesetmodels.PennsieveInstanceID {
	return ss.PennsieveID
}

type SampleSubjectLink link

type SampleSubjectInstance struct {
	SampleSubjectLink
	SampleSubject
}
type SavedSampleSubjectInstance struct {
	PennsieveID changesetmodels.PennsieveInstanceID `json:"-"`
	SampleSubjectInstance
}

func (sl SavedSampleSubjectInstance) GetPennsieveID() changesetmodels.PennsieveInstanceID {
	return sl.PennsieveID
}
