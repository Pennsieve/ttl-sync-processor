package metadatatest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubjectCopy(t *testing.T) {
	original := NewSubjectBuilder().WithSpeciesSynonyms(3).Build()
	copied := SubjectCopy(original)
	assert.NotSame(t, &original, &copied)
	assert.NotSame(t, &original.SpeciesSynonyms[0], &copied.SpeciesSynonyms[0])

	copied.SpeciesSynonyms = append(copied.SpeciesSynonyms, "cat")
	assert.NotEqual(t, original.SpeciesSynonyms, copied.SpeciesSynonyms)
}
