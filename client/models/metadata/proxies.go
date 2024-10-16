package metadata

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

type Proxy struct {
	// EntityID is the external ID of the record that is linked to the given PackageNodeID
	// So probably either a Subject or Sample ID
	EntityID      changesetmodels.ExternalInstanceID `json:"entity_id"`
	PackageNodeID string                             `json:"package_node_id"`
}

type SavedProxy struct {
	// PennsieveID is the ID in Pennsieve of the proxy instance
	PennsieveID changesetmodels.PennsieveInstanceID `json:"pennsieve_id"`
	// RecordID is the ID in Pennsieve of the record to which the package is linked
	RecordID changesetmodels.PennsieveInstanceID `json:"record_id"`
	Proxy
}
