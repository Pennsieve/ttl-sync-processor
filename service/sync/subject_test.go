package sync

import (
	"github.com/google/uuid"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
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
	changes, err := ComputeSubjectChanges(emptySchema, []metadata.SavedSubject{}, []metadata.Subject{})
	require.NoError(t, err)
	// Nil changes means no updates required.
	require.Nil(t, changes)
}

func subjectModelDoesNotExist(t *testing.T) {
	newSubject := metadatatest.NewSubjectBuilder().Build()
	newSubject2 := metadatatest.NewSubjectBuilder().WithSpeciesSynonyms(2).Build()
	changes, err := ComputeSubjectChanges(emptySchema,
		[]metadata.SavedSubject{},
		[]metadata.Subject{newSubject, newSubject2})
	require.NoError(t, err)
	require.NotNil(t, changes)

	var modelCreate *changesetmodels.ModelCreate
	require.IsType(t, modelCreate, changes)
	modelCreate = changes.(*changesetmodels.ModelCreate)

	assert.NotNil(t, modelCreate.Create)
	assert.Equal(t, metadata.SubjectModelName, modelCreate.Create.Model.Name)
	assert.Len(t, modelCreate.Create.Properties, 3)

	assert.NotNil(t, modelCreate.Records)

	assert.Len(t, modelCreate.Records, 2)
	// The Create for newSubject
	{
		assert.Equal(t, newSubject.ExternalID(), modelCreate.Records[0].ExternalID)
		values := modelCreate.Records[0].Values
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
		assert.Equal(t, newSubject2.ExternalID(), modelCreate.Records[1].ExternalID)
		values := modelCreate.Records[1].Values
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
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).Build())

	newSubject := metadatatest.NewSubjectBuilder().Build()
	newSubject2 := metadatatest.NewSubjectBuilder().WithSpeciesSynonyms(2).Build()

	changes, err := ComputeSubjectChanges(schemaData,
		[]metadata.SavedSubject{},
		[]metadata.Subject{newSubject, newSubject2})
	require.NoError(t, err)
	require.NotNil(t, changes)

	var modelUpdate *changesetmodels.ModelUpdate
	require.IsType(t, modelUpdate, changes)
	modelUpdate = changes.(*changesetmodels.ModelUpdate)

	expectedModel, _ := schemaData.ModelByName(metadata.SubjectModelName)
	assert.Equal(t, expectedModel.ID, modelUpdate.ID.String())

	assert.NotNil(t, modelUpdate.Records)
	assert.Empty(t, modelUpdate.Records.Update)
	assert.Empty(t, modelUpdate.Records.Delete)

	assert.Len(t, modelUpdate.Records.Create, 2)
	// The Create for newSubject
	{
		assert.Equal(t, newSubject.ExternalID(), modelUpdate.Records.Create[0].ExternalID)
		values := modelUpdate.Records.Create[0].Values
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
		assert.Equal(t, newSubject2.ExternalID(), modelUpdate.Records.Create[1].ExternalID)
		values := modelUpdate.Records.Create[1].Values
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
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).Build())

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
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).Build())

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
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).Build())

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

	var modelUpdate *changesetmodels.ModelUpdate
	require.IsType(t, modelUpdate, changes)
	modelUpdate = changes.(*changesetmodels.ModelUpdate)

	expectedModel, _ := schemaData.ModelByName(metadata.SubjectModelName)
	assert.Equal(t, expectedModel.ID, modelUpdate.ID.String())

	assert.NotNil(t, modelUpdate.Records)
	assert.Empty(t, modelUpdate.Records.Create)
	assert.Empty(t, modelUpdate.Records.Delete)

	assert.Len(t, modelUpdate.Records.Update, 2)
	// The Update for originalSubject
	{
		values := findRecordUpdateByPennsieveID(t, modelUpdate.Records.Update, originalSubjectSaved.PennsieveID).Values
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
		values := findRecordUpdateByPennsieveID(t, modelUpdate.Records.Update, originalSubject2Saved.PennsieveID).Values
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
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).Build())

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

	var modelUpdate *changesetmodels.ModelUpdate
	require.IsType(t, modelUpdate, changes)
	modelUpdate = changes.(*changesetmodels.ModelUpdate)

	expectedModel, _ := schemaData.ModelByName(metadata.SubjectModelName)
	assert.Equal(t, expectedModel.ID, modelUpdate.ID.String())

	assert.NotNil(t, modelUpdate.Records)
	assert.Empty(t, modelUpdate.Records.Create)
	assert.Empty(t, modelUpdate.Records.Update)

	assert.Len(t, modelUpdate.Records.Delete, 1)
	assert.Contains(t, modelUpdate.Records.Delete, deletedSubjectSaved.PennsieveID)
}
