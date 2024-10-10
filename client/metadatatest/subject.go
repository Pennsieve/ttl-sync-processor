package metadatatest

import (
	"github.com/google/uuid"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"slices"
)

type SubjectBuilder struct {
	s *metadata.Subject
}

func NewSubjectBuilder() *SubjectBuilder {
	return &SubjectBuilder{s: &metadata.Subject{ID: NewExternalInstanceID(), Species: uuid.NewString()}}
}

func (b *SubjectBuilder) WithSpeciesSynonyms(count int) *SubjectBuilder {
	for i := 0; i < count; i++ {
		b.s.SpeciesSynonyms = append(b.s.SpeciesSynonyms, uuid.NewString())
	}
	return b
}

func (b *SubjectBuilder) Build() metadata.Subject {
	return *b.s
}

func NewSavedSubject(subject metadata.Subject) metadata.SavedSubject {
	return metadata.SavedSubject{
		PennsieveID: changesetmodels.PennsieveInstanceID(uuid.NewString()),
		Subject:     subject,
	}
}

// Copy returns a deep copy.
func SubjectCopy(original metadata.Subject) metadata.Subject {
	// make use of pass by value. original is already a copy of the actual
	// original. But we need to clone the slice since that will still refer to
	// the caller's underlying array
	original.SpeciesSynonyms = slices.Clone(original.SpeciesSynonyms)
	return original
}
