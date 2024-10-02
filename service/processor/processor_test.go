package processor

import (
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
	"github.com/pennsieve/ttl-sync-processor/service/mappings/fromrecord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestCurationExportSyncProcessor_Run(t *testing.T) {
	currentLogLevel := logging.Level.Level()
	logging.Level.Set(slog.LevelDebug)
	t.Cleanup(func() {
		logging.Level.Set(currentLogLevel)
	})
	integrationID := uuid.NewString()
	inputDirectory := "testdata/input"
	outputDirectory := t.TempDir()

	processor, err := NewCurationExportSyncProcessor(integrationID, inputDirectory, outputDirectory)
	require.NoError(t, err)

	require.NoError(t, processor.Run())

	assert.FileExists(t, processor.ChangesetFilePath())
}

func TestCurationExportSyncProcessor_ReadExistingPennsieveMetadata(t *testing.T) {
	integrationID := uuid.NewString()
	inputDirectory := "testdata/input"
	outputDirectory := t.TempDir()

	processor, err := NewCurationExportSyncProcessor(integrationID, inputDirectory, outputDirectory)
	require.NoError(t, err)

	schemaData := fromrecord.SchemaData{}
	existingMetadata, err := processor.ExistingPennsieveMetadata(schemaData)
	require.NoError(t, err)
	assert.NotNil(t, existingMetadata)
	assert.Len(t, existingMetadata.Contributors, 5)

	assert.Equal(t, "d77470bb-f39d-49ee-aa17-783e128cdfa0", schemaData[metadata.ContributorModelName])

}

func TestCurationExportSyncProcessor_ReadCurationExport(t *testing.T) {
	integrationID := uuid.NewString()
	inputDirectory := "testdata/input"
	outputDirectory := t.TempDir()

	processor, err := NewCurationExportSyncProcessor(integrationID, inputDirectory, outputDirectory)
	require.NoError(t, err)

	export, err := processor.ReadCurationExport()
	require.NoError(t, err)
	assert.NotNil(t, export)
	assert.Len(t, export.Contributors, 3)

}
