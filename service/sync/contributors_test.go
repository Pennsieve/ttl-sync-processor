package sync

import (
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestComputeContributorsChanges(t *testing.T) {
	for scenario, test := range map[string]func(tt *testing.T){
		"handle edge cases without panic":       everythingEmpty,
		"no existing model":                     modelDoesNotExist,
		"model exists, but no existing records": modelExistsButNoExistingRecords,
		"no changes":                            noChanges,
		"order changes":                         orderChange,
		"remove contributor":                    removeContributor,
		"add contributor":                       addContributor,
	} {
		t.Run(scenario, func(t *testing.T) {
			test(t)
		})
	}
}

func everythingEmpty(t *testing.T) {
	changes, err := ComputeContributorsChanges(emptySchema, []metadata.Contributor{}, []metadata.Contributor{})
	require.NoError(t, err)
	// Nil changes means no updates required.
	require.Nil(t, changes)
}

func modelDoesNotExist(t *testing.T) {
	newContrib := metadatatest.NewContributorBuilder().WithMiddleInitial().Build()
	changes, err := ComputeContributorsChanges(emptySchema, []metadata.Contributor{}, []metadata.Contributor{newContrib})
	require.NoError(t, err)
	require.NotNil(t, changes)

	assert.Empty(t, changes.ID)
	assert.NotNil(t, changes.Create)
	assert.Equal(t, metadata.ContributorModelName, changes.Create.Model.Name)
	assert.Len(t, changes.Create.Properties, 6)

	assert.NotNil(t, changes.Records)
	assert.True(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Update)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Create, 1)
	values := changes.Records.Create[0].Values
	// Only contains first and last names and middle initial because other values are empty
	assert.Len(t, values, 3)
	valuesByName := map[string]changesetmodels.RecordValue{}
	for _, v := range values {
		valuesByName[v.Name] = v
	}
	assert.Contains(t, valuesByName, metadata.LastNameKey)
	assert.Equal(t, newContrib.LastName, valuesByName[metadata.LastNameKey].Value)

	assert.Contains(t, valuesByName, metadata.MiddleInitialKey)
	assert.Equal(t, newContrib.MiddleInitial, valuesByName[metadata.MiddleInitialKey].Value)

	assert.NotContains(t, valuesByName, metadata.NodeIDKey)

}

func modelExistsButNoExistingRecords(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.ContributorModelName, metadata.ContributorDisplayName).Build())
	newContrib := metadatatest.NewContributorBuilder().WithMiddleInitial().Build()
	newContrib2 := metadatatest.NewContributorBuilder().WithNodeID().Build()

	changes, err := ComputeContributorsChanges(schemaData, []metadata.Contributor{}, []metadata.Contributor{newContrib, newContrib2})
	require.NoError(t, err)
	require.NotNil(t, changes)

	expectedModel, _ := schemaData.ModelByName(metadata.ContributorModelName)
	assert.Equal(t, expectedModel.ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.NotNil(t, changes.Records)
	assert.True(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Update)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Create, 2)
	// First Contributor
	{
		values := changes.Records.Create[0].Values
		// Only contains first and last names and middle initial because other values are empty
		assert.Len(t, values, 3)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.LastNameKey)
		assert.Equal(t, newContrib.LastName, valuesByName[metadata.LastNameKey].Value)

		assert.Contains(t, valuesByName, metadata.MiddleInitialKey)
		assert.Equal(t, newContrib.MiddleInitial, valuesByName[metadata.MiddleInitialKey].Value)

		assert.NotContains(t, valuesByName, metadata.NodeIDKey)
	}

	//Second Contributor
	{
		values := changes.Records.Create[1].Values
		// Only contains first and last names and node id because other values are empty
		assert.Len(t, values, 3)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.LastNameKey)
		assert.Equal(t, newContrib2.LastName, valuesByName[metadata.LastNameKey].Value)

		assert.Contains(t, valuesByName, metadata.NodeIDKey)
		assert.Equal(t, newContrib2.NodeID, valuesByName[metadata.NodeIDKey].Value)

		assert.NotContains(t, valuesByName, metadata.MiddleInitialKey)
	}

}

func noChanges(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.ContributorModelName, metadata.ContributorDisplayName).Build())

	contrib := metadatatest.NewContributorBuilder().WithMiddleInitial().Build()
	contrib2 := metadatatest.NewContributorBuilder().WithNodeID().Build()

	changes, err := ComputeContributorsChanges(
		schemaData,
		[]metadata.Contributor{contrib, contrib2},
		[]metadata.Contributor{contrib, contrib2},
	)
	require.NoError(t, err)
	assert.Nil(t, changes)
}

func orderChange(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.ContributorModelName, metadata.ContributorDisplayName).Build())

	contrib := metadatatest.NewContributorBuilder().WithMiddleInitial().Build()
	contrib2 := metadatatest.NewContributorBuilder().WithNodeID().Build()

	changes, err := ComputeContributorsChanges(
		schemaData,
		[]metadata.Contributor{contrib2, contrib},
		[]metadata.Contributor{contrib, contrib2},
	)
	require.NoError(t, err)
	require.NotNil(t, changes)

	expectedModel, _ := schemaData.ModelByName(metadata.ContributorModelName)
	assert.Equal(t, expectedModel.ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.NotNil(t, changes.Records)
	assert.True(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Update)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Create, 2)
	// First Contributor
	{
		values := changes.Records.Create[0].Values
		// Only contains first and last names and middle initial because other values are empty
		assert.Len(t, values, 3)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.LastNameKey)
		assert.Equal(t, contrib.LastName, valuesByName[metadata.LastNameKey].Value)

		assert.Contains(t, valuesByName, metadata.MiddleInitialKey)
		assert.Equal(t, contrib.MiddleInitial, valuesByName[metadata.MiddleInitialKey].Value)

		assert.NotContains(t, valuesByName, metadata.NodeIDKey)
	}

	//Second Contributor
	{
		values := changes.Records.Create[1].Values
		// Only contains first and last names and node id because other values are empty
		assert.Len(t, values, 3)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.LastNameKey)
		assert.Equal(t, contrib2.LastName, valuesByName[metadata.LastNameKey].Value)

		assert.Contains(t, valuesByName, metadata.NodeIDKey)
		assert.Equal(t, contrib2.NodeID, valuesByName[metadata.NodeIDKey].Value)

		assert.NotContains(t, valuesByName, metadata.MiddleInitialKey)
	}
}

func removeContributor(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.ContributorModelName, metadata.ContributorDisplayName).Build())

	contrib := metadatatest.NewContributorBuilder().WithMiddleInitial().Build()
	contrib2 := metadatatest.NewContributorBuilder().WithNodeID().Build()
	contrib3 := metadatatest.NewContributorBuilder().WithDegree().Build()

	changes, err := ComputeContributorsChanges(
		schemaData,
		[]metadata.Contributor{contrib3, contrib2, contrib},
		[]metadata.Contributor{contrib, contrib2},
	)
	require.NoError(t, err)
	require.NotNil(t, changes)

	expectedModel, _ := schemaData.ModelByName(metadata.ContributorModelName)
	assert.Equal(t, expectedModel.ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.NotNil(t, changes.Records)
	assert.True(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Update)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Create, 2)
	// First Contributor
	{
		values := changes.Records.Create[0].Values
		// Only contains first and last names and middle initial because other values are empty
		assert.Len(t, values, 3)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.LastNameKey)
		assert.Equal(t, contrib.LastName, valuesByName[metadata.LastNameKey].Value)

		assert.Contains(t, valuesByName, metadata.MiddleInitialKey)
		assert.Equal(t, contrib.MiddleInitial, valuesByName[metadata.MiddleInitialKey].Value)

		assert.NotContains(t, valuesByName, metadata.NodeIDKey)
	}

	//Second Contributor
	{
		values := changes.Records.Create[1].Values
		// Only contains first and last names and node id because other values are empty
		assert.Len(t, values, 3)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.LastNameKey)
		assert.Equal(t, contrib2.LastName, valuesByName[metadata.LastNameKey].Value)

		assert.Contains(t, valuesByName, metadata.NodeIDKey)
		assert.Equal(t, contrib2.NodeID, valuesByName[metadata.NodeIDKey].Value)

		assert.NotContains(t, valuesByName, metadata.MiddleInitialKey)
	}
}

func addContributor(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().WithModel(metadata.ContributorModelName, metadata.ContributorDisplayName).Build())

	contrib := metadatatest.NewContributorBuilder().WithMiddleInitial().Build()
	contrib2 := metadatatest.NewContributorBuilder().WithNodeID().Build()

	changes, err := ComputeContributorsChanges(
		schemaData,
		[]metadata.Contributor{contrib},
		[]metadata.Contributor{contrib, contrib2},
	)
	require.NoError(t, err)
	require.NotNil(t, changes)

	expectedModel, _ := schemaData.ModelByName(metadata.ContributorModelName)
	assert.Equal(t, expectedModel.ID, changes.ID)
	assert.Nil(t, changes.Create)

	assert.NotNil(t, changes.Records)
	assert.True(t, changes.Records.DeleteAll)
	assert.Empty(t, changes.Records.Update)
	assert.Empty(t, changes.Records.Delete)

	assert.Len(t, changes.Records.Create, 2)
	// First Contributor
	{
		values := changes.Records.Create[0].Values
		// Only contains first and last names and middle initial because other values are empty
		assert.Len(t, values, 3)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.LastNameKey)
		assert.Equal(t, contrib.LastName, valuesByName[metadata.LastNameKey].Value)

		assert.Contains(t, valuesByName, metadata.MiddleInitialKey)
		assert.Equal(t, contrib.MiddleInitial, valuesByName[metadata.MiddleInitialKey].Value)

		assert.NotContains(t, valuesByName, metadata.NodeIDKey)
	}

	//Second Contributor
	{
		values := changes.Records.Create[1].Values
		// Only contains first and last names and node id because other values are empty
		assert.Len(t, values, 3)
		valuesByName := map[string]changesetmodels.RecordValue{}
		for _, v := range values {
			valuesByName[v.Name] = v
		}
		assert.Contains(t, valuesByName, metadata.LastNameKey)
		assert.Equal(t, contrib2.LastName, valuesByName[metadata.LastNameKey].Value)

		assert.Contains(t, valuesByName, metadata.NodeIDKey)
		assert.Equal(t, contrib2.NodeID, valuesByName[metadata.NodeIDKey].Value)

		assert.NotContains(t, valuesByName, metadata.MiddleInitialKey)
	}
}
