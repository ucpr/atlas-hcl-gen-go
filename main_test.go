package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toCamelCase(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "snake case pattern",
			in:   "created_at",
			out:  "CreatedAt",
		},
		{
			name: "camel case pattern",
			in:   "CreatedAt",
			out:  "CreatedAt",
		},
		{
			name: "snake case and camel case pattern",
			in:   "created_At",
			out:  "CreatedAt",
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := toCamelCase(tt.in)
			assert.Equal(t, tt.out, got)
		})
	}
}
