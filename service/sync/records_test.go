package sync

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringSliceEqual(t *testing.T) {
	for _, testCase := range []struct {
		scenario      string
		slice1        []string
		slice2        []string
		expectedEqual bool
	}{
		{
			"both nil",
			nil,
			nil,
			true,
		},
		{
			"both empty",
			[]string{},
			[]string{},
			true,
		},
		{
			"non-empty equal",
			[]string{"a", "b"},
			[]string{"a", "b"},
			true,
		},
		{
			"one empty",
			[]string{},
			[]string{"a", "b"},
			false,
		},
		{
			"non-empty unequal",
			[]string{"a", "b"},
			[]string{"a", "b", "c"},
			false,
		},
	} {
		t.Run(testCase.scenario, func(t *testing.T) {
			if testCase.expectedEqual {
				assert.True(t, stringSliceEqual(testCase.slice1, testCase.slice2))
			} else {
				assert.False(t, stringSliceEqual(testCase.slice1, testCase.slice2))
			}
		})
	}
}
