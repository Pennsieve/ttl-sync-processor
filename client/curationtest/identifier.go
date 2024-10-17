package curationtest

import (
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
)

type IdentifierBuilder struct {
	i *curation.EmbeddedIdentifier
}

func NewIdentifierBuilder() *IdentifierBuilder {
	return &IdentifierBuilder{i: &curation.EmbeddedIdentifier{
		ID:     uuid.NewString(),
		Label:  uuid.NewString(),
		Type:   "identifier",
		System: uuid.NewString(),
	}}
}

func (b *IdentifierBuilder) WithSynonyms(count int) *IdentifierBuilder {
	for i := 0; i < count; i++ {
		b.i.Synonyms = append(b.i.Synonyms, uuid.NewString())
	}
	return b
}

func (b *IdentifierBuilder) WithDescription() *IdentifierBuilder {
	b.i.Description = uuid.NewString()
	return b
}

func (b *IdentifierBuilder) Build() curation.EmbeddedIdentifier {
	return *b.i
}
