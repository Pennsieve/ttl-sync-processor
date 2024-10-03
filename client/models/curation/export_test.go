package curation_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/client/curationtest"
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	datasetID := uuid.NewString()
	simpleSubject, err := curationtest.NewSubjectBuilder().WithSimpleSpecies().Build()
	require.NoError(t, err)
	complexSubject, err := curationtest.NewSubjectBuilder().WithSubjectIdentifier(0).Build()
	require.NoError(t, err)
	complexSubjectWithSynonyms, err := curationtest.NewSubjectBuilder().WithSubjectIdentifier(3).Build()
	require.NoError(t, err)

	curationExport := curation.NewDatasetExport(datasetID).
		WithContributors(
			curationtest.NewContributorBuilder().Build(),
			curationtest.NewContributorBuilder().WithRoles(1).Build(),
			curationtest.NewContributorBuilder().WithRoles(2).WithAffiliation().Build(),
			curationtest.NewContributorBuilder().WithRoles(3).WithORCID().Build(),
			curationtest.NewContributorBuilder().WithRoles(4).WithORCID().WithAffiliation().Build(),
		).
		WithDirStructureEntries(
			curation.NewDirStructureEntry(uuid.NewString(), uuid.NewString()),
			curation.NewDirStructureEntry(uuid.NewString(), uuid.NewString()),
		).
		WithSpecimenDirs(
			*curation.NewSpecimenDirs().
				WithRecord(uuid.NewString(), curation.SampleRecordType, uuid.NewString()).
				WithRecord(uuid.NewString(), curation.SubjectRecordType, uuid.NewString()),
		).
		WithSubjects(
			simpleSubject,
			complexSubject,
			complexSubjectWithSynonyms,
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
	var fromFile *curation.DatasetExport
	require.NoError(t, json.NewDecoder(actualFile).Decode(&fromFile))

	assert.Equal(t, curationExport, fromFile)
}
