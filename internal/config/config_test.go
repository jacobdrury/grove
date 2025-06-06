package config

import (
	"testing"
)

func TestBranchResolver(t *testing.T) {
	br := BranchResolver{
		BranchPrefixAliases: map[BranchPrefixAlias]BranchPrefix{"u": "user1"},
		BranchDelimiter:     "/",
	}

	branches := []string{
		"user1/fm-331-asdf-asdf",
		"user1/fm-432-asdf-test",
		"user1/fm-432-asdf-test-2",
		"johndoe/asdf-asdf-2",
	}

	tests := []struct {
		in  string
		out string
	}{
		{in: "example", out: "example"},
		{in: "u/fm-331", out: "user1/fm-331-asdf-asdf"},
		{in: "u/fm-432", out: "user1/fm-432-asdf-test"},
		{in: "u/fm-432-asdf-test", out: "user1/fm-432-asdf-test"},
		{in: "u/fm-432-asdf-test-2", out: "user1/fm-432-asdf-test-2"},
		{in: "u/fm-554-asdfasdf", out: "user1/fm-554-asdfasdf"},
		{in: "n/fm-432", out: "n/fm-432"},
		{in: "john/asdf-asdf", out: "john/asdf-asdf"},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()

			if resolved := br.Resolve(tc.in, branches); resolved != tc.out {
				t.Errorf("'%v' -> '%v' != '%v'", tc.in, resolved, tc.out)
			}
		})
	}
}
