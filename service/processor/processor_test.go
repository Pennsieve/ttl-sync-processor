package processor

import (
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
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

	existingMetadata, err := processor.ExistingPennsieveMetadata()
	require.NoError(t, err)
	assert.NotNil(t, existingMetadata)
	assert.Len(t, existingMetadata.Contributors, 5)

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
