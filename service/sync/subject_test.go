package sync

import (
	"github.com/google/uuid"
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
		"handle edge cases without panic":       emptyChangesetSubject,
		"model does not exist":                  subjectModelDoesNotExist,
		"model exists, but no existing records": subjectModelExistsButNoExistingRecords,
		"no changes":                            noSubjectChanges,
		"order does not matter":                 subjectOrderDoesNotMatter,
		"updated subject":                       updateSubject,
		"deleted subject":                       deleteSubject,
	} {
		t.Run(scenario, func(t *testing.T) {
			test(t)
		})
	}
}

func emptyChangesetSubject(t *testing.T) {
	changes, err := ComputeSubjectChanges(map[string]schema.Element{}, []metadata.SavedSubject{}, []metadata.Subject{})
	require.NoError(t, err)
	// Nil changes means no updates required.
	require.Nil(t, changes)
}

func subjectModelDoesNotExist(t *testing.T) {
	newSubject := metadatatest.NewSubjectBuilder().Build()
	newSubject2 := metadatatest.NewSubjectBuilder().WithSpeciesSynonyms(2).Build()
	changes, err := ComputeSubjectChanges(map[string]schema.Element{},
		[]metadata.SavedSubject{},
		[]metadata.Subject{newSubject, newSubject2})
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

	assert.Len(t, changes.Records.Create, 2)
	// The Create for newSubject
	{
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

	// The Create for newSubject2
	{
		values := changes.Records.Create[1].Values
		// Only contains ID and species because other values are empty
		assert.Len(t, values, 3)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.SubjectIDKey)
		assert.Equal(t, newSubject2.ID, valuesByName[metadata.SubjectIDKey].Value)

		assert.Contains(t, valuesByName, metadata.SpeciesKey)
		assert.Equal(t, newSubject2.Species, valuesByName[metadata.SpeciesKey].Value)

		assert.Contains(t, valuesByName, metadata.SpeciesSynonymsKey)
		assert.Equal(t, newSubject2.SpeciesSynonyms, valuesByName[metadata.SpeciesSynonymsKey].Value)
	}

}

func subjectModelExistsButNoExistingRecords(t *testing.T) {
	schemaData := newTestSchemaData().WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName)

	newSubject := metadatatest.NewSubjectBuilder().Build()
	newSubject2 := metadatatest.NewSubjectBuilder().WithSpeciesSynonyms(2).Build()

	changes, err := ComputeSubjectChanges(schemaData,
		[]metadata.SavedSubject{},
		[]metadata.Subject{newSubject, newSubject2})
	require.NoError(t, err)
	require.NotNil(t, changes)

	assert.Equal(t, schemaData[metadata.SubjectModelName].ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.NotNil(t, changes.Records)
	assert.False(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Update)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Create, 2)
	// The Create for newSubject
	{
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

	// The Create for newSubject2
	{
		values := changes.Records.Create[1].Values
		// Only contains ID and species because other values are empty
		assert.Len(t, values, 3)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.SubjectIDKey)
		assert.Equal(t, newSubject2.ID, valuesByName[metadata.SubjectIDKey].Value)

		assert.Contains(t, valuesByName, metadata.SpeciesKey)
		assert.Equal(t, newSubject2.Species, valuesByName[metadata.SpeciesKey].Value)

		assert.Contains(t, valuesByName, metadata.SpeciesSynonymsKey)
		assert.Equal(t, newSubject2.SpeciesSynonyms, valuesByName[metadata.SpeciesSynonymsKey].Value)
	}
}

func noSubjectChanges(t *testing.T) {
	schemaData := newTestSchemaData().WithModel(metadata.ContributorModelName, metadata.ContributorDisplayName)

	subject1 := metadatatest.NewSubjectBuilder().Build()
	subject2 := metadatatest.NewSubjectBuilder().WithSpeciesSynonyms(3).Build()

	savedSubject1 := metadatatest.NewSavedSubject(subject1)
	savedSubject2 := metadatatest.NewSavedSubject(subject2)

	changes, err := ComputeSubjectChanges(
		schemaData,
		[]metadata.SavedSubject{savedSubject1, savedSubject2},
		[]metadata.Subject{subject1, subject2},
	)
	require.NoError(t, err)
	assert.Nil(t, changes)
}

func subjectOrderDoesNotMatter(t *testing.T) {
	schemaData := newTestSchemaData().WithModel(metadata.ContributorModelName, metadata.ContributorDisplayName)

	subject1 := metadatatest.NewSubjectBuilder().Build()
	subject2 := metadatatest.NewSubjectBuilder().WithSpeciesSynonyms(3).Build()

	savedSubject1 := metadatatest.NewSavedSubject(subject1)
	savedSubject2 := metadatatest.NewSavedSubject(subject2)

	changes, err := ComputeSubjectChanges(
		schemaData,
		[]metadata.SavedSubject{savedSubject1, savedSubject2},
		[]metadata.Subject{subject2, subject1},
	)
	require.NoError(t, err)
	assert.Nil(t, changes)
}

func updateSubject(t *testing.T) {
	schemaData := newTestSchemaData().WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName)

	originalSubject := metadatatest.NewSubjectBuilder().Build()
	originalSubject2 := metadatatest.NewSubjectBuilder().WithSpeciesSynonyms(2).Build()
	unchangedSubject := metadatatest.NewSubjectBuilder().WithSpeciesSynonyms(3).Build()

	originalSubjectSaved := metadatatest.NewSavedSubject(originalSubject)
	originalSubject2Saved := metadatatest.NewSavedSubject(originalSubject2)
	unchangedSubjectSaved := metadatatest.NewSavedSubject(unchangedSubject)

	updatedSubject := metadatatest.SubjectCopy(originalSubject)
	updatedSubject.Species = uuid.NewString()

	updatedSubject2 := metadatatest.SubjectCopy(originalSubject2)
	updatedSubject2.SpeciesSynonyms[0] = uuid.NewString()

	changes, err := ComputeSubjectChanges(schemaData,
		[]metadata.SavedSubject{originalSubjectSaved, originalSubject2Saved, unchangedSubjectSaved},
		[]metadata.Subject{unchangedSubject, updatedSubject2, updatedSubject})
	require.NoError(t, err)
	require.NotNil(t, changes)

	assert.Equal(t, schemaData[metadata.SubjectModelName].ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.NotNil(t, changes.Records)
	assert.False(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Create)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Update, 2)
	// The Update for originalSubject
	{
		values := findRecordUpdateByPennsieveID(t, changes.Records.Update, originalSubjectSaved.PennsieveID).Values
		// Only contains ID and species because other values are empty
		assert.Len(t, values, 3)

		// ID not updated
		id := findValueByName(t, values, metadata.SubjectIDKey)
		assert.Equal(t, originalSubject.ID, id.Value)

		// Species updated
		species := findValueByName(t, values, metadata.SpeciesKey)
		assert.Equal(t, updatedSubject.Species, species.Value)

		// Synonyms not updated
		synonyms := findValueByName(t, values, metadata.SpeciesSynonymsKey)
		assert.Equal(t, originalSubject.SpeciesSynonyms, synonyms.Value)
	}

	// The Update for originalSubject2
	{
		values := findRecordUpdateByPennsieveID(t, changes.Records.Update, originalSubject2Saved.PennsieveID).Values
		// Only contains ID and species because other values are empty
		assert.Len(t, values, 3)

		// ID not updated
		id := findValueByName(t, values, metadata.SubjectIDKey)
		assert.Equal(t, originalSubject2.ID, id.Value)

		// Species not updated
		species := findValueByName(t, values, metadata.SpeciesKey)
		assert.Equal(t, originalSubject2.Species, species.Value)

		// Synonyms updated
		synonyms := findValueByName(t, values, metadata.SpeciesSynonymsKey)
		assert.Equal(t, updatedSubject2.SpeciesSynonyms, synonyms.Value)
	}
}

func deleteSubject(t *testing.T) {
	schemaData := newTestSchemaData().WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName)

	keptSubject1 := metadatatest.NewSubjectBuilder().Build()
	deletedSubject := metadatatest.NewSubjectBuilder().WithSpeciesSynonyms(2).Build()
	keptSubject2 := metadatatest.NewSubjectBuilder().WithSpeciesSynonyms(3).Build()

	keptSubject1Saved := metadatatest.NewSavedSubject(keptSubject1)
	deletedSubjectSaved := metadatatest.NewSavedSubject(deletedSubject)
	keptSubject2Saved := metadatatest.NewSavedSubject(keptSubject2)

	changes, err := ComputeSubjectChanges(schemaData,
		[]metadata.SavedSubject{keptSubject1Saved, deletedSubjectSaved, keptSubject2Saved},
		[]metadata.Subject{keptSubject2, keptSubject1})
	require.NoError(t, err)
	require.NotNil(t, changes)

	assert.Equal(t, schemaData[metadata.SubjectModelName].ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.NotNil(t, changes.Records)
	assert.False(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Create)
	assert.Empty(t, changes.Records.Update)

	assert.Len(t, changes.Records.Delete, 1)
	assert.Contains(t, changes.Records.Delete, deletedSubjectSaved.PennsieveID)
}
