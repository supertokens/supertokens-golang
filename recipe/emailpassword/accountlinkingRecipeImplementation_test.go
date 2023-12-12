/* Copyright (c) 2023, VRAI Labs and/or its affiliates. All rights reserved.
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

// we have this file here case we cannot put it in supertokens or unittesting
// package due to cyclic imports.

func TestGetOldestUsersFirst(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	telemetry := false
	supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:   "Testing",
			Origin:    "http://localhost:3000",
			APIDomain: "http://localhost:3001",
		},
		Telemetry: &telemetry,
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	})

	user1, err := SignUp("public", "test@gmail.com", "testPass123")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = SignUp("public", "test1@gmail.com", "testPass123")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = SignUp("public", "test2@gmail.com", "testPass123")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = SignUp("public", "test3@gmail.com", "testPass123")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = SignUp("public", "test4@gmail.com", "testPass123")
	if err != nil {
		t.Error(err)
		return
	}

	{
		paginationResult, err := supertokens.GetUsersOldestFirst("public", nil, nil, nil, nil)
		if err != nil {
			t.Error(err)
			return
		}

		email := "test@gmail.com"
		assert.Nil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 5)
		assert.True(t, paginationResult.Users[0].LoginMethods[0].HasSameEmailAs(&email))
		assert.Equal(t, paginationResult.Users[0].ID, user1.OK.User.ID)
		assert.Equal(t, paginationResult.Users[0].ID, paginationResult.Users[0].LoginMethods[0].RecipeUserID.GetAsString())
	}

	{
		limit := 1
		paginationResult, err := supertokens.GetUsersOldestFirst("public", nil, &limit, nil, nil)
		if err != nil {
			t.Error(err)
			return
		}

		assert.NotNil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 1)
		assert.Equal(t, paginationResult.Users[0].Emails[0], "test@gmail.com")

		paginationResult, err = supertokens.GetUsersOldestFirst("public", paginationResult.NextPaginationToken, &limit, nil, nil)
		if err != nil {
			t.Error(err)
			return
		}

		assert.NotNil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 1)
		assert.Equal(t, paginationResult.Users[0].Emails[0], "test1@gmail.com")

		limit = 5
		paginationResult, err = supertokens.GetUsersOldestFirst("public", paginationResult.NextPaginationToken, &limit, nil, nil)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Nil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 3)
		assert.Equal(t, paginationResult.Users[0].Emails[0], "test2@gmail.com")
	}

	{
		paginationToken := "invalid"
		_, err := supertokens.GetUsersOldestFirst("public", &paginationToken, nil, nil, nil)
		if err != nil {
			assert.Contains(t, err.Error(), "invalid pagination token")
		} else {
			assert.Fail(t, "pagination token invalid should fail")
		}
	}
}
