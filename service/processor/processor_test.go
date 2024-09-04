package processor

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/client/modelstest"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCurationExportSyncProcessor_Run(t *testing.T) {
	integrationID := uuid.NewString()
	inputDirectory := t.TempDir()
	outputDirectory := t.TempDir()

	datasetID := uuid.NewString()

	processor := NewCurationExportSyncProcessor(integrationID, inputDirectory, outputDirectory)

	// Setup Input directory with pre-requisites
	// curation-export file
	curationExport := modelstest.NewDatasetCurationExportBuilder(datasetID).Build()
	curationExportFile, err := os.Create(processor.curationExportPath())
	require.NoError(t, err)
	require.NoError(t, json.NewEncoder(curationExportFile).Encode(curationExport))

	require.NoError(t, processor.Run())
}
