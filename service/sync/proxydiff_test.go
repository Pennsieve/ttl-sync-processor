package sync

import (
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings/fromrecord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
)

func TestComputeProxyChanges(t *testing.T) {
	for scenario, test := range map[string]func(t *testing.T){
		"handle edge case without panic":                 emptyChangesProxies,
		"proxy schema does not exist":                    proxySchemaDoesNotExist,
		"proxy schema exists, but no existing instances": proxySchemaExistsButNoInstances,
		"no changes":             noProxyChanges,
		"order does not matter":  proxyOrderDoesNotMatter,
		"deleted proxy":          proxyDeleted,
		"change record id":       proxyChangeRecordID,
		"change package node id": proxyChangeNodeID,
	} {
		t.Run(scenario, func(t *testing.T) {
			ExistingRecordStore = fromrecord.NewRecordIDStore()
			recordIDMap = make(fromrecord.RecordIDMap)

			test(t)
		})
	}
}

func emptyChangesProxies(t *testing.T) {
	changes, err := ComputeProxyChanges(emptySchema, []metadata.SavedProxy{}, []metadata.Proxy{})
	require.NoError(t, err)
	assert.Nil(t, changes)
}

func proxySchemaDoesNotExist(t *testing.T) {
	proxy1 := metadatatest.NewProxy(metadata.SampleModelName)
	proxy2 := metadatatest.NewProxy(metadata.SubjectModelName)

	changes, err := ComputeProxyChanges(emptySchema, []metadata.SavedProxy{}, []metadata.Proxy{proxy1, proxy2})
	require.NoError(t, err)

	assert.NotNil(t, changes)
	assert.True(t, changes.CreateProxyRelationshipSchema)

	assert.Len(t, changes.RecordChanges, 2)

	changes1 := FindProxyRecordChangeByProxyKey(changes.RecordChanges, proxy1.ProxyKey)
	require.NotNil(t, changes1)
	assert.Empty(t, changes1.InstanceIDDeletes)
	assert.Contains(t, changes1.NodeIDCreates, proxy1.PackageNodeID)

	changes2 := FindProxyRecordChangeByProxyKey(changes.RecordChanges, proxy2.ProxyKey)
	require.NotNil(t, changes2)
	assert.Empty(t, changes2.InstanceIDDeletes)
	assert.Contains(t, changes2.NodeIDCreates, proxy2.PackageNodeID)
}

func proxySchemaExistsButNoInstances(t *testing.T) {
	proxy1 := metadatatest.NewProxy(metadata.SampleModelName)
	proxy2 := metadatatest.NewProxy(metadata.SubjectModelName)

	changes, err := ComputeProxyChanges(fullSchema(), []metadata.SavedProxy{}, []metadata.Proxy{proxy1, proxy2})
	require.NoError(t, err)

	assert.NotNil(t, changes)
	assert.False(t, changes.CreateProxyRelationshipSchema)

	assert.Len(t, changes.RecordChanges, 2)

	changes1 := FindProxyRecordChangeByProxyKey(changes.RecordChanges, proxy1.ProxyKey)
	require.NotNil(t, changes1)
	assert.Empty(t, changes1.InstanceIDDeletes)
	assert.Len(t, changes1.NodeIDCreates, 1)
	assert.Contains(t, changes1.NodeIDCreates, proxy1.PackageNodeID)

	changes2 := FindProxyRecordChangeByProxyKey(changes.RecordChanges, proxy2.ProxyKey)
	require.NotNil(t, changes2)
	assert.Empty(t, changes2.InstanceIDDeletes)
	assert.Len(t, changes2.NodeIDCreates, 1)
	assert.Contains(t, changes2.NodeIDCreates, proxy2.PackageNodeID)
}

func noProxyChanges(t *testing.T) {
	proxy1 := metadatatest.NewProxy(metadata.SampleModelName)
	proxy2 := metadatatest.NewProxy(metadata.SubjectModelName)

	savedProxy1 := metadatatest.NewSavedProxy(proxy1)
	savedProxy2 := metadatatest.NewSavedProxy(proxy2)

	changes, err := ComputeProxyChanges(fullSchema(), []metadata.SavedProxy{savedProxy1, savedProxy2}, []metadata.Proxy{proxy1, proxy2})
	require.NoError(t, err)

	assert.Nil(t, changes)
}

func proxyOrderDoesNotMatter(t *testing.T) {
	proxy1 := metadatatest.NewProxy(metadata.SampleModelName)
	proxy2 := metadatatest.NewProxy(metadata.SubjectModelName)

	savedProxy1 := metadatatest.NewSavedProxy(proxy1)
	savedProxy2 := metadatatest.NewSavedProxy(proxy2)

	changes, err := ComputeProxyChanges(fullSchema(), []metadata.SavedProxy{savedProxy2, savedProxy1}, []metadata.Proxy{proxy1, proxy2})
	require.NoError(t, err)

	assert.Nil(t, changes)
}

func proxyDeleted(t *testing.T) {
	keepProxy1 := metadatatest.NewProxy(metadata.SampleModelName)
	keepProxy2 := metadatatest.NewProxy(metadata.SubjectModelName)

	savedKeepProxy1 := metadatatest.NewSavedProxy(keepProxy1)
	savedKeepProxy2 := metadatatest.NewSavedProxy(keepProxy2)
	savedDeletedProxy := metadatatest.NewSavedProxy(metadatatest.NewProxy(metadata.SampleModelName))

	changes, err := ComputeProxyChanges(fullSchema(), []metadata.SavedProxy{savedKeepProxy2, savedDeletedProxy, savedKeepProxy1}, []metadata.Proxy{keepProxy1, keepProxy2})
	require.NoError(t, err)

	assert.NotNil(t, changes)
	assert.False(t, changes.CreateProxyRelationshipSchema)

	assert.Len(t, changes.RecordChanges, 1)

	recordChanges := changes.RecordChanges[0]
	assert.Equal(t, savedDeletedProxy.ModelName, recordChanges.ModelName)
	assert.Equal(t, savedDeletedProxy.TargetExternalID, recordChanges.RecordExternalID)
	assert.Empty(t, recordChanges.NodeIDCreates)
	assert.Len(t, recordChanges.InstanceIDDeletes, 1)
	assert.Contains(t, recordChanges.InstanceIDDeletes, savedDeletedProxy.GetPennsieveID())
}

func proxyChangeRecordID(t *testing.T) {
	oldProxy := metadatatest.NewSavedProxy(metadatatest.NewProxy(metadata.SubjectModelName))
	newProxy := metadatatest.NewProxyBuilder().WithPackageNodeID(oldProxy.PackageNodeID).Build(metadata.SubjectModelName)

	changes, err := ComputeProxyChanges(fullSchema(), []metadata.SavedProxy{oldProxy}, []metadata.Proxy{newProxy})
	require.NoError(t, err)
	assert.NotNil(t, changes)
	assert.False(t, changes.CreateProxyRelationshipSchema)
	assert.Len(t, changes.RecordChanges, 2)

	deleteChanges := FindProxyRecordChangeByProxyKey(changes.RecordChanges, oldProxy.ProxyKey)
	assert.NotNil(t, deleteChanges)
	assert.Empty(t, deleteChanges.NodeIDCreates)
	assert.Len(t, deleteChanges.InstanceIDDeletes, 1)
	assert.Contains(t, deleteChanges.InstanceIDDeletes, oldProxy.GetPennsieveID())

	createChanges := FindProxyRecordChangeByProxyKey(changes.RecordChanges, newProxy.ProxyKey)
	assert.NotNil(t, createChanges)
	assert.Empty(t, createChanges.InstanceIDDeletes)
	assert.Len(t, createChanges.NodeIDCreates, 1)
	assert.Contains(t, createChanges.NodeIDCreates, newProxy.PackageNodeID)
}

func proxyChangeNodeID(t *testing.T) {
	oldProxy := metadatatest.NewSavedProxy(metadatatest.NewProxy(metadata.SampleModelName))
	newProxy := metadatatest.NewProxyBuilder().WithEntityID(oldProxy.TargetExternalID).Build(metadata.SampleModelName)

	changes, err := ComputeProxyChanges(fullSchema(), []metadata.SavedProxy{oldProxy}, []metadata.Proxy{newProxy})
	require.NoError(t, err)
	assert.NotNil(t, changes)
	assert.False(t, changes.CreateProxyRelationshipSchema)
	assert.Len(t, changes.RecordChanges, 1)

	recordChanges := changes.RecordChanges[0]
	assert.Equal(t, metadata.SampleModelName, recordChanges.ModelName)
	assert.Equal(t, oldProxy.TargetExternalID, recordChanges.RecordExternalID)

	assert.Len(t, recordChanges.NodeIDCreates, 1)
	assert.Contains(t, recordChanges.NodeIDCreates, newProxy.PackageNodeID)

	assert.Len(t, recordChanges.InstanceIDDeletes, 1)
	assert.Contains(t, recordChanges.InstanceIDDeletes, oldProxy.GetPennsieveID())

}

func FindProxyRecordChangeByProxyKey(proxyRecordChanges []changesetmodels.ProxyRecordChanges, key metadata.ProxyKey) *changesetmodels.ProxyRecordChanges {
	index := slices.IndexFunc(proxyRecordChanges, func(changes changesetmodels.ProxyRecordChanges) bool {
		return changes.ModelName == key.ModelName && changes.RecordExternalID == key.TargetExternalID
	})
	if index == -1 {
		return nil
	}
	return &proxyRecordChanges[index]
}
