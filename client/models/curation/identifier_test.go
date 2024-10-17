package curation

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewIdentifier_NoSynonyms(t *testing.T) {
	identifier := newDescriptionlessIdentifier(uuid.NewString(), uuid.NewString(), uuid.NewString())
	bytes, err := json.Marshal(identifier)
	require.NoError(t, err)
	assert.NotContains(t, string(bytes), "synonyms")
	assert.NotContains(t, string(bytes), "description")
}
