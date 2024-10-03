package fromrecord

import (
	"fmt"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"log/slog"
)

func ToDatasetMetadata(reader *metadataclient.Reader) (*metadata.DatasetMetadata, error) {
	contributors, err := MapRecords(reader, metadata.ContributorModelName, ToContributor)
	if err != nil {
		return nil, err
	}
	subjects, err := MapRecords(reader, metadata.SubjectModelName, ToSubject)
	if err != nil {
		return nil, err
	}
	existing := &metadata.DatasetMetadata{
		Contributors: contributors,
		Subjects:     subjects,
	}
	return existing, nil
}

func MapRecords[T any](reader *metadataclient.Reader, modelName string, mapping Mapping[T]) ([]T, error) {
	model, modelExists := reader.ModelNamesToSchemaElements[modelName]
	if !modelExists {
		logger.Warn("model does not exist", slog.String("modelName", modelName))
		return []T{}, nil
	}
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
