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
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
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

	canCreatePrimaryUserResponse, err := supertokens.CanCreatePrimaryUser(user2.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, canCreatePrimaryUserResponse.OK)
	assert.Equal(t, canCreatePrimaryUserResponse.RecipeUserIdAlreadyLinkedWithPrimaryUserIdError.PrimaryUserId, user1.ID)

	createPrimaryUserResponse, err := supertokens.CreatePrimaryUser(user2.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, createPrimaryUserResponse.OK)
	assert.Equal(t, createPrimaryUserResponse.RecipeUserIdAlreadyLinkedWithPrimaryUserIdError.PrimaryUserId, user1.ID)
}

func TestMakePrimaryFailCauseAccountInfoAlreadyAssociatedWithAnotherPrimaryUser(t *testing.T) {
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
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     "",
										ClientSecret: "",
									},
								},
							},
						},
					},
				},
			}),
		},
	})

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	epuser1 := convertEpUserToSuperTokensUser(epuser.OK.User)
	assert.False(t, epuser1.IsPrimaryUser)

	tpuser, err := thirdparty.ManuallyCreateOrUpdateUser("public", "google", "abc", "test@gmail.com")
	if err != nil {
		t.Error(err)
		return
	}

	tpUser1 := convertTpUserToSuperTokensUser(tpuser.OK.User)

	_, err = supertokens.CreatePrimaryUser(tpUser1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}

	canCreatePrimaryUserResult, err := supertokens.CanCreatePrimaryUser(epuser1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, canCreatePrimaryUserResult.OK)
	assert.Equal(t, canCreatePrimaryUserResult.AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError.PrimaryUserId, tpUser1.ID)

	createPrimaryUserResponse, err := supertokens.CreatePrimaryUser(epuser1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, createPrimaryUserResponse.OK)
	assert.Equal(t, createPrimaryUserResponse.AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError.PrimaryUserId, tpUser1.ID)
}

func TestLinkAccountsSuccess(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	telemetry := false
	var primaryUserInCallback supertokens.User
	var newAccountInfoInCallback supertokens.RecipeLevelUser
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
			session.Init(nil),
			Init(nil),
			supertokens.InitAccountLinking(&supertokens.AccountLinkingTypeInput{
				OnAccountLinked: func(user supertokens.User, newAccountUser supertokens.RecipeLevelUser, userContext supertokens.UserContext) error {
					primaryUserInCallback = user
					newAccountInfoInCallback = newAccountUser
					return nil
				},
			}),
		},
	})

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user1 := convertEpUserToSuperTokensUser(epuser.OK.User)
	assert.False(t, user1.IsPrimaryUser)
	supertokens.CreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)

	epuser2, err := SignUp("public", "test2@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user2 := convertEpUserToSuperTokensUser(epuser2.OK.User)
	assert.False(t, user2.IsPrimaryUser)

	// we create a new session to check that the session has not been revoked
	// when we link accounts, cause these users are already linked.
	session.CreateNewSessionWithoutRequestResponse("public", user2.ID, nil, nil, nil)
	sessions, err := session.GetAllSessionHandlesForUser(user2.ID, nil)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Len(t, sessions, 1)

	{
		linkAccountResponse, err := supertokens.CanLinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
		if err != nil {
			t.Error(err)
			return
		}
		assert.False(t, linkAccountResponse.OK.AccountsAlreadyLinked)
	}

	linkAccountResponse, err := supertokens.LinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.False(t, linkAccountResponse.OK.AccountsAlreadyLinked)

	linkedUser, err := supertokens.GetUser(user1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, *linkedUser, primaryUserInCallback)
	assert.Equal(t, *linkedUser, linkAccountResponse.OK.User)

	assert.Equal(t, newAccountInfoInCallback.RecipeID, supertokens.EmailPasswordRID)
	assert.Equal(t, *newAccountInfoInCallback.Email, "test2@gmail.com")
	sessions, err = session.GetAllSessionHandlesForUser(user2.LoginMethods[0].RecipeUserID.GetAsString(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Len(t, sessions, 0)
}

func TestLinkAccountsSuccessAlreadyLinked(t *testing.T) {
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
			session.Init(nil),
			Init(nil),
			supertokens.InitAccountLinking(&supertokens.AccountLinkingTypeInput{}),
		},
	})

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user1 := convertEpUserToSuperTokensUser(epuser.OK.User)
	assert.False(t, user1.IsPrimaryUser)
	supertokens.CreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)

	epuser2, err := SignUp("public", "test2@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user2 := convertEpUserToSuperTokensUser(epuser2.OK.User)
	assert.False(t, user2.IsPrimaryUser)

	linkAccountResponse, err := supertokens.LinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.False(t, linkAccountResponse.OK.AccountsAlreadyLinked)

	{
		canLinkAccountResponse, err := supertokens.CanLinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
		if err != nil {
			t.Error(err)
			return
		}
		assert.True(t, canLinkAccountResponse.OK.AccountsAlreadyLinked)
	}

	linkAccountResponse, err = supertokens.LinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, linkAccountResponse.OK.AccountsAlreadyLinked)
}

func TestLinkAccountsFailureAlreadyLinkedWithAnotherPrimaryUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	telemetry := false
	callbackCalled := false
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
			session.Init(nil),
			Init(nil),
			supertokens.InitAccountLinking(&supertokens.AccountLinkingTypeInput{
				OnAccountLinked: func(user supertokens.User, newAccountUser supertokens.RecipeLevelUser, userContext supertokens.UserContext) error {
					callbackCalled = true
					return nil
				},
			}),
		},
	})

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user1 := convertEpUserToSuperTokensUser(epuser.OK.User)
	assert.False(t, user1.IsPrimaryUser)
	supertokens.CreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)

	epuser2, err := SignUp("public", "test2@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user2 := convertEpUserToSuperTokensUser(epuser2.OK.User)
	assert.False(t, user2.IsPrimaryUser)

	linkAccountResponse, err := supertokens.LinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.False(t, linkAccountResponse.OK.AccountsAlreadyLinked)
	assert.True(t, callbackCalled)

	callbackCalled = false

	epuser3, err := SignUp("public", "test3@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user3 := convertEpUserToSuperTokensUser(epuser3.OK.User)
	assert.False(t, user3.IsPrimaryUser)
	supertokens.CreatePrimaryUser(user3.LoginMethods[0].RecipeUserID)

	{
		canLinkAccountResponse, err := supertokens.CanLinkAccounts(user2.LoginMethods[0].RecipeUserID, user3.ID)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Nil(t, canLinkAccountResponse.OK)
		assert.NotNil(t, canLinkAccountResponse.RecipeUserIdAlreadyLinkedWithAnotherPrimaryUserIdError)
		assert.Equal(t, canLinkAccountResponse.RecipeUserIdAlreadyLinkedWithAnotherPrimaryUserIdError.PrimaryUserId, user1.ID)
	}

	linkAccountResponse, err = supertokens.LinkAccounts(user2.LoginMethods[0].RecipeUserID, user3.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, linkAccountResponse.OK)
	assert.NotNil(t, linkAccountResponse.RecipeUserIdAlreadyLinkedWithAnotherPrimaryUserIdError)
	assert.Equal(t, linkAccountResponse.RecipeUserIdAlreadyLinkedWithAnotherPrimaryUserIdError.PrimaryUserId, user1.ID)

	assert.False(t, callbackCalled)
}

func TestLinkAccountsFailureInputUserIdNotAPrimaryUser(t *testing.T) {
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
			session.Init(nil),
			Init(nil),
			supertokens.InitAccountLinking(nil),
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

	{
		linkAccountResponse, err := supertokens.CanLinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Nil(t, linkAccountResponse.OK)
		assert.NotNil(t, linkAccountResponse.InputUserIsNotAPrimaryUserError)
	}

	linkAccountResponse, err := supertokens.LinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, linkAccountResponse.OK)
	assert.NotNil(t, linkAccountResponse.InputUserIsNotAPrimaryUserError)
}

func TestLinkAccountFailureAccountInfoAlreadyAssociatedWithAnotherPrimaryUser(t *testing.T) {
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
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     "",
										ClientSecret: "",
									},
								},
							},
						},
					},
				},
			}),
		},
	})

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	epuser1 := convertEpUserToSuperTokensUser(epuser.OK.User)
	assert.False(t, epuser1.IsPrimaryUser)
	supertokens.CreatePrimaryUser(epuser1.LoginMethods[0].RecipeUserID)

	tpuser, err := thirdparty.ManuallyCreateOrUpdateUser("public", "google", "abc", "test@gmail.com")
	if err != nil {
		t.Error(err)
		return
	}

	tpUser1 := convertTpUserToSuperTokensUser(tpuser.OK.User)

	epuser, err = SignUp("public", "test2@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	epuser2 := convertEpUserToSuperTokensUser(epuser.OK.User)
	supertokens.CreatePrimaryUser(epuser2.LoginMethods[0].RecipeUserID)

	{
		linkAccountResponse, err := supertokens.CanLinkAccounts(tpUser1.LoginMethods[0].RecipeUserID, epuser2.ID)
		if err != nil {
			t.Error(err)
			return
		}
		assert.NotNil(t, linkAccountResponse.AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError)
		assert.Equal(t, linkAccountResponse.AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError.PrimaryUserId, epuser1.ID)
	}

	linkAccountResponse, err := supertokens.LinkAccounts(tpUser1.LoginMethods[0].RecipeUserID, epuser2.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.NotNil(t, linkAccountResponse.AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError)
	assert.Equal(t, linkAccountResponse.AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError.PrimaryUserId, epuser1.ID)

}

func TestUnlinkAccountsSuccess(t *testing.T) {
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
			session.Init(nil),
			Init(nil),
			supertokens.InitAccountLinking(nil),
		},
	})

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user1 := convertEpUserToSuperTokensUser(epuser.OK.User)
	assert.False(t, user1.IsPrimaryUser)
	supertokens.CreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)

	epuser2, err := SignUp("public", "test2@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user2 := convertEpUserToSuperTokensUser(epuser2.OK.User)
	assert.False(t, user2.IsPrimaryUser)

	linkAccountResponse, err := supertokens.LinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.False(t, linkAccountResponse.OK.AccountsAlreadyLinked)

	session.CreateNewSessionWithoutRequestResponse("public", user2.ID, nil, nil, nil)
	sessions, err := session.GetAllSessionHandlesForUser(user2.ID, nil)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Len(t, sessions, 1)

	unlinkResponse, err := supertokens.UnlinkAccounts(user2.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.False(t, unlinkResponse.WasRecipeUserDeleted)
	assert.True(t, unlinkResponse.WasLinked)

	primaryUser, err := supertokens.GetUser(user1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Len(t, primaryUser.LoginMethods, 1)
	assert.True(t, primaryUser.IsPrimaryUser)

	recipeUser, err := supertokens.GetUser(user2.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Len(t, recipeUser.LoginMethods, 1)
	assert.False(t, recipeUser.IsPrimaryUser)

	sessions, err = session.GetAllSessionHandlesForUser(user2.LoginMethods[0].RecipeUserID.GetAsString(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Len(t, sessions, 0)
}

func TestUnlinkAccountsWithDeleteUserSuccess(t *testing.T) {
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
			session.Init(nil),
			Init(nil),
			supertokens.InitAccountLinking(nil),
		},
	})

	epuser, err := SignUp("public", "test@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user1 := convertEpUserToSuperTokensUser(epuser.OK.User)
	assert.False(t, user1.IsPrimaryUser)
	supertokens.CreatePrimaryUser(user1.LoginMethods[0].RecipeUserID)

	epuser2, err := SignUp("public", "test2@gmail.com", "pass123")
	if err != nil {
		t.Error(err)
		return
	}

	user2 := convertEpUserToSuperTokensUser(epuser2.OK.User)
	assert.False(t, user2.IsPrimaryUser)

	linkAccountResponse, err := supertokens.LinkAccounts(user2.LoginMethods[0].RecipeUserID, user1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.False(t, linkAccountResponse.OK.AccountsAlreadyLinked)

	unlinkResponse, err := supertokens.UnlinkAccounts(user1.LoginMethods[0].RecipeUserID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, unlinkResponse.WasRecipeUserDeleted)
	assert.True(t, unlinkResponse.WasLinked)
}

func TestDeleteUser(t *testing.T) {
	telemetry := false
	configValue := supertokens.TypeInput{
		Telemetry: &telemetry,
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

	res, err := SignUp("public", "test@example.com", "1234abcd")
	if err != nil {
		t.Error(err.Error())
	}
	reponseBeforeDeletingUser, err := supertokens.GetUsersOldestFirst("public", nil, nil, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 1, len(reponseBeforeDeletingUser.Users))
	err = supertokens.DeleteUser(res.OK.User.ID, true)
	if err != nil {
		t.Error(err.Error())
	}
	responseAfterDeletingUser, err := supertokens.GetUsersOldestFirst("public", nil, nil, nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 0, len(responseAfterDeletingUser.Users))
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

// TODO: remove this function
func convertTpUserToSuperTokensUser(tpuser tpmodels.User) supertokens.User {
	rUId, err := supertokens.NewRecipeUserID(tpuser.ID)
	if err != nil {
		panic(err.Error())
	}
	return supertokens.User{
		ID:            tpuser.ID,
		TimeJoined:    tpuser.TimeJoined,
		IsPrimaryUser: false,
		TenantIDs:     tpuser.TenantIds,
		Emails:        []string{tpuser.Email},
		PhoneNumbers:  []string{},
		ThirdParty: []supertokens.ThirdParty{
			{
				ID:     tpuser.ThirdParty.ID,
				UserID: tpuser.ThirdParty.UserID,
			},
		},
		LoginMethods: []supertokens.LoginMethods{
			{
				Verified: false,
				RecipeLevelUser: supertokens.RecipeLevelUser{
					TenantIDs:    tpuser.TenantIds,
					TimeJoined:   tpuser.TimeJoined,
					RecipeUserID: rUId,
					AccountInfoWithRecipeID: supertokens.AccountInfoWithRecipeID{
						RecipeID: supertokens.EmailPasswordRID,
						AccountInfo: supertokens.AccountInfo{
							Email: &tpuser.Email,
							ThirdParty: &supertokens.ThirdParty{
								ID:     tpuser.ThirdParty.ID,
								UserID: tpuser.ThirdParty.UserID,
							},
						},
					},
				},
			},
		},
	}
}
