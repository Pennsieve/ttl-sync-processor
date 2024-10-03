package fromcuration

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToDatasetMetadata(t *testing.T) {
	curationExportPath := "testdata/curation-export.json"

	datasetExport, err := UnmarshalDatasetExport(curationExportPath)
	require.NoError(t, err)

	datasetMetadata, err := ToDatasetMetadata(datasetExport)
	require.NoError(t, err)
	assert.Len(t, datasetMetadata.Contributors, 3)
	assert.Len(t, datasetMetadata.Subjects, 4)
}
