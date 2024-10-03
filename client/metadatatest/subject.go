package metadatatest

import (
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

type SubjectBuilder struct {
	s *metadata.Subject
}

func NewSubjectBuilder() *SubjectBuilder {
	return &SubjectBuilder{s: &metadata.Subject{ID: uuid.NewString(), Species: uuid.NewString()}}
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
