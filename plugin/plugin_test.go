// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/validator"
)

func TestStringInSlice(t *testing.T) {
	cases := []struct {
		stringArg      string
		sliceArg       []string
		expectedResult bool
	}{
		{
			stringArg:      "test",
			sliceArg:       []string{"element", "test", "another element"},
			expectedResult: true,
		},
		{
			stringArg:      "not element",
			sliceArg:       []string{"element", "test", "another element"},
			expectedResult: false,
		},
		{
			stringArg:      "",
			sliceArg:       []string{"element", "test", "another element"},
			expectedResult: false,
		},
		{
			stringArg:      "test",
			sliceArg:       []string{""},
			expectedResult: false,
		},
	}

	for tcIdx, tc := range cases {
		actualResult := stringInSlice(tc.stringArg, tc.sliceArg)
		if actualResult != tc.expectedResult {
			t.Fatalf("Test case #%d: unexpected result\nexpected: %t\nactual: %t",
				tcIdx+1, tc.expectedResult, actualResult)
		}
	}
}

func TestPlugin(t *testing.T) {
	privilegedUsers := []string{
		"octopus",
		"admin",
	}
	userPermissions := map[string]string{
		"johndoe": "uat",
		"lucifer": "uat;prod",
	}

	cases := []struct {
		input          *validator.Request
		expectedResult error
	}{
		{
			input: &validator.Request{
				Build: drone.Build{
					Event: "push",
				},
			},
			expectedResult: nil,
		},
		{
			input: &validator.Request{
				Build: drone.Build{
					Event:   "promote",
					Trigger: "octopus",
				},
			},
			expectedResult: nil,
		},
		{
			input: &validator.Request{
				Build: drone.Build{
					Event:   "promote",
					Trigger: "intruder",
				},
			},
			expectedResult: validator.ErrSkip,
		},
		{
			input: &validator.Request{
				Build: drone.Build{
					Event:   "promote",
					Trigger: "johndoe",
					Deploy:  "uat",
				},
			},
			expectedResult: nil,
		},
		{
			input: &validator.Request{
				Build: drone.Build{
					Event:   "promote",
					Trigger: "johndoe",
					Deploy:  "prod",
				},
			},
			expectedResult: validator.ErrSkip,
		},
		{
			input: &validator.Request{
				Build: drone.Build{
					Event:   "promote",
					Trigger: "lucifer",
					Deploy:  "prod",
				},
			},
			expectedResult: nil,
		},
		{
			input: &validator.Request{
				Build: drone.Build{
					Event:   "promote",
					Trigger: "lucifer",
					Deploy:  "prod",
				},
			},
			expectedResult: nil,
		},
	}

	plugin := New(privilegedUsers, userPermissions)

	for tcIdx, tc := range cases {
		actualResult := plugin.Validate(context.Background(), tc.input)
		if actualResult != tc.expectedResult {
			t.Fatalf("Test case #%d: unexpected result\nexpected: %s\nactual: %s",
				tcIdx+1, tc.expectedResult, actualResult)
		}
	}
}
