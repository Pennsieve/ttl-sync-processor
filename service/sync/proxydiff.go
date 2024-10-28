package sync

import (
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"log/slog"
)

func ComputeProxyChanges(
	schemaData *metadataclient.Schema,
	old []metadata.SavedProxy,
	new []metadata.Proxy) (*changesetmodels.ProxyChanges, error) {
	proxyChanges, err := addProxyInstanceChanges(old, new)
	if err != nil {
		return nil, err
	}
	if proxyChanges == nil {
		logger.Info("no change to proxies")
		return nil, nil
	}
	setProxySchemaCreate(proxyChanges, schemaData)
	createCount, deleteCount := proxyChanges.Summary()
	logger.Info("proxy change summary",
		slog.Int("createCount", createCount),
		slog.Int("deleteCount", deleteCount),
	)
	return proxyChanges, nil
}
func addProxyInstanceChanges(old []metadata.SavedProxy, new []metadata.Proxy) (*changesetmodels.ProxyChanges, error) {
	changesByProxyKey := map[metadata.ProxyKey]changesetmodels.ProxyRecordChanges{}

	oldByID := map[changesetmodels.ExternalInstanceID]metadata.SavedExternalIDer{}
	// anything remaining in oldToDelete after iterating over new below should be deleted.
	oldToDelete := map[changesetmodels.ExternalInstanceID]metadata.SavedProxy{}
	for _, oldInstance := range old {
		oldByID[oldInstance.Proxy.ExternalID()] = oldInstance
		oldToDelete[oldInstance.Proxy.ExternalID()] = oldInstance
	}

	for _, newInstance := range new {
		newID := newInstance.ExternalID()
		if _, keep := oldToDelete[newID]; keep {
			// don't delete if keep == true
			delete(oldToDelete, newID)
		}
		if _, exists := oldByID[newID]; !exists {
			proxyChanges, entryInitialized := changesByProxyKey[newInstance.ProxyKey]
			if !entryInitialized {
				modelName := newInstance.ModelName
				targetExternalID := newInstance.TargetExternalID
				if targetPennsieveID, found := ExistingRecordStore.GetPennsieve(modelName, targetExternalID); found {
					recordIDMap.Add(modelName, targetExternalID, targetPennsieveID)
				}
				proxyChanges.ModelName = modelName
				proxyChanges.RecordExternalID = targetExternalID
			}
			proxyChanges.NodeIDCreates = append(proxyChanges.NodeIDCreates, newInstance.PackageNodeID)
			changesByProxyKey[newInstance.ProxyKey] = proxyChanges
		}
	}

	if len(oldToDelete) > 0 {
		for _, toDelete := range oldToDelete {
			proxyChanges, entryInitialized := changesByProxyKey[toDelete.ProxyKey]
			if !entryInitialized {
				proxyChanges.ModelName = toDelete.ModelName
				proxyChanges.RecordExternalID = toDelete.TargetExternalID
			}
			proxyChanges.InstanceIDDeletes = append(proxyChanges.InstanceIDDeletes, toDelete.GetPennsieveID())
			changesByProxyKey[toDelete.ProxyKey] = proxyChanges
		}
	}
	if len(changesByProxyKey) == 0 {
		return nil, nil
	}
	var allProxyChanges []changesetmodels.ProxyRecordChanges
	for _, changes := range changesByProxyKey {
		allProxyChanges = append(allProxyChanges, changes)
	}
	return &changesetmodels.ProxyChanges{
		RecordChanges: allProxyChanges,
	}, nil
}

func setProxySchemaCreate(proxyChanges *changesetmodels.ProxyChanges, schemaData *metadataclient.Schema) {
	if proxySchema := schemaData.Proxy(); proxySchema != nil {
		logger.Info("proxySchema exists", slog.String("proxySchemaName", proxySchema.Name), slog.String("proxySchemaID", proxySchema.ID))
		proxyChanges.CreateProxyRelationshipSchema = false
	} else {
		logger.Info("proxySchema must be created")
		proxyChanges.CreateProxyRelationshipSchema = true
	}
}
