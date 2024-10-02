package fromcuration

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToContributor(t *testing.T) {
	inputDirectory := "testdata/curation-export.json"

	datasetExport, err := UnmarshalDatasetExport(inputDirectory)
	require.NoError(t, err)

	exported, err := MapSlice(datasetExport.Contributors, ToContributor)
	require.NoError(t, err)
	assert.NotNil(t, exported)
	assert.Len(t, exported, 3)

	contrib1 := exported[0]
	assert.Equal(t, "Elara", contrib1.FirstName)
	assert.Equal(t, "Lauridsen", contrib1.LastName)
	assert.Empty(t, contrib1.Degree)
	assert.Empty(t, contrib1.NodeID)
	assert.Empty(t, contrib1.MiddleInitial)
	assert.Empty(t, contrib1.ORCID)

	contrib2 := exported[1]
	assert.Equal(t, "Ajay", contrib2.FirstName)
	assert.Equal(t, "Carstensen", contrib2.LastName)
	assert.Empty(t, contrib2.Degree)
	assert.Equal(t, "N:user:3478dd52-e229-4e95-ab23-c6bd1e3d4d25", contrib2.NodeID)
	assert.Empty(t, contrib2.MiddleInitial)
	assert.Equal(t, "https://orcid.org/a1482559-3881-4466-b98f-d4240d9d467d", contrib2.ORCID)

	contrib3 := exported[2]
	assert.Equal(t, "Yordanka", contrib3.FirstName)
	assert.Equal(t, "Vukoja", contrib3.LastName)
	assert.Empty(t, contrib3.Degree)
	assert.Equal(t, "N:user:09b5d733-4c9e-4d3c-9842-1d4aead7f0b7", contrib3.NodeID)
	assert.Equal(t, "Ó", contrib3.MiddleInitial)
	assert.Empty(t, contrib3.ORCID)
}
