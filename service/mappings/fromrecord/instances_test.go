package fromrecord

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToDatasetMetadata(t *testing.T) {
	inputDirectory := "testdata/input"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)

	existingMetadata, err := ToSavedDatasetMetadata(reader)
	require.NoError(t, err)
	assert.NotNil(t, existingMetadata)
	assert.Len(t, existingMetadata.Contributors, 5)
	assert.Len(t, existingMetadata.Subjects, 2)
	assert.Len(t, existingMetadata.Samples, 2)
	assert.Len(t, existingMetadata.SampleSubjects, 2)
}
func TestToDatasetMetadata_NoModels(t *testing.T) {
	inputDirectory := "testdata/input_no_model"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)
	existingMetadata, err := ToSavedDatasetMetadata(reader)
	require.NoError(t, err)
	assert.Empty(t, existingMetadata.Contributors)

}

func TestToDatasetMetadata_NoRecords(t *testing.T) {
	inputDirectory := "testdata/input_no_records"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)
	existingMetadata, err := ToSavedDatasetMetadata(reader)
	require.NoError(t, err)
	assert.Empty(t, existingMetadata.Contributors)

}
