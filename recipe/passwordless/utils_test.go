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
// passwordless recipe, independent of the phonenumbers library major version.
func TestDefaultValidatePhoneNumber(t *testing.T) {
	invalidMsg := "Phone number is invalid"

	testCases := []struct {
		input    string
		expected *string
	}{
		{"+16502530000", nil},       // canonical E164
		{"+1 650 253 0000", nil},    // formatted, but carries a country code
		{"6502530000", &invalidMsg}, // no country code and no default region
		{"+1154", &invalidMsg},      // parses, but is not a valid number
		{"", &invalidMsg},           // unparseable
	}

	for _, tc := range testCases {
		err := DefaultValidatePhoneNumber(tc.input, "public")
		if tc.expected == nil {
			assert.Nil(t, err, "expected %v to be valid", tc.input)
		} else {
			assert.NotNil(t, err, "expected %v to be invalid", tc.input)
			assert.Equal(t, *tc.expected, *err, "unexpected error message for %v", tc.input)
		}
	}

	// non-string input returns the development bug message
	err := DefaultValidatePhoneNumber(1234567, "public")
	assert.NotNil(t, err)
	assert.Equal(t, "Development bug: Please make sure the email field yields a string", *err)
}

// Pins the E164 normalisation contract CreateCodePOST relies on: a phone
// number accepted by the default validation must format to its canonical
// E164 form before being stored.
func TestValidPhoneNumberIsNormalisedToE164(t *testing.T) {
	parsed, err := phonenumbers.Parse("+1 650 253 0000", "")
	assert.Nil(t, err)
	assert.True(t, phonenumbers.IsValidNumber(parsed))
	assert.Equal(t, "+16502530000", phonenumbers.Format(parsed, phonenumbers.E164))
}
