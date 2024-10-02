package sync

import (
	"fmt"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
)

var logger = logging.PackageLogger("sync")

func ComputeChangeset(schemaData map[string]schema.Element, old, new *metadata.DatasetMetadata) (*changesetmodels.Dataset, error) {
	datasetChanges := &changesetmodels.Dataset{}
	contributorChanges, err := ComputeContributorsChanges(schemaData, old, new)
	if err != nil {
		return nil, err
	}
	if contributorChanges != nil {
		datasetChanges.Models = append(datasetChanges.Models, *contributorChanges)
	}
	return datasetChanges, nil
}

func newStringPropertyCreate(propertyName, displayName string) (changesetmodels.PropertyCreate, error) {
	propCreate := &changesetmodels.PropertyCreate{
		DisplayName: displayName,
		Name:        propertyName,
	}
	if err := propCreate.SetDataType(datatypes.StringType); err != nil {
		return *propCreate, fmt.Errorf("error setting data type of %s %s to %s: %w", propertyName,
			displayName, datatypes.StringType, err)
	}
	return *propCreate, nil
}

func newStringConceptTitlePropertyCreate(propertyName, displayName string) (changesetmodels.PropertyCreate, error) {
	propCreate := &changesetmodels.PropertyCreate{
		DisplayName:  displayName,
		Name:         propertyName,
		ConceptTitle: true,
		Required:     true,
	}
	if err := propCreate.SetDataType(datatypes.StringType); err != nil {
		return *propCreate, fmt.Errorf("error setting data type of %s %s to %s: %w", propertyName,
			displayName, datatypes.StringType, err)
	}
	return *propCreate, nil
}

type stringPropertyCreator func(propertyName, displayName string) (changesetmodels.PropertyCreate, error)

func appendStringPropertyCreate(creates []changesetmodels.PropertyCreate, propertyName, displayName string, propCreator stringPropertyCreator, errs *[]error) []changesetmodels.PropertyCreate {
	create, err := propCreator(propertyName, displayName)
	if err != nil {
		*errs = append(*errs, err)
		return creates
	}
	creates = append(creates, create)
	return creates

}

// appendStringRecordValue only appends a new changesetmodels.RecordValue if len(value) > 0
func appendStringRecordValue(values []changesetmodels.RecordValue, name string, value string) []changesetmodels.RecordValue {
	if len(value) > 0 {
		return append(values, changesetmodels.RecordValue{
			Value: value,
			Name:  name,
		})
	}
	return values
}
