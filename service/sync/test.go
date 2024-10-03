package sync

import (
	"github.com/google/uuid"
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
)

func newModelSchemaElement(modelName, displayName string) schema.Element {
	return schema.Element{
		ID:          uuid.NewString(),
		Type:        string(schema.ModelType),
		Name:        modelName,
		DisplayName: displayName,
	}
}

type testSchemaData map[string]schema.Element

func newTestSchemaData() testSchemaData {
	return map[string]schema.Element{}
}

func (d testSchemaData) WithModel(modelName, displayName string) testSchemaData {
	d[modelName] = newModelSchemaElement(modelName, displayName)
	return d
}
