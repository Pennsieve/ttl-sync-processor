package fromrecord

import (
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"log/slog"
)

func MapProxies[T metadata.SavedExternalIDer](reader *metadataclient.Reader, modelName string, savedRecords []T) ([]metadata.SavedProxy, error) {
	_, modelExists := reader.Schema.ModelByName(modelName)
	if !modelExists {
		logger.Warn("model does not exist", slog.String("modelName", modelName))
		return []metadata.SavedProxy{}, nil
	}
	byRecordID, err := reader.GetProxiesForModel(modelName)
	if err != nil {
		return nil, fmt.Errorf("error reading existing %s proxies: %w", modelName, err)
	}
	savedByRecordID := map[changesetmodels.PennsieveInstanceID]metadata.SavedExternalIDer{}
	for _, savedRecord := range savedRecords {
		savedByRecordID[savedRecord.GetPennsieveID()] = savedRecord
	}
	var proxies []metadata.SavedProxy
	for recordID, proxyInstances := range byRecordID {
		pennsieveRecordID := changesetmodels.PennsieveInstanceID(recordID)
		record, foundRecord := savedByRecordID[pennsieveRecordID]
		if !foundRecord {
			return nil, fmt.Errorf("error mapping existing %s proxies: no %s record found with ID %s", modelName, modelName, recordID)
		}
		for _, proxyInstance := range proxyInstances {
			proxies = append(proxies, metadata.SavedProxy{
				PennsieveID: changesetmodels.PennsieveInstanceID(proxyInstance.ID),
				RecordID:    pennsieveRecordID,
				Proxy: metadata.Proxy{
					ProxyKey: metadata.ProxyKey{
						ModelName:        modelName,
						TargetExternalID: record.ExternalID(),
					},
					PackageNodeID: proxyInstance.Content.NodeID,
				},
			})
		}
	}
	return proxies, nil
}
