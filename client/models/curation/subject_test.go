package curation

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSubject_GetSpecies(t *testing.T) {
	for testName, testParams := range map[string]struct {
		subject         string
		expectedSpecies any
	}{
		"missing":            {`{"subject_id": "subject-1"}`, ""},
		"empty":              {`{"subject_id": "subject-1", "species": ""}`, ""},
		"string species":     {`{"subject_id": "subject-2", "species": "dog"}`, "dog"},
		"identifier species": {`{"subject_id": "subject-3", "species": {"id": "species-id", "label": "dog", "synonyms": ["canine"], "system": "my-ontology", "type": "identifier"}}`, SpeciesIdentifier(newDescriptionlessIdentifier("species-id", "dog", "my-ontology", "canine"))},
		"unexpected format":  {`{"subject_id": "subject-4", "species": {"A": "not expected", "B": 5}}`, SpeciesIdentifier{}},
	} {

		t.Run(testName, func(t *testing.T) {
			var subject Subject
			require.NoError(t, json.Unmarshal([]byte(testParams.subject), &subject))

			actual, err := subject.GetSpecies()
			require.NoError(t, err)
			assert.Equal(t, testParams.expectedSpecies, actual)
		})
	}
}
