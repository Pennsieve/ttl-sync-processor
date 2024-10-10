package metadatatest

import (
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

type SampleBuilder struct {
	s *metadata.Sample
}

func NewSampleBuilder() *SampleBuilder {
	return &SampleBuilder{s: &metadata.Sample{ID: NewExternalInstanceID()}}
}

func (b *SampleBuilder) Build() metadata.Sample {
	return *b.s
}

func NewSavedSample(sample metadata.Sample) metadata.SavedSample {
	return metadata.SavedSample{
		PennsieveID: NewPennsieveInstanceID(),
		Sample:      sample,
	}
}
