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
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
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

func TestGetNewestUsersFirst(t *testing.T) {
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

	_, err := SignUp("public", "test@gmail.com", "testPass123")
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
	user5, err := SignUp("public", "test4@gmail.com", "testPass123")
	if err != nil {
		t.Error(err)
		return
	}

	{
		paginationResult, err := supertokens.GetUsersNewestFirst("public", nil, nil, nil, nil)
		if err != nil {
			t.Error(err)
			return
		}

		email := "test4@gmail.com"
		assert.Nil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 5)
		assert.True(t, paginationResult.Users[0].LoginMethods[0].HasSameEmailAs(&email))
		assert.Equal(t, paginationResult.Users[0].ID, user5.OK.User.ID)
		assert.Equal(t, paginationResult.Users[0].ID, paginationResult.Users[0].LoginMethods[0].RecipeUserID.GetAsString())
	}

	{
		limit := 1
		paginationResult, err := supertokens.GetUsersNewestFirst("public", nil, &limit, nil, nil)
		if err != nil {
			t.Error(err)
			return
		}

		assert.NotNil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 1)
		assert.Equal(t, paginationResult.Users[0].Emails[0], "test4@gmail.com")

		paginationResult, err = supertokens.GetUsersNewestFirst("public", paginationResult.NextPaginationToken, &limit, nil, nil)
		if err != nil {
			t.Error(err)
			return
		}

		assert.NotNil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 1)
		assert.Equal(t, paginationResult.Users[0].Emails[0], "test3@gmail.com")

		limit = 5
		paginationResult, err = supertokens.GetUsersNewestFirst("public", paginationResult.NextPaginationToken, &limit, nil, nil)
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
		_, err := supertokens.GetUsersNewestFirst("public", &paginationToken, nil, nil, nil)
		if err != nil {
			assert.Contains(t, err.Error(), "invalid pagination token")
		} else {
			assert.Fail(t, "pagination token invalid should fail")
		}
	}
}

func TestGetOldestUsersFirstWithSearchParams(t *testing.T) {
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

	_, err := SignUp("public", "test@gmail.com", "testPass123")
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
	_, err = SignUp("public", "john@gmail.com", "testPass123")
	if err != nil {
		t.Error(err)
		return
	}

	{
		paginationResult, err := supertokens.GetUsersOldestFirst("public", nil, nil, nil, map[string]string{
			"email": "doe",
		})
		if err != nil {
			t.Error(err)
			return
		}

		assert.Nil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 0)
	}

	{
		paginationResult, err := supertokens.GetUsersOldestFirst("public", nil, nil, nil, map[string]string{
			"email": "john",
		})
		if err != nil {
			t.Error(err)
			return
		}

		assert.Nil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 1)
		assert.Equal(t, paginationResult.Users[0].Emails[0], "john@gmail.com")
		assert.Len(t, paginationResult.Users[0].LoginMethods, 1)
		assert.Len(t, paginationResult.Users[0].Emails, 1)
		assert.Len(t, paginationResult.Users[0].PhoneNumbers, 0)
		assert.Len(t, paginationResult.Users[0].ThirdParty, 0)
	}
}

func TestGetNewestUsersFirstWithSearchParams(t *testing.T) {
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

	_, err := SignUp("public", "test@gmail.com", "testPass123")
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
	_, err = SignUp("public", "john@gmail.com", "testPass123")
	if err != nil {
		t.Error(err)
		return
	}

	{
		paginationResult, err := supertokens.GetUsersNewestFirst("public", nil, nil, nil, map[string]string{
			"email": "doe",
		})
		if err != nil {
			t.Error(err)
			return
		}

		assert.Nil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 0)
	}

	{
		paginationResult, err := supertokens.GetUsersNewestFirst("public", nil, nil, nil, map[string]string{
			"email": "john",
		})
		if err != nil {
			t.Error(err)
			return
		}

		assert.Nil(t, paginationResult.NextPaginationToken)
		assert.Len(t, paginationResult.Users, 1)
		assert.Equal(t, paginationResult.Users[0].Emails[0], "john@gmail.com")
		assert.Len(t, paginationResult.Users[0].LoginMethods, 1)
		assert.Len(t, paginationResult.Users[0].Emails, 1)
		assert.Len(t, paginationResult.Users[0].PhoneNumbers, 0)
		assert.Len(t, paginationResult.Users[0].ThirdParty, 0)
	}
}

func TestGetUser(t *testing.T) {
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

	ogUser, err := SignUp("public", "test@gmail.com", "testPass123")
	if err != nil {
		t.Error(err)
		return
	}

	user, err := supertokens.GetUser(ogUser.OK.User.ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, user.ID, ogUser.OK.User.ID)
	assert.Equal(t, user.Emails[0], "test@gmail.com")
	assert.Len(t, user.LoginMethods, 1)
	assert.Len(t, user.Emails, 1)
	assert.Len(t, user.PhoneNumbers, 0)
	assert.Len(t, user.ThirdParty, 0)
	email := "test@gmail.com"
	assert.True(t, user.LoginMethods[0].HasSameEmailAs(&email))
	assert.Equal(t, user.ID, user.LoginMethods[0].RecipeUserID.GetAsString())
	assert.Equal(t, supertokens.EmailPasswordRID, user.LoginMethods[0].RecipeID)

	user, err = supertokens.GetUser("random")
	if err != nil {
		t.Error(err)
		return
	}

	assert.Nil(t, user)
}

func TestMakePrimaryUserSuccess(t *testing.T) {
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

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}
	user1 := convertEpUserToSuperTokensUser(epuser.OK.User)

	assert.False(t, user1.IsPrimaryUser)

	response, err := supertokens.CreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, response.OK.User.IsPrimaryUser)
	assert.False(t, response.OK.WasAlreadyAPrimaryUser)

	assert.Equal(t, user1.ID, response.OK.User.ID)
	assert.Equal(t, user1.Emails[0], response.OK.User.Emails[0])
	assert.Len(t, response.OK.User.LoginMethods, 1)

	refetchedUser, err := supertokens.GetUser(user1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, *refetchedUser, response.OK.User)
}

func TestMakePrimaryUserSuccessAlreadyPrimaryUser(t *testing.T) {
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

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}
	user1 := convertEpUserToSuperTokensUser(epuser.OK.User)

	assert.False(t, user1.IsPrimaryUser)

	response, err := supertokens.CreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, response.OK.User.IsPrimaryUser)
	assert.False(t, response.OK.WasAlreadyAPrimaryUser)

	response2, err := supertokens.CreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, response2.OK.User.IsPrimaryUser)
	assert.True(t, response2.OK.WasAlreadyAPrimaryUser)
	assert.Equal(t, response2.OK.User.ID, response.OK.User.ID)
}

func TestCanMakePrimaryUser(t *testing.T) {
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

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user1 := convertEpUserToSuperTokensUser(epuser.OK.User)

	assert.False(t, user1.IsPrimaryUser)

	response, err := supertokens.CanCreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.False(t, response.OK.WasAlreadyAPrimaryUser)

	_, err = supertokens.CreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}

	response, err = supertokens.CanCreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.True(t, response.OK.WasAlreadyAPrimaryUser)
}

func TestMakePrimaryFailCauseAlreadyLinkedToAnotherAccount(t *testing.T) {
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

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user1 := convertEpUserToSuperTokensUser(epuser.OK.User)
	assert.False(t, user1.IsPrimaryUser)

	epuser2, err := SignUp("public", "test2@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user2 := convertEpUserToSuperTokensUser(epuser2.OK.User)

	assert.False(t, user2.IsPrimaryUser)

	_, err = supertokens.CreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = supertokens.LinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
	if err != nil {
		t.Error(err)
		return
	}

	createPrimaryUserResponse, err := supertokens.CreatePrimaryUser(user2.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, createPrimaryUserResponse.OK)
	assert.Equal(t, createPrimaryUserResponse.RecipeUserIdAlreadyLinkedWithPrimaryUserIdError.PrimaryUserId, user1.ID)
}

// TODO: remove this function
func convertEpUserToSuperTokensUser(epuser epmodels.User) supertokens.User {
	rUId, err := supertokens.NewRecipeUserID(epuser.ID)
	if err != nil {
		panic(err.Error())
	}
	return supertokens.User{
		ID:            epuser.ID,
		TimeJoined:    epuser.TimeJoined,
		IsPrimaryUser: false,
		TenantIDs:     epuser.TenantIds,
		Emails:        []string{epuser.Email},
		PhoneNumbers:  []string{},
		ThirdParty:    []supertokens.ThirdParty{},
		LoginMethods: []supertokens.LoginMethods{
			{
				Verified: false,
				RecipeLevelUser: supertokens.RecipeLevelUser{
					TenantIDs:    epuser.TenantIds,
					TimeJoined:   epuser.TimeJoined,
					RecipeUserID: rUId,
					AccountInfoWithRecipeID: supertokens.AccountInfoWithRecipeID{
						RecipeID: supertokens.EmailPasswordRID,
						AccountInfo: supertokens.AccountInfo{
							Email: &epuser.Email,
						},
					},
				},
			},
		},
	}
}
