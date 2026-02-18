package main

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"os/exec"
)

var lsFullPath string

func init() {
	// on newer Linux ls path is not just "/bin/ls", might be "/usr/bin/ls"
	var err error
	lsFullPath, err = exec.LookPath("ls")
	if err != nil {
		panic(err)
	}
}

func TestParseArgs(t *testing.T) {
	var (
		path = func(p string) func(*config) bool {
			return func(c *config) bool { return c.path == p }
		}
		args = func(a ...string) func(*config) bool {
			return func(c *config) bool { return slices.Equal(c.args, a) }
		}
		strict = func(x bool) func(*config) bool {
			return func(c *config) bool { return c.strict == x }
		}
		v = func(x bool) func(*config) bool {
			return func(c *config) bool { return c.verbose == x }
		}
		help = func(pattern string) func(*config) bool {
			return func(c *config) bool {
				if c.help == nil {
					return false
				}
				b := new(strings.Builder)
				c.help(b)
				return regexp.MustCompile(pattern).MatchString(b.String())
			}
		}
		none = func(*config) bool { return false }
		all  = func(ff ...func(*config) bool) func(*config) bool {
			return func(c *config) bool {
				for _, f := range ff {
					if !f(c) {
						return false
					}
				}
				return true
			}
		}
	)

	var tab = []struct {
		input string
		errs  string
		want  func(*config) bool
	}{
		{"", "expecting program", none},
		{"--", "expecting program", none},
		{"-h", "", help("Usage:")},
		{"--help", "", help("Usage:")},
		{"-q", "unknown flag -q", none},
		{"--quux", "unknown flag --quux", none},
		{"-q -h", "", help("Usage:")},
		{"-q --help", "", help("Usage:")},
		{"./nonexistent", "no such file or dir", none},
		{"ls", "", all(path(lsFullPath), args(), strict(false), v(false))},
		{"ls 1 2", "", all(path(lsFullPath), args("1", "2"), strict(false), v(false))},
		{"-s", "expecting program", none},
		{"-s --", "expecting program", none},
		{"-s ls", "", all(path(lsFullPath), args(), strict(true), v(false))},
		{"--strict ls", "", all(path(lsFullPath), args(), strict(true))},
		{"-s -- ls", "", all(path(lsFullPath), args(), strict(true))},
		{"--strict -- ls", "", all(path(lsFullPath), args(), strict(true))},
		{"-v ls", "", all(path(lsFullPath), args(), strict(false), v(true))},
		{"--verbose ls", "", all(path(lsFullPath), args(), strict(false), v(true))},
		{"-v -s ls", "", all(path(lsFullPath), args(), strict(true), v(true))},
	}

	for i, tc := range tab {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			c, err := parseArgs(strings.Fields(tc.input))

			switch {
			case err != nil && tc.errs == "":
				t.Errorf("unexpected error: %s", err)

			case err != nil && tc.errs != "":
				if !regexp.MustCompile(tc.errs).MatchString(err.Error()) {
					t.Errorf("expected error with %q; having: %q", tc.errs, err)
				}

			case err == nil && tc.errs != "":
				t.Errorf("expected error with %q", tc.errs)

			default:
				if !(tc.want(&c)) {
					t.Errorf("predicate failed on %+v", c)
				}
			}
		})
	}
}
