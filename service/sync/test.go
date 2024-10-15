package sync

import (
	"github.com/google/uuid"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
)

func newSchemaElement(elementName, elementDisplayName string, elementType schema.Type) schema.Element {
	return schema.Element{
		ID:          uuid.NewString(),
		Type:        string(elementType),
		Name:        elementName,
		DisplayName: elementDisplayName,
	}
}

type testSchemaData struct {
	elements                []schema.Element
	proxyRelationshipSchema *schema.NullableRelationship
}

func newTestSchemaData() *testSchemaData {
	return &testSchemaData{}
}

func (d *testSchemaData) WithModel(modelName, displayName string) *testSchemaData {
	d.elements = append(d.elements, newSchemaElement(modelName, displayName, schema.ModelType))
	return d
}

func (d *testSchemaData) WithLinkedProperty(linkedPropertyName, displayName string) *testSchemaData {
	d.elements = append(d.elements, newSchemaElement(linkedPropertyName, displayName, schema.LinkedPropertyType))
	return d
}

func (d *testSchemaData) WithProxyRelationshipSchema() *testSchemaData {
	d.proxyRelationshipSchema = &schema.NullableRelationship{
		ID:          uuid.NewString(),
		Name:        schema.ProxyName,
		DisplayName: schema.ProxyDisplayName,
	}
	return d
}

func (d *testSchemaData) Build() ([]schema.Element, *schema.NullableRelationship) {
	return d.elements, d.proxyRelationshipSchema
}

var emptySchema = metadataclient.NewSchema(nil, nil)

func findValueByName(t *testing.T, values []changesetmodels.RecordValue, name string) changesetmodels.RecordValue {
	index := slices.IndexFunc(values, func(value changesetmodels.RecordValue) bool {
		return value.Name == name
	})
	require.GreaterOrEqual(t, index, 0)
	return values[index]
}

func findRecordUpdateByPennsieveID(t *testing.T, updates []changesetmodels.RecordUpdate, pennsieveID changesetmodels.PennsieveInstanceID) changesetmodels.RecordUpdate {
	index := slices.IndexFunc(updates, func(update changesetmodels.RecordUpdate) bool {
		return update.PennsieveID == pennsieveID
	})
	require.GreaterOrEqual(t, index, 0)
	return updates[index]
}
