package fromrecord

import (
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
)

type RecordIDKey struct {
	ModelName        string
	ExternalRecordID changesetmodels.ExternalInstanceID
}

type RecordIDMap map[RecordIDKey]changesetmodels.PennsieveInstanceID
type inverseRecordIDMap map[changesetmodels.PennsieveInstanceID]*RecordIDKey

func (m RecordIDMap) Add(modelName string, externalRecordID changesetmodels.ExternalInstanceID, recordID changesetmodels.PennsieveInstanceID) *RecordIDKey {
	key := RecordIDKey{
		ModelName:        modelName,
		ExternalRecordID: externalRecordID,
	}
	m[key] = recordID
	return &key
}

// RecordIDStore maps (modelName, externalRecordID) -> PennsieveInstanceID and the inverse
type RecordIDStore struct {
	externalToPennsieve RecordIDMap
	// pennsieveToExternal is the inverse of externalToPennsieve
	pennsieveToExternal inverseRecordIDMap
}

func NewRecordIDStore() *RecordIDStore {
	return &RecordIDStore{
		externalToPennsieve: make(RecordIDMap),
		pennsieveToExternal: make(inverseRecordIDMap),
	}
}

func (m *RecordIDStore) Add(modelName string, externalRecordID changesetmodels.ExternalInstanceID, recordID changesetmodels.PennsieveInstanceID) {
	key := m.externalToPennsieve.Add(modelName, externalRecordID, recordID)
	m.pennsieveToExternal[recordID] = key
}

func (m *RecordIDStore) GetPennsieve(modelName string, externalRecordID changesetmodels.ExternalInstanceID) (changesetmodels.PennsieveInstanceID, bool) {
	id, found := m.externalToPennsieve[RecordIDKey{
		ModelName:        modelName,
		ExternalRecordID: externalRecordID,
	}]
	return id, found
}

func (m *RecordIDStore) GetExternal(recordID changesetmodels.PennsieveInstanceID) *RecordIDKey {
	return m.pennsieveToExternal[recordID]
}

func (m *RecordIDStore) Len() int {
	return len(m.externalToPennsieve)
}
