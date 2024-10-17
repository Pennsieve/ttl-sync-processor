package metadatatest

import (
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

type SampleBuilder struct {
	s *metadata.Sample
}

func NewSampleBuilder() *SampleBuilder {
	return &SampleBuilder{s: &metadata.Sample{ID: NewExternalInstanceID(), PrimaryKey: uuid.NewString()}}
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

func SampleCopy(original metadata.Sample) metadata.Sample {
	// make use of pass by value. original is already a copy of the actual
	// original.
	// If a slice is added to Sample, will have to clone it as is done in SubjectCopy
	// If a slice is added to Sample, will have to clone it as is done in SubjectCopy
	return original
}
