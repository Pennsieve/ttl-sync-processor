package fromrecord

import (
	"fmt"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"log/slog"
)

func ToDatasetMetadata(reader *metadataclient.Reader, schemaData SchemaData) (*metadata.DatasetMetadata, error) {
	contributors, err := MapRecords(reader, schemaData, metadata.ContributorModelName, ToContributor)
	if err != nil {
		return nil, err
	}
	existing := &metadata.DatasetMetadata{
		Contributors: contributors,
	}
	return existing, nil
}

func MapRecords[T any](reader *metadataclient.Reader, schemaData SchemaData, modelName string, mapping Mapping[T]) ([]T, error) {
	model, modelExists := reader.ModelNamesToSchemaElements[modelName]
	if !modelExists {
		schemaData[modelName] = &changesetmodels.ModelCreate{Name: modelName}
		logger.Warn("model does not exist", slog.String("modelName", modelName))
		return []T{}, nil
	}
	schemaData[modelName] = model.ID
	logger.Info("reading existing records", slog.String("modelName", modelName),
		slog.String("modelID", model.ID))
	records, err := reader.GetRecordsForModel(modelName)
	if err != nil {
		return nil, fmt.Errorf("error reading %s records: %w", modelName, err)
	}
	mapped, err := MapSlice(records, mapping)
	if err != nil {
		return nil, fmt.Errorf("error marshalling %s records: %w", modelName, err)
	}
	logger.Info("read existing records", slog.String("modelName", modelName),
		slog.Int("count", len(mapped)))
	return mapped, nil
}
