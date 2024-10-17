package fromrecord

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToContributor(t *testing.T) {
	inputDirectory := "testdata/input"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)

	existingContributors, err := MapRecords(reader, metadata.ContributorModelName, ToContributor)
	require.NoError(t, err)
	assert.Len(t, existingContributors, 5)

	contrib1 := existingContributors[0]
	assert.Equal(t, "Elara", contrib1.FirstName)
	assert.Equal(t, "Lauridsen", contrib1.LastName)
	assert.Empty(t, contrib1.Degree)
	assert.Empty(t, contrib1.NodeID)
	assert.Empty(t, contrib1.MiddleInitial)
	assert.Empty(t, contrib1.ORCID)

	contrib2 := existingContributors[1]
	assert.Equal(t, "Yordanka", contrib2.FirstName)
	assert.Equal(t, "Vukoja", contrib2.LastName)
	assert.Equal(t, "PHD", contrib2.Degree)
	assert.Empty(t, contrib2.NodeID)
	assert.Equal(t, "T", contrib2.MiddleInitial)
	assert.Empty(t, contrib2.ORCID)

	contrib5 := existingContributors[4]
	assert.Equal(t, "Ajay", contrib5.FirstName)
	assert.Equal(t, "Carstensen", contrib5.LastName)
	assert.Empty(t, contrib5.Degree)
	assert.Equal(t, "N:user:3478dd52-e229-4e95-ab23-c6bd1e3d4d25", contrib5.NodeID)
	assert.Empty(t, contrib5.MiddleInitial)
	assert.Equal(t, "https://orcid.org/a1482559-3881-4466-b98f-d4240d9d467d", contrib5.ORCID)
}
