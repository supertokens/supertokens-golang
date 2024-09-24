/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package emailpassword

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultEmailValidator(t *testing.T) {
	assert.Nil(t, defaultEmailValidator("test@supertokens.io", "public"))
	assert.Nil(t, defaultEmailValidator("nsdafa@gmail.com", "public"))
	assert.Nil(t, defaultEmailValidator("fewf3r_fdkj@gmaildsfa.co.uk", "public"))
	assert.Nil(t, defaultEmailValidator("dafk.adfa@gmail.com", "public"))
	assert.Nil(t, defaultEmailValidator("skjlblc3f3@fnldsks.co", "public"))
	assert.Equal(t, "Email is invalid", *defaultEmailValidator("sdkjfnas34@gmail.com.c", "public"))
	assert.Equal(t, "Email is invalid", *defaultEmailValidator("d@c", "public"))
	assert.Equal(t, "Email is invalid", *defaultEmailValidator("fasd", "public"))
	assert.Equal(t, "Email is invalid", *defaultEmailValidator("dfa@@@abc.com", "public"))
	assert.Equal(t, "Email is invalid", *defaultEmailValidator("", "public"))
	assert.Equal(t, "Field is not optional", *defaultEmailValidator(nil, "public"))
}

func TestDefaultPasswordValidator(t *testing.T) {
	assert.Nil(t, defaultPasswordValidator("dsknfkf38H", "public"))
	assert.Nil(t, defaultPasswordValidator("lasdkf*787~sdfskj", "public"))
	assert.Nil(t, defaultPasswordValidator("L0493434505", "public"))
	assert.Nil(t, defaultPasswordValidator("3453342422L", "public"))
	assert.Nil(t, defaultPasswordValidator("1sdfsdfsdfsd", "public"))
	assert.Nil(t, defaultPasswordValidator("dksjnlvsnl2", "public"))
	assert.Nil(t, defaultPasswordValidator("abcgftr8", "public"))
	assert.Nil(t, defaultPasswordValidator("  d3    ", "public"))
	assert.Nil(t, defaultPasswordValidator("abc!@#$%^&*()gftr8", "public"))
	assert.Nil(t, defaultPasswordValidator("    dskj3", "public"))
	assert.Nil(t, defaultPasswordValidator("    dsk  3", "public"))
	assert.Equal(t, "Field is not optional", *defaultPasswordValidator(nil, "public"))

	assert.Equal(t, "Password must contain at least 8 characters, including a number", *defaultPasswordValidator("asd", "public"))

	assert.Equal(t, "Password's length must be lesser than 100 characters", *defaultPasswordValidator("asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4", "public"))

	assert.Equal(t, "Password must contain at least one number", *defaultPasswordValidator("ascdvsdfvsIUOO", "public"))

	assert.Equal(t, "Password must contain at least one alphabet", *defaultPasswordValidator("234235234523", "public"))
}

func TestInvalidAPIInputForFormFields(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	objToJson := func(obj interface{}) []byte {
		jsonBytes, err := json.Marshal(obj)
		assert.NoError(t, err)
		return jsonBytes
	}

	testCases := []struct {
		input      interface{}
		expected   string
		fieldError bool
	}{
		{
			input:      map[string]interface{}{},
			expected:   "Missing input param: formFields",
			fieldError: false,
		},
		{
			input: map[string]interface{}{
				"formFields": "abcd",
			},
			expected:   "formFields must be an array",
			fieldError: false,
		},
		{
			input: map[string]interface{}{
				"formFields": []string{"hello"},
			},
			expected:   "formFields must be an array of objects containing id and value of type string",
			fieldError: false,
		},
		{
			input: map[string]interface{}{
				"formFields": []map[string]interface{}{
					{
						"hello": "world",
					},
					{
						"world": "hello",
					},
				},
			},
			expected:   "Field is not optional",
			fieldError: true,
		},
		{
			input: map[string]interface{}{
				"formFields": []map[string]interface{}{
					{
						"id": 1,
					},
				},
			},
			expected:   "formFields must be an array of objects containing id and value of type string",
			fieldError: false,
		},
		{
			input: map[string]interface{}{
				"formFields": []map[string]interface{}{
					{
						"id":    1,
						"value": "world",
					},
				},
			},
			expected:   "formFields must be an array of objects containing id and value of type string",
			fieldError: false,
		},
		{
			input: map[string]interface{}{
				"formFields": []map[string]interface{}{
					{
						"id":    "hello",
						"value": 1,
					},
				},
			},
			expected:   "formFields must be an array of objects containing id and value of type string",
			fieldError: false,
		},
		{
			input: map[string]interface{}{
				"formFields": []map[string]interface{}{
					{
						"value": 1,
					},
				},
			},
			expected:   "formFields must be an array of objects containing id and value of type string",
			fieldError: false,
		},
		{
			input: map[string]interface{}{
				"formFields": []map[string]interface{}{
					{
						"id": "hello",
					},
					{
						"id": "world",
					},
				},
			},
			expected:   "Field is not optional",
			fieldError: true,
		},
		{
			input: map[string]interface{}{
				"formFields": []map[string]interface{}{
					{
						"value": "hello",
					},
					{
						"value": "world",
					},
				},
			},
			expected:   "Field is not optional",
			fieldError: true,
		},
	}

	APIs := []string{
		"/auth/signup",
		"/auth/signin",
	}

	for _, testCase := range testCases {
		for _, api := range APIs {
			resp, err := http.Post(testServer.URL+api, "application/json", bytes.NewBuffer(objToJson(testCase.input)))
			assert.NoError(t, err)

			if testCase.fieldError {
				assert.Equal(t, 200, resp.StatusCode)
			} else {
				assert.Equal(t, 400, resp.StatusCode)
			}
			dataInBytes1, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err.Error())
			}
			resp.Body.Close()
			var data map[string]interface{}
			err = json.Unmarshal(dataInBytes1, &data)
			if err != nil {
				t.Error(err.Error())
			}

			if testCase.fieldError {
				assert.Equal(t, "FIELD_ERROR", data["status"].(string))

				for _, formField := range data["formFields"].([]interface{}) {
					errorMessage := formField.(map[string]interface{})["error"]
					assert.Equal(t, testCase.expected, errorMessage)
				}
			} else {
				assert.Equal(t, testCase.expected, data["message"])
			}

		}
	}
}
