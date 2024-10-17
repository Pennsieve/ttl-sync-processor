package metadatatest

import (
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

type ProxyBuilder struct {
	entityID      *changesetmodels.ExternalInstanceID
	packageNodeID *string
}

func NewProxyBuilder() *ProxyBuilder {
	return &ProxyBuilder{}
}

func (b *ProxyBuilder) WithEntityID(entityID changesetmodels.ExternalInstanceID) *ProxyBuilder {
	b.entityID = &entityID
	return b
}

func (b *ProxyBuilder) WithPackageNodeID(nodeID string) *ProxyBuilder {
	b.packageNodeID = &nodeID
	return b
}

func (b *ProxyBuilder) Build(modelName string) metadata.Proxy {
	proxy := metadata.Proxy{
		ProxyKey: metadata.ProxyKey{
			ModelName: modelName,
		},
		PackageNodeID: "",
	}
	if b.entityID == nil {
		proxy.EntityID = NewExternalInstanceID()
	} else {
		proxy.EntityID = *b.entityID
	}
	if b.packageNodeID == nil {
		proxy.PackageNodeID = NewCollectionNodeID()
	} else {
		proxy.PackageNodeID = *b.packageNodeID
	}
	return proxy
}

func NewProxy(modelName string) metadata.Proxy {
	return NewProxyBuilder().Build(modelName)
}

func NewSavedProxy(proxy metadata.Proxy) metadata.SavedProxy {
	return metadata.SavedProxy{
		PennsieveID: NewPennsieveInstanceID(),
		RecordID:    NewPennsieveInstanceID(),
		Proxy:       proxy,
	}
}
