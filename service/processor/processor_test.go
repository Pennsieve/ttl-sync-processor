package processor

import (
	"encoding/json"
	"github.com/google/uuid"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
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

	// Check changes contents
	changesetFile, err := os.Open(processor.ChangesetFilePath())
	defer changesetFile.Close()
	var changeset changesetmodels.Dataset
	require.NoError(t, json.NewDecoder(changesetFile).Decode(&changeset))

	modelChanges := changeset.Models
	assert.Len(t, modelChanges, 1)
	contributorChanges := modelChanges[0]
	// model exists, so ID should be present and Create nil
	assert.Equal(t, "d77470bb-f39d-49ee-aa17-783e128cdfa0", contributorChanges.ID)
	assert.Nil(t, contributorChanges.Create)

	assert.True(t, contributorChanges.Records.DeleteAll)
	assert.Len(t, contributorChanges.Records.Create, 3)
	assert.Empty(t, contributorChanges.Records.Update)
	assert.Empty(t, contributorChanges.Records.Delete)
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
