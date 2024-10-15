package sync

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
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
			test(t)
		})
	}
}

func emptyChangesSampleSubject(t *testing.T) {
	changes, err := ComputeSampleSubjectChanges(emptySchema, []metadata.SavedSampleSubject{}, []metadata.SampleSubject{})
	require.NoError(t, err)
	assert.Nil(t, changes)
}
func sampleSubjectLinkSchemaDoesNotExist(t *testing.T) {
	link1 := metadatatest.NewSampleSubject()
	link2 := metadatatest.NewSampleSubject()

	changes, err := ComputeSampleSubjectChanges(emptySchema, []metadata.SavedSampleSubject{}, []metadata.SampleSubject{link1, link2})
	require.NoError(t, err)

	// no link schema, so no ID and we need to supply the correct info to create schema
	assert.Empty(t, changes.ID)
	assert.NotNil(t, changes.Create)
	assert.Equal(t, metadata.SampleSubjectLinkName, changes.Create.Name)
	assert.Equal(t, metadata.SampleSubjectLinkDisplayName, changes.Create.DisplayName)
	assert.Equal(t, metadata.SampleModelName, changes.Create.FromModelName)
	assert.Equal(t, metadata.SubjectModelName, changes.Create.ToModelName)
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

	expectedLinkSchema, _ := schemaData.LinkedPropertyByName(metadata.SampleSubjectLinkName)
	assert.Equal(t, expectedLinkSchema.ID, changes.ID)
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

	changes, err := ComputeSampleSubjectChanges(schemaData, []metadata.SavedSampleSubject{savedLink1, savedLink2}, []metadata.SampleSubject{link1, link2})
	require.NoError(t, err)
	assert.Nil(t, changes)
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

	changes, err := ComputeSampleSubjectChanges(schemaData, []metadata.SavedSampleSubject{savedLink2, savedLink1}, []metadata.SampleSubject{link1, link2})
	require.NoError(t, err)
	assert.Nil(t, changes)
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

	expectedLinkSchema, _ := schemaData.LinkedPropertyByName(metadata.SampleSubjectLinkName)
	assert.Equal(t, expectedLinkSchema.ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.Empty(t, changes.Instances.Create)
	assert.Len(t, changes.Instances.Delete, 1)
}

func sampleSubjectChangeSampleID(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().
		WithModel(metadata.SampleModelName, metadata.SampleDisplayName).
		WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).
		WithLinkedProperty(metadata.SampleSubjectLinkName, metadata.SampleSubjectLinkDisplayName).Build(),
	)

	oldLink := metadatatest.NewSavedSampleSubject(metadatatest.NewSampleSubject())
	newLink := metadatatest.NewSampleSubjectBuilder().WithSubjectID(oldLink.SubjectID).Build()

	changes, err := ComputeSampleSubjectChanges(schemaData, []metadata.SavedSampleSubject{oldLink}, []metadata.SampleSubject{newLink})
	require.NoError(t, err)

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
}

func sampleSubjectChangeSubjectID(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().
		WithModel(metadata.SampleModelName, metadata.SampleDisplayName).
		WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).
		WithLinkedProperty(metadata.SampleSubjectLinkName, metadata.SampleSubjectLinkDisplayName).Build(),
	)

	oldLink := metadatatest.NewSavedSampleSubject(metadatatest.NewSampleSubject())
	newLink := metadatatest.NewSampleSubjectBuilder().WithSampleID(oldLink.SampleID).Build()

	changes, err := ComputeSampleSubjectChanges(schemaData, []metadata.SavedSampleSubject{oldLink}, []metadata.SampleSubject{newLink})
	require.NoError(t, err)

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
}
