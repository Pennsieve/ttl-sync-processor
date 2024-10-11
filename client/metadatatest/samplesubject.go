package metadatatest

import (
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

type SampleSubjectBuilder struct {
	sampleID  *changesetmodels.ExternalInstanceID
	subjectID *changesetmodels.ExternalInstanceID
}

func NewSampleSubjectBuilder() *SampleSubjectBuilder {
	return &SampleSubjectBuilder{}
}

func (b *SampleSubjectBuilder) WithSampleID(sampleID changesetmodels.ExternalInstanceID) *SampleSubjectBuilder {
	b.sampleID = &sampleID
	return b
}
func (b *SampleSubjectBuilder) WithSubjectID(subjectID changesetmodels.ExternalInstanceID) *SampleSubjectBuilder {
	b.subjectID = &subjectID
	return b
}

func (b *SampleSubjectBuilder) Build() metadata.SampleSubject {
	var sampleID, subjectID changesetmodels.ExternalInstanceID
	if b.sampleID != nil {
		sampleID = *b.sampleID
	} else {
		sampleID = NewExternalInstanceID()
	}
	if b.subjectID != nil {
		subjectID = *b.subjectID
	} else {
		subjectID = NewExternalInstanceID()
	}
	link := metadata.SampleSubject{
		SampleID:  sampleID,
		SubjectID: subjectID,
	}
	return link
}

func NewSampleSubject() metadata.SampleSubject {
	return metadata.SampleSubject{
		SampleID:  NewExternalInstanceID(),
		SubjectID: NewExternalInstanceID(),
	}
}

func NewSavedSampleSubject(sampleSubject metadata.SampleSubject) metadata.SavedSampleSubject {
	return metadata.SavedSampleSubject{
		PennsieveID: NewPennsieveInstanceID(),
		Link: metadata.Link{
			From: NewPennsieveInstanceID(),
			To:   NewPennsieveInstanceID(),
		},
		SampleSubject: sampleSubject,
	}
}
