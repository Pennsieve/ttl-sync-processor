package fromrecord

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToDatasetMetadata_NoModels(t *testing.T) {
	inputDirectory := "testdata/input_no_model"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)
	schemaData := SchemaData{}
	existingMetadata, err := ToDatasetMetadata(reader, schemaData)
	require.NoError(t, err)
	assert.Empty(t, existingMetadata.Contributors)

	assert.Contains(t, schemaData, metadata.ContributorModelName)
	modelCreate := schemaData[metadata.ContributorModelName]
	assert.IsType(t, &changesetmodels.ModelCreate{}, modelCreate)
	assert.Equal(t, metadata.ContributorModelName, modelCreate.(*changesetmodels.ModelCreate).Name)
}

func TestToDatasetMetadata_NoRecords(t *testing.T) {
	inputDirectory := "testdata/input_no_records"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)
	schemaData := SchemaData{}
	existingMetadata, err := ToDatasetMetadata(reader, schemaData)
	require.NoError(t, err)
	assert.Empty(t, existingMetadata.Contributors)

	assert.Contains(t, schemaData, metadata.ContributorModelName)
	assert.Equal(t, "d77470bb-f39d-49ee-aa17-783e128cdfa0", schemaData[metadata.ContributorModelName])
}
