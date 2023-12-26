/*
Copyright 2019 PlanetScale Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package names

import (
	"strings"
	"testing"
)

// TestJoin checks determinism and uniqueness.
func TestJoin(t *testing.T) {
	// Check that it starts with the parts joined by '-'.
	if got, want := Join(DefaultConstraints, "one", "two", "three"), "one-two-three-"; !strings.HasPrefix(got, want) {
		t.Errorf("got %q, want prefix %q", got, want)
	}

	// Check determinism and uniqueness.
	table := []struct {
		name        string
		a, b        []string
		shouldEqual bool
	}{
		{
			name:        "same parts, same order",
			a:           []string{"one", "two", "three"},
			b:           []string{"one", "two", "three"},
			shouldEqual: true,
		},
		{
			name:        "same parts, different order",
			a:           []string{"one", "two", "three"},
			b:           []string{"one", "three", "two"},
			shouldEqual: false,
		},
		{
			name:        "different parts",
			a:           []string{"one", "two", "three"},
			b:           []string{"one", "two", "four"},
			shouldEqual: false,
		},
		{
			name:        "substring moved to adjacent part",
			a:           []string{"one-two", "three-four"},
			b:           []string{"one", "two-three-four"},
			shouldEqual: false,
		},
		{
			name:        "one part split into two parts",
			a:           []string{"one-two", "three-four"},
			b:           []string{"one-two", "three", "four"},
			shouldEqual: false,
		},
	}
	for _, test := range table {
		if got := Join(DefaultConstraints, test.a...) == Join(DefaultConstraints, test.b...); got != test.shouldEqual {
			t.Errorf("Join: %s: got %v; want %v", test.name, got, test.shouldEqual)
		}
	}
}

// TestJoinSalt checks that the salt affects the hash without appearing in the name.
func TestJoinSalt(t *testing.T) {
	salt := []string{"salt1", "salt2"}
	parts := []string{"hello", "world"}
	want := "hello-world-462f1b88"
	if got := JoinSalt(DefaultConstraints, salt, parts...); got != want {
		t.Errorf("JoinSalt(%v, %v) = %q, want %q", salt, parts, got, want)
	}

	salt = []string{"salt1-salt2"}
	parts = []string{"hello", "world"}
	want = "hello-world-c65378ee"
	if got := JoinSalt(DefaultConstraints, salt, parts...); got != want {
		t.Errorf("JoinSalt(%v, %v) = %q, want %q", salt, parts, got, want)
	}
}

func TestJoinConstraints(t *testing.T) {
	cons := Constraints{
		MaxLength:      50,
		ValidFirstChar: isLowercaseLetter,
	}

	table := []struct {
		input      []string
		wantPrefix string
	}{
		{
			input:      []string{"UpperCase-Letters", "Are-Lowercased"},
			wantPrefix: "uppercase-letters-are-lowercased-",
		},
		{
			input:      []string{"allowed-symbols", "dont---change"},
			wantPrefix: "allowed-symbols-dont---change-",
		},
		{
			input:      []string{"disallowed_symbols", "are.replaced"},
			wantPrefix: "disallowed-symbols-are-replaced-",
		},
		{
			input:      []string{"really-really-ridiculously-long-inputs-are-truncated"},
			wantPrefix: "really-really-ridiculously-long-inputs----",
		},
		{
			input:      []string{"-disallowed first chars", "-are prefixed"},
			wantPrefix: "x-disallowed-first-chars--are-prefixed-",
		},
		{
			input:      []string{"Transformed first char", "is ok"},
			wantPrefix: "transformed-first-char-is-ok-",
		},
	}

	for _, test := range table {
		got := Join(cons, test.input...)
		if !strings.HasPrefix(got, test.wantPrefix) {
			t.Errorf("Join(%v) = %q; want prefix %q", test.input, got, test.wantPrefix)
		}
	}
}

// TestJoinConstraintsMaxLength checks that values are truncated to fit
// within the max length.
func TestJoinConstraintsMaxLength(t *testing.T) {
	cons := Constraints{
		MaxLength:      25,
		ValidFirstChar: isLowercaseAlphanumeric,
	}

	// The total length after truncation should be equal to MaxLength.
	out := Join(cons, strings.Repeat("a", 20), strings.Repeat("b", 20))
	if len(out) != cons.MaxLength {
		t.Errorf("len(%q) = %v; want %v", out, len(out), cons.MaxLength)
	}

	// The outputs should still be unique thanks to the hash suffix,
	// even if the truncated portion is the same because the difference between
	// inputs is at the end that gets cut off.
	out1 := Join(cons, strings.Repeat("a", 20), strings.Repeat("b", 100)+"1")
	out2 := Join(cons, strings.Repeat("a", 20), strings.Repeat("b", 100)+"2")
	if out1 == out2 {
		t.Errorf("got same output for two different inputs: %v", out1)
	}
}

// TestJoinWithConstraintsTransform checks that outputs are still
// distinguishable (thanks to the hash) even if inputs differ only in ways that
// are otherwise invisible after transformation.
func TestJoinConstraintsTransform(t *testing.T) {
	cons := Constraints{
		MaxLength:      50,
		ValidFirstChar: isLowercaseAlphanumeric,
	}

	// The outputs should still be unique thanks to the hash suffix,
	// even if the differences between inputs are otherwise invisible after
	// transformation.
	out1 := Join(cons, "disallowed_symbol")
	out2 := Join(cons, "disallowed/symbol")
	if out1 == out2 {
		t.Errorf("got same output for two different inputs: %v", out1)
	}
}
