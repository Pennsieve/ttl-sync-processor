package models_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/client/models"
	"github.com/pennsieve/ttl-sync-processor/client/modelstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	datasetID := uuid.NewString()
	curationExport := models.NewDatasetCurationExport(datasetID).
		WithContributors(
			modelstest.NewContributorBuilder().Build(),
			modelstest.NewContributorBuilder().WithRoles(1).Build(),
			modelstest.NewContributorBuilder().WithRoles(2).WithAffiliation().Build(),
			modelstest.NewContributorBuilder().WithRoles(3).WithORCID().Build(),
			modelstest.NewContributorBuilder().WithRoles(4).WithORCID().WithAffiliation().Build(),
		).
		WithDirStructureEntries(
			models.NewDirStructureEntry(uuid.NewString(), uuid.NewString()),
			models.NewDirStructureEntry(uuid.NewString(), uuid.NewString()),
		).
		WithSpecimenDirs(
			*models.NewSpecimenDirs().
				WithRecord(uuid.NewString(), models.SampleRecordType, uuid.NewString()).
				WithRecord(uuid.NewString(), models.SubjectRecordType, uuid.NewString()),
		)

	directory := t.TempDir()
	path := filepath.Join(directory, "curation-export.json")
	file, err := os.Create(path)
	require.NoError(t, err)
	defer file.Close()

	require.NoError(t, json.NewEncoder(file).Encode(curationExport))
	require.NoError(t, file.Sync())

	actualFile, err := os.Open(path)
	require.NoError(t, err)
	defer actualFile.Close()
	var fromFile *models.DatasetCurationExport
	require.NoError(t, json.NewDecoder(actualFile).Decode(&fromFile))

	assert.Equal(t, curationExport, fromFile)
}
