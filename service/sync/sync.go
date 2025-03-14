package sync

import (
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
	"github.com/pennsieve/ttl-sync-processor/service/mappings/fromrecord"
)

var logger = logging.PackageLogger("sync")

// ExistingRecordStore is populated by the processor as it maps the existing, raw metadata to the models
// used for the sync. It will hold a mapping of all existing records of the form
// (modelName, recordExternalID) -> recordPennsieveID
// This store will have to be set and reset for tests.
var ExistingRecordStore *fromrecord.RecordIDStore

// recordIDMap is populated by this package and then passed off to the post-processor file we are creating in this package.
// Like ExistingRecordStore it is map (modelName, recordExternalID) -> recordPennsieveID, but it will only contain entries
// for existing records that take part in a link or package proxy "created" by this sync. Since we are only passing the post-processor
// external IDs, the post-processor will need this map to find the corresponding pennsieve record ids.
// It will also have to be set and reset for tests
var recordIDMap = make(fromrecord.RecordIDMap)

// ComputeChangeset is the entrypoint for computing the changes necessary to sync the dataset's Pennsieve metadata with that
// found in the curation export file.
// schemaData is the current state of the Pennsieve metadata schema as downloaded by the metadata pre-processor
// old is the current state of the Pennsieve metadata instances, i.e., records, as downloaded by the metadata pre-processor
// new is the state of the exported curation metadata as downloaded by the external file downloader
func ComputeChangeset(schemaData *metadataclient.Schema, old *metadata.SavedDatasetMetadata, new *metadata.DatasetMetadata) (*changesetmodels.Dataset, error) {
	datasetChanges := newDatasetChanges(schemaData)
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

	appendRecordIDMap(datasetChanges)

	return datasetChanges, nil
}

// modelChangeComputer instances are functions that compute the changes needed for a particular model.
// It should return nil if no changes are needed.
// It should return a ModelCreate or non-nil *ModelCreate if the model does not exist and needs to be created.
// It should return a ModelUpdate or non-nil *ModelUpdate if the model already exists.
// Any other return type will cause an error.
// OLD is the type of existing records for the model, for example metadata.Contributor or metadata.SavedSubject
// NEW is the type of new records for the model, for example metadata.Contributor or metadata.Subject
type modelChangeComputer[OLD, NEW any] func(schemaData *metadataclient.Schema, old []OLD, new []NEW) (any, error)

func appendModelChanges[OLD, NEW any](datasetChanges *changesetmodels.Dataset, schemaData *metadataclient.Schema, old []OLD, new []NEW, computer modelChangeComputer[OLD, NEW]) error {
	modelChanges, err := computer(schemaData, old, new)
	if err != nil {
		return err
	}
	switch changes := modelChanges.(type) {
	case nil:
	//Do nothing, no changes for this model
	case *changesetmodels.ModelCreate:
		datasetChanges.Models.Creates = append(datasetChanges.Models.Creates, *changes)
	case changesetmodels.ModelCreate:
		datasetChanges.Models.Creates = append(datasetChanges.Models.Creates, changes)
	case *changesetmodels.ModelUpdate:
		datasetChanges.Models.Updates = append(datasetChanges.Models.Updates, *changes)
	case changesetmodels.ModelUpdate:
		datasetChanges.Models.Updates = append(datasetChanges.Models.Updates, changes)
	default:
		return fmt.Errorf("unkown result type from modelChangeComputer: %T", changes)
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

func appendRecordIDMap(datasetChanges *changesetmodels.Dataset) {
	nested := map[string]map[changesetmodels.ExternalInstanceID]changesetmodels.PennsieveInstanceID{}
	for key, id := range recordIDMap {
		idMap, found := nested[key.ModelName]
		if !found {
			idMap = make(map[changesetmodels.ExternalInstanceID]changesetmodels.PennsieveInstanceID)
			nested[key.ModelName] = idMap
		}
		idMap[key.ExternalRecordID] = id
	}

	for modelName, idMap := range nested {
		datasetChanges.RecordIDMaps = append(datasetChanges.RecordIDMaps, changesetmodels.RecordIDMap{
			ModelName:           modelName,
			ExternalToPennsieve: idMap,
		})
	}
}

func newDatasetChanges(schema *metadataclient.Schema) *changesetmodels.Dataset {
	existingModelIDs := make(map[string]changesetmodels.PennsieveSchemaID, len(schema.ModelIDsByName()))
	for name, id := range schema.ModelIDsByName() {
		existingModelIDs[name] = changesetmodels.PennsieveSchemaID(id)
	}
	return &changesetmodels.Dataset{ExistingModelIDMap: existingModelIDs}
}
