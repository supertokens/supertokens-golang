/* Copyright (c) 2026, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package passwordless

import (
	"testing"

	"github.com/nyaruka/phonenumbers/v2"
	"github.com/stretchr/testify/assert"
)

// Pins the behaviour of the default phone number validation used by the
// passwordless recipe, so that a phonenumbers upgrade cannot silently change
// which numbers the recipe accepts.
//
// Fixtures are chosen to be independent of libphonenumber metadata releases:
// +1 650 253 0000 is a long standing test number, and the +1 000 area code is
// permanently unassignable under the NANP, so neither flips validity when
// metadata is refreshed.
func TestDefaultValidatePhoneNumber(t *testing.T) {
	const invalidMsg = "Phone number is invalid"

	testCases := []struct {
		name  string
		input string
		want  *string
	}{
		{"canonical E164", "+16502530000", nil},
		{"formatted, with country code", "+1 650 253 0000", nil},
		{"national format, no default region", "6502530000", ptr(invalidMsg)},
		{"well formed but unassignable", "+1 000 000 0000", ptr(invalidMsg)},
		{"too short to be a number", "+1", ptr(invalidMsg)},
		{"empty", "", ptr(invalidMsg)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := DefaultValidatePhoneNumber(tc.input, "public")
			if tc.want == nil {
				assert.Nil(t, err, "expected %q to be valid", tc.input)
				return
			}
			if assert.NotNil(t, err, "expected %q to be invalid", tc.input) {
				assert.Equal(t, *tc.want, *err)
			}
		})
	}
}

func TestDefaultValidatePhoneNumberRejectsNonString(t *testing.T) {
	err := DefaultValidatePhoneNumber(1234567, "public")

	if assert.NotNil(t, err) {
		assert.Equal(t, "Development bug: Please make sure the phone number field yields a string", *err)
	}
}

// Pins the E164 normalisation contract that CreateCodePOST relies on: a phone
// number accepted by the default validation is stored in its canonical E164
// form, whatever separators the caller used.
func TestValidPhoneNumberIsNormalisedToE164(t *testing.T) {
	for _, input := range []string{"+16502530000", "+1 650 253 0000", "+1 (650) 253-0000"} {
		t.Run(input, func(t *testing.T) {
			assert.Nil(t, DefaultValidatePhoneNumber(input, "public"))

			parsed, err := phonenumbers.Parse(input, "")
			assert.Nil(t, err)
			assert.True(t, phonenumbers.IsValidNumber(parsed))
			assert.Equal(t, "+16502530000", phonenumbers.Format(parsed, phonenumbers.E164))
		})
	}
}

// CreateCodePOST falls back to the trimmed input when Parse fails, which only
// happens for a number a custom ValidatePhoneNumber accepted. Pins that such
// input is still rejected by Parse, so the fallback branch stays reachable.
func TestUnparseablePhoneNumberStaysUnparseable(t *testing.T) {
	_, err := phonenumbers.Parse("6502530000", "")

	assert.NotNil(t, err)
}

func ptr(s string) *string {
	return &s
}
