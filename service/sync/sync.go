package sync

import (
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
)

var logger = logging.PackageLogger("sync")

// ComputeChangeset is the entrypoint for computing the changes necessary to sync the dataset's Pennsieve metadata with that
// found in the curation export file.
// schemaData is the current state of the Pennsieve metadata schema as downloaded by the metadata pre-processor
// old is the current state of the Pennsieve metadata instances, i.e., records, as downloaded by the metadata pre-processor
// new is the state of the exported curation metadata as downloaded by the external file downloader
func ComputeChangeset(schemaData *metadataclient.Schema, old *metadata.SavedDatasetMetadata, new *metadata.DatasetMetadata) (*changesetmodels.Dataset, error) {
	datasetChanges := &changesetmodels.Dataset{}
	if err := appendModelChanges(datasetChanges, schemaData, old.Contributors, new.Contributors, ComputeContributorsChanges); err != nil {
		return nil, err
	}
	if err := appendModelChanges(datasetChanges, schemaData, old.Subjects, new.Subjects, ComputeSubjectChanges); err != nil {
		return nil, err
	}
	if err := appendModelChanges(datasetChanges, schemaData, old.Samples, new.Samples, ComputeSampleChanges); err != nil {
		return nil, err
	}

	if err := appendLinkedPropertyChanges(datasetChanges, schemaData, old.SampleSubjects, new.SampleSubjects, ComputeSampleSubjectChanges); err != nil {
		return nil, err
	}
	proxyChanges, err := ComputeProxyChanges(schemaData, old.Proxies, new.Proxies)
	if err != nil {
		return nil, err
	}
	datasetChanges.Proxies = proxyChanges
	return datasetChanges, nil
}

// modelChangeComputer instances are functions that compute the changes needed for a particular model.
// It should return nil if no changes are needed.
// OLD is the type of existing records for the model, for example metadata.Contributor or metadata.SavedSubject
// NEW is the type of new records for the model, for example metadata.Contributor or metadata.Subject
type modelChangeComputer[OLD, NEW any] func(schemaData *metadataclient.Schema, old []OLD, new []NEW) (*changesetmodels.ModelChanges, error)

func appendModelChanges[OLD, NEW any](datasetChanges *changesetmodels.Dataset, schemaData *metadataclient.Schema, old []OLD, new []NEW, computer modelChangeComputer[OLD, NEW]) error {
	modelChanges, err := computer(schemaData, old, new)
	if err != nil {
		return err
	}
	if modelChanges != nil {
		datasetChanges.Models = append(datasetChanges.Models, *modelChanges)
	}
	return nil
}

type linkedPropertyChangeComputer[OLD metadata.SavedExternalLink, NEW metadata.ExternalLink] func(schemaData *metadataclient.Schema, old []OLD, new []NEW) (*changesetmodels.LinkedPropertyChanges, error)

func appendLinkedPropertyChanges[OLD metadata.SavedExternalLink, NEW metadata.ExternalLink](datasetChanges *changesetmodels.Dataset, schemaData *metadataclient.Schema, old []OLD, new []NEW, computer linkedPropertyChangeComputer[OLD, NEW]) error {
	linkedPropertyChanges, err := computer(schemaData, old, new)
	if err != nil {
		return err
	}
	if linkedPropertyChanges != nil {
		datasetChanges.LinkedProperties = append(datasetChanges.LinkedProperties, *linkedPropertyChanges)
	}
	return nil
}
