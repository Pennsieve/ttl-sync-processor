package metadatatest

import (
	"github.com/google/uuid"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

type SampleBuilder struct {
	s *metadata.Sample
}

func NewSampleBuilder() *SampleBuilder {
	return &SampleBuilder{s: &metadata.Sample{ID: uuid.NewString()}}
}

func (b *SampleBuilder) Build() metadata.Sample {
	return *b.s
}

func NewSavedSample(sample metadata.Sample) metadata.SavedSample {
	return metadata.SavedSample{
		PennsieveID: changesetmodels.PennsieveInstanceID(uuid.NewString()),
		Sample:      sample,
	}
}
