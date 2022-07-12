package ver

import (
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestGenericVersion_Less(t *testing.T) {
	tests := []struct {
		name             string
		left             string
		right            string
		expectLeftError  bool
		expectRightError bool
		expectLess       bool
		expectGreater    bool
	}{
		{
			name:          "normal",
			left:          "v1.1.1",
			right:         "2",
			expectLess:    true,
			expectGreater: false,
		},
		{
			name:          "short is less",
			left:          "v1.1.1",
			right:         "v1.1",
			expectLess:    false,
			expectGreater: true,
		},
		{
			name:          "equal",
			left:          "v1.1.1",
			right:         "1.1.1",
			expectLess:    false,
			expectGreater: false,
		},
		{
			name:          "no digit",
			left:          "alpha",
			right:         "beta",
			expectLess:    false,
			expectGreater: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, err := NewGeneric(tt.left)
			if tt.expectLeftError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			right, err := NewGeneric(tt.right)
			if tt.expectRightError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectLess, left.Less(right))
			require.Equal(t, tt.expectGreater, right.Less(left))
		})
	}
}

func TestGenericVersion_Sort(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
		expect   []string
	}{
		{
			name:     "normal",
			versions: []string{"v1.2", "v2", "v1.3", "v0.0.1-1", "v1.0", "0.0.1-2"},
			expect:   []string{"v0.0.1-1", "0.0.1-2", "v1.0", "v1.2", "v1.3", "v2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var versions GenericVersions
			for i := range tt.versions {
				gv, err := NewGeneric(tt.versions[i])
				require.NoError(t, err)
				versions = append(versions, *gv)
			}

			sort.Sort(versions)
			var actual []string
			for i := range versions {
				actual = append(actual, versions[i].String())
			}

			require.Equal(t, tt.expect, actual)
		})
	}
}
