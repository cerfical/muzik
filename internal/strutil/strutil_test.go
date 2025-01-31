package strutil_test

import (
	"testing"

	"github.com/cerfical/muzik/internal/strutil"
	"github.com/stretchr/testify/assert"
)

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty_string_remains_empty", "", ""},
		{"single_word_capitalize_first_letter", "string", "String"},
		{"single_word_other_letters_no_change", "sTRiNG", "STRiNG"},
		{"multiple_words_capitalize_first_letter_of_first_word", "example sTRING", "Example sTRING"},
		{"with_spaces_prefix_no_change", " string", " string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, strutil.Capitalize(tt.s))
		})
	}
}

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
