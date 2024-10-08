package fromrecord

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToSubject(t *testing.T) {

	inputDirectory := "testdata/input"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)

	subjects, err := MapRecords(reader, metadata.SubjectModelName, ToSubject)
	require.NoError(t, err)
	assert.Len(t, subjects, 2)

	subject1 := subjects[0]
	assert.Equal(t, changesetmodels.PennsieveInstanceID("b5ad14ab-f9e7-480b-b929-8e56db504181"), subject1.PennsieveID)
	assert.Equal(t, "f1027e6e-17cf-45d7-8b57-4c91bfd93fce", subject1.ID)
	assert.Equal(t, "Rattus norvegicus", subject1.Species)
	assert.Equal(t, []string{
		"brown rat",
		"Norway rat",
		"rat",
		"Mus norvegicus",
		"rats",
	},
		subject1.SpeciesSynonyms,
	)

	subject2 := subjects[1]
	assert.Equal(t, changesetmodels.PennsieveInstanceID("c823ae0b-0c83-48be-9b0e-16690f6e675e"), subject2.PennsieveID)
	assert.Equal(t, "a9ea0803-651b-4f0e-bda5-e2430e235e94", subject2.ID)
	assert.Equal(t, "canis lupus familiaris", subject2.Species)
	assert.Empty(t, subject2.SpeciesSynonyms)

}
