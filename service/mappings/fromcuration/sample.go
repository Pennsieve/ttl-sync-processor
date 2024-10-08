package fromcuration

import (
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

// ToSample is a mappings.Mapping from curation.Sample to metadata.Sample
func ToSample(exportedSample curation.Sample) (metadata.Sample, error) {
	return metadata.Sample{ID: exportedSample.ID}, nil
}

func ToSampleSubjectLink(exportedSample curation.Sample) (metadata.SampleSubject, error) {
	return metadata.SampleSubject{
		SampleID:  exportedSample.ID,
		SubjectID: exportedSample.SubjectID,
	}, nil
}
