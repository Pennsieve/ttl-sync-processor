package metadata

import (
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
)

const SampleSubjectLinkName = SubjectModelName
const SampleSubjectLinkDisplayName = SubjectDisplayName

// SampleSubject is what we can get from the curation export. Only has the
// external, curation ids, no Pennsieve ids.
type SampleSubject struct {
	SampleID  changesetmodels.ExternalInstanceID
	SubjectID changesetmodels.ExternalInstanceID
}

func (l SampleSubject) ExternalID() changesetmodels.ExternalInstanceID {
	return changesetmodels.ExternalInstanceID(fmt.Sprintf("%s:%s", l.SampleID, l.SubjectID))
}

func (l SampleSubject) FromExternalID() changesetmodels.ExternalInstanceID {
	return l.SampleID
}

func (l SampleSubject) ToExternalID() changesetmodels.ExternalInstanceID {
	return l.SubjectID
}

type SavedSampleSubject struct {
	PennsieveID changesetmodels.PennsieveInstanceID `json:"-"`
	Link
	SampleSubject
}

func (sl SavedSampleSubject) GetPennsieveID() changesetmodels.PennsieveInstanceID {
	return sl.PennsieveID
}
