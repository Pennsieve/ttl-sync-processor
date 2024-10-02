package fromrecord

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToDatasetMetadata_NoModels(t *testing.T) {
	inputDirectory := "testdata/input_no_model"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)
	existingMetadata, err := ToDatasetMetadata(reader)
	require.NoError(t, err)
	assert.Empty(t, existingMetadata.Contributors)

}

func TestToDatasetMetadata_NoRecords(t *testing.T) {
	inputDirectory := "testdata/input_no_records"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)
	existingMetadata, err := ToDatasetMetadata(reader)
	require.NoError(t, err)
	assert.Empty(t, existingMetadata.Contributors)

}
