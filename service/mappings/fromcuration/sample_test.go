package fromcuration

import (
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
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

	metadatatest.AssertExternalInstanceIDEqual(t, "09d2a327-be38-403a-884d-a4d1d98b732c", samples[0].ID)
	metadatatest.AssertExternalInstanceIDEqual(t, "09d2a327-be38-403a-884d-a4d1d98b732c", samples[0].ExternalID())
	assert.Equal(t, "f61ed5a9-5a69-49f7-9113-6447ee9e668b", samples[0].PrimaryKey)

	metadatatest.AssertExternalInstanceIDEqual(t, "80e9d14c-188e-476d-99ae-464da9e68bc3", samples[1].ID)
	metadatatest.AssertExternalInstanceIDEqual(t, "80e9d14c-188e-476d-99ae-464da9e68bc3", samples[1].ExternalID())
	assert.Equal(t, "ef07f30e-ea93-4221-86e2-0be3a278ff55", samples[1].PrimaryKey)

	metadatatest.AssertExternalInstanceIDEqual(t, "70a38ffa-15f6-4fdc-87ea-2f3a9b19f995", samples[2].ID)
	metadatatest.AssertExternalInstanceIDEqual(t, "70a38ffa-15f6-4fdc-87ea-2f3a9b19f995", samples[2].ExternalID())
	assert.Equal(t, "b07f1212-72c6-407e-b7b0-2020c9a35b19", samples[2].PrimaryKey)

}
