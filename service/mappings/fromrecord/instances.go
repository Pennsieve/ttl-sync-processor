package fromrecord

import (
	"fmt"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings"
	"log/slog"
)

func ToSavedDatasetMetadata(reader *metadataclient.Reader) (*metadata.SavedDatasetMetadata, error) {
	existing := &metadata.SavedDatasetMetadata{}
	var err error
	if existing.Contributors, err = MapRecords(reader, metadata.ContributorModelName, ToContributor); err != nil {
		return nil, err
	}
	if existing.Subjects, err = MapRecords(reader, metadata.SubjectModelName, ToSubject); err != nil {
		return nil, err
	}
	if existing.Samples, err = MapRecords(reader, metadata.SampleModelName, ToSample); err != nil {
		return nil, err
	}
	sampleSubjectMapping := NewSampleStoreMapping(existing.Samples, existing.Subjects)
	if existing.SampleSubjects, err = MapLinkedProperties(reader, metadata.SampleSubjectLinkName, sampleSubjectMapping); err != nil {
		return nil, err
	}
	return existing, nil
}

func MapRecords[To any](reader *metadataclient.Reader, modelName string, mapping mappings.Mapping[instance.Record, To]) ([]To, error) {
	model, modelExists := reader.ModelNamesToSchemaElements[modelName]
	if !modelExists {
		logger.Warn("model does not exist", slog.String("modelName", modelName))
		return []To{}, nil
	}
	logger.Info("reading existing records", slog.String("modelName", modelName),
		slog.String("modelID", model.ID))
	records, err := reader.GetRecordsForModel(modelName)
	if err != nil {
		return nil, fmt.Errorf("error reading %s records: %w", modelName, err)
	}
	mapped, err := mappings.MapSlice[instance.Record, To](records, mapping)
	if err != nil {
		return nil, fmt.Errorf("error marshalling %s records: %w", modelName, err)
	}
	logger.Info("read existing records", slog.String("modelName", modelName),
		slog.Int("count", len(mapped)))
	return mapped, nil
}

func MapLinkedProperties[To any](reader *metadataclient.Reader, linkedPropertyName string, mapping mappings.Mapping[instance.LinkedProperty, To]) ([]To, error) {
	linkedProperty, linkedPropertyExists := reader.LinkedPropNamesToSchemaElements[linkedPropertyName]
	if !linkedPropertyExists {
		logger.Warn("linkedProperty does not exist", slog.String("linkedPropertyName", linkedPropertyName))
		return []To{}, nil
	}
	logger.Info("reading existing linked properties", slog.String("linkedPropertyName", linkedPropertyName),
		slog.String("linkedPropertyID", linkedProperty.ID))
	linkedPropertyInstances, err := reader.GetLinkInstancesForProperty(linkedPropertyName)
	if err != nil {
		return nil, fmt.Errorf("error reading %s linked property instances: %w", linkedPropertyName, err)
	}
	mapped, err := mappings.MapSlice[instance.LinkedProperty, To](linkedPropertyInstances, mapping)
	if err != nil {
		return nil, fmt.Errorf("error marshalling %s linked property instances: %w", linkedPropertyName, err)
	}
	logger.Info("read existing linked property instances", slog.String("linkedPropertyName", linkedPropertyName),
		slog.Int("count", len(mapped)))
	return mapped, nil
}
