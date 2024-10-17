package metadatatest

import (
	"fmt"
	"github.com/google/uuid"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/stretchr/testify/assert"
)

func NewExternalInstanceID() changesetmodels.ExternalInstanceID {
	return changesetmodels.ExternalInstanceID(uuid.NewString())
}

func NewPennsieveInstanceID() changesetmodels.PennsieveInstanceID {
	return changesetmodels.PennsieveInstanceID(uuid.NewString())
}

func NewCollectionNodeID() string {
	return fmt.Sprintf("N:collection:%s", uuid.NewString())
}

func AssertExternalInstanceIDEqual(t assert.TestingT, expectedExternalID string, actualExternalID changesetmodels.ExternalInstanceID) bool {
	return assert.Equal(t, changesetmodels.ExternalInstanceID(expectedExternalID), actualExternalID)
}

func AssertPennsieveInstanceIDEqual(t assert.TestingT, expectedPennsieveID string, actualPennsieveID changesetmodels.PennsieveInstanceID) bool {
	return assert.Equal(t, changesetmodels.PennsieveInstanceID(expectedPennsieveID), actualPennsieveID)
}
