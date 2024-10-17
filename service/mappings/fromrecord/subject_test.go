package fromrecord

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
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
	assert.Len(t, subjects, 3)

	subject1 := subjects[0]
	metadatatest.AssertPennsieveInstanceIDEqual(t, "b5ad14ab-f9e7-480b-b929-8e56db504181", subject1.PennsieveID)
	metadatatest.AssertExternalInstanceIDEqual(t, "f1027e6e-17cf-45d7-8b57-4c91bfd93fce", subject1.ID)
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
	metadatatest.AssertPennsieveInstanceIDEqual(t, "c823ae0b-0c83-48be-9b0e-16690f6e675e", subject2.PennsieveID)
	metadatatest.AssertExternalInstanceIDEqual(t, "a9ea0803-651b-4f0e-bda5-e2430e235e94", subject2.ID)
	assert.Equal(t, "canis lupus familiaris", subject2.Species)
	assert.Empty(t, subject2.SpeciesSynonyms)

	subject3 := subjects[2]
	metadatatest.AssertPennsieveInstanceIDEqual(t, "a6725f6b-4504-490f-90bc-f21765d0cb07", subject3.GetPennsieveID())
	metadatatest.AssertExternalInstanceIDEqual(t, "dog-123", subject3.ExternalID())
	assert.Equal(t, "dog", subject3.Species)
	assert.Equal(t, []string{"dog",
		"pooch"}, subject3.SpeciesSynonyms)

}
