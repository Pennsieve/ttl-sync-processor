package fromcuration

import (
	"github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

func ToSampleSubjectLink(exportedSample curation.Sample) (metadata.SampleSubject, error) {
	return metadata.SampleSubject{
		SampleID:  models.ExternalInstanceID(exportedSample.ID),
		SubjectID: models.ExternalInstanceID(exportedSample.SubjectID),
	}, nil
}
