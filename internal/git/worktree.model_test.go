package git

import (
	"testing"
)

func TestWorkTreeListFormatPattern(t *testing.T) {
	tests := []struct {
		value  string
		name   string
		path   string
		head   string
		branch string
	}{
		{
			value:  "C:/git-repos/grove                   f93d89e [dev]",
			name:   "grove",
			path:   "C:/git-repos/grove",
			head:   "f93d89e",
			branch: "dev",
		},
		{
			value:  "/home/user/git-repos/test                   a3c8913 [main]",
			name:   "test",
			path:   "/home/user/git-repos/test",
			head:   "a3c8913",
			branch: "main",
		},
	}

	for _, tc := range tests {
		t.Run(tc.value, func(t *testing.T) {
			t.Parallel()

			wt := &WorkTree{}
			err := wt.Scan(tc.value)
			if err != nil {
				t.Fatal(err)
				return
			}

			if tc.name != wt.Name {
				t.Fatalf("%v != %v", tc.name, wt.Name)
			}

			if tc.path != wt.Path {
				t.Fatalf("%v != %v", tc.path, wt.Path)
			}

			if tc.head != wt.Head {
				t.Fatalf("%v != %v", tc.head, wt.Head)
			}

			if tc.branch != wt.Branch {
				t.Fatalf("%v != %v", tc.branch, wt.Branch)
			}
		})
	}
}
