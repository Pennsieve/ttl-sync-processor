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

	processor := NewCurationExportSyncProcessor(integrationID, inputDirectory, outputDirectory)

	require.NoError(t, processor.Run())
}

func TestCurationExportSyncProcessor_ReadExistingPennsieveMetadata(t *testing.T) {
	integrationID := uuid.NewString()
	inputDirectory := "testdata/input"
	outputDirectory := t.TempDir()

	processor := NewCurationExportSyncProcessor(integrationID, inputDirectory, outputDirectory)
	existingMetadata, err := processor.ReadExistingPennsieveMetadata()
	require.NoError(t, err)
	assert.NotNil(t, existingMetadata)
	assert.Len(t, existingMetadata.Contributors, 5)

	contrib1 := existingMetadata.Contributors[0]
	assert.Equal(t, "Elara", contrib1.FirstName)
	assert.Equal(t, "Lauridsen", contrib1.LastName)
	assert.Empty(t, contrib1.Degree)
	assert.Empty(t, contrib1.NodeID)
	assert.Empty(t, contrib1.MiddleInitial)
	assert.Empty(t, contrib1.ORCID)

	contrib2 := existingMetadata.Contributors[1]
	assert.Equal(t, "Yordanka", contrib2.FirstName)
	assert.Equal(t, "Vukoja", contrib2.LastName)
	assert.Equal(t, "PHD", contrib2.Degree)
	assert.Empty(t, contrib2.NodeID)
	assert.Equal(t, "T", contrib2.MiddleInitial)
	assert.Empty(t, contrib2.ORCID)

	contrib5 := existingMetadata.Contributors[4]
	assert.Equal(t, "Ajay", contrib5.FirstName)
	assert.Equal(t, "Carstensen", contrib5.LastName)
	assert.Empty(t, contrib5.Degree)
	assert.Equal(t, "N:user:3478dd52-e229-4e95-ab23-c6bd1e3d4d25", contrib5.NodeID)
	assert.Empty(t, contrib5.MiddleInitial)
	assert.Equal(t, "a1482559-3881-4466-b98f-d4240d9d467d", contrib5.ORCID)
}

func TestCurationExportSyncProcessor_ReadCurationExport(t *testing.T) {
	integrationID := uuid.NewString()
	inputDirectory := "testdata/input"
	outputDirectory := t.TempDir()

	processor := NewCurationExportSyncProcessor(integrationID, inputDirectory, outputDirectory)
	export, err := processor.ReadCurationExport()
	require.NoError(t, err)
	assert.NotNil(t, export)
	assert.Len(t, export.Contributors, 2)
}

func TestCurationExportSyncProcessor_ConvertCurationExport(t *testing.T) {
	integrationID := uuid.NewString()
	inputDirectory := "testdata/input"
	outputDirectory := t.TempDir()

	processor := NewCurationExportSyncProcessor(integrationID, inputDirectory, outputDirectory)
	export, err := processor.ReadCurationExport()
	require.NoError(t, err)

	exported, err := processor.ConvertCurationExport(export)
	require.NoError(t, err)
	assert.NotNil(t, exported)
	assert.Len(t, exported.Contributors, 2)
}
