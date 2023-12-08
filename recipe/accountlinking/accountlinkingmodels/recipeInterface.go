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

package accountlinkingmodels

import "github.com/supertokens/supertokens-golang/supertokens"

type UserPaginationResult struct {
	Users               []supertokens.User
	NextPaginationToken *string
}

type CanCreatePrimaryUserResponse struct {
	OK *struct {
		WasAlreadyAPrimaryUser bool
	}
	RecipeUserIdAlreadyLinkedWithPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
	AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
}

type CreatePrimaryUserResponse struct {
	OK *struct {
		User                   supertokens.User
		WasAlreadyAPrimaryUser bool
	}
	RecipeUserIdAlreadyLinkedWithPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
	AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
}

type CanLinkAccountResponse struct {
	OK *struct {
		AccountsAlreadyLinked bool
	}
	RecipeUserIdAlreadyLinkedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
	AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
	InputUserIsNotAPrimaryUserError *struct{}
}

type LinkAccountResponse struct {
	OK *struct {
		AccountsAlreadyLinked bool
		User                  supertokens.User
	}
	RecipeUserIdAlreadyLinkedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		User          supertokens.User
	}
	AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
	InputUserIsNotAPrimaryUserError *struct{}
}

type UnlinkAccountsResponse struct {
	WasRecipeUserDeleted bool
	WasLinked            bool
}

type RecipeInterface struct {
	GetUsersWithSearchParams *func(tenantID string, timeJoinedOrder string, paginationToken *string, limit *int, includeRecipeIds *[]string, searchParams map[string]string, userContext supertokens.UserContext) (UserPaginationResult, error)

	CanCreatePrimaryUser *func(recipeUserId supertokens.RecipeUserID, userContext supertokens.UserContext) (CanCreatePrimaryUserResponse, error)

	CreatePrimaryUser *func(recipeUserId supertokens.RecipeUserID, userContext supertokens.UserContext) (CreatePrimaryUserResponse, error)

	CanLinkAccounts *func(recipeUserId supertokens.RecipeUserID, primaryUserId string, userContext supertokens.UserContext) (CanLinkAccountResponse, error)

	LinkAccounts *func(recipeUserId supertokens.RecipeUserID, primaryUserId string, userContext supertokens.UserContext) (LinkAccountResponse, error)

	UnlinkAccounts *func(recipeUserId supertokens.RecipeUserID, userContext supertokens.UserContext) (UnlinkAccountsResponse, error)

	GetUser *func(userId string, userContext supertokens.UserContext) (*supertokens.User, error)

	ListUsersByAccountInfo *func(tenantID string, accountInfo supertokens.AccountInfo, doUnionOfAccountInfo bool, userContext supertokens.UserContext) ([]supertokens.User, error)

	DeleteUser *func(userId string, removeAllLinkedAccounts bool, userContext supertokens.UserContext) error
}
