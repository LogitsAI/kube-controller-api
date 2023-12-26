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

// Constraints specifies rules that the output of JoinWithConstraints must follow.
type Constraints struct {
	// MaxLength is the maximum length of the output, to be enforced after any
	// transformations and including the hash suffix. If a name has to be
	// truncated to fit within this maximum length, the hash at the end will be
	// preceded by a special truncation mark: "---" rather than the usual "-".
	//
	// MaxLength must be at least 12 because that's the shortest possible
	// truncated value (1 char + truncation mark + hash). Passing a value less
	// than 12 will result in a panic.
	MaxLength int
	// ValidFirstChar is a function that returns whether the given rune is
	// allowed as the first character in the output.
	ValidFirstChar func(r rune) bool
}

var (
	// DefaultConstraints are the name constraints for objects in Kubernetes
	// that don't have any special rules.
	DefaultConstraints = Constraints{
		MaxLength:      253,
		ValidFirstChar: isLowercaseAlphanumeric,
	}
	// ServiceConstraints are name constraints for Service objects.
	ServiceConstraints = Constraints{
		MaxLength:      63,
		ValidFirstChar: isLowercaseLetter,
	}
)

func isLowercaseLetter(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isUppercaseLetter(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLowercaseAlphanumeric(r rune) bool {
	return isLowercaseLetter(r) || isDigit(r)
}
