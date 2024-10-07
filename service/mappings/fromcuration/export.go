package fromcuration

import (
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

func ToDatasetMetadata(export *curation.DatasetExport) (*metadata.DatasetMetadata, error) {
	exported := &metadata.DatasetMetadata{}
	var err error
	exported.Contributors, err = MapSlice(export.Contributors, ToContributor)
	if err != nil {
		return nil, err
	}
	exported.Subjects, err = MapSlice(export.Subjects, ToSubject)
	if err != nil {
		return nil, err
	}
	exported.Samples, err = MapSlice(export.Samples, ToSample)
	if err != nil {
		return nil, err
	}
	return exported, nil
}
