package sync

import (
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestComputeSubjectChanges(t *testing.T) {
	for scenario, test := range map[string]func(tt *testing.T){
		"handle edge cases without panic": emptyChangesetSubject,
		"model does not exist":            subjectModelDoesNotExist,
	} {
		t.Run(scenario, func(t *testing.T) {
			test(t)
		})
	}
}

func emptyChangesetSubject(t *testing.T) {
	changes, err := ComputeSubjectChanges(map[string]schema.Element{}, &metadata.SavedDatasetMetadata{}, &metadata.DatasetMetadata{})
	require.NoError(t, err)
	// Nil changes means no updates required.
	require.Nil(t, changes)
}

func subjectModelDoesNotExist(t *testing.T) {
	newSubject := metadatatest.NewSubjectBuilder().Build()
	changes, err := ComputeSubjectChanges(map[string]schema.Element{},
		&metadata.SavedDatasetMetadata{},
		&metadata.DatasetMetadata{Subjects: []metadata.Subject{newSubject}})
	require.NoError(t, err)
	require.NotNil(t, changes)

	assert.Empty(t, changes.ID)
	assert.NotNil(t, changes.Create)
	assert.Equal(t, metadata.SubjectModelName, changes.Create.Model.Name)
	assert.Len(t, changes.Create.Properties, 3)

	assert.NotNil(t, changes.Records)
	assert.False(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Update)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Create, 1)
	values := changes.Records.Create[0].Values
	// Only contains ID and species because other values are empty
	assert.Len(t, values, 2)
	valuesByName := map[string]changesetmodels.RecordValue{}
	for _, v := range values {
		valuesByName[v.Name] = v
	}
	assert.Contains(t, valuesByName, metadata.SubjectIDKey)
	assert.Equal(t, newSubject.ID, valuesByName[metadata.SubjectIDKey].Value)

	assert.Contains(t, valuesByName, metadata.SpeciesKey)
	assert.Equal(t, newSubject.Species, valuesByName[metadata.SpeciesKey].Value)

	assert.NotContains(t, valuesByName, metadata.SpeciesSynonymsKey)

}
