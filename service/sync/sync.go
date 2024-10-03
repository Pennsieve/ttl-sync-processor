package sync

import (
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
)

var logger = logging.PackageLogger("sync")

func ComputeChangeset(schemaData map[string]schema.Element, old *metadata.SavedDatasetMetadata, new *metadata.DatasetMetadata) (*changesetmodels.Dataset, error) {
	datasetChanges := &changesetmodels.Dataset{}
	if err := appendModelChanges(datasetChanges, schemaData, old, new, ComputeContributorsChanges); err != nil {
		return nil, err
	}
	if err := appendModelChanges(datasetChanges, schemaData, old, new, ComputeSubjectChanges); err != nil {
		return nil, err
	}
	return datasetChanges, nil
}

type modelChangeComputer func(schemaData map[string]schema.Element, old *metadata.SavedDatasetMetadata, new *metadata.DatasetMetadata) (*changesetmodels.ModelChanges, error)

func appendModelChanges(datasetChanges *changesetmodels.Dataset, schemaData map[string]schema.Element, old *metadata.SavedDatasetMetadata, new *metadata.DatasetMetadata, computer modelChangeComputer) error {
	modelChanges, err := computer(schemaData, old, new)
	if err != nil {
		return err
	}
	if modelChanges != nil {
		datasetChanges.Models = append(datasetChanges.Models, *modelChanges)
	}
	return nil
}
