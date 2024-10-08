package fromcuration

import (
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings"
)

func ToDatasetMetadata(export *curation.DatasetExport) (*metadata.DatasetMetadata, error) {
	exported := &metadata.DatasetMetadata{}
	var err error
	exported.Contributors, err = mappings.MapSlice(export.Contributors, ToContributor)
	if err != nil {
		return nil, err
	}
	exported.Subjects, err = mappings.MapSlice(export.Subjects, ToSubject)
	if err != nil {
		return nil, err
	}
	exported.Samples, err = mappings.MapSlice(export.Samples, ToSample)
	if err != nil {
		return nil, err
	}
	if exported.SampleSubjects, err = mappings.MapSlice(export.Samples, ToSampleSubjectLink); err != nil {
		return nil, err
	}

	return exported, nil
}
