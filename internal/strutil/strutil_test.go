package strutil_test

import (
	"testing"

	"github.com/cerfical/muzik/internal/strutil"
	"github.com/stretchr/testify/assert"
)

func TestDedup(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		want []string
	}{
		{"empty_slice_remains_empty", []string{}, []string{}},
		{"consecutive_dups_are_removed", []string{"b1", "a1", "a1"}, []string{"a1", "b1"}},
		{"random_dups_are_removed", []string{"b1", "a1", "a1", "b1", "a1"}, []string{"a1", "b1"}},
		{"result_is_sorted", []string{"b", "h", "a"}, []string{"a", "b", "h"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, strutil.Dedup(tt.s))
		})
	}
}
