package curationtest

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
)

type SubjectBuilder struct {
	s    *curation.Subject
	errs []error
}

func NewSubjectBuilder() *SubjectBuilder {
	return &SubjectBuilder{s: &curation.Subject{ID: uuid.NewString()}}
}

func (b *SubjectBuilder) WithSimpleSpecies() *SubjectBuilder {
	if speciesBytes, err := json.Marshal(uuid.NewString()); err != nil {
		b.errs = append(b.errs, err)
	} else {
		b.s.Species = speciesBytes
	}
	return b
}

func (b *SubjectBuilder) WithSubjectIdentifier(synonymCount int) *SubjectBuilder {
	speciesIdentifier := NewIdentifierBuilder().WithSynonyms(synonymCount).Build()
	if speciesBytes, err := json.Marshal(speciesIdentifier); err != nil {
		b.errs = append(b.errs, err)
	} else {
		b.s.Species = speciesBytes
	}
	return b
}

func (b *SubjectBuilder) Build() (curation.Subject, error) {
	return *b.s, errors.Join(b.errs...)
}
