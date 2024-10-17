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

func TestComputeSampleChanges(t *testing.T) {
	for scenario, test := range map[string]func(tt *testing.T){
		"handle edge case without panic":        emptyChangesSample,
		"model does not exist":                  sampleModelDoesNotExist,
		"model exists, but no existing records": sampleModelExistsButNoExistingRecords,
		"no changes":                            noSampleChanges,
		"order does not matter":                 sampleOrderDoesNotMatter,
		"deleted sample":                        deleteSample,
		"update sample":                         updateSample,
	} {
		t.Run(scenario, func(t *testing.T) {
			test(t)
		})
	}
}

func emptyChangesSample(t *testing.T) {
	changes, err := ComputeSampleChanges(emptySchema, []metadata.SavedSample{}, []metadata.Sample{})
	require.NoError(t, err)
	// Nil changes means no updates required.
	require.Nil(t, changes)
}

func sampleModelDoesNotExist(t *testing.T) {
	newSample1 := metadatatest.NewSampleBuilder().Build()
	newSample2 := metadatatest.NewSampleBuilder().Build()
	changes, err := ComputeSampleChanges(emptySchema,
		[]metadata.SavedSample{},
		[]metadata.Sample{newSample1, newSample2})
	require.NoError(t, err)
	require.NotNil(t, changes)

	assert.Empty(t, changes.ID)
	assert.NotNil(t, changes.Create)
	assert.Equal(t, metadata.SampleModelName, changes.Create.Model.Name)
	assert.Len(t, changes.Create.Properties, 2)

	assert.NotNil(t, changes.Records)
	assert.False(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Update)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Create, 2)
	// The Create for newSample1
	{
		values := changes.Records.Create[0].Values
		// Only contains ID since that is the only property
		assert.Len(t, values, 2)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.SampleIDKey)
		assert.Equal(t, newSample1.ID, valuesByName[metadata.SampleIDKey].Value)

		assert.Contains(t, valuesByName, metadata.PrimaryKeyKey)
		assert.Equal(t, newSample1.PrimaryKey, valuesByName[metadata.PrimaryKeyKey].Value)

	}

	// The Create for newSample2
	{
		values := changes.Records.Create[1].Values
		// Only contains ID and species because other values are empty
		assert.Len(t, values, 2)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.SampleIDKey)
		assert.Equal(t, newSample2.ID, valuesByName[metadata.SampleIDKey].Value)

		assert.Contains(t, valuesByName, metadata.PrimaryKeyKey)
		assert.Equal(t, newSample2.PrimaryKey, valuesByName[metadata.PrimaryKeyKey].Value)

	}

}

func sampleModelExistsButNoExistingRecords(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.SampleModelName, metadata.SampleDisplayName).Build())

	newSample1 := metadatatest.NewSampleBuilder().Build()
	newSample2 := metadatatest.NewSampleBuilder().Build()

	changes, err := ComputeSampleChanges(schemaData,
		[]metadata.SavedSample{},
		[]metadata.Sample{newSample1, newSample2})
	require.NoError(t, err)
	require.NotNil(t, changes)

	expectedModel, _ := schemaData.ModelByName(metadata.SampleModelName)
	assert.Equal(t, expectedModel.ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.NotNil(t, changes.Records)
	assert.False(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Update)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Create, 2)
	// The Create for newSample1
	{
		values := changes.Records.Create[0].Values
		// Only contains ID
		assert.Len(t, values, 2)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.SampleIDKey)
		assert.Equal(t, newSample1.ID, valuesByName[metadata.SampleIDKey].Value)

		assert.Contains(t, valuesByName, metadata.PrimaryKeyKey)
		assert.Equal(t, newSample1.PrimaryKey, valuesByName[metadata.PrimaryKeyKey].Value)

	}

	// The Create for newSample2
	{
		values := changes.Records.Create[1].Values
		// Only contains ID
		assert.Len(t, values, 2)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.SampleIDKey)
		assert.Equal(t, newSample2.ID, valuesByName[metadata.SampleIDKey].Value)

		assert.Contains(t, valuesByName, metadata.PrimaryKeyKey)
		assert.Equal(t, newSample2.PrimaryKey, valuesByName[metadata.PrimaryKeyKey].Value)

	}
}

func noSampleChanges(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.SampleModelName, metadata.SampleDisplayName).Build())

	sample1 := metadatatest.NewSampleBuilder().Build()
	sample2 := metadatatest.NewSampleBuilder().Build()

	savedSample1 := metadatatest.NewSavedSample(sample1)
	savedSample2 := metadatatest.NewSavedSample(sample2)

	changes, err := ComputeSampleChanges(
		schemaData,
		[]metadata.SavedSample{savedSample1, savedSample2},
		[]metadata.Sample{sample1, sample2},
	)
	require.NoError(t, err)
	assert.Nil(t, changes)
}

func sampleOrderDoesNotMatter(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.SampleModelName, metadata.SampleDisplayName).Build())

	sample1 := metadatatest.NewSampleBuilder().Build()
	sample2 := metadatatest.NewSampleBuilder().Build()

	savedSample1 := metadatatest.NewSavedSample(sample1)
	savedSample2 := metadatatest.NewSavedSample(sample2)

	changes, err := ComputeSampleChanges(
		schemaData,
		[]metadata.SavedSample{savedSample2, savedSample1},
		[]metadata.Sample{sample1, sample2},
	)
	require.NoError(t, err)
	assert.Nil(t, changes)
}

func deleteSample(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.SampleModelName, metadata.SampleDisplayName).Build())

	keptSample1 := metadatatest.NewSampleBuilder().Build()
	keptSample2 := metadatatest.NewSampleBuilder().Build()

	keptSample1Saved := metadatatest.NewSavedSample(keptSample1)
	deletedSampleSaved := metadatatest.NewSavedSample(metadatatest.NewSampleBuilder().Build())
	keptSample2Saved := metadatatest.NewSavedSample(keptSample2)

	changes, err := ComputeSampleChanges(schemaData,
		[]metadata.SavedSample{keptSample1Saved, deletedSampleSaved, keptSample2Saved},
		[]metadata.Sample{keptSample2, keptSample1})
	require.NoError(t, err)
	require.NotNil(t, changes)

	expectedModel, _ := schemaData.ModelByName(metadata.SampleModelName)
	assert.Equal(t, expectedModel.ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.NotNil(t, changes.Records)
	assert.False(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Create)
	assert.Empty(t, changes.Records.Update)

	assert.Len(t, changes.Records.Delete, 1)
	assert.Contains(t, changes.Records.Delete, deletedSampleSaved.PennsieveID)
}

func updateSample(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.SampleModelName, metadata.SampleDisplayName).Build())

	originalSample := metadatatest.NewSampleBuilder().Build()
	originalSample2 := metadatatest.NewSampleBuilder().Build()
	unchangedSample := metadatatest.NewSampleBuilder().Build()

	originalSampleSaved := metadatatest.NewSavedSample(originalSample)
	originalSample2Saved := metadatatest.NewSavedSample(originalSample2)
	unchangedSampleSaved := metadatatest.NewSavedSample(unchangedSample)

	updatedSample := metadatatest.SampleCopy(originalSample)
	updatedSample.PrimaryKey = uuid.NewString()

	updatedSample2 := metadatatest.SampleCopy(originalSample2)
	updatedSample2.PrimaryKey = uuid.NewString()

	changes, err := ComputeSampleChanges(schemaData,
		[]metadata.SavedSample{originalSampleSaved, originalSample2Saved, unchangedSampleSaved},
		[]metadata.Sample{unchangedSample, updatedSample2, updatedSample})
	require.NoError(t, err)
	require.NotNil(t, changes)

	expectedModel, _ := schemaData.ModelByName(metadata.SampleModelName)
	assert.Equal(t, expectedModel.ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.NotNil(t, changes.Records)
	assert.False(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Create)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Update, 2)
	// The Update for originalSample
	{
		values := findRecordUpdateByPennsieveID(t, changes.Records.Update, originalSampleSaved.PennsieveID).Values
		assert.Len(t, values, 2)

		// ID not updated
		id := findValueByName(t, values, metadata.SampleIDKey)
		assert.Equal(t, originalSample.ID, id.Value)

		// PrimaryKey updated
		species := findValueByName(t, values, metadata.PrimaryKeyKey)
		assert.Equal(t, updatedSample.PrimaryKey, species.Value)

	}

	// The Update for originalSample2
	{
		values := findRecordUpdateByPennsieveID(t, changes.Records.Update, originalSample2Saved.PennsieveID).Values
		assert.Len(t, values, 2)

		// ID not updated
		id := findValueByName(t, values, metadata.SampleIDKey)
		assert.Equal(t, originalSample2.ID, id.Value)

		// PrimaryKey updated
		synonyms := findValueByName(t, values, metadata.PrimaryKeyKey)
		assert.Equal(t, updatedSample2.PrimaryKey, synonyms.Value)
	}
}
