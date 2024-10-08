package fromcuration

import (
	"github.com/pennsieve/ttl-sync-processor/service/mappings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToSubject(t *testing.T) {
	curationExportPath := "testdata/curation-export.json"

	datasetExport, err := UnmarshalDatasetExport(curationExportPath)
	require.NoError(t, err)

	exported, err := mappings.MapSlice(datasetExport.Subjects, ToSubject)
	require.NoError(t, err)
	assert.NotNil(t, exported)
	assert.Len(t, exported, 4)

	subject1 := exported[0]
	assert.Equal(t, "f1027e6e-17cf-45d7-8b57-4c91bfd93fce", subject1.ID)
	assert.Equal(t, "Rattus norvegicus", subject1.Species)
	assert.Equal(t, []string{
		"brown rat",
		"Norway rat",
		"rat",
		"Mus norvegicus",
		"rats",
	}, subject1.SpeciesSynonyms)

	subject2 := exported[1]
	assert.Equal(t, "a9ea0803-651b-4f0e-bda5-e2430e235e94", subject2.ID)
	assert.Equal(t, "canis lupus familiaris", subject2.Species)
	assert.Empty(t, subject2.SpeciesSynonyms)

	subject3 := exported[2]
	assert.Equal(t, "9b942959-2ebb-4680-b3f7-6d1cdbbe706d", subject3.ID)
	assert.Equal(t, "Canis familiaris", subject3.Species)
	assert.Empty(t, subject3.SpeciesSynonyms)

	subject4 := exported[3]
	assert.Equal(t, "5061b8e3-086b-4e87-b2f0-08f1a5b96679", subject4.ID)
	assert.Equal(t, "Canis familiaris", subject4.Species)
	assert.Empty(t, subject4.SpeciesSynonyms)

}
