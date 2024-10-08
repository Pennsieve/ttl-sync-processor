package fromcuration

import (
	"fmt"
	"github.com/pennsieve/ttl-sync-processor/service/mappings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToSample(t *testing.T) {
	curationExportPath := "testdata/curation-export.json"

	datasetExport, err := UnmarshalDatasetExport(curationExportPath)
	require.NoError(t, err)

	samples, err := mappings.MapSlice(datasetExport.Samples, ToSample)
	require.NoError(t, err)
	assert.Len(t, samples, 3)

	assert.Equal(t, "09d2a327-be38-403a-884d-a4d1d98b732c", samples[0].ID)
	assert.Equal(t, "09d2a327-be38-403a-884d-a4d1d98b732c", samples[0].GetID())

	assert.Equal(t, "80e9d14c-188e-476d-99ae-464da9e68bc3", samples[1].ID)
	assert.Equal(t, "80e9d14c-188e-476d-99ae-464da9e68bc3", samples[1].GetID())

	assert.Equal(t, "70a38ffa-15f6-4fdc-87ea-2f3a9b19f995", samples[2].ID)
	assert.Equal(t, "70a38ffa-15f6-4fdc-87ea-2f3a9b19f995", samples[2].GetID())

}

func TestToSampleSubjectLink(t *testing.T) {
	curationExportPath := "testdata/curation-export.json"

	datasetExport, err := UnmarshalDatasetExport(curationExportPath)
	require.NoError(t, err)

	sampleSubjectLink, err := mappings.MapSlice(datasetExport.Samples, ToSampleSubjectLink)
	require.NoError(t, err)
	assert.Len(t, sampleSubjectLink, 3)

	assert.Equal(t, "09d2a327-be38-403a-884d-a4d1d98b732c", sampleSubjectLink[0].SampleID)
	assert.Equal(t, "f1027e6e-17cf-45d7-8b57-4c91bfd93fce", sampleSubjectLink[0].SubjectID)
	assert.Equal(t, fmt.Sprintf("%s:%s", sampleSubjectLink[0].SampleID, sampleSubjectLink[0].SubjectID),
		sampleSubjectLink[0].GetID())

	assert.Equal(t, "80e9d14c-188e-476d-99ae-464da9e68bc3", sampleSubjectLink[1].SampleID)
	assert.Equal(t, "f1027e6e-17cf-45d7-8b57-4c91bfd93fce", sampleSubjectLink[1].SubjectID)
	assert.Equal(t, fmt.Sprintf("%s:%s", sampleSubjectLink[1].SampleID, sampleSubjectLink[1].SubjectID),
		sampleSubjectLink[1].GetID())

	assert.Equal(t, "70a38ffa-15f6-4fdc-87ea-2f3a9b19f995", sampleSubjectLink[2].SampleID)
	assert.Equal(t, "9b942959-2ebb-4680-b3f7-6d1cdbbe706d", sampleSubjectLink[2].SubjectID)
	assert.Equal(t, fmt.Sprintf("%s:%s", sampleSubjectLink[2].SampleID, sampleSubjectLink[2].SubjectID),
		sampleSubjectLink[2].GetID())

}
