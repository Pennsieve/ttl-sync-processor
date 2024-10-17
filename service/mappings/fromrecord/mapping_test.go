package fromrecord

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSafeStringSlice(t *testing.T) {
	for scenario, params := range map[string]struct {
		valueBytes []byte
		expected   []string
	}{
		"nil value":   {nil, []string{}},
		"empty value": {[]byte(`[]`), []string{}},
		"non-empty":   {[]byte(`["one", "two", "three"]`), []string{"one", "two", "three"}},
	} {
		t.Run(scenario, func(t *testing.T) {
			// Create an any value that is []any that can be converted to a slice of strings) since that is what gets unmarshalled
			// and passed to safeStringSlice
			var value any
			if params.valueBytes != nil {
				require.NoError(t, json.Unmarshal(params.valueBytes, &value))
			}
			strings := safeStringSlice(value)
			assert.Equal(t, params.expected, strings)
		})
	}
}
