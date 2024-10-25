package fromrecord

import (
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
)

type RecordIDKey struct {
	ModelName        string
	ExternalRecordID changesetmodels.ExternalInstanceID
}

// RecordIDStore maps (modelName, externalRecordID) -> PennsieveInstanceID and the inverse
type RecordIDStore struct {
	externalToPennsieve map[RecordIDKey]changesetmodels.PennsieveInstanceID
	// pennsieveToExternal is the inverse of externalToPennsieve
	pennsieveToExternal map[changesetmodels.PennsieveInstanceID]*RecordIDKey
}

func NewRecordIDStore() *RecordIDStore {
	return &RecordIDStore{
		externalToPennsieve: make(map[RecordIDKey]changesetmodels.PennsieveInstanceID),
		pennsieveToExternal: make(map[changesetmodels.PennsieveInstanceID]*RecordIDKey),
	}
}

func (m *RecordIDStore) Add(modelName string, externalRecordID changesetmodels.ExternalInstanceID, recordID changesetmodels.PennsieveInstanceID) {
	key := RecordIDKey{
		ModelName:        modelName,
		ExternalRecordID: externalRecordID,
	}
	m.externalToPennsieve[key] = recordID
	m.pennsieveToExternal[recordID] = &key
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
