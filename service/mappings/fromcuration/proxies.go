package fromcuration

import (
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"strings"
)

func MapProxies(samples []metadata.Sample, subjects []metadata.Subject, dirRecords []curation.Record, dirStructure []curation.DirStructureEntry) ([]metadata.Proxy, error) {
	// Set up the common data
	var sampleDirRecords []curation.Record
	var subjectDirRecords []curation.Record
	for _, dirRecord := range dirRecords {
		switch dirRecord.Type {
		case curation.SampleRecordType:
			sampleDirRecords = append(sampleDirRecords, dirRecord)
		case curation.SubjectRecordType:
			subjectDirRecords = append(subjectDirRecords, dirRecord)
		}
	}
	// Note that the remoteIDs from curation are not Node ids. They are missing the initial `N:`
	remoteIDByPath := make(map[string]string, len(dirStructure))
	for _, dirEntry := range dirStructure {
		remoteIDByPath[dirEntry.DatasetRelativePath] = dirEntry.RemoteID
	}

	// Group samples by primary key since that is the join key for samples
	sampleByPrimaryKey := make(map[string]metadata.ExternalIDer, len(samples))
	for _, sample := range samples {
		sampleByPrimaryKey[sample.PrimaryKey] = sample
	}
	proxies, err := mapProxiesOfType(metadata.SampleModelName, sampleDirRecords, sampleByPrimaryKey, remoteIDByPath)
	if err != nil {
		return nil, fmt.Errorf("error mapping sample proxies: %w", err)
	}

	// Group subjects by ID since that is the join key for subjects
	subjectByID := make(map[changesetmodels.ExternalInstanceID]metadata.ExternalIDer, len(subjects))
	for _, subject := range subjects {
		subjectByID[subject.ID] = subject
	}
	subjectProxies, err := mapProxiesOfType(metadata.SubjectModelName, subjectDirRecords, subjectByID, remoteIDByPath)
	if err != nil {
		return nil, fmt.Errorf("error mapping subject proxies: %w", err)
	}
	proxies = append(proxies, subjectProxies...)
	return proxies, nil
}

func mapProxiesOfType[JOIN ~string](modelName string, dirRecords []curation.Record, byJoinKey map[JOIN]metadata.ExternalIDer, remoteIDByPath map[string]string) ([]metadata.Proxy, error) {
	var proxies []metadata.Proxy

	for _, dirRecord := range dirRecords {
		specimenID := dirRecord.SpecimenID
		curationObject, foundSample := byJoinKey[JOIN(specimenID)]
		if !foundSample {
			return nil, fmt.Errorf("record with join key %s not found", specimenID)
		}
		for _, packagePath := range dirRecord.Dirs {
			remoteID, foundRemoteID := remoteIDByPath[packagePath]
			if !foundRemoteID {
				return nil, fmt.Errorf("remote id for package path %s not found", packagePath)
			}
			nodeID := remoteID
			if !strings.HasPrefix(nodeID, "N:") {
				nodeID = fmt.Sprintf("N:%s", nodeID)
			}
			proxies = append(proxies, metadata.Proxy{
				ProxyKey: metadata.ProxyKey{
					ModelName: modelName,
					EntityID:  curationObject.ExternalID(),
				},
				PackageNodeID: nodeID,
			})
		}
	}
	return proxies, nil
}
