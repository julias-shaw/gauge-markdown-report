/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/

package mdgen

import "testing"

// The Kind() methods exist so item-walkers can dispatch without type
// assertions. Tests pin them to the expected constant — if someone
// reorders the kind enum these break loud and early.

func TestStepKind(t *testing.T) {
	if got := (&step{}).Kind(); got != stepKind {
		t.Errorf("step.Kind() = %q, want %q", got, stepKind)
	}
}

func TestConceptKind(t *testing.T) {
	if got := (&concept{}).Kind(); got != conceptKind {
		t.Errorf("concept.Kind() = %q, want %q", got, conceptKind)
	}
}

func TestCommentKind(t *testing.T) {
	if got := (&comment{}).Kind(); got != commentKind {
		t.Errorf("comment.Kind() = %q, want %q", got, commentKind)
	}
}

func TestBuildErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  buildError
		want string
	}{
		{"parse-error-prefix", buildError{ErrorType: parseErrorType, Message: "bad token"}, "[Parse Error] bad token"},
		{"validation-error-prefix", buildError{ErrorType: validationErrorType, Message: "bad ref"}, "[Validation Error] bad ref"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildErrorIsParseError(t *testing.T) {
	if !(buildError{ErrorType: parseErrorType}).isParseError() {
		t.Error("parse error should report isParseError() == true")
	}
	if (buildError{ErrorType: validationErrorType}).isParseError() {
		t.Error("validation error should report isParseError() == false")
	}
}
