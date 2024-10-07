package fromcuration

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToSample(t *testing.T) {
	curationExportPath := "testdata/curation-export.json"

	datasetExport, err := UnmarshalDatasetExport(curationExportPath)
	require.NoError(t, err)

	samples, err := MapSlice(datasetExport.Samples, ToSample)
	require.NoError(t, err)
	assert.Len(t, samples, 3)

	assert.Equal(t, "09d2a327-be38-403a-884d-a4d1d98b732c", samples[0].ID)
	assert.Equal(t, "09d2a327-be38-403a-884d-a4d1d98b732c", samples[0].GetID())

	assert.Equal(t, "80e9d14c-188e-476d-99ae-464da9e68bc3", samples[1].ID)
	assert.Equal(t, "80e9d14c-188e-476d-99ae-464da9e68bc3", samples[1].GetID())

	assert.Equal(t, "70a38ffa-15f6-4fdc-87ea-2f3a9b19f995", samples[2].ID)
	assert.Equal(t, "70a38ffa-15f6-4fdc-87ea-2f3a9b19f995", samples[2].GetID())

}
