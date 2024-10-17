package metadata

import (
	"fmt"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
)

type ProxyKey struct {
	// ModelName is the name of the model of this proxy's record
	ModelName string `json:"model_name"`
	// EntityID is the external ID of the record that is linked to the given PackageNodeID
	// So probably either a Subject or Sample ID
	EntityID changesetmodels.ExternalInstanceID `json:"entity_id"`
}

type Proxy struct {
	ProxyKey
	PackageNodeID string `json:"package_node_id"`
}

func (p Proxy) ExternalID() changesetmodels.ExternalInstanceID {
	return changesetmodels.ExternalInstanceID(fmt.Sprintf("%s::%s::%s", p.ModelName, p.EntityID, p.PackageNodeID))
}

type SavedProxy struct {
	// PennsieveID is the ID in Pennsieve of the proxy instance
	PennsieveID changesetmodels.PennsieveInstanceID `json:"pennsieve_id"`
	// RecordID is the ID in Pennsieve of the record to which the package is linked
	RecordID changesetmodels.PennsieveInstanceID `json:"record_id"`
	Proxy
}

func (sp SavedProxy) GetPennsieveID() changesetmodels.PennsieveInstanceID {
	return sp.PennsieveID
}
