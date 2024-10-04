package sync

import (
	"github.com/google/uuid"
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
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

func findValueByName(t *testing.T, values []changesetmodels.RecordValue, name string) changesetmodels.RecordValue {
	index := slices.IndexFunc(values, func(value changesetmodels.RecordValue) bool {
		return value.Name == name
	})
	require.GreaterOrEqual(t, index, 0)
	return values[index]
}

func findRecordUpdateByPennsieveID(t *testing.T, updates []changesetmodels.RecordUpdate, pennsieveID changesetmodels.PennsieveRecordID) changesetmodels.RecordUpdate {
	index := slices.IndexFunc(updates, func(update changesetmodels.RecordUpdate) bool {
		return update.PennsieveID == pennsieveID
	})
	require.GreaterOrEqual(t, index, 0)
	return updates[index]
}
