package sync

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings/fromrecord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestComputeSampleSubjectChanges(t *testing.T) {
	for scenario, test := range map[string]func(t *testing.T){
		"handle edge case without panic":                emptyChangesSampleSubject,
		"link schema does not exist":                    sampleSubjectLinkSchemaDoesNotExist,
		"link schema exists, but no existing instances": sampleSubjectSchemaExistsButNoInstances,
		"no changes":                  noSampleSubjectChanges,
		"order does not matter":       sampleSubjectOrderDoesNotMatter,
		"deleted sample subject link": sampleSubjectDeleted,
		"change sample id (from)":     sampleSubjectChangeSampleID,
		"change subject id (to)":      sampleSubjectChangeSubjectID,
	} {
		t.Run(scenario, func(t *testing.T) {
			// Need to initialize these package-wide vars for each test
			// Individual test functions should populate ExistingRecordStore
			// to simulate different situations and then
			// check that recordIDMap has the expected entries.
			ExistingRecordStore = fromrecord.NewRecordIDStore()
			recordIDMap = make(fromrecord.RecordIDMap)

			test(t)
		})
	}
}

func emptyChangesSampleSubject(t *testing.T) {
	changes, err := ComputeSampleSubjectChanges(emptySchema, []metadata.SavedSampleSubject{}, []metadata.SampleSubject{})
	require.NoError(t, err)
	assert.Nil(t, changes)
	assert.Empty(t, recordIDMap)
}
func sampleSubjectLinkSchemaDoesNotExist(t *testing.T) {
	link1 := metadatatest.NewSampleSubject()
	link2 := metadatatest.NewSampleSubject()

	// This means that only the subject of link2 already existed in the
	// dataset's metadata
	link2ToPennsieveID := metadatatest.NewPennsieveInstanceID()
	ExistingRecordStore.Add(metadata.SubjectModelName, link2.ToExternalID(), link2ToPennsieveID)

	changes, err := ComputeSampleSubjectChanges(emptySchema, []metadata.SavedSampleSubject{}, []metadata.SampleSubject{link1, link2})
	require.NoError(t, err)

	// no link schema, so no ID and we need to supply the correct info to create schema
	assert.Empty(t, changes.ID)
	assert.NotNil(t, changes.Create)
	assert.Equal(t, metadata.SampleModelName, changes.FromModelName)
	assert.Equal(t, metadata.SubjectModelName, changes.ToModelName)
	assert.Equal(t, metadata.SampleSubjectLinkName, changes.Create.Name)
	assert.Equal(t, metadata.SampleSubjectLinkDisplayName, changes.Create.DisplayName)

	assert.Equal(t, 1, changes.Create.Position)

	// Need to create two link instances and no deletes
	assert.Empty(t, changes.Instances.Delete)
	assert.Len(t, changes.Instances.Create, 2)

	// First Create
	{
		create := changes.Instances.Create[0]
		assert.Equal(t, link1.FromExternalID(), create.FromExternalID)
		assert.Equal(t, link1.ToExternalID(), create.ToExternalID)
	}

	// Second Create
	{
		create := changes.Instances.Create[1]
		assert.Equal(t, link2.FromExternalID(), create.FromExternalID)
		assert.Equal(t, link2.ToExternalID(), create.ToExternalID)
	}
	assert.Len(t, recordIDMap, 1)
	assert.Equal(t, link2ToPennsieveID, recordIDMap[fromrecord.RecordIDKey{
		ModelName:        metadata.SubjectModelName,
		ExternalRecordID: link2.ToExternalID(),
	}])

}

func sampleSubjectSchemaExistsButNoInstances(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().
		WithModel(metadata.SampleModelName, metadata.SampleDisplayName).
		WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).
		WithLinkedProperty(metadata.SampleSubjectLinkName, metadata.SampleSubjectLinkDisplayName).Build(),
	)
	link1 := metadatatest.NewSampleSubject()
	link2 := metadatatest.NewSampleSubject()

	changes, err := ComputeSampleSubjectChanges(schemaData, []metadata.SavedSampleSubject{}, []metadata.SampleSubject{link1, link2})
	require.NoError(t, err)

	assert.Equal(t, metadata.SampleModelName, changes.FromModelName)
	assert.Equal(t, metadata.SubjectModelName, changes.ToModelName)

	expectedLinkSchema, _ := schemaData.LinkedPropertyByName(metadata.SampleSubjectLinkName)
	assert.Equal(t, expectedLinkSchema.ID, changes.ID.String())
	assert.Nil(t, changes.Create)

	// Need to create two link instances and no deletes
	assert.Empty(t, changes.Instances.Delete)
	assert.Len(t, changes.Instances.Create, 2)

	// First Create
	{
		create := changes.Instances.Create[0]
		assert.Equal(t, link1.FromExternalID(), create.FromExternalID)
		assert.Equal(t, link1.ToExternalID(), create.ToExternalID)
	}

	// Second Create
	{
		create := changes.Instances.Create[1]
		assert.Equal(t, link2.FromExternalID(), create.FromExternalID)
		assert.Equal(t, link2.ToExternalID(), create.ToExternalID)
	}

	assert.Empty(t, recordIDMap)
}

func noSampleSubjectChanges(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().
		WithModel(metadata.SampleModelName, metadata.SampleDisplayName).
		WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).
		WithLinkedProperty(metadata.SampleSubjectLinkName, metadata.SampleSubjectLinkDisplayName).Build(),
	)
	link1 := metadatatest.NewSampleSubject()
	link2 := metadatatest.NewSampleSubject()

	savedLink1 := metadatatest.NewSavedSampleSubject(link1)
	savedLink2 := metadatatest.NewSavedSampleSubject(link2)
	addSavedSampleSubjectRecordIDs(ExistingRecordStore, savedLink1)
	addSavedSampleSubjectRecordIDs(ExistingRecordStore, savedLink2)

	changes, err := ComputeSampleSubjectChanges(schemaData, []metadata.SavedSampleSubject{savedLink1, savedLink2}, []metadata.SampleSubject{link1, link2})
	require.NoError(t, err)
	assert.Nil(t, changes)
	assert.Empty(t, recordIDMap)

}

func sampleSubjectOrderDoesNotMatter(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().
		WithModel(metadata.SampleModelName, metadata.SampleDisplayName).
		WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).
		WithLinkedProperty(metadata.SampleSubjectLinkName, metadata.SampleSubjectLinkDisplayName).Build(),
	)
	link1 := metadatatest.NewSampleSubject()
	link2 := metadatatest.NewSampleSubject()

	savedLink1 := metadatatest.NewSavedSampleSubject(link1)
	savedLink2 := metadatatest.NewSavedSampleSubject(link2)
	addSavedSampleSubjectRecordIDs(ExistingRecordStore, savedLink1)
	addSavedSampleSubjectRecordIDs(ExistingRecordStore, savedLink2)

	changes, err := ComputeSampleSubjectChanges(schemaData, []metadata.SavedSampleSubject{savedLink2, savedLink1}, []metadata.SampleSubject{link1, link2})
	require.NoError(t, err)
	assert.Nil(t, changes)
	assert.Empty(t, recordIDMap)

}

func sampleSubjectDeleted(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().
		WithModel(metadata.SampleModelName, metadata.SampleDisplayName).
		WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).
		WithLinkedProperty(metadata.SampleSubjectLinkName, metadata.SampleSubjectLinkDisplayName).Build(),
	)
	link1 := metadatatest.NewSampleSubject()
	link2 := metadatatest.NewSampleSubject()

	savedLink1 := metadatatest.NewSavedSampleSubject(link1)
	deletedLink := metadatatest.NewSavedSampleSubject(metadatatest.NewSampleSubject())
	savedLink2 := metadatatest.NewSavedSampleSubject(link2)

	changes, err := ComputeSampleSubjectChanges(schemaData, []metadata.SavedSampleSubject{savedLink2, deletedLink, savedLink1}, []metadata.SampleSubject{link1, link2})
	require.NoError(t, err)

	assert.Equal(t, metadata.SampleModelName, changes.FromModelName)
	assert.Equal(t, metadata.SubjectModelName, changes.ToModelName)

	expectedLinkSchema, _ := schemaData.LinkedPropertyByName(metadata.SampleSubjectLinkName)
	assert.Equal(t, expectedLinkSchema.ID, changes.ID.String())
	assert.Nil(t, changes.Create)

	assert.Empty(t, changes.Instances.Create)
	assert.Len(t, changes.Instances.Delete, 1)
	assert.Empty(t, recordIDMap)

}

func sampleSubjectChangeSampleID(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().
		WithModel(metadata.SampleModelName, metadata.SampleDisplayName).
		WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).
		WithLinkedProperty(metadata.SampleSubjectLinkName, metadata.SampleSubjectLinkDisplayName).Build(),
	)

	oldLink := metadatatest.NewSavedSampleSubject(metadatatest.NewSampleSubject())
	newLink := metadatatest.NewSampleSubjectBuilder().WithSubjectID(oldLink.SubjectID).Build()
	addSavedSampleSubjectRecordIDs(ExistingRecordStore, oldLink)

	changes, err := ComputeSampleSubjectChanges(schemaData, []metadata.SavedSampleSubject{oldLink}, []metadata.SampleSubject{newLink})
	require.NoError(t, err)

	assert.Equal(t, metadata.SampleModelName, changes.FromModelName)
	assert.Equal(t, metadata.SubjectModelName, changes.ToModelName)

	assert.Nil(t, changes.Create)
	assert.NotEmpty(t, changes.ID)

	deletes := changes.Instances.Delete
	assert.Len(t, deletes, 1)
	assert.Equal(t, oldLink.FromPennsieveID(), deletes[0].FromRecordID)
	assert.Equal(t, oldLink.GetPennsieveID(), deletes[0].InstanceLinkedPropertyID)

	creates := changes.Instances.Create
	assert.Len(t, creates, 1)
	assert.Equal(t, newLink.FromExternalID(), creates[0].FromExternalID)
	assert.Equal(t, newLink.ToExternalID(), creates[0].ToExternalID)

	assert.Len(t, recordIDMap, 1)
	assert.Equal(t, oldLink.ToPennsieveID(),
		recordIDMap[fromrecord.RecordIDKey{ModelName: metadata.SubjectModelName, ExternalRecordID: oldLink.ToExternalID()}])
}

func sampleSubjectChangeSubjectID(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().
		WithModel(metadata.SampleModelName, metadata.SampleDisplayName).
		WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).
		WithLinkedProperty(metadata.SampleSubjectLinkName, metadata.SampleSubjectLinkDisplayName).Build(),
	)

	oldLink := metadatatest.NewSavedSampleSubject(metadatatest.NewSampleSubject())
	newLink := metadatatest.NewSampleSubjectBuilder().WithSampleID(oldLink.SampleID).Build()
	addSavedSampleSubjectRecordIDs(ExistingRecordStore, oldLink)

	changes, err := ComputeSampleSubjectChanges(schemaData, []metadata.SavedSampleSubject{oldLink}, []metadata.SampleSubject{newLink})
	require.NoError(t, err)

	assert.Equal(t, metadata.SampleModelName, changes.FromModelName)
	assert.Equal(t, metadata.SubjectModelName, changes.ToModelName)

	assert.Nil(t, changes.Create)
	assert.NotEmpty(t, changes.ID)

	deletes := changes.Instances.Delete
	assert.Len(t, deletes, 1)
	assert.Equal(t, oldLink.FromPennsieveID(), deletes[0].FromRecordID)
	assert.Equal(t, oldLink.GetPennsieveID(), deletes[0].InstanceLinkedPropertyID)

	creates := changes.Instances.Create
	assert.Len(t, creates, 1)
	assert.Equal(t, newLink.FromExternalID(), creates[0].FromExternalID)
	assert.Equal(t, newLink.ToExternalID(), creates[0].ToExternalID)

	assert.Len(t, recordIDMap, 1)
	assert.Equal(t, oldLink.FromPennsieveID(),
		recordIDMap[fromrecord.RecordIDKey{ModelName: metadata.SampleModelName, ExternalRecordID: oldLink.FromExternalID()}])
}

func addSavedSampleSubjectRecordIDs(existingRecordIDs *fromrecord.RecordIDStore, savedSampleSubject metadata.SavedSampleSubject) {
	existingRecordIDs.Add(metadata.SampleModelName, savedSampleSubject.FromExternalID(), savedSampleSubject.FromPennsieveID())
	existingRecordIDs.Add(metadata.SubjectModelName, savedSampleSubject.ToExternalID(), savedSampleSubject.ToPennsieveID())
}

func addSavedSampleRecordID(existingRecordIDs *fromrecord.RecordIDStore, savedSample metadata.SavedSample) {
	existingRecordIDs.Add(metadata.SampleModelName, savedSample.ExternalID(), savedSample.GetPennsieveID())
}

func addSavedSubjectRecordID(existingRecordIDs *fromrecord.RecordIDStore, savedSubject metadata.SavedSubject) {
	existingRecordIDs.Add(metadata.SubjectModelName, savedSubject.ExternalID(), savedSubject.GetPennsieveID())
}
