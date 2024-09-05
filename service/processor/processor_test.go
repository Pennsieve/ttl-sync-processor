package processor

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/client/models"
	"github.com/pennsieve/ttl-sync-processor/client/modelstest"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
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
	inputDirectory := t.TempDir()
	outputDirectory := t.TempDir()

	datasetID := uuid.NewString()

	processor := NewCurationExportSyncProcessor(integrationID, inputDirectory, outputDirectory)

	// Setup Input directory with pre-requisites
	// curation-export file
	subjectDirPath := uuid.NewString()
	sampleDirPath1 := uuid.NewString()
	sampleDirPath2 := uuid.NewString()
	curationExport := models.NewDatasetCurationExport(datasetID).
		WithContributors(
			modelstest.NewContributorBuilder().Build(),
			modelstest.NewContributorBuilder().WithRoles(2).WithAffiliation().Build(),
			modelstest.NewContributorBuilder().WithRoles(1).WithORCID().Build(),
		).
		WithDirStructureEntries(
			models.NewDirStructureEntry(subjectDirPath, uuid.NewString()),
			models.NewDirStructureEntry(sampleDirPath1, uuid.NewString()),
			models.NewDirStructureEntry(sampleDirPath2, uuid.NewString()),
		).
		WithSpecimenDirs(*models.NewSpecimenDirs().
			WithRecord(uuid.NewString(), models.SubjectRecordType, subjectDirPath).
			WithRecord(uuid.NewString(), models.SampleRecordType, sampleDirPath1, sampleDirPath2),
		)

	curationExportFile, err := os.Create(processor.curationExportPath())
	require.NoError(t, err)
	require.NoError(t, json.NewEncoder(curationExportFile).Encode(curationExport))

	require.NoError(t, processor.Run())
}
