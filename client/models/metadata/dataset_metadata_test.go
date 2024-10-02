package metadata

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSync_ContributorsHash(t *testing.T) {
	contrib1 := Contributor{
		FirstName:     uuid.NewString(),
		MiddleInitial: "W",
		LastName:      uuid.NewString(),
		Degree:        "PhD",
		ORCID:         uuid.NewString(),
	}
	contrib2 := Contributor{
		FirstName: uuid.NewString(),
		LastName:  uuid.NewString(),
		NodeID:    uuid.NewString(),
	}
	contrib3 := Contributor{
		FirstName:     uuid.NewString(),
		MiddleInitial: "M",
		LastName:      uuid.NewString(),
		ORCID:         uuid.NewString(),
		NodeID:        uuid.NewString(),
	}

	for scenario, params := range map[string]struct {
		contribs1   []Contributor
		contribs2   []Contributor
		expectEqual bool
	}{
		"one contrib":   {[]Contributor{contrib1}, []Contributor{contrib1}, true},
		"two contrib":   {[]Contributor{contrib1, contrib3}, []Contributor{contrib1, contrib3}, true},
		"three contrib": {[]Contributor{contrib3, contrib2, contrib1}, []Contributor{contrib3, contrib2, contrib1}, true},
		"unequal":       {[]Contributor{contrib1}, []Contributor{contrib2}, false},
		"order matters": {[]Contributor{contrib3, contrib2, contrib1}, []Contributor{contrib1, contrib2, contrib3}, false},
		"proper subset": {[]Contributor{contrib3, contrib2}, []Contributor{contrib3, contrib2, contrib1}, false},
	} {
		t.Run(scenario, func(t *testing.T) {
			hash1, err := DatasetMetadata{Contributors: params.contribs1}.ContributorsHash()
			require.NoError(t, err)
			hash2, err := DatasetMetadata{Contributors: params.contribs2}.ContributorsHash()
			require.NoError(t, err)
			if params.expectEqual {
				assert.Equal(t, hash1, hash2)
			} else {
				assert.NotEqual(t, hash1, hash2)
			}
		})
	}

}
