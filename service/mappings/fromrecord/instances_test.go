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

	idMap := NewRecordIDStore()
	existingMetadata, err := ToSavedDatasetMetadata(reader, idMap)

	require.NoError(t, err)
	assert.NotNil(t, existingMetadata)
	assert.Len(t, existingMetadata.Contributors, 5)
	assert.Len(t, existingMetadata.Subjects, 3)
	assert.Len(t, existingMetadata.Samples, 3)
	assert.Len(t, existingMetadata.SampleSubjects, 2)
	assert.Len(t, existingMetadata.Proxies, 3)
	assert.Equal(t, len(existingMetadata.Samples)+len(existingMetadata.Subjects), idMap.Len())
}
func TestToDatasetMetadata_NoModels(t *testing.T) {
	inputDirectory := "testdata/input_no_model"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)

	idMap := NewRecordIDStore()
	existingMetadata, err := ToSavedDatasetMetadata(reader, idMap)
	require.NoError(t, err)
	assert.Empty(t, existingMetadata.Contributors)
	assert.Empty(t, idMap.Len())

}

func TestToDatasetMetadata_NoRecords(t *testing.T) {
	inputDirectory := "testdata/input_no_records"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)

	idMap := NewRecordIDStore()
	existingMetadata, err := ToSavedDatasetMetadata(reader, idMap)
	require.NoError(t, err)
	assert.Empty(t, existingMetadata.Contributors)
	assert.Empty(t, idMap.Len())

}
