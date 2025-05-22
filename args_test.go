package main

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

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
		{"./nonexistent", "no such file or dir", none},
		{"ls", "", all(path("/bin/ls"), args(), strict(false))},
		{"ls 1 2", "", all(path("/bin/ls"), args("1", "2"), strict(false))},
		{"-s", "expecting program", none},
		{"-s --", "expecting program", none},
		{"-s ls", "", all(path("/bin/ls"), args(), strict(true))},
		{"--strict ls", "", all(path("/bin/ls"), args(), strict(true))},
		{"-s -- ls", "", all(path("/bin/ls"), args(), strict(true))},
		{"--strict -- ls", "", all(path("/bin/ls"), args(), strict(true))},
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
